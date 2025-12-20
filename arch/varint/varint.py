def decode(arr: bytearray) -> int:
    out = 0
    for b in arr[::-1]:
        out <<= 7
        out += b & 0x7F
    return out


# varint is little endian
def encode(n: int) -> bytearray:
    out = bytearray()
    while True:
        b = n & 0x7F  # 7 bits all 1, pretty understandable
        n >>= 7
        if n:
            out.append(b | 0x80)  # 8 bits: 1 followed by 0
        else:
            out.append(b)
            break
    return out


# steps by steps, you can do it. break the pattern. stay focus.
# plan:
# integer -> 7 bits padding -> bytes -> little endian -> return


def main():
    assert encode(0) == b"\x00"
    assert encode(2**21 - 1) == b"\xff\xff\x7f"
    assert encode(2**21) == b"\x80\x80\x80\x01"
    assert encode(2**28) == b"\x80\x80\x80\x80\x01"
    print("encode: OK")

    assert 0 == decode(b"\x00")
    assert 1   == decode(b"\x01")
    assert 10  == decode(b"\x0a")
    assert 63  == decode(b"\x3f")
    assert 127 == decode(b"\x7f")
    assert 150 == decode(b"\x96\x01")

    assert 2**21 - 1 == decode(b"\xff\xff\x7f")
    assert 2**21 == decode(b"\x80\x80\x80\x01")
    assert 2**28 == decode(b"\x80\x80\x80\x80\x01")
    print("decode: OK")


main()
