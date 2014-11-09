package msgp

import (
	"bytes"
	"io"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestSanity(t *testing.T) {
	if !isfixint(0) {
		t.Fatal("WUT.")
	}
}

func TestReadIntf(t *testing.T) {
	// NOTE: if you include cases
	// with, say, int32s, the test
	// will fail, b/c integers are
	// always read out as int64, and
	// unsigned integers as uint64

	var testCases = []interface{}{
		float64(128.032),
		float32(9082.092),
		int64(-40),
		uint64(9082981),
		time.Now(),
		"hello!",
		[]byte("hello!"),
		map[string]interface{}{
			"thing-1": "thing-1-value",
			"thing-2": int64(800),
			"thing-3": []byte("some inner bytes..."),
			"thing-4": false,
		},
	}

	var buf bytes.Buffer
	var v interface{}
	dec := NewReader(&buf)
	enc := NewWriter(&buf)

	for i, ts := range testCases {
		buf.Reset()
		err := enc.WriteIntf(ts)
		if err != nil {
			t.Errorf("Test case %d: %s", i, err)
			continue
		}
		err = enc.Flush()
		if err != nil {
			t.Fatal(err)
		}
		v, err = dec.ReadIntf()
		if err != nil {
			t.Errorf("Test case: %d: %s", i, err)
		}
		if !reflect.DeepEqual(v, ts) {
			t.Errorf("%v in; %v out", ts, v)
		}
	}

}

func TestReadMapHeader(t *testing.T) {
	tests := []struct {
		Sz uint32
	}{
		{0},
		{1},
		{tuint16},
		{tuint32},
	}

	var buf bytes.Buffer
	var sz uint32
	var err error
	wr := NewWriter(&buf)
	rd := NewReader(&buf)
	for i, test := range tests {
		buf.Reset()
		err = wr.WriteMapHeader(test.Sz)
		if err != nil {
			t.Fatal(err)
		}
		err = wr.Flush()
		if err != nil {
			t.Fatal(err)
		}
		sz, err = rd.ReadMapHeader()
		if err != nil {
			t.Errorf("Test case %d: got error %s", i, err)
		}
		if sz != test.Sz {
			t.Errorf("Test case %d: wrote size %d; got size %d", i, test.Sz, sz)
		}
	}
}

func TestReadArrayHeader(t *testing.T) {
	tests := []struct {
		Sz uint32
	}{
		{0},
		{1},
		{tuint16},
		{tuint32},
	}

	var buf bytes.Buffer
	var sz uint32
	var err error
	wr := NewWriter(&buf)
	rd := NewReader(&buf)
	for i, test := range tests {
		buf.Reset()
		err = wr.WriteArrayHeader(test.Sz)
		if err != nil {
			t.Fatal(err)
		}
		err = wr.Flush()
		if err != nil {
			t.Fatal(err)
		}
		sz, err = rd.ReadArrayHeader()
		if err != nil {
			t.Errorf("Test case %d: got error %s", i, err)
		}
		if sz != test.Sz {
			t.Errorf("Test case %d: wrote size %d; got size %d", i, test.Sz, sz)
		}
	}
}

func TestReadNil(t *testing.T) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	rd := NewReader(&buf)

	wr.WriteNil()
	wr.Flush()
	err := rd.ReadNil()
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadFloat64(t *testing.T) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	rd := NewReader(&buf)

	for i := 0; i < 100; i++ {
		buf.Reset()

		flt := (rand.Float64() - 0.5) * math.MaxFloat64
		err := wr.WriteFloat64(flt)
		if err != nil {
			t.Fatal(err)
		}
		err = wr.Flush()
		if err != nil {
			t.Fatal(err)
		}
		out, err := rd.ReadFloat64()
		if err != nil {
			t.Errorf("Error reading %f: %s", flt, err)
			continue
		}

		if out != flt {
			t.Errorf("Put in %f but got out %f", flt, out)
		}
	}
}

func TestReadFloat32(t *testing.T) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	rd := NewReader(&buf)

	for i := 0; i < 10000; i++ {
		buf.Reset()

		flt := (rand.Float32() - 0.5) * math.MaxFloat32
		err := wr.WriteFloat32(flt)
		if err != nil {
			t.Fatal(err)
		}
		err = wr.Flush()
		if err != nil {
			t.Fatal(err)
		}
		out, err := rd.ReadFloat32()
		if err != nil {
			t.Errorf("Error reading %f: %s", flt, err)
			continue
		}

		if out != flt {
			t.Errorf("Put in %f but got out %f", flt, out)
		}
	}
}

func TestReadInt64(t *testing.T) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	rd := NewReader(&buf)

	ints := []int64{-100000, -5000, -5, 0, 8, 240, int64(tuint16), int64(tuint32), int64(tuint64)}

	for i, num := range ints {
		buf.Reset()

		err := wr.WriteInt64(num)
		if err != nil {
			t.Fatal(err)
		}
		err = wr.Flush()
		if err != nil {
			t.Fatal(err)
		}
		out, err := rd.ReadInt64()
		if err != nil {
			t.Fatal(err)
		}
		if out != num {
			t.Errorf("Test case %d: put %d in and got %d out", i, num, out)
		}
	}
}

