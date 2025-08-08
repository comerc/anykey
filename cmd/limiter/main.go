package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	usecase "anykey/internal/limiter/usecase"
)

// requestCLI эмулирует внешний запрос, который занимает некоторое время
func requestCLI(id int) {
	time.Sleep(time.Duration(50+rand.Intn(150)) * time.Millisecond)
	fmt.Printf("done request %d at %s\n", id, time.Now().Format(time.RFC3339Nano))
}

// stdoutOnStart — удобный хук для CLI
func stdoutOnStart(id int, start time.Time, delta time.Duration) {
	fmt.Printf("START request %d at %s (Δ=%v)\n", id, start.Format(time.RFC3339Nano), delta)
}

func main() {
	const numRequests = 10
	const numWorkers = 3
	const minSpacing = 200 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usecase.RunWithPoolAndRateLimit(ctx, numRequests, numWorkers, minSpacing, requestCLI, stdoutOnStart)
}
