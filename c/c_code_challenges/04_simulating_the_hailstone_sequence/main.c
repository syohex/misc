#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
    if (argc < 2) {
        printf("Usage: %s num\n", argv[0]);
        return 1;
    }

    int n = atoi(argv[1]);
    printf("Enter value: %d\n", n);
    printf("Hailstone sequence:");

    int steps = 1;
    while (1) {
        printf(" %d ", n);
        if (n == 1) {
            printf("\n");
            break;
        }

        if (n % 2 == 0) {
            n /= 2;
        } else {
            n = 3 * n + 1;
        }
        ++steps;
    }
    printf("Sequence length: %d\n", steps);
    return 0;
}
