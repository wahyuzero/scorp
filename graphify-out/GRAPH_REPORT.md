# Graph Report - scorp  (2026-09-09)

## Corpus Check
- 288 files · ~208,698 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2280 nodes · 5832 edges · 115 communities (103 shown, 4 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 918 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f9e5a28c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- chat.go
- cost_router.go
- runSubagent
- metasearch_engines.go
- HomeDir
- rag_vector.go
- TaskPlan
- Scorp Agent (Go, ultra-light autonomous agent)
- LoadMCPConfig
- Benchmark
- StartDaemon
- getAgentSystemPrompt
- time.Time
- 🔍 Rincian Detail Uji per Fase
- HandleTelegramAction
- testing.T
- ScorpPath
- RegisterTool
- RunAgentSessionLoop
- session_search_fts5.go
- compaction_test.go
- collector_security.go
- context.Context
- GetAutonomyLevel
- collector_system.go
- RunPreToolUseHooks
- MCPServer
- api_commandcode.go
- CreateCheckpoint
- collector_system_native.go
- ToolCall
- TestIntegrityStatus
- CheckServerContracts
- client.go
- skills.go
- CallOpenCodeWithTools
- getSession
- PermissionDecision
- FormatHourlyReport
- checker.go
- confirmation.go
- SCORP — BRUTAL END-TO-END TEST PLAN
- ExecuteShell
- api_gemini.go
- SandboxActive
- maybeCompactHistory
- CallModel
- 🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering)
- GetStringArg
- CheckDenyRules
- eval.go
- session_ui.go
- startCLI
- TruncateStr
- patchReplace
- HandleUploadInAgentMode
- collector_docker.go
- RegisterProvider
- ResetNativeToolCache
- prepareNewTurnHistory
- 🔍 Temuan Kelemahan & Solusi Arsitektural Permanen
- TodoManager
- 🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam
- 🔍 Temuan Kelemahan & Solusi Arsitektural Permanen
- ExecuteReadURL
- testgate.go
- HandleModelCallback
- IsDangerousCommand
- os.File
- runLiveCase
- CredentialVault
- 🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks)
- IsContinuationDirective
- TestPhase6_AllTools
- .listenSSEStream
- install.sh
- deploy.sh
- collector_coolify.go
- rag_vector_test.go
- ExecuteTermuxAPI
- GenerateContextualSessionTitle
- 🔍 Temuan 10 Kelemahan Arsitektural & Solusi Permanen
- saveConfig
- ExecuteScript
- 🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark)
- 🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut
- registerMCPToolsAsNative
- ExecuteAutoLogin
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- Transpiler Phase 3b: Contract Verify & Graceful Degradation
- readInteractiveInput
- MCP Marketplace Build Plan & Execution Roadmap
- GetAllTools
- ExecuteTool
- LoadConfig
- inline.go
- TestRegisterPlugin
- TopProcess
- TransformToolDefinitions
- serviceBridgeRequests
- Laporan Eksekusi Komprehensif Uji End-to-End (E2E): Fase 1 – 15
- clarify.go
- TruncOutput
- Release Workflow (tag-triggered cross-compile + GitHub release)
- CLI Verify Deadloop (no timeout)

## God Nodes (most connected - your core abstractions)
1. `ModelConfig` - 89 edges
2. `HandleTelegramAction()` - 80 edges
3. `ChatMessage` - 75 edges
4. `RunAgentSessionLoop()` - 64 edges
5. `startCLI()` - 52 edges
6. `TruncateStr()` - 47 edges
7. `GetStringArg()` - 47 edges
8. `ScorpPath()` - 44 edges
9. `StartDaemon()` - 44 edges
10. `resumeAgentLoop()` - 40 edges

