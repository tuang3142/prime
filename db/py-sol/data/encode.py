import csv
import io
import os

PAGE_SIZE = 1024
COUNT_SIZE = 2
ENDIAN = 'big'

def encode_row(row: list[str]):
    schema = ('int', 'text', 'text')
    b = io.BytesIO()
    for typ, val in zip(schema, row):
        if typ == 'int':
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

def flush_to_file(out, buffer: bytearray, record_count: int) -> None:
    # write record count at start
    buffer[:COUNT_SIZE] = record_count.to_bytes(COUNT_SIZE, ENDIAN)
    out.write(buffer)
    
def encode_file(csv_path: str, out_path: str) -> None:
    with open(csv_path, 'r', newline='', encoding='UTF-8') as inp, \
         open(out_path, 'wb') as out:

        reader = csv.reader(inp)
        next(reader) # skip header

        record_end = PAGE_SIZE # start from the bottom up
        pos_start = COUNT_SIZE # skip some first (1 or 2) bytes cus they are saved for record count
        record_count = 0
        buffer = bytearray(PAGE_SIZE)

        for row in reader:
            record = encode_row(row)
            record_start = record_end - len(record)
            if record_start < pos_start + 2:
                flush_to_file(out, buffer, record_count)
                # reset page
                buffer = bytearray(PAGE_SIZE)
                record_end = PAGE_SIZE
                pos_start = COUNT_SIZE
                record_count = 0
                record_start = record_end - len(record)

            # even if this is the last row, we can still write it down below
            buffer[record_start:record_end] = record # write row byte data in the buffer
            buffer[pos_start:pos_start+2] = record_start.to_bytes(2, ENDIAN) # write starting position of row record

            record_count += 1
            pos_start += 2
            record_end = record_start

        # write remaning records (half-emptied buffer)
        if record_count > 0:
            flush_to_file(out, buffer, record_count)

encode_file('movies.csv', 'movies-paged.dat')

# crude testing with checking hexdump - C movies_slotted.dat, there should be innocene written
