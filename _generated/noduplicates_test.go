package _generated

import (
	"bytes"
	"errors"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

// --- helpers ---

func mustUnmarshal(t *testing.T, buf []byte, v msgp.Unmarshaler) {
	t.Helper()
	_, err := v.UnmarshalMsg(buf)
	if err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
}

func mustDecode(t *testing.T, buf []byte, v msgp.Decodable) {
	t.Helper()
	reader := msgp.NewReader(bytes.NewReader(buf))
	err := v.DecodeMsg(reader)
	if err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}
}

func expectDup(t *testing.T, err error) {
	t.Helper()
	if !errors.Is(err, msgp.ErrDuplicateEntry) {
		t.Fatalf("expected ErrDuplicateEntry, got %v", err)
	}
}

func expectNoDup(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------- NoDupSimple basic struct ----------

func TestNoDupSimple_RoundTrip(t *testing.T) {
	orig := NoDupSimple{Foo: "hello", Bar: 42, Baz: true}

	t.Run("MarshalUnmarshal", func(t *testing.T) {
		buf, err := orig.MarshalMsg(nil)
		expectNoDup(t, err)
		var out NoDupSimple
		mustUnmarshal(t, buf, &out)
		if out != orig {
			t.Fatalf("mismatch: got %+v, want %+v", out, orig)
		}
	})
	t.Run("EncodeDecode", func(t *testing.T) {
		var buf bytes.Buffer
		w := msgp.NewWriter(&buf)
		expectNoDup(t, orig.EncodeMsg(w))
		w.Flush()
		var out NoDupSimple
		mustDecode(t, buf.Bytes(), &out)
		if out != orig {
			t.Fatalf("mismatch: got %+v, want %+v", out, orig)
		}
	})
}

func TestNoDupSimple_DuplicateField(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 3)
		buf = msgp.AppendString(buf, "Foo")
		buf = msgp.AppendString(buf, "first")
		buf = msgp.AppendString(buf, "Bar")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "Foo") // duplicate
		buf = msgp.AppendString(buf, "second")
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupSimple
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupSimple
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupSimple_UnknownFieldsOK(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 3)
	buf = msgp.AppendString(buf, "Foo")
	buf = msgp.AppendString(buf, "hello")
	buf = msgp.AppendString(buf, "Unknown1")
	buf = msgp.AppendString(buf, "skip1")
	buf = msgp.AppendString(buf, "Unknown2")
	buf = msgp.AppendString(buf, "skip2")

	var out NoDupSimple
	mustUnmarshal(t, buf, &out)
	if out.Foo != "hello" {
		t.Fatalf("expected Foo=hello, got %q", out.Foo)
	}
}

func TestNoDupSimple_EmptyMessage(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 0)
	var out NoDupSimple
	mustUnmarshal(t, buf, &out)
}

func TestNoDupSimple_SingleField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 1)
	buf = msgp.AppendString(buf, "Foo")
	buf = msgp.AppendString(buf, "hello")
	var out NoDupSimple
	mustUnmarshal(t, buf, &out)
	if out.Foo != "hello" {
		t.Fatalf("expected Foo=hello, got %q", out.Foo)
	}
}

func TestNoDupSimple_AllFieldsDuplicated(t *testing.T) {
	// Each field appears exactly once - should work.
	buf := msgp.AppendMapHeader(nil, 3)
	buf = msgp.AppendString(buf, "Foo")
	buf = msgp.AppendString(buf, "f")
	buf = msgp.AppendString(buf, "Bar")
	buf = msgp.AppendInt(buf, 1)
	buf = msgp.AppendString(buf, "Baz")
	buf = msgp.AppendBool(buf, true)
	var out NoDupSimple
	mustUnmarshal(t, buf, &out)
}

func TestNoDupSimple_DuplicateSecondField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 3)
	buf = msgp.AppendString(buf, "Bar")
	buf = msgp.AppendInt(buf, 1)
	buf = msgp.AppendString(buf, "Baz")
	buf = msgp.AppendBool(buf, true)
	buf = msgp.AppendString(buf, "Bar") // dup of second field
	buf = msgp.AppendInt(buf, 2)
	var out NoDupSimple
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupSimple_DuplicateThirdField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "Baz")
	buf = msgp.AppendBool(buf, false)
	buf = msgp.AppendString(buf, "Baz") // dup of third field
	buf = msgp.AppendBool(buf, true)
	var out NoDupSimple
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

