package diff_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/pgdrift/internal/diff"
)

func buildBaselineResult() *diff.Result {
	return &diff.Result{
		Changes: []diff.Change{
			{Schema: "public", Table: "users", Kind: diff.ColumnTypeChanged, Field: "email"},
			{Schema: "public", Table: "orders", Kind: diff.TableAdded},
		},
	}
}

func TestNewBaseline_NilResult(t *testing.T) {
	b := diff.NewBaseline(nil)
	if b == nil {
		t.Fatal("expected non-nil baseline")
	}
	if len(b.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(b.Changes))
	}
}

func TestNewBaseline_FromResult(t *testing.T) {
	r := buildBaselineResult()
	b := diff.NewBaseline(r)
	if len(b.Changes) != 2 {
		t.Fatalf("expected 2 baseline changes, got %d", len(b.Changes))
	}
	if b.Changes[0].Table != "users" || b.Changes[0].Field != "email" {
		t.Errorf("unexpected first change: %+v", b.Changes[0])
	}
}

func TestSaveAndLoadBaseline(t *testing.T) {
	r := buildBaselineResult()
	b := diff.NewBaseline(r)

	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	if err := diff.SaveBaseline(b, path); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	loaded, err := diff.LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if len(loaded.Changes) != len(b.Changes) {
		t.Errorf("expected %d changes, got %d", len(b.Changes), len(loaded.Changes))
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := diff.LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := diff.LoadBaseline(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApplyBaseline_SuppressesKnown(t *testing.T) {
	r := buildBaselineResult()
	b := diff.NewBaseline(r)

	// add a new change not in baseline
	r.Changes = append(r.Changes, diff.Change{
		Schema: "public", Table: "products", Kind: diff.TableRemoved,
	})

	filtered := diff.ApplyBaseline(r, b)
	if len(filtered.Changes) != 1 {
		t.Fatalf("expected 1 new change after baseline, got %d", len(filtered.Changes))
	}
	if filtered.Changes[0].Table != "products" {
		t.Errorf("expected products change, got %+v", filtered.Changes[0])
	}
}

func TestApplyBaseline_NilInputs(t *testing.T) {
	if got := diff.ApplyBaseline(nil, nil); got != nil {
		t.Errorf("expected nil result for nil inputs")
	}
	r := buildBaselineResult()
	if got := diff.ApplyBaseline(r, nil); got != r {
		t.Errorf("expected original result when baseline is nil")
	}
}
