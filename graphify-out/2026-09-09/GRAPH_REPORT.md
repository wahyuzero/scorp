# Graph Report - scorp  (2026-09-09)

## Corpus Check
- 285 files · ~206,627 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2246 nodes · 5770 edges · 121 communities (107 shown, 6 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 917 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `fb80552d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- chat.go
- agent/autonomous.go
- runSubagent
- metasearch_engines.go
- HandleTelegramAction
- rag_vector.go
- TaskPlan
- Scorp Agent (Go, ultra-light autonomous agent)
- LoadMCPConfig
- Benchmark
- StartDaemon
- getAgentSystemPrompt
- time.Time
- 🔍 Rincian Detail Uji per Fase
- ToolCall
- testing.T
- ScorpPath
- ExecuteTool
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
- cost_router.go
- TestIntegrityStatus
- CheckServerContracts
- client.go
- skills.go
- config/hooks.go
- browser_session.go
- PermissionDecision
- FormatHourlyReport
- checker.go
- HandleConfirmation
- SCORP — BRUTAL END-TO-END TEST PLAN
- ExecuteShell
- api_gemini.go
- SandboxActive
- maybeCompactHistory
- ResolveAPIKey
- 🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering)
- GetStringArg
- CheckDenyRules
- runLiveCase
- ConfigManager
- startCLI
- TruncateStr
- patch.go
- CallOpenCodeWithTools
- collector_docker.go
- RegisterProvider
- ResetNativeToolCache
- prepareNewTurnHistory
- net/http.Request
- TodoManager
- 🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam
- CredentialVault
- read_url.go
- testgate.go
- HandleModelCallback
- IsDangerousCommand
- os.File
- sync.Mutex
- HasGreenTestRun
- 🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks)
- IsContinuationDirective
- GetAllTools
- .listenSSEStream
- install.sh
- deploy.sh
- collector_coolify.go
- TestPhase6_AllTools
- registry/registry.go
- GenerateContextualSessionTitle
- 🔍 Temuan 10 Kelemahan Arsitektural & Solusi Permanen
- handleTelegramMarketplaceInstall
- main
- 🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark)
- 🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut
- registerMCPToolsAsNative
- wireCLICallbacks
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- tools/monitor.go
- Transpiler Phase 3b: Contract Verify & Graceful Degradation
- handleMCPCommand
- MCP Marketplace Build Plan & Execution Roadmap
- ExecuteAutonomous
- registerHookProbeTool
- LoadConfig
- inline.go
- RecordToolReceipt
- renderStatusFooter
- ExecuteTermuxAPI
- TransformToolDefinitions
- ExecuteAutoLogin
- serviceBridgeRequests
- Laporan Eksekusi Komprehensif Uji End-to-End (E2E): Fase 1 – 15
- clarify.go
- StopMCPServers
- RedactSecrets
- acquireSessionLock
- Release Workflow (tag-triggered cross-compile + GitHub release)
- CLI Verify Deadloop (no timeout)
- TestToolSearchActivatesDeferredToolInStaticMode

## God Nodes (most connected - your core abstractions)
1. `ModelConfig` - 89 edges
2. `HandleTelegramAction()` - 80 edges
3. `ChatMessage` - 75 edges
4. `RunAgentSessionLoop()` - 64 edges
5. `startCLI()` - 52 edges
6. `TruncateStr()` - 47 edges
7. `GetStringArg()` - 47 edges
8. `StartDaemon()` - 44 edges
9. `ScorpPath()` - 43 edges
10. `resumeAgentLoop()` - 40 edges

## Surprising Connections (you probably didn't know these)
- `Durable Memory MEMORY.md (P1.7)` --implements--> `extractTaskMemory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/memory_md.go
- `MCP Contract Watch (P3.14)` --implements--> `CheckServerContracts()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → mcp/contract.go
- `Test-Integrity Gate (P0.3)` --implements--> `TestIntegrityStatus()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → tools/testgate.go
- `Auto-Mode Classifier (P3.13)` --implements--> `PermissionDecision()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/auto.go
- `Compaction Preservation (P2.10)` --implements--> `truncateToolResultsInHistory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/compaction.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **MCP Marketplace 5-Layer Security Stack** — docs_build_plan_mcp_marketplace_layer1_source_only_registry, docs_build_plan_mcp_marketplace_layer2_ast_prompt_scan, docs_build_plan_mcp_marketplace_layer3_hermetic_ci, docs_build_plan_mcp_marketplace_sha256_pinning, readme_outbound_secret_redactor [EXTRACTED 1.00]
- **Scorp Resilient State & Memory Workflow (PlanMode + TaskLedger + Checkpoint + MEMORY.md)** — docs_implementation_plan_scorp_plan_mode, docs_implementation_plan_scorp_task_ledger_persistence, docs_implementation_plan_scorp_checkpoint_rewind, docs_implementation_plan_scorp_durable_memory [EXTRACTED 1.00]
- **Scorp Trust & Safety Gate Stack (Sandbox + DenyRules + Hooks + AutoClassifier)** — docs_implementation_plan_scorp_sandbox_bwrap, docs_implementation_plan_scorp_deny_rule_engine, docs_implementation_plan_scorp_hooks_lifecycle, docs_implementation_plan_scorp_auto_mode_classifier [EXTRACTED 1.00]
- **Scorp Verification & Integrity Pipeline (TestGate + ClaimGate + EvalArena)** — docs_implementation_plan_scorp_test_integrity_gate, docs_implementation_plan_scorp_claim_gate, docs_implementation_plan_scorp_eval_arena [EXTRACTED 1.00]
- **AI Transpiler Pipeline (Probe -> Generate -> Build -> Verify)** — docs_build_plan_mcp_marketplace_transpiler, docs_build_plan_mcp_marketplace_transpiler_probe_phase, docs_build_plan_mcp_marketplace_transpiler_generate_phase, docs_build_plan_mcp_marketplace_transpiler_build_phase, docs_build_plan_mcp_marketplace_transpiler_verify_phase [EXTRACTED 1.00]
- **Tri-Option Install Convergence onto ~/.scorp/mcp.json** — docs_build_plan_mcp_marketplace_tri_option_install, docs_build_plan_mcp_marketplace_mcp_json_config, docs_build_plan_mcp_marketplace_mcp_manage, docs_build_plan_mcp_marketplace_watchdog, docs_build_plan_mcp_marketplace_tool_registry [EXTRACTED 1.00]

## Communities (121 total, 6 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (70): Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo, RegistryIndex (+62 more)

### Community 1 - "chat.go"
Cohesion: 0.06
Nodes (73): agentSession, appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), ClearStopRequest(), collectTableLines(), convertInlineMarkdown() (+65 more)

