package enc

import (
	"io"
)

// Raw is a []byte that is already messagepack-formatted
type Raw []byte

// Raw implements the MsgEncoder interface
func (r Raw) EncodeMsg(w io.Writer) (int, error) {
	return w.Write([]byte(r))
}

// Raw implements the MsgEncoder interface
func (r Raw) EncodeTo(wr *MsgWriter) (int, error) {
	return wr.Write([]byte(r))
}

// MarshalMsg implements the MsgEncoder interface
func (r Raw) MarshalMsg() ([]byte, error) {
	o := make([]byte, len(r))
	copy(o, []byte(r))
	return o, nil
}
