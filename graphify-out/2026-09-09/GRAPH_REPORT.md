# Graph Report - scorp  (2026-09-09)

## Corpus Check
- 282 files · ~204,531 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2224 nodes · 5739 edges · 123 communities (110 shown, 5 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 913 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ec344bb9`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- chat.go
- cost_router.go
- runSubagent
- metasearch_engines.go
- files.go
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
- wizard.go
- ExecuteSQL
- TestIntegrityStatus
- CheckServerContracts
- client.go
- skills.go
- config/hooks.go
- TruncOutput
- PermissionDecision
- CallAnthropicWithTools
- checker.go
- HandleConfirmation
- SCORP — BRUTAL END-TO-END TEST PLAN
- GetStringArg
- api_gemini.go
- SandboxActive
- HandleTelegramAction
- ResolveAPIKey
- 🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering)
- GetIntArg
- CheckDenyRules
- eval.go
- EscapeHTML
- startCLI
- TruncateStr
- patch.go
- runLiveCase
- SaveModelConfig
- RegisterProvider
- ResetNativeToolCache
- prepareNewTurnHistory
- net/http.Request
- TodoManager
- 🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam
- CredentialVault
- ExecuteReadURL
- testgate.go
- HandleModelCallback
- IsDangerousCommand
- os.File
- sync.Mutex
- v2_skills.go
- 🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks)
- IsContinuationDirective
- init
- .listenSSEStream
- install.sh
- deploy.sh
- RenameSession
- TestPhase6_AllTools
- registry/registry.go
- GenerateContextualSessionTitle
- ScorpDir
- handleTelegramMarketplaceInstall
- InitDefaultSOPs
- 🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark)
- 🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut
- registerMCPToolsAsNative
- wireCLICallbacks
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- main
- Transpiler Phase 3b: Contract Verify & Graceful Degradation
- handleMCPCommand
- MCP Marketplace Build Plan & Execution Roadmap
- model_catalog.go
- registerHookProbeTool
- LoadConfig
- inline.go
- ExecuteScript
- renderStatusFooter
- ExecuteTermuxAPI
- TransformToolDefinitions
- ExecuteAutoLogin
- serviceBridgeRequests
- startMCPServer
- clarify.go
- CleanupSessionsLoop
- RedactSecrets
- acquireSessionLock
- upload.go
- Release Workflow (tag-triggered cross-compile + GitHub release)
- CLI Verify Deadloop (no timeout)
- TestToolSearchActivatesDeferredToolInStaticMode
- CompactSessionHistory

## God Nodes (most connected - your core abstractions)
1. `ModelConfig` - 89 edges
2. `HandleTelegramAction()` - 80 edges
3. `ChatMessage` - 75 edges
4. `RunAgentSessionLoop()` - 64 edges
5. `startCLI()` - 52 edges
6. `TruncateStr()` - 47 edges
7. `GetStringArg()` - 47 edges
8. `StartDaemon()` - 44 edges
9. `ScorpPath()` - 42 edges
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

## Communities (123 total, 5 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (70): Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo, RegistryIndex (+62 more)

### Community 1 - "chat.go"
Cohesion: 0.15
Nodes (30): appendSessionHistory(), ClearStopRequest(), EnterAgentMode(), ExitAgentMode(), extractAndSaveMemory(), flushPendingMessages(), GetHistoryTokenEstimate(), getOrCreateSession() (+22 more)

### Community 2 - "cost_router.go"
Cohesion: 0.06
Nodes (62): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+54 more)

### Community 3 - "runSubagent"
Cohesion: 0.05
Nodes (59): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+51 more)

### Community 4 - "metasearch_engines.go"
Cohesion: 0.05
Nodes (52): FormatDuration(), net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost(), getClient() (+44 more)

### Community 5 - "files.go"
Cohesion: 0.28
Nodes (15): BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), HumanSize(), PathID(), RootsKeyboard() (+7 more)

### Community 6 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 7 - "TaskPlan"
Cohesion: 0.09
Nodes (42): PlanItem, ApprovePlan(), BeginPlanning(), CancelPlan(), EndPlanning(), AgentMessage, PlanningState(), RevisePlan() (+34 more)

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
Nodes (18): StopMCPServerMode(), Init(), StartServer(), StopServer(), runCommandLoop(), StartDaemon(), DeleteWebhook(), InitTelegram() (+10 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.06
Nodes (41): getSharedMemorySummary(), AppendMemoryMD(), extractTaskMemory(), AgentMessage, parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile() (+33 more)

### Community 13 - "time.Time"
Cohesion: 0.10
Nodes (31): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, time.Time, ScheduledTask, AddTask(), AddTaskEx() (+23 more)

### Community 14 - "🔍 Rincian Detail Uji per Fase"
Cohesion: 0.10
Nodes (19): 1️⃣1️⃣ Fase 11: Recursive Parent Lock & Deep Traversal, 1️⃣2️⃣ Fase 12: Synthetic Feedback & Unicode RTL/Emoji Edge Cases, 1️⃣3️⃣ Fase 13: Output Buffer Flooding & EPIPE Signals, 1️⃣4️⃣ Fase 14: Session Environment Persistence & Error Reporting, 1️⃣5️⃣ Fase 15: Subshell Mutation Isolation & Polyglot Compilation, 1️⃣ Fase 1: Tool Parsing & Thought Signature Contract, 2️⃣ Fase 2: Filesystem Whitelist & Bare-Metal Sandbox Toggle, 3️⃣ Fase 3: Subprocess Pipe Deadlock & Asynchronous Draining (+11 more)

### Community 15 - "ToolCall"
Cohesion: 0.14
Nodes (20): autoClassifyWithModel(), TestParseToolCalls(), extractFallbackToolCalls(), getCheapestModel(), RouteModelCostAware(), CostTracker, CallModelWithFallback(), findFirstVisionModel() (+12 more)

### Community 16 - "testing.T"
Cohesion: 0.07
Nodes (39): TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations() (+31 more)

### Community 17 - "ScorpPath"
Cohesion: 0.13
Nodes (24): init(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), HomeDir(), Hostname() (+16 more)

### Community 18 - "ExecuteTool"
Cohesion: 0.13
Nodes (18): TestAutoAllowlistPrefixAnchoring(), TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), ResetAutoAllowlist(), ResetAutoStats(), TestExecuteToolAutoPresetTrustedAndRecorded(), ExecuteTool(), init() (+10 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.13
Nodes (27): AgentMessage, confirmationDisplay(), ConsumeStopRequest(), sendScorpReply(), confirmKeyboard(), cleanToolCallTags(), getSessionSearchContext(), maxIterations() (+19 more)

### Community 20 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 21 - "compaction_test.go"
Cohesion: 0.12
Nodes (37): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, injectPreservationNote(), latestUserGoal(), maybeCompactHistory(), preservationNote() (+29 more)

### Community 22 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 23 - "context.Context"
Cohesion: 0.08
Nodes (19): context.Context, AnthropicProvider, CallCommandCode(), resolveOpenCodeKeyFromDisk(), AzureProvider, ChatMessage, ChatResponse, CohereProvider (+11 more)

### Community 24 - "GetAutonomyLevel"
Cohesion: 0.14
Nodes (24): ConfirmationRequired(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels() (+16 more)

### Community 25 - "collector_system.go"
Cohesion: 0.05
Nodes (73): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+65 more)

### Community 26 - "RunPreToolUseHooks"
Cohesion: 0.21
Nodes (20): PreToolUse & PostToolUse Hooks (P3.12), Community Praised Agent Patterns (2026), 2026 AI Coding Agent Competitor Landscape, hookPayload, appendHookContext(), runHookCommand(), RunPostToolUseHooks(), RunPreToolUseHooks() (+12 more)

### Community 27 - "MCPServer"
Cohesion: 0.20
Nodes (6): bufio.Scanner, encoding/json.Encoder, sync.Once, MCPServer, StopMCPServers(), StopWatchdogs()

### Community 28 - "api_commandcode.go"
Cohesion: 0.19
Nodes (12): buildCommandCodePayload(), CallCommandCodeStream(), createCommandCodeRequest(), resolveCommandCodeKeyFromDisk(), commandCodeMsg, commandCodeParams, commandCodePayload, CommandCodeProvider (+4 more)

### Community 29 - "CreateCheckpoint"
Cohesion: 0.26
Nodes (19): Checkpoint and Rewind (P1.6), caseCheckpoint(), CheckpointDiffStat(), checkpointRepoRoot(), ckptGit(), ckptRefFor(), CreateCheckpoint(), DeleteCheckpoint() (+11 more)

### Community 30 - "wizard.go"
Cohesion: 0.31
Nodes (17): AutoPopulateFromCatalog(), ProviderKeyEnv(), modelWizard, askAPIKey(), ClearModelWizard(), EnvFilePath(), finalizeModelKeySave(), finalizeProviderKeySave() (+9 more)

### Community 31 - "ExecuteSQL"
Cohesion: 0.83
Nodes (3): ExecuteSQL(), loadDBConnections(), dbConnection

### Community 32 - "TestIntegrityStatus"
Cohesion: 0.32
Nodes (15): IsTestRelatedPath(), IsTestRunCommand(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestIsTestRunCommand(), TestTestIntegrityStatus_FailingSuiteDoesNotCount() (+7 more)

### Community 33 - "CheckServerContracts"
Cohesion: 0.24
Nodes (14): caseMCPContract(), CheckServerContracts(), contractFile(), ContractStatusNotice(), ContractWarnings(), serverFingerprint(), SetContractPath(), contractTestServer() (+6 more)

### Community 34 - "client.go"
Cohesion: 0.17
Nodes (20): ACPRequest, encoding/json.RawMessage, FindMCPTool(), getExposedTools(), GetMCPTools(), MCPTool, handleMCPRequest(), MCPToolsForPrompt() (+12 more)

### Community 35 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 36 - "config/hooks.go"
Cohesion: 0.30
Nodes (13): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+5 more)

### Community 37 - "TruncOutput"
Cohesion: 0.28
Nodes (15): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+7 more)

### Community 38 - "PermissionDecision"
Cohesion: 0.26
Nodes (14): autoAllowlisted(), autoClassify(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand(), PermissionDecision(), setAutoMode(), TestAutoFallbackAfterRepeatedUncertain() (+6 more)

### Community 39 - "CallAnthropicWithTools"
Cohesion: 0.23
Nodes (12): anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic(), callAnthropicStream(), CallAnthropicWithTools() (+4 more)

### Community 40 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 41 - "HandleConfirmation"
Cohesion: 0.36
Nodes (10): TestStorePendingConfirmationArgs(), clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), StorePendingConfirmation() (+2 more)

### Community 42 - "SCORP — BRUTAL END-TO-END TEST PLAN"
Cohesion: 0.09
Nodes (20): 0. SCALE & SEVERITY DEFINITIONS, A. GATE-STACK ADVERSARIAL (S) — Layer Penetration Probes, B. LONG-HORIZON & CONTEXT COMPACTION (L), C. INFRASTRUCTURE CHAOS (M/L), D. CONCURRENCY & RACE CONDITIONS (M), E. TELEGRAM UX INTEGRITY (M), F. ADVERSARIAL SECURITY (M), G. EVALUATION & DEPLOYMENT GATE INTEGRITY (S) (+12 more)

### Community 43 - "GetStringArg"
Cohesion: 0.14
Nodes (18): init(), GetStringArg(), captureSessionExports(), cleanAbortedShm(), decodeUTF16(), decodeUTFString(), ExecuteListDir(), ExecuteReadFile() (+10 more)

### Community 44 - "api_gemini.go"
Cohesion: 0.17
Nodes (20): callGemini(), callGeminiStream(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), parseDataURL(), resolveGeminiBaseURL() (+12 more)

### Community 45 - "SandboxActive"
Cohesion: 0.25
Nodes (15): caseSandbox(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest(), SandboxStatusNotice(), SandboxVersion() (+7 more)

### Community 46 - "HandleTelegramAction"
Cohesion: 0.20
Nodes (18): RequestStop(), HandleTelegramAction(), GetPath(), BackAndRefreshKeyboard(), baseName(), EditMessage(), EditMessageByID(), MainMenuKeyboard() (+10 more)

### Community 47 - "ResolveAPIKey"
Cohesion: 0.16
Nodes (21): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestGemini_MultimodalParsing(), TestOpenCodeProvider_RegistrationAndKey(), CallModel(), CallModelStream(), GetProvider() (+13 more)

### Community 48 - "🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering)"
Cohesion: 0.18
Nodes (10): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Brutal, 🔍 3. Analisis Mendalam Kualitas Output & Arsitektur, 🏆 4. Kesimpulan Akhir Benchmark Brutal, A. Kompilasi Bahasa Go Multi-File (Kasus B2), B. Otomasi Server REST API Latar Belakang & Health Check (Kasus B7), 🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering), C. Fenomena Kebocoran Tag Simulasi pada ZeroClaw (Kasus B1, B3, B5, B8, B10, B12, B15) (+2 more)

### Community 49 - "GetIntArg"
Cohesion: 0.15
Nodes (10): GetBoolArg(), GetIntArg(), getUnclosedTags(), SplitMessage(), ExecuteCompose(), ExecuteGit(), ExecuteHTTP(), ExecuteLog() (+2 more)

### Community 50 - "CheckDenyRules"
Cohesion: 0.17
Nodes (16): TestExecuteToolDenyRulesFirst(), CheckDenyRules(), loadDenyRules(), ParseDenyRule(), ReloadDenyRules(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped() (+8 more)

### Community 51 - "eval.go"
Cohesion: 0.22
Nodes (12): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, Case, caseResult, CoreCases(), evalSandboxDir(), humanCount(), liveLabel() (+4 more)

### Community 52 - "EscapeHTML"
Cohesion: 0.21
Nodes (11): collectTableLines(), convertInlineMarkdown(), convertTableToList(), TestMarkdownToTelegramHTMLEscaping(), isSeparatorRow(), markdownToTelegramHTML(), parseTableRow(), safeIndex() (+3 more)

### Community 53 - "startCLI"
Cohesion: 0.27
Nodes (18): executeOneShot(), executeTurn(), handleCLISession(), handleCLISOP(), hasDebugFlag(), printBanner(), printCLIHelp(), printCostUsage() (+10 more)

### Community 54 - "TruncateStr"
Cohesion: 0.25
Nodes (23): TruncateStr(), applyAzureHeaders(), callAzure(), callAzureStream(), callAzureWithTools(), resolveAzureEndpoint(), CallCommandCodeWithTools(), applyOpenAIHeaders() (+15 more)

### Community 55 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "runLiveCase"
Cohesion: 0.31
Nodes (9): deploymentEnv(), LiveCases(), runLiveCase(), tail(), usageDelta(), usageSnapshot(), liveCase, TestUsageDeltaClampsNegative() (+1 more)

### Community 57 - "SaveModelConfig"
Cohesion: 0.21
Nodes (14): defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, ModelRouterConfig, ModelUsage, getProviderInfo(), HandleProviderCommand() (+6 more)

### Community 58 - "RegisterProvider"
Cohesion: 0.11
Nodes (20): init(), init(), init(), init(), init(), init(), init(), init() (+12 more)

### Community 59 - "ResetNativeToolCache"
Cohesion: 0.29
Nodes (13): ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), clearDynamicEnv(), registerTempTool(), TestDynamicToolTTL() (+5 more)

### Community 60 - "prepareNewTurnHistory"
Cohesion: 0.29
Nodes (8): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary()

### Community 61 - "net/http.Request"
Cohesion: 0.25
Nodes (13): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), TestGatewayEndpoints() (+5 more)

### Community 62 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 63 - "🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam"
Cohesion: 0.25
Nodes (7): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Ultra-Hardcore, 🔍 3. Temuan Kritis: Mengapa ZeroClaw Runtuh Total?, 💡 4. Analisis Komparasi Mendalam: Scorp vs PicoClaw, 🏆 5. Rekapitulasi Menyeluruh (Grand Total 65 Skenario Uji), 🌋 Benchmark Komparasi Ultra-Hardcore: 15 Tugas Tingkat Kernel, C Kripto & Rekayasa Arsitektur Mendalam, Scorp Agent vs PicoClaw vs ZeroClaw (Live Head-to-Head VPS Evaluation)

### Community 64 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 65 - "ExecuteReadURL"
Cohesion: 0.36
Nodes (7): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), TestReadURL_LocalMock(), truncateURLOutput(), tryRemoteScrape()

### Community 66 - "testgate.go"
Cohesion: 0.14
Nodes (22): OperationalClaim, GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), dedupeOpObjects(), extractOpObjects() (+14 more)

### Community 67 - "HandleModelCallback"
Cohesion: 0.20
Nodes (16): ApiFormatPickerKeyboard(), FallbackEditorKeyboard(), ModelInfoKeyboard(), ModelMenuKeyboard(), ModelPickerKeyboard(), addToFallback(), FallbackText(), FormatAPIKeysList() (+8 more)

### Community 68 - "IsDangerousCommand"
Cohesion: 0.33
Nodes (7): TestIsDangerousCommand(), devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), inspectBase64Payload(), IsDangerousCommand(), isHarmlessDevSink()

### Community 69 - "os.File"
Cohesion: 0.21
Nodes (9): lockFileExclusive(), unlockFile(), lockFileExclusive(), unlockFile(), os.File, fileLockExclusive(), fileUnlock(), fileLockExclusive() (+1 more)

### Community 70 - "sync.Mutex"
Cohesion: 0.50
Nodes (4): autoStats, sync.Mutex, getChatLock(), StartTestEndpoint()

### Community 71 - "v2_skills.go"
Cohesion: 0.27
Nodes (9): SkillMeta, ActivateSkill(), ListSkillsOverview(), LoadAllSkills(), ParseSkillMetadata(), ReadSkillBody(), scanLegacyJSONSkills(), scanSkillsDirectory() (+1 more)

### Community 72 - "🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks)"
Cohesion: 0.29
Nodes (6): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Agentic Berat, 🔍 3. Kesimpulan Utama Uji Komparasi Baru, 🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks), ⚙️ Lingkungan & Konfigurasi Pengujian, Scorp Agent vs PicoClaw vs ZeroClaw (Fresh API Keys & Fast-Path Heuristic)

### Community 73 - "IsContinuationDirective"
Cohesion: 0.40
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 75 - "init"
Cohesion: 0.26
Nodes (10): init(), unregisterMCPNativeTools(), GetAllTools(), countActiveTools(), countDeferredTools(), ExecuteToolCall(), ExecuteToolList(), ExecuteToolSearch() (+2 more)

### Community 76 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 78 - "deploy.sh"
Cohesion: 0.53
Nodes (4): die(), ok(), deploy.sh script, step()

### Community 79 - "RenameSession"
Cohesion: 0.18
Nodes (21): ClearChatSession(), FlushSessionHistory(), historyFilePath(), saveHistoryToDisk(), DeleteSession(), ListSessions(), RenameSession(), sanitizeSessionID() (+13 more)

### Community 80 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 81 - "registry/registry.go"
Cohesion: 0.15
Nodes (13): executeMCPServerTool(), echoPlugin, RegisterPlugin(), TestRegisterPlugin(), ExecuteToolByName(), GenerateSystemPromptDescriptions(), GetTool(), GetToolsByCategory() (+5 more)

### Community 82 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 83 - "ScorpDir"
Cohesion: 0.25
Nodes (9): apiKey, ScorpDir(), TestConfigPaths_ScorpDir(), TestConfigPaths_ScreenshotsDir(), envKeyName, installSystemdService(), maskString(), RunQuickstart() (+1 more)

### Community 84 - "handleTelegramMarketplaceInstall"
Cohesion: 0.60
Nodes (5): handleTelegramMarketplaceInstall(), MarketplaceTriOptionKeyboard(), reportInstallOutcome(), BackButtonKeyboard(), SendMessage()

### Community 85 - "InitDefaultSOPs"
Cohesion: 0.42
Nodes (8): SOP, Dir(), GetSOP(), InitDefaultSOPs(), ListSOPs(), SaveSOP(), TestSOPLifecycle(), ExecuteSOP()

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

### Community 99 - "main"
Cohesion: 0.20
Nodes (8): RegisterAutonomous(), StartGateway(), isCLIMode(), main(), FormatModelListWithHealth(), FormatUsageStats(), InitModelUsage(), SwitchModel()

### Community 100 - "Transpiler Phase 3b: Contract Verify & Graceful Degradation"
Cohesion: 0.20
Nodes (14): Multi-Provider LLM Gateway (gateway/gateway.go + models.json), mark3labs/mcp-go SDK, mcp-fetch Reference Port (web scraping/markdown), mcp-filesystem Reference Port (scoped file access), mcp-sqlite Reference Port (local DB inspection), wahyuzero/scorp-mcp-registry (public registry repo), AI Transpiler (Self-Hosted Codegen), Transpiler Phase 3a: Sandbox Build (+6 more)

### Community 101 - "handleMCPCommand"
Cohesion: 0.31
Nodes (7): handleMCPCommand(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste()

### Community 102 - "MCP Marketplace Build Plan & Execution Roadmap"
Cohesion: 0.21
Nodes (13): MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry, Security Layer 2: AST & Prompt Injection Scan, Security Layer 3: Hermetic CI/CD, ~/.scorp/mcp.json Central Config Entry, mcp_manage Agent Tool (mcp/manage.go), Scorp MCP Marketplace, MCP Marketplace Blueprint (docs/MCP_MARKETPLACE_BLUEPRINT.md) (+5 more)

### Community 103 - "model_catalog.go"
Cohesion: 0.31
Nodes (8): CatalogEntry, CatalogModels(), HasCatalog(), ProviderHasAPIKey(), RemoveProviderModels(), ProviderListKeyboard(), formatProvidersList(), AllProviderNames()

### Community 104 - "registerHookProbeTool"
Cohesion: 0.73
Nodes (5): registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks()

### Community 105 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 106 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 107 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 108 - "renderStatusFooter"
Cohesion: 0.53
Nodes (5): GetDailyTotalUSD(), getContextPill(), getGitStatus(), getShortCwd(), renderStatusFooter()

### Community 109 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 110 - "TransformToolDefinitions"
Cohesion: 0.39
Nodes (7): ChatRequest, TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions(), ToolSchema

### Community 112 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 113 - "startMCPServer"
Cohesion: 0.33
Nodes (7): ProbeServer(), startMCPServer(), MCPServerConfig, MCPServer, isRemoteMCP(), startSSEServer(), TestRemoteMCPServer_HTTP()

### Community 114 - "clarify.go"
Cohesion: 0.24
Nodes (10): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+2 more)

### Community 115 - "CleanupSessionsLoop"
Cohesion: 0.33
Nodes (5): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), cleanupAgentSessions(), TGDocument

### Community 117 - "acquireSessionLock"
Cohesion: 0.40
Nodes (3): acquireSessionLock(), TestAcquireSessionLock(), sessionLockFile

### Community 118 - "upload.go"
Cohesion: 0.50
Nodes (4): contentPart, imageURL, TestBase64Encode(), base64Encode()

### Community 119 - "Release Workflow (tag-triggered cross-compile + GitHub release)"
Cohesion: 0.40
Nodes (6): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection

### Community 120 - "CLI Verify Deadloop (no timeout)"
Cohesion: 0.60
Nodes (6): CLI Verify Deadloop (no timeout), complete_task Explicit Tool Contract, Incident Report: Runaway CLI Verify Session (2026-09-05), Session Lock (cli_lock.go kernel advisory file lock), Heartbeat / Stall Detection, Agent Turn Timeout (context.WithTimeout, SCORP_MAX_TURN_TIMEOUT)

### Community 121 - "TestToolSearchActivatesDeferredToolInStaticMode"
Cohesion: 0.70
Nodes (4): deferredToolActive(), osUnsetDynamic(), registerSearchTarget(), TestToolSearchActivatesDeferredToolInStaticMode()

### Community 122 - "CompactSessionHistory"
Cohesion: 0.83
Nodes (3): CompactSessionHistory(), FormatCompactStats(), CompactStats

## Knowledge Gaps
- **85 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+80 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 238 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `files.go`, `TaskPlan`, `StartDaemon`, `time.Time`, `RunAgentSessionLoop`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`, `wizard.go`, `CheckServerContracts`, `HandleConfirmation`, `SandboxActive`, `HandleModelCallback`, `v2_skills.go`, `RenameSession`, `handleTelegramMarketplaceInstall`, `InitDefaultSOPs`, `registerMCPToolsAsNative`, `main`, `clarify.go`, `CompactSessionHistory`?**
  _High betweenness centrality (0.086) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `rag_vector.go`, `TaskPlan`, `getAgentSystemPrompt`, `ToolCall`, `ExecuteTool`, `compaction_test.go`, `GetAutonomyLevel`, `CreateCheckpoint`, `TestIntegrityStatus`, `PermissionDecision`, `HandleConfirmation`, `GetStringArg`, `HandleTelegramAction`, `EscapeHTML`, `startCLI`, `TruncateStr`, `ResetNativeToolCache`, `prepareNewTurnHistory`, `testgate.go`, `IsDangerousCommand`, `IsContinuationDirective`, `RenameSession`, `GenerateContextualSessionTitle`, `wireCLICallbacks`, `ExecuteTermuxAPI`?**
  _High betweenness centrality (0.083) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `cost_router.go`, `files.go`, `rag_vector.go`, `LoadMCPConfig`, `getAgentSystemPrompt`, `time.Time`, `ScorpPath`, `ExecuteTool`, `collector_system.go`, `MCPServer`, `wizard.go`, `client.go`, `skills.go`, `HandleConfirmation`, `HandleTelegramAction`, `SaveModelConfig`, `IsDangerousCommand`, `sync.Mutex`, `v2_skills.go`, `handleTelegramMarketplaceInstall`, `InitDefaultSOPs`, `main`, `CleanupSessionsLoop`?**
  _High betweenness centrality (0.063) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _85 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.06044303797468355 - nodes in this community are weakly interconnected._
- **Should `chat.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14717741935483872 - nodes in this community are weakly interconnected._