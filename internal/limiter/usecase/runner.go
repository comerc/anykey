package usecase

import (
	"context"
	"sync"
	"time"

	domain "anykey/internal/limiter/domain"
)

// RunWithPoolAndRateLimit запускает n задач в m воркерах, строго соблюдая интервал между стартапами
// requester выполняет полезную работу; onStart вызывается перед запуском requester
func RunWithPoolAndRateLimit(ctx context.Context, n int, m int, minInterval time.Duration, requester func(int), onStart func(id int, start time.Time, delta time.Duration)) {
	rl := &domain.RateLimiter{MinSpacing: minInterval}

	jobs := make(chan int)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for id := range jobs {
			start, delta, err := rl.Wait(ctx)
			if err != nil {
				return
			}
			if onStart != nil {
				onStart(id, start, delta)
			}
			requester(id)
		}
	}

	if m <= 0 {
		panic("workers must be > 0")
	}
	if n < 0 {
		panic("n must be >= 0")
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
