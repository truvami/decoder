package selfupdate

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const (
	owner = "truvami"
	repo  = "decoder"
)

// HTTP client with a short timeout to keep checks non-blocking by default.
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// execPathFunc allows tests to override os.Executable().
var execPathFunc = os.Executable

// Overridable endpoints for tests.
var ghAPIReleasesURL = "https://api.github.com/repos/%s/%s/releases?per_page=20"
var ghDownloadBaseURL = "https://github.com/%s/%s/releases/download/%s"

// SetClientTimeout allows callers to override the HTTP client's timeout, e.g. for long-running updates.
func SetClientTimeout(d time.Duration) {
	httpClient.Timeout = d
}

type ghRelease struct {
	TagName    string `json:"tag_name"`
	Prerelease bool   `json:"prerelease"`
	Draft      bool   `json:"draft"`
}

// normalizeTag ensures tags use a leading v, as goreleaser publishes using vX.Y.Z.
func normalizeTag(tag string) string {
	if tag == "" {
		return ""
	}
	if !strings.HasPrefix(tag, "v") && semver.IsValid("v"+tag) {
		return "v" + tag
	}
	// Some tags may already have "v".
	if semver.IsValid(tag) {
		return tag
	}
	// If not valid semver even with v prefix, just return as-is.
	return tag
}

// normalizeVersion converts a raw version (e.g. "1.2.3" or "dev") to a normalized semver tag ("v1.2.3").
// Non-semver values will return empty string.
func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	// if already a tag like v1.2.3
	if semver.IsValid(v) {
		return v
	}
	// try with v prefix
	if semver.IsValid("v" + v) {
		return "v" + v
	}
	return ""
}

// LatestTag returns the newest release tag.
// If includePrerelease is false, only stable (non-prerelease) releases are considered.
// If includePrerelease is true, the newest non-draft release is returned regardless of prerelease flag.
func LatestTag(ctx context.Context, includePrerelease bool) (string, bool, error) {
	// Use the releases API to have full control over prerelease filtering:
	// https://api.github.com/repos/{owner}/{repo}/releases?per_page=20
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(ghAPIReleasesURL, owner, repo), nil)
	if err != nil {
		return "", false, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "truvami-decoder-selfupdate")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", false, fmt.Errorf("github releases http status %d", resp.StatusCode)
	}

	var releases []ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", false, err
	}

	// Releases are returned sorted by creation date descending.
	for _, r := range releases {
		if r.Draft {
			continue
		}
		if !includePrerelease && r.Prerelease {
			continue
		}
		tag := normalizeTag(r.TagName)
		if tag == "" {
			continue
		}
		// Ensure it's valid semver for meaningful comparisons downstream.
		if !semver.IsValid(tag) {
			continue
		}
		// Also treat tags that contain "-rc" as prereleases unless --next is used.
		if !includePrerelease && strings.Contains(strings.ToLower(tag), "-rc") {
			continue
		}
		return tag, r.Prerelease, nil
	}

	return "", false, errors.New("no suitable release found")
}

// CheckForUpdate compares the current version with the latest release tag.
// It returns the latest tag and whether an update is available.
// If current is not a valid semver, it will only return hasUpdate=true when a valid latest tag exists.
func CheckForUpdate(ctx context.Context, current string, includePrerelease bool) (latest string, hasUpdate bool, err error) {
	latest, _, err = LatestTag(ctx, includePrerelease)
	if err != nil {
		return "", false, err
	}
	cur := normalizeVersion(current)

	// If current is invalid (e.g., "dev"), treat any valid latest as newer.
	if cur == "" {
		return latest, true, nil
	}

	// Compare semver (expects leading v).
	if semver.Compare(latest, cur) > 0 {
		return latest, true, nil
	}
	return latest, false, nil
}

// UpdateToLatest finds the correct latest tag (stable or next) and performs the update.
func UpdateToLatest(ctx context.Context, current string, includePrerelease bool) (newTag string, err error) {
	latest, _, err := LatestTag(ctx, includePrerelease)
	if err != nil {
		return "", err
	}
	cur := normalizeVersion(current)
	if cur != "" && semver.Compare(latest, cur) <= 0 {
		return "", fmt.Errorf("already up to date (current=%s latest=%s)", cur, latest)
	}
	if err := downloadAndReplace(ctx, latest); err != nil {
		return "", err
	}
	return latest, nil
}

