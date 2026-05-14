package diff

import "time"

// TrendPoint represents a single drift measurement at a point in time.
type TrendPoint struct {
	Timestamp   time.Time        `json:"timestamp"`
	TotalChanges int             `json:"total_changes"`
	BySeverity  map[string]int   `json:"by_severity"`
	ByKind      map[string]int   `json:"by_kind"`
}

// Trend holds an ordered series of drift measurements.
type Trend struct {
	Points []TrendPoint `json:"points"`
}

// NewTrend creates an empty Trend.
func NewTrend() *Trend {
	return &Trend{Points: []TrendPoint{}}
}

// Record appends a new TrendPoint derived from the given Result.
// If result is nil, a zero-value point is appended.
func (t *Trend) Record(result *Result) {
	pt := TrendPoint{
		Timestamp:    time.Now().UTC(),
		TotalChanges: 0,
		BySeverity:   map[string]int{},
		ByKind:       map[string]int{},
	}
	if result != nil {
		pt.TotalChanges = len(result.Changes)
		for _, c := range result.Changes {
			sev := Severity(c.Kind).String()
			pt.BySeverity[sev]++
			pt.ByKind[string(c.Kind)]++
		}
	}
	t.Points = append(t.Points, pt)
}

// Delta returns the change in total drift count between the last two points.
// Returns 0 if fewer than two points exist.
func (t *Trend) Delta() int {
	if len(t.Points) < 2 {
		return 0
	}
	last := t.Points[len(t.Points)-1]
	prev := t.Points[len(t.Points)-2]
	return last.TotalChanges - prev.TotalChanges
}

// Latest returns the most recent TrendPoint, or nil if empty.
func (t *Trend) Latest() *TrendPoint {
	if len(t.Points) == 0 {
		return nil
	}
	p := t.Points[len(t.Points)-1]
	return &p
}
