package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/your-org/pgdrift/internal/diff"
)

// WriteMaturity writes a human-readable maturity report to w.
func WriteMaturity(w io.Writer, r *diff.MaturityReport) {
	if r == nil || len(r.Entries) == 0 {
		fmt.Fprintln(w, "maturity: no data available")
		return
	}

	entries := make([]diff.MaturityEntry, len(r.Entries))
	copy(entries, r.Entries)
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Schema != entries[j].Schema {
			return entries[i].Schema < entries[j].Schema
		}
		return entries[i].Table < entries[j].Table
	})

	fmt.Fprintf(w, "Schema Maturity Report (min_runs=%d)\n", r.MinRuns)
	fmt.Fprintf(w, "As of: %s\n\n", r.AsOf.Format("2006-01-02 15:04:05"))

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SCHEMA\tTABLE\tLEVEL\tCHANGES")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", e.Schema, e.Table, e.Level.String(), e.Changes)
	}
	_ = tw.Flush()
}

// WriteMaturityJSON writes the maturity report as JSON to w.
func WriteMaturityJSON(w io.Writer, r *diff.MaturityReport) error {
	if r == nil {
		_, err := fmt.Fprintln(w, "{}")
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
