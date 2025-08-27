package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// plan:
// open line by line, skip the first line because of heading
// for each type, encode it. int: 4 byte, string: 4 byte length + remaining

// int32, int64... hummm
var schema = []string{"uint32", "string", "string"}

func encode(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	r := csv.NewReader(f)
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to read heading: %v", err)
	}

	// TODO: wtf is this, i still dont' understand this
	out, err := os.OpenFile(encoded, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // bitwise - interesting
	if err != nil {
		return fmt.Errorf("Failed to open %v: %v", out, err)
	}
	defer out.Close()
	for range 10 { // read the first 10 line of data, will change to smt else
		record, err := r.Read()
		if err != nil {
			return fmt.Errorf("failed to read record: %v", err)
		}
		if len(record) != len(schema) {
			return fmt.Errorf("record %v doesn't match schema %v", record, schema)
		}
		for i := range schema {
			switch schema[i] {
			case "uint32":
				intVal, err := strconv.Atoi(record[i])
				if err != nil {
					fmt.Errorf("failed to convert %v to int: %v", record[i], err)
				}
				buf := make([]byte, 4)
				binary.LittleEndian.PutUint32(buf, uint32(intVal))
				if _, err := out.Write(buf); err != nil {
					return fmt.Errorf("failed to write %v to output: %v", buf, err)
				}
				fmt.Printf("printed %v to binary %v", intVal, buf)
			case "string":
				// todo
				break
			default:
				return fmt.Errorf("undefined type %v", schema[i])
			}
		}
	}

	return nil
}

// todo tmr: decode what i've just encoded
// i just need to find an effective way to deal with stress. i can deal with uncomfortable feeling, but feelling stressed all the time is not good
// talking to other peopl def. help. it is a way to socialize

func decode(file string) error {
	return fmt.Errorf("unimplemented")
}

func main() {

}
