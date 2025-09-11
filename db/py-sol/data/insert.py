import io
import os

COUNT_SIZE = 2
PAGE_SIZE = 1024
ENDIAN = "big"

schema = ['uint32', 'text', 'text']

def encode_row(row: list[str]):
    b = io.BytesIO()
    for typ, val in zip(schema, row):
        if typ == 'uint32':
            b.write(int(val).to_bytes(4, ENDIAN, signed=False)) # write id, big endian 4 bytes
        elif typ == 'text':
            data = val.encode('utf8')
            # 1 byte length = 255 maximum; cap or raise
            if len(data) > 255: # 1 byte: 0->255
                data = data[:255]
            b.write(len(data).to_bytes(1)) 
            b.write(data)
        else:
            raise ValueError(f'unknown type: {typ}')
    return b.getvalue()

def insert(inp, recs):
    with open(inp, 'r+b') as f:
        for rec in recs:
            insert_single(f, rec)

def insert_single(ip, rec):
    ip.seek(-PAGE_SIZE, os.SEEK_END)

    page = bytearray(ip.read(PAGE_SIZE))
    rec_count = int.from_bytes(page[:COUNT_SIZE])
    pos = rec_count * 2 # last position of the pos array
    rec_end = int.from_bytes(page[pos:pos+2])

    b = encode_row(rec)
    rec_start = rec_end - len(b)
    new_page = False
    if rec_start < pos + 2: # need a new page
        page = bytearray(PAGE_SIZE)
        rec_count = 0
        pos = 2
        rec_end = PAGE_SIZE
        rec_start = rec_end - len(b)
        new_page = True

    page[rec_start:rec_end] = b
    rec_count += 1
    page[:2] = rec_count.to_bytes(2) # big, signed = false
    page[rec_count * 2: rec_count * 2 + 2] = rec_start.to_bytes(2) # big, signed = false

    if not new_page: ip.seek(-PAGE_SIZE, os.SEEK_END)
    ip.write(page)


file = "movies-paged-copy.dat"
# insert(file, [
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
#     [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"],
# ])
# with open(file, 'r+b') as f:
#     insert_single(f, [131265, "Superman (2025)", "Commedy|Drama|Fantasy|Sci-Fi"])

# test
with open(file, "rb") as f:
    f.seek(-PAGE_SIZE, os.SEEK_END)
    buffer = f.read(PAGE_SIZE)
    rec_count = int.from_bytes(buffer[:2])
    print("record count: ", rec_count)
    p = 2
    rec_end = PAGE_SIZE
    for i in range(rec_count):
        rec_start = int.from_bytes(buffer[p : p + 2])
        rec = buffer[rec_start:rec_end]
        row = []
        for typ in schema:
            if typ == 'uint32':
                val = int.from_bytes(rec[:4]) # 4 bytes = 32 bit -> unsigned 32 bit int
                row.append(val)
                rec = rec[4:]
            if typ == 'text':
                l = int.from_bytes(rec[:1])
                row.append(rec[1:1+l].decode('utf-8'))
                rec = rec[1+l:]
        print(row)
        p += 2
        rec_end = rec_start

