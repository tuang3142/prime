package main

import "fmt"

func mergeSort(arr []int) {
	var sort func(l, r int)
	sort = func(l, r int) {
		if l >= r {
			return
		}
		m := (l + r) / 2
		sort(l, m)
		sort(m+1, r)

		var tmp []int
		i, j := l, m+1
		for i <= m && j <= r {
			if arr[i] < arr[j] {
				tmp = append(tmp, arr[i])
				i++
			} else {
				tmp = append(tmp, arr[j])
				j++
			}
		}
		for i <= m {
			tmp = append(tmp, arr[i])
			i++
		}
		for j <= r {
			tmp = append(tmp, arr[j])
			j++
		}
		i, j = l, 0
		for ; i <= r; i++ {
			arr[i] = tmp[j]
			j++
		}
	}

	sort(0, len(arr)-1)
}

func quickSort(arr []int) {
	var sort func(int, int)
	sort = func(l, r int) {
		if l >= r {
			return
		}
		m := (l + r) / 2
		p := arr[m]
		for i, j := l, r; i < j; {
			for i < j && arr[i] <= p {
				i++
			}
			for j > i && arr[j] > p {
				j--
			}
			arr[i], arr[j] = arr[j], arr[i]
			i, j = i+1, j-1
		}
		sort(l, m)
		sort(m+1, r)
	}
	sort(0, len(arr)-1)
}

func main() {
	arr := []int{5, 4, 3, 1, 2}
	quickSort(arr)
	for _, n := range arr {
		fmt.Printf("%d ", n)
	}
	print("\n")
}
