package tools

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestExecuteSQLSingleUnnamedConnection pins the J20/F-27 fix: a config with
// exactly one connection under an arbitrary key must be used automatically.
// Before the fix the tool rejected the config (it only looked up the
// "connection"/"default" names) and the model burned 10+ retries.
func TestExecuteSQLSingleUnnamedConnection(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	cfgPath := filepath.Join(dir, "db_connections.json")
	if err := os.WriteFile(cfgPath, []byte(`{
		"sessions": {"type": "sqlite", "connection": "`+dbPath+`"}
	}`), 0644); err != nil {
		t.Fatal(err)
	}

	saved := dbConnectionsFile
	dbConnectionsFile = cfgPath
	defer func() { dbConnectionsFile = saved }()

	// Create the table directly.
	if _, err := os.Stat(dbPath); err == nil {
		os.Remove(dbPath)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Skipf("sqlite unavailable: %v", err)
	}
	defer db.Close()
	if _, err := db.Exec("CREATE TABLE t_f27 (id INTEGER)"); err != nil {
		t.Fatalf("setup create failed: %v", err)
	}
	if _, err := db.Exec("INSERT INTO t_f27 (id) VALUES (7)"); err != nil {
		t.Fatalf("setup insert failed: %v", err)
	}

	out, ok := ExecuteSQL(map[string]interface{}{"query": "SELECT id FROM t_f27"}, 0)
	if !ok {
		t.Fatalf("single unnamed connection must be used automatically, out=%.200q", out)
	}
	if !strings.Contains(out, "7") {
		t.Fatalf("expected row value 7, got %q", out)
	}
}

// TestExecuteSQLMultiConnectionErrorNamesConnections keeps the error honest
// when several named connections exist and none matches.
func TestExecuteSQLMultiConnectionErrorNamesConnections(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "db_connections.json")
	cfg := map[string]dbConnection{
		"alpha": {Type: "sqlite", Connection: filepath.Join(dir, "a.db")},
		"beta":  {Type: "sqlite", Connection: filepath.Join(dir, "b.db")},
	}
	data, _ := json.Marshal(cfg)
	if err := os.WriteFile(cfgPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	saved := dbConnectionsFile
	dbConnectionsFile = cfgPath
	defer func() { dbConnectionsFile = saved }()

	out, ok := ExecuteSQL(map[string]interface{}{"query": "SELECT 1"}, 0)
	if ok {
		t.Fatal("ambiguous multi-connection config must not silently pick one")
	}
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Fatalf("error must name the available connections, got %q", out)
	}
}

// TestLoadReceiptsQuarantinesCorruptFile pins the C3/F-11 fix: a corrupted
// receipts.json must be quarantined (old evidence preserved for post-mortem)
// instead of silently overwritten on the next save.
func TestLoadReceiptsQuarantinesCorruptFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	// Reset package-level load state so this test re-reads from disk.
	savedLoaded := receiptsLoaded
	savedReceipts := recentReceipts
	receiptsLoaded = false
	recentReceipts = nil
	defer func() {
		receiptsLoaded = savedLoaded
		recentReceipts = savedReceipts
	}()

	receiptsPath := filepath.Join(home, ".scorp", "receipts.json")
	if err := os.MkdirAll(filepath.Dir(receiptsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(receiptsPath, []byte("{{{ not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	loadReceiptsLocked()

	matches, err := filepath.Glob(receiptsPath + ".corrupt-*")
	if err != nil || len(matches) != 1 {
		t.Fatalf("corrupted receipts file must be quarantined (matches=%v err=%v)", matches, err)
	}
	data, _ := os.ReadFile(matches[0])
	if string(data) != "{{{ not json" {
		t.Fatalf("quarantined copy must preserve the raw corrupted bytes, got %q", data)
	}
}
