package ioextras

import (
	"io"
	"errors"
)

// Returned by IOCombo methods when the specified operation is not supported by the underlying
// implementation.
var Unsupported = errors.New("Unsupported operation")

// IOCombo wraps any I/O primitives interface in the way the resulting object provides
// all of the I/O privimites provided by io package as well as by ioextra.
//
// Call on an unsupported operation returns Unsupported error.
//
type IOCombo struct {
	Reader io.Reader
	ReaderAt io.ReaderAt
	ContextualReader ContextualReader
	Writer io.Writer
	WriterAt io.WriterAt
	ContextualWriter ContextualWriter
	Seeker io.Seeker
	Closer io.Closer
	Flusher Flusher
	Sized Sized
	Named Named
}

func (w *IOCombo) Read(b []byte) (int, error) {
	if w.Reader == nil {
		return 0, Unsupported
	}
	return w.Reader.Read(b)
}

func (w *IOCombo) ReadAt(b []byte, o int64) (int, error) {
	if w.ReaderAt == nil {
		return 0, Unsupported
	}
	return w.ReaderAt.ReadAt(b, o)
}

func (w *IOCombo) ReadWithCtx(b []byte, ctx interface{}) (int, error) {
	if w.ContextualReader == nil {
		return 0, Unsupported
	}
	return w.ContextualReader.ReadWithCtx(b, ctx)
}

func (w *IOCombo) Write(b []byte) (int, error) {
	if w.Writer == nil {
		return 0, Unsupported
	}
	return w.Writer.Write(b)
}

func (w *IOCombo) WriteWithCtx(b []byte, ctx interface{}) (int, error) {
	if w.ContextualWriter == nil {
		return 0, Unsupported
	}
	return w.ContextualWriter.WriteWithCtx(b, ctx)
}

func (w *IOCombo) WriteAt(b []byte, o int64) (int, error) {
	if w.WriterAt == nil {
		return 0, Unsupported
	}
	return w.WriterAt.WriteAt(b, o)
}

func (w *IOCombo) Seek(o int64, whence int) (int64, error) {
	if w.Seeker == nil {
		return 0, Unsupported
	}
	return w.Seeker.Seek(o, whence)
}

func (w *IOCombo) Flush() error {
	if w.Flusher == nil {
		return Unsupported
	}
	return w.Flusher.Flush()
}

func (w *IOCombo) Close() error {
	if w.Flusher != nil {
		err := w.Flusher.Flush()
		if err != nil {
			return err
		}
	}
	if w.Closer != nil {
		return w.Closer.Close()
	}
	return nil
}

func (w *IOCombo) Size() (int64, error) {
	if w.Sized == nil {
		return 0, Unsupported
	}
	return w.Sized.Size()
}

func (w *IOCombo) Name() string {
	if w.Named == nil {
		return ""
	}
	return w.Named.Name()
}
