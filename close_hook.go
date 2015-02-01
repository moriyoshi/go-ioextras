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
)

// CloseHook wraps a io.Closer so that it will call the specified callback function on closing the
// channel.  The callback will only be called on an explicit close operation by the user.
type CloseHook struct {
	IOCombo
	Callback func(s io.Closer)
}

func (c *CloseHook) Close() error {
	err := c.Closer.Close()
	if err != nil {
		return err
	}
	c.Callback(c.Closer)
	return nil
}

// Creates a new CloseHook from a io.Closer and a callback that will be called on Close().
func NewCloseHook(c io.Closer, callback func(s io.Closer)) *CloseHook {
	reader, _ := c.(io.Reader)
	contextualReader, _ := c.(ContextualReader)
	writer, _ := c.(io.Writer)
	contextualWriter, _ := c.(ContextualWriter)
	readerAt, _ := c.(io.ReaderAt)
	writerAt, _ := c.(io.WriterAt)
	seeker, _ := c.(io.Seeker)
	flusher, _ := c.(Flusher)
	sized, _ := c.(Sized)
	named, _ := c.(Named)
	return &CloseHook{
		IOCombo: IOCombo {
			Reader: reader,
			ReaderAt: readerAt,
			ContextualReader: contextualReader,
			Writer: writer,
			WriterAt: writerAt,
			ContextualWriter: contextualWriter,
			Seeker: seeker,
			Closer: c,
			Flusher: flusher,
			Sized: sized,
			Named: named,
		},
		Callback: callback,
	}
}
