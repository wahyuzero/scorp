package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Hourly Report
// ──────────────────────────────────────────────

func formatHourlyReport(sys SystemData, docker DockerData, coolify CoolifyData,
	sec SecuritySummary, stor StorageData, net NetworkData) string {

	now := time.Now().UTC().Format("02 Jan 2006, 15:04 UTC")
	sections := []string{
		fmt.Sprintf("📊 <b>VPS Status Report</b>\n🕐 %s\n%s", now, strings.Repeat("─", 28)),
		sectionSystem(sys),
		sectionDocker(docker),
		sectionCoolify(coolify),
		sectionSecurity(sec),
		sectionStorage(stor),
		sectionNetwork(net),
	}
	return strings.Join(sections, "\n\n")
}

// ──────────────────────────────────────────────
// Quick Status
// ──────────────────────────────────────────────

func formatStatusResponse(sys SystemData, docker DockerData, stor StorageData) string {
	now := time.Now().UTC().Format("15:04 UTC")

	gdrive := "❌"
	if stor.GDriveMount.Mounted {
		gdrive = "✅"
	}
	s3 := "❌"
	if stor.S3Gateway.Reachable {
		s3 = "✅"
	}

	unhealthy := ""
	if docker.Summary.Unhealthy > 0 {
		unhealthy = fmt.Sprintf(" ⚠️ %d unhealthy", docker.Summary.Unhealthy)
	}

	return fmt.Sprintf("⚡ <b>Quick Status</b> (%s)\n\n"+
		"CPU:  %s %.1f%%\n"+
		"RAM:  %s %.1fG/%.1fG\n"+
		"Disk: %s %.1fG/%.1fG\n"+
		"Up:   %s\n\n"+
		"🐳 Containers: %d/%d running%s\n"+
		"📁 GDrive: %s | S3: %s",
		now,
		bar(sys.CPUPercent), sys.CPUPercent,
		bar(sys.RAMPercent), sys.RAMUsedGB, sys.RAMTotalGB,
		bar(sys.DiskPercent), sys.DiskUsedGB, sys.DiskTotalGB,
		sys.Uptime,
		docker.Summary.Running, docker.Summary.Total, unhealthy,
		gdrive, s3,
	)
}

// ──────────────────────────────────────────────
// Sections
// ──────────────────────────────────────────────

func sectionSystem(d SystemData) string {
	return fmt.Sprintf("<b>⚡ System</b>\n"+
		"CPU:  %s %.1f%%\n"+
		"Load: %.1f, %.1f, %.1f\n"+
		"RAM:  %s %.1fG / %.1fG (%.1f%%)\n"+
		"Swap: %.1fG / %.1fG\n"+
		"Disk: %s %.1fG / %.1fG (%.1f%%)\n"+
		"Net:  ⬆️ %.2fG ⬇️ %.2fG\n"+
		"Up:   %s",
		bar(d.CPUPercent), d.CPUPercent,
		d.LoadAvg[0], d.LoadAvg[1], d.LoadAvg[2],
		bar(d.RAMPercent), d.RAMUsedGB, d.RAMTotalGB, d.RAMPercent,
		d.SwapUsedGB, d.SwapTotalGB,
		bar(d.DiskPercent), d.DiskUsedGB, d.DiskTotalGB, d.DiskPercent,
		d.NetSentGB, d.NetRecvGB,
		d.Uptime,
	)
}

