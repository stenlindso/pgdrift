package report

import (
	"fmt"
	"strings"

	"github.com/pgdrift/internal/diff"
)

// Severity returns a human-readable severity label for a change kind.
func Severity(kind diff.ChangeKind) string {
	switch kind {
	case diff.TableAdded, diff.TableRemoved:
		return "HIGH"
	case diff.ColumnAdded, diff.ColumnRemoved:
		return "MEDIUM"
	case diff.ColumnTypeChanged, diff.ColumnNullableChanged, diff.ColumnDefaultChanged:
		return "LOW"
	default:
		return "UNKNOWN"
	}
}

// Summary returns a one-line summary of a diff.Result.
func Summary(r diff.Result) string {
	if !r.HasDrift() {
		return "No schema drift detected."
	}
	high, medium, low := countBySeverity(r.Changes)
	parts := []string{}
	if high > 0 {
		parts = append(parts, fmt.Sprintf("%d high", high))
	}
	if medium > 0 {
		parts = append(parts, fmt.Sprintf("%d medium", medium))
	}
	if low > 0 {
		parts = append(parts, fmt.Sprintf("%d low", low))
	}
	return fmt.Sprintf("Schema drift detected: %d change(s) [%s].",
		len(r.Changes), strings.Join(parts, ", "))
}

func countBySeverity(changes []diff.Change) (high, medium, low int) {
	for _, c := range changes {
		switch Severity(c.Kind) {
		case "HIGH":
			high++
		case "MEDIUM":
			medium++
		default:
			low++
		}
	}
	return
}
