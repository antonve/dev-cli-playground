package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type rewriteTransport struct{ echo *httptest.Server }

func (r rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	clone.URL.Scheme = "http"
	clone.URL.Host = strings.TrimPrefix(r.echo.URL, "http://")
	return http.DefaultTransport.RoundTrip(clone)
}
func TestInternalHeaderPropagation(t *testing.T) {
	echo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/readyz" {
			w.WriteHeader(200)
			return
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"service":"echo-api","branch":"test","revision":"x","pod":"p","observedRoute":"` + r.Header.Get("x-dev-branch") + `","time":"2026-01-01T00:00:00Z"}`))
	}))
	defer echo.Close()
	client := &http.Client{Transport: rewriteTransport{echo: echo}}
	req := httptest.NewRequest("GET", "/api/diagnostic", nil)
	req.Header.Set("x-dev-branch", "alice-feature-12345678")
	w := httptest.NewRecorder()
	handler(client, "test").ServeHTTP(w, req)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"observedRoute":"alice-feature-12345678"`) {
		t.Fatalf("%d %s", w.Code, w.Body.String())
	}
}
