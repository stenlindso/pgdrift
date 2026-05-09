package diff

import "testing"

func baseChange() Change {
	return Change{
		Schema:   "public",
		Table:    "users",
		Field:    "email",
		Kind:     ChangeKindColumnTypeChanged,
		OldValue: "text",
		NewValue: "varchar(255)",
	}
}

func TestAnnotate_AddsAnnotation(t *testing.T) {
	c := Annotate(baseChange(), "reason", "migration")
	if len(c.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(c.Annotations))
	}
	if c.Annotations[0].Key != "reason" || c.Annotations[0].Value != "migration" {
		t.Errorf("unexpected annotation: %+v", c.Annotations[0])
	}
}

func TestAnnotate_DoesNotMutateOriginal(t *testing.T) {
	orig := baseChange()
	_ = Annotate(orig, "k", "v")
	if len(orig.Annotations) != 0 {
		t.Error("original change was mutated")
	}
}

func TestGetAnnotation_Found(t *testing.T) {
	c := Annotate(baseChange(), "env", "prod")
	v, ok := GetAnnotation(c, "env")
	if !ok || v != "prod" {
		t.Errorf("expected (prod, true), got (%s, %v)", v, ok)
	}
}

func TestGetAnnotation_NotFound(t *testing.T) {
	c := baseChange()
	_, ok := GetAnnotation(c, "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestAnnotateResult_NilResult(t *testing.T) {
	out := AnnotateResult(nil, "k", "v", func(Change) bool { return true })
	if out != nil {
		t.Error("expected nil result")
	}
}

func TestAnnotateResult_MatchingChanges(t *testing.T) {
	r := &Result{
		Changes: []Change{
			baseChange(),
			{Schema: "public", Table: "orders", Kind: ChangeKindTableAdded},
		},
	}
	out := AnnotateResult(r, "tag", "flagged", func(c Change) bool {
		return c.Kind == ChangeKindColumnTypeChanged
	})
	if len(out.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(out.Changes))
	}
	_, ok0 := GetAnnotation(out.Changes[0], "tag")
	_, ok1 := GetAnnotation(out.Changes[1], "tag")
	if !ok0 {
		t.Error("expected first change to be annotated")
	}
	if ok1 {
		t.Error("expected second change not to be annotated")
	}
}
