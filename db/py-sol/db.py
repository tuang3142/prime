import csv
import os
import io

COUNT_SIZE = 2
PAGE_SIZE = 1024
ENDIAN = "big"

# todo: what i did here is pretty correct. need to debug. don't stupidly ask ai for help. this is the gym, this is practice, don't make it easy
# TODO finish this exer
# write on empty test file with custom schema to test it
# todo: write this in go, with just two node: insert and read (filescanpaged)
class Insert(object):
    def __init__(self, ip, schema):
        self.ip = ip
        self.schema = schema

    def encode_row(self, row: list[str]):
        b = io.BytesIO()
        for typ, val in zip(self.schema, row):
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

    def next(self):
        print('next is called')
        if self.child is None: return

        row = self.child.next()
        if row is None: return

        print('row inside insert: ', row)
        self.ip.seek(-PAGE_SIZE, os.SEEK_END)

        page = bytearray(self.ip.read(PAGE_SIZE))
        rec_count = int.from_bytes(page[:COUNT_SIZE])
        pos = rec_count * 2 # last position of the pos array
        rec_end = int.from_bytes(page[pos:pos+2])

        b = self.encode_row(row)
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

        if not new_page: self.ip.seek(-PAGE_SIZE, os.SEEK_END)
        self.ip.write(page)

        return row


class FileScanPaged(object):
    PAGE_SIZE = 1024 # page byte size
    REC_COUNT_SIZE = 2 # 2 bytes for record count, saved at the beginning of the page
    SLOT_SIZE = 2 # 2 bytes per record index in a page

    def __init__(self, file:str):
        self.file = open(file, "rb")
        self.rows = []
    
    def close(self):
        if self.file:
            self.file.close()

    def next(self):
        if len(self.rows) == 0:
            self.get_next_rows()
            if len(self.rows) == 0:
                self.close()
                return

        row = self.rows[0]
        self.rows = self.rows[1:]
        return row
    
    def get_next_rows(self):
        buffer = self.file.read(self.PAGE_SIZE)
        if len(buffer) == 0:
            return []
        
        rec_count = int.from_bytes(buffer[:self.SLOT_SIZE])
        p = self.REC_COUNT_SIZE
        rec_end = self.PAGE_SIZE
        for _ in range(rec_count):
            rec_start = int.from_bytes(buffer[p : p + self.SLOT_SIZE])
            rec = buffer[rec_start:rec_end]
            row = []
            for typ in ['uint32', 'text', 'text']:
                if typ == 'uint32':
                    val = int.from_bytes(rec[:4]) # 4 bytes = 32 bit -> unsigned 32 bit int
                    row.append(val)
                    rec = rec[4:]
                elif typ == 'text':
                    text_len = int.from_bytes(rec[:1])
                    text_start, text_end = 1, 1 + text_len
                    row.append(rec[text_start:text_end].decode('utf-8'))
                    rec = rec[text_end:]
                else:
                    raise ValueError(f'Unknown type {typ}')
            p += 2
            rec_end = rec_start
            self.rows.append(row)
        return []


class FileScan(object):
    def __init__(self, reader):
        self.reader = reader

    def next(self):
        try:
            row = next(self.reader)
            return row
        except StopIteration:
            pass
        


class MemoryScan(object):
    """
    Yield all records from the given "table" in memory.

    This is really just for testing... in the future our scan nodes
    will read from disk.
    """
    def __init__(self, table):
        self.table = table
        self.idx = 0

    def next(self):
        if self.idx >= len(self.table):
            return

        x = self.table[self.idx]
        self.idx += 1
        return x


class Projection(object):
    """
    Map the child records using the given map function, e.g. to return a subset
    of the fields.
    """
    def __init__(self, proj):
        self.proj = proj

    def next(self):
        row = self.child.next()
        if row is not None:
            return self.proj(row)


class Selection(object):
    """
    Filter the child records using the given predicate function.

    Yes it's confusing to call this "selection" as it's unrelated to SELECT in
    SQL, and is more like the WHERE clause. We keep the naming to be consistent
    with the literature.
    """
    def __init__(self, select):
        self.select = select

    def next(self):
        while True:
            row = self.child.next() # keep on looping, don't return the false result
            if row == None:
                return
            if self.select(row):
                return row


class Limit(object):
    """
    Return only as many as the limit, then stop
    """
    def __init__(self, limit):
        self.limit = limit

    def next(self):
        if self.limit == 0:
            return
        self.limit -= 1
        return self.child.next()


class Sort(object):
    """
    Sort based on the given key function
    """
    def __init__(self, key, desc=False):
        self.key = key
        self.desc = desc
        self.table = []
        self.idx = 0
        self.sorted = False

    def next(self):
        if not self.sorted:
            r = self.child.next()
            while r is not None:
                self.table.append(r)
                r = self.child.next()
            if self.desc:
                self.table.sort(key=lambda x: -self.key(x))
            else:
                self.table.sort(key=self.key)
            self.sorted = True

        if self.idx >= len(self.table):
            return
        r = self.table[self.idx]
        self.idx+=1
        return r

        

def Q(*nodes):
    """
    Construct a linked list of executor nodes from the given arguments,
    starting with a root node, and adding references to each child
    """
    ns = iter(nodes)
    parent = root = next(ns)
    for n in ns:
        parent.child = n
        parent = n
    return root


def run(q):
    """
    Run the given query to completion by calling `next` on the (presumed) root
    """
    while True:
        x = q.next()
        if x is None:
            break
        yield x


