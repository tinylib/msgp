// Package setof allows serializing sets map[T]struct{} as arrays.
//
// Nil maps are preserved as a nil value on stream.
//
// A deterministic, sorted version is available, with slightly lower performance.

package setof
