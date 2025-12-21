#include <assert.h>
#include <stdio.h>

extern int binconv(char* bits);

int main(void) {
	printf("%d\n", binconv("0"));
	printf("%d\n", binconv("1"));
	printf("%d\n", binconv("101"));
	printf("%d\n", binconv("001"));

	assert(binconv("0") == 0);
	assert(binconv("1") == 1);
	assert(binconv("10") == 2);
	assert(binconv("11") == 3);
	assert(binconv("1000") == 8);
	assert(binconv("1111") == 15);
	assert(binconv("000000000000000000000000000000000000001") == 1);
	assert(binconv("0000000000000000000000000000000000000000000000000000000000000000000000000000000000") == 0);

	printf("OK\n");
}
