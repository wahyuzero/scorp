# Graph Report - scorp  (2026-09-09)

## Corpus Check
- 279 files · ~200,140 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2191 nodes · 5695 edges · 119 communities (104 shown, 7 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 912 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4ea9635b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- chat.go
- cost_router.go
- execute.go
- metasearch_engines.go
- HomeDir
- rag_vector.go
- ClearTaskPlan
- Scorp Agent (Go, ultra-light autonomous agent)
- LoadMCPConfig
- Benchmark
- StartDaemon
- getAgentSystemPrompt
- time.Time
- acp.go
- CallModelWithFallback
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
- OpenCodeProvider
- ExecuteSQL
- TestIntegrityStatus
- CheckServerContracts
- client.go
- skills.go
- config/hooks.go
- TruncOutput
- PermissionDecision
- runSubagent
- checker.go
- HandleConfirmation
- SCORP — BRUTAL END-TO-END TEST PLAN
- GetStringArg
- api_gemini.go
- SandboxActive
- telegram.go
- ResolveAPIKey
- 🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering)
- GetIntArg
- eval/core.go
- eval.go
- runPlanningTurns
- startCLI
- TruncateStr
- patch.go
- runLiveCase
- SetTaskPlan
- RegisterProvider
- bg.go
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
- HandleTelegramAction
- TestPhase6_AllTools
- registry/registry.go
- GenerateContextualSessionTitle
- HasGreenTestRun
- handleTelegramMarketplaceInstall
- main
- 🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark)
- 🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut
- registerMCPToolsAsNative
- wireCLICallbacks
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- TaskPlan
- Transpiler Phase 3b: Contract Verify & Graceful Degradation
- readInteractiveInput
- MCP Marketplace Build Plan & Execution Roadmap
- TestSteeringQueue
- ExecuteTool
- LoadConfig
- os/exec.Cmd
- ExecuteScript
- renderStatusFooter
- ExecuteTermuxAPI
- TransformToolDefinitions
- ExecuteAutoLogin
- serviceBridgeRequests
- Auto-Mode Classifier (P3.13)
- clarify.go
- ShareToMarketplace
- RedactSecrets
- Release Workflow (tag-triggered cross-compile + GitHub release)
- CLI Verify Deadloop (no timeout)

## God Nodes (most connected - your core abstractions)
1. `ModelConfig` - 89 edges
2. `HandleTelegramAction()` - 80 edges
3. `ChatMessage` - 75 edges
4. `RunAgentSessionLoop()` - 63 edges
5. `startCLI()` - 51 edges
6. `TruncateStr()` - 47 edges
7. `GetStringArg()` - 47 edges
8. `StartDaemon()` - 44 edges
9. `ScorpPath()` - 42 edges
10. `resumeAgentLoop()` - 40 edges

## Surprising Connections (you probably didn't know these)
- `Durable Memory MEMORY.md (P1.7)` --implements--> `extractTaskMemory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/memory_md.go
- `Deny-Rule Engine (P0.2)` --implements--> `CheckDenyRules()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → config/deny.go
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

## Communities (119 total, 7 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (70): Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo, RegistryIndex (+62 more)

### Community 1 - "chat.go"
Cohesion: 0.07
Nodes (61): agentSession, appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown(), convertTableToList() (+53 more)

### Community 2 - "cost_router.go"
Cohesion: 0.07
Nodes (56): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+48 more)

### Community 3 - "execute.go"
Cohesion: 0.19
Nodes (18): runOpenCodeCLI(), AgentMessage, TestParseDelegateParams(), TestParseDelegateParamsCapsAndDefaults(), delegateBatchParams, delegateResult, delegateTaskParams, DefaultSubagentTools() (+10 more)

### Community 4 - "metasearch_engines.go"
Cohesion: 0.05
Nodes (51): FormatDuration(), net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost(), getClient() (+43 more)

