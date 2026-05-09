package report

import (
	"fmt"
	"io"

	"github.com/yourorg/pgdrift/internal/diff"
)

// WriteAnnotations writes a section listing all annotated changes to w.
// Only changes that carry at least one annotation are included.
func WriteAnnotations(w io.Writer, r *diff.Result) error {
	if r == nil || !r.HasDrift() {
		_, err := fmt.Fprintln(w, "No annotated changes.")
		return err
	}

	printed := false
	for _, c := range r.Changes {
		if len(c.Annotations) == 0 {
			continue
		}
		if !printed {
			fmt.Fprintln(w, "Annotated changes:")
			printed = true
		}
		fmt.Fprintf(w, "  %s\n", c.String())
		for _, a := range c.Annotations {
			fmt.Fprintf(w, "    [%s] %s\n", a.Key, a.Value)
		}
	}
	if !printed {
		fmt.Fprintln(w, "No annotated changes.")
	}
	return nil
}

// AnnotationSummary returns a map from annotation key to the count of changes
// carrying that key.
func AnnotationSummary(r *diff.Result) map[string]int {
	out := make(map[string]int)
	if r == nil {
		return out
	}
	for _, c := range r.Changes {
		seen := make(map[string]bool)
		for _, a := range c.Annotations {
			if !seen[a.Key] {
				out[a.Key]++
				seen[a.Key] = true
			}
		}
	}
	return out
}
