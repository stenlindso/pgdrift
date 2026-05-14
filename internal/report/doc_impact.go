// Package report provides formatting and output utilities for pgdrift results.
//
// Impact reporting (impact.go) translates a diff.ImpactReport into human-readable
// text or structured JSON output.
//
// # Text format
//
// WriteImpact writes a concise summary showing the overall impact level followed
// by a per-change breakdown with the assessed level and a plain-English reason:
//
//	 overall impact: critical
//	 ---
//	   [critical] public.orders column type changed — type change may break existing queries
//
// # JSON format
//
// WriteImpactJSON emits a JSON object suitable for machine consumption:
//
//	{
//	  "overall": "critical",
//	  "changes": [
//	    {
//	      "change": "...",
//	      "impact": "critical",
//	      "reason": "type change may break existing queries or application code"
//	    }
//	  ]
//	}
package report
