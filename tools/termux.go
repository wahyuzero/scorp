package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"scorp-agent/internal/helpers"
)

// ──────────────────────────────────────────────
// Deep Termux & Mobile Performance Optimization
// Direct Termux:API integration (WakeLock, Notification, Battery, Clipboard, Toast)
// ──────────────────────────────────────────────

// IsTermux detects if the agent is running in an Android Termux environment.
func IsTermux() bool {
	if prefix := os.Getenv("PREFIX"); strings.Contains(prefix, "com.termux") {
		return true
	}
	if _, err := os.Stat("/data/data/com.termux"); err == nil {
		return true
	}
	return false
}

// ExecuteTermuxAPI executes Termux API commands on Android devices.
func ExecuteTermuxAPI(args map[string]interface{}) (string, bool) {
	action := helpers.GetStringArg(args, "action", "")
	if action == "" {
		return "Error: 'action' is required (notification, toast, clipboard_get, clipboard_set, battery, vibrate, wake_lock, wake_unlock)", false
	}

	if !IsTermux() {
		// If running in development/desktop/VPS, provide clear non-fatal status
		switch action {
		case "notification", "toast":
			text := helpers.GetStringArg(args, "content", helpers.GetStringArg(args, "text", ""))
			return fmt.Sprintf("[Simulated Termux %s]: %s (not on Android Termux)", action, text), true
		case "battery":
			return `{"status": "AC", "percentage": 100, "health": "GOOD", "temperature": 25.0} (simulated desktop)`, true
		case "wake_lock", "wake_unlock":
			return fmt.Sprintf("[Simulated Termux %s] (no-op on desktop/VPS)", action), true
		default:
			return fmt.Sprintf("Termux action '%s' is only available inside Android Termux environment", action), false
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch action {
	case "notification":
		title := helpers.GetStringArg(args, "title", "Scorp Agent")
		content := helpers.GetStringArg(args, "content", "")
		cmd := exec.CommandContext(ctx, "termux-notification", "--title", title, "--content", content)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprintf("Error sending notification: %v (%s)", err, string(out)), false
		}
		return "✅ Termux notification sent", true

	case "toast":
		text := helpers.GetStringArg(args, "text", "")
		cmd := exec.CommandContext(ctx, "termux-toast", text)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprintf("Error displaying toast: %v (%s)", err, string(out)), false
		}
		return "✅ Termux toast displayed", true

	case "clipboard_get":
		cmd := exec.CommandContext(ctx, "termux-clipboard-get")
		out, err := cmd.Output()
		if err != nil {
			return fmt.Sprintf("Error reading clipboard: %v", err), false
		}
		return string(out), true

	case "clipboard_set":
		text := helpers.GetStringArg(args, "text", "")
		cmd := exec.CommandContext(ctx, "termux-clipboard-set", text)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprintf("Error writing clipboard: %v (%s)", err, string(out)), false
		}
		return "✅ Text copied to Android clipboard", true

	case "battery":
		cmd := exec.CommandContext(ctx, "termux-battery-status")
		out, err := cmd.Output()
		if err != nil {
			return fmt.Sprintf("Error reading battery status: %v", err), false
		}
		return string(out), true

	case "vibrate":
		duration := helpers.GetIntArg(args, "duration_ms", 200)
		cmd := exec.CommandContext(ctx, "termux-vibrate", "-d", fmt.Sprintf("%d", duration))
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprintf("Error vibrating device: %v (%s)", err, string(out)), false
		}
		return "✅ Haptic vibration triggered", true

	case "wake_lock":
		return AcquireTermuxWakeLock(), true

	case "wake_unlock":
		return ReleaseTermuxWakeLock(), true

	default:
		return fmt.Sprintf("Unknown Termux action: %s", action), false
	}
}

// AcquireTermuxWakeLock acquires a wake lock on Android Termux to prevent CPU sleep.
func AcquireTermuxWakeLock() string {
	if !IsTermux() {
		return "Wake lock skipped (not in Termux)"
	}
	cmd := exec.Command("termux-wake-lock")
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("Failed to acquire wake lock: %v (%s)", err, errBuf.String())
	}
	return "✅ Termux wake lock acquired (CPU will stay active)"
}

// ReleaseTermuxWakeLock releases the wake lock on Android Termux.
func ReleaseTermuxWakeLock() string {
	if !IsTermux() {
		return "Wake lock release skipped (not in Termux)"
	}
	cmd := exec.Command("termux-wake-unlock")
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("Failed to release wake lock: %v (%s)", err, errBuf.String())
	}
	return "✅ Termux wake lock released"
}

// SendTermuxNotification sends a native Android notification if running in Termux.
func SendTermuxNotification(title, content string) {
	if !IsTermux() {
		return
	}
	cmd := exec.Command("termux-notification", "--title", title, "--content", content)
	_ = cmd.Run()
}
