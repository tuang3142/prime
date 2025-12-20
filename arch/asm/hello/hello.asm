; ----------------------------------------------------------------------------------------
; Writes "Hello, World" to the console using only system calls. Runs on 64-bit Linux only.
; To assemble and run:
;
;     nasm -felf64 hello.asm && ld hello.o && ./a.out
; ----------------------------------------------------------------------------------------

          global    _start ; function definition, i think

          section   .text ; this is "code, text"
_start:   mov       rax, 1                  ; system call for write
          mov       rdi, 1                  ; file handle 1 is stdout
          mov       rsi, msg                ; address of string to output
          mov       rdx, 14                 ; number of bytes: the string hello and end of line
          syscall                           ; invoke operating system to do the write
          mov       rax, 60                 ; system call for exit
          xor       rdi, rdi                ; exit code 0
          syscall                           ; invoke operating system to exit

          section   .data ; this is the value of the variable;
msg:  db        "Hello, World!", 10      ; note the newline at the end

; system call = "built-in" function
; TODO: what is this system call? how can i build a mental model about this system call? where is the first turtle?
; section .text = code
; does c translate to asm?