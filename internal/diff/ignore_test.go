package diff

import (
	"testing"
)

func buildIgnoreResult() *Result {
	return &Result{
		Changes: []Change{
			{Schema: "public", Table: "users", Column: "email", Kind: KindColumnTypeChanged},
			{Schema: "public", Table: "users", Column: "name", Kind: KindColumnAdded},
			{Schema: "audit", Table: "logs", Column: "", Kind: KindTableAdded},
		},
	}
}

func TestIgnoreList_NilFilter(t *testing.T) {
	r := buildIgnoreResult()
	var il *IgnoreList
	out := il.Apply(r)
	if out != r {
		t.Fatal("expected original result when IgnoreList is nil")
	}
}

func TestIgnoreList_NoRules(t *testing.T) {
	r := buildIgnoreResult()
	il := NewIgnoreList(nil)
	out := il.Apply(r)
	if len(out.Changes) != len(r.Changes) {
		t.Fatalf("expected %d changes, got %d", len(r.Changes), len(out.Changes))
	}
}

func TestIgnoreList_ExactMatch(t *testing.T) {
	r := buildIgnoreResult()
	il := NewIgnoreList([]IgnoreRule{
		{Schema: "public", Table: "users", Column: "email", Kind: KindColumnTypeChanged},
	})
	out := il.Apply(r)
	if len(out.Changes) != 2 {
		t.Fatalf("expected 2 changes after suppression, got %d", len(out.Changes))
	}
}

func TestIgnoreList_WildcardTable(t *testing.T) {
	r := buildIgnoreResult()
	// suppress all changes in schema "public" regardless of table/column/kind
	il := NewIgnoreList([]IgnoreRule{
		{Schema: "public"},
	})
	out := il.Apply(r)
	if len(out.Changes) != 1 {
		t.Fatalf("expected 1 change after suppressing public schema, got %d", len(out.Changes))
	}
	if out.Changes[0].Schema != "audit" {
		t.Errorf("expected remaining change to be in audit schema, got %s", out.Changes[0].Schema)
	}
}

func TestIgnoreList_WildcardKind(t *testing.T) {
	r := buildIgnoreResult()
	il := NewIgnoreList([]IgnoreRule{
		{Schema: "public", Table: "users", Column: "name"},
	})
	out := il.Apply(r)
	if len(out.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(out.Changes))
	}
}

func TestIgnoreList_MatchesNilResult(t *testing.T) {
	il := NewIgnoreList([]IgnoreRule{{Schema: "public"}})
	out := il.Apply(nil)
	if out != nil {
		t.Fatal("expected nil result when input is nil")
	}
}
