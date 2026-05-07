package diff

import (
	"testing"

	"github.com/pgdrift/pgdrift/internal/filter"
)

func buildFilterResult() Result {
	return Result{
		Changes: []Change{
			{Kind: TableAdded, Schema: "public", Table: "users"},
			{Kind: TableAdded, Schema: "public", Table: "orders"},
			{Kind: TableAdded, Schema: "audit", Table: "logs"},
			{Kind: ColumnAdded, Schema: "public", Table: "users", Column: "email"},
		},
	}
}

func TestFilterResult_NilFilter(t *testing.T) {
	r := buildFilterResult()
	out := FilterResult(r, nil)
	if len(out.Changes) != len(r.Changes) {
		t.Fatalf("expected %d changes, got %d", len(r.Changes), len(out.Changes))
	}
}

func TestFilterResult_ExcludeSchema(t *testing.T) {
	r := buildFilterResult()
	f := filter.New(filter.Options{ExcludeSchemas: []string{"audit"}})
	out := FilterResult(r, f)
	for _, ch := range out.Changes {
		if ch.Schema == "audit" {
			t.Errorf("expected audit schema to be excluded, got change: %+v", ch)
		}
	}
	if len(out.Changes) != 3 {
		t.Fatalf("expected 3 changes after excluding audit, got %d", len(out.Changes))
	}
}

func TestFilterResult_ExcludeTable(t *testing.T) {
	r := buildFilterResult()
	f := filter.New(filter.Options{ExcludeTables: []string{"orders"}})
	out := FilterResult(r, f)
	for _, ch := range out.Changes {
		if ch.Table == "orders" {
			t.Errorf("expected orders table to be excluded, got change: %+v", ch)
		}
	}
}

func TestFilterResult_IncludeSchema(t *testing.T) {
	r := buildFilterResult()
	f := filter.New(filter.Options{IncludeSchemas: []string{"public"}})
	out := FilterResult(r, f)
	for _, ch := range out.Changes {
		if ch.Schema != "public" {
			t.Errorf("expected only public schema, got: %+v", ch)
		}
	}
	if len(out.Changes) != 3 {
		t.Fatalf("expected 3 changes for public schema, got %d", len(out.Changes))
	}
}

func TestFilterResult_EmptyResult(t *testing.T) {
	r := Result{}
	f := filter.New(filter.Options{ExcludeSchemas: []string{"public"}})
	out := FilterResult(r, f)
	if out.HasDrift() {
		t.Error("expected no drift in empty result")
	}
}
