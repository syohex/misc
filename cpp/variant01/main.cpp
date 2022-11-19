#include <iostream>
#include <variant>
#include <vector>

struct Dog {
    void say() const {
        std::cout << "I'm a dog" << std::endl;
    }
};

struct Cat {
    void say() const {
        std::cout << "I'm a cat" << std::endl;
    }
};

struct SayVisitor {
    void operator()(const Dog &dog) {
        dog.say();
    }

    void operator()(const Cat &cat) {
        cat.say();
    }
};

using Animal = std::variant<Dog, Cat>;

int main() {
    std::vector<Animal> animals;
    animals.push_back(Dog{});
    animals.push_back(Cat{});

    for (const auto &animal : animals) {
        std::visit(SayVisitor{}, animal);
    }

    return 0;
}
