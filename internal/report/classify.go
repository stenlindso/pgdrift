package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pgdrift/pgdrift/internal/diff"
)

// WriteClassification writes a human-readable classification report to w.
func WriteClassification(w io.Writer, cr *diff.ClassificationReport) {
	if cr == nil {
		fmt.Fprintln(w, "Classification: unavailable")
		return
	}
	fmt.Fprintf(w, "Risk Classification : %s\n", cr.RiskClass)
	fmt.Fprintf(w, "Summary             : %s\n", cr.Summary)
	if cr.Result != nil {
		fmt.Fprintf(w, "Total Changes       : %d\n", len(cr.Result.Changes))
	}
}

// classificationJSON is the serialisable form of a ClassificationReport.
type classificationJSON struct {
	RiskClass    string `json:"risk_class"`
	Summary      string `json:"summary"`
	TotalChanges int    `json:"total_changes"`
}

// WriteClassificationJSON writes a JSON-encoded classification report to w.
func WriteClassificationJSON(w io.Writer, cr *diff.ClassificationReport) error {
	out := classificationJSON{}
	if cr != nil {
		out.RiskClass = cr.RiskClass.String()
		out.Summary = cr.Summary
		if cr.Result != nil {
			out.TotalChanges = len(cr.Result.Changes)
		}
	} else {
		out.RiskClass = "unknown"
		out.Summary = "no classification available"
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
