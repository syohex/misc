#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <string.h>
#include <stdbool.h>

int main() {
    char *playlist[] = {"Like a Rolling Stone",
                        "Satisfaction",
                        "Imagine",
                        "What's Going On",
                        "Respect",
                        "Good Vibrations",
                        "Johnny B. Goode",
                        "Hey Jude",
                        "What'd I Say",
                        "Smells Like Teen Spirit",
                        "My Generation",
                        "Yesterday",
                        "Blowin' in the Wind",
                        "Purple Haze",
                        "London Calling",
                        "I Want to Hold Your Hand",
                        "Maybellene",
                        "Hound Dog",
                        "Let It Be",
                        "A Change Is Gonna Come"};
    size_t length = sizeof(playlist) / sizeof(playlist[0]);
    size_t played[15];
    for (size_t i = 0; i < 15; ++i) {
        played[i] = length;
    }

    srand((unsigned int)time(NULL));

    puts("Playlist:");

    size_t played_index = 0;
    for (int i = 0; i < 100; ++i) {
        while (true) {
            size_t index = rand() % length;
            size_t p = played_index;
            bool ok = true;

            for (size_t j = 0; j < 15; ++j) {
                if (index == played[(p + j) % 15]) {
                    ok = false;
                    break;
                }
            }

            if (ok) {
                printf("  %s [%d]\n", playlist[index], i + 1);
                played[played_index] = index;
                played_index = (played_index + 1) % 15;
                break;
            }
        }
    }

    return 0;
}
