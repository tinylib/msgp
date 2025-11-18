package _generated

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestSliceLimitEnforcement(t *testing.T) {
	data := UnlimitedData{}

	// Test slice limit with DecodeMsg (using big_slice which is []string)
	t.Run("DecodeMsg_SliceLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "big_slice")
		buf = msgp.AppendArrayHeader(buf, 150) // Exceeds limit of 100

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded, got %v", err)
		}
	})

	// Test slice limit with UnmarshalMsg
	t.Run("UnmarshalMsg_SliceLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "big_slice")
		buf = msgp.AppendArrayHeader(buf, 150) // Exceeds limit of 100

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded, got %v", err)
		}
	})

	// Test that slices within limit work fine
	t.Run("SliceWithinLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "big_slice")
		buf = msgp.AppendArrayHeader(buf, 50) // Within limit
		for i := 0; i < 50; i++ {
			buf = msgp.AppendString(buf, "test")
		}

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Errorf("Unexpected error for slice within limit: %v", err)
		}
	})
}

func TestMapLimitEnforcement(t *testing.T) {
	data := LimitTestData{}

	// Test map limit with DecodeMsg
	t.Run("DecodeMsg_MapLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "small_map")
		buf = msgp.AppendMapHeader(buf, 60) // Exceeds limit of 50

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded, got %v", err)
		}
	})

	// Test map limit with UnmarshalMsg
	t.Run("UnmarshalMsg_MapLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "small_map")
		buf = msgp.AppendMapHeader(buf, 60) // Exceeds limit of 50

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded, got %v", err)
		}
	})

	// Test that maps within limit work fine
	t.Run("MapWithinLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "small_map")
		buf = msgp.AppendMapHeader(buf, 3) // Within limit
		buf = msgp.AppendString(buf, "a")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "b")
		buf = msgp.AppendInt(buf, 2)
		buf = msgp.AppendString(buf, "c")
		buf = msgp.AppendInt(buf, 3)

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Errorf("Unexpected error for map within limit: %v", err)
		}
	})
}

func TestFixedArraysNotLimited(t *testing.T) {
	// Test that fixed arrays are not subject to limits
	// BigArray [1000]int should work even though 1000 > 100 (array limit)
	data := UnlimitedData{}

	t.Run("FixedArray_DecodeMsg", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "big_array")
		buf = msgp.AppendArrayHeader(buf, 1000) // Fixed array size, should not be limited
		for i := 0; i < 1000; i++ {
			buf = msgp.AppendInt(buf, i)
		}

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != nil {
			t.Errorf("Fixed arrays should not be limited, got error: %v", err)
		}
	})

	t.Run("FixedArray_UnmarshalMsg", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "big_array")
		buf = msgp.AppendArrayHeader(buf, 1000) // Fixed array size, should not be limited
		for i := 0; i < 1000; i++ {
			buf = msgp.AppendInt(buf, i)
		}

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Errorf("Fixed arrays should not be limited, got error: %v", err)
		}
	})
}

func TestSliceLimitsApplied(t *testing.T) {
	// Test that dynamic slices are subject to limits
	data := UnlimitedData{}

	t.Run("Slice_ExceedsLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "big_slice")
		buf = msgp.AppendArrayHeader(buf, 150) // Exceeds array limit of 100

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for slice, got %v", err)
		}
	})

	t.Run("Slice_WithinLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "big_slice")
		buf = msgp.AppendArrayHeader(buf, 50) // Within array limit of 100
		for i := 0; i < 50; i++ {
			buf = msgp.AppendString(buf, "test")
		}

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Errorf("Unexpected error for slice within limit: %v", err)
		}
	})
}

