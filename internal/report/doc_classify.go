// Package report provides formatted output for pgdrift analysis results.
//
// # Classification Reports
//
// The classify sub-feature assigns a broad RiskClass to a diff.Result:
//
//   - RiskClassNone     – no schema drift detected
//   - RiskClassLow      – minor, non-breaking structural changes
//   - RiskClassModerate – changes that may require review before deployment
//   - RiskClassCritical – breaking changes or potential data-loss risk
//
// Use WriteClassification to emit a plain-text summary, or
// WriteClassificationJSON to emit a machine-readable JSON object.
//
// Example:
//
//	cr := diff.Classify(result)
//	report.WriteClassification(os.Stdout, cr)
//	report.WriteClassificationJSON(os.Stdout, cr)
package report
