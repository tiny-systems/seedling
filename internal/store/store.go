// Package store keeps the short codes. In memory for now — growing a real
// backend is exactly the kind of issue this repo exists for.
package store

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
)

// ErrInvalidURL is returned by ValidateURL when the given string is not an
// absolute http(s) URL with a host.
var ErrInvalidURL = errors.New("url must be an absolute http:// or https:// URL")

// ValidateURL reports whether raw is an absolute URL with an http or https
// scheme and a host. It rejects unparsable strings, relative references,
// and other schemes such as javascript:.
func ValidateURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidURL
	}
	if u.Host == "" {
		return ErrInvalidURL
	}
	return nil
}

// Store maps short codes to URLs.
type Store struct {
	mu    sync.Mutex
	next  int
	links map[string]string
}

// New returns an empty store.
func New() *Store {
	return &Store{links: map[string]string{}}
}

// Add saves a URL and returns its short code.
func (s *Store) Add(url string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	code := fmt.Sprintf("s%d", s.next)
	s.links[code] = url
	return code
}

// Get resolves a short code.
func (s *Store) Get(code string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	url, ok := s.links[code]
	return url, ok
}
