package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// Release represents a GitHub release.
type Release struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Body    string  `json:"body"`
	Assets  []Asset `json:"assets"`
}

// Asset is a downloadable file in a release.
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// client is a shared HTTP client with a reasonable timeout.
var client = &http.Client{Timeout: 30 * time.Second}

// getRepo returns the configured GitHub repo (owner/repo).
// Priority: GitHubRepo var → GITHUB_REPO env → empty.
func getRepo() string {
	if GitHubRepo != "" {
		return GitHubRepo
	}
	if r := os.Getenv("GITHUB_REPO"); r != "" {
		return r
	}
	return ""
}

// FetchLatestRelease queries the GitHub API for the latest release.
// Returns nil (no error) if no repo is configured.
func FetchLatestRelease() (*Release, error) {
	repo := getRepo()
	if repo == "" {
		return nil, nil
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	// Optional: use token for higher rate limits
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("no releases found")
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api %d: %s", resp.StatusCode, string(body))
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("parse release: %w", err)
	}
	return &rel, nil
}

// IsNewer compares versions. Returns true if remote > local.
// Handles semver (v1.2.3), short (v1.2), and "dev" (always newer).
func IsNewer(local, remote string) bool {
	if local == "dev" || local == "" {
		return true // dev builds should always offer update
	}
	local = strings.TrimPrefix(local, "v")
	remote = strings.TrimPrefix(remote, "v")
	if local == remote {
		return false
	}
	// Simple comparison: split by . and compare numerically
	lp := parseVer(local)
	rp := parseVer(remote)
	for i := 0; i < 3; i++ {
		if rp[i] > lp[i] {
			return true
		}
		if rp[i] < lp[i] {
			return false
		}
	}
	return false
}

func parseVer(v string) [3]int {
	var parts [3]int
	split := strings.Split(v, ".")
	for i, s := range split {
		if i >= 3 {
			break
		}
		// Extract leading digits
		n := 0
		for _, c := range s {
			if c >= '0' && c <= '9' {
				n = n*10 + int(c-'0')
			} else {
				break
			}
		}
		parts[i] = n
	}
	return parts
}

// CheckForUpdate fetches latest release and compares with current Version.
// Returns the release if an update is available, nil otherwise.
func CheckForUpdate() (*Release, error) {
	rel, err := FetchLatestRelease()
	if err != nil {
		return nil, err
	}
	if rel == nil {
		return nil, nil
	}
	if IsNewer(Version, rel.TagName) {
		return rel, nil
	}
	return nil, nil
}

// FindAssetForArch finds the release asset matching the current OS/arch.
// Naming convention: scorp_{os}_{arch} or scorp-{os}-{arch}
func FindAssetForArch(rel *Release) *Asset {
	goos := runtime.GOOS   // "linux", "darwin"
	goarch := runtime.GOARCH // "arm64", "amd64"

	for i, a := range rel.Assets {
		name := strings.ToLower(a.Name)
		if strings.Contains(name, goos) && strings.Contains(name, goarch) {
			return &rel.Assets[i]
		}
	}
	return nil
}

// isTermux checks if we're running inside Termux.
// Termux builds use Android's Bionic libc, which is incompatible with
// the glibc binaries produced by GitHub Actions CI.
func isTermux() bool {
	if prefix := os.Getenv("PREFIX"); strings.Contains(prefix, "com.termux") {
		return true
	}
	if _, err := os.Stat("/data/data/com.termux"); err == nil {
		return true
	}
	return false
}

// DownloadAsset downloads a release asset to a temp file.
// Returns the path to the downloaded file.
func DownloadAsset(asset *Asset) (string, error) {
	tmpFile, err := os.CreateTemp("", "scorp-update-*")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tmpFile.Close()

	resp, err := client.Get(asset.BrowserDownloadURL)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("download %d: %s", resp.StatusCode, resp.Status)
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("write: %w", err)
	}

	return tmpFile.Name(), nil
}
