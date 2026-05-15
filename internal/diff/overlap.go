package diff

// OverlapReport describes tables and columns that appear in both schemas
// but carry conflicting or duplicated definitions across logical groups.
type OverlapReport struct {
	// SharedTables lists table names present in both schemas.
	SharedTables []string `json:"shared_tables"`
	// ConflictingColumns maps table name to columns whose type differs.
	ConflictingColumns map[string][]string `json:"conflicting_columns"`
	// TotalShared is the count of shared tables.
	TotalShared int `json:"total_shared"`
	// TotalConflicts is the total number of conflicting columns.
	TotalConflicts int `json:"total_conflicts"`
}

// DetectOverlap compares two Results and identifies tables and columns that
// exist in both but have type conflicts, providing a higher-level overlap view
// beyond the raw change list.
func DetectOverlap(source, target *Result) *OverlapReport {
	report := &OverlapReport{
		ConflictingColumns: make(map[string][]string),
	}
	if source == nil || target == nil {
		return report
	}

	// Index source changes by table for quick lookup.
	sourceByTable := make(map[string][]*Change)
	for _, c := range source.Changes {
		sourceByTable[c.Table] = append(sourceByTable[c.Table], c)
	}

	seen := make(map[string]bool)
	for _, c := range target.Changes {
		if _, ok := sourceByTable[c.Table]; ok {
			if !seen[c.Table] {
				report.SharedTables = append(report.SharedTables, c.Table)
				seen[c.Table] = true
			}
			if c.Kind == KindColumnTypeChanged {
				report.ConflictingColumns[c.Table] = append(
					report.ConflictingColumns[c.Table], c.Column,
				)
				report.TotalConflicts++
			}
		}
	}
	report.TotalShared = len(report.SharedTables)
	return report
}
