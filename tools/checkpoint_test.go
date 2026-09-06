package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupCkptRepo makes a temp git repo and chdirs the test into it (t.Chdir
// restores the previous directory afterwards).
func setupCkptRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Chdir(dir)
	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v (%s)", strings.Join(args, " "), err, out)
		}
	}
	run("init", "-q")
	run("config", "user.email", "test@scorp")
	run("config", "user.name", "test")
	if err := os.WriteFile(filepath.Join(dir, "seed.txt"), []byte("v1\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", "-A")
	run("commit", "-q", "-m", "seed")
	return dir
}

func writeFile(t *testing.T, name, content string) {
	t.Helper()
	wd, _ := os.Getwd()
	if err := os.WriteFile(filepath.Join(wd, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestCheckpointLifecycle(t *testing.T) {
	dir := setupCkptRepo(t)

	// checkpoint 1: pre-turn state with seed.txt = v1
	ck1, created, err := CreateCheckpoint("sess-ckpt-test")
	if err != nil || !created {
		t.Fatalf("checkpoint 1 not created: created=%v err=%v", created, err)
	}

	// turn 1 modifies the file
	writeFile(t, "seed.txt", "v2\n")

	// checkpoint 2: pre-turn-2 state (still v2)
	ck2, created, err := CreateCheckpoint("sess-ckpt-test")
	if err != nil || !created {
		t.Fatalf("checkpoint 2 not created: created=%v err=%v", created, err)
	}
	if ck2.Ref == ck1.Ref {
		t.Fatal("checkpoints must have distinct refs")
	}

	// turn 2 modifies again
	writeFile(t, "seed.txt", "v3\n")

	cks, err := ListCheckpoints("sess-ckpt-test")
	if err != nil || len(cks) != 2 {
		t.Fatalf("expected 2 checkpoints, got %d (err=%v)", len(cks), err)
	}
	if cks[0].Ref != ck2.Ref {
		t.Fatal("checkpoints must be listed newest-first")
	}

	// /undo: restore the LATEST checkpoint = pre-turn-2 state (v2)
	n, err := RestoreCheckpoint("sess-ckpt-test", cks[0].Ref)
	if err != nil {
		t.Fatalf("restore failed: %v", err)
	}
	if n < 1 {
		t.Fatalf("expected at least 1 restored file, got %d", n)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "seed.txt"))
	if string(got) != "v2\n" {
		t.Fatalf("undo must restore pre-turn state v2, got %q", got)
	}

	if err := DeleteCheckpoint("sess-ckpt-test", cks[0].Ref); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	cks, _ = ListCheckpoints("sess-ckpt-test")
	if len(cks) != 1 {
		t.Fatalf("expected 1 checkpoint after delete, got %d", len(cks))
	}
}

func TestCheckpointRejectsForeignRef(t *testing.T) {
	setupCkptRepo(t)
	if _, err := RestoreCheckpoint("s", "refs/heads/main"); err == nil {
		t.Fatal("restoring a non-checkpoint ref must be rejected")
	}
	if err := DeleteCheckpoint("s", "refs/heads/main"); err == nil {
		t.Fatal("deleting a non-checkpoint ref must be rejected")
	}
}

func TestCheckpointNonRepoNoop(t *testing.T) {
	t.Chdir(t.TempDir())
	_, created, err := CreateCheckpoint("whatever")
	if err != nil {
		t.Fatalf("non-repo checkpoint must not error, got %v", err)
	}
	if created {
		t.Fatal("non-repo checkpoint must report created=false")
	}
}

func TestCheckpointPruneCap(t *testing.T) {
	setupCkptRepo(t)
	for i := 0; i < ckptMaxPerSession+5; i++ {
		writeFile(t, "seed.txt", strings.Repeat("x", 1)+string(rune('a'+i%26))+"\n")
		if _, created, err := CreateCheckpoint("prune-sess"); err != nil || !created {
			t.Fatalf("checkpoint %d failed: %v", i, err)
		}
	}
	refs, err := listCheckpointRefs("prune-sess", checkpointRepoRoot())
	if err != nil {
		t.Fatal(err)
	}
	if len(refs) > ckptMaxPerSession {
		t.Fatalf("cap exceeded: %d refs", len(refs))
	}
}