func downloadAndReplace(ctx context.Context, tag string) error {
	osName, archName, format := inferPlatform()
	name := fmt.Sprintf("%s-%s-%s-%s.%s", repo, strings.TrimPrefix(tag, "v"), osName, archName, format)
	chk := fmt.Sprintf("%s-%s-checksums.txt", repo, strings.TrimPrefix(tag, "v"))

	base := fmt.Sprintf(ghDownloadBaseURL, owner, repo, tag)
	archiveURL := fmt.Sprintf("%s/%s", base, name)
	checksumURL := fmt.Sprintf("%s/%s", base, chk)

	// Download checksum and archive to temp dir
	tmpDir, err := os.MkdirTemp("", "decoder-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, name)
	checksumPath := filepath.Join(tmpDir, chk)

	if err := httpDownload(ctx, archiveURL, archivePath); err != nil {
		return fmt.Errorf("download archive: %w", err)
	}
	if err := httpDownload(ctx, checksumURL, checksumPath); err != nil {
		return fmt.Errorf("download checksum: %w", err)
	}
	if err := verifyChecksum(archivePath, checksumPath, filepath.Base(archivePath)); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}

	// Extract binary to temp file
	binName := repo
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	extractedPath := filepath.Join(tmpDir, binName)

	if format == "zip" {
		if err := unzipSingleBinary(archivePath, binName, extractedPath); err != nil {
			return err
		}
	} else {
		if err := untarGzSingleBinary(archivePath, binName, extractedPath); err != nil {
			return err
		}
	}

	// Ensure it is executable on unix
	if runtime.GOOS != "windows" {
		_ = os.Chmod(extractedPath, 0o755)
	}

	// Replace current executable
	dest, err := execPathFunc()
	if err != nil {
		return fmt.Errorf("cannot find current executable: %w", err)
	}
	// Ensure same directory for atomic rename
	destDir := filepath.Dir(dest)
	tmpNew := filepath.Join(destDir, binName+".tmp")

	// Move the new binary into place in the dest dir first
	if err := copyFile(extractedPath, tmpNew, 0o755); err != nil {
		return fmt.Errorf("preparing new binary: %w", err)
	}

	// Best-effort atomic replacement
	if runtime.GOOS == "windows" {
		// Windows cannot rename over a running executable.
		finalPath := filepath.Join(destDir, binName+".new")
		if err := os.Rename(tmpNew, finalPath); err != nil {
			return fmt.Errorf("prepare new binary on windows: %w", err)
		}
		// Inform the caller that the new binary is placed next to the current one.
		return &winPendingReplace{NewPath: finalPath}
	}

	// On POSIX, rename is atomic when same filesystem
	if err := os.Rename(tmpNew, dest); err != nil {
		// Try fallback: write to .new and let user move it
		_ = os.Remove(tmpNew)
		finalPath := filepath.Join(destDir, binName+".new")
		if err2 := copyFile(extractedPath, finalPath, 0o755); err2 != nil {
			return fmt.Errorf("failed to replace binary (%v) and to place .new binary (%v)", err, err2)
		}
		return fmt.Errorf("failed to replace running binary; placed new binary at %s", finalPath)
	}

	return nil
}

type winPendingReplace struct {
	NewPath string
}

func (w *winPendingReplace) Error() string {
	return fmt.Sprintf("new binary written to %s; please exit and replace the current binary manually", w.NewPath)
}

func inferPlatform() (osName, archName, format string) {
	osName = runtime.GOOS
	archName = runtime.GOARCH
	format = "tar.gz"
	if osName == "windows" {
		format = "zip"
	}
	// goreleaser uses darwin/linux/windows and amd64/arm64, which match runtime values
	return
}

func httpDownload(ctx context.Context, url, dest string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "truvami-decoder-selfupdate")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http status %d for %s", resp.StatusCode, url)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func verifyChecksum(archivePath, checksumPath, archiveName string) error {
	f, err := os.Open(checksumPath)
	if err != nil {
		return err
	}
	defer f.Close()

	var want string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		// typical format: "<sha256>  <filename>"
		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.HasSuffix(line, archiveName) {
			want = fields[0]
			break
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}
	if want == "" {
		return fmt.Errorf("checksum for %s not found in %s", archiveName, checksumPath)
	}

	got, err := fileSHA256(archivePath)
	if err != nil {
		return err
	}
	if !strings.EqualFold(want, got) {
		return fmt.Errorf("checksum mismatch: want %s got %s", want, got)
	}
	return nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func unzipSingleBinary(zipPath, binName, outPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// artifacts are wrapped in a directory, so match by leaf name
		if filepath.Base(f.Name) != binName {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		out, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, rc); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("binary %s not found in zip", binName)
}

func untarGzSingleBinary(tgzPath, binName, outPath string) error {
	f, err := os.Open(tgzPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.FileInfo().IsDir() {
			continue
		}
		if filepath.Base(hdr.Name) != binName {
			continue
		}
		out, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, tr); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("binary %s not found in archive", binName)
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	if mode != 0 {
		_ = os.Chmod(dst, mode)
	}
	return out.Sync()
}
