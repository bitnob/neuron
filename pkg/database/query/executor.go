package query

import (
	"context"
	"database/sql"
	"time"
)

// Cache interface for query caching
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// Query represents a database query with caching options
type Query struct {
	SQL      string
	Args     []interface{}
	CacheKey string
	CacheTTL time.Duration
}

type Executor struct {
	db    *sql.DB
	cache Cache
}

func NewExecutor(db *sql.DB, cache Cache) *Executor {
	return &Executor{
		db:    db,
		cache: cache,
	}
}

func (e *Executor) Execute(ctx context.Context, query Query) (*sql.Rows, error) {
	if query.CacheKey != "" && e.cache != nil {
		if cached, err := e.cache.Get(ctx, query.CacheKey); err == nil {
			return cached.(*sql.Rows), nil
		}
	}

	rows, err := e.db.QueryContext(ctx, query.SQL, query.Args...)
	if err != nil {
		return nil, err
	}

	if query.CacheKey != "" && e.cache != nil {
		e.cache.Set(ctx, query.CacheKey, rows, query.CacheTTL)
	}

	return rows, nil
}
