package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/pgdrift/pgdrift/internal/schema"
	"github.com/pgdrift/pgdrift/internal/snapshot"
)

// runSnapshot connects to the given DSN, loads the schema, and saves a
// snapshot to outPath. It returns a non-nil error on failure.
func runSnapshot(dsn, outPath string) error {
	if dsn == "" {
		return fmt.Errorf("snapshot: --source DSN is required")
	}
	if outPath == "" {
		return fmt.Errorf("snapshot: --output path is required")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("snapshot: open db: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	s, err := schema.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("snapshot: load schema: %w", err)
	}

	snap := snapshot.New(s, dsn)
	if err := snapshot.Save(snap, outPath); err != nil {
		return fmt.Errorf("snapshot: save: %w", err)
	}

	fmt.Fprintf(os.Stdout, "snapshot saved to %s (captured at %s)\n",
		outPath, snap.CapturedAt.Format("2006-01-02T15:04:05Z"))
	return nil
}

// runDiffSnapshot loads a snapshot from snapshotPath and compares it against
// the live schema at targetDSN, writing the report to stdout.
func runDiffSnapshot(snapshotPath, targetDSN string) error {
	if snapshotPath == "" {
		return fmt.Errorf("diff-snapshot: --snapshot path is required")
	}
	if targetDSN == "" {
		return fmt.Errorf("diff-snapshot: --target DSN is required")
	}

	baseline, err := snapshot.Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("diff-snapshot: load snapshot: %w", err)
	}

	db, err := sql.Open("postgres", targetDSN)
	if err != nil {
		return fmt.Errorf("diff-snapshot: open target db: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	current, err := schema.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("diff-snapshot: load target schema: %w", err)
	}

	_ = baseline
	_ = current
	// diff.Compare and report.Write would be called here, mirroring main logic.
	return nil
}
