package msgp

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"runtime"
	"strconv"
	"testing"
)

// Fuzz Reader tests the Reader methods with fuzzing inputs.
func FuzzReader(f *testing.F) {
	addFuzzDataFromZip(f, "testdata/FuzzCorpus.zip")

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 || len(data) > 1<<20 {
			return
		}
		var tmp5 [5]byte
		r := NewReader(bytes.NewReader(data))
		// Set some rather low limits for the reader to avoid excessive memory usage.
		r.SetMaxRecursionDepth(1000)
		r.SetMaxElements(10000)
		r.SetMaxStringLength(1000)

		reset := func() {
			t.Helper()
			// Reset the test state if needed
			r.Reset(bytes.NewReader(data))
		}
		// Uses reset when data is read off the stream.
		r.BufferSize()
		r.Buffered()
		r.CopyNext(io.Discard)
		reset()
		r.IsNil()
		r.NextType()
		r.Read(tmp5[:])
		reset()
		r.ReadArrayHeader()
		reset()
		r.ReadBool()
		reset()
		r.ReadByte()
		reset()
		r.ReadBytes(tmp5[:])
		reset()
		r.ReadComplex128()
		reset()
		r.ReadComplex64()
		reset()
		r.ReadDuration()
		reset()
		r.ReadExactBytes(tmp5[:])
		reset()
		r.ReadExtension(testExt{})
		reset()
		r.ReadExtensionRaw()
		reset()
		r.ReadFloat32()
		reset()
		r.ReadFloat64()
		reset()
		r.ReadFull(tmp5[:])
		reset()
		r.ReadInt()
		reset()
		r.ReadInt16()
		reset()
		r.ReadInt32()
		reset()
		r.ReadInt64()
		reset()
		r.ReadInt8()
		reset()
		r.ReadIntf()
		reset()
		runtime.GC()
		r.ReadJSONNumber()
		reset()
		r.ReadMapHeader()
		reset()
		r.ReadMapKey(tmp5[:])
		reset()
		runtime.GC()
		r.ReadMapKeyPtr()
		reset()
		r.ReadMapStrIntf(map[string]any{})
		reset()
		r.ReadNil()
		reset()
		r.ReadString()
		reset()
		r.ReadStringAsBytes(tmp5[:])
		reset()
		r.ReadStringHeader()
		reset()
		r.ReadTime()
		reset()
		r.ReadTimeUTC()
		reset()
		r.ReadUint()
		reset()
		r.ReadUint16()
		reset()
		r.ReadUint32()
		reset()
		r.ReadUint64()
		reset()
		r.ReadUint64()
		reset()
		r.ReadUint8()
		reset()
		r.ReadBinaryUnmarshal(encodingReader{})
		reset()
		r.ReadTextUnmarshal(encodingReader{})
		reset()
		r.ReadTextUnmarshalString(encodingReader{})
		reset()
		for range ReadArray(r, r.ReadInt) {
		}
		reset()
		iter, errfn := ReadMap(r, r.ReadString, r.ReadUint64)
		for range iter {
		}
		errfn()
		reset()

		r.Skip()
		reset()

		r.WriteToJSON(io.Discard)
	})
}

// Test all byte reading functions with fuzzing inputs.
// Since we are reading bytes memory allocations should be well controlled.
func FuzzReadBytes(f *testing.F) {
	addFuzzDataFromZip(f, "testdata/FuzzCorpus.zip")

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 || len(data) > 1<<20 {
			return
		}
		var tmp5 [5]byte

		ReadArrayHeaderBytes(data)
		ReadBoolBytes(data)
		ReadByteBytes(data)
		ReadBytesBytes(data, tmp5[:])
		ReadBytesZC(data)
		ReadComplex128Bytes(data)
		ReadComplex64Bytes(data)
		ReadDurationBytes(data)
		ReadExactBytes(data, tmp5[:])
		ReadExtensionBytes(data, testExt{})
		ReadFloat32Bytes(data)
		ReadFloat64Bytes(data)
		ReadInt8Bytes(data)
		ReadInt16Bytes(data)
		ReadInt32Bytes(data)
		ReadInt64Bytes(data)
		ReadIntBytes(data)
		ReadIntfBytes(data)
		ReadJSONNumberBytes(data)
		ReadMapHeaderBytes(data)
		ReadMapKeyZC(data)
		ReadMapStrIntfBytes(data, map[string]any{})
		ReadNilBytes(data)
		ReadStringBytes(data)
		ReadStringAsBytes(data, tmp5[:])
		ReadStringBytes(data)
		ReadStringZC(data)
		ReadTimeBytes(data)
		ReadTimeUTCBytes(data)
		ReadUint8Bytes(data)
		ReadUint16Bytes(data)
		ReadUint32Bytes(data)
		ReadUint64Bytes(data)
		ReadUintBytes(data)
		UnmarshalAsJSON(io.Discard, data)
		iter, errfn := ReadArrayBytes(data, ReadIntBytes)
		for range iter {
		}
		errfn()
		iter2, errfn2 := ReadMapBytes(data, ReadStringBytes, ReadUint64Bytes)
		for range iter2 {
		}
		errfn2()
		Skip(data)
	})
}

