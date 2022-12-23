#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>

static void print_board(int size, int row, int col) {
    for (int i = 0; i < size; ++i) {
        printf("#");
    }
    printf("\n");
    for (int i = 0; i < size; ++i) {
        for (int j = 0; j < size; ++j) {
            if (i == row && col == j) {
                printf("K");
            } else {
                printf(".");
            }
        }
        printf("\n");
    }
    for (int i = 0; i < size; ++i) {
        printf("#");
    }
    printf("\n");
}

int move() {
    return rand() % 3 - 1;
}

int main() {
    const int size = 8;
    int row = 5;
    int col = 4;

    srand((unsigned int)time(NULL));

    for (int i = 1;; ++i) {
        print_board(size, row, col);

        while (true) {
            int r = move();
            int c = move();
            if (r == 0 && c == 0) {
                continue;
            }

            row += r;
            col += c;
            break;
        }

        if (!(row >= 0 && row < size && col >= 0 && col < size)) {
            printf("## finished %d turns\n", i);
            break;
        }
    }

    return 0;
}
