#!/bin/bash
set -euo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"
BIN="$DIR/scorp"
INSTALL_BIN="/usr/local/bin/scorp"
SERVICE="scorp"
SVC_FILE="/etc/systemd/system/${SERVICE}.service"
ENV_FILE="$DIR/.env"
MODELS_FILE="$HOME/.scorp/models.json"

G='\033[0;32m'; Y='\033[0;33m'; R='\033[0;31m'; C='\033[0;36m'; B='\033[1m'; N='\033[0m'
ok()   { echo -e "${G}✓${N} $1"; }
warn() { echo -e "${Y}⚠${N}  $1"; }
die()  { echo -e "${R}✗${N} $1"; exit 1; }
ask()  { local val; read -rp "$(echo -e "${C}$1${N}")" val; echo "$val"; }
askd() { local val; read -rp "$(echo -e "${C}$1${N} [${2:-}]:")" val; echo "${val:-$2}"; }
confirm() { local val; read -rp "$(echo -e "${C}$1${N} [y/N]:")" val; [[ "$val" =~ ^[Yy] ]]; }
confirmY() { local val; read -rp "$(echo -e "${C}$1${N} [Y/n]:")" val; [[ ! "$val" =~ ^[Nn] ]]; }

echo -e "${B}=== scorp installer ===${N}"

TAGS="fts5"
for arg in "$@"; do
    case "$arg" in
        --minimal) TAGS="$TAGS nobrowser" ;;
    esac
done

command -v git >/dev/null || die "git not installed"
[ -f "$DIR/go.mod" ] || die "run from project root (go.mod not found)"
export PATH="$PATH:/usr/local/go/bin"
command -v go >/dev/null 2>&1 || die "Go not installed"

EXISTING_TOKEN=""
if [ -f "$ENV_FILE" ]; then
    EXISTING_TOKEN=$(grep -oP '^TELEGRAM_BOT_TOKEN=\K.*' "$ENV_FILE" 2>/dev/null || echo "")
fi

# ═══ Step 1: LLM Provider ═══
echo ""
echo -e "${B}[1/3] LLM Provider${N}"

SKIP_LLM=false
if [ -f "$MODELS_FILE" ]; then
    EXISTING_MODEL=$(grep -oP '"default_model":\s*"\K[^"]+' "$MODELS_FILE" 2>/dev/null || echo "")
    if [ -n "$EXISTING_MODEL" ]; then
        ok "models.json exists (model: $EXISTING_MODEL)"
        confirm "Reconfigure LLM?" || SKIP_LLM=true
    fi
fi

PROVIDER=""; MODEL_ID=""; KEY_ENV=""; BASE_URL=""; API_KEY=""
if [ "$SKIP_LLM" = false ]; then
    echo "  Select provider:"
    echo "    1) Z.AI / GLM        (free tier, recommended)"
    echo "    2) DeepSeek"
    echo "    3) OpenAI"
    echo "    4) Groq              (free tier, fast)"
    echo "    5) OpenRouter        (100+ models)"
    echo "    6) Google Gemini"
    echo "    7) Anthropic Claude"
    echo "    8) Ollama            (local, no API key)"
    echo "    9) Custom"
    PCHOICE=$(askd "  Provider" "1")

    case "$PCHOICE" in
        1) PROVIDER="zai-coding";   MODEL_ID="glm-4.7";                  KEY_ENV="ZAI_CODING_API_KEY"; BASE_URL="https://api.z.ai/api/coding/paas/v4" ;;
        2) PROVIDER="deepseek";     MODEL_ID="deepseek-chat";             KEY_ENV="DEEPSEEK_API_KEY";  BASE_URL="https://api.deepseek.com/v1" ;;
        3) PROVIDER="openai";       MODEL_ID="gpt-4o-mini";              KEY_ENV="OPENAI_API_KEY";    BASE_URL="https://api.openai.com/v1" ;;
        4) PROVIDER="groq";         MODEL_ID="llama-3.3-70b-versatile";  KEY_ENV="GROQ_API_KEY";      BASE_URL="https://api.groq.com/openai/v1" ;;
        5) PROVIDER="openrouter";   MODEL_ID="deepseek/deepseek-chat";   KEY_ENV="OPENROUTER_API_KEY"; BASE_URL="https://openrouter.ai/api/v1" ;;
        6) PROVIDER="gemini";       MODEL_ID="gemini-2.0-flash";         KEY_ENV="GOOGLE_API_KEY";    BASE_URL="https://generativelanguage.googleapis.com/v1beta" ;;
        7) PROVIDER="anthropic";    MODEL_ID="claude-sonnet-4-20250514"; KEY_ENV="ANTHROPIC_API_KEY"; BASE_URL="https://api.anthropic.com" ;;
        8) PROVIDER="ollama";       MODEL_ID="llama3.2";                  KEY_ENV="";                   BASE_URL="http://localhost:11434/v1" ;;
        9)
            PROVIDER="custom"
            MODEL_ID=$(ask "  Model ID:")
            KEY_ENV=$(ask "  Env var name for key (e.g. MY_API_KEY):")
            BASE_URL=$(ask "  Base URL:")
            ;;
        *) die "Invalid choice" ;;
    esac

    if [ -n "$KEY_ENV" ] && [ "$PROVIDER" != "ollama" ]; then
        EXISTING_KEY=""
        if [ -f "$ENV_FILE" ]; then
            EXISTING_KEY=$(grep -oP "^${KEY_ENV}=\K.*" "$ENV_FILE" 2>/dev/null || echo "")
        fi
        if [ -n "$EXISTING_KEY" ]; then
            ok "API key exists: ${EXISTING_KEY:0:8}..."
            confirm "Change key?" || API_KEY="$EXISTING_KEY"
        fi
        if [ -z "$API_KEY" ]; then
            API_KEY=$(ask "  API Key ($KEY_ENV):")
            [ -z "$API_KEY" ] && warn "No key - set it later in .env"
        fi
    fi