### Community 5 - "HomeDir"
Cohesion: 0.23
Nodes (18): HomeDir(), PythonSitePackages(), UploadsDir(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo() (+10 more)

### Community 6 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 7 - "ClearTaskPlan"
Cohesion: 0.32
Nodes (12): ApprovePlan(), CancelPlan(), TestApprovePlanWithoutLedger(), TestCancelPlanDropsLedger(), ClearTaskPlan(), execTaskPlanTool(), GetTaskPlan(), TestTaskPlanConcurrentAccess() (+4 more)

### Community 8 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.16
Nodes (19): Scorp vs PicoClaw vs ZeroClaw Architecture Comparison, ZeroClaw Defense-in-Depth Sandboxing, Scorp Embedded Native Multi-Engine Metasearch, Scorp Local Simhash & Vector RAG (0 external DB), ZeroClaw (ZeroClaw Labs, Rust), ZeroClaw Third-Party Search Stack (Tavily/Brave/Jina/SearXNG), CLI Slash Commands (/help /models /model /mode /session /sop /receipts /tools /cost /clear /exit), Cryptographic SHA-256 Tool Execution Receipts (+11 more)

### Community 9 - "LoadMCPConfig"
Cohesion: 0.38
Nodes (11): MCPConfigFilePath(), LoadMCPConfig(), ReloadMCPServers(), sanitizeMCPName(), AddServerEntry(), ExecuteMCPManage(), mcpManageAdd(), mcpManageList() (+3 more)

### Community 10 - "Benchmark"
Cohesion: 0.10
Nodes (37): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+29 more)

