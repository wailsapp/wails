//go:build windows

#include "textflag.h"
#include "go_asm.h"

// callDouble runs on the g0 stack with the C calling convention:
// runtime.cgocall hands it *doubleCall in R0. It loads a1-a4 into R0-R3 and
// v into D0 — the Windows ARM64 ABI numbers integer and floating-point
// argument registers independently, so the double is D0 whatever precedes
// it — calls fn and stores the HRESULT into hr. R19 keeps the block pointer
// across the call (callee-saved); the $16 frame saves LR and gives R19 its
// slot, as the runtime's asmstdcall does.
GLOBL ·callDoubleABI0(SB), NOPTR|RODATA, $8
DATA ·callDoubleABI0(SB)/8, $callDouble(SB)
TEXT callDouble(SB), NOSPLIT, $16
	MOVD	R19, 16(RSP)
	MOVD	R0, R19

	MOVD	doubleCall_a1(R19), R0
	MOVD	doubleCall_a2(R19), R1
	MOVD	doubleCall_a3(R19), R2
	MOVD	doubleCall_a4(R19), R3
	FMOVD	doubleCall_v(R19), F0
	MOVD	doubleCall_fn(R19), R12
	BL	(R12)

	MOVD	R0, doubleCall_hr(R19)
	MOVD	16(RSP), R19
	RET
