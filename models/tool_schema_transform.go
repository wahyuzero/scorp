package models

import (
	"strings"

	"scorp-agent/registry"
)

// ──────────────────────────────────────────────
// Tool Schema Transformation
// Adapts JSON schemas for sensitive providers (e.g. Gemini, DeepSeek, GLM)
// that reject complex JSON Schema constructs ($schema, anyOf, etc.)
// ──────────────────────────────────────────────

const (
	ToolSchemaTransformOff    = ""
	ToolSchemaTransformSimple = "simple"
)

// NormalizeToolSchemaTransform returns canonical transform name
func NormalizeToolSchemaTransform(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "simple", "basic", "strict", "flat":
		return ToolSchemaTransformSimple
	default:
		return ToolSchemaTransformOff
	}
}

// TransformToolDefinitions transforms tool parameters schema according to mode
func TransformToolDefinitions(tools []registry.ToolSchema, transform string) []registry.ToolSchema {
	mode := NormalizeToolSchemaTransform(transform)
	if mode == ToolSchemaTransformOff || len(tools) == 0 {
		return tools
	}

	out := make([]registry.ToolSchema, len(tools))
	for i, tool := range tools {
		out[i] = tool
		if tool.Type == "function" && tool.Function.Parameters != nil {
			out[i].Function.Parameters = SanitizeToolSchema(tool.Function.Parameters, mode)
		}
	}
	return out
}

// SanitizeToolSchema cleans up schema keywords that sensitive backends reject
func SanitizeToolSchema(schema map[string]interface{}, mode string) map[string]interface{} {
	if schema == nil {
		return map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		}
	}

	result := make(map[string]interface{})
	for k, v := range schema {
		// Strip unsupported keywords
		switch k {
		case "$schema", "title", "id", "$id", "definitions", "$defs":
			continue
		case "additionalProperties":
			if _, isBool := v.(bool); !isBool {
				continue
			}
			result[k] = v
		case "properties":
			if propsMap, ok := v.(map[string]interface{}); ok {
				cleanProps := make(map[string]interface{})
				for propKey, propVal := range propsMap {
					if propSchema, ok := propVal.(map[string]interface{}); ok {
						cleanProps[propKey] = sanitizeProperty(propSchema, mode)
					} else {
						cleanProps[propKey] = propVal
					}
				}
				result["properties"] = cleanProps
			}
		default:
			result[k] = v
		}
	}

	if _, hasType := result["type"]; !hasType {
		result["type"] = "object"
	}
	if _, hasProps := result["properties"]; !hasProps && result["type"] == "object" {
		result["properties"] = map[string]interface{}{}
	}

	return result
}

func sanitizeProperty(prop map[string]interface{}, mode string) map[string]interface{} {
	if prop == nil {
		return map[string]interface{}{"type": "string"}
	}

	clean := make(map[string]interface{})
	for k, v := range prop {
		switch k {
		case "$schema", "title", "id", "$id", "definitions", "$defs":
			continue
		case "anyOf", "oneOf", "allOf":
			// Simplify composite schemas to first branch's type
			if branches, ok := v.([]interface{}); ok && len(branches) > 0 {
				if first, ok := branches[0].(map[string]interface{}); ok {
					if t, ok := first["type"]; ok {
						clean["type"] = t
					}
				}
			}
		case "properties":
			if subProps, ok := v.(map[string]interface{}); ok {
				nested := make(map[string]interface{})
				for subKey, subVal := range subProps {
					if subSchema, ok := subVal.(map[string]interface{}); ok {
						nested[subKey] = sanitizeProperty(subSchema, mode)
					} else {
						nested[subKey] = subVal
					}
				}
				clean["properties"] = nested
			}
		case "items":
			if itemsSchema, ok := v.(map[string]interface{}); ok {
				clean["items"] = sanitizeProperty(itemsSchema, mode)
			} else {
				clean["items"] = v
			}
		default:
			clean[k] = v
		}
	}

	if _, hasType := clean["type"]; !hasType {
		clean["type"] = "string"
	}

	return clean
}
