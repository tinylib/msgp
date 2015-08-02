//+build !noasm

#include "textflag.h"

// getUnix(b []byte) (sec int64, nsec int32)
TEXT ·getUnix(SB),NOSPLIT,$0-36
	MOVQ   b+0(FP), AX
	MOVQ   0(AX), CX
	MOVL   8(AX), DX
	BSWAPQ CX 
	BSWAPL DX 
	MOVQ   CX, sec+24(FP)
	MOVL   DX, nsec+32(FP)
	RET

// putUnix(b []byte, sec int64, nsec int32)
TEXT ·putUnix(SB),NOSPLIT,$0-36
	MOVQ   b+0(FP), AX 
	MOVQ   sec+24(FP), CX 
	MOVL   nsec+32(FP), DX 
	BSWAPQ CX
	BSWAPL DX 
	MOVQ   CX, 0(AX)
	MOVL   DX, 8(AX)
	RET

// putMint64(b []byte, i int64)
TEXT ·putMint64(SB),NOSPLIT,$0-32
	MOVQ   b+0(FP), AX
	MOVQ   i+24(FP), CX 
	BSWAPQ CX 
	MOVB   $0xd3, 0(AX)
	MOVQ   CX, 1(AX)
	RET

// putMuint64(b []byte, u uint64)
TEXT ·putMuint64(SB),NOSPLIT,$0-32
	MOVQ   b+0(FP), AX 
	MOVQ   u+24(FP), CX 
	BSWAPQ CX 
	MOVB   $0xcf, 0(AX)
	MOVQ   CX, 1(AX)
	RET

// getMint64(b []byte) int64
TEXT ·getMint64(SB),NOSPLIT,$0-32
	MOVQ   b+0(FP), AX 
	MOVQ   1(AX), CX 
	BSWAPQ CX 
	MOVQ   CX, ret+24(FP)
	RET 

// getMuint64(b []byte) uint64
TEXT ·getMuint64(SB),NOSPLIT,$0-32
	JMP ·getMint64(SB)

TEXT ·putMint32(SB),NOSPLIT,$0-28
	MOVQ   b+0(FP), AX 
	MOVL   b+24(FP), CX 
	BSWAPL CX 
	MOVB   $0xd2, 0(AX)
	MOVL   CX, 1(AX)
	RET   

TEXT ·getMint32(SB),NOSPLIT,$0-28
	MOVQ   b+0(FP), AX 
	MOVL   1(AX), CX 
	BSWAPL CX
	MOVL   CX, ret+24(FP)
	RET 

TEXT ·prefixu64(SB),NOSPLIT,$0-33
	MOVQ   b+0(FP), AX 
	MOVQ   sz+24(FP), CX 
	MOVB   pre+32(FP), DX 
	BSWAPQ CX 
	MOVB   DX, 0(AX)
	MOVQ   CX, 1(AX)
	RET 

TEXT ·prefixu32(SB),NOSPLIT,$0-29
	MOVQ   b+0(FP), AX 
	MOVL   sz+24(FP), CX 
	MOVB   pre+28(FP), DX 
	BSWAPL CX
	MOVB   DX, 0(AX)
	MOVL   CX, 1(AX)
	RET 

TEXT ·put232(SB),NOSPLIT,$0-32
	MOVQ   p+0(FP), AX 
	MOVL   a+24(FP), CX 
	MOVL   b+28(FP), DX 
	BSWAPL CX 
	BSWAPL DX 
	MOVL   CX, 0(AX)
	MOVL   DX, 4(AX)
	RET

TEXT ·put264(SB),NOSPLIT,$0-40
	MOVQ   p+0(FP), AX
	MOVQ   a+24(FP), CX 
	MOVQ   b+32(FP), DX 
	BSWAPQ CX 
	BSWAPQ DX 
	MOVQ   CX, 0(AX)
	MOVQ   DX, 8(AX)
	RET

TEXT ·putMapHdr(SB),NOSPLIT,$0-24
	MOVQ p+0(FP), AX
	MOVQ sz+8(FP), DI
	MOVQ $0x80, BX
	MOVQ $0xde, SI
	JMP  put4bit<>(SB)

TEXT ·putArrayHdr(SB),NOSPLIT,$0-24
	MOVQ p+0(FP), AX
	MOVQ sz+8(FP), DI
	MOVQ $0x90, BX
	MOVQ $0xdc, SI
	JMP  put4bit<>(SB)

// putHdr for maps and arrays
// that do not have 8-bit prefixes
// ptr in AX, size in DI, desc in SI
// fixbits in BX
TEXT put4bit<>(SB),NOSPLIT,$0-24
	CMPQ DI, $16
	JAE  put16
	ORQ  BX, DI
	MOVB DI, (AX)
	MOVQ $1, ret+16(FP)
	RET
put16:
	CMPQ   DI, $0xffff
	JA     put32
	BSWAPL DI
	SHRQ   $8, DI
	ORQ    DI, SI
	MOVQ   $3, ret+16(FP)
	MOVQ   SI, (AX)
	RET
put32:
	BSWAPL DI
	INCQ   SI
	SHLQ   $8, DI
	ORQ    DI, SI
	MOVQ   $5, ret+16(FP)
	MOVQ   SI, (AX)
	RET

// putStrHdr(p *byte, sz uint32) int
TEXT ·putStrHdr(SB),NOSPLIT,$0-24
	MOVQ p+0(FP), AX
	MOVQ sz+8(FP), DI
	CMPQ DI, $32
	JAE  putstr
	ORQ  $0xa0, DI
	MOVB DI, (AX)
	MOVQ $1, ret+16(FP)
	RET
putstr:
	MOVQ $0xd9, SI
	JMP  putHdr<>(SB)

// putBinHdr(p *byte, sz int) int
TEXT ·putBinHdr(SB),NOSPLIT,$0-24
	MOVQ p+0(FP), AX
	MOVQ sz+8(FP), DI
	MOVL $0xc4, SI
	JMP  putHdr<>(SB)

// putHdr is the body for storing
// either a big-ending uint8, 
// uint16, or uint32.
//
// compare jumps are forward under
// the assumption that smaller
// headers are more likely than 
// larger.
TEXT putHdr<>(SB),NOSPLIT,$0-24
	CMPL   DI, $0xff
	JA     test16
	MOVL   $2, CX
	SHLL   $8, DI
	ORL    DI, SI
	MOVQ   CX, ret+16(FP)
	MOVL   SI, (AX)
	RET
test16:
	CMPL   DI, $0xffff
	JA     test32
	MOVL   $3, CX
	BSWAPL DI
	INCQ   SI
	SHRQ   $8, DI
	ORQ    DI, SI
	MOVQ   CX, ret+16(FP)
	MOVL   SI, (AX)
	RET
test32:
	MOVL   $5, CX
	BSWAPL DI
	ADDQ   $2, SI
	SHLQ   $8, DI
	ORQ    DI, SI
	MOVQ   CX, ret+16(FP)
	MOVQ   SI, (AX)
	RET
