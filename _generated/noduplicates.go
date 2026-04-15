//msgp:noduplicates

package _generated

//go:generate msgp

// NoDupSimple is a basic struct with a few fields.
type NoDupSimple struct {
	Foo string `msg:"Foo"`
	Bar int    `msg:"Bar"`
	Baz bool   `msg:"Baz"`
}

// NoDupMap is a standalone map type.
type NoDupMap map[string]string

// NoDupNested has a struct with a map field.
type NoDupNested struct {
	Name  string            `msg:"Name"`
	Items map[string]int    `msg:"Items"`
	Tags  map[string]string `msg:"Tags"`
}

// NoDupManyFields exercises the bitmask across byte boundaries (>8 fields).
type NoDupManyFields struct {
	F01 string `msg:"f01"`
	F02 string `msg:"f02"`
	F03 string `msg:"f03"`
	F04 string `msg:"f04"`
	F05 string `msg:"f05"`
	F06 string `msg:"f06"`
	F07 string `msg:"f07"`
	F08 string `msg:"f08"`
	F09 string `msg:"f09"`
	F10 string `msg:"f10"`
}

// NoDupMapInt is a map with int values.
type NoDupMapInt map[string]int

// NoDupRichTypes covers pointer, slice, bytes, float, and nested struct fields.
type NoDupRichTypes struct {
	Str      string           `msg:"str"`
	Num      float64          `msg:"num"`
	Data     []byte           `msg:"data"`
	PtrStr   *string          `msg:"ptr_str"`
	Strings  []string         `msg:"strings"`
	Inner    NoDupInner       `msg:"inner"`
	PtrInner *NoDupInner      `msg:"ptr_inner"`
	MapSlice map[string][]int `msg:"map_slice"`
}

// NoDupInner is used as a nested struct field.
type NoDupInner struct {
	X int    `msg:"x"`
	Y string `msg:"y"`
}

// NoDupOmitempty has omitempty fields mixed with regular fields.
type NoDupOmitempty struct {
	Required string `msg:"required"`
	Optional string `msg:"optional,omitempty"`
	Flag     bool   `msg:"flag,omitempty"`
	Count    int    `msg:"count"`
}

// NoDupAllownil has allownil fields.
type NoDupAllownil struct {
	Name    string            `msg:"name"`
	Vals    []byte            `msg:"vals,allownil"`
	Items   []int             `msg:"items,allownil"`
	Mapping map[string]string `msg:"mapping,allownil"`
}

// NoDupMapComplex is a map with complex values.
type NoDupMapComplex map[string]NoDupInner

// NoDup17Fields tests 17 fields (3 bytes in bitmask).
type NoDup17Fields struct {
	A01 int `msg:"a01"`
	A02 int `msg:"a02"`
	A03 int `msg:"a03"`
	A04 int `msg:"a04"`
	A05 int `msg:"a05"`
	A06 int `msg:"a06"`
	A07 int `msg:"a07"`
	A08 int `msg:"a08"`
	A09 int `msg:"a09"`
	A10 int `msg:"a10"`
	A11 int `msg:"a11"`
	A12 int `msg:"a12"`
	A13 int `msg:"a13"`
	A14 int `msg:"a14"`
	A15 int `msg:"a15"`
	A16 int `msg:"a16"`
	A17 int `msg:"a17"`
}
