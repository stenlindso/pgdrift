package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/your-org/pgdrift/internal/diff"
)

// WriteDriftRate writes a human-readable drift rate report to w.
func WriteDriftRate(w io.Writer, dr *diff.DriftRate) {
	if dr == nil {
		fmt.Fprintln(w, "drift rate: no data")
		return
	}
	fmt.Fprintf(w, "Drift Rate Report\n")
	fmt.Fprintf(w, "  Window:       %s\n", formatWindow(dr))
	fmt.Fprintf(w, "  Total Runs:   %d\n", dr.TotalRuns)
	fmt.Fprintf(w, "  Drift Runs:   %d\n", dr.DriftRuns)
	fmt.Fprintf(w, "  Change Count: %d\n", dr.ChangeCount)
	fmt.Fprintf(w, "  Rate:         %.1f%%\n", dr.Rate()*100)
	fmt.Fprintf(w, "  Frequency:    %s\n", dr.Label())
}

// WriteDriftRateJSON writes the drift rate as JSON to w.
func WriteDriftRateJSON(w io.Writer, dr *diff.DriftRate) error {
	type payload struct {
		WindowStart time.Time `json:"window_start"`
		WindowEnd   time.Time `json:"window_end"`
		TotalRuns   int       `json:"total_runs"`
		DriftRuns   int       `json:"drift_runs"`
		ChangeCount int       `json:"change_count"`
		Rate        float64   `json:"rate"`
		Label       string    `json:"label"`
	}

	if dr == nil {
		dr = &diff.DriftRate{}
	}

	p := payload{
		WindowStart: dr.WindowStart,
		WindowEnd:   dr.WindowEnd,
		TotalRuns:   dr.TotalRuns,
		DriftRuns:   dr.DriftRuns,
		ChangeCount: dr.ChangeCount,
		Rate:        dr.Rate(),
		Label:       dr.Label(),
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

func formatWindow(dr *diff.DriftRate) string {
	if dr.WindowStart.IsZero() || dr.WindowEnd.IsZero() {
		return "unknown"
	}
	dur := dr.WindowDuration().Round(time.Minute)
	return fmt.Sprintf("%s → %s (%s)",
		dr.WindowStart.Format(time.RFC3339),
		dr.WindowEnd.Format(time.RFC3339),
		dur,
	)
}
