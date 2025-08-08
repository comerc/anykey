package usecase

import (
	"anykey/internal/jsonf/domain"
	"encoding/json"
	"testing"
)

func TestDedupeKeepOrder(t *testing.T) {
	in := []string{"age", "name", "age", "id", "name"}
	got := DedupeKeepOrder(in)
	want := []string{"age", "name", "id"}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order mismatch at %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestFilterObjects_Basic(t *testing.T) {
	objs := []domain.Object{
		{
			"id":   json.RawMessage("2"),
			"name": json.RawMessage(`"Alice"`),
			"age":  json.RawMessage("30"),
		},
		{
			"age":  json.RawMessage("25"),
			"name": json.RawMessage(`"Bob"`),
			"id":   json.RawMessage("1"),
		},
	}
	keep := []string{"name", "age"}
	out := FilterObjects(objs, keep)

	// Проверяем, что в выходе только нужные ключи
	for i, o := range out {
		if _, ok := o["name"]; !ok {
			t.Fatalf("obj %d: missing name", i)
		}
		if _, ok := o["age"]; !ok {
			t.Fatalf("obj %d: missing age", i)
		}
		if _, ok := o["id"]; ok {
			t.Fatalf("obj %d: unexpected id present", i)
		}
	}
}
