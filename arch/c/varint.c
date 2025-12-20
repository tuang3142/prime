#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

uint64_t decode(uint8_t* arr) {
  uint64_t out = 0, shift = 0;
  for (size_t i = 0;; i++) {
    uint8_t byte = arr[i];
    uint8_t payload = byte & 0x7F;  // get the last 7 bits
    out |= (payload << shift);
    if ((byte >> 7) == 0) {  // break if continuation bit == 0
      break;
    }
    shift += 7;
  }

  return out;
}

uint8_t* encode(uint64_t n) {
  if (n < 0) return NULL;

  // find the chunk size of n
  uint64_t cp = n;
  size_t cnt = 0;
  do {
    cp >>= 7;
    cnt += 1;
  } while (cp);

  uint8_t* out = calloc(cnt, sizeof *out);
  if (!out) return NULL;

  for (size_t i = 0; i < cnt; i++) {
    uint8_t b = n & 0x7F;
    n >>= 7;
    out[i] = n ? (b | 0x80) : b;
  }

  return out;
}

int main() {
  uint8_t* p;
  p = encode(0);
  assert(p[0] == 0x00);
  p = encode(127);
  assert(p[0] == 0x7F);
  p = encode(150);
  assert(p[0] == 0x96);
  assert(p[1] == 0x01);

  printf("OK encode\n");

  uint8_t a0[] = {0x00};
  assert(decode(a0) == 0);
  uint8_t a1[] = {0x80, 0x01};
  assert(decode(a1) == 128);
  uint8_t a3[] = {0x96, 0x01};
  assert(decode(a3) == 150);
  uint8_t a2[] = {0xAC, 0x02};
  assert(decode(a2) == 300);

  for (int i = 0; i < 1 << 31; i++) {
    assert(decode(encode(i)) == i);
  }

  uint64_t big = 1ULL << 32;
  assert(decode(encode(big)) == big);
  printf("OK decode\n");
}