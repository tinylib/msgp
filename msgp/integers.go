package msgp

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
