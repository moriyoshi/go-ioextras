package ioextras

// Defines a writer interface accompanied by the opaque context information
type ContextualWriter interface {
	WriteWithCtx([]byte, interface{}) (int, error)
}

// Defines a reader interface accompanied by the opaque context information
type ContextualReader interface {
	ReadWithCtx([]byte, interface{}) (int, error)
}

// Flusher is an I/O channel (or stream) that provides `Flush` operation.
type Flusher interface {
	Flush() error
}

// Sized is an I/O concept for a blob of a specific size.  An I/O channel (or stream) backed by such a blob may also have this.
type Sized interface {
	Size() (int64, error)
}

// Named is an I/O concept for a named blob.  An I/O channel (or stream) backed by such a blob may also have this.
type Named interface {
	Name() string
}
