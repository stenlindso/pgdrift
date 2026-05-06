// Package filter provides schema and table name filtering for pgdrift.
//
// It supports inclusion and exclusion lists for both schema names and table
// names, allowing users to narrow the scope of drift detection to relevant
// parts of their PostgreSQL databases.
//
// Usage:
//
//	opts := filter.ParseOptions(includeSchemas, excludeSchemas, includeTables, excludeTables)
//	f := filter.New(opts)
//	if f.AllowTable("users") {
//	    // compare this table
//	}
package filter
