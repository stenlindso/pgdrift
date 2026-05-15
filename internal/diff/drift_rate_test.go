package diff

import (
	"testing"
	"time"
)

func buildDriftRateChangelog(entries []ChangelogEntry) *Changelog {
	return &Changelog{Entries: entries}
}

func TestComputeDriftRate_NilChangelog(t *testing.T) {
	dr := ComputeDriftRate(nil, 24*time.Hour)
	if dr.TotalRuns != 0 || dr.DriftRuns != 0 {
		t.Errorf("expected empty rate, got %+v", dr)
	}
}

func TestComputeDriftRate_NoEntriesInWindow(t *testing.T) {
	old := time.Now().UTC().Add(-72 * time.Hour)
	cl := buildDriftRateChangelog([]ChangelogEntry{
		{RecordedAt: old, ChangeCount: 3},
	})
	dr := ComputeDriftRate(cl, 24*time.Hour)
	if dr.TotalRuns != 0 {
		t.Errorf("expected 0 runs in window, got %d", dr.TotalRuns)
	}
}

func TestComputeDriftRate_AllDrift(t *testing.T) {
	now := time.Now().UTC()
	cl := buildDriftRateChangelog([]ChangelogEntry{
		{RecordedAt: now.Add(-1 * time.Hour), ChangeCount: 2},
		{RecordedAt: now.Add(-2 * time.Hour), ChangeCount: 1},
	})
	dr := ComputeDriftRate(cl, 24*time.Hour)
	if dr.TotalRuns != 2 {
		t.Errorf("expected 2 total runs, got %d", dr.TotalRuns)
	}
	if dr.DriftRuns != 2 {
		t.Errorf("expected 2 drift runs, got %d", dr.DriftRuns)
	}
	if dr.Rate() != 1.0 {
		t.Errorf("expected rate 1.0, got %f", dr.Rate())
	}
}

func TestComputeDriftRate_MixedRuns(t *testing.T) {
	now := time.Now().UTC()
	cl := buildDriftRateChangelog([]ChangelogEntry{
		{RecordedAt: now.Add(-1 * time.Hour), ChangeCount: 0},
		{RecordedAt: now.Add(-2 * time.Hour), ChangeCount: 1},
		{RecordedAt: now.Add(-3 * time.Hour), ChangeCount: 0},
		{RecordedAt: now.Add(-4 * time.Hour), ChangeCount: 2},
	})
	dr := ComputeDriftRate(cl, 24*time.Hour)
	if dr.TotalRuns != 4 {
		t.Errorf("expected 4 runs, got %d", dr.TotalRuns)
	}
	if dr.DriftRuns != 2 {
		t.Errorf("expected 2 drift runs, got %d", dr.DriftRuns)
	}
	if dr.ChangeCount != 3 {
		t.Errorf("expected 3 total changes, got %d", dr.ChangeCount)
	}
}

func TestDriftRate_Label(t *testing.T) {
	cases := []struct {
		total, drift int
		want         string
	}{
		{0, 0, "stable"},
		{10, 0, "stable"},
		{10, 2, "infrequent"},
		{10, 5, "moderate"},
		{10, 9, "frequent"},
	}
	for _, c := range cases {
		dr := &DriftRate{TotalRuns: c.total, DriftRuns: c.drift}
		if got := dr.Label(); got != c.want {
			t.Errorf("Label() = %q, want %q (total=%d drift=%d)", got, c.want, c.total, c.drift)
		}
	}
}
