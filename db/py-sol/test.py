import unittest
import tempfile
import os

from import_csv import import_csv
from executor import (
    run,
    Q,
    Count,
    Limit,
    Projection,
    Selection,
    Sort,
    MemoryScan,
    FileScan,
    insert,
)

birds = (
    ('amerob', 'American Robin', 0.077, True),
    ('baleag', 'Bald Eagle', 4.74, True),
    ('eursta', 'European Starling', 0.082, True),
    ('barswa', 'Barn Swallow', 0.019, True),
    ('ostric1', 'Ostrich', 104.0, False),
    ('emppen1', 'Emperor Penguin', 23.0, False),
    ('rufhum', 'Rufous Hummingbird', 0.0034, True),
    ('comrav', 'Common Raven', 1.2, True),
    ('wanalb', 'Wandering Albatross', 8.5, False),
    ('norcar', 'Northern Cardinal', 0.045, True)
)


def rq(*nodes):
    return tuple(run(Q(*nodes)))


class TestInMemory(unittest.TestCase):

    def test_memory_scan(self):
        self.assertEqual(rq(MemoryScan(birds)), birds)

    def test_limit(self):
        # Limit to any number of records. Limiting to 0 returns empty; Limiting to > max returns all
        for i in range(len(birds) + 2):
            self.assertEqual(rq(Limit(i), MemoryScan(birds)), birds[:i])

    def test_count(self):
        # Correct count, including after limit
        for i in range(len(birds) + 1):
            self.assertEqual(rq(Count(), Limit(i), MemoryScan(birds))[0], (i,))

    def test_selection(self):
        # Select everything
        self.assertEqual(rq(Selection(lambda x: True), MemoryScan(birds)), birds)
        # Select nothing
        self.assertEqual(rq(Selection(lambda x: False), MemoryScan(birds)), tuple())
        # Select non-US birds
        self.assertEqual(rq(Selection(lambda x: not x[3]), MemoryScan(birds)), (
            ('ostric1', 'Ostrich', 104.0, False),
            ('emppen1', 'Emperor Penguin', 23.0, False),
            ('wanalb', 'Wandering Albatross', 8.5, False),
        ))

    def test_projection(self):
        # Identity projection x -> x
        self.assertEqual(rq(Projection(lambda x: x), MemoryScan(birds)), birds)
        # Trivial projection
        self.assertEqual(rq(Projection(lambda x: ('hello',)), MemoryScan(birds)), len(birds) * (('hello',),))
        # Project to two fields
        self.assertEqual(rq(Projection(lambda x: (x[0], x[3])), MemoryScan(birds)), (
            ('amerob', True),
            ('baleag', True),
            ('eursta', True),
            ('barswa', True),
            ('ostric1', False),
            ('emppen1', False),
            ('rufhum', True),
            ('comrav', True),
            ('wanalb', False),
            ('norcar', True)
        ))
        # Selection then projection: ids of non-US birds
        self.assertEqual(rq(
            Projection(lambda x: (x[0],)),
            Selection(lambda x: not x[3]),
            MemoryScan(birds)
        ), (
            ('ostric1',),
            ('emppen1',),
            ('wanalb',),
        ))
        # Projection, then selection, should also work
        self.assertEqual(rq(
            Selection(lambda x: not x[1]),
            Projection(lambda x: (x[0], x[3])),
            MemoryScan(birds)
        ), (
            ('ostric1', False),
            ('emppen1', False),
            ('wanalb', False),
        ))

    def test_sort(self):
        # id and weight of 3 heaviest birds
        self.assertEqual(rq(
            Projection(lambda x: (x[0], x[2])),
            Limit(3),
            Sort(lambda x: x[2], desc=True),
            MemoryScan(birds),
        ), (
            ('ostric1', 104.0),
            ('emppen1', 23.0),
            ('wanalb', 8.5),
        ))


class TestFileInsertAndScan(unittest.TestCase):
    """
    Insert records into a new file, using our native file format,
    and read them back again
    """
    def setUp(self):
        _, self.temp_file = tempfile.mkstemp()

    def tearDown(self):
        os.unlink(self.temp_file)

    def test_all(self):
        test_schema = ('uint32', 'text')
        # Write one record and read it back again
        insert(self.temp_file, test_schema, ((1, 'foo'),))
        self.assertEqual(rq(FileScan(self.temp_file, test_schema)), ((1, 'foo'),))

        # Write a subsequent and read BOTH back again
        insert(self.temp_file, test_schema, ((2, 'bar'),))
        self.assertEqual(rq(FileScan(self.temp_file, test_schema)), (
            (1, 'foo'),
            (2, 'bar'),
        ))
        self.assertEqual(rq(Count(), FileScan(self.temp_file, test_schema)), (
            (2,),
        ))

        # - Write 100s of records (bulk insert) across a page boundary and count all
        insert(self.temp_file, test_schema, ((3, 'etc'),) * 101)
        self.assertEqual(rq(FileScan(self.temp_file, test_schema)), (
            (1, 'foo'),
            (2, 'bar'),
        ) + ((3, 'etc'),) * 101)

        # - Write yet another record and count all
        insert(self.temp_file, test_schema, ((4, 'final'),))
        self.assertEqual(rq(Count(), FileScan(self.temp_file, test_schema)), (
            (104,),
        ))


# @unittest.skip('Skipping movielens tests')
class TestMovielens(unittest.TestCase):
    """
    A set of tests using MovieLens 20m dataset, including import from original CSV.

    Skip if you don't have the files by uncommenting the @unittest.skip decorator above
    """
    def setUp(self):
        _, self.temp_file = tempfile.mkstemp()

    def tearDown(self):
        os.unlink(self.temp_file)

    def test_movies(self):
        movies_schema = ('uint32', 'text', 'text')
        num_written = import_csv(
            os.path.expanduser('~/Downloads/ml-20m/movies.csv'),
            self.temp_file,
            schema=movies_schema
        )
        self.assertEqual(num_written, 27278)
        self.assertEqual(rq(Count(), FileScan(self.temp_file, movies_schema)), ((27278,),))

        # Find by id
        self.assertEqual(rq(
            Projection(lambda x: (x[1],)),
            Selection(lambda x: x[0] == 5000),
            FileScan(self.temp_file, movies_schema)
        ), (('Medium Cool (1969)',),))


if __name__ == '__main__':
    unittest.main()
