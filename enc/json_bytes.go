package enc

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
)

// UnmarshalAsJSON takes raw messagepack and writes
// it as JSON to 'w'.
func UnmarshalAsJSON(w io.Writer, msg []byte) ([]byte, error) {
	if jsw, ok := w.(jsWriter); ok {
		return writeNext(jsw, msg)
	}
	b := bufio.NewWriterSize(w, 512)
	msg, err := writeNext(b, msg)
	if err != nil {
		return msg, err
	}
	err = b.Flush()
	return msg, err
}

func writeNext(w jsWriter, msg []byte) ([]byte, error) {
	if len(msg) < 1 {
		return nil, ErrShortBytes
	}
	k := getKind(msg[0])
	switch k {
	case karray:
		return rwArrayBytes(w, msg)
	case kmap:
		return rwMapBytes(w, msg)
	case kbool:
		return rwBoolBytes(w, msg)
	case kint:
		return rwIntBytes(w, msg)
	case kuint:
		return rwUintBytes(w, msg)
	case kfloat32:
		return rwFloatBytes(w, msg, false)
	case kfloat64:
		return rwFloatBytes(w, msg, true)
	case knull:
		return rwNullBytes(w, msg)
	case kstring:
		return rwStringBytes(w, msg)
	case kbytes:
		return rwBytesBytes(w, msg)
	default:
		return msg, errors.New("bad encoding; invalid type")
	}
}

func rwArrayBytes(w jsWriter, msg []byte) ([]byte, error) {
	sz, msg, err := ReadArrayHeaderBytes(msg)
	if err != nil {
		return msg, err
	}
	err = w.WriteByte('[')
	if err != nil {
		return msg, err
	}
	for i := uint32(0); i < sz; i++ {
		if i != 0 {
			err = w.WriteByte(',')
			if err != nil {
				return msg, err
			}
		}
		msg, err = writeNext(w, msg)
		if err != nil {
			return msg, err
		}
	}
	err = w.WriteByte(']')
	return msg, err
}

func rwMapBytes(w jsWriter, msg []byte) ([]byte, error) {
	sz, msg, err := ReadMapHeaderBytes(msg)
	if err != nil {
		return msg, err
	}
	err = w.WriteByte('{')
	if err != nil {
		return msg, err
	}
	for i := uint32(0); i < sz; i++ {
		if i != 0 {
			err = w.WriteByte(',')
			if err != nil {
				return msg, err
			}
		}
		msg, err = rwStringBytes(w, msg)
		if err != nil {
			return msg, err
		}
		err = w.WriteByte(':')
		if err != nil {
			return msg, err
		}
		msg, err = writeNext(w, msg)
		if err != nil {
			return msg, err
		}
	}
	err = w.WriteByte('}')
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

func rwBytesBytes(w jsWriter, msg []byte) ([]byte, error) {
	bts, msg, err := ReadBytesZC(msg)
	if err != nil {
		return msg, err
	}
	s := make([]byte, base64.StdEncoding.EncodedLen(len(bts)))
	base64.StdEncoding.Encode(s, bts)
	err = w.WriteByte('"')
	if err != nil {
		return msg, err
	}
	_, err = w.Write(s)
	if err != nil {
		return msg, err
	}
	err = w.WriteByte('"')
	return msg, err
}

func rwNullBytes(w jsWriter, msg []byte) ([]byte, error) {
	msg, err := ReadNilBytes(msg)
	if err != nil {
		return msg, err
	}
	_, err = w.WriteString("null")
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

func rwIntBytes(w jsWriter, msg []byte) ([]byte, error) {
	i, msg, err := ReadInt64Bytes(msg)
	if err != nil {
		return msg, err
	}
	_, err = w.WriteString(strconv.FormatInt(i, 10))
	return msg, err
}

func rwUintBytes(w jsWriter, msg []byte) ([]byte, error) {
	u, msg, err := ReadUint64Bytes(msg)
	if err != nil {
		return msg, err
	}
	_, err = w.WriteString(strconv.FormatUint(u, 10))
	return msg, err
}

func rwFloatBytes(w jsWriter, msg []byte, f64 bool) ([]byte, error) {
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
		return msg, err
	}
	_, err = w.WriteString(strconv.FormatFloat(f, 'f', -1, sz))
	return msg, err
}
