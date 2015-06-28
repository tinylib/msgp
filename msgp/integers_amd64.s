
#include "textflag.h"

TEXT ·getUnix(SB),NOSPLIT,$0-36
	MOVQ   b+0(FP), AX
	MOVQ   0(AX), CX
	MOVL   8(AX), DX
	BSWAPQ CX 
	BSWAPL DX 
	MOVQ   CX, sec+24(FP)
	MOVL   DX, nsec+32(FP)
	RET


TEXT ·putUnix(SB),NOSPLIT,$0-36
	MOVQ   b+0(FP), AX 
	MOVQ   sec+24(FP), CX 
	MOVL   nsec+32(FP), DX 
	BSWAPQ CX
	BSWAPL DX 
	MOVQ   CX, 0(AX)
	MOVL   DX, 8(AX)
	RET
