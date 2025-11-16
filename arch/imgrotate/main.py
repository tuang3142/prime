import sys

def run(in_path: str, out_path: str = "out.bmp") -> None:
    data = bytearray(open(in_path, "rb").read())

    pixels_offset = int.from_bytes(data[10:14], "little", signed=False)
    w = int.from_bytes(data[18:22], "little", signed=True)
    h = int.from_bytes(data[22:26], "little", signed=True)

    # per-row padding (24-bit BMP): rows are width*3 bytes, padded to a multiple of 4
    sp = 4 - ((w * 3) % 4)
    tp = 4 - ((h * 3) % 4)

    pixels = bytearray()
    for ty in range(w):
        for tx in range(h):
            sy = tx
            sx = w - ty - 1
            n = pixels_offset + 3 * (sy * w + sx) + sy * sp
            pixels.extend(data[n:n + 3])
        pixels.extend(b"\x00" * tp)

    # swap width/height in header copy
    header = bytearray(data[:pixels_offset])
    header[18:22], header[22:26] = header[22:26], header[18:22]

    open(out_path, "wb").write(header + pixels)

if __name__ == "__main__":
    run(in_path="testdata/2_in.bmp")
