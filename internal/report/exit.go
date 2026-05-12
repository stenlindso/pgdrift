package report

import "github.com/your-org/pgdrift/internal/diff"

// ExitCode returns an appropriate OS exit code based on the drift result
// and a minimum severity threshold.
//
// Exit codes:
//   0 - no drift detected above threshold
//   1 - drift detected above threshold
//   2 - internal error (reserved for callers)
func ExitCode(result *diff.Result, threshold diff.SeverityLevel) int {
	if result == nil || !result.HasDrift() {
		return 0
	}
	for _, ch := range result.Changes {
		if diff.Severity(ch.Kind) >= threshold {
			return 1
		}
	}
	return 0
}

// ExitCodeStrict returns exit code 1 if any drift exists, regardless of severity.
func ExitCodeStrict(result *diff.Result) int {
	if result != nil && result.HasDrift() {
		return 1
	}
	return 0
}

// ExitCodeForSeverity returns exit code 1 if any change in the result matches
// exactly the given severity level. This is useful when callers want to trigger
// on a specific severity rather than a minimum threshold.
func ExitCodeForSeverity(result *diff.Result, level diff.SeverityLevel) int {
	if result == nil || !result.HasDrift() {
		return 0
	}
	for _, ch := range result.Changes {
		if diff.Severity(ch.Kind) == level {
			return 1
		}
	}
	return 0
}
