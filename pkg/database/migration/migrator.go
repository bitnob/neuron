package migration

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Migration struct {
	ID        string
	Version   int64
	Name      string
	Up        func(*sql.Tx) error
	Down      func(*sql.Tx) error
	CreatedAt time.Time
}

type Migrator struct {
	db          *sql.DB
	migrations  []Migration
	tableName   string
	schemaName  string
	lockTimeout time.Duration
}

func NewMigrator(db *sql.DB, options ...MigratorOption) *Migrator {
	m := &Migrator{
		db:          db,
		tableName:   "migrations",
		schemaName:  "public",
		lockTimeout: 5 * time.Minute,
		migrations:  make([]Migration, 0),
	}

	for _, opt := range options {
		opt(m)
	}

	return m
}

func (m *Migrator) Initialize(ctx context.Context) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s (
			id SERIAL PRIMARY KEY,
			version BIGINT NOT NULL,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`, m.schemaName, m.tableName)

	_, err := m.db.ExecContext(ctx, query)
	return err
}

func (m *Migrator) Up(ctx context.Context) error {
	return m.runMigrations(ctx, true)
}

func (m *Migrator) Down(ctx context.Context) error {
	return m.runMigrations(ctx, false)
}
