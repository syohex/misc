#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
    size_t rows;
    size_t cols;
    size_t size;
    char **data;
} matrix_t;

static matrix_t init_matrix(size_t rows, size_t cols) {
    matrix_t ret = {
        .rows = rows,
        .cols = cols,
        .data = NULL,
    };

    ret.size = rows >= cols ? rows : cols;
    ret.data = malloc(sizeof(char *) * ret.size);
    for (size_t i = 0; i < ret.size; ++i) {
        ret.data[i] = malloc(sizeof(char) * ret.size);
    }

    return ret;
}

static void destroy_matrix(matrix_t *matrix) {
    for (size_t i = 0; i < matrix->size; ++i) {
        free(matrix->data[i]);
    }
    free(matrix->data);
}

static void print_matrix(matrix_t *matrix) {
    for (size_t i = 0; i < matrix->rows; ++i) {
        printf("[");
        for (size_t j = 0; j < matrix->cols; ++j) {
            printf(" %c ", matrix->data[i][j]);
        }
        printf("]\n");
    }
}

static void rotating_matrix(matrix_t *matrix) {
    size_t tmp = matrix->rows;
    matrix->rows = matrix->cols;
    matrix->cols = tmp;

    for (size_t i = 0; i < matrix->size; ++i) {
        for (size_t j = i; j < matrix->size; ++j) {
            char tmp = matrix->data[i][j];
            matrix->data[i][j] = matrix->data[j][i];
            matrix->data[j][i] = tmp;
        }
    }

    for (size_t i = 0; i < matrix->size; ++i) {
        for (size_t j = 0; j < matrix->size / 2; ++j) {
            char tmp = matrix->data[i][j];
            size_t k = matrix->cols - j - 1;
            matrix->data[i][j] = matrix->data[i][k];
            matrix->data[i][k] = tmp;
        }
    }
}

int main(void) {
    {
        matrix_t m = init_matrix(5, 5);
        memcpy(m.data[0], "cbvtb", 5);
        memcpy(m.data[1], "tegvm", 5);
        memcpy(m.data[2], "pvvka", 5);
        memcpy(m.data[3], "hmzoi", 5);
        memcpy(m.data[4], "xuvtt", 5);

        printf("## Original\n");
        print_matrix(&m);
        rotating_matrix(&m);

        printf("## Rotated 90\n");
        print_matrix(&m);

        printf("## Rotated 180\n");
        rotating_matrix(&m);
        print_matrix(&m);

        destroy_matrix(&m);
    }
    {
        matrix_t m = init_matrix(3, 4);
        memcpy(m.data[0], "1234", 4);
        memcpy(m.data[1], "2345", 4);
        memcpy(m.data[2], "3456", 4);

        printf("## Original\n");
        print_matrix(&m);
        rotating_matrix(&m);

        printf("## Rotated 90\n");
        print_matrix(&m);

        printf("## Rotated 180\n");
        rotating_matrix(&m);
        print_matrix(&m);

        destroy_matrix(&m);
    }

    return 0;
}
