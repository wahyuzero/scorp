# scorp

**Small agent, big tools.** A tiny Go binary that connects an LLM to your system — shell, files, browser, network, code — accessible from terminal or Telegram.

Not a framework, not a platform. Just one binary that runs on anything and does things when you ask it to.

---

## What it does

- **Agent, not chatbot** — every message goes through a tool-use loop. Ask it to check disk space, it runs `df -h`. Ask it to fix a config, it patches the file.
- **Shell, files, code** — read/write/patch files, run commands, execute Python
- **Browser** — headless Chrome via chromedp (navigate, click, scrape, screenshot)
- **Network** — fetch URLs, search the web, make HTTP requests
- **CLI + Telegram** — interactive REPL by default, optional Telegram bot for remote access
- **Monitoring** — CPU/RAM/disk alerts, SSH login detection (toggle on/off)
- **Multi-provider** — OpenAI, Anthropic, Gemini, Groq, DeepSeek, OpenRouter, Z.ai, Ollama, LM Studio
- **MCP client** — connect to Model Context Protocol servers
- **Self-update** — check GitHub releases, update with one command

Some features work better than others.

---

## Try it

```bash
git clone https://github.com/wahyuzero/scorp.git
cd scorp
./install.sh
```

The installer walks you through provider setup and optional Telegram configuration.

CLI mode without Telegram:

```bash
scorp          # interactive REPL
```

---

## Configuration

### .env

Copy `.env.example` to `.env` and fill in what you need:

```
TELEGRAM_BOT_TOKEN=...     # optional, skip for CLI-only
OPENAI_API_KEY=sk-...      # set only what you use
GITHUB_REPO=wahyuzero/scorp    # for self-update checks
```

### models.json

Located at `~/.scorp/models.json`. Defines which LLM models to use.

**OpenAI example:**

```json
{
  "default_model": "gpt-4o-mini",
  "agent_model": "gpt-4o-mini",
  "models": {
    "gpt-4o-mini": {
      "provider": "openai",
      "model": "gpt-4o-mini",
      "key_env": "OPENAI_API_KEY",
      "base_url": "https://api.openai.com/v1",
      "max_tokens": 4096,
      "api": "openai"
    }
  },
  "routing_rules": {
    "agent": "gpt-4o-mini",
    "chat": "gpt-4o-mini"
  },
  "fallback_on_error": ["rate_limit", "timeout", "server_error"]
}
```

**Anthropic Claude example:**

```json
{
  "default_model": "claude-sonnet-4-20250514",
  "agent_model": "claude-sonnet-4-20250514",
  "models": {
    "claude-sonnet-4-20250514": {
      "provider": "anthropic",
      "model": "claude-sonnet-4-20250514",
      "key_env": "ANTHROPIC_API_KEY",
      "base_url": "https://api.anthropic.com",
      "max_tokens": 4096,
      "api": "anthropic"
    }
  },
  "routing_rules": {
    "agent": "claude-sonnet-4-20250514",
    "chat": "claude-sonnet-4-20250514"
  },
  "fallback_on_error": ["rate_limit", "timeout", "server_error"]
}
```

**Local (Ollama, no API key):**

```json
{
  "default_model": "llama3.2",
  "agent_model": "llama3.2",
  "models": {
    "llama3.2": {
      "provider": "ollama",
      "model": "llama3.2",
      "key_env": "",
      "base_url": "http://localhost:11434/v1",
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

Any text that doesn't start with `/` is sent to the agent with full tool access.

---

## Self-update

Set `GITHUB_REPO` in `.env` to enable update checks. Scorp checks for new GitHub releases on startup and notifies you.

```
scorp update       # CLI
/update            # Telegram
```

The update process:
1. Fetches latest release from GitHub API
2. Downloads the binary matching your OS/arch
3. Replaces the running binary (backup kept at `scorp.bak`)
4. Restarts the service if running under systemd

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
├── tools/            # 48 tool implementations
├── models/           # LLM provider config
├── updater/          # self-update logic
├── config/           # paths + env
├── session/          # session DB (SQLite + FTS5)
└── install.sh
```

---

## Status

Personal project, developed sporadically. No guarantees of stability or support. Use at your own risk.

---

## License

MIT — see [LICENSE](LICENSE).
