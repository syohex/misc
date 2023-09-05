#include <assert.h>
#include <stdio.h>
#include <string.h>

void center_text(int width, const char *text) {
    size_t len = strlen(text);
    size_t spaces = 0;
    if (len <= width) {
        spaces = (width - len) / 2;
    }

    for (size_t i = 0; i < spaces; ++i) {
        printf(" ");
    }
    printf("%s\n", text);
}

int main() {
    char *titles[] = {
        "March Sales",
        "My First Project",
        "Centering output is so much fun!",
        "This title is very long, just to see whether the code can handle such a long title",
    };

    for (int i = 0; i < 4; ++i) {
        center_text(80, titles[i]);
    }
    return 0;
}
