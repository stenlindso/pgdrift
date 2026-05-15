package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/your-org/pgdrift/internal/diff"
)

// WriteStale writes a human-readable staleness report to w.
func WriteStale(w io.Writer, r *diff.StalenessReport) {
	if r == nil {
		fmt.Fprintln(w, "staleness: no report available")
		return
	}

	fmt.Fprintf(w, "staleness: %s\n", r.Level)
	fmt.Fprintf(w, "  message : %s\n", r.Message)

	if !r.CapturedAt.IsZero() {
		fmt.Fprintf(w, "  captured: %s\n", r.CapturedAt.Format(time.RFC3339))
		fmt.Fprintf(w, "  age     : %s\n", formatAge(r.Age))
	}
}

// WriteStaleJSON writes a JSON-encoded staleness report to w.
func WriteStaleJSON(w io.Writer, r *diff.StalenessReport) error {
	type payload struct {
		Level      string `json:"level"`
		Message    string `json:"message"`
		CapturedAt string `json:"captured_at,omitempty"`
		AgeSeconds int64  `json:"age_seconds,omitempty"`
	}

	if r == nil {
		return json.NewEncoder(w).Encode(payload{Level: "unknown", Message: "no report"})
	}

	p := payload{
		Level:      r.Level.String(),
		Message:    r.Message,
		AgeSeconds: int64(r.Age.Seconds()),
	}
	if !r.CapturedAt.IsZero() {
		p.CapturedAt = r.CapturedAt.Format(time.RFC3339)
	}
	return json.NewEncoder(w).Encode(p)
}

func formatAge(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
