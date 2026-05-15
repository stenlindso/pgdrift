package diff

// RiskClass represents a broad risk classification for a set of changes.
type RiskClass int

const (
	RiskClassNone     RiskClass = iota // No changes detected
	RiskClassLow                       // Minor, non-breaking changes
	RiskClassModerate                  // Potentially disruptive changes
	RiskClassCritical                  // Breaking or data-loss risk changes
)

// String returns a human-readable label for the RiskClass.
func (r RiskClass) String() string {
	switch r {
	case RiskClassNone:
		return "none"
	case RiskClassLow:
		return "low"
	case RiskClassModerate:
		return "moderate"
	case RiskClassCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// ClassifyResult derives an overall RiskClass from a diff Result.
// It inspects every change's severity and escalates accordingly.
func ClassifyResult(r *Result) RiskClass {
	if r == nil || len(r.Changes) == 0 {
		return RiskClassNone
	}

	worst := RiskClassLow
	for _, c := range r.Changes {
		sev := Severity(c)
		switch {
		case sev == SeverityHigh && worst < RiskClassCritical:
			worst = RiskClassCritical
		case sev == SeverityMedium && worst < RiskClassModerate:
			worst = RiskClassModerate
		}
	}
	return worst
}

// ClassificationReport bundles a Result with its computed RiskClass.
type ClassificationReport struct {
	Result    *Result
	RiskClass RiskClass
	Summary   string
}

// Classify builds a ClassificationReport for the given Result.
func Classify(r *Result) *ClassificationReport {
	rc := ClassifyResult(r)
	return &ClassificationReport{
		Result:    r,
		RiskClass: rc,
		Summary:   classificationSummary(rc, r),
	}
}

func classificationSummary(rc RiskClass, r *Result) string {
	if r == nil || len(r.Changes) == 0 {
		return "No schema drift detected."
	}
	switch rc {
	case RiskClassCritical:
		return "Critical drift detected: breaking changes present."
	case RiskClassModerate:
		return "Moderate drift detected: review recommended before deployment."
	default:
		return "Low-risk drift detected: minor structural changes only."
	}
}
