#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

# TODO: need to learn about asccii text and 0x1f?
int ispangram(char* s) {
  uint32_t bitmap = 0;
  char c;
  while ((c=*s++) != '\0') {
    if (c < '@') {
      continue;
    }
    int shift = c - (c >= 'a' ? 'a' : 'A');
    bitmap |= (1 << shift);
  }

  return bitmap == 0x3FFFFFF;
}

int main() {
  size_t len;
  size_t read;
  char* line = NULL;
  while ((read = getline(&line, &len, stdin)) != -1) {
    if (ispangram(line)) {
      printf("%s", line);
    }
  }

  if (ferror(stdin)) {
    fprintf(stderr, "Error reading from stdin");
  }

  free(line);
  fprintf(stderr, "ok\n");
}
