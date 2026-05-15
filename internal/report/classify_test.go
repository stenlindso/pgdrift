package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pgdrift/pgdrift/internal/diff"
)

func makeClassReport(rc diff.RiskClass, changes int) *diff.ClassificationReport {
	r := &diff.Result{}
	for i := 0; i < changes; i++ {
		r.Changes = append(r.Changes, diff.Change{
			Kind:  diff.KindTableRemoved,
			Table: "t",
		})
	}
	return &diff.ClassificationReport{
		Result:    r,
		RiskClass: rc,
		Summary:   "test summary",
	}
}

func TestWriteClassification_Nil(t *testing.T) {
	var buf bytes.Buffer
	WriteClassification(&buf, nil)
	if !strings.Contains(buf.String(), "unavailable") {
		t.Errorf("expected 'unavailable' in output, got: %s", buf.String())
	}
}

func TestWriteClassification_WithReport(t *testing.T) {
	var buf bytes.Buffer
	cr := makeClassReport(diff.RiskClassCritical, 3)
	WriteClassification(&buf, cr)
	out := buf.String()
	if !strings.Contains(out, "critical") {
		t.Errorf("expected 'critical' in output, got: %s", out)
	}
	if !strings.Contains(out, "3") {
		t.Errorf("expected change count '3' in output, got: %s", out)
	}
}

func TestWriteClassificationJSON_Nil(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteClassificationJSON(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["risk_class"] != "unknown" {
		t.Errorf("expected risk_class=unknown, got %v", out["risk_class"])
	}
}

func TestWriteClassificationJSON_WithReport(t *testing.T) {
	var buf bytes.Buffer
	cr := makeClassReport(diff.RiskClassModerate, 2)
	if err := WriteClassificationJSON(&buf, cr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["risk_class"] != "moderate" {
		t.Errorf("expected moderate, got %v", out["risk_class"])
	}
	if int(out["total_changes"].(float64)) != 2 {
		t.Errorf("expected total_changes=2, got %v", out["total_changes"])
	}
}
