# basically, turn everything in to binary, in a way that we can decode it
import csv
import io

# uint31 = unsigned int 32 - 2 ^ 32
schema = ('uint32', 'text', 'text')

path = 'db/dbms-py/db/'
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
                    out.write(int(val).to_bytes(4, 'little')) # 4 bytes, little indian (im not racist sorry)
                elif typ == 'text':
                    val = val.encode('utf8')
                    out.write(len(val).to_bytes(1)) # 1 bytes, string, this is kinda like padding for the string
                    out.write(val)
                else:
                    raise ValueError('unknown type')

with open(path + 'movies.dat', 'rb') as f:
    for  _ in range(1000):
        for typ in schema:
            row = []
            if typ == 'uint32':
                val = int.from_bytes(f.read(4), 'little')
            elif typ == 'text':
                l = int.from_bytes(f.read(1))
                val = f.read(l)
            print(val, end=' ')
            print('\n')
            row.append(val) 

# todo: decode this things. i should be able to do this well
# todo: interface, filescan, jump to a location # wow wow wow - this is how psql works?

