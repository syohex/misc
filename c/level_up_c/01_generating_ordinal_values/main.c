#include <stdio.h>
#include <assert.h>
#include <string.h>

const char *ordinal(int v) {
    if (v == 11 || v == 12 || v == 13) {
        return "th";
    }

    switch (v % 10) {
    case 1:
        return "st";
    case 2:
        return "nd";
    case 3:
        return "rd";
    default:
        return "th";
    }
}

int main() {
    assert(strcmp(ordinal(11), "th") == 0);
    assert(strcmp(ordinal(12), "th") == 0);
    assert(strcmp(ordinal(13), "th") == 0);
    assert(strcmp(ordinal(1), "st") == 0);
    assert(strcmp(ordinal(2), "nd") == 0);
    assert(strcmp(ordinal(3), "rd") == 0);

    for (int i = 4; i <= 10; ++i) {
        assert(strcmp(ordinal(i), "th") == 0);
    }

    printf("ok\n");
    return 0;
}
