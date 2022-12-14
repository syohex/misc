#include <stdio.h>
#include <string.h>
#include <stdlib.h>

char *left(const char *s, size_t len) {
    char *p = (char *)malloc(len + 1);

    for (size_t i = 0; i < len && s[i] != '\0'; ++i) {
        p[i] = s[i];
    }

    p[len] = '\0';
    return p;
}

char *right(const char *s, size_t len) {
    char *p = (char *)malloc(len + 1);
    size_t length = strlen(s);
    size_t start = length - len;

    for (size_t i = 0; i < len && s[i] != '\0'; ++i) {
        p[i] = s[start + i];
    }

    p[len] = '\0';
    return p;
}

char *mid(const char *s, size_t offset, size_t len) {
    char *p = (char *)malloc(len + 1);
    for (size_t i = 0; i < len && s[i] != '\0'; ++i) {
        p[i] = s[offset - 1 + i];
    }

    p[len] = '\0';
    return p;
}

int main() {
    char text[] = "Once upon a time, there was a string";

    char *p1 = left(text, 16);
    printf("## left(16)=%s\n", p1);

    char *p2 = right(text, 16);
    printf("## right(16)=%s\n", p2);

    char *p3 = mid(text, 13, 11);
    printf("## mid(13)=%s\n", p3);

    free(p1);
    free(p2);
    free(p3);
    return 0;
}
