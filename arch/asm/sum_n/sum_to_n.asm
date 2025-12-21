section .text

global sum_to_n

; n is in rdi
; result is in rax

; 0(1)
sum_to_n:
	mov  rax, rdi
	inc  rax
	imul rax, rdi
	shr  rax, 1

; O(N)
; .loop:
; 	cmp rbx, rdi
; 	jg  .done
; 	add rax, rbx
; 	add rbx, 1
; 	jmp .loop

.done:
	ret
