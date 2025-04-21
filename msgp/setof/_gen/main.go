package main

import (
	"fmt"
	"os"
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
		for k := range dst {
			delete(dst, k)
		}
	}
	for i := uint32(0); i < sz; i++ {
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
		for k := range dst {
			delete(dst, k)
		}
	}
	for i := uint32(0); i < sz; i++ {
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
	if len(s) == 0 {
		return msgp.ArrayHeaderSize
	}
	size := msgp.ArrayHeaderSize
	size += len(s) * msgp.StringPrefixSize
	return size
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
	bytes = msgp.AppendArrayHeader(bytes, uint32(len(s)))
	for k := range s {
		bytes = msgp.AppendString(bytes, string(k))
	}
	return bytes, nil
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
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

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
	bytes = msgp.AppendArrayHeader(bytes, uint32(len(s)))
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		bytes = msgp.AppendString(bytes, k)
	}
	return bytes, nil
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
}

var replacers = []replacer{
	{
		GoType:      "string",
		PackageName: "Foo",
		DecodeValue: "ReadString",
		EncodeValue: "WriteString",
		AppendValue: "AppendString",
		KeyLen:      "size += len(s) * msgp.StringPrefixSize",
		Sorter:      "sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })",
	},
	{
		GoType:      "string",
		PackageName: "String",
		DecodeValue: "ReadString",
		EncodeValue: "WriteString",
		AppendValue: "AppendString",
		KeyLen:      "for key := range s {\n\t\t\tsize += msgp.StringPrefixSize + len(key)\n\t\t}",
		Sorter:      "sort.Strings(keys)",
	},
	{
		GoType:      "int",
		PackageName: "Int",
		DecodeValue: "ReadInt",
		EncodeValue: "WriteInt",
		AppendValue: "AppendInt",
		KeyLen:      "size += len(s) * msgp.IntSize",
	},
	{
		GoType:      "uint",
		PackageName: "Uint",
		DecodeValue: "ReadUint",
		EncodeValue: "WriteUint",
		AppendValue: "AppendUint",
		KeyLen:      "size += len(s) * msgp.UintSize",
	},
	{
		GoType:      "byte",
		PackageName: "Byte",
		DecodeValue: "ReadByte",
		EncodeValue: "WriteByte",
		AppendValue: "AppendByte",
		KeyLen:      "size += len(s) * msgp.ByteSize",
	},
	{
		GoType:      "int8",
		PackageName: "Int8",
		DecodeValue: "ReadInt8",
		EncodeValue: "WriteInt8",
		AppendValue: "AppendInt8",
		KeyLen:      "size += len(s) * msgp.Int8Size",
	},
	{
		GoType:      "uint8",
		PackageName: "Uint8",
		DecodeValue: "ReadUint8",
		EncodeValue: "WriteUint8",
		AppendValue: "AppendUint8",
		KeyLen:      "size += len(s) * msgp.Uint8Size",
	},
	{
		GoType:      "int16",
		PackageName: "Int16",
		DecodeValue: "ReadInt16",
		EncodeValue: "WriteInt16",
		AppendValue: "AppendInt16",
		KeyLen:      "size += len(s) * msgp.Int16Size",
	},
	{
		GoType:      "uint16",
		PackageName: "Uint16",
		DecodeValue: "ReadUint16",
		EncodeValue: "WriteUint16",
		AppendValue: "AppendUint16",
		KeyLen:      "size += len(s) * msgp.Uint16Size",
	},
	{
		GoType:      "int32",
		PackageName: "Int32",
		DecodeValue: "ReadInt32",
		EncodeValue: "WriteInt32",
		AppendValue: "AppendInt32",
		KeyLen:      "size += len(s) * msgp.Int32Size",
	},
	{
		GoType:      "uint32",
		PackageName: "Uint32",
		DecodeValue: "ReadUint32",
		EncodeValue: "WriteUint32",
		AppendValue: "AppendUint32",
		KeyLen:      "size += len(s) * msgp.Uint32Size",
	},
	{
		GoType:      "int64",
		PackageName: "Int64",
		DecodeValue: "ReadInt64",
		EncodeValue: "WriteInt64",
		AppendValue: "AppendInt64",
		KeyLen:      "size += len(s) * msgp.Int64Size",
	},
	{
		GoType:      "uint64",
		PackageName: "Uint64",
		DecodeValue: "ReadUint64",
		EncodeValue: "WriteUint64",
		AppendValue: "AppendUint64",
		KeyLen:      "size += len(s) * msgp.Uint64Size",
	},
	{
		GoType:      "float64",
		PackageName: "Float64",
		DecodeValue: "ReadFloat64",
		EncodeValue: "WriteFloat",
		AppendValue: "AppendFloat",
		KeyLen:      "size += len(s) * msgp.Float64Size",
	},
	{
		GoType:      "float32",
		PackageName: "Float32",
		DecodeValue: "ReadFloat32",
		EncodeValue: "WriteFloat32",
		AppendValue: "AppendFloat32",
		KeyLen:      "size += len(s) * msgp.Float32Size",
	},
}

func main() {
	out, err := os.Create("../generated.go")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	fmt.Fprintln(out, `// Code generated by ./gen/main.go. DO NOT EDIT.

package setof

import (
	"sort"

	"github.com/tinylib/msgp/msgp"
)`)

	base := replacers[0]
	for _, r := range replacers[1:] {
		replaced := unsorted + template
		replaced = strings.ReplaceAll(replaced, base.GoType, r.GoType)
		replaced = strings.ReplaceAll(replaced, base.PackageName, r.PackageName)
		replaced = strings.ReplaceAll(replaced, base.EncodeValue, r.EncodeValue)
		replaced = strings.ReplaceAll(replaced, base.DecodeValue, r.DecodeValue)
		replaced = strings.ReplaceAll(replaced, base.AppendValue, r.AppendValue)
		replaced = strings.ReplaceAll(replaced, base.KeyLen, r.KeyLen)

		fmt.Fprintln(out, replaced)

		replaced = sorted + template
		replaced = strings.ReplaceAll(replaced, base.GoType, r.GoType)
		replaced = strings.ReplaceAll(replaced, base.PackageName, r.PackageName+"Sorted")
		replaced = strings.ReplaceAll(replaced, base.EncodeValue, r.EncodeValue)
		replaced = strings.ReplaceAll(replaced, base.DecodeValue, r.DecodeValue)
		replaced = strings.ReplaceAll(replaced, base.AppendValue, r.AppendValue)
		replaced = strings.ReplaceAll(replaced, base.KeyLen, r.KeyLen)
		if r.Sorter != "" {
			replaced = strings.ReplaceAll(replaced, base.Sorter, r.Sorter)
		}

		fmt.Fprintln(out, replaced)
	}
}
