package main

import "testing"

func TestLFUCache(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	if v, _ := c.Get(1); v != 1 {
		t.Errorf("Get(1) = %d", v)
	}
	c.Put(3, 3) // evicts 2: freq(1)=2, freq(2)=1
	if _, ok := c.Get(2); ok {
		t.Error("2 should be evicted (lowest frequency)")
	}
	if v, _ := c.Get(3); v != 3 {
		t.Errorf("Get(3) = %d", v)
	}
	c.Put(4, 4) // freq: 1->2, 3->2; tie broken by recency: 1 is older
	if _, ok := c.Get(1); ok {
		t.Error("1 should be evicted (tie on freq, least recent)")
	}
	if v, _ := c.Get(3); v != 3 {
		t.Errorf("Get(3) = %d", v)
	}
	if v, _ := c.Get(4); v != 4 {
		t.Errorf("Get(4) = %d", v)
	}
}

func TestLFUUpdateCountsAsUse(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(1, 10) // update bumps 1's frequency
	c.Put(3, 3)  // must evict 2
	if _, ok := c.Get(2); ok {
		t.Error("2 should be evicted")
	}
	if v, _ := c.Get(1); v != 10 {
		t.Errorf("Get(1) = %d, want 10", v)
	}
}

func TestLFUZeroCapacity(t *testing.T) {
	c := NewLFUCache(0)
	c.Put(1, 1)
	if _, ok := c.Get(1); ok {
		t.Error("zero-capacity cache stores nothing")
	}
}
