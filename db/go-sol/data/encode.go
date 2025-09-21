package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func Encode(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", fileName, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	if headers, err := r.Read(); err != nil {
		return fmt.Errorf("failed to read headers: %v", err)
	} else {
		fmt.Printf("%v\n", headers)
	}

	out := [][]string{}
	for row, err := r.Read(); ; row, err = r.Read() {
		o := []string{}
		if err != nil {
			if err == io.EOF {
				fmt.Println("Read completed!")
				for _, o := range out {
					fmt.Println(o)
				}
				return nil
			}
			return fmt.Errorf("failed to read: %v", err)
		}
		for i, typ := range schema {
			switch typ {
			case "uint32":
				o = append(o, row[i])
			case "text":
				o = append(o, row[i])
			default:
				return fmt.Errorf("unknown type: %v", typ)
			}

		}
		out = append(out, o)
	}
}
