package main

import (
	adapter "anykey/internal/jsonf/adapter/cli"
	"anykey/internal/jsonf/usecase"
	"fmt"
	"io"
	"os"
)

// никаких прокси и совместимости — только вызовы адаптеров и юзкейсов

func main() {
	fields := os.Args[1:]
	if len(fields) == 0 {
		fmt.Fprintln(os.Stderr, "usage: jsonf FIELD [FIELD ...]  # reads JSON array from stdin and outputs filtered array")
		os.Exit(2)
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read stdin:", err)
		os.Exit(1)
	}

	objects, err := adapter.ParseJSONArray(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	keepOrder := usecase.DedupeKeepOrder(fields)
	filtered := usecase.FilterObjects(objects, keepOrder)
	if err := adapter.WriteObjectsOrdered(os.Stdout, filtered, keepOrder); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
