#include <stdio.h>

long make_change(double dollar) {
    long cents = dollar * 100;
    long ret = 0;
    long coins[] = {100, 50, 25, 10, 5, 1};
    for (size_t i = 0; i < 6; ++i) {
        ret += cents / coins[i];
        cents %= coins[i];
    }

    return ret;
}

int main() {
    double dollars[] = {0.49, 1.27, 0.75, 1.31, 0.83};
    for (size_t i = 0; i < 5; ++i) {
        printf("$%g = %ld coins\n", dollars[i], make_change(dollars[i]));
    }

    return 0;
}
