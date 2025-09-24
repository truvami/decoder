package selfupdate

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSetClientTimeout(t *testing.T) {
	old := httpClient.Timeout
	defer func() { httpClient.Timeout = old }()

	SetClientTimeout(123 * time.Millisecond)
	if httpClient.Timeout != 123*time.Millisecond {
		t.Fatalf("timeout not applied: got %v", httpClient.Timeout)
	}
}

func TestHttpDownload_Success(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "OK")
	}))
	defer s.Close()

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "x.bin")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := httpDownload(ctx, s.URL+"/asset", dest); err != nil {
		t.Fatalf("httpDownload error: %v", err)
	}
	b, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("read dest: %v", err)
	}
	if string(b) != "OK" {
		t.Fatalf("unexpected content %q", string(b))
	}
}

func TestVerifyChecksum_Success(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "file.tar.gz")
	if err := os.WriteFile(archive, []byte("ABC"), 0o644); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	sum := sha256.Sum256([]byte("ABC"))
	chk := filepath.Join(tmp, "checksums.txt")
	if err := os.WriteFile(chk, []byte(fmt.Sprintf("%x  %s\n", sum[:], filepath.Base(archive))), 0o644); err != nil {
		t.Fatalf("write checksum: %v", err)
	}
	if err := verifyChecksum(archive, chk, filepath.Base(archive)); err != nil {
		t.Fatalf("verifyChecksum unexpected error: %v", err)
	}
}

func TestWinPendingReplaceErrorString(t *testing.T) {
	w := &winPendingReplace{NewPath: "/tmp/decoder.new"}
	if msg := w.Error(); !strings.Contains(msg, "/tmp/decoder.new") {
		t.Fatalf("error string does not contain path: %q", msg)
	}
}

func TestCopyFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.bin")
	dst := filepath.Join(tmp, "dst.bin")
	if err := os.WriteFile(src, []byte("CONTENT"), 0o600); err != nil {
		t.Fatalf("write src: %v", err)
	}
	if err := copyFile(src, dst, 0o644); err != nil {
		t.Fatalf("copyFile error: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	if string(got) != "CONTENT" {
		t.Fatalf("unexpected dst content %q", string(got))
	}
}
