package enc

import (
	"bytes"
	"math"
	"testing"
)

var (
	tint8          = 126                  // cannot be most fix* types
	tint16         = 150                  // cannot be int8
	tint32         = math.MaxInt16 + 100  // cannot be int16
	tint64         = math.MaxInt32 + 100  // cannot be int32
	tuint16 uint32 = 300                  // cannot be uint8
	tuint32 uint32 = math.MaxUint16 + 100 // cannot be int16
	tuint64 uint64 = math.MaxUint32 + 100 // cannot be int32
)

func TestWriteMapHeader(t *testing.T) {
	tests := []struct {
		Sz       uint32
		Outbytes []byte
	}{
		{0, []byte{mfixmap}},
		{1, []byte{mfixmap | byte(1)}},
		{100, []byte{mmap16, byte(uint16(100) >> 8), byte(uint16(100))}},
		{tuint32,
			[]byte{mmap32,
				byte(tuint32 >> 24),
				byte(tuint32 >> 16),
				byte(tuint32 >> 8),
				byte(tuint32),
			},
		},
	}

	var buf bytes.Buffer
	var err error
	var n int
	for _, test := range tests {
		buf.Reset()
		n, err = WriteMapHeader(&buf, test.Sz)
		if err != nil {
			t.Error(err)
		}
		if n != len(test.Outbytes) {
			t.Errorf("expected to write %d bytes; wrote %d", len(test.Outbytes), n)
		}
		if !bytes.Equal(buf.Bytes(), test.Outbytes) {
			t.Errorf("Expected bytes %x; got %x", test.Outbytes, buf.Bytes())
		}
	}
}

func TestWriteArrayHeader(t *testing.T) {
	tests := []struct {
		Sz       uint32
		Outbytes []byte
	}{
		{0, []byte{mfixarray}},
		{1, []byte{mfixarray | byte(1)}},
		{tuint16, []byte{marray16, byte(tuint16 >> 8), byte(tuint16)}},
		{tuint32, []byte{marray32, byte(tuint32 >> 24), byte(tuint32 >> 16), byte(tuint32 >> 8), byte(tuint32)}},
	}

	var buf bytes.Buffer
	var err error
	var n int
	for _, test := range tests {
		buf.Reset()
		n, err = WriteArrayHeader(&buf, test.Sz)
		if err != nil {
			t.Error(err)
		}
		if n != len(test.Outbytes) {
			t.Errorf("expected to write %d bytes; wrote %d", len(test.Outbytes), n)
		}
		if !bytes.Equal(buf.Bytes(), test.Outbytes) {
			t.Errorf("Expected bytes %x; got %x", test.Outbytes, buf.Bytes())
		}
	}
}
