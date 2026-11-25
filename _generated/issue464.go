package _generated

//go:generate msgp

type Issue464 struct {
	Bytes      *[]byte            `msg:",zerocopy"`
	BytesPP    **[]byte           `msg:",zerocopy"`
	SlicePtr   []*[]byte          `msg:",zerocopy"`
	MapPtr     map[string]*[]byte `msg:",zerocopy"`
	NotZC      *[]byte
	LimitBytes *[]byte `msg:",zerocopy,limit=10"`
}
