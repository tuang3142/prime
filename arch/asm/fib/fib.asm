section .text
global  fib
global  sum_n
global slen

slen:
	movzx rcx, byte [rdi]
	test rcx, rcx
	jz .base
	inc rdi
	call slen
	inc rax
	ret
	

.base:
	xor rax, rax
	ret

fib:
	;   trick: push and pop rdi value to the stack
	;   return value is stored in rax
	;   TODO: p2 - how do asm hanle very big number? like over 64 bit? turn them to string?
	;   TODO: p2 - watch stack alignment
	cmp rdi, 1; base: n <= 1
	jbe .base

	push rdi; move it to stack; ok ok i get it: if you keep pushing in the stack without popping, you'll get stack overflow
	sub  rdi, 1
	call fib; fib(n - 1)
	pop  rdi
	push rax

	sub  rdi, 2
	call fib; fib(n - 2)
	pop  rbx
	add  rax, rbx; fib(n-1) + fib(n-2)
	ret

.base:
	;   return n if n <= 1
	mov rax, rdi
	ret

	; bonus: sum from 1 to n using recussion.
	; This is a more gentle introduction to recurrsion in asm.
	; I also want to check if I set and reset rax correctly.

sum_n:
	test rdi, rdi
	jz   .base
	push rdi
	dec  rdi
	call sum_n
	pop  rdi
	add  rax, rdi
	ret

.base:
	;   this how to reset rax: do it at base
	xor rax, rax
	ret
