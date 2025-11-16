#include <stdio.h>

int main() {
  char c[] = "vamp fox held quartz duck just by wing";
  int seen[26];
  for (char* i = c; *i != '\0'; i++) {
    if ('a' <= *i && *i <= 'z') {
      printf("%c ", *i);
      seen[*i - 'a'] = 1;
    }
  }
  printf("\n");
  int cnt = 0;
  for (int i = 0; i < 26; i++) {
    cnt += seen[i];
  }
  printf("%s\n", cnt == 26 ? "true" : "false");
}