func sectionDocker(d DockerData) string {
	header := fmt.Sprintf("<b>🐳 Docker</b> (%d/%d running)", d.Summary.Running, d.Summary.Total)
	if d.Summary.Unhealthy > 0 {
		header += fmt.Sprintf(" ⚠️ %d unhealthy", d.Summary.Unhealthy)
	}
	lines := []string{header}

	sort.Slice(d.Containers, func(i, j int) bool {
		return d.Containers[i].Name < d.Containers[j].Name
	})

	for _, c := range d.Containers {
		icon := "✅"
		if c.Status != "running" {
			icon = "❌"
		} else if c.Health == "unhealthy" {
			icon = "⚠️"
		}
		stats := ""
		cpu := ""
		mem := ""
		if c.CPUPercent > 0 {
			cpu = fmt.Sprintf("%.1f%%", c.CPUPercent)
		}
		if c.MemoryMB > 0 {
			mem = fmt.Sprintf("%.0fM", c.MemoryMB)
		}
		if cpu != "" || mem != "" {
			stats = fmt.Sprintf(" (%s, %s)", cpu, mem)
		}
		lines = append(lines, fmt.Sprintf("  %s <code>%s</code>%s", icon, c.Name, stats))
	}
	return strings.Join(lines, "\n")
}

func sectionCoolify(d CoolifyData) string {
	if !d.Available {
		return "<b>☁️ Coolify</b>\n  ⚠️ API unavailable"
	}

	header := "<b>☁️ Coolify</b>"
	if d.Version != "" {
		header = fmt.Sprintf("<b>☁️ Coolify</b> (v%s)", d.Version)
	}
	lines := []string{header}

	for _, s := range d.Servers {
		icon := "✅"
		if !s.Reachable {
			icon = "❌"
		}
		lines = append(lines, fmt.Sprintf("  %s Server: %s (%s)", icon, s.Name, s.IP))
	}

	if len(d.Applications) > 0 {
		lines = append(lines, fmt.Sprintf("  Apps (%d):", len(d.Applications)))
		for _, a := range d.Applications {
			icon := "❌"
			if a.Status == "running" && a.Health == "healthy" {
				icon = "✅"
			} else if a.Status == "running" {
				icon = "🟡"
			}
			health := ""
			if a.Health != "" {
				health = fmt.Sprintf(" [%s]", a.Health)
			}
			fqdn := ""
			if a.FQDN != "" {
				fqdn = fmt.Sprintf(" → %s", a.FQDN)
			}
			lines = append(lines, fmt.Sprintf("    %s %s%s%s", icon, a.Name, health, fqdn))
		}
	}

	if len(d.Databases) > 0 {
		lines = append(lines, fmt.Sprintf("  Databases (%d):", len(d.Databases)))
		for _, db := range d.Databases {
			icon := "✅"
			if db.Status != "running" {
				icon = "⚠️"
			}
			lines = append(lines, fmt.Sprintf("    %s %s (%s)", icon, db.Name, db.Type))
		}
	}

	if len(d.Services) > 0 {
		lines = append(lines, fmt.Sprintf("  Services (%d):", len(d.Services)))
		for _, s := range d.Services {
			icon := "✅"
			if s.Status != "running" {
				icon = "⚠️"
			}
			lines = append(lines, fmt.Sprintf("    %s %s", icon, s.Name))
		}
	}

	if len(lines) == 1 {
		lines = append(lines, "  No resources found via API")
	}
	return strings.Join(lines, "\n")
}

