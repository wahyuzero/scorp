#!/usr/bin/env bash
# ============================================================
# VPS Backup Script — vps-monitor-go
# Daily:   DB dumps + docker volumes + coolify state + configs
# Weekly:  + projects, custom docker images
# Monthly: + full home, all docker images
# Notifikasi via Telegram (bot token from .env)
# ============================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"

# ── Load .env ────────────────────────────────────────────────
if [[ -f "$ENV_FILE" ]]; then
    set -a
    # shellcheck disable=SC1090
    source "$ENV_FILE"
    set +a
fi

# Variables from .env: TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID
TG_TOKEN="${TELEGRAM_BOT_TOKEN:-}"
TG_CHAT="${TELEGRAM_CHAT_ID:-}"

# ── Config ───────────────────────────────────────────────────
BACKUP_MODE="${1:-daily}"        # daily | weekly | monthly | test
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DATE_SHORT=$(date +%Y%m%d)
TMPDIR="/tmp/vps-backup-${TIMESTAMP}"
RCLONE_BASE="gdrive:backups/vps"

# Skip if load too high
MAX_LOAD=2.0

# ── Get DB passwords from Docker containers ──────────────────
get_env() {
    local container="$1" var="$2"
    docker inspect "$container" --format '{{range .Config.Env}}{{println .}}{{end}}' 2>/dev/null \
        | grep "^${var}=" | head -1 | cut -d= -f2-
}

MYSQL_ROOT_SS="$(get_env sistem-surat-mysql MYSQL_ROOT_PASSWORD)"
MYSQL_ROOT_WIV="$(get_env wiv0y55b1f5dzbyzlagyas1k MYSQL_ROOT_PASSWORD)"
PG_USER_COOLIFY="$(get_env coolify-db POSTGRES_USER)"
PG_PASS_COOLIFY="$(get_env coolify-db POSTGRES_PASSWORD)"
PG_DB_COOLIFY="$(get_env coolify-db POSTGRES_DB)"
PG_USER_N8N="$(get_env postgresql-s4q9sws8se3p3g0c621qbq9b POSTGRES_USER)"
PG_PASS_N8N="$(get_env postgresql-s4q9sws8se3p3g0c621qbq9b POSTGRES_PASSWORD)"
PG_DB_N8N="$(get_env postgresql-s4q9sws8se3p3g0c621qbq9b POSTGRES_DB)"
PG_USER_TUNNEL="$(get_env tunnel-postgres POSTGRES_USER)"
PG_PASS_TUNNEL="$(get_env tunnel-postgres POSTGRES_PASSWORD)"
PG_DB_TUNNEL="$(get_env tunnel-postgres POSTGRES_DB)"

# ── Helpers ──────────────────────────────────────────────────
human_size() {
    local bytes=$1
    if (( bytes >= 1073741824 )); then
        echo "$(echo "scale=1; $bytes/1073741824" | bc)GB"
    elif (( bytes >= 1048576 )); then
        echo "$(echo "scale=1; $bytes/1048576" | bc)MB"
    else
        echo "$(( bytes / 1024 ))KB"
    fi
}

tg_send() {
    local text="$1"
    if [[ -z "$TG_TOKEN" || -z "$TG_CHAT" ]]; then
        echo "[WARN] No TG_TOKEN/TG_CHAT, skipping notification"
        return
    fi
    curl -sf -X POST "https://api.telegram.org/bot${TG_TOKEN}/sendMessage" \
        -H "Content-Type: application/json" \
        -d "$(jq -n --arg chat_id "$TG_CHAT" --arg text "$text" \
            '{chat_id: $chat_id, text: $text, parse_mode: "HTML", disable_web_page_preview: true}')" \
        > /dev/null 2>&1 || echo "[WARN] TG send failed"
}

cleanup() {
    sudo rm -rf "$TMPDIR" 2>/dev/null || true
}
trap cleanup EXIT

# ── Pre-check ────────────────────────────────────────────────
START_EPOCH=$(date +%s)
START_TIME=$(date '+%Y-%m-%d %H:%M:%S UTC')

