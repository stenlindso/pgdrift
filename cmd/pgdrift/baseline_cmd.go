package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/pgdrift/internal/diff"
	"github.com/your-org/pgdrift/internal/schema"
)

// runSaveBaseline connects to source and target, computes the diff, and saves
// the result as a baseline file to suppress known drift in future runs.
func runSaveBaseline(args []string) error {
	fs := flag.NewFlagSet("save-baseline", flag.ContinueOnError)
	srcDSN := fs.String("source", "", "source database DSN")
	dstDSN := fs.String("target", "", "target database DSN")
	out := fs.String("out", "baseline.json", "output baseline file path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *srcDSN == "" || *dstDSN == "" {
		return fmt.Errorf("save-baseline: --source and --target are required")
	}

	src, err := schema.Load(*srcDSN)
	if err != nil {
		return fmt.Errorf("save-baseline: load source: %w", err)
	}
	dst, err := schema.Load(*dstDSN)
	if err != nil {
		return fmt.Errorf("save-baseline: load target: %w", err)
	}

	result := diff.Compare(src, dst)
	baseline := diff.NewBaseline(result)
	if err := diff.SaveBaseline(baseline, *out); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "baseline saved to %s (%d changes recorded)\n", *out, len(baseline.Changes))
	return nil
}

// runDiffWithBaseline runs a diff and filters out changes present in the
// baseline file before printing the report.
func runDiffWithBaseline(args []string) error {
	fs := flag.NewFlagSet("diff-baseline", flag.ContinueOnError)
	srcDSN := fs.String("source", "", "source database DSN")
	dstDSN := fs.String("target", "", "target database DSN")
	baselinePath := fs.String("baseline", "baseline.json", "baseline file path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *srcDSN == "" || *dstDSN == "" {
		return fmt.Errorf("diff-baseline: --source and --target are required")
	}

	src, err := schema.Load(*srcDSN)
	if err != nil {
		return fmt.Errorf("diff-baseline: load source: %w", err)
	}
	dst, err := schema.Load(*dstDSN)
	if err != nil {
		return fmt.Errorf("diff-baseline: load target: %w", err)
	}

	result := diff.Compare(src, dst)

	baseline, err := diff.LoadBaseline(*baselinePath)
	if err != nil {
		return fmt.Errorf("diff-baseline: load baseline: %w", err)
	}

	filtered := diff.ApplyBaseline(result, baseline)
	if !filtered.HasDrift() {
		fmt.Fprintln(os.Stdout, "no new drift detected beyond baseline")
		return nil
	}
	for _, c := range filtered.Changes {
		fmt.Fprintln(os.Stdout, c.String())
	}
	return nil
}
