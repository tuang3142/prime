#include <assert.h>
#include <stdio.h>

extern int fib();


int main(void) {
  assert(sum(123, 321) == 444);

  printf("OK\n");
}
