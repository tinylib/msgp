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
