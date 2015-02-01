package ioextras

import (
	"io"
)

// CloseHook wraps a io.Closer so that it will call the specified callback function on closing the
// channel.  The callback will only be called on an explicit close operation by the user.
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

// Creates a new CloseHook from a io.Closer and a callback that will be called on Close().
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
