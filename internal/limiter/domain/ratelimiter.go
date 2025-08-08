package domain

import (
	"context"
	"sync"
	"time"
)

// RateLimiter управляет запуском запросов с ограничением на интервал между запросами
type RateLimiter struct {
	Mu         sync.Mutex
	LastStart  time.Time
	MinSpacing time.Duration
}

// Wait блокируется до тех пор, пока запуск нового запроса не нарушит минимальный интервал.
// Возвращает фактическое время старта и дельту до предыдущего старта.
func (r *RateLimiter) Wait(ctx context.Context) (time.Time, time.Duration, error) {
	for {
		// Быстрая отмена, если контекст уже завершён
		if err := ctx.Err(); err != nil {
			return time.Time{}, 0, err
		}

		r.Mu.Lock()
		now := time.Now()
		target := r.LastStart.Add(r.MinSpacing)
		if r.LastStart.IsZero() || !now.Before(target) {
			var delta time.Duration
			if !r.LastStart.IsZero() {
				delta = now.Sub(r.LastStart)
			}
			r.LastStart = now
			r.Mu.Unlock()
			return now, delta, nil
		}
		wait := target.Sub(now)
		r.Mu.Unlock()

		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return time.Time{}, 0, ctx.Err()
		case <-timer.C:
		}
	}
}
