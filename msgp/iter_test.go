package msgp

import (
	"bytes"
	"fmt"
	"maps"
	"math"
	"testing"
	"time"
)

// Example: reading an array of ints using ReadArray with a *Reader.
// It prints each element in order.
func ExampleReadArray() {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	// Write an array [10, 20, 30] using WriteArray
	_ = WriteArray(w, []int{10, 20, 30}, w.WriteInt)
	_ = w.Flush()

	r := NewReader(&buf)

	seq := ReadArray(r, r.ReadInt)
	seq(func(v int, err error) bool {
		if err != nil {
			fmt.Println("err:", err)
			return false
		}
		fmt.Println(v)
		return true
	})

	// Output:
	// 10
	// 20
	// 30
}

// Example: Writing and array with WriteArray,
// then reading back using ReadArray with a *Reader.
// It prints each element in order.
func ExampleWriteArray() {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	// Write an array [10, 20, 30] using WriteArray
	_ = WriteArray(w, []int8{10, 20, 30}, w.WriteInt8)
	_ = w.Flush()

	r := NewReader(&buf)

	seq := ReadArray(r, r.ReadInt)
	seq(func(v int, err error) bool {
		if err != nil {
			fmt.Println("err:", err)
			return false
		}
		fmt.Println(v)
		return true
	})

	// Output:
	// 10
	// 20
	// 30
}

// Example: reading a map[string]int using ReadMap with a *Reader.
// It prints key=value pairs in the same order they were written.
func ExampleReadMap() {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	// Write a map {"a":1, "b":2} using WriteMap - we use the sorted version so output is predictable.
	_ = WriteMapSorted(w, map[string]int{"a": 1, "b": 2}, w.WriteString, w.WriteInt)
	_ = w.Flush()

	r := NewReader(&buf)

	seq, done := ReadMap(r, r.ReadString, r.ReadInt)
	seq(func(k string, v int) bool {
		fmt.Printf("%s=%d\n", k, v)
		return true
	})
	if err := done(); err != nil {
		fmt.Println("err:", err)
	}

	// Output:
	// a=1
	// b=2
}

// Example: writing a map using WriteMap (non-sorted).
// Uses a single-entry map to keep output deterministic.
// Use WriteMapSorted to write a sorted map.
func ExampleWriteMap() {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	// Write a map {"only":1} using WriteMap
	_ = WriteMap(w, map[string]int{"only": 1}, w.WriteString, w.WriteInt)
	_ = w.Flush()

	r := NewReader(&buf)
	seq, done := ReadMap(r, r.ReadString, r.ReadInt)
	seq(func(k string, v int) bool {
		fmt.Printf("%s=%d\n", k, v)
		return true
	})
	if err := done(); err != nil {
		fmt.Println("err:", err)
	}

	// Output:
	// only=1
}

// Example: reading an array of strings from a byte slice using ReadArrayBytes.
// It prints each element and then checks for a final error from the returned closure.
func ExampleReadArrayBytes() {
	var b []byte
	// Append ["x","y","z"] using AppendArray
	b = AppendArray(b, []string{"x", "y", "z"}, AppendString)

	seq, finish := ReadArrayBytes(b, ReadStringBytes)
	seq(func(s string) bool {
		fmt.Println(s)
		return true
	})
	if _, err := finish(); err != nil {
		fmt.Println("err:", err)
	}

	// Output:
	// x
	// y
	// z
}

// Example: reading a map[string]float64 from a byte slice using ReadMapBytes.
// It prints key=value pairs and then checks the remaining bytes/error from the returned closure.
func ExampleReadMapBytes() {
	var b []byte
	// Append {"pi":3.14, "e":2.718} using AppendMap - we use the sorted version for the example
	b = AppendMapSorted(b, map[string]float64{"pi": 3.14, "e": 2.718}, AppendString, AppendFloat64)

	seq, finish := ReadMapBytes(b, ReadStringBytes, ReadFloat64Bytes)
	seq(func(k string, v float64) bool {
		fmt.Printf("%s=%.3f\n", k, v)
		return true
	})
	if _, err := finish(); err != nil {
		fmt.Println("err:", err)
	}

	// Output:
	// e=2.718
	// pi=3.140
}

// Example: slice roundtrip with struct elements using WriteArray/ReadArray.
// Uses testDec as the element type, with EncoderTo/DecoderFrom helpers.
func ExampleReadArray_struct() {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	in := []testDec{{A: 1, B: "x"}, {A: 2, B: "y"}}
	// Write []testDec using EncoderTo as the per-element writer.
	_ = WriteArray(w, in, EncoderTo(w, testDec{}))
	_ = w.Flush()

	r := NewReader(&buf)
	// Read []testDec using DecoderFrom as the per-element reader.
	seq := ReadArray(r, DecoderFrom(r, testDec{}))

	seq(func(v testDec, err error) bool {
		if err != nil {
			fmt.Println("err:", err)
			return false
		}
		fmt.Printf("%d %s\n", v.A, v.B)
		return true
	})

	// Output:
	// 1 x
	// 2 y
}

// Example: map roundtrip with struct values using WriteMapSorted/ReadMap.
// Uses testDec as the value type and sorts keys for deterministic output.
// Employs EncoderTo for values and DecoderFrom for values when reading.
func ExampleReadMap_struct() {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	in := map[string]testDec{
		"a": {A: 1, B: "x"},
		"b": {A: 2, B: "y"},
	}
	// Write map[string]testDec using sorted keys for stable example output.
	_ = WriteMapSorted(w, in, w.WriteString, EncoderTo(w, testDec{}))
	_ = w.Flush()

	r := NewReader(&buf)
	seq, done := ReadMap(r, r.ReadString, DecoderFrom(r, testDec{}))

	seq(func(k string, v testDec) bool {
		fmt.Printf("%s=%d,%s\n", k, v.A, v.B)
		return true
	})
	if err := done(); err != nil {
		fmt.Println("err:", err)
	}

	// Output:
	// a=1,x
	// b=2,y
}

// Example: slice roundtrip with struct elements in a byte slice using AppendArray/ReadArrayBytes.
// Uses testDec as the element type, with EncoderToBytes/DecoderFromBytes helpers.
func ExampleReadArrayBytes_struct() {
	in := []testDec{{A: 1, B: "x"}, {A: 2, B: "y"}}
	var b []byte

	// Append []testDec using EncoderToBytes as the per-element appender.
	b = AppendArray(b, in, EncoderToBytes(testDec{}))

	// Read back using DecoderFromBytes as the per-element reader.
	seq, finish := ReadArrayBytes(b, DecoderFromBytes(testDec{}))

	seq(func(v testDec) bool {
		fmt.Printf("%d %s\n", v.A, v.B)
		return true
	})
	if _, err := finish(); err != nil {
		fmt.Println("err:", err)
	}

	// Output:
	// 1 x
	// 2 y
}

