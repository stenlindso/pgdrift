package diff

import (
	"testing"
)

func buildClassifyResult(kinds ...ChangeKind) *Result {
	r := &Result{}
	for _, k := range kinds {
		r.Changes = append(r.Changes, Change{
			Kind:   k,
			Schema: "public",
			Table:  "t",
		})
	}
	return r
}

func TestRiskClass_String(t *testing.T) {
	cases := []struct {
		rc   RiskClass
		want string
	}{
		{RiskClassNone, "none"},
		{RiskClassLow, "low"},
		{RiskClassModerate, "moderate"},
		{RiskClassCritical, "critical"},
		{RiskClass(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.rc.String(); got != tc.want {
			t.Errorf("RiskClass(%d).String() = %q, want %q", tc.rc, got, tc.want)
		}
	}
}

func TestClassifyResult_NilResult(t *testing.T) {
	if got := ClassifyResult(nil); got != RiskClassNone {
		t.Errorf("expected RiskClassNone for nil result, got %v", got)
	}
}

func TestClassifyResult_NoChanges(t *testing.T) {
	if got := ClassifyResult(&Result{}); got != RiskClassNone {
		t.Errorf("expected RiskClassNone for empty result, got %v", got)
	}
}

func TestClassifyResult_LowRisk(t *testing.T) {
	r := buildClassifyResult(KindIndexAdded)
	if got := ClassifyResult(r); got != RiskClassLow {
		t.Errorf("expected RiskClassLow, got %v", got)
	}
}

func TestClassifyResult_ModerateRisk(t *testing.T) {
	r := buildClassifyResult(KindColumnAdded)
	if got := ClassifyResult(r); got < RiskClassLow {
		t.Errorf("expected at least RiskClassLow, got %v", got)
	}
}

func TestClassifyResult_CriticalRisk(t *testing.T) {
	r := buildClassifyResult(KindTableRemoved)
	if got := ClassifyResult(r); got != RiskClassCritical {
		t.Errorf("expected RiskClassCritical for table removed, got %v", got)
	}
}

func TestClassify_ReturnsSummary(t *testing.T) {
	report := Classify(nil)
	if report == nil {
		t.Fatal("expected non-nil ClassificationReport")
	}
	if report.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestClassify_CriticalSummaryText(t *testing.T) {
	r := buildClassifyResult(KindColumnTypeChanged)
	report := Classify(r)
	if report.RiskClass != RiskClassCritical {
		t.Errorf("expected critical, got %v", report.RiskClass)
	}
	if report.Summary == "" {
		t.Error("expected non-empty summary")
	}
}
