package _generated

import(
	"github.com/philhofer/msgp/enc"
	"io"
	"bytes"

)

func (z *TestType) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	_, err := z.EncodeMsg(&buf)
	return buf.Bytes(), err
}

func (z *TestType) Unmarshal(b []byte) error {
	_, err := z.DecodeMsg(bytes.NewReader(b))
	return err
}

func (z *TestType) EncodeMsg(w io.Writer) (n int, err error) {
	var nn int
	en := enc.NewEncoder(w)
	_ = nn
	_ = en

	if z == nil {
		nn, err = en.WriteNil()
		n += nn
		if err != nil {
			return
		}
	} else {

		nn, err = en.WriteMapHeader(uint32(4))
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteString("float")
		n += nn
		if err != nil {
			return
		}

		if z.F == nil {
			nn, err = en.WriteNil()
			n += nn
			if err != nil {
				return
			}
		} else {

			nn, err = en.WriteFloat64(*z.F)
			n += nn
			if err != nil {
				return
			}
		}
		nn, err = en.WriteString("elements")
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteMapStrStr(z.Els)
		n += nn
		if err != nil {
			return
		}
		nn, err = en.WriteString("object")
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteMapHeader(uint32(2))
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteString("value_a")
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteString(z.Obj.ValueA)
		n += nn
		if err != nil {
			return
		}
		nn, err = en.WriteString("value_b")
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteBytes(z.Obj.ValueB)
		n += nn
		if err != nil {
			return
		}
		nn, err = en.WriteString("child")
		n += nn
		if err != nil {
			return
		}

		if z.Child == nil {
			nn, err = en.WriteNil()
			n += nn
			if err != nil {
				return
			}
		} else {

			nn, err = en.WriteIdent(z.Child)
			n += nn
			if err != nil {
				return
			}
		}
	}
	return
}

func (z *TestType) DecodeMsg(r io.Reader) (n int, err error) {
	var sz uint32
	var nn int
	var field []byte
	dc := enc.NewDecoder(r)
	_ = sz
	_ = nn
	_ = field

	if dc.IsNil() {
		nn, err = dc.ReadNil()
		n += nn
		if err != nil {
			return
		}
		z = nil
	} else {
		if z == nil {
			z = new(TestType)
		}

		sz, nn, err = dc.ReadMapHeader()
		n += nn
		if err != nil {
			return
		}
		for xplz := uint32(0); xplz < sz; xplz++ {
			field, nn, err = dc.ReadStringAsBytes(field)
			n += nn
			if err != nil {
				return
			}
			switch enc.UnsafeString(field) {

			case "float":
				if dc.IsNil() {
					nn, err = dc.ReadNil()
					n += nn
					if err != nil {
						return
					}
					z.F = nil
				} else {
					if z.F == nil {
						z.F = new(float64)
					}

					*z.F, nn, err = dc.ReadFloat64()

					n += nn
					if err != nil {
						return
					}

				}

			case "elements":

				if z.Els == nil {
					z.Els = make(map[string]string)
				}
				nn, err = dc.ReadMapStrStr(z.Els)

				n += nn
				if err != nil {
					return
				}

			case "object":
				sz, nn, err = dc.ReadMapHeader()
				n += nn
				if err != nil {
					return
				}
				for xplz := uint32(0); xplz < sz; xplz++ {
					field, nn, err = dc.ReadStringAsBytes(field)
					n += nn
					if err != nil {
						return
					}
					switch enc.UnsafeString(field) {

					case "value_a":

						z.Obj.ValueA, nn, err = dc.ReadString()

						n += nn
						if err != nil {
							return
						}

					case "value_b":

						z.Obj.ValueB, nn, err = dc.ReadBytes(z.Obj.ValueB)

						n += nn
						if err != nil {
							return
						}

					}
				}

			case "child":
				if dc.IsNil() {
					nn, err = dc.ReadNil()
					n += nn
					if err != nil {
						return
					}
					z.Child = nil
				} else {
					if z.Child == nil {
						z.Child = new(TestType)
					}

					nn, err = dc.ReadIdent(z.Child)

					n += nn
					if err != nil {
						return
					}

				}

			}
		}

	}

	enc.Done(dc)
	return
}

func (z *TestFast) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	_, err := z.EncodeMsg(&buf)
	return buf.Bytes(), err
}

func (z *TestFast) Unmarshal(b []byte) error {
	_, err := z.DecodeMsg(bytes.NewReader(b))
	return err
}

func (z *TestFast) EncodeMsg(w io.Writer) (n int, err error) {
	var nn int
	en := enc.NewEncoder(w)
	_ = nn
	_ = en

	if z == nil {
		nn, err = en.WriteNil()
		n += nn
		if err != nil {
			return
		}
	} else {

		nn, err = en.WriteMapHeader(uint32(4))
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteString("Lat")
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteFloat64(z.Lat)
		n += nn
		if err != nil {
			return
		}
		nn, err = en.WriteString("Long")
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteFloat64(z.Long)
		n += nn
		if err != nil {
			return
		}
		nn, err = en.WriteString("Alt")
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteFloat64(z.Alt)
		n += nn
		if err != nil {
			return
		}
		nn, err = en.WriteString("Data")
		n += nn
		if err != nil {
			return
		}

		nn, err = en.WriteBytes(z.Data)
		n += nn
		if err != nil {
			return
		}
	}
	return
}

func (z *TestFast) DecodeMsg(r io.Reader) (n int, err error) {
	var sz uint32
	var nn int
	var field []byte
	dc := enc.NewDecoder(r)
	_ = sz
	_ = nn
	_ = field

	if dc.IsNil() {
		nn, err = dc.ReadNil()
		n += nn
		if err != nil {
			return
		}
		z = nil
	} else {
		if z == nil {
			z = new(TestFast)
		}

		sz, nn, err = dc.ReadMapHeader()
		n += nn
		if err != nil {
			return
		}
		for xplz := uint32(0); xplz < sz; xplz++ {
			field, nn, err = dc.ReadStringAsBytes(field)
			n += nn
			if err != nil {
				return
			}
			switch enc.UnsafeString(field) {

			case "Lat":

				z.Lat, nn, err = dc.ReadFloat64()

				n += nn
				if err != nil {
					return
				}

			case "Long":

				z.Long, nn, err = dc.ReadFloat64()

				n += nn
				if err != nil {
					return
				}

			case "Alt":

				z.Alt, nn, err = dc.ReadFloat64()

				n += nn
				if err != nil {
					return
				}

			case "Data":

				z.Data, nn, err = dc.ReadBytes(z.Data)

				n += nn
				if err != nil {
					return
				}

			}
		}

	}

	enc.Done(dc)
	return
}

