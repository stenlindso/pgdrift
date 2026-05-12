package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved diff result used as a reference point
// for suppressing known drift in future comparisons.
type Baseline struct {
	CreatedAt time.Time        `json:"created_at"`
	Changes   []BaselineChange `json:"changes"`
}

// BaselineChange is a compact representation of a Change stored in a baseline.
type BaselineChange struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
	Kind   string `json:"kind"`
	Field  string `json:"field,omitempty"`
}

// NewBaseline creates a Baseline from an existing Result.
func NewBaseline(r *Result) *Baseline {
	if r == nil {
		return &Baseline{CreatedAt: time.Now()}
	}
	changes := make([]BaselineChange, 0, len(r.Changes))
	for _, c := range r.Changes {
		changes = append(changes, BaselineChange{
			Schema: c.Schema,
			Table:  c.Table,
			Kind:   string(c.Kind),
			Field:  c.Field,
		})
	}
	return &Baseline{CreatedAt: time.Now(), Changes: changes}
}

// SaveBaseline writes a Baseline to a JSON file at the given path.
func SaveBaseline(b *Baseline, path string) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// LoadBaseline reads a Baseline from a JSON file at the given path.
func LoadBaseline(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// ApplyBaseline removes changes from r that are already present in the baseline.
func ApplyBaseline(r *Result, b *Baseline) *Result {
	if r == nil || b == nil {
		return r
	}
	known := make(map[BaselineChange]struct{}, len(b.Changes))
	for _, bc := range b.Changes {
		known[bc] = struct{}{}
	}
	filtered := make([]Change, 0, len(r.Changes))
	for _, c := range r.Changes {
		key := BaselineChange{Schema: c.Schema, Table: c.Table, Kind: string(c.Kind), Field: c.Field}
		if _, suppressed := known[key]; !suppressed {
			filtered = append(filtered, c)
		}
	}
	return &Result{Changes: filtered}
}
