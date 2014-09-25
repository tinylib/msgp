package _generated

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

// this will work if we compile
func TestBuild(t *testing.T) {}

func BenchmarkMarshal(b *testing.B) {
	v := &TestBench{
		Name:     "A-Name",
		BirthDay: time.Now(),
		Phone:    "1015559098",
		Siblings: 5,
		Spouse:   false,
		Money:    10982.0,
		Tags: map[string]string{
			"tag_a": "a_tag_value",
		},
		Aliases: []string{
			"Another_Alias",
			"Yet_Another_Alias",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Marshal()
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	v := &TestBench{
		Name:     "A-Name",
		BirthDay: time.Now(),
		Phone:    "1015559098",
		Siblings: 5,
		Spouse:   false,
		Money:    10982.0,
		Tags: map[string]string{
			"tag_a": "a_tag_value",
		},
		Aliases: []string{
			"Another_Alias",
			"Yet_Another_Alias",
		},
	}

	var buf bytes.Buffer
	v.WriteTo(&buf)
	bts := buf.Bytes()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Unmarshal(bts)
	}

}

func Test1EncodeDecode(t *testing.T) {
	f := 32.00
	tt := &TestType{
		F: &f,
		Els: map[string]string{
			"thing_one": "one",
			"thing_two": "two",
		},
		Obj: struct {
			ValueA string `msg:"value_a"`
			ValueB []byte `msg:"value_b"`
		}{
			ValueA: "here's the first inner value",
			ValueB: []byte("here's the second inner value"),
		},
		Child: nil,
		Time:  time.Now(),
	}

	var buf bytes.Buffer

	_, err := tt.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}

	tnew := new(TestType)

	_, err = tnew.ReadFrom(&buf)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(tt, tnew) {
		t.Logf("in: %v", tt)
		t.Logf("out: %v", tnew)
		t.Fatal("objects not equal")
	}
}
