package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Cooldown tracking
var (
	alertCooldowns = make(map[string]time.Time)
	alertMu        sync.Mutex
	alertCallCount int
)

func canFire(key string) bool {
	alertMu.Lock()
	defer alertMu.Unlock()

	now := time.Now()
	cooldown := time.Duration(cfg.AlertCooldown) * time.Second

	// Periodic cleanup every 100 calls
	alertCallCount++
	if alertCallCount >= 100 {
		alertCallCount = 0
		for k, t := range alertCooldowns {
			if now.Sub(t) > cooldown {
				delete(alertCooldowns, k)
			}
		}
	}

	if last, ok := alertCooldowns[key]; ok {
		if now.Sub(last) < cooldown {
			return false
		}
	}
	alertCooldowns[key] = now
	return true
}

func checkSystemAlerts(d SystemData) []string {
	var alerts []string

	if d.CPUPercent >= cfg.CPUThreshold && canFire("cpu_high") {
		topStr := formatTopProcs(d.TopProcesses)
		alerts = append(alerts, fmt.Sprintf("🔴 <b>HIGH CPU USAGE</b>\n"+
			"CPU: <b>%.1f%%</b> (threshold: %.0f%%)\n"+
			"Load: %.1f, %.1f, %.1f\n"+
			"Top: %s",
			d.CPUPercent, cfg.CPUThreshold,
			d.LoadAvg[0], d.LoadAvg[1], d.LoadAvg[2], topStr))
	}

	if d.RAMPercent >= cfg.RAMThreshold && canFire("ram_high") {
		alerts = append(alerts, fmt.Sprintf("🟠 <b>HIGH MEMORY USAGE</b>\n"+
			"RAM: <b>%.1fG / %.1fG (%.1f%%)</b>\n"+
			"Available: %.1fG",
			d.RAMUsedGB, d.RAMTotalGB, d.RAMPercent, d.RAMAvailGB))
	}

	if d.DiskPercent >= cfg.DiskThreshold && canFire("disk_high") {
		alerts = append(alerts, fmt.Sprintf("🟡 <b>HIGH DISK USAGE</b>\n"+
			"Disk: <b>%.1fG / %.1fG (%.1f%%)</b>",
			d.DiskUsedGB, d.DiskTotalGB, d.DiskPercent))
	}

	if d.LoadAvg[0] >= cfg.LoadThreshold && canFire("load_high") {
		alerts = append(alerts, fmt.Sprintf("⚡ <b>HIGH LOAD AVERAGE</b>\n"+
			"Load: <b>%.1f</b> (threshold: %.1f, cores: %d)",
			d.LoadAvg[0], cfg.LoadThreshold, d.CPUCount))
	}

	swapMB := d.SwapUsedGB * 1024
	if swapMB >= cfg.SwapThresholdMB && canFire("swap_high") {
		alerts = append(alerts, fmt.Sprintf("🟣 <b>SWAP USAGE DETECTED</b>\n"+
			"Swap: <b>%.1fG / %.1fG</b>\n"+
			"This indicates RAM pressure.",
			d.SwapUsedGB, d.SwapTotalGB))
	}

	return alerts
}

func checkDockerAlerts(d DockerData) []string {
	var alerts []string
	for _, c := range d.Containers {
		if c.Status != "running" && canFire(fmt.Sprintf("container_down_%s", c.Name)) {
			alerts = append(alerts, fmt.Sprintf("🐳 <b>CONTAINER DOWN</b>\n"+
				"Name: <code>%s</code>\nStatus: %s\nImage: %s",
				c.Name, c.Status, c.Image))
		}
		if c.Health == "unhealthy" && canFire(fmt.Sprintf("container_unhealthy_%s", c.Name)) {
			alerts = append(alerts, fmt.Sprintf("🐳 <b>CONTAINER UNHEALTHY</b>\n"+
				"Name: <code>%s</code>\nHealth: %s",
				c.Name, c.Health))
		}
	}
	return alerts
}

func checkStorageAlerts(d StorageData) []string {
	var alerts []string
	if !d.GDriveMount.Mounted && canFire("gdrive_unmounted") {
		alerts = append(alerts, fmt.Sprintf("📁 <b>GDRIVE UNMOUNTED</b>\n"+
			"Path: <code>%s</code>\nThe rclone Google Drive mount is down!",
			d.GDriveMount.Path))
	}
	if !d.S3Gateway.Reachable && canFire("s3_unreachable") {
		alerts = append(alerts, fmt.Sprintf("📦 <b>S3 GATEWAY DOWN</b>\n"+
			"URL: <code>%s</code>\nError: %s",
			d.S3Gateway.URL, d.S3Gateway.Error))
	}
	return alerts
}

func checkNetworkAlerts(d NetworkData) []string {
	var alerts []string
	for _, p := range d.NewPorts {
		if canFire(fmt.Sprintf("new_port_%d", p.Port)) {
			alerts = append(alerts, fmt.Sprintf("🔓 <b>NEW PORT DETECTED</b>\n"+
				"Port: <b>%d</b>\nProcess: <code>%s</code>\nAddress: %s",
				p.Port, p.Process, p.Address))
		}
	}
	if d.Traefik.Error5xx >= 10 && canFire("traefik_5xx") {
		alerts = append(alerts, fmt.Sprintf("⚠️ <b>HIGH 5xx ERROR RATE</b>\n"+
			"5xx errors: <b>%d</b> in last hour\nTotal requests: %d",
			d.Traefik.Error5xx, d.Traefik.TotalRequests))
	}
	return alerts
}

func formatTopProcs(procs []TopProcess) string {
	if len(procs) == 0 {
		return "N/A"
	}
	var lines []string
	for i := 0; i < 3 && i < len(procs); i++ {
		p := procs[i]
		lines = append(lines, fmt.Sprintf("  %s (CPU: %.1f%%, MEM: %.1f%%)", p.Name, p.CPU, p.Mem))
	}
	return strings.Join(lines, "\n")
}