LOAD1=$(awk '{print $1}' /proc/loadavg)
IS_OVER=$(echo "$LOAD1 > $MAX_LOAD" | bc -l 2>/dev/null || echo 0)
if [[ "$IS_OVER" == "1" ]]; then
    tg_send "⏭️ <b>Backup Skipped</b>
🕐 ${START_TIME}
📊 Load average ${LOAD1} > ${MAX_LOAD}
Backup ditunda karena server sedang sibuk."
    echo "SKIP: load ${LOAD1} > ${MAX_LOAD}"
    exit 0
fi

# ── Start notification ───────────────────────────────────────
tg_send "🔄 <b>Backup ${BACKUP_MODE^^} Dimulai</b>
🕐 ${START_TIME}
📊 Load: ${LOAD1}
📦 Mode: ${BACKUP_MODE}
⏳ Proses backup sedang berjalan..."

mkdir -p "$TMPDIR"
ERRORS=()

# ── 1. Database Dumps ───────────────────────────────────────
echo "[1/9] Dumping databases..."
DB_DIR="$TMPDIR/databases"
mkdir -p "$DB_DIR"

# MySQL: sistem_surat
if docker exec sistem-surat-mysql mysqldump -u root -p"${MYSQL_ROOT_SS}" \
    --single-transaction --routines --triggers \
    sistem_surat 2>/dev/null | gzip > "$DB_DIR/sistem_surat.sql.gz"; then
    echo "  ✅ sistem_surat dumped"
else
    ERRORS+=("sistem_surat MySQL dump failed")
    echo "  ❌ sistem_surat FAILED"
fi

# MySQL: default (wiv container)
if docker exec wiv0y55b1f5dzbyzlagyas1k mysqldump -u root -p"${MYSQL_ROOT_WIV}" \
    --single-transaction --routines --triggers \
    default 2>/dev/null | gzip > "$DB_DIR/wiv_default.sql.gz"; then
    echo "  ✅ wiv_default dumped"
else
    ERRORS+=("wiv MySQL dump failed")
    echo "  ❌ wiv_default FAILED"
fi

# PostgreSQL: coolify
if docker exec coolify-db pg_dump -U "${PG_USER_COOLIFY}" "${PG_DB_COOLIFY}" 2>/dev/null \
    | gzip > "$DB_DIR/coolify.sql.gz"; then
    echo "  ✅ coolify dumped"
else
    ERRORS+=("coolify PG dump failed")
    echo "  ❌ coolify FAILED"
fi

# PostgreSQL: n8n
if docker exec postgresql-s4q9sws8se3p3g0c621qbq9b pg_dump -U "${PG_USER_N8N}" "${PG_DB_N8N}" 2>/dev/null \
    | gzip > "$DB_DIR/n8n.sql.gz"; then
    echo "  ✅ n8n dumped"
else
    ERRORS+=("n8n PG dump failed")
    echo "  ❌ n8n FAILED"
fi

# PostgreSQL: tunnel_service
if docker exec tunnel-postgres pg_dump -U "${PG_USER_TUNNEL}" "${PG_DB_TUNNEL}" 2>/dev/null \
    | gzip > "$DB_DIR/tunnel.sql.gz"; then
    echo "  ✅ tunnel dumped"
else
    ERRORS+=("tunnel PG dump failed")
    echo "  ❌ tunnel FAILED"
fi

# ── 2. Docker Volumes ────────────────────────────────────────
echo "[2/9] Dumping Docker volumes..."
DVOL_DIR="$TMPDIR/docker-volumes"
mkdir -p "$DVOL_DIR"

for vol in $(docker volume ls --format '{{.Name}}'); do
    # Skip cache/ephemeral volumes
    if echo "$vol" | grep -qiE '(cache|tmp)'; then
        echo "  ⏭️ $vol (cache — skipped)"
        continue
    fi
    if docker run --rm -v "$vol:/data" -v "$DVOL_DIR:/backup" alpine \
        tar czf "/backup/${vol}.tar.gz" -C /data . 2>/dev/null; then
        echo "  ✅ $vol"
    else
        ERRORS+=("docker volume $vol tar failed")
        echo "  ❌ $vol FAILED"
    fi
