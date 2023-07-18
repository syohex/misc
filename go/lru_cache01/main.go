package main

import "fmt"

type ListNode struct {
	prev  *ListNode
	next  *ListNode
	key   int
	value int
}

func NewListNode(key int, value int) *ListNode {
	return &ListNode{
		prev:  nil,
		next:  nil,
		key:   key,
		value: value,
	}
}

type LRUCache struct {
	head     *ListNode
	tail     *ListNode
	capacity int
	cache    map[int]*ListNode
}

func NewLRUCache(capacity int) *LRUCache {
	c := &LRUCache{
		head:     NewListNode(-1, -1),
		tail:     NewListNode(-1, -1),
		capacity: capacity,
		cache:    make(map[int]*ListNode),
	}

	c.head.next = c.tail
	c.tail.prev = c.head

	return c
}

func (c *LRUCache) Add(key int, value int) {
	if val, ok := c.cache[key]; ok {
		c.remove(val)
	}

	node := NewListNode(key, value)
	c.cache[key] = node
	c.pushBack(node)

	if len(c.cache) > c.capacity {
		removed := c.head.next
		c.head.next = removed.next
		delete(c.cache, removed.key)
	}
}

func (c *LRUCache) Get(key int) int {
	v, ok := c.cache[key]
	if !ok {
		return -1
	}

	c.remove(v)
	c.pushBack(v)

	return v.value
}

func (c *LRUCache) pushBack(node *ListNode) {
	prev := c.tail.prev
	prev.next = node
	node.prev = prev
	node.next = c.tail
	c.tail.prev = node
}

func (c *LRUCache) remove(node *ListNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *LRUCache) Print() {
	p := c.head.next
	fmt.Print("[")
	for p != nil {
		fmt.Printf(" (%d => %d)", p.key, p.value)
		p = p.next
	}
	fmt.Print("]\n")
}

func main() {
	c := NewLRUCache(2)
	c.Add(1, 1)
	c.Print()
	c.Add(2, 2)
	c.Print()
	fmt.Printf("Get(1)=%d(=1)\n", c.Get(1))
	c.Add(3, 3)
	c.Add(4, 4)
	fmt.Printf("Get(2)=%d(=-1)\n", c.Get(2))
	c.Add(4, 5)
	fmt.Printf("Get(3)=%d(=3)\n", c.Get(3))
	fmt.Printf("Get(4)=%d(=5)\n", c.Get(4))
}
