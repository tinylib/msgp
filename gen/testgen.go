package gen

import (
	"bytes"
	"io"
)

// WriteMarshalUnmarshalTests writes tests for s.MarshalMsg and s.UnmarshalMsg, using
// buf as scratch space
func WriteMarshalUnmarshalTests(w io.Writer, s *Struct, buf *bytes.Buffer) error {
	return execAndFormat(marshalTestTemplate, w, s, buf)
}

// WriteEncodeDecodeTests writes tests for s.EncodeMsg and s.DecodeMsg, using
// buf as scratch space
func WriteEncodeDecodeTests(w io.Writer, s *Struct, buf *bytes.Buffer) error {
	return execAndFormat(encodeTestTemplate, w, s, buf)
}
