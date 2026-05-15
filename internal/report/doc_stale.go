// Package report provides formatting and output helpers for pgdrift results.
//
// # Staleness Reporting
//
// The staleness sub-feature assesses how out-of-date a schema snapshot is
// relative to the current wall-clock time.
//
// Use [AssessStaleness] (in package diff) to compute a [diff.StalenessReport]
// from a snapshot's capture timestamp, then render it with:
//
//	// Human-readable text
//	 report.WriteStale(os.Stdout, stalenessReport)
//
//	// Machine-readable JSON
//	 report.WriteStaleJSON(os.Stdout, stalenessReport)
//
// Staleness levels:
//
//	- fresh    – captured within the warning threshold (default: 24 h)
//	- warning  – between warning and critical thresholds (default: 24–72 h)
//	- critical – older than the critical threshold, or missing timestamp
package report
