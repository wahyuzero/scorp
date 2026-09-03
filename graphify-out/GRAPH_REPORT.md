# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 122 files · ~103,528 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1174 nodes · 2912 edges · 40 communities (31 shown, 2 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 266 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `78253b41`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- TruncateStr
- handleAction
- GetStringArg
- chat.go
- config_paths.go
- client.go
- testing.T
- rag_vector.go
- collector_system.go
- ConfigMgr
- runSubagent
- collector_system_native.go
- wizard.go
- startCLI
- time.Time
- session_search_fts5.go
- collector_security.go
- RunAgentLoop
- runSelfReview
- skills.go
- checker.go
- collector_coolify.go
- install.sh
- LoopController
- scorp-agent
- AI Coding Agent Modernization Roadmap & Blueprint (v2.0)
- scorp
- Scorp vs ZeroClaw: Competitive Analysis, Strategic Positioning, & Growth Playbook
- CleanupSessionsLoop
- 9. Pilar Upgrade 7: Ultra-Low-RAM Web Engine (< 5MB RAM di VPS 512MB & Termux)
- getAgentSystemPrompt
- 🚀 4. Taktik Pertumbuhan Komunitas untuk Scorp (Playbook 10k+ Stars)
- 5. Pilar Upgrade 3: Surgical Diff Editing (Bukan Full Overwrite)

## God Nodes (most connected - your core abstractions)
1. `handleAction()` - 59 edges
2. `main()` - 54 edges
3. `GetStringArg()` - 40 edges
4. `TruncateStr()` - 36 edges
5. `RunAgentLoop()` - 35 edges
6. `init()` - 30 edges
7. `HandleModelCallback()` - 29 edges
8. `startCLI()` - 28 edges
9. `ModelConfig` - 28 edges
10. `resumeAgentLoop()` - 23 edges