### Community 2 - "agent/autonomous.go"
Cohesion: 0.15
Nodes (23): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), makeDecision() (+15 more)

### Community 3 - "runSubagent"
Cohesion: 0.05
Nodes (59): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+51 more)

### Community 4 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (40): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+32 more)

### Community 5 - "HandleTelegramAction"
Cohesion: 0.15
Nodes (27): RequestStop(), FormatUsageStats(), HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo() (+19 more)

### Community 6 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 7 - "TaskPlan"
Cohesion: 0.09
Nodes (43): PlanItem, ApprovePlan(), BeginPlanning(), CancelPlan(), EndPlanning(), AgentMessage, PlanningState(), presentPlanTelegram() (+35 more)

### Community 8 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.16
Nodes (19): Scorp vs PicoClaw vs ZeroClaw Architecture Comparison, ZeroClaw Defense-in-Depth Sandboxing, Scorp Embedded Native Multi-Engine Metasearch, Scorp Local Simhash & Vector RAG (0 external DB), ZeroClaw (ZeroClaw Labs, Rust), ZeroClaw Third-Party Search Stack (Tavily/Brave/Jina/SearXNG), CLI Slash Commands (/help /models /model /mode /session /sop /receipts /tools /cost /clear /exit), Cryptographic SHA-256 Tool Execution Receipts (+11 more)

