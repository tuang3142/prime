#include <assert.h>
#include <stdio.h>
#include <string.h>

extern int fib(int n);
extern int sum_n(int n);
extern int slen(char *s);

int f(int n) {
  if (n <= 1) {
    return n;
  }
  return f(n - 1) + f(n - 2);
}

int main(void) {
  assert(sum_n(0) == 0);
  assert(sum_n(1) == 1);
  assert(sum_n(5) == 15);
  assert(sum_n(10) == 55);
  assert(sum_n(100) == 5050);

  for (int i = 0; i < 35; i++) {
    assert(fib(i) == f(i));
  }

  assert(slen("1234") == 4);
  assert(slen("12345") == 5);
  assert(slen("1") == 1);
  assert(slen("") == 0);
  char str[10000];
  assert(slen(str) == 0);
  memset(str, 'A', 9999);
  str[9999] = '\0';
  assert(slen(str) == 9999);

  printf("OK\n");
}
