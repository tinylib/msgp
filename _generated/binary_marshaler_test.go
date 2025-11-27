package _generated

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestBinaryMarshalerDirective(t *testing.T) {
	original := &BinaryTestType{Value: "test_data"}

	// Test marshal
	data, err := original.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Test unmarshal
	result := &BinaryTestType{}
	_, err = result.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result.Value != original.Value {
		t.Errorf("Expected %s, got %s", original.Value, result.Value)
	}
}

func TestTextMarshalerBinDirective(t *testing.T) {
	original := &TextBinTestType{Value: "test_data"}

	// Test marshal
	data, err := original.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Test unmarshal
	result := &TextBinTestType{}
	_, err = result.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result.Value != original.Value {
		t.Errorf("Expected %s, got %s", original.Value, result.Value)
	}
}

func TestTextMarshalerStringDirective(t *testing.T) {
	original := &TextStringTestType{Value: "test_data"}

	// Test marshal
	data, err := original.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Test unmarshal
	result := &TextStringTestType{}
	_, err = result.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result.Value != original.Value {
		t.Errorf("Expected %s, got %s", original.Value, result.Value)
	}
}

func TestBinaryMarshalerRoundTrip(t *testing.T) {
	tests := []string{"", "hello", "world with spaces", "unicode: 测试"}

	for _, testVal := range tests {
		original := &BinaryTestType{Value: testVal}

		// Test marshal/unmarshal
		buf, err := original.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("MarshalMsg failed for %q: %v", testVal, err)
		}

		result := &BinaryTestType{}
		_, err = result.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("UnmarshalMsg failed for %q: %v", testVal, err)
		}

		if result.Value != testVal {
			t.Errorf("Round trip failed: expected %q, got %q", testVal, result.Value)
		}
	}
}

func TestTextMarshalerSize(t *testing.T) {
	testType := &TextBinTestType{Value: "test"}

	// Get the size
	size := testType.Msgsize()

	// Marshal and check the actual size matches estimate
	data, err := testType.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("MarshalMsg failed: %v", err)
	}

	if len(data) > size {
		t.Errorf("Msgsize underestimated: estimated %d, actual %d", size, len(data))
	}

	if size > len(data)+100 { // Allow some reasonable overhead
		t.Errorf("Msgsize too conservative: estimated %d, actual %d", size, len(data))
	}
}

