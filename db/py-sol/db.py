import csv

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
    artifacts = (
        ('elderwd', 'Elder Wand', 99.0, False),
        ('philstn', "Philosopher's Stone", 95.0, False),
        ('horcrux1', "Tom Riddle's Diary", 90.0, True),
        ('horcrux2', "Marvolo Gaunt's Ring", 88.0, True),
        ('cloak', 'Invisibility Cloak', 70.0, False),
        ('sword', 'Sword of Gryffindor', 85.0, False),
        ('mapmara', "Marauder's Map", 40.0, False),
        ('timeurn', 'Time-Turner', 75.0, False),
        ('basilisk', 'Basilisk Fang', 60.0, False),
        ('lktpnsn', 'Locket of Slytherin', 87.0, True),
    )

    schema = (
        ('id', str),
        ('name', str),
        ('power', float),
        ('is_dark', bool),
    )

    assert tuple(run(Q(
        Projection(lambda x: (x[0],)),
        Selection(lambda x: x[3]),
        MemoryScan(artifacts),
    ))) == (
        ('horcrux1',),
        ('horcrux2',),
        ('lktpnsn',),
    )
    print('ok 1')

    assert tuple(run(Q(
        Projection(lambda x: (x[0], x[2])),
        Limit(3),
        Sort(lambda x: x[2], desc=True),
        MemoryScan(artifacts),
    ))) == (
        ('elderwd', 99.0),
        ('philstn', 95.0),
        ('horcrux1', 90.0),
    )
    print('ok 2')

    assert tuple(run(Q(
        Projection(lambda x: (x[1], x[2])),
        Limit(3),
        Sort(lambda x: x[2]),
        Selection(lambda x: (not x[3]) and (60.0 <= x[2] < 90.0)),
        MemoryScan(artifacts),
    ))) == (
        ('Basilisk Fang', 60.0),
        ('Invisibility Cloak', 70.0),
        ('Time-Turner', 75.0),
    )
    print('ok 3')

    assert tuple(run(Q(
        Projection(lambda x: (x[0],)),
        Selection(lambda x: ' of ' in x[1]),
        MemoryScan(artifacts),
    ))) == (
        ('sword',),
        ('lktpnsn',),
    )
    print('ok 4')

    assert tuple(run(Q(
        Projection(lambda x: (x[0],)),
        Limit(5),
        Sort(lambda x: x[1]),
        MemoryScan(artifacts),
    ))) == (
        ('basilisk',),
        ('elderwd',),
        ('cloak',),
        ('lktpnsn',),
        ('mapmara',),
    )
    print('ok 5')

    assert tuple(run(Q(
        Projection(lambda x: (x[0],)),
        Selection(lambda x: x[2] > 120.0),
        MemoryScan(artifacts),
    ))) == ()
    print('ok 6')

    assert tuple(run(Q(
        Projection(lambda x: (x[1],)),
        Selection(lambda x: 'wand' in x[1].lower()),
        MemoryScan(artifacts),
    ))) == (('Elder Wand',),)
    print('ok 7')


    with open('data/movies.csv', 'r', newline='', encoding='UTF-8') as inp:
        r = csv.reader(inp)
        next(inp) # skip header

        schema = (
            ('id', int),
            ('title', str),
            ('genres', str),
        )
        t = tuple(run(Q(
            Projection(lambda x: (x[1],)),
            Limit(2),
            FileScan(r),
        )))
        assert t == (
            ('Toy Story (1995)',),
            ('Jumanji (1995)',),
        ), f'got = {t}'
        print('ok 8')

    schema = (
        ('id', int),
        ('title', str),
        ('genres', str),
    )
    t = tuple(run(Q(
        Projection(lambda x: (x[1],)),
        Limit(3),
        Sort(lambda x: x[0], desc=True),
        FileScanPaged('data/movies-paged.dat'),
    )))
    assert t == (
        ('Innocence (2014)',),
        ('Rentun Ruusu (2001)',),
        ('The Pirates (2014)',)
    ), f'got = {t}'
    print('ok 9')