## Surprising Connections (you probably didn't know these)
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `ExecuteSQL()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/db.go → agent/loop.go
- `ExecuteShell()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/exec.go → agent/loop.go
- `ExecuteGit()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/git.go → agent/loop.go
- `ExecuteShell()` --calls--> `IsDangerousCommand()`  [INFERRED]
  tools/exec.go → agent/prompt.go

## Import Cycles
- None detected.

## Communities (40 total, 2 thin omitted)

### Community 0 - "TruncateStr"
Cohesion: 0.06
Nodes (91): context.Context, TruncateStr(), anthropicRequest, anthropicResponse, anthropicTool, callAnthropic(), CallAnthropicWithTools(), buildCommandCodePayload() (+83 more)

### Community 1 - "handleAction"
Cohesion: 0.06
Nodes (74): IsUserActive(), net/http.Request, net/http.ResponseWriter, commandLoop(), formatBasicAlert(), handleAction(), hourlyLoop(), isCLIMode() (+66 more)

### Community 2 - "GetStringArg"
Cohesion: 0.06
Nodes (62): init(), init(), ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill() (+54 more)

### Community 3 - "chat.go"
Cohesion: 0.13
Nodes (38): agentAutoStop(), appendSessionHistory(), ClearChatSession(), collectTableLines(), convertInlineMarkdown(), convertTableToList(), EnterAgentMode(), ExitAgentMode() (+30 more)

### Community 4 - "config_paths.go"
Cohesion: 0.05
Nodes (63): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+55 more)

### Community 5 - "client.go"
Cohesion: 0.05
Nodes (58): RegisterAutonomous(), init(), init(), init(), init(), bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage (+50 more)

### Community 6 - "testing.T"
Cohesion: 0.07
Nodes (54): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+46 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (59): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), checkGDriveMount() (+51 more)

### Community 9 - "ConfigMgr"
Cohesion: 0.06
Nodes (53): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+45 more)

### Community 10 - "runSubagent"
Cohesion: 0.06
Nodes (57): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+49 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.07
Nodes (46): TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), FormatDuration(), TopProcess, CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg() (+38 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (43): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), SaveModelConfig() (+35 more)

### Community 13 - "startCLI"
Cohesion: 0.18
Nodes (23): HasPendingConfirmation(), executeOneShot(), executeTurn(), formatFinalResponse(), formatTerminalText(), hasDebugFlag(), isTerminal(), printBanner() (+15 more)

### Community 14 - "time.Time"
Cohesion: 0.18
Nodes (24): time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult() (+16 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 17 - "RunAgentLoop"
Cohesion: 0.15
Nodes (29): AgentMessage, sendScorpReply(), contentPart, imageURL, base64Encode(), browserToolsWereUsed(), buildThinkingMessage(), clearPendingConfirmation() (+21 more)

### Community 18 - "runSelfReview"
Cohesion: 0.20
Nodes (14): memoryFact, AgentMessage, runSelfReview(), TestSelfReviewIntegration(), LoadJSON(), SaveJSON(), SaveJSONPerm(), deleteMemory() (+6 more)

### Community 19 - "skills.go"
Cohesion: 0.22
Nodes (14): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+6 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 31 - "AI Coding Agent Modernization Roadmap & Blueprint (v2.0)"
Cohesion: 0.13
Nodes (15): 10. Roadmap Eksekusi Teknis, 1. Lanskap Model AI Modern & Verifikasi Data (2026), 2. Evaluasi Arsitektur Scorp & Termagent Saat Ini, 3. Pilar Upgrade 1: Multi-Model & Fallback Engine, 4. Pilar Upgrade 2: Prompt Caching & Token Economics, 6. Pilar Upgrade 4: Standardisasi MCP (Model Context Protocol), 7. Pilar Upgrade 5: Subagent & Task Delegation System, 8. Pilar Upgrade 6: Deep Termux & Mobile Performance Optimization (+7 more)

### Community 32 - "scorp"
Cohesion: 0.17
Nodes (12): Build, Commands, Configuration, .env, License, models.json, Project layout, scorp (+4 more)

### Community 33 - "Scorp vs ZeroClaw: Competitive Analysis, Strategic Positioning, & Growth Playbook"
Cohesion: 0.20
Nodes (10): 🔍 1. Mengapa ZeroClaw Berhasil Meraih 32k+ Stars?, 🎯 2. Kelemahan Fatal ZeroClaw di Lingkungan Mobile & Potato Hardware, 🦂 3. Positioning Strategis Scorp: "The Mobile & Edge Native AI Agent", A. Narasi "Anti-Python & Anti-Docker" yang Sangat Kuat, B. Arsitektur Modular Crate (Community Flywheel Effect), C. Multi-Channel Ubiquity (Bukan Cuma Terminal CLI), 📊 Matriks Fitur Lengkap: Scorp vs ZeroClaw, 📌 Ringkasan Eksekutif (+2 more)

### Community 34 - "CleanupSessionsLoop"
Cohesion: 0.33
Nodes (5): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), cleanupAgentSessions(), TGDocument

### Community 35 - "9. Pilar Upgrade 7: Ultra-Low-RAM Web Engine (< 5MB RAM di VPS 512MB & Termux)"
Cohesion: 0.33
Nodes (6): 💡 3 Solusi Arsitektur Low-RAM Browser:, 9. Pilar Upgrade 7: Ultra-Low-RAM Web Engine (< 5MB RAM di VPS 512MB & Termux), 🛑 Masalah Utama:, 🌐 Opsi 1 (Default): Zero-RAM HTTP Fetch + Readability Parser (~2MB RAM), ☁️ Opsi 2: Remote / Offloaded Browser API (RAM Lokal = 0 MB!), ⚙️ Opsi 3: Chromium Lokal Teroptimasi Ekstrem (~80MB–120MB RAM)

### Community 36 - "getAgentSystemPrompt"
Cohesion: 0.40
Nodes (3): getSharedMemorySummary(), getAgentSystemPrompt(), GetMemorySummary()

### Community 37 - "🚀 4. Taktik Pertumbuhan Komunitas untuk Scorp (Playbook 10k+ Stars)"
Cohesion: 0.40
Nodes (5): 🚀 4. Taktik Pertumbuhan Komunitas untuk Scorp (Playbook 10k+ Stars), Taktik 1: Standardisasi Interface Plugin Go (Memancing Kontributor), Taktik 2: Aktifkan Mode Daemon Telegram 24/7, Taktik 3: Demo Visual yang Menghipnotis (GIF / Video di README), Taktik 4: Strategi Peluncuran Global

### Community 39 - "5. Pilar Upgrade 3: Surgical Diff Editing (Bukan Full Overwrite)"
Cohesion: 0.67
Nodes (3): 5. Pilar Upgrade 3: Surgical Diff Editing (Bukan Full Overwrite), 🛑 Masalah Cara Lama (Full File Rewrite):, ✂️ Solusi Modern (Chunk Replacement / Exact Match):

## Knowledge Gaps
- **50 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+45 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 113 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `handleAction` to `TruncateStr`, `CleanupSessionsLoop`, `config_paths.go`, `client.go`, `rag_vector.go`, `collector_system.go`, `ConfigMgr`, `runSubagent`, `collector_system_native.go`, `wizard.go`, `startCLI`, `time.Time`, `collector_security.go`, `RunAgentLoop`, `runSelfReview`, `skills.go`, `checker.go`?**
  _High betweenness centrality (0.167) - this node is a cross-community bridge._
- **Why does `handleAction()` connect `handleAction` to `TruncateStr`, `chat.go`, `config_paths.go`, `client.go`, `collector_system.go`, `collector_system_native.go`, `wizard.go`, `time.Time`, `collector_security.go`, `RunAgentLoop`, `skills.go`, `checker.go`, `collector_coolify.go`?**
  _High betweenness centrality (0.154) - this node is a cross-community bridge._
- **Why does `RunAgentLoop()` connect `RunAgentLoop` to `TruncateStr`, `handleAction`, `GetStringArg`, `chat.go`, `getAgentSystemPrompt`, `rag_vector.go`, `runSubagent`, `startCLI`, `time.Time`?**
  _High betweenness centrality (0.090) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `main()` (e.g. with `hasDebugFlag()` and `startCLI()`) actually correct?**
  _`main()` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _50 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `TruncateStr` be split into smaller, more focused modules?**
  _Cohesion score 0.06059405940594059 - nodes in this community are weakly interconnected._
- **Should `handleAction` be split into smaller, more focused modules?**
  _Cohesion score 0.05617283950617284 - nodes in this community are weakly interconnected._