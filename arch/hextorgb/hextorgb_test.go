package hextorgb

import "testing"

func Test_HexToRGB_Simple(t *testing.T) {
	testCases := []struct {
		hexString string
		wantRGB   []int
	}{
		{"00ff00", []int{0, 255, 0}},
		{"fe030a", []int{254, 3, 10}},
		{"0f0def", []int{15, 13, 239}},
		{"0F0DEF", []int{15, 13, 239}},
		{"0000cc", []int{0, 0, 204}},
	}

	for _, tc := range testCases {
		t.Run(tc.hexString, func(t *testing.T) {
			gotRGB := HexToRGB(tc.hexString)

			for i, val := range gotRGB {
				if val != tc.wantRGB[i] {
					t.Fatalf("HexToRGB(%s)=%v, want %v", tc.hexString, gotRGB, tc.wantRGB)
				}
			}
		})
	}
}

func Test_HexToRGB_Advance(t *testing.T) {
	testCases := []struct {
		hexString string
		wantRGB   []int
	}{
		{"0000FFC0", []int{0, 0, 255, 192}},
		{"123", []int{17, 34, 51}},
		{"00f8", []int{0, 0, 255, 136}},
	}

	for _, tc := range testCases {
		t.Run(tc.hexString, func(t *testing.T) {
			gotRGB := HexToRGB(tc.hexString)

			for i, val := range gotRGB {
				if val != tc.wantRGB[i] {
					t.Fatalf("HexToRGB(%s)=%v, want %v", tc.hexString, gotRGB, tc.wantRGB)
				}
			}
		})
	}
}