// Example: map roundtrip with struct values in a byte slice using AppendMapSorted/ReadMapBytes.
// Uses testDec as the value type and sorts keys for deterministic output.
// Employs EncoderToBytes for values and DecoderFromBytes for values when reading.
func ExampleReadMapBytes_struct() {
	in := map[string]testDec{
		"a": {A: 1, B: "x"},
		"b": {A: 2, B: "y"},
	}
	var b []byte

	// Append map[string]testDec with sorted keys for stable example output.
	b = AppendMapSorted(b, in, AppendString, EncoderToBytes(testDec{}))

	// Read back using DecoderFromBytes for values.
	seq, finish := ReadMapBytes(b, ReadStringBytes, DecoderFromBytes(testDec{}))

	seq(func(k string, v testDec) bool {
		fmt.Printf("%s=%d,%s\n", k, v.A, v.B)
		return true
	})
	if _, err := finish(); err != nil {
		fmt.Println("err:", err)
	}

	// Output:
	// a=1,x
	// b=2,y
}

// Example: appending a map to a byte slice using AppendMap (non-sorted).
// Uses a single-entry map to keep output deterministic.
// Use AppendMapSorted to write a sorted map.
func ExampleAppendMap() {
	var b []byte

	// Append {"only":1} using AppendMap
	b = AppendMap(b, map[string]int{"only": 1}, AppendString, AppendInt)

	// Read back and print
	seq, finish := ReadMapBytes(b, ReadStringBytes, ReadIntBytes)
	seq(func(k string, v int) bool {
		fmt.Printf("%s=%d\n", k, v)
		return true
	})
	if _, err := finish(); err != nil {
		fmt.Println("err:", err)
	}

	// Output:
	// only=1
}

var nilMsg = AppendNil(nil)

// collectSeq2 collects values from an iter.Seq2[V, error] into a slice.
// It stops at the first non-nil error and returns it together with the collected values.
func collectSeq2[V any](seq func(func(V, error) bool)) (vals []V, err error) {
	seq(func(v V, e error) bool {
		if e != nil {
			err = e
			return false
		}
		vals = append(vals, v)
		return true
	})
	return
}

