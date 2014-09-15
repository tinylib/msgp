package enc

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
)

type kind byte

const (
	invalid kind = iota
	kmap
	karray
	kint
	kuint
	kfloat32
	kfloat64
	kextension
	kbytes
	kstring
	knull
)

// CopyToJson reads MessagePack-encoded elements from 'src' and
// writes them as json to 'dst'.
func CopyToJSON(dst io.Writer, src io.Reader) (n int, err error) {
	w := bufio.NewWriter(dst)
	r := NewReader(src)

	switch r.nextKind() {
	case kmap:
		return rwArray(w, r)
	case karray:
		return rwMap(w, r)
	default:
		return 0, errors.New("enc: 'src' must represent a map or array")
	}
}

func (m *MsgReader) nextKind() kind {
	var lead []byte
	lead, _ = m.r.Peek(1)
	if len(lead) == 0 {
		return invalid
	}
	switch lead[0] {
	case mmap16, mmap32:
		return kmap
	case marray16, marray32:
		return karray
	case mfloat32:
		return kfloat32
	case mfloat64:
		return kfloat64
	case mint8, mint16, mint32, mint64:
		return kint
	case muint8, muint16, muint32, muint64:
		return kuint
	case mfixext1, mfixext2, mfixext4, mfixext8, mfixext16, mext8, mext16, mext32:
		return kextension
	case mstr8, mstr16, mstr32:
		return kstring
	case mbin8, mbin16, mbin32:
		return kbytes
	case mnil:
		return knull
	default:

		if lead[0]&mfixint == mfixint {
			return kint
		}
		if lead[0]&mfixstr == mfixstr {
			return kstring
		}
		if lead[0]&mnfixint == mnfixint {
			return kint
		}
		if lead[0]&mfixarray == mfixarray {
			return karray
		}
		if lead[0]&mfixmap == mfixmap {
			return kmap
		}

		return invalid
	}
}

func rwMap(dst *bufio.Writer, src *MsgReader) (n int, err error) {
	var comma bool
	var sz uint32

	sz, _, err = src.ReadMapHeader()
	if err != nil {
		return
	}

	if sz == 0 {
		return dst.WriteString("{}")
	}

	var scratch []byte
	err = dst.WriteByte('{')
	if err != nil {
		return
	}
	n += 1
	var s string
	var nn int
	for i := uint32(0); i < sz; i++ {
		if comma {
			err = dst.WriteByte(',')
			if err != nil {
				return
			}
			n += 1
		}

		s, _, err = src.ReadString()
		scratch = strconv.AppendQuote(scratch[0:0], s)

		scratch = append(scratch, ':')
		nn, err = dst.Write(scratch)
		n += nn
		if err != nil {
			return
		}

		switch src.nextKind() {
		case knull:
			nn, err = dst.WriteString("null")
		case kint:
			nn, err = rwInt(dst, src, kint)
		case kuint:
			nn, err = rwInt(dst, src, kuint)
		case kstring:
			nn, err = rwString(dst, src)
		case kbytes:
			nn, err = rwBytes(dst, src)
		case kfloat32:
			nn, err = rwFloat(dst, src, kfloat32)
		case kfloat64:
			nn, err = rwFloat(dst, src, kfloat64)
		case karray:
			nn, err = rwArray(dst, src)
		case kmap:
			nn, err = rwMap(dst, src)
		case kextension:
			nn, err = rwExtension(dst, src)
		default:
			return n, errors.New("bad token in src")
		}
		n += nn
		if err != nil {
			return
		}
		if !comma {
			comma = true
		}
	}

	err = dst.WriteByte('}')
	if err != nil {
		return
	}
	n += 1
	return
}

