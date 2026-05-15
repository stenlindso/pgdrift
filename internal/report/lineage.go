package report

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/your-org/pgdrift/internal/diff"
)

// WriteLineage writes a human-readable lineage table to w.
func WriteLineage(w io.Writer, l *diff.Lineage) {
	if l == nil {
		fmt.Fprintln(w, "no lineage data")
		return
	}
	entries := l.Entries()
	if len(entries) == 0 {
		fmt.Fprintln(w, "no lineage entries recorded")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tFINGERPRINT\tCHANGES\tSCORE")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%.2f\n",
			e.Timestamp.Format(time.RFC3339),
			truncate(e.Fingerprint, 12),
			e.ChangeCount,
			e.DriftScore,
		)
	}
	_ = tw.Flush()

	if l.Stable() {
		fmt.Fprintln(w, "\nSchema is stable across all recorded snapshots.")
	} else {
		fmt.Fprintln(w, "\nSchema drift detected in lineage.")
	}
}

// WriteLineageJSON writes the lineage entries as a JSON array to w.
func WriteLineageJSON(w io.Writer, l *diff.Lineage) error {
	type jsonEntry struct {
		Timestamp   string  `json:"timestamp"`
		Fingerprint string  `json:"fingerprint"`
		ChangeCount int     `json:"change_count"`
		DriftScore  float64 `json:"drift_score"`
	}

	var rows []jsonEntry
	if l != nil {
		for _, e := range l.Entries() {
			rows = append(rows, jsonEntry{
				Timestamp:   e.Timestamp.Format(time.RFC3339),
				Fingerprint: e.Fingerprint,
				ChangeCount: e.ChangeCount,
				DriftScore:  e.DriftScore,
			})
		}
	}
	if rows == nil {
		rows = []jsonEntry{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
