// +build linux

package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Native Linux /proc collectors (replaces gopsutil)
// ──────────────────────────────────────────────

// Background CPU sampler using /proc/stat
var (
	nativeCPUPercent   float64
	nativeCPUMu        sync.RWMutex
)

func startCPUSampler() {
	go func() {
		prevIdle, prevTotal := readCPUStat()
		for {
			time.Sleep(3 * time.Second)
			idle, total := readCPUStat()
			deltaIdle := idle - prevIdle
			deltaTotal := total - prevTotal
			if deltaTotal > 0 {
				nativeCPUMu.Lock()
				nativeCPUPercent = 100.0 * (1.0 - float64(deltaIdle)/float64(deltaTotal))
				nativeCPUMu.Unlock()
			}
			prevIdle, prevTotal = idle, total
		}
	}()
}

// readCPUStat returns (idle, total) ticks from /proc/stat
func readCPUStat() (int64, int64) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			return 0, 0
		}
		var total int64
		for _, f := range fields[1:] {
			v, _ := strconv.ParseInt(f, 10, 64)
			total += v
		}
		idle, _ := strconv.ParseInt(fields[4], 10, 64)
		return idle, total
	}
	return 0, 0
}

// getCPUCount reads CPU count from /proc/cpuinfo
func getCPUCount() int {
	count := 0
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return 1
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "processor") {
			count++
		}
	}
	if count == 0 {
		count = 1
	}
	return count
}

// getLoadAvg reads /proc/loadavg
func getLoadAvg() [3]float64 {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return [3]float64{}
	}
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return [3]float64{}
	}
	var avg [3]float64
	avg[0], _ = strconv.ParseFloat(fields[0], 64)
	avg[1], _ = strconv.ParseFloat(fields[1], 64)
	avg[2], _ = strconv.ParseFloat(fields[2], 64)
	return avg
}

// getMemInfo reads /proc/meminfo for RAM + Swap
func getMemInfo() (totalKB, availKB, usedKB, swapTotalKB, swapFreeKB uint64) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, 0, 0
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "MemTotal:"):
			totalKB = parseMemInfoValue(line)
		case strings.HasPrefix(line, "MemAvailable:"):
			availKB = parseMemInfoValue(line)
		case strings.HasPrefix(line, "SwapTotal:"):
			swapTotalKB = parseMemInfoValue(line)
		case strings.HasPrefix(line, "SwapFree:"):
			swapFreeKB = parseMemInfoValue(line)
		}
	}
	usedKB = totalKB - availKB
	return
}

func parseMemInfoValue(line string) uint64 {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0
	}
	v, _ := strconv.ParseUint(fields[1], 10, 64)
	return v
}

// getUptime reads /proc/uptime
func getUptime() time.Duration {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return 0
	}
	secs, _ := strconv.ParseFloat(fields[0], 64)
	return time.Duration(secs) * time.Second
}

// getNetBytes reads /proc/net/dev and sums all interfaces
func getNetBytes() (sent, recv uint64) {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return 0, 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		iface := strings.TrimRight(fields[0], ":")
		if iface == "lo" {
			continue
		}
		r, _ := strconv.ParseUint(fields[1], 10, 64)
		s, _ := strconv.ParseUint(fields[9], 10, 64)
		recv += r
		sent += s
	}
	return
}

// getDiskUsage runs `df -B1 /` and parses the output
func getDiskUsage() (total, used, pct uint64) {
	out, err := exec.Command("df", "-B1", "/").Output()
	if err != nil {
		return 0, 0, 0
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return 0, 0, 0
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return 0, 0, 0
	}
	total, _ = strconv.ParseUint(fields[1], 10, 64)
	used, _ = strconv.ParseUint(fields[2], 10, 64)
	pctStr := strings.TrimRight(fields[4], "%")
	p, _ := strconv.ParseUint(pctStr, 10, 64)
	pct = p
	return
}

// ──────────────────────────────────────────────
// Process list using /proc (replaces gopsutil process)
// ──────────────────────────────────────────────

var (
	procCacheTime time.Time
	procCacheMu   sync.RWMutex
	procCache     []nativeTopProcess
)

type nativeTopProcess struct {
	PID  int32
	Name string
	CPU  float64
	Mem  float64 // percent
}

