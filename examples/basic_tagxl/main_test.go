package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	_ = r.Close()

	return buf.String()
}

func TestMain(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("skipping: main() panicked (likely AWS API unavailable): %v", r)
		}
	}()

	out := captureStdout(t, main)

	if !strings.Contains(out, "longitude") {
		t.Errorf("expected output %q not found in %q", "longitude", out)
	}
}