func TestReadNumberArray_Int(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	want := []int{1, -2, 3, 0, 42}
	if err := w.WriteArrayHeader(uint32(len(want))); err != nil {
		t.Fatalf("WriteArrayHeader: %v", err)
	}
	for _, v := range want {
		if err := w.WriteInt(v); err != nil {
			t.Fatalf("WriteInt: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	r := NewReader(&buf)
	got, err := collectSeq2(ReadArray(r, r.ReadInt))
	if err != nil {
		t.Fatalf("iteration error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %v want %v", i, got[i], want[i])
		}
	}
}

func TestReadNumberArray_Float64(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	want := []float64{0, 1.5, -2.25, 1e9}
	if err := w.WriteArrayHeader(uint32(len(want))); err != nil {
		t.Fatalf("WriteArrayHeader: %v", err)
	}
	for _, v := range want {
		if err := w.WriteFloat64(v); err != nil {
			t.Fatalf("WriteFloat64: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	r := NewReader(&buf)
	got, err := collectSeq2(ReadArray(r, r.ReadFloat64))
	if err != nil {
		t.Fatalf("iteration error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %v want %v", i, got[i], want[i])
		}
	}
}

func TestReadArray_String(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	want := []string{"", "a", "hello", "世界"}
	if err := w.WriteArrayHeader(uint32(len(want))); err != nil {
		t.Fatalf("WriteArrayHeader: %v", err)
	}
	for _, v := range want {
		if err := w.WriteString(v); err != nil {
			t.Fatalf("WriteString: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	r := NewReader(&buf)
	got, err := collectSeq2(ReadArray(r, r.ReadString))
	if err != nil {
		t.Fatalf("iteration error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %q want %q", i, got[i], want[i])
		}
	}
}

func TestReadArray_Bool(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	want := []bool{true, false, true}
	if err := w.WriteArrayHeader(uint32(len(want))); err != nil {
		t.Fatalf("WriteArrayHeader: %v", err)
	}
	for _, v := range want {
		if err := w.WriteBool(v); err != nil {
			t.Fatalf("WriteBool: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	r := NewReader(&buf)
	got, err := collectSeq2(ReadArray(r, r.ReadBool))
	if err != nil {
		t.Fatalf("iteration error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %v want %v", i, got[i], want[i])
		}
	}
}

// A decodable type to exercise the default branch of ReadArray.
type testDec struct {
	A int
	B string
}

func (t *testDec) MarshalMsg(i []byte) ([]byte, error) {
	i = AppendInt(i, t.A)
	i = AppendString(i, t.B)
	return i, nil
}

func (t *testDec) UnmarshalMsg(i []byte) ([]byte, error) {
	var err error
	if t.A, i, err = ReadIntBytes(i); err != nil {
		return nil, err
	}
	if t.B, i, err = ReadStringBytes(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (t *testDec) EncodeMsg(w *Writer) error {
	if err := w.WriteInt(t.A); err != nil {
		return err
	}
	return w.WriteString(t.B)
}

func (t *testDec) DecodeMsg(r *Reader) error {
	var err error
	if t.A, err = r.ReadInt(); err != nil {
		return err
	}
	t.B, err = r.ReadString()
	return err
}

func TestReadArray_Decodable(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	want := []testDec{
		{A: 1, B: "x"},
		{A: -5, B: "yz"},
	}
	if err := w.WriteArrayHeader(uint32(len(want))); err != nil {
		t.Fatalf("WriteArrayHeader: %v", err)
	}
	for i := range want {
		if err := (&want[i]).EncodeMsg(w); err != nil {
			t.Fatalf("EncodeMsg: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	r := NewReader(&buf)
	got, err := collectSeq2(ReadArray[testDec](r, func() (testDec, error) {
		var t testDec
		if err := t.DecodeMsg(r); err != nil {
			return testDec{}, err
		}
		return t, nil
	}))
	if err != nil {
		t.Fatalf("iteration error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].A != want[i].A || got[i].B != want[i].B {
			t.Fatalf("index %d: got %+v want %+v", i, got[i], want[i])
		}
	}
}

func TestReadArray_TimeAndDuration(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	now := time.Unix(1700000000, 123456789).UTC()
	durs := []time.Duration{0, time.Second, -5 * time.Millisecond}

	// time.Time
	if err := w.WriteArrayHeader(2); err != nil {
		t.Fatalf("WriteArrayHeader: %v", err)
	}
	if err := w.WriteTime(now); err != nil {
		t.Fatalf("WriteTime: %v", err)
	}
	if err := w.WriteTime(now.Add(time.Minute)); err != nil {
		t.Fatalf("WriteTime: %v", err)
	}

	// time.Duration
	if err := w.WriteArrayHeader(uint32(len(durs))); err != nil {
		t.Fatalf("WriteArrayHeader: %v", err)
	}
	for _, d := range durs {
		if err := w.WriteDuration(d); err != nil {
			t.Fatalf("WriteDuration: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	r := NewReader(&buf)
	timesGot, err := collectSeq2(ReadArray[time.Time](r, r.ReadTime))
	if err != nil {
		t.Fatalf("times iteration error: %v", err)
	}
	if len(timesGot) != 2 || !timesGot[0].Equal(now) || !timesGot[1].Equal(now.Add(time.Minute)) {
		t.Fatalf("times mismatch: got %v", timesGot)
	}

	dursGot, err := collectSeq2(ReadArray[time.Duration](r, r.ReadDuration))
	if err != nil {
		t.Fatalf("durations iteration error: %v", err)
	}
	if len(dursGot) != len(durs) {
		t.Fatalf("length mismatch: got %d want %d", len(dursGot), len(durs))
	}
	for i := range durs {
		if dursGot[i] != durs[i] {
			t.Fatalf("index %d: got %v want %v", i, dursGot[i], durs[i])
		}
	}
}

func TestReadNumberArrayBytes_Uint16(t *testing.T) {
	var msg []byte
	want := []uint16{0, 1, 255, 256, 65535}

	msg = AppendArrayHeader(msg, uint32(len(want)))
	for _, v := range want {
		msg = AppendUint16(msg, v)
	}

	seq, tail := ReadArrayBytes(msg, ReadUint16Bytes)
	var got []uint16
	for v := range seq {
		got = append(got, v)
	}
	remain, err := tail()
	if err != nil {
		t.Fatalf("tail err: %v", err)
	}
	if len(remain) != 0 {
		t.Fatalf("expected no remaining bytes, got %d", len(remain))
	}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %v want %v", i, got[i], want[i])
		}
	}
}

func TestReadNumberArrayBytes_ErrOnTruncated(t *testing.T) {
	var msg []byte
	// 2 elements, but we will truncate the second one
	msg = AppendArrayHeader(msg, 2)
	msg = AppendInt32(msg, 123)
	full := AppendInt32(msg, 456)

	// Truncate to cause an error when reading the second element.
	trunc := full[:len(full)-2]

	seq, tail := ReadArrayBytes(trunc, ReadInt32Bytes)
	var got []int32
	for v := range seq {
		got = append(got, v)
		// The second element should fail before yielding
	}
	if len(got) != 1 || got[0] != 123 {
		t.Fatalf("expected to read only first element (123), got %v", got)
	}
	remain, err := tail()
	if err == nil {
		t.Fatalf("expected an error from tail() on truncated input")
	}
	// remain can be partial bytes; ensure it's from the truncated buffer
	if len(remain) != 0 {
		// Not strictly required, but checks contract that tail returns the remaining unread slice.
		_ = remain
	}
}

func TestReadArray_ErrorOnTooFewElements(t *testing.T) {
	// Array header says 2, but only 1 element provided.
	var buf bytes.Buffer
	w := NewWriter(&buf)
	if err := w.WriteArrayHeader(2); err != nil {
		t.Fatalf("WriteArrayHeader: %v", err)
	}
	if err := w.WriteInt(7); err != nil {
		t.Fatalf("WriteInt: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	r := NewReader(&buf)
	var got []int
	var firstErr error
	ReadArray(r, r.ReadInt)(func(v int, err error) bool {
		if err != nil {
			firstErr = err
			return false
		}
		got = append(got, v)
		return true
	})
	if firstErr == nil {
		t.Fatalf("expected error due to missing second element, got nil")
	}
	if len(got) != 1 || got[0] != 7 {
		t.Fatalf("unexpected values read before error: %v", got)
	}
}

// approxEqual checks approximate equality for float32/float64
func approxEqual[T ~float32 | ~float64](a, b T) bool {
	af := float64(a)
	bf := float64(b)
	const eps = 1e-6
	return math.Abs(af-bf) <= eps*(1+math.Max(math.Abs(af), math.Abs(bf)))
}

func TestRoundtripNumberArray_AllTypes(t *testing.T) {
	type testcase[V any] struct {
		name  string
		vals  []V
		write func(w *Writer, v V) error
	}
	now := time.Now() // not used here, avoid unused import warnings if refactoring
	_ = now

	tests := []any{
		testcase[uint]{name: "uint", vals: []uint{0, 1, 255, 1 << 20, math.MaxUint32}, write: func(w *Writer, v uint) error { return w.WriteUint(v) }},
		testcase[uint8]{name: "uint8", vals: []uint8{0, 1, 127, 128, 255}, write: func(w *Writer, v uint8) error { return w.WriteUint8(v) }},
		testcase[uint16]{name: "uint16", vals: []uint16{0, 255, 256, 65535}, write: func(w *Writer, v uint16) error { return w.WriteUint16(v) }},
		testcase[uint32]{name: "uint32", vals: []uint32{0, 65535, 1 << 20, math.MaxUint32}, write: func(w *Writer, v uint32) error { return w.WriteUint32(v) }},
		testcase[uint64]{name: "uint64", vals: []uint64{0, 1, 1 << 40, math.MaxUint32 + 1, math.MaxUint64 >> 1}, write: func(w *Writer, v uint64) error { return w.WriteUint64(v) }},
		testcase[int]{name: "int", vals: []int{0, 1, -1, 1 << 20, -(1 << 20)}, write: func(w *Writer, v int) error { return w.WriteInt(v) }},
		testcase[int8]{name: "int8", vals: []int8{0, 1, -1, 127, -128}, write: func(w *Writer, v int8) error { return w.WriteInt8(v) }},
		testcase[int16]{name: "int16", vals: []int16{0, 1, -1, 32767, -32768}, write: func(w *Writer, v int16) error { return w.WriteInt16(v) }},
		testcase[int32]{name: "int32", vals: []int32{0, 1, -1, math.MaxInt32, math.MinInt32}, write: func(w *Writer, v int32) error { return w.WriteInt32(v) }},
		testcase[int64]{name: "int64", vals: []int64{0, 1, -1, math.MaxInt32 + 1, math.MinInt32 - 1}, write: func(w *Writer, v int64) error { return w.WriteInt64(v) }},
		testcase[float32]{name: "float32", vals: []float32{0, 1.5, -2.25, 3.14159, 1e20}, write: func(w *Writer, v float32) error { return w.WriteFloat32(v) }},
		testcase[float64]{name: "float64", vals: []float64{0, 1.5, -2.25, math.Pi, 1e308}, write: func(w *Writer, v float64) error { return w.WriteFloat64(v) }},
	}

	for _, anytc := range tests {
		switch tc := anytc.(type) {
		case testcase[uint]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadUint))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if got[i] != tc.vals[i] {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
			r = NewReader(bytes.NewReader(nilMsg))
			got, err = collectSeq2(ReadArray(r, r.ReadUint))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != 0 {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), 0)
			}

		case testcase[uint8]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadUint8))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if got[i] != tc.vals[i] {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case testcase[uint16]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadUint16))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if got[i] != tc.vals[i] {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case testcase[uint32]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadUint32))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if got[i] != tc.vals[i] {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case testcase[uint64]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadUint64))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if got[i] != tc.vals[i] {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case testcase[int]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadInt))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if got[i] != tc.vals[i] {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case testcase[int8]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadInt8))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if got[i] != tc.vals[i] {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case testcase[int16]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadInt16))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if got[i] != tc.vals[i] {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case testcase[int32]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadInt32))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if got[i] != tc.vals[i] {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case testcase[int64]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadInt64))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if got[i] != tc.vals[i] {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case testcase[float32]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadFloat32))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if !approxEqual(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case testcase[float64]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("WriteArrayHeader %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadFloat64))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if !approxEqual(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		}
	}
}

func TestRoundtripArray_AllTypes(t *testing.T) {
	type regCase[V any] struct {
		name  string
		vals  []V
		write func(*Writer, V) error
		eq    func(a, b V) bool
	}

	now := time.Unix(1700000000, 123456789).UTC()
	later := now.Add(5 * time.Minute)

	rcases := []any{
		regCase[bool]{name: "bool", vals: []bool{true, false, true}, write: func(w *Writer, v bool) error { return w.WriteBool(v) }, eq: func(a, b bool) bool { return a == b }},
		regCase[string]{name: "string", vals: []string{"", "a", "hello", "世界"}, write: func(w *Writer, v string) error { return w.WriteString(v) }, eq: func(a, b string) bool { return a == b }},
		regCase[[]byte]{name: "bytes", vals: [][]byte{nil, {}, {0x00}, {0x01, 0x02, 0x03}}, write: func(w *Writer, v []byte) error { return w.WriteBytes(v) }, eq: func(a, b []byte) bool {
			if len(a) != len(b) {
				return false
			}
			for i := range a {
				if a[i] != b[i] {
					return false
				}
			}
			return true
		}},
		regCase[time.Time]{name: "time", vals: []time.Time{now, later}, write: func(w *Writer, v time.Time) error { return w.WriteTime(v) }, eq: func(a, b time.Time) bool { return a.Equal(b) }},
		regCase[time.Duration]{name: "duration", vals: []time.Duration{0, time.Second, -5 * time.Millisecond}, write: func(w *Writer, v time.Duration) error { return w.WriteDuration(v) }, eq: func(a, b time.Duration) bool { return a == b }},
		regCase[complex64]{name: "complex64", vals: []complex64{0, 1 + 2i, -3.5 + 4.25i}, write: func(w *Writer, v complex64) error { return w.WriteComplex64(v) }, eq: func(a, b complex64) bool { return a == b }},
		regCase[complex128]{name: "complex128", vals: []complex128{0, 1 + 2i, -3.5 + 4.25i}, write: func(w *Writer, v complex128) error { return w.WriteComplex128(v) }, eq: func(a, b complex128) bool { return a == b }},
		regCase[testDec]{name: "decoder", vals: []testDec{{A: 1, B: "abc"}, {A: 2, B: "def"}}, write: func(w *Writer, v testDec) error { return v.EncodeMsg(w) }, eq: func(a, b testDec) bool { return a == b }},
	}

	for _, anytc := range rcases {
		switch tc := anytc.(type) {
		case regCase[bool]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadBool))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
			r.Reset(bytes.NewReader(nilMsg))
			got, err = collectSeq2(ReadArray(r, r.ReadBool))
			if len(got) != 0 {
				t.Fatalf("%s len: got %d want 0", tc.name, len(got))
			}
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}

		case regCase[string]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadString))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %q want %q", tc.name, i, got[i], tc.vals[i])
				}
			}
		case regCase[[]byte]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, func() ([]byte, error) {
				return r.ReadBytes(nil)
			}))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case regCase[time.Time]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadTime))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case regCase[time.Duration]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadDuration))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case regCase[complex64]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadComplex64))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case regCase[complex128]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray(r, r.ReadComplex128))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case regCase[testDec]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteArrayHeader(uint32(len(tc.vals))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for _, v := range tc.vals {
				if err := tc.write(w, v); err != nil {
					t.Fatalf("%s write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			got, err := collectSeq2(ReadArray[testDec](r, DecoderFrom(r, testDec{})))
			if err != nil {
				t.Fatalf("%s iterate: %v", tc.name, err)
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: got %d want %d", tc.name, len(got), len(tc.vals))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}

		}
	}
}

func TestRoundtripNumberArrayBytes_AllTypes(t *testing.T) {
	type tb[V any] struct {
		name   string
		vals   []V
		append func([]byte, V) []byte
		eq     func(a, b V) bool
	}

	tests := []any{
		tb[uint]{name: "uint", vals: []uint{0, 1, 255, 1 << 20}, append: func(b []byte, v uint) []byte { return AppendUint(b, v) }, eq: func(a, b uint) bool { return a == b }},
		tb[uint8]{name: "uint8", vals: []uint8{0, 1, 127, 128, 255}, append: func(b []byte, v uint8) []byte { return AppendUint8(b, v) }, eq: func(a, b uint8) bool { return a == b }},
		tb[uint16]{name: "uint16", vals: []uint16{0, 255, 256, 65535}, append: func(b []byte, v uint16) []byte { return AppendUint16(b, v) }, eq: func(a, b uint16) bool { return a == b }},
		tb[uint32]{name: "uint32", vals: []uint32{0, 65535, 1 << 20}, append: func(b []byte, v uint32) []byte { return AppendUint32(b, v) }, eq: func(a, b uint32) bool { return a == b }},
		tb[uint64]{name: "uint64", vals: []uint64{0, 1, 1 << 40}, append: func(b []byte, v uint64) []byte { return AppendUint64(b, v) }, eq: func(a, b uint64) bool { return a == b }},
		tb[int]{name: "int", vals: []int{0, 1, -1, 1 << 20, -(1 << 20)}, append: func(b []byte, v int) []byte { return AppendInt(b, v) }, eq: func(a, b int) bool { return a == b }},
		tb[int8]{name: "int8", vals: []int8{0, 1, -1, 127, -128}, append: func(b []byte, v int8) []byte { return AppendInt8(b, v) }, eq: func(a, b int8) bool { return a == b }},
		tb[int16]{name: "int16", vals: []int16{0, 1, -1, 32767, -32768}, append: func(b []byte, v int16) []byte { return AppendInt16(b, v) }, eq: func(a, b int16) bool { return a == b }},
		tb[int32]{name: "int32", vals: []int32{0, 1, -1, math.MaxInt32, math.MinInt32}, append: func(b []byte, v int32) []byte { return AppendInt32(b, v) }, eq: func(a, b int32) bool { return a == b }},
		tb[int64]{name: "int64", vals: []int64{0, 1, -1, math.MaxInt32 + 1, math.MinInt32 - 1}, append: func(b []byte, v int64) []byte { return AppendInt64(b, v) }, eq: func(a, b int64) bool { return a == b }},
		tb[float32]{name: "float32", vals: []float32{0, 1.5, -2.25, 3.14159, 1e20}, append: func(b []byte, v float32) []byte { return AppendFloat32(b, v) }, eq: func(a, b float32) bool { return approxEqual(a, b) }},
		tb[float64]{name: "float64", vals: []float64{0, 1.5, -2.25, math.Pi, 1e308}, append: func(b []byte, v float64) []byte { return AppendFloat(b, v) }, eq: func(a, b float64) bool { return approxEqual(a, b) }},
	}

	for _, anytc := range tests {
		switch tc := anytc.(type) {
		case tb[uint]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadUintBytes)
			var got []uint
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
			seq, tail = ReadArrayBytes(nilMsg, ReadUintBytes)
			for range seq {
				t.Fatalf("%s: got entries on nil", tc.name)
			}
			remain, err = tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}

		case tb[uint8]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadUint8Bytes)
			var got []uint8
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case tb[uint16]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadUint16Bytes)
			var got []uint16
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case tb[uint32]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadUint32Bytes)
			var got []uint32
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case tb[uint64]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadUint64Bytes)
			var got []uint64
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case tb[int]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadIntBytes)
			var got []int
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case tb[int8]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadInt8Bytes)
			var got []int8
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case tb[int16]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadInt16Bytes)
			var got []int16
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case tb[int32]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadInt32Bytes)
			var got []int32
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case tb[int64]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadInt64Bytes)
			var got []int64
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case tb[float32]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadFloat32Bytes)
			var got []float32
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case tb[float64]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadFloat64Bytes)
			var got []float64
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len: %d", tc.name, len(got))
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		}
	}
}

