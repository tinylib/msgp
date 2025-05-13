package _generated

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/tinylib/msgp/msgp"
)

func TestNewTime(t *testing.T) {
	value := NewTime{
		T:     time.Now().UTC(),
		Array: []time.Time{time.Now().UTC(), time.Now().UTC()},
		Map: map[string]time.Time{
			"a": time.Now().UTC(),
		},
	}
	encoded, err := value.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	checkExtMinusOne(t, encoded)
	var got NewTime
	_, err = got.UnmarshalMsg(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !value.Equal(got) {
		t.Errorf("UnmarshalMsg got %v want %v", value, got)
	}
	if got.T.Location() != time.Local {
		t.Errorf("DecodeMsg got %v want %v", got.T.Location(), time.Local)
	}

	var buf bytes.Buffer
	w := msgp.NewWriter(&buf)
	err = value.EncodeMsg(w)
	if err != nil {
		t.Fatal(err)
	}
	w.Flush()
	checkExtMinusOne(t, buf.Bytes())

	got = NewTime{}
	r := msgp.NewReader(&buf)
	err = got.DecodeMsg(r)
	if err != nil {
		t.Fatal(err)
	}
	if !value.Equal(got) {
		t.Errorf("DecodeMsg got %v want %v", value, got)
	}
	if got.T.Location() != time.Local {
		t.Errorf("DecodeMsg got %v want %v", got.T.Location(), time.Local)
	}
}

func checkExtMinusOne(t *testing.T, b []byte) {
	r := msgp.NewReader(bytes.NewBuffer(b))
	_, err := r.ReadMapHeader()
	if err != nil {
		t.Fatal(err)
	}
	key, err := r.ReadMapKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	for !bytes.Equal(key, []byte("T")) {
		key, err = r.ReadMapKey(nil)
		if err != nil {
			t.Fatal(err)
		}
	}
	n, _, err := r.ReadExtensionRaw()
	if err != nil {
		t.Fatal(err)
	}
	if n != -1 {
		t.Fatalf("got %v want -1", n)
	}
	t.Log("Was -1 extension")
}

func TestNewTimeRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	runs := int(1e6)
	if testing.Short() {
		runs = 1e4
	}
	for i := 0; i < runs; i++ {
		nanos := rng.Int63n(999999999 + 1)
		secs := rng.Uint64()
		// Tweak the distribution, so we get more than average number of
		// length 4 and 8 timestamps.
		if rng.Intn(5) == 0 {
			secs %= uint64(time.Now().Unix())
			if rng.Intn(2) == 0 {
				nanos = 0
			}
		}

		value := NewTime{
			T: time.Unix(int64(secs), nanos),
		}
		encoded, err := value.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		var got NewTime
		_, err = got.UnmarshalMsg(encoded)
		if err != nil {
			t.Fatal(err)
		}
		if !value.Equal(got) {
			t.Fatalf("UnmarshalMsg got %v want %v", value, got)
		}
		var buf bytes.Buffer
		w := msgp.NewWriter(&buf)
		err = value.EncodeMsg(w)
		if err != nil {
			t.Fatal(err)
		}
		w.Flush()
		got = NewTime{}
		r := msgp.NewReader(&buf)
		err = got.DecodeMsg(r)
		if err != nil {
			t.Fatal(err)
		}
		if !value.Equal(got) {
			t.Fatalf("DecodeMsg got %v want %v", value, got)
		}
	}
}