// ---------- NoDupMap standalone map ----------

func TestNoDupMap_RoundTrip(t *testing.T) {
	orig := NoDupMap{"a": "1", "b": "2", "c": "3"}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupMap
	mustUnmarshal(t, buf, &out)
	if len(out) != len(orig) {
		t.Fatalf("length mismatch: got %d, want %d", len(out), len(orig))
	}
}

func TestNoDupMap_DuplicateKey(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 3)
		buf = msgp.AppendString(buf, "key1")
		buf = msgp.AppendString(buf, "val1")
		buf = msgp.AppendString(buf, "key2")
		buf = msgp.AppendString(buf, "val2")
		buf = msgp.AppendString(buf, "key1")
		buf = msgp.AppendString(buf, "val3")
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupMap
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupMap
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupMap_EmptyAndSingle(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 0)
		var out NoDupMap
		mustUnmarshal(t, buf, &out)
	})
	t.Run("Single", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "k")
		buf = msgp.AppendString(buf, "v")
		var out NoDupMap
		mustUnmarshal(t, buf, &out)
		if out["k"] != "v" {
			t.Fatal("mismatch")
		}
	})
}

// ---------- NoDupMapInt ----------

func TestNoDupMapInt_DuplicateKey(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 3)
		buf = msgp.AppendString(buf, "a")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "b")
		buf = msgp.AppendInt(buf, 2)
		buf = msgp.AppendString(buf, "a")
		buf = msgp.AppendInt(buf, 3)
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupMapInt
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupMapInt
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

// ---------- NoDupMapComplex: map[string]NoDupInner ----------

func TestNoDupMapComplex_DuplicateOuterKey(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 2)
		// first entry
		buf = msgp.AppendString(buf, "obj1")
		buf = msgp.AppendMapHeader(buf, 1)
		buf = msgp.AppendString(buf, "x")
		buf = msgp.AppendInt(buf, 10)
		// duplicate outer key
		buf = msgp.AppendString(buf, "obj1")
		buf = msgp.AppendMapHeader(buf, 1)
		buf = msgp.AppendString(buf, "x")
		buf = msgp.AppendInt(buf, 20)
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupMapComplex
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupMapComplex
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupMapComplex_DuplicateInnerField(t *testing.T) {
	// Outer keys are unique but the inner struct has a duplicate field.
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "obj1")
		buf = msgp.AppendMapHeader(buf, 3)
		buf = msgp.AppendString(buf, "x")
		buf = msgp.AppendInt(buf, 10)
		buf = msgp.AppendString(buf, "y")
		buf = msgp.AppendString(buf, "hi")
		buf = msgp.AppendString(buf, "x") // dup inner field
		buf = msgp.AppendInt(buf, 20)
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupMapComplex
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupMapComplex
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupMapComplex_RoundTrip(t *testing.T) {
	orig := NoDupMapComplex{
		"a": {X: 1, Y: "hello"},
		"b": {X: 2, Y: "world"},
	}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupMapComplex
	mustUnmarshal(t, buf, &out)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

// ---------- NoDupNested: struct with map fields ----------

func TestNoDupNested_RoundTrip(t *testing.T) {
	orig := NoDupNested{
		Name:  "test",
		Items: map[string]int{"a": 1, "b": 2},
		Tags:  map[string]string{"x": "y"},
	}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupNested
	mustUnmarshal(t, buf, &out)
	if out.Name != orig.Name {
		t.Fatalf("Name mismatch")
	}
}

func TestNoDupNested_DuplicateStructField(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 3)
		buf = msgp.AppendString(buf, "Name")
		buf = msgp.AppendString(buf, "first")
		buf = msgp.AppendString(buf, "Items")
		buf = msgp.AppendMapHeader(buf, 0)
		buf = msgp.AppendString(buf, "Name")
		buf = msgp.AppendString(buf, "second")
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupNested
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupNested
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupNested_DuplicateMapKey(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "Items")
		buf = msgp.AppendMapHeader(buf, 3)
		buf = msgp.AppendString(buf, "x")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "y")
		buf = msgp.AppendInt(buf, 2)
		buf = msgp.AppendString(buf, "x")
		buf = msgp.AppendInt(buf, 3)
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupNested
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupNested
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

// ---------- NoDupManyFields: 10 fields, 2-byte bitmask ----------