func rwArray(dst *bufio.Writer, src *MsgReader) (n int, err error) {
	err = dst.WriteByte('[')
	if err != nil {
		return
	}
	var sz uint32
	var nn int
	sz, _, err = src.ReadArrayHeader()
	if err != nil {
		return
	}
	var comma bool
	for i := uint32(0); i < sz; i++ {
		if comma {
			err = dst.WriteByte(',')
			if err != nil {
				return
			}
			n += 1
		}

		switch src.nextKind() {
		case knull:
			nn, err = dst.WriteString("null")
		case kint:
			nn, err = rwInt(dst, src, kint)
		case kuint:
			nn, err = rwInt(dst, src, kuint)
		case kstring:
			nn, err = rwString(dst, src)
		case kbytes:
			nn, err = rwBytes(dst, src)
		case kfloat32:
			nn, err = rwFloat(dst, src, kfloat32)
		case kfloat64:
			nn, err = rwFloat(dst, src, kfloat64)
		case karray:
			nn, err = rwArray(dst, src)
		case kmap:
			nn, err = rwMap(dst, src)
		case kextension:
			nn, err = rwExtension(dst, src)
		default:
			return n, errors.New("bad token in src")
		}

		n += nn
		if err != nil {
			return
		}

		if !comma {
			comma = true
		}
	}

	err = dst.WriteByte(']')
	if err != nil {
		return
	}
	n += 1
	return
}

func rwFloat(dst *bufio.Writer, src *MsgReader, k kind) (n int, err error) {
	var bts []byte
	if k == kfloat32 {
		var f float64
		f, _, err = src.ReadFloat64()
		if err != nil {
			return
		}
		bts = strconv.AppendFloat(src.leader[0:0], f, 'f', -1, 64)
	} else {
		var f float32
		f, _, err = src.ReadFloat32()
		if err != nil {
			return
		}
		bts = strconv.AppendFloat(src.leader[0:0], float64(f), 'f', -1, 32)
	}

	return dst.Write(bts)
}

func rwInt(dst *bufio.Writer, src *MsgReader, k kind) (n int, err error) {
	var bts []byte
	if k == kint {
		var i int64
		i, _, err = src.ReadInt64()
		if err != nil {
			return
		}
		bts = strconv.AppendInt(src.leader[0:0], i, 10)
	} else {
		var u uint64
		u, _, err = src.ReadUint64()
		if err != nil {
			return
		}
		bts = strconv.AppendUint(src.leader[0:0], u, 10)
	}
	return dst.Write(bts)
}

func rwExtension(dst *bufio.Writer, src *MsgReader) (n int, err error) {
	var nn int
	err = dst.WriteByte('{')
	if err != nil {
		return
	}
	n += 1

	e := Extension{}

	nn, err = src.ReadExtension(&e)
	n += nn
	if err != nil {
		return
	}

	nn, err = dst.WriteString(`"type:"`)
	n += nn
	if err != nil {
		return
	}

	tb := strconv.AppendInt(src.leader[:], int64(e.Type), 10)
	nn, err = dst.Write(tb)
	n += nn
	if err != nil {
		return
	}

	nn, err = dst.WriteString(`,"data":"`)
	n += nn
	if err != nil {
		return
	}
	enc := base64.NewEncoder(base64.StdEncoding, dst)

	nn, err = enc.Write(e.Data)
	n += nn
	if err != nil {
		return
	}
	err = enc.Close()
	if err != nil {
		return
	}
	nn, err = dst.WriteString(`"}`)
	n += nn
	return
}

func rwString(dst *bufio.Writer, src *MsgReader) (n int, err error) {
	var s string
	s, _, err = src.ReadString()
	if err != nil {
		return
	}
	src.scratch = strconv.AppendQuote(src.scratch, s)
	return dst.Write(src.scratch)
}

func rwBytes(dst *bufio.Writer, src *MsgReader) (n int, err error) {
	var nn int
	err = dst.WriteByte('"')
	if err != nil {
		return
	}
	n += 1
	src.scratch, _, err = src.ReadBytes(src.scratch)
	if err != nil {
		return
	}
	enc := base64.NewEncoder(base64.StdEncoding, dst)
	nn, err = enc.Write(src.scratch)
	n += nn
	if err != nil {
		return
	}
	err = enc.Close()
	if err != nil {
		return
	}
	err = dst.WriteByte('"')
	if err != nil {
		return
	}
	n += 1
	return
}