func TestReadUint64(t *testing.T) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	rd := NewReader(&buf)

	ints := []uint64{0, 8, 240, uint64(tuint16), uint64(tuint32), uint64(tuint64)}

	for i, num := range ints {
		buf.Reset()

		err := wr.WriteUint64(num)
		if err != nil {
			t.Fatal(err)
		}
		err = wr.Flush()
		if err != nil {
			t.Fatal(err)
		}
		out, err := rd.ReadUint64()
		if out != num {
			t.Errorf("Test case %d: put %d in and got %d out", i, num, out)
		}
	}
}

func TestReadBytes(t *testing.T) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	rd := NewReader(&buf)

	sizes := []int{0, 1, 225, int(tuint32)}
	var scratch []byte
	for i, size := range sizes {
		buf.Reset()
		bts := RandBytes(size)

		err := wr.WriteBytes(bts)
		if err != nil {
			t.Fatal(err)
		}
		err = wr.Flush()
		if err != nil {
			t.Fatal(err)
		}

		out, err := rd.ReadBytes(scratch)
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
			continue
		}

		if !bytes.Equal(bts, out) {
			t.Errorf("Test case %d: Bytes not equal.", i)
		}

	}
}

func TestReadString(t *testing.T) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	rd := NewReader(&buf)

	sizes := []int{0, 1, 225, int(tuint32)}
	for i, size := range sizes {
		buf.Reset()
		in := string(RandBytes(size))

		err := wr.WriteString(in)
		if err != nil {
			t.Fatal(err)
		}
		err = wr.Flush()
		if err != nil {
			t.Fatal(err)
		}

		out, err := rd.ReadString()
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}
		if out != in {
			t.Errorf("Test case %d: strings not equal.", i)
		}

	}
}

func TestReadComplex64(t *testing.T) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	rd := NewReader(&buf)

	for i := 0; i < 100; i++ {
		buf.Reset()
		f := complex(rand.Float32()*math.MaxFloat32, rand.Float32()*math.MaxFloat32)

		wr.WriteComplex64(f)
		err := wr.Flush()
		if err != nil {
			t.Fatal(err)
		}

		out, err := rd.ReadComplex64()
		if err != nil {
			t.Error(err)
			continue
		}

		if out != f {
			t.Errorf("Wrote %f; read %f", f, out)
		}

	}
}

func TestReadComplex128(t *testing.T) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	rd := NewReader(&buf)

	for i := 0; i < 10; i++ {
		buf.Reset()
		f := complex(rand.Float64()*math.MaxFloat64, rand.Float64()*math.MaxFloat64)

		wr.WriteComplex128(f)
		err := wr.Flush()
		if err != nil {
			t.Fatal(err)
		}

		out, err := rd.ReadComplex128()
		if err != nil {
			t.Error(err)
			continue
		}
		if out != f {
			t.Errorf("Wrote %f; read %f", f, out)
		}

	}
}

func TestTime(t *testing.T) {
	var buf bytes.Buffer
	now := time.Now()
	en := NewWriter(&buf)
	dc := NewReader(&buf)

	err := en.WriteTime(now)
	if err != nil {
		t.Fatal(err)
	}
	err = en.Flush()
	if err != nil {
		t.Fatal(err)
	}

	out, err := dc.ReadTime()
	if err != nil {
		t.Fatal(err)
	}
	if !now.Equal(out) {
		t.Fatalf("%s in; %s out", now, out)
	}
}

func TestSkip(t *testing.T) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	rd := NewReader(&buf)

	wr.WriteMapHeader(4)
	wr.WriteString("key_1")
	wr.WriteBytes([]byte("value_1"))
	wr.WriteString("key_2")
	wr.WriteFloat64(2.0)
	wr.WriteString("key_3")
	wr.WriteComplex128(3.0i)
	wr.WriteString("key_4")
	wr.WriteInt64(49080432189)
	wr.Flush()

	// this should skip the whole map
	err := rd.Skip()
	if err != nil {
		t.Fatal(err)
	}

	_, err = rd.NextType()
	if err != io.EOF {
		t.Errorf("expected %q; got %q", io.EOF, err)
	}

}

func BenchmarkSkip(b *testing.B) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)

	wr.WriteMapHeader(4)
	wr.WriteString("key_1")
	wr.WriteBytes([]byte("value_1"))
	wr.WriteString("key_2")
	wr.WriteFloat64(2.0)
	wr.WriteString("key_3")
	wr.WriteComplex128(3.0i)
	wr.WriteString("key_4")
	wr.WriteInt64(49080432189)
	wr.Flush()

	bts := buf.Bytes()
	b.SetBytes(int64(len(bts)))

	brd := bytes.NewReader(bts)
	rd := NewReader(brd)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		brd.Seek(0, 0)
		err := rd.Skip()
		if err != nil {
			b.Fatal(err)
		}
	}
}
