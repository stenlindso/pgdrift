package diff

import "fmt"

// Change describes a single schema drift item.
type Change struct {
	Schema      string
	Table       string
	Field       string
	Kind        ChangeKind
	OldValue    string
	NewValue    string
	Annotations []Annotation
}

// String returns a human-readable summary of the change.
func (c Change) String() string {
	if c.Field != "" {
		return fmt.Sprintf("%s.%s.%s: %s (%s -> %s)", c.Schema, c.Table, c.Field, c.Kind, c.OldValue, c.NewValue)
	}
	return fmt.Sprintf("%s.%s: %s", c.Schema, c.Table, c.Kind)
}

// Result holds all detected changes between two schemas.
type Result struct {
	Changes []Change
}

// HasDrift reports whether any changes were detected.
func (r *Result) HasDrift() bool {
	if r == nil {
		return false
	}
	return len(r.Changes) > 0
}

// ByKind returns all changes that match the given kind.
func (r *Result) ByKind(k ChangeKind) []Change {
	var out []Change
	for _, c := range r.Changes {
		if c.Kind == k {
			out = append(out, c)
		}
	}
	return out
}
