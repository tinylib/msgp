package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const template = `
// DecodeMsg decodes the message from the reader.
func (s *Foo) DecodeMsg(reader *msgp.Reader) error {
	if reader.IsNil() {
		*s = nil
		return reader.Skip()
	}
	sz, err := reader.ReadArrayHeader()
	if err != nil {
		return err
	}
	dst := *s
	if dst == nil {
		dst = make(Foo, sz)
	} else {
		clear(dst)
	}
	for range sz {
		var k string
		k, err = reader.ReadString()
		if err != nil {
			return err
		}
		dst[string(k)] = struct{}{}
	}
	*s = dst
	return nil
}

// UnmarshalMsg decodes the message from the bytes.
func (s *Foo) UnmarshalMsg(bytes []byte) ([]byte, error) {
	if msgp.IsNil(bytes) {
		*s = nil
		return bytes[msgp.NilSize:], nil
	}
	// Read the array header
	sz, bytes, err := msgp.ReadArrayHeaderBytes(bytes)
	if err != nil {
		return nil, err
	}
	dst := *s
	if dst == nil {
		dst = make(Foo, sz)
	} else {
		clear(dst)
	}
	for range sz {
		var k string
		k, bytes, err = msgp.ReadStringBytes(bytes)
		if err != nil {
			return nil, err
		}
		dst[string(k)] = struct{}{}
	}
	*s = dst
	return bytes, nil
}

// Msgsize returns the maximum size of the message.
func (s Foo) Msgsize() int {
	if s == nil {
		return msgp.NilSize
	}
	size := msgp.ArrayHeaderSize
	size += len(s) * msgp.StringPrefixSize
	return size
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Foo) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*s = nil
		return nil
	}
	dst := *s
	if dst != nil {
		clear(dst)
	} else {
		dst = make(Foo)
	}
	err := jsonArrayIter(data, func(raw []byte) error {
		k, parseErr := jsonParseQuoted(raw)
		if parseErr != nil {
			return parseErr
		}
		dst[k] = struct{}{}
		return nil
	})
	if err != nil {
		return err
	}
	*s = dst
	return nil
}

// FooFromSlice creates a Foo from a slice.
func FooFromSlice(s []string) Foo {
	if s == nil {
		return nil
	}
	dst := make(Foo, len(s))
	for _, v := range s {
		dst[v] = struct{}{}
	}
	return dst
}
`

const unsorted = `
// Foo is a set of strings that will be stored as an array.
// Elements are not sorted and the order of elements is not guaranteed.
type Foo map[string]struct{}

// EncodeMsg encodes the message to the writer.
func (s Foo) EncodeMsg(writer *msgp.Writer) error {
	if s == nil {
		return writer.WriteNil()
	}
	err := writer.WriteArrayHeader(uint32(len(s)))
	if err != nil {
		return err
	}
	for k := range s {
		err = writer.WriteString(k)
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalMsg encodes the message to the bytes.
func (s Foo) MarshalMsg(bytes []byte) ([]byte, error) {
	if s == nil {
		return msgp.AppendNil(bytes), nil
	}
	if len(s) == 0 {
		return msgp.AppendArrayHeader(bytes, 0), nil
	}
	bytes = ensure(bytes, s.Msgsize())
	bytes = msgp.AppendArrayHeader(bytes, uint32(len(s)))
	for k := range s {
		bytes = msgp.AppendString(bytes, string(k))
	}
	return bytes, nil
}

// AsSlice returns the set as a slice.
func (s Foo) AsSlice() []string {
	if s == nil {
		return nil
	}
	dst := make([]string, 0, len(s))
	for k := range s {
		dst = append(dst, k)
	}
	return dst
}

// MarshalJSON implements json.Marshaler.
func (s Foo) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	dst := make([]byte, 0, 2+len(s)*10)
	dst = append(dst, '[')
	first := true
	for k := range s {
		if !first {
			dst = append(dst, ',')
		}
		first = false
		dst = jsonAppendQuote(dst, k)
	}
	dst = append(dst, ']')
	return dst, nil
}
`

