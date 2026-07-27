#!/usr/bin/env bash
# ══════════════════════════════════════════════════════════════
#  install-scorp.sh — Scorp Agent Installer
#  Repo   : https://github.com/wahyuzero/scorp
#  Fokus  : Agent-first, build from source, systemd service
# ══════════════════════════════════════════════════════════════
set -euo pipefail

# ── Colors ──────────────────────────────────────────────────
G='\033[0;32m'; Y='\033[0;33m'; R='\033[0;31m'
C='\033[0;36m'; B='\033[1m';    N='\033[0m'
ok()   { echo -e "${G}✓${N} $1"; }
warn() { echo -e "${Y}⚠${N}  $1"; }
die()  { echo -e "${R}✗${N} $1"; exit 1; }
ask()  { local v; read -rp "$(echo -e "${C}  $1: ${N}")" v; echo "$v"; }

# ── Constants ───────────────────────────────────────────────
REPO_URL="https://github.com/wahyuzero/scorp"
SRC_DIR="/root/scorp_src"
WORK_DIR="/root/scorp"
INSTALL_BIN="/usr/local/bin/scorp"
ENV_FILE="$WORK_DIR/.env"
MODELS_FILE="$HOME/.scorp/models.json"
SVC_FILE="/etc/systemd/system/scorp.service"
GO_BIN="/usr/local/go/bin/go"

echo ""
echo -e "${B}╔══════════════════════════════════╗${N}"
echo -e "${B}║     Scorp Agent Installer        ║${N}"
echo -e "${B}╚══════════════════════════════════╝${N}"
echo ""

# ── Step 0: Prasyarat ────────────────────────────────────────
echo -e "${B}[0/4] Checking prerequisites...${N}"

command -v git  >/dev/null 2>&1 || die "git tidak terinstall. Jalankan: apt install git"
command -v curl >/dev/null 2>&1 || die "curl tidak terinstall. Jalankan: apt install curl"

# Cari Go yang benar (prioritas /usr/local/go/bin/go)
if [ -x "$GO_BIN" ]; then
    GO_VER=$("$GO_BIN" version | awk '{print $3}' | sed 's/go//')
    ok "Go ditemukan: $GO_VER ($GO_BIN)"
elif command -v go >/dev/null 2>&1; then
    GO_BIN=$(command -v go)
    GO_VER=$(go version | awk '{print $3}' | sed 's/go//')
    GO_MINOR=$(echo "$GO_VER" | cut -d. -f2)
    if [ "$GO_MINOR" -lt 21 ]; then
        die "Go $GO_VER terlalu lama (minimal 1.21). Install: https://go.dev/dl/"
    fi
    ok "Go ditemukan: $GO_VER"
else
    die "Go tidak terinstall. Install: https://go.dev/dl/"
fi

# CGO untuk fts5
if command -v gcc >/dev/null 2>&1; then
    ok "gcc ditemukan (fts5 enabled)"
    BUILD_TAGS="fts5"
    CGO_FLAG="CGO_ENABLED=1"
else
    warn "gcc tidak ada — build tanpa fts5"
    BUILD_TAGS=""
    CGO_FLAG="CGO_ENABLED=0"
fi

# ── Step 1: Credentials ──────────────────────────────────────
echo ""
echo -e "${B}[1/4] Credentials...${N}"

OLD_TG_TOKEN=""; OLD_TG_CHATID=""; OLD_DEEPSEEK=""; OLD_AGENTON=""

if [ -f "$ENV_FILE" ]; then
    OLD_TG_TOKEN=$(grep -oP '(?<=^TELEGRAM_BOT_TOKEN=).+' "$ENV_FILE" 2>/dev/null || true)
    OLD_TG_CHATID=$(grep -oP '(?<=^TELEGRAM_CHAT_ID=).+' "$ENV_FILE" 2>/dev/null || true)
    OLD_DEEPSEEK=$(grep -oP '(?<=^DEEPSEEK_API_KEY=).+' "$ENV_FILE" 2>/dev/null || true)
    OLD_AGENTON=$(grep -oP '(?<=^AGENTON_API_KEY=).+' "$ENV_FILE" 2>/dev/null || true)
    ok "Credentials lama ditemukan di $ENV_FILE"
fi

# Telegram Token
if [ -n "$OLD_TG_TOKEN" ]; then
    echo -e "  Token lama: ${C}${OLD_TG_TOKEN:0:15}...${N}"
    read -rp "  $(echo -e "${C}Gunakan token lama? [Y/n]: ${N}")" USE_OLD
    [[ "$USE_OLD" =~ ^[Nn] ]] && TG_TOKEN=$(ask "Bot Token baru") || TG_TOKEN="$OLD_TG_TOKEN"