func TestNoDupManyFields_RoundTrip(t *testing.T) {
	orig := NoDupManyFields{
		F01: "a", F02: "b", F03: "c", F04: "d", F05: "e",
		F06: "f", F07: "g", F08: "h", F09: "i", F10: "j",
	}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupManyFields
	mustUnmarshal(t, buf, &out)
	if out != orig {
		t.Fatalf("mismatch")
	}
}

func TestNoDupManyFields_DupFirstField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "f01")
	buf = msgp.AppendString(buf, "a")
	buf = msgp.AppendString(buf, "f01")
	buf = msgp.AppendString(buf, "b")
	var out NoDupManyFields
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupManyFields_DupField8(t *testing.T) {
	// f08 is bit 7 of byte 0 - the last bit in the first byte.
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "f08")
	buf = msgp.AppendString(buf, "a")
	buf = msgp.AppendString(buf, "f08")
	buf = msgp.AppendString(buf, "b")
	var out NoDupManyFields
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupManyFields_DupField9(t *testing.T) {
	// f09 is bit 0 of byte 1 - crosses the byte boundary.
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 3)
		buf = msgp.AppendString(buf, "f09")
		buf = msgp.AppendString(buf, "first")
		buf = msgp.AppendString(buf, "f01")
		buf = msgp.AppendString(buf, "ok")
		buf = msgp.AppendString(buf, "f09")
		buf = msgp.AppendString(buf, "second")
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupManyFields
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupManyFields
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupManyFields_DupLastField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "f10")
	buf = msgp.AppendString(buf, "a")
	buf = msgp.AppendString(buf, "f10")
	buf = msgp.AppendString(buf, "b")
	var out NoDupManyFields
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

// ---------- NoDup17Fields: 17 fields, 3-byte bitmask ----------

func TestNoDup17Fields_RoundTrip(t *testing.T) {
	orig := NoDup17Fields{
		A01: 1, A02: 2, A03: 3, A04: 4, A05: 5, A06: 6, A07: 7, A08: 8,
		A09: 9, A10: 10, A11: 11, A12: 12, A13: 13, A14: 14, A15: 15, A16: 16, A17: 17,
	}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDup17Fields
	mustUnmarshal(t, buf, &out)
	if out != orig {
		t.Fatalf("mismatch")
	}
}

func TestNoDup17Fields_DupField17(t *testing.T) {
	// a17 is bit 0 of byte 2.
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "a17")
	buf = msgp.AppendInt(buf, 1)
	buf = msgp.AppendString(buf, "a17")
	buf = msgp.AppendInt(buf, 2)
	var out NoDup17Fields
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDup17Fields_DupField16(t *testing.T) {
	// a16 is bit 7 of byte 1 - last bit of second byte.
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "a16")
	buf = msgp.AppendInt(buf, 1)
	buf = msgp.AppendString(buf, "a16")
	buf = msgp.AppendInt(buf, 2)
	var out NoDup17Fields
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

// ---------- NoDupRichTypes: mixed field types ----------

func TestNoDupRichTypes_RoundTrip(t *testing.T) {
	s := "ptr"
	orig := NoDupRichTypes{
		Str:      "hello",
		Num:      3.14,
		Data:     []byte{1, 2, 3},
		PtrStr:   &s,
		Strings:  []string{"a", "b"},
		Inner:    NoDupInner{X: 10, Y: "y"},
		PtrInner: &NoDupInner{X: 20, Y: "z"},
		MapSlice: map[string][]int{"k": {1, 2}},
	}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupRichTypes
	mustUnmarshal(t, buf, &out)
	if out.Str != orig.Str || out.Num != orig.Num {
		t.Fatalf("basic field mismatch")
	}
	if *out.PtrStr != *orig.PtrStr {
		t.Fatalf("PtrStr mismatch")
	}
}

