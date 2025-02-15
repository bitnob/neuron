package query

type Builder struct {
	table      string
	selections []string
	conditions []Condition
	joins      []Join
	groupBy    []string
	having     []Condition
	orderBy    []Order
	limit      *int
	offset     *int
	params     []interface{}
}

func NewBuilder() *Builder {
	return &Builder{
		selections: make([]string, 0),
		conditions: make([]Condition, 0),
		joins:      make([]Join, 0),
		groupBy:    make([]string, 0),
		having:     make([]Condition, 0),
		orderBy:    make([]Order, 0),
		params:     make([]interface{}, 0),
	}
}

func (b *Builder) Select(columns ...string) *Builder {
	b.selections = append(b.selections, columns...)
	return b
}

func (b *Builder) From(table string) *Builder {
	b.table = table
	return b
}

func (b *Builder) Where(condition string, args ...interface{}) *Builder {
	b.conditions = append(b.conditions, Condition{
		SQL:  condition,
		Args: args,
	})
	return b
}
