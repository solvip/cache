package cache

type Cache interface {
	Get(key string) (value int, ok bool)
	Put(key string, value int)
	Statistics() (hits, misses, evictions int64)
}