func TestNoDupRichTypes_DupStringField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "str")
	buf = msgp.AppendString(buf, "first")
	buf = msgp.AppendString(buf, "str")
	buf = msgp.AppendString(buf, "second")
	var out NoDupRichTypes
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupRichTypes_DupNumField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "num")
	buf = msgp.AppendFloat64(buf, 1.0)
	buf = msgp.AppendString(buf, "num")
	buf = msgp.AppendFloat64(buf, 2.0)
	var out NoDupRichTypes
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupRichTypes_DupBytesField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "data")
	buf = msgp.AppendBytes(buf, []byte{1})
	buf = msgp.AppendString(buf, "data")
	buf = msgp.AppendBytes(buf, []byte{2})
	var out NoDupRichTypes
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupRichTypes_DupPointerField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "ptr_str")
	buf = msgp.AppendString(buf, "first")
	buf = msgp.AppendString(buf, "ptr_str")
	buf = msgp.AppendString(buf, "second")
	var out NoDupRichTypes
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupRichTypes_DupSliceField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "strings")
	buf = msgp.AppendArrayHeader(buf, 1)
	buf = msgp.AppendString(buf, "a")
	buf = msgp.AppendString(buf, "strings")
	buf = msgp.AppendArrayHeader(buf, 1)
	buf = msgp.AppendString(buf, "b")
	var out NoDupRichTypes
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupRichTypes_DupNestedStructField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "inner")
	buf = msgp.AppendMapHeader(buf, 1)
	buf = msgp.AppendString(buf, "x")
	buf = msgp.AppendInt(buf, 1)
	buf = msgp.AppendString(buf, "inner")
	buf = msgp.AppendMapHeader(buf, 1)
	buf = msgp.AppendString(buf, "x")
	buf = msgp.AppendInt(buf, 2)
	var out NoDupRichTypes
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupRichTypes_DupInsideNestedStruct(t *testing.T) {
	// inner struct itself has a duplicate field
	buf := msgp.AppendMapHeader(nil, 1)
	buf = msgp.AppendString(buf, "inner")
	buf = msgp.AppendMapHeader(buf, 3)
	buf = msgp.AppendString(buf, "x")
	buf = msgp.AppendInt(buf, 1)
	buf = msgp.AppendString(buf, "y")
	buf = msgp.AppendString(buf, "ok")
	buf = msgp.AppendString(buf, "x") // dup inside inner
	buf = msgp.AppendInt(buf, 2)
	var out NoDupRichTypes
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
	// Error should contain the full context path "Inner/X"
	errStr := err.Error()
	if !bytes.Contains([]byte(errStr), []byte("Inner")) {
		t.Fatalf("error should reference 'Inner', got: %s", errStr)
	}
	if !bytes.Contains([]byte(errStr), []byte("X")) {
		t.Fatalf("error should reference 'X', got: %s", errStr)
	}
}

func TestNoDupRichTypes_DupPtrStructField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "ptr_inner")
	buf = msgp.AppendMapHeader(buf, 1)
	buf = msgp.AppendString(buf, "x")
	buf = msgp.AppendInt(buf, 1)
	buf = msgp.AppendString(buf, "ptr_inner")
	buf = msgp.AppendMapHeader(buf, 1)
	buf = msgp.AppendString(buf, "x")
	buf = msgp.AppendInt(buf, 2)
	var out NoDupRichTypes
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupRichTypes_DupMapSliceKey(t *testing.T) {
	// map_slice field has a duplicate key
	buf := msgp.AppendMapHeader(nil, 1)
	buf = msgp.AppendString(buf, "map_slice")
	buf = msgp.AppendMapHeader(buf, 2)
	buf = msgp.AppendString(buf, "k")
	buf = msgp.AppendArrayHeader(buf, 1)
	buf = msgp.AppendInt(buf, 1)
	buf = msgp.AppendString(buf, "k") // dup map key
	buf = msgp.AppendArrayHeader(buf, 1)
	buf = msgp.AppendInt(buf, 2)
	var out NoDupRichTypes
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

// ---------- NoDupOmitempty: interaction with omitempty ----------

func TestNoDupOmitempty_RoundTrip(t *testing.T) {
	orig := NoDupOmitempty{Required: "r", Optional: "o", Flag: true, Count: 5}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupOmitempty
	mustUnmarshal(t, buf, &out)
	if out != orig {
		t.Fatalf("mismatch: got %+v, want %+v", out, orig)
	}
}

func TestNoDupOmitempty_DupRequiredField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "required")
	buf = msgp.AppendString(buf, "first")
	buf = msgp.AppendString(buf, "required")
	buf = msgp.AppendString(buf, "second")
	var out NoDupOmitempty
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupOmitempty_DupOmitemptyField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "optional")
	buf = msgp.AppendString(buf, "first")
	buf = msgp.AppendString(buf, "optional")
	buf = msgp.AppendString(buf, "second")
	var out NoDupOmitempty
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupOmitempty_PartialFieldsOK(t *testing.T) {
	// Only some fields present - no duplicates.
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "required")
	buf = msgp.AppendString(buf, "val")
	buf = msgp.AppendString(buf, "count")
	buf = msgp.AppendInt(buf, 42)
	var out NoDupOmitempty
	mustUnmarshal(t, buf, &out)
	if out.Required != "val" || out.Count != 42 {
		t.Fatalf("mismatch: got %+v", out)
	}
}

