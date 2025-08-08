package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Request emulates an external request that takes some time
func Request(id int) {
	// Simulate variable latency
	time.Sleep(time.Duration(50+rand.Intn(150)) * time.Millisecond)
	fmt.Printf("done request %d at %s\n", id, time.Now().Format(time.RFC3339Nano))
}

// runWithPoolAndRateLimit runs n requests using m workers and ensures no more than
// one Request starts within each 200ms window across the whole program.
type RateLimiter struct {
	mu         sync.Mutex
	lastStart  time.Time
	minSpacing time.Duration
}

// Wait blocks until запуск нового запроса не нарушит минимальный интервал.
// Возвращает фактическое время старта и дельту до предыдущего старта.
func (r *RateLimiter) Wait(ctx context.Context) (time.Time, time.Duration, error) {
	for {
		// Быстрый выход, если контекст отменён
		select {
		case <-ctx.Done():
			return time.Time{}, 0, ctx.Err()
		default:
		}

		r.mu.Lock()
		now := time.Now()
		target := r.lastStart.Add(r.minSpacing)
		// Если интервал соблюдён — фиксируем новый старт и выходим
		if r.lastStart.IsZero() || !now.Before(target) {
			var delta time.Duration
			if !r.lastStart.IsZero() {
				delta = now.Sub(r.lastStart)
			}
			r.lastStart = now
			r.mu.Unlock()
			return now, delta, nil
		}
		// Иначе — нужно подождать до target
		wait := target.Sub(now)
		r.mu.Unlock()

		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return time.Time{}, 0, ctx.Err()
		case <-timer.C:
			// После ожидания повторим цикл и перепроверим условия под мьютексом
		}
	}
}

func runWithPoolAndRateLimit(ctx context.Context, n int, m int, minInterval time.Duration) {
	rl := &RateLimiter{minSpacing: minInterval}

	jobs := make(chan int)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for id := range jobs {
			start, delta, err := rl.Wait(ctx)
			if err != nil {
				return
			}
			fmt.Printf("START request %d at %s (Δ=%v)\n", id, start.Format(time.RFC3339Nano), delta)
			Request(id)
		}
	}

	wg.Add(m)
	for i := 0; i < m; i++ {
		go worker()
	}

	go func() {
		defer close(jobs)
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case jobs <- i:
			}
		}
	}()

	wg.Wait()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const numRequests = 10
	const numWorkers = 3
	const minSpacing = 200 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	runWithPoolAndRateLimit(ctx, numRequests, numWorkers, minSpacing)
}