### Community 9 - "LoadMCPConfig"
Cohesion: 0.31
Nodes (13): MCPConfigFilePath(), LoadMCPConfig(), rebuildMCPToolList(), ReloadMCPServers(), sanitizeMCPName(), StartMCPServers(), AddServerEntry(), ExecuteMCPManage() (+5 more)

### Community 10 - "Benchmark"
Cohesion: 0.09
Nodes (39): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+31 more)

### Community 11 - "StartDaemon"
Cohesion: 0.14
Nodes (26): StopMCPServerMode(), Init(), StartServer(), StopServer(), InitModelUsage(), runCommandLoop(), StartDaemon(), AnswerCallback() (+18 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.06
Nodes (46): getSharedMemorySummary(), AppendMemoryMD(), extractTaskMemory(), AgentMessage, parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile() (+38 more)

### Community 13 - "time.Time"
Cohesion: 0.16
Nodes (28): time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), dispatchDueTasks(), ExecuteSchedule(), gateScheduledShell() (+20 more)

### Community 14 - "🔍 Rincian Detail Uji per Fase"
Cohesion: 0.10
Nodes (19): 1️⃣1️⃣ Fase 11: Recursive Parent Lock & Deep Traversal, 1️⃣2️⃣ Fase 12: Synthetic Feedback & Unicode RTL/Emoji Edge Cases, 1️⃣3️⃣ Fase 13: Output Buffer Flooding & EPIPE Signals, 1️⃣4️⃣ Fase 14: Session Environment Persistence & Error Reporting, 1️⃣5️⃣ Fase 15: Subshell Mutation Isolation & Polyglot Compilation, 1️⃣ Fase 1: Tool Parsing & Thought Signature Contract, 2️⃣ Fase 2: Filesystem Whitelist & Bare-Metal Sandbox Toggle, 3️⃣ Fase 3: Subprocess Pipe Deadlock & Asynchronous Draining (+11 more)

### Community 15 - "ToolCall"
Cohesion: 0.22
Nodes (8): TestParseToolCalls(), formatMessagesForCLI(), init(), ClaudeCliProvider, ToolCall, ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls()

### Community 16 - "testing.T"
Cohesion: 0.07
Nodes (43): TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestIsDangerousCommand() (+35 more)

### Community 17 - "ScorpPath"
Cohesion: 0.11
Nodes (27): init(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), HomeDir(), Hostname() (+19 more)

### Community 18 - "ExecuteTool"
Cohesion: 0.13
Nodes (18): TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), ResetAutoStats(), TestExecuteToolAutoPresetTrustedAndRecorded(), ExecuteTool(), init(), init(), init() (+10 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.24
Nodes (16): AgentMessage, confirmationDisplay(), ConsumeStopRequest(), confirmKeyboard(), cleanToolCallTags(), getSessionSearchContext(), maxIterations(), maxTurnTimeout() (+8 more)

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
Cohesion: 0.09
Nodes (15): context.Context, AnthropicProvider, AzureProvider, ChatMessage, CohereProvider, DashScopeProvider, FastInferenceProvider, GeminiProvider (+7 more)

### Community 24 - "GetAutonomyLevel"
Cohesion: 0.16
Nodes (19): GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels(), TestConfirmationRequired() (+11 more)

### Community 25 - "collector_system.go"
Cohesion: 0.23
Nodes (18): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+10 more)

### Community 26 - "RunPreToolUseHooks"
Cohesion: 0.21
Nodes (20): PreToolUse & PostToolUse Hooks (P3.12), Community Praised Agent Patterns (2026), 2026 AI Coding Agent Competitor Landscape, hookPayload, appendHookContext(), runHookCommand(), RunPostToolUseHooks(), RunPreToolUseHooks() (+12 more)

