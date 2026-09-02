package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCookieSelectionRejectsHeader(t *testing.T) {
	r := router{}
	req := httptest.NewRequest("GET", "https://example/", nil)
	req.Header.Set("x-dev-branch", "spoof")
	if got := r.selected(req); got != "" {
		t.Fatalf("trusted header: %q", got)
	}
	req.AddCookie(&http.Cookie{Name: cookieName, Value: "alice-feature-12345678"})
	if got := r.selected(req); got != "alice-feature-12345678" {
		t.Fatalf("cookie: %q", got)
	}
}
func TestSelectCookieScoped(t *testing.T) {
	r := router{}
	req := httptest.NewRequest("GET", "https://example/__dev/select?branch=alice-feature-12345678", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c := w.Result().Cookies()[0]
	if c.Path != "/" || !c.HttpOnly || !c.Secure {
		t.Fatalf("cookie not scoped securely: %#v", c)
	}
}
