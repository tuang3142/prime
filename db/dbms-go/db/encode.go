// i am a capable person and gritty. i can improve myself. i'wil start with studying everyday for at
// least 30 minutes, afterwork, following course at csprimer.com
// todo: watch the video again, why encode that way?
package main

import (
	"encoding/binary"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

const csvFile = "movies.csv"
const encoded = "movies.dat"

var schema = []string{"int", "string", "string"}

func encode() {
	in, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open %v: %v", csvFile, err)
	}
	defer in.Close()

	reader := csv.NewReader(in)

	if _, err := reader.Read(); err != nil { // skip the heading
		log.Fatalf("Failed to read csv: %v", err)
	}

	out, err := os.OpenFile(encoded, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // bitwise - interesting
	if err != nil {
		log.Fatalf("Failed to open %v: %v", out, err)
	}
	defer out.Close()

	// why encode this way...
	// encode each column depends on the type (int32, string, string)
	for range 10 { // testing, read the first 10 line
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF { // end of file
				break
			}
			log.Fatalf("Failed to read csv: %v", err)
		}
		for i, stringVal := range row {
			dataType := schema[i]
			switch dataType {
			case "int": // int or uint32
				intVal, err := strconv.Atoi(stringVal)
				if err != nil {
					log.Fatalf("Failed to convert %v to int", stringVal)
				}
				buf := make([]byte, 4)
				binary.LittleEndian.PutUint32(buf, uint32(intVal)) // why uint32

				if _, err := out.Write(buf); err != nil {
					log.Fatalf("Failed to write %v: %v", buf, err)
				}
				log.Printf("Written %v as %v", stringVal, buf)
			case "string":
				// first write the length of the string, for decoding purpose
				// use only 1 byte (8 bits) to store the length tho
				buf := make([]byte, 1) // 8 bit, store up to 2^8
				binary.LittleEndian.PutUint32(buf, uint32(len(stringVal)))
				if _, err := out.Write(buf); err != nil {
					log.Fatalf("Failed to write %v: %v", buf, err)
				}
				// write string into n byte, however long the string is
				if _, err := out.Write([]byte(stringVal)); err != nil {
					log.Fatalf("Failed to write %v: %v", stringVal, err)
				}
				log.Printf("Written %v as %v", stringVal, buf)
				// when storing string, a text, it needs more than 4 byte

			default:
				log.Fatalf("Undefined type %v", dataType)
			}
		}
	}
}

func decode() {
	out, err := os.Open(encoded)
	if err != nil {
		log.Fatalf("Failed to open %v: %v", encoded, err)
	}
	defer out.Close()

	for range 8 { // being conservative
		for _, typ := range schema {
			switch typ {
			case "int":
				buf := make([]byte, 4)
				io.ReadFull(out, buf)
				println(binary.LittleEndian.Uint32(buf))
			case "string":
				buf := make([]byte, 1) // read the length of the string
				io.ReadFull(out, buf)
				l := binary.LittleEndian.Uint32(buf)
				println(l)
				buf = make([]byte, l)
				io.ReadFull(out, buf)
				println(string(buf))
			default:
				log.Fatalf("Undefined type %v", typ)
			}
		}
	}
}

func main() {
	encode() // good so far, i need to check if i have written this correctly however
	// this is every good habit builder. keep this up. i need a streak.
	decode() // ok there is some error but i am smarter a little bit // good.
}
