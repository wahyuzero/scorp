# Graph Report - scorp  (2026-09-06)

## Corpus Check
- 187 files · ~137,593 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1642 nodes · 3918 edges · 81 communities (68 shown, 2 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 511 edges (avg confidence: 0.85)
- Token cost: 225,276 input · 0 output

## Community Hubs (Navigation)
- Chat Providers & Tool Call Parsing
- Telegram Chat Session State
- Tool Bootstrap & Safety Confirmation
- MCP Marketplace Manifest Model
- MCP JSON-RPC Client Primitives
- HTTP Transport Pool
- Model Config & Catalog IO
- Tool Registration (Bootstrap init)
- Local RAG Engine
- ACP Subagent Bridge
- AI Transpiler Codegen & Sandbox
- Telegram Session & Document Types
- Agent Prompt & Memory Context
- Session Search & Exec Primitives
- System Metrics Collector (Native)
- Builtin Skills & Daemon Loop
- Agent Unit Tests (HTML/Prompt)
- Phase6 Tool Integration Tests
- Security Collector (Brute Force/VNC)
- Config Path Helpers
- Autonomous Loop Engine
- Scorp vs PicoClaw vs ZeroClaw Comparison
- Telegram File Manager UI
- Compaction Tests
- CLI Session & One-Shot Execution
- Network/Storage Collector
- Skill Management CRUD
- Updater & Release Assets
- TUI Input & Slash Commands
- Reporter Rendering Types
- Misc: config_manager.go
- Misc: StartGateway
- Misc: cost_router.go
- Misc: collector_docker.go
- Misc: collector_native_test.go
- Misc: MCP Marketplace Build Plan & Exe
- Misc: gateway.go
- Misc: collector_coolify.go
- Misc: Multi-Provider LLM Gateway (gate
- Misc: clarify.go
- Misc: tools/monitor.go
- Misc: todo.go
- Misc: SaveAutonomousConfig
- Misc: cli_callbacks.go
- Misc: Config
- Misc: COMPARISON_SCORP_PICOCLAW_ZEROCL
- Misc: inline.go
- Misc: confirmation.go
- Misc: continuation.go
- Misc: autonomy.go
- Misc: INCIDENT_20260905_runaway_cli_ve
- Misc: 📱 Build for Android Termux
- Misc: receipts.go
- Misc: steering.go
- Misc: cli_lock.go
- Misc: probe_test.go
- Misc: CGO + FTS5 Build Requirement
- Misc: compaction.go
- Misc: safety.go
- Misc: CLI Verify Deadloop (no timeout)
- Misc: install.sh
- Misc: LoadAutonomousConfig
- Misc: config_helper.go
- Misc: rag/paths.go
- Misc: contentPart
- Misc: config_paths_test.go
- Misc: api_commandcode_test.go
- Misc: session/paths.go
- Misc: termux_test.go
- Misc: scorp-agent

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 65 edges
2. `startCLI()` - 45 edges
3. `GetStringArg()` - 43 edges
4. `RunAgentSessionLoop()` - 41 edges
5. `StartDaemon()` - 41 edges
6. `ModelConfig` - 36 edges
7. `TruncateStr()` - 35 edges
8. `init()` - 30 edges
9. `ChatMessage` - 30 edges
10. `HandleModelCallback()` - 29 edges

## Surprising Connections (you probably didn't know these)
- `ZeroClaw Defense-in-Depth Sandboxing` --semantically_similar_to--> `3-Tier Autonomy & Security Sandbox (readonly/supervised/yolo)`  [INFERRED] [semantically similar]
  docs/COMPARISON_SCORP_PICOCLAW_ZEROCLAW.md → README.md
- `ExecuteAutonomous()` --calls--> `RunAutonomousCycle()`  [INFERRED]
  tools/autonomous.go → agent/autonomous.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `serviceBridgeRequests()` --calls--> `ExecuteTool()`  [INFERRED]
  tools/exec_code.go → agent/prompt.go
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **MCP Marketplace 5-Layer Security Stack** — docs_build_plan_mcp_marketplace_layer1_source_only_registry, docs_build_plan_mcp_marketplace_layer2_ast_prompt_scan, docs_build_plan_mcp_marketplace_layer3_hermetic_ci, docs_build_plan_mcp_marketplace_sha256_pinning, readme_outbound_secret_redactor [EXTRACTED 1.00]
- **AI Transpiler Pipeline (Probe -> Generate -> Build -> Verify)** — docs_build_plan_mcp_marketplace_transpiler, docs_build_plan_mcp_marketplace_transpiler_probe_phase, docs_build_plan_mcp_marketplace_transpiler_generate_phase, docs_build_plan_mcp_marketplace_transpiler_build_phase, docs_build_plan_mcp_marketplace_transpiler_verify_phase [EXTRACTED 1.00]
- **Tri-Option Install Convergence onto ~/.scorp/mcp.json** — docs_build_plan_mcp_marketplace_tri_option_install, docs_build_plan_mcp_marketplace_mcp_json_config, docs_build_plan_mcp_marketplace_mcp_manage, docs_build_plan_mcp_marketplace_watchdog, docs_build_plan_mcp_marketplace_tool_registry [EXTRACTED 1.00]

## Communities (81 total, 2 thin omitted)

### Community 0 - "Chat Providers & Tool Call Parsing"
Cohesion: 0.05
Nodes (94): scorpChat(), TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool (+86 more)

### Community 1 - "Telegram Chat Session State"
Cohesion: 0.05
Nodes (86): AgentMessage, appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown(), convertTableToList() (+78 more)

### Community 2 - "Tool Bootstrap & Safety Confirmation"
Cohesion: 0.05
Nodes (62): StorePendingConfirmation(), init(), init(), ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract() (+54 more)

### Community 3 - "MCP Marketplace Manifest Model"
Cohesion: 0.06
Nodes (75): Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo, RegistryIndex (+67 more)

### Community 4 - "MCP JSON-RPC Client Primitives"
Cohesion: 0.06
Nodes (58): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, sync.Mutex, TruncOutputTool() (+50 more)

### Community 5 - "HTTP Transport Pool"
Cohesion: 0.06
Nodes (40): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+32 more)

### Community 6 - "Model Config & Catalog IO"
Cohesion: 0.07
Nodes (56): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+48 more)

