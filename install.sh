#!/bin/bash
# install.sh — Build and install VPS Monitor (Go) as a systemd service

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY="$SCRIPT_DIR/vps-monitor-go"
SERVICE_NAME="vps-monitor-go"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

echo "=== Building VPS Monitor (Go) ==="
export PATH=$PATH:/usr/local/go/bin
cd "$SCRIPT_DIR"
go build -o "$BINARY" .
echo "✅ Built: $BINARY ($(du -h "$BINARY" | cut -f1))"

echo ""
echo "=== Installing systemd service ==="

sudo tee "$SERVICE_FILE" > /dev/null << EOF
[Unit]
Description=VPS Monitor (Go)
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=simple
ExecStart=$BINARY
WorkingDirectory=$SCRIPT_DIR
Restart=always
RestartSec=10
User=root
Environment=PATH=/usr/local/bin:/usr/bin:/bin:/usr/local/go/bin

# Reduce OOM killer probability
OOMScoreAdjust=-500

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable "$SERVICE_NAME"

echo "✅ Service installed: $SERVICE_NAME"
echo ""
echo "=== Usage ==="
echo "  Start:   sudo systemctl start $SERVICE_NAME"
echo "  Stop:    sudo systemctl stop $SERVICE_NAME" 
echo "  Status:  sudo systemctl status $SERVICE_NAME"
echo "  Logs:    sudo journalctl -u $SERVICE_NAME -f"
echo ""
echo "⚠️  Make sure to stop the Python version first:"
echo "  sudo systemctl stop vps-monitor"
echo ""
