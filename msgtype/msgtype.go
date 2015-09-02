// package msgtype defines a number
// of useful types for use with
// the msgp MessagePack library.
package msgtype

import (
	"github.com/tinylib/msgp/msgp"
	"io"
)

// IOvec is a scatter list that
// is written and read as a contiguous
// MessagePack 'bin' object, but is
// represented in memory as multiple
// blocks of memory. It can be used
// to prevent a collection of memory
// regions from being re-allocated
// and copied in order to represent
// them as one object on the wire.
type IOvec [][]byte

func (iv IOvec) Msgsize() int {
	return msgp.BytesPrefixSize + iv.size()
}

func (iv IOvec) size() (sz int) {
	for _, v := range iv {
		sz += len(v)
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (iv *IOvec) DecodeMsg(r *msgp.Reader) error {
	sz, err := r.ReadBytesHeader()
	if err != nil {
		return err
	}
	// fill existing slices
	for i, v := range *iv {
		if sz > 0 {
			c := cap(v)
			if c > int(sz) {
				c = int(sz)
			}
			(*iv)[i] = (*iv)[i][:c]
			_, err := r.ReadFull((*iv)[i])
			if err != nil {
				return err
			}
			sz -= uint32(c)
		} else {
			(*iv)[i] = (*iv)[i][:0]
		}
	}
	// push any remaining data
	// onto the end of the iovec
	if sz > 0 {
		add := make([]byte, sz)
		_, err := r.ReadFull(add)
		if err != nil {
			return err
		}
		*iv = append(*iv, add)
	}
	return nil
}

// EncodeMsg implements msgp.Encodable
func (iv IOvec) EncodeMsg(w *msgp.Writer) error {
	err := w.WriteBytesHeader(uint32(iv.size()))
	if err != nil {
		return err
	}
	for _, v := range iv {
		if len(v) > 0 {
			if _, err := w.Write(v); err != nil {
				return err
			}
		}
	}
	return nil
}

// MarshalMsg implements msgp.Marshaler
func (iv IOvec) MarshalMsg(in []byte) ([]byte, error) {
	sz := iv.size()
	out := msgp.Require(in, msgp.BytesPrefixSize+sz)
	out = msgp.AppendBytesHeader(in, uint32(sz))
	for _, v := range iv {
		if len(v) > 0 {
			out = append(out, v...)
		}
	}
	return out, nil
}

// UnmarshalMsg implements msgp.Unmarshaler
func (iv *IOvec) UnmarshalMsg(in []byte) ([]byte, error) {
	sz, in, err := msgp.ReadBytesHeader(in)
	if err != nil {
		return nil, err
	}
	if sz > uint32(len(in)) {
		return nil, msgp.ErrShortBytes
	}
	for i := range *iv {
		l := uint32(cap((*iv)[i]))
		if l >= sz {
			(*iv)[i] = (*iv)[i][:sz]
			in = in[:copy((*iv)[i], in)]
			for j := range (*iv)[i+1:] {
				(*iv)[j+i+1] = (*iv)[j+i+1][:0]
			}
			return in, nil
		} else {
			(*iv)[i] = (*iv)[i][:l]
			in = in[:copy((*iv)[i], in)]
			sz -= l
		}
	}
	if sz > 0 {
		last := make([]byte, sz)
		in = in[:copy(last, in)]
		*iv = append(*iv, last)
	}
	return in, nil
}

// BytesSink is a destination into
// which 'bin'-encoded msgpack can
// be written.
type BytesSink struct {
	io.Writer
}

// DecodeMsg implements msgp.Decodable
func (b BytesSink) DecodeMsg(r *msgp.Reader) error {
	sz, err := r.ReadBytesHeader()
	if err != nil {
		return err
	}
	_, err = io.CopyN(b.Writer, r.R, int64(sz))
	return err
}

// UnmarshalMsg implements msgp.Unmarshaler
func (b BytesSink) UnmarshalMsg(in []byte) ([]byte, error) {
	sz, o, err := msgp.ReadBytesHeader(in)
	if err != nil {
		return in, err
	}
	if uint32(len(o)) < sz {
		return in, msgp.ErrShortBytes
	}
	_, err = b.Writer.Write(o[:sz])
	if err != nil {
		return in, err
	}
	return o[sz:], nil
}

// BytesSource is a source from which
// 'bin'-encoded msgpack can be read.
type BytesSource struct {
	R io.Reader // R is the reader from which data is read
	N uint32    // N is the number of bytes to read
}

func (b *BytesSource) EncodeMsg(w *msgp.Writer) error {
	err := w.WriteBytesHeader(b.N)
	if err != nil {
		return err
	}
	_, err = w.ReadFrom(&io.LimitedReader{b.R, int(b.N)})
	return err
}

func (b *BytesSource) MarshalMsg(d []byte) (o []byte, err error) {
	o = msgp.Require(d, b.N+msgp.BytesPrefixSize)
	o = msgp.AppendBytesHeader(d, b.N)
	n := len(o)
	o = o[:n+len(d)]
	_, err = io.ReadFull(b.R, o[n:])
	return
}
