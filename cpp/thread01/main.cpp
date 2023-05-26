#include <thread>
#include <functional>
#include <iostream>

struct Thread {
    Thread(int a, int val) : a_(a) {
        t_ = std::thread(std::bind(&Thread::Method, this, std::placeholders::_1), val);
    }

    ~Thread() {
        t_.join();
    }

    void Method(int val) {
        std::cout << "a=" << a_ << ", arg=" << val << std::endl;
    }

    int a_;
    std::thread t_;
};

int main() {
    Thread t(42, 99);
    return 0;
}
