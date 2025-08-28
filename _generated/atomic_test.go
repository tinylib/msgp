package _generated

import (
	"bytes"
	"testing"

	"sync/atomic"

	"github.com/tinylib/msgp/msgp"
)

// Helper: compare two AtomicTest values deeply.
func equalAtomicTest(t *testing.T, a, b *AtomicTest) {
	t.Helper()

	// Top-level atomics
	if a.A.Load() != b.A.Load() {
		t.Fatalf("A mismatch: %d != %d", a.A.Load(), b.A.Load())
	}
	if a.B.Load() != b.B.Load() {
		t.Fatalf("B mismatch: %d != %d", a.B.Load(), b.B.Load())
	}
	if a.C.Load() != b.C.Load() {
		t.Fatalf("C mismatch: %d != %d", a.C.Load(), b.C.Load())
	}
	if a.D.Load() != b.D.Load() {
		t.Fatalf("D mismatch: %d != %d", a.D.Load(), b.D.Load())
	}
	if a.E.Load() != b.E.Load() {
		t.Fatalf("E mismatch: %t != %t", a.E.Load(), b.E.Load())
	}

	// Omitted (non-pointer) atomics
	if a.Omitted.A.Load() != b.Omitted.A.Load() {
		t.Fatalf("Omitted.A mismatch: %d != %d", a.Omitted.A.Load(), b.Omitted.A.Load())
	}
	if a.Omitted.B.Load() != b.Omitted.B.Load() {
		t.Fatalf("Omitted.B mismatch: %d != %d", a.Omitted.B.Load(), b.Omitted.B.Load())
	}
	if a.Omitted.C.Load() != b.Omitted.C.Load() {
		t.Fatalf("Omitted.C mismatch: %d != %d", a.Omitted.C.Load(), b.Omitted.C.Load())
	}
	if a.Omitted.D.Load() != b.Omitted.D.Load() {
		t.Fatalf("Omitted.D mismatch: %d != %d", a.Omitted.D.Load(), b.Omitted.D.Load())
	}
	if a.Omitted.E.Load() != b.Omitted.E.Load() {
		t.Fatalf("Omitted.E mismatch: %t != %t", a.Omitted.E.Load(), b.Omitted.E.Load())
	}

	// Pointer section
	switch {
	case (a.Ptr.A == nil) != (b.Ptr.A == nil):
		t.Fatalf("Ptr.A nil mismatch: %v vs %v", a.Ptr.A == nil, b.Ptr.A == nil)
	case a.Ptr.A != nil && a.Ptr.A.Load() != b.Ptr.A.Load():
		t.Fatalf("Ptr.A mismatch: %d != %d", a.Ptr.A.Load(), b.Ptr.A.Load())
	}
	switch {
	case (a.Ptr.B == nil) != (b.Ptr.B == nil):
		t.Fatalf("Ptr.B nil mismatch: %v vs %v", a.Ptr.B == nil, b.Ptr.B == nil)
	case a.Ptr.B != nil && a.Ptr.B.Load() != b.Ptr.B.Load():
		t.Fatalf("Ptr.B mismatch: %d != %d", a.Ptr.B.Load(), b.Ptr.B.Load())
	}
	switch {
	case (a.Ptr.C == nil) != (b.Ptr.C == nil):
		t.Fatalf("Ptr.C nil mismatch: %v vs %v", a.Ptr.C == nil, b.Ptr.C == nil)
	case a.Ptr.C != nil && a.Ptr.C.Load() != b.Ptr.C.Load():
		t.Fatalf("Ptr.C mismatch: %d != %d", a.Ptr.C.Load(), b.Ptr.C.Load())
	}
	switch {
	case (a.Ptr.D == nil) != (b.Ptr.D == nil):
		t.Fatalf("Ptr.D nil mismatch: %v vs %v", a.Ptr.D == nil, b.Ptr.D == nil)
	case a.Ptr.D != nil && a.Ptr.D.Load() != b.Ptr.D.Load():
		t.Fatalf("Ptr.D mismatch: %d != %d", a.Ptr.D.Load(), b.Ptr.D.Load())
	}
	switch {
	case (a.Ptr.E == nil) != (b.Ptr.E == nil):
		t.Fatalf("Ptr.E nil mismatch: %v vs %v", a.Ptr.E == nil, b.Ptr.E == nil)
	case a.Ptr.E != nil && a.Ptr.E.Load() != b.Ptr.E.Load():
		t.Fatalf("Ptr.E mismatch: %t != %t", a.Ptr.E.Load(), b.Ptr.E.Load())
	}

	// Ptr.Omitted pointers
	switch {
	case (a.Ptr.Omitted.A == nil) != (b.Ptr.Omitted.A == nil):
		t.Fatalf("Ptr.Omitted.A nil mismatch: %v vs %v", a.Ptr.Omitted.A == nil, b.Ptr.Omitted.A == nil)
	case a.Ptr.Omitted.A != nil && a.Ptr.Omitted.A.Load() != b.Ptr.Omitted.A.Load():
		t.Fatalf("Ptr.Omitted.A mismatch: %d != %d", a.Ptr.Omitted.A.Load(), b.Ptr.Omitted.A.Load())
	}
	switch {
	case (a.Ptr.Omitted.B == nil) != (b.Ptr.Omitted.B == nil):
		t.Fatalf("Ptr.Omitted.B nil mismatch: %v vs %v", a.Ptr.Omitted.B == nil, b.Ptr.Omitted.B == nil)
	case a.Ptr.Omitted.B != nil && a.Ptr.Omitted.B.Load() != b.Ptr.Omitted.B.Load():
		t.Fatalf("Ptr.Omitted.B mismatch: %d != %d", a.Ptr.Omitted.B.Load(), b.Ptr.Omitted.B.Load())
	}
	switch {
	case (a.Ptr.Omitted.C == nil) != (b.Ptr.Omitted.C == nil):
		t.Fatalf("Ptr.Omitted.C nil mismatch: %v vs %v", a.Ptr.Omitted.C == nil, b.Ptr.Omitted.C == nil)
	case a.Ptr.Omitted.C != nil && a.Ptr.Omitted.C.Load() != b.Ptr.Omitted.C.Load():
		t.Fatalf("Ptr.Omitted.C mismatch: %d != %d", a.Ptr.Omitted.C.Load(), b.Ptr.Omitted.C.Load())
	}
	switch {
	case (a.Ptr.Omitted.D == nil) != (b.Ptr.Omitted.D == nil):
		t.Fatalf("Ptr.Omitted.D nil mismatch: %v vs %v", a.Ptr.Omitted.D == nil, b.Ptr.Omitted.D == nil)
	case a.Ptr.Omitted.D != nil && a.Ptr.Omitted.D.Load() != b.Ptr.Omitted.D.Load():
		t.Fatalf("Ptr.Omitted.D mismatch: %d != %d", a.Ptr.Omitted.D.Load(), b.Ptr.Omitted.D.Load())
	}
	switch {
	case (a.Ptr.Omitted.E == nil) != (b.Ptr.Omitted.E == nil):
		t.Fatalf("Ptr.Omitted.E nil mismatch: %v vs %v", a.Ptr.Omitted.E == nil, b.Ptr.Omitted.E == nil)
	case a.Ptr.Omitted.E != nil && a.Ptr.Omitted.E.Load() != b.Ptr.Omitted.E.Load():
		t.Fatalf("Ptr.Omitted.E mismatch: %t != %t", a.Ptr.Omitted.E.Load(), b.Ptr.Omitted.E.Load())
	}

	// F slice
	if len(a.F) != len(b.F) {
		t.Fatalf("F length mismatch: %d != %d", len(a.F), len(b.F))
	}
	for i := range a.F {
		if a.F[i].Load() != b.F[i].Load() {
			t.Fatalf("F[%d] mismatch: %d != %d", i, a.F[i].Load(), b.F[i].Load())
		}
	}

	// G slice of pointers
	if len(a.G) != len(b.G) {
		t.Fatalf("G length mismatch: %d != %d", len(a.G), len(b.G))
	}
	for i := range a.G {
		if (a.G[i] == nil) != (b.G[i] == nil) {
			t.Fatalf("G[%d] nil mismatch: %v vs %v", i, a.G[i] == nil, b.G[i] == nil)
		}
		if a.G[i] != nil && a.G[i].Load() != b.G[i].Load() {
			t.Fatalf("G[%d] mismatch: %d != %d", i, a.G[i].Load(), b.G[i].Load())
		}
	}

	// H map[string]*atomic.Int32
	if len(a.H) != len(b.H) {
		t.Fatalf("H length mismatch: %d != %d", len(a.H), len(b.H))
	}
	for k, va := range a.H {
		vb, ok := b.H[k]
		if !ok {
			t.Fatalf("H missing key %q in b", k)
		}
		if (va == nil) != (vb == nil) {
			t.Fatalf("H[%q] nil mismatch: %v vs %v", k, va == nil, vb == nil)
		}
		if va != nil && va.Load() != vb.Load() {
			t.Fatalf("H[%q] mismatch: %d != %d", k, va.Load(), vb.Load())
		}
	}
}