else
    echo -e "  Dapatkan token dari ${C}@BotFather${N} di Telegram"
    TG_TOKEN=$(ask "Bot Token (kosong = CLI-only mode)")
fi

# Chat ID
TG_CHATID=""
if [ -n "$TG_TOKEN" ]; then
    if [ -n "$OLD_TG_CHATID" ]; then
        echo -e "  Chat ID lama: ${C}$OLD_TG_CHATID${N}"
        read -rp "  $(echo -e "${C}Gunakan Chat ID lama? [Y/n]: ${N}")" USE_OLD_ID
        [[ "$USE_OLD_ID" =~ ^[Nn] ]] && TG_CHATID=$(ask "Chat ID baru") || TG_CHATID="$OLD_TG_CHATID"
    else
        TG_CHATID=$(ask "Chat ID (Enter = auto-detect)")
        if [ -z "$TG_CHATID" ]; then
            echo -e "  ${C}Auto-detecting...${N}"
            TG_CHATID=$(curl -s "https://api.telegram.org/bot${TG_TOKEN}/getUpdates?limit=1&offset=-1" \
                | grep -oP '"id":\K[0-9]+' | head -1 || true)
            [ -n "$TG_CHATID" ] && ok "Chat ID: $TG_CHATID" || TG_CHATID=$(ask "Chat ID (manual)")
        fi
    fi
fi

# DeepSeek Key
if [ -n "$OLD_DEEPSEEK" ]; then
    echo -e "  DeepSeek key lama: ${C}${OLD_DEEPSEEK:0:12}...${N}"
    read -rp "  $(echo -e "${C}Gunakan key lama? [Y/n]: ${N}")" USE_DS
    [[ "$USE_DS" =~ ^[Nn] ]] && DEEPSEEK_KEY=$(ask "DeepSeek API Key") || DEEPSEEK_KEY="$OLD_DEEPSEEK"
else
    DEEPSEEK_KEY=$(ask "DeepSeek API Key (dari platform.deepseek.com)")
fi

# AgentON Key
AGENTON_KEY="${OLD_AGENTON:-}"
if [ -z "$AGENTON_KEY" ]; then
    AGENTON_KEY=$(ask "AgentON API Key (opsional, Enter skip)")
fi
[ -n "$AGENTON_KEY" ] && ok "AgentON key tersimpan" || true

# ── Step 2: Clone / Update ───────────────────────────────────
echo ""
echo -e "${B}[2/4] Source code...${N}"

if [ -d "$SRC_DIR/.git" ]; then
    echo "  Update repo..."
    git -C "$SRC_DIR" fetch origin
    git -C "$SRC_DIR" reset --hard origin/main 2>&1 | tail -2
    ok "Source diupdate"
else
    echo "  Cloning $REPO_URL..."
    git clone "$REPO_URL" "$SRC_DIR" 2>&1 | tail -3
    ok "Clone berhasil"
fi

# ── FIX BUG UPSTREAM: go.mod versi 3-part tidak valid ──
GOMOD_VER=$(grep '^go ' "$SRC_DIR/go.mod" | awk '{print $2}')
GOMOD_PARTS=$(echo "$GOMOD_VER" | awk -F. '{print NF}')
if [ "$GOMOD_PARTS" -gt 2 ]; then
    warn "go.mod: versi '$GOMOD_VER' tidak valid (bug upstream). Auto-fix ke '1.24'..."
    sed -i "s/^go $GOMOD_VER$/go 1.24/" "$SRC_DIR/go.mod"
    ok "go.mod diperbaiki"
fi

# ── Step 3: Build ────────────────────────────────────────────
echo ""
echo -e "${B}[3/4] Building...${N}"

cd "$SRC_DIR"

echo "  go mod tidy..."
"$GO_BIN" mod tidy 2>&1 | grep -v '^$' | head -15 || true

echo "  Compiling..."
if [ -n "$BUILD_TAGS" ]; then
    CGO_ENABLED=1 "$GO_BIN" build -tags "$BUILD_TAGS" -ldflags="-s -w" -trimpath -o scorp .
else
    CGO_ENABLED=0 "$GO_BIN" build -ldflags="-s -w" -trimpath -o scorp .
fi

BIN_SIZE=$(du -h "$SRC_DIR/scorp" | cut -f1)
ok "Build sukses: $BIN_SIZE"

# Stop service lama jika jalan
systemctl stop scorp 2>/dev/null || true

