package main

import (
	"anykey/internal/limiter/usecase"
	"context"
	"fmt"
	"math/rand"
	"time"
)

// requestCLI эмулирует внешний запрос, который занимает некоторое время
func requestCLI(id int) {
	time.Sleep(time.Duration(50+rand.Intn(150)) * time.Millisecond)
	fmt.Printf("done request %d at %s\n", id, time.Now().Format(time.RFC3339Nano))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const numRequests = 10
	const numWorkers = 3
	const minSpacing = 200 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usecase.RunWithPoolAndRateLimit(ctx, numRequests, numWorkers, minSpacing, requestCLI, usecase.StdoutOnStart)
}
