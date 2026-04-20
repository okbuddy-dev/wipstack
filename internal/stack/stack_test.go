package stack

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStackIsEmpty(t *testing.T) {
	s := New()
	if !s.IsEmpty() {
		t.Error("new stack should be empty")
	}
	if len(s.Items()) != 0 {
		t.Error("new stack should have no items")
	}
}

func TestPushAndItems(t *testing.T) {
	s := New()
	s.Push("first")
	s.Push("second")
	got := s.Items()
	want := []string{"first", "second"}
	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("item %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestPop(t *testing.T) {
	s := New()
	s.Push("first")
	s.Push("second")
	item, ok := s.Pop()
	if !ok {
		t.Fatal("expected ok=true from pop on non-empty stack")
	}
	if item != "second" {
		t.Errorf("expected %q, got %q", "second", item)
	}
	remaining := s.Items()
	if len(remaining) != 1 || remaining[0] != "first" {
		t.Errorf("expected [first], got %v", remaining)
	}
}

func TestPopEmpty(t *testing.T) {
	s := New()
	item, ok := s.Pop()
	if ok {
		t.Error("expected ok=false from pop on empty stack")
	}
	if item != "" {
		t.Errorf("expected empty string, got %q", item)
	}
}

func TestClear(t *testing.T) {
	s := New()
	s.Push("a")
	s.Push("b")
	cleared := s.Clear()
	if len(cleared) != 2 || cleared[0] != "a" || cleared[1] != "b" {
		t.Errorf("expected [a b], got %v", cleared)
	}
	if !s.IsEmpty() {
		t.Error("stack should be empty after clear")
	}
}

func TestClearEmpty(t *testing.T) {
	s := New()
	cleared := s.Clear()
	if len(cleared) != 0 {
		t.Errorf("expected empty slice, got %v", cleared)
	}
}

func TestItemsReturnsCopy(t *testing.T) {
	s := New()
	s.Push("a")
	items := s.Items()
	items[0] = "modified"
	if s.Items()[0] != "a" {
		t.Error("Items() should return a copy, not a reference")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stack.json")

	s := New()
	s.Push("first")
	s.Push("second")

	err := s.Save(path)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	got := loaded.Items()
	want := []string{"first", "second"}
	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("item %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nosuch.json")

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load of missing file should not error, got: %v", err)
	}
	if !s.IsEmpty() {
		t.Error("Load of missing file should return empty stack")
	}
}

func TestSaveCreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "deep", "stack.json")

	s := New()
	s.Push("test")
	err := s.Save(path)
	if err != nil {
		t.Fatalf("Save should create parent dirs, got: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file was not created")
	}
}

func TestLoadCorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not valid json{"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("Load of corrupt file should return error")
	}
}

func TestSaveAtomicOverwrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stack.json")

	s := New()
	s.Push("first")
	s.Save(path)

	s2 := New()
	s2.Push("alpha")
	s2.Push("beta")
	s2.Save(path)

	loaded, _ := Load(path)
	items := loaded.Items()
	if len(items) != 2 || items[0] != "alpha" || items[1] != "beta" {
		t.Errorf("expected [alpha beta], got %v", items)
	}
}