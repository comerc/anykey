//go:build goexperiment.synctest
// +build goexperiment.synctest

package usecase

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestRunWithPoolAndRateLimit(t *testing.T) {
	synctest.Run(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		start := time.Now()
		RunWithPoolAndRateLimit(ctx, 5, 3, 200*time.Millisecond, func(int) {}, nil)
		elapsed := time.Since(start)

		wantMin := 4 * 200 * time.Millisecond
		if elapsed < wantMin {
			t.Fatalf("elapsed too small with synctest: %v < %v", elapsed, wantMin)
		}
	})
}
