package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/pgdrift/internal/diff"
)

func makeMatches() []diff.RuleMatch {
	return []diff.RuleMatch{
		{
			Change: diff.Change{Kind: diff.ChangeKindTableAdded, Schema: "public", Table: "audit"},
			Rule:   diff.Rule{Action: diff.RuleActionWarn, Note: "review new tables"},
		},
		{
			Change: diff.Change{Kind: diff.ChangeKindColumnTypeChanged, Schema: "public", Table: "users", Object: "email"},
			Rule:   diff.Rule{Action: diff.RuleActionError},
		},
	}
}

func TestWriteRuleMatches_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteRuleMatches(&buf, nil)
	if !strings.Contains(buf.String(), "no rules matched") {
		t.Errorf("expected 'no rules matched', got: %s", buf.String())
	}
}

func TestWriteRuleMatches_WithMatches(t *testing.T) {
	var buf bytes.Buffer
	matches := makeMatches()
	WriteRuleMatches(&buf, matches)
	out := buf.String()

	if !strings.Contains(out, "2 rule(s) matched") {
		t.Errorf("expected count line, got: %s", out)
	}
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", out)
	}
	if !strings.Contains(out, "review new tables") {
		t.Errorf("expected note in output, got: %s", out)
	}
}

func TestWriteRuleMatchesJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := WriteRuleMatchesJSON(&buf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty array, got %d items", len(out))
	}
}

func TestWriteRuleMatchesJSON_WithMatches(t *testing.T) {
	var buf bytes.Buffer
	matches := makeMatches()
	err := WriteRuleMatchesJSON(&buf, matches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 items, got %d", len(out))
	}
	if out[0]["action"] != "warn" {
		t.Errorf("expected warn action, got %v", out[0]["action"])
	}
	if out[0]["note"] != "review new tables" {
		t.Errorf("expected note, got %v", out[0]["note"])
	}
	if out[1]["object"] != "email" {
		t.Errorf("expected object=email, got %v", out[1]["object"])
	}
}
