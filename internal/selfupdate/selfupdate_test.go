package selfupdate

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestNormalizeTag(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in  string
		out string
	}{
		{"1.2.3", "v1.2.3"},
		{"v1.2.3", "v1.2.3"},
		{"", ""},
		{"dev", "dev"}, // non-semver: returned as-is
	}
	for _, c := range cases {
		if got := normalizeTag(c.in); got != c.out {
			t.Fatalf("normalizeTag(%q)=%q want %q", c.in, got, c.out)
		}
	}
}

func TestNormalizeVersion(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in  string
		out string
	}{
		{"1.2.3", "v1.2.3"},
		{"v1.2.3", "v1.2.3"},
		{"", ""},
		{"dev", ""},
	}
	for _, c := range cases {
		if got := normalizeVersion(c.in); got != c.out {
			t.Fatalf("normalizeVersion(%q)=%q want %q", c.in, got, c.out)
		}
	}
}

func TestLatestTag_FilteringStableVsRC(t *testing.T) {

	origAPI := ghAPIReleasesURL
	defer func() { ghAPIReleasesURL = origAPI }()

	// Make template accept owner/repo segments
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/releases") {
			http.NotFound(w, r)
			return
		}
		// Return in "newest first" order
		resp := []map[string]any{
			{"tag_name": "v2.0.0-rc1", "prerelease": true, "draft": false},
			{"tag_name": "v1.5.0", "prerelease": false, "draft": false},
			{"tag_name": "v1.4.0", "prerelease": false, "draft": true}, // draft: excluded
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer s.Close()

	ghAPIReleasesURL = s.URL + "/%s/%s/releases?per_page=20"

	ctx := context.Background()
	tag, pre, err := LatestTag(ctx, false) // stable only
	if err != nil {
		t.Fatalf("LatestTag error: %v", err)
	}
	if tag != "v1.5.0" || pre {
		t.Fatalf("LatestTag stable expected v1.5.0, got %s prerelease=%v", tag, pre)
	}

	tag, pre, err = LatestTag(ctx, true) // include prerelease
	if err != nil {
		t.Fatalf("LatestTag error: %v", err)
	}
	if tag != "v2.0.0-rc1" || !pre {
		t.Fatalf("LatestTag prerelease expected v2.0.0-rc1 prerelease, got %s prerelease=%v", tag, pre)
	}
}

func TestCheckForUpdate(t *testing.T) {

	origAPI := ghAPIReleasesURL
	defer func() { ghAPIReleasesURL = origAPI }()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/releases") {
			http.NotFound(w, r)
			return
		}
		resp := []map[string]any{
			{"tag_name": "v1.1.0", "prerelease": false, "draft": false},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer s.Close()

	ghAPIReleasesURL = s.URL + "/%s/%s/releases?per_page=20"

	latest, has, err := CheckForUpdate(context.Background(), "v1.0.0", false)
	if err != nil {
		t.Fatalf("CheckForUpdate error: %v", err)
	}
	if latest != "v1.1.0" || !has {
		t.Fatalf("expected update to v1.1.0, has=%v got %s", has, latest)
	}

	latest, has, err = CheckForUpdate(context.Background(), "dev", false)
	if err != nil {
		t.Fatalf("CheckForUpdate error: %v", err)
	}
	if latest != "v1.1.0" || !has {
		t.Fatalf("dev should always consider valid latest as newer; got latest=%s has=%v", latest, has)
	}
}

func TestDownloadAndReplace_TarGz(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("tar.gz replacement test targets POSIX flow")
	}

	origBase := ghDownloadBaseURL
	origExec := execPathFunc
	defer func() {
		ghDownloadBaseURL = origBase
		execPathFunc = origExec
	}()

	// Prepare temp dir and "current executable"
	tmpDir := t.TempDir()
	binName := "decoder"
	destPath := filepath.Join(tmpDir, binName)
	if err := os.WriteFile(destPath, []byte("OLD_BIN"), 0o755); err != nil {
		t.Fatalf("write dest: %v", err)
	}
	execPathFunc = func() (string, error) { return destPath, nil }

	// Create archive that contains the new binary
	osName, archName, format := inferPlatform()
	if format != "tar.gz" {
		t.Skip("runtime format not tar.gz; skipping")
	}
	tag := "v1.2.3"
	archiveName := repo + "-1.2.3-" + osName + "-" + archName + ".tar.gz"
	checksumName := repo + "-1.2.3-checksums.txt"

	archiveBytes := mustBuildTarGz(t, binName, []byte("NEW_BIN"))
	sha := sha256.Sum256(archiveBytes)
	checksum := hex.EncodeToString(sha[:]) + "  " + archiveName + "\n"

	// Test server for both checksum and archive
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/"+archiveName):
			_, _ = w.Write(archiveBytes)
			return
		case strings.HasSuffix(r.URL.Path, "/"+checksumName):
			_, _ = io.WriteString(w, checksum)
			return
		default:
			http.NotFound(w, r)
		}
	}))
	defer s.Close()

	ghDownloadBaseURL = s.URL + "/%s/%s/%s"

	// Execute replacement
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := downloadAndReplace(ctx, tag); err != nil {
		t.Fatalf("downloadAndReplace error: %v", err)
	}

	// Verify contents replaced
	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read dest after replace: %v", err)
	}
	if string(got) != "NEW_BIN" {
		t.Fatalf("expected binary content NEW_BIN, got %q", string(got))
	}
}