### Community 27 - "MCPServer"
Cohesion: 0.16
Nodes (14): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, sync.Once, FindMCPTool(), MCPServer, MCPTool, ProbeServer() (+6 more)

### Community 28 - "api_commandcode.go"
Cohesion: 0.17
Nodes (15): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeStream(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), resolveCommandCodeKeyFromDisk(), commandCodeMsg (+7 more)

### Community 29 - "CreateCheckpoint"
Cohesion: 0.26
Nodes (19): Checkpoint and Rewind (P1.6), caseCheckpoint(), CheckpointDiffStat(), checkpointRepoRoot(), ckptGit(), ckptRefFor(), CreateCheckpoint(), DeleteCheckpoint() (+11 more)

### Community 30 - "collector_system_native.go"
Cohesion: 0.10
Nodes (35): FormatDuration(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+27 more)

### Community 31 - "cost_router.go"
Cohesion: 0.19
Nodes (18): ConfigMgr(), defaultCostConfig(), formatCostReport(), FormatDailyCostSummary(), getCheapestModel(), handleCostCommand(), init(), isBudgetExceeded() (+10 more)

### Community 32 - "TestIntegrityStatus"
Cohesion: 0.38
Nodes (13): IsTestRelatedPath(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestTestIntegrityStatus_FailingSuiteDoesNotCount(), TestTestIntegrityStatus_GreenRunMustBeAfterEdit(), TestTestIntegrityStatus_NoTouches() (+5 more)

### Community 33 - "CheckServerContracts"
Cohesion: 0.24
Nodes (14): caseMCPContract(), CheckServerContracts(), contractFile(), ContractStatusNotice(), ContractWarnings(), serverFingerprint(), SetContractPath(), contractTestServer() (+6 more)

### Community 34 - "client.go"
Cohesion: 0.19
Nodes (18): encoding/json.RawMessage, executeMCPServerTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError(), sendMCPResult() (+10 more)

### Community 35 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 36 - "config/hooks.go"
Cohesion: 0.30
Nodes (13): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+5 more)

### Community 37 - "browser_session.go"
Cohesion: 0.33
Nodes (13): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+5 more)

### Community 38 - "PermissionDecision"
Cohesion: 0.22
Nodes (17): TestAutoAllowlistPrefixAnchoring(), autoAllowlisted(), autoClassify(), autoClassifyWithModel(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand(), PermissionDecision() (+9 more)

### Community 39 - "FormatHourlyReport"
Cohesion: 0.18
Nodes (16): NetworkData, SystemData, CollectSystem(), GetTopProcesses(), TopProcess, PortInfo, Bar(), bar() (+8 more)

### Community 40 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 41 - "HandleConfirmation"
Cohesion: 0.36
Nodes (10): TestStorePendingConfirmationArgs(), clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), StorePendingConfirmation() (+2 more)

### Community 42 - "SCORP — BRUTAL END-TO-END TEST PLAN"
Cohesion: 0.09
Nodes (20): 0. SCALE & SEVERITY DEFINITIONS, A. GATE-STACK ADVERSARIAL (S) — Layer Penetration Probes, B. LONG-HORIZON & CONTEXT COMPACTION (L), C. INFRASTRUCTURE CHAOS (M/L), D. CONCURRENCY & RACE CONDITIONS (M), E. TELEGRAM UX INTEGRITY (M), F. ADVERSARIAL SECURITY (M), G. EVALUATION & DEPLOYMENT GATE INTEGRITY (S) (+12 more)

### Community 43 - "ExecuteShell"
Cohesion: 0.13
Nodes (22): init(), Sandbox Shell Bubblewrap (P0.1), WriteFileAtomic(), captureSessionExports(), cleanAbortedShm(), decodeUTF16(), decodeUTFString(), ExecuteListDir() (+14 more)

### Community 44 - "api_gemini.go"
Cohesion: 0.23
Nodes (19): callGemini(), callGeminiStream(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), parseDataURL(), resolveGeminiBaseURL() (+11 more)

