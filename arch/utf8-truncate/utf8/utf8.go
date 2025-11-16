package main

// UTF-8 can be up to 4 bytes, covering up to 21 bits
// U+07FF (11 bits) -> 2 bytes (1 continuation)
// U+FFFF (16 bits) -> 3 bytes (2 continuations)
// U+10FFFF (21 bits) -> 4 bytes (3 continuations)
func encode(c int) []int {
	if c <= 0b1111111 { // 0 to 127
		return []int{c}
	}
	cont := []int{}
	// Stop when c is reduced to the bits needed for the leading byte (not for the tail)
	for {
		part := 0b10000000 | (c & 0b111111) // take the last 6 bits (10xxxxxx)
		cont = append([]int{part}, cont...)
		c >>= 6 // Shift off the 6 bits

		// Break condition: c is small enough for the leading byte (what I was wrong earlier: c is big enough to fit the tail
		// For 2-byte (5 remaining bits): c <= 0b11111 (31)
		if len(cont) == 1 && c <= 0b11111 { // 5 remaining bits
			break
		}
		// For 3-byte (4 remaining bits): c <= 0b1111 (15)
		if len(cont) == 2 && c <= 0b1111 { // 4 remaining bits
			break
		}
		// For 4-byte (3 remaining bits): c <= 0b111 (7)
		if len(cont) == 3 && c <= 0b111 { // 3 remaining bits
			break
		}
	}

	lead := map[int]int{
		1: 0b11000000,
		2: 0b11100000,
		3: 0b11110000,
	}[len(cont)] | c

	return append([]int{lead}, cont...)
}

func decode(b []int) int {
	switch {
	case b[0]&0b10000000 == 0:
		return int(b[0]) // 1-byte (ASCII)

	case b[0]&0b11100000 == 0b11000000:
		return (int(b[0]&0b00011111) << 6) |
			int(b[1]&0b00111111)

	case b[0]&0b11110000 == 0b11100000:
		return (int(b[0]&0b00001111) << 12) |
			(int(b[1]&0b00111111) << 6) |
			int(b[2]&0b00111111)

	case b[0]&0b11111000 == 0b11110000:
		return (int(b[0]&0b00000111) << 18) |
			(int(b[1]&0b00111111) << 12) |
			(int(b[2]&0b00111111) << 6) |
			int(b[3]&0b00111111)
	}
	return 0
}

// rune = int32
func main() {
	chars := []rune{
		'A', 'z', // 1-byte (ASCII)
		'ñ', 'ø', // 2-byte
		'你', 'Ж', // 3-byte
		'🙂', '𐍈', // 4-byte
		'ß', 'Ж', 'Δ', '你', 'あ', '한', 'א', 'ع', '𐍈', 'ệ', 'ê', 'ấ',
	}
	for _, ch := range chars {
		if int(ch) != decode(encode(int(ch))) {
			print(ch, " diff\n")
			return
		}
	}
	print("ok :)\n")
}
