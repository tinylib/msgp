package msgp

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"
	"strconv"
	"time"
)

// UnmarshalAsJSON takes raw messagepack and writes
// it as JSON to 'w'. If an error is returned, the
// bytes not translated will also be returned. If
// no errors are encountered, the length of the returned
// slice will be zero.
func UnmarshalAsJSON(w io.Writer, msg []byte) ([]byte, error) {
	var (
		scratch []byte
		cast    bool
		dst     jsWriter
		err     error
	)
	if jsw, ok := w.(jsWriter); ok {
		dst = jsw
		cast = true
	} else {
		dst = bufio.NewWriterSize(w, 512)
	}
	for len(msg) > 0 && err == nil {
		msg, scratch, err = writeNext(dst, msg, scratch)
	}
	if !cast && err == nil {
		err = dst.(*bufio.Writer).Flush()
	}
	return msg, err
}

func writeNext(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) {
	if len(msg) < 1 {
		return nil, scratch, ErrShortBytes
	}
	t := getType(msg[0])
	var err error
	switch t {
	case ArrayType:
		return rwArrayBytes(w, msg, scratch)
	case MapType:
		return rwMapBytes(w, msg, scratch)
	case BoolType:
		msg, err = rwBoolBytes(w, msg)
		return msg, scratch, err
	case IntType:
		return rwIntBytes(w, msg, scratch)
	case UintType:
		return rwUintBytes(w, msg, scratch)
	case Float32Type:
		return rwFloatBytes(w, msg, false, scratch)
	case Float64Type:
		return rwFloatBytes(w, msg, true, scratch)
	case NilType:
		msg, err = rwNullBytes(w, msg)
		return msg, scratch, err
	case StrType:
		msg, err = rwStringBytes(w, msg)
		return msg, scratch, err
	case BinType:
		return rwBytesBytes(w, msg, scratch)
	case ExtensionType:
		return rwExtensionBytes(w, msg, scratch)
	default:
		return msg, scratch, InvalidPrefixError(msg[0])
	}
}

func rwArrayBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) {
	sz, msg, err := ReadArrayHeaderBytes(msg)
	if err != nil {
		return msg, scratch, err
	}
	err = w.WriteByte('[')
	if err != nil {
		return msg, scratch, err
	}
	for i := uint32(0); i < sz; i++ {
		if i != 0 {
			err = w.WriteByte(',')
			if err != nil {
				return msg, scratch, err
			}
		}
		msg, scratch, err = writeNext(w, msg, scratch)
		if err != nil {
			return msg, scratch, err
		}
	}
	err = w.WriteByte(']')
	return msg, scratch, err
}

func rwMapBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) {
	sz, msg, err := ReadMapHeaderBytes(msg)
	if err != nil {
		return msg, scratch, err
	}
	err = w.WriteByte('{')
	if err != nil {
		return msg, scratch, err
	}
	for i := uint32(0); i < sz; i++ {
		if i != 0 {
			err = w.WriteByte(',')
			if err != nil {
				return msg, scratch, err
			}
		}
		msg, err = rwMapKeyBytes(w, msg)
		if err != nil {
			return msg, scratch, err
		}
		err = w.WriteByte(':')
		if err != nil {
			return msg, scratch, err
		}
		msg, scratch, err = writeNext(w, msg, scratch)
		if err != nil {
			return msg, scratch, err
		}
	}
	err = w.WriteByte('}')
	return msg, scratch, err
}

func rwMapKeyBytes(w jsWriter, msg []byte) ([]byte, error) {
	str, msg, err := ReadMapKeyZC(msg)
	if err != nil {
		return msg, err
	}
	_, err = rwquoted(w, str)
	return msg, err
}

func rwStringBytes(w jsWriter, msg []byte) ([]byte, error) {
	str, msg, err := ReadStringZC(msg)
	if err != nil {
		return msg, err
	}
	_, err = rwquoted(w, str)
	return msg, err
}

func rwBytesBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) {
	bts, msg, err := ReadBytesZC(msg)
	if err != nil {
		return msg, scratch, err
	}
	l := base64.StdEncoding.EncodedLen(len(bts))
	if cap(scratch) >= l {
		scratch = scratch[0:l]
	} else {
		scratch = make([]byte, l)
	}
	base64.StdEncoding.Encode(scratch, bts)
	err = w.WriteByte('"')
	if err != nil {
		return msg, scratch, err
	}
	_, err = w.Write(scratch)
	if err != nil {
		return msg, scratch, err
	}
	err = w.WriteByte('"')
	return msg, scratch, err
}

func rwNullBytes(w jsWriter, msg []byte) ([]byte, error) {
	msg, err := ReadNilBytes(msg)
	if err != nil {
		return msg, err
	}
	_, err = w.Write(null)
	return msg, err
}

func rwBoolBytes(w jsWriter, msg []byte) ([]byte, error) {
	b, msg, err := ReadBoolBytes(msg)
	if err != nil {
		return msg, err
	}
	if b {
		_, err = w.WriteString("true")
		return msg, err
	}
	_, err = w.WriteString("false")
	return msg, err
}

func rwIntBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) {
	i, msg, err := ReadInt64Bytes(msg)
	if err != nil {
		return msg, scratch, err
	}
	scratch = strconv.AppendInt(scratch[0:0], i, 10)
	_, err = w.Write(scratch)
	return msg, scratch, err
}

func rwUintBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) {
	u, msg, err := ReadUint64Bytes(msg)
	if err != nil {
		return msg, scratch, err
	}
	scratch = strconv.AppendUint(scratch[0:0], u, 10)
	_, err = w.Write(scratch)
	return msg, scratch, err
}

func rwFloatBytes(w jsWriter, msg []byte, f64 bool, scratch []byte) ([]byte, []byte, error) {
	var f float64
	var err error
	var sz int
	if f64 {
		sz = 64
		f, msg, err = ReadFloat64Bytes(msg)
	} else {
		sz = 32
		var v float32
		v, msg, err = ReadFloat32Bytes(msg)
		f = float64(v)
	}
	if err != nil {
		return msg, scratch, err
	}
	scratch = strconv.AppendFloat(scratch, f, 'f', -1, sz)
	_, err = w.Write(scratch)
	return msg, scratch, err
}

func rwExtensionBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) {
	var err error
	var et int8
	et, err = peekExtension(msg)
	if err != nil {
		return msg, scratch, err
	}

	// if it's time.Time
	if et == TimeExtension {
		var tm time.Time
		tm, msg, err = ReadTimeBytes(msg)
		if err != nil {
			return msg, scratch, err
		}
		bts, err := tm.MarshalJSON()
		if err != nil {
			return msg, scratch, err
		}
		_, err = w.Write(bts)
		return msg, scratch, err
	}

	// if the extension is registered,
	// use its canonical JSON form
	if f, ok := extensionReg[et]; ok {
		e := f()
		msg, err = ReadExtensionBytes(msg, e)
		if err != nil {
			return msg, scratch, err
		}
		bts, err := json.Marshal(e)
		if err != nil {
			return msg, scratch, err
		}
		_, err = w.Write(bts)
		return msg, scratch, err
	}

	// otherwise, write `{"type": <num>, "data": "<base64data>"}`
	r := RawExtension{}
	r.Type = et
	msg, err = ReadExtensionBytes(msg, &r)
	if err != nil {
		return msg, scratch, err
	}
	scratch, err = writeExt(w, r, scratch)
	return msg, scratch, err
}

func writeExt(w jsWriter, r RawExtension, scratch []byte) ([]byte, error) {
	_, err := w.WriteString(`{"type":`)
	if err != nil {
		return scratch, err
	}
	scratch = strconv.AppendInt(scratch[0:0], int64(r.Type), 10)
	_, err = w.Write(scratch)
	if err != nil {
		return scratch, err
	}
	_, err = w.WriteString(`,"data":"`)
	if err != nil {
		return scratch, err
	}
	l := base64.StdEncoding.EncodedLen(len(r.Data))
	if cap(scratch) >= l {
		scratch = scratch[0:l]
	} else {
		scratch = make([]byte, l)
	}
	base64.StdEncoding.Encode(scratch, r.Data)
	_, err = w.Write(scratch)
	if err != nil {
		return scratch, err
	}
	_, err = w.WriteString(`"}`)
	return scratch, err
}
