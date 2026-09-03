package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInternalHeaderPropagation(t *testing.T) {
	echo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "echo-api.internal" {
			t.Fatalf("unexpected internal host %q", r.Host)
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"service":"echo-api","branch":"test","revision":"x","pod":"p","observedRoute":"` + r.Header.Get("x-dev-branch") + `","time":"2026-01-01T00:00:00Z"}`))
	}))
	defer echo.Close()
	client := echo.Client()
	req := httptest.NewRequest("GET", "/api/diagnostic", nil)
	req.Header.Set("x-dev-branch", "alice-feature-12345678")
	w := httptest.NewRecorder()
	handler(client, echo.URL).ServeHTTP(w, req)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"observedRoute":"alice-feature-12345678"`) {
		t.Fatalf("%d %s", w.Code, w.Body.String())
	}
}
