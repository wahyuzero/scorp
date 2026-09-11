# Graph Report - scorp  (2026-09-11)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 2348 nodes · 5987 edges · 138 communities (117 shown, 12 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 944 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a4742f91`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- collector_system.go
- runSubagent
- HandleModelCallback
- TaskPlan
- context.Context
- getAgentSystemPrompt
- testing.T
- Benchmark
- readInputEngine
- compaction_test.go
- testgate.go
- RunPreToolUseHooks
- rag_vector.go
- RegisterTool
- ScorpPath
- RegisterProvider
- StartDaemon
- GetAutonomyLevel
- ResolveAPIKey
- GetStringArg
- collector_security.go
- resolver.go
- setSession
- HandleTelegramAction
- session_search_fts5.go
- ModelConfig
- PermissionDecision
- CreateCheckpoint
- client.go
- agent/autonomous.go
- SCORP — BRUTAL END-TO-END TEST PLAN
- chat.go
- ToolCall
- RAGIndex
- scheduler.go
- metasearch_engines.go
- time.Time
- MCPServer
- startCLI
- 🔍 Rincian Detail Uji per Fase
- api_gemini.go
- CallModelWithFallback
- getSession
- TruncOutput
- Scorp Agent (Go, ultra-light autonomous agent)
- RunAgentSessionLoop
- main
- GetIntArg
- skills.go
- TruncateStr
- confirmation.go
- os.File
- CheckDenyRules
- SandboxActive
- checker.go
- CheckServerContracts
- ResetNativeToolCache
- ConfigManager
- LoadMCPConfig
- net/http.Request
- Transpiler Phase 3b: Contract Verify & Graceful Degradation
- runLiveCase
- cost_router.go
- TestIntegrityStatus
- session_ui.go
- MCP Marketplace Build Plan & Execution Roadmap
- eval.go
- ComputeSimhash
- init
- 🔍 Temuan 10 Kelemahan Arsitektural & Solusi Permanen
- time.Duration
- EscapeHTML
- NewDefaultMetaSearchAggregator
- patchReplace
- 🔍 Temuan Kelemahan & Solusi Arsitektural Permanen
- 🔍 Temuan Kelemahan & Solusi Arsitektural Permanen
- 🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering)
- clarify.go
- CredentialVault
- tools/monitor.go
- TodoManager
- wireCLICallbacks
- inline.go
- ExecuteTool
- TestPhase6_AllTools
- IsDangerousCommand
- net/http.Client
- startMCPServer
- GenerateNativeToolsSchema
- metasearch_test.go
- ExecuteReadURL
- ConfigMgr
- 🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam
- ExecuteTermuxAPI
- ExecuteAutonomous
- 🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut
- 🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks)
- 🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark)
- eval/core.go
- ExecuteScript
- Release Workflow (tag-triggered cross-compile + GitHub release)
- IsContinuationDirective
- acquireSessionLock
- CollectSystem
- LoadConfig
- CLI Verify Deadloop (no timeout)
- .listenSSEStream
- registerMCPToolsAsNative
- RestartServer
- deploy.sh
- GenerateContextualSessionTitle
- Laporan Eksekusi Komprehensif Uji End-to-End (E2E): Fase 1 – 15
- echoPlugin
- serviceBridgeRequests
- os.FileMode
- RegisterPlugin
- BraveSearchEngine
- ExecuteSQL
- GitHubEngine
- TavilySearchEngine
- WikipediaEngine
- RedactSecrets
- ShareToMarketplace
- ExecuteAutoLogin
- install.sh
- MCP Deferred-by-Default (P2.9)
- MCPServer
- scorp-agent

## God Nodes (most connected - your core abstractions)
1. `ModelConfig` - 89 edges
2. `HandleTelegramAction()` - 81 edges
3. `ChatMessage` - 75 edges
4. `RunAgentSessionLoop()` - 64 edges
5. `startCLI()` - 53 edges
6. `GetStringArg()` - 47 edges
7. `TruncateStr()` - 47 edges
8. `StartDaemon()` - 45 edges
9. `ScorpPath()` - 44 edges
10. `resumeAgentLoop()` - 40 edges

## Surprising Connections (you probably didn't know these)
- `MCP Contract Watch (P3.14)` --implements--> `CheckServerContracts()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → mcp/contract.go
- `Durable Memory MEMORY.md (P1.7)` --implements--> `extractTaskMemory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/memory_md.go
- `Test-Integrity Gate (P0.3)` --implements--> `TestIntegrityStatus()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → tools/testgate.go
- `Compaction Preservation (P2.10)` --implements--> `truncateToolResultsInHistory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/compaction.go
- `StartDaemon()` --calls--> `StartCPUSampler()`  [INFERRED]
  telegram/daemon.go → collectors/collector_system_other.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **MCP Marketplace 5-Layer Security Stack** — docs_build_plan_mcp_marketplace_layer1_source_only_registry, docs_build_plan_mcp_marketplace_layer2_ast_prompt_scan, docs_build_plan_mcp_marketplace_layer3_hermetic_ci, docs_build_plan_mcp_marketplace_sha256_pinning, readme_outbound_secret_redactor [EXTRACTED 1.00]
