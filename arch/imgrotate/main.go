// Package main does smt cool
package main

import (
	"encoding/binary"
	"os"
)

func pad(n uint32) uint32 {
	return (-n * 3) % 4
}

func main() {
	data, err := os.ReadFile("testdata/1_in.bmp")
	// data, err := os.ReadFile("testdata/2_in.bmp") // stretch goal
	if err != nil {
		panic(err)
	}

	offset := binary.LittleEndian.Uint32(data[10:14])
	width := binary.LittleEndian.Uint32(data[18:22])
	height := binary.LittleEndian.Uint32(data[22:26])

	var pixels []byte
	d_padding, p_padding := pad(width), pad(height)

	// idea: draw 2D arrays on paper -> do rotation logic by hand -> profit
	// easy exercise first, stretch later (it's tricky, but worth a try)
	for py := range width {
		for px := range height {
			dy := px
			dx := width - py - 1
			d := offset + 3*(dy*width+dx)
			d += dy * d_padding // add padding for row dy to correct index of data
			pixels = append(pixels, data[d:d+3]...)
		}
		pixels = append(pixels, make([]byte, p_padding)...) // add padding for target
	}

	header := data[:offset]
	// swap w and h in header
	tmp := make([]byte, 4)
	copy(tmp, header[18:22])
	copy(header[18:22], header[22:26])
	copy(header[22:26], tmp)

	out := append(header, pixels...)

	if err := os.WriteFile("out.bmp", out, 0644); err != nil {
		panic(err)
	}
}
