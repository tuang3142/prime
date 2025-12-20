# asm guide

original doc: [avenger, assemble](https://github.com/hackclub/some-assembly-required/blob/main/guide/introduction.md)

## why do this?

- better understanding of how computer work at the lowest level -> better programmer
- fun brain teaser; but also kinda like magic
- challenging, like a hard workout -> good for the brain

## what is a cpu

- CPU: brain; central processing unit
- CPU is fast but stupid. It can only:
  - read values  (number? in a register perhaps?)
  - write values
  - simple math: add, subtract, and, or, xor, not

### talk to a cpu

- computer: me like number (in fact, me can only understand number).
- hooman: me like text.
- assembly (hooman "readable" machine instruction) -> assembler -> numbero

- side notes: using annalogy is a good way to build a mental model to understand something
  - this guy used a lot: car, food, etc.

## how CPU work?

### ram vs register

- register: smaller, small efficient memory (smallter but better than ram)
- register vs RAM: RAM → register → compute → register → RAM.
  - imagine in the warehouse, the CPU has a long arm to grab the box from RAM to put it in to register. register is much closer to the CPU than in the RAM.

### decoder

- `add` r12 4: add 4 to the value at register 12 (lets say its 4)
  - in base 10: `1 12 4`.
  - CPU knows how to map `1` to `add`.
  - CPU -> send (1, 4, 4) to ALU -> ALU return 8 -> CPU store 8 to r12.

### physic (TODO - I skipped this part)

- 1 = on (high voltage, low resistant), 0 = off (high resistant, low voltage)
- bus???

## writing code

### data

- 1 & 0: integer, character (actually, float can be converted to bit), memory address,
- register: variable. 16 of them, "typically":
  - data: rax: accomulator - rbx: base - rcx: counter - rdx: data
  - index: rsi: source index - rdi: destinatioin inde x
  - r8-r15: addtional register
  - rsp, rbp: stack pointer, base pointer
  - rip: instruction pointer
- 32bit vs 64bit: refers to the size of the register

### instruction

- `mov`: move; `mov rax 5`: move 5 to the rax register.
  - more: `add`, `sub`, `mul`
- instruction pointer: keep track of the current line of code which the computer is executing
  - simply, instructions are stored in memory; the instruction pointer keep jumpping around between memoery addresss. in the example bellow, lets assume 1-5 is the address in RAM.

    ```assembly
      1. mov rcx 10
      2. jmp .newLine  ; jump to add 5, since jmp redirects RIP to address after .newLine
      3. sub rcx -10 
      4. .newLine
      5.   mul rcx 4
    ```

  - looks like label is a bookmark but it will not get executed

### flags

- zero flag `ZF`: add, sub, etc. writes to this flag; jump instruction reads this flag