- **Scorp Resilient State & Memory Workflow (PlanMode + TaskLedger + Checkpoint + MEMORY.md)** — docs_implementation_plan_scorp_plan_mode, docs_implementation_plan_scorp_task_ledger_persistence, docs_implementation_plan_scorp_checkpoint_rewind, docs_implementation_plan_scorp_durable_memory [EXTRACTED 1.00]
- **Scorp Trust & Safety Gate Stack (Sandbox + DenyRules + Hooks + AutoClassifier)** — docs_implementation_plan_scorp_sandbox_bwrap, docs_implementation_plan_scorp_deny_rule_engine, docs_implementation_plan_scorp_hooks_lifecycle, docs_implementation_plan_scorp_auto_mode_classifier [EXTRACTED 1.00]
- **Scorp Verification & Integrity Pipeline (TestGate + ClaimGate + EvalArena)** — docs_implementation_plan_scorp_test_integrity_gate, docs_implementation_plan_scorp_claim_gate, docs_implementation_plan_scorp_eval_arena [EXTRACTED 1.00]
- **AI Transpiler Pipeline (Probe -> Generate -> Build -> Verify)** — docs_build_plan_mcp_marketplace_transpiler, docs_build_plan_mcp_marketplace_transpiler_probe_phase, docs_build_plan_mcp_marketplace_transpiler_generate_phase, docs_build_plan_mcp_marketplace_transpiler_build_phase, docs_build_plan_mcp_marketplace_transpiler_verify_phase [EXTRACTED 1.00]
- **Tri-Option Install Convergence onto ~/.scorp/mcp.json** — docs_build_plan_mcp_marketplace_tri_option_install, docs_build_plan_mcp_marketplace_mcp_json_config, docs_build_plan_mcp_marketplace_mcp_manage, docs_build_plan_mcp_marketplace_watchdog, docs_build_plan_mcp_marketplace_tool_registry [EXTRACTED 1.00]

## Communities (138 total, 12 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (74): Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo, RegistryIndex (+66 more)

### Community 1 - "collector_system.go"
Cohesion: 0.05
Nodes (72): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+64 more)

### Community 2 - "runSubagent"
Cohesion: 0.05
Nodes (60): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+52 more)

### Community 3 - "HandleModelCallback"
Cohesion: 0.07
Nodes (62): apiKey, envKeyName, CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog() (+54 more)

### Community 4 - "TaskPlan"
Cohesion: 0.07
Nodes (51): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary() (+43 more)

### Community 5 - "context.Context"
Cohesion: 0.09
Nodes (20): context.Context, AnthropicProvider, applyOpenAIHeaders(), buildOpenAIRequestBody(), CallOpenAI(), CallOpenAIStream(), CallOpenAIWithTools(), formatOpenAIMessages() (+12 more)

### Community 6 - "getAgentSystemPrompt"
Cohesion: 0.06
Nodes (45): getSharedMemorySummary(), AppendMemoryMD(), extractTaskMemory(), AgentMessage, parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile() (+37 more)

### Community 7 - "testing.T"
Cohesion: 0.07
Nodes (42): TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestParseToolCalls() (+34 more)

### Community 8 - "Benchmark"
Cohesion: 0.10
Nodes (37): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+29 more)

### Community 9 - "readInputEngine"
Cohesion: 0.10
Nodes (37): io.Reader, io.Writer, GetDailyTotalUSD(), chunkReader, SlashCommand, clampLineWidth(), filterCommands(), isWideRune() (+29 more)

### Community 10 - "compaction_test.go"
Cohesion: 0.12
Nodes (37): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, injectPreservationNote(), latestUserGoal(), maybeCompactHistory(), preservationNote() (+29 more)

### Community 11 - "testgate.go"
Cohesion: 0.09
Nodes (34): Evidence-Based Claim Gate (P4.16), MCP Contract Watch (P3.14), Test-Integrity Gate (P0.3), Agent Failure Modes and Critiques (2026), caseClaimGate(), TestHasGreenTestRun(), TestLooksLikeTestPassClaim(), OperationalClaim (+26 more)

### Community 12 - "RunPreToolUseHooks"
Cohesion: 0.14
Nodes (31): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+23 more)

### Community 13 - "rag_vector.go"
Cohesion: 0.11
Nodes (20): activeToolCall, getBoolArg(), getFloatArg(), hybridResult, EnsureVectorRAG(), InitVectorRAG(), RagVecIngest(), RagVecList() (+12 more)

