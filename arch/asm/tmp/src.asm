section .text

global foo

foo:
	xor rax, rax

.loop:
	movzx rcx, byte [rdi]
	cmp   rcx, 0
	je    .loop_end
	inc   rdi
	sub   rcx, '0'
	add   rax, rcx
	jmp   .loop

.loop_end:
	; cmp rax, 10
	; sete al
	; movzx rax, al
	ret