func TestRoundtripArrayBytes_AllTypes(t *testing.T) {
	type rb[V any] struct {
		name   string
		vals   []V
		append func([]byte, V) []byte
		eq     func(a, b V) bool
	}

	now := time.Unix(1700000000, 123456789).UTC()
	later := now.Add(7 * time.Second)

	rtests := []any{
		rb[bool]{name: "bool", vals: []bool{true, false, true}, append: func(b []byte, v bool) []byte { return AppendBool(b, v) }, eq: func(a, b bool) bool { return a == b }},
		rb[string]{name: "string", vals: []string{"", "hi", "世界"}, append: func(b []byte, v string) []byte { return AppendString(b, v) }, eq: func(a, b string) bool { return a == b }},
		rb[[]byte]{name: "bytes", vals: [][]byte{nil, {}, {0x00}, {0x01, 0x02}}, append: func(b []byte, v []byte) []byte { return AppendBytes(b, v) }, eq: func(a, b []byte) bool {
			if len(a) != len(b) {
				return false
			}
			for i := range a {
				if a[i] != b[i] {
					return false
				}
			}
			return true
		}},
		rb[time.Time]{name: "time", vals: []time.Time{now, later}, append: func(b []byte, v time.Time) []byte { return AppendTime(b, v) }, eq: func(a, b time.Time) bool { return a.Equal(b) }},
		rb[time.Duration]{name: "duration", vals: []time.Duration{0, time.Second, -3 * time.Millisecond}, append: func(b []byte, v time.Duration) []byte { return AppendDuration(b, v) }, eq: func(a, b time.Duration) bool { return a == b }},
		rb[complex64]{name: "complex64", vals: []complex64{0, 1 + 2i, -3.5 + 4.25i}, append: func(b []byte, v complex64) []byte { return AppendComplex64(b, v) }, eq: func(a, b complex64) bool { return a == b }},
		rb[complex128]{name: "complex128", vals: []complex128{0, 1 + 2i, -3.5 + 4.25i}, append: func(b []byte, v complex128) []byte { return AppendComplex128(b, v) }, eq: func(a, b complex128) bool { return a == b }},
		rb[testDec]{name: "unmarshal", vals: []testDec{{A: 1, B: "abc"}, {A: 2, B: "def"}}, append: func(b []byte, v testDec) []byte { b, _ = v.MarshalMsg(b); return b }, eq: func(a, b testDec) bool { return a == b }},
	}

	for _, anytc := range rtests {
		switch tc := anytc.(type) {
		case rb[bool]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadBoolBytes)
			var got []bool
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len mismatch", tc.name)
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
			seq, tail = ReadArrayBytes(nilMsg, ReadBoolBytes)
			for range seq {
				t.Fatalf("%s: got entries on nil", tc.name)
			}
			remain, err = tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
		case rb[string]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadStringBytes)
			var got []string
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len mismatch", tc.name)
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %q want %q", tc.name, i, got[i], tc.vals[i])
				}
			}
		case rb[[]byte]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes[[]byte](msg, func(i []byte) ([]byte, []byte, error) {
				return ReadBytesBytes(i, nil)
			})
			var got [][]byte
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len mismatch", tc.name)
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case rb[time.Time]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes[time.Time](msg, ReadTimeBytes)
			var got []time.Time
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len mismatch", tc.name)
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case rb[time.Duration]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes[time.Duration](msg, ReadDurationBytes)
			var got []time.Duration
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len mismatch", tc.name)
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case rb[complex64]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadComplex64Bytes)
			var got []complex64
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len mismatch", tc.name)
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case rb[complex128]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes(msg, ReadComplex128Bytes)
			var got []complex128
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len mismatch", tc.name)
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}
		case rb[testDec]:
			msg := AppendArrayHeader(nil, uint32(len(tc.vals)))
			for _, v := range tc.vals {
				msg = tc.append(msg, v)
			}
			seq, tail := ReadArrayBytes[testDec](msg, DecoderFromBytes(testDec{}))
			var got []testDec
			for v := range seq {
				got = append(got, v)
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain: %d", tc.name, len(remain))
			}
			if len(got) != len(tc.vals) {
				t.Fatalf("%s len mismatch", tc.name)
			}
			for i := range got {
				if !tc.eq(got[i], tc.vals[i]) {
					t.Fatalf("%s[%d]: got %v want %v", tc.name, i, got[i], tc.vals[i])
				}
			}

		}
	}
}

