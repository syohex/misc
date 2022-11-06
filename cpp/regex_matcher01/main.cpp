#include <cassert>
#include <string>

namespace {

bool match_here(const char *regexp, const char *text);

bool match_star(int c, const char *regexp, const char *text) {
    do {
        if (match_here(regexp, text)) {
            return true;
        }
    } while (*text != '\0' && (*text++ == c || c == '.'));

    return false;
}

bool match_here(const char *regexp, const char *text) {
    if (regexp[0] == '\0') {
        return true;
    }

    if (regexp[1] == '*') {
        return match_star(regexp[0], regexp + 2, text);
    }

    if (regexp[0] == '$' && regexp[1] == '\0') {
        return *text == '\0';
    }

    if (*text != '\0' && (regexp[0] == '.' || regexp[0] == *text)) {
        return match_here(regexp + 1, text + 1);
    }

    return false;
}

bool match(const char *regexp, const char *text) {
    if (regexp[0] == '^') {
        return match_here(regexp + 1, text);
    }

    do {
        if (match_here(regexp, text)) {
            return true;
        }
    } while (*text++ != '\0');

    return 0;
}

} // namespace

int main() {
    assert(match("^apple$", "apple"));
    assert(match("^apple", "apple orange melon"));
    return 0;
}
