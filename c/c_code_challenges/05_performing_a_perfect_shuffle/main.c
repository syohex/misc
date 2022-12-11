#include <stdio.h>
#include <string.h>

void perfect_shuffle() {
    char orig[] = "abcdefghijklmnopqrstuvwxyz";
    char cards[26];

    memcpy(cards, orig, 26);

    for (int i = 1;; ++i) {
        printf("[%d] %s\n", i, cards);

        char tmp[26];
        int k = 0, m = 13;
        for (int j = 0; j < 26; ++j) {
            if (j % 2 == 0) {
                tmp[j] = cards[k];
                ++k;
            } else {
                tmp[j] = cards[m];
                ++m;
            }
        }

        memcpy(cards, tmp, 26);

        if (memcmp(orig, cards, 26) == 0) {
            printf("[%d] %s\n", i, cards);
            break;
        }
    }
}

int main() {
    perfect_shuffle();
    return 0;
}