const sorted = `
// Foo is a set of strings that will be stored as an array.
// Elements are sorted and the order of elements is guaranteed.
type Foo map[string]struct{}

// EncodeMsg encodes the message to the writer.
func (s Foo) EncodeMsg(writer *msgp.Writer) error {
	if s == nil {
		return writer.WriteNil()
	}
	err := writer.WriteArrayHeader(uint32(len(s)))
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b string) int {
		if a < b {
			return -1
		}
		return 1
	})

	for _, k := range keys {
		err = writer.WriteString(k)
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalMsg encodes the message to the bytes.
func (s Foo) MarshalMsg(bytes []byte) ([]byte, error) {
	if s == nil {
		return msgp.AppendNil(bytes), nil
	}
	if len(s) == 0 {
		return msgp.AppendArrayHeader(bytes, 0), nil
	}
	bytes = ensure(bytes, s.Msgsize())
	bytes = msgp.AppendArrayHeader(bytes, uint32(len(s)))
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b string) int {
		if a < b {
			return -1
		}
		return 1
	})
	for _, k := range keys {
		bytes = msgp.AppendString(bytes, k)
	}
	return bytes, nil
}

// AsSlice returns the set as a sorted slice.
func (s Foo) AsSlice() []string {
	if s == nil {
		return nil
	}
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b string) int {
		if a < b {
			return -1
		}
		return 1
	})
	return keys
}

// MarshalJSON implements json.Marshaler.
func (s Foo) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	keys := s.AsSlice()
	dst := make([]byte, 0, 2+len(keys)*10)
	dst = append(dst, '[')
	for i, k := range keys {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = jsonAppendQuote(dst, k)
	}
	dst = append(dst, ']')
	return dst, nil
}
`

