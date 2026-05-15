package diff

import "time"

// DriftRate measures how frequently schema drift occurs over a time window.
type DriftRate struct {
	WindowStart time.Time
	WindowEnd   time.Time
	TotalRuns   int
	DriftRuns   int
	ChangeCount int
}

// Rate returns the fraction of runs that contained drift (0.0–1.0).
func (d *DriftRate) Rate() float64 {
	if d.TotalRuns == 0 {
		return 0.0
	}
	return float64(d.DriftRuns) / float64(d.TotalRuns)
}

// WindowDuration returns the duration of the measurement window.
func (d *DriftRate) WindowDuration() time.Duration {
	return d.WindowEnd.Sub(d.WindowStart)
}

// Label returns a human-readable drift frequency label.
func (d *DriftRate) Label() string {
	r := d.Rate()
	switch {
	case r == 0:
		return "stable"
	case r < 0.25:
		return "infrequent"
	case r < 0.60:
		return "moderate"
	default:
		return "frequent"
	}
}

// ComputeDriftRate calculates drift rate from a Changelog over the given window.
func ComputeDriftRate(cl *Changelog, window time.Duration) *DriftRate {
	if cl == nil || len(cl.Entries) == 0 {
		return &DriftRate{}
	}

	now := time.Now().UTC()
	cutoff := now.Add(-window)

	dr := &DriftRate{
		WindowStart: cutoff,
		WindowEnd:   now,
	}

	for _, e := range cl.Entries {
		if e.RecordedAt.Before(cutoff) {
			continue
		}
		dr.TotalRuns++
		if e.ChangeCount > 0 {
			dr.DriftRuns++
			dr.ChangeCount += e.ChangeCount
		}
	}

	return dr
}
