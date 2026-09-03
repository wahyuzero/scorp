# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 160 files · ~110,424 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1311 nodes · 3192 edges · 57 communities (46 shown, 2 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 373 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `af9049f2`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- TruncateStr
- HandleTelegramAction
- testing.T
- getSession
- ConfigManager
- client.go
- compaction_test.go
- rag_vector.go
- collector_system.go
- collector_security.go
- runSubagent
- uptime.go
- wizard.go
- startCLI
- time.Time
- session_search_fts5.go
- collector_system_native.go
- ConfigMgr
- ScorpPath
- skills.go
- checker.go
- GetProvider
- install.sh
- RunAgentSessionLoop
- scorp-agent
- AI Coding Agent Modernization Roadmap & Blueprint (v2.0)
- scorp
- Scorp vs ZeroClaw: Competitive Analysis, Strategic Positioning, & Growth Playbook
- RegisterTool
- 9. Pilar Upgrade 7: Ultra-Low-RAM Web Engine (< 5MB RAM di VPS 512MB & Termux)
- IsRateLimitError
- 🚀 4. Taktik Pertumbuhan Komunitas untuk Scorp (Playbook 10k+ Stars)
- 5. Pilar Upgrade 3: Surgical Diff Editing (Bukan Full Overwrite)
- chat.go
- GetStringArg
- GetAllTools
- init
- TestSessionManager
- patch.go
- StorePendingConfirmation
- TestSteeringQueue
- helpers.go
- read_url.go
- ExecuteTermuxAPI
- collector_native_test.go
- ExecuteTodo
- RedactSecrets
- HandleUploadInAgentMode

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 54 edges
2. `GetStringArg()` - 43 edges
3. `StartDaemon()` - 40 edges
4. `RunAgentSessionLoop()` - 36 edges
5. `TruncateStr()` - 36 edges
6. `ModelConfig` - 36 edges
7. `startCLI()` - 34 edges
8. `init()` - 32 edges
9. `ChatMessage` - 29 edges
10. `HandleModelCallback()` - 29 edges

