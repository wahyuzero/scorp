package sop

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"scorp-agent/config"
)

// ──────────────────────────────────────────────
// Standard Operating Procedures (SOP) Engine (ZeroClaw Parity)
// Enables declarative, repeatable playbooks and routines
// ──────────────────────────────────────────────

// SOP defines a Standard Operating Procedure
type SOP struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	Prompt      string   `json:"prompt"`
}

var (
	sopsMu sync.RWMutex
)

// Dir returns the directory storing SOP playbooks (~/.scorp/sops)
func Dir() string {
	d := config.ScorpPath("sops")
	_ = os.MkdirAll(d, 0755)
	return d
}

// InitDefaultSOPs creates built-in template playbooks if none exist
func InitDefaultSOPs() {
	sopsMu.Lock()
	defer sopsMu.Unlock()

	dir := Dir()

	defaults := []SOP{
		{
			Name:        "health_audit",
			Description: "Comprehensive server health & resource audit",
			Steps: []string{
				"1. Inspect CPU, Memory, Disk, and Network via system_info",
				"2. Check running processes and services with process",
				"3. Check status of docker containers if available",
				"4. Summarize health status and flag any resources >85% usage",
			},
			Prompt: "Execute SOP 'health_audit': Check CPU, RAM, Disk, top processes, and docker services. Produce a clean, structured health report.",
		},
		{
			Name:        "code_review",
			Description: "Codebase status, git diff, and test validation",
			Steps: []string{
				"1. Check git status and recent commits via git",
				"2. Run test suites with shell (e.g. go test ./... or npm test)",
				"3. Identify modified files and verify test integrity",
				"4. Provide concise status and actionable recommendations",
			},
			Prompt: "Execute SOP 'code_review': Inspect git status, run project test suite via shell, and report code readiness.",
		},
		{
			Name:        "site_check",
			Description: "Verify web service reachability and content via low-RAM web engine",
			Steps: []string{
				"1. Read target endpoint using read_url (<5MB RAM)",
				"2. Verify HTTP response and title header",
				"3. Report availability and page excerpt",
			},
			Prompt: "Execute SOP 'site_check': Use read_url to inspect target website availability and return reader excerpt.",
		},
	}

	for _, s := range defaults {
		p := filepath.Join(dir, s.Name+".json")
		if _, err := os.Stat(p); os.IsNotExist(err) {
			data, _ := json.MarshalIndent(s, "", "  ")
			_ = os.WriteFile(p, data, 0644)
		}
	}
}

// ListSOPs returns all available SOP playbooks
func ListSOPs() []SOP {
	sopsMu.RLock()
	defer sopsMu.RUnlock()

	dir := Dir()
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var list []SOP
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			p := filepath.Join(dir, f.Name())
			if data, err := os.ReadFile(p); err == nil {
				var s SOP
				if err := json.Unmarshal(data, &s); err == nil {
					list = append(list, s)
				}
			}
		}
	}
	return list
}

// GetSOP retrieves a specific SOP by name
func GetSOP(name string) (*SOP, error) {
	sopsMu.RLock()
	defer sopsMu.RUnlock()

	cleanName := strings.TrimSuffix(name, ".json")
	p := filepath.Join(Dir(), cleanName+".json")
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("SOP '%s' not found: %w", name, err)
	}

	var s SOP
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("error parsing SOP: %w", err)
	}
	return &s, nil
}

// SaveSOP creates or updates an SOP playbook
func SaveSOP(sop SOP) error {
	sopsMu.Lock()
	defer sopsMu.Unlock()

	if sop.Name == "" {
		return fmt.Errorf("SOP name cannot be empty")
	}

	p := filepath.Join(Dir(), sop.Name+".json")
	data, err := json.MarshalIndent(sop, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}
