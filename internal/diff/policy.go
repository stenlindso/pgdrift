package diff

import (
	"fmt"
	"strings"
)

// PolicyAction defines what to do when a policy rule matches.
type PolicyAction string

const (
	PolicyWarn PolicyAction = "warn"
	PolicyFail PolicyAction = "fail"
	PolicyIgnore PolicyAction = "ignore"
)

// PolicyRule maps a change kind to an action.
type PolicyRule struct {
	Kind   ChangeKind
	Action PolicyAction
}

// Policy holds a set of rules that govern how changes are treated.
type Policy struct {
	rules map[ChangeKind]PolicyAction
}

// NewPolicy creates a Policy from a slice of rules.
// Later rules for the same kind override earlier ones.
func NewPolicy(rules []PolicyRule) *Policy {
	p := &Policy{rules: make(map[ChangeKind]PolicyAction, len(rules))}
	for _, r := range rules {
		p.rules[r.Kind] = r.Action
	}
	return p
}

// ActionFor returns the PolicyAction for a given ChangeKind.
// Defaults to PolicyWarn if no rule is defined.
func (p *Policy) ActionFor(kind ChangeKind) PolicyAction {
	if p == nil {
		return PolicyWarn
	}
	if action, ok := p.rules[kind]; ok {
		return action
	}
	return PolicyWarn
}

// ApplyPolicy annotates each change in the result with the policy action
// and returns whether any change triggered a PolicyFail action.
func ApplyPolicy(result *Result, policy *Policy) (hasFail bool, summary string) {
	if result == nil || policy == nil {
		return false, "no policy applied"
	}

	counts := map[PolicyAction]int{}
	for i, c := range result.Changes {
		action := policy.ActionFor(c.Kind)
		counts[action]++
		if action == PolicyFail {
			hasFail = true
		}
		annotated := Annotate(c, "policy", string(action))
		result.Changes[i] = annotated
	}

	parts := make([]string, 0, 3)
	for _, a := range []PolicyAction{PolicyFail, PolicyWarn, PolicyIgnore} {
		if n := counts[a]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s:%d", a, n))
		}
	}
	if len(parts) == 0 {
		return false, "no changes"
	}
	return hasFail, strings.Join(parts, " ")
}
