// i am a capable person and gritty. i can improve myself. i'wil start with studying everyday for at
// least 30 minutes, afterwork, following course at csprimer.com
// todo: watch the video again, why encode that way?
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

const csvFile = "movies.csv"

var schema = []string{"int", "string", "string"}

func main() {
	f, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open %v: %v", csvFile, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	if _, err := reader.Read(); err != nil { // skip the heading
		log.Fatalf("Failed to read csv: %v", err)
	}

	// why encode this way...
	// encode each column depends on the type (int32, string, string)
	for i := 0; i < 10; i++ { // testing, read the first 10 line
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF { // end of file
				break
			}
			log.Fatalf("Failed to read csv: %v", err)
		}
		for j, val := range row {
			dataType := schema[j]
			switch dataType {
			case "int":
				v, err := strconv.Atoi(val)
				if err != nil {
					log.Fatalf("Failed to convert %v to int", val)
				}
				buf := make([]byte, 4)
				binary.LittleEndian.PutUint32(buf, v) // todo: continue from here. i need to learn to love programming again.
			case "text":
			default:
				log.Fatalf("Undefined type %v", dataType)
			}

			fmt.Println(s)
		}
	}

	// write it to a file (movie.dat)
}
