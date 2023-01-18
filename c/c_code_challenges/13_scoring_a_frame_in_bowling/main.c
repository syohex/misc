#include <stdio.h>
#include <stdlib.h>
#include <time.h>

typedef struct {
    char ball1;
    char ball2;
    int throw1;
    int throw2;
    int score;
} Frame;

int throw_ball(int pins) {
    return rand() % (pins + 1);
}

Frame create_frame(int pins) {
    Frame f;
    f.throw1 = throw_ball(pins);

    if (f.throw1 == 10) {
        f.ball1 = 'X';
        f.ball2 = ' ';
        f.throw2 = 0;
        f.score = 10;
    } else {
        f.ball1 = f.throw1 == 0 ? '-' : '0' + f.throw1;
        f.throw2 = throw_ball(pins - f.throw1);
        if (f.throw1 + f.throw2 == 10) {
            f.ball2 = '/';
        } else {
            f.ball2 = '0' + f.throw2;
        }
        f.score = f.throw1 + f.throw2;
    }

    return f;
}

int main() {
    srand(time(NULL));

    Frame f = create_frame(10);
    printf("| %c|%c|\n", f.ball1, f.ball2);
    printf("|%3d |\n", f.score);

    return 0;
}
