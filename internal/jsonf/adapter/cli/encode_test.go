package cli

import (
	"anykey/internal/jsonf/domain"
	"bytes"
	"encoding/json"
	"testing"
)

func TestWriteObjectsOrdered_OrderAndMissing(t *testing.T) {
	objs := []domain.Object{
		{"id": json.RawMessage("1"), "name": json.RawMessage(`"X"`)},
	}
	keep := []string{"name", "age"}
	var buf bytes.Buffer
	if err := WriteObjectsOrdered(&buf, objs, keep); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	want := `[{"name":"X"}]`
	if got != want {
		t.Fatalf("unexpected output\n got: %s\nwant: %s", got, want)
	}
}