func eqNum[T comparable](a, b T) bool { return a == b }

func TestReadMap_AllNumberTypes_SameKeyValueTypes(t *testing.T) {
	type numCase[T any] struct {
		name  string
		keys  []T
		vals  []T
		write func(*Writer, T) error
		eq    func(a, b T) bool
	}
	// Equality helpers
	eqF32 := func(a, b float32) bool { return approxEqual(a, b) }
	eqF64 := func(a, b float64) bool { return approxEqual(a, b) }

	cases := []any{
		numCase[uint]{name: "uint", keys: []uint{1, 2, 3}, vals: []uint{10, 20, 30}, write: func(w *Writer, v uint) error { return w.WriteUint(v) }, eq: eqNum[uint]},
		numCase[uint8]{name: "uint8", keys: []uint8{1, 2, 255}, vals: []uint8{9, 8, 7}, write: func(w *Writer, v uint8) error { return w.WriteUint8(v) }, eq: eqNum[uint8]},
		numCase[uint16]{name: "uint16", keys: []uint16{1, 300, 65535}, vals: []uint16{100, 200, 300}, write: func(w *Writer, v uint16) error { return w.WriteUint16(v) }, eq: eqNum[uint16]},
		numCase[uint32]{name: "uint32", keys: []uint32{1, 1 << 20, 4000000000}, vals: []uint32{42, 43, 44}, write: func(w *Writer, v uint32) error { return w.WriteUint32(v) }, eq: eqNum[uint32]},
		numCase[uint64]{name: "uint64", keys: []uint64{1, 1 << 40, 1<<50 + 123}, vals: []uint64{5, 6, 7}, write: func(w *Writer, v uint64) error { return w.WriteUint64(v) }, eq: eqNum[uint64]},
		numCase[int]{name: "int", keys: []int{-1, 0, 2}, vals: []int{100, -100, 0}, write: func(w *Writer, v int) error { return w.WriteInt(v) }, eq: eqNum[int]},
		numCase[int8]{name: "int8", keys: []int8{-128, 0, 127}, vals: []int8{1, -1, 0}, write: func(w *Writer, v int8) error { return w.WriteInt8(v) }, eq: eqNum[int8]},
		numCase[int16]{name: "int16", keys: []int16{-32768, 0, 32767}, vals: []int16{2, -2, 3}, write: func(w *Writer, v int16) error { return w.WriteInt16(v) }, eq: eqNum[int16]},
		numCase[int32]{name: "int32", keys: []int32{-1, 0, 1}, vals: []int32{7, 8, 9}, write: func(w *Writer, v int32) error { return w.WriteInt32(v) }, eq: eqNum[int32]},
		numCase[int64]{name: "int64", keys: []int64{-1, 0, 1<<40 + 5}, vals: []int64{9, 8, 7}, write: func(w *Writer, v int64) error { return w.WriteInt64(v) }, eq: eqNum[int64]},
		numCase[float32]{name: "float32", keys: []float32{-2.5, 0, 3.25}, vals: []float32{1.5, -0.25, 10}, write: func(w *Writer, v float32) error { return w.WriteFloat32(v) }, eq: eqF32},
		numCase[float64]{name: "float64", keys: []float64{-2.5, 0, 3.25}, vals: []float64{1.5, -0.25, 10}, write: func(w *Writer, v float64) error { return w.WriteFloat64(v) }, eq: eqF64},
	}

	for _, anytc := range cases {
		switch tc := anytc.(type) {
		case numCase[uint]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadUint, r.ReadUint)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
			// Test nil
			r = NewReader(bytes.NewReader(nilMsg))
			seq, tail = ReadMap(r, r.ReadUint, r.ReadUint)
			for k, v := range seq {
				t.Fatalf("nil %s: got key %v val %v", tc.name, k, v)
			}
			if err := tail(); err != nil {
				t.Fatalf("nil %s: tail: %v", tc.name, err)
			}
		case numCase[uint8]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadUint8, r.ReadUint8)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[uint16]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadUint16, r.ReadUint16)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[uint32]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadUint32, r.ReadUint32)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[uint64]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadUint64, r.ReadUint64)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[int]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadInt, r.ReadInt)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[int8]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadInt8, r.ReadInt8)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[int16]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadInt16, r.ReadInt16)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[int32]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadInt32, r.ReadInt32)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[int64]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadInt64, r.ReadInt64)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[float32]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadFloat32, r.ReadFloat32)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[float64]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadFloat64, r.ReadFloat64)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size: got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v: got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		}
	}
}

