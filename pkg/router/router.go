package router

// FastRouter implements high-performance routing
type FastRouter struct {
	trees map[string]*node
	pool  sync.Pool
	cache *PathCache
}

// PathCache implements route caching
type PathCache struct {
	cache *lru.Cache
	mu    sync.RWMutex
}
