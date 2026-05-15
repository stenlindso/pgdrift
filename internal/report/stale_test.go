package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/your-org/pgdrift/internal/diff"
)

func makeStaleReport(level diff.StalenessLevel, age time.Duration) *diff.StalenessReport {
	return &diff.StalenessReport{
		CapturedAt: time.Now().Add(-age),
		Age:        age,
		Level:      level,
		Message:    level.String() + " message",
	}
}

func TestWriteStale_Nil(t *testing.T) {
	var buf bytes.Buffer
	WriteStale(&buf, nil)
	if !strings.Contains(buf.String(), "no report") {
		t.Errorf("expected 'no report', got %q", buf.String())
	}
}

func TestWriteStale_Fresh(t *testing.T) {
	var buf bytes.Buffer
	WriteStale(&buf, makeStaleReport(diff.StalenessNone, 2*time.Hour))
	out := buf.String()
	if !strings.Contains(out, "fresh") {
		t.Errorf("expected 'fresh' in output, got %q", out)
	}
	if !strings.Contains(out, "captured") {
		t.Errorf("expected 'captured' in output, got %q", out)
	}
}

func TestWriteStale_Critical(t *testing.T) {
	var buf bytes.Buffer
	WriteStale(&buf, makeStaleReport(diff.StalenessCritical, 100*time.Hour))
	if !strings.Contains(buf.String(), "critical") {
		t.Errorf("expected 'critical' in output")
	}
}

func TestWriteStaleJSON_Nil(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteStaleJSON(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["level"] != "unknown" {
		t.Errorf("expected level=unknown, got %v", m["level"])
	}
}

func TestWriteStaleJSON_WithData(t *testing.T) {
	var buf bytes.Buffer
	r := makeStaleReport(diff.StalenessWarning, 30*time.Hour)
	if err := WriteStaleJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["level"] != "warning" {
		t.Errorf("expected level=warning, got %v", m["level"])
	}
	if _, ok := m["captured_at"]; !ok {
		t.Error("expected captured_at field")
	}
	if _, ok := m["age_seconds"]; !ok {
		t.Error("expected age_seconds field")
	}
}
