package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/pgdrift/internal/diff"
)

func annotatedResult() *diff.Result {
	c1 := diff.Change{
		Schema:   "public",
		Table:    "users",
		Field:    "email",
		Kind:     diff.ChangeKindColumnTypeChanged,
		OldValue: "text",
		NewValue: "varchar",
	}
	c1 = diff.Annotate(c1, "reason", "intentional")
	c1 = diff.Annotate(c1, "ticket", "PROJ-42")

	c2 := diff.Change{
		Schema: "public",
		Table:  "orders",
		Kind:   diff.ChangeKindTableAdded,
	}
	return &diff.Result{Changes: []diff.Change{c1, c2}}
}

func TestWriteAnnotations_NilResult(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteAnnotations(&buf, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No annotated") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestWriteAnnotations_WithAnnotations(t *testing.T) {
	var buf bytes.Buffer
	r := annotatedResult()
	if err := WriteAnnotations(&buf, r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Annotated changes:") {
		t.Error("missing header")
	}
	if !strings.Contains(out, "reason") || !strings.Contains(out, "intentional") {
		t.Error("missing reason annotation")
	}
	if !strings.Contains(out, "PROJ-42") {
		t.Error("missing ticket annotation")
	}
	// unannotated change should not appear
	if strings.Contains(out, "orders") {
		t.Error("unannotated change should not appear")
	}
}

func TestAnnotationSummary_Empty(t *testing.T) {
	m := AnnotationSummary(nil)
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestAnnotationSummary_Counts(t *testing.T) {
	r := annotatedResult()
	m := AnnotationSummary(r)
	if m["reason"] != 1 {
		t.Errorf("expected reason=1, got %d", m["reason"])
	}
	if m["ticket"] != 1 {
		t.Errorf("expected ticket=1, got %d", m["ticket"])
	}
	if _, ok := m["missing"]; ok {
		t.Error("unexpected key in summary")
	}
}
