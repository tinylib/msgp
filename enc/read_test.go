package enc

import (
	"bytes"
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
