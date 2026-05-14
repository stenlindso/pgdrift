package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/example/pgdrift/internal/diff"
)

// WriteRuleMatches writes a human-readable summary of ruleset evaluation to w.
func WriteRuleMatches(w io.Writer, matches []diff.RuleMatch) {
	if len(matches) == 0 {
		fmt.Fprintln(w, "ruleset: no rules matched")
		return
	}

	fmt.Fprintf(w, "ruleset: %d rule(s) matched\n", len(matches))
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, m := range matches {
		fmt.Fprintln(w, m.String())
	}
}

// ruleMatchJSON is the JSON representation of a RuleMatch.
type ruleMatchJSON struct {
	Table  string `json:"table"`
	Object string `json:"object,omitempty"`
	Kind   string `json:"kind"`
	Action string `json:"action"`
	Note   string `json:"note,omitempty"`
}

// WriteRuleMatchesJSON writes rule matches as a JSON array to w.
func WriteRuleMatchesJSON(w io.Writer, matches []diff.RuleMatch) error {
	out := make([]ruleMatchJSON, 0, len(matches))
	for _, m := range matches {
		out = append(out, ruleMatchJSON{
			Table:  m.Change.Table,
			Object: m.Change.Object,
			Kind:   string(m.Change.Kind),
			Action: string(m.Rule.Action),
			Note:   m.Rule.Note,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
