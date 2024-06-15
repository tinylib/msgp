package msgp

import (
	"encoding/binary"
	"testing"
)

func BenchmarkIntegers(b *testing.B) {
	bytes64 := make([]byte, 9)
	bytes32 := make([]byte, 5)
	bytes16 := make([]byte, 3)

	b.Run("Int64", func(b *testing.B) {
		b.Run("Put", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				putMint64(bytes64, -1234567890123456789)
			}
		})
		b.Run("Get", func(b *testing.B) {
			putMint64(bytes64, -1234567890123456789)
			for i := 0; i < b.N; i++ {
				getMint64(bytes64)
			}
		})
	})
	b.Run("Int32", func(b *testing.B) {
		b.Run("Put", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				putMint32(bytes32, -123456789)
			}
		})
		b.Run("Get", func(b *testing.B) {
			putMint32(bytes32, -123456789)
			for i := 0; i < b.N; i++ {
				getMint32(bytes32)
			}
		})
	})
	b.Run("Int16", func(b *testing.B) {
		b.Run("Put", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				putMint16(bytes16, -12345)
			}
		})
		b.Run("Get", func(b *testing.B) {
			putMint16(bytes16, -12345)
			for i := 0; i < b.N; i++ {
				getMint16(bytes16)
			}
		})
	})

	b.Run("Uint64", func(b *testing.B) {
		b.Run("Put", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				putMuint64(bytes64, 1234567890123456789)
			}
		})
		b.Run("Get", func(b *testing.B) {
			putMuint64(bytes64, 1234567890123456789)
			for i := 0; i < b.N; i++ {
				getMuint64(bytes64)
			}
		})
	})
	b.Run("Uint32", func(b *testing.B) {
		b.Run("Put", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				putMuint32(bytes32, 123456789)
			}
		})
		b.Run("Get", func(b *testing.B) {
			putMuint32(bytes32, 123456789)
			for i := 0; i < b.N; i++ {
				getMuint32(bytes32)
			}
		})
	})
	b.Run("Uint16", func(b *testing.B) {
		b.Run("Put", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				putMuint16(bytes16, 12345)
			}
		})
		b.Run("Get", func(b *testing.B) {
			putMuint16(bytes16, 12345)
			for i := 0; i < b.N; i++ {
				getMuint16(bytes16)
			}
		})
	})
}

func BenchmarkIntegersUnix(b *testing.B) {
	bytes := make([]byte, 12)
	var sec int64 = 1609459200
	var nsec int32 = 123456789

	b.Run("Get", func(b *testing.B) {
		binary.BigEndian.PutUint64(bytes, uint64(sec))
		binary.BigEndian.PutUint32(bytes[8:], uint32(nsec))
		for i := 0; i < b.N; i++ {
			getUnix(bytes)
		}
	})

	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			putUnix(bytes, sec, nsec)
		}
	})
}

func BenchmarkIntegersPrefix(b *testing.B) {
	bytesU16 := make([]byte, 3)
	bytesU32 := make([]byte, 5)
	bytesU64 := make([]byte, 9)

	b.Run("u16", func(b *testing.B) {
		var pre byte = 0x01
		var sz uint16 = 12345
		for i := 0; i < b.N; i++ {
			prefixu16(bytesU16, pre, sz)
		}
	})
	b.Run("u32", func(b *testing.B) {
		var pre byte = 0x02
		var sz uint32 = 123456789
		for i := 0; i < b.N; i++ {
			prefixu32(bytesU32, pre, sz)
		}
	})
	b.Run("u64", func(b *testing.B) {
		var pre byte = 0x03
		var sz uint64 = 1234567890123456789
		for i := 0; i < b.N; i++ {
			prefixu64(bytesU64, pre, sz)
		}
	})
}
