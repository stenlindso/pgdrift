package diff

import (
	"testing"
)

func buildChangelogResult(kinds ...ChangeKind) *Result {
	r := &Result{}
	for i, k := range kinds {
		table := fmt.Sprintf("public.table%d", i)
		r.Changes = append(r.Changes, Change{
			Schema: "public",
			Table:  fmt.Sprintf("table%d", i),
			Kind:   k,
			Object: table,
		})
	}
	return r
}

func TestNewChangelog_Empty(t *testing.T) {
	c := NewChangelog()
	if c == nil {
		t.Fatal("expected non-nil changelog")
	}
	if c.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", c.Len())
	}
}

func TestChangelog_RecordNil(t *testing.T) {
	c := NewChangelog()
	c.Record(nil, "")
	if c.Len() != 0 {
		t.Errorf("expected 0 entries after nil record, got %d", c.Len())
	}
}

func TestChangelog_RecordAddsEntry(t *testing.T) {
	c := NewChangelog()
	r := buildChangelogResult(KindColumnAdded, KindTableRemoved)
	c.Record(r, "v1.2")
	if c.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", c.Len())
	}
	e := c.Latest()
	if e == nil {
		t.Fatal("expected non-nil latest entry")
	}
	if e.Label != "v1.2" {
		t.Errorf("expected label v1.2, got %q", e.Label)
	}
	if e.Summary.TotalChanges != 2 {
		t.Errorf("expected 2 total changes in summary, got %d", e.Summary.TotalChanges)
	}
}

func TestChangelog_Latest_Empty(t *testing.T) {
	c := NewChangelog()
	if c.Latest() != nil {
		t.Error("expected nil latest on empty changelog")
	}
}

func TestChangelog_TopChanged(t *testing.T) {
	c := NewChangelog()

	r1 := &Result{
		Changes: []Change{
			{Schema: "public", Table: "users", Kind: KindColumnAdded, Object: "public.users"},
			{Schema: "public", Table: "orders", Kind: KindColumnAdded, Object: "public.orders"},
		},
	}
	r2 := &Result{
		Changes: []Change{
			{Schema: "public", Table: "users", Kind: KindColumnTypeChanged, Object: "public.users"},
		},
	}

	c.Record(r1, "")
	c.Record(r2, "")

	top := c.TopChanged(2)
	if len(top) == 0 {
		t.Fatal("expected at least one result")
	}
	if top[0] != "public.users (2)" {
		t.Errorf("expected public.users (2) as top, got %q", top[0])
	}
}

func TestChangelog_TopChanged_NilChangelog(t *testing.T) {
	var c *Changelog
	if c.TopChanged(5) != nil {
		t.Error("expected nil from nil changelog")
	}
}
