package enc

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"unicode/utf8"
)

type kind byte

const (
	invalid kind = iota
	kmap
	kbool
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

var null = []byte("null")
var hex = []byte("0123456789abcdef")

// this is the interface
// used to write json
type jsWriter interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

// AsJSON reads MessagePack from src and
// translates it to JSON
func AsJSON(src io.Reader) io.Reader {
	rd, wr := io.Pipe()
	go func() {
		_, err := CopyToJSON(wr, src)
		if err != nil {
			rd.CloseWithError(err)
		} else {
			rd.Close()
		}
	}()
	return rd
}

// CopyToJson reads a single MessagePack-encoded message from 'src' and
// writes it as JSON to 'dst'. It returns the number of bytes written,
// and any errors encountered in the process.
func CopyToJSON(dst io.Writer, src io.Reader) (n int64, err error) {
	return DecodeToJSON(dst, NewDecoder(src))
}

func DecodeToJSON(dst io.Writer, src *MsgReader) (n int64, err error) {
	var w jsWriter
	var cast bool
	if jsw, ok := dst.(jsWriter); ok {
		w = jsw
		cast = true
	} else {
		w = bufio.NewWriterSize(dst, 256)
	}
	var k kind
	k, err = src.nextKind()
	if err != nil {
		return
	}
	var nn int
	switch k {
	case kmap:
		nn, err = rwMap(w, src)
	case karray:
		nn, err = rwArray(w, src)
	default:
		return 0, errors.New("enc: 'src' must represent a map or array")
	}
	n = int64(nn)
	if err != nil {
		return
	}
	if !cast {
		err = w.(*bufio.Writer).Flush()
	}
	return
}

func (m *MsgReader) nextKind() (kind, error) {
	var lead []byte
	var err error
	lead, err = m.r.Peek(1)
	if err != nil {
		return invalid, err
	}
	l := lead[0]
	if isfixmap(l) {
		return kmap, nil
	}
	if isfixint(l) {
		return kint, nil
	}
	if isnfixint(l) {
		return kint, nil
	}
	if isfixstr(l) {
		return kstring, nil
	}
	if isfixarray(l) {
		return karray, nil
	}
	switch l {
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
	case mfalse, mtrue:
		return kbool, nil
	default:
		return invalid, nil
	}
}

func rwMap(dst jsWriter, src *MsgReader) (n int, err error) {
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
	var nn int
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

		src.scratch, _, err = src.ReadStringAsBytes(src.scratch)
		if err != nil {
			return
		}

		nn, err = rwquoted(dst, src.scratch)
		n += nn
		if err != nil {
			return
		}
		err = dst.WriteByte(':')
		if err != nil {
			return
		}
		n++
		k, err = src.nextKind()
		if err != nil {
			return
		}
		switch k {
		case knull:
			nn, err = dst.Write(null)
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
		case kbool:
			nn, err = rwBool(dst, src)
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

func rwArray(dst jsWriter, src *MsgReader) (n int, err error) {
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
		case kbool:
			nn, err = rwBool(dst, src)
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

func rwFloat(dst jsWriter, src *MsgReader, k kind) (n int, err error) {
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

func rwInt(dst jsWriter, src *MsgReader, k kind) (n int, err error) {
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

func rwBool(dst jsWriter, src *MsgReader) (int, error) {
	b, _, err := src.ReadBool()
	if err != nil {
		return 0, err
	}
	if b {
		return dst.WriteString("true")
	}
	return dst.WriteString("false")
}

func rwExtension(dst jsWriter, src *MsgReader) (n int, err error) {
	var nn int
	err = dst.WriteByte('{')
	if err != nil {
		return
	}
	n += 1

	e := Extension{}

	e, nn, err = src.ReadExtension()
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

func rwString(dst jsWriter, src *MsgReader) (n int, err error) {
	src.scratch, _, err = src.ReadStringAsBytes(src.scratch)
	if err != nil {
		return
	}
	return rwquoted(dst, src.scratch)
}

func rwBytes(dst jsWriter, src *MsgReader) (n int, err error) {
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

// see: encoding/json/encode.go:(*encodeState).stringbytes()
func rwquoted(dst jsWriter, s []byte) (n int, err error) {
	var nn int
	err = dst.WriteByte('"')
	if err != nil {
		return
	}
	n += 1
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' {
				i++
				continue
			}
			if start < i {
				nn, err = dst.Write(s[start:i])
				n += nn
				if err != nil {
					return
				}
			}
			switch b {
			case '\\', '"':
				err = dst.WriteByte('\\')
				if err != nil {
					return
				}
				n++
				err = dst.WriteByte(b)
				if err != nil {
					return
				}
				n++
			case '\n':
				err = dst.WriteByte('\\')
				if err != nil {
					return
				}
				n++
				err = dst.WriteByte('n')
				if err != nil {
					return
				}
				n++
			case '\r':
				err = dst.WriteByte('\\')
				if err != nil {
					return
				}
				n++
				err = dst.WriteByte('r')
				if err != nil {
					return
				}
				n++
			default:
				nn, err = dst.WriteString(`\u00`)
				n += nn
				if err != nil {
					return
				}
				err = dst.WriteByte(hex[b>>4])
				if err != nil {
					return
				}
				n++
				err = dst.WriteByte(hex[b&0xF])
				if err != nil {
					return
				}
				n++
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				nn, err = dst.Write(s[start:i])
				n += nn
				if err != nil {
					return
				}
				nn, err = dst.WriteString(`\ufffd`)
				n += nn
				if err != nil {
					return
				}
				i += size
				start = i
				continue
			}
		}
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				nn, err = dst.Write(s[start:i])
				n += nn
				if err != nil {
					return
				}
				nn, err = dst.WriteString(`\u202`)
				n += nn
				if err != nil {
					return
				}
				err = dst.WriteByte(hex[c&0xF])
				if err != nil {
					return
				}
				n++
			}
		}
		i += size
	}
	if start < len(s) {
		nn, err = dst.Write(s[start:])
		n += nn
		if err != nil {
			return
		}
	}
	err = dst.WriteByte('"')
	if err != nil {
		return
	}
	n++
	return
}