const testTemplate = `
func Test{{.TypeName}}_RoundTrip(t *testing.T) {
	set := make({{.TypeName}})
	{{.PopulateSet}}

	// Test EncodeMsg/DecodeMsg
	var buf bytes.Buffer
	writer := msgp.NewWriter(&buf)
	err := set.EncodeMsg(writer)
	if err != nil {
		t.Fatalf("EncodeMsg failed: %v", err)
	}
	writer.Flush()

	reader := msgp.NewReader(&buf)
	var decoded {{.TypeName}}
	err = decoded.DecodeMsg(reader)
	if err != nil {
		t.Fatalf("DecodeMsg failed: %v", err)
	}

	if len(set) != len(decoded) {
		t.Fatalf("length mismatch: expected %d, got %d", len(set), len(decoded))
	}

	for k := range set {
		if _, ok := decoded[k]; !ok {
			t.Fatalf("missing key: %v", k)
		}
	}

	// Test MarshalMsg/UnmarshalMsg
	data, err := set.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("MarshalMsg failed: %v", err)
	}

	var unmarshaled {{.TypeName}}
	_, err = unmarshaled.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("UnmarshalMsg failed: %v", err)
	}

	if len(set) != len(unmarshaled) {
		t.Fatalf("length mismatch: expected %d, got %d", len(set), len(unmarshaled))
	}

	for k := range set {
		if _, ok := unmarshaled[k]; !ok {
			t.Fatalf("missing key: %v", k)
		}
	}
}

func Test{{.TypeName}}_AsSlice(t *testing.T) {
	set := make({{.TypeName}})
	{{.PopulateSet}}

	slice := set.AsSlice()
	if len(slice) != len(set) {
		t.Fatalf("length mismatch: expected %d, got %d", len(set), len(slice))
	}

	found := make(map[{{.GoType}}]bool)
	for _, v := range slice {
		found[v] = true
	}

	for k := range set {
		if !found[k] {
			t.Fatalf("missing key in slice: %v", k)
		}
	}
}

func Test{{.TypeName}}_FromSlice(t *testing.T) {
	slice := []{{.GoType}}{{{.SliceValues}}}
	set := {{.TypeName}}FromSlice(slice)

	if len(set) != len(slice) {
		t.Fatalf("length mismatch: expected %d, got %d", len(slice), len(set))
	}

	for _, v := range slice {
		if _, ok := set[v]; !ok {
			t.Fatalf("missing key: %v", v)
		}
	}
}

func Test{{.TypeName}}_NilHandling(t *testing.T) {
	var nilSet {{.TypeName}}

	// Test nil encoding
	var buf bytes.Buffer
	writer := msgp.NewWriter(&buf)
	err := nilSet.EncodeMsg(writer)
	if err != nil {
		t.Fatalf("EncodeMsg failed for nil: %v", err)
	}
	writer.Flush()

	// Test nil decoding
	reader := msgp.NewReader(&buf)
	var decoded {{.TypeName}}
	err = decoded.DecodeMsg(reader)
	if err != nil {
		t.Fatalf("DecodeMsg failed for nil: %v", err)
	}

	if decoded != nil {
		t.Fatal("expected nil, got non-nil")
	}

	// Test nil marshaling
	data, err := nilSet.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("MarshalMsg failed for nil: %v", err)
	}

	// Test nil unmarshaling
	var unmarshaled {{.TypeName}}
	_, err = unmarshaled.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("UnmarshalMsg failed for nil: %v", err)
	}

	if unmarshaled != nil {
		t.Fatal("expected nil, got non-nil")
	}

	// Test AsSlice on nil
	slice := nilSet.AsSlice()
	if slice != nil {
		t.Fatal("expected nil slice, got non-nil")
	}

	// Test FromSlice with nil
	fromNilSlice := {{.TypeName}}FromSlice(nil)
	if fromNilSlice != nil {
		t.Fatal("expected nil from nil slice, got non-nil")
	}
}

func Test{{.TypeName}}_EmptySet(t *testing.T) {
	set := make({{.TypeName}})

	// Test empty set encoding
	var buf bytes.Buffer
	writer := msgp.NewWriter(&buf)
	err := set.EncodeMsg(writer)
	if err != nil {
		t.Fatalf("EncodeMsg failed for empty: %v", err)
	}
	writer.Flush()

	// Test empty set decoding
	reader := msgp.NewReader(&buf)
	var decoded {{.TypeName}}
	err = decoded.DecodeMsg(reader)
	if err != nil {
		t.Fatalf("DecodeMsg failed for empty: %v", err)
	}

	if len(decoded) != 0 {
		t.Fatalf("expected empty set, got length %d", len(decoded))
	}

	// Test empty set marshaling
	data, err := set.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("MarshalMsg failed for empty: %v", err)
	}

	// Test empty set unmarshaling
	var unmarshaled {{.TypeName}}
	_, err = unmarshaled.UnmarshalMsg(data)
	if err != nil {
		t.Fatalf("UnmarshalMsg failed for empty: %v", err)
	}

	if len(unmarshaled) != 0 {
		t.Fatalf("expected empty set, got length %d", len(unmarshaled))
	}
}

func Test{{.TypeName}}_JSONRoundTrip(t *testing.T) {
	set := make({{.TypeName}})
	{{.PopulateSet}}

	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Verify JSON is an array
	if len(data) == 0 || data[0] != '[' || data[len(data)-1] != ']' {
		t.Fatalf("expected JSON array, got: %s", data)
	}

	// Verify it unmarshals as a generic JSON array
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		t.Fatalf("JSON is not a valid array: %v", err)
	}
	if len(arr) != len(set) {
		t.Fatalf("JSON array length mismatch: expected %d, got %d", len(set), len(arr))
	}

	var decoded {{.TypeName}}
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if len(set) != len(decoded) {
		t.Fatalf("length mismatch: expected %d, got %d", len(set), len(decoded))
	}

	for k := range set {
		if _, ok := decoded[k]; !ok {
			t.Fatalf("missing key: %v", k)
		}
	}
}

func Test{{.TypeName}}_JSONNil(t *testing.T) {
	var nilSet {{.TypeName}}

	data, err := json.Marshal(nilSet)
	if err != nil {
		t.Fatalf("MarshalJSON nil failed: %v", err)
	}
	if string(data) != "null" {
		t.Fatalf("expected null, got %s", data)
	}

	var decoded {{.TypeName}}
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("UnmarshalJSON nil failed: %v", err)
	}
	if decoded != nil {
		t.Fatal("expected nil, got non-nil")
	}
}

func Test{{.TypeName}}_JSONEmpty(t *testing.T) {
	set := make({{.TypeName}})

	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("MarshalJSON empty failed: %v", err)
	}
	if string(data) != "[]" {
		t.Fatalf("expected [], got %s", data)
	}

	var decoded {{.TypeName}}
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("UnmarshalJSON empty failed: %v", err)
	}
	if len(decoded) != 0 {
		t.Fatalf("expected empty set, got length %d", len(decoded))
	}
}
`

