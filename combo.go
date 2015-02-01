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
	"io"
	"errors"
)

// Returned by IOCombo methods when the specified operation is not supported by the underlying
// implementation.
var Unsupported = errors.New("Unsupported operation")

// IOCombo combines objects implementing I/O primitives interfaces in the way it would virtually
// provide all of the I/O privimites supported by the given objects.
//
// Typical usage is to provide extra primitives interfaces such as io.Closer for a single I/O
// privimitives interface like io.Reader.  In that case, ioutil.NopCloser() can be used instead
// unless any operation is needed in a call to Close() method.
//
// Call on an unsupported operation returns the Unsupported error (that is a singleton).
//
type IOCombo struct {
	Reader io.Reader
	ReaderAt io.ReaderAt
	ContextualReader ContextualReader
	Writer io.Writer
	WriterAt io.WriterAt
	ContextualWriter ContextualWriter
	Seeker io.Seeker
	Closer io.Closer
	Flusher Flusher
	Sized Sized
	Named Named
}

func (w *IOCombo) Read(b []byte) (int, error) {
	if w.Reader == nil {
		return 0, Unsupported
	}
	return w.Reader.Read(b)
}

func (w *IOCombo) ReadAt(b []byte, o int64) (int, error) {
	if w.ReaderAt == nil {
		return 0, Unsupported
	}
	return w.ReaderAt.ReadAt(b, o)
}

func (w *IOCombo) ReadWithCtx(b []byte, ctx interface{}) (int, error) {
	if w.ContextualReader == nil {
		return 0, Unsupported
	}
	return w.ContextualReader.ReadWithCtx(b, ctx)
}

func (w *IOCombo) Write(b []byte) (int, error) {
	if w.Writer == nil {
		return 0, Unsupported
	}
	return w.Writer.Write(b)
}

func (w *IOCombo) WriteWithCtx(b []byte, ctx interface{}) (int, error) {
	if w.ContextualWriter == nil {
		return 0, Unsupported
	}
	return w.ContextualWriter.WriteWithCtx(b, ctx)
}

func (w *IOCombo) WriteAt(b []byte, o int64) (int, error) {
	if w.WriterAt == nil {
		return 0, Unsupported
	}
	return w.WriterAt.WriteAt(b, o)
}

func (w *IOCombo) Seek(o int64, whence int) (int64, error) {
	if w.Seeker == nil {
		return 0, Unsupported
	}
	return w.Seeker.Seek(o, whence)
}

func (w *IOCombo) Flush() error {
	if w.Flusher == nil {
		return Unsupported
	}
	return w.Flusher.Flush()
}

func (w *IOCombo) Close() error {
	if w.Flusher != nil {
		err := w.Flusher.Flush()
		if err != nil {
			return err
		}
	}
	if w.Closer != nil {
		return w.Closer.Close()
	}
	return nil
}

func (w *IOCombo) Size() (int64, error) {
	if w.Sized == nil {
		return 0, Unsupported
	}
	return w.Sized.Size()
}

func (w *IOCombo) Name() string {
	if w.Named == nil {
		return ""
	}
	return w.Named.Name()
}
