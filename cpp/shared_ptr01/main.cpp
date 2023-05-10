#include <memory>
#include <iostream>

struct Test {
    Test() = default;
    ~Test() {
        std::cout << "destructor" << std::endl;
    }
};

int main() {
    std::shared_ptr<Test> a(new Test());
    auto b = a;

    std::cout << "Call reset1" << std::endl;
    a.reset();

    std::cout << "Call reset2" << std::endl;
    b.reset();

    return 0;
}
