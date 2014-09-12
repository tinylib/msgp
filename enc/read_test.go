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
