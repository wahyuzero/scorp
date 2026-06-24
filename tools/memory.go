package tools

import (
	"scorp-agent/internal/helpers"
	"scorp-agent/config"
	"fmt"
	"log"
	"strings"
	"sync"
)

// ── Memory (persistent key-value store with in-process cache) ──

var memoryFilePathResolved = config.MemoryFilePath() // resolved once at init

var (
	memCache     map[string]string
	memCacheMu   sync.RWMutex
	memLoaded    bool
	persistMu    sync.Mutex // serializes disk writes to prevent race
)

// initMemoryCache loads memory.json into in-process cache at startup
func InitMemoryCache() {
	memCacheMu.Lock()
	defer memCacheMu.Unlock()

	var mem map[string]string
	if err := config.LoadJSON(memoryFilePathResolved, &mem); err != nil {
		log.Printf("[memory] Load error: %v", err)
		memCache = make(map[string]string)
	} else if mem != nil {
		memCache = mem
	} else {
		memCache = make(map[string]string)
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
func SetMemory(key, value string) {
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

	// Persist synchronously to prevent race conditions
	persistMemory(snapshot)
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

	persistMemory(snapshot)
}

// listMemory returns a copy of the entire memory map
func ListMemory() map[string]string {
	memCacheMu.RLock()
	defer memCacheMu.RUnlock()
	result := make(map[string]string, len(memCache))
	for k, v := range memCache {
		result[k] = v
	}
	return result
}

// getMemorySummary returns a compact summary for prompt injection (max 500 chars)
func GetMemorySummary() string {
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

// persistMemory writes the full memory map to disk (serialized to prevent race)
func persistMemory(mem map[string]string) {
	persistMu.Lock()
	defer persistMu.Unlock()
	path := memoryFilePathResolved
	if err := config.SaveJSON(path, mem); err != nil {
		log.Printf("[memory] Save error: %v", err)
	}
}

// executeMemory handles the "memory" tool call from the agent loop
func ExecuteMemory(args map[string]interface{}) (string, bool) {
	action := helpers.GetStringArg(args, "action", "")
	key := helpers.GetStringArg(args, "key", "")
	value := helpers.GetStringArg(args, "value", "")

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
		SetMemory(key, value)
		return fmt.Sprintf("Memory[%s] = %s (saved)", key, value), true
	case "list":
		mem := ListMemory()
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
