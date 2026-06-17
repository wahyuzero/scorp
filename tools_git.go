package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// ──────────────────────────────────────────────
// Git Tool — structured git operations
// ──────────────────────────────────────────────

func executeGit(args map[string]interface{}, chatID int64) (string, bool) {
	action := getStringArg(args, "action", "")
	repo := getStringArg(args, "repo", ".")
	_ = repo // repo is used via -C flag

	if action == "" {
		return "Error: 'action' argument is required (status, log, diff, commit, branch, stash, pull, push)", false
	}

	timeout := getIntArg(args, "timeout", 30)
	if timeout > 120 {
		timeout = 120
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	switch action {
	case "status":
		cmd = exec.CommandContext(ctx, "git", "-C", repo, "status", "--short", "--branch")

	case "log":
		count := getIntArg(args, "count", 10)
		if count > 50 {
			count = 50
		}
		cmd = exec.CommandContext(ctx, "git", "-C", repo, "log", "--oneline",
			fmt.Sprintf("-%d", count), "--graph", "--decorate")

	case "diff":
		staged := getBoolArg(args, "staged", false)
		if staged {
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "diff", "--staged", "--stat")
		} else {
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "diff", "--stat")
		}

	case "branch":
		subAction := getStringArg(args, "sub_action", "list")
		switch subAction {
		case "list":
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "branch", "-a", "-v")
		case "create":
			name := getStringArg(args, "name", "")
			if name == "" {
				return "Error: 'name' argument required for branch create", false
			}
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "checkout", "-b", name)
		case "switch":
			name := getStringArg(args, "name", "")
			if name == "" {
				return "Error: 'name' argument required for branch switch", false
			}
			cmd = exec.CommandContext(ctx, "git", "-C", repo, "checkout", name)
		default:
			return "Error: unknown branch sub_action: " + subAction, false
		}

	case "stash":
		subAction := getStringArg(args, "sub_action", "push")
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
		remote := getStringArg(args, "remote", "origin")
		branch := getStringArg(args, "branch", "")
		force := getBoolArg(args, "force", false)

		if force {
			chatIDStr := fmt.Sprintf("%d", chatID)
			storePendingConfirmation(chatIDStr, "git", fmt.Sprintf("git push --force %s %s (repo: %s)", remote, branch, repo), nil)
			return "⚠️ FORCE PUSH requires confirmation.\nPlease confirm to proceed.", false
		}

		pushArgs := []string{"-C", repo, "push", remote}
		if branch != "" {
			pushArgs = append(pushArgs, branch)
		}
		cmd = exec.CommandContext(ctx, "git", pushArgs...)

	case "commit":
		message := getStringArg(args, "message", "")
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
		ref := getStringArg(args, "ref", "HEAD")
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
		return fmt.Sprintf("Git %s failed: %v\n%s", action, err, truncOutput(result, maxToolOutput)), false
	}

	if result == "" {
		return fmt.Sprintf("git %s: OK (no output)", action), true
	}

	return truncOutput(result, maxToolOutput), true
}
