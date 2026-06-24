package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration values.
type Config struct {
	TelegramBotToken string
	TelegramChatID   string
	TelegramWebhookURL string // optional: if set, use webhook mode instead of polling

	CoolifyAPIURL   string
	CoolifyAPIToken string

	CPUThreshold   float64
	RAMThreshold   float64
	DiskThreshold  float64
	LoadThreshold  float64
	SwapThresholdMB float64

	AlertCooldown  int // seconds
	ReportInterval int // seconds

	// Monitoring toggles (agent-first: monitoring is optional)
	MonitoringEnabled     bool // Resource alerts (CPU/RAM/disk) every 30s
	SecurityAlertsEnabled bool // Real-time SSH login/fail alerts
	ScheduledReportsEnabled bool // Periodic full status reports

	RcloneGDriveMount  string
	RcloneS3GatewayURL string

	SSHAlertWhitelist []string
	SSHAlertDedup     int // seconds
}

var Cfg Config

func LoadConfig() error {
	// Load .env file (ignore error if not found — env vars may be set directly)
	godotenv.Load()

	Cfg.TelegramBotToken = EnvStr("TELEGRAM_BOT_TOKEN", "")
	Cfg.TelegramChatID = EnvStr("TELEGRAM_CHAT_ID", "")
	Cfg.TelegramWebhookURL = EnvStr("TELEGRAM_WEBHOOK_URL", "")
	if Cfg.TelegramBotToken == "" || Cfg.TelegramChatID == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID are required")
	}

	Cfg.CoolifyAPIURL = EnvStr("COOLIFY_API_URL", "http://localhost:8000")
	Cfg.CoolifyAPIToken = EnvStr("COOLIFY_API_TOKEN", "")

	Cfg.CPUThreshold = EnvFloat("CPU_THRESHOLD", 90)
	Cfg.RAMThreshold = EnvFloat("RAM_THRESHOLD", 85)
	Cfg.DiskThreshold = EnvFloat("DISK_THRESHOLD", 85)
	Cfg.LoadThreshold = EnvFloat("LOAD_THRESHOLD", 4.0)
	Cfg.SwapThresholdMB = EnvFloat("SWAP_THRESHOLD_MB", 512)

	Cfg.AlertCooldown = EnvInt("ALERT_COOLDOWN", 600)
	Cfg.ReportInterval = EnvInt("REPORT_INTERVAL", 3600)

	// Monitoring toggles — defaults: monitoring ON, security ON, reports OFF
	Cfg.MonitoringEnabled = EnvBool("MONITORING_ENABLED", true)
	Cfg.SecurityAlertsEnabled = EnvBool("SECURITY_ALERTS_ENABLED", true)
	Cfg.ScheduledReportsEnabled = EnvBool("SCHEDULED_REPORTS_ENABLED", false)

	Cfg.RcloneGDriveMount = EnvStr("RCLONE_GDRIVE_MOUNT", "/data/coolify/storage/gdrive")
	Cfg.RcloneS3GatewayURL = EnvStr("RCLONE_S3_GATEWAY_URL", "http://localhost:9900")

	wl := EnvStr("SSH_ALERT_WHITELIST", "10.0.2.,127.0.0.1")
	for _, s := range strings.Split(wl, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			Cfg.SSHAlertWhitelist = append(Cfg.SSHAlertWhitelist, s)
		}
	}
	Cfg.SSHAlertDedup = EnvInt("SSH_ALERT_DEDUP", 120)

	return nil
}

func EnvStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func EnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func EnvFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

func EnvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	}
	return def
}