### Community 45 - "SandboxActive"
Cohesion: 0.25
Nodes (15): caseSandbox(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest(), SandboxStatusNotice(), SandboxVersion() (+7 more)

### Community 46 - "maybeCompactHistory"
Cohesion: 0.25
Nodes (17): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, injectPreservationNote(), latestUserGoal(), maybeCompactHistory(), preservationNote() (+9 more)

### Community 47 - "ResolveAPIKey"
Cohesion: 0.10
Nodes (34): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestGemini_MultimodalParsing(), TestOpenCodeProvider_RegistrationAndKey(), ChatResponse, CallModel(), CallModelStream() (+26 more)

### Community 48 - "🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering)"
Cohesion: 0.18
Nodes (10): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Brutal, 🔍 3. Analisis Mendalam Kualitas Output & Arsitektur, 🏆 4. Kesimpulan Akhir Benchmark Brutal, A. Kompilasi Bahasa Go Multi-File (Kasus B2), B. Otomasi Server REST API Latar Belakang & Health Check (Kasus B7), 🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering), C. Fenomena Kebocoran Tag Simulasi pada ZeroClaw (Kasus B1, B3, B5, B8, B10, B12, B15) (+2 more)

### Community 49 - "GetStringArg"
Cohesion: 0.10
Nodes (25): init(), ConfirmationRequired(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage(), TruncOutput() (+17 more)

### Community 50 - "CheckDenyRules"
Cohesion: 0.16
Nodes (17): TestExecuteToolDenyRulesFirst(), CheckDenyRules(), loadDenyRules(), ParseDenyRule(), ReloadDenyRules(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped() (+9 more)

### Community 51 - "runLiveCase"
Cohesion: 0.14
Nodes (21): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, Case, caseResult, CoreCases(), evalSandboxDir(), humanCount(), liveLabel() (+13 more)

### Community 52 - "ConfigManager"
Cohesion: 0.19
Nodes (6): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, os.FileMode, TestScheduledShellGateAndGroupKill()

### Community 53 - "startCLI"
Cohesion: 0.28
Nodes (17): executeOneShot(), executeTurn(), handleCLISOP(), hasDebugFlag(), printBanner(), printCLIHelp(), printCostUsage(), printCronTasks() (+9 more)

### Community 54 - "TruncateStr"
Cohesion: 0.22
Nodes (25): TruncateStr(), anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic(), callAnthropicStream() (+17 more)

### Community 55 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "CallOpenCodeWithTools"
Cohesion: 0.25
Nodes (9): CallOpenCode(), CallOpenCodeStream(), CallOpenCodeWithTools(), resolveOpenCodeBaseURL(), resolveOpenCodeKeyFromDisk(), resolveOpenCodeSessionID(), RecordCostWithCache(), OpenCodeProvider (+1 more)

### Community 57 - "collector_docker.go"
Cohesion: 0.27
Nodes (11): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+3 more)

### Community 58 - "RegisterProvider"
Cohesion: 0.11
Nodes (19): init(), init(), init(), init(), init(), init(), init(), init() (+11 more)

### Community 59 - "ResetNativeToolCache"
Cohesion: 0.29
Nodes (13): ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), clearDynamicEnv(), registerTempTool(), TestDynamicToolTTL() (+5 more)

### Community 60 - "prepareNewTurnHistory"
Cohesion: 0.29
Nodes (8): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary()

### Community 61 - "net/http.Request"
Cohesion: 0.28
Nodes (12): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), TestGatewayEndpoints() (+4 more)

### Community 62 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 63 - "🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam"
Cohesion: 0.25
Nodes (7): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Ultra-Hardcore, 🔍 3. Temuan Kritis: Mengapa ZeroClaw Runtuh Total?, 💡 4. Analisis Komparasi Mendalam: Scorp vs PicoClaw, 🏆 5. Rekapitulasi Menyeluruh (Grand Total 65 Skenario Uji), 🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam, Scorp Agent vs PicoClaw vs ZeroClaw (Live Head-to-Head VPS Evaluation)

