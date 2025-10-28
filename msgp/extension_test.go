package msgp

import (
	"bytes"
	"math/rand"
	"testing"
)

var extSizes = [...]int{0, 1, 2, 4, 8, 16, int(tint8), int(tuint16), int(tuint32)}

func randomExt() RawExtension {
	e := RawExtension{}
	e.Type = int8(rand.Int())
	e.Data = RandBytes(extSizes[rand.Intn(len(extSizes))])
	return e
}

func TestReadWriteExtension(t *testing.T) {
	var buf bytes.Buffer
	en := NewWriter(&buf)
	dc := NewReader(&buf)

	t.Run("interface", func(t *testing.T) {
		for range 25 {
			buf.Reset()
			e := randomExt()
			en.WriteExtension(&e)
			en.Flush()
			err := dc.ReadExtension(&e)
			if err != nil {
				t.Errorf("error with extension (length %d): %s", len(buf.Bytes()), err)
			}
		}
	})

	t.Run("raw", func(t *testing.T) {
		for range 25 {
			buf.Reset()
			e := randomExt()
			en.WriteExtensionRaw(e.Type, e.Data)
			en.Flush()
			typ, payload, err := dc.ReadExtensionRaw()
			if err != nil {
				t.Errorf("error with extension (length %d): %s", len(buf.Bytes()), err)
			}
			if typ != e.Type || !bytes.Equal(payload, e.Data) {
				t.Errorf("extension mismatch: %d %x != %d %x", typ, payload, e.Type, e.Data)
			}
		}
	})
}

func TestReadWriteLargeExtensionRaw(t *testing.T) {
	var buf bytes.Buffer
	en := NewWriter(&buf)
	dc := NewReader(&buf)

	largeExt := RawExtension{Type: int8(rand.Int()), Data: RandBytes(int(tuint32))}

	err := en.WriteExtensionRaw(largeExt.Type, largeExt.Data)
	if err != nil {
		t.Errorf("error with large extension write: %s", err)
	}
	// write a nil as a marker
	err = en.WriteNil()
	if err != nil {
		t.Errorf("error with large extension write: %s", err)
	}
	en.Flush()

	typ, payload, err := dc.ReadExtensionRaw()
	if err != nil {
		t.Errorf("error with large extension read: %s", err)
	}
	if typ != largeExt.Type || !bytes.Equal(payload, largeExt.Data) {
		t.Errorf("large extension mismatch: %d %x != %d %x", typ, payload, largeExt.Type, largeExt.Data)
	}
	err = dc.ReadNil()
	if err != nil {
		t.Errorf("error with large extension read: %s", err)
	}
}

func TestExtensionRawStackBuffer(t *testing.T) {
	var buf bytes.Buffer
	en := NewWriter(&buf)
	dc := NewReader(&buf)

	const bufSize = 128
	e := &RawExtension{Type: int8(rand.Int()), Data: RandBytes(bufSize)}

	allocs := testing.AllocsPerRun(100, func() {
		buf.Reset()

		var staticBuf [bufSize]byte
		slc := e.Data[:rand.Intn(bufSize)]
		copy(staticBuf[:], slc)

		err := en.WriteExtensionRaw(e.Type, staticBuf[:len(slc)])
		if err != nil {
			t.Errorf("error writing extension: %s", err)
		}
		en.Flush()

		typ, payload, err := dc.ReadExtensionRaw()
		if err != nil {
			t.Errorf("error reading extension: %s", err)
		}
		if typ != e.Type || !bytes.Equal(payload, slc) {
			t.Errorf("extension mismatch: %d %x != %d %x", typ, payload, e.Type, slc)
		}
	})

	if allocs != 0 {
		t.Errorf("using stack allocated buffer with WriteExtensionRaw caused %f allocations", allocs)
	}
}

func TestReadWriteExtensionBytes(t *testing.T) {
	var bts []byte

	for range 24 {
		e := randomExt()
		bts, _ = AppendExtension(bts[0:0], &e)
		_, err := ReadExtensionBytes(bts, &e)
		if err != nil {
			t.Errorf("error with extension (length %d): %s", len(bts), err)
		}
	}
}

func TestAppendAndWriteCompatibility(t *testing.T) {
	var bts []byte
	var buf bytes.Buffer
	en := NewWriter(&buf)

	for range 24 {
		buf.Reset()
		e := randomExt()
		bts, _ = AppendExtension(bts[0:0], &e)
		en.WriteExtension(&e)
		en.Flush()

		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("the outputs are different:\n\t%x\n\t%x", buf.Bytes(), bts)
		}

		_, err := ReadExtensionBytes(bts, &e)
		if err != nil {
			t.Errorf("error with extension (length %d): %s", len(bts), err)
		}
	}
}

func BenchmarkExtensionReadWrite(b *testing.B) {
	var buf bytes.Buffer
	en := NewWriter(&buf)
	dc := NewReader(&buf)

	b.Run("interface", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			buf.Reset()

			e := randomExt()
			err := en.WriteExtension(&e)
			if err != nil {
				b.Errorf("error writing extension: %s", err)
			}
			en.Flush()

			err = dc.ReadExtension(&e)
			if err != nil {
				b.Errorf("error reading extension: %s", err)
			}
		}
	})

	b.Run("raw", func(b *testing.B) {
		// this should have zero allocations
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			buf.Reset()

			e := randomExt()
			err := en.WriteExtensionRaw(e.Type, e.Data)
			if err != nil {
				b.Errorf("error writing extension: %s", err)
			}
			en.Flush()

			typ, payload, err := dc.ReadExtensionRaw()
			if err != nil {
				b.Errorf("error reading extension: %s", err)
			}
			if typ != e.Type || !bytes.Equal(payload, e.Data) {
				b.Errorf("extension mismatch: %d %x != %d %x", typ, payload, e.Type, e.Data)
			}

			buf.Reset()
		}
	})
}