const benchTemplate = `
func Benchmark{{.TypeName}}_EncodeMsg(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			set := make({{.TypeName}})
			{{.GeneratePopulateCode}}

			var buf bytes.Buffer
			writer := msgp.NewWriter(&buf)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf.Reset()
				writer.Reset(&buf)
				set.EncodeMsg(writer)
				writer.Flush()
			}
		})
	}
}

func Benchmark{{.TypeName}}_DecodeMsg(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			set := make({{.TypeName}})
			{{.GeneratePopulateCode}}

			var buf bytes.Buffer
			writer := msgp.NewWriter(&buf)
			set.EncodeMsg(writer)
			writer.Flush()
			encoded := buf.Bytes()

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reader := msgp.NewReader(bytes.NewReader(encoded))
				var decoded {{.TypeName}}
				decoded.DecodeMsg(reader)
			}
		})
	}
}

func Benchmark{{.TypeName}}_MarshalMsg(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			set := make({{.TypeName}})
			{{.GeneratePopulateCode}}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := set.MarshalMsg(nil)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func Benchmark{{.TypeName}}_UnmarshalMsg(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			set := make({{.TypeName}})
			{{.GeneratePopulateCode}}

			data, _ := set.MarshalMsg(nil)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var decoded {{.TypeName}}
				_, err := decoded.UnmarshalMsg(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func Benchmark{{.TypeName}}_AsSlice(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			set := make({{.TypeName}})
			{{.GeneratePopulateCode}}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = set.AsSlice()
			}
		})
	}
}

func Benchmark{{.TypeName}}_FromSlice(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			{{.GenerateSliceCode}}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = {{.TypeName}}FromSlice(slice)
			}
		})
	}
}

func Benchmark{{.TypeName}}_MarshalJSON(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			set := make({{.TypeName}})
			{{.GeneratePopulateCode}}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := json.Marshal(set)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func Benchmark{{.TypeName}}_UnmarshalJSON(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			set := make({{.TypeName}})
			{{.GeneratePopulateCode}}

			data, _ := json.Marshal(set)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var decoded {{.TypeName}}
				err := json.Unmarshal(data, &decoded)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
`

type replacer struct {
	GoType      string // 'string'
	PackageName string // 'Foo'
	DecodeValue string // 'ReadString'
	EncodeValue string // 'WriteString'
	AppendValue string // 'AppendString'
	KeyLen      string // 'size += msgp.StringPrefixSize'
	Sorter      string // 'sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })'

	JSONAppendKey string // 'jsonAppendQuote(dst, k)'
	JSONParseKey  string // 'jsonParseQuoted(raw)'
}

type testGen struct {
	TypeName    string
	GoType      string
	PopulateSet string
	SliceValues string
	Size        string
}

