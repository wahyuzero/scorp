package config

import (
	"os"
	"path/filepath"
)

// ──────────────────────────────────────────────
// Centralized path resolution — replaces hardcoded /home/ubuntu paths
// ──────────────────────────────────────────────

// homeDir returns the user's home directory from $HOME env.
func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	if h, err := os.UserHomeDir(); err == nil {
		return h
	}
	return "/tmp" // last resort
}

// scorpDir returns the scorp config directory (~/.scorp/).
func ScorpDir() string {
	return filepath.Join(HomeDir(), ".scorp")
}

// scorpPath joins sub-paths under the scorp-agent config directory.
func ScorpPath(parts ...string) string {
	all := append([]string{ScorpDir()}, parts...)
	return filepath.Join(all...)
}

// uploadsDir returns the uploads directory.
func UploadsDir() string {
	return filepath.Join(HomeDir(), "uploads")
}

// historyDirPath returns the conversation history directory.
func HistoryDirPath() string {
	return ScorpPath("history")
}

// memoryFilePath returns the memory.json path.
func MemoryFilePath() string {
	return ScorpPath("memory.json")
}

// skillsDirPath returns the skills directory.
func SkillsDirPath() string {
	return ScorpPath("skills")
}

// screenshotsDir returns the browser screenshots directory.
func ScreenshotsDir() string {
	return ScorpPath("screenshots")
}

// ragDirPath returns the RAG vector index directory path.
func RagDirPath() string {
	return filepath.Join(HomeDir(), ".scorp", "rag")
}

// ragVectorDBPath returns the vector RAG index file path.
func RagVectorDBPath() string {
	return filepath.Join(RagDirPath(), "vector_index.json")
}

// ragIndexPath returns the TF-IDF RAG index file path.
func RagIndexPath() string {
	return filepath.Join(RagDirPath(), "index.json")
}

// browserSessionsDir returns the browser sessions directory (persistent user data).
func BrowserSessionsDir() string {
	return ScorpPath("browser_sessions")
}

// browserSessionPath returns the user data dir path for a specific session.
func BrowserSessionPath(sessionID string) string {
	return filepath.Join(BrowserSessionsDir(), sessionID)
}

// ──────────────────────────────────────────────
// Additional config file paths
// ──────────────────────────────────────────────

// dbConnectionsPath returns the DB connections config file.
func DBConnectionsPath() string {
	return ScorpPath("db_connections.json")
}

// mcpConfigFilePath returns the MCP servers config file.
func MCPConfigFilePath() string {
	return ScorpPath("mcp.json")
}

// costConfigFilePath returns the cost optimization config.
func CostConfigFilePath() string {
	return ScorpPath("cost_config.json")
}

// costLogFilePath returns the daily cost log.
func CostLogFilePath() string {
	return ScorpPath("cost_daily.json")
}

// schedulerFilePath returns the scheduler state file.
func SchedulerFilePath() string {
	return ScorpPath("scheduler.json")
}

// modelCfgFilePath returns the model router config.
func ModelCfgFilePath() string {
	return ScorpPath("models.json")
}

// modelUsageFilePath returns the model usage stats.
func ModelUsageFilePath() string {
	return ScorpPath("model_usage.json")
}

// vncLogPath returns the VNC log path (if exists).
func VNCLogPath() string {
	// Try common VNC display :3 first, then :1
	candidates := []string{
		filepath.Join(HomeDir(), ".vnc", Hostname()+":3.log"),
		filepath.Join(HomeDir(), ".vnc", Hostname()+":1.log"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	// Default
	return filepath.Join(HomeDir(), ".vnc", Hostname()+":3.log")
}

// pythonSitePackages returns the user's Python site-packages path.
func PythonSitePackages() string {
	return filepath.Join(HomeDir(), ".local", "lib", "python3.12", "site-packages")
}

// projectDir returns the scorp-agent project directory.
func ProjectDir() string {
	// Check if running from source
	src := filepath.Join(HomeDir(), "projects", "vps-monitor-go")
	if _, err := os.Stat(filepath.Join(src, "go.mod")); err == nil {
		return src
	}
	// Fallback: current working directory
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return src
}

// hostname returns the system hostname (for VNC log paths).
func Hostname() string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "localhost"
	}
	return h
}
