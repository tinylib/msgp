package enc

import (
	"bytes"
	"math"
	"math/rand"
	"testing"
)

func BenchmarkReadWriteBytes(b *testing.B) {
	bts := RandBytes(256)
	var buf bytes.Buffer
	wr := MsgWriter{w: &buf}
	rd := MsgReader{r: &buf}

	var ns []byte
	b.SetBytes(256)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		wr.WriteBytes(bts)
		ns, _, _ = rd.ReadBytes(ns)
	}

}

func BenchmarkReadWriteFloat64(b *testing.B) {
	flt := rand.Float64()
	var buf bytes.Buffer

	wr := MsgWriter{w: &buf}
	rd := MsgReader{r: &buf}
	b.SetBytes(9)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		nw, _ := wr.WriteFloat64(flt)
		nflt, nr, err := rd.ReadFloat64()
		if err != nil {
			b.Fatal(err)
		}
		if nflt != flt {
			b.Fatalf("%f in; %f out", flt, nflt)
		}
		if nw != nr {
			b.Fatalf("Wrote %d bytes; read %d bytes", nw, nr)
		}
	}
}

func BenchmarkReadWriteInt64(b *testing.B) {
	var buf bytes.Buffer
	wr := MsgWriter{w: &buf}
	rd := MsgReader{r: &buf}
	ints := make([]int64, b.N)
	for i := range ints {
		ints[i] = rand.Int63n(math.MaxInt64) - (math.MaxInt64)
	}
	b.SetBytes(9)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		wr.WriteInt64(ints[i])
		no, _, _ := rd.ReadInt64()
		if no != ints[i] {
			b.Fatalf("%d in; %d out", ints[i], no)
		}
	}
}
