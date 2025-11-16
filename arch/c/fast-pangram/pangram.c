#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

bool ispangram(char* s) {
  int seen[26] = {0};
  for (; *s != '\0'; s++) {
    if ('a' <= *s && *s <= 'z') {
      seen[*s - 'a']++;
    }
    if ('A' <= *s && *s <= 'Z') {
      seen[*s - 'A']++;
    }
  }
  int cnt = 0;
  for (int i = 0; i < 26; i++) {
    cnt += seen[i];
  }
  return cnt == 26;
}

int main() {
  size_t len;
  size_t read;
  // ssize_t read;
  char* line = NULL;
  while ((read = getline(&line, &len, stdin)) != -1) {
    if (ispangram(line)) printf("%s", line);
  }

  if (ferror(stdin)) fprintf(stderr, "Error reading from stdin");

  free(line);
  fprintf(stderr, "ok\n");
}
