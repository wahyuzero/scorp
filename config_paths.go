package main

import (
	"os"
	"path/filepath"
)

// ──────────────────────────────────────────────
// Centralized path resolution — replaces hardcoded /home/ubuntu paths
// ──────────────────────────────────────────────

// homeDir returns the user's home directory from $HOME env.
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	if h, err := os.UserHomeDir(); err == nil {
		return h
	}
	return "/tmp" // last resort
}

// scorpDir returns the scorp-agent config directory (~/.scorp-agent/).
func scorpDir() string {
	return filepath.Join(homeDir(), ".scorp-agent")
}

// scorpPath joins sub-paths under the scorp-agent config directory.
func scorpPath(parts ...string) string {
	all := append([]string{scorpDir()}, parts...)
	return filepath.Join(all...)
}

// uploadsDir returns the uploads directory.
func uploadsDir() string {
	return filepath.Join(homeDir(), "uploads")
}

// historyDirPath returns the conversation history directory.
func historyDirPath() string {
	return scorpPath("history")
}

// memoryFilePath returns the memory.json path.
func memoryFilePath() string {
	return scorpPath("memory.json")
}

// skillsDirPath returns the skills directory.
func skillsDirPath() string {
	return scorpPath("skills")
}

// screenshotsDir returns the browser screenshots directory.
func screenshotsDir() string {
	return scorpPath("screenshots")
}

// sttPythonPath returns the STT Python binary path.
func sttPythonPath() string {
	return filepath.Join(homeDir(), ".hermes", "hermes-agent", "venv", "bin", "python3")
}

// ttsBinPath returns the TTS (edge-tts) binary path.
// edge-tts is installed in user's hermes venv (not root's).
func ttsBinPath() string {
	// Check if installed in ubuntu's venv
	ubuntuPath := "/home/ubuntu/.hermes/hermes-agent/venv/bin/edge-tts"
	if _, err := os.Stat(ubuntuPath); err == nil {
		return ubuntuPath
	}
	// Fallback to homeDir-based path
	return filepath.Join(homeDir(), ".hermes", "hermes-agent", "venv", "bin", "edge-tts")
}

// sttScriptPath returns the STT script path.
func sttScriptPath() string {
	return filepath.Join(homeDir(), "projects", "vps-monitor-go", "stt.py")
}

// ragDirPath returns the RAG vector index directory path.
func ragDirPath() string {
	return filepath.Join(homeDir(), ".scorp-agent", "rag")
}

// ragVectorDBPath returns the vector RAG index file path.
func ragVectorDBPath() string {
	return filepath.Join(ragDirPath(), "vector_index.json")
}

// ragIndexPath returns the TF-IDF RAG index file path.
func ragIndexPath() string {
	return filepath.Join(ragDirPath(), "index.json")
}

// hermesConfigPath_ returns the Hermes config file path.
func hermesConfigPath_() string {
	return filepath.Join(homeDir(), ".hermes", "config.yaml")
}

// omhConfigPath_ returns the OMH config file path.
func omhConfigPath_() string {
	return filepath.Join(homeDir(), ".hermes", "plugins", "omh", "config.yaml")
}

// browserSessionsDir returns the browser sessions directory (persistent user data).
func browserSessionsDir() string {
	return scorpPath("browser_sessions")
}

// browserSessionPath returns the user data dir path for a specific session.
func browserSessionPath(sessionID string) string {
	return filepath.Join(browserSessionsDir(), sessionID)
}