// readProcessList reads /proc for all processes
func readProcessList() []nativeTopProcess {
	procCacheMu.RLock()
	if time.Since(procCacheTime) < 5*time.Second && procCache != nil {
		cp := procCache
		procCacheMu.RUnlock()
		return cp
	}
	procCacheMu.RUnlock()

	// Get total RAM from /proc/meminfo for memory percent calculation
	totalRAMKB, _, _, _, _ := getMemInfo()

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}

	// Get current CPU ticks for calculation
	_, totalTicks := readCPUStat()

	var procs []nativeTopProcess
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.ParseInt(e.Name(), 10, 32)
		if err != nil {
			continue
		}

		name := readProcName(int32(pid))
		if name == "" {
			continue
		}
		procTicks := readProcStat(int32(pid))
		uptimeTicks := float64(totalTicks)
		cpuPct := 0.0
		if uptimeTicks > 0 {
			cpuPct = 100.0 * float64(procTicks) / uptimeTicks
		}

		rssKB := readProcRSS(int32(pid))
		memPct := 0.0
		if totalRAMKB > 0 {
			memPct = 100.0 * float64(rssKB) / float64(totalRAMKB)
		}

		procs = append(procs, nativeTopProcess{
			PID:  int32(pid),
			Name: name,
			CPU:  cpuPct,
			Mem:  memPct,
		})

		// Limit to first 200 processes to avoid memory issues
		if len(procs) >= 200 {
			break
		}
	}

	// Sort by CPU descending (insertion sort for top N)
	sortByCPUDesc(procs)

	procCacheMu.Lock()
	procCache = procs
	procCacheTime = time.Now()
	procCacheMu.Unlock()

	return procs
}

func readProcName(pid int32) string {
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(int(pid)), "comm"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// readProcStat returns (utime+stime, starttime) in ticks
func readProcStat(pid int32) (totalTicks int64) {
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(int(pid)), "stat"))
	if err != nil {
		return 0
	}
	// Find the last ')' in the comm field to handle spaces in comm name
	line := string(data)
	idx := strings.LastIndex(line, ")")
	if idx < 0 {
		return 0
	}
	rest := strings.Fields(line[idx+2:]) // skip ") "
	if len(rest) < 20 {
		return 0
	}
	utime, _ := strconv.ParseInt(rest[11], 10, 64)
	stime, _ := strconv.ParseInt(rest[12], 10, 64)
	return utime + stime
}

func readProcRSS(pid int32) uint64 {
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(int(pid)), "status"))
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "VmRSS:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				v, _ := strconv.ParseUint(fields[1], 10, 64)
				return v
			}
		}
	}
	return 0
}

func sortByCPUDesc(procs []nativeTopProcess) {
	for i := 0; i < len(procs); i++ {
		for j := i + 1; j < len(procs); j++ {
			if procs[j].CPU > procs[i].CPU {
				procs[i], procs[j] = procs[j], procs[i]
			}
		}
	}
}

// getTopProcesses returns top N processes by CPU
func getTopProcesses(n int) []TopProcess {
	procs := readProcessList()
	if len(procs) == 0 {
		return nil
	}
	if n > len(procs) {
		n = len(procs)
	}
	top := make([]TopProcess, n)
	for i := 0; i < n; i++ {
		top[i] = TopProcess{
			PID:  procs[i].PID,
			Name: procs[i].Name,
			CPU:  procs[i].CPU,
			Mem:  procs[i].Mem,
		}
	}
	return top
}

// collectSystem — native Linux implementation
func collectSystem() SystemData {
	var d SystemData

	nativeCPUMu.RLock()
	d.CPUPercent = nativeCPUPercent
	nativeCPUMu.RUnlock()

	d.CPUCount = getCPUCount()
	d.LoadAvg = getLoadAvg()

	totalKB, availKB, usedKB, swapTotalKB, swapFreeKB := getMemInfo()

	const kbToGB = 1024 * 1024
	d.RAMTotalGB = float64(totalKB) / float64(kbToGB)
	d.RAMUsedGB = float64(usedKB) / float64(kbToGB)
	d.RAMAvailGB = float64(availKB) / float64(kbToGB)
	if totalKB > 0 {
		d.RAMPercent = 100.0 * float64(usedKB) / float64(totalKB)
	}

	swapUsedKB := swapTotalKB - swapFreeKB
	d.SwapTotalGB = float64(swapTotalKB) / float64(kbToGB)
	d.SwapUsedGB = float64(swapUsedKB) / float64(kbToGB)
	if swapTotalKB > 0 {
		d.SwapPercent = 100.0 * float64(swapUsedKB) / float64(swapTotalKB)
	}

	diskTotal, diskUsed, diskPct := getDiskUsage()
	d.DiskTotalGB = float64(diskTotal) / float64(kbToGB)
	d.DiskUsedGB = float64(diskUsed) / float64(kbToGB)
	d.DiskPercent = float64(diskPct)

	sent, recv := getNetBytes()
	d.NetSentGB = float64(sent) / float64(kbToGB)
	d.NetRecvGB = float64(recv) / float64(kbToGB)

	uptime := getUptime()
	d.Uptime = formatDuration(uptime)

	d.TopProcesses = getTopProcesses(5)

	return d
}

// getProcesses returns cached process list (interface compat for other callers)
func getProcesses() []nativeTopProcess {
	return readProcessList()
}
