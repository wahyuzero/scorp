package collectors

import (
	"scorp-agent/config"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Security types
type SecurityEvent struct {
	Type   string
	User   string
	IP     string
	Port   string
	Method string
	Time   string
	Cmd    string
}

type SecuritySummary struct {
	SSHLoginsCount int
	SSHFailedCount int
	SSHUniqueIPs   int
	TotalBannedIPs int
	VNCConnections []SecurityEvent
	Fail2ban       Fail2banInfo
	BruteForce     []BruteForceAlert
	RecentFailed   []RecentFailedEntry
}

// RecentFailedEntry is a single failed SSH attempt with geo info.
type RecentFailedEntry struct {
	IP   string
	User string
	Time string
	Geo  GeoInfo
}

type Fail2banInfo struct {
	Active bool
	Jails  map[string]JailInfo
}

type JailInfo struct {
	Banned      int
	TotalBanned int
	Failed      int
}

type BruteForceAlert struct {
	IP        string
	Attempts  int
	WindowMin int
}

// ──────────────────────────────────────────────
// Dedup + brute force tracking
// ──────────────────────────────────────────────

var (
	seenEvents     = make(map[string]time.Time)
	seenMu         sync.Mutex
	dedupCallCount int

	failedAttempts   = make(map[string][]time.Time)
	failedAttemptsMu sync.Mutex

	EventChan = make(chan SecurityEvent, 100)

	// failedSSHBuffer accumulates failed SSH events for hourly digest
	failedSSHBuffer   []RecentFailedEntry
	failedSSHBufferMu sync.Mutex
)

func isDuplicate(key string) bool {
	seenMu.Lock()
	defer seenMu.Unlock()

	now := time.Now()
	dedupDuration := time.Duration(config.Cfg.SSHAlertDedup) * time.Second

	// Periodic cleanup
	dedupCallCount++
	if dedupCallCount >= 50 {
		dedupCallCount = 0
		cutoff := now.Add(-dedupDuration)
		for k, t := range seenEvents {
			if t.Before(cutoff) {
				delete(seenEvents, k)
			}
		}
	}

	if t, ok := seenEvents[key]; ok && now.Sub(t) < dedupDuration {
		return true
	}
	seenEvents[key] = now
	return false
}

func trackFailed(ip string) {
	failedAttemptsMu.Lock()
	defer failedAttemptsMu.Unlock()

	now := time.Now()
	cutoff := now.Add(-10 * time.Minute)
	var recent []time.Time
	for _, t := range failedAttempts[ip] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	failedAttempts[ip] = append(recent, now)
}

func CheckBruteForce() []BruteForceAlert {
	failedAttemptsMu.Lock()
	defer failedAttemptsMu.Unlock()

	cutoff := time.Now().Add(-5 * time.Minute)
	var alerts []BruteForceAlert
	for ip, times := range failedAttempts {
		var recent int
		for _, t := range times {
			if t.After(cutoff) {
				recent++
			}
		}
		if recent >= 5 {
			key := fmt.Sprintf("brute_force:%s", ip)
			if !isDuplicate(key) {
				alerts = append(alerts, BruteForceAlert{
					IP:        ip,
					Attempts:  recent,
					WindowMin: 5,
				})
			}
		}
	}
	return alerts
}

// ──────────────────────────────────────────────
// Journal watcher (goroutine)
// ──────────────────────────────────────────────

var (
	sshAcceptRe   = regexp.MustCompile(`Accepted\s+(\w+)\s+for\s+(\w+)\s+from\s+([\d.]+)\s+port\s+(\d+)`)
	sshFailedRe   = regexp.MustCompile(`from\s+([\d.]+)`)
	sshUserRe     = regexp.MustCompile(`for\s+(?:invalid user\s+)?(\w+)`)
	invalidUserRe = regexp.MustCompile(`Invalid user\s+(\S+)\s+from\s+([\d.]+)`)
	timeRe        = regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T[\d:]+)`)
	fail2banBannedRe  = regexp.MustCompile(`Currently banned:\s*(\d+)`)
	fail2banTotalRe   = regexp.MustCompile(`Total banned:\s*(\d+)`)
	fail2banFailedRe  = regexp.MustCompile(`Currently failed:\s*(\d+)`)
	acceptRe          = regexp.MustCompile(`Accepted\s+\w+\s+for\s+(\w+)\s+from\s+([\d.]+)`)
)

func WatchJournal(done <-chan struct{}) {
	for {
		// Check for shutdown before starting new subprocess
		select {
		case <-done:
			log.Println("[security] Journal watcher stopped (shutdown)")
			return
		default:
		}

		cmd := exec.Command("journalctl", "--follow", "--no-pager", "-o", "short-iso",
			"-t", "sshd", "--since", "now")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("[security] journal pipe error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if err := cmd.Start(); err != nil {
			log.Printf("[security] journal start error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("[security] Journal watcher started (PID %d)", cmd.Process.Pid)

		// Kill subprocess on shutdown signal
		go func() {
			<-done
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		}()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			events := parseJournalLine(line)
			for _, ev := range events {
				// Skip whitelisted
				whitelisted := false
				for _, prefix := range config.Cfg.SSHAlertWhitelist {
					if strings.HasPrefix(ev.IP, prefix) {
						whitelisted = true
						break
					}
				}
				if whitelisted {
					continue
				}

				if ev.Type == "ssh_login" {
					// Successful logins → real-time alert
					key := fmt.Sprintf("%s:%s:%s", ev.Type, ev.User, ev.IP)
					if !isDuplicate(key) {
						EventChan <- ev
					}
				} else if ev.Type == "ssh_failed" {
					// Failed attempts → buffer for digest (geo enriched async)
					entry := RecentFailedEntry{
						IP:   ev.IP,
						User: ev.User,
						Time: ev.Time,
					}
					failedSSHBufferMu.Lock()
					failedSSHBuffer = append(failedSSHBuffer, entry)
					if len(failedSSHBuffer) > 200 {
						failedSSHBuffer = failedSSHBuffer[len(failedSSHBuffer)-200:]
					}
					failedSSHBufferMu.Unlock()
				}
			}
		}

		cmd.Wait()
		// Check if we're shutting down
		select {
		case <-done:
			log.Println("[security] Journal watcher stopped (shutdown)")
			return
		default:
		}
		log.Println("[security] Journal watcher process ended, restarting in 5s...")
		time.Sleep(5 * time.Second)
	}
}

func parseJournalLine(line string) []SecurityEvent {
	var events []SecurityEvent
	ts := extractTime(line)

	if strings.Contains(line, "Accepted") {
		m := sshAcceptRe.FindStringSubmatch(line)
		if m != nil {
			events = append(events, SecurityEvent{
				Type: "ssh_login", Method: m[1], User: m[2], IP: m[3], Port: m[4], Time: ts,
			})
		}
	} else if strings.Contains(line, "Failed password") || strings.Contains(strings.ToLower(line), "authentication failure") {
		ip := "unknown"
		if m := sshFailedRe.FindStringSubmatch(line); m != nil {
			ip = m[1]
		}
		user := "unknown"
		if m := sshUserRe.FindStringSubmatch(line); m != nil {
			user = m[1]
		}
		events = append(events, SecurityEvent{Type: "ssh_failed", User: user, IP: ip, Time: ts})
		trackFailed(ip)
	} else if strings.Contains(line, "Invalid user") {
		m := invalidUserRe.FindStringSubmatch(line)
		if m != nil {
			events = append(events, SecurityEvent{Type: "ssh_failed", User: m[1], IP: m[2], Time: ts})
			trackFailed(m[2])
		}
	}

	return events
}

func extractTime(line string) string {
	m := timeRe.FindStringSubmatch(line)
	if m != nil {
		return m[1]
	}
	return ""
}

// ──────────────────────────────────────────────
// VNC monitoring
// ──────────────────────────────────────────────

var vncLastPos int64

func CheckVNCConnections() []SecurityEvent {
	vncLog := config.VNCLogPath()
	var events []SecurityEvent

	fi, err := os.Stat(vncLog)
	if err != nil {
		return events
	}

	size := fi.Size()
	if size < vncLastPos {
		vncLastPos = 0
	}
	if vncLastPos == 0 {
		vncLastPos = size
		return events
	}

	f, err := os.Open(vncLog)
	if err != nil {
		return events
	}
	defer f.Close()

	f.Seek(vncLastPos, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "SConnection: Client needs protocol") {
			key := fmt.Sprintf("vnc_connect:%s", time.Now().Format("15:04"))
			if !isDuplicate(key) {
				events = append(events, SecurityEvent{
					Type: "vnc_connection",
					Time: time.Now().Format("2006-01-02 15:04:05"),
				})
			}
		}
	}
	vncLastPos = size

	return events
}

// ──────────────────────────────────────────────
// Fail2ban + SSH summary (for reports)
// ──────────────────────────────────────────────

func getFail2banStatus() Fail2banInfo {
	info := Fail2banInfo{Jails: make(map[string]JailInfo)}

	out, err := exec.Command("sudo", "fail2ban-client", "status").Output()
	if err != nil {
		return info
	}
	info.Active = true

	jailRe := regexp.MustCompile(`Jail list:\s*(.*)`)
	m := jailRe.FindStringSubmatch(string(out))
	if m == nil {
		return info
	}

	for _, jail := range strings.Split(m[1], ",") {
		jail = strings.TrimSpace(jail)
		if jail == "" {
			continue
		}

		jailOut, err := exec.Command("sudo", "fail2ban-client", "status", jail).Output()
		if err != nil {
			continue
		}
		jStr := string(jailOut)

		var ji JailInfo
		if bm := fail2banBannedRe.FindStringSubmatch(jStr); bm != nil {
			fmt.Sscanf(bm[1], "%d", &ji.Banned)
		}
		if tm := fail2banTotalRe.FindStringSubmatch(jStr); tm != nil {
			fmt.Sscanf(tm[1], "%d", &ji.TotalBanned)
		}
		if fm := fail2banFailedRe.FindStringSubmatch(jStr); fm != nil {
			fmt.Sscanf(fm[1], "%d", &ji.Failed)
		}
		info.Jails[jail] = ji
	}

	return info
}

func GetRecentSSHSummary(sinceMinutes int) (logins, failed, uniqueIPs int) {
	since := time.Now().Add(-time.Duration(sinceMinutes) * time.Minute).Format("2006-01-02 15:04:05")
	out, err := exec.Command("journalctl", "-t", "sshd", "--since", since, "--no-pager", "-o", "short-iso").Output()
	if err != nil {
		return
	}

	seen := make(map[string]bool)
	ips := make(map[string]bool)

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "Accepted") {
			m := acceptRe.FindStringSubmatch(line)
			if m != nil {
				ts := extractTime(line)
				dedup := fmt.Sprintf("%s:%s:%s", m[1], m[2], ts[:16])
				if !seen[dedup] {
					seen[dedup] = true
					logins++
					ips[m[2]] = true
				}
			}
		} else if strings.Contains(line, "Failed") || strings.Contains(line, "Invalid user") {
			failed++
		}
	}
	uniqueIPs = len(ips)
	return
}

func CollectSecurity() SecuritySummary {
	logins, failed, uniqueIPs := GetRecentSSHSummary(60)
	f2b := getFail2banStatus()
	bf := CheckBruteForce()
	recent := DrainFailedSSHBuffer()

	totalBanned := 0
	for _, j := range f2b.Jails {
		totalBanned += j.Banned
	}

	return SecuritySummary{
		SSHLoginsCount: logins,
		SSHFailedCount: failed,
		SSHUniqueIPs:   uniqueIPs,
		TotalBannedIPs: totalBanned,
		Fail2ban:       f2b,
		BruteForce:     bf,
		RecentFailed:   recent,
	}
}

// collectSecurityWithPeek returns security data without draining the buffer
func CollectSecurityWithPeek() SecuritySummary {
	logins, failed, uniqueIPs := GetRecentSSHSummary(60)
	f2b := getFail2banStatus()
	bf := CheckBruteForce()
	recent := PeekFailedSSHBuffer()

	totalBanned := 0
	for _, j := range f2b.Jails {
		totalBanned += j.Banned
	}

	return SecuritySummary{
		SSHLoginsCount: logins,
		SSHFailedCount: failed,
		SSHUniqueIPs:   uniqueIPs,
		TotalBannedIPs: totalBanned,
		Fail2ban:       f2b,
		BruteForce:     bf,
		RecentFailed:   recent,
	}
}

// drainFailedSSHBuffer returns the last 10 events with geo, then clears buffer.
func DrainFailedSSHBuffer() []RecentFailedEntry {
	failedSSHBufferMu.Lock()
	buf := failedSSHBuffer
	failedSSHBuffer = nil
	failedSSHBufferMu.Unlock()

	return enrichLast10(buf)
}

// peekFailedSSHBuffer returns the last 10 events with geo WITHOUT clearing.
func PeekFailedSSHBuffer() []RecentFailedEntry {
	failedSSHBufferMu.Lock()
	buf := make([]RecentFailedEntry, len(failedSSHBuffer))
	copy(buf, failedSSHBuffer)
	failedSSHBufferMu.Unlock()

	return enrichLast10(buf)
}

// enrichLast10 takes the last 10 events from the buffer and enriches with geo.
func enrichLast10(buf []RecentFailedEntry) []RecentFailedEntry {
	if len(buf) == 0 {
		return nil
	}

	// Take last 10 (most recent)
	start := 0
	if len(buf) > 10 {
		start = len(buf) - 10
	}
	recent := buf[start:]

	// Reverse so newest is first
	result := make([]RecentFailedEntry, len(recent))
	for i, e := range recent {
		e.Geo = LookupIP(e.IP)
		result[len(recent)-1-i] = e
	}

	log.Printf("[security] Returning %d recent failed SSH events (buffer had %d)", len(result), len(buf))
	return result
}
