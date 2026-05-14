package diff

import "fmt"

// DriftScore represents a numeric health score for a schema comparison.
// 100 means no drift; lower values indicate more severe drift.
type DriftScore struct {
	Value   int    // 0–100
	Grade   string // A, B, C, D, F
	Summary string
}

// Score computes a drift health score from a Result.
// Each change deducts points based on its severity:
//   - High:   10 points
//   - Medium:  5 points
//   - Low:     2 points
func Score(r *Result) DriftScore {
	if r == nil || !r.HasDrift() {
		return DriftScore{Value: 100, Grade: "A", Summary: "No drift detected"}
	}

	deduction := 0
	for _, c := range r.Changes {
		switch Severity(c.Kind) {
		case High:
			deduction += 10
		case Medium:
			deduction += 5
		case Low:
			deduction += 2
		}
	}

	value := 100 - deduction
	if value < 0 {
		value = 0
	}

	grade := scoreGrade(value)
	summary := fmt.Sprintf("%d change(s) detected; score reduced by %d point(s)", len(r.Changes), deduction)

	return DriftScore{Value: value, Grade: grade, Summary: summary}
}

func scoreGrade(v int) string {
	switch {
	case v >= 90:
		return "A"
	case v >= 75:
		return "B"
	case v >= 60:
		return "C"
	case v >= 40:
		return "D"
	default:
		return "F"
	}
}
