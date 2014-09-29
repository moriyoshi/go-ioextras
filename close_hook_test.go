package ioextras

import (
	"io"
	"testing"
)

type dummyCloser2 struct{}

func (c *dummyCloser2) Close() error {
	return nil
}

func TestCloseHook(t *testing.T) {
	h := NewCloseHook(&dummyCloser2{}, func(c io.Closer) {
		_, ok := c.(*dummyCloser2)
		if !ok {
			t.Fail()
		}
	})
	err := h.Close()
	if err != nil { t.Fail() }
}
