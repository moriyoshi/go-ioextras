package ioextras

type ContextualWriter interface {
	WriteWithCtx([]byte, interface{}) (int, error)
}

type ContextualReader interface {
	ReadWithCtx([]byte, interface{}) (int, error)
}

type Flusher interface {
	Flush() error
}

type Sized interface {
	Size() (int64, error)
}

type Named interface {
	Name() string
}
