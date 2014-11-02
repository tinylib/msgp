package enc

import (
	"errors"
)

var (
	ErrFieldNotFound = errors.New("field not found")
)

// Locate returns a []byte pointing to the field
// in a messagepack map with the provided key. (The returned []byte
// points to a sub-slice of 'raw'; Locate does no allocations.) If the
// key doesn't exist in the map, a zero-length []byte will be returned.
func Locate(key string, raw []byte) []byte {
	s, n := locate(raw, key)
	return raw[s:n]
}

// Replace takes a key ("key") in a messagepack map ("raw")
// and replaces its value with the one provided and returns
// the new []byte. The returned []byte may point to the same
// memory as "raw". Replace makes no effort to evaluate the validity
// of the contents of 'val'. It may use up to the full capacity of 'raw.'
func Replace(key string, raw []byte, val []byte) ([]byte, error) {
	start, end := locate(raw, key)
	if start == end {
		return nil, ErrFieldNotFound
	}
	return replace(raw, start, end, val, true), nil
}

// CopyReplace works similarly to Replace except that the returned
// byte slice does not point to the same memory as 'raw'.
func CopyReplace(key string, raw []byte, val []byte) ([]byte, error) {
	start, end := locate(raw, key)
	if start == end {
		return nil, ErrFieldNotFound
	}
	return replace(raw, start, end, val, false), nil
}

func replace(raw []byte, start int, end int, val []byte, inplace bool) []byte {
	ll := end - start // length of segment to replace
	lv := len(val)

	if inplace {
		extra := lv - ll

		// fastest case: we're doing
		// a 1:1 replacement
		if extra == 0 {
			copy(raw[start:], val)
			return raw

		} else if extra < 0 {
			// 'val' smaller than replaced value
			// copy in place and shift back

			x := copy(raw[start:], val)
			y := copy(raw[start+x:], raw[end:])
			return raw[:start+x+y]

		} else if extra < cap(raw)-len(raw) {
			// 'val' less than (cap-len) extra bytes
			// copy in place and shift forward
			raw = raw[0 : len(raw)+extra]
			// shift end forward
			copy(raw[end+extra:], raw[end:])
			copy(raw[start:], val)
			return raw
		}
	}

	// we have to allocate new space
	out := make([]byte, len(raw)+len(val)-ll)
	x := copy(out, raw[:start])
	y := copy(out[x:], val)
	copy(out[x+y:], raw[end:])
	return out
}

// locate does a naive O(n) search for the map key; returns start, end
// (returns 0,0 on error)
func locate(raw []byte, key string) (start int, end int) {
	var (
		sz    uint32
		bts   []byte
		field []byte
		err   error
	)
	sz, bts, err = ReadMapHeaderBytes(raw)
	if err != nil {
		return
	}

	// loop and locate field
	for i := uint32(0); i < sz; i++ {
		field, bts, err = ReadStringZC(bts)
		if err != nil {
			return
		}
		if UnsafeString(field) == key {
			// start location
			start = len(raw) - len(bts)
			bts, err = Skip(bts)
			if err != nil {
				return
			}
			end = len(raw) - len(bts)
			return
		}
		bts, err = Skip(bts)
		if err != nil {
			return
		}
	}
	return
}
