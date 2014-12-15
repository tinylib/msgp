package msgp

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"
	"strconv"
	"time"
	"unicode/utf8"
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

// CopyToJSON reads a single MessagePack-encoded message from 'src' and
// writes it as JSON to 'dst' until 'src' returns EOF. It returns the
// number of bytes written and any errors encountered.
func CopyToJSON(dst io.Writer, src io.Reader) (n int64, err error) {
	r := NewReader(src)
	n, err = r.WriteToJSON(dst)
	FreeR(r)
	return
}

// WriteToJSON translates MessagePack from 'r' and writes it as
// JSON to 'w' until the underlying reader returns io.EOF. It returns
// the number of bytes written, and an error if it stopped before EOF.
func (r *Reader) WriteToJSON(w io.Writer) (n int64, err error) {
	var j jsWriter
	var cast bool
	if jsw, ok := w.(jsWriter); ok {
		j = jsw
		cast = true
	} else {
		j = bufio.NewWriterSize(w, 512)
	}
	var nn int
	for err == nil {
		nn, err = rwNext(j, r)
		n += int64(nn)
	}
	if err != io.EOF {
		if !cast {
			j.(*bufio.Writer).Flush()
		}
		return
	}
	err = nil
	if !cast {
		err = j.(*bufio.Writer).Flush()
	}
	return
}

func rwNext(w jsWriter, src *Reader) (int, error) {
	t, err := src.NextType()
	if err != nil {
		return 0, err
	}
	switch t {
	case NilType:
		return w.Write(null)
	case BoolType:
		return rwBool(w, src)
	case IntType, UintType:
		return rwInt(w, src, t)
	case StrType:
		return rwString(w, src)
	case BinType:
		return rwBytes(w, src)
	case ExtensionType, Complex64Type, Complex128Type, TimeType:
		return rwExtension(w, src, t)
	case ArrayType:
		return rwArray(w, src)
	case Float32Type, Float64Type:
		return rwFloat(w, src, t)
	case MapType:
		return rwMap(w, src)
	default:
		// we shouldn't get here; NextType() errors if the type is not recognized
		return 0, fatal
	}
}

func rwMap(dst jsWriter, src *Reader) (n int, err error) {
	var comma bool
	var sz uint32

	sz, err = src.ReadMapHeader()
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
	n++
	var nn int
	for i := uint32(0); i < sz; i++ {
		if comma {
			err = dst.WriteByte(',')
			if err != nil {
				return
			}
			n++
		}

		nn, err = rwString(dst, src)
		n += nn
		if err != nil {
			if tperr, ok := err.(TypeError); ok && tperr.Encoded == BinType {
				nn, err = rwBytes(dst, src)
				n += nn
				if err != nil {
					return
				}
			} else {
				return
			}
		}

		err = dst.WriteByte(':')
		if err != nil {
			return
		}
		n++
		nn, err = rwNext(dst, src)
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
	n++
	return
}

func rwArray(dst jsWriter, src *Reader) (n int, err error) {
	err = dst.WriteByte('[')
	if err != nil {
		return
	}
	var sz uint32
	var nn int
	sz, err = src.ReadArrayHeader()
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
			n++
		}
		nn, err = rwNext(dst, src)
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
	n++
	return
}

func rwFloat(dst jsWriter, src *Reader, t Type) (n int, err error) {
	if t == Float64Type {
		var f float64
		f, err = src.ReadFloat64()
		if err != nil {
			return
		}
		src.scratch = strconv.AppendFloat(src.scratch[0:0], f, 'f', -1, 64)
	} else {
		var f float32
		f, err = src.ReadFloat32()
		if err != nil {
			return
		}
		src.scratch = strconv.AppendFloat(src.scratch[0:0], float64(f), 'f', -1, 32)
	}
	return dst.Write(src.scratch)
}

func rwInt(dst jsWriter, src *Reader, t Type) (n int, err error) {
	if t == IntType {
		var i int64
		i, err = src.ReadInt64()
		if err != nil {
			return
		}
		src.scratch = strconv.AppendInt(src.scratch[0:0], i, 10)
	} else {
		var u uint64
		u, err = src.ReadUint64()
		if err != nil {
			return
		}
		src.scratch = strconv.AppendUint(src.scratch[0:0], u, 10)
	}
	return dst.Write(src.scratch)
}

func rwBool(dst jsWriter, src *Reader) (int, error) {
	b, err := src.ReadBool()
	if err != nil {
		return 0, err
	}
	if b {
		return dst.WriteString("true")
	}
	return dst.WriteString("false")
}

func rwExtension(dst jsWriter, src *Reader, t Type) (n int, err error) {
	// time.Time is a json.Marshaler
	if t == TimeExtension {
		var t time.Time
		var bts []byte
		t, err = src.ReadTime()
		if err != nil {
			return
		}
		bts, err = t.MarshalJSON()
		if err != nil {
			return
		}
		return dst.Write(bts)
	}

	var et int8
	et, err = src.peekExtensionType()
	if err != nil {
		return 0, err
	}

	// registered extensions can override
	// the JSON encoding
	if j, ok := extensionReg[et]; ok {
		var bts []byte
		e := j()
		err = src.ReadExtension(e)
		if err != nil {
			return
		}
		bts, err = json.Marshal(e)
		if err != nil {
			return
		}
		return dst.Write(bts)
	}

	e := RawExtension{}
	e.Type = et
	err = src.ReadExtension(&e)
	if err != nil {
		return
	}

	var nn int
	err = dst.WriteByte('{')
	if err != nil {
		return
	}
	n++

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

func rwString(dst jsWriter, src *Reader) (n int, err error) {
	var p []byte
	p, err = src.r.Peek(1)
	if err != nil {
		return
	}
	lead := p[0]
	var read int

	if isfixstr(lead) {
		read = int(rfixstr(lead))
		src.r.Skip(1)
		goto write
	}

	switch lead {
	case mstr8:
		p, err = src.r.Next(2)
		if err != nil {
			return
		}
		read = int(uint8(p[1]))
	case mstr16:
		p, err = src.r.Next(3)
		if err != nil {
			return
		}
		read = int(big.Uint16(p[1:]))
	case mstr32:
		p, err = src.r.Next(5)
		if err != nil {
			return
		}
		read = int(big.Uint32(p[1:]))
	default:
		err = badPrefix(StrType, lead)
		return
	}
write:
	p, err = src.r.Next(read)
	if err != nil {
		return
	}
	n, err = rwquoted(dst, p)
	return
}

func rwBytes(dst jsWriter, src *Reader) (n int, err error) {
	var nn int
	err = dst.WriteByte('"')
	if err != nil {
		return
	}
	n++
	src.scratch, err = src.ReadBytes(src.scratch[0:0])
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
	n++
	return
}

// see: encoding/json/encode.go:(*encodeState).stringbytes()
func rwquoted(dst jsWriter, s []byte) (n int, err error) {
	var nn int
	err = dst.WriteByte('"')
	if err != nil {
		return
	}
	n++
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
