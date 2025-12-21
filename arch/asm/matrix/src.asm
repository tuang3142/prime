section .text

global matrix_get

matrix_get:
	imul rcx, rdx             ; i * n
	add  rcx, r8              ; + j
	mov  rax, [rdi + 4 * rcx] ; m + i * n + j, but TODO: how does assembly know the number 4? (p2)
    // TODO: transpose	 (p1)

	ret
