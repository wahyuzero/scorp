package agent

import (
	"path/filepath"
	"strings"
)

// IsDangerousCommand checks if a shell command line contains genuinely destructive operations.
// Avoids naive substring matching so innocent commands like 'echo "chmod 777"' or 'grep "rm -rf"' don't trigger.
func IsDangerousCommand(cmd string) bool {
	trimmed := strings.TrimSpace(cmd)
	if trimmed == "" {
		return false
	}

	// Check fork bomb before splitting
	if strings.Contains(trimmed, ":(){ :|:& };:") {
		return true
	}

	// Split compound commands by semicolon, newline, or logical AND/OR
	parts := strings.FieldsFunc(trimmed, func(r rune) bool {
		return r == ';' || r == '\n'
	})

	for _, part := range parts {
		subCmd := strings.TrimSpace(part)
		if subCmd == "" {
			continue
		}

		// Handle pipe or logical AND/OR
		subParts := strings.Split(subCmd, "&&")
		var allSubTokens []string
		for _, sp := range subParts {
			for _, orPart := range strings.Split(sp, "||") {
				for _, pipePart := range strings.Split(orPart, "|") {
					allSubTokens = append(allSubTokens, strings.TrimSpace(pipePart))
				}
			}
		}

		for _, tokenCmd := range allSubTokens {
			if tokenCmd == "" {
				continue
			}

			// Global device overwrite check — but redirection sinks like
			// /dev/null, /dev/stdout, /dev/stderr, /dev/tty are harmless and
			// extremely common (`2>/dev/null`), so they never confirm-gate.
			if devOverwriteTarget(tokenCmd) != "" && !isHarmlessDevSink(devOverwriteTarget(tokenCmd)) {
				return true
			}

			fields := strings.Fields(tokenCmd)
			if len(fields) == 0 {
				continue
			}

			bin := strings.ToLower(filepath.Base(fields[0]))

			// 1. Destructive File Deletion: rm
			if bin == "rm" {
				hasRecursive := false
				hasForce := false
				for _, rawArg := range fields[1:] {
					arg := strings.ToLower(rawArg)
					if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
						if strings.Contains(arg, "r") {
							hasRecursive = true
						}
						if strings.Contains(arg, "f") {
							hasForce = true
						}
					}
					if arg == "--recursive" {
						hasRecursive = true
					}
					if arg == "--force" {
						hasForce = true
					}

					if arg == "/" || arg == "/*" || strings.HasPrefix(arg, "/etc") || strings.HasPrefix(arg, "/boot") {
						return true
					}
				}
				if hasRecursive && hasForce {
					return true
				}
			}

			// 2. Dangerous filesystem & disk formats
			if bin == "mkfs" || strings.HasPrefix(bin, "mkfs.") || bin == "fdisk" || bin == "parted" {
				return true
			}

			// 3. Raw byte device overwrites (dd) — writing to /dev/null (and
			// other standard sinks) is harmless; only real devices gate.
			if bin == "dd" {
				for _, rawArg := range fields[1:] {
					arg := strings.ToLower(rawArg)
					if strings.HasPrefix(arg, "of=/dev/") {
						target := strings.TrimPrefix(arg, "of=/dev/")
						if !isHarmlessDevSink(target) {
							return true
						}
					}
				}
			}

			// 4. Insecure permission overrides (chmod 777)
			if bin == "chmod" {
				for _, rawArg := range fields[1:] {
					arg := strings.ToLower(rawArg)
					if arg == "777" || arg == "-r 777" || arg == "a+rwx" {
						return true
					}
				}
			}

			// 5. Mass process kills
			if bin == "killall" || bin == "pkill" {
				return true
			}
			if bin == "kill" {
				for _, rawArg := range fields[1:] {
					arg := strings.ToLower(rawArg)
					if arg == "-9" || arg == "-kill" {
						return true
					}
				}
			}

			// 6. System service teardown
			if bin == "systemctl" {
				for _, rawArg := range fields[1:] {
					arg := strings.ToLower(rawArg)
					if arg == "stop" || arg == "disable" || arg == "mask" {
						return true
					}
				}
			}

			// 7. Package and container removals
			if bin == "apt" {
				for _, rawArg := range fields[1:] {
					arg := strings.ToLower(rawArg)
					if arg == "remove" || arg == "purge" {
						return true
					}
				}
			}
			if bin == "pip" {
				for _, rawArg := range fields[1:] {
					arg := strings.ToLower(rawArg)
					if arg == "uninstall" {
						return true
					}
				}
			}
			if bin == "docker" {
				for j, rawArg := range fields[1:] {
					arg := strings.ToLower(rawArg)
					if arg == "rm" || arg == "rmi" {
						return true
					}
					if arg == "compose" && j+1 < len(fields[1:]) {
						nextArg := strings.ToLower(fields[1:][j+1])
						if nextArg == "down" {
							return true
						}
					}
				}
			}
			if bin == "docker-compose" {
				for _, rawArg := range fields[1:] {
					arg := strings.ToLower(rawArg)
					if arg == "down" {
						return true
					}
				}
			}

			// 8. Database destructive operations
			lowerToken := strings.ToLower(tokenCmd)
			if strings.HasPrefix(lowerToken, "drop database ") || strings.HasPrefix(lowerToken, "drop table ") ||
				strings.HasPrefix(lowerToken, "delete from ") || strings.HasPrefix(lowerToken, "truncate ") {
				return true
			}
		}
	}

	return false
}

// devOverwriteTarget extracts the device path a redirect writes to, if any
// (e.g. "2>/dev/null" → "null", "dd of=/dev/sda" is handled separately).
func devOverwriteTarget(token string) string {
	lower := strings.ToLower(token)
	for _, marker := range []string{"> /dev/", ">/dev/"} {
		if idx := strings.Index(lower, marker); idx >= 0 {
			rest := strings.TrimSpace(lower[idx+len(marker):])
			end := strings.IndexAny(rest, " \t\"';&)")
			if end >= 0 {
				rest = rest[:end]
			}
			if rest != "" {
				return rest
			}
		}
	}
	return ""
}

// isHarmlessDevSink reports whether the device is a standard stream sink.
func isHarmlessDevSink(name string) bool {
	switch name {
	case "null", "stdout", "stderr", "tty", "zero", "full":
		return true
	}
	return false
}
