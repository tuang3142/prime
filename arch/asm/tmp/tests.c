#include <assert.h>
#include <stdio.h>

extern char foo(char* s);

int main(void) {
  printf("%d\n", foo("1234"));
  printf("OK\n");
}
