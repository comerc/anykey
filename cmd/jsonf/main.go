package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// filterJSONArrayFields reads a JSON array of objects and keeps only specified fields
func filterJSONArrayFields(input []byte, keepFields []string) ([]byte, error) {
	if len(keepFields) == 0 {
		return nil, fmt.Errorf("no fields specified; pass field names as arguments")
	}

	// Unmarshal into slice of maps with RawMessage to preserve original JSON values
	var objects []map[string]json.RawMessage
	if err := json.Unmarshal(input, &objects); err != nil {
		return nil, fmt.Errorf("input must be a JSON array of objects: %w", err)
	}

	// Build a set for quick lookup and also keep the order of fields as provided
	keepOrder := make([]string, 0, len(keepFields))
	keepSet := make(map[string]struct{}, len(keepFields))
	for _, f := range keepFields {
		if _, exists := keepSet[f]; exists {
			continue
		}
		keepSet[f] = struct{}{}
		keepOrder = append(keepOrder, f)
	}

	// Build output manually to preserve field order as provided in keepOrder
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, obj := range objects {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('{')
		wrote := false
		for _, field := range keepOrder {
			val, ok := obj[field]
			if !ok {
				continue
			}
			if wrote {
				buf.WriteByte(',')
			}
			// Write key as proper JSON string
			keyBytes, _ := json.Marshal(field)
			buf.Write(keyBytes)
			buf.WriteByte(':')
			buf.Write(val)
			wrote = true
		}
		buf.WriteByte('}')
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

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

	output, err := filterJSONArrayFields(input, fields)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Stdout.Write(output)
}
