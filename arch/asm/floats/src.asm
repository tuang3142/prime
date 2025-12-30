section .text

global mul
global cone_volume

mul:
	mulsd xmm0, xmm1
	ret

cone_volume:
	mulsd xmm0, xmm0
	mulsd xmm0, xmm1
	movsd xmm1, [rel pi64]; xmm1 = pi
	mulsd xmm0, xmm1
	movsd xmm1, [rel inv3]; xmm1 = 1/3
	mulsd xmm0, xmm1
	ret

section .rodata
pi64:    dq 0x400921FB54442D18    ; 3.141592653589793 (hex IEEE-754)
inv3:    dq 0x3FD5555555555555    ; 1/3 (double)
