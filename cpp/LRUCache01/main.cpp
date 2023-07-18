#include <cassert>
#include <unordered_map>
#include <list>
#include <optional>
#include <cstdio>

template <typename K, typename V>
class LRUCache {
  public:
    LRUCache(size_t capacity) : capacity_(capacity) {
    }

    void Add(const K &key, const V &value) {
        if (m_.contains(key)) {
            q_.erase(m_[key]);
        }

        q_.push_front({key, value});
        m_[key] = q_.begin();

        if (m_.size() > capacity_) {
            auto &pair = q_.back();
            m_.erase(pair.first);
            q_.pop_back();
        }
    }

    std::optional<V> Get(const K &key) {
        if (!m_.contains(key)) {
            return std::nullopt;
        }

        const auto pair = *m_[key];
        q_.erase(m_[key]);
        q_.push_front(pair);
        m_[key] = q_.begin();
        return pair.second;
    }

    void Remove(const K &key) {
        q_.erase(m_[key]);
        m_.erase(key);
    }

  private:
    size_t capacity_;
    std::unordered_map<K, typename std::list<std::pair<K, V>>::iterator> m_;
    std::list<std::pair<K, V>> q_;
};

int main() {
    {
        LRUCache<int, int> c(2);
        c.Add(1, 1);
        c.Add(2, 2);
        assert(c.Get(1).value() == 1);
        c.Add(3, 3);
        assert(!c.Get(2).has_value());
        c.Add(4, 4);
        assert(!c.Get(1).has_value());
        assert(c.Get(3).value() == 3);
        assert(c.Get(4).value() == 4);
    }
    return 0;
}
