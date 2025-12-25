section .text
global  binconv

binconv:
	xor rax, rax; accum = 0

.loop:
	movzx rcx, byte [rdi]; (c = *s)
	cmp   rcx, 0; ((c = *s) == NULL)
	je    .end
	add   rdi, 1; s++
	sub   rcx, '0'; subtract '0' from *s
	shl   rax, 1; accum <<= 1
	add   rax, rcx; accum += c
	jmp   .loop

.end:
	ret

	; Also since you have a C version, you may wish to run it through godbolt to see how close it is to your assembly :)
	; todo: sanitize input, throw error and such

