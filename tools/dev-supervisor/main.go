package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type supervisor struct {
	mu                          sync.Mutex
	base, next, current, health string
	cmd                         *exec.Cmd
	done                        chan error
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp := dst + ".tmp"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0555)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	if err := os.Chmod(tmp, 0555); err != nil {
		return err
	}
	return os.Rename(tmp, dst)
}
func (s *supervisor) start(path string) error {
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		return err
	}
	s.cmd = cmd
	s.done = make(chan error, 1)
	done := s.done
	go func() {
		err := cmd.Wait()
		log.Printf("child %s exited: %v", filepath.Base(path), err)
		done <- err
		close(done)
	}()
	return nil
}
func (s *supervisor) stop() {
	if s.cmd == nil || s.cmd.Process == nil {
		return
	}
	_ = s.cmd.Process.Signal(syscall.SIGTERM)
	select {
	case <-s.done:
	case <-time.After(2 * time.Second):
		_ = s.cmd.Process.Kill()
		<-s.done
	}
}
func (s *supervisor) ready(ctx context.Context) bool {
	client := http.Client{Timeout: 250 * time.Millisecond}
	t := time.NewTicker(150 * time.Millisecond)
	defer t.Stop()
	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, s.health, nil)
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return true
			}
		}
		select {
		case <-ctx.Done():
			return false
		case <-t.C:
		}
	}
}
func (s *supervisor) reload(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := os.Stat(s.next); err != nil {
		http.Error(w, "next binary missing", 400)
		return
	}
	backup := s.current + ".backup"
	if _, err := os.Stat(s.current); err == nil {
		_ = copyFile(s.current, backup)
	}
	if err := copyFile(s.next, s.current); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	s.stop()
	if err := s.start(s.current); err != nil {
		s.rollback(backup)
		http.Error(w, err.Error(), 500)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ok := s.ready(ctx)
	cancel()
	if !ok {
		s.stop()
		s.rollback(backup)
		http.Error(w, "new binary failed readiness; rolled back", 500)
		return
	}
	_, _ = fmt.Fprintln(w, "ready")
}
func (s *supervisor) rollback(backup string) {
	candidate := s.base
	if _, err := os.Stat(backup); err == nil {
		candidate = backup
	}
	if err := s.start(candidate); err != nil {
		log.Printf("rollback failed: %v", err)
	}
}
func main() {
	base := flag.String("binary", "", "base executable")
	next := flag.String("next", "", "staged executable")
	health := flag.String("health", "", "readiness URL")
	flag.Parse()
	if *base == "" || *next == "" || *health == "" {
		log.Fatal(errors.New("binary, next, and health are required"))
	}
	s := &supervisor{base: *base, next: *next, current: *next + ".current", health: *health}
	if err := s.start(s.base); err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/reload", s.reload)
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })
	log.Fatal(http.ListenAndServe(":9000", mux))
}
