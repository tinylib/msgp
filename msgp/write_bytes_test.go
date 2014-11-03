package msgp

import (
	"bytes"
	"testing"
)

func TestAppendMapHeader(t *testing.T) {
	szs := []uint32{0, 1, uint32(tint8), uint32(tint16), tuint32}
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	for _, sz := range szs {
		buf.Reset()
		en.WriteMapHeader(sz)
		bts = AppendMapHeader(bts[0:0], sz)

		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for size %d, encoder wrote %q and append wrote %q", sz, buf.Bytes(), bts)
		}
	}
}

func TestAppendArrayHeader(t *testing.T) {
	szs := []uint32{0, 1, uint32(tint8), uint32(tint16), tuint32}
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	for _, sz := range szs {
		buf.Reset()
		en.WriteArrayHeader(sz)
		bts = AppendArrayHeader(bts[0:0], sz)

		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for size %d, encoder wrote %q and append wrote %q", sz, buf.Bytes(), bts)
		}
	}
}

func TestAppendNil(t *testing.T) {
	var bts []byte
	bts = AppendNil(bts[0:0])
	if bts[0] != mnil {
		t.Fatal("bts[0] is not 'nil'")
	}
}

func TestAppendFloat64(t *testing.T) {
	f := float64(3.14159)
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	en.WriteFloat64(f)
	bts = AppendFloat64(bts[0:0], f)
	if !bytes.Equal(buf.Bytes(), bts) {
		t.Errorf("for float %f, encoder wrote %q; append wrote %q", f, buf.Bytes(), bts)
	}
}

func TestAppendFloat32(t *testing.T) {
	f := float32(3.14159)
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	en.WriteFloat32(f)
	bts = AppendFloat32(bts[0:0], f)
	if !bytes.Equal(buf.Bytes(), bts) {
		t.Errorf("for float %f, encoder wrote %q; append wrote %q", f, buf.Bytes(), bts)
	}
}

func TestAppendInt64(t *testing.T) {
	is := []int64{0, 1, -5, -50, int64(tint16), int64(tint32), int64(tint64)}
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	for _, i := range is {
		buf.Reset()
		en.WriteInt64(i)
		bts = AppendInt64(bts[0:0], i)
		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for int64 %d, encoder wrote %q; append wrote %q", buf.Bytes(), bts)
		}
	}
}

func TestAppendUint64(t *testing.T) {
	us := []uint64{0, 1, uint64(tuint16), uint64(tuint32), tuint64}
	var buf bytes.Buffer
	en := NewWriter(&buf)
	var bts []byte

	for _, u := range us {
		buf.Reset()
		en.WriteUint64(u)
		bts = AppendUint64(bts[0:0], u)
		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for uint64 %d, encoder wrote %q; append wrote %q", buf.Bytes(), bts)
		}
	}
}

func TestAppendBytes(t *testing.T) {
	sizes := []int{0, 1, 225, int(tuint32)}
	var buf bytes.Buffer
	en := NewWriter(&buf)
	var bts []byte

	for _, sz := range sizes {
		buf.Reset()
		b := RandBytes(sz)
		en.WriteBytes(b)
		bts = AppendBytes(b[0:0], b)
		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for bytes of length %d, encoder wrote %d bytes and append wrote %d bytes", sz, buf.Len(), len(bts))
		}
	}
}

func TestAppendString(t *testing.T) {
	sizes := []int{0, 1, 225, int(tuint32)}
	var buf bytes.Buffer
	en := NewWriter(&buf)
	var bts []byte

	for _, sz := range sizes {
		buf.Reset()
		s := string(RandBytes(sz))
		en.WriteString(s)
		bts = AppendString(bts[0:0], s)
		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for string of length %d, encoder wrote %d bytes and append wrote %d bytes", sz, buf.Len(), len(bts))
		}
	}
}

func TestAppendBool(t *testing.T) {
	vs := []bool{true, false}
	var buf bytes.Buffer
	en := NewWriter(&buf)
	var bts []byte

	for _, v := range vs {
		buf.Reset()
		en.WriteBool(v)
		bts = AppendBool(bts[0:0], v)
		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for %t, encoder wrote %q and append wrote %q", v, buf.Bytes(), bts)
		}
	}
}
