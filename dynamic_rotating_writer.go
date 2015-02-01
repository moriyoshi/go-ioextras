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
	"fmt"
	"io"
	"os"
	"sync"
)

// IDBuilder is supposed to return a string that determines the rotation cycle.
// Rotation occurs when it returns an identifier different from the one obtained during
// the previous call to WriteWithCtx() method.
// ctx argument is the value passed to WriteWithCtx() method as the second argument.
type IDBuilder func(w io.Writer, ctx interface{}) string

// HeadPathGenerator is supposed to return the path to the file data is written to.
// Generally the path is a static location.
// ctx argument is the value passed to WriteWithCtx() method as the second argument.
type HeadPathGenerator func(ID string, ctx interface{}) string

// RotationCallback is called when rotation occurs.  Typically the callback would rename the
// file specified by path to something like xxx.1 so that the generated files are rotating.
type RotationCallback func(ID, path string, ctx interface{}) error

// DynamicRotatingWriter is an io.Writer that writes the data to the file determined by
// HeadPathGenerator until the next rotation cycle.  The rotation cycle is determined by the return
// value of IDBuilder.  The file is created by WriterFactory.  RotationCallback will be called
// on rotation.
type DynamicRotatingWriter struct {
	IDBuilder            IDBuilder
	WriterFactory        WriterFactory
	HeadPathGenerator    HeadPathGenerator
	RotationCallback     RotationCallback
	CloseErrorReportChan chan<- CloserErrorPair
	mtx                  sync.Mutex
	currentID            string
	currentWriter        io.Writer
	currentPath          string
	closed               bool
}

// Write() method of io.Writer() interface.  This simply calls WriteWithCtx() with the second argument
// being nil.
func (w *DynamicRotatingWriter) Write(b []byte) (int, error) {
	return w.WriteWithCtx(b, nil)
}

// WriteWithCtx() method of ContextualWriter interface.  This function is designed to be reentrant,
// Because of the dynamic nature of 
func (w *DynamicRotatingWriter) WriteWithCtx(b []byte, ctx interface{}) (int, error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.closed {
		return 0, io.EOF
	}
	id := w.IDBuilder(w, ctx)
	if id != w.currentID || w.currentWriter == nil {
		if w.currentWriter != nil {
			c, ok := (w.currentWriter).(io.Closer)
			if ok {
				err := c.Close()
				if err != nil {
					w.CloseErrorReportChan <- CloserErrorPair{c, err}
				}
			}
			w.currentWriter = nil
		}
		if w.currentPath != "" {
			if w.RotationCallback != nil {
				err := w.RotationCallback(w.currentID, w.currentPath, ctx)
				if err != nil {
					return 0, err
				}
			}
			w.currentPath = ""
		}
		path := w.HeadPathGenerator(id, ctx)
		wr, err := w.WriterFactory(path, ctx)
		if err != nil {
			return 0, err
		}
		w.currentWriter = wr
		w.currentPath = path
		w.currentID = id
	}
	return w.currentWriter.Write(b)
}

func (w *DynamicRotatingWriter) Close() error {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.closed {
		return nil
	}
	if w.currentWriter != nil {
		c, ok := (w.currentWriter).(io.Closer)
		if ok {
			err := c.Close()
			if err != nil {
				w.CloseErrorReportChan <- CloserErrorPair{c, err}
			}
		}
		w.currentWriter = nil
	}
	w.closed = true
	close(w.CloseErrorReportChan)
	return nil
}

// Creates a new DynamicRotatingWriter. Pass StandardWriterFactory as writerFactory if you aren't
// interested in any contextual information passed as the second argument of WriteWithCtx() when
// opening the file.  closeErrorReportChan is a channel that will asynchronously receive errors
// that occur during closing files that have been opened bby writerFastory.  closeErrorReportChan
// can be nil.
func NewDynamicRotatingWriter(idBuilder IDBuilder, writerFactory WriterFactory, headPathGenerator HeadPathGenerator, rotationCallback RotationCallback, closeErrorReportChan chan<- CloserErrorPair) *DynamicRotatingWriter {
	if closeErrorReportChan == nil {
		closeErrorReportChan_ := make(chan CloserErrorPair)
		closeErrorReportChan = (chan<- CloserErrorPair)(closeErrorReportChan_)
		go func() {
			for _ = range closeErrorReportChan_ {
			} // just ignoring errors
		}()
	}
	return &DynamicRotatingWriter{
		IDBuilder:            idBuilder,
		WriterFactory:        writerFactory,
		HeadPathGenerator:    headPathGenerator,
		RotationCallback:     rotationCallback,
		CloseErrorReportChan: closeErrorReportChan,
		currentID:            "",
		currentWriter:        nil,
		currentPath:          "",
	}
}

func makeRotatedPath(path string, n int) string {
	return fmt.Sprintf("%s.%d", path, n)
}

func makeRoom(basePath string, n int, maxFiles int) (string, error) {
	path := makeRotatedPath(basePath, n)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return path, nil
		}
		return "", err
	}
	if n+1 >= maxFiles {
		err = os.Remove(path)
	} else {
		var path_ string
		path_, err = makeRoom(basePath, n+1, maxFiles)
		if err != nil {
			return "", err
		}
		err = os.Rename(path, path_)
	}
	return path, err
}

// This is a factory function that returns a typical implementation of RotationCallback which
// move the file specified by path argument to the one suffixed by ".1" after moving the file
// in the destination path to that with the suffix changed to what the number part is incremented by one
// (".2" for ".1").  Renaming is done accordingly until at most maxFile number of files remain.
func SerialRotationCallbackFactory(maxFiles int) RotationCallback {
	return func(id string, path string, _ interface{}) error {
		newPath, err := makeRoom(path, 0, maxFiles)
		if err != nil {
			return err
		}
		return os.Rename(path, newPath)
	}
}