### Community 14 - "RegisterTool"
Cohesion: 0.09
Nodes (24): init(), init(), init(), init(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema(), TestIsRateLimitError(), TestRegisterPlugin() (+16 more)

### Community 15 - "ScorpPath"
Cohesion: 0.11
Nodes (29): init(), setupCLILogging(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), HomeDir() (+21 more)

### Community 16 - "RegisterProvider"
Cohesion: 0.09
Nodes (22): init(), init(), formatMessagesForCLI(), init(), init(), init(), init(), init() (+14 more)

### Community 17 - "StartDaemon"
Cohesion: 0.13
Nodes (26): Init(), StartServer(), StopServer(), runCommandLoop(), StartDaemon(), BackAndRefreshKeyboard(), baseName(), DeleteWebhook() (+18 more)

### Community 18 - "GetAutonomyLevel"
Cohesion: 0.15
Nodes (23): ConfirmationRequired(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels() (+15 more)

### Community 19 - "ResolveAPIKey"
Cohesion: 0.15
Nodes (22): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestGemini_MultimodalParsing(), TestOpenCodeProvider_RegistrationAndKey(), CallModel(), CallModelStream(), GetProvider() (+14 more)

### Community 20 - "GetStringArg"
Cohesion: 0.14
Nodes (20): init(), GetStringArg(), WriteFileAtomic(), captureSessionExports(), decodeUTF16(), decodeUTFString(), ExecuteListDir(), ExecuteReadFile() (+12 more)

### Community 21 - "collector_security.go"
Cohesion: 0.17
Nodes (25): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+17 more)

### Community 22 - "resolver.go"
Cohesion: 0.15
Nodes (25): net.Conn, net.Dialer, net.Resolver, BuildCandidateList(), deduplicate(), DialDNS(), DiscoverPlatformDNS(), getprop() (+17 more)

### Community 23 - "setSession"
Cohesion: 0.14
Nodes (25): appendSessionHistory(), EnterAgentMode(), GetHistoryTokenEstimate(), getOrCreateSession(), getSessionHistory(), getSessionMap(), AgentMessage, loadHistoryFromDisk() (+17 more)

### Community 24 - "HandleTelegramAction"
Cohesion: 0.18
Nodes (24): RequestStop(), FormatUsageStats(), HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo() (+16 more)

### Community 25 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (15): homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult, SessionResult, InitSessionDB(), SearchSessions() (+7 more)

### Community 26 - "ModelConfig"
Cohesion: 0.24
Nodes (18): anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic(), callAnthropicStream(), CallAnthropicWithTools() (+10 more)

### Community 27 - "PermissionDecision"
Cohesion: 0.18
Nodes (20): TestAutoAllowlistPrefixAnchoring(), TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), autoAllowlisted(), autoClassify(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand() (+12 more)

### Community 28 - "CreateCheckpoint"
Cohesion: 0.23
Nodes (21): Checkpoint and Rewind (P1.6), Community Praised Agent Patterns (2026), 2026 AI Coding Agent Competitor Landscape, caseCheckpoint(), CheckpointDiffStat(), checkpointRepoRoot(), ckptGit(), ckptRefFor() (+13 more)

### Community 29 - "client.go"
Cohesion: 0.17
Nodes (21): encoding/json.RawMessage, executeMCPServerTool(), FindMCPTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError() (+13 more)

### Community 30 - "agent/autonomous.go"
Cohesion: 0.18
Nodes (19): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), makeDecision() (+11 more)

### Community 31 - "SCORP — BRUTAL END-TO-END TEST PLAN"
Cohesion: 0.09
Nodes (20): 0. SCALE & SEVERITY DEFINITIONS, A. GATE-STACK ADVERSARIAL (S) — Layer Penetration Probes, B. LONG-HORIZON & CONTEXT COMPACTION (L), C. INFRASTRUCTURE CHAOS (M/L), D. CONCURRENCY & RACE CONDITIONS (M), E. TELEGRAM UX INTEGRITY (M), F. ADVERSARIAL SECURITY (M), G. EVALUATION & DEPLOYMENT GATE INTEGRITY (S) (+12 more)

### Community 32 - "chat.go"
Cohesion: 0.16
Nodes (18): cleanupChatSessions(), CleanupSessionsLoop(), collectTableLines(), convertInlineMarkdown(), convertTableToList(), ExitAgentMode(), extractAndSaveMemory(), flushPendingMessages() (+10 more)

### Community 33 - "ToolCall"
Cohesion: 0.18
Nodes (15): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeStream(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), commandCodeMsg, commandCodeParams (+7 more)

### Community 34 - "RAGIndex"
Cohesion: 0.17
Nodes (12): homeDir(), ragDirPath(), ragIndexPath(), ragVectorDBPath(), computeTF(), InitRAG(), RagIndexAdd(), RagIndexList() (+4 more)

