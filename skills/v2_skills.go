package skills

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"scorp-agent/config"
)

// ──────────────────────────────────────────────
// Open-Standard SKILL.md Specification
// Supports Level 1 (Metadata: Name, Description) for compact index,
// and Level 2 (Body: Markdown instructions) loaded on-demand via progressive disclosure.
// ──────────────────────────────────────────────

type SkillMeta struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category,omitempty"`
	Emoji       string   `json:"emoji,omitempty"`
	SourcePath  string   `json:"source_path,omitempty"` // Path to SKILL.md
	Scope       string   `json:"scope"`                 // "project" or "global"
	Keywords    []string `json:"keywords,omitempty"`
}

var (
	registryMu sync.RWMutex
	skillIndex = make(map[string]SkillMeta) // name -> metadata

	// Active loaded skill contents (with turn TTL)
	activeSkillsMu sync.RWMutex
	activeSkills   = make(map[string]string) // name -> full markdown body
	activeSkillTTL = make(map[string]int)    // name -> turns left
)

// LoadAllSkills scans both Project-level (.scorp/skills, .agents/skills)
// and Global-level (~/.scorp/skills) directories. Project skills override global ones.
func LoadAllSkills() {
	registryMu.Lock()
	defer registryMu.Unlock()

	skillIndex = make(map[string]SkillMeta)

	// 1. Scan Global Skills (~/.scorp/skills/)
	globalDir := config.SkillsDirPath()
	scanSkillsDirectory(globalDir, "global")

	// 2. Scan Project Skills (.scorp/skills/ and .agents/skills/)
	cwd, err := os.Getwd()
	if err == nil {
		scanSkillsDirectory(filepath.Join(cwd, ".scorp", "skills"), "project")
		scanSkillsDirectory(filepath.Join(cwd, ".agents", "skills"), "project")
	}

	// 3. Fallback: Also scan legacy JSON skills
	scanLegacyJSONSkills(globalDir)
}

func scanSkillsDirectory(dir, scope string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skillPath := filepath.Join(dir, entry.Name(), "SKILL.md")
		if _, err := os.Stat(skillPath); err == nil {
			meta, err := ParseSkillMetadata(skillPath, scope)
			if err == nil && meta.Name != "" {
				skillIndex[meta.Name] = meta
			}
		}
	}
}

// ParseSkillMetadata parses frontmatter from a SKILL.md file without reading the whole body
func ParseSkillMetadata(filePath, scope string) (SkillMeta, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return SkillMeta{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inFrontmatter := false
	meta := SkillMeta{
		SourcePath: filePath,
		Scope:      scope,
		Emoji:      "🎯",
		Category:   "general",
	}

	dirName := filepath.Base(filepath.Dir(filePath))
	meta.Name = dirName

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				// End of frontmatter
				break
			}
		}

		if inFrontmatter && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.ToLower(strings.TrimSpace(parts[0]))
			val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")

			switch key {
			case "name":
				if val != "" {
					meta.Name = val
				}
			case "description":
				meta.Description = val
			case "category":
				meta.Category = val
			case "emoji":
				meta.Emoji = val
			case "keywords":
				rawKeywords := strings.Trim(val, "[]")
				for _, kw := range strings.Split(rawKeywords, ",") {
					clean := strings.TrimSpace(kw)
					if clean != "" {
						meta.Keywords = append(meta.Keywords, clean)
					}
				}
			}
		}
	}

	if meta.Description == "" {
		meta.Description = fmt.Sprintf("Specialized instructions for %s", meta.Name)
	}

	return meta, scanner.Err()
}

// ReadSkillBody loads the full Markdown instruction body of a skill
func ReadSkillBody(skillName string) (string, error) {
	registryMu.RLock()
	meta, exists := skillIndex[skillName]
	registryMu.RUnlock()

	if !exists {
		return "", fmt.Errorf("skill '%s' not found", skillName)
	}

	data, err := os.ReadFile(meta.SourcePath)
	if err != nil {
		return "", fmt.Errorf("read skill file: %w", err)
	}

	// Strip frontmatter if present
	content := string(data)
	if strings.HasPrefix(content, "---") {
		idx := strings.Index(content[3:], "---")
		if idx != -1 {
			content = strings.TrimSpace(content[3+idx+3:])
		}
	}

	return content, nil
}

// ActivateSkill marks a skill active with a turn-based TTL (default 4 turns)
func ActivateSkill(skillName string, ttlTurns int) (string, error) {
	body, err := ReadSkillBody(skillName)
	if err != nil {
		return "", err
	}

	if ttlTurns <= 0 {
		ttlTurns = 4
	}

	activeSkillsMu.Lock()
	activeSkills[skillName] = body
	activeSkillTTL[skillName] = ttlTurns
	activeSkillsMu.Unlock()

	return body, nil
}

// TickActiveSkills decrements TTL and purges expired skills
func TickActiveSkills() {
	activeSkillsMu.Lock()
	defer activeSkillsMu.Unlock()

	for name, ttl := range activeSkillTTL {
		if ttl <= 1 {
			delete(activeSkills, name)
			delete(activeSkillTTL, name)
		} else {
			activeSkillTTL[name] = ttl - 1
		}
	}
}

// GetActiveSkillsContext returns markdown injection for all currently active skills
func GetActiveSkillsContext() string {
	activeSkillsMu.RLock()
	defer activeSkillsMu.RUnlock()

	if len(activeSkills) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n## ACTIVE SKILLS (Specialized Instructions)\n")
	for name, body := range activeSkills {
		sb.WriteString(fmt.Sprintf("\n### [Skill: %s]\n%s\n", name, body))
	}
	return sb.String()
}

// FormatSkillsIndexForSystemPrompt builds Level-1 compact index for system prompt (~30 tokens/skill)
func FormatSkillsIndexForSystemPrompt() string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	if len(skillIndex) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n## AVAILABLE SKILLS\n")
	sb.WriteString("To load full specialized instructions for a task, call tool 'activate_skill' with the skill name.\n")
	for name, meta := range skillIndex {
		sb.WriteString(fmt.Sprintf("- %s %s (%s): %s\n", meta.Emoji, name, meta.Scope, meta.Description))
	}
	return sb.String()
}

// ListSkillsOverview returns user-readable summary for CLI and Telegram
func ListSkillsOverview() string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	if len(skillIndex) == 0 {
		return "🎯 <b>No skills installed.</b>\nCreate one at <code>~/.scorp/skills/&lt;name&gt;/SKILL.md</code>"
	}

	var sb strings.Builder
	sb.WriteString("🎯 <b>Available Skills</b>\n\n")
	for name, meta := range skillIndex {
		sb.WriteString(fmt.Sprintf("%s <b>%s</b> <i>[%s]</i>\n", meta.Emoji, name, meta.Scope))
		sb.WriteString(fmt.Sprintf("   %s\n", meta.Description))
		sb.WriteString(fmt.Sprintf("   ↳ Path: <code>%s</code>\n\n", meta.SourcePath))
	}
	sb.WriteString("💡 <i>Gunakan: <code>/skill &lt;name&gt;</code> atau panggil saat dibutuhkan.</i>")
	return sb.String()
}

// scanLegacyJSONSkills scans backward-compatible .json skill files
func scanLegacyJSONSkills(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".json")
		if _, exists := skillIndex[name]; !exists {
			skillIndex[name] = SkillMeta{
				Name:        name,
				Description: fmt.Sprintf("Legacy skill: %s", name),
				SourcePath:  filepath.Join(dir, entry.Name()),
				Scope:       "global",
				Emoji:       "📜",
			}
		}
	}
}
