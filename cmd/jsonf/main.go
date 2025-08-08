package main

import (
	"fmt"
	"os"

	adapter "anykey/internal/jsonf/adapter"
)

func main() {
	fields := os.Args[1:]
	if len(fields) == 0 {
		fmt.Fprintln(os.Stderr, "usage: jsonf FIELD [FIELD ...]  # reads JSON array from stdin and outputs filtered array")
		os.Exit(2)
	}

	if err := adapter.StreamFilterAndWrite(os.Stdin, fields, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