fi

# ═══ Step 2: Telegram (optional) ═══
echo ""
echo -e "${B}[2/3] Telegram (optional)${N}"

TG_TOKEN=""; TG_CHATID=""; USE_TELEGRAM=false

if [ -n "$EXISTING_TOKEN" ]; then
    ok "Telegram already configured: ${EXISTING_TOKEN:0:12}..."
    if confirmY "Keep Telegram?"; then
        USE_TELEGRAM=true
        TG_TOKEN="$EXISTING_TOKEN"
        EXISTING_CHATID=$(grep -oP '^TELEGRAM_CHAT_ID=\K.*' "$ENV_FILE" 2>/dev/null || echo "")
        TG_CHATID="$EXISTING_CHATID"
    fi
else
    echo "  Skip = CLI mode (default)."
    if confirm "Set up Telegram bot now?"; then
        USE_TELEGRAM=true
    fi
fi

if [ "$USE_TELEGRAM" = true ] && [ -z "$TG_TOKEN" ]; then
    echo -e "  Get token from ${C}@BotFather${N} on Telegram"
    TG_TOKEN=$(ask "  Bot Token:")
    if [ -z "$TG_TOKEN" ]; then
        warn "No token - falling back to CLI mode"
        USE_TELEGRAM=false
    fi
fi

if [ "$USE_TELEGRAM" = true ] && [ -z "$TG_CHATID" ]; then
    echo -e "  Chat ID - get from ${C}@userinfobot${N}, or Enter to auto-detect"
    TG_CHATID=$(ask "  Your Chat ID:")
    if [ -z "$TG_CHATID" ]; then
        echo -e "  ${C}Auto-detecting (send any message to your bot first)...${N}"
        TG_CHATID=$(curl -s "https://api.telegram.org/bot${TG_TOKEN}/getUpdates?limit=1&offset=-1" \
            | grep -oP '"chat":\{"id":\K[0-9]+' | head -1)
        if [ -n "$TG_CHATID" ]; then
            ok "Detected Chat ID: $TG_CHATID"
        else
            warn "Could not auto-detect."
            TG_CHATID=$(ask "  Enter Chat ID manually:")
        fi
    fi
fi

# ═══ Step 3: Options ═══
echo ""
echo -e "${B}[3/3] Options${N}"

MON_ON="false"; SEC_ON="false"
if [ "$USE_TELEGRAM" = true ]; then
    confirm "Enable resource monitoring (CPU/RAM/Disk)?" && MON_ON="true"
    confirm "Enable security alerts (SSH login)?" && SEC_ON="true"
fi

COOLIFY_URL=$(askd "  Coolify URL" "")
COOLIFY_TOKEN=""
if [ -n "$COOLIFY_URL" ]; then
    COOLIFY_TOKEN=$(ask "  Coolify API Token:")
fi

# ═══ Summary ═══
echo ""
echo -e "${B}Summary:${N}"
if [ "$USE_TELEGRAM" = true ]; then
    echo "  Mode: Telegram + CLI"
else
    echo "  Mode: CLI only"
fi
echo "  Provider: ${PROVIDER:-existing} / ${MODEL_ID:-existing}"
[ "$USE_TELEGRAM" = true ] && echo "  Monitoring: $MON_ON | Security: $SEC_ON"
[ -n "$COOLIFY_URL" ] && echo "  Coolify: $COOLIFY_URL"
echo ""