### Community 35 - "scheduler.go"
Cohesion: 0.24
Nodes (19): ScheduledTask, AddTask(), AddTaskEx(), dispatchDueTasks(), ExecuteSchedule(), FormatTasksList(), GetTask(), LoadTasks() (+11 more)

### Community 36 - "metasearch_engines.go"
Cohesion: 0.14
Nodes (11): BingEngine, DuckDuckGoHTMLEngine, DuckDuckGoLiteEngine, ghSearchResp, GoogleCSEEngine, NewBingEngine(), NewDuckDuckGoHTMLEngine(), NewDuckDuckGoLiteEngine() (+3 more)

### Community 37 - "time.Time"
Cohesion: 0.17
Nodes (12): agentSession, TGDocument, time.Time, probeDNS(), ResilientDNSConn, AgentMessage, AutonomousConfig, AutonomousLogEntry (+4 more)

### Community 38 - "MCPServer"
Cohesion: 0.14
Nodes (10): autoStats, bufio.Scanner, context.CancelFunc, encoding/json.Encoder, sync.Mutex, sync.Once, MCPServer, ServerWatchdog (+2 more)

### Community 39 - "startCLI"
Cohesion: 0.26
Nodes (18): executeOneShot(), executeTurn(), formatTerminalText(), handleCLISession(), handleCLISOP(), handleMCPCommand(), printBanner(), printCLIHelp() (+10 more)

### Community 40 - "🔍 Rincian Detail Uji per Fase"
Cohesion: 0.10
Nodes (19): 1️⃣1️⃣ Fase 11: Recursive Parent Lock & Deep Traversal, 1️⃣2️⃣ Fase 12: Synthetic Feedback & Unicode RTL/Emoji Edge Cases, 1️⃣3️⃣ Fase 13: Output Buffer Flooding & EPIPE Signals, 1️⃣4️⃣ Fase 14: Session Environment Persistence & Error Reporting, 1️⃣5️⃣ Fase 15: Subshell Mutation Isolation & Polyglot Compilation, 1️⃣ Fase 1: Tool Parsing & Thought Signature Contract, 2️⃣ Fase 2: Filesystem Whitelist & Bare-Metal Sandbox Toggle, 3️⃣ Fase 3: Subprocess Pipe Deadlock & Asynchronous Draining (+11 more)

### Community 41 - "api_gemini.go"
Cohesion: 0.23
Nodes (19): callGemini(), callGeminiStream(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), parseDataURL(), resolveGeminiBaseURL() (+11 more)

### Community 42 - "CallModelWithFallback"
Cohesion: 0.18
Nodes (16): AgentMessage, autoClassifyWithModel(), ChatResponse, getCheapestModel(), RouteModelCostAware(), CostTracker, CallModelWithFallback(), findFirstVisionModel() (+8 more)

### Community 43 - "getSession"
Cohesion: 0.22
Nodes (17): ClearChatSession(), FlushSessionHistory(), getSession(), historyFilePath(), historyWriterLoop(), init(), IsAgentMode(), IsLoopActive() (+9 more)

### Community 44 - "TruncOutput"
Cohesion: 0.28
Nodes (15): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+7 more)

### Community 45 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.16
Nodes (19): Scorp vs PicoClaw vs ZeroClaw Architecture Comparison, ZeroClaw Defense-in-Depth Sandboxing, Scorp Embedded Native Multi-Engine Metasearch, Scorp Local Simhash & Vector RAG (0 external DB), ZeroClaw (ZeroClaw Labs, Rust), ZeroClaw Third-Party Search Stack (Tavily/Brave/Jina/SearXNG), CLI Slash Commands (/help /models /model /mode /session /sop /receipts /tools /cost /clear /exit), Cryptographic SHA-256 Tool Execution Receipts (+11 more)

### Community 46 - "RunAgentSessionLoop"
Cohesion: 0.23
Nodes (16): ClearStopRequest(), ConsumeStopRequest(), confirmKeyboard(), cleanToolCallTags(), getSessionSearchContext(), maxIterations(), maxTurnTimeout(), resumeAgentLoop() (+8 more)

### Community 47 - "main"
Cohesion: 0.20
Nodes (14): RegisterAutonomous(), hasDebugFlag(), StartGateway(), isCLIMode(), main(), InitModelUsage(), SOP, Dir() (+6 more)

### Community 48 - "GetIntArg"
Cohesion: 0.15
Nodes (10): GetBoolArg(), GetIntArg(), getUnclosedTags(), SplitMessage(), ExecuteCompose(), ExecuteGit(), ExecuteHTTP(), ExecuteLog() (+2 more)

### Community 49 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 50 - "TruncateStr"
Cohesion: 0.23
Nodes (11): confirmationDisplay(), TruncateStr(), CallOpenCode(), CallOpenCodeStream(), CallOpenCodeWithTools(), resolveOpenCodeBaseURL(), resolveOpenCodeKeyFromDisk(), resolveOpenCodeSessionID() (+3 more)

