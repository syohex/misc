#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
    /* Shakespear's 18th Sonnet */
    static char *text = "Shall I compare thee to a summer's day? \
Thou art more lovely and more temperate:\n\
Rough winds do shake the darling buds of May, \
And summer's lease hath all too short a date;\n\
Sometime too hot the eye of heaven shines, \
And often is his gold complexion dimm'd;\n\
And every fair from fair sometime declines, \
By chance or nature's changing course untrimm'd;\n\
But thy eternal summer shall not fade, \
Nor lose possession of that fair thou ow'st;\n\
Nor shall death brag thou wander'st in his shade, \
When in eternal lines to time thou grow'st:\n\
So long as men can breathe or eyes can see, \
So long lives this, and this gives life to thee.";

    int width = 40;
    if (argc >= 2) {
        width = atoi(argv[1]);
        if (width < 16 || width > 100) {
            width = 40;
        }
    }

    char *start = text;
    char *end = text;
    while (*end != '\0') {
        if (*end == '\n') {
            for (char *p = start; p != end; ++p) {
                putchar(*p);
            }
            putchar('\n');

            ++end;
            start = end;
        } else {
            ++end;
            if (end == start + width) {
                while (*end != ' ') {
                    --end;
                    if (start == end) {
                        end += width;
                        break;
                    }
                }

                for (char *p = start; p != end; ++p) {
                    putchar(*p);
                }
                if (*end != '\n') {
                    putchar('\n');
                }
                ++end;
                start = end;
            }
        }
    }

    for (char *p = start; p != end; ++p) {
        putchar(*p);
    }
    if (*end != '\n') {
        putchar('\n');
    }

    return 0;
}
