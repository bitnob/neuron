package cache

// MultiLevelCache implements a tiered caching system
type MultiLevelCache struct {
	l1    *ristretto.Cache // Memory cache
	l2    *redis.Client    // Distributed cache
	stats *CacheStats
}
