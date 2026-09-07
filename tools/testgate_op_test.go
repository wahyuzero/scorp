package tools

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// seedOpReceipts installs receipts into the gate window and restores prior
// state afterwards.
func seedOpReceipts(t *testing.T, receipts []ToolReceipt) {
	t.Helper()
	savedLoaded := receiptsLoaded
	savedReceipts := recentReceipts
	savedBoundary := testGateBoundary
	receiptsLoaded = true
	recentReceipts = receipts
	testGateBoundary = time.Now().Add(-time.Hour)
	t.Cleanup(func() {
		receiptsLoaded = savedLoaded
		recentReceipts = savedReceipts
		testGateBoundary = savedBoundary
	})
}

func opReceipt(tool string, success bool, meta map[string]string) ToolReceipt {
	if meta == nil {
		meta = map[string]string{}
	}
	return ToolReceipt{Tool: tool, Success: success, Timestamp: time.Now(), Meta: meta}
}

func TestLooksLikeOperationalClaims(t *testing.T) {
	cases := []struct {
		text      string
		wantCount int
		wantClass string
		wantObj   string
	}{
		{"Task t1 deleted.", 1, "delete", "t1"},
		{"The file `/tmp/a/c14v3-runs.txt` was removed successfully.", 1, "delete", "/tmp/a/c14v3-runs.txt"},
		{"File /tmp/x.txt dihapus.", 1, "delete", "/tmp/x.txt"},
		{"Task created: ID t1.", 1, "create", "t1"},
		{"Config file /etc/app.conf has been written.", 1, "create", "/etc/app.conf"},
		{"The service scorp has been restarted.", 1, "lifecycle", "scorp"},
		{"The bug was fixed in `main.go`.", 1, "lifecycle", "main.go"},
	}
	for _, tc := range cases {
		got := LooksLikeOperationalClaims(tc.text)
		if len(got) != tc.wantCount {
			t.Errorf("%q: want %d claims, got %d (%v)", tc.text, tc.wantCount, len(got), got)
			continue
		}
		if got[0].Class != tc.wantClass {
			t.Errorf("%q: class=%s want %s", tc.text, got[0].Class, tc.wantClass)
		}
		if !strings.Contains(got[0].Object, tc.wantObj) && !strings.Contains(tc.wantObj, got[0].Object) {
			t.Errorf("%q: object=%q want containing %q", tc.text, got[0].Object, tc.wantObj)
		}
	}
}

func TestLooksLikeOperationalClaimsSkipsVague(t *testing.T) {
	vague := []string{
		"The report was created and everything looks good.",   // no concrete object... "report" not matched
		"I have inspected the results.",                       // no op verb
		"The deletion process is designed to be safe.",         // present-tense generic
		"catatan: jangan menghapus file apa pun.",              // future/imperative
	}
	for _, s := range vague {
		if got := LooksLikeOperationalClaims(s); len(got) != 0 {
			t.Errorf("vague text %q must not produce claims, got %v", s, got)
		}
	}
}

func TestUnverifiedOperationalClaims(t *testing.T) {
	t.Run("delete backed by successful rm receipt", func(t *testing.T) {
		seedOpReceipts(t, []ToolReceipt{
			opReceipt("shell", true, map[string]string{"cmd": "rm -rf /tmp/c14v3-runs.txt"}),
		})
		got := UnverifiedOperationalClaims("The file /tmp/c14v3-runs.txt was removed successfully.")
		if len(got) != 0 {
			t.Fatalf("rm receipt must back the delete claim, got %v", got)
		}
	})

	t.Run("delete claim with only a read receipt is unverified", func(t *testing.T) {
		seedOpReceipts(t, []ToolReceipt{
			opReceipt("shell", true, map[string]string{"cmd": "cat /tmp/c14v3-runs.txt"}),
		})
		got := UnverifiedOperationalClaims("The file /tmp/c14v3-runs.txt was deleted.")
		if len(got) != 1 {
			t.Fatalf("read-only receipt must NOT back a delete claim, got %v", got)
		}
	})

	t.Run("failed rm does not back the claim", func(t *testing.T) {
		seedOpReceipts(t, []ToolReceipt{
			opReceipt("shell", false, map[string]string{"cmd": "rm -rf /tmp/still-here.txt"}),
		})
		got := UnverifiedOperationalClaims("/tmp/still-here.txt was deleted.")
		if len(got) != 1 {
			t.Fatalf("failed rm must not back a delete claim, got %v", got)
		}
	})

	t.Run("schedule_manage delete backed via action+id meta", func(t *testing.T) {
		seedOpReceipts(t, []ToolReceipt{
			opReceipt("schedule_manage", true, map[string]string{"action": "delete", "id": "t1"}),
		})
		got := UnverifiedOperationalClaims("Task t1 deleted.")
		if len(got) != 0 {
			t.Fatalf("schedule_manage delete receipt must back the claim, got %v", got)
		}
	})

	t.Run("create backed by write_file path receipt", func(t *testing.T) {
		seedOpReceipts(t, []ToolReceipt{
			opReceipt("write_file", true, map[string]string{"path": "/tmp/done.txt"}),
		})
		got := UnverifiedOperationalClaims("File /tmp/done.txt was created.")
		if len(got) != 0 {
			t.Fatalf("write_file receipt must back the create claim, got %v", got)
		}
	})

	t.Run("no receipts at all → unverified", func(t *testing.T) {
		seedOpReceipts(t, nil)
		got := UnverifiedOperationalClaims("Task t1 deleted.")
		if len(got) != 1 {
			t.Fatalf("no receipts must leave the claim unverified, got %v", got)
		}
	})
}

func TestRecordToolReceiptCapturesStructuredArgs(t *testing.T) {
	savedLoaded := receiptsLoaded
	savedReceipts := recentReceipts
	receiptsLoaded = true
	recentReceipts = nil
	defer func() {
		receiptsLoaded = savedLoaded
		recentReceipts = savedReceipts
	}()

	RecordToolReceipt("schedule_manage", map[string]interface{}{
		"action": "delete",
		"id":     "t1",
		"query":  "secret-key-hunter", // exercises redaction of extra meta
	}, "ok", true)

	r := recentReceipts[len(recentReceipts)-1]
	if r.Meta["action"] != "delete" || r.Meta["id"] != "t1" {
		t.Fatalf("structured args must land in meta, got %v", r.Meta)
	}
}

func TestOperationalClaimsBounded(t *testing.T) {
	text := ""
	for i := 0; i < 40; i++ {
		text += fmt.Sprintf("File /tmp/obj-%d.txt was deleted. ", i)
	}
	if got := LooksLikeOperationalClaims(text); len(got) > 6 {
		t.Fatalf("claims must be bounded to 6, got %d", len(got))
	}
}