func TestReadMap_AllRegularTypes_SameKeyValueTypes(t *testing.T) {
	type regCase[T any] struct {
		name  string
		keys  []T
		vals  []T
		write func(*Writer, T) error
		eq    func(a, b T) bool
	}
	eqBool := func(a, b bool) bool { return a == b }
	eqStr := func(a, b string) bool { return a == b }
	eqBytes := func(a, b []byte) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
	eqTime := func(a, b time.Time) bool { return a.Equal(b) }
	eqDur := func(a, b time.Duration) bool { return a == b }
	eqC64 := func(a, b complex64) bool { return a == b }
	eqC128 := func(a, b complex128) bool { return a == b }

	now := time.Unix(1700000000, 123456789).UTC()
	later := now.Add(123 * time.Second)

	cases := []any{
		regCase[bool]{name: "bool", keys: []bool{false, true}, vals: []bool{false, true, false}, write: func(w *Writer, v bool) error { return w.WriteBool(v) }, eq: eqBool},
		regCase[string]{name: "string", keys: []string{"a", "b", "c"}, vals: []string{"x", "y", "z"}, write: func(w *Writer, v string) error { return w.WriteString(v) }, eq: eqStr},
		// Note: []byte keys are not comparable in Go; we validate by pair matching instead of using a map.
		regCase[[]byte]{name: "bytes", keys: [][]byte{{0x01}, {0x02, 0x03}, nil}, vals: [][]byte{{0x09}, {}, {0xFF}}, write: func(w *Writer, v []byte) error { return w.WriteBytes(v) }, eq: eqBytes},
		regCase[time.Time]{name: "time", keys: []time.Time{now, later}, vals: []time.Time{later, now}, write: func(w *Writer, v time.Time) error { return w.WriteTime(v) }, eq: eqTime},
		regCase[time.Duration]{name: "duration", keys: []time.Duration{0, time.Second, -5 * time.Millisecond}, vals: []time.Duration{time.Minute, -time.Second, 0}, write: func(w *Writer, v time.Duration) error { return w.WriteDuration(v) }, eq: eqDur},
		regCase[complex64]{name: "complex64", keys: []complex64{1 + 2i, -3 + 4.5i, 0}, vals: []complex64{9 - 2i, 3 - 4.5i, 7}, write: func(w *Writer, v complex64) error { return w.WriteComplex64(v) }, eq: eqC64},
		regCase[complex128]{name: "complex128", keys: []complex128{1 + 2i, -3 + 4.5i, 0}, vals: []complex128{9 - 2i, 3 - 4.5i, 7}, write: func(w *Writer, v complex128) error { return w.WriteComplex128(v) }, eq: eqC128},
	}

	for _, anytc := range cases {
		switch tc := anytc.(type) {
		case regCase[bool]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadBool, r.ReadBool)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s: key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
			r = NewReader(bytes.NewReader(nilMsg))
			seq, tail = ReadMap(r, r.ReadBool, r.ReadBool)
			for k, v := range seq {
				t.Fatalf("%s:expected ni results, got %v:%v", tc.name, k, v)
			}
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}

		case regCase[string]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadString, r.ReadString)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s: key %q got %q want %q", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case regCase[[]byte]:
			// For []byte keys (not comparable), validate by pair presence.
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}

			type pair struct{ k, v []byte }
			expected := make([]pair, len(tc.keys))
			for i := range tc.keys {
				expected[i] = pair{append([]byte(nil), tc.keys[i]...), append([]byte(nil), tc.vals[i]...)}
			}

			r := NewReader(&buf)
			seq, tail := ReadMap(r, func() ([]byte, error) {
				return r.ReadBytes(nil)
			}, func() ([]byte, error) {
				return r.ReadBytes(nil)
			})
			var got []pair
			for k, v := range seq {
				kk := append([]byte(nil), k...)
				vv := append([]byte(nil), v...)
				got = append(got, pair{kk, vv})
			}
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(expected) {
				t.Fatalf("%s size mismatch got %d want %d", tc.name, len(got), len(expected))
			}

			// Match expected pairs
			match := func(p pair, set []pair) bool {
				for _, q := range set {
					if eqBytes(p.k, q.k) && eqBytes(p.v, q.v) {
						return true
					}
				}
				return false
			}
			for _, p := range expected {
				if !match(p, got) {
					t.Fatalf("%s missing pair key=%v val=%v", tc.name, p.k, p.v)
				}
			}
		case regCase[time.Time]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadTime, r.ReadTime)
			got := make(map[time.Time]time.Time, len(tc.keys))
			for k, v := range seq {
				got[k.UTC()] = v.UTC()
			}
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s: key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case regCase[time.Duration]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadDuration, r.ReadDuration)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Logf("got: %#v", got)
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s: key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case regCase[complex64]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadComplex64, r.ReadComplex64)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s: key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case regCase[complex128]:
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.WriteMapHeader(uint32(len(tc.keys))); err != nil {
				t.Fatalf("hdr %s: %v", tc.name, err)
			}
			for i := range tc.keys {
				if err := tc.write(w, tc.keys[i]); err != nil {
					t.Fatalf("%s key write: %v", tc.name, err)
				}
				if err := tc.write(w, tc.vals[i]); err != nil {
					t.Fatalf("%s val write: %v", tc.name, err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			r := NewReader(&buf)
			seq, tail := ReadMap(r, r.ReadComplex128, r.ReadComplex128)
			got := maps.Collect(seq)
			if err := tail(); err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s: key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		}
	}
}

