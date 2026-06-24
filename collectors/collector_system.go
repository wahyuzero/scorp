package collectors

import (
	"scorp-agent/config"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SystemData holds system metrics.
type SystemData struct {
	CPUPercent   float64
	CPUCount     int
	LoadAvg      [3]float64
	RAMTotalGB   float64
	RAMUsedGB    float64
	RAMAvailGB   float64
	RAMPercent   float64
	SwapTotalGB  float64
	SwapUsedGB   float64
	SwapPercent  float64
	DiskTotalGB  float64
	DiskUsedGB   float64
	DiskPercent  float64
	NetSentGB    float64
	NetRecvGB    float64
	Uptime       string
	TopProcesses []TopProcess
}

type TopProcess struct {
	PID  int32
	Name string
	CPU  float64
	Mem  float64
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

func FormatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	secs := int(d.Seconds()) % 60
	if days > 0 {
		return fmt.Sprintf("%d days, %d:%02d:%02d", days, hours, mins, secs)
	}
	return fmt.Sprintf("%d:%02d:%02d", hours, mins, secs)
}

// ──────────────────────────────────────────────
// Network collector (ports, connections, traefik)
// ──────────────────────────────────────────────

type NetworkData struct {
	ListeningPorts int
	NewPorts       []PortInfo
	Connections    ConnectionInfo
	Traefik        TraefikInfo
}

type PortInfo struct {
	Port    int
	Address string
	Process string
}

type ConnectionInfo struct {
	Total       int
	UniqueIPs   int
	ExternalIPs []string
	ByPort      map[string]int
}

type TraefikInfo struct {
	Error4xx      int
	Error5xx      int
	TotalRequests int
}

var knownPorts map[int]bool

func CollectNetwork() NetworkData {
	var d NetworkData

	ports := getListeningPorts()
	d.ListeningPorts = len(ports)
	d.NewPorts = detectNewPorts(ports)
	d.Connections = getEstablishedConnections()
	d.Traefik = getTraefikErrors()

	return d
}

// Pre-compiled regexes for network collectors
var (
	netPortRe   = regexp.MustCompile(`:(\d+)$`)
	netProcRe   = regexp.MustCompile(`"([^"]+)"`)
	netRemoteRe = regexp.MustCompile(`([\d.]+):(\d+)$`)
	netLocalRe  = regexp.MustCompile(`:(\d+)$`)
	netStatusRe = regexp.MustCompile(`" (\d{3}) `)
)

func getListeningPorts() []PortInfo {
	var ports []PortInfo
	out, err := exec.Command("ss", "-tlnp").Output()
	if err != nil {
		return ports
	}

	for _, line := range strings.Split(string(out), "\n")[1:] {
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		addr := parts[3]
		m := netPortRe.FindStringSubmatch(addr)
		if m == nil {
			continue
		}
		port, _ := strconv.Atoi(m[1])
		procName := "unknown"
		if len(parts) >= 6 {
			pm := netProcRe.FindStringSubmatch(parts[len(parts)-1])
			if pm != nil {
				procName = pm[1]
			}
		}
		ports = append(ports, PortInfo{Port: port, Address: addr, Process: procName})
	}
	return ports
}

func detectNewPorts(current []PortInfo) []PortInfo {
	currentSet := make(map[int]bool)
	for _, p := range current {
		currentSet[p.Port] = true
	}

	if knownPorts == nil {
		knownPorts = currentSet
		return nil
	}

	var newPorts []PortInfo
	for _, p := range current {
		if !knownPorts[p.Port] {
			newPorts = append(newPorts, p)
		}
	}
	knownPorts = currentSet
	return newPorts
}

func getEstablishedConnections() ConnectionInfo {
	info := ConnectionInfo{ByPort: make(map[string]int)}
	ips := make(map[string]bool)

	out, err := exec.Command("ss", "-tnp", "state", "established").Output()
	if err != nil {
		return info
	}

	remoteRe := netRemoteRe
	localRe := netLocalRe

	for _, line := range strings.Split(string(out), "\n")[1:] {
		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}
		info.Total++

		if m := remoteRe.FindStringSubmatch(parts[4]); m != nil {
			ips[m[1]] = true
		}
		if m := localRe.FindStringSubmatch(parts[3]); m != nil {
			info.ByPort[m[1]]++
		}
	}

	info.UniqueIPs = len(ips)
	for ip := range ips {
		if len(info.ExternalIPs) < 20 {
			info.ExternalIPs = append(info.ExternalIPs, ip)
		}
	}
	return info
}

