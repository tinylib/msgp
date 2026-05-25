package setof

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func FuzzStringJSON(f *testing.F) {
	seeds := []string{
		"",
		"hello",
		"a\"b\\c",
		string([]byte{0x00, 0x01, 0x02}),
		string([]byte{0x00, 0x7F}),
		"\b\f\n\r\t",
		string([]byte{0x0B, 0x1F}),
		"こんにちは",
		string([]byte{0xFF, 0xFE}),
		string([]byte{0xC0, 0xC1}),
		string([]byte{0xED, 0xA0, 0x80}),
		string([]byte{0x80}),
		"🎉",
		"\"quoted\"",
		"line1\nline2",
		"</script>",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		set := String{s: {}}

		ourData, err := json.Marshal(set)
		if err != nil {
			t.Fatalf("our MarshalJSON failed: %v", err)
		}

		var fromOurs []string
		if err := json.Unmarshal(ourData, &fromOurs); err != nil {
			t.Fatalf("stdlib couldn't decode our JSON %q: %v", ourData, err)
		}
		if len(fromOurs) != 1 {
			t.Fatalf("expected 1 element, got %d", len(fromOurs))
		}

		stdData, err := json.Marshal([]string{s})
		if err != nil {
			t.Fatalf("stdlib MarshalJSON failed: %v", err)
		}
		var fromStd []string
		if err := json.Unmarshal(stdData, &fromStd); err != nil {
			t.Fatalf("stdlib couldn't decode its own JSON: %v", err)
		}

		if !reflect.DeepEqual(fromOurs, fromStd) {
			t.Errorf("decoded value differs from stdlib:\n  ours   = %q (% x)\n  stdlib = %q (% x)\n  ourJSON  = %s\n  stdJSON  = %s",
				fromOurs[0], fromOurs[0], fromStd[0], fromStd[0], ourData, stdData)
		}

		var rt String
		if err := json.Unmarshal(ourData, &rt); err != nil {
			t.Fatalf("our UnmarshalJSON failed: %v", err)
		}
		if len(rt) != 1 {
			t.Fatalf("round-trip wrong length: %d", len(rt))
		}
		for k := range rt {
			if k != fromOurs[0] {
				t.Errorf("round-trip lost data: got %q, want %q", k, fromOurs[0])
			}
		}

		var rtFromStd String
		if err := json.Unmarshal(stdData, &rtFromStd); err != nil {
			t.Fatalf("our UnmarshalJSON failed on stdlib JSON %q: %v", stdData, err)
		}
		if len(rtFromStd) != 1 {
			t.Fatalf("stdlib JSON wrong length: %d", len(rtFromStd))
		}
		for k := range rtFromStd {
			if k != fromStd[0] {
				t.Errorf("our decode of stdlib JSON wrong: got %q, want %q", k, fromStd[0])
			}
		}
	})
}

func FuzzIntJSON(f *testing.F) {
	for _, n := range []int{0, 1, -1, 42, -42, math.MaxInt, math.MinInt, math.MaxInt32, math.MinInt32} {
		f.Add(n)
	}
	f.Fuzz(func(t *testing.T, n int) {
		set := Int{n: {}}
		data, err := json.Marshal(set)
		if err != nil {
			t.Fatalf("MarshalJSON failed for %d: %v", n, err)
		}

		var fromStd []int
		if err := json.Unmarshal(data, &fromStd); err != nil {
			t.Fatalf("stdlib couldn't decode our JSON %q: %v", data, err)
		}
		if len(fromStd) != 1 || fromStd[0] != n {
			t.Errorf("stdlib decoded %v, want [%d]", fromStd, n)
		}

		var rt Int
		if err := json.Unmarshal(data, &rt); err != nil {
			t.Fatalf("UnmarshalJSON failed for %q: %v", data, err)
		}
		if _, ok := rt[n]; !ok || len(rt) != 1 {
			t.Errorf("round-trip lost %d: got %v", n, rt)
		}
	})
}

func FuzzFloat64JSON(f *testing.F) {
	for _, x := range []float64{0, 1, -1, 0.5, -0.5, 1.5e10, -1.5e-10, math.MaxFloat64, math.SmallestNonzeroFloat64} {
		f.Add(x)
	}
	f.Fuzz(func(t *testing.T, x float64) {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			t.Skip("NaN/Inf not representable in JSON")
		}
		set := Float64{x: {}}
		data, err := json.Marshal(set)
		if err != nil {
			t.Fatalf("MarshalJSON failed for %g: %v", x, err)
		}

		var rt Float64
		if err := json.Unmarshal(data, &rt); err != nil {
			t.Fatalf("UnmarshalJSON failed for %q: %v", data, err)
		}
		if _, ok := rt[x]; !ok || len(rt) != 1 {
			t.Errorf("round-trip lost %g: got %v (JSON=%s)", x, rt, data)
		}
	})
}
