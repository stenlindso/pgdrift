package diff

// ImpactLevel represents the estimated operational impact of a schema change.
type ImpactLevel int

const (
	ImpactNone   ImpactLevel = iota
	ImpactLow                // e.g. adding a nullable column
	ImpactMedium             // e.g. adding a non-nullable column with default
	ImpactHigh               // e.g. dropping a column or table
	ImpactCritical           // e.g. type change on indexed/PK column
)

func (l ImpactLevel) String() string {
	switch l {
	case ImpactLow:
		return "low"
	case ImpactMedium:
		return "medium"
	case ImpactHigh:
		return "high"
	case ImpactCritical:
		return "critical"
	default:
		return "none"
	}
}

// ImpactReport holds the assessed impact for each change in a Result.
type ImpactReport struct {
	Changes []ImpactedChange
	Overall ImpactLevel
}

// ImpactedChange pairs a Change with its assessed ImpactLevel.
type ImpactedChange struct {
	Change  Change
	Impact  ImpactLevel
	Reason  string
}

// AssessImpact evaluates each change in the result and returns an ImpactReport.
func AssessImpact(r *Result) *ImpactReport {
	if r == nil {
		return &ImpactReport{}
	}
	report := &ImpactReport{}
	overall := ImpactNone
	for _, c := range r.Changes {
		ic := ImpactedChange{Change: c, Impact: impactForChange(c)}
		ic.Reason = impactReason(ic.Impact, c)
		if ic.Impact > overall {
			overall = ic.Impact
		}
		report.Changes = append(report.Changes, ic)
	}
	report.Overall = overall
	return report
}

func impactForChange(c Change) ImpactLevel {
	switch c.Kind {
	case KindTableRemoved, KindColumnRemoved:
		return ImpactHigh
	case KindColumnTypeChanged:
		return ImpactCritical
	case KindColumnAdded:
		return ImpactMedium
	case KindTableAdded:
		return ImpactLow
	default:
		return ImpactLow
	}
}

func impactReason(level ImpactLevel, c Change) string {
	switch level {
	case ImpactCritical:
		return "type change may break existing queries or application code"
	case ImpactHigh:
		return "removal may cause runtime errors in dependent code"
	case ImpactMedium:
		return "new column may require migration or application update"
	case ImpactLow:
		return "additive change with minimal risk"
	default:
		return ""
	}
}
