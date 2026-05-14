package diff

import (
	"testing"
)

func buildRulesetResult() *Result {
	return &Result{
		Changes: []Change{
			{Kind: ChangeKindColumnTypeChanged, Schema: "public", Table: "users", Object: "email"},
			{Kind: ChangeKindTableAdded, Schema: "public", Table: "audit_log"},
			{Kind: ChangeKindColumnAdded, Schema: "public", Table: "orders", Object: "total"},
		},
	}
}

func TestRuleset_NilRuleset(t *testing.T) {
	var rs *Ruleset
	result := buildRulesetResult()
	matches, out := rs.Evaluate(result)
	if len(matches) != 0 {
		t.Errorf("expected no matches, got %d", len(matches))
	}
	if out != result {
		t.Error("expected original result returned unchanged")
	}
}

func TestRuleset_NilResult(t *testing.T) {
	rs := NewRuleset()
	matches, out := rs.Evaluate(nil)
	if len(matches) != 0 {
		t.Errorf("expected no matches, got %d", len(matches))
	}
	if out != nil {
		t.Error("expected nil result returned")
	}
}

func TestRuleset_NoRules(t *testing.T) {
	rs := NewRuleset()
	result := buildRulesetResult()
	matches, out := rs.Evaluate(result)
	if len(matches) != 0 {
		t.Errorf("expected no matches, got %d", len(matches))
	}
	if len(out.Changes) != 3 {
		t.Errorf("expected 3 changes, got %d", len(out.Changes))
	}
}

func TestRuleset_WarnRule(t *testing.T) {
	rs := NewRuleset()
	rs.Add(Rule{Kind: ChangeKindTableAdded, Action: RuleActionWarn, Note: "review new tables"})

	result := buildRulesetResult()
	matches, out := rs.Evaluate(result)

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Rule.Action != RuleActionWarn {
		t.Errorf("expected warn action, got %s", matches[0].Rule.Action)
	}
	// warn does not skip, so change remains
	if len(out.Changes) != 3 {
		t.Errorf("expected 3 changes kept, got %d", len(out.Changes))
	}
}

func TestRuleset_SkipRule(t *testing.T) {
	rs := NewRuleset()
	rs.Add(Rule{Kind: ChangeKindColumnTypeChanged, Action: RuleActionSkip})

	result := buildRulesetResult()
	matches, out := rs.Evaluate(result)

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if len(out.Changes) != 2 {
		t.Errorf("expected 2 changes after skip, got %d", len(out.Changes))
	}
}

func TestRuleset_TableScopedRule(t *testing.T) {
	rs := NewRuleset()
	rs.Add(Rule{Kind: ChangeKindColumnAdded, Table: "orders", Action: RuleActionError, Note: "orders schema is locked"})

	result := buildRulesetResult()
	matches, out := rs.Evaluate(result)

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Change.Table != "orders" {
		t.Errorf("expected match on orders, got %s", matches[0].Change.Table)
	}
	if len(out.Changes) != 3 {
		t.Errorf("expected 3 changes (error does not skip), got %d", len(out.Changes))
	}
}

func TestRuleMatch_String(t *testing.T) {
	m := RuleMatch{
		Change: Change{Kind: ChangeKindTableAdded, Schema: "public", Table: "logs"},
		Rule:   Rule{Action: RuleActionWarn, Note: "check this"},
	}
	s := m.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
	if !containsStr(s, "WARN") {
		t.Errorf("expected WARN in output, got: %s", s)
	}
	if !containsStr(s, "check this") {
		t.Errorf("expected note in output, got: %s", s)
	}
}
