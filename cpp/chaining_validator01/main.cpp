#include <cassert>
#include <vector>
#include <string>

namespace {

class Validator {
  public:
    virtual ~Validator() = default;
    virtual bool validate(const std::string &input) const = 0;
};

class ChainValidator {
  public:
    ~ChainValidator() {
        for (auto *validator : validators_) {
            delete validator;
        }
    }

    bool validate(const std::string &input) const {
        for (const auto &validator : validators_) {
            if (!validator->validate(input)) {
                return false;
            }
        }

        return true;
    }

    void addValidator(Validator *validator) {
        validators_.push_back(validator);
    }

  private:
    std::vector<Validator *> validators_;
};

class NotEmptyValidator : public Validator {
  public:
    bool validate(const std::string &input) const override {
        return !input.empty();
    }
};

class HasPrefixValidator : public Validator {
  public:
    explicit HasPrefixValidator(std::string prefix) : prefix_(std::move(prefix)) {
    }

    bool validate(const std::string &input) const override {
        return input.find(prefix_) == 0;
    }

  private:
    std::string prefix_;
};

} // namespace

int main() {
    ChainValidator cv;
    cv.addValidator(new NotEmptyValidator());
    cv.addValidator(new HasPrefixValidator("foo"));

    assert(cv.validate("fooBar"));
    assert(!cv.validate(""));
    assert(!cv.validate("bar"));
    assert(!cv.validate("FOO_BAR"));
    return 0;
}
