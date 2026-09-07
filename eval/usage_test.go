package eval

import (
	"os"
	"path/filepath"
	"testing"
)

// TestUsageDeltaClampsNegative pins the calibration contract: model_usage.json
// is cumulative, so a reset mid-run must clamp to zero, and the delta is the
// per-task usage.
func TestUsageDeltaClampsNegative(t *testing.T) {
	before := Usage{In: 1000, Cached: 5000, Out: 200, Calls: 10}
	after := Usage{In: 1500, Cached: 4000, Out: 350, Calls: 13} // cached reset mid-run
	d := usageDelta(before, after)
	if d.In != 500 || d.Out != 150 || d.Calls != 3 {
		t.Fatalf("delta wrong: %+v", d)
	}
	if d.Cached != 0 {
		t.Fatalf("negative cached drift must clamp to 0, got %d", d.Cached)
	}
	if d.Fresh() != 650 || d.Total() != 650 {
		t.Fatalf("fresh/total wrong: %+v", d)
	}
}

// TestUsageSnapshotParsesModelUsage pins the field mapping: input/cached/
// output are summed per model but kept separate; unrelated fields are ignored.
func TestUsageSnapshotParsesModelUsage(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	scorpDir := filepath.Join(home, ".scorp")
	if err := os.MkdirAll(scorpDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `{
	  "m1": {"model": "m1", "input_tokens": 1200, "cached_tokens": 800, "output_tokens": 90, "calls": 3, "last_used": "x"},
	  "m2": {"model": "m2", "input_tokens": 100,  "cached_tokens": 0,   "output_tokens": 10,  "calls": 1, "last_used": "y"}
	}`
	if err := os.WriteFile(filepath.Join(scorpDir, "model_usage.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	u := usageSnapshot()
	if u.In != 1300 || u.Cached != 800 || u.Out != 100 || u.Calls != 4 {
		t.Fatalf("snapshot wrong: %+v", u)
	}
	if u.Fresh() != 1400 || u.Total() != 2200 {
		t.Fatalf("fresh/total wrong: fresh=%d total=%d", u.Fresh(), u.Total())
	}

	// Missing file → zero snapshot (live case then reports zero delta).
	t.Setenv("HOME", filepath.Join(home, "nonexistent"))
	if got := usageSnapshot(); got != (Usage{}) {
		t.Fatalf("missing file must yield zero snapshot, got %+v", got)
	}
}
