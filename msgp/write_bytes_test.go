package msgp

import (
	"bytes"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestIssue116(t *testing.T) {
	data := AppendInt64(nil, math.MinInt64)
	i, _, err := ReadInt64Bytes(data)
	if err != nil {
		t.Fatal(err)
	}
	if i != math.MinInt64 {
		t.Errorf("put %d in and got %d out", int64(math.MinInt64), i)
	}

	var buf bytes.Buffer

	w := NewWriter(&buf)
	w.WriteInt64(math.MinInt64)
	w.Flush()
	i, err = NewReader(&buf).ReadInt64()
	if err != nil {
		t.Fatal(err)
	}
	if i != math.MinInt64 {
		t.Errorf("put %d in and got %d out", int64(math.MinInt64), i)
	}
}

func TestAppendMapHeader(t *testing.T) {
	szs := []uint32{0, 1, uint32(tint8), uint32(tint16), tuint32}
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	for _, sz := range szs {
		buf.Reset()
		en.WriteMapHeader(sz)
		en.Flush()
		bts = AppendMapHeader(bts[0:0], sz)

		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for size %d, encoder wrote %q and append wrote %q", sz, buf.Bytes(), bts)
		}
	}
}

func BenchmarkAppendMapHeader(b *testing.B) {
	buf := make([]byte, 0, 9)
	N := b.N / 4
	b.ReportAllocs()
	b.ResetTimer()
	for range N {
		AppendMapHeader(buf[:0], 0)
		AppendMapHeader(buf[:0], uint32(tint8))
		AppendMapHeader(buf[:0], tuint16)
		AppendMapHeader(buf[:0], tuint32)
	}
}

func TestAppendArrayHeader(t *testing.T) {
	szs := []uint32{0, 1, uint32(tint8), uint32(tint16), tuint32}
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	for _, sz := range szs {
		buf.Reset()
		en.WriteArrayHeader(sz)
		en.Flush()
		bts = AppendArrayHeader(bts[0:0], sz)

		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for size %d, encoder wrote %q and append wrote %q", sz, buf.Bytes(), bts)
		}
	}
}

func BenchmarkAppendArrayHeader(b *testing.B) {
	buf := make([]byte, 0, 9)
	N := b.N / 4
	b.ReportAllocs()
	b.ResetTimer()
	for range N {
		AppendArrayHeader(buf[:0], 0)
		AppendArrayHeader(buf[:0], uint32(tint8))
		AppendArrayHeader(buf[:0], tuint16)
		AppendArrayHeader(buf[:0], tuint32)
	}
}

func TestAppendBytesHeader(t *testing.T) {
	szs := []uint32{0, 1, uint32(tint8), uint32(tint16), tuint32, math.MaxUint32}
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	for _, sz := range szs {
		buf.Reset()
		en.WriteBytesHeader(sz)
		en.Flush()
		bts = AppendBytesHeader(bts[0:0], sz)

		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for size %d, encoder wrote %q and append wrote %q", sz, buf.Bytes(), bts)
		}
	}
}

func BenchmarkAppendBytesHeader(b *testing.B) {
	buf := make([]byte, 0, 9)
	N := b.N / 4
	b.ReportAllocs()
	b.ResetTimer()
	for range N {
		AppendBytesHeader(buf[:0], 0)
		AppendBytesHeader(buf[:0], uint32(tint8))
		AppendBytesHeader(buf[:0], tuint16)
		AppendBytesHeader(buf[:0], tuint32)
	}
}

func TestAppendNil(t *testing.T) {
	var bts []byte
	bts = AppendNil(bts[0:0])
	if bts[0] != mnil {
		t.Fatal("bts[0] is not 'nil'")
	}
}

func TestAppendFloat(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	const n = 1e7
	src := make([]float64, n)
	for i := range src {
		// ~50% full float64, 50% converted from float32.
		if rng.Uint32()&1 == 1 {
			src[i] = rng.NormFloat64()
		} else {
			src[i] = float64(math.MaxFloat32 * (0.5 - rng.Float32()))
		}
	}

	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	for _, f := range src {
		en.WriteFloat(f)
		bts = AppendFloat(bts, f)
	}
	en.Flush()
	if buf.Len() != len(bts) {
		t.Errorf("encoder wrote %d; append wrote %d bytes", buf.Len(), len(bts))
	}
	t.Logf("%f bytes/value", float64(buf.Len())/n)
	a, b := bts, buf.Bytes()
	for i := range a {
		if a[i] != b[i] {
			t.Errorf("mismatch at byte %d, %d != %d", i, a[i], b[i])
			break
		}
	}

	for i, want := range src {
		var got float64
		var err error
		got, a, err = ReadFloat64Bytes(a)
		if err != nil {
			t.Fatal(err)
		}
		if want != got {
			t.Errorf("value #%d: want %v; got %v", i, want, got)
		}
	}
}

