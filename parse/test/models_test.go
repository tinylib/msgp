package test

import (
	"fmt"
	"testing"

	"github.com/tinylib/msgp/parse"
)

func TestParseType(t *testing.T) {
	f, err := parse.File("models.go", false)
	if err != nil {
		panic(err)
	}
	for k, v := range f.Aliases {
		if k != "AliasType" || fmt.Sprintf("%v", v) != "string" {
			t.Errorf("Incorrect alias %s: %s", k, v)
		}
	}
}

func TestStructParse(t *testing.T) {
	s := &TestStruct{
		Alias: "alias",
		Str:   "string",
	}
	data, err := s.MarshalMsg(nil)
	if err != nil {
		t.Error(err)
		return
	}
	d := &TestStruct{}
	_, err = d.UnmarshalMsg(data)
	if err != nil {
		t.Error(err)
		return
	}
	if d.Str != "string" {
		t.Errorf("d.Str should be 'string' but actual value is '%s'", d.Str)
	}
	if d.Alias != "alias" {
		t.Errorf("d.Alias should be 'alias' but actual value is '%s'", d.Alias)
	}
}
