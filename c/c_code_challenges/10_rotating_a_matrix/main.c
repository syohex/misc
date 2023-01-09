#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>

static void print_matrix(char matrix[5][5], size_t size) {
    for (size_t i = 0; i < size; ++i) {
        printf("[");
        for (size_t j = 0; j < size; ++j) {
            printf(" %c ", matrix[i][j]);
        }
        printf("]\n");
    }
}

static void rotating_matrix(char matrix[5][5], size_t size) {
    for (size_t i = 0; i < size; ++i) {
        for (size_t j = i; j < size; ++j) {
            char tmp = matrix[i][j];
            matrix[i][j] = matrix[j][i];
            matrix[j][i] = tmp;
        }
    }

    for (size_t i = 0; i < size; ++i) {
        for (size_t j = 0; j < size / 2; ++j) {
            char tmp = matrix[i][j];
            matrix[i][j] = matrix[i][size - j - 1];
            matrix[i][size - j - 1] = tmp;
        }
    }
}

void check(char matrix[5][5], char expected[5][5], size_t size) {
    for (size_t i = 0; i < size; ++i) {
        for (size_t j = i; j < size; ++j) {
            if (matrix[i][j] != expected[i][j]) {
                printf("unexpected matrix[%zd][%zd]=%c(expected: %c)\n", i, j, matrix[i][j], expected[i][j]);
                exit(1);
            }
        }
    }
}

int main(void) {
    // clang-format off
    char matrix[5][5] = {
        {'c', 'b', 'v', 't', 'b'},
        {'t', 'e', 'g', 'v', 'm'},
        {'p', 'v', 'v', 'k', 'a'},
        {'h', 'm', 'z', 'o', 'i'},
        {'x', 'u', 'v', 't', 't'},
    };
    char expected[5][5] = {
        {'x', 'h', 'p', 't', 'c'},
        {'u', 'm', 'v', 'e', 'b'},
        {'v', 'z', 'v', 'g', 'v'},
        {'t', 'o', 'k', 'v', 't'},
        {'t', 'i', 'a', 'm', 'b'},
    };
    // clang-format on

    printf("## Original\n");
    print_matrix(matrix, 5);
    rotating_matrix(matrix, 5);

    printf("## Rotated\n");
    print_matrix(matrix, 5);

    check(matrix, expected, 5);
    return 0;
}
