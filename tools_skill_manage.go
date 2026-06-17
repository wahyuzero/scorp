package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ──────────────────────────────────────────────
// skill_manage tool — CRUD for skills
// ──────────────────────────────────────────────

func executeSkillManage(args map[string]interface{}) (string, bool) {
	action, _ := args["action"].(string)
	name, _ := args["name"].(string)
	content, _ := args["content"].(string)

	switch action {
	case "list":
		return executeSkillManageList()
	case "view":
		return executeSkillManageView(name)
	case "create":
		return executeSkillManageCreate(name, content)
	case "update":
		return executeSkillManageUpdate(name, content)
	case "delete":
		return executeSkillManageDelete(name)
	default:
		return fmt.Sprintf("Unknown action: %s. Use: list, view, create, update, delete", action), false
	}
}

func executeSkillManageList() (string, bool) {
	allSkills := getAllSkills()
	if len(allSkills) == 0 {
		return "No skills found.", true
	}

	var sb strings.Builder
	sb.WriteString("🎯 <b>Skills</b>\n\n")
	for key, skill := range allSkills {
		sb.WriteString(fmt.Sprintf("%s <b>%s</b> (<code>%s</code>)\n", skill.Emoji, skill.Name, key))
		sb.WriteString(fmt.Sprintf("   Category: %s\n", skill.Category))
		sb.WriteString(fmt.Sprintf("   %s\n\n", skill.Description))
	}
	return sb.String(), true
}

func executeSkillManageView(name string) (string, bool) {
	if name == "" {
		return "Missing 'name' parameter", false
	}
	skill, ok := getSkillByName(name)
	if !ok {
		return fmt.Sprintf("Skill not found: %s", name), false
	}

	data, _ := json.MarshalIndent(skill, "", "  ")
	return fmt.Sprintf("📋 <b>Skill: %s</b>\n<pre>%s</pre>", name, string(data)), true
}

func executeSkillManageCreate(name, content string) (string, bool) {
	if name == "" || content == "" {
		return "Missing 'name' or 'content' parameter", false
	}

	key := sanitizeSkillName(name)
	if _, ok := getSkillByName(key); ok {
		return fmt.Sprintf("Skill already exists: %s already exists", key), false
	}

	var skill Skill
	if err := json.Unmarshal([]byte(content), &skill); err != nil {
		return fmt.Sprintf("Invalid JSON: %v", err), false
	}

	// Override name from parameter
	skill.Name = name

	if err := saveSkill(key, skill); err != nil {
		return fmt.Sprintf("Failed to save: %v", err), false
	}

	return fmt.Sprintf("✅ Skill created: %s (%s)", skill.Name, key), true
}

func executeSkillManageUpdate(name, content string) (string, bool) {
	if name == "" || content == "" {
		return "Missing 'name' or 'content' parameter", false
	}

	key := sanitizeSkillName(name)
	_, ok := getSkillByName(key)
	if !ok {
		return fmt.Sprintf("Skill not found: %s", name), false
	}

	var skill Skill
	if err := json.Unmarshal([]byte(content), &skill); err != nil {
		return fmt.Sprintf("Invalid JSON: %v", err), false
	}

	skill.Name = name

	if err := saveSkill(key, skill); err != nil {
		return fmt.Sprintf("Failed to update: %v", err), false
	}

	return fmt.Sprintf("✅ Skill updated: %s (%s)", skill.Name, key), true
}

func executeSkillManageDelete(name string) (string, bool) {
	if name == "" {
		return "Missing 'name' parameter", false
	}

	key := sanitizeSkillName(name)
	if err := deleteSkill(key); err != nil {
		return fmt.Sprintf("Failed to delete: %v", err), false
	}

	return fmt.Sprintf("✅ Skill deleted: %s", key), true
}