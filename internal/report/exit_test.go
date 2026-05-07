package report_test

import (
	"testing"

	"github.com/your-org/pgdrift/internal/diff"
	"github.com/your-org/pgdrift/internal/report"
)

func emptyResult() *diff.Result {
	return &diff.Result{}
}

func resultWithKind(kind diff.ChangeKind) *diff.Result {
	return &diff.Result{
		Changes: []diff.Change{
			{Kind: kind, Table: "public.users"},
		},
	}
}

func TestExitCode_NoDrift(t *testing.T) {
	code := report.ExitCode(emptyResult(), diff.SeverityLow)
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
}

func TestExitCode_NilResult(t *testing.T) {
	code := report.ExitCode(nil, diff.SeverityLow)
	if code != 0 {
		t.Errorf("expected 0 for nil result, got %d", code)
	}
}

func TestExitCode_BelowThreshold(t *testing.T) {
	result := resultWithKind(diff.ChangeKindColumnDefault)
	code := report.ExitCode(result, diff.SeverityHigh)
	if code != 0 {
		t.Errorf("expected 0 when severity below threshold, got %d", code)
	}
}

func TestExitCode_AtThreshold(t *testing.T) {
	result := resultWithKind(diff.ChangeKindTableAdded)
	code := report.ExitCode(result, diff.SeverityHigh)
	if code != 1 {
		t.Errorf("expected 1 when severity meets threshold, got %d", code)
	}
}

func TestExitCode_AboveThreshold(t *testing.T) {
	result := resultWithKind(diff.ChangeKindTableRemoved)
	code := report.ExitCode(result, diff.SeverityLow)
	if code != 1 {
		t.Errorf("expected 1 when severity above threshold, got %d", code)
	}
}

func TestExitCodeStrict_NoDrift(t *testing.T) {
	if report.ExitCodeStrict(emptyResult()) != 0 {
		t.Error("expected 0 for no drift")
	}
}

func TestExitCodeStrict_WithDrift(t *testing.T) {
	result := resultWithKind(diff.ChangeKindColumnDefault)
	if report.ExitCodeStrict(result) != 1 {
		t.Error("expected 1 for any drift in strict mode")
	}
}

func TestExitCodeStrict_Nil(t *testing.T) {
	if report.ExitCodeStrict(nil) != 0 {
		t.Error("expected 0 for nil result")
	}
}
