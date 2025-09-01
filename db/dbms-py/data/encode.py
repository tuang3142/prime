# TODO: slotted pages (sounds like segment tree)
# we have an tree of integer that store the poisiton of each item
# that way we can query quicker and faster (?) and not to load everything into memoery (right?)

# basically, turn everything in to binary, in a way that we can decode it
import csv
import os

# uint31 = unsigned int 32 - 2 ^ 32
schema = ('uint32', 'text', 'text')

path = 'db/dbms-py/db/'
def encode():
    with open(path + 'movies.csv', 'r') as f:
        # out is an binary encoded file?, convert from string (csv rows) to binary
        # out is written in binary
        with open(path + 'movies.dat', 'wb') as out:
            reader = csv.reader(f)
            print(next(reader)) # header
            for i, row in enumerate(reader):
                for typ, val in zip(schema, row):
                    if typ == 'uint32':
                        val = int(val)
                        # write number column (id) into bytes
                        out.write(int(val).to_bytes(4, 'little')) # 4 bytes, little indian (im not racist sorry)
                    elif typ == 'text':
                        val = val.encode('utf8')
                        out.write(len(val).to_bytes(1)) # write length of the string
                        out.write(val)
                    else:
                        raise ValueError('unknown type')

def decode():
    path = '' # hack - to run this program locally
    with open(path + 'movies.dat', 'rb') as f:
        size = os.fstat(f.fileno()).st_size
        while f.tell() < size:
            for typ in schema:
                row = []
                if typ == 'uint32':
                    b = f.read(4)
                    val = int.from_bytes(b, 'little')
                elif typ == 'text':
                    b = f.read(1)
                    l = int.from_bytes(b)
                    val = f.read(l)
                # print(val)
                row.append(val) 

