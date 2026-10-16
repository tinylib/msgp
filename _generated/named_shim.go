package _generated

//go:generate msgp

//msgp:shim *NamedShimOriginal as:NamedShimWire using:namedShimToWire/NamedShimFromWire
//msgp:ignore NamedShimOriginal

// The original type has promoted string serializers. The shim must use the
// generated map serializer instead, including when decoding into a zero value.
type NamedShimOriginal struct {
	NamedShimString
}

type NamedShimString string

type NamedShimWire struct {
	Value *string
}

func namedShimToWire(v *NamedShimOriginal) *NamedShimWire {
	return &NamedShimWire{Value: (*string)(&v.NamedShimString)}
}

// NamedShimFromWire is the inverse of namedShimToWire.
func NamedShimFromWire(v *NamedShimWire) *NamedShimOriginal {
	return &NamedShimOriginal{NamedShimString: NamedShimString(*v.Value)}
}

type NamedShimContainer struct {
	Value   NamedShimOriginal
	Pointer *NamedShimOriginal
	Plain   NamedShimString
	Slice   []NamedShimOriginal
	Array   [1]NamedShimOriginal
	Map     map[string]NamedShimOriginal
}
