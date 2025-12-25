section .text
global  pangram

pangram:
	xor rax, rax
	mov rbx, 0x03ffffff; bit map for 26 chars

.loop:
	movzx rcx, byte [rdi]; loop through *s
	cmp   rcx, 0; break if '\0'
	je    .loop_end
	inc   rdi; move the pointer

	or    rcx, 32; turn on the 5th bit, force lower case
	sub   rcx, 'a' ; TODO: understand oz soltion?
	cmp   rcx, 26; if > 26, continue
	jae    .loop
	mov   rdx, 1
	shl   rdx, cl
	or    rax, rdx
	jmp   .loop

.loop_end:
	cmp   rax, rbx
	sete  al
	movzx rax, al
	ret
