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

const (
	pageSize = 1024
	slotSize = 2
)

var schema = []string{"uint32", "text", "text"}

func encode(input, output string) error {
	inp, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	r := csv.NewReader(inp)
	_, err = r.Read() // skip header
	if err != nil {
		return fmt.Errorf("failed to read header: %v", err)
	}

	out, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // TODO: understand this?
	if err != nil {
		return fmt.Errorf("Failed to open %v: %v", out, err)
	}
	defer out.Close()

	page := make([]byte, pageSize)
	recEnd := pageSize
	recCount := 0
	slotStart := slotSize
	for record, err := r.Read(); err != nil; record, err = r.Read() {
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to read record: %v", err)
		}

		recByte := []byte{}
		for i, typ := range schema {
			switch typ {
			case "uint32":
				intVal, err := strconv.Atoi(record[i])
				if err != nil {
					return fmt.Errorf("failed to convert %v to int: %v", record[i], err)
				}
				buf := make([]byte, 4)
				binary.LittleEndian.PutUint32(buf, uint32(intVal))
				recByte = append(recByte, buf...)
			case "text":
				l := len(record[i])
				buf := make([]byte, 1)
				binary.LittleEndian.PutUint32(buf, uint32(l))
				recByte = append(recByte, buf...)
				recByte = append(recByte, []byte(record[i])...)
			default:
				return fmt.Errorf("undefined type: %v", typ)
			}
		}
		recStart := recEnd - len(recByte)
		if recStart < slotStart+slotSize {
			if _, err := out.Write(page); err != nil {
				return fmt.Errorf("failed to write to output: %v", err)
			}
			page = make([]byte, pageSize)
			recEnd = pageSize
			slotStart = slotSize
			recCount = 0

			recStart = recEnd - len(recByte)
		}
		copy(page[recStart:recEnd], recByte)

		recCount += 1
		binary.LittleEndian.PutUint16(page[:slotSize], uint16(recCount))
		binary.LittleEndian.PutUint16(page[slotStart:slotStart+slotSize], uint16(recStart))

		slotStart += 2
		recEnd = recStart
	}

	return nil
}

func lastNBytes(f *os.File, n int) {

}

// TODO: testing. This is as far as I can do. I can still study polish after this

func decode(input string) error {
	ip, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", input, err)
	}
	defer ip.Close()

	info, err := ip.Stat()
	if err != nil {
		return err
	}

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
			case "text":
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