func TestNestedArrayLimits(t *testing.T) {
	// Test limits on nested arrays within maps
	data := UnlimitedData{}

	t.Run("NestedArray_ExceedsLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "big_map")
		buf = msgp.AppendMapHeader(buf, 1) // Within map limit
		buf = msgp.AppendString(buf, "key")
		buf = msgp.AppendArrayHeader(buf, 150) // Nested array exceeds limit of 100

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for nested array, got %v", err)
		}
	})

	t.Run("NestedArray_WithinLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "big_map")
		buf = msgp.AppendMapHeader(buf, 1) // Within map limit
		buf = msgp.AppendString(buf, "key")
		buf = msgp.AppendArrayHeader(buf, 50) // Nested array within limit
		for i := 0; i < 50; i++ {
			buf = msgp.AppendInt(buf, i)
		}

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Errorf("Unexpected error for nested array within limit: %v", err)
		}
	})
}

func TestMapExceedsLimit(t *testing.T) {
	data := UnlimitedData{}

	t.Run("Map_ExceedsLimit", func(t *testing.T) {
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "big_map")
		buf = msgp.AppendMapHeader(buf, 60) // Exceeds map limit of 50

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for map, got %v", err)
		}
	})
}

func TestStructLevelLimits(t *testing.T) {
	// Test that the struct-level map limits are enforced
	data := LimitTestData{}

	t.Run("StructMap_ExceedsLimit", func(t *testing.T) {
		// Create a struct with too many fields
		buf := msgp.AppendMapHeader(nil, 60) // Exceeds map limit of 50

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for struct map, got %v", err)
		}
	})
}

func TestNormalOperationWithinLimits(t *testing.T) {
	// Test that normal operation works when everything is within limits
	data := LimitTestData{}

	// Create valid data
	data.SmallArray = [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	data.LargeSlice = []byte("test data")
	data.SmallMap = map[string]int{"a": 1, "b": 2, "c": 3}

	t.Run("RoundTrip_Marshal_Unmarshal", func(t *testing.T) {
		// Test MarshalMsg -> UnmarshalMsg
		buf, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("MarshalMsg failed: %v", err)
		}

		var result LimitTestData
		_, err = result.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("UnmarshalMsg failed: %v", err)
		}

		// Verify data integrity
		if result.SmallArray != data.SmallArray {
			t.Errorf("SmallArray mismatch: got %v, want %v", result.SmallArray, data.SmallArray)
		}
		if !bytes.Equal(result.LargeSlice, data.LargeSlice) {
			t.Errorf("LargeSlice mismatch: got %v, want %v", result.LargeSlice, data.LargeSlice)
		}
		if len(result.SmallMap) != len(data.SmallMap) {
			t.Errorf("SmallMap length mismatch: got %d, want %d", len(result.SmallMap), len(data.SmallMap))
		}
	})

	t.Run("RoundTrip_Encode_Decode", func(t *testing.T) {
		// Test EncodeMsg -> DecodeMsg
		var buf bytes.Buffer
		writer := msgp.NewWriter(&buf)
		err := data.EncodeMsg(writer)
		if err != nil {
			t.Fatalf("EncodeMsg failed: %v", err)
		}
		writer.Flush()

		var result LimitTestData
		reader := msgp.NewReader(&buf)
		err = result.DecodeMsg(reader)
		if err != nil {
			t.Fatalf("DecodeMsg failed: %v", err)
		}

		// Verify data integrity
		if result.SmallArray != data.SmallArray {
			t.Errorf("SmallArray mismatch: got %v, want %v", result.SmallArray, data.SmallArray)
		}
		if !bytes.Equal(result.LargeSlice, data.LargeSlice) {
			t.Errorf("LargeSlice mismatch: got %v, want %v", result.LargeSlice, data.LargeSlice)
		}
		if len(result.SmallMap) != len(data.SmallMap) {
			t.Errorf("SmallMap length mismatch: got %d, want %d", len(result.SmallMap), len(data.SmallMap))
		}
	})
}

