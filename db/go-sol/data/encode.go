package data

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

var schema = []string{"uint32", "text", "text"}

func Encode(inpFile, outFile string) error {
	f, err := os.Open(inpFile)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", inpFile, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	if headers, err := r.Read(); err != nil {
		return fmt.Errorf("failed to read headers: %v", err)
	} else {
		fmt.Printf("%v\n", headers)
	}

	out := make([]byte, 0)
	for {
		row, err := r.Read()
		if err == io.EOF {
			fmt.Printf("read completed! writing to %v...\n", outFile)
			if werr := os.WriteFile(outFile, out, 0644); werr != nil {
				return fmt.Errorf("failed to write to %v: %v", outFile, werr)
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to read: %v", err)
		}

		if len(row) != len(schema) {
			return fmt.Errorf("row has %d fields, want %d", len(row), len(schema))
		}

		for i, typ := range schema {
			switch typ {
			case "uint32":
				u, err := strconv.ParseUint(row[i], 10, 32)
				if err != nil {
					return fmt.Errorf("failed to parse %q as uint32: %v", row[i], err)
				}
				var buf [4]byte
				binary.BigEndian.PutUint32(buf[:], uint32(u))
				out = append(out, buf[:]...)

			case "text":
				b := []byte(row[i]) // UTF-8 bytes
				var lenbuf [4]byte
				binary.BigEndian.PutUint32(lenbuf[:], uint32(len(b)))
				out = append(out, lenbuf[:]...)
				out = append(out, b...)

			default:
				return fmt.Errorf("unknown type: %v", typ)
			}
		}
	}
}

type Row []any

// thanks chatgpt but test tmr pls
// also just watch video next time (1h)
func DecodeFile(path string) ([]Row, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var rows []Row
	for {
		row := make(Row, 0, len(schema))
		for _, typ := range schema {
			switch typ {
			case "uint32":
				var buf [4]byte
				if _, err := io.ReadFull(f, buf[:]); err != nil {
					return nil, err
				}
				val := binary.BigEndian.Uint32(buf[:])
				row = append(row, val)
			case "text":
				var lenBuf [4]byte
				if _, err := io.ReadFull(f, lenBuf[:]); err != nil {
					return nil, err
				}
				size := binary.BigEndian.Uint32(lenBuf[:])
				b := make([]byte, size)
				if _, err := io.ReadFull(f, b); err != nil {
					return nil, err
				}
				row = append(row, string(b))

			default:
				return nil, fmt.Errorf("unknown type: %q", typ)
			}
		}
		rows = append(rows, row)
	}
}
