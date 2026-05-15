package diff

import (
	"fmt"
	"time"
)

// LineageEntry records a single schema state at a point in time.
type LineageEntry struct {
	Timestamp  time.Time
	Fingerprint string
	ChangeCount int
	DriftScore  float64
}

// Lineage tracks the historical progression of schema fingerprints over time,
// allowing callers to detect when a schema diverged and how it has evolved.
type Lineage struct {
	entries []LineageEntry
}

// NewLineage returns an empty Lineage.
func NewLineage() *Lineage {
	return &Lineage{}
}

// Record appends a new entry derived from the given Result.
// If result is nil the call is a no-op.
func (l *Lineage) Record(result *Result, fp string, ts time.Time) {
	if result == nil {
		return
	}
	s := Score(result)
	l.entries = append(l.entries, LineageEntry{
		Timestamp:   ts,
		Fingerprint: fp,
		ChangeCount: len(result.Changes),
		DriftScore:  s.Score,
	})
}

// Entries returns a copy of all recorded entries in chronological order.
func (l *Lineage) Entries() []LineageEntry {
	out := make([]LineageEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

// DivergencePoint returns the earliest entry where the fingerprint first
// differed from the initial recorded fingerprint, or an error if there are
// fewer than two entries.
func (l *Lineage) DivergencePoint() (LineageEntry, error) {
	if len(l.entries) < 2 {
		return LineageEntry{}, fmt.Errorf("lineage: need at least 2 entries, have %d", len(l.entries))
	}
	baseline := l.entries[0].Fingerprint
	for _, e := range l.entries[1:] {
		if e.Fingerprint != baseline {
			return e, nil
		}
	}
	return LineageEntry{}, fmt.Errorf("lineage: no divergence detected")
}

// Stable reports whether all recorded fingerprints are identical.
func (l *Lineage) Stable() bool {
	if len(l.entries) == 0 {
		return true
	}
	fp := l.entries[0].Fingerprint
	for _, e := range l.entries[1:] {
		if e.Fingerprint != fp {
			return false
		}
	}
	return true
}
