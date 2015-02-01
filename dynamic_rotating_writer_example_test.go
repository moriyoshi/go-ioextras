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

package ioextras_test

import (
	"github.com/moriyoshi/go-ioextras"
	"io"
	"log"
	"os"
	"strconv"
)

func ExampleDynamicRotatingWriter() {
	currentId := 0
	maxSize := int64(4)
	w := ioextras.NewDynamicRotatingWriter(
		// IDBuilder
		func(f io.Writer, _ interface{}) string {
			_f, ok := f.(*os.File)
			if !ok {
				panic("WTF?")
			}
			fi, err := _f.Stat()
			if err != nil {
				return strconv.Itoa(currentId)
			}
			if fi.Size() > maxSize {
				currentId++
			}
			return strconv.Itoa(currentId)
		},
		// WriterFactory
		ioextras.StandardWriterFactory,
		// HeadPathGenerator
		func(id string, _ interface{}) string {
			return "/tmp/demo.log"
		},
		// RotationCallback
		ioextras.SerialRotationCallbackFactory(3),
		// CloseErrorReportChan
		nil,
	)
	var err error

	_ ,err = w.Write([]byte("test"))
	if err != nil {
		log.Fatalf("error writing data: %v", err)
	}
	_, err = w.Write([]byte("test"))
	if err != nil {
		log.Fatalf("error writing data: %v", err)
	}
	_, err = w.Write([]byte("test"))
	if err != nil {
		log.Fatalf("error writing data: %v", err)
	}
	_, err = w.Write([]byte("test"))
	if err != nil {
		log.Fatalf("error writing data: %v", err)
	}
}