done

# ── 3. Coolify State ────────────────────────────────────────
echo "[3/9] Archiving Coolify state..."
if sudo tar czf "$TMPDIR/coolify-state.tar.gz" \
    -C /data/coolify \
    --exclude='storage/gdrive' \
    . 2>/dev/null; then
    echo "  ✅ coolify-state (excl gdrive mount)"
else
    ERRORS+=("coolify-state tar failed")
    echo "  ❌ coolify-state FAILED"
fi

# ── 4. Hermes Agent ─────────────────────────────────────────
echo "[4/9] Archiving Hermes agent (essential data only)..."
if ionice -c2 -n7 tar czf "$TMPDIR/hermes.tar.gz" -C /home/ubuntu \
    .hermes/config.yaml .hermes/.env .hermes/memories .hermes/skills \
    .hermes/sessions .hermes/logs .hermes/cron .hermes/plans .hermes/plugins \
    2>/dev/null; then
    echo "  ✅ hermes (essentials)"
else
    ERRORS+=("hermes tar failed")
    echo "  ❌ hermes FAILED"
fi

# ── 5. App Configs + System Configs ──────────────────────────
echo "[5/9] Archiving app + system configs..."
mkdir -p "$TMPDIR/app-configs"
# Rclone config (CRITICAL for restore)
cp -a /home/ubuntu/.config/rclone "$TMPDIR/app-configs/rclone" 2>/dev/null || true
cp -a /etc/docker/daem.* "$TMPDIR/app-configs/" 2>/dev/null || true
cp -a /etc/rclone "$TMPDIR/app-configs/rclone-etc" 2>/dev/null || true
# System configs for restore
mkdir -p "$TMPDIR/system-configs"
# SSH keys
cp -a /home/ubuntu/.ssh "$TMPDIR/system-configs/ssh-home" 2>/dev/null || true
sudo cp -a /root/.ssh "$TMPDIR/system-configs/ssh-root" 2>/dev/null || true
# Crontab dump
crontab -l > "$TMPDIR/system-configs/crontab-ubuntu" 2>/dev/null || true
sudo crontab -l > "$TMPDIR/system-configs/crontab-root" 2>/dev/null || true
# Docker daemon config
sudo cp -a /etc/docker "$TMPDIR/system-configs/docker-etc" 2>/dev/null || true
# Tunnel service
sudo tar czf "$TMPDIR/system-configs/tunnel-service.tar.gz" -C /opt tunnel-service 2>/dev/null || true
# Combined
tar czf "$TMPDIR/app-configs.tar.gz" -C "$TMPDIR" app-configs 2>/dev/null && echo "  ✅ app-configs" || ERRORS+=("app-configs tar failed")
sudo tar czf "$TMPDIR/system-configs.tar.gz" -C "$TMPDIR" system-configs 2>/dev/null && echo "  ✅ system-configs (SSH, crontab, docker, tunnel)" || ERRORS+=("system-configs tar failed")
rm -rf "$TMPDIR/app-configs" 2>/dev/null || true
sudo rm -rf "$TMPDIR/system-configs" 2>/dev/null || true

# ── 6. Weekly extras ────────────────────────────────────────
if [[ "$BACKUP_MODE" == "weekly" || "$BACKUP_MODE" == "monthly" ]]; then
    echo "[6/9] Weekly/Monthly: archiving projects & extras..."
    ionice -c2 -n7 tar czf "$TMPDIR/projects.tar.gz" \
        -C /home/ubuntu projects \
        --exclude='*/node_modules' \
        --exclude='*/venv' \
        --exclude='*/__pycache__' \
        --exclude='*/.git' \
        --exclude='open-codesign' \
        2>/dev/null && echo "  ✅ projects" || ERRORS+=("projects tar failed")

    ionice -c2 -n7 tar czf "$TMPDIR/extras.tar.gz" \
        -C /home/ubuntu \
        .cloakbrowser \
        .opencode \
        mqtt-api \
        2>/dev/null || true
    echo "  ✅ extras (best-effort)"
