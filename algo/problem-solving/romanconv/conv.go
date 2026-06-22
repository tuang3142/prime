package main

import (
	"errors"
	"fmt"
)

var num []int = {1, 5, 10, 50, 100, 500, 1000}
var letter []string = {"I", "V", "X", "L", "C", "D", "M"}


func convHelper(int d) {
	// var base int
	// var unit int
	// switch d {
	// case 1 <= d && d <= 5:
	// 	base, unit = 0, 0
	//
	// case 5 <= d && d <= 10:
	// 	base, unit = 0
	// }

	if base > d {
		return unit + base
	}

	return base + unit
}

func conv(int n) (string, error) {
	if n >= 4000 {
		return "", errors.New("unsupported")
	}
	pow := 1
	for n > 0 {
		digit = (n % 10) * pow
		n /= 10
		pow *= 10
		ret = ret + convHelper(digit)
	}
	return ret, nil
}

func main() {
	tc := []struct {
		int    n
		string want
	}{
		{0, "N"}, // Nulla
		{18, "XVIII"},
		{371, "CCCLXXI"},
		{3999, "MMMCMXCIX"},
	}

	for t, _ := range tc {
		got := conv(t.n)
		if got != t.want {
			fmt.Printf("conv(%d) = %s, want %s", t.n, got, want)
		}
	}
	print("OK")
}
