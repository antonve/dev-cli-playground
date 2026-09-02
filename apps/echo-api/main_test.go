package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEchoIdentity(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/echo", nil)
	r.Header.Set("x-dev-branch", "alice-feature-12345678")
	w := httptest.NewRecorder()
	handler().ServeHTTP(w, r)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"observedRoute":"alice-feature-12345678"`) {
		t.Fatalf("%d %s", w.Code, w.Body.String())
	}
}
