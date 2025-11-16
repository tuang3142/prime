# code is easy but what we need to understand 
# how utf-8 work. i jsut need to watch the explaination video and take notes.
# multi byte encoding: things start with 0b10 (0x10)
# python is so shit. code this in go pls
# what i dont understand is how python read bytes. it is so confusion. each of these should be byte
# but they read to int? wtf? that is not how i wanted it to behave

def truncate(s, n):
    if n >= len(s): return s
    while n > 0 and (s[n] & 0xc0) == 0x80:
        n -= 1
    return s[:n]

with open("out", "wb") as out:
    for line in open("in", "rb"):
        line = truncate(line[1:-1], line[0]) + b'\n'
        out.write(line)