#include <assert.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static char *add_commas(int n) {
    char buffer[32];
    sprintf(buffer, "%d", n);
    int len = strlen(buffer);

    char tmp[32];
    size_t k = 0;
    for (int i = len - 1, j = 0; i >= 0; --i, ++j) {
        if (j != 0 && j % 3 == 0) {
            tmp[k++] = ',';
        }
        tmp[k++] = buffer[i];
    }

    char *ret = (char *)malloc(32);
    if (ret == NULL) {
        perror("malloc");
        return NULL;
    }

    for (int i = 0; i < k; ++i) {
        ret[i] = tmp[k - 1 - i];
    }

    ret[k] = '\0';
    return ret;
}

int main() {
    int values[] = {123, 1899, 48266, 123456, 9876543, 10100100, 5, 500000, 99000111, 83};
    char *expected[] = {"123", "1,899", "48,266", "123,456", "9,876,543", "10,100,100", "5", "500,000", "99,000,111", "83"};

    for (int i = 0; i < 10; ++i) {
        char *ret = add_commas(values[i]);
        if (strcmp(ret, expected[i]) != 0) {
            printf("## expected: %s(got: %s)\n", expected[i], ret);
            abort();
        }

        free(ret);
    }

    printf("## OK\n");
    return 0;
}
