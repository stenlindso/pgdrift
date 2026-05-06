package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pgdrift/pgdrift/internal/diff"
	"github.com/pgdrift/pgdrift/internal/report"
)

func makeResult(changes ...diff.Change) *diff.Result {
	r := &diff.Result{}
	for _, c := range changes {
		r.Add(c)
	}
	return r
}

func TestWriteText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatText)
	if err := w.Write(makeResult()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No schema drift") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestWriteText_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatText)
	result := makeResult(diff.Change{
		Object:     "table:users",
		ChangeType: diff.ChangeAdded,
		Detail:     "table added",
	})
	if err := w.Write(result); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "added") || !strings.Contains(out, "table:users") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestWriteJSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatJSON)
	if err := w.Write(makeResult()); err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["drift_detected"] != false {
		t.Errorf("expected drift_detected=false")
	}
}

func TestWriteJSON_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatJSON)
	result := makeResult(diff.Change{
		Object:     "column:users.email",
		ChangeType: diff.ChangeAltered,
		Detail:     "nullable changed",
	})
	if err := w.Write(result); err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["drift_detected"] != true {
		t.Errorf("expected drift_detected=true")
	}
	changes, ok := out["changes"].([]interface{})
	if !ok || len(changes) != 1 {
		t.Errorf("expected 1 change, got %v", out["changes"])
	}
}
