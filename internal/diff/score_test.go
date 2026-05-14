package diff

import "testing"

func buildScoreResult(kinds ...ChangeKind) *Result {
	r := &Result{}
	for _, k := range kinds {
		r.Changes = append(r.Changes, Change{
			Kind:   k,
			Schema: "public",
			Table:  "t",
		})
	}
	return r
}

func TestScore_NilResult(t *testing.T) {
	s := Score(nil)
	if s.Value != 100 || s.Grade != "A" {
		t.Fatalf("expected 100/A, got %d/%s", s.Value, s.Grade)
	}
}

func TestScore_NoDrift(t *testing.T) {
	s := Score(&Result{})
	if s.Value != 100 || s.Grade != "A" {
		t.Fatalf("expected 100/A, got %d/%s", s.Value, s.Grade)
	}
}

func TestScore_SingleHighChange(t *testing.T) {
	s := Score(buildScoreResult(KindTableRemoved))
	if s.Value != 90 {
		t.Fatalf("expected 90, got %d", s.Value)
	}
	if s.Grade != "A" {
		t.Fatalf("expected grade A, got %s", s.Grade)
	}
}

func TestScore_MultipleHighChanges(t *testing.T) {
	s := Score(buildScoreResult(
		KindTableRemoved, KindTableRemoved, KindTableRemoved,
		KindTableRemoved, KindTableRemoved,
	))
	// 5 * 10 = 50 deducted → 50
	if s.Value != 50 {
		t.Fatalf("expected 50, got %d", s.Value)
	}
	if s.Grade != "D" {
		t.Fatalf("expected grade D, got %s", s.Grade)
	}
}

func TestScore_FloorAtZero(t *testing.T) {
	kinds := make([]ChangeKind, 20)
	for i := range kinds {
		kinds[i] = KindTableRemoved // 20 * 10 = 200 deducted
	}
	s := Score(buildScoreResult(kinds...))
	if s.Value != 0 {
		t.Fatalf("expected 0, got %d", s.Value)
	}
	if s.Grade != "F" {
		t.Fatalf("expected grade F, got %s", s.Grade)
	}
}

func TestScore_MixedSeverities(t *testing.T) {
	// 1 high (10) + 1 medium (5) + 1 low (2) = 17 → 83
	s := Score(buildScoreResult(KindTableRemoved, KindColumnTypeChanged, KindColumnAdded))
	if s.Value != 83 {
		t.Fatalf("expected 83, got %d", s.Value)
	}
	if s.Grade != "B" {
		t.Fatalf("expected grade B, got %s", s.Grade)
	}
}

func TestScore_SummaryNonEmpty(t *testing.T) {
	s := Score(buildScoreResult(KindTableRemoved))
	if s.Summary == "" {
		t.Fatal("expected non-empty summary")
	}
}
