package ioextras

import (
	"testing"
	"strconv"
	"os"
	"io/ioutil"
	"path/filepath"
)

func TestDynamicRotatingWriter (t *testing.T) {
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
	if c != 4 { t.Logf("%d (%d)", c, count); t.Fail() }
}