func generateTestValues(goType string, size int) (populateSet, sliceValues string) {
	var populate []string
	var values []string

	switch goType {
	case "string":
		for i := 0; i < size; i++ {
			val := fmt.Sprintf(`"val%d"`, i)
			populate = append(populate, fmt.Sprintf("set[%s] = struct{}{}", val))
			values = append(values, val)
		}
	case "int8":
		for i := 0; i < size; i++ {
			val := fmt.Sprintf("%d", (i%256)-128) // int8 range: -128 to 127
			populate = append(populate, fmt.Sprintf("set[%s] = struct{}{}", val))
			values = append(values, val)
		}
	case "int16":
		for i := 0; i < size; i++ {
			val := fmt.Sprintf("%d", (i%65536)-32768) // int16 range: -32768 to 32767
			populate = append(populate, fmt.Sprintf("set[%s] = struct{}{}", val))
			values = append(values, val)
		}
	case "int", "int32", "int64":
		for i := 0; i < size; i++ {
			val := fmt.Sprintf("%d", i)
			populate = append(populate, fmt.Sprintf("set[%s] = struct{}{}", val))
			values = append(values, val)
		}
	case "byte":
		for i := 0; i < size; i++ {
			val := fmt.Sprintf("%d", i%256) // Prevent byte overflow
			populate = append(populate, fmt.Sprintf("set[%s] = struct{}{}", val))
			values = append(values, val)
		}
	case "uint8":
		for i := 0; i < size; i++ {
			val := fmt.Sprintf("%d", i%256) // Prevent uint8 overflow
			populate = append(populate, fmt.Sprintf("set[%s] = struct{}{}", val))
			values = append(values, val)
		}
	case "uint16":
		for i := 0; i < size; i++ {
			val := fmt.Sprintf("%d", i%65536) // Prevent uint16 overflow
			populate = append(populate, fmt.Sprintf("set[%s] = struct{}{}", val))
			values = append(values, val)
		}
	case "uint", "uint32", "uint64":
		for i := 0; i < size; i++ {
			val := fmt.Sprintf("%d", i)
			populate = append(populate, fmt.Sprintf("set[%s] = struct{}{}", val))
			values = append(values, val)
		}
	case "float32", "float64":
		for i := 0; i < size; i++ {
			val := fmt.Sprintf("%d.0", i)
			populate = append(populate, fmt.Sprintf("set[%s] = struct{}{}", val))
			values = append(values, val)
		}
	}

	return strings.Join(populate, "\n\t"), strings.Join(values, ", ")
}

func generateDynamicPopulateCode(goType string) string {
	switch goType {
	case "string":
		return `for i := 0; i < size; i++ {
				set[fmt.Sprintf("val%d", i)] = struct{}{}
			}`
	case "int8":
		return `for i := 0; i < size; i++ {
				set[int8((i%256)-128)] = struct{}{}
			}`
	case "int16":
		return `for i := 0; i < size; i++ {
				set[int16((i%65536)-32768)] = struct{}{}
			}`
	case "int", "int32", "int64":
		return `for i := 0; i < size; i++ {
				set[` + goType + `(i)] = struct{}{}
			}`
	case "byte":
		return `for i := 0; i < size; i++ {
				set[byte(i%256)] = struct{}{}
			}`
	case "uint8":
		return `for i := 0; i < size; i++ {
				set[uint8(i%256)] = struct{}{}
			}`
	case "uint16":
		return `for i := 0; i < size; i++ {
				set[uint16(i%65536)] = struct{}{}
			}`
	case "uint", "uint32", "uint64":
		return `for i := 0; i < size; i++ {
				set[` + goType + `(i)] = struct{}{}
			}`
	case "float32", "float64":
		return `for i := 0; i < size; i++ {
				set[` + goType + `(i)] = struct{}{}
			}`
	}
	return ""
}

func generateDynamicSliceCode(goType string) string {
	switch goType {
	case "string":
		return `slice := make([]string, size)
			for i := 0; i < size; i++ {
				slice[i] = fmt.Sprintf("val%d", i)
			}`
	case "int8":
		return `slice := make([]int8, size)
			for i := 0; i < size; i++ {
				slice[i] = int8((i%256)-128)
			}`
	case "int16":
		return `slice := make([]int16, size)
			for i := 0; i < size; i++ {
				slice[i] = int16((i%65536)-32768)
			}`
	case "int", "int32", "int64":
		return `slice := make([]` + goType + `, size)
			for i := 0; i < size; i++ {
				slice[i] = ` + goType + `(i)
			}`
	case "byte":
		return `slice := make([]byte, size)
			for i := 0; i < size; i++ {
				slice[i] = byte(i%256)
			}`
	case "uint8":
		return `slice := make([]uint8, size)
			for i := 0; i < size; i++ {
				slice[i] = uint8(i%256)
			}`
	case "uint16":
		return `slice := make([]uint16, size)
			for i := 0; i < size; i++ {
				slice[i] = uint16(i%65536)
			}`
	case "uint", "uint32", "uint64":
		return `slice := make([]` + goType + `, size)
			for i := 0; i < size; i++ {
				slice[i] = ` + goType + `(i)
			}`
	case "float32", "float64":
		return `slice := make([]` + goType + `, size)
			for i := 0; i < size; i++ {
				slice[i] = ` + goType + `(i)
			}`
	}
	return ""
}

