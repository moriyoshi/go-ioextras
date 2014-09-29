package ioextras

import (
	"io"
	"os"
	"fmt"
	"sync"
)

type IDBuilder func(ctx interface{}) string
type HeadPathGenerator func(ctx interface{}) string

type RotationEvent struct {
	ID string
	Path string
}

type RotationCallback func(ID, Path string, ctx interface{}) error

type DynamicRotatingWriter struct {
	IDBuilder IDBuilder
	WriterFactory WriterFactory
	HeadPathGenerator HeadPathGenerator
	RotationCallback RotationCallback
	CloseErrorReportChan chan<- CloserErrorPair
	mtx sync.Mutex
	currentID string
	currentWriter io.Writer
	currentPath string
	closed bool
}

func (w *DynamicRotatingWriter) Write(b []byte) (int, error) {
	return w.WriteWithCtx(b, nil)
}

func (w *DynamicRotatingWriter) WriteWithCtx(b []byte, ctx interface{}) (int, error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.closed {
		return 0, io.EOF
	}
	id := w.IDBuilder(ctx)
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
		path := w.HeadPathGenerator(ctx)
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

func NewDynamicRotatingWriter(idBuilder IDBuilder, writerFactory WriterFactory, headPathGenerator HeadPathGenerator, rotationCallback RotationCallback, closeErrorReportChan chan<-CloserErrorPair) *DynamicRotatingWriter {
	if closeErrorReportChan == nil {
		closeErrorReportChan_ := make(chan CloserErrorPair)
		closeErrorReportChan = (chan<- CloserErrorPair)(closeErrorReportChan_)
		go func() {
			for _ = range closeErrorReportChan_ {} // just ignoring errors
		}()
	}
	return &DynamicRotatingWriter {
		IDBuilder: idBuilder,
		WriterFactory: writerFactory,
		HeadPathGenerator: headPathGenerator,
		RotationCallback: rotationCallback,
		CloseErrorReportChan: closeErrorReportChan,
		currentID: "",
		currentWriter: nil,
		currentPath: "",
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
	if n + 1 >= maxFiles {
		err = os.Remove(path)
	} else {
		var path_ string
		path_, err = makeRoom(basePath, n + 1, maxFiles)
		if err != nil {
			return "", err
		}
		err = os.Rename(path, path_)
	}
	return path, err
}


func SerialRotationCallbackFactory(maxFiles int) RotationCallback {
	return func (id string, path string, _ interface{}) error {
		newPath, err := makeRoom(path, 0, maxFiles)
		if err != nil {
			return err
		}
		return os.Rename(path, newPath)
	}
}