func getTraefikErrors() TraefikInfo {
	// Cache Traefik errors for 5 minutes to avoid frequent docker logs calls
	cachedTraefikMu.RLock()
	if time.Since(cachedTraefikTime) < traefikCacheTTL && cachedTraefikErrors != nil {
		result := *cachedTraefikErrors
		cachedTraefikMu.RUnlock()
		return result
	}
	cachedTraefikMu.RUnlock()

	var info TraefikInfo
	out, err := exec.Command("docker", "logs", "coolify-proxy", "--tail", "200", "--since", "1h").CombinedOutput()
	if err != nil {
		return info
	}

	statusRe := netStatusRe
	for _, line := range strings.Split(string(out), "\n") {
		m := statusRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		info.TotalRequests++
		code, _ := strconv.Atoi(m[1])
		if code >= 400 && code < 500 {
			info.Error4xx++
		} else if code >= 500 {
			info.Error5xx++
		}
	}

	cachedTraefikMu.Lock()
	cachedTraefikErrors = &info
	cachedTraefikTime = time.Now()
	cachedTraefikMu.Unlock()

	return info
}

var (
	cachedTraefikErrors *TraefikInfo
	cachedTraefikTime   time.Time
	cachedTraefikMu     sync.RWMutex
	traefikCacheTTL     = 5 * time.Minute
)

// ──────────────────────────────────────────────
// Storage collector
// ──────────────────────────────────────────────

type StorageData struct {
	GDriveMount   GDriveInfo
	S3Gateway     S3Info
	RcloneService RcloneInfo
	DockerDisk    []DockerDFItem
}

type GDriveInfo struct {
	Path       string
	Mounted    bool
	Accessible bool
}

type S3Info struct {
	URL        string
	Reachable  bool
	StatusCode int
	Error      string
}

type RcloneInfo struct {
	GDriveService string
	CacheInfo     string
}

type DockerDFItem struct {
	Type        string
	Total       string
	Active      string
	Size        string
	Reclaimable string
}

func CollectStorage() StorageData {
	var d StorageData
	var wg sync.WaitGroup
	wg.Add(4)
	go func() { defer wg.Done(); d.GDriveMount = checkGDriveMount() }()
	go func() { defer wg.Done(); d.S3Gateway = checkS3Gateway() }()
	go func() { defer wg.Done(); d.RcloneService = getRcloneServiceStatus() }()
	go func() { defer wg.Done(); d.DockerDisk = getDockerVolumeSizes() }()
	wg.Wait()
	return d
}

func checkGDriveMount() GDriveInfo {
	info := GDriveInfo{Path: config.Cfg.RcloneGDriveMount}
	fi, err := os.Stat(config.Cfg.RcloneGDriveMount)
	if err != nil || !fi.IsDir() {
		return info
	}

	// Check if it's a mount point
	_, err = exec.Command("mountpoint", "-q", config.Cfg.RcloneGDriveMount).CombinedOutput()
	info.Mounted = err == nil

	if info.Mounted {
		_, err = os.ReadDir(config.Cfg.RcloneGDriveMount)
		info.Accessible = err == nil
	}
	return info
}

func checkS3Gateway() S3Info {
	info := S3Info{URL: config.Cfg.RcloneS3GatewayURL}
	client := httpShort
	resp, err := client.Get(config.Cfg.RcloneS3GatewayURL)
	if err != nil {
		info.Error = err.Error()
		return info
	}
	defer resp.Body.Close()
	info.Reachable = true
	info.StatusCode = resp.StatusCode
	return info
}

func getRcloneServiceStatus() RcloneInfo {
	info := RcloneInfo{GDriveService: "unknown"}

	out, err := exec.Command("systemctl", "is-active", "rclone-gdrive").Output()
	if err == nil {
		info.GDriveService = strings.TrimSpace(string(out))
	}

	cachePath := filepath.Join(config.HomeDir(), ".cache", "rclone")
	if _, err := os.Stat(cachePath); err == nil {
		out, err := exec.Command("du", "-sh", cachePath).Output()
		if err == nil && len(out) > 0 {
			parts := strings.Split(string(out), "\t")
			if len(parts) > 0 {
				info.CacheInfo = strings.TrimSpace(parts[0])
			}
		}
	}
	return info
}

func getDockerVolumeSizes() []DockerDFItem {
	var items []DockerDFItem
	out, err := exec.Command("docker", "system", "df").Output()
	if err != nil {
		return items
	}
	for _, line := range strings.Split(string(out), "\n")[1:] {
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		item := DockerDFItem{
			Type:   parts[0],
			Total:  parts[1],
			Active: parts[2],
			Size:   parts[3],
		}
		if len(parts) > 4 {
			item.Reclaimable = parts[4]
		} else {
			item.Reclaimable = "N/A"
		}
		items = append(items, item)
	}
	return items
}
