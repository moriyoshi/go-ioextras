// Copyright (c) 2014-2015 Moriyoshi Koizumi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ioextras

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
)

// RandomAccessStore models an I/O channel for a blob.
type RandomAccessStore interface {
	io.ReaderAt
	io.WriterAt
	io.Closer
}

// SizedRandomAccessStore models an I/O channel for a blob is supposed to have a specific size.
type SizedRandomAccessStore interface {
	RandomAccessStore
	Size() (int64, error)
}

// NamedRandomAccessStore models an I/O channel for a named blob (probably a file)
type NamedRandomAccessStore interface {
	RandomAccessStore
	Name() string
}

// RandomAccessStoreFactory creates a RandomAccessStore instance.
type RandomAccessStoreFactory interface {
	RandomAccessStore() (RandomAccessStore, error)
}

// SeekerWrapper wraps a RandomAccessStore for it to behave as io.Seeker.
type SeekerWrapper struct {
	s  RandomAccessStore
	sk io.Seeker
	ns NamedRandomAccessStore
}

// Just delegates the operation to the underlying RandomAccessStore's ReadAt() method.
func (s *SeekerWrapper) ReadAt(p []byte, offset int64) (int, error) { return s.s.ReadAt(p, offset) }

// Just delegates the operation to the underlying RandomAccessStore's WriteAt() method.
func (s *SeekerWrapper) WriteAt(p []byte, offset int64) (int, error) { return s.s.WriteAt(p, offset) }

// Just delegates the operation to the underlying RandomAccessStore's Close() method.
func (s *SeekerWrapper) Close() error { return s.s.Close() }

// If the underlying RandomAccessStore also provides NamedRandomAccessStore, delegates the call to its Name() method.  Otherwise, returns an empty string.
func (s *SeekerWrapper) Name() string {
	if s.ns != nil {
		return s.ns.Name()
	} else {
		return ""
	}
}

// Just delegates the operation to the underlying RandomAccessStore's Size() method.
func (s *SeekerWrapper) Size() (int64, error) {
	return s.sk.Seek(0, os.SEEK_END)
}

// Creates a new SeekerWrapper instance.
func NewSeekerWrapper(s RandomAccessStore) *SeekerWrapper {
	ns, _ := s.(NamedRandomAccessStore)
	return &SeekerWrapper{s, s.(io.Seeker), ns}
}

// StoreReadWriter wraps a RandomAccessStore for it to behave like io.Reader or io.Writer.
type StoreReadWriter struct {
	Store    RandomAccessStore
	Position int64
	Size     int64
}

func (rw *StoreReadWriter) Write(p []byte) (int, error) {
	n, err := rw.Store.WriteAt(p, rw.Position)
	rw.Position += int64(n)
	return n, err
}

func (rw *StoreReadWriter) Read(p []byte) (int, error) {
	n, err := rw.Store.ReadAt(p, rw.Position)
	if err == io.EOF {
		rw.Size = rw.Position + int64(n)
	}
	rw.Position += int64(n)
	return n, err
}

func (rw *StoreReadWriter) Close() error { return nil }

func (rw *StoreReadWriter) Seek(pos int64, whence int) (int64, error) {
	switch whence {
	case os.SEEK_SET:
		rw.Position = pos
	case os.SEEK_CUR:
		rw.Position += pos
	case os.SEEK_END:
		if rw.Size < 0 {
			return -1, errors.New("trying to seek to EOF while the store size is not known")
		}
		rw.Position = rw.Size + pos
	}
	return rw.Position, nil
}

// MemoryRandomAccessStore implements a RandomAccessStore backed by a byte slice.
type MemoryRandomAccessStore struct {
	buf []byte
}

func (s *MemoryRandomAccessStore) WriteAt(p []byte, offset int64) (int, error) {
	err := (error)(nil)
	o := int(offset)
	e := o + len(p)
	if e > len(s.buf) {
		if e <= cap(s.buf) {
			s.buf = s.buf[0:e]
		} else {
			newBuf := make([]byte, e, cap(s.buf)*2)
			copy(newBuf, s.buf)
			s.buf = newBuf
		}
	}
	n := e - o
	copy(s.buf[o:e], p)
	return n, err
}

func (s *MemoryRandomAccessStore) ReadAt(p []byte, offset int64) (int, error) {
	err := (error)(nil)
	o := int(offset)
	e := o + len(p)
	if e > len(s.buf) {
		e = len(s.buf)
		err = io.EOF
	}
	n := e - o
	copy(p, s.buf[o:e])
	return n, err
}

func (s *MemoryRandomAccessStore) Size() (int64, error) {
	return int64(len(s.buf)), nil
}

func (s *MemoryRandomAccessStore) Close() error { return nil }

func NewMemoryRandomAccessStore() *MemoryRandomAccessStore {
	return &MemoryRandomAccessStore{
		buf: make([]byte, 0, 16),
	}
}

type MemoryRandomAccessStoreFactory struct{}

func (ras *MemoryRandomAccessStoreFactory) RandomAccessStore() (RandomAccessStore, error) {
	return NewMemoryRandomAccessStore(), nil
}

// TempFileRandomAccessStoreFactory implements a RandomAccessStore backed by a temporary file
// that is created by ioutil.TempFIle.
// If the RandomAccessStore is closed, the underlying temporary file is deleted accordingly.
type TempFileRandomAccessStoreFactory struct {
	Dir    string
	Prefix string
	GCChan chan *os.File
}

func (ras *TempFileRandomAccessStoreFactory) RandomAccessStore() (RandomAccessStore, error) {
	f, err := ioutil.TempFile(ras.Dir, ras.Prefix)
	if err != nil {
		return nil, err
	}
	f_ := (RandomAccessStore)(f)
	if ras.GCChan != nil {
		c := ras.GCChan
		f_ = NewCloseHook(
			f_,
			func(s io.Closer) {
				c <- s.(*os.File)
			},
		)
	}
	return NewSeekerWrapper(f_), nil
}