else
    echo "[6/9] Daily mode — skipping projects"
fi

# ── 7. Custom Docker Images (weekly/monthly) ─────────────────
if [[ "$BACKUP_MODE" == "weekly" || "$BACKUP_MODE" == "monthly" ]]; then
    echo "[7/9] Saving custom Docker images..."
    IMG_DIR="$TMPDIR/docker-images"
    mkdir -p "$IMG_DIR"
    # Images that are custom-built (not from public registries)
    CUSTOM_IMAGES=(
        "x10fjrc7dyd74dkogd1m3dvp:56f478a70eaf2b2401c4fa7b7d671ff7f03bc82e"
        "xnvshnbpkfjzgbxphbcgpqo5:240256db8ec719375c57f9e8f99150c1619a7134"
        "q8c0kc0sogos8cgwscgsccgg:bf3cb6959ab7cd56bf2c5dce2ed84dd0ff718552"
    )
    for img in "${CUSTOM_IMAGES[@]}"; do
        safe_name=$(echo "$img" | tr ':/' '-')
        if docker save "$img" | gzip > "$IMG_DIR/${safe_name}.tar.gz" 2>/dev/null; then
            echo "  ✅ $safe_name"
        else
            ERRORS+=("docker save $img failed")
            echo "  ❌ $safe_name FAILED"
        fi
 done
else
    echo "[7/9] Daily mode — skipping docker images"
fi

# ── 8. Monthly full home backup ─────────────────────────────
if [[ "$BACKUP_MODE" == "monthly" ]]; then
    echo "[8/9] Monthly: full home backup (excl caches)..."
    ionice -c2 -n7 tar czf "$TMPDIR/home-full.tar.gz" \
        -C /home/ubuntu \
        --exclude='.npm' \
        --exclude='.npm-global' \
        --exclude='.bun' \
        --exclude='.gradle' \
        --exclude='.local' \
        --exclude='.cache' \
        --exclude='.cloakbrowser' \
        --exclude='android-sdk' \
        --exclude='scraping-env' \
        --exclude='snap' \
        --exclude='*/node_modules' \
        --exclude='*/venv' \
        --exclude='*/__pycache__' \
        --exclude='*/.git' \
        . 2>/dev/null && echo "  ✅ home-full" || ERRORS+=("home-full tar failed")
else
    echo "[8/9] Skipping full home backup (daily/weekly mode)"
fi

# ── 9. Summary of what's backed up ──────────────────────────
echo "[9/9] Building archive summary..."
echo "Backup contents:" > "$TMPDIR/backup-manifest.txt"
echo "  Mode: ${BACKUP_MODE}" >> "$TMPDIR/backup-manifest.txt"
echo "  Date: ${DATE_SHORT}" >> "$TMPDIR/backup-manifest.txt"
echo "  Host: $(hostname)" >> "$TMPDIR/backup-manifest.txt"
ls -lh "$TMPDIR/"*.tar.gz 2>/dev/null | awk '{print "  " $5, $9}' >> "$TMPDIR/backup-manifest.txt" || true

# ── Compress final archive ───────────────────────────────────
echo "Compressing final archive..."
FINAL_NAME="vps-backup-${BACKUP_MODE}-${DATE_SHORT}.tar.gz"
FINAL_PATH="/tmp/${FINAL_NAME}"

ionice -c2 -n7 tar czf "$FINAL_PATH" -C "$TMPDIR" . 2>/dev/null

FINAL_SIZE=$(stat -c%s "$FINAL_PATH" 2>/dev/null || echo 0)

