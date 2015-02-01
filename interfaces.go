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

// Defines a writer interface accompanied by the opaque context information
type ContextualWriter interface {
	WriteWithCtx([]byte, interface{}) (int, error)
}

// Defines a reader interface accompanied by the opaque context information
type ContextualReader interface {
	ReadWithCtx([]byte, interface{}) (int, error)
}

// Flusher is an I/O channel (or stream) that provides `Flush` operation.
type Flusher interface {
	Flush() error
}

// Sized is an I/O concept for a blob of a specific size.  An I/O channel (or stream) backed by such a blob may also have this.
type Sized interface {
	Size() (int64, error)
}

// Named is an I/O concept for a named blob.  An I/O channel (or stream) backed by such a blob may also have this.
type Named interface {
	Name() string
}