## Surprising Connections (you probably didn't know these)
- `Durable Memory MEMORY.md (P1.7)` --implements--> `extractTaskMemory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/memory_md.go
- `Auto-Mode Classifier (P3.13)` --implements--> `PermissionDecision()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/auto.go
- `Compaction Preservation (P2.10)` --implements--> `truncateToolResultsInHistory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/compaction.go
- `Plan Mode Workflow (P1.4)` --implements--> `BeginPlanning()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/planmode.go
- `Persistent Task Ledger (P1.5)` --implements--> `savePlanToDisk()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/taskplan.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **MCP Marketplace 5-Layer Security Stack** — docs_build_plan_mcp_marketplace_layer1_source_only_registry, docs_build_plan_mcp_marketplace_layer2_ast_prompt_scan, docs_build_plan_mcp_marketplace_layer3_hermetic_ci, docs_build_plan_mcp_marketplace_sha256_pinning, readme_outbound_secret_redactor [EXTRACTED 1.00]
- **Scorp Resilient State & Memory Workflow (PlanMode + TaskLedger + Checkpoint + MEMORY.md)** — docs_implementation_plan_scorp_plan_mode, docs_implementation_plan_scorp_task_ledger_persistence, docs_implementation_plan_scorp_checkpoint_rewind, docs_implementation_plan_scorp_durable_memory [EXTRACTED 1.00]
- **Scorp Trust & Safety Gate Stack (Sandbox + DenyRules + Hooks + AutoClassifier)** — docs_implementation_plan_scorp_sandbox_bwrap, docs_implementation_plan_scorp_deny_rule_engine, docs_implementation_plan_scorp_hooks_lifecycle, docs_implementation_plan_scorp_auto_mode_classifier [EXTRACTED 1.00]
- **Scorp Verification & Integrity Pipeline (TestGate + ClaimGate + EvalArena)** — docs_implementation_plan_scorp_test_integrity_gate, docs_implementation_plan_scorp_claim_gate, docs_implementation_plan_scorp_eval_arena [EXTRACTED 1.00]
- **AI Transpiler Pipeline (Probe -> Generate -> Build -> Verify)** — docs_build_plan_mcp_marketplace_transpiler, docs_build_plan_mcp_marketplace_transpiler_probe_phase, docs_build_plan_mcp_marketplace_transpiler_generate_phase, docs_build_plan_mcp_marketplace_transpiler_build_phase, docs_build_plan_mcp_marketplace_transpiler_verify_phase [EXTRACTED 1.00]
- **Tri-Option Install Convergence onto ~/.scorp/mcp.json** — docs_build_plan_mcp_marketplace_tri_option_install, docs_build_plan_mcp_marketplace_mcp_json_config, docs_build_plan_mcp_marketplace_mcp_manage, docs_build_plan_mcp_marketplace_watchdog, docs_build_plan_mcp_marketplace_tool_registry [EXTRACTED 1.00]

## Communities (115 total, 4 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (72): handleMCPCommand(), Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo (+64 more)

### Community 1 - "chat.go"
Cohesion: 0.11
Nodes (25): ClearStopRequest(), collectTableLines(), convertInlineMarkdown(), convertTableToList(), ExitAgentMode(), extractAndSaveMemory(), flushPendingMessages(), GetHistoryTokenEstimate() (+17 more)

### Community 2 - "cost_router.go"
Cohesion: 0.07
Nodes (56): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+48 more)

### Community 3 - "runSubagent"
Cohesion: 0.05
Nodes (59): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+51 more)

### Community 4 - "metasearch_engines.go"
Cohesion: 0.05
Nodes (50): net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost(), getClient(), getShortClient() (+42 more)

### Community 5 - "HomeDir"
Cohesion: 0.25
Nodes (17): HomeDir(), ProjectDir(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), HumanSize() (+9 more)

### Community 6 - "rag_vector.go"
Cohesion: 0.07
Nodes (38): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+30 more)

### Community 7 - "TaskPlan"
Cohesion: 0.09
Nodes (42): PlanItem, ApprovePlan(), BeginPlanning(), CancelPlan(), EndPlanning(), AgentMessage, PlanningState(), RevisePlan() (+34 more)

### Community 8 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.16
Nodes (19): Scorp vs PicoClaw vs ZeroClaw Architecture Comparison, ZeroClaw Defense-in-Depth Sandboxing, Scorp Embedded Native Multi-Engine Metasearch, Scorp Local Simhash & Vector RAG (0 external DB), ZeroClaw (ZeroClaw Labs, Rust), ZeroClaw Third-Party Search Stack (Tavily/Brave/Jina/SearXNG), CLI Slash Commands (/help /models /model /mode /session /sop /receipts /tools /cost /clear /exit), Cryptographic SHA-256 Tool Execution Receipts (+11 more)

### Community 9 - "LoadMCPConfig"
Cohesion: 0.38
Nodes (11): MCPConfigFilePath(), LoadMCPConfig(), ReloadMCPServers(), sanitizeMCPName(), AddServerEntry(), ExecuteMCPManage(), mcpManageAdd(), mcpManageList() (+3 more)

### Community 10 - "Benchmark"
Cohesion: 0.09
Nodes (39): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+31 more)

### Community 11 - "StartDaemon"
Cohesion: 0.14
Nodes (25): UploadsDir(), runCommandLoop(), StartDaemon(), BackAndRefreshKeyboard(), baseName(), DeleteWebhook(), EditMessage(), EditMessageByID() (+17 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.07
Nodes (38): getSharedMemorySummary(), AppendMemoryMD(), extractTaskMemory(), AgentMessage, parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile() (+30 more)

### Community 13 - "time.Time"
Cohesion: 0.08
Nodes (40): presentPlanTelegram(), CM(), InitConfigManager(), NewConfigManager(), ConfigManager, time.Time, EscapeHTML(), ScheduledTask (+32 more)

### Community 14 - "🔍 Rincian Detail Uji per Fase"
Cohesion: 0.10
Nodes (19): 1️⃣1️⃣ Fase 11: Recursive Parent Lock & Deep Traversal, 1️⃣2️⃣ Fase 12: Synthetic Feedback & Unicode RTL/Emoji Edge Cases, 1️⃣3️⃣ Fase 13: Output Buffer Flooding & EPIPE Signals, 1️⃣4️⃣ Fase 14: Session Environment Persistence & Error Reporting, 1️⃣5️⃣ Fase 15: Subshell Mutation Isolation & Polyglot Compilation, 1️⃣ Fase 1: Tool Parsing & Thought Signature Contract, 2️⃣ Fase 2: Filesystem Whitelist & Bare-Metal Sandbox Toggle, 3️⃣ Fase 3: Subprocess Pipe Deadlock & Asynchronous Draining (+11 more)

### Community 15 - "HandleTelegramAction"
Cohesion: 0.28
Nodes (12): HandleTelegramAction(), GetPath(), handleTelegramMarketplaceInstall(), MarketplaceTriOptionKeyboard(), reportInstallOutcome(), BackButtonKeyboard(), MainMenuKeyboard(), MonitorMenuKeyboard() (+4 more)

### Community 16 - "testing.T"
Cohesion: 0.06
Nodes (44): TestBuildThinkingMessage(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations(), TestParseToolCalls() (+36 more)

### Community 17 - "ScorpPath"
Cohesion: 0.12
Nodes (23): init(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), Hostname(), MemoryFilePath() (+15 more)

### Community 18 - "RegisterTool"
Cohesion: 0.11
Nodes (17): TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), ResetAutoStats(), init(), init(), init(), init(), GenerateSystemPromptDescriptions() (+9 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.15
Nodes (25): AgentMessage, confirmationDisplay(), ConsumeStopRequest(), setLoopActive(), confirmKeyboard(), cleanToolCallTags(), getSessionSearchContext(), maxIterations() (+17 more)

### Community 20 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 21 - "compaction_test.go"
Cohesion: 0.19
Nodes (20): AgentMessage, makeHistory(), makeToolResult(), mkCompactionHistory(), TestEstimateHistoryTokens(), TestPrune_BoundaryAges(), TestPrune_DockerScenario_41Messages(), TestPrune_EmptyHistory() (+12 more)

### Community 22 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 23 - "context.Context"
Cohesion: 0.10
Nodes (18): context.Context, AnthropicProvider, applyOpenAIHeaders(), buildOpenAIRequestBody(), CallOpenAI(), CallOpenAIWithTools(), formatOpenAIMessages(), AzureProvider (+10 more)

### Community 24 - "GetAutonomyLevel"
Cohesion: 0.14
Nodes (24): ConfirmationRequired(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels() (+16 more)

### Community 25 - "collector_system.go"
Cohesion: 0.23
Nodes (18): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+10 more)

### Community 26 - "RunPreToolUseHooks"
Cohesion: 0.13
Nodes (33): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+25 more)

### Community 27 - "MCPServer"
Cohesion: 0.14
Nodes (14): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, sync.Once, FindMCPTool(), MCPServer, MCPTool, ProbeServer() (+6 more)

### Community 28 - "api_commandcode.go"
Cohesion: 0.14
Nodes (18): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeStream(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), resolveCommandCodeKeyFromDisk(), ChatRequest (+10 more)

### Community 29 - "CreateCheckpoint"
Cohesion: 0.26
Nodes (19): Checkpoint and Rewind (P1.6), caseCheckpoint(), CheckpointDiffStat(), checkpointRepoRoot(), ckptGit(), ckptRefFor(), CreateCheckpoint(), DeleteCheckpoint() (+11 more)

### Community 30 - "collector_system_native.go"
Cohesion: 0.21
Nodes (18): FormatDuration(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+10 more)

### Community 31 - "ToolCall"
Cohesion: 0.21
Nodes (17): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), ToolCallResp, isVisionModelName() (+9 more)

### Community 32 - "TestIntegrityStatus"
Cohesion: 0.29
Nodes (16): Test-Integrity Gate (P0.3), IsTestRelatedPath(), IsTestRunCommand(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestIsTestRunCommand() (+8 more)

### Community 33 - "CheckServerContracts"
Cohesion: 0.20
Nodes (16): MCP Contract Watch (P3.14), Agent Failure Modes and Critiques (2026), caseMCPContract(), CheckServerContracts(), contractFile(), ContractStatusNotice(), ContractWarnings(), serverFingerprint() (+8 more)

### Community 34 - "client.go"
Cohesion: 0.15
Nodes (21): ACPRequest, encoding/json.RawMessage, executeMCPServerTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError() (+13 more)

### Community 35 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 36 - "CallOpenCodeWithTools"
Cohesion: 0.21
Nodes (10): CallOpenCode(), CallOpenCodeStream(), CallOpenCodeWithTools(), resolveOpenCodeBaseURL(), resolveOpenCodeKeyFromDisk(), resolveOpenCodeSessionID(), RecordCostWithCache(), CostTracker (+2 more)

### Community 37 - "getSession"
Cohesion: 0.19
Nodes (25): appendSessionHistory(), ClearChatSession(), EnterAgentMode(), FlushSessionHistory(), getOrCreateSession(), getSession(), getSessionHistory(), getSessionMap() (+17 more)

### Community 38 - "PermissionDecision"
Cohesion: 0.21
Nodes (18): TestAutoAllowlistPrefixAnchoring(), autoAllowlisted(), autoClassify(), autoClassifyWithModel(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand(), PermissionDecision() (+10 more)

### Community 39 - "FormatHourlyReport"
Cohesion: 0.27
Nodes (13): NetworkData, SystemData, PortInfo, Bar(), bar(), FormatHourlyReport(), FormatStatusResponse(), SectionCoolify() (+5 more)

### Community 40 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 41 - "confirmation.go"
Cohesion: 0.26
Nodes (15): TestStorePendingConfirmationArgs(), clearPendingConfirmation(), clearPendingConfirmationFromDisk(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation() (+7 more)

### Community 42 - "SCORP — BRUTAL END-TO-END TEST PLAN"
Cohesion: 0.09
Nodes (20): 0. SCALE & SEVERITY DEFINITIONS, A. GATE-STACK ADVERSARIAL (S) — Layer Penetration Probes, B. LONG-HORIZON & CONTEXT COMPACTION (L), C. INFRASTRUCTURE CHAOS (M/L), D. CONCURRENCY & RACE CONDITIONS (M), E. TELEGRAM UX INTEGRITY (M), F. ADVERSARIAL SECURITY (M), G. EVALUATION & DEPLOYMENT GATE INTEGRITY (S) (+12 more)

### Community 43 - "ExecuteShell"
Cohesion: 0.13
Nodes (22): init(), Sandbox Shell Bubblewrap (P0.1), WriteFileAtomic(), captureSessionExports(), cleanAbortedShm(), decodeUTF16(), decodeUTFString(), ExecuteListDir() (+14 more)

### Community 44 - "api_gemini.go"
Cohesion: 0.17
Nodes (19): callGemini(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), parseDataURL(), resolveGeminiBaseURL(), geminiContent (+11 more)

### Community 45 - "SandboxActive"
Cohesion: 0.25
Nodes (15): caseSandbox(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest(), SandboxStatusNotice(), SandboxVersion() (+7 more)

### Community 46 - "maybeCompactHistory"
Cohesion: 0.25
Nodes (17): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, injectPreservationNote(), latestUserGoal(), maybeCompactHistory(), preservationNote() (+9 more)

### Community 47 - "CallModel"
Cohesion: 0.16
Nodes (16): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestGemini_MultimodalParsing(), TestOpenCodeProvider_RegistrationAndKey(), CallModel(), CallModelStream(), GetProvider() (+8 more)

### Community 48 - "🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering)"
Cohesion: 0.18
Nodes (10): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Brutal, 🔍 3. Analisis Mendalam Kualitas Output & Arsitektur, 🏆 4. Kesimpulan Akhir Benchmark Brutal, A. Kompilasi Bahasa Go Multi-File (Kasus B2), B. Otomasi Server REST API Latar Belakang & Health Check (Kasus B7), 🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering), C. Fenomena Kebocoran Tag Simulasi pada ZeroClaw (Kasus B1, B3, B5, B8, B10, B12, B15) (+2 more)

### Community 49 - "GetStringArg"
Cohesion: 0.09
Nodes (23): TestGetBoolArg(), init(), PythonSitePackages(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage() (+15 more)

### Community 50 - "CheckDenyRules"
Cohesion: 0.20
Nodes (15): CheckDenyRules(), loadDenyRules(), ParseDenyRule(), ReloadDenyRules(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped(), TestCheckDenyRulesMatches() (+7 more)

### Community 51 - "eval.go"
Cohesion: 0.22
Nodes (12): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, Case, caseResult, CoreCases(), evalSandboxDir(), humanCount(), liveLabel() (+4 more)

### Community 52 - "session_ui.go"
Cohesion: 0.31
Nodes (11): FormatCompactStats(), CompactStats, BuildSessionMenuKeyboard(), FormatSessionMenuText(), GetActiveSessionID(), HandleSessionCallback(), init(), loadTgSessionMapping() (+3 more)

### Community 53 - "startCLI"
Cohesion: 0.05
Nodes (65): RequestStop(), RegisterAutonomous(), formatFileSize(), wireCLICallbacks(), executeOneShot(), executeTurn(), formatFinalResponse(), formatTerminalText() (+57 more)

### Community 54 - "TruncateStr"
Cohesion: 0.21
Nodes (26): TruncateStr(), anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic(), callAnthropicStream() (+18 more)

### Community 55 - "patchReplace"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "HandleUploadInAgentMode"
Cohesion: 0.20
Nodes (10): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), contentPart, imageURL, TestBase64Encode(), cleanupAgentSessions(), TGDocument (+2 more)

### Community 57 - "collector_docker.go"
Cohesion: 0.27
Nodes (11): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+3 more)

### Community 58 - "RegisterProvider"
Cohesion: 0.09
Nodes (22): init(), init(), formatMessagesForCLI(), init(), init(), init(), init(), init() (+14 more)

### Community 59 - "ResetNativeToolCache"
Cohesion: 0.20
Nodes (19): ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), clearDynamicEnv(), registerTempTool(), TestDynamicToolTTL() (+11 more)

### Community 60 - "prepareNewTurnHistory"
Cohesion: 0.29
Nodes (8): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary()

### Community 61 - "🔍 Temuan Kelemahan & Solusi Arsitektural Permanen"
Cohesion: 0.18
Nodes (10): 1. Unbounded Concurrent Execution pada Scheduler (`scheduler/scheduler.go`), 2. Lazy Initialization Panic pada Vector RAG (`rag/rag_vector.go`), 3. Non-Atomic RAG Index Disk Persistence (`rag/rag.go` & `rag/rag_vector.go`), 4. SimHash Zero-Variance Whitespace Hallucination Pollution (`rag/rag_vector.go`), 5. Node Modules & Git Trash Folder Ingestion di RAG (`rag/rag_vector.go`), 6. Task Plan Atomic Serialization Under High Rate (`agent/taskplan.go`), Laporan Audit Lanjutan Flaw Hunting & Penguatan Sistem (Fase 21), 📊 Ringkasan Hasil Pengujian Fase 21 pada VPS (+2 more)

### Community 62 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 63 - "🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam"
Cohesion: 0.25
Nodes (7): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Ultra-Hardcore, 🔍 3. Temuan Kritis: Mengapa ZeroClaw Runtuh Total?, 💡 4. Analisis Komparasi Mendalam: Scorp vs PicoClaw, 🏆 5. Rekapitulasi Menyeluruh (Grand Total 65 Skenario Uji), 🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam, Scorp Agent vs PicoClaw vs ZeroClaw (Live Head-to-Head VPS Evaluation)

### Community 64 - "🔍 Temuan Kelemahan & Solusi Arsitektural Permanen"
Cohesion: 0.18
Nodes (10): 1. Search Code Flag Injection & Catastrophic Traversal (`tools/search.go`), 2. Symlink Traversal Whitelist Escape Bypass (`tools/exec.go`), 3. Persistent Memory Concurrent Race & Truncation Hazard (`tools/memory.go`), 4. Unbounded Real-Time Steering Queue Flooding (`agent/steering.go`), 5. Telegram File Browser Path Mapping Memory Leak (`telegram/files.go`), 6. Ephemeral Pending Confirmation Loss on Restart (`agent/confirmation.go`), Laporan Audit Lanjutan Flaw Hunting & Penguatan Sistem (Fase 20), 📊 Ringkasan Hasil Pengujian Fase 20 pada VPS (+2 more)

### Community 65 - "ExecuteReadURL"
Cohesion: 0.36
Nodes (7): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), TestReadURL_LocalMock(), truncateURLOutput(), tryRemoteScrape()

### Community 66 - "testgate.go"
Cohesion: 0.12
Nodes (24): OperationalClaim, GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), RedactSecrets(), TestRedactSecrets() (+16 more)

### Community 67 - "HandleModelCallback"
Cohesion: 0.08
Nodes (56): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+48 more)

### Community 68 - "IsDangerousCommand"
Cohesion: 0.33
Nodes (7): TestIsDangerousCommand(), devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), inspectBase64Payload(), IsDangerousCommand(), isHarmlessDevSink()

### Community 69 - "os.File"
Cohesion: 0.21
Nodes (13): lockFileExclusive(), unlockFile(), lockFileExclusive(), unlockFile(), os.File, FileLockExclusive(), fileLockExclusive(), FileUnlock() (+5 more)

### Community 70 - "runLiveCase"
Cohesion: 0.31
Nodes (9): deploymentEnv(), LiveCases(), runLiveCase(), tail(), usageDelta(), usageSnapshot(), liveCase, TestUsageDeltaClampsNegative() (+1 more)

### Community 71 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 72 - "🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks)"
Cohesion: 0.29
Nodes (6): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Agentic Berat, 🔍 3. Kesimpulan Utama Uji Komparasi Baru, 🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks), ⚙️ Lingkungan & Konfigurasi Pengujian, Scorp Agent vs PicoClaw vs ZeroClaw (Fresh API Keys & Fast-Path Heuristic)

### Community 73 - "IsContinuationDirective"
Cohesion: 0.40
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 75 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 76 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 78 - "deploy.sh"
Cohesion: 0.53
Nodes (4): die(), ok(), deploy.sh script, step()

### Community 79 - "collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 80 - "rag_vector_test.go"
Cohesion: 0.29
Nodes (7): hammingDistance(), simhashSimilarity(), TestSimHash_ComputeSimhash(), TestSimHash_HammingDistance(), TestSimHash_Similarity(), TestSimHash_SmartChunk(), TestSimHash_Tokenize()

### Community 81 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 82 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 83 - "🔍 Temuan 10 Kelemahan Arsitektural & Solusi Permanen"
Cohesion: 0.17
Nodes (11): 1. In-Memory Session Environment Non-Persistence across Independent CLI Invocations, 2. Truncation Hazard pada `write_file`, `replace_file_content`, & Bridge Protocol, 3. CRLF (`\r\n`) Line-Ending Mismatch pada Surgical Chunk Replacement, 4. 64-Bit Cryptographic Receipt Hash Collision Blindspot, 5. Infinite Traversal Hang pada Symlink Cycles & Inode Tracking, 6. Binary Non-UTF8 Stream Corruption di Shell Output Buffer, 7. Over-Inspection Loop pada Skenario Modifikasi File, Laporan Audit Kelemahan Arsitektural Non-Kecepatan & Penguatan Sistem (Fase 16 – 18) (+3 more)

### Community 84 - "saveConfig"
Cohesion: 0.43
Nodes (6): apiKey, envKeyName, installSystemdService(), maskString(), RunQuickstart(), saveConfig()

### Community 85 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 86 - "🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark)"
Cohesion: 0.29
Nodes (6): 📊 1. Ringkasan Hasil Eksekusi, ⏱️ 2. Tabel Waktu Respon Per Tugas (Detik), 🔍 3. Analisis Peningkatan Scorp, 🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark), ⚙️ Lingkungan & Konfigurasi Pengujian, Scorp Agent vs PicoClaw vs ZeroClaw (Fresh API Keys & Fast-Path Heuristic)

### Community 87 - "🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut"
Cohesion: 0.29
Nodes (6): 🏛️ 15 Skenario Uji Ekstrem, 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Waktu Eksekusi 15 Skenario Ekstrem (Detik), 🔍 3. Analisis Mendalam Keandalan Arsitektur, 🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut, Scorp Agent vs PicoClaw vs ZeroClaw (Head-to-Head Live VPS Evaluation)

### Community 89 - "registerMCPToolsAsNative"
Cohesion: 0.15
Nodes (15): autoStats, sync.Mutex, TruncOutputTool(), buildArgDefsFromInputSchema(), MCPToolsDeferred(), rebuildMCPToolList(), registerMCPToolsAsNative(), StartMCPServers() (+7 more)

### Community 100 - "Transpiler Phase 3b: Contract Verify & Graceful Degradation"
Cohesion: 0.20
Nodes (14): Multi-Provider LLM Gateway (gateway/gateway.go + models.json), mark3labs/mcp-go SDK, mcp-fetch Reference Port (web scraping/markdown), mcp-filesystem Reference Port (scoped file access), mcp-sqlite Reference Port (local DB inspection), wahyuzero/scorp-mcp-registry (public registry repo), AI Transpiler (Self-Hosted Codegen), Transpiler Phase 3a: Sandbox Build (+6 more)

### Community 101 - "readInteractiveInput"
Cohesion: 0.23
Nodes (11): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste(), getContextPill() (+3 more)

### Community 102 - "MCP Marketplace Build Plan & Execution Roadmap"
Cohesion: 0.21
Nodes (13): MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry, Security Layer 2: AST & Prompt Injection Scan, Security Layer 3: Hermetic CI/CD, ~/.scorp/mcp.json Central Config Entry, mcp_manage Agent Tool (mcp/manage.go), Scorp MCP Marketplace, MCP Marketplace Blueprint (docs/MCP_MARKETPLACE_BLUEPRINT.md) (+5 more)

### Community 103 - "GetAllTools"
Cohesion: 0.33
Nodes (8): unregisterMCPNativeTools(), GetAllTools(), getChatLock(), StartTestEndpoint(), countActiveTools(), countDeferredTools(), ExecuteToolList(), ExecuteToolSearch()

### Community 104 - "ExecuteTool"
Cohesion: 0.44
Nodes (7): TestExecuteToolDenyRulesFirst(), registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks(), ExecuteTool()

### Community 105 - "LoadConfig"
Cohesion: 0.33
Nodes (9): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), Init(), StartServer() (+1 more)

### Community 106 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 107 - "TestRegisterPlugin"
Cohesion: 0.24
Nodes (5): echoPlugin, RegisterPlugin(), TestRegisterPlugin(), ToolPlugin, ToolPluginWithSchema

### Community 108 - "TopProcess"
Cohesion: 0.40
Nodes (3): CollectSystem(), GetTopProcesses(), TopProcess

### Community 110 - "TransformToolDefinitions"
Cohesion: 0.60
Nodes (5): TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions()

### Community 112 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 113 - "Laporan Eksekusi Komprehensif Uji End-to-End (E2E): Fase 1 – 15"
Cohesion: 0.40
Nodes (4): 🔍 Analisis Temuan & Ketiadaan Isu, Laporan Eksekusi Komprehensif Uji End-to-End (E2E): Fase 1 – 15, 🎯 Ringkasan Eksekusi Suite, 📊 Tabel Hasil Uji Sekuensial Fase 1 – 15

### Community 114 - "clarify.go"
Cohesion: 0.24
Nodes (10): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+2 more)

### Community 118 - "TruncOutput"
Cohesion: 0.34
Nodes (14): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+6 more)

### Community 119 - "Release Workflow (tag-triggered cross-compile + GitHub release)"
Cohesion: 0.40
Nodes (6): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection

### Community 120 - "CLI Verify Deadloop (no timeout)"
Cohesion: 0.60
Nodes (6): CLI Verify Deadloop (no timeout), complete_task Explicit Tool Contract, Incident Report: Runaway CLI Verify Session (2026-09-05), Session Lock (cli_lock.go kernel advisory file lock), Heartbeat / Stall Detection, Agent Turn Timeout (context.WithTimeout, SCORP_MAX_TURN_TIMEOUT)

## Knowledge Gaps
- **113 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+108 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 272 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `HomeDir`, `TaskPlan`, `StartDaemon`, `time.Time`, `RunAgentSessionLoop`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`, `CheckServerContracts`, `getSession`, `FormatHourlyReport`, `confirmation.go`, `SandboxActive`, `session_ui.go`, `startCLI`, `collector_docker.go`, `HandleModelCallback`, `collector_coolify.go`, `registerMCPToolsAsNative`, `clarify.go`?**
  _High betweenness centrality (0.100) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `rag_vector.go`, `TaskPlan`, `getAgentSystemPrompt`, `time.Time`, `HandleTelegramAction`, `RegisterTool`, `GetAutonomyLevel`, `CreateCheckpoint`, `ToolCall`, `TestIntegrityStatus`, `getSession`, `PermissionDecision`, `confirmation.go`, `maybeCompactHistory`, `GetStringArg`, `startCLI`, `TruncateStr`, `ResetNativeToolCache`, `prepareNewTurnHistory`, `testgate.go`, `IsDangerousCommand`, `IsContinuationDirective`, `ExecuteTermuxAPI`, `GenerateContextualSessionTitle`, `ExecuteTool`?**
  _High betweenness centrality (0.071) - this node is a cross-community bridge._
- **Why does `TruncateStr()` connect `TruncateStr` to `chat.go`, `cost_router.go`, `runSubagent`, `CallOpenCodeWithTools`, `getSession`, `PermissionDecision`, `TaskPlan`, `getAgentSystemPrompt`, `api_gemini.go`, `time.Time`, `GetStringArg`, `RunAgentSessionLoop`, `TruncOutput`, `context.Context`, `registerMCPToolsAsNative`, `MCPServer`, `api_commandcode.go`?**
  _High betweenness centrality (0.063) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _113 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.05730238025271819 - nodes in this community are weakly interconnected._
- **Should `chat.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10846560846560846 - nodes in this community are weakly interconnected._