# ── Test mode: skip upload, just report ──────────────────────
if [[ "$BACKUP_MODE" == "test" ]]; then
    END_EPOCH=$(date +%s)
    DURATION=$(( END_EPOCH - START_EPOCH ))
    MINS=$(( DURATION / 60 ))
    SECS=$(( DURATION % 60 ))

    # List what's in the archive
    CONTENTS=$(tar tzf "$FINAL_PATH" 2>/dev/null | head -30 | tr '\n' ', ')
    echo "TEST: Archive contents: ${CONTENTS}"

    tg_send "🧪 <b>Backup TEST Selesai</b>
🕐 ${START_TIME}
⏱️ Durasi: ${MINS}m ${SECS}s
📦 Size: $(human_size $FINAL_SIZE)
📊 Load: ${LOAD1}
📋 Isi: databases, docker-volumes, coolify-state, hermes, system-configs
$( [[ ${#ERRORS[@]} -gt 0 ]] && echo "⚠️ Errors: ${ERRORS[*]}" || echo "✅ Tanpa error" )
⚡ <i>Upload SKIPPED (test mode)</i>"

    echo "========================================="
    echo "TEST backup done, no upload"
    echo "Size: $(human_size $FINAL_SIZE)"
    rm -f "$FINAL_PATH"
    exit 0
fi

# ── Upload to Google Drive ───────────────────────────────────
echo "Uploading to Google Drive..."
UPLOAD_OK=0
if rclone copy "$FINAL_PATH" "${RCLONE_BASE}/${BACKUP_MODE}/" \
    --no-check-certificate \
    --transfers 1 \
    --retries 3 \
    --retries-sleep 5s 2>&1; then
    UPLOAD_OK=1
    echo "Upload success"
else
    ERRORS+=("rclone upload failed")
    echo "Upload FAILED"
fi

# ── Cleanup old backups ──────────────────────────────────────
echo "Cleaning old backups..."
if [[ "$BACKUP_MODE" == "daily" ]]; then
    rclone delete "${RCLONE_BASE}/daily/" --min-age 7d --rmdirs 2>/dev/null || true
elif [[ "$BACKUP_MODE" == "weekly" ]]; then
    rclone delete "${RCLONE_BASE}/weekly/" --min-age 28d --rmdirs 2>/dev/null || true
elif [[ "$BACKUP_MODE" == "monthly" ]]; then
    rclone delete "${RCLONE_BASE}/monthly/" --min-age 84d --rmdirs 2>/dev/null || true
fi

# ── Cleanup local ────────────────────────────────────────────
rm -f "$FINAL_PATH"

# ── Final notification ───────────────────────────────────────
END_EPOCH=$(date +%s)
DURATION=$(( END_EPOCH - START_EPOCH ))
MINS=$(( DURATION / 60 ))
SECS=$(( DURATION % 60 ))

if [[ ${#ERRORS[@]} -eq 0 && $UPLOAD_OK -eq 1 ]]; then
    STATUS_ICON="✅"
    STATUS_TEXT="BERHASIL"
else
    STATUS_ICON="⚠️"
    STATUS_TEXT="SELESAI (dengan error)"
fi

ERR_LINES=""
if [[ ${#ERRORS[@]} -gt 0 ]]; then
    ERR_LINES=$'\n'"⚠️ <b>Errors (${#ERRORS[@]}):</b>"
    for e in "${ERRORS[@]}"; do
        ERR_LINES+=$'\n'"  ❌ ${e}"
    done
fi

tg_send "${STATUS_ICON} <b>Backup ${BACKUP_MODE^^} ${STATUS_TEXT}</b>

🕐 ${START_TIME} → $(date '+%H:%M:%S UTC')
⏱️ Durasi: ${MINS}m ${SECS}s
📦 Size: $(human_size $FINAL_SIZE)
☁️ GDrive: backups/vps/${BACKUP_MODE}/${FINAL_NAME}
📊 Load: ${LOAD1}${ERR_LINES}"

echo "========================================="
echo "Backup ${BACKUP_MODE} ${STATUS_TEXT}"
echo "Size: $(human_size $FINAL_SIZE)"
echo "Duration: ${MINS}m ${SECS}s"
echo "========================================="
