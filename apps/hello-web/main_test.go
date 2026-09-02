package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServesReactModule(t *testing.T) {
	w := httptest.NewRecorder()
	server{fallback: embedded}.ServeHTTP(w, httptest.NewRequest("GET", "/src/app.tsx", nil))
	if w.Code != 200 || !strings.Contains(w.Header().Get("content-type"), "javascript") || !strings.Contains(w.Body.String(), "React") {
		t.Fatalf("%d %s %s", w.Code, w.Header(), w.Body.String())
	}
}