### Community 51 - "confirmation.go"
Cohesion: 0.26
Nodes (15): TestStorePendingConfirmationArgs(), clearPendingConfirmation(), clearPendingConfirmationFromDisk(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation() (+7 more)

### Community 52 - "os.File"
Cohesion: 0.21
Nodes (13): lockFileExclusive(), unlockFile(), lockFileExclusive(), unlockFile(), os.File, FileLockExclusive(), fileLockExclusive(), FileUnlock() (+5 more)

### Community 53 - "CheckDenyRules"
Cohesion: 0.20
Nodes (15): CheckDenyRules(), loadDenyRules(), ParseDenyRule(), ReloadDenyRules(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped(), TestCheckDenyRulesMatches() (+7 more)

### Community 54 - "SandboxActive"
Cohesion: 0.25
Nodes (15): caseSandbox(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest(), SandboxStatusNotice(), SandboxVersion() (+7 more)

### Community 55 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 56 - "CheckServerContracts"
Cohesion: 0.24
Nodes (14): caseMCPContract(), CheckServerContracts(), contractFile(), ContractStatusNotice(), ContractWarnings(), serverFingerprint(), SetContractPath(), contractTestServer() (+6 more)

### Community 57 - "ResetNativeToolCache"
Cohesion: 0.29
Nodes (13): ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), clearDynamicEnv(), registerTempTool(), TestDynamicToolTTL() (+5 more)

### Community 58 - "ConfigManager"
Cohesion: 0.23
Nodes (4): CM(), InitConfigManager(), NewConfigManager(), ConfigManager

### Community 59 - "LoadMCPConfig"
Cohesion: 0.28
Nodes (14): MCPConfigFilePath(), LoadMCPConfig(), rebuildMCPToolList(), ReloadMCPServers(), sanitizeMCPName(), StartMCPServers(), unregisterMCPNativeTools(), AddServerEntry() (+6 more)

### Community 60 - "net/http.Request"
Cohesion: 0.28
Nodes (12): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), TestGatewayEndpoints() (+4 more)

### Community 61 - "Transpiler Phase 3b: Contract Verify & Graceful Degradation"
Cohesion: 0.20
Nodes (14): Multi-Provider LLM Gateway (gateway/gateway.go + models.json), mark3labs/mcp-go SDK, mcp-fetch Reference Port (web scraping/markdown), mcp-filesystem Reference Port (scoped file access), mcp-sqlite Reference Port (local DB inspection), wahyuzero/scorp-mcp-registry (public registry repo), AI Transpiler (Self-Hosted Codegen), Transpiler Phase 3a: Sandbox Build (+6 more)

### Community 62 - "runLiveCase"
Cohesion: 0.22
Nodes (12): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, evalSandboxDir(), deploymentEnv(), LiveCases(), runLiveCase(), tail(), usageDelta() (+4 more)

### Community 63 - "cost_router.go"
Cohesion: 0.25
Nodes (13): defaultCostConfig(), formatCostReport(), FormatDailyCostSummary(), handleCostCommand(), init(), isBudgetExceeded(), isOffPeak(), LoadCostConfig() (+5 more)

### Community 64 - "TestIntegrityStatus"
Cohesion: 0.38
Nodes (13): IsTestRelatedPath(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestTestIntegrityStatus_FailingSuiteDoesNotCount(), TestTestIntegrityStatus_GreenRunMustBeAfterEdit(), TestTestIntegrityStatus_NoTouches() (+5 more)

### Community 65 - "session_ui.go"
Cohesion: 0.31
Nodes (11): FormatCompactStats(), CompactStats, BuildSessionMenuKeyboard(), FormatSessionMenuText(), GetActiveSessionID(), HandleSessionCallback(), init(), loadTgSessionMapping() (+3 more)

### Community 66 - "MCP Marketplace Build Plan & Execution Roadmap"
Cohesion: 0.21
Nodes (13): MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry, Security Layer 2: AST & Prompt Injection Scan, Security Layer 3: Hermetic CI/CD, ~/.scorp/mcp.json Central Config Entry, mcp_manage Agent Tool (mcp/manage.go), Scorp MCP Marketplace, MCP Marketplace Blueprint (docs/MCP_MARKETPLACE_BLUEPRINT.md) (+5 more)

### Community 67 - "eval.go"
Cohesion: 0.29
Nodes (9): Case, caseResult, CoreCases(), humanCount(), liveLabel(), report(), Run(), TestRunnerAggregationAndFilter() (+1 more)

### Community 68 - "ComputeSimhash"
Cohesion: 0.21
Nodes (12): tokenize(), ComputeSimhash(), hammingDistance(), simhashSimilarity(), TestSimHash_ComputeSimhash(), TestSimHash_HammingDistance(), TestSimHash_Similarity(), TestSimHash_SmartChunk() (+4 more)

### Community 69 - "init"
Cohesion: 0.29
Nodes (9): init(), GetAllTools(), countActiveTools(), countDeferredTools(), ExecuteToolCall(), ExecuteToolList(), ExecuteToolSearch(), callVisionModel() (+1 more)

### Community 70 - "🔍 Temuan 10 Kelemahan Arsitektural & Solusi Permanen"
Cohesion: 0.17
Nodes (11): 1. In-Memory Session Environment Non-Persistence across Independent CLI Invocations, 2. Truncation Hazard pada `write_file`, `replace_file_content`, & Bridge Protocol, 3. CRLF (`\r\n`) Line-Ending Mismatch pada Surgical Chunk Replacement, 4. 64-Bit Cryptographic Receipt Hash Collision Blindspot, 5. Infinite Traversal Hang pada Symlink Cycles & Inode Tracking, 6. Binary Non-UTF8 Stream Corruption di Shell Output Buffer, 7. Over-Inspection Loop pada Skenario Modifikasi File, Laporan Audit Kelemahan Arsitektural Non-Kecepatan & Penguatan Sistem (Fase 16 – 18) (+3 more)

### Community 71 - "time.Duration"
Cohesion: 0.35
Nodes (11): time.Duration, AddUptimeTarget(), checkTarget(), ExecuteUptime(), ListUptimeTargets(), RemoveUptimeTarget(), runUptimeCheck(), UptimeLoop() (+3 more)

### Community 72 - "EscapeHTML"
Cohesion: 0.27
Nodes (10): EscapeHTML(), gateScheduledShell(), isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage(), ShellTask() (+2 more)

### Community 73 - "NewDefaultMetaSearchAggregator"
Cohesion: 0.32
Nodes (8): deduplicateAndRank(), GetMetaSearchAggregator(), NewDefaultMetaSearchAggregator(), normalizeSearchURL(), TestNormalizeSearchURL(), MetaSearchAggregator, SearchEngine, SearchResult

### Community 74 - "patchReplace"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 75 - "🔍 Temuan Kelemahan & Solusi Arsitektural Permanen"
Cohesion: 0.18
Nodes (10): 1. Search Code Flag Injection & Catastrophic Traversal (`tools/search.go`), 2. Symlink Traversal Whitelist Escape Bypass (`tools/exec.go`), 3. Persistent Memory Concurrent Race & Truncation Hazard (`tools/memory.go`), 4. Unbounded Real-Time Steering Queue Flooding (`agent/steering.go`), 5. Telegram File Browser Path Mapping Memory Leak (`telegram/files.go`), 6. Ephemeral Pending Confirmation Loss on Restart (`agent/confirmation.go`), Laporan Audit Lanjutan Flaw Hunting & Penguatan Sistem (Fase 20), 📊 Ringkasan Hasil Pengujian Fase 20 pada VPS (+2 more)

### Community 76 - "🔍 Temuan Kelemahan & Solusi Arsitektural Permanen"
Cohesion: 0.18
Nodes (10): 1. Unbounded Concurrent Execution pada Scheduler (`scheduler/scheduler.go`), 2. Lazy Initialization Panic pada Vector RAG (`rag/rag_vector.go`), 3. Non-Atomic RAG Index Disk Persistence (`rag/rag.go` & `rag/rag_vector.go`), 4. SimHash Zero-Variance Whitespace Hallucination Pollution (`rag/rag_vector.go`), 5. Node Modules & Git Trash Folder Ingestion di RAG (`rag/rag_vector.go`), 6. Task Plan Atomic Serialization Under High Rate (`agent/taskplan.go`), Laporan Audit Lanjutan Flaw Hunting & Penguatan Sistem (Fase 21), 📊 Ringkasan Hasil Pengujian Fase 21 pada VPS (+2 more)

### Community 77 - "🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering)"
Cohesion: 0.18
Nodes (10): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Brutal, 🔍 3. Analisis Mendalam Kualitas Output & Arsitektur, 🏆 4. Kesimpulan Akhir Benchmark Brutal, A. Kompilasi Bahasa Go Multi-File (Kasus B2), B. Otomasi Server REST API Latar Belakang & Health Check (Kasus B7), 🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering), C. Fenomena Kebocoran Tag Simulasi pada ZeroClaw (Kasus B1, B3, B5, B8, B10, B12, B15) (+2 more)

### Community 78 - "clarify.go"
Cohesion: 0.27
Nodes (9): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+1 more)

