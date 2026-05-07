package diff

import (
	"testing"
)

func TestSeverityLevel_String(t *testing.T) {
	tests := []struct {
		level    SeverityLevel
		want     string
	}{
		{SeverityLow, "LOW"},
		{SeverityMedium, "MEDIUM"},
		{SeverityHigh, "HIGH"},
		{SeverityLevel(99), "UNKNOWN"},
	}
	for _, tc := range tests {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("SeverityLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestSeverity_ByKind(t *testing.T) {
	tests := []struct {
		kind ChangeKind
		want SeverityLevel
	}{
		{ChangeKindTableAdded, SeverityLow},
		{ChangeKindTableRemoved, SeverityHigh},
		{ChangeKindColumnAdded, SeverityLow},
		{ChangeKindColumnRemoved, SeverityHigh},
		{ChangeKindColumnTypeChanged, SeverityHigh},
		{ChangeKindColumnNullabilityChanged, SeverityMedium},
		{ChangeKindColumnDefaultChanged, SeverityMedium},
	}
	for _, tc := range tests {
		got := Severity(tc.kind)
		if got != tc.want {
			t.Errorf("Severity(%v) = %v, want %v", tc.kind, got, tc.want)
		}
	}
}

func TestSeverity_HigherThanLow(t *testing.T) {
	if SeverityHigh <= SeverityLow {
		t.Error("expected SeverityHigh > SeverityLow")
	}
	if SeverityMedium <= SeverityLow {
		t.Error("expected SeverityMedium > SeverityLow")
	}
	if SeverityHigh <= SeverityMedium {
		t.Error("expected SeverityHigh > SeverityMedium")
	}
}
