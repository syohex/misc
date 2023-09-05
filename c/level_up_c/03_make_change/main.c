#include <stdio.h>

// 0.49 1.27

int make_change(double value) {
    long amount = value * 100;
    int ret = 0;

    while (amount >= 25) {
        amount -= 25;
        ++ret;
    }

    while (amount >= 10) {
        amount -= 10;
        ++ret;
    }

    while (amount >= 5) {
        amount -= 10;
        ++ret;
    }

    return ret + amount;
}

int main(void) {
    double amount[] = {0.49, 1.27, 0.75, 1.31, 0.83};

    for (int i = 0; i < 5; ++i) {
        printf("## changes(%g) = %d\n", amount[i], make_change(amount[i]));
    }

    return 0;
}
