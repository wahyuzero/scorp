package main

import (
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
	skills     = make(map[string]Skill)
	skillsMu   sync.RWMutex
	skillsPath = skillsDirPath() + "/" // resolve at init
)

// loadSkills reads all .json files from skills directory
func loadSkills() {
	skillsMu.Lock()
	defer skillsMu.Unlock()

	// Reset skills map
	skills = make(map[string]Skill)

	entries, err := os.ReadDir(skillsPath)
	if err != nil {
		// Directory might not exist yet, that's ok
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

		// Use filename (without .json) as key, but also support name field
		key := strings.TrimSuffix(entry.Name(), ".json")
		skills[key] = skill
	}
}

// listSkills returns a formatted list of all available skills
func listSkills() string {
	skillsMu.RLock()
	defer skillsMu.RUnlock()

	var sb strings.Builder
	sb.WriteString("🎯 <b>Available Skills</b>\n\n")
	for id, skill := range skills {
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

// handleSkillCommand processes /skill commands
// Returns the skill prompt prefix and the user's task
func handleSkillCommand(text string) (skillPrompt string, task string, found bool) {
	// Parse: /skill <skill_name> <task>
	parts := strings.SplitN(text, " ", 3)
	if len(parts) < 2 {
		return "", "", false
	}

	skillName := strings.ToLower(parts[1])
	skillsMu.RLock()
	skill, ok := skills[skillName]
	skillsMu.RUnlock()

	if !ok {
		return "", "", false
	}

	task = ""
	if len(parts) >= 3 {
		task = parts[2]
	}

	return skill.Prompt, task, true
}

// getSkillPromptForMessage checks if agent message mentions a skill
// and returns additional context if matched
func getSkillPromptForMessage(message string) string {
	lower := strings.ToLower(message)
	skillsMu.RLock()
	defer skillsMu.RUnlock()

	for id, skill := range skills {
		// Check if user mentions "pakai skill X" or "use skill X"
		if strings.Contains(lower, "skill "+id) ||
			strings.Contains(lower, "pakai "+id) ||
			strings.Contains(lower, "mode "+id) {
			return fmt.Sprintf("\n\n[Active Skill: %s %s]\n%s", skill.Emoji, skill.Name, skill.Prompt)
		}
	}
	return ""
}

// getSkillByName returns a skill by its key name
func getSkillByName(name string) (Skill, bool) {
	skillsMu.RLock()
	defer skillsMu.RUnlock()
	skill, ok := skills[name]
	return skill, ok
}

// getAllSkills returns a copy of all skills
func getAllSkills() map[string]Skill {
	skillsMu.RLock()
	defer skillsMu.RUnlock()

	result := make(map[string]Skill, len(skills))
	for k, v := range skills {
		result[k] = v
	}
	return result
}

// saveSkill writes a skill to its JSON file
func saveSkill(key string, skill Skill) error {
	skillsMu.Lock()
	defer skillsMu.Unlock()

	// Ensure directory exists
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

	// Update in-memory cache
	skills[key] = skill
	return nil
}

// deleteSkill removes a skill file and from memory
func deleteSkill(key string) error {
	skillsMu.Lock()
	defer skillsMu.Unlock()

	// Don't allow deleting built-in skills
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

	delete(skills, key)
	return nil
}

// sanitizeSkillName converts a name to valid filename
func sanitizeSkillName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	// Remove any non-alphanumeric/underscore
	var sb strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}