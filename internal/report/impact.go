package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/your-org/pgdrift/internal/diff"
)

// WriteImpact writes a human-readable impact report to w.
func WriteImpact(w io.Writer, rep *diff.ImpactReport) {
	if rep == nil || len(rep.Changes) == 0 {
		fmt.Fprintln(w, "impact: no changes assessed")
		return
	}
	fmt.Fprintf(w, "overall impact: %s\n", rep.Overall)
	fmt.Fprintln(w, "---")
	for _, ic := range rep.Changes {
		fmt.Fprintf(w, "  [%s] %s — %s\n", ic.Impact, ic.Change, ic.Reason)
	}
}

// WriteImpactJSON writes the impact report as JSON to w.
func WriteImpactJSON(w io.Writer, rep *diff.ImpactReport) error {
	type jsonChange struct {
		Change string `json:"change"`
		Impact string `json:"impact"`
		Reason string `json:"reason"`
	}
	type jsonReport struct {
		Overall string       `json:"overall"`
		Changes []jsonChange `json:"changes"`
	}

	out := jsonReport{Overall: "none"}
	if rep != nil {
		out.Overall = rep.Overall.String()
		for _, ic := range rep.Changes {
			out.Changes = append(out.Changes, jsonChange{
				Change: ic.Change.String(),
				Impact: ic.Impact.String(),
				Reason: ic.Reason,
			})
		}
	}
	if out.Changes == nil {
		out.Changes = []jsonChange{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
