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

// CopyToJson reads a single MsgPack-encoded message from 'src' and
// writes it as JSON to 'dst'. It returns the number of bytes written,
// and any errors encountered in the process.
func CopyToJSON(dst io.Writer, src io.Reader) (n int, err error) {
	w := bufio.NewWriterSize(dst, 16)
	r := NewReader(src)
	var k kind
	k, err = r.nextKind()
	if err != nil {
		return
	}
	switch k {
	case kmap:
		n, err = rwMap(w, r)
	case karray:
		n, err = rwArray(w, r)
	default:
		return 0, errors.New("enc: 'src' must represent a map or array")
	}
	if err != nil {
		return
	}
	err = w.Flush()
	return
}

func (m *MsgReader) nextKind() (kind, error) {
	var lead []byte
	var err error
	lead, err = m.r.Peek(1)
	if err != nil {
		return invalid, err
	}
	switch lead[0] {
	case mmap16, mmap32:
		return kmap, nil
	case marray16, marray32:
		return karray, nil
	case mfloat32:
		return kfloat32, nil
	case mfloat64:
		return kfloat64, nil
	case mint8, mint16, mint32, mint64:
		return kint, nil
	case muint8, muint16, muint32, muint64:
		return kuint, nil
	case mfixext1, mfixext2, mfixext4, mfixext8, mfixext16, mext8, mext16, mext32:
		return kextension, nil
	case mstr8, mstr16, mstr32:
		return kstring, nil
	case mbin8, mbin16, mbin32:
		return kbytes, nil
	case mnil:
		return knull, nil
	default:
		if isfixmap(lead[0]) {
			return kmap, nil
		}
		if isfixint(lead[0]) || isnfixint(lead[0]) {
			return kint, nil
		}
		if isfixstr(lead[0]) {
			return kstring, nil
		}
		if isfixarray(lead[0]) {
			return karray, nil
		}
		return invalid, nil
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

	err = dst.WriteByte('{')
	if err != nil {
		return
	}
	n += 1
	var s string
	var nn int
	var bts []byte
	for i := uint32(0); i < sz; i++ {
		if comma {
			err = dst.WriteByte(',')
			if err != nil {
				return
			}
			n += 1
		}

		var k kind
		k, err = src.nextKind()
		if err != nil {
			return
		}
		if k != kstring {
			return n, errors.New("map keys must be strings")
		}

		s, _, err = src.ReadString()
		if err != nil {
			return
		}
		bts = strconv.AppendQuote(bts[0:0], s)

		bts = append(bts, ':')
		nn, err = dst.Write(bts)
		n += nn
		if err != nil {
			return
		}

		k, err = src.nextKind()
		if err != nil {
			return
		}
		switch k {
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
		var k kind
		k, err = src.nextKind()
		if err != nil {
			return
		}
		switch k {
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
	if k == kfloat64 {
		var f float64
		f, _, err = src.ReadFloat64()
		if err != nil {
			return
		}
		src.scratch = strconv.AppendFloat(src.scratch[0:0], f, 'f', -1, 64)
	} else {
		var f float32
		f, _, err = src.ReadFloat32()
		if err != nil {
			return
		}
		src.scratch = strconv.AppendFloat(src.scratch[0:0], float64(f), 'f', -1, 32)
	}

	return dst.Write(src.scratch)
}

func rwInt(dst *bufio.Writer, src *MsgReader, k kind) (n int, err error) {
	if k == kint {
		var i int64
		i, _, err = src.ReadInt64()
		if err != nil {
			return
		}
		src.scratch = strconv.AppendInt(src.scratch[0:0], i, 10)
	} else {
		var u uint64
		u, _, err = src.ReadUint64()
		if err != nil {
			return
		}
		src.scratch = strconv.AppendUint(src.scratch[0:0], u, 10)
	}
	return dst.Write(src.scratch)
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

	src.scratch = strconv.AppendInt(src.scratch[0:0], int64(e.Type), 10)
	nn, err = dst.Write(src.scratch)
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
	var s []byte
	s, _, err = src.ReadStringAsBytes(src.scratch)
	if err != nil {
		return
	}
	src.scratch = strconv.AppendQuote(src.scratch[0:0], UnsafeString(s))
	return dst.Write(src.scratch)
}

func rwBytes(dst *bufio.Writer, src *MsgReader) (n int, err error) {
	var nn int
	err = dst.WriteByte('"')
	if err != nil {
		return
	}
	n += 1
	src.scratch, _, err = src.ReadBytes(src.scratch[0:0])
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
