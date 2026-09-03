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
	ReceiptID string    `json:"receipt_id"`
	Tool      string    `json:"tool"`
	ArgsHash  string    `json:"args_hash"`
	OutHash   string    `json:"out_hash"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
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

// RecordToolReceipt logs a cryptographic receipt of tool execution
func RecordToolReceipt(toolName string, args map[string]interface{}, output string, ok bool) ToolReceipt {
	argsJSON, _ := json.Marshal(args)
	argsH := sha256.Sum256(argsJSON)
	argsHash := hex.EncodeToString(argsH[:])

	outH := sha256.Sum256([]byte(output))
	outHash := hex.EncodeToString(outH[:])

	now := time.Now()
	rawID := fmt.Sprintf("%s:%s:%s:%d", toolName, argsHash, outHash, now.UnixNano())
	idH := sha256.Sum256([]byte(rawID))
	receiptID := hex.EncodeToString(idH[:8]) // 16 hex chars

	receipt := ToolReceipt{
		ReceiptID: receiptID,
		Tool:      toolName,
		ArgsHash:  argsHash[:16],
		OutHash:   outHash[:16],
		Success:   ok,
		Timestamp: now,
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
