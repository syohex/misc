#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <assert.h>

struct Heap {
    int *data;
    size_t items;
    size_t capacity;
};

struct Heap *HeapInit(size_t capacity) {
    struct Heap *h = malloc(sizeof(struct Heap));
    if (h == NULL) {
        return NULL;
    }

    h->data = malloc(capacity * sizeof(int));
    if (h->data == NULL) {
        free(h);
        return NULL;
    }

    h->items = 0;
    h->capacity = capacity;

    return h;
}

void HeapAdd(struct Heap *heap, int data) {
    heap->data[heap->items] = data;
    ++heap->items;

    int index = heap->items - 1;
    while (index > 0) {
        int parent = (index - 1) / 2;
        if (heap->data[parent] >= data) {
            break;
        }

        heap->data[index] = heap->data[parent];
        index = parent;
    }

    heap->data[index] = data;
}

int HeapTop(struct Heap *heap) {
    return heap->data[0];
}

void HeapPop(struct Heap *heap) {
    if (heap->items == 0) {
        return;
    }

    int v = heap->data[heap->items - 1];
    --heap->items;

    size_t index = 0;
    while (index * 2 + 1 < heap->items) {
        size_t child1 = index * 2 + 1;
        size_t child2 = index * 2 + 2;

        if (child2 < heap->items && heap->data[child2] > heap->data[child1]) {
            child1 = child2;
        }

        if (heap->data[child1] <= v) {
            break;
        }

        heap->data[index] = heap->data[child1];
        index = child1;
    }

    heap->data[index] = v;
}

void HeapDestroy(struct Heap *heap) {
    assert(heap != NULL);

    free(heap->data);
    free(heap);
}

int main(void) {
    {
        struct Heap *h = HeapInit(16);
        for (int i = 0; i < 16; ++i) {
            HeapAdd(h, i * i);
        }

        while (h->items != 0) {
            int ret = HeapTop(h);
            HeapPop(h);
            printf("## pop %d\n", ret);
        }

        HeapDestroy(h);
    }
    return 0;
}