### Community 7 - "Tool Registration (Bootstrap init)"
Cohesion: 0.05
Nodes (42): RegisterAutonomous(), init(), init(), init(), init(), executeMCPServerTool(), unregisterMCPNativeTools(), TestCallModelWithToolsNilModel() (+34 more)

### Community 8 - "Local RAG Engine"
Cohesion: 0.06
Nodes (37): activeToolCall, hybridResult, computeTF(), InitRAG(), RagIndexAdd(), RagIndexList(), RagIndexRemove(), RagIndexSearch() (+29 more)

### Community 9 - "ACP Subagent Bridge"
Cohesion: 0.08
Nodes (39): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+31 more)

### Community 10 - "AI Transpiler Codegen & Sandbox"
Cohesion: 0.11
Nodes (33): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+25 more)

### Community 11 - "Telegram Session & Document Types"
Cohesion: 0.11
Nodes (33): agentSession, TGDocument, time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule() (+25 more)

### Community 12 - "Agent Prompt & Memory Context"
Cohesion: 0.08
Nodes (30): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+22 more)

### Community 13 - "Session Search & Exec Primitives"
Cohesion: 0.08
Nodes (25): bytes.Buffer, io.WriteCloser, os/exec.Cmd, ExecuteSessionSearch(), SessionResult, SessionResult, SearchSessions(), EscapeFTS5Query() (+17 more)

### Community 14 - "System Metrics Collector (Native)"
Cohesion: 0.13
Nodes (29): TestCollectorNative_CollectSystem_Structure(), FormatDuration(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes() (+21 more)

### Community 15 - "Builtin Skills & Daemon Loop"
Cohesion: 0.12
Nodes (28): scorp-agent/agent.TGDocument, StopMCPServerMode(), EnsureBuiltinSkills(), runCommandLoop(), StartDaemon(), AnswerCallback(), BackAndRefreshKeyboard(), baseName() (+20 more)

### Community 16 - "Agent Unit Tests (HTML/Prompt)"
Cohesion: 0.10
Nodes (25): TestMarkdownToTelegramHTMLEscaping(), TestBase64Encode(), TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg() (+17 more)

### Community 17 - "Phase6 Tool Integration Tests"
Cohesion: 0.14
Nodes (17): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession() (+9 more)

### Community 18 - "Security Collector (Brute Force/VNC)"
Cohesion: 0.17
Nodes (25): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+17 more)

