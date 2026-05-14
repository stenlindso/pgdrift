package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/your-org/pgdrift/internal/diff"
)

// WriteScore writes a human-readable drift score to w.
func WriteScore(w io.Writer, s diff.DriftScore) error {
	_, err := fmt.Fprintf(w,
		"Drift Score: %d/100 (Grade: %s)\n%s\n",
		s.Value, s.Grade, s.Summary,
	)
	return err
}

// WriteScoreJSON writes the drift score as a JSON object to w.
func WriteScoreJSON(w io.Writer, s diff.DriftScore) error {
	type jsonScore struct {
		Value   int    `json:"value"`
		Grade   string `json:"grade"`
		Summary string `json:"summary"`
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonScore{
		Value:   s.Value,
		Grade:   s.Grade,
		Summary: s.Summary,
	})
}
