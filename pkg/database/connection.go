package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"sync"
	"time"
)

// ConnectionConfig holds database connection configuration
type ConnectionConfig struct {
	MaxOpen     int
	MaxIdle     int
	MaxLifetime time.Duration
	MaxIdleTime time.Duration
}

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

func (p *ConnectionPool) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	conn, err := p.Get(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		p.Put(conn)
		return nil, err
	}
	return tx, nil
}

func (p *ConnectionPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.connections)
	for conn := range p.connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *ConnectionPool) Get(ctx context.Context) (*sql.Conn, error) {
	select {
	case conn := <-p.connections:
		return conn, nil
	default:
		db := sql.OpenDB(p)
		conn, err := db.Conn(ctx)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

// Implement sql.Connector interface
func (p *ConnectionPool) Connect(ctx context.Context) (driver.Conn, error) {
	return nil, errors.New("not implemented")
}

func (p *ConnectionPool) Driver() driver.Driver {
	return nil
}

func (p *ConnectionPool) Put(conn *sql.Conn) {
	select {
	case p.connections <- conn:
		// Connection returned to pool
	default:
		// Pool is full, close the connection
		conn.Close()
	}
}
