package diff

import (
	"testing"
	"time"
)

func buildMaturityChangelog() *Changelog {
	cl := NewChangelog()
	now := time.Now()

	// entry 1: public.orders changed
	r1 := &Result{
		Changes: []Change{
			{Schema: "public", Table: "orders", Kind: ColumnTypeChanged},
		},
	}
	cl.Entries = append(cl.Entries, ChangelogEntry{
		RecordedAt: now.Add(-2 * time.Hour),
		Changes:    r1.Changes,
	})

	// entry 2: public.orders changed again + public.users changed
	r2 := &Result{
		Changes: []Change{
			{Schema: "public", Table: "orders", Kind: ColumnAdded},
			{Schema: "public", Table: "users", Kind: TableAdded},
		},
	}
	cl.Entries = append(cl.Entries, ChangelogEntry{
		RecordedAt: now.Add(-1 * time.Hour),
		Changes:    r2.Changes,
	})

	return cl
}

func TestAssessMaturity_NilChangelog(t *testing.T) {
	report := AssessMaturity(nil, 3)
	if report == nil {
		t.Fatal("expected non-nil report")
	}
	if len(report.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(report.Entries))
	}
}

func TestAssessMaturity_BelowMinRuns(t *testing.T) {
	cl := buildMaturityChangelog()
	report := AssessMaturity(cl, 10) // minRuns > len(entries)
	for _, e := range report.Entries {
		if e.Level != MaturityNew {
			t.Errorf("expected MaturityNew for %s.%s, got %s", e.Schema, e.Table, e.Level)
		}
	}
}

func TestAssessMaturity_Fluctuating(t *testing.T) {
	cl := buildMaturityChangelog()
	report := AssessMaturity(cl, 1) // minRuns satisfied

	found := false
	for _, e := range report.Entries {
		if e.Schema == "public" && e.Table == "orders" {
			found = true
			if e.Level != MaturityFluctuating {
				t.Errorf("expected Fluctuating for orders, got %s", e.Level)
			}
			if e.Changes != 2 {
				t.Errorf("expected 2 changes, got %d", e.Changes)
			}
		}
	}
	if !found {
		t.Error("orders entry not found in report")
	}
}

func TestMaturityLevel_String(t *testing.T) {
	cases := []struct {
		level MaturityLevel
		want  string
	}{
		{MaturityNew, "new"},
		{MaturityStable, "stable"},
		{MaturityFluctuating, "fluctuating"},
		{MaturityUnknown, "unknown"},
	}
	for _, c := range cases {
		if got := c.level.String(); got != c.want {
			t.Errorf("String() = %q, want %q", got, c.want)
		}
	}
}