### Community 19 - "Config Path Helpers"
Cohesion: 0.16
Nodes (24): BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), HistoryDirPath(), HomeDir(), Hostname() (+16 more)

### Community 20 - "Autonomous Loop Engine"
Cohesion: 0.18
Nodes (19): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), makeDecision() (+11 more)

### Community 21 - "Scorp vs PicoClaw vs ZeroClaw Comparison"
Cohesion: 0.14
Nodes (21): Scorp vs PicoClaw vs ZeroClaw Architecture Comparison, ZeroClaw Defense-in-Depth Sandboxing, Scorp Embedded Native Multi-Engine Metasearch, Scorp Local Simhash & Vector RAG (0 external DB), PicoClaw (Sipeed, Go), PicoClaw provider:auto Search Fallback, ZeroClaw (ZeroClaw Labs, Rust), ZeroClaw Third-Party Search Stack (Tavily/Brave/Jina/SearXNG) (+13 more)

### Community 22 - "Telegram File Manager UI"
Cohesion: 0.21
Nodes (19): HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath(), HumanSize() (+11 more)

### Community 23 - "Compaction Tests"
Cohesion: 0.22
Nodes (18): AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens(), TestPrune_BoundaryAges(), TestPrune_DockerScenario_41Messages(), TestPrune_EmptyHistory(), TestPrune_HeadTailFormat() (+10 more)

### Community 24 - "CLI Session & One-Shot Execution"
Cohesion: 0.24
Nodes (18): executeOneShot(), executeTurn(), handleCLISession(), handleCLISOP(), hasDebugFlag(), handleMCPCommand(), printBanner(), printCLIHelp() (+10 more)

### Community 25 - "Network/Storage Collector"
Cohesion: 0.23
Nodes (18): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+10 more)

### Community 26 - "Skill Management CRUD"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 27 - "Updater & Release Assets"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 28 - "TUI Input & Slash Commands"
Cohesion: 0.21
Nodes (12): GetHistoryTokenEstimate(), GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste() (+4 more)

### Community 29 - "Reporter Rendering Types"
Cohesion: 0.27
Nodes (14): NetworkData, SystemData, scorp-agent/collectors.DockerData, PortInfo, Bar(), bar(), FormatHourlyReport(), FormatStatusResponse() (+6 more)

### Community 30 - "Misc: config_manager.go"
Cohesion: 0.23
Nodes (4): CM(), InitConfigManager(), NewConfigManager(), ConfigManager

### Community 31 - "Misc: StartGateway"
Cohesion: 0.29
Nodes (11): StartGateway(), isCLIMode(), main(), SOP, Dir(), GetSOP(), InitDefaultSOPs(), ListSOPs() (+3 more)

### Community 32 - "Misc: cost_router.go"
Cohesion: 0.25
Nodes (13): defaultCostConfig(), formatCostReport(), FormatDailyCostSummary(), handleCostCommand(), init(), isBudgetExceeded(), isOffPeak(), LoadCostConfig() (+5 more)

### Community 33 - "Misc: collector_docker.go"
Cohesion: 0.27
Nodes (11): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+3 more)

### Community 34 - "Misc: collector_native_test.go"
Cohesion: 0.21
Nodes (12): TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single(), TestCollectorNative_StartCPUSampler_DoesNotBlock() (+4 more)

### Community 35 - "Misc: MCP Marketplace Build Plan & Exe"
Cohesion: 0.21
Nodes (13): MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry, Security Layer 2: AST & Prompt Injection Scan, Security Layer 3: Hermetic CI/CD, ~/.scorp/mcp.json Central Config Entry, mcp_manage Agent Tool (mcp/manage.go), Scorp MCP Marketplace, MCP Marketplace Blueprint (docs/MCP_MARKETPLACE_BLUEPRINT.md) (+5 more)

### Community 36 - "Misc: gateway.go"
Cohesion: 0.35
Nodes (11): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), TestGatewayEndpoints() (+3 more)

