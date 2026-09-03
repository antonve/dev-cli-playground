package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/antonve/dev-cli-playground/internal/identity"
)

type diagnostic struct {
	Self       identity.Value `json:"self"`
	Echo       identity.Value `json:"echo"`
	EchoTarget string         `json:"echoTarget"`
	Fallback   bool           `json:"fallback"`
}

func handler(client *http.Client, gateway string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200); _, _ = w.Write([]byte("ok\n")) })
	mux.HandleFunc("/api/diagnostic", func(w http.ResponseWriter, r *http.Request) {
		route := r.Header.Get("x-dev-branch")
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, gateway+"/echo", nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if route != "" {
			req.Header.Set("x-dev-branch", route)
		}
		req.Host = "echo-api.internal"
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("echo request: %v", err), 502)
			return
		}
		defer resp.Body.Close()
		var echo identity.Value
		if err := json.NewDecoder(resp.Body).Decode(&echo); err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(diagnostic{Self: identity.Current("hello-api", route), Echo: echo, EchoTarget: "echo-api.internal", Fallback: route != "" && echo.Branch == "base"})
	})
	return mux
}
func main() {
	gateway := os.Getenv("DEV_INTERNAL_GATEWAY")
	if gateway == "" {
		gateway = "http://dev-cli-gateway"
	}
	client := &http.Client{Timeout: 2 * time.Second}
	log.Fatal(http.ListenAndServe(":8081", handler(client, gateway)))
}
