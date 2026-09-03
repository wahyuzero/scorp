package tools

import (
	"testing"
)

func TestRecordToolReceipt(t *testing.T) {
	args := map[string]interface{}{
		"target_file": "test.txt",
		"content":     "sample content",
	}

	receipt := RecordToolReceipt("write_file", args, "Success", true)
	if receipt.ReceiptID == "" {
		t.Fatalf("Expected non-empty receipt ID")
	}
	if receipt.Tool != "write_file" {
		t.Errorf("Expected tool 'write_file', got '%s'", receipt.Tool)
	}
	if !receipt.Success {
		t.Errorf("Expected success to be true")
	}

	recent := GetRecentReceipts()
	if len(recent) == 0 {
		t.Fatalf("Expected at least 1 recent receipt")
	}

	found := false
	for _, r := range recent {
		if r.ReceiptID == receipt.ReceiptID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected recorded receipt to be in recent receipts list")
	}
}