### Community 64 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 65 - "read_url.go"
Cohesion: 0.60
Nodes (5): ReadURL(), scrapeFirecrawl(), scrapeTavily(), truncateURLOutput(), tryRemoteScrape()

### Community 66 - "testgate.go"
Cohesion: 0.20
Nodes (16): OperationalClaim, dedupeOpObjects(), extractOpObjects(), isWriteToolName(), LooksLikeOperationalClaims(), opReceipt(), seedOpReceipts(), TestLooksLikeOperationalClaims() (+8 more)

### Community 67 - "HandleModelCallback"
Cohesion: 0.07
Nodes (62): apiKey, envKeyName, CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog() (+54 more)

### Community 68 - "IsDangerousCommand"
Cohesion: 0.39
Nodes (6): devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), inspectBase64Payload(), IsDangerousCommand(), isHarmlessDevSink()

### Community 69 - "os.File"
Cohesion: 0.21
Nodes (9): lockFileExclusive(), unlockFile(), lockFileExclusive(), unlockFile(), os.File, fileLockExclusive(), fileUnlock(), fileLockExclusive() (+1 more)

### Community 70 - "sync.Mutex"
Cohesion: 0.50
Nodes (4): autoStats, sync.Mutex, getChatLock(), StartTestEndpoint()

### Community 71 - "HasGreenTestRun"
Cohesion: 0.19
Nodes (12): Evidence-Based Claim Gate (P4.16), MCP Contract Watch (P3.14), Test-Integrity Gate (P0.3), Agent Failure Modes and Critiques (2026), caseClaimGate(), TestHasGreenTestRun(), TestLooksLikeTestPassClaim(), HasGreenTestRun() (+4 more)

### Community 72 - "🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks)"
Cohesion: 0.29
Nodes (6): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Agentic Berat, 🔍 3. Kesimpulan Utama Uji Komparasi Baru, 🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks), ⚙️ Lingkungan & Konfigurasi Pengujian, Scorp Agent vs PicoClaw vs ZeroClaw (Fresh API Keys & Fast-Path Heuristic)

### Community 73 - "IsContinuationDirective"
Cohesion: 0.40
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 75 - "GetAllTools"
Cohesion: 0.60
Nodes (5): unregisterMCPNativeTools(), GetAllTools(), countActiveTools(), countDeferredTools(), ExecuteToolList()

### Community 76 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 78 - "deploy.sh"
Cohesion: 0.53
Nodes (4): die(), ok(), deploy.sh script, step()

### Community 79 - "collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 80 - "TestPhase6_AllTools"
Cohesion: 0.26
Nodes (14): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession() (+6 more)

### Community 81 - "registry/registry.go"
Cohesion: 0.17
Nodes (11): echoPlugin, RegisterPlugin(), TestRegisterPlugin(), ExecuteToolByName(), GenerateSystemPromptDescriptions(), GetTool(), GetToolsByCategory(), ArgDef (+3 more)

### Community 82 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 83 - "🔍 Temuan 10 Kelemahan Arsitektural & Solusi Permanen"
Cohesion: 0.17
Nodes (11): 1. In-Memory Session Environment Non-Persistence across Independent CLI Invocations, 2. Truncation Hazard pada `write_file`, `replace_file_content`, & Bridge Protocol, 3. CRLF (`\r\n`) Line-Ending Mismatch pada Surgical Chunk Replacement, 4. 64-Bit Cryptographic Receipt Hash Collision Blindspot, 5. Infinite Traversal Hang pada Symlink Cycles & Inode Tracking, 6. Binary Non-UTF8 Stream Corruption di Shell Output Buffer, 7. Over-Inspection Loop pada Skenario Modifikasi File, Laporan Audit Kelemahan Arsitektural Non-Kecepatan & Penguatan Sistem (Fase 16 – 18) (+3 more)

