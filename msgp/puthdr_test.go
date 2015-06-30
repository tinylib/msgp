//+build amd64,!noasm

package msgp

import (
	"bytes"
	"testing"
)

func BenchmarkPut4BinHeader(b *testing.B) {
	out := make([]byte, 8)
	ptr := &out[0]
	b.SetBytes(11)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		putBinHdr(ptr, 0)
		putBinHdr(ptr, 180)
		putBinHdr(ptr, 65000)
		putBinHdr(ptr, 1<<23)
	}
}

func TestPutHeader(t *testing.T) {

	out := make([]byte, 8)
	ptr := &out[0]

	cases := []struct {
		sz     uint32
		output []byte
	}{
		{
			sz:     0,
			output: []byte{mbin8, 0},
		},
		{
			sz:     1,
			output: []byte{mbin8, 1},
		},
		{
			sz:     180,
			output: []byte{mbin8, 180},
		},
		{
			sz:     400,
			output: []byte{mbin16, 0x1, 0x90},
		},
		{
			sz:     65000,
			output: []byte{mbin16, 0xfd, 0xe8},
		},
		{
			sz:     1 << 23,
			output: []byte{mbin32, 0, 0x80, 0, 0},
		},
		{
			sz:     1 << 30,
			output: []byte{mbin32, 0x40, 0, 0, 0},
		},
	}

	for i, c := range cases {
		w := putBinHdr(ptr, int(c.sz))
		if w < 8 && !bytes.Equal(c.output, out[:w]) {
			t.Errorf("case %d: expected %#v; got %#v", i, c.output, out)
		} else if w > 8 {
			t.Errorf("case %d: put %d in and got length %d...", i, c.sz, w)
		}
	}
}