if __name__ == '__main__':
    # artifacts = (
    #     ('elderwd', 'Elder Wand', 99.0, False),
    #     ('philstn', "Philosopher's Stone", 95.0, False),
    #     ('horcrux1', "Tom Riddle's Diary", 90.0, True),
    #     ('horcrux2', "Marvolo Gaunt's Ring", 88.0, True),
    #     ('cloak', 'Invisibility Cloak', 70.0, False),
    #     ('sword', 'Sword of Gryffindor', 85.0, False),
    #     ('mapmara', "Marauder's Map", 40.0, False),
    #     ('timeurn', 'Time-Turner', 75.0, False),
    #     ('basilisk', 'Basilisk Fang', 60.0, False),
    #     ('lktpnsn', 'Locket of Slytherin', 87.0, True),
    # )
    # t = tuple(run(Q(
    #     # Projection(lambda x: (x[0],)),
    #     # Selection(lambda x: x[3]),
    #     MemoryScan(artifacts),
    # )))
    # assert t == (
    #     ('horcrux1',),
    #     ('horcrux2',),
    #     ('lktpnsn',),
    # )
    # print('ok 1')

    # assert tuple(run(Q(
    #     Projection(lambda x: (x[0], x[2])),
    #     Limit(3),
    #     Sort(lambda x: x[2], desc=True),
    #     MemoryScan(artifacts),
    # ))) == (
    #     ('elderwd', 99.0),
    #     ('philstn', 95.0),
    #     ('horcrux1', 90.0),
    # )
    # print('ok 2')

    # assert tuple(run(Q(
    #     Projection(lambda x: (x[1], x[2])),
    #     Limit(3),
    #     Sort(lambda x: x[2]),
    #     Selection(lambda x: (not x[3]) and (60.0 <= x[2] < 90.0)),
    #     MemoryScan(artifacts),
    # ))) == (
    #     ('Basilisk Fang', 60.0),
    #     ('Invisibility Cloak', 70.0),
    #     ('Time-Turner', 75.0),
    # )
    # print('ok 3')

    # assert tuple(run(Q(
    #     Projection(lambda x: (x[0],)),
    #     Selection(lambda x: ' of ' in x[1]),
    #     MemoryScan(artifacts),
    # ))) == (
    #     ('sword',),
    #     ('lktpnsn',),
    # )
    # print('ok 4')

    # assert tuple(run(Q(
    #     Projection(lambda x: (x[0],)),
    #     Limit(5),
    #     Sort(lambda x: x[1]),
    #     MemoryScan(artifacts),
    # ))) == (
    #     ('basilisk',),
    #     ('elderwd',),
    #     ('cloak',),
    #     ('lktpnsn',),
    #     ('mapmara',),
    # )
    # print('ok 5')

    # assert tuple(run(Q(
    #     Projection(lambda x: (x[0],)),
    #     Selection(lambda x: x[2] > 120.0),
    #     MemoryScan(artifacts),
    # ))) == ()
    # print('ok 6')

    # assert tuple(run(Q(
    #     Projection(lambda x: (x[1],)),
    #     Selection(lambda x: 'wand' in x[1].lower()),
    #     MemoryScan(artifacts),
    # ))) == (('Elder Wand',),)
    # print('ok 7')


    # with open('data/movies.csv', 'r', newline='', encoding='UTF-8') as inp:
    #     r = csv.reader(inp)
    #     next(inp) # skip header

    #     t = tuple(run(Q(
    #         Projection(lambda x: (x[1],)),
    #         Limit(2),
    #         FileScan(r),
    #     )))
    #     assert t == (
    #         ('Toy Story (1995)',),
    #         ('Jumanji (1995)',),
    #     ), f'got = {t}'
    #     print('ok 8')

    # # t = tuple(run(Q(
    # #     Projection(lambda x: (x[1],)),
    # #     Limit(3),
    # #     Sort(lambda x: x[0], desc=True),
    # #     FileScanPaged('data/movies-paged.dat'),
    # # )))
    # # assert t == (
    # #     ('Innocence (2014)',),
    # #     ('Rentun Ruusu (2001)',),
    # #     ('The Pirates (2014)',)
    # # ), f'got = {t}'
    # # print('ok 9')

    schema = ['uint32', 'text', 'text']
    movies = [
        [131267, 'The Roses (2025)', 'Commedy|Dark|Fantasy|Sci-Fi'],
        [131267, 'The Roses (2025)', 'Commedy|Dark|Fantasy|Sci-Fi'],
        [131267, 'The Roses (2025)', 'Commedy|Dark|Fantasy|Sci-Fi'],
        [131267, 'The Flower (2025)', 'Commedy|Dark|Fantasy|Sci-Fi'],
    ]
    file = "data/movies-paged-copy.dat"
    with open(file, 'r+b') as ip:
        run(Q(
            Insert(ip, schema),
            # Selection(lambda x: 'Flower' in x[1]),
            MemoryScan(movies)))
    last_row = []
    with open(file, 'r+b') as f:
        f.seek(-PAGE_SIZE, os.SEEK_END)
        buffer = f.read(PAGE_SIZE)
        rec_count = int.from_bytes(buffer[:2])
        print('rec_count ', rec_count)
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
                elif typ == 'text':
                    l = int.from_bytes(rec[:1])
                    row.append(rec[1:1+l].decode('utf-8'))
                    rec = rec[1+l:]
                else:
                    raise ValueError('unknown type: ', typ)
            last_row = row
            print(row)
            p += 2
            rec_end = rec_start
    assert last_row == [131267, 'The Roses (2025)', 'Commedy|Dark|Fantasy|Sci-Fi'], f'{last_row}'
    print('ok - test_insert')
