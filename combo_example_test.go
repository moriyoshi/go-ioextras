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
	"log"
	"io"
	"strings"
	"github.com/moriyoshi/go-ioextras"
)

type myCloser int
func (myCloser) Close() error {
	log.Printf("File is being closed")
	return nil
}

func ExampleIOCombo() {
	combo := &ioextras.IOCombo {
		Reader: strings.NewReader("test"),
		Closer: myCloser(0),
	}
	// an IOCombo should have provide io.ReadCloser interface.
	rc := io.ReadCloser(combo)

	b := make([]byte, 4)
	_, err := rc.Read(b)
	if err != nil {
		log.Fatalf("Error during reading: %v", err)
	}
	rc.Close()
}
