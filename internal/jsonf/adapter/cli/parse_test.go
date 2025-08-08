package cli

import "testing"

func TestParseJSONArray_Valid(t *testing.T) {
	in := []byte(`[{"a":1},{"b":2}]`)
	objs, err := ParseJSONArray(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(objs) != 2 {
		t.Fatalf("want 2 objects, got %d", len(objs))
	}
}

func TestParseJSONArray_Invalid(t *testing.T) {
	in := []byte(`{"a":1}`)
	if _, err := ParseJSONArray(in); err == nil {
		t.Fatalf("expected error for non-array input")
	}
}
