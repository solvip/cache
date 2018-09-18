package lru

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestGetWhenEmpty(t *testing.T) {
	c := New(10)

	if value, ok := c.Get("hello"); ok {
		t.Fatalf("expected c.Get to return _, false; instead got %v, %v", value, ok)
	}
}

func TestPutWhenEmpty(t *testing.T) {
	c := New(10)

	key := "hello"
	c.Put(key, 1)
	assert_cache_entry_exists(t, c, key, 1)
}

func TestPutOverwritesKey(t *testing.T) {
	c := New(3)

	c.Put("key-1", 1)
	c.Put("key-1", 2)
	c.Put("key-1", 3)
	c.Put("key-1", 4)
	assert_cache_entry_exists(t, c, "key-1", 4)

	c.Put("key-2", 2)
	c.Put("key-3", 3)
	c.Put("key-4", 4)
	assert_cache_entry_absent(t, c, "key-1")
	assert_cache_entry_exists(t, c, "key-2", 2)
	assert_cache_entry_exists(t, c, "key-3", 3)
	assert_cache_entry_exists(t, c, "key-4", 4)

}

func TestEvictions(t *testing.T) {
	c := New(2)

	c.Put("1", 1)
	c.Put("2", 2)
	c.Put("3", 3)

	if value, ok := c.Get("1"); ok {
		t.Fatalf("Expected 1 to be evicted; instead Get() returned %v, %v", value, ok)
	}

	assert_cache_entry_exists(t, c, "2", 2)
	assert_cache_entry_exists(t, c, "3", 3)

	c.Put("4", 4)
	assert_cache_entry_absent(t, c, "2")
	assert_cache_entry_exists(t, c, "3", 3)
	assert_cache_entry_exists(t, c, "4", 4)
}

func TestStatistics(t *testing.T) {
	c := New(10)

	for i := 0; i < 100; i++ {
		c.Put(strconv.Itoa(i), i)
	}

	for i := 0; i < 55; i++ {
		c.Get("no_such_key")
	}

	for i := 90; i < 100; i++ {
		c.Get(strconv.Itoa(i))
	}

	expectedEvicts := 90
	expectedHits := 10
	expectedMisses := 55

	hits, misses, evicts := c.Statistics()
	if hits != expectedHits || misses != expectedMisses || evicts != expectedEvicts {
		t.Fatalf("Expected c.Statistics = %v, %v, %v; instead got %v, %v, %v",
			expectedHits, expectedMisses, expectedEvicts, hits, misses, evicts)
	}
}

func assert_cache_entry_exists(t *testing.T, c *LRU, key string, expectedValue int) {
	if value, ok := c.Get(key); value != expectedValue || !ok {
		t.Fatalf("Expected c.Get(%s) = %v, true; instead got %v, %v", key, expectedValue, value, ok)
	}
}

func assert_cache_entry_absent(t *testing.T, c *LRU, key string) {
	if value, ok := c.Get(key); ok {
		t.Fatalf("Expected c.Get(%s) = _, false; instead got %v, %v", key, value, ok)
	}
}
