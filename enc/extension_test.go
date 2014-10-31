package enc

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
	en := NewEncoder(&buf)
	dc := NewDecoder(&buf)

	for i := 0; i < 25; i++ {
		buf.Reset()
		e := randomExt()
		en.WriteExtension(&e)
		_, err := dc.ReadExtension(&e)
		if err != nil {
			t.Errorf("error with extension %v: %s", e, err)
		}
	}
}
