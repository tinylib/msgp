package enc

import (
	"io"
)

// Raw is a []byte that is already messagepack-formatted
type Raw []byte

// Raw implements MsgEncoder.EncodeMsg
func (r Raw) EncodeMsg(w io.Writer) (int, error) {
	return w.Write([]byte(r))
}

// Raw implements MsgEncoder.EncodeTo
func (r Raw) EncodeTo(wr *MsgWriter) (int, error) {
	return wr.Write([]byte(r))
}

// MarshalMsg implements MsgEncoder.MarshalMsg
func (r Raw) MarshalMsg() ([]byte, error) {
	o := make([]byte, len(r))
	copy(o, []byte(r))
	return o, nil
}

// AppendMsg implements MsgEncoder.Append
func (r Raw) AppendMsg(b []byte) ([]byte, error) {
	return append(b, []byte(r)...), nil
}
