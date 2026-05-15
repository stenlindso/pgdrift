package diff

import (
	"testing"
	"time"
)

func TestStalenessLevel_String(t *testing.T) {
	cases := []struct {
		level StalenessLevel
		want  string
	}{
		{StalenessNone, "fresh"},
		{StalenessWarning, "warning"},
		{StalenessCritical, "critical"},
		{StalenessLevel(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.level.String(); got != c.want {
			t.Errorf("String() = %q, want %q", got, c.want)
		}
	}
}

func TestAssessStaleness_ZeroTime(t *testing.T) {
	r := AssessStaleness(time.Time{}, nil)
	if r.Level != StalenessCritical {
		t.Errorf("expected critical for zero time, got %s", r.Level)
	}
}

func TestAssessStaleness_Fresh(t *testing.T) {
	now := time.Now().Add(-1 * time.Hour)
	r := AssessStaleness(now, nil)
	if r.Level != StalenessNone {
		t.Errorf("expected fresh, got %s", r.Level)
	}
	if r.Age <= 0 {
		t.Error("expected positive age")
	}
}

func TestAssessStaleness_Warning(t *testing.T) {
	ts := time.Now().Add(-30 * time.Hour)
	r := AssessStaleness(ts, nil)
	if r.Level != StalenessWarning {
		t.Errorf("expected warning, got %s", r.Level)
	}
}

func TestAssessStaleness_Critical(t *testing.T) {
	ts := time.Now().Add(-100 * time.Hour)
	r := AssessStaleness(ts, nil)
	if r.Level != StalenessCritical {
		t.Errorf("expected critical, got %s", r.Level)
	}
}

func TestAssessStaleness_CustomOptions(t *testing.T) {
	opts := &StalenessOptions{
		WarningAfter:  1 * time.Hour,
		CriticalAfter: 2 * time.Hour,
	}
	ts := time.Now().Add(-90 * time.Minute)
	r := AssessStaleness(ts, opts)
	if r.Level != StalenessWarning {
		t.Errorf("expected warning with custom opts, got %s", r.Level)
	}
}

func TestAssessStaleness_MessageSet(t *testing.T) {
	ts := time.Now().Add(-1 * time.Minute)
	r := AssessStaleness(ts, nil)
	if r.Message == "" {
		t.Error("expected non-empty message")
	}
}