func TestReadMapBytes_AllNumberTypes_SameKeyValueTypes(t *testing.T) {
	type numCase[T any] struct {
		name   string
		keys   []T
		vals   []T
		append func([]byte, T) []byte
		eq     func(a, b T) bool
	}

	eqF32 := func(a, b float32) bool { return approxEqual(a, b) }
	eqF64 := func(a, b float64) bool { return approxEqual(a, b) }

	cases := []any{
		numCase[uint]{name: "uint", keys: []uint{1, 2, 3}, vals: []uint{10, 20, 30}, append: AppendUint, eq: eqNum[uint]},
		numCase[uint8]{name: "uint8", keys: []uint8{1, 2, 255}, vals: []uint8{9, 8, 7}, append: AppendUint8, eq: eqNum[uint8]},
		numCase[uint16]{name: "uint16", keys: []uint16{1, 300, 65535}, vals: []uint16{100, 200, 300}, append: AppendUint16, eq: eqNum[uint16]},
		numCase[uint32]{name: "uint32", keys: []uint32{1, 1 << 20, 4000000000}, vals: []uint32{42, 43, 44}, append: AppendUint32, eq: eqNum[uint32]},
		numCase[uint64]{name: "uint64", keys: []uint64{1, 1 << 40, 1<<50 + 123}, vals: []uint64{5, 6, 7}, append: AppendUint64, eq: eqNum[uint64]},
		numCase[int]{name: "int", keys: []int{-1, 0, 2}, vals: []int{100, -100, 0}, append: AppendInt, eq: eqNum[int]},
		numCase[int8]{name: "int8", keys: []int8{-128, 0, 127}, vals: []int8{1, -1, 0}, append: AppendInt8, eq: eqNum[int8]},
		numCase[int16]{name: "int16", keys: []int16{-32768, 0, 32767}, vals: []int16{2, -2, 3}, append: AppendInt16, eq: eqNum[int16]},
		numCase[int32]{name: "int32", keys: []int32{-1, 0, 1}, vals: []int32{7, 8, 9}, append: AppendInt32, eq: eqNum[int32]},
		numCase[int64]{name: "int64", keys: []int64{-1, 0, 1<<40 + 5}, vals: []int64{9, 8, 7}, append: AppendInt64, eq: eqNum[int64]},
		numCase[float32]{name: "float32", keys: []float32{-2.5, 0, 3.25}, vals: []float32{1.5, -0.25, 10}, append: AppendFloat32, eq: eqF32},
		numCase[float64]{name: "float64", keys: []float64{-2.5, 0, 3.25}, vals: []float64{1.5, -0.25, 10}, append: AppendFloat, eq: eqF64},
	}

	for _, anytc := range cases {
		switch tc := anytc.(type) {
		case numCase[uint]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadUintBytes, ReadUintBytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
			seq, tail = ReadMapBytes(nilMsg, ReadUintBytes, ReadUintBytes)
			for k, v := range seq {
				t.Fatalf("%s: got %v:%v want nothing", tc.name, k, v)
			}
			remain, err = tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
		case numCase[uint8]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadUint8Bytes, ReadUint8Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[uint16]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadUint16Bytes, ReadUint16Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[uint32]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadUint32Bytes, ReadUint32Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[uint64]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadUint64Bytes, ReadUint64Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[int]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadIntBytes, ReadIntBytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[int8]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadInt8Bytes, ReadInt8Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[int16]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadInt16Bytes, ReadInt16Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[int32]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadInt32Bytes, ReadInt32Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[int64]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadInt64Bytes, ReadInt64Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[float32]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadFloat32Bytes, ReadFloat32Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case numCase[float64]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadFloat64Bytes, ReadFloat64Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size got %d want %d", tc.name, len(got), len(tc.keys))
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		}
	}
}

