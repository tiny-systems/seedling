// Seedling is a tiny link shortener — and the playground where tiny
// (https://github.com/tiny-systems/tiny) shows itself off: file an issue,
// label it `tiny`, and a coding-agent session on a Kubernetes cluster
// picks it up and sends back a pull request.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tiny-systems/seedling/internal/store"
)

func newMux(s *store.Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		url := strings.TrimSpace(r.FormValue("url"))
		if url == "" {
			http.Error(w, "url is required", http.StatusBadRequest)
			return
		}
		code := s.Add(url)
		fmt.Fprintf(w, "http://%s/%s\n", r.Host, code)
	})

	mux.HandleFunc("GET /{code}", func(w http.ResponseWriter, r *http.Request) {
		url, ok := s.Get(r.PathValue("code"))
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	})

	return mux
}

func main() {
	s := store.New()
	mux := newMux(s)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("seedling listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
