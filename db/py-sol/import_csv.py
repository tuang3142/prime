import argparse
import csv

from executor import insert

STR2PY = {
    'text': str,
    'uint32': int,
}


def import_csv(source, target, skip_header=True, schema=None):
    """
    Bulk import the entirety of the CSV to the target data file
    """
    # Touch target: fine to append if already exists
    with open(target, 'a') as f:
        pass

    # Read all from the CSV and just run through our insert function
    with open(source, 'r') as f:
        reader = csv.reader(f)
        if skip_header:
            next(reader)

        # TODO do this in chunks in case it doesn't fit in memory
        content = tuple(reader)
        if schema is None:
            schema = ('text',) * len(content[0])

        rows = [
            [STR2PY[typ](val) for typ, val in zip(schema, c)]
            for c in content
        ]

        return insert(target, schema, rows)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Convert from a CSV file to a heap file')
    parser.add_argument('source')
    parser.add_argument('target')
    parser.add_argument('-s', '--schema', help='comma separated string of column types e.g \'uint32,text,txt\'')
    parser.add_argument('--noheader', action='store_true', help='don\'t skip the first line')
    args = parser.parse_args()

    count = import_csv(
        args.source,
        args.target,
        skip_header=not args.noheader,
        schema=args.schema.split(',') if args.schema else None
    )
    print(f'Inserted {count} rows')
