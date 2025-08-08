//go:build goexperiment.synctest
// +build goexperiment.synctest

package domain

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestRateLimiter_StrictSpacing(t *testing.T) {
	synctest.Run(func() {
		rl := &RateLimiter{MinSpacing: 200 * time.Millisecond}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		var prev time.Time
		for i := 0; i < 5; i++ {
			start, delta, err := rl.Wait(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if i == 0 && delta != 0 {
				t.Fatalf("first delta=%v, want 0", delta)
			}
			if i > 0 && start.Sub(prev) < 200*time.Millisecond {
				t.Fatalf("interval=%v < 200ms", start.Sub(prev))
			}
			prev = start
		}
	})
}
