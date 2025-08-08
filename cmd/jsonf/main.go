package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// dedupeKeepOrder возвращает порядок полей без дубликатов
func dedupeKeepOrder(keepFields []string) []string {
	dedup := make([]string, 0, len(keepFields))
	seen := make(map[string]struct{}, len(keepFields))
	for _, f := range keepFields {
		if _, ok := seen[f]; ok {
			continue
		}
		seen[f] = struct{}{}
		dedup = append(dedup, f)
	}
	return dedup
}

// parseJSONArray парсит вход в массив объектов
func parseJSONArray(input []byte) ([]map[string]json.RawMessage, error) {
	var objects []map[string]json.RawMessage
	if err := json.Unmarshal(input, &objects); err != nil {
		return nil, fmt.Errorf("input must be a JSON array of objects: %w", err)
	}
	return objects, nil
}

// encodeObjectsOrdered кодирует только указанные поля и сохраняет их порядок
func writeObjectsOrdered(w io.Writer, objects []map[string]json.RawMessage, keepOrder []string) error {
	bw := bufio.NewWriter(w)
	bw.WriteByte('[')
	for i, obj := range objects {
		if i > 0 {
			bw.WriteByte(',')
		}
		bw.WriteByte('{')
		wrote := false
		for _, field := range keepOrder {
			if val, ok := obj[field]; ok {
				if wrote {
					bw.WriteByte(',')
				}
				keyBytes, _ := json.Marshal(field)
				bw.Write(keyBytes)
				bw.WriteByte(':')
				bw.Write(val)
				wrote = true
			}
		}
		bw.WriteByte('}')
	}
	bw.WriteByte(']')
	return bw.Flush()
}

// filterJSONArrayFields читает JSON-массив объектов и сохраняет только указанные поля
func filterJSONArrayFields(input []byte, keepFields []string) ([]byte, error) {
	if len(keepFields) == 0 {
		return nil, fmt.Errorf("no fields specified; pass field names as arguments")
	}

	objects, err := parseJSONArray(input)
	if err != nil {
		return nil, err
	}
	keepOrder := dedupeKeepOrder(keepFields)
	var buf bytes.Buffer
	if err := writeObjectsOrdered(&buf, objects, keepOrder); err != nil {
		return nil, err
	}
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

	objects, err := parseJSONArray(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	keepOrder := dedupeKeepOrder(fields)
	if err := writeObjectsOrdered(os.Stdout, objects, keepOrder); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
