package varint

import (
	"fmt"
	"strings"
)

const (
	byteSize = 8
	hexChar  = "0123456789abcdef"
)

func Encode(n uint64) (out string) {
	if n == 0 {
		return "00"
	}

	bin := toBin(n)
	size := byteSize - 1
	for i := len(bin); i > 0; i -= size {
		var byt string
		if i-size > 0 {
			byt = "1" + bin[i-size:i]
		} else {
			byt = "0" + pad(bin[:i], size)
		}
		out += binToHex(byt)
	}

	return out
}

func Decode(hex string) (n uint64) {
	if hex == "00" {
		return 0
	}
	hex = strings.ToLower(hex)
	hexToBinMap := map[rune]string{
		'0': "0000", '1': "0001", '2': "0010", '3': "0011",
		'4': "0100", '5': "0101", '6': "0110", '7': "0111",
		'8': "1000", '9': "1001", 'a': "1010", 'b': "1011",
		'c': "1100", 'd': "1101", 'e': "1110", 'f': "1111",
	}
	bin := ""
	for i := 0; i < len(hex); i += 2 {
		b := hexToBinMap[rune(hex[i])] + hexToBinMap[rune(hex[i+1])]
		bin = b[1:] + bin // skip the MSBit
	}

	for bin[0] == '0' {
		bin = bin[1:]
	}

	for i := 0; i < len(bin); i++ {
		n = (n << 1) | uint64(bin[i]-'0')
	}

	return n
}

func binToHex(bits string) (out string) {
	for i := 0; i+4 <= len(bits); i += 4 {
		val := 0
		for _, bit := range bits[i : i+4] {
			val = val<<1 + int(bit-'0')
		}
		out += string(hexChar[val])
	}

	return out
}

func toBin(n uint64) (bin string) {
	for {
		if n == 1 {
			bin = fmt.Sprint(n) + bin
			break
		}
		bin = fmt.Sprint(n%2) + bin
		n /= 2
	}
	return bin
}

func pad(bits string, size int) string {
	for len(bits)%size != 0 {
		bits = "0" + bits
	}
	return bits
}
