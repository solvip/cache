package cache

type Cache interface {
	Get(key string) (value interface{}, ok bool)
	Put(key string, value interface{})
	Statistics() (hits, misses, evictions int64)
}
