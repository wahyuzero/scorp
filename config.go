package main

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

	RcloneGDriveMount  string
	RcloneS3GatewayURL string

	SSHAlertWhitelist []string
	SSHAlertDedup     int // seconds
}

var cfg Config

func loadConfig() error {
	// Load .env file (ignore error if not found — env vars may be set directly)
	godotenv.Load()

	cfg.TelegramBotToken = envStr("TELEGRAM_BOT_TOKEN", "")
	cfg.TelegramChatID = envStr("TELEGRAM_CHAT_ID", "")
	cfg.TelegramWebhookURL = envStr("TELEGRAM_WEBHOOK_URL", "")
	if cfg.TelegramBotToken == "" || cfg.TelegramChatID == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID are required")
	}

	cfg.CoolifyAPIURL = envStr("COOLIFY_API_URL", "http://localhost:8000")
	cfg.CoolifyAPIToken = envStr("COOLIFY_API_TOKEN", "")

	cfg.CPUThreshold = envFloat("CPU_THRESHOLD", 90)
	cfg.RAMThreshold = envFloat("RAM_THRESHOLD", 85)
	cfg.DiskThreshold = envFloat("DISK_THRESHOLD", 85)
	cfg.LoadThreshold = envFloat("LOAD_THRESHOLD", 4.0)
	cfg.SwapThresholdMB = envFloat("SWAP_THRESHOLD_MB", 512)

	cfg.AlertCooldown = envInt("ALERT_COOLDOWN", 600)
	cfg.ReportInterval = envInt("REPORT_INTERVAL", 3600)

	cfg.RcloneGDriveMount = envStr("RCLONE_GDRIVE_MOUNT", "/data/coolify/storage/gdrive")
	cfg.RcloneS3GatewayURL = envStr("RCLONE_S3_GATEWAY_URL", "http://localhost:9900")

	wl := envStr("SSH_ALERT_WHITELIST", "10.0.2.,127.0.0.1")
	for _, s := range strings.Split(wl, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			cfg.SSHAlertWhitelist = append(cfg.SSHAlertWhitelist, s)
		}
	}
	cfg.SSHAlertDedup = envInt("SSH_ALERT_DEDUP", 120)

	return nil
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}
