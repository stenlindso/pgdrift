package diff

import (
	"testing"
	"time"
)

func buildLineageResult(n int) *Result {
	r := &Result{}
	for i := 0; i < n; i++ {
		r.Changes = append(r.Changes, Change{
			Kind:   KindColumnTypeChanged,
			Schema: "public",
			Table:  "users",
			Object: "email",
		})
	}
	return r
}

func TestNewLineage_Empty(t *testing.T) {
	l := NewLineage()
	if len(l.Entries()) != 0 {
		t.Fatal("expected empty entries")
	}
	if !l.Stable() {
		t.Fatal("empty lineage should be stable")
	}
}

func TestLineage_RecordNil(t *testing.T) {
	l := NewLineage()
	l.Record(nil, "abc", time.Now())
	if len(l.Entries()) != 0 {
		t.Fatal("nil result should not be recorded")
	}
}

func TestLineage_RecordAddsEntry(t *testing.T) {
	l := NewLineage()
	r := buildLineageResult(2)
	l.Record(r, "fp1", time.Now())
	entries := l.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Fingerprint != "fp1" {
		t.Errorf("unexpected fingerprint: %s", entries[0].Fingerprint)
	}
	if entries[0].ChangeCount != 2 {
		t.Errorf("expected change count 2, got %d", entries[0].ChangeCount)
	}
}

func TestLineage_Stable_AllSame(t *testing.T) {
	l := NewLineage()
	r := buildLineageResult(0)
	l.Record(r, "fp1", time.Now())
	l.Record(r, "fp1", time.Now())
	if !l.Stable() {
		t.Fatal("expected stable lineage")
	}
}

func TestLineage_Stable_Different(t *testing.T) {
	l := NewLineage()
	r := buildLineageResult(1)
	l.Record(r, "fp1", time.Now())
	l.Record(r, "fp2", time.Now())
	if l.Stable() {
		t.Fatal("expected unstable lineage")
	}
}

func TestLineage_DivergencePoint(t *testing.T) {
	l := NewLineage()
	r := buildLineageResult(1)
	t0 := time.Now()
	t1 := t0.Add(time.Hour)
	l.Record(r, "fp1", t0)
	l.Record(r, "fp1", t1)
	l.Record(r, "fp2", t1.Add(time.Hour))

	e, err := l.DivergencePoint()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Fingerprint != "fp2" {
		t.Errorf("expected fp2, got %s", e.Fingerprint)
	}
}

func TestLineage_DivergencePoint_NoDivergence(t *testing.T) {
	l := NewLineage()
	r := buildLineageResult(0)
	l.Record(r, "fp1", time.Now())
	l.Record(r, "fp1", time.Now())
	_, err := l.DivergencePoint()
	if err == nil {
		t.Fatal("expected error when no divergence")
	}
}

func TestLineage_DivergencePoint_TooFewEntries(t *testing.T) {
	l := NewLineage()
	_, err := l.DivergencePoint()
	if err == nil {
		t.Fatal("expected error for empty lineage")
	}
}