### Community 11 - "StartDaemon"
Cohesion: 0.25
Nodes (9): Init(), StartServer(), StopServer(), runCommandLoop(), StartDaemon(), InitTelegram(), SendMessageGetID(), SendMessageGetIDWithKeyboard() (+1 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.07
Nodes (39): getSharedMemorySummary(), AppendMemoryMD(), extractTaskMemory(), AgentMessage, parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile() (+31 more)

### Community 13 - "time.Time"
Cohesion: 0.07
Nodes (48): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, time.Time, EscapeHTML(), ScheduledTask, AddTask() (+40 more)

### Community 14 - "acp.go"
Cohesion: 0.22
Nodes (11): checkACPAvailable(), launchACP(), listAvailableACP(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams, ACPMessagePart (+3 more)

### Community 15 - "CallModelWithFallback"
Cohesion: 0.13
Nodes (21): autoClassifyWithModel(), TestParseToolCalls(), ChatRequest, ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModelWithFallback(), findFirstVisionModel() (+13 more)

### Community 16 - "testing.T"
Cohesion: 0.07
Nodes (39): TestBuildThinkingMessage(), TestGetFloatArg(), TestGetInt64Arg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations(), TestToolDescription(), TestTruncOutput() (+31 more)

### Community 17 - "ScorpPath"
Cohesion: 0.12
Nodes (24): init(), setupCLILogging(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), Hostname() (+16 more)

### Community 18 - "RegisterTool"
Cohesion: 0.13
Nodes (14): TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), ResetAutoStats(), init(), init(), init(), init(), TestValidateSubagentTools() (+6 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.19
Nodes (21): AgentMessage, confirmationDisplay(), ClearStopRequest(), ConsumeStopRequest(), sendScorpReply(), confirmKeyboard(), AgentMessage, StorePendingConfirmation() (+13 more)

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
Cohesion: 0.10
Nodes (17): context.Context, AnthropicProvider, CallCommandCode(), CallOpenAIStream(), AzureProvider, ChatMessage, CohereProvider, DashScopeProvider (+9 more)

### Community 24 - "GetAutonomyLevel"
Cohesion: 0.18
Nodes (19): ConfirmationRequired(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels() (+11 more)

### Community 25 - "collector_system.go"
Cohesion: 0.05
Nodes (73): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+65 more)

### Community 26 - "RunPreToolUseHooks"
Cohesion: 0.21
Nodes (20): PreToolUse & PostToolUse Hooks (P3.12), Community Praised Agent Patterns (2026), 2026 AI Coding Agent Competitor Landscape, hookPayload, appendHookContext(), runHookCommand(), RunPostToolUseHooks(), RunPreToolUseHooks() (+12 more)

### Community 27 - "MCPServer"
Cohesion: 0.14
Nodes (14): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, sync.Once, FindMCPTool(), MCPServer, MCPTool, ProbeServer() (+6 more)

### Community 28 - "api_commandcode.go"
Cohesion: 0.21
Nodes (14): buildCommandCodePayload(), CallCommandCodeStream(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), resolveCommandCodeKeyFromDisk(), commandCodeMsg, commandCodeParams (+6 more)

### Community 29 - "CreateCheckpoint"
Cohesion: 0.26
Nodes (19): Checkpoint and Rewind (P1.6), caseCheckpoint(), CheckpointDiffStat(), checkpointRepoRoot(), ckptGit(), ckptRefFor(), CreateCheckpoint(), DeleteCheckpoint() (+11 more)

### Community 31 - "ExecuteSQL"
Cohesion: 0.25
Nodes (6): ExecuteSQL(), loadDBConnections(), dbConnection, TestLiveWebSearch(), ExecuteWebFetch(), ExecuteWebSearch()

### Community 32 - "TestIntegrityStatus"
Cohesion: 0.32
Nodes (15): IsTestRelatedPath(), IsTestRunCommand(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestIsTestRunCommand(), TestTestIntegrityStatus_FailingSuiteDoesNotCount() (+7 more)

### Community 33 - "CheckServerContracts"
Cohesion: 0.20
Nodes (16): MCP Contract Watch (P3.14), Agent Failure Modes and Critiques (2026), caseMCPContract(), CheckServerContracts(), contractFile(), ContractStatusNotice(), ContractWarnings(), serverFingerprint() (+8 more)

### Community 34 - "client.go"
Cohesion: 0.14
Nodes (22): ACPRequest, encoding/json.RawMessage, executeMCPServerTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError() (+14 more)

### Community 35 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 36 - "config/hooks.go"
Cohesion: 0.30
Nodes (13): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+5 more)

### Community 37 - "TruncOutput"
Cohesion: 0.34
Nodes (14): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+6 more)

### Community 38 - "PermissionDecision"
Cohesion: 0.22
Nodes (17): TestAutoAllowlistPrefixAnchoring(), autoAllowlisted(), autoClassify(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand(), PermissionDecision(), ResetAutoAllowlist() (+9 more)

### Community 39 - "runSubagent"
Cohesion: 0.24
Nodes (10): cleanupSubagentSandbox(), createSubagentSandbox(), defaultIsolation(), formatIsolationInfo(), getSubagentIsolation(), isSubagentToolBlocked(), registerSubagentIsolation(), unregisterSubagentIsolation() (+2 more)

### Community 40 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 41 - "HandleConfirmation"
Cohesion: 0.46
Nodes (7): TestStorePendingConfirmationArgs(), clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation

### Community 42 - "SCORP — BRUTAL END-TO-END TEST PLAN"
Cohesion: 0.09
Nodes (20): 0. SCALE & SEVERITY DEFINITIONS, A. GATE-STACK ADVERSARIAL (S) — Layer Penetration Probes, B. LONG-HORIZON & CONTEXT COMPACTION (L), C. INFRASTRUCTURE CHAOS (M/L), D. CONCURRENCY & RACE CONDITIONS (M), E. TELEGRAM UX INTEGRITY (M), F. ADVERSARIAL SECURITY (M), G. EVALUATION & DEPLOYMENT GATE INTEGRITY (S) (+12 more)

### Community 43 - "GetStringArg"
Cohesion: 0.17
Nodes (13): init(), GetStringArg(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteSystemInfo(), ExecuteWriteFile(), formatSendFileSize() (+5 more)

### Community 44 - "api_gemini.go"
Cohesion: 0.23
Nodes (19): callGemini(), callGeminiStream(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), parseDataURL(), resolveGeminiBaseURL() (+11 more)

### Community 45 - "SandboxActive"
Cohesion: 0.25
Nodes (15): caseSandbox(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest(), SandboxStatusNotice(), SandboxVersion() (+7 more)

### Community 46 - "telegram.go"
Cohesion: 0.11
Nodes (27): GetPath(), AnswerCallback(), BackAndRefreshKeyboard(), baseName(), DeleteWebhook(), EditMessage(), EditMessageByID(), MainMenuKeyboard() (+19 more)

### Community 47 - "ResolveAPIKey"
Cohesion: 0.13
Nodes (23): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestGemini_MultimodalParsing(), TestOpenCodeProvider_RegistrationAndKey(), MockCustomProvider, CallModel(), CallModelStream() (+15 more)

### Community 48 - "🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering)"
Cohesion: 0.18
Nodes (10): 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Hasil 15 Skenario Brutal, 🔍 3. Analisis Mendalam Kualitas Output & Arsitektur, 🏆 4. Kesimpulan Akhir Benchmark Brutal, A. Kompilasi Bahasa Go Multi-File (Kasus B2), B. Otomasi Server REST API Latar Belakang & Health Check (Kasus B7), 🏛️ Benchmark Komparasi: 15 Tugas Agentic Berat & Kompleks (Deep Brutal Engineering), C. Fenomena Kebocoran Tag Simulasi pada ZeroClaw (Kasus B1, B3, B5, B8, B10, B12, B15) (+2 more)

### Community 49 - "GetIntArg"
Cohesion: 0.17
Nodes (10): TestGetBoolArg(), TestGetIntArg(), GetBoolArg(), GetIntArg(), getUnclosedTags(), SplitMessage(), ExecuteCompose(), ExecuteGit() (+2 more)

### Community 50 - "eval/core.go"
Cohesion: 0.18
Nodes (17): CheckDenyRules(), loadDenyRules(), ParseDenyRule(), ReloadDenyRules(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped(), TestCheckDenyRulesMatches() (+9 more)

### Community 51 - "eval.go"
Cohesion: 0.29
Nodes (9): Case, caseResult, CoreCases(), humanCount(), liveLabel(), report(), Run(), TestRunnerAggregationAndFilter() (+1 more)

### Community 52 - "runPlanningTurns"
Cohesion: 0.29
Nodes (11): BeginPlanning(), EndPlanning(), AgentMessage, PlanningState(), presentPlanTelegram(), RevisePlan(), RunPlanningLoop(), runPlanningTurns() (+3 more)

### Community 53 - "startCLI"
Cohesion: 0.26
Nodes (18): executeOneShot(), executeTurn(), formatTerminalText(), handleCLISession(), handleCLISOP(), handleMCPCommand(), printBanner(), printCLIHelp() (+10 more)

### Community 54 - "TruncateStr"
Cohesion: 0.17
Nodes (32): TruncateStr(), anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic(), callAnthropicStream() (+24 more)

### Community 55 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "runLiveCase"
Cohesion: 0.22
Nodes (12): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, evalSandboxDir(), deploymentEnv(), LiveCases(), runLiveCase(), tail(), usageDelta() (+4 more)

### Community 57 - "SetTaskPlan"
Cohesion: 0.45
Nodes (11): mkTestPlan(), simulateRestart(), TestClearTaskPlanRemovesPersistedFile(), TestPlanFilePathSanitizesSessionID(), TestPlanPersistenceRoundtripAcrossRestart(), TestStalePlanExpiredOnLoad(), TestUpdateItemStatusPersistsThroughRestart(), uniqueSess() (+3 more)

### Community 58 - "RegisterProvider"
Cohesion: 0.09
Nodes (22): init(), init(), formatMessagesForCLI(), init(), init(), init(), init(), init() (+14 more)

### Community 59 - "bg.go"
Cohesion: 0.38
Nodes (11): bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait(), bgWrite(), closeStdin(), ExecuteBgProcess() (+3 more)

### Community 60 - "prepareNewTurnHistory"
Cohesion: 0.29
Nodes (8): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary()

### Community 61 - "net/http.Request"
Cohesion: 0.33
Nodes (11): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), StartGateway() (+3 more)

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
Cohesion: 0.07
Nodes (62): apiKey, ProjectDir(), envKeyName, CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider (+54 more)

### Community 68 - "IsDangerousCommand"
Cohesion: 0.36
Nodes (6): TestIsDangerousCommand(), devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), IsDangerousCommand(), isHarmlessDevSink()

### Community 69 - "os.File"
Cohesion: 0.19
Nodes (8): acquireSessionLock(), TestAcquireSessionLock(), lockFileExclusive(), unlockFile(), lockFileExclusive(), unlockFile(), os.File, sessionLockFile

### Community 70 - "sync.Mutex"
Cohesion: 0.25
Nodes (8): autoStats, bytes.Buffer, io.WriteCloser, sync.Mutex, CostTracker, getChatLock(), StartTestEndpoint(), BGProcess

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
Cohesion: 0.29
Nodes (9): init(), GetAllTools(), countActiveTools(), countDeferredTools(), ExecuteToolCall(), ExecuteToolList(), ExecuteToolSearch(), callVisionModel() (+1 more)

### Community 76 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 78 - "deploy.sh"
Cohesion: 0.53
Nodes (4): die(), ok(), deploy.sh script, step()

### Community 79 - "HandleTelegramAction"
Cohesion: 0.35
Nodes (12): RequestStop(), FormatUsageStats(), HandleTelegramAction(), BuildSessionMenuKeyboard(), FormatSessionMenuText(), GetActiveSessionID(), HandleSessionCallback(), init() (+4 more)

### Community 80 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 81 - "registry/registry.go"
Cohesion: 0.10
Nodes (28): ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), clearDynamicEnv(), registerTempTool(), TestDynamicToolTTL() (+20 more)

### Community 82 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 83 - "HasGreenTestRun"
Cohesion: 0.31
Nodes (8): Evidence-Based Claim Gate (P4.16), Test-Integrity Gate (P0.3), caseClaimGate(), TestHasGreenTestRun(), TestLooksLikeTestPassClaim(), HasGreenTestRun(), LooksLikeTestPassClaim(), MarkTaskBoundary()

### Community 84 - "handleTelegramMarketplaceInstall"
Cohesion: 0.60
Nodes (5): handleTelegramMarketplaceInstall(), MarketplaceTriOptionKeyboard(), reportInstallOutcome(), BackButtonKeyboard(), SendMessage()

### Community 85 - "main"
Cohesion: 0.21
Nodes (13): RegisterAutonomous(), hasDebugFlag(), isCLIMode(), main(), InitModelUsage(), SOP, Dir(), GetSOP() (+5 more)

### Community 86 - "🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark)"
Cohesion: 0.29
Nodes (6): 📊 1. Ringkasan Hasil Eksekusi, ⏱️ 2. Tabel Waktu Respon Per Tugas (Detik), 🔍 3. Analisis Peningkatan Scorp, 🥊 Benchmark Komparasi: 20 Tugas Ringan (Light Tasks Benchmark), ⚙️ Lingkungan & Konfigurasi Pengujian, Scorp Agent vs PicoClaw vs ZeroClaw (Fresh API Keys & Fast-Path Heuristic)

### Community 87 - "🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut"
Cohesion: 0.29
Nodes (6): 🏛️ 15 Skenario Uji Ekstrem, 📊 1. Ringkasan Hasil Global, ⏱️ 2. Tabel Waktu Eksekusi 15 Skenario Ekstrem (Detik), 🔍 3. Analisis Mendalam Keandalan Arsitektur, 🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut, Scorp Agent vs PicoClaw vs ZeroClaw (Head-to-Head Live VPS Evaluation)

### Community 89 - "registerMCPToolsAsNative"
Cohesion: 0.16
Nodes (14): TruncOutputTool(), buildArgDefsFromInputSchema(), MCPToolsDeferred(), rebuildMCPToolList(), registerMCPToolsAsNative(), StartMCPServers(), StopMCPServers(), TestMCPToolsDeferredEnvParsing() (+6 more)

### Community 90 - "wireCLICallbacks"
Cohesion: 0.29
Nodes (7): formatFileSize(), wireCLICallbacks(), formatFinalResponse(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 99 - "TaskPlan"
Cohesion: 0.36
Nodes (4): PlanItem, TaskPlan, loadPlanFromDisk(), renderPlanItems()

### Community 100 - "Transpiler Phase 3b: Contract Verify & Graceful Degradation"
Cohesion: 0.20
Nodes (14): Multi-Provider LLM Gateway (gateway/gateway.go + models.json), mark3labs/mcp-go SDK, mcp-fetch Reference Port (web scraping/markdown), mcp-filesystem Reference Port (scoped file access), mcp-sqlite Reference Port (local DB inspection), wahyuzero/scorp-mcp-registry (public registry repo), AI Transpiler (Self-Hosted Codegen), Transpiler Phase 3a: Sandbox Build (+6 more)

### Community 101 - "readInteractiveInput"
Cohesion: 0.43
Nodes (6): SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste()

### Community 102 - "MCP Marketplace Build Plan & Execution Roadmap"
Cohesion: 0.21
Nodes (13): MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry, Security Layer 2: AST & Prompt Injection Scan, Security Layer 3: Hermetic CI/CD, ~/.scorp/mcp.json Central Config Entry, mcp_manage Agent Tool (mcp/manage.go), Scorp MCP Marketplace, MCP Marketplace Blueprint (docs/MCP_MARKETPLACE_BLUEPRINT.md) (+5 more)

### Community 103 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 104 - "ExecuteTool"
Cohesion: 0.44
Nodes (7): TestExecuteToolDenyRulesFirst(), registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks(), ExecuteTool()

### Community 105 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 106 - "os/exec.Cmd"
Cohesion: 0.38
Nodes (5): os/exec.Cmd, KillProcessGroup(), SetProcessGroup(), KillProcessGroup(), SetProcessGroup()

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
Cohesion: 0.60
Nodes (5): TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions()

### Community 112 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 113 - "Auto-Mode Classifier (P3.13)"
Cohesion: 0.50
Nodes (4): Auto-Mode Classifier (P3.13), Deny-Rule Engine (P0.2), Durable Memory MEMORY.md (P1.7), Claude Code 2026 Reference Architecture

### Community 114 - "clarify.go"
Cohesion: 0.27
Nodes (9): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID() (+1 more)

### Community 119 - "Release Workflow (tag-triggered cross-compile + GitHub release)"
Cohesion: 0.40
Nodes (6): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection

### Community 120 - "CLI Verify Deadloop (no timeout)"
Cohesion: 0.60
Nodes (6): CLI Verify Deadloop (no timeout), complete_task Explicit Tool Contract, Incident Report: Runaway CLI Verify Session (2026-09-05), Session Lock (cli_lock.go kernel advisory file lock), Heartbeat / Stall Detection, Agent Turn Timeout (context.WithTimeout, SCORP_MAX_TURN_TIMEOUT)

## Knowledge Gaps
- **68 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+63 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 220 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `HomeDir`, `ClearTaskPlan`, `StartDaemon`, `time.Time`, `RunAgentSessionLoop`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`, `CheckServerContracts`, `HandleConfirmation`, `SandboxActive`, `telegram.go`, `runPlanningTurns`, `HandleModelCallback`, `v2_skills.go`, `handleTelegramMarketplaceInstall`, `main`, `registerMCPToolsAsNative`, `TestSteeringQueue`, `clarify.go`?**
  _High betweenness centrality (0.099) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `rag_vector.go`, `ClearTaskPlan`, `getAgentSystemPrompt`, `time.Time`, `CallModelWithFallback`, `RegisterTool`, `compaction_test.go`, `GetAutonomyLevel`, `CreateCheckpoint`, `TestIntegrityStatus`, `PermissionDecision`, `GetStringArg`, `startCLI`, `TruncateStr`, `prepareNewTurnHistory`, `testgate.go`, `IsDangerousCommand`, `IsContinuationDirective`, `HandleTelegramAction`, `registry/registry.go`, `GenerateContextualSessionTitle`, `HasGreenTestRun`, `wireCLICallbacks`, `TaskPlan`, `TestSteeringQueue`, `ExecuteTool`, `ExecuteTermuxAPI`?**
  _High betweenness centrality (0.090) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `chat.go`, `cost_router.go`, `HomeDir`, `rag_vector.go`, `getAgentSystemPrompt`, `time.Time`, `ScorpPath`, `RunAgentSessionLoop`, `collector_system.go`, `client.go`, `skills.go`, `telegram.go`, `HandleModelCallback`, `IsDangerousCommand`, `sync.Mutex`, `v2_skills.go`, `handleTelegramMarketplaceInstall`, `main`, `registerMCPToolsAsNative`, `ExecuteTool`?**
  _High betweenness centrality (0.057) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _68 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.06044303797468355 - nodes in this community are weakly interconnected._
- **Should `chat.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06649616368286446 - nodes in this community are weakly interconnected._