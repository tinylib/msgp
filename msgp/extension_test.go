package msgp

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

var extSizes = [...]int{0, 1, 2, 4, 8, 16, int(tint8), int(tuint16), int(tuint32)}

func randomExt() RawExtension {
	e := RawExtension{}
	e.Type = int8(rand.Int())
	e.Data = RandBytes(extSizes[rand.Intn(len(extSizes))])
	return e
}

func TestReadWriteExtension(t *testing.T) {
	rand.Seed(time.Now().Unix())

	var buf bytes.Buffer
	en := NewWriter(&buf)
	dc := NewReader(&buf)

	t.Run("interface", func(t *testing.T) {
		for i := 0; i < 25; i++ {
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
		for i := 0; i < 25; i++ {
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

func TestReadWriteExtensionBytes(t *testing.T) {
	var bts []byte
	rand.Seed(time.Now().Unix())

	for i := 0; i < 24; i++ {
		e := randomExt()
		bts, _ = AppendExtension(bts[0:0], &e)
		_, err := ReadExtensionBytes(bts, &e)
		if err != nil {
			t.Errorf("error with extension (length %d): %s", len(bts), err)
		}
	}
}

func TestAppendAndWriteCompatibility(t *testing.T) {
	rand.Seed(time.Now().Unix())

	var bts []byte
	var buf bytes.Buffer
	en := NewWriter(&buf)

	for i := 0; i < 24; i++ {
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
