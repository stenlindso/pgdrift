package report

import (
	"testing"

	"github.com/pgdrift/internal/diff"
)

func TestSeverity(t *testing.T) {
	cases := []struct {
		kind     diff.ChangeKind
		wantSev  string
	}{
		{diff.TableAdded, "HIGH"},
		{diff.TableRemoved, "HIGH"},
		{diff.ColumnAdded, "MEDIUM"},
		{diff.ColumnRemoved, "MEDIUM"},
		{diff.ColumnTypeChanged, "LOW"},
		{diff.ColumnNullableChanged, "LOW"},
		{diff.ColumnDefaultChanged, "LOW"},
	}
	for _, tc := range cases {
		got := Severity(tc.kind)
		if got != tc.wantSev {
			t.Errorf("Severity(%v) = %q, want %q", tc.kind, got, tc.wantSev)
		}
	}
}

func TestSummary_NoDrift(t *testing.T) {
	r := diff.Result{}
	got := Summary(r)
	if got != "No schema drift detected." {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSummary_WithChanges(t *testing.T) {
	r := diff.Result{
		Changes: []diff.Change{
			{Kind: diff.TableAdded, Schema: "public", Table: "orders"},
			{Kind: diff.ColumnRemoved, Schema: "public", Table: "users", Column: "email"},
			{Kind: diff.ColumnTypeChanged, Schema: "public", Table: "users", Column: "age"},
		},
	}
	got := Summary(r)
	if got == "No schema drift detected." {
		t.Fatal("expected drift summary, got no-drift message")
	}
	for _, sub := range []string{"3 change(s)", "1 high", "1 medium", "1 low"} {
		if !containsSub(got, sub) {
			t.Errorf("summary %q missing %q", got, sub)
		}
	}
}

func TestSummary_OnlyHigh(t *testing.T) {
	r := diff.Result{
		Changes: []diff.Change{
			{Kind: diff.TableAdded, Schema: "public", Table: "logs"},
			{Kind: diff.TableRemoved, Schema: "public", Table: "tmp"},
		},
	}
	got := Summary(r)
	if !containsSub(got, "2 high") {
		t.Errorf("expected '2 high' in %q", got)
	}
	if containsSub(got, "medium") || containsSub(got, "low") {
		t.Errorf("unexpected severity labels in %q", got)
	}
}

func containsSub(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubHelper(s, sub))
}

func containsSubHelper(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
