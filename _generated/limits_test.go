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
