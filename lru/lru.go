package lru

// LRU implements a least-recently-used cache
type LRU struct {
	capacity int

	m map[string]*node

	// The most recently used item is always at the head of the list
	// and the tail contains the item to evict
	cache list

	hits      int
	misses    int
	evictions int
}

// New - allocate a new LRU cache having capacity for `capacity` items.
func New(capacity int) *LRU {
	return &LRU{
		capacity: capacity,
		m:        make(map[string]*node),
	}
}

func (lru *LRU) Statistics() (hits, misses, evictions int) {
	return lru.hits, lru.misses, lru.evictions
}

func (lru *LRU) Get(key string) (int, bool) {
	node := lru.m[key]
	if node == nil {
		lru.misses++
		return 0, false
	}

	lru.cache.moveToHead(node)
	lru.hits++

	return node.value, true
}

func (lru *LRU) Put(key string, value int) {
	n := lru.m[key]
	if n != nil {
		// The node is already in the cache.
		// We simply need to update it's value and move it to the head
		n.value = value
		lru.cache.moveToHead(n)

		return
	}

	// If we're at capacity, we need to evict the LRU item
	if len(lru.m) == lru.capacity {
		evicted := lru.cache.dropTail()
		delete(lru.m, evicted.key)
		lru.evictions++
	}

	n = &node{key: key, value: value}
	lru.m[key] = n
	lru.cache.moveToHead(n)

	return
}
