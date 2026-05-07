// Package snapshot handles persistence of PostgreSQL schema snapshots.
//
// A snapshot captures the state of a database schema at a point in time and
// serialises it to JSON on disk. Saved snapshots can later be loaded and
// compared against a live database (or another snapshot) using the diff
// package to detect schema drift without requiring two live database
// connections simultaneously.
//
// Typical usage:
//
//	// Capture and save
//	s, _ := schemaloader.Load(ctx, db)
//	snap := snapshot.New(s, dsn)
//	snapshot.Save(snap, "baseline.json")
//
//	// Later: load and compare
//	baseline, _ := snapshot.Load("baseline.json")
//	current, _ := schemaloader.Load(ctx, db)
//	result := diff.Compare(baseline.Schema, current)
package snapshot
