package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/your-org/pgdrift/internal/diff"
)

func makeDriftRate(total, driftRuns, changes int) *diff.DriftRate {
	now := time.Now().UTC()
	return &diff.DriftRate{
		WindowStart: now.Add(-24 * time.Hour),
		WindowEnd:   now,
		TotalRuns:   total,
		DriftRuns:   driftRuns,
		ChangeCount: changes,
	}
}

func TestWriteDriftRate_Nil(t *testing.T) {
	var buf bytes.Buffer
	WriteDriftRate(&buf, nil)
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data', got: %s", buf.String())
	}
}

func TestWriteDriftRate_WithData(t *testing.T) {
	var buf bytes.Buffer
	dr := makeDriftRate(10, 4, 7)
	WriteDriftRate(&buf, dr)
	out := buf.String()

	for _, want := range []string{"Drift Rate", "Total Runs", "Drift Runs", "40.0%", "moderate"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestWriteDriftRateJSON_Nil(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteDriftRateJSON(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["label"] != "stable" {
		t.Errorf("expected label 'stable', got %v", m["label"])
	}
}

func TestWriteDriftRateJSON_WithData(t *testing.T) {
	var buf bytes.Buffer
	dr := makeDriftRate(8, 7, 15)
	if err := WriteDriftRateJSON(&buf, dr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["total_runs"].(float64) != 8 {
		t.Errorf("expected total_runs=8, got %v", m["total_runs"])
	}
	if m["drift_runs"].(float64) != 7 {
		t.Errorf("expected drift_runs=7, got %v", m["drift_runs"])
	}
	if m["label"] != "frequent" {
		t.Errorf("expected label 'frequent', got %v", m["label"])
	}
}
