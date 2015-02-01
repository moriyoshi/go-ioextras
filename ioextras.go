// ioextras provides a number of interfaces that semantically reinforces the Go's standard io
// primitives interfaces with a number of extra concepts, in addition to utility functions.
//
// An I/O channel is an abstraction of a stateful object (that is either taken care of by the operating
// system or implemented in userland) from/to which information is retrieved / stored through a set of
// I/O primitives, such as Read() and Write().
//
// A blob is opaque data which can be randomly accessed by io.ReaderAt and io.WriterAt while a stream
// can be accessed by io.Reader and io.Writer, additionally with io.Seeker.
//
package ioextras
