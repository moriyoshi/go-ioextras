package ioextras

import (
	"os"
	"io"
	"sync"
	"sync/atomic"
)

type PathBuilder func(interface{}) (string, error)
type WriterFactory func(string, interface{}) (io.Writer, error)


type CloserErrorPair struct {
	Closer io.Closer
	Error error
}

type writerEntry struct {
	path string
	w io.Writer
	refs int64
	closeErrorReportChan chan<- CloserErrorPair
}

type StaticRotatingWriter struct {
	PathBuilder PathBuilder
	WriterFactory WriterFactory
	CloseErrorReportChan chan<- CloserErrorPair
	writersMtx sync.Mutex
	writers map[string]*writerEntry
}

func (w *writerEntry) addRef() {
	atomic.AddInt64(&w.refs, 1)
}

func (w *writerEntry) delRef() bool {
	refs := atomic.AddInt64(&w.refs, -1)
	if refs == 0 {
		c, ok := w.w.(io.Closer)
		if ok {
			err := c.Close()
			if err != nil {
				if w.closeErrorReportChan != nil {
					w.closeErrorReportChan <- CloserErrorPair { c, err }
				}
			}
		}
		return true
	} else if refs < 0 {
		panic("something went wrong!")
	}
	return false
}

func (w *StaticRotatingWriter) Write(b []byte) (int, error) {
	return w.WriteWithCtx(b, nil)
}

func (w *StaticRotatingWriter) WriteWithCtx(b []byte, ctx interface{}) (int, error) {
	path, err := w.PathBuilder(ctx)
	if err != nil {
		return 0, err
	}
	we, err := func(path string) (*writerEntry, error) {
		w.writersMtx.Lock()
		defer w.writersMtx.Unlock()
		we, ok := w.writers[path]
		if !ok {
			writersToBeRemoved := make([]*writerEntry, 0, len(w.writers))
			for _, we_ := range w.writers {
				if we_.delRef() {
					writersToBeRemoved = append(writersToBeRemoved, we_)
				}
			}
			for _, we_ := range writersToBeRemoved {
				delete(w.writers, we_.path)
			}
			w_, err := w.WriterFactory(path, ctx)
			if err != nil {
				return nil, err
			}
			we = &writerEntry {
				path: path,
				w: w_,
				refs: 1,
				closeErrorReportChan: (chan<-CloserErrorPair)(w.CloseErrorReportChan),
			}
			w.writers[path] = we
		}
		we.addRef()
		return we, nil
	}(path)
	if err != nil {
		return 0, err
	}
	defer func() {
		if we.delRef() {
			w.writersMtx.Lock()
			defer w.writersMtx.Unlock()
			delete(w.writers, we.path)
		}
	}()
	return we.w.Write(b)
}

func (w *StaticRotatingWriter) Close() error {
	w.writersMtx.Lock()
	defer w.writersMtx.Unlock()
	writersToBeRemoved := make([]*writerEntry, 0, len(w.writers))
	for _, we := range w.writers {
		if we.delRef() {
			writersToBeRemoved = append(writersToBeRemoved, we)
		}
	}
	for _, we := range writersToBeRemoved {
		delete(w.writers, we.path)
	}
	close(w.CloseErrorReportChan)
	return nil
}

func NewStaticRotatingWriter(pathBuilder PathBuilder, writerFactory WriterFactory, closeErrorReportChan chan<-CloserErrorPair) *StaticRotatingWriter {
	if closeErrorReportChan == nil {
		closeErrorReportChan_ := make(chan CloserErrorPair)
		closeErrorReportChan = (chan<- CloserErrorPair)(closeErrorReportChan_)
		go func() {
			for _ = range closeErrorReportChan_ {} // just ignoring errors
		}()
	}
	return &StaticRotatingWriter {
		PathBuilder: pathBuilder,
		WriterFactory: writerFactory,
		CloseErrorReportChan: closeErrorReportChan,
		writersMtx: sync.Mutex{},
		writers: make(map[string]*writerEntry),
	}
}

func StandardWriterFactory(path string, _ interface{}) (io.Writer, error) {
	return os.OpenFile(path, os.O_CREATE | os.O_WRONLY | os.O_APPEND, os.FileMode(0666))
}
