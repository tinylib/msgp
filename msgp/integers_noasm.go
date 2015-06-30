//+build noasm

package msgp

import (
	"math"
)

func getUnix(b []byte) (sec int64, nsec int32) {
	sec = (int64(b[0]) << 56) | (int64(b[1]) << 48) |
		(int64(b[2]) << 40) | (int64(b[3]) << 32) |
		(int64(b[4]) << 24) | (int64(b[5]) << 16) |
		(int64(b[6]) << 8) | (int64(b[7]))
	nsec = (int32(b[8]) << 24) | (int32(b[9]) << 16) | (int32(b[10]) << 8) | (int32(b[11]))
	return
}

func putUnix(b []byte, sec int64, nsec int32) {
	b[0] = byte(sec >> 56)
	b[1] = byte(sec >> 48)
	b[2] = byte(sec >> 40)
	b[3] = byte(sec >> 32)
	b[4] = byte(sec >> 24)
	b[5] = byte(sec >> 16)
	b[6] = byte(sec >> 8)
	b[7] = byte(sec)
	b[8] = byte(nsec >> 24)
	b[9] = byte(nsec >> 16)
	b[10] = byte(nsec >> 8)
	b[11] = byte(nsec)
}

func putMint64(b []byte, i int64) {
	b[0] = mint64
	b[1] = byte(i >> 56)
	b[2] = byte(i >> 48)
	b[3] = byte(i >> 40)
	b[4] = byte(i >> 32)
	b[5] = byte(i >> 24)
	b[6] = byte(i >> 16)
	b[7] = byte(i >> 8)
	b[8] = byte(i)
}

func getMint64(b []byte) int64 {
	return (int64(b[1]) << 56) | (int64(b[2]) << 48) |
		(int64(b[3]) << 40) | (int64(b[4]) << 32) |
		(int64(b[5]) << 24) | (int64(b[6]) << 16) |
		(int64(b[7]) << 8) | (int64(b[8]))
}

func putMint32(b []byte, i int32) {
	b[0] = mint32
	b[1] = byte(i >> 24)
	b[2] = byte(i >> 16)
	b[3] = byte(i >> 8)
	b[4] = byte(i)
}

func getMint32(b []byte) int32 {
	return (int32(b[1]) << 24) | (int32(b[2]) << 16) | (int32(b[3]) << 8) | (int32(b[4]))
}

func putMuint64(b []byte, u uint64) {
	b[0] = muint64
	b[1] = byte(u >> 56)
	b[2] = byte(u >> 48)
	b[3] = byte(u >> 40)
	b[4] = byte(u >> 32)
	b[5] = byte(u >> 24)
	b[6] = byte(u >> 16)
	b[7] = byte(u >> 8)
	b[8] = byte(u)
}

func getMuint64(b []byte) uint64 {
	return (uint64(b[1]) << 56) | (uint64(b[2]) << 48) |
		(uint64(b[3]) << 40) | (uint64(b[4]) << 32) |
		(uint64(b[5]) << 24) | (uint64(b[6]) << 16) |
		(uint64(b[7]) << 8) | (uint64(b[8]))
}

func prefixu64(b []byte, sz uint64, pre byte) {
	b[0] = pre
	b[1] = byte(sz >> 56)
	b[2] = byte(sz >> 48)
	b[3] = byte(sz >> 40)
	b[4] = byte(sz >> 32)
	b[5] = byte(sz >> 24)
	b[6] = byte(sz >> 16)
	b[7] = byte(sz >> 8)
	b[8] = byte(sz)
}

func prefixu32(b []byte, sz uint32, pre byte) {
	b[0] = pre
	b[1] = byte(sz >> 24)
	b[2] = byte(sz >> 16)
	b[3] = byte(sz >> 8)
	b[4] = byte(sz)
}

func put264(p []byte, a float64, b float64) {
	big.PutUint64(p, math.Float64bits(a))
	big.PutUint64(p[8:], math.Float64bits(b))
}

func put232(p []byte, a float32, b float32) {
	big.PutUint32(p, math.Float32bits(a))
	big.PutUint32(p[4:], math.Float32bits(b))
}