func TestReadMapBytes_AllRegularTypes_SameKeyValueTypes(t *testing.T) {
	type regCase[T any] struct {
		name   string
		keys   []T
		vals   []T
		append func([]byte, T) []byte
		eq     func(a, b T) bool
	}

	eqBool := func(a, b bool) bool { return a == b }
	eqStr := func(a, b string) bool { return a == b }
	eqBytes := func(a, b []byte) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
	eqTime := func(a, b time.Time) bool { return a.Equal(b) }
	eqDur := func(a, b time.Duration) bool { return a == b }
	eqC64 := func(a, b complex64) bool { return a == b }
	eqC128 := func(a, b complex128) bool { return a == b }

	now := time.Unix(1700000000, 123456789).UTC()
	later := now.Add(99 * time.Second)

	cases := []any{
		regCase[bool]{name: "bool", keys: []bool{false, true}, vals: []bool{false, true}, append: AppendBool, eq: eqBool},
		regCase[string]{name: "string", keys: []string{"a", "b", "c"}, vals: []string{"x", "y", "z"}, append: AppendString, eq: eqStr},
		regCase[[]byte]{name: "bytes", keys: [][]byte{{0x01}, {0x02, 0x03}, nil}, vals: [][]byte{{0x09}, {}, {0xFF}}, append: AppendBytes, eq: eqBytes},
		regCase[time.Time]{name: "time", keys: []time.Time{now, later}, vals: []time.Time{later, now}, append: AppendTime, eq: eqTime},
		regCase[time.Duration]{name: "duration", keys: []time.Duration{0, time.Second, -5 * time.Millisecond}, vals: []time.Duration{time.Minute, -time.Second, 0}, append: AppendDuration, eq: eqDur},
		regCase[complex64]{name: "complex64", keys: []complex64{1 + 2i, -3 + 4.5i, 0}, vals: []complex64{9 - 2i, 3 - 4.5i, 7}, append: AppendComplex64, eq: eqC64},
		regCase[complex128]{name: "complex128", keys: []complex128{1 + 2i, -3 + 4.5i, 0}, vals: []complex128{9 - 2i, 3 - 4.5i, 7}, append: AppendComplex128, eq: eqC128},
	}

	for _, anytc := range cases {
		switch tc := anytc.(type) {
		case regCase[bool]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadBoolBytes, ReadBoolBytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
			seq, tail = ReadMapBytes(nilMsg, ReadBoolBytes, ReadBoolBytes)
			for k, v := range seq {
				t.Fatalf("%s key %v:%v want nothing", tc.name, k, v)
			}
			remain, err = tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}

		case regCase[string]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadStringBytes, ReadStringBytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %q got %q want %q", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case regCase[[]byte]:
			// For []byte keys (not comparable), validate by pair presence.
			type pair struct{ k, v []byte }
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, append([]byte(nil), tc.keys[i]...))
				msg = tc.append(msg, append([]byte(nil), tc.vals[i]...))
			}
			expected := make([]pair, len(tc.keys))
			for i := range tc.keys {
				expected[i] = pair{append([]byte(nil), tc.keys[i]...), append([]byte(nil), tc.vals[i]...)}
			}

			seq, tail := ReadMapBytes[[]byte, []byte](msg, ReadBytesZC, ReadBytesZC)
			var got []pair
			for k, v := range seq {
				kk := append([]byte(nil), k...)
				vv := append([]byte(nil), v...)
				got = append(got, pair{kk, vv})
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(expected) {
				t.Fatalf("%s count got %d want %d", tc.name, len(got), len(expected))
			}
			match := func(p pair, set []pair) bool {
				for _, q := range set {
					if eqBytes(p.k, q.k) && eqBytes(p.v, q.v) {
						return true
					}
				}
				return false
			}
			for _, p := range expected {
				if !match(p, got) {
					t.Fatalf("%s missing pair key=%v val=%v", tc.name, p.k, p.v)
				}
			}
		case regCase[time.Time]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadTimeBytes, ReadTimeBytes)
			got := make(map[time.Time]time.Time, len(tc.keys))
			for k, v := range seq {
				got[k.UTC()] = v.UTC()
			}
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case regCase[time.Duration]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadDurationBytes, ReadDurationBytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case regCase[complex64]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadComplex64Bytes, ReadComplex64Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		case regCase[complex128]:
			msg := AppendMapHeader(nil, uint32(len(tc.keys)))
			for i := range tc.keys {
				msg = tc.append(msg, tc.keys[i])
				msg = tc.append(msg, tc.vals[i])
			}
			seq, tail := ReadMapBytes(msg, ReadComplex128Bytes, ReadComplex128Bytes)
			got := maps.Collect(seq)
			remain, err := tail()
			if err != nil {
				t.Fatalf("%s tail: %v", tc.name, err)
			}
			if len(remain) != 0 {
				t.Fatalf("%s remain %d", tc.name, len(remain))
			}
			if len(got) != len(tc.keys) {
				t.Fatalf("%s size mismatch", tc.name)
			}
			for i := range tc.keys {
				if !tc.eq(got[tc.keys[i]], tc.vals[i]) {
					t.Fatalf("%s key %v got %v want %v", tc.name, tc.keys[i], got[tc.keys[i]], tc.vals[i])
				}
			}
		}
	}
}

func TestReadMapBytes_TailErrorOnTruncated(t *testing.T) {
	// Prepare a map with 2 entries: {1: 10, 2: 20}
	// Then truncate after encoding the second key to induce an error while reading the second value.
	msg := AppendMapHeader(nil, 2)
	msg = AppendInt(msg, 1)
	msg = AppendInt(msg, 10)
	msg = AppendInt(msg, 2)
	full := AppendInt(msg, 20)
	trunc := full[:len(full)-2] // truncate some bytes from the last value

	seq, tail := ReadMapBytes(trunc, ReadIntBytes, ReadIntBytes)
	got := maps.Collect(seq)
	// We expect only the first pair to be read
	if len(got) != 1 || got[1] != 10 {
		t.Fatalf("expected only first pair (1:10), got %v", got)
	}
	remain, err := tail()
	if err == nil {
		t.Fatalf("expected tail error due to truncation")
	}
	_ = remain // remaining bytes content is not strictly asserted here
}
