package ioextras

import (
	"io"
)

type CloseHook struct {
	IOCombo
	Callback func(s io.Closer)
}

func (c *CloseHook) Close() error {
	err := c.Closer.Close()
	if err != nil {
		return err
	}
	c.Callback(c.Closer)
	return nil
}

func NewCloseHook(c io.Closer, callback func(s io.Closer)) *CloseHook {
	reader, _ := c.(io.Reader)
	contextualReader, _ := c.(ContextualReader)
	writer, _ := c.(io.Writer)
	contextualWriter, _ := c.(ContextualWriter)
	readerAt, _ := c.(io.ReaderAt)
	writerAt, _ := c.(io.WriterAt)
	seeker, _ := c.(io.Seeker)
	flusher, _ := c.(Flusher)
	sized, _ := c.(Sized)
	named, _ := c.(Named)
	return &CloseHook{
		IOCombo: IOCombo {
			Reader: reader,
			ReaderAt: readerAt,
			ContextualReader: contextualReader,
			Writer: writer,
			WriterAt: writerAt,
			ContextualWriter: contextualWriter,
			Seeker: seeker,
			Closer: c,
			Flusher: flusher,
			Sized: sized,
			Named: named,
		},
		Callback: callback,
	}
}