cp "$SRC_DIR/scorp" "$INSTALL_BIN"
chmod +x "$INSTALL_BIN"
ok "Binary: $INSTALL_BIN"

# ── Step 4: Config & Service ─────────────────────────────────
echo ""
echo -e "${B}[4/4] Konfigurasi...${N}"

mkdir -p "$WORK_DIR" "$HOME/.scorp"

# Tulis .env
{
    echo "# Scorp Agent — generated $(date '+%Y-%m-%d')"
    echo ""
    if [ -n "$TG_TOKEN" ]; then
        echo "# Telegram"
        echo "TELEGRAM_BOT_TOKEN=$TG_TOKEN"
        [ -n "$TG_CHATID" ] && echo "TELEGRAM_CHAT_ID=$TG_CHATID"
        echo ""
    fi
    [ -n "$DEEPSEEK_KEY" ] && echo "DEEPSEEK_API_KEY=$DEEPSEEK_KEY"
    [ -n "$AGENTON_KEY"  ] && echo "AGENTON_API_KEY=$AGENTON_KEY"
    echo ""
    echo "GITHUB_REPO=wahyuzero/scorp"
    echo "MONITORING_ENABLED=true"
    echo "SECURITY_ALERTS_ENABLED=true"
    echo "SCHEDULED_REPORTS_ENABLED=false"
} > "$ENV_FILE"
chmod 600 "$ENV_FILE"
ok ".env → $ENV_FILE"

# Tulis models.json jika belum ada
if [ ! -f "$MODELS_FILE" ] && [ -n "$DEEPSEEK_KEY" ]; then
    cat > "$MODELS_FILE" << MEOF
{
  "default_model": "deepseek-chat",
  "agent_model": "deepseek-chat",
  "delegation_model": "",
  "premium_model": "",
  "models": {
    "deepseek-chat": {
      "provider": "deepseek",
      "model": "deepseek-chat",
      "key_env": "DEEPSEEK_API_KEY",
      "base_url": "https://api.deepseek.com/v1",
      "max_tokens": 4096,
      "api": "openai"
    }
  },
  "routing_rules": { "agent": "deepseek-chat", "chat": "deepseek-chat" },
  "fallback_on_error": ["rate_limit","timeout","server_error","auth_error","network_error"]
}
MEOF
    ok "models.json → $MODELS_FILE"
fi

# Systemd service
if [ -n "$TG_TOKEN" ] && systemctl is-system-running >/dev/null 2>&1; then
    cat > "$SVC_FILE" << EOF
[Unit]
Description=Scorp - Personal AI Agent
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=$INSTALL_BIN
WorkingDirectory=$WORK_DIR
Restart=always
RestartSec=10
User=root
Environment=PATH=/usr/local/bin:/usr/bin:/bin
Environment=HOME=$HOME
EnvironmentFile=$ENV_FILE
OOMScoreAdjust=-500

[Install]
WantedBy=multi-user.target
EOF
    systemctl daemon-reload
    systemctl enable scorp 2>/dev/null || true
    systemctl restart scorp
    sleep 2

    if systemctl is-active --quiet scorp; then
        ok "Service aktif! PID: $(systemctl show -p MainPID --value scorp)"
    else
        warn "Service gagal start → cek: journalctl -u scorp -e"
    fi
else
    rm -f "$SVC_FILE" 2>/dev/null || true
    systemctl daemon-reload 2>/dev/null || true
    ok "Mode CLI (tanpa systemd)"
fi

# ── Done ─────────────────────────────────────────────────────
echo ""
echo -e "${B}╔══════════════════════════════════╗${N}"
echo -e "${B}║       Instalasi Selesai! ✓       ║${N}"
echo -e "${B}╚══════════════════════════════════╝${N}"
echo ""
echo -e "  Binary  : ${C}$INSTALL_BIN${N}"
echo -e "  Config  : ${C}$ENV_FILE${N}"
echo -e "  Models  : ${C}$MODELS_FILE${N}"
echo -e "  Source  : ${C}$SRC_DIR${N}"
echo ""

if systemctl is-active --quiet scorp 2>/dev/null; then
    echo -e "  ${G}● Bot berjalan 24/7 via systemd${N}"
    echo ""
    echo "  Logs    : journalctl -u scorp -f"
    echo "  Stop    : systemctl stop scorp"
    echo "  Restart : systemctl restart scorp"
else
    echo "  Jalankan : scorp"
fi

echo ""
echo "  Update  : scorp update"
echo "  Version : scorp version"
echo ""