func TestStructWithMarshalerFields(t *testing.T) {
	// Create a test struct with all field types populated
	original := &TestStruct{
		// Direct values
		BinaryValue:     BinaryTestType{Value: "binary_val"},
		TextBinValue:    TextBinTestType{Value: "text_bin_val"},
		TextStringValue: TextStringTestType{Value: "text_str_val"},

		// Pointers
		BinaryPtr:     &BinaryTestType{Value: "binary_ptr"},
		TextBinPtr:    &TextBinTestType{Value: "text_bin_ptr"},
		TextStringPtr: &TextStringTestType{Value: "text_str_ptr"},

		// Slices
		BinarySlice: []BinaryTestType{
			{Value: "bin_slice_0"},
			{Value: "bin_slice_1"},
		},
		TextBinSlice: []TextBinTestType{
			{Value: "text_bin_slice_0"},
			{Value: "text_bin_slice_1"},
		},
		TextStringSlice: []TextStringTestType{
			{Value: "text_str_slice_0"},
			{Value: "text_str_slice_1"},
		},

		// Arrays
		BinaryArray: [3]BinaryTestType{
			{Value: "bin_array_0"},
			{Value: "bin_array_1"},
			{Value: "bin_array_2"},
		},
		TextBinArray: [2]TextBinTestType{
			{Value: "text_bin_array_0"},
			{Value: "text_bin_array_1"},
		},
		TextStringArray: [4]TextStringTestType{
			{Value: "text_str_array_0"},
			{Value: "text_str_array_1"},
			{Value: "text_str_array_2"},
			{Value: "text_str_array_3"},
		},

		// Maps
		BinaryMap: map[string]BinaryTestType{
			"key1": {Value: "bin_map_val1"},
			"key2": {Value: "bin_map_val2"},
		},
		TextBinMap: map[string]TextBinTestType{
			"key1": {Value: "text_bin_map_val1"},
			"key2": {Value: "text_bin_map_val2"},
		},
		TextStringMap: map[string]TextStringTestType{
			"key1": {Value: "text_str_map_val1"},
			"key2": {Value: "text_str_map_val2"},
		},

		// Nested types
		NestedPtrSlice: []*BinaryTestType{
			{Value: "nested_ptr_0"},
			{Value: "nested_ptr_1"},
		},
		SliceOfArrays: [][2]TextBinTestType{
			{{Value: "slice_arr_0_0"}, {Value: "slice_arr_0_1"}},
			{{Value: "slice_arr_1_0"}, {Value: "slice_arr_1_1"}},
		},
		MapOfSlices: map[string][]BinaryTestType{
			"slice1": {{Value: "map_slice_val_0"}, {Value: "map_slice_val_1"}},
			"slice2": {{Value: "map_slice_val_2"}},
		},
	}

	// Test marshal/unmarshal
	data, err := original.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("MarshalMsg failed: %v", err)
	}

	result := &TestStruct{}
	_, err = result.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("UnmarshalMsg failed: %v", err)
	}

	// Verify direct values
	if result.BinaryValue.Value != original.BinaryValue.Value {
		t.Errorf("BinaryValue mismatch: expected %s, got %s", original.BinaryValue.Value, result.BinaryValue.Value)
	}
	if result.TextBinValue.Value != original.TextBinValue.Value {
		t.Errorf("TextBinValue mismatch: expected %s, got %s", original.TextBinValue.Value, result.TextBinValue.Value)
	}
	if result.TextStringValue.Value != original.TextStringValue.Value {
		t.Errorf("TextStringValue mismatch: expected %s, got %s", original.TextStringValue.Value, result.TextStringValue.Value)
	}

	// Verify pointers
	if result.BinaryPtr == nil || result.BinaryPtr.Value != original.BinaryPtr.Value {
		t.Errorf("BinaryPtr mismatch")
	}
	if result.TextBinPtr == nil || result.TextBinPtr.Value != original.TextBinPtr.Value {
		t.Errorf("TextBinPtr mismatch")
	}
	if result.TextStringPtr == nil || result.TextStringPtr.Value != original.TextStringPtr.Value {
		t.Errorf("TextStringPtr mismatch")
	}

	// Verify slices
	if len(result.BinarySlice) != len(original.BinarySlice) {
		t.Errorf("BinarySlice length mismatch: expected %d, got %d", len(original.BinarySlice), len(result.BinarySlice))
	} else {
		for i := range result.BinarySlice {
			if result.BinarySlice[i].Value != original.BinarySlice[i].Value {
				t.Errorf("BinarySlice[%d] mismatch: expected %s, got %s", i, original.BinarySlice[i].Value, result.BinarySlice[i].Value)
			}
		}
	}

	// Verify arrays
	for i := range result.BinaryArray {
		if result.BinaryArray[i].Value != original.BinaryArray[i].Value {
			t.Errorf("BinaryArray[%d] mismatch: expected %s, got %s", i, original.BinaryArray[i].Value, result.BinaryArray[i].Value)
		}
	}

	// Verify maps
	if len(result.BinaryMap) != len(original.BinaryMap) {
		t.Errorf("BinaryMap length mismatch: expected %d, got %d", len(original.BinaryMap), len(result.BinaryMap))
	} else {
		for k, v := range original.BinaryMap {
			if resultVal, exists := result.BinaryMap[k]; !exists || resultVal.Value != v.Value {
				t.Errorf("BinaryMap[%s] mismatch: expected %s, got %s", k, v.Value, resultVal.Value)
			}
		}
	}

	// Verify nested structures
	if len(result.NestedPtrSlice) != len(original.NestedPtrSlice) {
		t.Errorf("NestedPtrSlice length mismatch")
	} else {
		for i := range result.NestedPtrSlice {
			if result.NestedPtrSlice[i] == nil || result.NestedPtrSlice[i].Value != original.NestedPtrSlice[i].Value {
				t.Errorf("NestedPtrSlice[%d] mismatch", i)
			}
		}
	}
}

