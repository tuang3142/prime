import os

COUNT_SIZE = 2
PAGE_SIZE = 1024
ENDIAN = "big"

# todo: two bug: 
# first row 0, page end got weird
# check csv file vs dat size

schema = ['uint32', 'text', 'text']
with open("movies-paged.dat", "rb") as f:
    while True:
        buffer = f.read(PAGE_SIZE)
        if len(buffer) == 0:
            break
        assert len(buffer) == PAGE_SIZE
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
