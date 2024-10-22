package _generated

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/tinylib/msgp/msgp"
)

func TestClearOmitted(t *testing.T) {
	cleared := ClearOmitted0{}
	encoded, err := cleared.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	vPtr := NamedStringCO("value")
	filled := ClearOmitted0{
		AStruct:             ClearOmittedA{A: "something"},
		BStruct:             ClearOmittedA{A: "somthing"},
		AStructPtr:          &ClearOmittedA{A: "something"},
		AExt:                OmitZeroExt{25},
		APtrNamedStr:        &vPtr,
		ANamedStruct:        NamedStructCO{A: "value"},
		APtrNamedStruct:     &NamedStructCO{A: "sdf"},
		EmbeddableStructCO:  EmbeddableStructCO{"value"},
		EmbeddableStructCO2: EmbeddableStructCO2{"value"},
		ATime:               time.Now(),
		ASlice:              []int{1, 2, 3},
		AMap:                map[string]int{"1": 1},
		ABin:                []byte{1, 2, 3},
		ClearOmittedTuple:   ClearOmittedTuple{FieldA: "value"},
		AInt:                42,
		AString:             "value",
		Adur:                time.Second,
		AJSON:               json.Number(`43.0000000000002`),
	}
	dst := filled
	_, err = dst.UnmarshalMsg(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(dst, cleared) {
		t.Errorf("\n got=%#v\nwant=%#v", dst, cleared)
	}
	// Reset
	dst = filled
	err = dst.DecodeMsg(msgp.NewReader(bytes.NewReader(encoded)))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(dst, cleared) {
		t.Errorf("\n got=%#v\nwant=%#v", dst, cleared)
	}

	// Check that fields aren't accidentally zeroing fields.
	wantJson, err := json.Marshal(filled)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err = filled.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	dst = ClearOmitted0{}
	_, err = dst.UnmarshalMsg(encoded)
	if err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, wantJson) {
		t.Errorf("\n got=%#v\nwant=%#v", string(got), string(wantJson))
	}
	// Reset
	dst = ClearOmitted0{}
	err = dst.DecodeMsg(msgp.NewReader(bytes.NewReader(encoded)))
	if err != nil {
		t.Fatal(err)
	}
	got, err = json.Marshal(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, wantJson) {
		t.Errorf("\n got=%#v\nwant=%#v", string(got), string(wantJson))
	}
	t.Log("OK - got", string(got))
}