// Build a populated AtomicTest for roundtrip testing.
func makeAtomicTestSample() *AtomicTest {
	var at AtomicTest

	// Top-level
	at.A.Store(-123)
	at.B.Store(456)
	at.C.Store(-789)
	at.D.Store(1011)
	at.E.Store(true)

	// Omitted (present even if zero, but we'll set some non-zero)
	at.Omitted.A.Store(1)
	at.Omitted.B.Store(2)
	// leave C zero to exercise omitempty + clearing path
	at.Omitted.D.Store(4)
	at.Omitted.E.Store(false) // explicit false

	// Ptr section
	at.Ptr.A = new(atomic.Int64)
	at.Ptr.A.Store(-77)
	at.Ptr.B = new(atomic.Uint64)
	at.Ptr.B.Store(88)
	// leave C nil to exercise pointer omitempty
	at.Ptr.D = new(atomic.Uint32)
	at.Ptr.D.Store(99)
	at.Ptr.E = new(atomic.Bool)
	at.Ptr.E.Store(true)

	// Ptr.Omitted
	at.Ptr.Omitted.A = new(atomic.Int64)
	at.Ptr.Omitted.A.Store(5)
	// leave B nil
	at.Ptr.Omitted.C = new(atomic.Int32)
	at.Ptr.Omitted.C.Store(-6)
	// leave D nil
	at.Ptr.Omitted.E = new(atomic.Bool)
	at.Ptr.Omitted.E.Store(false)

	// F slice (values)
	at.F = make([]atomic.Int32, 3)
	at.F[0].Store(10)
	at.F[1].Store(-20)
	at.F[2].Store(30)

	// G slice (pointers)
	at.G = make([]*atomic.Int32, 3)
	at.G[0] = new(atomic.Int32)
	at.G[0].Store(100)
	// leave G[1] nil
	at.G[2] = new(atomic.Int32)
	at.G[2].Store(-300)

	// H map
	at.H = make(map[string]*atomic.Int32)
	v1 := new(atomic.Int32)
	v1.Store(7)
	at.H["x"] = v1
	// include a nil pointer in the map as well
	at.H["nil"] = nil

	return &at
}

func TestAtomicTest_Roundtrip_EncodeDecode(t *testing.T) {
	orig := makeAtomicTestSample()

	var buf bytes.Buffer
	w := msgp.NewWriter(&buf)
	if err := orig.EncodeMsg(w); err != nil {
		t.Fatalf("EncodeMsg failed: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush failed: %v", err)
	}

	r := msgp.NewReader(bytes.NewReader(buf.Bytes()))
	var out AtomicTest
	if err := out.DecodeMsg(r); err != nil {
		t.Fatalf("DecodeMsg failed: %v", err)
	}

	equalAtomicTest(t, orig, &out)
}

func TestAtomicTest_Roundtrip_MarshalUnmarshal(t *testing.T) {
	orig := makeAtomicTestSample()

	bts, err := orig.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("MarshalMsg failed: %v", err)
	}

	var out AtomicTest
	left, err := out.UnmarshalMsg(bts)
	if err != nil {
		t.Fatalf("UnmarshalMsg failed: %v", err)
	}
	if len(left) != 0 {
		t.Fatalf("leftover bytes after UnmarshalMsg: %d", len(left))
	}

	equalAtomicTest(t, orig, &out)
}
