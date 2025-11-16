#include <assert.h>
#include <stdio.h>

int bitcount(unsigned int n) {
	int cnt = 0;
	while (n > 0) {
		cnt++;
		n &= n - 1;
	}
	return cnt;
}

int main() {
  assert(bitcount(0) == 0);
  assert(bitcount(1) == 1);
  assert(bitcount(2) == 1);
  assert(bitcount(3) == 2);
  assert(bitcount(8) == 1);

  assert(bitcount(0b11111111) == 8);
  assert(bitcount(0b11001101) == 5);

  int b = bitcount(0xffffffff);
  assert(b == 32);

  printf("OK\n");
}