// ---------- NoDupAllownil: interaction with allownil ----------

func TestNoDupAllownil_RoundTrip(t *testing.T) {
	orig := NoDupAllownil{
		Name:    "test",
		Vals:    []byte{1, 2},
		Items:   []int{3, 4},
		Mapping: map[string]string{"k": "v"},
	}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupAllownil
	mustUnmarshal(t, buf, &out)
	if out.Name != orig.Name {
		t.Fatal("name mismatch")
	}
}

func TestNoDupAllownil_NilFieldsThenDup(t *testing.T) {
	// Send nil for allownil fields, then a duplicate struct field.
	buf := msgp.AppendMapHeader(nil, 3)
	buf = msgp.AppendString(buf, "vals")
	buf = msgp.AppendNil(buf)
	buf = msgp.AppendString(buf, "name")
	buf = msgp.AppendString(buf, "hello")
	buf = msgp.AppendString(buf, "vals") // dup
	buf = msgp.AppendNil(buf)
	var out NoDupAllownil
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupAllownil_DupNameField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "name")
	buf = msgp.AppendString(buf, "first")
	buf = msgp.AppendString(buf, "name")
	buf = msgp.AppendString(buf, "second")
	var out NoDupAllownil
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupAllownil_DupMapKey(t *testing.T) {
	// The "mapping" field's map has a duplicate key.
	buf := msgp.AppendMapHeader(nil, 1)
	buf = msgp.AppendString(buf, "mapping")
	buf = msgp.AppendMapHeader(buf, 2)
	buf = msgp.AppendString(buf, "k")
	buf = msgp.AppendString(buf, "v1")
	buf = msgp.AppendString(buf, "k") // dup inside map
	buf = msgp.AppendString(buf, "v2")
	var out NoDupAllownil
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupAllownil_AllNilNoDup(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 4)
	buf = msgp.AppendString(buf, "name")
	buf = msgp.AppendString(buf, "")
	buf = msgp.AppendString(buf, "vals")
	buf = msgp.AppendNil(buf)
	buf = msgp.AppendString(buf, "items")
	buf = msgp.AppendNil(buf)
	buf = msgp.AppendString(buf, "mapping")
	buf = msgp.AppendNil(buf)
	var out NoDupAllownil
	mustUnmarshal(t, buf, &out)
}

// ---------- Per-type directive (noduplicates_specific.go) ----------

func TestNoDupTargeted_DuplicateField(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 3)
		buf = msgp.AppendString(buf, "alpha")
		buf = msgp.AppendString(buf, "first")
		buf = msgp.AppendString(buf, "beta")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "alpha") // dup
		buf = msgp.AppendString(buf, "second")
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupTargeted
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupTargeted
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupTargeted_RoundTrip(t *testing.T) {
	orig := NoDupTargeted{Alpha: "a", Beta: 1, Gamma: true}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupTargeted
	mustUnmarshal(t, buf, &out)
	if out != orig {
		t.Fatal("mismatch")
	}
}

func TestNoDupTargetedMap_DuplicateKey(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 2)
		buf = msgp.AppendString(buf, "x")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "x") // dup
		buf = msgp.AppendInt(buf, 2)
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupTargetedMap
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupTargetedMap
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupUntargeted_DuplicatesAllowed(t *testing.T) {
	// NoDupUntargeted is NOT targeted by the directive - duplicates must be accepted.
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 3)
		buf = msgp.AppendString(buf, "alpha")
		buf = msgp.AppendString(buf, "first")
		buf = msgp.AppendString(buf, "beta")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "alpha") // dup - should be silently accepted
		buf = msgp.AppendString(buf, "second")
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupUntargeted
		_, err := out.UnmarshalMsg(mkBuf())
		expectNoDup(t, err)
		if out.Alpha != "second" {
			t.Fatalf("expected last value, got %q", out.Alpha)
		}
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupUntargeted
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectNoDup(t, out.DecodeMsg(r))
		if out.Alpha != "second" {
			t.Fatalf("expected last value, got %q", out.Alpha)
		}
	})
}

