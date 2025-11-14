package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestGenericsMarshal(t *testing.T) {
	x := GenericTest[Fixed, *Fixed]{}
	x.B = &Fixed{A: 1.5}
	x.C = append(x.C, Fixed{A: 2.5})
	x.D = map[string]Fixed{"hello": Fixed{A: 2.5}}
	x.E.A = Fixed{A: 2.5}
	x.F = append(x.F, GenericTest2[Fixed, *Fixed, string]{A: Fixed{A: 3.5}})
	x.G = map[string]GenericTest2[Fixed, *Fixed, string]{"hello": {A: Fixed{A: 3.5}}}

	bts, err := x.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	got := GenericTest[Fixed, *Fixed]{}
	got.B = x.B // We must initialize this.
	*got.B = Fixed{}
	bts, err = got.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(x, got) {
		t.Errorf("\n got=%#v\nwant=%#v", got, x)
	}
}

func TestGenericsEncode(t *testing.T) {
	x := GenericTest[Fixed, *Fixed]{}
	x.B = &Fixed{A: 1.5}
	x.C = append(x.C, Fixed{A: 2.5})
	x.D = map[string]Fixed{"hello": Fixed{A: 2.5}}
	x.E.A = Fixed{A: 2.5}
	x.F = append(x.F, GenericTest2[Fixed, *Fixed, string]{A: Fixed{A: 3.5}})
	x.G = map[string]GenericTest2[Fixed, *Fixed, string]{"hello": {A: Fixed{A: 3.5}}}

	var buf bytes.Buffer
	w := msgp.NewWriter(&buf)
	err := x.EncodeMsg(w)
	if err != nil {
		t.Fatal(err)
	}
	w.Flush()
	got := GenericTest[Fixed, *Fixed]{}
	got.B = x.B // We must initialize this.
	*got.B = Fixed{}
	r := msgp.NewReader(&buf)
	err = got.DecodeMsg(r)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(x, got) {
		t.Errorf("\n got=%#v\nwant=%#v", got, x)
	}
}

// Test for generic alias types (slices, maps, arrays) - these verify the type parameter fix
func TestGenericTest3List(t *testing.T) {
	// Test GenericTest3List
	list := GenericTest3List[Fixed, Fixed, *Fixed, *Fixed]{
		{A: Fixed{A: 1.5}, B: Fixed{A: 2.5}},
		{A: Fixed{A: 3.5}, B: Fixed{A: 4.5}},
	}

	// Test marshaling
	data, err := list.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("GenericTest3List.MarshalMsg failed: %v", err)
	}

	// Test unmarshaling
	var list2 GenericTest3List[Fixed, Fixed, *Fixed, *Fixed]
	_, err = list2.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("GenericTest3List.UnmarshalMsg failed: %v", err)
	}

	// Verify round-trip
	if len(list2) != 2 || list2[0].A.A != 1.5 || list2[0].B.A != 2.5 {
		t.Errorf("GenericTest3List round-trip failed")
	}
}

func TestGenericTest3Map(t *testing.T) {
	// Test GenericTest3Map
	testMap := GenericTest3Map[Fixed, Fixed, *Fixed, *Fixed]{
		"key1": {A: Fixed{A: 1.5}, B: Fixed{A: 2.5}},
		"key2": {A: Fixed{A: 3.5}, B: Fixed{A: 4.5}},
	}

	// Test marshaling
	data, err := testMap.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("GenericTest3Map.MarshalMsg failed: %v", err)
	}

	// Test unmarshaling
	var testMap2 GenericTest3Map[Fixed, Fixed, *Fixed, *Fixed]
	_, err = testMap2.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("GenericTest3Map.UnmarshalMsg failed: %v", err)
	}

	// Verify round-trip
	if len(testMap2) != 2 || testMap2["key1"].A.A != 1.5 || testMap2["key1"].B.A != 2.5 {
		t.Errorf("GenericTest3Map round-trip failed")
	}
}

func TestGenericTest3Array(t *testing.T) {
	// Test GenericTest3Array
	testArray := GenericTest3Array[Fixed, Fixed, *Fixed, *Fixed]{
		{A: Fixed{A: 1.5}, B: Fixed{A: 2.5}},
		{A: Fixed{A: 3.5}, B: Fixed{A: 4.5}},
		{A: Fixed{A: 5.5}, B: Fixed{A: 6.5}},
	}

	// Test marshaling
	data, err := testArray.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("GenericTest3Array.MarshalMsg failed: %v", err)
	}

	// Test unmarshaling
	var testArray2 GenericTest3Array[Fixed, Fixed, *Fixed, *Fixed]
	_, err = testArray2.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("GenericTest3Array.UnmarshalMsg failed: %v", err)
	}

	// Verify round-trip
	if testArray2[0].A.A != 1.5 || testArray2[0].B.A != 2.5 || testArray2[2].A.A != 5.5 {
		t.Errorf("GenericTest3Array round-trip failed")
	}
}
