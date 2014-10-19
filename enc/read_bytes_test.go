package enc

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestReadMapHeaderBytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := []uint32{0, 1, 5, 49082}

	for i, v := range tests {
		buf.Reset()
		en.WriteMapHeader(v)

		out, left, err := ReadMapHeaderBytes(buf.Bytes())
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}

		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}

		if out != v {
			t.Errorf("%d in; %d out", v, out)
		}
	}
}

func TestReadArrayHeaderBytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := []uint32{0, 1, 5, 49082}

	for i, v := range tests {
		buf.Reset()
		en.WriteArrayHeader(v)

		out, left, err := ReadArrayHeaderBytes(buf.Bytes())
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}

		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}

		if out != v {
			t.Errorf("%d in; %d out", v, out)
		}
	}
}

func TestReadNilBytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)
	en.WriteNil()

	left, err := ReadNilBytes(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(left) != 0 {
		t.Errorf("expected 0 bytes left; found %d", len(left))
	}

}

func TestReadFloat64Bytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)
	en.WriteFloat64(3.14159)

	out, left, err := ReadFloat64Bytes(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(left) != 0 {
		t.Errorf("expected 0 bytes left; found %d", len(left))
	}
	if out != 3.14159 {
		t.Errorf("%f in; %f out", 3.14159, out)
	}
}

func TestReadFloat32Bytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)
	en.WriteFloat32(3.1)

	out, left, err := ReadFloat32Bytes(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(left) != 0 {
		t.Errorf("expected 0 bytes left; found %d", len(left))
	}
	if out != 3.1 {
		t.Errorf("%f in; %f out", 3.1, out)
	}
}

func TestReadBoolBytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := []bool{true, false}

	for i, v := range tests {
		buf.Reset()
		en.WriteBool(v)
		out, left, err := ReadBoolBytes(buf.Bytes())

		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}

		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}

		if out != v {
			t.Errorf("%t in; %t out", v, out)
		}
	}
}

func TestReadInt64Bytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := []int64{-5, -30, 0, 1, 127, 300, 40921, 34908219}

	for i, v := range tests {
		buf.Reset()
		en.WriteInt64(v)
		out, left, err := ReadInt64Bytes(buf.Bytes())

		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}

		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}

		if out != v {
			t.Errorf("%d in; %d out", v, out)
		}
	}
}

func TestReadUint64Bytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := []uint64{0, 1, 127, 300, 40921, 34908219}

	for i, v := range tests {
		buf.Reset()
		en.WriteUint64(v)
		out, left, err := ReadUint64Bytes(buf.Bytes())

		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}

		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}

		if out != v {
			t.Errorf("%d in; %d out", v, out)
		}
	}
}

func TestReadBytesBytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := [][]byte{[]byte{}, []byte("some bytes"), []byte("some more bytes")}
	var scratch []byte

	for i, v := range tests {
		buf.Reset()
		en.WriteBytes(v)
		out, left, err := ReadBytesBytes(buf.Bytes(), scratch)
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}
		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}
		if !bytes.Equal(out, v) {
			t.Errorf("%q in; %q out", v, out)
		}
	}
}

func TestReadZCBytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := [][]byte{[]byte{}, []byte("some bytes"), []byte("some more bytes")}

	for i, v := range tests {
		buf.Reset()
		en.WriteBytes(v)
		out, left, err := ReadBytesZC(buf.Bytes())
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}
		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}
		if !bytes.Equal(out, v) {
			t.Errorf("%q in; %q out", v, out)
		}
	}
}

func TestReadZCString(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := []string{"", "hello", "here's another string......"}

	for i, v := range tests {
		buf.Reset()
		en.WriteString(v)

		out, left, err := ReadStringZC(buf.Bytes())
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}
		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}
		if string(out) != v {
			t.Errorf("%q in; %q out", v, out)
		}
	}
}

func TestReadStringBytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := []string{"", "hello", "here's another string......"}

	for i, v := range tests {
		buf.Reset()
		en.WriteString(v)

		out, left, err := ReadStringBytes(buf.Bytes())
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}
		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}
		if out != v {
			t.Errorf("%q in; %q out", v, out)
		}
	}
}

func TestReadComplex128Bytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := []complex128{complex(0, 0), complex(12.8, 32.0)}

	for i, v := range tests {
		buf.Reset()
		en.WriteComplex128(v)

		out, left, err := ReadComplex128Bytes(buf.Bytes())
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}
		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}
		if out != v {
			t.Errorf("%f in; %f out", v, out)
		}
	}
}

func TestReadComplex64Bytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := []complex64{complex(0, 0), complex(12.8, 32.0)}

	for i, v := range tests {
		buf.Reset()
		en.WriteComplex64(v)

		out, left, err := ReadComplex64Bytes(buf.Bytes())
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}
		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}
		if out != v {
			t.Errorf("%f in; %f out", v, out)
		}
	}
}

func TestReadTimeBytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	now := time.Now()
	en.WriteTime(now)
	out, left, err := ReadTimeBytes(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if len(left) != 0 {
		t.Errorf("expected 0 bytes left; found %d", len(left))
	}
	if !now.Equal(out) {
		t.Errorf("%s in; %s out", now, out)
	}
}

func TestReadIntfBytes(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)

	tests := make([]interface{}, 0, 10)
	tests = append(tests, float64(3.5))
	tests = append(tests, int64(-49082))
	tests = append(tests, uint64(34908))
	tests = append(tests, string("hello!"))
	tests = append(tests, []byte("blah."))
	tests = append(tests, map[string]interface{}{
		"key_one": 3.5,
		"key_two": "hi.",
	})

	for i, v := range tests {
		buf.Reset()
		en.WriteIntf(v)

		out, left, err := ReadIntfBytes(buf.Bytes())
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}
		if len(left) != 0 {
			t.Errorf("expected 0 bytes left; found %d", len(left))
		}
		if !reflect.DeepEqual(v, out) {
			t.Errorf("%v in; %v out", v, out)
		}
	}

}
