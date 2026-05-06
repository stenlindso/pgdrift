package diff_test

import (
	"testing"

	"github.com/pgdrift/pgdrift/internal/diff"
)

func TestResult_HasDrift_Empty(t *testing.T) {
	r := &diff.Result{}
	if r.HasDrift() {
		t.Error("expected no drift on empty result")
	}
}

func TestResult_HasDrift_WithChanges(t *testing.T) {
	r := &diff.Result{}
	r.Add(diff.Change{
		Object:     "table:foo",
		ChangeType: diff.ChangeAdded,
		Detail:     "table added",
	})
	if !r.HasDrift() {
		t.Error("expected drift after adding a change")
	}
}

func TestChange_String(t *testing.T) {
	c := diff.Change{
		Object:     "table:users",
		ChangeType: diff.ChangeRemoved,
		Detail:     "table removed",
	}
	s := c.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
	for _, want := range []string{"removed", "table:users", "table removed"} {
		if !containsStr(s, want) {
			t.Errorf("expected %q in %q", want, s)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}())
}