func sectionSecurity(d SecuritySummary) string {
	lines := []string{
		"<b>🔐 Security</b> (last hour)",
	}

	// Summary totals
	lines = append(lines, fmt.Sprintf("  ✅ Successful logins: <b>%d</b>", d.SSHLoginsCount))
	lines = append(lines, fmt.Sprintf("  ❌ Failed attempts: <b>%d</b>", d.SSHFailedCount))
	lines = append(lines, fmt.Sprintf("  🌐 Unique IPs: %d", d.SSHUniqueIPs))
	lines = append(lines, fmt.Sprintf("  🚫 Banned IPs: %d", d.TotalBannedIPs))

	if d.Fail2ban.Active {
		for jail, info := range d.Fail2ban.Jails {
			lines = append(lines, fmt.Sprintf("  🛡 fail2ban [%s]: %d banned, %d failed (total: %d)",
				jail, info.Banned, info.Failed, info.TotalBanned))
		}
	}

	if len(d.BruteForce) > 0 {
		lines = append(lines, fmt.Sprintf("  🚨 Brute force: %d IPs detected!", len(d.BruteForce)))
	}

	// Last 10 individual failed attempts
	if len(d.RecentFailed) > 0 {
		lines = append(lines, fmt.Sprintf("\n  <b>📋 Last %d failed attempts:</b>", len(d.RecentFailed)))
		for _, e := range d.RecentFailed {
			ts := e.Time
			if len(ts) > 16 {
				ts = ts[11:16] // Extract HH:MM from ISO timestamp
			}
			user := e.User
			if user == "" || user == "unknown" {
				user = "?"
			}
			geo := ""
			if e.Geo.Country != "" {
				geo = fmt.Sprintf(" %s %s", flag(e.Geo.CC), e.Geo.Country)
				if e.Geo.City != "" {
					geo += ", " + e.Geo.City
				}
			}
			lines = append(lines, fmt.Sprintf("    %s <code>%s</code> → %s%s", ts, e.IP, user, geo))
		}
	} else {
		lines = append(lines, "\n  ✨ No recent failed attempts in buffer")
	}

	return strings.Join(lines, "\n")
}

func sectionStorage(d StorageData) string {
	gdIcon := "❌"
	if d.GDriveMount.Mounted && d.GDriveMount.Accessible {
		gdIcon = "✅"
	}
	s3Icon := "❌"
	if d.S3Gateway.Reachable {
		s3Icon = "✅"
	}
	svcIcon := "❌"
	if d.RcloneService.GDriveService == "active" {
		svcIcon = "✅"
	}

	gdStatus := "UNMOUNTED"
	if d.GDriveMount.Mounted {
		gdStatus = "OK"
	}
	s3Status := "DOWN"
	if d.S3Gateway.Reachable {
		s3Status = "OK"
	}

	lines := []string{
		"<b>📁 Storage</b>",
		fmt.Sprintf("  %s GDrive mount: %s", gdIcon, gdStatus),
		fmt.Sprintf("  %s Rclone service: %s", svcIcon, d.RcloneService.GDriveService),
		fmt.Sprintf("  %s S3 Gateway: %s", s3Icon, s3Status),
	}

	if d.RcloneService.CacheInfo != "" {
		lines = append(lines, fmt.Sprintf("  Cache: %s", d.RcloneService.CacheInfo))
	}

	if len(d.DockerDisk) > 0 {
		lines = append(lines, "  Docker disk:")
		for _, item := range d.DockerDisk {
			lines = append(lines, fmt.Sprintf("    %s: %s (reclaimable: %s)", item.Type, item.Size, item.Reclaimable))
		}
	}

	return strings.Join(lines, "\n")
}

func sectionNetwork(d NetworkData) string {
	lines := []string{
		"<b>🌐 Network</b>",
		fmt.Sprintf("  Listening ports: %d", d.ListeningPorts),
		fmt.Sprintf("  Active connections: %d (%d unique IPs)", d.Connections.Total, d.Connections.UniqueIPs),
	}

	if d.Traefik.TotalRequests > 0 {
		lines = append(lines, fmt.Sprintf("  Traefik: %d req | 4xx: %d | 5xx: %d",
			d.Traefik.TotalRequests, d.Traefik.Error4xx, d.Traefik.Error5xx))
	}

	if len(d.NewPorts) > 0 {
		ports := make([]string, 0)
		for _, p := range d.NewPorts {
			ports = append(ports, fmt.Sprintf("%d", p.Port))
		}
		lines = append(lines, fmt.Sprintf("  ⚠️ New ports: %s", strings.Join(ports, ", ")))
	}

	return strings.Join(lines, "\n")
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

func bar(percent float64) string {
	width := 10
	filled := int(float64(width) * percent / 100)
	if filled > width {
		filled = width
	}
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", width-filled) + "]"
}
