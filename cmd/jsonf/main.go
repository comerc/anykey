package main

import (
	adapter "anykey/internal/jsonf/adapter/cli"
	"anykey/internal/jsonf/usecase"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// dedupeKeepOrder возвращает порядок полей без дубликатов
// dedupeKeepOrder проксирует usecase для совместимости с существующими тестами
func dedupeKeepOrder(keepFields []string) []string { return usecase.DedupeKeepOrder(keepFields) }

// parseJSONArray парсит вход в массив объектов
// parseJSONArray проксирует адаптер для совместимости
func parseJSONArray(input []byte) ([]map[string]json.RawMessage, error) {
	objs, err := adapter.ParseJSONArray(input)
	if err != nil {
		return nil, err
	}
	// адаптер возвращает тип domain.Object, совместимый с map[string]json.RawMessage
	return objs, nil
}

// encodeObjectsOrdered кодирует только указанные поля и сохраняет их порядок
// writeObjectsOrdered проксирует адаптер
func writeObjectsOrdered(w io.Writer, objects []map[string]json.RawMessage, keepOrder []string) error {
	return adapter.WriteObjectsOrdered(w, objects, keepOrder)
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