## Surprising Connections (you probably didn't know these)
- `ExecuteShell()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/exec.go → agent/confirmation.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `ExecuteShell()` --calls--> `IsDangerousCommand()`  [INFERRED]
  tools/exec.go → agent/prompt.go
- `serviceBridgeRequests()` --calls--> `ExecuteTool()`  [INFERRED]
  tools/exec_code.go → agent/prompt.go
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go

## Import Cycles
- None detected.

## Communities (57 total, 2 thin omitted)

### Community 0 - "TruncateStr"
Cohesion: 0.06
Nodes (90): context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic(), CallAnthropicWithTools() (+82 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.06
Nodes (67): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), HomeDir(), ProjectDir() (+59 more)

### Community 2 - "testing.T"
Cohesion: 0.10
Nodes (25): TestBase64Encode(), TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg() (+17 more)

### Community 3 - "getSession"
Cohesion: 0.18
Nodes (22): agentAutoStop(), appendSessionHistory(), EnterAgentMode(), ExitAgentMode(), flushPendingMessages(), getOrCreateSession(), getSession(), getSessionHistory() (+14 more)

### Community 4 - "ConfigManager"
Cohesion: 0.14
Nodes (15): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts() (+7 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (47): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, buildArgDefsFromInputSchema(), executeMCPServerTool() (+39 more)

### Community 6 - "compaction_test.go"
Cohesion: 0.19
Nodes (22): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+14 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (52): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+44 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "runSubagent"
Cohesion: 0.06
Nodes (56): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+48 more)

### Community 11 - "uptime.go"
Cohesion: 0.11
Nodes (27): FormatDuration(), getUptime(), net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost() (+19 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (43): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), SaveModelConfig() (+35 more)

### Community 13 - "startCLI"
Cohesion: 0.07
Nodes (50): clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation, RegisterAutonomous() (+42 more)

### Community 14 - "time.Time"
Cohesion: 0.12
Nodes (34): time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult() (+26 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "collector_system_native.go"
Cohesion: 0.20
Nodes (18): TestCollectorNative_StartCPUSampler_DoesNotBlock(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+10 more)

### Community 17 - "ConfigMgr"
Cohesion: 0.07
Nodes (57): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+49 more)

### Community 18 - "ScorpPath"
Cohesion: 0.05
Nodes (62): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+54 more)

### Community 19 - "skills.go"
Cohesion: 0.08
Nodes (34): getSharedMemorySummary(), memoryFact, getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage, runSelfReview() (+26 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "RunAgentSessionLoop"
Cohesion: 0.17
Nodes (18): AgentMessage, confirmKeyboard(), countStepsInMessage(), AgentMessage, looksLikeContinuation(), mentionsBrowserTask(), screenshotWasTaken(), cleanToolCallTags() (+10 more)

### Community 31 - "AI Coding Agent Modernization Roadmap & Blueprint (v2.0)"
Cohesion: 0.13
Nodes (15): 10. Roadmap Eksekusi Teknis, 1. Lanskap Model AI Modern & Verifikasi Data (2026), 2. Evaluasi Arsitektur Scorp & Termagent Saat Ini, 3. Pilar Upgrade 1: Multi-Model & Fallback Engine, 4. Pilar Upgrade 2: Prompt Caching & Token Economics, 6. Pilar Upgrade 4: Standardisasi MCP (Model Context Protocol), 7. Pilar Upgrade 5: Subagent & Task Delegation System, 8. Pilar Upgrade 6: Deep Termux & Mobile Performance Optimization (+7 more)

### Community 32 - "scorp"
Cohesion: 0.17
Nodes (12): Build, Commands, Configuration, .env, License, models.json, Project layout, scorp (+4 more)

### Community 33 - "Scorp vs ZeroClaw: Competitive Analysis, Strategic Positioning, & Growth Playbook"
Cohesion: 0.20
Nodes (10): 🔍 1. Mengapa ZeroClaw Berhasil Meraih 32k+ Stars?, 🎯 2. Kelemahan Fatal ZeroClaw di Lingkungan Mobile & Potato Hardware, 🦂 3. Positioning Strategis Scorp: "The Mobile & Edge Native AI Agent", A. Narasi "Anti-Python & Anti-Docker" yang Sangat Kuat, B. Arsitektur Modular Crate (Community Flywheel Effect), C. Multi-Channel Ubiquity (Bukan Cuma Terminal CLI), 📊 Matriks Fitur Lengkap: Scorp vs ZeroClaw, 📌 Ringkasan Eksekutif (+2 more)

### Community 34 - "RegisterTool"
Cohesion: 0.08
Nodes (23): init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), echoPlugin, RegisterPlugin() (+15 more)

### Community 35 - "9. Pilar Upgrade 7: Ultra-Low-RAM Web Engine (< 5MB RAM di VPS 512MB & Termux)"
Cohesion: 0.33
Nodes (6): 💡 3 Solusi Arsitektur Low-RAM Browser:, 9. Pilar Upgrade 7: Ultra-Low-RAM Web Engine (< 5MB RAM di VPS 512MB & Termux), 🛑 Masalah Utama:, 🌐 Opsi 1 (Default): Zero-RAM HTTP Fetch + Readability Parser (~2MB RAM), ☁️ Opsi 2: Remote / Offloaded Browser API (RAM Lokal = 0 MB!), ⚙️ Opsi 3: Chromium Lokal Teroptimasi Ekstrem (~80MB–120MB RAM)

### Community 36 - "IsRateLimitError"
Cohesion: 0.40
Nodes (4): IsRateLimitError(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema(), TestIsRateLimitError()

### Community 37 - "🚀 4. Taktik Pertumbuhan Komunitas untuk Scorp (Playbook 10k+ Stars)"
Cohesion: 0.40
Nodes (5): 🚀 4. Taktik Pertumbuhan Komunitas untuk Scorp (Playbook 10k+ Stars), Taktik 1: Standardisasi Interface Plugin Go (Memancing Kontributor), Taktik 2: Aktifkan Mode Daemon Telegram 24/7, Taktik 3: Demo Visual yang Menghipnotis (GIF / Video di README), Taktik 4: Strategi Peluncuran Global

### Community 39 - "5. Pilar Upgrade 3: Surgical Diff Editing (Bukan Full Overwrite)"
Cohesion: 0.67
Nodes (3): 5. Pilar Upgrade 3: Surgical Diff Editing (Bukan Full Overwrite), 🛑 Masalah Cara Lama (Full File Rewrite):, ✂️ Solusi Modern (Chunk Replacement / Exact Match):

### Community 40 - "chat.go"
Cohesion: 0.16
Nodes (19): cleanupChatSessions(), CleanupSessionsLoop(), collectTableLines(), convertInlineMarkdown(), convertTableToList(), extractAndSaveMemory(), historyWriterLoop(), init() (+11 more)

### Community 41 - "GetStringArg"
Cohesion: 0.23
Nodes (13): init(), GetStringArg(), TruncOutput(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteShell(), ExecuteSystemInfo() (+5 more)

### Community 42 - "GetAllTools"
Cohesion: 0.22
Nodes (13): ActivateToolWithTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), TestDynamicToolTTL(), TickToolTTL(), GetAllTools(), ResetNativeToolCache() (+5 more)

### Community 43 - "init"
Cohesion: 0.21
Nodes (10): init(), GetBoolArg(), GetIntArg(), ExecuteCompose(), ExecuteHTTP(), ExecuteLog(), ExecuteReadURL(), decodeURL() (+2 more)

### Community 44 - "TestSessionManager"
Cohesion: 0.32
Nodes (12): ClearChatSession(), historyFilePath(), saveHistoryToDisk(), DeleteSession(), ListSessions(), RenameSession(), sanitizeSessionID(), SessionExists() (+4 more)

### Community 45 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 46 - "StorePendingConfirmation"
Cohesion: 0.28
Nodes (6): StorePendingConfirmation(), ExecuteSQL(), loadDBConnections(), dbConnection, ExecuteGit(), ExecuteProcess()

### Community 47 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 48 - "helpers.go"
Cohesion: 0.40
Nodes (5): GetInt64Arg(), GetStringSliceArg(), getUnclosedTags(), SplitMessage(), TruncOutputTool()

### Community 49 - "read_url.go"
Cohesion: 0.60
Nodes (5): ReadURL(), scrapeFirecrawl(), scrapeTavily(), truncateURLOutput(), tryRemoteScrape()

### Community 50 - "ExecuteTermuxAPI"
Cohesion: 0.73
Nodes (5): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification()

### Community 51 - "collector_native_test.go"
Cohesion: 0.24
Nodes (10): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+2 more)

### Community 52 - "ExecuteTodo"
Cohesion: 0.60
Nodes (5): ExecuteTodo(), formatTodoList(), formatTodoListLocked(), getStringArgFromMap(), TodoItem

### Community 61 - "HandleUploadInAgentMode"
Cohesion: 0.28
Nodes (7): agentSession, contentPart, imageURL, RunAgentLoop(), TGDocument, base64Encode(), HandleUploadInAgentMode()

## Knowledge Gaps
- **50 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+45 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 140 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `TruncateStr`, `getSession`, `client.go`, `collector_system.go`, `collector_security.go`, `TestSessionManager`, `startCLI`, `time.Time`, `TestSteeringQueue`, `collector_system_native.go`, `wizard.go`, `HandleUploadInAgentMode`?**
  _High betweenness centrality (0.129) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `TruncateStr`, `ConfigManager`, `client.go`, `rag_vector.go`, `chat.go`, `collector_system.go`, `runSubagent`, `wizard.go`, `startCLI`, `StorePendingConfirmation`, `time.Time`, `collector_system_native.go`, `ConfigMgr`, `skills.go`, `RunAgentSessionLoop`?**
  _High betweenness centrality (0.102) - this node is a cross-community bridge._
- **Why does `GetStringArg()` connect `GetStringArg` to `TruncateStr`, `testing.T`, `client.go`, `runSubagent`, `init`, `GetAllTools`, `patch.go`, `time.Time`, `StorePendingConfirmation`, `helpers.go`, `ConfigMgr`, `ScorpPath`, `skills.go`, `startCLI`, `ExecuteTermuxAPI`, `RunAgentSessionLoop`, `uptime.go`?**
  _High betweenness centrality (0.076) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _50 weakly-connected nodes found - possible documentation gaps or missing edges._