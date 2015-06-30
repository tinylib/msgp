package msgp

/* ----------------------------------
	integer encoding utilities
	(inline-able)

	TODO(tinylib): there are faster,
	albeit non-portable solutions
	to the code below. implement
	byteswap?
   ---------------------------------- */

func putMint16(b []byte, i int16) {
	b[0] = mint16
	b[1] = byte(i >> 8)
	b[2] = byte(i)
}

func getMint16(b []byte) (i int16) {
	return (int16(b[1]) << 8) | int16(b[2])
}

func putMint8(b []byte, i int8) {
	b[0] = mint8
	b[1] = byte(i)
}

func getMint8(b []byte) (i int8) {
	return int8(b[1])
}

func putMuint32(b []byte, u uint32) {
	b[0] = muint32
	b[1] = byte(u >> 24)
	b[2] = byte(u >> 16)
	b[3] = byte(u >> 8)
	b[4] = byte(u)
}

func getMuint32(b []byte) uint32 {
	return (uint32(b[1]) << 24) | (uint32(b[2]) << 16) | (uint32(b[3]) << 8) | (uint32(b[4]))
}

func putMuint16(b []byte, u uint16) {
	b[0] = muint16
	b[1] = byte(u >> 8)
	b[2] = byte(u)
}

func getMuint16(b []byte) uint16 {
	return (uint16(b[1]) << 8) | uint16(b[2])
}

func putMuint8(b []byte, u uint8) {
	b[0] = muint8
	b[1] = byte(u)
}

func getMuint8(b []byte) uint8 {
	return uint8(b[1])
}

/* -----------------------------
		prefix utilities
   ----------------------------- */

// write prefix and uint8
func prefixu8(b []byte, sz uint8, pre byte) {
	b[0] = pre
	b[1] = byte(sz)
}

// write prefix and big-endian uint16
func prefixu16(b []byte, sz uint16, pre byte) {
	b[0] = pre
	b[1] = byte(sz >> 8)
	b[2] = byte(sz)
}
