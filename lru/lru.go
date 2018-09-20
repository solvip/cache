package lru

// LRU implements a least-recently-used cache
type LRU struct {
	// freeNodes initially represents our backing array of cache nodes,
	// which we allocate at startup
	freeNodes []node

	// m maps a cache key to the node responsible
	// for the value associated with key
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
		freeNodes: make([]node, capacity),
		m:         make(map[string]*node),
	}
}

func (lru *LRU) Statistics() (hits, misses, evictions int) {
	return lru.hits, lru.misses, lru.evictions
}

func (lru *LRU) Get(key string) (interface{}, bool) {
	node := lru.m[key]
	if node == nil {
		lru.misses++
		return nil, false
	}

	lru.cache.moveToHead(node)
	lru.hits++

	return node.value, true
}

func (lru *LRU) Put(key string, value interface{}) {
	var n *node

	if n = lru.m[key]; n != nil {
		// The node is already in the cache.
		// We simply need to update it's value and move it to the head
		n.value = value
		lru.cache.moveToHead(n)

		return
	}

	// If we're at capacity, we need to evict the LRU item
	// We reuse the evicted node as the new head node if possible
	// If we're not at capacity; we pick a free node instead
	if len(lru.freeNodes) == 0 {
		n = lru.cache.dropTail()
		delete(lru.m, n.key)
		lru.evictions++
	} else {
		n = &lru.freeNodes[0]
		lru.freeNodes = lru.freeNodes[1:]
	}

	n.key = key
	n.value = value

	lru.m[key] = n

	lru.cache.moveToHead(n)

	return
}
