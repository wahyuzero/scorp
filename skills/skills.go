package skills

import (
	"scorp-agent/config"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// Dynamic Skills System — File-based skills
// ──────────────────────────────────────────────

type Skill struct {
	Name              string   `json:"name"`
	Emoji             string   `json:"emoji"`
	Description       string   `json:"description"`
	Category          string   `json:"category"`
	Prompt            string   `json:"prompt"`
	Examples          []string `json:"examples"`
	AutoLoadKeywords  []string `json:"auto_load_keywords"`
}

var (
	cache      = make(map[string]Skill)
	cacheMu    sync.RWMutex
	skillsPath = config.SkillsDirPath() + "/" // resolve at init
)

// Load reads all .json files from skills directory
func Load() {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	cache = make(map[string]Skill)

	entries, err := os.ReadDir(skillsPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		path := filepath.Join(skillsPath, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var skill Skill
		if err := json.Unmarshal(data, &skill); err != nil {
			continue
		}

		key := strings.TrimSuffix(entry.Name(), ".json")
		cache[key] = skill
	}
}

// List returns a formatted list of all available skills
func List() string {
	cacheMu.RLock()
	defer cacheMu.RUnlock()

	var sb strings.Builder
	sb.WriteString("🎯 <b>Available Skills</b>\n\n")
	for id, skill := range cache {
		sb.WriteString(fmt.Sprintf("%s <b>%s</b> — <code>/skill %s</code>\n", skill.Emoji, skill.Name, id))
		sb.WriteString(fmt.Sprintf("   %s\n", skill.Description))
		if len(skill.Examples) > 0 {
			sb.WriteString(fmt.Sprintf("   💡 <i>%s</i>\n", skill.Examples[0]))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("📝 <b>Usage:</b> <code>/skill docker cek status container</code>\n")
	sb.WriteString("Atau di agent mode: <code>pakai skill docker</code>")
	return sb.String()
}

// HandleCommand processes /skill commands
// Returns the skill prompt prefix and the user's task
func HandleCommand(text string) (skillPrompt string, task string, found bool) {
	parts := strings.SplitN(text, " ", 3)
	if len(parts) < 2 {
		return "", "", false
	}

	skillName := strings.ToLower(parts[1])
	cacheMu.RLock()
	skill, ok := cache[skillName]
	cacheMu.RUnlock()

	if !ok {
		return "", "", false
	}

	task = ""
	if len(parts) >= 3 {
		task = parts[2]
	}

	return skill.Prompt, task, true
}

// GetPromptForMessage checks if agent message mentions a skill
// and returns additional context if matched
func GetPromptForMessage(message string) string {
	lower := strings.ToLower(message)
	cacheMu.RLock()
	defer cacheMu.RUnlock()

	for id, skill := range cache {
		if strings.Contains(lower, "skill "+id) ||
			strings.Contains(lower, "pakai "+id) ||
			strings.Contains(lower, "mode "+id) {
			return fmt.Sprintf("\n\n[Active Skill: %s %s]\n%s", skill.Emoji, skill.Name, skill.Prompt)
		}
	}
	return ""
}

// GetByName returns a skill by its key name
func GetByName(name string) (Skill, bool) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	skill, ok := cache[name]
	return skill, ok
}

// GetAll returns a copy of all skills
func GetAll() map[string]Skill {
	cacheMu.RLock()
	defer cacheMu.RUnlock()

	result := make(map[string]Skill, len(cache))
	for k, v := range cache {
		result[k] = v
	}
	return result
}

// Save writes a skill to its JSON file
func Save(key string, skill Skill) error {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if err := os.MkdirAll(skillsPath, 0755); err != nil {
		return err
	}

	path := filepath.Join(skillsPath, key+".json")
	data, err := json.MarshalIndent(skill, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	cache[key] = skill
	return nil
}

// Delete removes a skill file and from memory
func Delete(key string) error {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	builtin := map[string]bool{
		"docker": true, "git": true, "network": true, "coolify": true,
		"backup": true, "security": true, "disk": true, "performance": true,
	}
	if builtin[key] {
		return fmt.Errorf("cannot delete built-in skill: %s", key)
	}

	path := filepath.Join(skillsPath, key+".json")
	if err := os.Remove(path); err != nil {
		return err
	}

	delete(cache, key)
	return nil
}

// SanitizeName converts a name to valid filename
func SanitizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	var sb strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
