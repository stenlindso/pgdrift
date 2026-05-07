package diff

// SeverityLevel represents how critical a schema change is.
type SeverityLevel int

const (
	// SeverityLow indicates a backward-compatible change (e.g. new table, new nullable column).
	SeverityLow SeverityLevel = iota
	// SeverityMedium indicates a potentially breaking change (e.g. column default changed).
	SeverityMedium
	// SeverityHigh indicates a breaking change (e.g. column removed, type changed).
	SeverityHigh
)

// String returns a human-readable label for the severity level.
func (s SeverityLevel) String() string {
	switch s {
	case SeverityLow:
		return "LOW"
	case SeverityMedium:
		return "MEDIUM"
	case SeverityHigh:
		return "HIGH"
	default:
		return "UNKNOWN"
	}
}

// Severity returns the SeverityLevel for a given ChangeKind.
func Severity(kind ChangeKind) SeverityLevel {
	switch kind {
	case ChangeKindTableAdded:
		return SeverityLow
	case ChangeKindTableRemoved:
		return SeverityHigh
	case ChangeKindColumnAdded:
		return SeverityLow
	case ChangeKindColumnRemoved:
		return SeverityHigh
	case ChangeKindColumnTypeChanged:
		return SeverityHigh
	case ChangeKindColumnNullabilityChanged:
		return SeverityMedium
	case ChangeKindColumnDefaultChanged:
		return SeverityMedium
	default:
		return SeverityLow
	}
}
