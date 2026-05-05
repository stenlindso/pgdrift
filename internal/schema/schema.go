package schema

// Table represents a PostgreSQL table with its columns.
type Table struct {
	Name    string
	Columns map[string]Column
}

// Column represents a column in a PostgreSQL table.
type Column struct {
	Name       string
	DataType   string
	IsNullable bool
	Default    *string
}

// Schema represents the full schema snapshot of a database.
type Schema struct {
	Tables map[string]Table
}

// NewSchema initializes an empty Schema.
func NewSchema() *Schema {
	return &Schema{
		Tables: make(map[string]Table),
	}
}

// AddTable adds a Table to the Schema.
func (s *Schema) AddTable(t Table) {
	s.Tables[t.Name] = t
}

// TableNames returns a sorted list of table names in the schema.
func (s *Schema) TableNames() []string {
	names := make([]string, 0, len(s.Tables))
	for name := range s.Tables {
		names = append(names, name)
	}
	return names
}
