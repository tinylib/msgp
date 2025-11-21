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

func TestAllowNilFunctionality(t *testing.T) {
	// Test allownil functionality with DecodeMsg and UnmarshalMsg

	t.Run("AllowNil_NilValues_DecodeMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Create message with nil values for allownil fields
		buf := msgp.AppendMapHeader(nil, 8)

		// NilSlice - send nil
		buf = msgp.AppendString(buf, "nil_slice")
		buf = msgp.AppendNil(buf)

		// NilTightSlice - send nil
		buf = msgp.AppendString(buf, "nil_tight_slice")
		buf = msgp.AppendNil(buf)

		// NilLooseSlice - send nil
		buf = msgp.AppendString(buf, "nil_loose_slice")
		buf = msgp.AppendNil(buf)

		// NilMap - send nil
		buf = msgp.AppendString(buf, "nil_map")
		buf = msgp.AppendNil(buf)

		// NilTightMap - send nil
		buf = msgp.AppendString(buf, "nil_tight_map")
		buf = msgp.AppendNil(buf)

		// NilLooseMap - send nil
		buf = msgp.AppendString(buf, "nil_loose_map")
		buf = msgp.AppendNil(buf)

		// RegularSlice - send empty slice
		buf = msgp.AppendString(buf, "regular_slice")
		buf = msgp.AppendBytes(buf, []byte{})

		// RegularMap - send empty map
		buf = msgp.AppendString(buf, "regular_map")
		buf = msgp.AppendMapHeader(buf, 0)

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != nil {
			t.Fatalf("DecodeMsg failed with allownil fields: %v", err)
		}

		// Verify nil fields are properly handled
		if data.NilSlice != nil {
			t.Errorf("Expected NilSlice to be nil, got %v", data.NilSlice)
		}
		if data.NilTightSlice != nil {
			t.Errorf("Expected NilTightSlice to be nil, got %v", data.NilTightSlice)
		}
		if data.NilLooseSlice != nil {
			t.Errorf("Expected NilLooseSlice to be nil, got %v", data.NilLooseSlice)
		}
		if data.NilMap != nil {
			t.Errorf("Expected NilMap to be nil, got %v", data.NilMap)
		}
		if data.NilTightMap != nil {
			t.Errorf("Expected NilTightMap to be nil, got %v", data.NilTightMap)
		}
		if data.NilLooseMap != nil {
			t.Errorf("Expected NilLooseMap to be nil, got %v", data.NilLooseMap)
		}

		// Regular fields should be empty but allocated
		if data.RegularSlice == nil {
			t.Error("Expected RegularSlice to be allocated (empty), got nil")
		}
		if len(data.RegularSlice) != 0 {
			t.Errorf("Expected RegularSlice to be empty, got length %d", len(data.RegularSlice))
		}
		if data.RegularMap == nil {
			t.Error("Expected RegularMap to be allocated (empty), got nil")
		}
		if len(data.RegularMap) != 0 {
			t.Errorf("Expected RegularMap to be empty, got length %d", len(data.RegularMap))
		}
	})

	t.Run("AllowNil_NilValues_UnmarshalMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Create same message for unmarshal test
		buf := msgp.AppendMapHeader(nil, 8)
		buf = msgp.AppendString(buf, "nil_slice")
		buf = msgp.AppendNil(buf)
		buf = msgp.AppendString(buf, "nil_tight_slice")
		buf = msgp.AppendNil(buf)
		buf = msgp.AppendString(buf, "nil_loose_slice")
		buf = msgp.AppendNil(buf)
		buf = msgp.AppendString(buf, "nil_map")
		buf = msgp.AppendNil(buf)
		buf = msgp.AppendString(buf, "nil_tight_map")
		buf = msgp.AppendNil(buf)
		buf = msgp.AppendString(buf, "nil_loose_map")
		buf = msgp.AppendNil(buf)
		buf = msgp.AppendString(buf, "regular_slice")
		buf = msgp.AppendBytes(buf, []byte{})
		buf = msgp.AppendString(buf, "regular_map")
		buf = msgp.AppendMapHeader(buf, 0)

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("UnmarshalMsg failed with allownil fields: %v", err)
		}

		// Same verification as DecodeMsg test
		if data.NilSlice != nil {
			t.Errorf("Expected NilSlice to be nil, got %v", data.NilSlice)
		}
		if data.RegularSlice == nil {
			t.Error("Expected RegularSlice to be allocated (empty), got nil")
		}
	})

	t.Run("AllowNil_WithData_RespectLimits", func(t *testing.T) {
		data := AllowNilTestData{}

		// Test that allownil fields still respect limits when they have data
		buf := msgp.AppendMapHeader(nil, 3)

		// NilTightSlice with data exceeding limit (5)
		buf = msgp.AppendString(buf, "nil_tight_slice")
		buf = msgp.AppendArrayHeader(buf, 10) // Exceeds field limit of 5
		for i := 0; i < 10; i++ {
			buf = msgp.AppendInt(buf, i)
		}

		// NilLooseSlice with data within limit (200)
		buf = msgp.AppendString(buf, "nil_loose_slice")
		buf = msgp.AppendArrayHeader(buf, 150) // Within field limit of 200
		for i := 0; i < 150; i++ {
			buf = msgp.AppendString(buf, "test")
		}

		// NilTightMap with data exceeding limit (3)
		buf = msgp.AppendString(buf, "nil_tight_map")
		buf = msgp.AppendMapHeader(buf, 5) // Exceeds field limit of 3
		for i := 0; i < 5; i++ {
			buf = msgp.AppendString(buf, fmt.Sprintf("key%d", i))
			buf = msgp.AppendInt(buf, i)
		}

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for allownil field exceeding limit, got %v", err)
		}
	})

	t.Run("AllowNil_WithData_WithinLimits", func(t *testing.T) {
		data := AllowNilTestData{}

		// Test that allownil fields work normally when within limits
		buf := msgp.AppendMapHeader(nil, 3)

		// NilTightSlice with data within limit (5)
		buf = msgp.AppendString(buf, "nil_tight_slice")
		buf = msgp.AppendArrayHeader(buf, 3) // Within field limit of 5
		for i := 0; i < 3; i++ {
			buf = msgp.AppendInt(buf, i)
		}

		// NilSlice with data within file limit (100)
		buf = msgp.AppendString(buf, "nil_slice")
		buf = msgp.AppendBytes(buf, make([]byte, 50)) // Within file limit of 100

		// NilTightMap with data within limit (3)
		buf = msgp.AppendString(buf, "nil_tight_map")
		buf = msgp.AppendMapHeader(buf, 2) // Within field limit of 3
		buf = msgp.AppendString(buf, "key1")
		buf = msgp.AppendInt(buf, 1)
		buf = msgp.AppendString(buf, "key2")
		buf = msgp.AppendInt(buf, 2)

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("Unexpected error for allownil fields within limits: %v", err)
		}

		// Verify data was properly loaded
		if len(data.NilTightSlice) != 3 {
			t.Errorf("Expected NilTightSlice length 3, got %d", len(data.NilTightSlice))
		}
		if len(data.NilSlice) != 50 {
			t.Errorf("Expected NilSlice length 50, got %d", len(data.NilSlice))
		}
		if len(data.NilTightMap) != 2 {
			t.Errorf("Expected NilTightMap length 2, got %d", len(data.NilTightMap))
		}
	})
}

