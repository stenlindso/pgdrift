package diff_test

import (
	"testing"

	"github.com/pgdrift/pgdrift/internal/diff"
	"github.com/pgdrift/pgdrift/internal/schema"
)

func buildSchema(tables map[string][]schema.Column) *schema.Schema {
	s := schema.NewSchema()
	for name, cols := range tables {
		t := &schema.Table{Name: name, Columns: cols}
		s.AddTable(t)
	}
	return s
}

func TestCompare_NoDrift(t *testing.T) {
	cols := []schema.Column{{Name: "id", DataType: "integer", Nullable: false}}
	src := buildSchema(map[string][]schema.Column{"users": cols})
	tgt := buildSchema(map[string][]schema.Column{"users": cols})

	result := diff.Compare(src, tgt)
	if result.HasDrift() {
		t.Errorf("expected no drift, got %d changes", len(result.Changes))
	}
}

func TestCompare_TableAdded(t *testing.T) {
	src := buildSchema(map[string][]schema.Column{})
	tgt := buildSchema(map[string][]schema.Column{
		"orders": {{Name: "id", DataType: "integer"}},
	})

	result := diff.Compare(src, tgt)
	if !result.HasDrift() {
		t.Fatal("expected drift")
	}
	if result.Changes[0].ChangeType != diff.ChangeAdded {
		t.Errorf("expected ChangeAdded, got %s", result.Changes[0].ChangeType)
	}
}

func TestCompare_TableRemoved(t *testing.T) {
	src := buildSchema(map[string][]schema.Column{
		"users": {{Name: "id", DataType: "integer"}},
	})
	tgt := buildSchema(map[string][]schema.Column{})

	result := diff.Compare(src, tgt)
	if len(result.Changes) != 1 || result.Changes[0].ChangeType != diff.ChangeRemoved {
		t.Errorf("expected one ChangeRemoved, got %+v", result.Changes)
	}
}

func TestCompare_ColumnTypeChanged(t *testing.T) {
	src := buildSchema(map[string][]schema.Column{
		"users": {{Name: "age", DataType: "integer"}},
	})
	tgt := buildSchema(map[string][]schema.Column{
		"users": {{Name: "age", DataType: "bigint"}},
	})

	result := diff.Compare(src, tgt)
	if !result.HasDrift() {
		t.Fatal("expected drift")
	}
	if result.Changes[0].ChangeType != diff.ChangeAltered {
		t.Errorf("expected ChangeAltered, got %s", result.Changes[0].ChangeType)
	}
}

func TestCompare_ColumnNullableChanged(t *testing.T) {
	src := buildSchema(map[string][]schema.Column{
		"users": {{Name: "email", DataType: "text", Nullable: false}},
	})
	tgt := buildSchema(map[string][]schema.Column{
		"users": {{Name: "email", DataType: "text", Nullable: true}},
	})

	result := diff.Compare(src, tgt)
	if !result.HasDrift() {
		t.Fatal("expected drift for nullable change")
	}
}
