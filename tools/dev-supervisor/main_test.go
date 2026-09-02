package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFileAtomic(t *testing.T) {
	d := t.TempDir()
	src := filepath.Join(d, "src")
	dst := filepath.Join(d, "dst")
	if err := os.WriteFile(src, []byte("working"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := copyFile(src, dst); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(dst)
	if err != nil || string(b) != "working" {
		t.Fatalf("%q %v", b, err)
	}
}
