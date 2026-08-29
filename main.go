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

	"github.com/tiny-systems/seedling/internal/store"
)

func main() {
	linksFile := os.Getenv("LINKS_FILE")
	if linksFile == "" {
		linksFile = "./links.json"
	}
	s, err := store.Load(linksFile)
	if err != nil {
		log.Fatalf("loading links from %s: %v", linksFile, err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue("url")
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

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("seedling listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
