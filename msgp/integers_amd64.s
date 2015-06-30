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

#define STRDESC 0xd9

// putStrHdr(p *byte, sz uint32) int
TEXT ·putStrHdr(SB),NOSPLIT,$0-24
	MOVQ p+0(FP), AX
	MOVQ sz+8(FP), DI
	CMPL DI, $32
	JB   fixstr
	MOVQ $STRDESC, SI
	JMP  putHdr<>(SB)
fixstr:
	ORL $0xa0, DI
	MOVB DI, (AX)
	MOVQ $1, ret+16(FP)
	RET

#define BINDESC 0xc4

// putBinHdr(p *byte, sz int) int
TEXT ·putBinHdr(SB),NOSPLIT,$0-24
	MOVQ p+0(FP), AX
	MOVQ sz+8(FP), DI
	MOVQ $BINDESC, SI
	JMP  putHdr<>(SB)

// ptr in AX, size in DI, desc in SI, ret in BX
TEXT putHdr<>(SB),NOSPLIT,$0-24
	BSRQ DI, CX
	JZ   zero
	ANDQ $0x18, CX  // CX is now 0, 8, 16, 24
	MOVQ $hdrtab<>(SB), DX
	ADDQ CX, DX
	JMP  (DX)      // jumps into hdr_{8/16/32}<>(SB)
zero:
	MOVQ $2, BX
	MOVQ SI, (AX)
	MOVQ BX, ret+16(FP)
	RET            

TEXT hdr_8<>(SB),NOSPLIT,$0-24
	MOVQ $2, BX
	SHLQ $8, DI
	ORQ  DI, SI
	MOVQ BX, ret+16(FP)
	MOVQ SI, (AX)
	RET

TEXT hdr_16<>(SB),NOSPLIT,$0-24
	MOVQ   $3, BX
	BSWAPL DI
	INCQ   SI
	SHRQ   $8, DI
	ORQ    DI, SI
	MOVQ   BX, ret+16(FP)
	MOVQ   SI, (AX)
	RET

TEXT hdr_32<>(SB),NOSPLIT,$0-24
	MOVQ   $5, BX
	BSWAPL DI
	SHLQ   $8, DI
	ADDQ   $2, SI
	ORQ    DI, SI
	MOVQ   BX, ret+16(FP)
	MOVQ   SI, (AX)
	RET

// jump table for header
DATA hdrtab<>+0(SB)/8, $hdr_8<>(SB)
DATA hdrtab<>+8(SB)/8, $hdr_16<>(SB)
DATA hdrtab<>+16(SB)/8, $hdr_32<>(SB)
DATA hdrtab<>+24(SB)/8, $hdr_32<>(SB)
GLOBL hdrtab<>(SB),RODATA,$36
