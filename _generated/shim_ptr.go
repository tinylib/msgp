package _generated

import "errors"

//go:generate msgp

//msgp:shim *ShimPtrStr as:string using:shimPtrStrToString/stringToShimPtrStr
//msgp:shim *ShimPtrBytes as:[]byte using:shimPtrBytesToBytes/bytesToShimPtrBytes
//msgp:shim *ShimPtrErr as:string using:shimPtrErrToString/stringToShimPtrErr witherr:true
//msgp:ignore ShimPtrBytes

// Shims of *T take and return a pointer. Fields of type T and *T use them.
type ShimPtrStr struct{ A string }

func shimPtrStrToString(t *ShimPtrStr) string { return t.A }
func stringToShimPtrStr(s string) *ShimPtrStr {
	if s == "" {
		return nil // decodes as the zero value
	}
	return &ShimPtrStr{A: s}
}

type ShimPtrBytes struct{ B string }

func shimPtrBytesToBytes(t *ShimPtrBytes) []byte { return []byte(t.B) }
func bytesToShimPtrBytes(b []byte) *ShimPtrBytes { return &ShimPtrBytes{B: string(b)} }

type ShimPtrErr struct{ C string }

func shimPtrErrToString(t *ShimPtrErr) string { return t.C }
func stringToShimPtrErr(s string) (*ShimPtrErr, error) {
	if s == "bad" {
		return nil, errors.New("bad")
	}
	return &ShimPtrErr{C: s}, nil
}

type ShimPtr struct {
	S  ShimPtrStr
	B  ShimPtrBytes
	E  ShimPtrErr
	MS map[string]ShimPtrStr
	PS *ShimPtrStr
	PE *ShimPtrErr
}
