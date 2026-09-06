package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"scorp-agent/config"
)

// ──────────────────────────────────────────────
// Cryptographic Tool Execution Receipts (ZeroClaw Parity)
// Every executed tool generates a SHA-256 verifiable receipt
// ──────────────────────────────────────────────

type ToolReceipt struct {
	ReceiptID string            `json:"receipt_id"`
	Tool      string            `json:"tool"`
	ArgsHash  string            `json:"args_hash"`
	OutHash   string            `json:"out_hash"`
	Success   bool              `json:"success"`
	Timestamp time.Time         `json:"timestamp"`
	Meta      map[string]string `json:"meta,omitempty"` // non-secret facts for gate queries (cmd, path)
}

var (
	recentReceipts   []ToolReceipt
	recentReceiptsMu sync.RWMutex
	maxReceipts      = 50
	receiptsLoaded   bool
)

func loadReceiptsLocked() {
	if receiptsLoaded {
		return
	}
	receiptsLoaded = true
	p := config.ScorpPath("receipts.json")
	if data, err := os.ReadFile(p); err == nil {
		_ = json.Unmarshal(data, &recentReceipts)
	}
}

func saveReceiptsLocked() {
	p := config.ScorpPath("receipts.json")
	if data, err := json.Marshal(recentReceipts); err == nil {
		_ = os.WriteFile(p, data, 0644)
	}
}

// RecordToolReceipt logs a cryptographic receipt of tool execution. Optional
// extraMeta entries (e.g. auto_decision from the P3.13 classifier) are merged
// into the receipt meta.
func RecordToolReceipt(toolName string, args map[string]interface{}, output string, ok bool, extraMeta ...map[string]string) ToolReceipt {
	argsJSON, _ := json.Marshal(args)
	argsH := sha256.Sum256(argsJSON)
	argsHash := hex.EncodeToString(argsH[:])

	outH := sha256.Sum256([]byte(output))
	outHash := hex.EncodeToString(outH[:])

	now := time.Now()
	rawID := fmt.Sprintf("%s:%s:%s:%d", toolName, argsHash, outHash, now.UnixNano())
	idH := sha256.Sum256([]byte(rawID))
	receiptID := hex.EncodeToString(idH[:8]) // 16 hex chars

	// Meta carries the few plaintext facts gates need (test-integrity gate,
	// audit greps) without weakening the hash scheme: the hashes above still
	// pin the exact args/output. Secrets are redacted before storage.
	meta := map[string]string{}
	if cmd, ok := args["command"].(string); ok && cmd != "" {
		if len(cmd) > 300 {
			cmd = cmd[:300]
		}
		meta["cmd"] = RedactSecrets(cmd)
	}
	for _, k := range []string{"path", "target_file", "file", "file_path"} {
		if v, ok := args[k].(string); ok && v != "" {
			meta["path"] = v
			break
		}
	}
	for _, extra := range extraMeta {
		for k, v := range extra {
			if v != "" {
				meta[k] = v
			}
		}
	}

	receipt := ToolReceipt{
		ReceiptID: receiptID,
		Tool:      toolName,
		ArgsHash:  argsHash[:16],
		OutHash:   outHash[:16],
		Success:   ok,
		Timestamp: now,
		Meta:      meta,
	}

	recentReceiptsMu.Lock()
	defer recentReceiptsMu.Unlock()
	loadReceiptsLocked()
	recentReceipts = append(recentReceipts, receipt)
	if len(recentReceipts) > maxReceipts {
		recentReceipts = recentReceipts[len(recentReceipts)-maxReceipts:]
	}
	saveReceiptsLocked()

	return receipt
}

// GetRecentReceipts returns recently recorded receipts
func GetRecentReceipts() []ToolReceipt {
	recentReceiptsMu.Lock()
	defer recentReceiptsMu.Unlock()
	loadReceiptsLocked()
	cp := make([]ToolReceipt, len(recentReceipts))
	copy(cp, recentReceipts)
	return cp
}
