package diff

import (
	"time"
)

// StalenessLevel indicates how out-of-date a schema snapshot is.
type StalenessLevel int

const (
	StalenessNone    StalenessLevel = iota // checked recently
	StalenessWarning                        // moderately stale
	StalenessCritical                       // very stale
)

func (s StalenessLevel) String() string {
	switch s {
	case StalenessNone:
		return "fresh"
	case StalenessWarning:
		return "warning"
	case StalenessCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// StalenessReport describes how stale a snapshot is relative to now.
type StalenessReport struct {
	CapturedAt time.Time
	Age        time.Duration
	Level      StalenessLevel
	Message    string
}

// StalenessOptions controls thresholds for staleness detection.
type StalenessOptions struct {
	WarningAfter  time.Duration // default: 24h
	CriticalAfter time.Duration // default: 72h
}

func defaultStalenessOptions() StalenessOptions {
	return StalenessOptions{
		WarningAfter:  24 * time.Hour,
		CriticalAfter: 72 * time.Hour,
	}
}

// AssessStaleness returns a StalenessReport for the given snapshot capture time.
func AssessStaleness(capturedAt time.Time, opts *StalenessOptions) *StalenessReport {
	if capturedAt.IsZero() {
		return &StalenessReport{
			Level:   StalenessCritical,
			Message: "snapshot has no capture timestamp",
		}
	}

	o := defaultStalenessOptions()
	if opts != nil {
		if opts.WarningAfter > 0 {
			o.WarningAfter = opts.WarningAfter
		}
		if opts.CriticalAfter > 0 {
			o.CriticalAfter = opts.CriticalAfter
		}
	}

	age := time.Since(capturedAt)
	r := &StalenessReport{
		CapturedAt: capturedAt,
		Age:        age,
	}

	switch {
	case age >= o.CriticalAfter:
		r.Level = StalenessCritical
		r.Message = "snapshot is critically stale"
	case age >= o.WarningAfter:
		r.Level = StalenessWarning
		r.Message = "snapshot may be outdated"
	default:
		r.Level = StalenessNone
		r.Message = "snapshot is fresh"
	}
	return r
}
