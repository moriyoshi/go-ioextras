package ioextras

import (
	"testing"
	"io"
	"sync"
	"bytes"
	"strconv"
	"sync/atomic"
)

type dummyCloser struct {
	results *[][]byte
	w *bytes.Buffer
	l sync.Mutex
	wg *sync.WaitGroup
}

func (c *dummyCloser) Close() error {
	c.l.Lock()
	defer c.l.Unlock()
	*c.results = append(*c.results, c.w.Bytes())
	if c.wg != nil {
		c.wg.Done()
	}
	return nil
}

func TestStaticRotatingWriter(t *testing.T) {
	count := -1
	writers := make([]*bytes.Buffer, 0)
	results := make([][]byte, 0)
	w := NewStaticRotatingWriter(
		func(ctx interface{}) (string, error) {
			count += 1
			return strconv.Itoa(count / 2), nil
		},
		func (path string, ctx interface{}) (io.Writer, error) {
			w := &bytes.Buffer{}
			writers = append(writers, w)
			return &IOCombo{Writer: w, Closer: &dummyCloser{&results, w, sync.Mutex{}, nil}}, nil
		},
		nil,
	)
	w.Write([]byte("aaa"))
	w.Write([]byte("bbb"))
	w.Write([]byte("ccc"))
	w.Write([]byte("ddd"))
	w.Write([]byte("eee"))
	w.Write([]byte("fff"))
	t.Logf("len(writers)=%d", len(writers))
	if len(writers) != 3 { t.Fail() }
	t.Logf("len(results)=%d", len(results))
	if len(results) != 2 { t.Fail() }
	if !bytes.Equal(writers[0].Bytes(), []byte("aaabbb")) { t.Fail() }
	if !bytes.Equal(writers[0].Bytes(), results[0]) { t.Fail() }
	if !bytes.Equal(writers[1].Bytes(), []byte("cccddd")) { t.Fail() }
	if !bytes.Equal(writers[1].Bytes(), results[1]) { t.Fail() }
	if !bytes.Equal(writers[2].Bytes(), []byte("eeefff")) { t.Fail() }
}

type fancyWriter struct {
	w io.Writer
	condFulfilled *bool
	cond *sync.Cond
}

func (w *fancyWriter) Write(b []byte) (int, error) {
	w.cond.L.Lock()
	if !*w.condFulfilled {
		w.cond.Wait()
	}
	w.cond.L.Unlock()
	return w.w.Write(b)
}

func TestStaticRotatingWriterConcurrent(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count := int64(-1)
			writers := make([]*bytes.Buffer, 0)
			writersMtx := sync.Mutex{}
			results := make([][]byte, 0)
			condFulfilled := false
			cond := sync.Cond{L: &sync.Mutex{}}
			wg := sync.WaitGroup{}
			w := NewStaticRotatingWriter(
				func(ctx interface{}) (string, error) {
					c := atomic.AddInt64(&count, 1)
					return strconv.Itoa(int(c)), nil
				},
				func (path string, ctx interface{}) (io.Writer, error) {
					writersMtx.Lock()
					defer writersMtx.Unlock()
					w := &bytes.Buffer{}
					writers = append(writers, w)
					return &IOCombo{Writer: &fancyWriter{w, &condFulfilled, &cond}, Closer: &dummyCloser{&results, w, sync.Mutex{}, &wg}}, nil
				},
				nil,
			)
			wg.Add(5)
			go func() { w.Write([]byte("aaa")) }()
			go func() { w.Write([]byte("bbb")) }()
			go func() { w.Write([]byte("ccc")) }()
			go func() { w.Write([]byte("ddd")) }()
			go func() { w.Write([]byte("eee")) }()
			go func() { w.Write([]byte("fff")) }()
			cond.L.Lock()
			condFulfilled = true
			cond.Broadcast()
			cond.L.Unlock()
			wg.Wait()
			t.Logf("len(writers)=%d", len(writers))
			if len(writers) != 6 { t.Fail() }
			t.Logf("len(results)=%d", len(results))
			if len(results) != 5 { t.Fail() }
		}()
	}
	wg.Wait()
}
