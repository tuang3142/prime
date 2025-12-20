#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

// TODO: ensure no memeory leak with valgrind
#define STARTING_CAPACITY 8
#define nptr ((void*)0)

typedef struct {
  void** items;
  int len;
  int cap;
} DA;

DA* New() {
  DA* da = malloc(sizeof(DA));
  if (!da) return NULL;

  da->items = malloc(STARTING_CAPACITY * sizeof(void*));
  if (!da->items) {
    free(da);
    return NULL;
  }

  da->len = 0;
  da->cap = STARTING_CAPACITY;

  return da;
};

void Free(DA* da) {
  free(da->items);
  free(da);
}

void Append(DA* da, void* x) {
  if (da->len >= da->cap) {
    da->cap *= 2;
    da->items = realloc(da->items, da->cap * sizeof(void*));
  }
  da->items[da->len] = x;
  da->len++;
}

void* Pop(DA* da) {
  if (da->len == 0) {
    return NULL;
  }

  void* x = da->items[da->len - 1];
  da->len--;

  if (da->len < da->cap / 4) {
    da->cap /= 2;
    da->items = realloc(da->items, da->cap * sizeof(void*));
  }

  return x;
}

void* Get(DA* da, int i) {
  if (i >= da->len) return NULL;

  return da->items[i];
}

void Set(DA* da, void* x, int i) {
  if (i >= da->len) return;

  da->items[i] = x;
}

int Size(DA* da) { return da->len; }

int main() {
  DA* da = New();

  assert(Size(da) == 0);

  // basic push and pop test
  int x = 5;
  float y = 12.4;
  Append(da, &x);
  Append(da, &y);
  assert(Size(da) == 2);

  assert(Pop(da) == &y);
  assert(Size(da) == 1);

  assert(Pop(da) == &x);
  assert(Size(da) == 0);
  assert(Pop(da) == NULL);

  // basic set/get test
  Append(da, &x);
  Set(da, &y, 0);
  assert(Get(da, 0) == &y);
  Pop(da);
  assert(Size(da) == 0);

  // expansion test
  DA* da2 = New();  // use another DA to show it doesn't get overriden
  Append(da2, &x);
  int i, n = 100 * STARTING_CAPACITY, arr[n];
  for (i = 0; i < n; i++) {
    arr[i] = i;
    Append(da, &arr[i]);
  }
  assert(Size(da) == n);
  for (i = 0; i < n; i++) {
    assert(Get(da, i) == &arr[i]);
  }
  for (; n; n--) Pop(da);
  assert(Size(da) == 0);
  assert(Pop(da2) == &x);  // this will fail if da doesn't expand

  Free(da);
  Free(da2);
  printf("OK\n");
}

// TODO: test this
// typedef struct {
//   int x;
//   float y;
// } Foo;
// Foo* food = malloc(10 * sizeof(Foo));
// for (int i = 0; i < 10; i++) {
//   food[i]->x = i;
//   food[i]->y = float(i) * 0.2;
// }
