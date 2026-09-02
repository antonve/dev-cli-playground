package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const cookieName = "dev_branch"

var routePattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

type router struct {
	namespace string
	client    *http.Client
}

func (r router) selected(req *http.Request) string {
	c, err := req.Cookie(cookieName)
	if err != nil || !routePattern.MatchString(c.Value) {
		return ""
	}
	return c.Value
}
func (r router) endpoint(ctx context.Context, service string, port int, route string) (string, bool) {
	base := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", service, r.namespace, port)
	if route == "" {
		return base, true
	}
	overlay := fmt.Sprintf("http://%s-dev-%s.%s.svc.cluster.local:%d", service, route, r.namespace, port)
	probe, _ := http.NewRequestWithContext(ctx, http.MethodGet, overlay+"/readyz", nil)
	resp, err := r.client.Do(probe)
	if err == nil {
		resp.Body.Close()
		if resp.StatusCode == 200 {
			return overlay, false
		}
	}
	return base, true
}
func (r router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/__dev/select" {
		route := req.URL.Query().Get("branch")
		if route != "" && !routePattern.MatchString(route) {
			http.Error(w, "invalid route key", 400)
			return
		}
		max := 28800
		if route == "" {
			max = -1
		}
		http.SetCookie(w, &http.Cookie{Name: cookieName, Value: route, Path: "/", MaxAge: max, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode})
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	if req.URL.Path == "/__dev/current" {
		route := r.selected(req)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"route": route})
		return
	}
	service, port := "hello-web", 8080
	if strings.HasPrefix(req.URL.Path, "/api/") {
		service, port = "hello-api", 8081
	}
	route := r.selected(req)
	endpoint, _ := r.endpoint(req.Context(), service, port, route)
	target, _ := url.Parse(endpoint)
	proxy := httputil.NewSingleHostReverseProxy(target)
	original := proxy.Director
	proxy.Director = func(out *http.Request) {
		original(out)
		out.Header.Del("x-dev-branch")
		if route != "" {
			out.Header.Set("x-dev-branch", route)
		}
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, err error) {
		http.Error(w, "upstream unavailable: "+err.Error(), http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, req)
}
func main() {
	ns := os.Getenv("DEV_NAMESPACE")
	if ns == "" {
		ns = "dev-cli-playground"
	}
	log.Fatal(http.ListenAndServe(":8080", router{namespace: ns, client: &http.Client{Timeout: 350 * time.Millisecond}}))
}