func TestMarshalLimitEnforcement(t *testing.T) {
	// Test marshal-time limit enforcement with MarshalLimitTestData
	// This struct has marshal:true with arrays:30 maps:20

	t.Run("MarshalMsg_SliceLimit", func(t *testing.T) {
		data := MarshalLimitTestData{
			TestSlice: make([]string, 40), // Exceeds array limit of 30
		}
		// Fill the slice
		for i := range data.TestSlice {
			data.TestSlice[i] = "test"
		}

		_, err := data.MarshalMsg(nil)
		if err == nil {
			t.Error("Expected error for slice exceeding marshal limit, got nil")
		}
	})

	t.Run("MarshalMsg_MapLimit", func(t *testing.T) {
		data := MarshalLimitTestData{
			TestMap: make(map[string]int, 25), // Exceeds map limit of 20
		}
		// Fill the map
		for i := 0; i < 25; i++ {
			data.TestMap[fmt.Sprintf("key%d", i)] = i
		}

		_, err := data.MarshalMsg(nil)
		if err == nil {
			t.Error("Expected error for map exceeding marshal limit, got nil")
		}
	})

	t.Run("EncodeMsg_SliceLimit", func(t *testing.T) {
		data := MarshalLimitTestData{
			TestSlice: make([]string, 40), // Exceeds array limit of 30
		}
		// Fill the slice
		for i := range data.TestSlice {
			data.TestSlice[i] = "test"
		}

		var buf bytes.Buffer
		writer := msgp.NewWriter(&buf)
		err := data.EncodeMsg(writer)
		if err == nil {
			t.Error("Expected error for slice exceeding marshal limit, got nil")
		}
	})

	t.Run("EncodeMsg_MapLimit", func(t *testing.T) {
		data := MarshalLimitTestData{
			TestMap: make(map[string]int, 25), // Exceeds map limit of 20
		}
		// Fill the map
		for i := 0; i < 25; i++ {
			data.TestMap[fmt.Sprintf("key%d", i)] = i
		}

		var buf bytes.Buffer
		writer := msgp.NewWriter(&buf)
		err := data.EncodeMsg(writer)
		if err == nil {
			t.Error("Expected error for map exceeding marshal limit, got nil")
		}
	})

	t.Run("MarshalWithinLimits", func(t *testing.T) {
		data := MarshalLimitTestData{
			SmallArray: [5]int{1, 2, 3, 4, 5},
			TestSlice:  []string{"a", "b", "c"},        // Within limit of 30
			TestMap:    map[string]int{"x": 1, "y": 2}, // Within limit of 20
		}

		// Test MarshalMsg
		_, err := data.MarshalMsg(nil)
		if err != nil {
			t.Errorf("Unexpected error for data within marshal limits: %v", err)
		}

		// Test EncodeMsg
		var buf bytes.Buffer
		writer := msgp.NewWriter(&buf)
		err = data.EncodeMsg(writer)
		if err != nil {
			t.Errorf("Unexpected error for data within marshal limits: %v", err)
		}
	})

	t.Run("FixedArraysNotLimited_Marshal", func(t *testing.T) {
		// Fixed arrays should not be subject to marshal limits
		data := MarshalLimitTestData{
			SmallArray: [5]int{1, 2, 3, 4, 5}, // Fixed array size
		}

		_, err := data.MarshalMsg(nil)
		if err != nil {
			t.Errorf("Fixed arrays should not be limited during marshal, got error: %v", err)
		}
	})
}

