package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/antonve/dev-cli-playground/internal/identity"
)

var routePattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

type diagnostic struct {
	Self       identity.Value `json:"self"`
	Echo       identity.Value `json:"echo"`
	EchoTarget string         `json:"echoTarget"`
	Fallback   bool           `json:"fallback"`
}

func target(ctx context.Context, client *http.Client, route, namespace string) (string, bool) {
	base := "http://echo-api." + namespace + ".svc.cluster.local:8082"
	if !routePattern.MatchString(route) {
		return base, true
	}
	overlay := "http://echo-api-dev-" + route + "." + namespace + ".svc.cluster.local:8082"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, overlay+"/readyz", nil)
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return overlay, false
		}
	}
	return base, true
}
func handler(client *http.Client, namespace string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200); _, _ = w.Write([]byte("ok\n")) })
	mux.HandleFunc("/api/diagnostic", func(w http.ResponseWriter, r *http.Request) {
		route := r.Header.Get("x-dev-branch")
		endpoint, fallback := target(r.Context(), client, route, namespace)
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, endpoint+"/echo", nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if route != "" {
			req.Header.Set("x-dev-branch", route)
		}
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
		_ = json.NewEncoder(w).Encode(diagnostic{Self: identity.Current("hello-api", route), Echo: echo, EchoTarget: endpoint, Fallback: fallback})
	})
	return mux
}
func main() {
	ns := os.Getenv("DEV_NAMESPACE")
	if ns == "" {
		ns = "dev-cli-playground"
	}
	client := &http.Client{Timeout: 2 * time.Second}
	log.Fatal(http.ListenAndServe(":8081", handler(client, ns)))
}
