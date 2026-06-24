#!/bin/bash
set -euo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"
BIN="$DIR/scorp"
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

# ─── Detect environment ───
IS_TERMUX=false
SUDO=""
PREFIX_BIN=""
HAS_SYSTEMD=false

if [ -n "${PREFIX:-}" ] && [[ "$PREFIX" == */com.termux* ]]; then
    IS_TERMUX=true
    PREFIX_BIN="$PREFIX/bin"
elif [ -d "/data/data/com.termux" ]; then
    IS_TERMUX=true
    PREFIX_BIN="${PREFIX:-/data/data/com.termux/files/usr}/bin"
fi

if [ "$IS_TERMUX" = false ]; then
    HAS_SYSTEMD=$(systemctl is-system-running >/dev/null 2>&1 && echo "yes" || echo "no")
    if command -v sudo >/dev/null 2>&1; then
        SUDO="sudo"
    fi
fi

# Install path
if [ "$IS_TERMUX" = true ]; then
    INSTALL_BIN="$PREFIX_BIN/scorp"
else
    INSTALL_BIN="/usr/local/bin/scorp"
fi

SERVICE="scorp"
SVC_FILE="/etc/systemd/system/${SERVICE}.service"

echo -e "${B}=== scorp installer ===${N}"
[ "$IS_TERMUX" = true ] && echo -e "${Y}Termux detected${N} (${PREFIX_BIN})"
[ "$HAS_SYSTEMD" = "yes" ] && echo "systemd detected"

TAGS="fts5"
for arg in "$@"; do
    case "$arg" in
        --minimal) TAGS="$TAGS nobrowser" ;;
    esac
done

command -v git >/dev/null || die "git not installed"
[ -f "$DIR/go.mod" ] || die "run from project root (go.mod not found)"

# Add Go to PATH (standard Linux install location)
export PATH="$PATH:/usr/local/go/bin"
command -v go >/dev/null 2>&1 || die "Go not installed (run: pkg install golang on Termux, or install from go.dev)"

# Use grep -E instead of -P for portability (Termux grep may lack PCRE)
grep_word() {
    grep -oE "$1" "$2" 2>/dev/null | head -1 || echo ""
}

EXISTING_TOKEN=""
if [ -f "$ENV_FILE" ]; then
    EXISTING_TOKEN=$(grep_word '^TELEGRAM_BOT_TOKEN=.*' "$ENV_FILE")
    EXISTING_TOKEN="${EXISTING_TOKEN#TELEGRAM_BOT_TOKEN=}"
fi

# ═══ Step 1: LLM Provider ═══
echo ""
echo -e "${B}[1/3] LLM Provider${N}"

SKIP_LLM=false
if [ -f "$MODELS_FILE" ]; then
    EXISTING_MODEL=$(grep_word '"default_model":\s*"[^"]+"' "$MODELS_FILE")
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
            EXISTING_KEY=$(grep_word "^${KEY_ENV}=.*" "$ENV_FILE")
            EXISTING_KEY="${EXISTING_KEY#${KEY_ENV}=}"
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
        EXISTING_CHATID=$(grep_word '^TELEGRAM_CHAT_ID=.*' "$ENV_FILE")
        EXISTING_CHATID="${EXISTING_CHATID#TELEGRAM_CHAT_ID=}"
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
            | grep_word '"chat":{"id":[0-9]+' | grep -oE '[0-9]+')
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
[ "$IS_TERMUX" = true ] && echo -e "  ${Y}Platform: Termux${N} (build from source, no systemd)"
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

# ─── Install binary ───
if [ "$IS_TERMUX" = true ]; then
    # Termux: no sudo needed, install to $PREFIX/bin
    cp "$BIN" "$INSTALL_BIN"
    ok "Installed: $INSTALL_BIN"
elif [ "$HAS_SYSTEMD" = "yes" ]; then
    $SUDO cp "$BIN" "$INSTALL_BIN"
    ok "Installed: $INSTALL_BIN"
else
    # Non-Termux, no systemd — try local bin
    LOCAL_BIN="$HOME/.local/bin"
    mkdir -p "$LOCAL_BIN"
    cp "$BIN" "$LOCAL_BIN/scorp"
    INSTALL_BIN="$LOCAL_BIN/scorp"
    ok "Installed: $INSTALL_BIN"
    warn "Add to PATH: export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

# ─── Service setup (systemd only) ───
if [ "$HAS_SYSTEMD" = "yes" ]; then
    if [ "$USE_TELEGRAM" = true ]; then
        $SUDO tee "$SVC_FILE" > /dev/null << EOF
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
        $SUDO systemctl daemon-reload
        $SUDO systemctl enable "$SERVICE" 2>/dev/null
        $SUDO systemctl restart "$SERVICE"
        sleep 2
        if systemctl is-active --quiet "$SERVICE"; then
            ok "Service running! PID $(systemctl show -p MainPID --value "$SERVICE")"
        else
            die "Service failed. Check: journalctl -u $SERVICE -e"
        fi
    else
        $SUDO systemctl stop "$SERVICE" 2>/dev/null || true
        $SUDO systemctl disable "$SERVICE" 2>/dev/null || true
        $SUDO rm -f "$SVC_FILE" 2>/dev/null || true
        ok "CLI mode ready"
    fi
fi

echo ""
echo -e "${B}=== Done ===${N}"
echo ""
if [ "$IS_TERMUX" = true ]; then
    if [ "$USE_TELEGRAM" = true ]; then
        echo "  Bot running in foreground. For background:"
        echo "    nohup scorp &"
        echo "    (or use tmux: tmux new -s scorp, run scorp, Ctrl-B D)"
    else
        echo "  Run 'scorp' to start chatting."
        echo "  Config:  ~/.scorp/models.json"
        echo "  Env:     $ENV_FILE"
    fi
elif [ "$HAS_SYSTEMD" = "yes" ] && [ "$USE_TELEGRAM" = true ]; then
    echo "  Bot running 24/7 via systemd. CLI also available: run 'scorp'"
    echo ""
    echo "  Logs:    journalctl -u $SERVICE -f"
    echo "  Stop:    systemctl stop $SERVICE"
    echo "  Restart: systemctl restart $SERVICE"
else
    echo "  Run 'scorp' to start chatting."
    echo "  Config:  ~/.scorp/models.json"
    echo "  Env:     $ENV_FILE"
fi
echo ""
echo "  Update:  scorp update"
echo "  Version: scorp version"
echo ""
