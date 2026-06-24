package updater

// Version is the current build version. Overridden at build time via:
//   -ldflags="-X scorp-agent/updater.Version=v1.0.0"
var Version = "dev"

// GitHubRepo is the repo used for update checks. Override via
// GITHUB_REPO env var (format: "owner/repo").
var GitHubRepo = ""
