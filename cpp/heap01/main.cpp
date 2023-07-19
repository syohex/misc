#include <cassert>
#include <vector>

class Heap {
  public:
    void Push(int x) {
        heap_.push_back(x);

        int i = heap_.size() - 1;
        while (i > 0) {
            int parent = (i - 1) / 2;
            if (heap_[parent] >= x) {
                break;
            }

            heap_[i] = heap_[parent];
            i = parent;
        }

        heap_[i] = x;
    }

    int Top() {
        if (heap_.empty()) {
            return -1;
        }

        return heap_[0];
    }

    void Pop() {
        if (heap_.empty()) {
            return;
        }

        int v = heap_.back();
        heap_.pop_back();

        int i = 0;
        int len = heap_.size();
        while (i * 2 + 1 < len) {
            int child1 = i * 2 + 1;
            int child2 = i * 2 + 2;
            if (child2 < len && heap_[child2] > heap_[child1])  {
                child1 = child2;
            }

            if (heap_[child1] <= v) {
                break;
            }

            heap_[i] = heap_[child1];
            i = child1;
        }

        heap_[i] = v;
    }

  private:
    std::vector<int> heap_;
};

int main() {
    Heap h;
    h.Push(1);
    h.Push(2);
    h.Push(3);
    h.Push(4);

    assert(h.Top() == 4);
    h.Pop();
    assert(h.Top() == 3);
    h.Pop();
    assert(h.Top() == 2);
    h.Pop();
    assert(h.Top() == 1);
    return 0;
}
