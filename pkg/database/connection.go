package database

import (
	"database/sql"
	"sync"
	"time"
)

type ConnectionPool struct {
	mu          sync.RWMutex
	connections chan *sql.Conn
	maxOpen     int
	maxIdle     int
	maxLifetime time.Duration
	maxIdleTime time.Duration
}

func NewConnectionPool(config ConnectionConfig) *ConnectionPool {
	return &ConnectionPool{
		connections: make(chan *sql.Conn, config.MaxOpen),
		maxOpen:     config.MaxOpen,
		maxIdle:     config.MaxIdle,
		maxLifetime: config.MaxLifetime,
		maxIdleTime: config.MaxIdleTime,
	}
}