### Community 79 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 80 - "tools/monitor.go"
Cohesion: 0.38
Nodes (10): ExecuteMonitor(), InitMonitor(), loadMonitorTargets(), monitorCheckOne(), monitorLoop(), ragIngestText(), sanitizeFilename(), saveMonitorTargets() (+2 more)

### Community 81 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 82 - "wireCLICallbacks"
Cohesion: 0.29
Nodes (7): formatFileSize(), wireCLICallbacks(), formatFinalResponse(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 83 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 84 - "ExecuteTool"
Cohesion: 0.44
Nodes (7): TestExecuteToolDenyRulesFirst(), registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks(), ExecuteTool()

### Community 85 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 86 - "IsDangerousCommand"
Cohesion: 0.33
Nodes (7): TestIsDangerousCommand(), devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), inspectBase64Payload(), IsDangerousCommand(), isHarmlessDevSink()

### Community 87 - "net/http.Client"
Cohesion: 0.36
Nodes (8): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport()

### Community 88 - "startMCPServer"
Cohesion: 0.33
Nodes (7): ProbeServer(), startMCPServer(), MCPServerConfig, MCPServer, isRemoteMCP(), startSSEServer(), TestRemoteMCPServer_HTTP()

### Community 89 - "GenerateNativeToolsSchema"
Cohesion: 0.33
Nodes (8): ChatRequest, TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions(), GenerateNativeToolsSchema(), ToolSchema

