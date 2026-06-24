package collectors

import (
	"scorp-agent/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

// CoolifyData holds Coolify API data.
type CoolifyData struct {
	Available    bool
	Version      string
	Servers      []CoolifyServer
	Applications []CoolifyApp
	Databases    []CoolifyDB
	Services     []CoolifyService
}

type CoolifyServer struct {
	Name      string
	IP        string
	Reachable bool
	Usable    bool
}

type CoolifyApp struct {
	Name      string
	FullName  string
	FQDN      string
	Status    string
	Health    string
	UUID      string
	BuildPack string
	Desc      string
}

type CoolifyDB struct {
	Name   string
	Type   string
	Status string
	Health string
}

type CoolifyService struct {
	Name   string
	Status string
	Health string
}

func CollectCoolify() CoolifyData {
	var d CoolifyData

	// Check version (fast fail)
	versionBody := coolifyGet("/version")
	if versionBody == nil {
		return d
	}
	d.Available = true
	d.Version = strings.Trim(string(versionBody), "\"\n ")

	// Concurrent API calls
	var wg sync.WaitGroup
	var servers, apps, dbs, services json.RawMessage
	var mu sync.Mutex

	fetch := func(path string, target *json.RawMessage) {
		defer wg.Done()
		body := coolifyGet(path)
		if body != nil {
			mu.Lock()
			*target = body
			mu.Unlock()
		}
	}

	wg.Add(4)
	go fetch("/servers", &servers)
	go fetch("/applications", &apps)
	go fetch("/databases", &dbs)
	go fetch("/services", &services)
	wg.Wait()

	// Parse servers
	if servers != nil {
		var list []map[string]interface{}
		if json.Unmarshal(servers, &list) == nil {
			for _, s := range list {
				d.Servers = append(d.Servers, CoolifyServer{
					Name:      jsonStr(s, "name"),
					IP:        jsonStr(s, "ip"),
					Reachable: jsonBool(s, "is_reachable"),
					Usable:    jsonBool(s, "is_usable"),
				})
			}
		}
	}

	// Parse applications
	if apps != nil {
		var list []map[string]interface{}
		if json.Unmarshal(apps, &list) == nil {
			for _, a := range list {
				state, health := parseStatus(jsonStr(a, "status"))
				d.Applications = append(d.Applications, CoolifyApp{
					Name:      cleanAppName(jsonStr(a, "name")),
					FullName:  jsonStr(a, "name"),
					FQDN:      jsonStr(a, "fqdn"),
					Status:    state,
					Health:    health,
					UUID:      jsonStr(a, "uuid"),
					BuildPack: jsonStr(a, "build_pack"),
					Desc:      jsonStr(a, "description"),
				})
			}
		}
	}

	// Parse databases
	if dbs != nil {
		var list []map[string]interface{}
		if json.Unmarshal(dbs, &list) == nil {
			for _, db := range list {
				state, health := parseStatus(jsonStr(db, "status"))
				dbType := jsonStr(db, "type")
				if dbType == "" {
					dbType = jsonStr(db, "database_type")
				}
				d.Databases = append(d.Databases, CoolifyDB{
					Name:   jsonStr(db, "name"),
					Type:   dbType,
					Status: state,
					Health: health,
				})
			}
		}
	}

	// Parse services
	if services != nil {
		var list []map[string]interface{}
		if json.Unmarshal(services, &list) == nil {
			for _, s := range list {
				state, health := parseStatus(jsonStr(s, "status"))
				d.Services = append(d.Services, CoolifyService{
					Name:   jsonStr(s, "name"),
					Status: state,
					Health: health,
				})
			}
		}
	}

	return d
}

func coolifyGet(path string) json.RawMessage {
	url := config.Cfg.CoolifyAPIURL + "/api/v1" + path
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+config.Cfg.CoolifyAPIToken)
	req.Header.Set("Accept", "application/json")

	client := httpShort
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[coolify] API error %s: %v", path, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode != 404 {
			log.Printf("[coolify] API %s returned %d", path, resp.StatusCode)
		}
		return nil
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		// Try as plain text (version endpoint)
		return json.RawMessage(fmt.Sprintf(`"%s"`, path))
	}
	return raw
}

func cleanAppName(raw string) string {
	if raw == "" {
		return "unknown"
	}
	name := raw
	if i := strings.Index(name, ":"); i != -1 {
		name = name[:i]
	}
	if i := strings.LastIndex(name, "/"); i != -1 {
		name = name[i+1:]
	}
	return name
}

func parseStatus(status string) (string, string) {
	if status == "" {
		return "unknown", ""
	}
	parts := strings.SplitN(status, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

func jsonStr(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func jsonBool(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	b, _ := v.(bool)
	return b
}
