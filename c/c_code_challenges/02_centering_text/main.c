#include <stdio.h>
#include <string.h>

void center_text(size_t width, char *text) {
    size_t len = strlen(text);
    if (len >= width) {
        puts(text);
        return;
    }

    size_t prefix_spaces = (width / 2) - (len / 2);
    size_t i = 0;
    for (; i < prefix_spaces; ++i) {
        printf(" ");
    }

    printf("%s", text);
    i += len;
    for (; i < width; ++i) {
        printf(" ");
    }

    printf("\n");
}

int main() {
    char *title[] = {
        "March Sales",
        "My First Project",
        "Centering output is so much fun!",
        "This title is very long, just to see whether the code can handle such a long title",
        NULL,
    };

    for (int i = 0; title[i] != NULL; ++i) {
        center_text(80, title[i]);
    }
    return 0;
}
