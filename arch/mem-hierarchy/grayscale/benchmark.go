package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// create2D creates rows x cols 2D slice filled with values.
func create2D(rows, cols int) [][]int {
	data := make([][]int, rows)
	for i := 0; i < rows; i++ {
		row := make([]int, cols)
		for j := 0; j < cols; j++ {
			row[j] = i*cols + j
		}
		data[i] = row
	}
	return data
}

// traverseRowMajor sums elements with row-major order (cache-friendly).
func traverseRowMajor(a [][]int) int64 {
	var sum int64
	for i := 0; i < len(a); i++ {
		row := a[i]
		for j := 0; j < len(row); j++ {
			sum += int64(row[j])
		}
	}
	return sum
}

// traverseColMajor sums elements with column-major order (cache-unfriendly).
func traverseColMajor(a [][]int) int64 {
	if len(a) == 0 {
		return 0
	}
	cols := len(a[0])
	rows := len(a)
	var sum int64
	for j := 0; j < cols; j++ {
		for i := 0; i < rows; i++ {
			sum += int64(a[i][j])
		}
	}
	return sum
}

// measure runs f repeatedly and reports wall time, avg per run, and allocations.
func measure(name string, f func() int64, runs int) {
	var mStart, mEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&mStart)

	start := time.Now()
	// optional CPU profile
	cpuFile, _ := os.Create(fmt.Sprintf("%s_cpu.prof", name))
	pprof.StartCPUProfile(cpuFile)
	var last int64
	for k := 0; k < runs; k++ {
		last = f()
	}
	pprof.StopCPUProfile()
	cpuFile.Close()

	elapsed := time.Since(start)

	runtime.GC()
	runtime.ReadMemStats(&mEnd)

	fmt.Printf("=== %s ===\n", name)
	fmt.Printf("Result checksum: %d\n", last)
	fmt.Printf("Total runs: %d\n", runs)
	fmt.Printf("Total time: %v\n", elapsed)
	fmt.Printf("Avg time/run: %v\n", elapsed/time.Duration(runs))
	fmt.Printf("Allocations (bytes) before: %d, after: %d, delta: %d\n",
		mStart.TotalAlloc, mEnd.TotalAlloc, int64(mEnd.TotalAlloc)-int64(mStart.TotalAlloc))
	fmt.Println()
}

func main() {
	// reasonable default: 8192x8192 might be too large for some machines.
	rows, cols := 2048, 2048
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("Matrix size: %dx%d\n\n", rows, cols)

	a := create2D(rows, cols)

	// warm-up
	_ = traverseRowMajor(a)

	// choose runs so total time is measurable
	runs := 5

	measure("row_major", func() int64 { return traverseRowMajor(a) }, runs)
	measure("col_major", func() int64 { return traverseColMajor(a) }, runs)
}