func TestUnzipSingleBinary(t *testing.T) {
	t.Parallel()
	// Build a zip with a nested file; extraction should match by base name
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "a.zip")
	outPath := filepath.Join(tmpDir, "decoder")
	binName := "decoder"

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	w, err := zw.Create("nested/dir/" + binName)
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := w.Write([]byte("ZIP_BIN")); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	if err := os.WriteFile(zipPath, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("write zip: %v", err)
	}

	if err := unzipSingleBinary(zipPath, binName, outPath); err != nil {
		t.Fatalf("unzipSingleBinary: %v", err)
	}
	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read out: %v", err)
	}
	if string(got) != "ZIP_BIN" {
		t.Fatalf("unexpected content %q", string(got))
	}
}

func TestVerifyChecksumMismatch(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	archive := filepath.Join(tmpDir, "file.tar.gz")
	if err := os.WriteFile(archive, []byte("X"), 0o644); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	sum := filepath.Join(tmpDir, "checksums.txt")
	if err := os.WriteFile(sum, []byte("deadbeef  file.tar.gz\n"), 0o644); err != nil {
		t.Fatalf("write checksum: %v", err)
	}
	if err := verifyChecksum(archive, sum, "file.tar.gz"); err == nil {
		t.Fatalf("expected checksum mismatch error")
	}
}

// Helper to build a tar.gz containing a single file named binName with given content in a nested directory.
func mustBuildTarGz(t *testing.T, binName string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	dir := "decoder-1.2.3-" + runtime.GOOS + "-" + runtime.GOARCH
	// write header for file
	hdr := &tar.Header{
		Name: filepath.ToSlash(filepath.Join(dir, binName)),
		Mode: 0o755,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("tar header: %v", err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatalf("tar write: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("tar close: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}
	return buf.Bytes()
}

func TestLatestTag_NoSuitableRelease(t *testing.T) {
	origAPI := ghAPIReleasesURL
	defer func() { ghAPIReleasesURL = origAPI }()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/releases") {
			http.NotFound(w, r)
			return
		}
		// Only draft and rc while includePrerelease=false should both be skipped.
		resp := []map[string]any{
			{"tag_name": "v2.0.0-rc1", "prerelease": true, "draft": false},
			{"tag_name": "v1.5.0", "prerelease": false, "draft": true},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer s.Close()

	ghAPIReleasesURL = s.URL + "/%s/%s/releases?per_page=20"

	if _, _, err := LatestTag(context.Background(), false); err == nil {
		t.Fatalf("expected error when no suitable release is found")
	}
}

func TestUpdateToLatest_AlreadyUpToDate(t *testing.T) {
	origAPI := ghAPIReleasesURL
	defer func() { ghAPIReleasesURL = origAPI }()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/releases") {
			http.NotFound(w, r)
			return
		}
		resp := []map[string]any{
			{"tag_name": "v1.1.0", "prerelease": false, "draft": false},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer s.Close()

	ghAPIReleasesURL = s.URL + "/%s/%s/releases?per_page=20"

	if _, err := UpdateToLatest(context.Background(), "v1.1.0", false); err == nil {
		t.Fatalf("expected already up to date error")
	}
}

func TestHttpDownload_Non200(t *testing.T) {
	tmp := t.TempDir()
	dest := filepath.Join(tmp, "x.bin")

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := httpDownload(ctx, s.URL+"/x", dest); err == nil {
		t.Fatalf("expected error on non-200 response")
	}
}

func TestVerifyChecksum_NotFound(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "file.tar.gz")
	if err := os.WriteFile(archive, []byte("X"), 0o644); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	sum := filepath.Join(tmp, "checksums.txt")
	// Checksum file has unrelated line only.
	if err := os.WriteFile(sum, []byte("deadbeef  other.tar.gz\n"), 0o644); err != nil {
		t.Fatalf("write checksum: %v", err)
	}
	if err := verifyChecksum(archive, sum, "file.tar.gz"); err == nil {
		t.Fatalf("expected error when checksum entry is missing")
	}
}
