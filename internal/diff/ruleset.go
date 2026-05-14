package diff

import (
	"fmt"
	"strings"
)

// RuleAction defines what happens when a rule matches.
type RuleAction string

const (
	RuleActionWarn  RuleAction = "warn"
	RuleActionError RuleAction = "error"
	RuleActionSkip  RuleAction = "skip"
)

// Rule defines a single matching rule applied to a Change.
type Rule struct {
	Kind   ChangeKind
	Table  string // empty means any table
	Action RuleAction
	Note   string
}

// Ruleset holds an ordered list of Rules evaluated against diff results.
type Ruleset struct {
	rules []Rule
}

// NewRuleset creates an empty Ruleset.
func NewRuleset() *Ruleset {
	return &Ruleset{}
}

// Add appends a Rule to the Ruleset.
func (rs *Ruleset) Add(r Rule) {
	rs.rules = append(rs.rules, r)
}

// RuleMatch describes a rule that matched a Change.
type RuleMatch struct {
	Change Change
	Rule   Rule
}

// String returns a human-readable description of the match.
func (m RuleMatch) String() string {
	base := fmt.Sprintf("[%s] %s", strings.ToUpper(string(m.Rule.Action)), m.Change.String())
	if m.Rule.Note != "" {
		return base + " — " + m.Rule.Note
	}
	return base
}

// Evaluate applies all rules to the result and returns matched rules.
// Changes matched by a RuleActionSkip rule are excluded from the returned result.
func (rs *Ruleset) Evaluate(result *Result) ([]RuleMatch, *Result) {
	if rs == nil || result == nil {
		return nil, result
	}

	var matches []RuleMatch
	var kept []Change

	for _, ch := range result.Changes {
		matched := false
		for _, rule := range rs.rules {
			if rule.Kind != ch.Kind {
				continue
			}
			if rule.Table != "" && !strings.EqualFold(rule.Table, ch.Table) {
				continue
			}
			matches = append(matches, RuleMatch{Change: ch, Rule: rule})
			if rule.Action == RuleActionSkip {
				matched = true
			}
			break
		}
		if !matched {
			kept = append(kept, ch)
		}
	}

	filtered := &Result{Changes: kept}
	return matches, filtered
}
