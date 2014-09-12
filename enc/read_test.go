package enc

import (
	"bytes"
	"math"
	"math/rand"
	"testing"
)

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
	var n int
	var nr int
	var sz uint32
	var err error
	for i, test := range tests {
		buf.Reset()
		n, err = WriteMapHeader(&buf, test.Sz)
		if err != nil {
			t.Fatal(err)
		}
		sz, nr, err = ReadMapHeader(&buf)
		if err != nil {
			t.Errorf("Test case %d: got error %s", i, err)
		}
		if nr != n {
			t.Errorf("Test case %d: wrote %d bytes but read %d", i, n, nr)
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
	var n int
	var nr int
	var sz uint32
	var err error
	for i, test := range tests {
		buf.Reset()
		n, err = WriteArrayHeader(&buf, test.Sz)
		if err != nil {
			t.Fatal(err)
		}
		sz, nr, err = ReadArrayHeader(&buf)
		if err != nil {
			t.Errorf("Test case %d: got error %s", i, err)
		}
		if nr != n {
			t.Errorf("Test case %d: wrote %d bytes but read %d", i, n, nr)
		}
		if sz != test.Sz {
			t.Errorf("Test case %d: wrote size %d; got size %d", i, test.Sz, sz)
		}
	}
}

func TestReadNil(t *testing.T) {
	var buf bytes.Buffer

	n, err := WriteNil(&buf)
	if err != nil {
		t.Fatal(err)
	}
	nr, err := ReadNil(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != nr {
		t.Errorf("Wrote %d bytes; read %d", n, nr)
	}
}

func TestReadFloat64(t *testing.T) {
	var buf bytes.Buffer

	for i := 0; i < 10000; i++ {
		buf.Reset()

		flt := (rand.Float64() - 0.5) * math.MaxFloat64
		n, err := WriteFloat64(&buf, flt)
		if err != nil {
			t.Fatal(err)
		}
		out, nr, err := ReadFloat64(&buf)
		if err != nil {
			t.Errorf("Error reading %f: %s", flt, err)
			continue
		}

		if nr != n {
			t.Errorf("Wrote %d bytes but read %d", n, nr)
		}

		if out != flt {
			t.Errorf("Put in %f but got out %f", flt, out)
		}
	}
}

func TestReadFloat32(t *testing.T) {
	var buf bytes.Buffer

	for i := 0; i < 10000; i++ {
		buf.Reset()

		flt := (rand.Float32() - 0.5) * math.MaxFloat32
		n, err := WriteFloat32(&buf, flt)
		if err != nil {
			t.Fatal(err)
		}
		out, nr, err := ReadFloat32(&buf)
		if err != nil {
			t.Errorf("Error reading %f: %s", flt, err)
			continue
		}

		if nr != n {
			t.Errorf("Wrote %d bytes but read %d", n, nr)
		}

		if out != flt {
			t.Errorf("Put in %f but got out %f", flt, out)
		}
	}
}

func TestReadInt64(t *testing.T) {
	var buf bytes.Buffer

	for i := 0; i < 100000; i++ {
		buf.Reset()

		num := rand.Int63n(math.MaxInt64) - (math.MaxInt64 / 2)

		n, err := WriteInt64(&buf, num)
		if err != nil {
			t.Fatal(err)
		}
		out, nr, err := ReadInt64(&buf)
		if nr != n {
			t.Errorf("Wrote %d bytes; read %d", n, nr)
		}
		if out != num {
			t.Errorf("Put %d in and got %d out", num, out)
		}
	}
}

func TestReadUint64(t *testing.T) {
	var buf bytes.Buffer

	for i := 0; i < 10000; i++ {
		buf.Reset()

		num := uint64(rand.Int63n(math.MaxInt64))

		n, err := WriteUint64(&buf, num)
		if err != nil {
			t.Fatal(err)
		}
		out, nr, err := ReadUint64(&buf)
		if nr != n {
			t.Errorf("Wrote %d bytes; read %d", n, nr)
		}
		if out != num {
			t.Errorf("Put %d in and got %d out", num, out)
		}
	}
}

func TestReadBytes(t *testing.T) {
	var buf bytes.Buffer

	sizes := []int{0, 1, 225, int(tuint32)}
	var scratch []byte
	for i, size := range sizes {
		buf.Reset()
		bts := RandBytes(size)

		n, err := WriteBytes(&buf, bts)
		if err != nil {
			t.Fatal(err)
		}

		out, nr, err := ReadBytes(&buf, scratch)
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
			continue
		}

		if n != nr {
			t.Errorf("test case %d: wrote %d bytes; read %d bytes", i, n, nr)
		}

		if !bytes.Equal(bts, out) {
			t.Error("Test case %d: Bytes not equal.", i)
		}

	}
}

func TestReadStrings(t *testing.T) {
	var buf bytes.Buffer

	sizes := []int{0, 1, 225, int(tuint32)}
	for i, size := range sizes {
		buf.Reset()
		in := string(RandBytes(size))

		n, err := WriteString(&buf, in)
		if err != nil {
			t.Fatal(err)
		}

		out, nr, err := ReadString(&buf)
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
			continue
		}

		if n != nr {
			t.Errorf("test case %d: wrote %d bytes; read %d bytes", i, n, nr)
		}

		if out != in {
			t.Error("Test case %d: strings not equal.", i)
		}

	}
}
