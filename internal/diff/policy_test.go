package diff

import (
	"testing"
)

func buildPolicyResult(kinds ...ChangeKind) *Result {
	r := &Result{}
	for _, k := range kinds {
		r.Changes = append(r.Changes, Change{
			Kind:   k,
			Schema: "public",
			Table:  "users",
		})
	}
	return r
}

func TestNewPolicy_DefaultAction(t *testing.T) {
	p := NewPolicy(nil)
	if got := p.ActionFor(KindColumnAdded); got != PolicyWarn {
		t.Errorf("expected warn, got %s", got)
	}
}

func TestNewPolicy_NilPolicy(t *testing.T) {
	var p *Policy
	if got := p.ActionFor(KindColumnAdded); got != PolicyWarn {
		t.Errorf("expected warn for nil policy, got %s", got)
	}
}

func TestNewPolicy_CustomRule(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Kind: KindColumnTypeChanged, Action: PolicyFail},
		{Kind: KindColumnAdded, Action: PolicyIgnore},
	})

	if got := p.ActionFor(KindColumnTypeChanged); got != PolicyFail {
		t.Errorf("expected fail, got %s", got)
	}
	if got := p.ActionFor(KindColumnAdded); got != PolicyIgnore {
		t.Errorf("expected ignore, got %s", got)
	}
	if got := p.ActionFor(KindTableAdded); got != PolicyWarn {
		t.Errorf("expected warn (default), got %s", got)
	}
}

func TestApplyPolicy_NilResult(t *testing.T) {
	p := NewPolicy(nil)
	hasFail, summary := ApplyPolicy(nil, p)
	if hasFail {
		t.Error("expected no fail for nil result")
	}
	if summary != "no policy applied" {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestApplyPolicy_NoChanges(t *testing.T) {
	p := NewPolicy(nil)
	result := buildPolicyResult()
	hasFail, summary := ApplyPolicy(result, p)
	if hasFail {
		t.Error("expected no fail")
	}
	if summary != "no changes" {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestApplyPolicy_FailOnTypeChange(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Kind: KindColumnTypeChanged, Action: PolicyFail},
	})
	result := buildPolicyResult(KindColumnTypeChanged, KindColumnAdded)
	hasFail, summary := ApplyPolicy(result, p)
	if !hasFail {
		t.Error("expected hasFail=true")
	}
	if summary != "fail:1 warn:1" {
		t.Errorf("unexpected summary: %q", summary)
	}
}

func TestApplyPolicy_AnnotatesChanges(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Kind: KindTableRemoved, Action: PolicyFail},
	})
	result := buildPolicyResult(KindTableRemoved)
	ApplyPolicy(result, p)

	val, ok := GetAnnotation(result.Changes[0], "policy")
	if !ok {
		t.Fatal("expected policy annotation")
	}
	if val != "fail" {
		t.Errorf("expected annotation value 'fail', got %q", val)
	}
}
