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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestDynamicRotatingWriter(t *testing.T) {
	count := 0
	baseDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Logf("%v", err)
		t.FailNow()
	}
	defer os.RemoveAll(baseDir)
	w := NewDynamicRotatingWriter(
		func(_ interface{}) string {
			count += 1
			return strconv.Itoa(count)
		},
		StandardWriterFactory,
		func(_ interface{}) string {
			return filepath.Join(baseDir, "TEST")
		},
		SerialRotationCallbackFactory(3),
		nil,
	)
	w.Write([]byte("aaa\n"))
	w.Write([]byte("bbb\n"))
	w.Write([]byte("ccc\n"))
	w.Write([]byte("ddd\n"))
	w.Write([]byte("eee\n"))
	w.Write([]byte("fff\n"))
	w.Write([]byte("ggg\n"))
	c := 0
	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			c += 1
		}
		return err
	})
	if c != 4 {
		t.Logf("%d (%d)", c, count)
		t.Fail()
	}
}
