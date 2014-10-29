package enc

// The sizes provided
// are the worst-case
// encoded sizes for
// each type. For variable-
// length types ([]byte, string),
// the total encoded size is
// the prefix size plus the
// length of the object.
const (
	Int64Size      = 9
	IntSize        = Int64Size
	UintSize       = IntSize
	Int8Size       = IntSize
	Int16Size      = IntSize
	Int32Size      = IntSize
	Uint8Size      = IntSize
	Uint16Size     = IntSize
	Uint32Size     = IntSize
	Uint64Size     = IntSize
	Float64Size    = 9
	Float32Size    = 5
	Complex64Size  = 10
	Complex128Size = 18

	TimeSize = 18
	BoolSize = 1
	NilSize  = 1

	MapHeaderSize   = 5
	ArrayHeaderSize = 5

	BytesPrefixSize  = 5
	StringPrefixSize = 5
)