func TestStructWithOmitEmptyFields(t *testing.T) {
	// Test struct with nil pointers (should use omitempty/omitzero)
	original := &TestStruct{
		BinaryValue:     BinaryTestType{Value: "present"},
		TextBinValue:    TextBinTestType{Value: "also_present"},
		TextStringValue: TextStringTestType{Value: "string_present"},
		// Leave pointers nil
		BinaryPtr:     nil,
		TextBinPtr:    nil, // this has omitempty tag
		TextStringPtr: nil, // this has omitzero tag
		// Empty slices and maps
		BinarySlice: []BinaryTestType{},
		BinaryMap:   map[string]BinaryTestType{},
	}

	// Test marshal/unmarshal
	data, err := original.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("MarshalMsg failed: %v", err)
	}

	result := &TestStruct{}
	_, err = result.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("UnmarshalMsg failed: %v", err)
	}

	// Verify values are preserved and nils are handled correctly
	if result.BinaryValue.Value != original.BinaryValue.Value {
		t.Errorf("BinaryValue mismatch")
	}
	if result.BinaryPtr != nil {
		t.Errorf("Expected BinaryPtr to be nil")
	}
	if result.TextBinPtr != nil {
		t.Errorf("Expected TextBinPtr to be nil")
	}
	if result.TextStringPtr != nil {
		t.Errorf("Expected TextStringPtr to be nil")
	}
}

func TestStructSizeEstimation(t *testing.T) {
	testStruct := &TestStruct{
		BinaryValue: BinaryTestType{Value: "test"},
		BinarySlice: []BinaryTestType{
			{Value: "slice1"},
			{Value: "slice2"},
		},
		BinaryMap: map[string]BinaryTestType{
			"key": {Value: "value"},
		},
	}

	// Get size estimate
	size := testStruct.Msgsize()

	// Marshal and verify size estimate
	data, err := testStruct.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("MarshalMsg failed: %v", err)
	}

	if len(data) > size {
		t.Errorf("Msgsize underestimated: estimated %d, actual %d", size, len(data))
	}
}

func TestMarshalerErrorHandling(t *testing.T) {
	// Test marshaling with marshaler that can fail
	errorType := &ErrorTestType{ShouldError: true}

	_, err := errorType.MarshalMsg(nil)
	if err == nil {
		t.Error("Expected marshal error but got none")
	}

	// Test in struct context
	testStruct := &TestStruct{
		BinaryValue: BinaryTestType{Value: "good"},
	}

	// This should succeed
	_, err = testStruct.MarshalMsg(nil)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

func TestTextMarshalerAsStringPositioning(t *testing.T) {
	// Test that as:string works when placed in middle of type list
	middle := &TestTextMarshalerStringMiddle{Value: "test_middle"}

	// Marshal
	data, err := middle.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var decoded TestTextMarshalerStringMiddle
	_, err = decoded.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Value != "middle:test_middle" {
		t.Errorf("Expected 'middle:test_middle', got '%s'", decoded.Value)
	}

	// Test that as:string works when placed at end of type list
	end := &TestTextMarshalerStringEnd{Value: "test_end"}

	// Marshal
	data, err = end.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var decodedEnd TestTextMarshalerStringEnd
	_, err = decodedEnd.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decodedEnd.Value != "end:test_end" {
		t.Errorf("Expected 'end:test_end', got '%s'", decodedEnd.Value)
	}
}

func TestTextAppenderAsStringPositioning(t *testing.T) {
	// Test that as:string works with textappend directive
	appender := &TestTextAppenderStringPos{Value: "test_append"}

	// Marshal
	data, err := appender.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var decoded TestTextAppenderStringPos
	_, err = decoded.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Value != "append:test_append" {
		t.Errorf("Expected 'append:test_append', got '%s'", decoded.Value)
	}
}

func TestBinaryAppenderDirective(t *testing.T) {
	for _, size := range []int{10, 100, 1000, 100000} {

		// Test BinaryAppender interface
		appender := &BinaryAppenderType{Value: strings.Repeat("a", size)}

		// Test round trip
		data, err := appender.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var decoded BinaryAppenderType
		_, err = decoded.UnmarshalMsg(data)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if want := "binappend:" + appender.Value; decoded.Value != want {
			t.Errorf("Expected '%s', got '%s'", want, decoded.Value)
		}

		// Test encode/decode
		var buf bytes.Buffer
		writer := msgp.NewWriter(&buf)
		err = appender.EncodeMsg(writer)
		if err != nil {
			t.Fatalf("Failed to encode: %v", err)
		}
		writer.Flush()

		reader := msgp.NewReader(&buf)
		var decoded2 BinaryAppenderType
		err = decoded2.DecodeMsg(reader)
		if err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}

		if want := "binappend:" + appender.Value; decoded2.Value != want {
			t.Errorf("Expected '%s', got '%s'", want, decoded2.Value)
		}
	}
}

