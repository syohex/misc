#include <assert.h>
#include <stdbool.h>

bool leapyear(int n) {
    if (n % 400 == 0) {
        return true;
    }
    if (n % 100 == 0) {
        return false;
    }

    return n % 4 == 0;
}

int main() {
    assert(leapyear(2004));
    assert(leapyear(2000));
    assert(!leapyear(2100));
    assert(leapyear(1996));
    assert(!leapyear(100));
    return 0;
}