# ═══ Build & Install ═══
echo -e "${B}Building...${N}"

echo "# scorp - generated $(date)" > "$ENV_FILE"
if [ "$USE_TELEGRAM" = true ]; then
    echo "TELEGRAM_BOT_TOKEN=$TG_TOKEN" >> "$ENV_FILE"
    echo "TELEGRAM_CHAT_ID=$TG_CHATID" >> "$ENV_FILE"
fi
if [ "$SKIP_LLM" = false ] && [ -n "$KEY_ENV" ] && [ -n "$API_KEY" ]; then
    echo "$KEY_ENV=$API_KEY" >> "$ENV_FILE"
fi
echo "MONITORING_ENABLED=$MON_ON" >> "$ENV_FILE"
echo "SECURITY_ALERTS_ENABLED=$SEC_ON" >> "$ENV_FILE"
echo "SCHEDULED_REPORTS_ENABLED=false" >> "$ENV_FILE"
if [ -n "$COOLIFY_URL" ]; then
    echo "COOLIFY_API_URL=$COOLIFY_URL" >> "$ENV_FILE"
    [ -n "$COOLIFY_TOKEN" ] && echo "COOLIFY_API_TOKEN=$COOLIFY_TOKEN" >> "$ENV_FILE"
fi
ok ".env"

mkdir -p "$HOME/.scorp"
if [ "$SKIP_LLM" = false ]; then
    cat > "$MODELS_FILE" << MEOF
{
  "default_model": "$MODEL_ID",
  "agent_model": "$MODEL_ID",
  "delegation_model": "",
  "premium_model": "",
  "models": {
    "$MODEL_ID": {
      "provider": "$PROVIDER",
      "model": "$MODEL_ID",
      "key_env": "$KEY_ENV",
      "base_url": "$BASE_URL",
      "max_tokens": 4096,
      "api": "openai"
    }
  },
  "routing_rules": { "agent": "$MODEL_ID", "chat": "$MODEL_ID" },
  "fallback_on_error": ["rate_limit","timeout","server_error","auth_error","network_error"]
}
MEOF
    ok "models.json ($MODEL_ID)"
fi

echo "Building ($TAGS)..."
CGO_ENABLED=1 go build -tags "$TAGS" -ldflags="-s -w" -trimpath -o "$BIN" . || die "Build failed"
ok "Built: $(du -h "$BIN" | cut -f1)"

sudo cp "$BIN" "$INSTALL_BIN"
ok "Installed: $INSTALL_BIN"

if [ "$USE_TELEGRAM" = true ]; then
    sudo tee "$SVC_FILE" > /dev/null << EOF
[Unit]
Description=Scorp - Personal AI Agent
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=simple
ExecStart=$INSTALL_BIN
WorkingDirectory=$DIR
Restart=always
RestartSec=10
User=root
Environment=PATH=/usr/local/bin:/usr/bin:/bin:/usr/local/go/bin
Environment=HOME=$HOME
EnvironmentFile=-$DIR/.env
OOMScoreAdjust=-500

[Install]
WantedBy=multi-user.target
EOF
    sudo systemctl daemon-reload
    sudo systemctl enable "$SERVICE" 2>/dev/null
    sudo systemctl restart "$SERVICE"
    sleep 2
    if systemctl is-active --quiet "$SERVICE"; then
        ok "Service running! PID $(systemctl show -p MainPID --value "$SERVICE")"
    else
        die "Service failed. Check: journalctl -u $SERVICE -e"
    fi
else
    sudo systemctl stop "$SERVICE" 2>/dev/null || true
    sudo systemctl disable "$SERVICE" 2>/dev/null || true
    sudo rm -f "$SVC_FILE" 2>/dev/null || true
    ok "CLI mode ready"
fi

echo ""
echo -e "${B}=== Done ===${N}"
echo ""
if [ "$USE_TELEGRAM" = true ]; then
    echo "  Bot running 24/7 via systemd. CLI also available: run 'scorp'"
    echo ""
    echo "  Logs:    sudo journalctl -u $SERVICE -f"
    echo "  Stop:    sudo systemctl stop $SERVICE"
    echo "  Restart: sudo systemctl restart $SERVICE"
else
    echo "  Run 'scorp' to start chatting."
    echo "  Config:  ~/.scorp/models.json"
    echo "  Env:     $ENV_FILE"
fi
echo ""
echo "  Update:  scorp update  (or /update in Telegram)"
echo "  Version: scorp version"
echo ""
