package relay

import (
	"context"
	"sync/atomic"
	"time"
)

func WatchHold(ctx context.Context, d time.Duration, fired *atomic.Bool) error {
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		fired.Store(true)
		return nil
	}
}