### Community 90 - "metasearch_test.go"
Cohesion: 0.31
Nodes (6): SearchResult, TestDeduplicateAndRank(), TestLiveWebSearch(), TestMetaSearchAggregator_ConcurrentSuccess(), TestMetaSearchAggregator_FaultTolerance(), MockSearchEngine

### Community 91 - "ExecuteReadURL"
Cohesion: 0.36
Nodes (7): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), TestReadURL_LocalMock(), truncateURLOutput(), tryRemoteScrape()

### Community 92 - "ConfigMgr"
Cohesion: 0.36
Nodes (8): LoadAutonomousConfig(), saveAutonomousConfigLocked(), SetKillSwitch(), setupTestPaths(), TestPhase7_ConfigPersistence(), TestPhase7_KillSwitch(), ConfigMgr(), saveCostTracker()

### Community 93 - "🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam"
Cohesion: 0.25
Nodes (7): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Ultra-Hardcore, 🔍 3. Temuan Kritis: Mengapa ZeroClaw Runtuh Total?, 💡 4. Analisis Komparasi Mendalam: Scorp vs PicoClaw, 🏆 5. Rekapitulasi Menyeluruh (Grand Total 65 Skenario Uji), 🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam, Scorp Agent vs PicoClaw vs ZeroClaw (Live Head-to-Head VPS Evaluation)

### Community 94 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 95 - "ExecuteAutonomous"
Cohesion: 0.52
Nodes (6): SaveAutonomousConfig(), autoShowActions(), autoShowConfig(), autoShowLog(), autoStatus(), ExecuteAutonomous()

### Community 96 - "🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut"
Cohesion: 0.29
Nodes (6): 🏛️ 15 Skenario Uji Ekstrem, 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Waktu Eksekusi 15 Skenario Ekstrem (Detik), 🔍 3. Analisis Mendalam Keandalan Arsitektur, 🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut, Scorp Agent vs PicoClaw vs ZeroClaw (Head-to-Head Live VPS Evaluation)

### Community 97 - "🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks)"
Cohesion: 0.29
Nodes (6): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Agentic Berat, 🔍 3. Kesimpulan Utama Uji Komparasi Baru, 🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks), ⚙️ Lingkungan & Konfigurasi Pengujian, Scorp Agent vs PicoClaw vs ZeroClaw (Fresh API Keys & Fast-Path Heuristic)

### Community 98 - "🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark)"
Cohesion: 0.29
Nodes (6): 📊 1. Ringkasan Hasil Eksekusi, ⏱️ 2. Tabel Waktu Respon Per Tugas (Detik), 🔍 3. Analisis Peningkatan Scorp, 🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark), ⚙️ Lingkungan & Konfigurasi Pengujian, Scorp Agent vs PicoClaw vs ZeroClaw (Fresh API Keys & Fast-Path Heuristic)

### Community 99 - "eval/core.go"
Cohesion: 0.43
Nodes (6): caseAutoClassifier(), caseDenyInvalidSkipped(), caseHooksBlockAndContext(), caseLedgerClear(), caseLedgerPersisted(), withEnv()

### Community 100 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 101 - "Release Workflow (tag-triggered cross-compile + GitHub release)"
Cohesion: 0.40
Nodes (6): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection

### Community 102 - "IsContinuationDirective"
Cohesion: 0.40
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 104 - "acquireSessionLock"
Cohesion: 0.40
Nodes (3): acquireSessionLock(), TestAcquireSessionLock(), sessionLockFile

### Community 105 - "CollectSystem"
Cohesion: 0.33
Nodes (5): CollectSystem(), GetTopProcesses(), StartCPUSampler(), SystemData, TopProcess

