; add.asm
; assemble: nasm -f elf64 add.asm
; link:     ld -o add add.o
; run:      ./add ; echo $?

global _start

section .text

_start:
    mov rax, 5      ; first number
    mov rbx, 7      ; second number
    add rax, rbx    ; rax = 12

    ; exit(status = rax)
    mov rdi, rax    ; exit code
    mov rax, 60     ; sys_exit
    syscall