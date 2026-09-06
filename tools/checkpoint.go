package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Checkpoint / Rewind (P1.6)
//
// "FINALLY checkpoints!" — Cursor & Claude Code made it baseline UX. Design:
// before each tool-bearing turn, the agent records a SHADOW COMMIT of the
// working repo using git plumbing (temp GIT_INDEX_FILE + add -A + write-tree
// + commit-tree) under refs/scorp/ckpt/<session>/<nanos> — the user's index,
// working tree, and history are never touched. /undo restores the latest
// pre-turn state and walks back one checkpoint per invocation, capped at 20
// per session. Non-git working directories are simply skipped (no-op).
// ──────────────────────────────────────────────

const (
	ckptRefPrefix     = "refs/scorp/ckpt/"
	ckptMaxPerSession = 20
	ckptGitTimeout    = 15 * time.Second
)

var ckptRefSanitizer = regexp.MustCompile(`[^A-Za-z0-9._-]`)

// CheckpointInfo describes one recorded checkpoint ref.
type CheckpointInfo struct {
	Ref      string
	RepoRoot string
	Time     string
	Subject  string
}

// ckptGit runs one git command with a hard timeout. Env entries (e.g.
// GIT_INDEX_FILE, author ident) are appended to the process environment.
func ckptGit(dir string, env []string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ckptGitTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func checkpointRepoRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	out, err := ckptGit(wd, nil, "rev-parse", "--show-toplevel")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

func ckptRefFor(sessionID string, ts time.Time) string {
	safe := ckptRefSanitizer.ReplaceAllString(strings.TrimSpace(sessionID), "_")
	if len(safe) > 30 {
		safe = safe[len(safe)-30:]
	}
	if safe == "" {
		safe = "default"
	}
	return fmt.Sprintf("%s%s/%020d", ckptRefPrefix, safe, ts.UnixNano())
}

// CreateCheckpoint records the repo's current working state as a shadow
// commit. It runs even when the tree is clean so the undo chain stays
// complete (every turn has a pre-state). Returns created=false when the
// working directory is not a git repo.
func CreateCheckpoint(sessionID string) (CheckpointInfo, bool, error) {
	root := checkpointRepoRoot()
	if root == "" {
		return CheckpointInfo{}, false, nil
	}
	ts := time.Now().UTC()
	ref := ckptRefFor(sessionID, ts)

	tmpIdx, err := os.CreateTemp("", "scorp-ckpt-idx-")
	if err != nil {
		return CheckpointInfo{}, false, err
	}
	tmpIdx.Close()
	os.Remove(tmpIdx.Name()) // git must create the index file itself
	defer os.Remove(tmpIdx.Name())
	idxEnv := []string{"GIT_INDEX_FILE=" + tmpIdx.Name()}

	if _, err := ckptGit(root, idxEnv, "add", "-A"); err != nil {
		return CheckpointInfo{}, false, fmt.Errorf("git add: %v", err)
	}
	treeOut, err := ckptGit(root, idxEnv, "write-tree")
	if err != nil {
		return CheckpointInfo{}, false, fmt.Errorf("git write-tree: %v", err)
	}
	tree := strings.TrimSpace(treeOut)

	ident := []string{
		"GIT_AUTHOR_NAME=scorp", "GIT_AUTHOR_EMAIL=scorp@checkpoint",
		"GIT_COMMITTER_NAME=scorp", "GIT_COMMITTER_EMAIL=scorp@checkpoint",
	}
	msg := fmt.Sprintf("checkpoint: %s | %s", sessionID, ts.Format(time.RFC3339))
	shaOut, err := ckptGit(root, ident, "commit-tree", tree, "-m", msg)
	if err != nil {
		return CheckpointInfo{}, false, fmt.Errorf("git commit-tree: %v", err)
	}
	sha := strings.TrimSpace(shaOut)
	if _, err := ckptGit(root, nil, "update-ref", ref, sha); err != nil {
		return CheckpointInfo{}, false, fmt.Errorf("git update-ref: %v", err)
	}

	pruneCheckpoints(sessionID, root)
	return CheckpointInfo{Ref: ref, RepoRoot: root, Time: ts.Format(time.RFC3339), Subject: msg}, true, nil
}

// ListCheckpoints returns the session's checkpoints, newest first.
func ListCheckpoints(sessionID string) ([]CheckpointInfo, error) {
	root := checkpointRepoRoot()
	if root == "" {
		return nil, nil
	}
	safe := ckptRefSanitizer.ReplaceAllString(strings.TrimSpace(sessionID), "_")
	pattern := ckptRefPrefix + safe + "/*"
	out, err := ckptGit(root, nil, "for-each-ref", "--sort=-refname",
		"--format=%(refname)%09%(committerdate:iso8601)%09%(contents:subject)", pattern)
	if err != nil {
		return nil, err
	}
	var infos []CheckpointInfo
	for _, line := range strings.Split(out, "\n") {
		parts := strings.Split(strings.TrimRight(line, "\r"), "\t")
		if len(parts) < 3 || !strings.HasPrefix(parts[0], ckptRefPrefix) {
			continue
		}
		infos = append(infos, CheckpointInfo{Ref: parts[0], RepoRoot: root, Time: parts[1], Subject: parts[2]})
	}
	return infos, nil
}

// RestoreCheckpoint materializes a checkpoint's tree into the working
// directory (overlay: files from the checkpoint are written; files created
// after it remain). Returns the number of files in the checkpoint tree.
func RestoreCheckpoint(sessionID, ref string) (int, error) {
	root := checkpointRepoRoot()
	if root == "" {
		return 0, fmt.Errorf("not a git repository")
	}
	if !strings.HasPrefix(ref, ckptRefPrefix) {
		return 0, fmt.Errorf("ref is not a scorp checkpoint")
	}
	tmpIdx, err := os.CreateTemp("", "scorp-ckpt-idx-")
	if err != nil {
		return 0, err
	}
	tmpIdx.Close()
	os.Remove(tmpIdx.Name()) // git must create the index file itself
	defer os.Remove(tmpIdx.Name())
	idxEnv := []string{"GIT_INDEX_FILE=" + tmpIdx.Name()}

	if _, err := ckptGit(root, idxEnv, "read-tree", ref); err != nil {
		return 0, fmt.Errorf("git read-tree: %v", err)
	}
	if _, err := ckptGit(root, idxEnv, "checkout-index", "-a", "-f"); err != nil {
		return 0, fmt.Errorf("git checkout-index: %v", err)
	}
	lsOut, err := ckptGit(root, nil, "ls-tree", "-r", "--name-only", ref)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, l := range strings.Split(lsOut, "\n") {
		if strings.TrimSpace(l) != "" {
			n++
		}
	}
	return n, nil
}

// DeleteCheckpoint drops one checkpoint ref (used after a successful undo so
// repeated /undo walks back through the chain).
func DeleteCheckpoint(_ string, ref string) error {
	root := checkpointRepoRoot()
	if root == "" {
		return fmt.Errorf("not a git repository")
	}
	if !strings.HasPrefix(ref, ckptRefPrefix) {
		return fmt.Errorf("ref is not a scorp checkpoint")
	}
	_, err := ckptGit(root, nil, "update-ref", "-d", ref)
	return err
}

// CheckpointDiffStat summarizes what changed in the working tree since a
// checkpoint (tracked files only).
func CheckpointDiffStat(ref string) string {
	root := checkpointRepoRoot()
	if root == "" {
		return ""
	}
	out, err := ckptGit(root, nil, "diff", "--stat", ref)
	if err != nil {
		return ""
	}
	out = strings.TrimSpace(out)
	if len(out) > 600 {
		out = out[:600] + "\n..."
	}
	return out
}

func pruneCheckpoints(sessionID, root string) {
	refs, err := listCheckpointRefs(sessionID, root)
	if err != nil || len(refs) <= ckptMaxPerSession {
		return
	}
	// refs sorted newest-first — drop the oldest beyond the cap
	for _, ref := range refs[ckptMaxPerSession:] {
		_, _ = ckptGit(root, nil, "update-ref", "-d", ref)
	}
}

func listCheckpointRefs(sessionID, root string) ([]string, error) {
	safe := ckptRefSanitizer.ReplaceAllString(strings.TrimSpace(sessionID), "_")
	out, err := ckptGit(root, nil, "for-each-ref", "--sort=-refname", "--format=%(refname)", ckptRefPrefix+safe+"/*")
	if err != nil {
		return nil, err
	}
	var refs []string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, ckptRefPrefix) {
			refs = append(refs, strings.TrimSpace(line))
		}
	}
	return refs, nil
}