func TestNoDupUntargetedMap_DuplicatesAllowed(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 2)
		buf = msgp.AppendString(buf, "x")
		buf = msgp.AppendString(buf, "v1")
		buf = msgp.AppendString(buf, "x") // dup - should be silently accepted
		buf = msgp.AppendString(buf, "v2")
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupUntargetedMap
		_, err := out.UnmarshalMsg(mkBuf())
		expectNoDup(t, err)
		if out["x"] != "v2" {
			t.Fatalf("expected last value, got %q", out["x"])
		}
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupUntargetedMap
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectNoDup(t, out.DecodeMsg(r))
		if out["x"] != "v2" {
			t.Fatalf("expected last value, got %q", out["x"])
		}
	})
}

// ---------- Error quality ----------

func TestNoDup_ErrorContainsFieldName(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "Foo")
	buf = msgp.AppendString(buf, "first")
	buf = msgp.AppendString(buf, "Foo")
	buf = msgp.AppendString(buf, "second")
	var out NoDupSimple
	_, err := out.UnmarshalMsg(buf)
	if err == nil {
		t.Fatal("expected error")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("Foo")) {
		t.Fatalf("error should reference field 'Foo', got: %s", err.Error())
	}
}

func TestNoDup_ErrorIsUnwrappable(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "Foo")
	buf = msgp.AppendString(buf, "first")
	buf = msgp.AppendString(buf, "Foo")
	buf = msgp.AppendString(buf, "second")
	var out NoDupSimple
	_, err := out.UnmarshalMsg(buf)
	if err == nil {
		t.Fatal("expected error")
	}
	if cause := msgp.Cause(err); cause != msgp.ErrDuplicateEntry {
		t.Fatalf("Cause should be ErrDuplicateEntry, got %v", cause)
	}
	if !errors.Is(err, msgp.ErrDuplicateEntry) {
		t.Fatal("errors.Is should match ErrDuplicateEntry")
	}
}

func TestNoDupMap_ErrorContainsKey(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "mykey")
	buf = msgp.AppendString(buf, "v1")
	buf = msgp.AppendString(buf, "mykey")
	buf = msgp.AppendString(buf, "v2")
	var out NoDupMap
	_, err := out.UnmarshalMsg(buf)
	if err == nil {
		t.Fatal("expected error")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("mykey")) {
		t.Fatalf("error should reference key 'mykey', got: %s", err.Error())
	}
}

func TestNoDupNested_MapErrorContainsContext(t *testing.T) {
	// A duplicate inside the Items map should reference the field path.
	buf := msgp.AppendMapHeader(nil, 1)
	buf = msgp.AppendString(buf, "Items")
	buf = msgp.AppendMapHeader(buf, 2)
	buf = msgp.AppendString(buf, "dupkey")
	buf = msgp.AppendInt(buf, 1)
	buf = msgp.AppendString(buf, "dupkey")
	buf = msgp.AppendInt(buf, 2)
	var out NoDupNested
	_, err := out.UnmarshalMsg(buf)
	if err == nil {
		t.Fatal("expected error")
	}
	errStr := err.Error()
	if !bytes.Contains([]byte(errStr), []byte("Items")) {
		t.Fatalf("error should reference 'Items', got: %s", errStr)
	}
	if !bytes.Contains([]byte(errStr), []byte("dupkey")) {
		t.Fatalf("error should reference 'dupkey', got: %s", errStr)
	}
}

// ---------- Binary key maps (noduplicates_binkey.go) ----------

func TestNoDupBinKeyInt_RoundTrip(t *testing.T) {
	orig := NoDupBinKeyInt{1: "a", 2: "b", 3: "c"}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupBinKeyInt
	mustUnmarshal(t, buf, &out)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestNoDupBinKeyInt_DuplicateKey(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 3)
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "a")
		buf = msgp.AppendInt(buf, 2)
		buf = msgp.AppendString(buf, "b")
		buf = msgp.AppendInt(buf, 1) // dup
		buf = msgp.AppendString(buf, "c")
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupBinKeyInt
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupBinKeyInt
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupBinKeyUint32_DuplicateKey(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 2)
		buf = msgp.AppendUint32(buf, 42)
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendUint32(buf, 42) // dup
		buf = msgp.AppendInt(buf, 2)
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupBinKeyUint32
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupBinKeyUint32
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupBinKeyArray_DuplicateKey(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 2)
		buf = msgp.AppendBytes(buf, []byte{0x01, 0x02})
		buf = msgp.AppendInt(buf, 10)
		buf = msgp.AppendBytes(buf, []byte{0x01, 0x02}) // dup
		buf = msgp.AppendInt(buf, 20)
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupBinKeyArray
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupBinKeyArray
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupBinKeyArray_UniqueKeys(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendBytes(buf, []byte{0x01, 0x02})
	buf = msgp.AppendInt(buf, 10)
	buf = msgp.AppendBytes(buf, []byte{0x03, 0x04})
	buf = msgp.AppendInt(buf, 20)
	var out NoDupBinKeyArray
	mustUnmarshal(t, buf, &out)
	if out[[2]byte{0x01, 0x02}] != 10 || out[[2]byte{0x03, 0x04}] != 20 {
		t.Fatalf("mismatch: %+v", out)
	}
}

