// Package main provides implementation and test for binary search algorithm.
// It includes two variant: recurrsive and non-recurrsive.
package main

import "fmt"

// bsRecurrsive returns the index of `target` in `arr`, or -1 if not found.
// It is assumed that `arr` is sorted ASC.
func bsRecurrsive(arr []int, target int) int {
	var f func(int, int) int
	f = func(l, r int) int {
		if l > r {
			return -1
		}
		m := (l + r) / 2
		if arr[m] == target {
			return m
		}
		if arr[m] > target {
			r = m - 1
		} else if arr[m] < target {
			l = m + 1
		}
		return f(l, r)
	}
	l, r := 0, len(arr)-1
	return f(l, r)
}

func main() {
	arr := []int{2, 4, 6, 8, 10, 11, 13, 14, 100}
	if got, want := bsRecurrsive(arr, 2), 0; got != want {
		fmt.Printf("bsRecurrsive(%v, 2) = %v, want %v\n", arr, got, want)
		return
	}
	if got, want := bsRecurrsive(arr, 100), 8; got != want {
		fmt.Printf("bsRecurrsive(%v, 100) = %v, want %v\n", arr, got, want)
		return
	}
	if got, want := bsRecurrsive(arr, 11), 5; got != want {
		fmt.Printf("bsRecurrsive(%v, 11) = %v, want %v\n", arr, got, want)
	}
	if got, want := bsRecurrsive(arr, 15), -1; got != want {
		fmt.Printf("bsRecurrsive(%v, 15) = %v, want %v\n", arr, got, want)
		return
	}

	arr = []int{}
	if got, want := bsRecurrsive(arr, 17), -1; got != want {
		fmt.Printf("bsRecurrsive(%v, 17) = %v, want %v\n", arr, got, want)
		return
	}

	println("ok")
}
