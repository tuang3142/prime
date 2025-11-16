#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

// ensure safe memory leak or smt?
const int INIT_SIZE = 10;  // could be 1?

typedef struct DA {
  int* arr;
  int sz;
} DA;

DA* New() {  // does it need to be new(void)
  // not size of DA*?
  DA* da = malloc(sizeof(DA));
  if (!da) return NULL;

  // why not int*, watch supliment video
  da->arr = calloc(INIT_SIZE, sizeof(int));
  if (!da->arr) {
    free(da);
    return NULL;
  }

  da->sz = 0;
};

void Append(DA* da, int n) {
  da->arr[da->sz] = n;
  // todo: alloc
  da->sz++;
}

int Pop(DA* da) {
  // do we need to "free" this memoery address?
  int n = da->arr[da->sz - 1];
  da->sz--;
  return n;
}

int Get(DA* da, int i) {
  if (i >= da->sz) return -1;

  return da->arr[i];
}

void Set(DA* da, int i, int n) {
  if (i >= da->sz) return;

  da->arr[i] = n;
}

int Size(DA* da) { return da->sz; }

void Free(DA* da) { free(da); }

int main() {
  DA* da = New();
  assert(Size(da) == 0);

  Append(da, 10);
  Append(da, 20);
  Append(da, 30);
  assert(Size(da) == 3);

  assert(Get(da, 0) == 10);
  assert(Get(da, 1) == 20);
  assert(Get(da, 2) == 30);

  int p = Pop(da);
  assert(p == 30);
  assert(Size(da) == 2);

  Set(da, 0, 777);
  assert(Get(da, 0) == 777);
  assert(Size(da) == 2);

  Free(da);

  printf("OK\n");
}