#include <span>
#include <vector>
#include <string>
#include <cstdio>

void print(std::span<char> s) {
    printf("[");
    for (char c : s) {
        printf(" %c ", c);
    }
    printf("]\n");
}

int main() {
    std::vector<char> v{'h', 'e', 'l', 'l', 'o'};
    std::string s("hello world");
    print(std::span(v));
    print(std::span(s));
    return 0;
}
