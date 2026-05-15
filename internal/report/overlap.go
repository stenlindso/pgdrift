package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/your-org/pgdrift/internal/diff"
)

// WriteOverlap writes a human-readable overlap report to w.
func WriteOverlap(w io.Writer, r *diff.OverlapReport) {
	if r == nil {
		fmt.Fprintln(w, "overlap: no report available")
		return
	}
	if r.TotalShared == 0 {
		fmt.Fprintln(w, "overlap: no shared tables detected")
		return
	}

	fmt.Fprintf(w, "overlap: %d shared table(s), %d conflicting column(s)\n",
		r.TotalShared, r.TotalConflicts)

	tables := make([]string, len(r.SharedTables))
	copy(tables, r.SharedTables)
	sort.Strings(tables)

	for _, tbl := range tables {
		cols, hasConflict := r.ConflictingColumns[tbl]
		if hasConflict && len(cols) > 0 {
			sort.Strings(cols)
			fmt.Fprintf(w, "  table %-30s conflicts: %v\n", tbl, cols)
		} else {
			fmt.Fprintf(w, "  table %-30s (no column conflicts)\n", tbl)
		}
	}
}

// WriteOverlapJSON writes the overlap report as JSON to w.
func WriteOverlapJSON(w io.Writer, r *diff.OverlapReport) error {
	if r == nil {
		r = &diff.OverlapReport{
			ConflictingColumns: make(map[string][]string),
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
