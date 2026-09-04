# 🦂 Scorp

**Small agent, big tools.** A tiny, ultra-fast Go binary that connects modern LLMs (DeepSeek V4, Gemini 2.0, Claude 3.5, OpenAI, Ollama) to your system — shell, files, web, network, code — accessible from terminal, Telegram, or embedded Web Dashboard.

Runs natively with a memory footprint of **under 15MB RAM** on Linux VPS, edge servers, and Android Termux.

> 📋 **Modernization Notes & Comparisons:**
> - [Scorp v2.0 Architecture & Upgrade Notes](MODERNIZATION_NOTES.md)
> - [Head-to-Head Comparison: Scorp vs PicoClaw vs ZeroClaw](docs/COMPARISON_SCORP_PICOCLAW_ZEROCLAW.md)

---

## ⚡ What Makes Scorp v2.0 Special

- **Agent-First, Not a Chatbot** — Every message runs through an autonomous ReAct tool loop. Ask it to fix code, and it surgically patches the exact lines.
- **Surgical Diff Editing** (`replace_file_content`) — Replaces exact or fuzzy content chunks instead of rewriting entire files. Saves 90%+ tokens and executes in sub-seconds.
- **Ultra-Low-RAM Web Engine** (`read_url`) — Read documentation and articles in clean Markdown using Firefox Readability (<5MB RAM), with automatic cloud offload fallback.
- **3-Tier Autonomy & Security Sandbox** — Choose between `readonly` (safe inspection/audit), `supervised` (prompts for dangerous commands), or `yolo` (full unattended autonomy).
- **Outbound Secret Redactor** — Automatically masks API keys, tokens, private keys, and environment secrets before tool outputs reach the LLM or chat logs.
- **Real-Time Steering Queue** — Intercept and redirect the agent mid-run without killing the session.
- **Standard Operating Procedures (SOP) Engine** — Reusable, declarative automation playbooks (`scorp sop run health_audit`).
- **Micro HTTP Gateway & Dashboard** — Built-in web UI on port 8080 running on just **~1.5 MB RAM**.
- **Model Context Protocol (MCP)** — Native Stdio and SSE/HTTP JSON-RPC 2.0 client to connect external MCP servers.
- **Multi-Session Management** — Manage, switch, rename, and delete conversation sessions (`scorp -s <name>` or `/session`).
- **Native Android Termux & Low-Resource VPS** — Direct Termux:API integration (WakeLock, Android notifications) and single static ARM64 binary (`make termux`).

---

## 🚀 Quick Start (Under 60 Seconds)

```bash
git clone https://github.com/wahyuzero/scorp.git
cd scorp
make
./scorp setup
```

The interactive wizard will guide you through choosing your AI provider (Command Code, DeepSeek, Gemini, OpenAI, Claude, Ollama) and safety mode.

### Launch Interactive CLI Chat:
```bash
./scorp
```

### Run One-Shot Tasks:
```bash
./scorp "read the first 10 lines of main.go and summarize"
```

### Launch Embedded Web Dashboard (< 2MB RAM):
```bash
./scorp gateway --port 8080
```
Open `http://localhost:8080` in your phone or PC browser.

---

## 🛠️ CLI Slash Commands

| Command | Description |
| :--- | :--- |
| `/help` | Show command help and usage examples |
| `/models` | List all configured models and key status |
| `/model <name>` | Switch active model on the fly |
| `/mode <level>` | Switch autonomy mode (`readonly`, `supervised`, `yolo`) |
| `/session` | List saved sessions and switch/rename/delete |
| `/sop [run]` | List or execute automated SOP playbooks |
| `/receipts` | View recent cryptographic SHA-256 tool execution receipts |
| `/tools` | List all registered tools |
| `/cost` | Display token usage and estimated spend |
| `/clear` | Reset current session history |
| `/exit` | Quit Scorp |

---

## 📱 Build for Android Termux

Compile an ultra-lightweight static ARM64 binary ready to run directly on Termux without external libc dependencies:

```bash
make termux
```
*(Produces `scorp-termux` ~16MB single binary).*

---

## 📄 License

MIT License.
