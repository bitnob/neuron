package integration

import (
	"context"
	"testing"
	"time"

	"your/framework/pkg/cache"
	"your/framework/pkg/config"
	"your/framework/pkg/database"
)

type IntegrationSuite struct {
	config  *config.Config
	cache   cache.Cache
	db      *database.ConnectionPool
	cleanup func()
}

func setupSuite(t *testing.T) *IntegrationSuite {
	// Load test configuration
	cfg, err := config.NewLoader(
		config.NewFileSource("../../config/test.yaml", 1),
	).LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Initialize cache
	cacheInstance := cache.NewMemoryCache(cache.Options{
		MaxEntries: 100,
		DefaultTTL: time.Minute,
	})

	// Initialize database
	pool := database.NewConnectionPool(database.ConnectionConfig{
		MaxOpen:     10,
		MaxIdle:     5,
		MaxLifetime: time.Hour,
	})

	return &IntegrationSuite{
		config: cfg,
		cache:  cacheInstance,
		db:     pool,
		cleanup: func() {
			// Cleanup resources
			pool.Close()
		},
	}
}

func TestIntegration(t *testing.T) {
	suite := setupSuite(t)
	defer suite.cleanup()

	// Run integration tests
	t.Run("CacheOperations", func(t *testing.T) {
		testCacheOperations(t, suite)
	})

	t.Run("DatabaseOperations", func(t *testing.T) {
		testDatabaseOperations(t, suite)
	})
}

func testCacheOperations(t *testing.T, suite *IntegrationSuite) {
	ctx := context.Background()
	key := "test-key"
	value := "test-value"

	err := suite.cache.Set(ctx, key, value, time.Minute)
	if err != nil {
		t.Errorf("Failed to set cache: %v", err)
	}

	got, err := suite.cache.Get(ctx, key)
	if err != nil {
		t.Errorf("Failed to get from cache: %v", err)
	}
	if got != value {
		t.Errorf("Cache get = %v, want %v", got, value)
	}
}
