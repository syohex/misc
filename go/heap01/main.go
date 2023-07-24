package main

import "fmt"

type Heap struct {
	heap []int
}

func NewHeap() *Heap {
	return &Heap{
		heap: []int{},
	}
}

func (h *Heap) Add(n int) {
	h.heap = append(h.heap, n)
	index := len(h.heap) - 1

	for index > 0 {
		parent := (index - 1) / 2
		if h.heap[parent] >= n {
			break
		}

		h.heap[index] = h.heap[parent]
		index = parent
	}

	h.heap[index] = n
}

func (h *Heap) Top() int {
	return h.heap[0]
}

func (h *Heap) Pop() {
	n := h.heap[len(h.heap)-1]
	h.heap = h.heap[:len(h.heap)-1]
	if len(h.heap) == 0 {
		return
	}

	index := 0

	for 2*index+1 < len(h.heap) {
		child1 := 2*index + 1
		child2 := 2*index + 2

		if child2 < len(h.heap) && h.heap[child1] < h.heap[child2] {
			child1 = child2
		}

		if n >= h.heap[child1] {
			break
		}

		h.heap[index] = h.heap[child1]
		index = child1
	}

	h.heap[index] = n
}

func main() {
	{
		h := NewHeap()
		for i := 1; i <= 10; i++ {
			h.Add(i * i)
		}

		fmt.Printf("### %v\n", h.heap)

		for i := 1; i <= 10; i++ {
			fmt.Printf("## pop %d\n", h.Top())
			h.Pop()
		}
	}

}
