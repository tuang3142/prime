#include <assert.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>

extern double cone_volume(double a, double b);
extern double mul(double a, double b);

int main(void) {
  assert(fabsf(mul(0, 0) - 0) < 1e-6f);
  assert(fabsf(mul(12, 3) - 36) < 1e-6f);
  assert(fabsf(mul(2.0f, 3.5f) - 7.0f) < 1e-6f);
  assert(fabsf(mul(-1.25f, 4.0f) + 5.0f) < 1e-6f);

  assert(fabsf(cone_volume(0.0f, 10.0f) - 0.0f) < 1e-6f);
  assert(fabsf(cone_volume(1.0f, 1.0f) - (M_PI * 1.0f * 1.0f * 1.0f / 3.0f)) <
         1e-6f);
  assert(fabsf(cone_volume(2.5f, 4.0f) - (M_PI * 2.5f * 2.5f * 4.0f / 3.0f)) <
         1e-5f);
  printf("OK\n");
}
