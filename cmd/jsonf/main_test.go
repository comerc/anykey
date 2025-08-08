package main

import (
	"strings"
	"testing"
)

func TestFilterJSONArrayFields_Basic(t *testing.T) {
	input := `[{"id":2,"name":"Alice","age":30},{"age":25,"name":"Bob","id":1}]`
	out, err := filterJSONArrayFields([]byte(input), []string{"name", "age"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := string(out)
	want := `[{"name":"Alice","age":30},{"name":"Bob","age":25}]`
	if got != want {
		t.Fatalf("unexpected output.\n got: %s\nwant: %s", got, want)
	}
}

func TestFilterJSONArrayFields_OrderAndDedup(t *testing.T) {
	input := `[{"id":1,"name":"X","age":10}]`
	out, err := filterJSONArrayFields([]byte(input), []string{"age", "name", "age"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != `[{"age":10,"name":"X"}]` {
		t.Fatalf("order/dedup failed, got: %s", string(out))
	}
}

func TestFilterJSONArrayFields_InvalidInput(t *testing.T) {
	_, err := filterJSONArrayFields([]byte(`{"not":"array"}`), []string{"a"})
	if err == nil || !strings.Contains(err.Error(), "JSON array") {
		t.Fatalf("expected JSON array error, got: %v", err)
	}
}

func TestFilterJSONArrayFields_NoFields(t *testing.T) {
	_, err := filterJSONArrayFields([]byte(`[]`), nil)
	if err == nil || !strings.Contains(err.Error(), "no fields specified") {
		t.Fatalf("expected no fields error, got: %v", err)
	}
}
