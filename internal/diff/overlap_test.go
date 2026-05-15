package diff

import (
	"testing"
)

func buildOverlapResult(changes []*Change) *Result {
	return &Result{Changes: changes}
}

func TestDetectOverlap_NilInputs(t *testing.T) {
	r := DetectOverlap(nil, nil)
	if r == nil {
		t.Fatal("expected non-nil report")
	}
	if r.TotalShared != 0 || r.TotalConflicts != 0 {
		t.Errorf("expected empty report, got %+v", r)
	}
}

func TestDetectOverlap_NoSharedTables(t *testing.T) {
	src := buildOverlapResult([]*Change{
		{Table: "orders", Kind: KindTableAdded},
	})
	tgt := buildOverlapResult([]*Change{
		{Table: "products", Kind: KindTableAdded},
	})
	r := DetectOverlap(src, tgt)
	if r.TotalShared != 0 {
		t.Errorf("expected 0 shared tables, got %d", r.TotalShared)
	}
}

func TestDetectOverlap_SharedTableNoConflict(t *testing.T) {
	src := buildOverlapResult([]*Change{
		{Table: "users", Kind: KindTableAdded},
	})
	tgt := buildOverlapResult([]*Change{
		{Table: "users", Kind: KindTableAdded},
	})
	r := DetectOverlap(src, tgt)
	if r.TotalShared != 1 {
		t.Errorf("expected 1 shared table, got %d", r.TotalShared)
	}
	if r.TotalConflicts != 0 {
		t.Errorf("expected 0 conflicts, got %d", r.TotalConflicts)
	}
}

func TestDetectOverlap_ConflictingColumn(t *testing.T) {
	src := buildOverlapResult([]*Change{
		{Table: "users", Column: "email", Kind: KindColumnTypeChanged},
	})
	tgt := buildOverlapResult([]*Change{
		{Table: "users", Column: "email", Kind: KindColumnTypeChanged},
	})
	r := DetectOverlap(src, tgt)
	if r.TotalShared != 1 {
		t.Errorf("expected 1 shared table, got %d", r.TotalShared)
	}
	if r.TotalConflicts != 1 {
		t.Errorf("expected 1 conflict, got %d", r.TotalConflicts)
	}
	cols, ok := r.ConflictingColumns["users"]
	if !ok || len(cols) != 1 || cols[0] != "email" {
		t.Errorf("unexpected conflicting columns: %v", r.ConflictingColumns)
	}
}

func TestDetectOverlap_UniqueSharedTables(t *testing.T) {
	src := buildOverlapResult([]*Change{
		{Table: "logs", Kind: KindTableAdded},
	})
	tgt := buildOverlapResult([]*Change{
		{Table: "logs", Column: "id", Kind: KindColumnTypeChanged},
		{Table: "logs", Column: "ts", Kind: KindColumnTypeChanged},
	})
	r := DetectOverlap(src, tgt)
	if r.TotalShared != 1 {
		t.Errorf("expected 1 unique shared table, got %d", r.TotalShared)
	}
	if r.TotalConflicts != 2 {
		t.Errorf("expected 2 conflicts, got %d", r.TotalConflicts)
	}
}
