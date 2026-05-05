package schema

import (
	"testing"
)

func TestNewSchema(t *testing.T) {
	s := NewSchema()
	if s == nil {
		t.Fatal("expected non-nil schema")
	}
	if len(s.Tables) != 0 {
		t.Errorf("expected empty tables, got %d", len(s.Tables))
	}
}

func TestAddTable(t *testing.T) {
	s := NewSchema()
	tbl := Table{
		Name:    "users",
		Columns: map[string]Column{
			"id": {Name: "id", DataType: "integer", IsNullable: false},
		},
	}
	s.AddTable(tbl)

	if _, ok := s.Tables["users"]; !ok {
		t.Error("expected table 'users' to be present")
	}
}

func TestTableNames(t *testing.T) {
	s := NewSchema()
	s.AddTable(Table{Name: "orders", Columns: map[string]Column{}})
	s.AddTable(Table{Name: "users", Columns: map[string]Column{}})

	names := s.TableNames()
	if len(names) != 2 {
		t.Errorf("expected 2 table names, got %d", len(names))
	}

	found := map[string]bool{}
	for _, n := range names {
		found[n] = true
	}
	if !found["users"] || !found["orders"] {
		t.Error("expected both 'users' and 'orders' in table names")
	}
}

func TestColumnNullable(t *testing.T) {
	col := Column{Name: "email", DataType: "text", IsNullable: true}
	if !col.IsNullable {
		t.Error("expected column to be nullable")
	}
}

func TestColumnDefault(t *testing.T) {
	defVal := "now()"
	col := Column{Name: "created_at", DataType: "timestamp", Default: &defVal}
	if col.Default == nil || *col.Default != "now()" {
		t.Error("expected default value 'now()'")
	}
}
