package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ── Memory (persistent key-value store with in-process cache) ──

var memoryFilePathResolved = memoryFilePath() // resolved once at init

var (
	memCache   map[string]string
	memCacheMu sync.RWMutex
	memLoaded  bool
)

// initMemoryCache loads memory.json into in-process cache at startup
func initMemoryCache() {
	memCacheMu.Lock()
	defer memCacheMu.Unlock()

	data, err := os.ReadFile(memoryFilePathResolved)
	if err != nil {
		memCache = make(map[string]string)
		memLoaded = true
		return
	}
	var mem map[string]string
	if err := json.Unmarshal(data, &mem); err != nil {
		log.Printf("[memory] Failed to parse memory.json: %v", err)
		memCache = make(map[string]string)
	} else {
		memCache = mem
	}
	memLoaded = true
	log.Printf("[memory] Loaded %d items from memory.json", len(memCache))
}

// getMemory returns a value from the in-process cache
func getMemory(key string) (string, bool) {
	memCacheMu.RLock()
	defer memCacheMu.RUnlock()
	v, ok := memCache[key]
	return v, ok
}

// setMemory sets a value in the cache and persists to disk
func setMemory(key, value string) {
	memCacheMu.Lock()
	if memCache == nil {
		memCache = make(map[string]string)
	}
	memCache[key] = value
	snapshot := make(map[string]string, len(memCache))
	for k, v := range memCache {
		snapshot[k] = v
	}
	memCacheMu.Unlock()

	// Persist asynchronously
	go persistMemory(snapshot)
}

// deleteMemory removes a key from cache and persists
func deleteMemory(key string) {
	memCacheMu.Lock()
	delete(memCache, key)
	snapshot := make(map[string]string, len(memCache))
	for k, v := range memCache {
		snapshot[k] = v
	}
	memCacheMu.Unlock()

	go persistMemory(snapshot)
}

// listMemory returns a copy of the entire memory map
func listMemory() map[string]string {
	memCacheMu.RLock()
	defer memCacheMu.RUnlock()
	result := make(map[string]string, len(memCache))
	for k, v := range memCache {
		result[k] = v
	}
	return result
}

// getMemorySummary returns a compact summary for prompt injection (max 500 chars)
func getMemorySummary() string {
	memCacheMu.RLock()
	defer memCacheMu.RUnlock()

	if len(memCache) == 0 {
		return ""
	}

	var sb strings.Builder
	totalLen := 0
	for key, val := range memCache {
		line := fmt.Sprintf("- %s: %s\n", key, val)
		if totalLen+len(line) > 500 {
			sb.WriteString("- ... (more in memory)\n")
			break
		}
		sb.WriteString(line)
		totalLen += len(line)
	}
	return sb.String()
}

// persistMemory writes the full memory map to disk
func persistMemory(mem map[string]string) {
	path := memoryFilePathResolved
	os.MkdirAll(filepath.Dir(path), 0755)
	data, err := json.MarshalIndent(mem, "", "  ")
	if err != nil {
		log.Printf("[memory] Failed to marshal memory: %v", err)
		return
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Printf("[memory] Failed to write memory.json: %v", err)
	}
}

// executeMemory handles the "memory" tool call from the agent loop
func executeMemory(args map[string]interface{}) (string, bool) {
	action := getStringArg(args, "action", "")
	key := getStringArg(args, "key", "")
	value := getStringArg(args, "value", "")

	switch action {
	case "get":
		if key == "" {
			return "Error: 'key' is required for get", false
		}
		if v, ok := getMemory(key); ok {
			return fmt.Sprintf("Memory[%s] = %s", key, v), true
		}
		return fmt.Sprintf("Memory[%s] not found", key), true
	case "set":
		if key == "" || value == "" {
			return "Error: 'key' and 'value' are required for set", false
		}
		setMemory(key, value)
		return fmt.Sprintf("Memory[%s] = %s (saved)", key, value), true
	case "list":
		mem := listMemory()
		if len(mem) == 0 {
			return "Memory is empty", true
		}
		var sb strings.Builder
		sb.WriteString("📝 Memory:\n")
		for k, v := range mem {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
		return sb.String(), true
	case "delete":
		if key == "" {
			return "Error: 'key' is required for delete", false
		}
		deleteMemory(key)
		return fmt.Sprintf("Memory[%s] deleted", key), true
	default:
		return "Error: action must be get, set, list, or delete", false
	}
}