### Community 84 - "handleTelegramMarketplaceInstall"
Cohesion: 0.60
Nodes (5): handleTelegramMarketplaceInstall(), MarketplaceTriOptionKeyboard(), reportInstallOutcome(), BackButtonKeyboard(), SendMessage()

### Community 85 - "main"
Cohesion: 0.23
Nodes (12): RegisterAutonomous(), StartGateway(), isCLIMode(), main(), SOP, Dir(), GetSOP(), InitDefaultSOPs() (+4 more)

### Community 86 - "🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark)"
Cohesion: 0.29
Nodes (6): 📊 1. Ringkasan Hasil Eksekusi, ⏱️ 2. Tabel Waktu Respon Per Tugas (Detik), 🔍 3. Analisis Peningkatan Scorp, 🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark), ⚙️ Lingkungan & Konfigurasi Pengujian, Scorp Agent vs PicoClaw vs ZeroClaw (Fresh API Keys & Fast-Path Heuristic)

### Community 87 - "🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut"
Cohesion: 0.29
Nodes (6): 🏛️ 15 Skenario Uji Ekstrem, 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Waktu Eksekusi 15 Skenario Ekstrem (Detik), 🔍 3. Analisis Mendalam Keandalan Arsitektur, 🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut, Scorp Agent vs PicoClaw vs ZeroClaw (Head-to-Head Live VPS Evaluation)

### Community 89 - "registerMCPToolsAsNative"
Cohesion: 0.22
Nodes (10): TruncOutputTool(), buildArgDefsFromInputSchema(), MCPToolsDeferred(), registerMCPToolsAsNative(), TestMCPToolsDeferredEnvParsing(), ServerWatchdog, GetServerHealthStatus(), MCPServer (+2 more)

