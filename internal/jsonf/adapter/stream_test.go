package cli

import (
	"bytes"
	"testing"
)

func TestStreamFilterAndWrite(t *testing.T) {
	input := `[{"id":2,"name":"Alice","age":30},{"age":25,"name":"Bob","id":1}]`
	var out bytes.Buffer
	if err := StreamFilterAndWrite(bytes.NewReader([]byte(input)), []string{"name", "age"}, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.String()
	want := `[{"name":"Alice","age":30},{"age":25,"name":"Bob"}]`
	if got != want {
		t.Fatalf("unexpected output.\n got: %s\nwant: %s", got, want)
	}
}