type testExt struct{}

func (e testExt) ExtensionType() int8 {
	return 42
}

func (e testExt) Len() int {
	return 3
}

var marshalExtContent = bytes.Repeat([]byte("x"), 3)

func (e testExt) MarshalBinaryTo(i []byte) error {
	if len(i) < e.Len() {
		return io.ErrShortBuffer
	}
	copy(i, marshalExtContent)
	return nil
}

func (e testExt) UnmarshalBinary(i []byte) error {
	if len(i) < e.Len() {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// addFuzzDataFromZip will read the supplied zip and add all as corpus for f.
// Byte slices only.
func addFuzzDataFromZip(f *testing.F, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		f.Fatal(err)
	}
	fi, err := file.Stat()
	if fi == nil {
		return
	}

	if err != nil {
		f.Fatal(err)
	}
	zr, err := zip.NewReader(file, fi.Size())
	if err != nil {
		f.Fatal(err)
	}
	for _, file := range zr.File {
		rc, err := file.Open()
		if err != nil {
			f.Fatal(err)
		}

		b, err := io.ReadAll(rc)
		if err != nil {
			f.Fatal(err)
		}
		rc.Close()

		if !bytes.HasPrefix(b, []byte("go test fuzz")) {
			continue
		}

		vals, err := unmarshalCorpusFile(b)
		if err != nil {
			f.Fatal(err)
		}
		for _, v := range vals {
			f.Add(v)
		}
	}
}

// unmarshalCorpusFile decodes corpus bytes into their respective values.
func unmarshalCorpusFile(b []byte) ([][]byte, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("cannot unmarshal empty string")
	}
	lines := bytes.Split(b, []byte("\n"))
	if len(lines) < 2 {
		return nil, fmt.Errorf("must include version and at least one value")
	}
	var vals = make([][]byte, 0, len(lines)-1)
	for _, line := range lines[1:] {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		v, err := parseCorpusValue(line)
		if err != nil {
			return nil, fmt.Errorf("malformed line %q: %v", line, err)
		}
		vals = append(vals, v)
	}
	return vals, nil
}

// parseCorpusValue
func parseCorpusValue(line []byte) ([]byte, error) {
	fs := token.NewFileSet()
	expr, err := parser.ParseExprFrom(fs, "(test)", line, 0)
	if err != nil {
		return nil, err
	}
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil, fmt.Errorf("expected call expression")
	}
	if len(call.Args) != 1 {
		return nil, fmt.Errorf("expected call expression with 1 argument; got %d", len(call.Args))
	}
	arg := call.Args[0]

	if arrayType, ok := call.Fun.(*ast.ArrayType); ok {
		if arrayType.Len != nil {
			return nil, fmt.Errorf("expected []byte or primitive type")
		}
		elt, ok := arrayType.Elt.(*ast.Ident)
		if !ok || elt.Name != "byte" {
			return nil, fmt.Errorf("expected []byte")
		}
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return nil, fmt.Errorf("string literal required for type []byte")
		}
		s, err := strconv.Unquote(lit.Value)
		if err != nil {
			return nil, err
		}
		return []byte(s), nil
	}
	return nil, fmt.Errorf("expected []byte")
}

type encodingReader struct{}

func (e encodingReader) UnmarshalBinary(data []byte) error {
	return nil
}

func (e encodingReader) UnmarshalText(data []byte) error {
	return nil
}
