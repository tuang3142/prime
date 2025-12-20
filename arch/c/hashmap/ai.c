#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define STARTING_BUCKETS 8
#define MAX_KEY_SIZE 8

typedef struct Node {
  void* val;
  char* key;
  struct Node* next;
} Node;

Node* Node_new(char* key, void* val) {
  Node* n = malloc(sizeof *n);
  n->key = key;
  n->val = val;
  return n;
}

void Node_free(Node* n) { free(n); }

// continue hasmap
typedef struct {
  Node** arr;
} Hashmap;

int hash_func(char* s) {
  int out = 0;
  for (; *s != '\0'; s++) {
    out += *s;
  }
  return out % STARTING_BUCKETS;
}

Hashmap* Hashmap_new() {
  Hashmap* h = malloc(sizeof *h);
  Node** arr = malloc(STARTING_BUCKETS * (sizeof(*arr)));
  h->arr = arr;
  return h;
}

void Hashmap_set(Hashmap* h, char* key, void* val) {
  int k = hash_func(key);
  Node* new_node = Node_new(strdup(key), val);
  Node* root = h->arr[k];
  do {
    if (!root) {
      h->arr[k] = new_node;
      return;
    }
    if (strcmp(root->key, key) == 0) {
      root->val = val;
      return;
    }
    if (!root->next) {
      root->next = new_node;
      return;
    }
    root = root->next;
  } while (1);
}

void* Hashmap_get(Hashmap* h, char* key) {
  int k = hash_func(key);
  Node* root = h->arr[k];
  while (root && strcmp(root->key, key)) {
    root = root->next;
  }
  return root ? root->val : NULL;
}

void Hashmap_delete(Hashmap* h, char* key) {
  int k = hash_func(key);
  Node* prev = h->arr[k];
  if (!prev) return;
  if (!strcmp(prev->key, key)) {
    h->arr[k] = prev->next;
    Node_free(prev);
    return;
  }
  Node* nxt = prev->next;
  while (nxt && strcmp(nxt->key, key)) {
    nxt = nxt->next;
  }
  if (nxt) {
    prev->next = nxt->next;
    free(nxt);
  }
}

void Hashmap_free(Hashmap* h) {
  free(h->arr);
  free(h);
}

int main() {
  Hashmap* h = Hashmap_new();

  // basic get/set functionality
  int a = 5;
  double b = 7.17;  // TODO: float, TODO: floating point in C
  Hashmap_set(h, "item a", &a);
  Hashmap_set(h, "item b", &b);
  assert(Hashmap_get(h, "item a") == &a);
  assert(Hashmap_get(h, "item b") == &b);

  // using the same key should override the previous value
  long c = 20;
  Hashmap_set(h, "item a", &c);
  assert(Hashmap_get(h, "item a") == &c);

  // basic delete functionality
  Hashmap_delete(h, "item a");
  assert(Hashmap_get(h, "item a") == NULL);

  // check handle collision
  assert(hash_func("ab") == hash_func("ba"));
  Hashmap_set(h, "ab", &a);
  Hashmap_set(h, "ba", &b);
  assert(Hashmap_get(h, "ab") == &a);
  assert(Hashmap_get(h, "ba") == &b);
  Hashmap_delete(h, "ab");
  assert(Hashmap_get(h, "ab") == NULL);
  assert(Hashmap_get(h, "ba") == &b);

  // handle collisions correctly
  // note: this doesn't necessarily test expansion
  int i, n = STARTING_BUCKETS * 10, ns[n];
  char key[MAX_KEY_SIZE];
  for (i = 0; i < n; i++) {
    ns[i] = i;
    sprintf(key, "item %d", i);
    Hashmap_set(h, key, &ns[i]);
  }
  for (i = 0; i < n; i++) {
    sprintf(key, "item %d", i);
    if (Hashmap_get(h, key) != &ns[i]) {
      printf("%s\n", key);
      return 0;
    }
  }

  Hashmap_free(h);

  printf("ok\n");
}

/*
   stretch goals:
   - expand the underlying array if we start to get a lot of collisions
   - support non-string keys
   - try different hash functions
   - switch from chaining to open addressing
   - use a sophisticated rehashing scheme to avoid clustered collisions
   - implement some features from Python dicts, such as reducing space use,
   maintaing key ordering etc. see
   https://www.youtube.com/watch?v=npw4s1QTmPg for ideas
   */