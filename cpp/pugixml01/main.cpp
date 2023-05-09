#include <pugixml.hpp>
#include <iostream>

int main() {
    pugi::xml_document doc;
    auto result = doc.load_file("note.xml");
    if (!result) {
        std::cerr << "Failed to parse xml: " << result.description() << std::endl;
        return 1;
    }

    auto nodes = doc.child("note");
    std::cout << nodes.name() << std::endl;

    for (auto n = nodes.first_child(); n; n = n.next_sibling()) {
        std::cout << "\t" << n.name() << std::endl;
    }

    return 0;
}
