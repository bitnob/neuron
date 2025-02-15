package migration

import (
	"fmt"
	"strings"
)

type SchemaBuilder struct {
	statements []string
}

func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		statements: make([]string, 0),
	}
}

func (s *SchemaBuilder) CreateTable(name string, fn func(*TableBuilder)) string {
	builder := NewTableBuilder(name)
	fn(builder)
	s.statements = append(s.statements, builder.Build())
	return builder.Build()
}

type TableBuilder struct {
	name        string
	columns     []string
	indexes     []string
	foreignKeys []string
}

func NewTableBuilder(name string) *TableBuilder {
	return &TableBuilder{
		name:        name,
		columns:     make([]string, 0),
		indexes:     make([]string, 0),
		foreignKeys: make([]string, 0),
	}
}

func (t *TableBuilder) Build() string {
	// TODO: Implement table creation SQL
	return fmt.Sprintf("CREATE TABLE %s (\n%s\n)", t.name, strings.Join(t.columns, ",\n"))
}

func (t *TableBuilder) Integer(name string, options ...ColumnOption) *TableBuilder {
	column := NewColumn(name, "INTEGER", options...)
	t.columns = append(t.columns, column.Build())
	return t
}

func (t *TableBuilder) String(name string, options ...ColumnOption) *TableBuilder {
	column := NewColumn(name, "VARCHAR", options...)
	t.columns = append(t.columns, column.Build())
	return t
}

// ColumnOption defines a column configuration option
type ColumnOption func(*Column)

// Column represents a database column
type Column struct {
	name     string
	dataType string
	nullable bool
	primary  bool
	unique   bool
	default_ string
}

func NewColumn(name, dataType string, options ...ColumnOption) *Column {
	c := &Column{
		name:     name,
		dataType: dataType,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func (c *Column) Build() string {
	parts := []string{c.name, c.dataType}
	if !c.nullable {
		parts = append(parts, "NOT NULL")
	}
	if c.primary {
		parts = append(parts, "PRIMARY KEY")
	}
	if c.unique {
		parts = append(parts, "UNIQUE")
	}
	if c.default_ != "" {
		parts = append(parts, "DEFAULT", c.default_)
	}
	return strings.Join(parts, " ")
}
