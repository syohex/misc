#include <assert.h>
#include <stdio.h>
#include <string.h>

static char *ordinal(int v) {
    if (v >= 11 && v <= 13) {
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

int main(void) {
    assert(strcmp(ordinal(51), "st") == 0);
    assert(strcmp(ordinal(82), "nd") == 0);
    assert(strcmp(ordinal(33), "rd") == 0);
    assert(strcmp(ordinal(40), "th") == 0);

    assert(strcmp(ordinal(11), "th") == 0);
    assert(strcmp(ordinal(12), "th") == 0);
    assert(strcmp(ordinal(13), "th") == 0);

    for (int i = 4; i <= 9; ++i) {
        assert(strcmp(ordinal(i), "th") == 0);
    }

    for (int i = 1; i <= 20; ++i) {
        printf("%d%s\n", i, ordinal(i));
    }

    return 0;
}