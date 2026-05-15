package diff

import "time"

// MaturityLevel describes how stable a schema element is over time.
type MaturityLevel int

const (
	MaturityUnknown MaturityLevel = iota
	MaturityNew                   // seen for fewer than MinRuns snapshots
	MaturityStable                // no changes over MinRuns snapshots
	MaturityFluctuating           // changed more than once recently
)

func (m MaturityLevel) String() string {
	switch m {
	case MaturityNew:
		return "new"
	case MaturityStable:
		return "stable"
	case MaturityFluctuating:
		return "fluctuating"
	default:
		return "unknown"
	}
}

// MaturityEntry records the maturity assessment for a single table.
type MaturityEntry struct {
	Schema   string        `json:"schema"`
	Table    string        `json:"table"`
	Level    MaturityLevel `json:"level"`
	Changes  int           `json:"changes"`
	Since    time.Time     `json:"since"`
}

// MaturityReport holds maturity assessments across all observed tables.
type MaturityReport struct {
	Entries   []MaturityEntry `json:"entries"`
	AsOf      time.Time       `json:"as_of"`
	MinRuns   int             `json:"min_runs"`
}

// AssessMaturity derives a MaturityReport from a Changelog.
// minRuns is the minimum number of recorded snapshots before a table is
// considered stable rather than new.
func AssessMaturity(cl *Changelog, minRuns int) *MaturityReport {
	if cl == nil {
		return &MaturityReport{AsOf: time.Now(), MinRuns: minRuns}
	}

	type tableKey struct{ schema, table string }
	changeCount := map[tableKey]int{}
	tableFirst := map[tableKey]time.Time{}

	for _, entry := range cl.Entries {
		for _, ch := range entry.Changes {
			k := tableKey{ch.Schema, ch.Table}
			changeCount[k]++
			if first, ok := tableFirst[k]; !ok || entry.RecordedAt.Before(first) {
				tableFirst[k] = entry.RecordedAt
			}
		}
	}

	report := &MaturityReport{
		AsOf:    time.Now(),
		MinRuns: minRuns,
	}

	for k, count := range changeCount {
		level := maturityLevel(count, len(cl.Entries), minRuns)
		report.Entries = append(report.Entries, MaturityEntry{
			Schema:  k.schema,
			Table:   k.table,
			Level:   level,
			Changes: count,
			Since:   tableFirst[k],
		})
	}
	return report
}

func maturityLevel(changes, totalRuns, minRuns int) MaturityLevel {
	if totalRuns < minRuns {
		return MaturityNew
	}
	if changes == 0 {
		return MaturityStable
	}
	if changes > 1 {
		return MaturityFluctuating
	}
	return MaturityNew
}
