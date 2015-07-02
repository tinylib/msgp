//+build amd64,!noasm

package msgp

import (
	"bytes"
	"testing"
)

func BenchmarkPut4StrHeader(b *testing.B) {
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

func BenchmarkPut4BinHeader(b *testing.B) {
	out := make([]byte, 8)
	ptr := &out[0]
	b.SetBytes(12)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		putBinHdr(ptr, 0)
		putBinHdr(ptr, 180)
		putBinHdr(ptr, 65000)
		putBinHdr(ptr, 1<<23)
	}
}

func BenchmarkPut4ArrayHeader(b *testing.B) {
	out := make([]byte, 8)
	ptr := &out[0]
	b.SetBytes(12)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		putArrayHdr(ptr, 0)
		putArrayHdr(ptr, 180)
		putArrayHdr(ptr, 65000)
		putArrayHdr(ptr, 1<<23)
	}
}

func BenchmarkPut4MapHeader(b *testing.B) {
	out := make([]byte, 8)
	ptr := &out[0]
	b.SetBytes(12)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		putMapHdr(ptr, 0)
		putMapHdr(ptr, 180)
		putMapHdr(ptr, 65000)
		putMapHdr(ptr, 1<<23)
	}
}

func TestPutMapHeader(t *testing.T) {
	out := make([]byte, 8)
	ptr := &out[0]

	cases := []struct {
		sz     int
		output []byte
	}{
		{sz: 0, output: []byte{mfixmap}},
		{sz: 1, output: []byte{mfixmap | 1}},
		{sz: 15, output: []byte{mfixmap | 15}},
		{sz: 255, output: []byte{mmap16, 0, 255}},
		{sz: 400, output: []byte{mmap16, 0x1, 0x90}},
		{sz: 65000, output: []byte{mmap16, 0xfd, 0xe8}},
		{sz: 1 << 23, output: []byte{mmap32, 0, 0x80, 0, 0}},
		{sz: 1 << 30, output: []byte{mmap32, 0x40, 0, 0, 0}},
	}

	for i, c := range cases {
		w := putMapHdr(ptr, c.sz)
		if w < 8 && !bytes.Equal(c.output, out[:w]) {
			t.Errorf("case %d: expected %#v; got %#v", i, c.output, out)
		} else if w > 8 {
			t.Errorf("case %d: put %d in and got length %d...", i, c.sz, w)
		}
	}
}

func TestPutArrayHeader(t *testing.T) {
	out := make([]byte, 8)
	ptr := &out[0]

	cases := []struct {
		sz     int
		output []byte
	}{
		{sz: 0, output: []byte{mfixarray}},
		{sz: 1, output: []byte{mfixarray | 1}},
		{sz: 15, output: []byte{mfixarray | 15}},
		{sz: 255, output: []byte{marray16, 0, 255}},
		{sz: 400, output: []byte{marray16, 0x1, 0x90}},
		{sz: 65000, output: []byte{marray16, 0xfd, 0xe8}},
		{sz: 1 << 23, output: []byte{marray32, 0, 0x80, 0, 0}},
		{sz: 1 << 30, output: []byte{marray32, 0x40, 0, 0, 0}},
	}

	for i, c := range cases {
		w := putArrayHdr(ptr, c.sz)
		if w < 8 && !bytes.Equal(c.output, out[:w]) {
			t.Errorf("case %d: expected %#v; got %#v", i, c.output, out)
		} else if w > 8 {
			t.Errorf("case %d: put %d in and got length %d...", i, c.sz, w)
		}
	}
}

func TestPutStrHeader(t *testing.T) {
	out := make([]byte, 8)
	ptr := &out[0]

	cases := []struct {
		sz     int
		output []byte
	}{
		{sz: 0, output: []byte{mfixstr}},
		{sz: 1, output: []byte{mfixstr | 1}},
		{sz: 31, output: []byte{mfixstr | 31}},
		{sz: 180, output: []byte{mstr8, 180}},
		{sz: 400, output: []byte{mstr16, 0x1, 0x90}},
		{sz: 65000, output: []byte{mstr16, 0xfd, 0xe8}},
		{sz: 1 << 23, output: []byte{mstr32, 0, 0x80, 0, 0}},
		{sz: 1 << 30, output: []byte{mstr32, 0x40, 0, 0, 0}},
	}

	for i, c := range cases {
		w := putStrHdr(ptr, c.sz)
		if w < 8 && !bytes.Equal(c.output, out[:w]) {
			t.Errorf("case %d: expected %#v; got %#v", i, c.output, out)
		} else if w > 8 {
			t.Errorf("case %d: put %d in and got length %d...", i, c.sz, w)
		}
	}
}

func TestPutBinHeader(t *testing.T) {

	out := make([]byte, 8)
	ptr := &out[0]

	cases := []struct {
		sz     int
		output []byte
	}{
		{sz: 0, output: []byte{mbin8, 0}},
		{sz: 1, output: []byte{mbin8, 1}},
		{sz: 180, output: []byte{mbin8, 180}},
		{sz: 400, output: []byte{mbin16, 0x1, 0x90}},
		{sz: 65000, output: []byte{mbin16, 0xfd, 0xe8}},
		{sz: 1 << 23, output: []byte{mbin32, 0, 0x80, 0, 0}},
		{sz: 1 << 30, output: []byte{mbin32, 0x40, 0, 0, 0}},
	}

	for i, c := range cases {
		w := putBinHdr(ptr, c.sz)
		if w < 8 && !bytes.Equal(c.output, out[:w]) {
			t.Errorf("case %d: expected %#v; got %#v", i, c.output, out)
		} else if w > 8 {
			t.Errorf("case %d: put %d in and got length %d...", i, c.sz, w)
		}
	}
}