### Community 37 - "Misc: collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 38 - "Misc: Multi-Provider LLM Gateway (gate"
Cohesion: 0.24
Nodes (12): Multi-Provider LLM Gateway (gateway/gateway.go + models.json), mark3labs/mcp-go SDK, mcp-fetch Reference Port (web scraping/markdown), mcp-filesystem Reference Port (scoped file access), mcp-sqlite Reference Port (local DB inspection), wahyuzero/scorp-mcp-registry (public registry repo), AI Transpiler (Self-Hosted Codegen), Transpiler Phase 3a: Sandbox Build (+4 more)

### Community 39 - "Misc: clarify.go"
Cohesion: 0.27
Nodes (9): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID() (+1 more)

### Community 40 - "Misc: tools/monitor.go"
Cohesion: 0.38
Nodes (10): ExecuteMonitor(), InitMonitor(), loadMonitorTargets(), monitorCheckOne(), monitorLoop(), ragIngestText(), sanitizeFilename(), saveMonitorTargets() (+2 more)

### Community 41 - "Misc: todo.go"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 42 - "Misc: SaveAutonomousConfig"
Cohesion: 0.33
Nodes (9): SaveAutonomousConfig(), saveAutonomousConfigLocked(), SetKillSwitch(), TestPhase7_KillSwitch(), autoShowActions(), autoShowConfig(), autoShowLog(), autoStatus() (+1 more)

### Community 43 - "Misc: cli_callbacks.go"
Cohesion: 0.31
Nodes (7): wireCLICallbacks(), formatFinalResponse(), formatTerminalText(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 44 - "Misc: Config"
Cohesion: 0.33
Nodes (9): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), Init(), StartServer() (+1 more)

### Community 45 - "Misc: COMPARISON_SCORP_PICOCLAW_ZEROCL"
Cohesion: 0.20
Nodes (8): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch)

### Community 46 - "Misc: inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 47 - "Misc: confirmation.go"
Cohesion: 0.39
Nodes (8): clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation, printStatus()

### Community 48 - "Misc: continuation.go"
Cohesion: 0.29
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 49 - "Misc: autonomy.go"
Cohesion: 0.43
Nodes (6): GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), SetAutonomyLevel(), TestAutonomyLevels(), AutonomyLevel

### Community 50 - "Misc: INCIDENT_20260905_runaway_cli_ve"
Cohesion: 0.25
Nodes (7): Akar Masalah (dugaan), Ciri-ciri Runaway yang Terdeteksi, Incident Report: Runaway CLI Verify Session, Pelajaran Umum, Penanganan, Rekomendasi Perbaikan & Status Implementasi, Ringkasan

### Community 51 - "Misc: 📱 Build for Android Termux"
Cohesion: 0.25
Nodes (8): 📱 Build for Android Termux, Launch Embedded Web Dashboard (< 2MB RAM):, Launch Interactive CLI Chat:, 📄 License, 🚀 Quick Start (Under 60 Seconds), Run One-Shot Tasks:, 🦂 Scorp, ⚡ What Makes Scorp v2.0 Special

### Community 52 - "Misc: receipts.go"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 53 - "Misc: steering.go"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 54 - "Misc: cli_lock.go"
Cohesion: 0.33
Nodes (4): acquireSessionLock(), TestAcquireSessionLock(), os.File, sessionLockFile

### Community 55 - "Misc: probe_test.go"
Cohesion: 0.38
Nodes (6): contains(), goldenFilesystemBinary(), TestDetectRuntimeMessages(), TestProbeAndVerifyGoldenRoundTrip(), TestProbeUpstreamNpx(), TestStripCodeFences()

### Community 56 - "Misc: CGO + FTS5 Build Requirement"
Cohesion: 0.40
Nodes (6): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection

### Community 57 - "Misc: compaction.go"
Cohesion: 0.53
Nodes (5): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, TestEstimateTokens()

### Community 58 - "Misc: safety.go"
Cohesion: 0.47
Nodes (4): devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), IsDangerousCommand(), isHarmlessDevSink()

### Community 60 - "Misc: CLI Verify Deadloop (no timeout)"
Cohesion: 0.60
Nodes (6): CLI Verify Deadloop (no timeout), complete_task Explicit Tool Contract, Incident Report: Runaway CLI Verify Session (2026-09-05), Session Lock (cli_lock.go kernel advisory file lock), Heartbeat / Stall Detection, Agent Turn Timeout (context.WithTimeout, SCORP_MAX_TURN_TIMEOUT)

