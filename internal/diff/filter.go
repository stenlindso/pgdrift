package diff

import "github.com/pgdrift/pgdrift/internal/filter"

// FilterResult returns a new Result containing only the changes that pass
// the provided filter. Schema and table names are checked against the filter
// rules so that ignored schemas/tables are excluded from drift reports.
func FilterResult(result Result, f *filter.Filter) Result {
	if f == nil {
		return result
	}

	filtered := Result{}
	for _, ch := range result.Changes {
		if !f.AllowSchema(ch.Schema) {
			continue
		}
		if ch.Table != "" && !f.AllowTable(ch.Table) {
			continue
		}
		filtered.Changes = append(filtered.Changes, ch)
	}
	return filtered
}