func generateTests(out *os.File, r replacer) {
	// Generate basic tests (using small set)
	populateSet, sliceValues := generateTestValues(r.GoType, 5)

	// Regular type tests
	testCode := testTemplate
	testCode = strings.ReplaceAll(testCode, "{{.TypeName}}", r.PackageName)
	testCode = strings.ReplaceAll(testCode, "{{.GoType}}", r.GoType)
	testCode = strings.ReplaceAll(testCode, "{{.PopulateSet}}", populateSet)
	testCode = strings.ReplaceAll(testCode, "{{.SliceValues}}", sliceValues)
	fmt.Fprintln(out, testCode)

	// Sorted type tests
	testCode = testTemplate
	testCode = strings.ReplaceAll(testCode, "{{.TypeName}}", r.PackageName+"Sorted")
	testCode = strings.ReplaceAll(testCode, "{{.GoType}}", r.GoType)
	testCode = strings.ReplaceAll(testCode, "{{.PopulateSet}}", populateSet)
	testCode = strings.ReplaceAll(testCode, "{{.SliceValues}}", sliceValues)
	fmt.Fprintln(out, testCode)

	// Generate consolidated benchmarks
	populateCode := generateDynamicPopulateCode(r.GoType)
	sliceCode := generateDynamicSliceCode(r.GoType)

	// Regular type benchmarks
	benchCode := benchTemplate
	benchCode = strings.ReplaceAll(benchCode, "{{.TypeName}}", r.PackageName)
	benchCode = strings.ReplaceAll(benchCode, "{{.GoType}}", r.GoType)
	benchCode = strings.ReplaceAll(benchCode, "{{.GeneratePopulateCode}}", populateCode)
	benchCode = strings.ReplaceAll(benchCode, "{{.GenerateSliceCode}}", sliceCode)
	fmt.Fprintln(out, benchCode)

	// Sorted type benchmarks
	benchCode = benchTemplate
	benchCode = strings.ReplaceAll(benchCode, "{{.TypeName}}", r.PackageName+"Sorted")
	benchCode = strings.ReplaceAll(benchCode, "{{.GoType}}", r.GoType)
	benchCode = strings.ReplaceAll(benchCode, "{{.GeneratePopulateCode}}", populateCode)
	benchCode = strings.ReplaceAll(benchCode, "{{.GenerateSliceCode}}", sliceCode)
	fmt.Fprintln(out, benchCode)
}

