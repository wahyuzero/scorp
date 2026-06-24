package updater

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// SelfUpdate performs a full self-update:
// 1. Check GitHub releases for a newer version
// 2. Download the matching binary for this OS/arch
// 3. Replace the running binary atomically
// 4. Restart the service (if running under systemd) or exit for manual restart
//
// Returns a human-readable message describing the outcome.
func SelfUpdate() (string, error) {
	// Step 1: Check for update
	rel, err := CheckForUpdate()
	if err != nil {
		return "", fmt.Errorf("check failed: %w", err)
	}
	if rel == nil {
		msg := fmt.Sprintf("Already up to date (%s)", Version)
		return msg, nil
	}

	// Step 2: Find matching asset
	asset := FindAssetForArch(rel)
	if asset == nil {
		// No prebuilt binary — try building from source
		return fmt.Sprintf("Update %s available, but no prebuilt binary for %s/%s.\n"+
			"Build manually: git pull && make", rel.TagName, runtime.GOOS, runtime.GOARCH), nil
	}

	// Termux check: CI binaries use glibc, Termux uses Bionic — incompatible
	if isTermux() {
		return fmt.Sprintf("🔔 Update %s available.\n"+
			"Termux detected — prebuilt binaries are incompatible (glibc vs Bionic).\n"+
			"Update manually:\n"+
			"  cd ~/scorp && git pull && ./install.sh", rel.TagName), nil
	}

	// Step 3: Download
	fmt.Printf("Downloading %s (%s)...\n", asset.Name, humanSize(asset.Size))
	tmpPath, err := DownloadAsset(asset)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpPath)

	// Step 4: Find current binary path
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("find executable: %w", err)
	}
	exePath, _ = filepath.EvalSymlinks(exePath)

	// Step 5: Atomic replace
	// Write downloaded binary to a temp file next to the target,
	// then rename (atomic on same filesystem)
	dir := filepath.Dir(exePath)
	newPath := filepath.Join(dir, ".scorp.new")
	if err := copyFile(tmpPath, newPath); err != nil {
		return "", fmt.Errorf("stage new binary: %w", err)
	}

	if err := os.Chmod(newPath, 0755); err != nil {
		os.Remove(newPath)
		return "", fmt.Errorf("chmod: %w", err)
	}

	// Backup current binary
	backup := exePath + ".bak"
	os.Remove(backup) // remove stale backup
	os.Rename(exePath, backup)

	if err := os.Rename(newPath, exePath); err != nil {
		// Restore backup
		os.Rename(backup, exePath)
		return "", fmt.Errorf("replace binary: %w", err)
	}

	// Keep backup as safety net — don't remove automatically.

	// Step 6: Restart
	msg := fmt.Sprintf("✅ Updated %s → %s", Version, rel.TagName)

	// Note: backup left at exePath.bak as safety net.
	// It will be overwritten on next update or can be removed manually.

	if runningUnderSystemd() {
		// Restart via systemctl
		cmd := exec.Command("systemctl", "restart", "scorp")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			msg += "\n⚠️ Binary updated but auto-restart failed. Run: sudo systemctl restart scorp"
		} else {
			msg += "\n🔄 Service restarted."
		}
	} else {
		msg += "\n💡 Restart scorp to use the new version."
	}

	return msg, nil
}

// CheckAndNotify is called on startup (async). It checks for updates
// and returns a notification message if one is available.
// Returns ("", nil) if no update or no repo configured.
func CheckAndNotify() (string, error) {
	rel, err := CheckForUpdate()
	if err != nil {
		return "", err
	}
	if rel == nil {
		return "", nil
	}

	asset := FindAssetForArch(rel)
	archNote := ""
	if asset != nil {
		archNote = fmt.Sprintf(" (%s available)", asset.Name)
	}

	changelog := ""
	if rel.Body != "" {
		// Trim changelog to first 500 chars
		cl := rel.Body
		if len(cl) > 500 {
			cl = cl[:500] + "..."
		}
		changelog = "\n\n" + cl
	}

	msg := fmt.Sprintf("🔔 Update available: %s%s\n"+
		"Current: %s\n"+
		"Run /update to install.%s",
		rel.TagName, archNote, Version, changelog)

	return msg, nil
}

// runningUnderSystemd checks if the process is managed by systemd.
func runningUnderSystemd() bool {
	// Check if we have a systemd scope
	if info, err := os.Stat("/run/systemd/system"); err == nil && info.IsDir() {
		// Check if our parent is systemd
		ppid := os.Getppid()
		if ppid == 1 {
			return true
		}
		// Check if the service name is in systemctl
		cmd := exec.Command("systemctl", "is-active", "scorp")
		out, _ := cmd.Output()
		if string(out) == "active\n" || string(out) == "active" {
			return true
		}
	}
	return false
}

// copyFile copies src to dst.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// humanSize formats bytes as a readable string.
func humanSize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%dB", n)
	}
	if n < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(n)/1024)
	}
	return fmt.Sprintf("%.1fMB", float64(n)/(1024*1024))
}

// runtime_GOOS and runtime_GOARCH wrappers removed — use runtime package directly.
