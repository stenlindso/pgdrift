package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/your-org/pgdrift/internal/diff"
)

func makeMaturityReport() *diff.MaturityReport {
	return &diff.MaturityReport{
		MinRuns: 3,
		AsOf:    time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Entries: []diff.MaturityEntry{
			{Schema: "public", Table: "orders", Level: diff.MaturityFluctuating, Changes: 4, Since: time.Now()},
			{Schema: "public", Table: "users", Level: diff.MaturityStable, Changes: 0, Since: time.Now()},
			{Schema: "audit", Table: "logs", Level: diff.MaturityNew, Changes: 1, Since: time.Now()},
		},
	}
}

func TestWriteMaturity_Nil(t *testing.T) {
	var buf bytes.Buffer
	WriteMaturity(&buf, nil)
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data' message, got: %s", buf.String())
	}
}

func TestWriteMaturity_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	WriteMaturity(&buf, makeMaturityReport())
	out := buf.String()

	for _, want := range []string{"orders", "fluctuating", "users", "stable", "logs", "new", "min_runs=3"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestWriteMaturityJSON_Nil(t *testing.T) {
	var buf bytes.Buffer
	err := WriteMaturityJSON(&buf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "{") {
		t.Errorf("expected JSON object, got: %s", buf.String())
	}
}

func TestWriteMaturityJSON_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	err := WriteMaturityJSON(&buf, makeMaturityReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"\"schema\"", "\"level\"", "\"changes\""} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output:\n%s", want, out)
		}
	}
}