### Community 106 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 107 - "CLI Verify Deadloop (no timeout)"
Cohesion: 0.60
Nodes (6): CLI Verify Deadloop (no timeout), complete_task Explicit Tool Contract, Incident Report: Runaway CLI Verify Session (2026-09-05), Session Lock (cli_lock.go kernel advisory file lock), Heartbeat / Stall Detection, Agent Turn Timeout (context.WithTimeout, SCORP_MAX_TURN_TIMEOUT)

### Community 108 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 109 - "registerMCPToolsAsNative"
Cohesion: 0.33
Nodes (5): TruncOutputTool(), buildArgDefsFromInputSchema(), MCPToolsDeferred(), registerMCPToolsAsNative(), TestMCPToolsDeferredEnvParsing()

### Community 110 - "RestartServer"
Cohesion: 0.40
Nodes (5): StopMCPServers(), GetServerHealthStatus(), RegisterWatchdog(), RestartServer(), StopWatchdogs()

### Community 111 - "deploy.sh"
Cohesion: 0.53
Nodes (4): die(), ok(), deploy.sh script, step()

### Community 112 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 113 - "Laporan Eksekusi Komprehensif Uji End-to-End (E2E): Fase 1 – 15"
Cohesion: 0.40
Nodes (4): 🔍 Analisis Temuan & Ketiadaan Isu, Laporan Eksekusi Komprehensif Uji End-to-End (E2E): Fase 1 – 15, 🎯 Ringkasan Eksekusi Suite, 📊 Tabel Hasil Uji Sekuensial Fase 1 – 15

### Community 115 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 116 - "os.FileMode"
Cohesion: 0.67
Nodes (3): SaveJSON(), SaveJSONPerm(), os.FileMode

### Community 117 - "RegisterPlugin"
Cohesion: 0.83
Nodes (3): RegisterPlugin(), ToolPlugin, ToolPluginWithSchema

### Community 120 - "ExecuteSQL"
Cohesion: 0.83
Nodes (3): ExecuteSQL(), loadDBConnections(), dbConnection

## Knowledge Gaps
- **113 isolated node(s):** `containerStats`, `🔍 Analisis Temuan & Ketiadaan Isu`, `🎯 Ringkasan Eksekusi Suite`, `📊 Tabel Hasil Uji Sekuensial Fase 1 – 15`, `hookPayload` (+108 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 285 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **12 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `collector_system.go`, `HandleModelCallback`, `TaskPlan`, `getAgentSystemPrompt`, `StartDaemon`, `GetAutonomyLevel`, `collector_security.go`, `setSession`, `CreateCheckpoint`, `scheduler.go`, `getSession`, `RunAgentSessionLoop`, `main`, `confirmation.go`, `SandboxActive`, `CheckServerContracts`, `session_ui.go`, `EscapeHTML`, `clarify.go`, `CollectSystem`, `RestartServer`?**
  _High betweenness centrality (0.110) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `TaskPlan`, `getAgentSystemPrompt`, `compaction_test.go`, `testgate.go`, `GetAutonomyLevel`, `GetStringArg`, `setSession`, `HandleTelegramAction`, `PermissionDecision`, `CreateCheckpoint`, `chat.go`, `startCLI`, `CallModelWithFallback`, `getSession`, `TruncateStr`, `confirmation.go`, `ResetNativeToolCache`, `TestIntegrityStatus`, `ComputeSimhash`, `EscapeHTML`, `wireCLICallbacks`, `ExecuteTool`, `IsDangerousCommand`, `ExecuteTermuxAPI`, `IsContinuationDirective`, `GenerateContextualSessionTitle`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Why does `TruncateStr()` connect `TruncateStr` to `runSubagent`, `HandleModelCallback`, `TaskPlan`, `context.Context`, `getAgentSystemPrompt`, `ResolveAPIKey`, `setSession`, `ModelConfig`, `PermissionDecision`, `agent/autonomous.go`, `chat.go`, `ToolCall`, `scheduler.go`, `MCPServer`, `api_gemini.go`, `getSession`, `TruncOutput`, `RunAgentSessionLoop`, `GetIntArg`, `EscapeHTML`, `ExecuteAutonomous`, `registerMCPToolsAsNative`?**
  _High betweenness centrality (0.055) - this node is a cross-community bridge._
- **Are the 25 inferred relationships involving `HandleTelegramAction()` (e.g. with `CollectSystem()` and `DirKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 25 INFERRED edges - model-reasoned connections that need verification._
- **What connects `containerStats`, `🔍 Analisis Temuan & Ketiadaan Isu`, `🎯 Ringkasan Eksekusi Suite` to the rest of the system?**
  _113 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.05714285714285714 - nodes in this community are weakly interconnected._
- **Should `collector_system.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05328005328005328 - nodes in this community are weakly interconnected._