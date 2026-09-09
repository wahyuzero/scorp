#!/usr/bin/env bash
set -e

# ─────────────────────────────────────────────────────────────────────────────
# Scorp Agent — One-Line Universal Installer
# Works on Linux (x86_64, aarch64, armv7), macOS (Apple Silicon & Intel),
# and Android (Termux).
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/wahyuzero/scorp/master/install.sh | bash
# ─────────────────────────────────────────────────────────────────────────────

REPO="wahyuzero/scorp"
FALLBACK_TAG="v0.7.1"

# ANSI Colors
BOLD="\033[1m"
GREEN="\033[1;32m"
CYAN="\033[1;36m"
YELLOW="\033[1;33m"
RED="\033[1;31m"
RESET="\033[0m"

echo -e "${CYAN}╔══════════════════════════════════════════════════════════════════════╗${RESET}"
echo -e "${CYAN}║                     🦂 SCORP AGENT INSTALLER                         ║${RESET}"
echo -e "${CYAN}║            Autonomous, Resilient & Lightweight AI Assistant          ║${RESET}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════════════╝${RESET}"
echo ""

# 1. Detect OS & Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Termux Detection
IS_TERMUX=false
if [ -n "$PREFIX" ] && [[ "$PREFIX" == *"com.termux"* ]]; then
    IS_TERMUX=true
elif [ -d "/data/data/com.termux" ]; then
    IS_TERMUX=true
fi

TARGET_OS=""
TARGET_ARCH=""

# Map OS
case "$OS" in
    linux*)
        if [ "$IS_TERMUX" = true ]; then
            TARGET_OS="android"
        else
            TARGET_OS="linux"
        fi
        ;;
    darwin*)
        TARGET_OS="darwin"
        ;;
    msys*|mingw*|cygwin*)
        TARGET_OS="windows"
        ;;
    *)
        echo -e "${RED}❌ Unsupported Operating System: $OS${RESET}"
        exit 1
        ;;
esac

# Map Architecture
case "$ARCH" in
    x86_64|amd64)
        TARGET_ARCH="amd64"
        ;;
    aarch64|arm64)
        TARGET_ARCH="arm64"
        ;;
    armv7l|armv7|armhf)
        TARGET_ARCH="armv7"
        ;;
    *)
        echo -e "${RED}❌ Unsupported CPU Architecture: $ARCH${RESET}"
        exit 1
        ;;
esac

ASSET_NAME="scorp_${TARGET_OS}_${TARGET_ARCH}"
if [ "$TARGET_OS" = "windows" ]; then
    ASSET_NAME="scorp_windows_amd64.exe"
fi

echo -e "🔎 Detected environment: ${BOLD}${TARGET_OS}/${TARGET_ARCH}${RESET} (Asset: ${CYAN}${ASSET_NAME}${RESET})"

# 2. Resolve Latest Release Tag
echo -n "📡 Checking latest release... "
LATEST_TAG=""
if command -v curl >/dev/null 2>&1; then
    LATEST_TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || true)
fi

if [ -z "$LATEST_TAG" ]; then
    LATEST_TAG="$FALLBACK_TAG"
    echo -e "${YELLOW}${LATEST_TAG} (fallback)${RESET}"
else
    echo -e "${GREEN}${LATEST_TAG}${RESET}"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${ASSET_NAME}"

# 3. Download Binary
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT
TMP_BIN="${TMP_DIR}/scorp"

echo -e "⬇️  Downloading Scorp ${LATEST_TAG}..."
if command -v curl >/dev/null 2>&1; then
    curl -fsSL --progress-bar "$DOWNLOAD_URL" -o "$TMP_BIN"
elif command -v wget >/dev/null 2>&1; then
    wget -q --show-progress "$DOWNLOAD_URL" -O "$TMP_BIN"
else
    echo -e "${RED}❌ Neither curl nor wget was found on your system.${RESET}"
    exit 1
fi

chmod +x "$TMP_BIN"

# 4. Determine Install Location
INSTALL_DIR=""
USE_SUDO=false

if [ "$IS_TERMUX" = true ]; then
    INSTALL_DIR="$PREFIX/bin"
elif [ "$(id -u)" -eq 0 ]; then
    INSTALL_DIR="/usr/local/bin"
elif [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
elif command -v sudo >/dev/null 2>&1; then
    INSTALL_DIR="/usr/local/bin"
    USE_SUDO=true
else
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

TARGET_BIN="${INSTALL_DIR}/scorp"

echo -e "📦 Installing binary to ${BOLD}${TARGET_BIN}${RESET}..."
if [ "$USE_SUDO" = true ]; then
    sudo mv "$TMP_BIN" "$TARGET_BIN"
    sudo chmod +x "$TARGET_BIN"
else
    mkdir -p "$INSTALL_DIR"
    mv "$TMP_BIN" "$TARGET_BIN"
    chmod +x "$TARGET_BIN"
fi

# Ensure install dir is in PATH if installed to ~/.local/bin
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    SHELL_RC=""
    if [ -n "$ZSH_VERSION" ] || [[ "$SHELL" == *"zsh"* ]]; then
        SHELL_RC="$HOME/.zshrc"
    else
        SHELL_RC="$HOME/.bashrc"
    fi

    if [ -f "$SHELL_RC" ]; then
        echo "export PATH=\"\$PATH:${INSTALL_DIR}\"" >> "$SHELL_RC"
        echo -e "${YELLOW}ℹ️  Added ${INSTALL_DIR} to PATH in ${SHELL_RC}. Run: source ${SHELL_RC}${RESET}"
    fi
    export PATH="$PATH:${INSTALL_DIR}"
fi

echo -e "${GREEN}✅ Scorp installed successfully!${RESET}"
echo ""

# 5. Dependency & Security Check
if [ "$TARGET_OS" = "linux" ]; then
    if ! command -v bwrap >/dev/null 2>&1; then
        echo -e "${YELLOW}💡 Note: bubblewrap is not installed.${RESET}"
        echo -e "   For root-filesystem isolation, you can install it anytime: ${BOLD}sudo apt install -y bubblewrap${RESET}"
        echo ""
    fi
fi

# 6. Interactive Setup Hand-off
RUN_SETUP=false
if [ -e /dev/tty ]; then
    # Redirect stdin from tty so we can read user prompt even in `curl | bash` pipe
    exec < /dev/tty
    echo -n -e "${BOLD}Would you like to run the interactive setup wizard now? [Y/n]: ${RESET}"
    read -r resp || resp="y"
    resp=$(echo "$resp" | tr '[:upper:]' '[:lower:]')
    if [ -z "$resp" ] || [ "$resp" = "y" ] || [ "$resp" = "yes" ]; then
        RUN_SETUP=true
    fi
fi

if [ "$RUN_SETUP" = true ]; then
    echo ""
    "$TARGET_BIN" setup
else
    echo -e "To configure Scorp anytime, simply run:"
    echo -e "  ${CYAN}${BOLD}scorp setup${RESET}"
    echo ""
    echo -e "To start interactive chat:"
    echo -e "  ${CYAN}${BOLD}scorp --cli${RESET}"
    echo ""
fi
