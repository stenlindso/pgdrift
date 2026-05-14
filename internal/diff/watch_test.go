package diff

import (
	"context"
	"errors"
	"testing"
	"time"
)

func fixedLoader(result *Result, err error) SchemaLoader {
	return func() (*Result, error) {
		return result, err
	}
}

func TestWatch_EmitsImmediately(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := WatchOptions{Interval: 500 * time.Millisecond, MaxRuns: 1}
	result := &Result{Changes: []Change{}}

	ch := Watch(ctx, opts, fixedLoader(result, nil))

	event, ok := <-ch
	if !ok {
		t.Fatal("expected at least one event")
	}
	if event.Err != nil {
		t.Fatalf("unexpected error: %v", event.Err)
	}
	if event.Result != result {
		t.Error("expected result to be forwarded")
	}
	if event.RunAt.IsZero() {
		t.Error("expected RunAt to be set")
	}
}

func TestWatch_RespectsMaxRuns(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opts := WatchOptions{Interval: 50 * time.Millisecond, MaxRuns: 3}

	calls := 0
	loader := func() (*Result, error) {
		calls++
		return &Result{}, nil
	}

	ch := Watch(ctx, opts, loader)
	for range ch {
	}

	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestWatch_PropagatesError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sentinel := errors.New("db unavailable")
	opts := WatchOptions{Interval: 50 * time.Millisecond, MaxRuns: 1}

	ch := Watch(ctx, opts, fixedLoader(nil, sentinel))

	event := <-ch
	if !errors.Is(event.Err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", event.Err)
	}
}

func TestWatch_CancelStopsLoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	opts := WatchOptions{Interval: 20 * time.Millisecond}
	ch := Watch(ctx, opts, fixedLoader(&Result{}, nil))

	// Drain the first event then cancel.
	<-ch
	cancel()

	// Channel must close within a reasonable time.
	select {
	case _, ok := <-ch:
		if ok {
			// drain any buffered event
			for range ch {
			}
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("channel did not close after context cancellation")
	}
}
