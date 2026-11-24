package _generated

//go:generate msgp

type PtrPtrInner struct{ A int }

type PtrPtr struct {
	F   **PtrPtrInner
	FFF ***PtrPtrInner
	M   map[string]**PtrPtrInner
	S   []**PtrPtrInner
	E   **PtrPtrExt // has its own methods
	ME  map[string]**PtrPtrExt
}
