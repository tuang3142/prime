#include <stdio.h>
#include <stdlib.h>

typedef struct Node {
  int val;
  struct Node* next;  // IMPORTANT: use 'struct Node*', not just 'Node*'
} Node;

Node* Node_new() {
  Node* node = malloc(sizeof *node);
  return node;
}

void Node_append(Node* n, int val) {
  Node* root = n;
  while (n && n->next) {
    n = n->next;
  }
  n->next = Node_new();
  n->next->val = val;
}

int main() {
}