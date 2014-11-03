package msgp

import (
	"bytes"
	"io"
)

// Raw is a []byte that is already messagepack-formatted
type Raw []byte

// EncodeMsg implements MsgEncoder.EncodeMsg
func (r Raw) EncodeMsg(w io.Writer) (int, error) {
	return w.Write([]byte(r))
}

// EncodeTo implements MsgEncoder.EncodeTo
func (r Raw) EncodeTo(wr *Writer) (int, error) {
	return wr.Write([]byte(r))
}

// MarshalMsg implements MsgMarshaler.MarshalMsg
// (it makes a copy of the raw bytes)
func (r Raw) MarshalMsg() ([]byte, error) {
	o := make([]byte, len(r))
	copy(o, []byte(r))
	return o, nil
}

// AppendMsg implements MsgMarshaler.AppendMsg
// (it appends the raw bytes to 'b')
func (r Raw) AppendMsg(b []byte) ([]byte, error) {
	return append(b, []byte(r)...), nil
}

// MarshalJSON implements the json.Marshaler interface.
// (this translates the raw messagepack directly into JSON)
func (r Raw) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(3 * len(r) / 2)
	_, err := UnmarshalAsJSON(&buf, []byte(r))
	return buf.Bytes(), err
}
