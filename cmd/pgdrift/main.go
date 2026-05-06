package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/pgdrift/internal/diff"
	"github.com/example/pgdrift/internal/report"
	"github.com/example/pgdrift/internal/schema"
)

func main() {
	sourceDSN := flag.String("source", "", "DSN for the source (baseline) database")
	targetDSN := flag.String("target", "", "DSN for the target database")
	format := flag.String("format", "text", "Output format: text or json")
	flag.Parse()

	if *sourceDSN == "" || *targetDSN == "" {
		fmt.Fprintln(os.Stderr, "error: --source and --target are required")
		flag.Usage()
		os.Exit(1)
	}

	source, err := schema.Load(*sourceDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading source schema: %v\n", err)
		os.Exit(1)
	}

	target, err := schema.Load(*targetDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading target schema: %v\n", err)
		os.Exit(1)
	}

	result := diff.Compare(source, target)

	w := report.NewWriter(os.Stdout)
	switch *format {
	case "json":
		err = w.WriteJSON(result)
	default:
		err = w.WriteText(result)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	if result.HasDrift() {
		os.Exit(2)
	}
}
