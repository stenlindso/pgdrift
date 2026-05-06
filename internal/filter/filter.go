package filter

import "strings"

// Options holds filtering configuration for schema comparison.
type Options struct {
	IncludeSchemas []string
	ExcludeSchemas []string
	IncludeTables  []string
	ExcludeTables  []string
}

// Filter applies inclusion/exclusion rules to schema and table names.
type Filter struct {
	opts Options
}

// New creates a new Filter with the given options.
func New(opts Options) *Filter {
	return &Filter{opts: opts}
}

// AllowSchema returns true if the given schema name passes the filter.
func (f *Filter) AllowSchema(schema string) bool {
	if len(f.opts.ExcludeSchemas) > 0 && containsIgnoreCase(f.opts.ExcludeSchemas, schema) {
		return false
	}
	if len(f.opts.IncludeSchemas) > 0 {
		return containsIgnoreCase(f.opts.IncludeSchemas, schema)
	}
	return true
}

// AllowTable returns true if the given table name passes the filter.
func (f *Filter) AllowTable(table string) bool {
	if len(f.opts.ExcludeTables) > 0 && containsIgnoreCase(f.opts.ExcludeTables, table) {
		return false
	}
	if len(f.opts.IncludeTables) > 0 {
		return containsIgnoreCase(f.opts.IncludeTables, table)
	}
	return true
}

func containsIgnoreCase(list []string, val string) bool {
	for _, item := range list {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}
