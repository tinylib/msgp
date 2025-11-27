package _generated

import (
	"encoding"
	"fmt"
)

//go:generate msgp -v

// BinaryTestType implements encoding.BinaryMarshaler and encoding.BinaryUnmarshaler
type BinaryTestType struct {
	Value string
}

func (t *BinaryTestType) MarshalBinary() ([]byte, error) {
	return []byte(t.Value), nil
}

func (t *BinaryTestType) UnmarshalBinary(data []byte) error {
	t.Value = string(data)
	return nil
}

// Verify it implements the interfaces
var _ encoding.BinaryMarshaler = (*BinaryTestType)(nil)
var _ encoding.BinaryUnmarshaler = (*BinaryTestType)(nil)

//msgp:binmarshal BinaryTestType

// TextBinTestType implements encoding.TextMarshaler and encoding.TextUnmarshaler
type TextBinTestType struct {
	Value string
}

func (t *TextBinTestType) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("text:%s", t.Value)), nil
}

func (t *TextBinTestType) UnmarshalText(data []byte) error {
	t.Value = string(data[5:]) // Remove "text:" prefix
	return nil
}

var _ encoding.TextMarshaler = (*TextBinTestType)(nil)
var _ encoding.TextUnmarshaler = (*TextBinTestType)(nil)

//msgp:textmarshal TextBinTestType

// TextStringTestType for testing as:string option
type TextStringTestType struct {
	Value string
}

func (t *TextStringTestType) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("stringtext:%s", t.Value)), nil
}

func (t *TextStringTestType) UnmarshalText(data []byte) error {
	t.Value = string(data[11:]) // Remove "stringtext:" prefix
	return nil
}

var _ encoding.TextMarshaler = (*TextStringTestType)(nil)
var _ encoding.TextUnmarshaler = (*TextStringTestType)(nil)

//msgp:textmarshal as:string TextStringTestType

// TestStruct contains various combinations of marshaler types
type TestStruct struct {
	// Direct values
	BinaryValue     BinaryTestType     `msg:"bin_val"`
	TextBinValue    TextBinTestType    `msg:"text_bin_val"`
	TextStringValue TextStringTestType `msg:"text_str_val"`

	// Pointers
	BinaryPtr     *BinaryTestType     `msg:"bin_ptr"`
	TextBinPtr    *TextBinTestType    `msg:"text_bin_ptr,omitempty"`
	TextStringPtr *TextStringTestType `msg:"text_str_ptr,omitempty"`

	// Slices
	BinarySlice     []BinaryTestType     `msg:"bin_slice"`
	TextBinSlice    []TextBinTestType    `msg:"text_bin_slice"`
	TextStringSlice []TextStringTestType `msg:"text_str_slice"`

	// Arrays
	BinaryArray     [3]BinaryTestType     `msg:"bin_array"`
	TextBinArray    [2]TextBinTestType    `msg:"text_bin_array"`
	TextStringArray [4]TextStringTestType `msg:"text_str_array"`

	// Maps with marshaler types as values
	BinaryMap     map[string]BinaryTestType     `msg:"bin_map"`
	TextBinMap    map[string]TextBinTestType    `msg:"text_bin_map"`
	TextStringMap map[string]TextStringTestType `msg:"text_str_map"`

	// Nested pointers and slices
	NestedPtrSlice []*BinaryTestType           `msg:"nested_ptr_slice"`
	SliceOfArrays  [][2]TextBinTestType        `msg:"slice_of_arrays"`
	MapOfSlices    map[string][]BinaryTestType `msg:"map_of_slices"`
}

//msgp:binmarshal ErrorTestType

// ErrorTestType for testing error conditions
type ErrorTestType struct {
	ShouldError bool
}

func (e *ErrorTestType) MarshalBinary() ([]byte, error) {
	if e.ShouldError {
		return nil, fmt.Errorf("intentional marshal error")
	}
	return []byte("ok"), nil
}

func (e *ErrorTestType) UnmarshalBinary(data []byte) error {
	if string(data) == "error" {
		return fmt.Errorf("intentional unmarshal error")
	}
	e.ShouldError = false
	return nil
}

// Test types for as:string positioning flexibility

// TestTextMarshalerStringMiddle for testing as:string in middle of type list
type TestTextMarshalerStringMiddle struct {
	Value string
}

func (t *TestTextMarshalerStringMiddle) MarshalText() ([]byte, error) {
	return []byte("middle:" + t.Value), nil
}

func (t *TestTextMarshalerStringMiddle) UnmarshalText(text []byte) error {
	t.Value = string(text)
	return nil
}

