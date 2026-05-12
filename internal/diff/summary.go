package diff

// Summary holds aggregated statistics about a diff Result.
type Summary struct {
	TotalChanges int
	ByKind       map[ChangeKind]int
	BySeverity   map[SeverityLevel]int
	AffectedTables []string
}

// Summarize computes a Summary from a Result.
func Summarize(r *Result) *Summary {
	if r == nil {
		return &Summary{
			ByKind:     make(map[ChangeKind]int),
			BySeverity: make(map[SeverityLevel]int),
		}
	}

	s := &Summary{
		ByKind:     make(map[ChangeKind]int),
		BySeverity: make(map[SeverityLevel]int),
	}

	seen := make(map[string]bool)

	for _, c := range r.Changes {
		s.TotalChanges++
		s.ByKind[c.Kind]++
		lvl := Severity(c.Kind)
		s.BySeverity[lvl]++

		key := c.Schema + "." + c.Table
		if !seen[key] && (c.Schema != "" || c.Table != "") {
			seen[key] = true
			s.AffectedTables = append(s.AffectedTables, key)
		}
	}

	return s
}

// HighestSeverity returns the highest SeverityLevel present in the summary.
// Returns SeverityLow if there are no changes.
func (s *Summary) HighestSeverity() SeverityLevel {
	if s.BySeverity[SeverityHigh] > 0 {
		return SeverityHigh
	}
	if s.BySeverity[SeverityMedium] > 0 {
		return SeverityMedium
	}
	return SeverityLow
}
