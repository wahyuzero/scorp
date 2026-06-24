package tools

import (
	"scorp-agent/internal/helpers"
	"scorp-agent/config"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Database Query Tool — SQLite, PostgreSQL, MySQL
// ──────────────────────────────────────────────

type dbConnection struct {
	Type       string `json:"type"`       // "sqlite", "postgres", "mysql"
	Connection string `json:"connection"` // file path or DSN
}

var dbConnectionsFile = config.DBConnectionsPath()

func loadDBConnections() map[string]dbConnection {
	conn := make(map[string]dbConnection)
	data, err := os.ReadFile(dbConnectionsFile)
	if err != nil {
		return conn
	}
	json.Unmarshal(data, &conn)
	return conn
}

func ExecuteSQL(args map[string]interface{}, chatID int64) (string, bool) {
	query := helpers.GetStringArg(args, "query", "")
	if query == "" {
		return "Error: 'query' argument is required", false
	}

	rowLimit := helpers.GetIntArg(args, "limit", 100)
	if rowLimit > 1000 {
		rowLimit = 1000
	}

	// Determine connection
	connName := helpers.GetStringArg(args, "connection", "default")
	dbType := helpers.GetStringArg(args, "db_type", "")
	dsn := helpers.GetStringArg(args, "dsn", "")

	if dbType == "" || dsn == "" {
		// Try named connection from config
		conns := loadDBConnections()
		if conn, ok := conns[connName]; ok {
			dbType = conn.Type
			dsn = conn.Connection
		} else if conn, ok := conns["default"]; ok {
			dbType = conn.Type
			dsn = conn.Connection
		} else {
			return "Error: no connection specified. Provide db_type + dsn, or configure ~/.scorp/db_connections.json", false
		}
	}

	// Check for write operations (require confirmation for non-SELECT)
	trimmedQuery := strings.TrimSpace(strings.ToUpper(query))
	isWrite := strings.HasPrefix(trimmedQuery, "INSERT") ||
		strings.HasPrefix(trimmedQuery, "UPDATE") ||
		strings.HasPrefix(trimmedQuery, "DELETE") ||
		strings.HasPrefix(trimmedQuery, "DROP") ||
		strings.HasPrefix(trimmedQuery, "CREATE") ||
		strings.HasPrefix(trimmedQuery, "ALTER") ||
		strings.HasPrefix(trimmedQuery, "TRUNCATE")

	if isWrite {
		chatIDStr := fmt.Sprintf("%d", chatID)
		StorePendingConfirmation(chatIDStr, "sql", query, nil)
		return fmt.Sprintf("⚠️ WRITE QUERY DETECTED:\n%s\n\nThis will modify the database. Please confirm execution.", query), false
	}

	// Open database
	var db *sql.DB
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch strings.ToLower(dbType) {
	case "sqlite":
		db, err = sql.Open("sqlite3", dsn)
	case "postgres":
		db, err = sql.Open("postgres", dsn)
	case "mysql":
		db, err = sql.Open("mysql", dsn)
	default:
		return fmt.Sprintf("Error: unsupported db_type '%s'. Use: sqlite, postgres, mysql", dbType), false
	}

	if err != nil {
		return fmt.Sprintf("Error opening database: %v", err), false
	}
	defer db.Close()

	// Add LIMIT if not present
	if !strings.Contains(strings.ToUpper(query), "LIMIT") {
		query = strings.TrimRight(query, ";") + fmt.Sprintf(" LIMIT %d;", rowLimit)
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Sprintf("Query error: %v", err), false
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return fmt.Sprintf("Error getting columns: %v", err), false
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 SQL Query: %s\n", query))
	sb.WriteString(fmt.Sprintf("Type: %s | Columns: %d\n\n", dbType, len(cols)))

	// Header
	sb.WriteString(strings.Join(cols, " | ") + "\n")
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	// Rows
	rowCount := 0
	for rows.Next() {
		if rowCount >= rowLimit {
			sb.WriteString(fmt.Sprintf("\n... (%d row limit reached)\n", rowLimit))
			break
		}
		values := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			continue
		}

		strVals := make([]string, len(cols))
		for i, v := range values {
			if v == nil {
				strVals[i] = "NULL"
			} else {
				strVals[i] = fmt.Sprintf("%v", v)
				if len(strVals[i]) > 50 {
					strVals[i] = strVals[i][:47] + "..."
				}
			}
		}
		sb.WriteString(strings.Join(strVals, " | ") + "\n")
		rowCount++
	}

	sb.WriteString(fmt.Sprintf("\n📊 %d row(s) returned\n", rowCount))
	return helpers.TruncOutput(sb.String(), helpers.MaxToolOutput), true
}
