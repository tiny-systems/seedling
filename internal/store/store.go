// Package store keeps the short codes. In memory for now — growing a real
// backend is exactly the kind of issue this repo exists for.
package store

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Store maps short codes to URLs.
type Store struct {
	mu    sync.Mutex
	next  int
	links map[string]string
	path  string
}

// New returns an empty, in-memory-only store.
func New() *Store {
	return &Store{links: map[string]string{}}
}

// Load reads links from path, if it exists, and returns a Store that
// persists future changes back to it. A missing file just starts empty —
// it's created on the first Add.
func Load(path string) (*Store, error) {
	s := &Store{links: map[string]string{}, path: path}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.links); err != nil {
		return nil, err
	}

	// Resume the code counter past the highest loaded "sN" so new codes
	// never collide with ones restored from disk.
	for code := range s.links {
		if n, err := strconv.Atoi(strings.TrimPrefix(code, "s")); err == nil && n > s.next {
			s.next = n
		}
	}
	return s, nil
}

// Add saves a URL and returns its short code.
func (s *Store) Add(url string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	code := fmt.Sprintf("s%d", s.next)
	s.links[code] = url
	s.save()
	return code
}

// Get resolves a short code.
func (s *Store) Get(code string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	url, ok := s.links[code]
	return url, ok
}

// save writes the current links to path, if the store was created with one.
// Called with mu held.
func (s *Store) save() {
	if s.path == "" {
		return
	}
	data, err := json.MarshalIndent(s.links, "", "  ")
	if err != nil {
		log.Printf("store: marshaling links: %v", err)
		return
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		log.Printf("store: writing %s: %v", s.path, err)
	}
}
