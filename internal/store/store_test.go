package store

import (
	"path/filepath"
	"testing"
)

func TestAddThenGet(t *testing.T) {
	s := New()
	code := s.Add("https://example.com")
	got, ok := s.Get(code)
	if !ok || got != "https://example.com" {
		t.Fatalf("Get(%q) = %q, %v", code, got, ok)
	}
}

func TestGetUnknownCode(t *testing.T) {
	if _, ok := New().Get("nope"); ok {
		t.Fatal("unknown code resolved")
	}
}

func TestLoadMissingFileStartsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "links.json")
	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load(%q) = %v", path, err)
	}
	if _, ok := s.Get("s1"); ok {
		t.Fatal("empty store resolved a code")
	}
}

func TestPersistsAcrossLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "links.json")

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load(%q) = %v", path, err)
	}
	code := s.Add("https://example.com")

	reloaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load(%q) = %v", path, err)
	}
	got, ok := reloaded.Get(code)
	if !ok || got != "https://example.com" {
		t.Fatalf("Get(%q) after reload = %q, %v", code, got, ok)
	}

	// New codes after a reload must not collide with restored ones.
	next := reloaded.Add("https://second.example.com")
	if next == code {
		t.Fatalf("Add after reload reused code %q", code)
	}
}
