package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// loadJSON reads a JSON file and unmarshals into target.
// Returns nil error if file doesn't exist (target stays zero-value).
func LoadJSON(path string, target interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // file doesn't exist yet — not an error
		}
		return err
	}
	return json.Unmarshal(data, target)
}

// saveJSON marshals data as pretty JSON, creates parent dirs, and writes file.
// Default permission 0644.
func SaveJSON(path string, data interface{}) error {
	return SaveJSONPerm(path, data, 0644)
}

// SaveJSONPerm marshals data as pretty JSON, creates parent dirs, writes with custom perm.
func SaveJSONPerm(path string, data interface{}, perm os.FileMode) error {
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