var replacers = []replacer{
	{
		GoType:        "string",
		PackageName:   "Foo",
		DecodeValue:   "ReadString",
		EncodeValue:   "WriteString",
		AppendValue:   "AppendString",
		KeyLen:        "size += len(s) * msgp.StringPrefixSize",
		Sorter:        "slices.SortFunc(keys, func(a, b string) int {\n\t\tif a < b {\n\t\t\treturn -1\n\t\t}\n\t\treturn 1\n\t})",
		JSONAppendKey: "jsonAppendQuote(dst, k)",
		JSONParseKey:  "jsonParseQuoted(raw)",
	},
	{
		GoType:        "string",
		PackageName:   "String",
		DecodeValue:   "ReadString",
		EncodeValue:   "WriteString",
		AppendValue:   "AppendString",
		KeyLen:        "for key := range s {\n\t\t\tsize += msgp.StringPrefixSize + len(key)\n\t\t}",
		Sorter:        "sort.Strings(keys)",
		JSONAppendKey: "jsonAppendQuote(dst, k)",
		JSONParseKey:  "jsonParseQuoted(raw)",
	},
	{
		GoType:        "int",
		PackageName:   "Int",
		DecodeValue:   "ReadInt",
		EncodeValue:   "WriteInt",
		AppendValue:   "AppendInt",
		KeyLen:        "size += len(s) * msgp.IntSize",
		JSONAppendKey: "strconv.AppendInt(dst, int64(k), 10)",
		JSONParseKey:  "jsonParseSigned[int](raw, 0)",
	},
	{
		GoType:        "uint",
		PackageName:   "Uint",
		DecodeValue:   "ReadUint",
		EncodeValue:   "WriteUint",
		AppendValue:   "AppendUint",
		KeyLen:        "size += len(s) * msgp.UintSize",
		JSONAppendKey: "strconv.AppendUint(dst, uint64(k), 10)",
		JSONParseKey:  "jsonParseUnsigned[uint](raw, 0)",
	},
	{
		GoType:        "byte",
		PackageName:   "Byte",
		DecodeValue:   "ReadByte",
		EncodeValue:   "WriteByte",
		AppendValue:   "AppendByte",
		KeyLen:        "size += len(s) * msgp.ByteSize",
		JSONAppendKey: "strconv.AppendUint(dst, uint64(k), 10)",
		JSONParseKey:  "jsonParseUnsigned[byte](raw, 8)",
	},
	{
		GoType:        "int8",
		PackageName:   "Int8",
		DecodeValue:   "ReadInt8",
		EncodeValue:   "WriteInt8",
		AppendValue:   "AppendInt8",
		KeyLen:        "size += len(s) * msgp.Int8Size",
		JSONAppendKey: "strconv.AppendInt(dst, int64(k), 10)",
		JSONParseKey:  "jsonParseSigned[int8](raw, 8)",
	},
	{
		GoType:        "uint8",
		PackageName:   "Uint8",
		DecodeValue:   "ReadUint8",
		EncodeValue:   "WriteUint8",
		AppendValue:   "AppendUint8",
		KeyLen:        "size += len(s) * msgp.Uint8Size",
		JSONAppendKey: "strconv.AppendUint(dst, uint64(k), 10)",
		JSONParseKey:  "jsonParseUnsigned[uint8](raw, 8)",
	},
	{
		GoType:        "int16",
		PackageName:   "Int16",
		DecodeValue:   "ReadInt16",
		EncodeValue:   "WriteInt16",
		AppendValue:   "AppendInt16",
		KeyLen:        "size += len(s) * msgp.Int16Size",
		JSONAppendKey: "strconv.AppendInt(dst, int64(k), 10)",
		JSONParseKey:  "jsonParseSigned[int16](raw, 16)",
	},
	{
		GoType:        "uint16",
		PackageName:   "Uint16",
		DecodeValue:   "ReadUint16",
		EncodeValue:   "WriteUint16",
		AppendValue:   "AppendUint16",
		KeyLen:        "size += len(s) * msgp.Uint16Size",
		JSONAppendKey: "strconv.AppendUint(dst, uint64(k), 10)",
		JSONParseKey:  "jsonParseUnsigned[uint16](raw, 16)",
	},
	{
		GoType:        "int32",
		PackageName:   "Int32",
		DecodeValue:   "ReadInt32",
		EncodeValue:   "WriteInt32",
		AppendValue:   "AppendInt32",
		KeyLen:        "size += len(s) * msgp.Int32Size",
		JSONAppendKey: "strconv.AppendInt(dst, int64(k), 10)",
		JSONParseKey:  "jsonParseSigned[int32](raw, 32)",
	},
	{
		GoType:        "uint32",
		PackageName:   "Uint32",
		DecodeValue:   "ReadUint32",
		EncodeValue:   "WriteUint32",
		AppendValue:   "AppendUint32",
		KeyLen:        "size += len(s) * msgp.Uint32Size",
		JSONAppendKey: "strconv.AppendUint(dst, uint64(k), 10)",
		JSONParseKey:  "jsonParseUnsigned[uint32](raw, 32)",
	},
	{
		GoType:        "int64",
		PackageName:   "Int64",
		DecodeValue:   "ReadInt64",
		EncodeValue:   "WriteInt64",
		AppendValue:   "AppendInt64",
		KeyLen:        "size += len(s) * msgp.Int64Size",
		JSONAppendKey: "strconv.AppendInt(dst, int64(k), 10)",
		JSONParseKey:  "jsonParseSigned[int64](raw, 64)",
	},
	{
		GoType:        "uint64",
		PackageName:   "Uint64",
		DecodeValue:   "ReadUint64",
		EncodeValue:   "WriteUint64",
		AppendValue:   "AppendUint64",
		KeyLen:        "size += len(s) * msgp.Uint64Size",
		JSONAppendKey: "strconv.AppendUint(dst, uint64(k), 10)",
		JSONParseKey:  "jsonParseUnsigned[uint64](raw, 64)",
	},
	{
		GoType:        "float64",
		PackageName:   "Float64",
		DecodeValue:   "ReadFloat64",
		EncodeValue:   "WriteFloat",
		AppendValue:   "AppendFloat",
		KeyLen:        "size += len(s) * msgp.Float64Size",
		JSONAppendKey: "strconv.AppendFloat(dst, float64(k), 'f', -1, 64)",
		JSONParseKey:  "jsonParseFloat[float64](raw, 64)",
	},
	{
		GoType:        "float32",
		PackageName:   "Float32",
		DecodeValue:   "ReadFloat32",
		EncodeValue:   "WriteFloat32",
		AppendValue:   "AppendFloat32",
		KeyLen:        "size += len(s) * msgp.Float32Size",
		JSONAppendKey: "strconv.AppendFloat(dst, float64(k), 'f', -1, 32)",
		JSONParseKey:  "jsonParseFloat[float32](raw, 32)",
	},
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		} else {
			// Format generated code
			if err := exec.Command("go", "fmt", "../generated.go").Run(); err != nil {
				fmt.Printf("Warning: failed to format generated.go: %v\n", err)
			}
			if err := exec.Command("go", "fmt", "../generated_test.go").Run(); err != nil {
				fmt.Printf("Warning: failed to format generated_test.go: %v\n", err)
			}
		}
	}()
	out, err := os.Create("../generated.go")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	fmt.Fprintln(out, `// Code generated by ./gen/main.go. DO NOT EDIT.

package setof

import (
	"slices"
	"sort"
	"strconv"

	"github.com/tinylib/msgp/msgp"
)`)

	base := replacers[0]
	for _, r := range replacers[1:] {
		replaced := unsorted + template
		replaced = strings.ReplaceAll(replaced, base.JSONAppendKey, r.JSONAppendKey)
		replaced = strings.ReplaceAll(replaced, base.JSONParseKey, r.JSONParseKey)
		replaced = strings.ReplaceAll(replaced, base.GoType, r.GoType)
		replaced = strings.ReplaceAll(replaced, base.PackageName, r.PackageName)
		replaced = strings.ReplaceAll(replaced, base.EncodeValue, r.EncodeValue)
		replaced = strings.ReplaceAll(replaced, base.DecodeValue, r.DecodeValue)
		replaced = strings.ReplaceAll(replaced, base.AppendValue, r.AppendValue)
		replaced = strings.ReplaceAll(replaced, base.KeyLen, r.KeyLen)

		fmt.Fprintln(out, replaced)

		replaced = sorted + template
		replaced = strings.ReplaceAll(replaced, base.JSONAppendKey, r.JSONAppendKey)
		replaced = strings.ReplaceAll(replaced, base.JSONParseKey, r.JSONParseKey)
		if r.Sorter != "" {
			replaced = strings.ReplaceAll(replaced, base.Sorter, r.Sorter)
		}
		replaced = strings.ReplaceAll(replaced, base.GoType, r.GoType)
		replaced = strings.ReplaceAll(replaced, base.PackageName, r.PackageName+"Sorted")
		replaced = strings.ReplaceAll(replaced, base.EncodeValue, r.EncodeValue)
		replaced = strings.ReplaceAll(replaced, base.DecodeValue, r.DecodeValue)
		replaced = strings.ReplaceAll(replaced, base.AppendValue, r.AppendValue)
		replaced = strings.ReplaceAll(replaced, base.KeyLen, r.KeyLen)

		fmt.Fprintln(out, replaced)
	}

	// Generate test file
	testOut, err := os.Create("../generated_test.go")
	if err != nil {
		panic(err)
	}
	defer testOut.Close()

	fmt.Fprintln(testOut, `// Code generated by ./gen/main.go. DO NOT EDIT.

package setof

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tinylib/msgp/msgp"
)`)

	// Generate tests for each type (skip the template base type)
	for _, r := range replacers[1:] {
		generateTests(testOut, r)
	}
}