func TestTextAppenderBinDirective(t *testing.T) {
	// Test TextAppender interface (stored as binary)
	appender := &TextAppenderBinType{Value: "test_text_bin"}

	// Test round trip
	data, err := appender.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded TextAppenderBinType
	_, err = decoded.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Value != "textbin:test_text_bin" {
		t.Errorf("Expected 'textbin:test_text_bin', got '%s'", decoded.Value)
	}

	// Test encode/decode
	var buf bytes.Buffer
	writer := msgp.NewWriter(&buf)
	err = appender.EncodeMsg(writer)
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}
	writer.Flush()

	reader := msgp.NewReader(&buf)
	var decoded2 TextAppenderBinType
	err = decoded2.DecodeMsg(reader)
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	if decoded2.Value != "textbin:test_text_bin" {
		t.Errorf("Expected 'textbin:test_text_bin', got '%s'", decoded2.Value)
	}
}

func TestBinaryAppenderErrorHandling(t *testing.T) {
	// Test error handling with BinaryAppender
	errorType := &ErrorBinaryAppenderType{ShouldError: true}

	_, err := errorType.MarshalMsg(nil)
	if err == nil {
		t.Error("Expected marshal error but got none")
	}

	// Test successful case
	successType := &ErrorBinaryAppenderType{ShouldError: false}
	data, err := successType.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	var decoded ErrorBinaryAppenderType
	_, err = decoded.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Value != "ok" {
		t.Errorf("Expected 'ok', got '%s'", decoded.Value)
	}
}

func TestTextAppenderErrorHandling(t *testing.T) {
	// Test error handling with TextAppender
	errorType := &ErrorTextAppenderType{ShouldError: true}

	_, err := errorType.MarshalMsg(nil)
	if err == nil {
		t.Error("Expected marshal error but got none")
	}

	// Test successful case
	successType := &ErrorTextAppenderType{ShouldError: false}
	data, err := successType.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	var decoded ErrorTextAppenderType
	_, err = decoded.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Value != "ok" {
		t.Errorf("Expected 'ok', got '%s'", decoded.Value)
	}
}

func TestAppenderTypesSizeEstimation(t *testing.T) {
	// Test that size estimation works for appender types
	binaryAppender := &BinaryAppenderType{Value: "size_test"}
	size := binaryAppender.Msgsize()

	data, err := binaryAppender.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	if len(data) > size {
		t.Errorf("Binary appender size underestimated: estimated %d, actual %d", size, len(data))
	}

	textAppender := &TextAppenderBinType{Value: "size_test"}
	size = textAppender.Msgsize()

	data, err = textAppender.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	if len(data) > size {
		t.Errorf("Text appender size underestimated: estimated %d, actual %d", size, len(data))
	}
}

func TestBinaryAppenderValueRoundTrip(t *testing.T) {
	tests := []string{"", "hello", "world with spaces", "unicode: 测试"}

	for _, testVal := range tests {
		original := BinaryAppenderValue{Value: testVal}

		buf, err := original.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("MarshalMsg failed for %q: %v", testVal, err)
		}

		var result BinaryAppenderValue
		_, err = result.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("UnmarshalMsg failed for %q: %v", testVal, err)
		}

		expected := "binappend:" + testVal
		if result.Value != expected {
			t.Errorf("Round trip failed: expected %q, got %q", expected, result.Value)
		}
	}
}

func TestTestAppendTextStringRoundTrip(t *testing.T) {
	tests := []string{"", "hello", "world with spaces", "unicode: 测试"}

	for _, testVal := range tests {
		original := TestAppendTextString{Value: testVal}

		buf, err := original.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("MarshalMsg failed for %q: %v", testVal, err)
		}

		var result TestAppendTextString
		_, err = result.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("UnmarshalMsg failed for %q: %v", testVal, err)
		}

		expected := "append:" + testVal
		if result.Value != expected {
			t.Errorf("Round trip failed: expected %q, got %q", expected, result.Value)
		}
	}
}

func TestTextAppenderBinValueRoundTrip(t *testing.T) {
	tests := []string{"", "hello", "world with spaces", "unicode: 测试"}

	for _, testVal := range tests {
		original := TextAppenderBinValue{Value: testVal}

		buf, err := original.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("MarshalMsg failed for %q: %v", testVal, err)
		}

		var result TextAppenderBinValue
		_, err = result.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("UnmarshalMsg failed for %q: %v", testVal, err)
		}

		expected := "textbin:" + testVal
		if result.Value != expected {
			t.Errorf("Round trip failed: expected %q, got %q", expected, result.Value)
		}
	}
}