func TestAllowNilZeroSizedSlices(t *testing.T) {
	// Test the zero-sized slice allocation behavior

	t.Run("ZeroSized_NilBytes_DecodeMsg", func(t *testing.T) {
		data := AllowNilZeroTestData{}

		// Send zero-sized byte slice
		buf := msgp.AppendMapHeader(nil, 2)
		buf = msgp.AppendString(buf, "nil_bytes")
		buf = msgp.AppendBytes(buf, []byte{}) // Zero-sized bytes
		buf = msgp.AppendString(buf, "regular_bytes")
		buf = msgp.AppendBytes(buf, []byte{})

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != nil {
			t.Fatalf("DecodeMsg failed: %v", err)
		}

		// Both should be allocated but empty (not nil)
		if data.NilBytes == nil {
			t.Error("Expected NilBytes to be allocated (empty), got nil - allownil should allocate zero-sized slice")
		}
		if len(data.NilBytes) != 0 {
			t.Errorf("Expected NilBytes to be empty, got length %d", len(data.NilBytes))
		}
		if data.RegularBytes == nil {
			t.Error("Expected RegularBytes to be allocated (empty), got nil")
		}
		if len(data.RegularBytes) != 0 {
			t.Errorf("Expected RegularBytes to be empty, got length %d", len(data.RegularBytes))
		}
	})

	t.Run("ZeroSized_NilBytes_UnmarshalMsg", func(t *testing.T) {
		data := AllowNilZeroTestData{}

		// Same test for unmarshal
		buf := msgp.AppendMapHeader(nil, 2)
		buf = msgp.AppendString(buf, "nil_bytes")
		buf = msgp.AppendBytes(buf, []byte{})
		buf = msgp.AppendString(buf, "regular_bytes")
		buf = msgp.AppendBytes(buf, []byte{})

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("UnmarshalMsg failed: %v", err)
		}

		// Verify zero-sized allocation
		if data.NilBytes == nil {
			t.Error("Expected NilBytes to be allocated (empty), got nil - allownil should allocate zero-sized slice")
		}
		if len(data.NilBytes) != 0 {
			t.Errorf("Expected NilBytes to be empty, got length %d", len(data.NilBytes))
		}
	})

	t.Run("ActualNil_vs_ZeroSized", func(t *testing.T) {
		data1 := AllowNilZeroTestData{}
		data2 := AllowNilZeroTestData{}

		// Test 1: Send actual nil
		buf1 := msgp.AppendMapHeader(nil, 1)
		buf1 = msgp.AppendString(buf1, "nil_bytes")
		buf1 = msgp.AppendNil(buf1)

		// Test 2: Send zero-sized slice
		buf2 := msgp.AppendMapHeader(nil, 1)
		buf2 = msgp.AppendString(buf2, "nil_bytes")
		buf2 = msgp.AppendBytes(buf2, []byte{})

		_, err1 := data1.UnmarshalMsg(buf1)
		if err1 != nil {
			t.Fatalf("UnmarshalMsg with nil failed: %v", err1)
		}

		_, err2 := data2.UnmarshalMsg(buf2)
		if err2 != nil {
			t.Fatalf("UnmarshalMsg with zero-sized failed: %v", err2)
		}

		// data1 received actual nil - should remain nil for allownil field
		if data1.NilBytes != nil {
			t.Errorf("Expected data1.NilBytes (from nil) to remain nil, got %v", data1.NilBytes)
		}
		// data2 received zero-sized slice - should be allocated empty slice
		if data2.NilBytes == nil {
			t.Error("Expected data2.NilBytes (from zero-sized) to be allocated, got nil")
		}

		// Both should be empty
		if len(data1.NilBytes) != 0 {
			t.Errorf("Expected data1.NilBytes length 0, got %d", len(data1.NilBytes))
		}
		if len(data2.NilBytes) != 0 {
			t.Errorf("Expected data2.NilBytes length 0, got %d", len(data2.NilBytes))
		}
	})
}

