package diff

import "fmt"

// ChangeType represents the kind of schema change detected.
type ChangeType string

const (
	ChangeTypeTableAdded    ChangeType = "table_added"
	ChangeTypeTableRemoved  ChangeType = "table_removed"
	ChangeTypeColumnAdded   ChangeType = "column_added"
	ChangeTypeColumnRemoved ChangeType = "column_removed"
	ChangeTypeColumnChanged ChangeType = "column_changed"
)

// Change describes a single schema drift item.
type Change struct {
	Type    ChangeType `json:"type"`
	Table   string     `json:"table"`
	Column  string     `json:"column,omitempty"`
	Message string     `json:"message"`
}

// String returns a human-readable representation of the change.
func (c Change) String() string {
	if c.Column != "" {
		return fmt.Sprintf("[%s] table=%s column=%s: %s", c.Type, c.Table, c.Column, c.Message)
	}
	return fmt.Sprintf("[%s] table=%s: %s", c.Type, c.Table, c.Message)
}

// Result holds the full set of detected drift between two schemas.
type Result struct {
	Changes []Change `json:"changes"`
}

// HasDrift returns true when at least one change was detected.
func (r *Result) HasDrift() bool {
	return len(r.Changes) > 0
}

// Add appends a change to the result.
func (r *Result) Add(c Change) {
	r.Changes = append(r.Changes, c)
}
