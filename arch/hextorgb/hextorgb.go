package hextorgb

import "strings"

func HexToRGB(h string) []int {
	h = strings.ToLower(h)
	if len(h) == 3 || len(h) == 4 {
		h = double(h)
	}
	if len(h) == 6 {
		return sixDigit(h)
	}
	if len(h) == 8 {
		return eightDigit(h)
	}
	panic("wrong hex format")
}

func double(h string) (ret string) {
	for _, r := range h {
		ret += string(r) + string(r)
	}
	return ret
}

func sixDigit(h string) []int {
	return []int{
		hexToInt(h[:2]),
		hexToInt(h[2:4]),
		hexToInt(h[4:6]),
	}
}

func eightDigit(h string) []int {
	return []int{
		hexToInt(h[:2]),
		hexToInt(h[2:4]),
		hexToInt(h[4:6]),
		hexToInt(h[6:8]),
	}
}

var mp = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7,
	'8': 8, '9': 9, 'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15,
}

func hexToInt(h string) int {
	fi := mp[rune(h[0])]
	se := mp[rune(h[1])]

	return fi*16 + se
}
