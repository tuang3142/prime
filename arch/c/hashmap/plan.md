here is how i will go about implementing hashmap:

- implement simple hash func
- implement storing only int
- handle collision
- implement more sofishticted has func
- implement storing void?

next sessin: finish hash map, submit it; fast pangram

fix (wut):

The key array is created on the stack (locally inside the function/loop).

In each iteration of the loop, key is overwritten with a new string (e.g., "item 0", "item 1", etc.).

The original Hashmap_set function simply takes the char* key argument and stores this pointer in the Node: n->key = key;.

Crucially, the pointer stored in the hashmap node is the address of the local char key[MAX_KEY_SIZE] buffer on the stack.

After the loop finishes, or even on the next iteration, the content of this single stack buffer (key) changes.

When you run the Hashmap_get loop, every single node's key pointer points to the same final string that was written to the stack buffer (e.g., "item 79"). When you call Hashmap_get for "item 10", it searches the chain, but every node thinks its key is "item 79" (or whatever the last value was).

The hashmap needs to store a copy of the key, not just the pointer to a temporary buffer.