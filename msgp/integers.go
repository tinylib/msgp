package msgp

import "unsafe"

/* -----------------------------
	integer encoding utilities
	(inline-able)

	TODO(philhofer): someday, when
	inline assembler is supported,
	we can use, for example
	__asm {
		movq   eax, i
		bswapq eax
	}
   ----------------------------- */

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

func getMint64(b []byte) (i int64) {
	i |= int64(b[1]) << 56
	i |= int64(b[2]) << 48
	i |= int64(b[3]) << 40
	i |= int64(b[4]) << 32
	i |= int64(b[5]) << 24
	i |= int64(b[6]) << 16
	i |= int64(b[7]) << 8
	i |= int64(b[8])
	return
}

func putMint32(b []byte, i int32) {
	b[0] = mint32
	b[1] = byte(i >> 24)
	b[2] = byte(i >> 16)
	b[3] = byte(i >> 8)
	b[4] = byte(i)
}

func getMint32(b []byte) (i int32) {
	i |= int32(b[1]) << 24
	i |= int32(b[2]) << 16
	i |= int32(b[3]) << 8
	i |= int32(b[4])
	return
}

func putMint16(b []byte, i int16) {
	b[0] = mint16
	b[1] = byte(i >> 8)
	b[2] = byte(i)
}

func getMint16(b []byte) (i int16) {
	i |= int16(b[1]) << 8
	i |= int16(b[2])
	return
}

func putMint8(b []byte, i int8) {
	b[0] = mint8
	b[1] = byte(i)
}

func getMint8(b []byte) (i int8) {
	i = int8(b[1])
	return
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

func getMuint64(b []byte) (u uint64) {
	u |= uint64(b[1]) << 56
	u |= uint64(b[2]) << 48
	u |= uint64(b[3]) << 40
	u |= uint64(b[4]) << 32
	u |= uint64(b[5]) << 24
	u |= uint64(b[6]) << 16
	u |= uint64(b[7]) << 8
	u |= uint64(b[8])
	return
}

func putMuint32(b []byte, u uint32) {
	b[0] = muint32
	b[1] = byte(u >> 24)
	b[2] = byte(u >> 16)
	b[3] = byte(u >> 8)
	b[4] = byte(u)
}

func getMuint32(b []byte) (u uint32) {
	u |= uint32(b[1]) << 24
	u |= uint32(b[2]) << 16
	u |= uint32(b[3]) << 8
	u |= uint32(b[4])
	return
}

func putMuint16(b []byte, u uint16) {
	b[0] = muint16
	b[1] = byte(u >> 8)
	b[2] = byte(u)
}

func getMuint16(b []byte) (u uint16) {
	u |= uint16(b[1]) << 8
	u |= uint16(b[2])
	return
}

func putMuint8(b []byte, u uint8) {
	b[0] = muint8
	b[1] = byte(u)
}

func getMuint8(b []byte) (u uint8) {
	u = uint8(b[1])
	return
}

func getUnix(b []byte) (sec int64, nsec int32) {
	sec |= int64(b[0]) << 56
	sec |= int64(b[1]) << 48
	sec |= int64(b[2]) << 40
	sec |= int64(b[3]) << 32
	sec |= int64(b[4]) << 24
	sec |= int64(b[5]) << 16
	sec |= int64(b[6]) << 8
	sec |= int64(b[7])
	nsec |= int32(b[8]) << 24
	nsec |= int32(b[9]) << 16
	nsec |= int32(b[10]) << 8
	nsec |= int32(b[11])
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

/* -----------------------------
		prefix utilities
   ----------------------------- */

// write prefix and uint8
func prefixu8(b []byte, pre byte, sz uint8) {
	b[0] = pre
	b[1] = byte(sz)
}

// write prefix and big-endian uint16
func prefixu16(b []byte, pre byte, sz uint16) {
	b[0] = pre
	b[1] = byte(sz >> 8)
	b[2] = byte(sz)
}

// write prefix and big-endian uint32
func prefixu32(b []byte, pre byte, sz uint32) {
	b[0] = pre
	b[1] = byte(sz >> 24)
	b[2] = byte(sz >> 16)
	b[3] = byte(sz >> 8)
	b[4] = byte(sz)
}

/* ---------------------------
	memory-copying utilities
	WARNING: gross code ahead
   --------------------------- */

// copy 4 unaligned bytes w/o branch
func memcpy4(to unsafe.Pointer, from unsafe.Pointer) {
	s, f := (*[4]byte)(to), (*[4]byte)(from)
	s[0] = f[0]
	s[1] = f[1]
	s[2] = f[2]
	s[3] = f[3]
}

// copy 8 unaligned bytes w/o branch
func memcpy8(to unsafe.Pointer, from unsafe.Pointer) {
	s, f := (*[8]byte)(to), (*[8]byte)(from)
	s[0] = f[0]
	s[1] = f[1]
	s[2] = f[2]
	s[3] = f[3]
	s[4] = f[4]
	s[5] = f[5]
	s[6] = f[6]
	s[7] = f[7]
}

// copy 16 unaligned bytes w/o branch
func memcpy16(to unsafe.Pointer, from unsafe.Pointer) {
	s, f := (*[16]byte)(to), (*[16]byte)(from)
	s[0] = f[0]
	s[1] = f[1]
	s[2] = f[2]
	s[3] = f[3]
	s[4] = f[4]
	s[5] = f[5]
	s[6] = f[6]
	s[7] = f[7]
	s[8] = f[8]
	s[9] = f[9]
	s[10] = f[10]
	s[11] = f[11]
	s[12] = f[12]
	s[13] = f[13]
	s[14] = f[14]
	s[15] = f[15]
}
