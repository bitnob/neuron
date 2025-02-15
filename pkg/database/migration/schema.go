package migration

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
