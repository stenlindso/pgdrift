package diff

import "sort"

// HotspotEntry represents a table and the number of drift changes recorded
// across all changelog entries.
type HotspotEntry struct {
	Schema  string
	Table   string
	Changes int
}

// HotspotReport holds the ranked list of tables with the most drift activity.
type HotspotReport struct {
	Entries []HotspotEntry
	Total   int // total change events considered
}

// DetectHotspots analyses a Changelog and returns the tables that have drifted
// most frequently, ranked in descending order of change count.
// minChanges can be used to filter out tables below a threshold (0 = include all).
func DetectHotspots(cl *Changelog, minChanges int) *HotspotReport {
	if cl == nil {
		return &HotspotReport{}
	}

	type key struct{ schema, table string }
	counts := make(map[key]int)
	total := 0

	for _, entry := range cl.Entries() {
		if entry.Result == nil {
			continue
		}
		for _, ch := range entry.Result.Changes {
			counts[key{ch.Schema, ch.Table}]++
			total++
		}
	}

	var entries []HotspotEntry
	for k, n := range counts {
		if n >= minChanges || minChanges == 0 {
			entries = append(entries, HotspotEntry{Schema: k.schema, Table: k.table, Changes: n})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Changes != entries[j].Changes {
			return entries[i].Changes > entries[j].Changes
		}
		if entries[i].Schema != entries[j].Schema {
			return entries[i].Schema < entries[j].Schema
		}
		return entries[i].Table < entries[j].Table
	})

	return &HotspotReport{Entries: entries, Total: total}
}
