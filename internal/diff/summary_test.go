package diff

import (
	"testing"
)

func buildSummaryResult(changes ...Change) *Result {
	return &Result{Changes: changes}
}

func TestSummarize_NilResult(t *testing.T) {
	s := Summarize(nil)
	if s.TotalChanges != 0 {
		t.Errorf("expected 0 total changes, got %d", s.TotalChanges)
	}
	if len(s.ByKind) != 0 {
		t.Errorf("expected empty ByKind map")
	}
}

func TestSummarize_NoChanges(t *testing.T) {
	s := Summarize(&Result{})
	if s.TotalChanges != 0 {
		t.Errorf("expected 0 total changes, got %d", s.TotalChanges)
	}
	if len(s.AffectedTables) != 0 {
		t.Errorf("expected no affected tables")
	}
}

func TestSummarize_CountsByKind(t *testing.T) {
	r := buildSummaryResult(
		Change{Kind: KindTableAdded, Schema: "public", Table: "users"},
		Change{Kind: KindTableAdded, Schema: "public", Table: "orders"},
		Change{Kind: KindColumnRemoved, Schema: "public", Table: "users"},
	)
	s := Summarize(r)
	if s.TotalChanges != 3 {
		t.Errorf("expected 3 total changes, got %d", s.TotalChanges)
	}
	if s.ByKind[KindTableAdded] != 2 {
		t.Errorf("expected 2 KindTableAdded, got %d", s.ByKind[KindTableAdded])
	}
	if s.ByKind[KindColumnRemoved] != 1 {
		t.Errorf("expected 1 KindColumnRemoved, got %d", s.ByKind[KindColumnRemoved])
	}
}

func TestSummarize_AffectedTablesUnique(t *testing.T) {
	r := buildSummaryResult(
		Change{Kind: KindColumnTypeChanged, Schema: "public", Table: "users"},
		Change{Kind: KindColumnRemoved, Schema: "public", Table: "users"},
		Change{Kind: KindTableAdded, Schema: "public", Table: "orders"},
	)
	s := Summarize(r)
	if len(s.AffectedTables) != 2 {
		t.Errorf("expected 2 unique affected tables, got %d: %v", len(s.AffectedTables), s.AffectedTables)
	}
}

func TestSummary_HighestSeverity_High(t *testing.T) {
	r := buildSummaryResult(
		Change{Kind: KindColumnRemoved, Schema: "public", Table: "users"},
	)
	s := Summarize(r)
	if s.HighestSeverity() != SeverityHigh {
		t.Errorf("expected SeverityHigh, got %s", s.HighestSeverity())
	}
}

func TestSummary_HighestSeverity_NoChanges(t *testing.T) {
	s := Summarize(&Result{})
	if s.HighestSeverity() != SeverityLow {
		t.Errorf("expected SeverityLow for empty result, got %s", s.HighestSeverity())
	}
}
