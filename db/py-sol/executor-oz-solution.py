import io
import os


class QueryNode(object):
    def close(self):
        if self.child:
            self.child.close()


class Root(QueryNode):
    def next(self):
        x = self.child.next()
        if x is None:
            self.child.close()
        return x


class MemoryScan(QueryNode):
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
            return None

        x = self.table[self.idx]
        self.idx += 1
        return x


PAGE_SIZE = 1024


class FileScan(QueryNode):
    """
    Yield all rows from the given heap file in our custom format!
    """
    def __init__(self, filename, schema):
        self.filename = filename
        self.schema = schema
        self.file = None
        self.buffer = None
        self.record_i = None  # `i` from print_dat
        self.end = None

    def next(self):
        if self.file is None:
            self.file = open(self.filename, 'rb')

        if self.buffer is None:
            self.buffer = self.file.read(PAGE_SIZE)
            if len(self.buffer) == 0:
                return None
            assert len(self.buffer) == PAGE_SIZE
            self.record_i = 0
            self.end = PAGE_SIZE

        num_records = int.from_bytes(self.buffer[0:2], 'little')
        start = int.from_bytes(self.buffer[(self.record_i + 1) * 2:(self.record_i + 2) * 2], 'little')
        rec = io.BytesIO(self.buffer[start:self.end])
        row = []
        for typ in self.schema:
            if typ == 'uint32':
                row.append(int.from_bytes(rec.read(4), 'little'))
            elif typ == 'text':
                n = int.from_bytes(rec.read(1))
                row.append(rec.read(n).decode('utf8'))
            else:
                raise ValueError('Unknown type')
        self.end = start

        self.record_i += 1
        if self.record_i == num_records:
            # reached last record on page
            self.buffer = None

        return tuple(row)

    def close(self):
        if self.file:
            self.file.close()


class Count(QueryNode):
    def __init__(self):
        self.done = False

    def next(self):
        if self.done:
            return None
        count = 0
        while True:
            x = self.child.next()
            if x is None:
                self.done = True
                return (count,)
            count = count + 1


class Projection(QueryNode):
    """
    Map the child records using the given map function, e.g. to return a subset
    of the fields.
    """
    def __init__(self, proj):
        self.proj = proj

    def next(self):
        x = self.child.next()
        if x is None:
            return None
        return self.proj(x)


class Selection(QueryNode):
    """
    Filter the child records using the given predicate function.

    Yes it's confusing to call this "selection" as it's unrelated to SELECT in
    SQL, and is more like the WHERE clause. We keep the naming to be consistent
    with the literature.
    """
    def __init__(self, predicate):
        self.predicate = predicate

    def next(self):
        while True:
            x = self.child.next()
            if x is None or self.predicate(x):
                return x


class Limit(QueryNode):
    """
    Return only as many as the limit, then stop
    """
    def __init__(self, n):
        self.remaining = n

    def next(self):
        if self.remaining == 0:
            return None
        self.remaining -= 1
        return self.child.next()


class Sort(QueryNode):
    """
    Sort based on the given key function
    """
    def __init__(self, key, desc=False):
        self.key = key
        self.desc = desc
        self.tuples = None
        self.idx = 0

    def next(self):
        if self.tuples is None:
            self.tuples = []
            while True:
                x = self.child.next()
                if x is None:
                    break
                self.tuples.append(x)
            self.tuples.sort(key=self.key, reverse=self.desc)

        if self.idx >= len(self.tuples):
            return None

        x = self.tuples[self.idx]
        self.idx += 1
        return x


def Q(*nodes):
    """
    Construct a linked list of executor nodes from the given arguments,
    starting with a root node, and adding references to each child
    """
    ns = iter(nodes)
    parent = root = Root()
    for n in ns:
        parent.child = n
        parent = n
    parent.child = None
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


def insert(table_filename, schema, rows):
    with open(table_filename, 'r+b') as out:
        file_size = out.seek(0, os.SEEK_END)
        if file_size == 0:
            buffer = bytearray(PAGE_SIZE)
            num_records = 0
            records_ptr = PAGE_SIZE
        else:
            assert file_size % PAGE_SIZE == 0
            out.seek(-PAGE_SIZE, os.SEEK_END)
            buffer = bytearray(out.read(PAGE_SIZE))
            out.seek(-PAGE_SIZE, os.SEEK_END)
            num_records = int.from_bytes(buffer[:2], 'little')
            records_ptr = int.from_bytes(buffer[num_records * 2:num_records * 2 + 2], 'little')
        for row in rows:
            record = io.BytesIO()
            for typ, val in zip(schema, row):
                if typ == 'uint32':
                    record.write(val.to_bytes(4, 'little'))
                elif typ == 'text':
                    bs = val.encode('utf8')
                    record.write(len(bs).to_bytes(1))
                    record.write(bs)
                else:
                    raise ValueError('Unknown type')
            record.seek(0)
            rec_bytes = record.read()
            start = records_ptr - len(rec_bytes)
            if start < (num_records + 2) * 2:
                out.write(buffer)
                buffer = bytearray(PAGE_SIZE)
                num_records = 0
                records_ptr = PAGE_SIZE
                start = records_ptr - len(rec_bytes)
            buffer[start:records_ptr] = rec_bytes
            num_records += 1
            buffer[0:2] = num_records.to_bytes(2, 'little')
            buffer[num_records * 2:(num_records + 1) * 2] = start.to_bytes(2, 'little')
            records_ptr = start
            assert len(buffer) == PAGE_SIZE
        if num_records > 0:
            out.write(buffer)
    # For now, we assume we wrote *all* rows if we got this far
    return len(rows)
