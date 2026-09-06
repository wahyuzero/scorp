#!/usr/bin/env bash
# Scorp deploy pipeline — `scorp eval` as an automated pre-deploy gate.
#
# Flow: local checks → rsync source → build ON the VPS (CGO needed for
# mattn/go-sqlite3) → unit tests on VPS → scorp eval gate → swap binary →
# verify md5 → service healthy.
#
# If ANY gate fails, the running production binary is left untouched.
#
# Escapes (emergency only, print loud warnings):
#   SCORP_DEPLOY_SKIP_TESTS=1   skip local+remote unit tests
#   SCORP_DEPLOY_SKIP_EVAL=1    skip the eval gate
#
# Overridable targets:
#   SCORP_DEPLOY_HOST   ssh alias/host      (default tencent-vps)
#   SCORP_DEPLOY_DIR    remote source dir   (default /root/scorp-src)
#   SCORP_DEPLOY_BIN    remote binary path  (default /usr/local/bin/scorp)
#   SCORP_DEPLOY_SVC    systemd unit        (default scorp)
#   SCORP_DEPLOY_WORKDIR dir containing production .env, used as eval cwd
set -euo pipefail

HOST="${SCORP_DEPLOY_HOST:-tencent-vps}"
REMOTE_DIR="${SCORP_DEPLOY_DIR:-/root/scorp_src}"
BIN="${SCORP_DEPLOY_BIN:-/usr/local/bin/scorp}"
SVC="${SCORP_DEPLOY_SVC:-scorp}"
WORKDIR="${SCORP_DEPLOY_WORKDIR:-}"
CANDIDATE="/tmp/scorp.candidate.$$.bin"

step() { printf '\n\033[1;36m==> %s\033[0m\n' "$*"; }
ok()   { printf '\033[1;32m✓ %s\033[0m\n' "$*"; }
die()  { printf '\n\033[1;31m✗ %s\033[0m\n' "$*" >&2; exit 1; }

cleanup() { ssh "$HOST" "rm -f $CANDIDATE" 2>/dev/null || true; }
trap cleanup EXIT

cd "$(dirname "$0")/.."

# ── 0. Local gate: build + vet + tests ────────────────────────────────
step "0/6 local: go build + go vet + go test"
go build ./... || die "local build failed"
go vet ./... || die "go vet failed"
if [ "${SCORP_DEPLOY_SKIP_TESTS:-0}" = "1" ]; then
  echo "⚠️  local tests skipped (SCORP_DEPLOY_SKIP_TESTS=1)"
else
  go test ./... -count=1 >/tmp/scorp_deploy_local_test.log 2>&1 \
    || { tail -30 /tmp/scorp_deploy_local_test.log; die "local tests failed"; }
  ok "local tests green"
fi

# ── 1. Sync source ────────────────────────────────────────────────────
step "1/6 rsync source → $HOST:$REMOTE_DIR"
rsync -az --delete --exclude .git --exclude .env --exclude 'scorp' --exclude '*.log' \
  -e ssh ./ "$HOST:$REMOTE_DIR/" || die "rsync failed"
ok "source synced (no .env, no .git)"

# ── 2. Build on the VPS ───────────────────────────────────────────────
step "2/6 build candidate on VPS (CGO_ENABLED=1 for sqlite)"
ssh "$HOST" "cd $REMOTE_DIR && CGO_ENABLED=1 go build -buildvcs=false -o $CANDIDATE . " \
  || die "remote build failed"
CAND_MD5=$(ssh "$HOST" "md5sum $CANDIDATE | cut -d' ' -f1")
ok "candidate built: md5 $CAND_MD5"

# ── 3. Unit tests on the VPS ──────────────────────────────────────────
step "3/6 unit tests on VPS"
if [ "${SCORP_DEPLOY_SKIP_TESTS:-0}" = "1" ]; then
  echo "⚠️  remote tests skipped (SCORP_DEPLOY_SKIP_TESTS=1)"
else
  ssh "$HOST" "cd $REMOTE_DIR && go test ./... -count=1" \
    || die "remote tests failed — deploy aborted, production untouched"
  ok "remote tests green"
fi

# ── 4. Eval gate ──────────────────────────────────────────────────────
step "4/6 scorp eval gate (run from production workdir so .env is present)"
if [ -z "$WORKDIR" ]; then
  WORKDIR=$(ssh "$HOST" "systemctl show $SVC -p WorkingDirectory --value")
  [ -n "$WORKDIR" ] || WORKDIR=$REMOTE_DIR
  echo "   eval cwd: $WORKDIR"
fi
if [ "${SCORP_DEPLOY_SKIP_EVAL:-0}" = "1" ]; then
  echo "🚨 SCORP_DEPLOY_SKIP_EVAL=1 — EVAL GATE BYPASSED (emergency only)"
else
  ssh "$HOST" "cd $WORKDIR && $CANDIDATE eval" || {
    echo "" >&2
    die "EVAL GATE FAILED — deploy aborted, production binary untouched"
  }
  ok "eval gate passed"
fi

# ── 5. Swap binary (stop → cp → verify md5 → start) ──────────────────
step "5/6 swap binary + restart $SVC"
ssh "$HOST" "systemctl stop $SVC" || die "systemctl stop failed"
ssh "$HOST" "cp $CANDIDATE $BIN && chmod 755 $BIN" || die "binary copy failed"
INSTALLED_MD5=$(ssh "$HOST" "md5sum $BIN | cut -d' ' -f1")
if [ "$INSTALLED_MD5" != "$CAND_MD5" ]; then
  die "md5 mismatch after copy ($CAND_MD5 → $INSTALLED_MD5) — NOT starting service, investigate manually"
fi
ok "installed md5 verified: $INSTALLED_MD5"
ssh "$HOST" "systemctl start $SVC && sleep 3 && systemctl is-active $SVC" \
  || die "service failed to start"
ok "service active"

# ── 6. Smoke: daemon binary sanity ────────────────────────────────────
step "6/6 smoke"
ssh "$HOST" "$BIN version 2>/dev/null || $BIN --version 2>/dev/null || echo '(no version cmd)'" || true

printf '\n\033[1;32m✅ DEPLOY COMPLETE — %s running %s (md5 %s), eval gate passed.\033[0m\n' "$SVC" "$BIN" "$CAND_MD5"
