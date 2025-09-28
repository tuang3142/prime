package varint

import (
	"fmt"
	"testing"
)

func Test_toBin(t *testing.T) {
	tcs := []struct {
		n    uint64
		want string
	}{
		{150, "10010110"},
		{21, "10101"},
		{1, "1"},
		{2, "10"},
		{3, "11"},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprint(tc.n), func(t *testing.T) {
			got := toBin(tc.n)
			if got != tc.want {
				t.Errorf("toBin(%v) = %v, want %v", tc.n, got, tc.want)
			}
		})
	}
}

func Test_pad(t *testing.T) {
	tcs := []struct {
		name string
		bits string
		want string
	}{
		{name: "pad_one", bits: "1000001", want: "01000001"},
		{name: "pad_two", bits: "100001", want: "00100001"},
		{name: "pad_none", bits: "10000001", want: "10000001"},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := pad(tc.bits, 8)
			if got != tc.want {
				t.Errorf("addPadding(%v) = %v, want %v", tc.bits, got, tc.want)
			}
		})
	}
}

func Test_binToHex(t *testing.T) {
	tcs := []struct {
		bits string
		want string
	}{
		{"00001101", "0d"},
		{"11111111", "ff"},
		{"00101010", "2a"},
		{"10010110", "96"},
	}

	for _, tc := range tcs {
		t.Run(tc.bits, func(t *testing.T) {
			got := binToHex(tc.bits)
			if got != tc.want {
				t.Errorf("toHex(%v) = %v, want %v", tc.bits, got, tc.want)
			}
		})
	}
}

func Test_Encode(t *testing.T) {
	val := uint64(150)
	want := "9601"
	got := Encode(val)

	if got != want {
		t.Errorf("Encode(%v) = %v, want %v", val, got, want)
	}
}

func Test_Decode(t *testing.T) {
	tests := []struct {
		hex  string
		want uint64
	}{
		{"01", 1},     // single-byte
		{"9601", 150}, // multi-byte
		{"ac02", 300}, // larger multi-byte
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			got := Decode(tt.hex)
			if got != tt.want {
				t.Errorf("Decode(%v) = %v, want %v", tt.hex, got, tt.want)
			}
		})
	}
}

func Test_EncodeThenDecode_LotsOfTime(t *testing.T) {
	for i := 0; i < 1048576; i++ {
		if uint64(i) != Decode(Encode(uint64(i))) {
			t.Fatalf("failed on %v", i)
		}
	}
}
