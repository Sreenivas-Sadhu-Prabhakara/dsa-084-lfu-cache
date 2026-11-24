package main

import "fmt"

type fnode struct {
	key, val, freq int
	prev, next     *fnode
}

// flist is a doubly linked list with sentinels holding all keys of one
// frequency, most recently touched at the front.
type flist struct {
	head, tail *fnode
	size       int
}

func newFlist() *flist {
	h, t := &fnode{}, &fnode{}
	h.next, t.prev = t, h
	return &flist{head: h, tail: t}
}

func (l *flist) pushFront(n *fnode) {
	n.next = l.head.next
	n.prev = l.head
	l.head.next.prev = n
	l.head.next = n
	l.size++
}

func (l *flist) remove(n *fnode) {
	n.prev.next = n.next
	n.next.prev = n.prev
	l.size--
}

// LFUCache evicts the least frequently used key, breaking frequency ties by
// least recent use. A map of per-frequency lists plus a running minimum
// frequency keeps every operation O(1).
type LFUCache struct {
	capacity, minFreq int
	items             map[int]*fnode
	freqs             map[int]*flist
}

// NewLFUCache returns a cache holding at most capacity entries.
func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{capacity: capacity, items: map[int]*fnode{}, freqs: map[int]*flist{}}
}

func (c *LFUCache) bucket(f int) *flist {
	if c.freqs[f] == nil {
		c.freqs[f] = newFlist()
	}
	return c.freqs[f]
}

// touch moves n from its frequency bucket to the next one up.
func (c *LFUCache) touch(n *fnode) {
	l := c.freqs[n.freq]
	l.remove(n)
	if l.size == 0 && c.minFreq == n.freq {
		c.minFreq++
	}
	n.freq++
	c.bucket(n.freq).pushFront(n)
}

// Get returns the value for key, counting the access.
func (c *LFUCache) Get(key int) (int, bool) {
	n, ok := c.items[key]
	if !ok {
		return 0, false
	}
	c.touch(n)
	return n.val, true
}

// Put inserts or updates key; on overflow it evicts the least frequently
// used key, oldest first within that frequency.
func (c *LFUCache) Put(key, val int) {
	if c.capacity <= 0 {
		return
	}
	if n, ok := c.items[key]; ok {
		n.val = val
		c.touch(n)
		return
	}
	if len(c.items) == c.capacity {
		l := c.freqs[c.minFreq]
		victim := l.tail.prev // least recently used at the minimum frequency
		l.remove(victim)
		delete(c.items, victim.key)
	}
	n := &fnode{key: key, val: val, freq: 1}
	c.items[key] = n
	c.bucket(1).pushFront(n)
	c.minFreq = 1
}

func main() {
	c := NewLFUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Get(1)              // key 1 now has frequency 2
	c.Put(3, 3)           // evicts key 2 (lowest frequency)
	fmt.Println(c.Get(2)) // miss
	fmt.Println(c.Get(3)) // 3
}
