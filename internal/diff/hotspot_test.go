package diff

import (
	"testing"
	"time"
)

func buildHotspotChangelog() *Changelog {
	cl := NewChangelog()

	// entry 1: two changes on orders, one on users
	r1 := &Result{
		Changes: []Change{
			{Kind: KindColumnTypeChanged, Schema: "public", Table: "orders"},
			{Kind: KindColumnAdded, Schema: "public", Table: "orders"},
			{Kind: KindColumnAdded, Schema: "public", Table: "users"},
		},
	}
	cl.Record(r1, time.Now())

	// entry 2: one more change on orders, two on products
	r2 := &Result{
		Changes: []Change{
			{Kind: KindColumnRemoved, Schema: "public", Table: "orders"},
			{Kind: KindTableAdded, Schema: "public", Table: "products"},
			{Kind: KindColumnAdded, Schema: "public", Table: "products"},
		},
	}
	cl.Record(r2, time.Now())

	return cl
}

func TestDetectHotspots_NilChangelog(t *testing.T) {
	rep := DetectHotspots(nil, 0)
	if rep == nil {
		t.Fatal("expected non-nil report")
	}
	if len(rep.Entries) != 0 {
		t.Errorf("expected no entries, got %d", len(rep.Entries))
	}
}

func TestDetectHotspots_RankedByCount(t *testing.T) {
	cl := buildHotspotChangelog()
	rep := DetectHotspots(cl, 0)

	if len(rep.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(rep.Entries))
	}
	if rep.Entries[0].Table != "orders" {
		t.Errorf("expected orders first, got %s", rep.Entries[0].Table)
	}
	if rep.Entries[0].Changes != 3 {
		t.Errorf("expected 3 changes for orders, got %d", rep.Entries[0].Changes)
	}
	if rep.Total != 6 {
		t.Errorf("expected total 6, got %d", rep.Total)
	}
}

func TestDetectHotspots_MinChangesFilter(t *testing.T) {
	cl := buildHotspotChangelog()
	rep := DetectHotspots(cl, 3)

	if len(rep.Entries) != 1 {
		t.Fatalf("expected 1 entry after filter, got %d", len(rep.Entries))
	}
	if rep.Entries[0].Table != "orders" {
		t.Errorf("expected orders, got %s", rep.Entries[0].Table)
	}
}

func TestDetectHotspots_EmptyChangelog(t *testing.T) {
	cl := NewChangelog()
	rep := DetectHotspots(cl, 0)
	if len(rep.Entries) != 0 {
		t.Errorf("expected no entries for empty changelog, got %d", len(rep.Entries))
	}
	if rep.Total != 0 {
		t.Errorf("expected total 0, got %d", rep.Total)
	}
}