func TestFieldLevelLimits(t *testing.T) {
	// Test field-level limit enforcement with FieldLimitTestData

	t.Run("SmallSlice_WithinLimit", func(t *testing.T) {
		data := FieldLimitTestData{
			SmallSlice: []int{1, 2, 3}, // Within limit of 5
		}

		// Marshal
		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		// Unmarshal
		var result FieldLimitTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err != nil {
			t.Fatalf("Unexpected unmarshal error: %v", err)
		}

		if len(result.SmallSlice) != len(data.SmallSlice) {
			t.Errorf("SmallSlice length mismatch: got %d, want %d", len(result.SmallSlice), len(data.SmallSlice))
		}
	})

	t.Run("SmallSlice_ExceedsLimit", func(t *testing.T) {
		data := FieldLimitTestData{
			SmallSlice: []int{1, 2, 3, 4, 5, 6, 7}, // Exceeds limit of 5
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		// Unmarshal should fail
		var result FieldLimitTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err == nil {
			t.Error("Expected error for SmallSlice exceeding limit, got nil")
		}
	})

	t.Run("SmallMap_WithinLimit", func(t *testing.T) {
		data := FieldLimitTestData{
			SmallMap: map[string]int{"a": 1, "b": 2}, // Within limit of 3
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		var result FieldLimitTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err != nil {
			t.Fatalf("Unexpected unmarshal error: %v", err)
		}

		if len(result.SmallMap) != len(data.SmallMap) {
			t.Errorf("SmallMap length mismatch: got %d, want %d", len(result.SmallMap), len(data.SmallMap))
		}
	})

	t.Run("SmallMap_ExceedsLimit", func(t *testing.T) {
		data := FieldLimitTestData{
			SmallMap: map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, // Exceeds limit of 3
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		// Unmarshal should fail
		var result FieldLimitTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err == nil {
			t.Error("Expected error for SmallMap exceeding limit, got nil")
		}
	})

	t.Run("FixedArrays_NotLimited_Field", func(t *testing.T) {
		// Fixed arrays should not be subject to field-level limits
		data := FieldLimitTestData{
			FixedArray: [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, // Fixed size, should not be limited
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error for fixed array: %v", err)
		}

		var result FieldLimitTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err != nil {
			t.Errorf("Fixed arrays should not be limited with field tags, got error: %v", err)
		}
	})

	t.Run("DecodeMsg_FieldLimits", func(t *testing.T) {
		data := FieldLimitTestData{
			SmallSlice: []int{1, 2, 3, 4, 5, 6}, // Exceeds limit of 5
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		// Test DecodeMsg path
		var result FieldLimitTestData
		reader := msgp.NewReader(bytes.NewReader(marshaled))
		err = result.DecodeMsg(reader)
		if err == nil {
			t.Error("Expected error for DecodeMsg with field exceeding limit, got nil")
		}
	})
}

func TestFieldVsFileLimitPrecedence(t *testing.T) {
	// Test precedence: field limits should override file limits
	// limits.go has: arrays:100 maps:50

	t.Run("TightSlice_FieldOverride", func(t *testing.T) {
		// File limit: arrays:100, field limit: 10 -> field limit should apply
		data := FieldOverrideTestData{
			TightSlice: make([]int, 15), // Exceeds field limit (10) but within file limit (100)
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		var result FieldOverrideTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err == nil {
			t.Error("Expected error for TightSlice exceeding field limit (10), got nil")
		}
	})

	t.Run("LooseSlice_FieldOverride", func(t *testing.T) {
		// File limit: arrays:100, field limit: 200 -> field limit should apply
		data := FieldOverrideTestData{
			LooseSlice: make([]string, 150), // Within field limit (200) but exceeds file limit (100)
		}
		for i := range data.LooseSlice {
			data.LooseSlice[i] = "test"
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		var result FieldOverrideTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err != nil {
			t.Errorf("Expected success for LooseSlice within field limit (200), got error: %v", err)
		}
	})

	t.Run("DefaultFields_UseFileLimit", func(t *testing.T) {
		// Fields without field limits should use file limits
		data := FieldOverrideTestData{
			DefaultSlice: make([]byte, 120), // Exceeds file limit (100)
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		var result FieldOverrideTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err == nil {
			t.Error("Expected error for DefaultSlice exceeding file limit (100), got nil")
		}
	})
}

func TestMarshalFieldVsFileLimitPrecedence(t *testing.T) {
	// Test precedence with marshal:true
	// marshal_limits.go has: arrays:30 maps:20 marshal:true

	t.Run("Marshal_TightSlice_FieldOverride", func(t *testing.T) {
		// File limit: arrays:30, field limit: 10 -> field limit should apply
		// Note: Since I only implemented field limits for unmarshal.go, marshal limits won't work yet
		data := MarshalFieldOverrideTestData{
			TightSlice: make([]int, 15), // Exceeds field limit (10) but within file limit (30)
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		var result MarshalFieldOverrideTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err == nil {
			t.Error("Expected error for TightSlice exceeding field limit (10), got nil")
		}
	})

	t.Run("Marshal_DefaultFields_UseFileLimit", func(t *testing.T) {
		// Fields without field limits should use file limits
		data := MarshalFieldOverrideTestData{
			DefaultSlice: make([]byte, 35), // Exceeds file limit (30)
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		var result MarshalFieldOverrideTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err == nil {
			t.Error("Expected error for DefaultSlice exceeding file limit (30), got nil")
		}
	})
}

func TestAliasedTypesWithFieldLimits(t *testing.T) {
	// Test field-level limits with aliased types

	t.Run("SmallAliasedSlice_WithinLimit", func(t *testing.T) {
		data := AliasedFieldLimitTestData{
			SmallAliasedSlice: AliasedSlice{"a", "b"}, // Within limit of 3
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		var result AliasedFieldLimitTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err != nil {
			t.Fatalf("Unexpected unmarshal error: %v", err)
		}

		if len(result.SmallAliasedSlice) != len(data.SmallAliasedSlice) {
			t.Errorf("SmallAliasedSlice length mismatch: got %d, want %d", len(result.SmallAliasedSlice), len(data.SmallAliasedSlice))
		}
	})

	t.Run("SmallAliasedSlice_ExceedsLimit", func(t *testing.T) {
		data := AliasedFieldLimitTestData{
			SmallAliasedSlice: AliasedSlice{"a", "b", "c", "d", "e"}, // Exceeds limit of 3
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		var result AliasedFieldLimitTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err == nil {
			t.Error("Expected error for SmallAliasedSlice exceeding limit (3), got nil")
		}
	})

	t.Run("SmallAliasedMap_ExceedsLimit", func(t *testing.T) {
		data := AliasedFieldLimitTestData{
			SmallAliasedMap: AliasedMap{
				"key1": true,
				"key2": false,
				"key3": true, // Exceeds limit of 2
			},
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		var result AliasedFieldLimitTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err == nil {
			t.Error("Expected error for SmallAliasedMap exceeding limit (2), got nil")
		}
	})

	t.Run("IntSliceAlias_ExceedsLimit", func(t *testing.T) {
		data := AliasedFieldLimitTestData{
			IntSliceAlias: make(AliasedIntSlice, 15), // Exceeds limit of 10
		}
		for i := range data.IntSliceAlias {
			data.IntSliceAlias[i] = i
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		var result AliasedFieldLimitTestData
		_, err = result.UnmarshalMsg(marshaled)
		if err == nil {
			t.Error("Expected error for IntSliceAlias exceeding limit (10), got nil")
		}
	})

	t.Run("DecodeMsg_AliasedTypes", func(t *testing.T) {
		data := AliasedFieldLimitTestData{
			SmallAliasedSlice: AliasedSlice{"a", "b", "c", "d"}, // Exceeds limit of 3
		}

		marshaled, err := data.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("Unexpected marshal error: %v", err)
		}

		// Test DecodeMsg path
		var result AliasedFieldLimitTestData
		reader := msgp.NewReader(bytes.NewReader(marshaled))
		err = result.DecodeMsg(reader)
		if err == nil {
			t.Error("Expected error for DecodeMsg with aliased slice exceeding limit, got nil")
		}
	})
}
