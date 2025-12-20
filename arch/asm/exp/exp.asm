section .data
    ; Buffer to store the ASCII string representation of the number
    NEWLINE_CHAR equ 0Ah ; ASCII value for newline
    BUFFER_SIZE equ 10 ; Max digits for a 64-bit number (plus newline)
    output_buffer: db "          ", NEWLINE_CHAR

section .text
global _start

_start:
    ; --- 1. Power Calculation (Your Original Code) ---
    mov rax, 8          ; our exponent (8)
    mov rbx, 2          ; our base (2)
    mov rcx, 1          ; our result (starts at 1 = 2^0)
    mov rdx, 0          ; our counter

.calculatePower:
    ; The 'mul' instruction on x86-64 only takes one operand.
    ; It multiplies RAX by the operand and stores the 128-bit result in RDX:RAX.
    ; Since we are multiplying RCX by RBX and storing the 64-bit result in RCX,
    ; we need to use the `imul` (integer multiply) instruction with three operands
    ; or explicitly move the result back.
    ; For simplicity and correctness with the existing loop structure, we'll use:

    mov rax, rcx          ; Move the current result (rcx) into rax for multiplication
    mul rbx               ; rax = rax * rbx. The 64-bit result is in rax.
    mov rcx, rax          ; Move the new result back to rcx

    inc rdx               ; increment counter
    cmp rdx, 8            ; Compare counter (rdx) with exponent (8)
    jl .calculatePower    ; jump to the beginning of the loop if rdx < 8

    ; At this point, RCX contains the result (256)

    ; --- 2. Convert Number (RCX) to ASCII String ---
    ; This process is analogous to repeatedly dividing the number by 10
    ; and taking the remainder as the next digit (in reverse order).

    mov r8, output_buffer + BUFFER_SIZE - 2 ; Pointer to the last writable digit position
    mov r9, 10             ; Divisor (10)

.integer_to_string_loop:
    mov rax, rcx           ; Value to divide is in RCX, but division needs it in RAX
    xor rdx, rdx           ; Clear RDX (high part of the dividend for 64-bit division)
    div r9                 ; Divide RDX:RAX by 10.
                           ; Quotient (New RCX) -> RAX
                           ; Remainder (Digit) -> RDX

    mov rcx, rax           ; Save the quotient (new number) back to RCX
    add rdx, '0'           ; Convert remainder (0-9) to ASCII digit ('0'-'9')
    mov [r8], dl           ; Store the ASCII digit at the current pointer position

    dec r8                 ; Move the pointer to the next position (leftward)
    cmp rax, 0             ; Check if the quotient is 0
    jnz .integer_to_string_loop ; If not 0, continue dividing

    ; --- 3. Write String to Terminal (Linux Syscall) ---

    mov rax, 1             ; Syscall number 1: 'write'
    mov rdi, 1             ; File descriptor 1: STDOUT (the terminal)
    mov rsi, r8            ; Pointer to the start of the string (where the first digit was written)
    ; Calculate the length of the string: (output_buffer + BUFFER_SIZE) - RSI
    mov rdx, output_buffer + BUFFER_SIZE
    sub rdx, rsi           ; Length of the string (including the newline character)
    syscall                ; Execute the write syscall

    ; --- 4. Exit the Program (Linux Syscall) ---
    mov rax, 60            ; Syscall number 60: 'exit'
    xor rdi, rdi           ; Exit code 0 (success)
    syscall                ; Execute the exit syscall