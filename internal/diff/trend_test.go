package diff

import (
	"testing"
)

func buildTrendResult(kinds ...ChangeKind) *Result {
	r := &Result{}
	for _, k := range kinds {
		r.Changes = append(r.Changes, Change{Kind: k, Table: "t", Object: "c"})
	}
	return r
}

func TestNewTrend_Empty(t *testing.T) {
	tr := NewTrend()
	if tr == nil {
		t.Fatal("expected non-nil trend")
	}
	if len(tr.Points) != 0 {
		t.Fatalf("expected 0 points, got %d", len(tr.Points))
	}
}

func TestTrend_RecordNilResult(t *testing.T) {
	tr := NewTrend()
	tr.Record(nil)
	if len(tr.Points) != 1 {
		t.Fatalf("expected 1 point, got %d", len(tr.Points))
	}
	if tr.Points[0].TotalChanges != 0 {
		t.Errorf("expected 0 changes for nil result")
	}
}

func TestTrend_RecordWithChanges(t *testing.T) {
	tr := NewTrend()
	res := buildTrendResult(KindColumnAdded, KindColumnAdded, KindTableRemoved)
	tr.Record(res)
	pt := tr.Points[0]
	if pt.TotalChanges != 3 {
		t.Errorf("expected 3 total changes, got %d", pt.TotalChanges)
	}
	if pt.ByKind[string(KindColumnAdded)] != 2 {
		t.Errorf("expected 2 column_added, got %d", pt.ByKind[string(KindColumnAdded)])
	}
	if pt.ByKind[string(KindTableRemoved)] != 1 {
		t.Errorf("expected 1 table_removed")
	}
}

func TestTrend_Delta(t *testing.T) {
	tr := NewTrend()
	tr.Record(buildTrendResult(KindColumnAdded, KindColumnAdded))
	tr.Record(buildTrendResult(KindColumnAdded, KindColumnAdded, KindColumnAdded, KindColumnAdded, KindColumnAdded))
	if tr.Delta() != 3 {
		t.Errorf("expected delta 3, got %d", tr.Delta())
	}
}

func TestTrend_Delta_Shrink(t *testing.T) {
	tr := NewTrend()
	tr.Record(buildTrendResult(KindColumnAdded, KindColumnAdded, KindColumnAdded))
	tr.Record(buildTrendResult(KindColumnAdded))
	if tr.Delta() != -2 {
		t.Errorf("expected delta -2, got %d", tr.Delta())
	}
}

func TestTrend_Delta_TooFewPoints(t *testing.T) {
	tr := NewTrend()
	if tr.Delta() != 0 {
		t.Errorf("expected 0 delta for empty trend")
	}
	tr.Record(buildTrendResult(KindColumnAdded))
	if tr.Delta() != 0 {
		t.Errorf("expected 0 delta for single point")
	}
}

func TestTrend_Latest_Empty(t *testing.T) {
	tr := NewTrend()
	if tr.Latest() != nil {
		t.Error("expected nil for empty trend")
	}
}

func TestTrend_Latest_ReturnsLast(t *testing.T) {
	tr := NewTrend()
	tr.Record(buildTrendResult(KindColumnAdded))
	tr.Record(buildTrendResult(KindTableAdded, KindTableAdded))
	latest := tr.Latest()
	if latest == nil {
		t.Fatal("expected non-nil latest")
	}
	if latest.TotalChanges != 2 {
		t.Errorf("expected latest to have 2 changes, got %d", latest.TotalChanges)
	}
}
