package _generated

import (
	"bytes"
	"io"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func fillErrorCtxAsMap() *ErrorCtxAsMap {
	v := &ErrorCtxAsMap{}
	v.Val = "foo"
	v.Child = &ErrorCtxMapChild{Val: "foo"}
	v.Children = []*ErrorCtxMapChild{{Val: "foo"}, {Val: "bar"}}
	v.Map = map[string]string{"foo": "bar", "baz": "qux"}
	v.Nest.Val = "foo"
	v.Nest.Child = &ErrorCtxMapChild{Val: "foo"}
	v.Nest.Children = []*ErrorCtxMapChild{{Val: "foo"}, {Val: "bar"}}
	v.Nest.Map = map[string]string{"foo": "bar", "baz": "qux"}
	v.Nest.Nest.Val = "foo"
	v.Nest.Nest.Child = &ErrorCtxMapChild{Val: "foo"}
	v.Nest.Nest.Children = []*ErrorCtxMapChild{{Val: "foo"}, {Val: "bar"}}
	v.Nest.Nest.Map = map[string]string{"foo": "bar", "baz": "qux"}
	return v
}

func fillErrorCtxAsTuple() *ErrorCtxAsTuple {
	v := &ErrorCtxAsTuple{}
	v.Val = "foo"
	v.Child = &ErrorCtxTupleChild{Val: "foo"}
	v.Children = []*ErrorCtxTupleChild{{Val: "foo"}, {Val: "bar"}}
	v.Map = map[string]string{"foo": "bar", "baz": "qux"}
	v.Nest.Val = "foo"
	v.Nest.Child = &ErrorCtxTupleChild{Val: "foo"}
	v.Nest.Children = []*ErrorCtxTupleChild{{Val: "foo"}, {Val: "bar"}}
	v.Nest.Map = map[string]string{"foo": "bar", "baz": "qux"}
	v.Nest.Nest.Val = "foo"
	v.Nest.Nest.Child = &ErrorCtxTupleChild{Val: "foo"}
	v.Nest.Nest.Children = []*ErrorCtxTupleChild{{Val: "foo"}, {Val: "bar"}}
	v.Nest.Nest.Map = map[string]string{"foo": "bar", "baz": "qux"}
	return v
}

type outBuf struct {
	*bytes.Buffer
	dodgifyString int
	strIdx        int
}

func (o *outBuf) Write(b []byte) (n int, err error) {
	ilen := len(b)
	if msgp.NextType(b) == msgp.StrType {
		if o.strIdx == o.dodgifyString {
			// Fool msgp into thinking this value is a fixint. msgp will throw
			// a type error for this value.
			b[0] = 1
		}
		o.strIdx++
	}
	_, err = o.Buffer.Write(b)
	return ilen, err
}

type strCounter int

func (o *strCounter) Write(b []byte) (n int, err error) {
	if msgp.NextType(b) == msgp.StrType {
		*o++
	}
	return len(b), nil
}

func countStrings(bts []byte) int {
	r := msgp.NewReader(bytes.NewReader(bts))
	strCounter := strCounter(0)
	for {
		_, err := r.CopyNext(&strCounter)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return int(strCounter)
}

func marshalErrorCtx(m msgp.Marshaler) []byte {
	bts, err := m.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}
	return bts
}

// dodgifyMsgpString will wreck the nth string in the msgpack blob
// so that it raises an error when decoded or unmarshaled.
func dodgifyMsgpString(bts []byte, idx int) []byte {
	r := msgp.NewReader(bytes.NewReader(bts))
	out := &outBuf{Buffer: &bytes.Buffer{}, dodgifyString: idx}
	for {
		_, err := r.CopyNext(out)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return out.Bytes()
}

func TestErrorCtxAsMapUnmarshal(t *testing.T) {
	bts := marshalErrorCtx(fillErrorCtxAsMap())
	cnt := countStrings(bts)

	var as []string
	for i := 0; i < cnt; i++ {
		dodgeBts := dodgifyMsgpString(bts, i)

		var ec ErrorCtxAsMap
		_, err := (&ec).UnmarshalMsg(dodgeBts)
		words := strings.Split(err.Error(), " ")
		last := words[len(words)-1]
		as = append(as, last)
	}

	// Map key iteration order is not consistent
	sort.Strings(expectedAsMap)
	sort.Strings(as)

	if !reflect.DeepEqual(expectedAsMap, as) {
		t.Fatal()
	}
}

func TestErrorCtxAsMapDecode(t *testing.T) {
	bts := marshalErrorCtx(fillErrorCtxAsMap())
	cnt := countStrings(bts)

	var as []string
	for i := 0; i < cnt; i++ {
		dodgeBts := dodgifyMsgpString(bts, i)

		r := msgp.NewReader(bytes.NewReader(dodgeBts))
		var ec ErrorCtxAsMap
		err := (&ec).DecodeMsg(r)
		words := strings.Split(err.Error(), " ")
		last := words[len(words)-1]
		as = append(as, last)
	}

	// Map key iteration order is not consistent
	sort.Strings(expectedAsMap)
	sort.Strings(as)

	if !reflect.DeepEqual(expectedAsMap, as) {
		t.Fatal()
	}
}

func TestErrorCtxAsTupleUnmarshal(t *testing.T) {
	bts := marshalErrorCtx(fillErrorCtxAsTuple())
	cnt := countStrings(bts)

	var as []string
	for i := 0; i < cnt; i++ {
		dodgeBts := dodgifyMsgpString(bts, i)

		var ec ErrorCtxAsTuple
		_, err := (&ec).UnmarshalMsg(dodgeBts)
		words := strings.Split(err.Error(), " ")
		last := words[len(words)-1]
		as = append(as, last)
	}

	// Map key iteration order is not consistent
	sort.Strings(expectedAsTuple)
	sort.Strings(as)

	if !reflect.DeepEqual(expectedAsTuple, as) {
		t.Fatal()
	}
}

func TestErrorCtxAsTupleDecode(t *testing.T) {
	bts := marshalErrorCtx(fillErrorCtxAsTuple())
	cnt := countStrings(bts)

	var as []string
	for i := 0; i < cnt; i++ {
		dodgeBts := dodgifyMsgpString(bts, i)

		r := msgp.NewReader(bytes.NewReader(dodgeBts))
		var ec ErrorCtxAsTuple
		err := (&ec).DecodeMsg(r)
		words := strings.Split(err.Error(), " ")
		last := words[len(words)-1]
		as = append(as, last)
	}

	// Map key iteration order is not consistent
	sort.Strings(expectedAsTuple)
	sort.Strings(as)

	if !reflect.DeepEqual(expectedAsTuple, as) {
		t.Fatal()
	}
}

var expectedAsTuple = []string{
	"ErrorCtxAsTuple/Val",
	"ErrorCtxAsTuple/Child/Val",
	"ErrorCtxAsTuple/Children/0/Val",
	"ErrorCtxAsTuple/Children/1/Val",
	"ErrorCtxAsTuple/Map",
	"ErrorCtxAsTuple/Map/baz",
	"ErrorCtxAsTuple/Map",
	"ErrorCtxAsTuple/Map/foo",
	"ErrorCtxAsTuple/Nest",
	"ErrorCtxAsTuple/Nest/Val",
	"ErrorCtxAsTuple/Nest",
	"ErrorCtxAsTuple/Nest/Child/Val",
	"ErrorCtxAsTuple/Nest",
	"ErrorCtxAsTuple/Nest/Children/0/Val",
	"ErrorCtxAsTuple/Nest/Children/1/Val",
	"ErrorCtxAsTuple/Nest",
	"ErrorCtxAsTuple/Nest/Map",
	"ErrorCtxAsTuple/Nest/Map/foo",
	"ErrorCtxAsTuple/Nest/Map",
	"ErrorCtxAsTuple/Nest/Map/baz",
	"ErrorCtxAsTuple/Nest",
	"ErrorCtxAsTuple/Nest/Nest",
	"ErrorCtxAsTuple/Nest/Nest/Val",
	"ErrorCtxAsTuple/Nest/Nest",
	"ErrorCtxAsTuple/Nest/Nest/Child/Val",
	"ErrorCtxAsTuple/Nest/Nest",
	"ErrorCtxAsTuple/Nest/Nest/Children/0/Val",
	"ErrorCtxAsTuple/Nest/Nest/Children/1/Val",
	"ErrorCtxAsTuple/Nest/Nest",
	"ErrorCtxAsTuple/Nest/Nest/Map",
	"ErrorCtxAsTuple/Nest/Nest/Map/foo",
	"ErrorCtxAsTuple/Nest/Nest/Map",
	"ErrorCtxAsTuple/Nest/Nest/Map/baz",
}

// there are a lot of extra errors in here at the struct level because we are
// not discriminating between dodgy struct field map key strings and
// values. dodgy struct field map keys have no field context available when
// they are read.
var expectedAsMap = []string{
	"ErrorCtxAsMap",
	"ErrorCtxAsMap/Val",
	"ErrorCtxAsMap",
	"ErrorCtxAsMap/Child",
	"ErrorCtxAsMap/Child/Val",
	"ErrorCtxAsMap",
	"ErrorCtxAsMap/Children/0",
	"ErrorCtxAsMap/Children/0/Val",
	"ErrorCtxAsMap/Children/1",
	"ErrorCtxAsMap/Children/1/Val",
	"ErrorCtxAsMap",
	"ErrorCtxAsMap/Map",
	"ErrorCtxAsMap/Map/foo",
	"ErrorCtxAsMap/Map",
	"ErrorCtxAsMap/Map/baz",
	"ErrorCtxAsMap",
	"ErrorCtxAsMap/Nest",
	"ErrorCtxAsMap/Nest/Val",
	"ErrorCtxAsMap/Nest",
	"ErrorCtxAsMap/Nest/Child",
	"ErrorCtxAsMap/Nest/Child/Val",
	"ErrorCtxAsMap/Nest",
	"ErrorCtxAsMap/Nest/Children/0",
	"ErrorCtxAsMap/Nest/Children/0/Val",
	"ErrorCtxAsMap/Nest/Children/1",
	"ErrorCtxAsMap/Nest/Children/1/Val",
	"ErrorCtxAsMap/Nest",
	"ErrorCtxAsMap/Nest/Map",
	"ErrorCtxAsMap/Nest/Map/foo",
	"ErrorCtxAsMap/Nest/Map",
	"ErrorCtxAsMap/Nest/Map/baz",
	"ErrorCtxAsMap/Nest",
	"ErrorCtxAsMap/Nest/Nest",
	"ErrorCtxAsMap/Nest/Nest/Val",
	"ErrorCtxAsMap/Nest/Nest",
	"ErrorCtxAsMap/Nest/Nest/Child",
	"ErrorCtxAsMap/Nest/Nest/Child/Val",
	"ErrorCtxAsMap/Nest/Nest",
	"ErrorCtxAsMap/Nest/Nest/Children/0",
	"ErrorCtxAsMap/Nest/Nest/Children/0/Val",
	"ErrorCtxAsMap/Nest/Nest/Children/1",
	"ErrorCtxAsMap/Nest/Nest/Children/1/Val",
	"ErrorCtxAsMap/Nest/Nest",
	"ErrorCtxAsMap/Nest/Nest/Map",
	"ErrorCtxAsMap/Nest/Nest/Map/baz",
	"ErrorCtxAsMap/Nest/Nest/Map",
	"ErrorCtxAsMap/Nest/Nest/Map/foo",
}
