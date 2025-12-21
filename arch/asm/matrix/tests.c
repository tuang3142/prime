#include <assert.h>
#include <stdio.h>

extern int matrix_get(int* ma, int m, int n, int i, int j);

// int matrix_get(int* ma, int m, int n, int i, int j) {
// 	return *(*(ma + i) + j);
// }

// how do i balance out speed vs understanding?
// i want to understand but i am progressing slowly
// how can i move faster
// minimum understanding: just try to solve problem
// (sort of) maxium: learn watever interest me, sometimes it is kind of a
// distraction
// i think i should set a limit: finish the program
// set a target during the week. solve todo during the weekend (1 sitting)
// yeah that could work. treat them like extra1. everyting else move to extra 2, do them when I finish the overal course.
// yeah, could be fun that way
// need a way to organize this - obs - obsediant
// between p1 and p2: prioritize problem solving (stretch exercise)
// feeling the progress dude. 
// know my pace (e.g. 1 problem = 1h)

// Contiguous typed data: values of the same type stored back-to-back in memory with no gaps (except alignment). Key points:
//
// - Arrays: T a[N] allocates N objects of type T contiguously. Address of a[i] = (char*)a + i * sizeof(T).
// - Multidimensional arrays: T a[R][C] is contiguous: elements laid out row-major — a[i][j] at offset (i*C + j)*sizeof(T).
// - Pointer arithmetic uses element size: for T* p, p + k advances k * sizeof(T) bytes.
// - Decay: array expressions often decay to pointer-to-first-element — T a[N] → T* when passed to functions.
// - int** is NOT a pointer to contiguous 2D int matrix; it's a pointer to pointers (often non-contiguous). int (*)[C] or int* (flat) represent contiguous 2D storage.
// - memcpy/mmap/read can copy raw bytes; when treating as typed data, respect alignment and object lifetimes.
// - Advantages: cache locality, simple indexing math (i*cols + j).
// - Pitfalls: casting between incompatible pointer types (e.g., int[][C] ↔ int**) is undefined behavior.

int main(void) {
	int m1[1][4] = {{1,2,3,4}};

	assert(matrix_get((int*)m1, 1, 4, 0, 0) == 1);
	assert(matrix_get((int*)m1, 1, 4, 0, 1) == 2);
	assert(matrix_get((int*)m1, 1, 4, 0, 2) == 3);
	assert(matrix_get((int*)m1, 1, 4, 0, 3) == 4);

    printf("OK\n");
}
