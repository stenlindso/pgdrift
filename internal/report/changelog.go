package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/pgdrift/internal/diff"
)

// WriteChangelog writes a human-readable changelog summary to w.
func WriteChangelog(w io.Writer, c *diff.Changelog) {
	if c == nil || c.Len() == 0 {
		fmt.Fprintln(w, "changelog: no entries recorded")
		return
	}
	fmt.Fprintf(w, "changelog: %d entries\n", c.Len())
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for i, e := range c.Entries {
		label := e.Label
		if label == "" {
			label = "(unlabeled)"
		}
		fmt.Fprintf(w, "[%d] %s  label=%s\n", i+1, e.Timestamp.Format("2006-01-02T15:04:05Z"), label)
		fmt.Fprintf(w, "    changes=%d  score=%.1f (%s)\n",
			e.Summary.TotalChanges, e.Score.Score, e.Score.Grade)
		if len(e.Summary.AffectedTables) > 0 {
			fmt.Fprintf(w, "    tables: %s\n", strings.Join(e.Summary.AffectedTables, ", "))
		}
	}
	fmt.Fprintln(w, strings.Repeat("-", 40))
	top := c.TopChanged(3)
	if len(top) > 0 {
		fmt.Fprintf(w, "most changed tables: %s\n", strings.Join(top, ", "))
	}
}

// WriteChangelogJSON writes the changelog as a JSON array to w.
func WriteChangelogJSON(w io.Writer, c *diff.Changelog) error {
	type entry struct {
		Timestamp string  `json:"timestamp"`
		Label     string  `json:"label,omitempty"`
		Changes   int     `json:"total_changes"`
		Score     float64 `json:"score"`
		Grade     string  `json:"grade"`
	}
	var rows []entry
	if c != nil {
		for _, e := range c.Entries {
			rows = append(rows, entry{
				Timestamp: e.Timestamp.Format("2006-01-02T15:04:05Z"),
				Label:     e.Label,
				Changes:   e.Summary.TotalChanges,
				Score:     e.Score.Score,
				Grade:     e.Score.Grade,
			})
		}
	}
	if rows == nil {
		rows = []entry{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