func TestAllowNilSecurityLimits(t *testing.T) {
	// Test that allownil fields still enforce security limits properly

	t.Run("AllowNil_BytesExceedLimit_UnmarshalMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Try to send bytes that exceed file limit (100)
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_slice") // Uses file limit (100)
		buf = msgp.AppendBytes(buf, make([]byte, 150)) // Exceeds limit

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for allownil bytes exceeding limit, got %v", err)
		}
	})

	t.Run("AllowNil_BytesExceedLimit_DecodeMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Try to send bytes that exceed file limit (100)
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_slice") // Uses file limit (100)
		buf = msgp.AppendBytes(buf, make([]byte, 150)) // Exceeds limit

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for allownil bytes exceeding limit, got %v", err)
		}
	})

	t.Run("AllowNil_SliceExceedFieldLimit_UnmarshalMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Try to send slice that exceeds field limit
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_tight_slice") // Field limit = 5
		buf = msgp.AppendArrayHeader(buf, 8) // Exceeds field limit
		for i := 0; i < 8; i++ {
			buf = msgp.AppendInt(buf, i)
		}

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for allownil slice exceeding field limit, got %v", err)
		}
	})

	t.Run("AllowNil_SliceExceedFieldLimit_DecodeMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Try to send slice that exceeds field limit
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_tight_slice") // Field limit = 5
		buf = msgp.AppendArrayHeader(buf, 8) // Exceeds field limit
		for i := 0; i < 8; i++ {
			buf = msgp.AppendInt(buf, i)
		}

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for allownil slice exceeding field limit, got %v", err)
		}
	})

	t.Run("AllowNil_MapExceedFieldLimit_UnmarshalMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Try to send map that exceeds field limit
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_tight_map") // Field limit = 3
		buf = msgp.AppendMapHeader(buf, 5) // Exceeds field limit
		for i := 0; i < 5; i++ {
			buf = msgp.AppendString(buf, fmt.Sprintf("key%d", i))
			buf = msgp.AppendInt(buf, i)
		}

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for allownil map exceeding field limit, got %v", err)
		}
	})

	t.Run("AllowNil_MapExceedFieldLimit_DecodeMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Try to send map that exceeds field limit
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_tight_map") // Field limit = 3
		buf = msgp.AppendMapHeader(buf, 5) // Exceeds field limit
		for i := 0; i < 5; i++ {
			buf = msgp.AppendString(buf, fmt.Sprintf("key%d", i))
			buf = msgp.AppendInt(buf, i)
		}

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for allownil map exceeding field limit, got %v", err)
		}
	})

	t.Run("AllowNil_HeaderFirstSecurity_UnmarshalMsg", func(t *testing.T) {
		// Verify that allownil still uses header-first security approach
		data := AllowNilTestData{}

		// Create a message that would cause memory exhaustion if read data-first
		// This test ensures that even allownil fields check limits before allocation
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_slice")

		// Manually craft a message with a huge size header but don't include the data
		// This should fail at the limit check, not during data reading
		buf = append(buf, 0xc6) // bin32 format
		buf = append(buf, 0x00, 0x00, 0x27, 0x10) // size = 10000 (exceeds limit of 100)
		// Don't append actual data - should fail before trying to read it

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded (header-first check), got %v", err)
		}
	})

	t.Run("AllowNil_HeaderFirstSecurity_DecodeMsg", func(t *testing.T) {
		// Verify that allownil still uses header-first security approach with DecodeMsg
		data := AllowNilTestData{}

		// Create a message that would cause memory exhaustion if read data-first
		// This test ensures that even allownil fields check limits before allocation
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_slice")

		// Manually craft a message with a huge size header but don't include the data
		// This should fail at the limit check, not during data reading
		buf = append(buf, 0xc6) // bin32 format
		buf = append(buf, 0x00, 0x00, 0x27, 0x10) // size = 10000 (exceeds limit of 100)
		// Don't append actual data - should fail before trying to read it

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded (header-first check), got %v", err)
		}
	})
}

