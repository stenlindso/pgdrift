package report

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/pgdrift/pgdrift/internal/diff"
)

// Format represents the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Writer renders a diff.Result to an io.Writer in the specified format.
type Writer struct {
	format Format
	out    io.Writer
}

// NewWriter creates a new report Writer.
func NewWriter(out io.Writer, format Format) *Writer {
	return &Writer{out: out, format: format}
}

// Write renders the diff result.
func (w *Writer) Write(result *diff.Result) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(result)
	default:
		return w.writeText(result)
	}
}

func (w *Writer) writeText(result *diff.Result) error {
	if !result.HasDrift() {
		_, err := fmt.Fprintln(w.out, "No schema drift detected.")
		return err
	}

	tw := tabwriter.NewWriter(w.out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "TYPE\tOBJECT\tDETAIL")
	_, _ = fmt.Fprintln(tw, "----\t------\t------")
	for _, c := range result.Changes {
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\n", c.ChangeType, c.Object, c.Detail)
	}
	return tw.Flush()
}

type jsonChange struct {
	Type   diff.ChangeType `json:"type"`
	Object string          `json:"object"`
	Detail string          `json:"detail"`
}

type jsonReport struct {
	DriftDetected bool         `json:"drift_detected"`
	Changes       []jsonChange `json:"changes"`
}

func (w *Writer) writeJSON(result *diff.Result) error {
	report := jsonReport{
		DriftDetected: result.HasDrift(),
		Changes:       make([]jsonChange, 0, len(result.Changes)),
	}
	for _, c := range result.Changes {
		report.Changes = append(report.Changes, jsonChange{
			Type:   c.ChangeType,
			Object: c.Object,
			Detail: c.Detail,
		})
	}
	enc := json.NewEncoder(w.out)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