### Community 61 - "Misc: install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 62 - "Misc: LoadAutonomousConfig"
Cohesion: 0.60
Nodes (5): LoadAutonomousConfig(), setupTestPaths(), TestPhase7_ConfigPersistence(), ConfigMgr(), saveCostTracker()

### Community 63 - "Misc: config_helper.go"
Cohesion: 0.50
Nodes (3): SaveJSON(), SaveJSONPerm(), os.FileMode

### Community 65 - "Misc: rag/paths.go"
Cohesion: 0.70
Nodes (4): homeDir(), ragDirPath(), ragIndexPath(), ragVectorDBPath()

### Community 66 - "Misc: contentPart"
Cohesion: 0.67
Nodes (3): contentPart, imageURL, base64Encode()

### Community 67 - "Misc: config_paths_test.go"
Cohesion: 0.50
Nodes (3): TestConfigPaths_RagVectorDBPath(), TestConfigPaths_ScorpDir(), TestConfigPaths_ScreenshotsDir()

### Community 68 - "Misc: api_commandcode_test.go"
Cohesion: 0.50
Nodes (3): TestBuildCommandCodePayload(), TestExtractFallbackToolCalls(), TestSwitchActiveModel()

### Community 70 - "Misc: session/paths.go"
Cohesion: 0.83
Nodes (3): homeDir(), scorpDir(), scorpPath()

## Knowledge Gaps
- **34 isolated node(s):** `Akar Masalah (dugaan)`, `Ciri-ciri Runaway yang Terdeteksi`, `Pelajaran Umum`, `Penanganan`, `Rekomendasi Perbaikan & Status Implementasi` (+29 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 175 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `Telegram File Manager UI` to `Chat Providers & Tool Call Parsing`, `Telegram Chat Session State`, `MCP Marketplace Manifest Model`, `MCP JSON-RPC Client Primitives`, `Misc: collector_coolify.go`, `Model Config & Catalog IO`, `Misc: clarify.go`, `Telegram Session & Document Types`, `Agent Prompt & Memory Context`, `System Metrics Collector (Native)`, `Misc: confirmation.go`, `Builtin Skills & Daemon Loop`, `Security Collector (Brute Force/VNC)`, `Misc: steering.go`, `Network/Storage Collector`, `Reporter Rendering Types`, `Misc: StartGateway`?**
  _High betweenness centrality (0.180) - this node is a cross-community bridge._
- **Why does `GetStringArg()` connect `Tool Bootstrap & Safety Confirmation` to `Misc: cost_router.go`, `Telegram Chat Session State`, `MCP JSON-RPC Client Primitives`, `Model Config & Catalog IO`, `Tool Registration (Bootstrap init)`, `Misc: tools/monitor.go`, `ACP Subagent Bridge`, `Telegram Session & Document Types`, `Agent Prompt & Memory Context`, `System Metrics Collector (Native)`, `Agent Unit Tests (HTML/Prompt)`, `Phase6 Tool Integration Tests`, `Misc: StartGateway`?**
  _High betweenness centrality (0.110) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `Telegram Chat Session State` to `Chat Providers & Tool Call Parsing`, `Tool Bootstrap & Safety Confirmation`, `Tool Registration (Bootstrap init)`, `Local RAG Engine`, `Telegram Session & Document Types`, `Agent Prompt & Memory Context`, `Misc: continuation.go`, `Misc: autonomy.go`, `Misc: steering.go`, `Telegram File Manager UI`, `CLI Session & One-Shot Execution`, `Misc: safety.go`?**
  _High betweenness centrality (0.095) - this node is a cross-community bridge._
- **Are the 22 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 22 INFERRED edges - model-reasoned connections that need verification._
- **Are the 7 inferred relationships involving `startCLI()` (e.g. with `wireCLICallbacks()` and `formatTerminalText()`) actually correct?**
  _`startCLI()` has 7 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Akar Masalah (dugaan)`, `Ciri-ciri Runaway yang Terdeteksi`, `Pelajaran Umum` to the rest of the system?**
  _34 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Chat Providers & Tool Call Parsing` be split into smaller, more focused modules?**
  _Cohesion score 0.051574212893553226 - nodes in this community are weakly interconnected._