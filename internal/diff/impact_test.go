package diff

import (
	"testing"
)

func buildImpactResult(kinds ...ChangeKind) *Result {
	r := &Result{}
	for _, k := range kinds {
		r.Changes = append(r.Changes, Change{
			Kind:   k,
			Schema: "public",
			Table:  "users",
		})
	}
	return r
}

func TestAssessImpact_NilResult(t *testing.T) {
	rep := AssessImpact(nil)
	if rep == nil {
		t.Fatal("expected non-nil report")
	}
	if rep.Overall != ImpactNone {
		t.Errorf("expected ImpactNone, got %v", rep.Overall)
	}
}

func TestAssessImpact_NoChanges(t *testing.T) {
	rep := AssessImpact(&Result{})
	if rep.Overall != ImpactNone {
		t.Errorf("expected ImpactNone, got %v", rep.Overall)
	}
	if len(rep.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(rep.Changes))
	}
}

func TestAssessImpact_ColumnTypeChangedIsCritical(t *testing.T) {
	rep := AssessImpact(buildImpactResult(KindColumnTypeChanged))
	if rep.Overall != ImpactCritical {
		t.Errorf("expected ImpactCritical, got %v", rep.Overall)
	}
	if rep.Changes[0].Impact != ImpactCritical {
		t.Errorf("expected change impact critical")
	}
}

func TestAssessImpact_TableRemovedIsHigh(t *testing.T) {
	rep := AssessImpact(buildImpactResult(KindTableRemoved))
	if rep.Overall != ImpactHigh {
		t.Errorf("expected ImpactHigh, got %v", rep.Overall)
	}
}

func TestAssessImpact_OverallIsMaximum(t *testing.T) {
	rep := AssessImpact(buildImpactResult(KindTableAdded, KindColumnTypeChanged, KindColumnRemoved))
	if rep.Overall != ImpactCritical {
		t.Errorf("expected ImpactCritical as overall, got %v", rep.Overall)
	}
	if len(rep.Changes) != 3 {
		t.Errorf("expected 3 changes, got %d", len(rep.Changes))
	}
}

func TestImpactLevel_String(t *testing.T) {
	cases := []struct {
		level ImpactLevel
		want  string
	}{
		{ImpactNone, "none"},
		{ImpactLow, "low"},
		{ImpactMedium, "medium"},
		{ImpactHigh, "high"},
		{ImpactCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("ImpactLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestAssessImpact_ReasonNotEmpty(t *testing.T) {
	rep := AssessImpact(buildImpactResult(KindColumnTypeChanged))
	if rep.Changes[0].Reason == "" {
		t.Error("expected non-empty reason for critical change")
	}
}
