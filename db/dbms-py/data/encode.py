# TODO: slotted pages (sounds like segment tree)

"""
Slotted pages: idea
- keep track of each record in the dat file: 
- keep an array (index), that store the begging of each record (optionally the legnth of it, but we can calculate)
- i drawed a picture, but it is like: pos = []int{p1, p2, ... pn} where pi is the starting position of record i, in the file ("byte position")
"""

import csv
import io

schema = ('uint32', 'text', 'text')

PAGE_SIZE = 1024

path = 'db/dbms-py/db/'
def encode():
    with open(path + 'movies.csv', 'r') as f:
        # out is an binary encoded file?, convert from string (csv rows) to binary
        # out is written in binary
        with open(path + 'movies_slotted.dat', 'wb') as out:
            buffer = bytearray(PAGE_SIZE) # temporary data write in memoory, periodically flushed to out (movies_slotted.dat)
            records_in_buffer = 0
            records_start = PAGE_SIZE # start from the bottom up, but why?
            reader = csv.reader(f)
            print(next(reader)) # skip header
            for i, row in enumerate(reader):
                if i == 2: break # why
                record = io.BytesIO()
                for typ, val in zip(schema, row):
                    if typ == 'uint32':
                        val = int(val)
                        record.write(int(val).to_bytes(4, 'little')) # write id, little endian 4 bytes
                    elif typ == 'text':
                        bs = val.encode('utf8')
                        record.write(len(bs).to_bytes(1)) # write length of the string, 256-sized string (2^8)
                        record.write(val)
                    else:
                        raise ValueError('unknown type')
                record.seek(0)
                rec_bytes = record.read() # records in bytes, why need this?

                # TODO write num recored in the first 2 bytes (records_in_buffer) - ok this is one of the resaon we write from bottom up - leave space for meta data like this
                       # each pagese has its own count of number of record!!!!!
                # TODO write starting position of each record, but hwhere, mb also at the beginning of the page
                # stop before fillingin the middle 
                # ok i see you: top: page size (number of record), next: array indidate position of each readcord in the page
                # meet in the middle. nice man. we see this every wereh
                buffer[records_start-len(rec_bytes):records_start] = rec_bytes # write byte data in the buffer
                records_start -= len(rec_bytes)

encode()
