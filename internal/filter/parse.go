package filter

import "strings"

// ParseOptions builds an Options struct from raw CLI flag strings.
// Each argument is expected to be a comma-separated list of names.
func ParseOptions(
	includeSchemas string,
	excludeSchemas string,
	includeTables string,
	excludeTables string,
) Options {
	return Options{
		IncludeSchemas: splitCSV(includeSchemas),
		ExcludeSchemas: splitCSV(excludeSchemas),
		IncludeTables:  splitCSV(includeTables),
		ExcludeTables:  splitCSV(excludeTables),
	}
}

func splitCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
