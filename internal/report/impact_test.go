package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/pgdrift/internal/diff"
)

func makeImpactReport(levels ...diff.ImpactLevel) *diff.ImpactReport {
	rep := &diff.ImpactReport{}
	overall := diff.ImpactNone
	for _, l := range levels {
		rep.Changes = append(rep.Changes, diff.ImpactedChange{
			Change: diff.Change{Kind: diff.KindColumnTypeChanged, Schema: "public", Table: "orders"},
			Impact: l,
			Reason: "test reason",
		})
		if l > overall {
			overall = l
		}
	}
	rep.Overall = overall
	return rep
}

func TestWriteImpact_NilReport(t *testing.T) {
	var buf bytes.Buffer
	WriteImpact(&buf, nil)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes' in output, got: %s", buf.String())
	}
}

func TestWriteImpact_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	WriteImpact(&buf, makeImpactReport(diff.ImpactHigh, diff.ImpactCritical))
	out := buf.String()
	if !strings.Contains(out, "overall impact: critical") {
		t.Errorf("expected overall impact critical, got: %s", out)
	}
	if !strings.Contains(out, "test reason") {
		t.Errorf("expected reason in output")
	}
}

func TestWriteImpactJSON_NilReport(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteImpactJSON(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"overall": "none"`) {
		t.Errorf("expected none overall in JSON, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), `"changes": []`) {
		t.Errorf("expected empty changes array")
	}
}

func TestWriteImpactJSON_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteImpactJSON(&buf, makeImpactReport(diff.ImpactCritical)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"overall": "critical"`) {
		t.Errorf("expected critical in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"impact": "critical"`) {
		t.Errorf("expected impact field in change")
	}
	if !strings.Contains(out, `"reason": "test reason"`) {
		t.Errorf("expected reason field in change")
	}
}
