package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// ──────────────────────────────────────────────
// Unified Config Manager
// ──────────────────────────────────────────────

// ConfigManager centralizes all config file operations.
// Replaces 12 individual load/save functions with single API.
type ConfigManager struct {
	mu      sync.RWMutex
	baseDir string
}

// NewConfigManager creates a new config manager.
func NewConfigManager(baseDir string) *ConfigManager {
	return &ConfigManager{baseDir: baseDir}
}

// OverrideBaseDir changes the base directory (for testing).
func (cm *ConfigManager) OverrideBaseDir(dir string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.baseDir = dir
}

// Path returns the full path for a config file.
func (cm *ConfigManager) Path(name string) string {
	return filepath.Join(cm.baseDir, name)
}

// Load reads and unmarshals a config file.
// If file doesn't exist, target remains zero-value (no error).
func (cm *ConfigManager) Load(name string, target interface{}) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	path := cm.Path(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, target)
}

// Save marshals and writes a config file.
// Creates parent directories if needed. Default perm 0644.
func (cm *ConfigManager) Save(name string, data interface{}) error {
	return cm.SavePerm(name, data, 0644)
}

// SavePerm saves with custom file permissions.
func (cm *ConfigManager) SavePerm(name string, data interface{}, perm os.FileMode) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	path := cm.Path(name)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonData, perm)
}

// Exists checks if a config file exists.
func (cm *ConfigManager) Exists(name string) bool {
	path := cm.Path(name)
	_, err := os.Stat(path)
	return err == nil
}

// Delete removes a config file.
func (cm *ConfigManager) Delete(name string) error {
	return os.Remove(cm.Path(name))
}

// List returns all config files in the directory.
func (cm *ConfigManager) List() ([]string, error) {
	entries, err := os.ReadDir(cm.baseDir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			files = append(files, e.Name())
		}
	}
	return files, nil
}

// Backup copies all config files to a backup directory.
func (cm *ConfigManager) Backup(backupDir string) error {
	return filepath.Walk(cm.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(cm.baseDir, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(backupDir, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, info.Mode())
	})
}

// Restore restores config files from a backup directory.
func (cm *ConfigManager) Restore(backupDir string) error {
	return filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(backupDir, path)
		if err != nil {
			return err
		}
		dst := cm.Path(rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, info.Mode())
	})
}

// ──────────────────────────────────────────────
// Global instance & convenience wrappers
// ──────────────────────────────────────────────

var configManager *ConfigManager

// InitConfigManager initializes the global config manager.
func InitConfigManager() {
	configManager = NewConfigManager(ScorpDir())
}

// CM returns the global config manager.
func CM() *ConfigManager {
	if configManager == nil {
		InitConfigManager()
	}
	return configManager
}

// ConfigMgr returns the global config manager (shorthand for CM()).
func ConfigMgr() *ConfigManager {
	return CM()
}