func TestNoDupBinKeyStruct_DupStructField(t *testing.T) {
	buf := msgp.AppendMapHeader(nil, 2)
	buf = msgp.AppendString(buf, "name")
	buf = msgp.AppendString(buf, "first")
	buf = msgp.AppendString(buf, "name") // dup
	buf = msgp.AppendString(buf, "second")
	var out NoDupBinKeyStruct
	_, err := out.UnmarshalMsg(buf)
	expectDup(t, err)
}

func TestNoDupBinKeyStruct_DupIntMapKey(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "int_map")
		buf = msgp.AppendMapHeader(buf, 2)
		buf = msgp.AppendInt(buf, 5)
		buf = msgp.AppendString(buf, "five")
		buf = msgp.AppendInt(buf, 5) // dup
		buf = msgp.AppendString(buf, "cinq")
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupBinKeyStruct
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupBinKeyStruct
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupBinKeyStruct_DupStrMapKey(t *testing.T) {
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "str_map")
		buf = msgp.AppendMapHeader(buf, 2)
		buf = msgp.AppendString(buf, "k")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "k") // dup
		buf = msgp.AppendInt(buf, 2)
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupBinKeyStruct
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupBinKeyStruct
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

// ---------- Shimmed key maps (noduplicates_shim.go) ----------

func TestNoDupShimMap_RoundTrip(t *testing.T) {
	orig := NoDupShimMap{"foo": 1, "bar": 2}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupShimMap
	mustUnmarshal(t, buf, &out)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestNoDupShimMap_DuplicateWireKey(t *testing.T) {
	// Same wire string sent twice.
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 2)
		buf = msgp.AppendString(buf, "foo")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "foo") // dup
		buf = msgp.AppendInt(buf, 2)
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupShimMap
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupShimMap
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupShimMap_DifferentWireKeysSameShimmedValue(t *testing.T) {
	// "Foo" and "foo" are different wire strings but noDupShimDec lowercases both,
	// so they resolve to the same map key. Should be detected as duplicate.
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 2)
		buf = msgp.AppendString(buf, "Foo")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "foo") // different wire key, same after shim
		buf = msgp.AppendInt(buf, 2)
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupShimMap
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupShimMap
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupShimStruct_DupShimMapKey(t *testing.T) {
	// Shimmed map field inside a struct with keys that collapse.
	mkBuf := func() []byte {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "data")
		buf = msgp.AppendMapHeader(buf, 2)
		buf = msgp.AppendString(buf, "Hello")
		buf = msgp.AppendString(buf, "world")
		buf = msgp.AppendString(buf, "hello") // collapses to same key
		buf = msgp.AppendString(buf, "again")
		return buf
	}
	t.Run("Unmarshal", func(t *testing.T) {
		var out NoDupShimStruct
		_, err := out.UnmarshalMsg(mkBuf())
		expectDup(t, err)
	})
	t.Run("Decode", func(t *testing.T) {
		var out NoDupShimStruct
		r := msgp.NewReader(bytes.NewReader(mkBuf()))
		expectDup(t, out.DecodeMsg(r))
	})
}

func TestNoDupShimStruct_RoundTrip(t *testing.T) {
	orig := NoDupShimStruct{
		Name: "test",
		Data: map[NoDupShimKey]string{"foo": "bar", "baz": "qux"},
	}
	buf, err := orig.MarshalMsg(nil)
	expectNoDup(t, err)
	var out NoDupShimStruct
	mustUnmarshal(t, buf, &out)
	if out.Name != "test" || len(out.Data) != 2 {
		t.Fatalf("mismatch: %+v", out)
	}
}
