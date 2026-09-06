package tools

import (
	"context"
	"fmt"
	"os/exec"
	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Git Tool — structured git operations
// ──────────────────────────────────────────────

// ExecuteGit runs git operations
func ExecuteGit(args map[string]interface{}, chatID int64) (string, bool) {
	action := helpers.GetStringArg(args, "action", "")
	if action == "" {
		action = helpers.GetStringArg(args, "command", "")
	}
	repo := helpers.GetStringArg(args, "repo", ".")

	// Parse flags and formats like "-C /path status" or "git status"
	fields := strings.Fields(action)
	if len(fields) > 0 {
		if fields[0] == "git" {
			fields = fields[1:]
		}
		for len(fields) >= 2 && fields[0] == "-C" {
			repo = fields[1]
			fields = fields[2:]
		}
		if len(fields) > 0 {
			action = fields[0]
		}
	}

	if action == "" {
		return "Error: 'action' argument is required (status, log, diff, commit, branch, stash, pull, push)", false
	}

	timeout := helpers.GetIntArg(args, "timeout", 30)
	if timeout > 120 {
		timeout = 120
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	// If complex compound arguments were passed for safe read-only operations, pass directly
	if len(fields) > 1 {
		subCmd := fields[0]
		if subCmd == "log" || subCmd == "status" || subCmd == "diff" || subCmd == "show" || subCmd == "branch" || subCmd == "remote" {
			gitArgs := append([]string{"-C", repo}, fields...)
			cmd = exec.CommandContext(ctx, "git", gitArgs...)
			output, err := cmd.CombinedOutput()
			result := string(output)
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Sprintf("Git command timed out after %ds", timeout), false
			}
			if err != nil {
				return fmt.Sprintf("Git command failed: %v\n%s", err, helpers.TruncOutput(result, helpers.MaxToolOutput)), false
			}
			return helpers.TruncOutput(result, helpers.MaxToolOutput), true
		}
	}

	switch action {
	case "status":
		cmd = exec.CommandContext(ctx, "git", "-C", repo, "status", "--short", "--branch")

	case "log":
		count := helpers.GetIntArg(args, "count", 10)
		if count > 50 {
			count = 50
		}
		cmd = exec.CommandContext(ctx, "git", "-C", repo, "log", "--oneline",
			fmt.Sprintf("-%d", count), "--graph", "--decorate")

	case "diff":
		staged := helpers.GetBoolArg(args, "staged", false)
		if staged {
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "diff", "--staged", "--stat")
		} else {
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "diff", "--stat")
		}

	case "branch":
		subAction := helpers.GetStringArg(args, "sub_action", "list")
		switch subAction {
		case "list":
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "branch", "-a", "-v")
		case "create":
			name := helpers.GetStringArg(args, "name", "")
			if name == "" {
				return "Error: 'name' argument required for branch create", false
			}
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "checkout", "-b", name)
		case "switch":
			name := helpers.GetStringArg(args, "name", "")
			if name == "" {
				return "Error: 'name' argument required for branch switch", false
			}
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "checkout", name)
		default:
			return "Error: unknown branch sub_action: " + subAction, false
		}

	case "stash":
		subAction := helpers.GetStringArg(args, "sub_action", "push")
		switch subAction {
		case "push":
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "stash")
		case "pop":
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "stash", "pop")
		case "list":
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "stash", "list")
		default:
			return "Error: unknown stash sub_action: " + subAction, false
		}

	case "pull":
		cmd = exec.CommandContext(ctx, "git", "-C", repo, "pull", "--ff-only")

	case "push":
		// Push requires confirmation — always treated as potentially dangerous
		remote := helpers.GetStringArg(args, "remote", "origin")
		branch := helpers.GetStringArg(args, "branch", "")
		force := helpers.GetBoolArg(args, "force", false)

		if force && config.ConfirmationRequired() {
			chatIDStr := fmt.Sprintf("%d", chatID)
			StorePendingConfirmation(chatIDStr, "git", fmt.Sprintf("git push --force %s %s (repo: %s)", remote, branch, repo), nil)
			return "⚠️ FORCE PUSH requires confirmation.\nPlease confirm to proceed.", false
		}

		pushArgs := []string{"-C", repo, "push", remote}
		if force {
			pushArgs = append(pushArgs, "--force")
		}
		if branch != "" {
			pushArgs = append(pushArgs, branch)
		}
		cmd = exec.CommandContext(ctx, "git", pushArgs...)

	case "commit":
		message := helpers.GetStringArg(args, "message", "")
		if message == "" {
			return "Error: 'message' argument required for commit", false
		}
		// Stage all + commit
		stageCmd := exec.CommandContext(ctx, "git", "-C", repo, "add", "-A")
		stageOut, stageErr := stageCmd.CombinedOutput()
		if stageErr != nil {
			return fmt.Sprintf("git add failed: %v\n%s", stageErr, string(stageOut)), false
		}
		cmd = exec.CommandContext(ctx, "git", "-C", repo, "commit", "-m", message)

	case "show":
		ref := helpers.GetStringArg(args, "ref", "HEAD")
		cmd = exec.CommandContext(ctx, "git", "-C", repo, "show", "--stat", ref)

	case "remote":
		cmd = exec.CommandContext(ctx, "git", "-C", repo, "remote", "-v")

	default:
		return fmt.Sprintf("Error: unknown git action '%s'. Available: status, log, diff, commit, branch, stash, pull, push, show, remote", action), false
	}

	output, err := cmd.CombinedOutput()
	result := string(output)

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Sprintf("Git command timed out after %ds", timeout), false
	}

	if err != nil {
		return fmt.Sprintf("Git %s failed: %v\n%s", action, err, helpers.TruncOutput(result, helpers.MaxToolOutput)), false
	}

	if result == "" {
		return fmt.Sprintf("git %s: OK (no output)", action), true
	}

	return helpers.TruncOutput(result, helpers.MaxToolOutput), true
}
