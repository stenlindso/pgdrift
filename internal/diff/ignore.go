package diff

// IgnoreRule defines a rule for suppressing specific drift changes.
type IgnoreRule struct {
	Schema string
	Table  string
	Column string
	Kind   ChangeKind
}

// IgnoreList holds a collection of rules used to suppress changes from a Result.
type IgnoreList struct {
	rules []IgnoreRule
}

// NewIgnoreList creates an IgnoreList from the provided rules.
func NewIgnoreList(rules []IgnoreRule) *IgnoreList {
	return &IgnoreList{rules: rules}
}

// Matches reports whether the given Change is suppressed by any rule in the list.
func (il *IgnoreList) Matches(c Change) bool {
	for _, r := range il.rules {
		if matchField(r.Schema, c.Schema) &&
			matchField(r.Table, c.Table) &&
			matchField(r.Column, c.Column) &&
			(r.Kind == "" || r.Kind == c.Kind) {
			return true
		}
	}
	return false
}

// Apply returns a new Result with all changes that match the IgnoreList removed.
func (il *IgnoreList) Apply(r *Result) *Result {
	if il == nil || r == nil {
		return r
	}
	out := &Result{}
	for _, c := range r.Changes {
		if !il.Matches(c) {
			out.Changes = append(out.Changes, c)
		}
	}
	return out
}

// matchField returns true when the rule field is empty (wildcard) or equals the value.
func matchField(ruleField, value string) bool {
	return ruleField == "" || ruleField == value
}
