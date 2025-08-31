package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// plan:
// open line by line, skip the first line because of heading
// for each type, encode it. int: 4 byte, string: 4 byte length + remaining

// int32, int64... hummm
var schema = []string{"uint32", "string", "string"}

func encode(input, output string) error {
	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	r := csv.NewReader(f)
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to read heading: %v", err)
	}

	out, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open %v: %v", out, err)
	}
	defer out.Close()
	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
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
					return fmt.Errorf("failed to convert %v to int: %v", record[i], err)
				}
				buf := make([]byte, 4)
				binary.LittleEndian.PutUint32(buf, uint32(intVal))
				if _, err := out.Write(buf); err != nil {
					return fmt.Errorf("failed to write %v to output: %v", intVal, err)
				}
				fmt.Printf("printed %v to binary %#v\n", intVal, buf)
			case "string":
				l := len(record[i])
				buf := make([]byte, 4)
				binary.LittleEndian.PutUint32(buf, uint32(l))
				if _, err := out.Write(buf); err != nil {
					return fmt.Errorf("failed to write %v to output: %v", l, err)
				}
				if _, err := out.Write([]byte(record[i])); err != nil {
					return fmt.Errorf("failed to write %v to output: %v", record[i], err)
				}
				fmt.Printf("printed %v to file\n", record[i])
				// fmt.Printf("byte: %#v", []byte(record[i])) - encode: char -> utf8 (int - id of sort) -> hex
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

func decode(input, output string) error {
	ip, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", input, err)
	}
	defer ip.Close()
	for {
		for _, typ := range schema {
			switch typ {
			case "uint32":
				buf := make([]byte, 4)
				// read the next 4 byte
				n, err := io.ReadFull(ip, buf)
				if err != nil {
					if err == io.EOF {
						return nil
					}
					return fmt.Errorf("failed to read uint32 %v", err, n)
				}
				fmt.Printf("%v ", binary.LittleEndian.Uint32(buf))
				// conver to in
			case "string":
				buf := make([]byte, 4)
				_, err := io.ReadFull(ip, buf)
				if err != nil {
					if err == io.EOF {
						return nil
					}
					return fmt.Errorf("failed to read uint32 (lenght) %v", err)
				}
				l := binary.LittleEndian.Uint32(buf)
				buf = make([]byte, l)
				_, err = io.ReadFull(ip, buf)
				if err != nil {
					if err == io.EOF {
						return nil
					}
					return fmt.Errorf("failed to read string %v", err)
				}
				fmt.Printf("%v ", string(buf))
			default:
				return fmt.Errorf("undefined type %v", typ)
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	input := "movies.csv"
	output := "movies.dat"
	if err := encode(input, output); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("encoding completed!\n")

	input = "movies.csv"
	output = "movies.dat"
	if err := decode(input, output); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("decoding completed!\n")
}
