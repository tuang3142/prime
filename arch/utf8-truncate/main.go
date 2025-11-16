package main

import (
	"fmt"
	"os"
)

func truncate(line []byte, n int) []byte {
	if n >= len(line) {
		return line
	}

	for i := n; i > 0 && ((line[i] & 0x0c) == 0x80); i-- {
	}

	return line[:n]
}

func main() {
	f, err := os.ReadFile("in")
	if err != nil {
		panic(err)
	}

	var line []byte
	for _, b := range f {
		if b != '\n' {
			line = append(line, b)
		}
		if b == '\n' {
			line = truncate(line[1:], int(line[0]))
			fmt.Printf("%s\n", line)
			line = []byte{} // reset
		}
	}
}
