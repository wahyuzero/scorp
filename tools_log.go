package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// ──────────────────────────────────────────────
// Log Tail Tool — docker logs, journalctl, file tail
// ──────────────────────────────────────────────

func executeLog(args map[string]interface{}) (string, bool) {
	source := getStringArg(args, "source", "")
	target := getStringArg(args, "target", "")
	lines := getIntArg(args, "lines", 50)
	follow := getBoolArg(args, "follow", false)
	duration := getIntArg(args, "duration", 10)

	if source == "" {
		return "Error: 'source' is required (docker, journal, file)", false
	}
	if target == "" {
		return "Error: 'target' is required (container name, unit name, or file path)", false
	}

	if lines > 500 {
		lines = 500
	}
	if duration > 30 {
		duration = 30
	}

	var cmd *exec.Cmd
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration+5)*time.Second)
	defer cancel()

	switch source {
	case "docker":
		args := []string{"logs", "--tail", fmt.Sprintf("%d", lines)}
		if follow {
			args = append(args, "-f")
		}
		args = append(args, target)
		cmd = exec.CommandContext(ctx, "docker", args...)

	case "journal":
		jArgs := []string{"-u", target, "--no-pager"}
		if follow {
			jArgs = append(jArgs, "-f")
		} else {
			jArgs = append(jArgs, "-n", fmt.Sprintf("%d", lines))
		}
		cmd = exec.CommandContext(ctx, "journalctl", jArgs...)

	case "file":
		tailArgs := []string{"-n", fmt.Sprintf("%d", lines)}
		if follow {
			tailArgs = append(tailArgs, "-f")
		}
		tailArgs = append(tailArgs, target)
		cmd = exec.CommandContext(ctx, "tail", tailArgs...)

	default:
		return fmt.Sprintf("Error: unknown source '%s'. Use: docker, journal, file", source), false
	}

	output, err := cmd.CombinedOutput()
	result := string(output)

	if ctx.Err() == context.DeadlineExceeded && follow {
		// Follow mode timeout is expected — return what we got
		return truncOutput(fmt.Sprintf("📋 %s logs for '%s' (followed %ds):\n\n%s", source, target, duration, result), maxToolOutput), true
	}

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Sprintf("Log fetch timed out after %ds", duration+5), false
	}

	if err != nil && !follow {
		return fmt.Sprintf("Log fetch failed: %v\n%s", err, truncOutput(result, maxToolOutput)), false
	}

	header := fmt.Sprintf("📋 %s logs for '%s' (last %d lines):\n\n", source, target, lines)
	return truncOutput(header+result, maxToolOutput), true
}
