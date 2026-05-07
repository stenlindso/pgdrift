// Package snapshot provides functionality to save and load schema snapshots
// to/from disk, enabling drift detection against a previously captured baseline.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pgdrift/pgdrift/internal/schema"
)

// Snapshot wraps a schema with metadata about when it was captured.
type Snapshot struct {
	CapturedAt time.Time     `json:"captured_at"`
	DSN        string        `json:"dsn,omitempty"`
	Schema     *schema.Schema `json:"schema"`
}

// New creates a new Snapshot from the given schema and DSN.
func New(s *schema.Schema, dsn string) *Snapshot {
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		DSN:        dsn,
		Schema:     s,
	}
}

// Save writes the snapshot as JSON to the given file path.
func Save(snap *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file %q: %w", path, err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	if snap.Schema == nil {
		return nil, fmt.Errorf("snapshot: file %q contains no schema", path)
	}
	return &snap, nil
}
