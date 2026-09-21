package store

import "testing"

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

func TestValidateURL(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"javascript scheme rejected", "javascript:alert(1)", true},
		{"not a url", "hello", true},
		{"valid https url", "https://example.com/x", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateURL(tc.url)
			if tc.wantErr && err == nil {
				t.Fatalf("ValidateURL(%q) = nil, want error", tc.url)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("ValidateURL(%q) = %v, want nil", tc.url, err)
			}
		})
	}
}
