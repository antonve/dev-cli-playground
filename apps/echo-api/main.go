package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/antonve/dev-cli-playground/internal/identity"
)

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(identity.Current("echo-api", r.Header.Get("x-dev-branch")))
	})
	return mux
}

func main() {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8082"
	}
	log.Printf("echo-api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler()))
}