func TestAllowNilZeroCopy(t *testing.T) {
	// Test that zerocopy+allownil combination works correctly with security fixes

	t.Run("ZeroCopy_AllowNil_NilValues_UnmarshalMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Create message with nil for zerocopy+allownil field
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_zc_slice")
		buf = msgp.AppendNil(buf)

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("UnmarshalMsg failed with zerocopy allownil nil value: %v", err)
		}

		// Field should remain nil
		if data.NilZCSlice != nil {
			t.Errorf("Expected NilZCSlice to remain nil, got %v", data.NilZCSlice)
		}
	})

	t.Run("ZeroCopy_AllowNil_EmptyData_UnmarshalMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Create message with empty data for zerocopy+allownil field
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_zc_slice")
		buf = msgp.AppendBytes(buf, []byte{})

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("UnmarshalMsg failed with zerocopy allownil empty data: %v", err)
		}

		// Field should be allocated but empty
		if data.NilZCSlice == nil {
			t.Error("Expected NilZCSlice to be allocated (empty), got nil")
		}
		if len(data.NilZCSlice) != 0 {
			t.Errorf("Expected NilZCSlice to be empty, got length %d", len(data.NilZCSlice))
		}
	})

	t.Run("ZeroCopy_AllowNil_WithData_UnmarshalMsg", func(t *testing.T) {
		data := AllowNilTestData{}
		testData := []byte{1, 2, 3, 4, 5}

		// Create message with actual data for zerocopy+allownil field
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_zc_slice")
		buf = msgp.AppendBytes(buf, testData)

		originalBuf := make([]byte, len(buf))
		copy(originalBuf, buf)

		_, err := data.UnmarshalMsg(buf)
		if err != nil {
			t.Fatalf("UnmarshalMsg failed with zerocopy allownil data: %v", err)
		}

		// Field should contain the data
		if data.NilZCSlice == nil {
			t.Error("Expected NilZCSlice to be allocated, got nil")
		}
		if len(data.NilZCSlice) != len(testData) {
			t.Errorf("Expected NilZCSlice length %d, got %d", len(testData), len(data.NilZCSlice))
		}
		for i, b := range testData {
			if data.NilZCSlice[i] != b {
				t.Errorf("Expected NilZCSlice[%d] = %d, got %d", i, b, data.NilZCSlice[i])
			}
		}
	})

	t.Run("ZeroCopy_AllowNil_SecurityLimits_UnmarshalMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Test that zerocopy+allownil still enforces limits
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_zc_tight_slice") // Field limit = 5

		// Create oversized data that should exceed field limit
		oversizedData := make([]byte, 10) // Exceeds field limit of 5
		buf = msgp.AppendBytes(buf, oversizedData)

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for zerocopy allownil field exceeding limit, got %v", err)
		}
	})

	t.Run("ZeroCopy_AllowNil_HeaderFirstSecurity_UnmarshalMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Test header-first security for zerocopy+allownil
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_zc_slice")

		// Manually craft message with huge header but no data
		buf = append(buf, 0xc6) // bin32 format
		buf = append(buf, 0x00, 0x00, 0x27, 0x10) // size = 10000 (exceeds limit)

		_, err := data.UnmarshalMsg(buf)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded (header-first check) for zerocopy allownil, got %v", err)
		}
	})

	// Add DecodeMsg versions for complete coverage
	t.Run("ZeroCopy_AllowNil_NilValues_DecodeMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Create message with nil for zerocopy+allownil field
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_zc_slice")
		buf = msgp.AppendNil(buf)

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != nil {
			t.Fatalf("DecodeMsg failed with zerocopy allownil nil value: %v", err)
		}

		// Field should remain nil
		if data.NilZCSlice != nil {
			t.Errorf("Expected NilZCSlice to remain nil, got %v", data.NilZCSlice)
		}
	})

	t.Run("ZeroCopy_AllowNil_EmptyData_DecodeMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Create message with empty data for zerocopy+allownil field
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_zc_slice")
		buf = msgp.AppendBytes(buf, []byte{})

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != nil {
			t.Fatalf("DecodeMsg failed with zerocopy allownil empty data: %v", err)
		}

		// Field should be allocated but empty
		if data.NilZCSlice == nil {
			t.Error("Expected NilZCSlice to be allocated (empty), got nil")
		}
		if len(data.NilZCSlice) != 0 {
			t.Errorf("Expected NilZCSlice to be empty, got length %d", len(data.NilZCSlice))
		}
	})

	t.Run("ZeroCopy_AllowNil_WithData_DecodeMsg", func(t *testing.T) {
		data := AllowNilTestData{}
		testData := []byte{1, 2, 3, 4, 5}

		// Create message with actual data for zerocopy+allownil field
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_zc_slice")
		buf = msgp.AppendBytes(buf, testData)

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != nil {
			t.Fatalf("DecodeMsg failed with zerocopy allownil data: %v", err)
		}

		// Field should contain the data
		if data.NilZCSlice == nil {
			t.Error("Expected NilZCSlice to be allocated, got nil")
		}
		if len(data.NilZCSlice) != len(testData) {
			t.Errorf("Expected NilZCSlice length %d, got %d", len(testData), len(data.NilZCSlice))
		}
		for i, b := range testData {
			if data.NilZCSlice[i] != b {
				t.Errorf("Expected NilZCSlice[%d] = %d, got %d", i, b, data.NilZCSlice[i])
			}
		}
	})

	t.Run("ZeroCopy_AllowNil_SecurityLimits_DecodeMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Test that zerocopy+allownil still enforces limits
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_zc_tight_slice") // Field limit = 5

		// Create oversized data that should exceed field limit
		oversizedData := make([]byte, 10) // Exceeds field limit of 5
		buf = msgp.AppendBytes(buf, oversizedData)

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded for zerocopy allownil field exceeding limit, got %v", err)
		}
	})

	t.Run("ZeroCopy_AllowNil_HeaderFirstSecurity_DecodeMsg", func(t *testing.T) {
		data := AllowNilTestData{}

		// Test header-first security for zerocopy+allownil
		buf := msgp.AppendMapHeader(nil, 1)
		buf = msgp.AppendString(buf, "nil_zc_slice")

		// Manually craft message with huge header but no data
		buf = append(buf, 0xc6) // bin32 format
		buf = append(buf, 0x00, 0x00, 0x27, 0x10) // size = 10000 (exceeds limit)

		reader := msgp.NewReader(bytes.NewReader(buf))
		err := data.DecodeMsg(reader)
		if err != msgp.ErrLimitExceeded {
			t.Errorf("Expected ErrLimitExceeded (header-first check) for zerocopy allownil, got %v", err)
		}
	})
}
