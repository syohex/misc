#include <stdio.h>

double golden_ratio(double a, double b) {
    if (b != 0) {
        return a + (1 / golden_ratio(a, b - 1));
    } else {
        return a;
    }
}

int main() {
    double ret = golden_ratio(1.0, 15);
    printf("## a=%g, b=%g, golden_ration=%g\n", 1.0, 15.0, ret);
    return 0;
}
