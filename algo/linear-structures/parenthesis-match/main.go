package main

import "fmt"

func IsBalanced(s string) bool {
	stack := make([]rune, 0, len(s))
	pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
	for _, r := range s {
		switch r {
		case '(', '[', '{':
			stack = append(stack, r)
		case ')', ']', '}':
			if len(stack) == 0 || stack[len(stack)-1] != pairs[r] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func assertEqual(got bool, want bool, msg string) {
	if got != want {
		panic(fmt.Sprintf("assertion failed: %s — got %v, want %v", msg, got, want))
	}
}

func main() {
	assertEqual(IsBalanced(""), true, `empty string`)
	assertEqual(IsBalanced("()"), true, `simple balanced "()"`)
	assertEqual(IsBalanced("([]{()})"), true, `nested mixed "([]{()})"`)
	assertEqual(IsBalanced("({[)"), false, `unbalanced missing close "({[)"`)
	assertEqual(IsBalanced("())"), false, `extra closing "())"`)

	println("all assertions passed")
}
