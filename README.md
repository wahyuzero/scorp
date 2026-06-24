# scorp

A small project — a personal AI agent for VPS management, written in Go.

It started as a monitoring script, grew into a Telegram bot with LLM integration, and is slowly becoming a general-purpose agent. Still rough around the edges, but it works (mostly).

Not affiliated with any company or project. Built for personal use, shared in case someone finds it useful.

---

## What it does

- **Chat with tools** — send a message, the AI decides which tools to use
- **CLI + Telegram** — interactive REPL by default, optional Telegram bot
- **VPS monitoring** — CPU/RAM/disk alerts, SSH login detection
- **Multiple LLM providers** — OpenAI, Anthropic, Gemini, Groq, DeepSeek, Z.ai, OpenRouter, local (Ollama/LM Studio)
- **MCP client** — connects to Model Context Protocol servers via stdio
- **Self-update** — checks GitHub releases, one-command update via `/update`

That's the idea anyway. Some features work better than others.

---

## Try it

```bash
git clone https://github.com/OWNER/scorp.git
cd scorp
./install.sh
```

The installer walks you through provider setup and optional Telegram configuration.

If you just want CLI mode without Telegram:

```bash
scorp          # interactive REPL
```

---

## Configuration

### .env (Telegram + API keys)

Copy `.env.example` to `.env` and fill in what you need:

```
TELEGRAM_BOT_TOKEN=...     # optional, skip for CLI-only
OPENAI_API_KEY=sk-...      # set only what you use
GITHUB_REPO=owner/scorp    # for self-update checks
```

### models.json

Located at `~/.scorp/models.json`. Defines which LLM models to use:

```json
{
  "default_model": "glm-4.7",
  "models": {
    "glm-4.7": {
      "provider": "zai-coding",
      "model": "glm-4.7",
      "key_env": "GLM_API_KEY",
      "base_url": "https://api.z.ai/api/coding/paas/v4",
      "max_tokens": 4096,
      "api": "openai"
    }
  }
}
```

---

## Commands

**CLI:**
- `scorp` — start REPL
- `scorp update` — check and install updates
- `scorp version` — show current version
- `scorp --cli` — force CLI mode

**Telegram:**
- `/help` — show commands
- `/update` — check & install updates
- `/version` — show version
- `/model` — change AI model
- `/start` — interactive menu

---

## Self-update

If `GITHUB_REPO` is set in `.env`, scorp checks for new GitHub releases on startup and notifies you. Run `/update` (Telegram) or `scorp update` (CLI) to download and install.

The update process:
1. Fetches latest release from GitHub API
2. Downloads the binary matching your OS/arch
3. Replaces the running binary (backup kept at `scorp.bak`)
4. Restarts the service if running under systemd

If no prebuilt binary is available for your platform, it tells you to `git pull && make`.

---

## Build

Requires Go 1.25+ and CGO (for SQLite).

```bash
make              # production build
make deploy       # build + install + restart service
make VERSION=v1.0.0 deploy  # with version tag
```

---

## Project layout

```
scorp/
├── main.go           # entry point
├── cli.go            # CLI REPL
├── telegram/         # bot transport
├── agent/            # agent loop, prompt, chat
├── bootstrap/        # tool registry
├── tools/            # tool implementations
├── models/           # LLM config
├── updater/          # self-update logic
├── config/           # paths + env
├── session/          # session DB (SQLite)
├── ...
└── install.sh
```

---

## Status

Personal project, developed sporadically. No guarantees of stability or support. Use at your own risk.

---

## License

MIT — see [LICENSE](LICENSE).

## Acknowledgements

Inspired by various open-source agent projects. Uses standard Go libraries and a few third-party packages listed in `go.mod`.