func BenchmarkAppendFloat(b *testing.B) {
	rng := rand.New(rand.NewSource(0))
	const n = 1 << 16
	src := make([]float64, n)
	for i := range src {
		// ~50% full float64, 50% converted from float32.
		if rng.Uint32()&1 == 1 {
			src[i] = rng.NormFloat64()
		} else {
			src[i] = float64(math.MaxFloat32 * (0.5 - rng.Float32()))
		}
	}
	buf := make([]byte, 0, 9)
	b.SetBytes(8)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AppendFloat(buf, src[i&(n-1)])
	}
}

func TestAppendFloat64(t *testing.T) {
	f := float64(3.14159)
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	en.WriteFloat64(f)
	en.Flush()
	bts = AppendFloat64(bts[0:0], f)
	if !bytes.Equal(buf.Bytes(), bts) {
		t.Errorf("for float %f, encoder wrote %q; append wrote %q", f, buf.Bytes(), bts)
	}
}

func BenchmarkAppendFloat64(b *testing.B) {
	f := float64(3.14159)
	buf := make([]byte, 0, 9)
	b.SetBytes(9)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AppendFloat64(buf[0:0], f)
	}
}

func TestAppendFloat32(t *testing.T) {
	f := float32(3.14159)
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	en.WriteFloat32(f)
	en.Flush()
	bts = AppendFloat32(bts[0:0], f)
	if !bytes.Equal(buf.Bytes(), bts) {
		t.Errorf("for float %f, encoder wrote %q; append wrote %q", f, buf.Bytes(), bts)
	}
}

func BenchmarkAppendFloat32(b *testing.B) {
	f := float32(3.14159)
	buf := make([]byte, 0, 5)
	b.SetBytes(5)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AppendFloat32(buf[0:0], f)
	}
}

func TestAppendInt64(t *testing.T) {
	is := []int64{0, 1, -5, -50, int64(tint16), int64(tint32), tint64}
	var buf bytes.Buffer
	en := NewWriter(&buf)

	var bts []byte
	for _, i := range is {
		buf.Reset()
		en.WriteInt64(i)
		en.Flush()
		bts = AppendInt64(bts[0:0], i)
		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for int64 %d, encoder wrote %q; append wrote %q", i, buf.Bytes(), bts)
		}
	}
}

func BenchmarkAppendInt64(b *testing.B) {
	is := []int64{0, 1, -5, -50, int64(tint16), int64(tint32), tint64}
	l := len(is)
	buf := make([]byte, 0, 9)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AppendInt64(buf[0:0], is[i%l])
	}
}

func TestAppendUint64(t *testing.T) {
	us := []uint64{0, 1, uint64(tuint16), uint64(tuint32), tuint64}
	var buf bytes.Buffer
	en := NewWriter(&buf)
	var bts []byte

	for _, u := range us {
		buf.Reset()
		en.WriteUint64(u)
		en.Flush()
		bts = AppendUint64(bts[0:0], u)
		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for uint64 %d, encoder wrote %q; append wrote %q", u, buf.Bytes(), bts)
		}
	}
}

func BenchmarkAppendUint64(b *testing.B) {
	us := []uint64{0, 1, 15, uint64(tuint16), uint64(tuint32), tuint64}
	buf := make([]byte, 0, 9)
	b.ReportAllocs()
	b.ResetTimer()
	l := len(us)
	for i := 0; i < b.N; i++ {
		AppendUint64(buf[0:0], us[i%l])
	}
}

func TestAppendBytes(t *testing.T) {
	sizes := []int{0, 1, 225, int(tuint32)}
	var buf bytes.Buffer
	en := NewWriter(&buf)
	var bts []byte

	for _, sz := range sizes {
		buf.Reset()
		b := RandBytes(sz)
		en.WriteBytes(b)
		en.Flush()
		bts = AppendBytes(b[0:0], b)
		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for bytes of length %d, encoder wrote %d bytes and append wrote %d bytes", sz, buf.Len(), len(bts))
		}
	}
}

func benchappendBytes(size uint32, b *testing.B) {
	bts := RandBytes(int(size))
	buf := make([]byte, 0, len(bts)+5)
	b.SetBytes(int64(len(bts) + 5))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AppendBytes(buf[0:0], bts)
	}
}

func BenchmarkAppend16Bytes(b *testing.B) { benchappendBytes(16, b) }

func BenchmarkAppend256Bytes(b *testing.B) { benchappendBytes(256, b) }

func BenchmarkAppend2048Bytes(b *testing.B) { benchappendBytes(2048, b) }

func TestAppendString(t *testing.T) {
	sizes := []int{0, 1, 225, int(tuint32)}
	var buf bytes.Buffer
	en := NewWriter(&buf)
	var bts []byte

	for _, sz := range sizes {
		buf.Reset()
		s := string(RandBytes(sz))
		en.WriteString(s)
		en.Flush()
		bts = AppendString(bts[0:0], s)
		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for string of length %d, encoder wrote %d bytes and append wrote %d bytes", sz, buf.Len(), len(bts))
			t.Errorf("WriteString prefix: %x", buf.Bytes()[0:5])
			t.Errorf("Appendstring prefix: %x", bts[0:5])
		}
	}
}

func benchappendString(size uint32, b *testing.B) {
	str := string(RandBytes(int(size)))
	buf := make([]byte, 0, len(str)+5)
	b.SetBytes(int64(len(str) + 5))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AppendString(buf[0:0], str)
	}
}

