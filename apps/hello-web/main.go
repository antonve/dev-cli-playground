package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//go:embed web/*
var embedded embed.FS

type server struct {
	root     string
	fallback fs.FS
}

func (s server) read(name string) ([]byte, error) {
	if s.root != "" {
		if b, err := os.ReadFile(filepath.Join(s.root, name)); err == nil {
			return b, nil
		}
	}
	return fs.ReadFile(s.fallback, "web/"+name)
}
func (s server) fingerprint() string {
	h := sha256.New()
	var names []string
	if s.root != "" {
		_ = filepath.WalkDir(s.root, func(p string, d fs.DirEntry, e error) error {
			if e == nil && !d.IsDir() {
				names = append(names, p)
			}
			return nil
		})
	}
	sort.Strings(names)
	for _, n := range names {
		b, _ := os.ReadFile(n)
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))
}
func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/readyz" {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok\n"))
		return
	}
	if r.URL.Path == "/__hmr" {
		w.Header().Set("content-type", "text/event-stream")
		w.Header().Set("cache-control", "no-cache")
		f, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "stream unsupported", 500)
			return
		}
		last := s.fingerprint()
		t := time.NewTicker(200 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-r.Context().Done():
				return
			case <-t.C:
				next := s.fingerprint()
				if next != last {
					last = next
					fmt.Fprintf(w, "data: %s\n\n", next)
					f.Flush()
				}
			}
		}
	}
	name := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/")
	if name == "." || name == "" {
		name = "index.html"
	}
	b, err := s.read(name)
	if err != nil {
		b, err = s.read("index.html")
	}
	if err != nil {
		http.NotFound(w, r)
		return
	}
	ext := filepath.Ext(name)
	if ext == ".tsx" {
		w.Header().Set("content-type", "text/javascript; charset=utf-8")
	} else if mt := mime.TypeByExtension(ext); mt != "" {
		w.Header().Set("content-type", mt)
	}
	w.Header().Set("cache-control", "no-store")
	_, _ = w.Write(b)
}
func main() {
	root := flag.String("root", "", "live source root")
	_ = flag.String("fallback-root", "", "compatibility flag")
	flag.Parse()
	log.Fatal(http.ListenAndServe(":8080", server{root: *root, fallback: embedded}))
}
