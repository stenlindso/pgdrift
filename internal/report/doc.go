// Package report provides writers and formatting utilities for pgdrift
// schema-drift results.
//
// Writers
//
// NewWriter returns a Writer that can render a diff.Result in either
// plain-text or JSON format depending on the format string passed at
// construction time ("text" or "json").
//
// Formatting helpers
//
// Severity maps a diff.ChangeKind to a human-readable severity label
// (HIGH / MEDIUM / LOW) that can be used in reports or CI integrations.
//
// Summary produces a concise one-line string describing the overall
// drift status of a diff.Result, suitable for log output or a terminal
// status line.
package report