// TestTextMarshalerStringEnd for testing as:string at end of type list
type TestTextMarshalerStringEnd struct {
	Value string
}

func (t *TestTextMarshalerStringEnd) MarshalText() ([]byte, error) {
	return []byte("end:" + t.Value), nil
}

func (t *TestTextMarshalerStringEnd) UnmarshalText(text []byte) error {
	t.Value = string(text)
	return nil
}

// TestTextAppenderStringPos for testing textappend with as:string positioning
type TestTextAppenderStringPos struct {
	Value string
}

func (t *TestTextAppenderStringPos) AppendText(dst []byte) ([]byte, error) {
	return append(dst, []byte("append:"+t.Value)...), nil
}

func (t *TestTextAppenderStringPos) UnmarshalText(text []byte) error {
	t.Value = string(text)
	return nil
}

// BinaryAppenderType implements encoding.BinaryAppender (Go 1.22+)
type BinaryAppenderType struct {
	Value string
}

func (t *BinaryAppenderType) AppendBinary(dst []byte) ([]byte, error) {
	return append(dst, []byte("binappend:"+t.Value)...), nil
}

func (t *BinaryAppenderType) UnmarshalBinary(data []byte) error {
	t.Value = string(data)
	return nil
}

// TextAppenderBinType implements encoding.TextAppender (stored as binary)
type TextAppenderBinType struct {
	Value string
}

func (t *TextAppenderBinType) AppendText(dst []byte) ([]byte, error) {
	return append(dst, []byte("textbin:"+t.Value)...), nil
}

func (t *TextAppenderBinType) UnmarshalText(text []byte) error {
	t.Value = string(text)
	return nil
}

// ErrorBinaryAppenderType for testing error conditions with BinaryAppender
type ErrorBinaryAppenderType struct {
	ShouldError bool
	Value       string
}

func (e *ErrorBinaryAppenderType) AppendBinary(dst []byte) ([]byte, error) {
	if e.ShouldError {
		return nil, fmt.Errorf("intentional append binary error")
	}
	return append(dst, []byte("ok")...), nil
}

func (e *ErrorBinaryAppenderType) UnmarshalBinary(data []byte) error {
	if string(data) == "error" {
		return fmt.Errorf("intentional unmarshal binary error")
	}
	e.ShouldError = false
	e.Value = string(data)
	return nil
}

// ErrorTextAppenderType for testing error conditions with TextAppender
type ErrorTextAppenderType struct {
	ShouldError bool
	Value       string
}

func (e *ErrorTextAppenderType) AppendText(dst []byte) ([]byte, error) {
	if e.ShouldError {
		return nil, fmt.Errorf("intentional append text error")
	}
	return append(dst, []byte("ok")...), nil
}

func (e *ErrorTextAppenderType) UnmarshalText(text []byte) error {
	if string(text) == "error" {
		return fmt.Errorf("intentional unmarshal text error")
	}
	e.ShouldError = false
	e.Value = string(text)
	return nil
}

//msgp:binappend BinaryAppenderType ErrorBinaryAppenderType
//msgp:textappend TextAppenderBinType ErrorTextAppenderType
//msgp:textmarshal TestTextMarshalerStringMiddle as:string TestTextMarshalerStringEnd
//msgp:textappend TestTextAppenderStringPos as:string

//msgp:binappend BinaryAppenderValue

// BinaryAppenderValue implements encoding.BinaryAppender (Go 1.22+)
type BinaryAppenderValue struct {
	Value string `msg:"-"`
}

func (t BinaryAppenderValue) AppendBinary(dst []byte) ([]byte, error) {
	return append(dst, []byte("binappend:"+t.Value)...), nil
}

func (t *BinaryAppenderValue) UnmarshalBinary(data []byte) error {
	t.Value = string(data)
	return nil
}

//msgp:textappend TestAppendTextString as:string

type TestAppendTextString struct {
	Value string `msg:"-"`
}

func (t TestAppendTextString) AppendText(dst []byte) ([]byte, error) {
	return append(dst, []byte("append:"+t.Value)...), nil
}

func (t *TestAppendTextString) UnmarshalText(text []byte) error {
	t.Value = string(text)
	return nil
}

//msgp:textappend TextAppenderBinValue

// TextAppenderBinValue implements encoding.TextAppender (stored as binary)
type TextAppenderBinValue struct {
	Value string `msg:"-"`
}

func (t TextAppenderBinValue) AppendText(dst []byte) ([]byte, error) {
	return append(dst, []byte("textbin:"+t.Value)...), nil
}

func (t *TextAppenderBinValue) UnmarshalText(text []byte) error {
	t.Value = string(text)
	return nil
}
