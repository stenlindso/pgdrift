package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/pgdrift/internal/diff"
)

// WriteTrend writes a human-readable trend summary to w.
// It shows each recorded point with its timestamp, total changes, and delta
// relative to the previous point.
func WriteTrend(w io.Writer, t *diff.Trend) error {
	if t == nil || len(t.Points) == 0 {
		_, err := fmt.Fprintln(w, "No trend data recorded.")
		return err
	}

	_, err := fmt.Fprintf(w, "Drift Trend (%d points)\n", len(t.Points))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, strings.Repeat("-", 40))
	if err != nil {
		return err
	}

	for i, pt := range t.Points {
		delta := ""
		if i > 0 {
			d := t.Points[i].TotalChanges - t.Points[i-1].TotalChanges
			switch {
			case d > 0:
				delta = fmt.Sprintf(" (+%d)", d)
			case d < 0:
				delta = fmt.Sprintf(" (%d)", d)
			default:
				delta = " (no change)"
			}
		}
		_, err = fmt.Fprintf(w, "[%s] changes: %d%s\n",
			pt.Timestamp.Format("2006-01-02 15:04:05"),
			pt.TotalChanges,
			delta,
		)
		if err != nil {
			return err
		}
		for kind, count := range pt.ByKind {
			_, err = fmt.Fprintf(w, "    %-30s %d\n", kind, count)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
