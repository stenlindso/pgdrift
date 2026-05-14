package diff

import (
	"fmt"
	"sort"
	"time"
)

// ChangelogEntry represents a single recorded diff event with metadata.
type ChangelogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Label     string    `json:"label,omitempty"`
	Summary   DriftSummary `json:"summary"`
	Score     ScoreResult  `json:"score"`
}

// Changelog holds an ordered list of diff entries over time.
type Changelog struct {
	Entries []ChangelogEntry `json:"entries"`
}

// NewChangelog returns an empty Changelog.
func NewChangelog() *Changelog {
	return &Changelog{}
}

// Record appends a new entry derived from the given Result.
// label is an optional human-readable tag (e.g. a deploy SHA or env name).
func (c *Changelog) Record(r *Result, label string) {
	if r == nil {
		return
	}
	entry := ChangelogEntry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Summary:   Summarize(r),
		Score:     Score(r),
	}
	c.Entries = append(c.Entries, entry)
}

// Len returns the number of recorded entries.
func (c *Changelog) Len() int {
	if c == nil {
		return 0
	}
	return len(c.Entries)
}

// Latest returns the most recent entry, or nil if the changelog is empty.
func (c *Changelog) Latest() *ChangelogEntry {
	if c == nil || len(c.Entries) == 0 {
		return nil
	}
	return &c.Entries[len(c.Entries)-1]
}

// TopChanged returns up to n table names that appear most frequently across
// all recorded entries.
func (c *Changelog) TopChanged(n int) []string {
	if c == nil {
		return nil
	}
	counts := map[string]int{}
	for _, e := range c.Entries {
		for _, t := range e.Summary.AffectedTables {
			counts[t]++
		}
	}
	type kv struct {
		Key   string
		Count int
	}
	var pairs []kv
	for k, v := range counts {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Key < pairs[j].Key
	})
	result := make([]string, 0, n)
	for i, p := range pairs {
		if i >= n {
			break
		}
		result = append(result, fmt.Sprintf("%s (%d)", p.Key, p.Count))
	}
	return result
}
