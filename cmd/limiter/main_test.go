package main

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter_StrictSpacing(t *testing.T) {
	rl := &RateLimiter{minSpacing: 200 * time.Millisecond}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// First start should be immediate
	start1, delta1, err := rl.Wait(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if delta1 != 0 {
		t.Fatalf("first delta should be 0, got %v", delta1)
	}

	// Second start should be >= 200ms after first
	start2, delta2, err := rl.Wait(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if delta2 < 200*time.Millisecond {
		t.Fatalf("delta should be >= 200ms, got %v", delta2)
	}
	if start2.Sub(start1) < 200*time.Millisecond {
		t.Fatalf("starts should be >= 200ms apart, got %v", start2.Sub(start1))
	}

	// Third start should again respect spacing from second
	start3, delta3, err := rl.Wait(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if delta3 < 200*time.Millisecond {
		t.Fatalf("delta should be >= 200ms, got %v", delta3)
	}
	if start3.Sub(start2) < 200*time.Millisecond {
		t.Fatalf("starts should be >= 200ms apart, got %v", start3.Sub(start2))
	}
}
