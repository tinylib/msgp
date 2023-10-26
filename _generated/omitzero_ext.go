package _generated

import (
	"github.com/tinylib/msgp/msgp"
)

// this has "external" types that will show up
// as generic IDENT during code generation

type OmitZeroExt struct {
	a int // custom type
}

// IsZero will return true if a is not positive
func (o OmitZeroExt) IsZero() bool { return o.a <= 0 }

// EncodeMsg implements msgp.Encodable
func (o OmitZeroExt) EncodeMsg(en *msgp.Writer) (err error) {
	if o.a > 0 {
		return en.WriteInt(o.a)
	}
	return en.WriteNil()
}

// DecodeMsg implements msgp.Decodable
func (o *OmitZeroExt) DecodeMsg(dc *msgp.Reader) (err error) {
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
func (o OmitZeroExt) MarshalMsg(b []byte) (ret []byte, err error) {
	ret = msgp.Require(b, o.Msgsize())
	if o.a > 0 {
		return msgp.AppendInt(ret, o.a), nil
	}
	return msgp.AppendNil(ret), nil
}

// UnmarshalMsg implements msgp.Unmarshaler
func (o *OmitZeroExt) UnmarshalMsg(bts []byte) (ret []byte, err error) {
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		return bts, err
	}
	o.a, bts, err = msgp.ReadIntBytes(bts)
	return bts, err
}

// Msgsize implements msgp.Msgsizer
func (o OmitZeroExt) Msgsize() (s int) {
	return msgp.IntSize
}

type OmitZeroExtPtr struct {
	a int // custom type
}

// IsZero will return true if a is nil or not positive
func (o *OmitZeroExtPtr) IsZero() bool { return o == nil || o.a <= 0 }

// EncodeMsg implements msgp.Encodable
func (o *OmitZeroExtPtr) EncodeMsg(en *msgp.Writer) (err error) {
	if o.a > 0 {
		return en.WriteInt(o.a)
	}
	return en.WriteNil()
}

// DecodeMsg implements msgp.Decodable
func (o *OmitZeroExtPtr) DecodeMsg(dc *msgp.Reader) (err error) {
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
func (o *OmitZeroExtPtr) MarshalMsg(b []byte) (ret []byte, err error) {
	ret = msgp.Require(b, o.Msgsize())
	if o.a > 0 {
		return msgp.AppendInt(ret, o.a), nil
	}
	return msgp.AppendNil(ret), nil
}

// UnmarshalMsg implements msgp.Unmarshaler
func (o *OmitZeroExtPtr) UnmarshalMsg(bts []byte) (ret []byte, err error) {
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		return bts, err
	}
	o.a, bts, err = msgp.ReadIntBytes(bts)
	return bts, err
}

// Msgsize implements msgp.Msgsizer
func (o *OmitZeroExtPtr) Msgsize() (s int) {
	return msgp.IntSize
}
