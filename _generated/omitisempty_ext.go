package _generated

import (
	"github.com/tinylib/msgp/msgp"
)

// this has "external" types that will show up
// has generic IDENT during code generation

type OmitIsEmptyExt struct {
	a int // custom type
}

// IsEmpty will return true if a is positive
func (o OmitIsEmptyExt) IsEmpty() bool { return o.a <= 0 }

// EncodeMsg implements msgp.Encodable
func (o OmitIsEmptyExt) EncodeMsg(en *msgp.Writer) (err error) {
	if o.a > 0 {
		return en.WriteInt(o.a)
	}
	return en.WriteNil()
}

// DecodeMsg implements msgp.Decodable
func (o *OmitIsEmptyExt) DecodeMsg(dc *msgp.Reader) (err error) {
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		o.a = 0
		return
	}
	o.a, err = dc.ReadInt()
	return err
}

// MarshalMsg implements msgp.Marshaler
func (o OmitIsEmptyExt) MarshalMsg(b []byte) (ret []byte, err error) {
	ret = msgp.Require(b, o.Msgsize())
	if o.a > 0 {
		return msgp.AppendInt(ret, o.a), nil
	}
	return msgp.AppendNil(ret), nil
}

// UnmarshalMsg implements msgp.Unmarshaler
func (o *OmitIsEmptyExt) UnmarshalMsg(bts []byte) (ret []byte, err error) {
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		return bts, err
	}
	o.a, bts, err = msgp.ReadIntBytes(bts)
	return bts, err
}

// Msgsize implements msgp.Msgsizer
func (o OmitIsEmptyExt) Msgsize() (s int) {
	return msgp.IntSize
}

type OmitIsEmptyExtPtr struct {
	a int // custom type
}

// IsEmpty will return true if a is positive
func (o *OmitIsEmptyExtPtr) IsEmpty() bool { return o == nil || o.a <= 0 }

// EncodeMsg implements msgp.Encodable
func (o *OmitIsEmptyExtPtr) EncodeMsg(en *msgp.Writer) (err error) {
	if o.a > 0 {
		return en.WriteInt(o.a)
	}
	return en.WriteNil()
}

// DecodeMsg implements msgp.Decodable
func (o *OmitIsEmptyExtPtr) DecodeMsg(dc *msgp.Reader) (err error) {
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		o.a = 0
		return
	}
	o.a, err = dc.ReadInt()
	return err
}

// MarshalMsg implements msgp.Marshaler
func (o *OmitIsEmptyExtPtr) MarshalMsg(b []byte) (ret []byte, err error) {
	ret = msgp.Require(b, o.Msgsize())
	if o.a > 0 {
		return msgp.AppendInt(ret, o.a), nil
	}
	return msgp.AppendNil(ret), nil
}

// UnmarshalMsg implements msgp.Unmarshaler
func (o *OmitIsEmptyExtPtr) UnmarshalMsg(bts []byte) (ret []byte, err error) {
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		return bts, err
	}
	o.a, bts, err = msgp.ReadIntBytes(bts)
	return bts, err
}

// Msgsize implements msgp.Msgsizer
func (o *OmitIsEmptyExtPtr) Msgsize() (s int) {
	return msgp.IntSize
}
