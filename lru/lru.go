package lru

import (
	"container/list"
)

type entry struct {
	key   string
	value interface{}
}

// LRU implements a least-recently-used cache
type LRU struct {
	capacity int

	// m maps a cache key to the node responsible
	// for the value associated with key
	m map[string]*list.Element

	// The most recently used item is always at the head of the list
	// and the tail contains the item to evict
	// cache list
	l list.List

	hits      int
	misses    int
	evictions int
}

// New - allocate a new LRU cache having capacity for `capacity` items.
func New(capacity int) *LRU {
	lru := &LRU{
		capacity: capacity,
		m:        make(map[string]*list.Element),
	}

	// // Preallocate the cache list
	// entries := make([]entry, capacity)
	// for i := 0; i < capacity; i++ {
	// 	lru.l.PushFront(&entries[i])
	// }

	return lru
}

func (lru *LRU) Statistics() (hits, misses, evictions int) {
	return lru.hits, lru.misses, lru.evictions
}

func (lru *LRU) Get(key string) (interface{}, bool) {
	elem := lru.m[key]
	if elem == nil {
		lru.misses++
		return nil, false
	}

	lru.l.MoveToFront(elem)
	lru.hits++

	return elem.Value.(*entry).value, true
}

func (lru *LRU) Put(key string, value interface{}) {
	elem := lru.m[key]
	if elem != nil {
		// The entry is already in the cache.
		lru.l.MoveToFront(elem)
		elem.Value.(*entry).value = value

		return
	}

	// We're still not at capacity; allocate a new node at the head
	// and continue
	if len(lru.m) < lru.capacity {
		elem = lru.l.PushFront(&entry{key: key, value: value})
		lru.m[key] = elem
		return
	}

	lru.evictions++
	elem = lru.l.Back()
	lru.l.MoveToFront(elem)

	delete(lru.m, elem.Value.(*entry).key) // delete old key
	lru.m[key] = elem

	elem.Value.(*entry).key = key
	elem.Value.(*entry).value = value

	return
}
