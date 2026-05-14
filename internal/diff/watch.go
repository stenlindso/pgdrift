package diff

import (
	"context"
	"time"
)

// WatchOptions configures the polling behaviour of Watch.
type WatchOptions struct {
	// Interval is how often the comparison is re-run.
	Interval time.Duration
	// MaxRuns limits the number of poll cycles; 0 means unlimited.
	MaxRuns int
}

// WatchEvent is emitted each time a comparison cycle completes.
type WatchEvent struct {
	RunAt  time.Time
	Result *Result
	Err    error
}

// SchemaLoader is a function that returns a fresh Result on demand.
type SchemaLoader func() (*Result, error)

// Watch repeatedly calls loader at the configured interval, sending each
// WatchEvent to the returned channel. The channel is closed when ctx is
// cancelled or MaxRuns cycles have completed.
func Watch(ctx context.Context, opts WatchOptions, loader SchemaLoader) <-chan WatchEvent {
	if opts.Interval <= 0 {
		opts.Interval = 30 * time.Second
	}

	ch := make(chan WatchEvent, 1)

	go func() {
		defer close(ch)

		runs := 0
		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		// Run immediately before waiting for the first tick.
		emit := func() bool {
			result, err := loader()
			event := WatchEvent{RunAt: time.Now(), Result: result, Err: err}
			select {
			case ch <- event:
			case <-ctx.Done():
				return false
			}
			runs++
			return opts.MaxRuns == 0 || runs < opts.MaxRuns
		}

		if !emit() {
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if !emit() {
					return
				}
			}
		}
	}()

	return ch
}