### Community 90 - "wireCLICallbacks"
Cohesion: 0.29
Nodes (8): formatFileSize(), wireCLICallbacks(), formatFinalResponse(), formatTerminalText(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 99 - "tools/monitor.go"
Cohesion: 0.38
Nodes (10): ExecuteMonitor(), InitMonitor(), loadMonitorTargets(), monitorCheckOne(), monitorLoop(), ragIngestText(), sanitizeFilename(), saveMonitorTargets() (+2 more)

### Community 100 - "Transpiler Phase 3b: Contract Verify & Graceful Degradation"
Cohesion: 0.20
Nodes (14): Multi-Provider LLM Gateway (gateway/gateway.go + models.json), mark3labs/mcp-go SDK, mcp-fetch Reference Port (web scraping/markdown), mcp-filesystem Reference Port (scoped file access), mcp-sqlite Reference Port (local DB inspection), wahyuzero/scorp-mcp-registry (public registry repo), AI Transpiler (Self-Hosted Codegen), Transpiler Phase 3a: Sandbox Build (+6 more)

### Community 101 - "handleMCPCommand"
Cohesion: 0.31
Nodes (7): handleMCPCommand(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste()

### Community 102 - "MCP Marketplace Build Plan & Execution Roadmap"
Cohesion: 0.21
Nodes (13): MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry, Security Layer 2: AST & Prompt Injection Scan, Security Layer 3: Hermetic CI/CD, ~/.scorp/mcp.json Central Config Entry, mcp_manage Agent Tool (mcp/manage.go), Scorp MCP Marketplace, MCP Marketplace Blueprint (docs/MCP_MARKETPLACE_BLUEPRINT.md) (+5 more)

### Community 103 - "ExecuteAutonomous"
Cohesion: 0.36
Nodes (8): LoadAutonomousConfig(), SaveAutonomousConfig(), TestPhase7_ConfigPersistence(), autoShowActions(), autoShowConfig(), autoShowLog(), autoStatus(), ExecuteAutonomous()

### Community 104 - "registerHookProbeTool"
Cohesion: 0.73
Nodes (5): registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks()

### Community 105 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 106 - "inline.go"
Cohesion: 0.33
Nodes (10): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+2 more)

### Community 107 - "RecordToolReceipt"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 108 - "renderStatusFooter"
Cohesion: 0.53
Nodes (5): GetDailyTotalUSD(), getContextPill(), getGitStatus(), getShortCwd(), renderStatusFooter()

### Community 109 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 110 - "TransformToolDefinitions"
Cohesion: 0.33
Nodes (8): ChatRequest, TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions(), ToolSchema, ToolSchemaFunc

### Community 112 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 113 - "Laporan Eksekusi Komprehensif Uji End-to-End (E2E): Fase 1 – 15"
Cohesion: 0.40
Nodes (4): 🔍 Analisis Temuan & Ketiadaan Isu, Laporan Eksekusi Komprehensif Uji End-to-End (E2E): Fase 1 – 15, 🎯 Ringkasan Eksekusi Suite, 📊 Tabel Hasil Uji Sekuensial Fase 1 – 15

### Community 114 - "clarify.go"
Cohesion: 0.27
Nodes (9): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID() (+1 more)

### Community 117 - "acquireSessionLock"
Cohesion: 0.40
Nodes (3): acquireSessionLock(), TestAcquireSessionLock(), sessionLockFile

### Community 119 - "Release Workflow (tag-triggered cross-compile + GitHub release)"
Cohesion: 0.40
Nodes (6): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection

### Community 120 - "CLI Verify Deadloop (no timeout)"
Cohesion: 0.60
Nodes (6): CLI Verify Deadloop (no timeout), complete_task Explicit Tool Contract, Incident Report: Runaway CLI Verify Session (2026-09-05), Session Lock (cli_lock.go kernel advisory file lock), Heartbeat / Stall Detection, Agent Turn Timeout (context.WithTimeout, SCORP_MAX_TURN_TIMEOUT)

### Community 121 - "TestToolSearchActivatesDeferredToolInStaticMode"
Cohesion: 0.70
Nodes (4): deferredToolActive(), osUnsetDynamic(), registerSearchTarget(), TestToolSearchActivatesDeferredToolInStaticMode()

## Knowledge Gaps
- **97 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+92 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 253 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `TaskPlan`, `StartDaemon`, `getAgentSystemPrompt`, `time.Time`, `RunAgentSessionLoop`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`, `CheckServerContracts`, `FormatHourlyReport`, `HandleConfirmation`, `SandboxActive`, `collector_docker.go`, `HandleModelCallback`, `collector_coolify.go`, `handleTelegramMarketplaceInstall`, `main`, `registerMCPToolsAsNative`, `clarify.go`?**
  _High betweenness centrality (0.085) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `HandleTelegramAction`, `rag_vector.go`, `TaskPlan`, `getAgentSystemPrompt`, `time.Time`, `ToolCall`, `ExecuteTool`, `GetAutonomyLevel`, `CreateCheckpoint`, `TestIntegrityStatus`, `PermissionDecision`, `HandleConfirmation`, `maybeCompactHistory`, `ResolveAPIKey`, `GetStringArg`, `startCLI`, `TruncateStr`, `ResetNativeToolCache`, `prepareNewTurnHistory`, `testgate.go`, `IsDangerousCommand`, `HasGreenTestRun`, `IsContinuationDirective`, `GenerateContextualSessionTitle`, `wireCLICallbacks`, `ExecuteTermuxAPI`?**
  _High betweenness centrality (0.079) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `chat.go`, `agent/autonomous.go`, `HandleTelegramAction`, `rag_vector.go`, `LoadMCPConfig`, `getAgentSystemPrompt`, `time.Time`, `ScorpPath`, `ExecuteTool`, `cost_router.go`, `client.go`, `skills.go`, `HandleConfirmation`, `ConfigManager`, `collector_docker.go`, `HandleModelCallback`, `IsDangerousCommand`, `sync.Mutex`, `handleTelegramMarketplaceInstall`, `main`, `StopMCPServers`?**
  _High betweenness centrality (0.058) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _97 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.06044303797468355 - nodes in this community are weakly interconnected._
- **Should `chat.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05660945498343872 - nodes in this community are weakly interconnected._