#include <assert.h>
#include <stdio.h>

extern int sum(int a, int b);


int main(void) {
  assert(sum(123, 321) == 444);

  printf("OK\n");
}