func BenchmarkAppend16String(b *testing.B) { benchappendString(16, b) }

func BenchmarkAppend256String(b *testing.B) { benchappendString(256, b) }

func BenchmarkAppend2048String(b *testing.B) { benchappendString(2048, b) }

func TestAppendBool(t *testing.T) {
	vs := []bool{true, false}
	var buf bytes.Buffer
	en := NewWriter(&buf)
	var bts []byte

	for _, v := range vs {
		buf.Reset()
		en.WriteBool(v)
		en.Flush()
		bts = AppendBool(bts[0:0], v)
		if !bytes.Equal(buf.Bytes(), bts) {
			t.Errorf("for %t, encoder wrote %q and append wrote %q", v, buf.Bytes(), bts)
		}
	}
}

func BenchmarkAppendBool(b *testing.B) {
	vs := []bool{true, false}
	buf := make([]byte, 0, 1)
	b.SetBytes(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AppendBool(buf[0:0], vs[i%2])
	}
}

func BenchmarkAppendTime(b *testing.B) {
	t := time.Now()
	b.SetBytes(15)
	buf := make([]byte, 0, 15)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AppendTime(buf[0:0], t)
	}
}

func BenchmarkAppendTimeExt(b *testing.B) {
	t := time.Now()
	buf := make([]byte, 0, 15)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AppendTimeExt(buf[0:0], t)
	}
}

// TestEncodeDecode does a back-and-forth test of encoding and decoding and compare the value with a given output.
func TestEncodeDecode(t *testing.T) {
	for _, tc := range []struct {
		name        string
		input       any
		output      any
		encodeError string
	}{
		{
			name:  "nil",
			input: nil,
		},
		{
			name:  "bool",
			input: true,
		},
		{
			name:  "int",
			input: int64(42),
		},
		{
			name:  "float",
			input: 3.14159,
		},
		{
			name:  "string",
			input: "hello",
		},
		{
			name:  "bytes",
			input: []byte("hello"),
		},
		{
			name:  "array-empty",
			input: []any{},
		},
		{
			name:  "array",
			input: []any{int64(1), int64(2), int64(3)},
		},
		{
			name:  "map-empty",
			input: map[string]any{},
		},
		{
			name:  "map",
			input: map[string]any{"a": int64(1), "b": int64(2)},
		},
		{
			name:  "map-interface",
			input: map[string]any{"a": int64(1), "b": "2"},
		},
		{
			name:   "map-string",
			input:  map[string]string{"a": "1", "b": "2"},
			output: map[string]any{"a": "1", "b": "2"},
		},
		{
			name:   "map-array",
			input:  map[string][]int64{"a": {1, 2}, "b": {3}},
			output: map[string]any{"a": []any{int64(1), int64(2)}, "b": []any{int64(3)}},
		},
		{
			name:   "map-map",
			input:  map[string]map[string]int64{"a": {"a": 1, "b": 2}, "b": {"c": 3}},
			output: map[string]any{"a": map[string]any{"a": int64(1), "b": int64(2)}, "b": map[string]any{"c": int64(3)}},
		},
		{
			name:   "array-map",
			input:  []any{map[string]any{"a": int64(1), "b": "2"}, map[string]int64{"c": 3}},
			output: []any{map[string]any{"a": int64(1), "b": "2"}, map[string]any{"c": int64(3)}},
		},
		{
			name:   "array-array",
			input:  []any{[]int64{1, 2}, []any{int64(3)}},
			output: []any{[]any{int64(1), int64(2)}, []any{int64(3)}},
		},
		{
			name:  "array-array-map",
			input: []any{[]any{int64(1), int64(2)}, map[string]any{"c": int64(3)}},
		},
		{
			name:  "map-array-map",
			input: map[string]any{"a": []any{int64(1), int64(2)}, "b": map[string]any{"c": int64(3)}},
		},
		{
			name:        "map-invalid-keys",
			input:       map[any]any{int64(1): int64(2)},
			encodeError: "msgp: map keys must be strings",
		},
		{
			name:        "map-nested-invalid-keys",
			input:       map[string]any{"a": map[int64]string{1: "2"}},
			encodeError: "msgp: map keys must be strings",
		},
		{
			name:        "invalid-type",
			input:       struct{}{},
			encodeError: "msgp: type \"struct {}\" not supported",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// If no output is given, use the input as output
			if tc.output == nil {
				tc.output = tc.input
			}

			buf, err := AppendIntf(nil, tc.input)
			if tc.encodeError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.encodeError) {
					t.Fatalf("expected encode error '%s' but got '%s'", tc.encodeError, err)
				}
				return
			}

			if tc.encodeError == "" && err != nil {
				t.Fatalf("expected no encode error but got '%s'", err.Error())
			}

			out, _, _ := ReadIntfBytes(buf)
			if err != nil {
				t.Fatalf("expected no decode error but got '%s'", err.Error())
			}

			if !reflect.DeepEqual(tc.output, out) {
				t.Fatalf("expected '%v' but got '%v'", tc.input, out)
			}
		})
	}
}
