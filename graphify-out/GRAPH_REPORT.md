# Graph Report - scorp  (2026-09-08)

## Corpus Check
- 273 files · ~191,604 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2139 nodes · 5628 edges · 113 communities (99 shown, 6 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 909 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ca038099`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- chat.go
- cost_router.go
- acp.go
- metasearch_engines.go
- RegisterTool
- rag_vector.go
- TaskPlan
- Scorp Agent (Go, ultra-light autonomous agent)
- registerMCPToolsAsNative
- Benchmark
- StartDaemon
- getAgentSystemPrompt
- time.Time
- HandleTelegramAction
- ToolCall
- testing.T
- ScorpPath
- CallOpenCodeStream
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
- eval/core.go
- ExecuteSQL
- TestIntegrityStatus
- CheckServerContracts
- client.go
- skills.go
- LoadMCPConfig
- helpers.go
- ExecuteTool
- runSubagent
- checker.go
- HandleConfirmation
- SCORP — BRUTAL END-TO-END TEST PLAN
- GetStringArg
- api_gemini.go
- ExecuteShell
- prepareNewTurnHistory
- CallModel
- StartTestEndpoint
- init
- CheckDenyRules
- collector_system_native.go
- FormatHourlyReport
- startCLI
- TruncateStr
- patch.go
- runLiveCase
- maybeCompactHistory
- RegisterProvider
- extractTaskMemory
- TruncOutput
- RenameSession
- CompactSessionHistory
- ExecuteTermuxAPI
- HandleUploadInAgentMode
- ExecuteReadURL
- testgate.go
- HandleModelCallback
- IsDangerousCommand
- os.File
- TodoManager
- collector_docker.go
- eval.go
- IsContinuationDirective
- TestGatewayEndpoints
- .listenSSEStream
- install.sh
- deploy.sh
- readInteractiveInput
- EscapeHTML
- collector_coolify.go
- isolation.go
- CredentialVault
- inline.go
- InitDefaultSOPs
- clarify.go
- TestPhase6_AllTools
- main
- wireCLICallbacks
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- serviceBridgeRequests
- TestSteeringQueue
- ScorpDir
- probe_test.go
- ExecuteScript
- registerHookProbeTool
- LoadConfig
- ScreenshotsDir
- providers.go
- metasearch_test.go
- TestGenerateNativeToolsSchema
- TransformToolDefinitions
- ExecuteAutoLogin
- ExecuteAnalyzeImage

## God Nodes (most connected - your core abstractions)
1. `ModelConfig` - 89 edges
2. `HandleTelegramAction()` - 80 edges
3. `ChatMessage` - 75 edges
4. `RunAgentSessionLoop()` - 63 edges
5. `startCLI()` - 51 edges
6. `TruncateStr()` - 47 edges
7. `GetStringArg()` - 47 edges
8. `ScorpPath()` - 42 edges
9. `StartDaemon()` - 42 edges
10. `resumeAgentLoop()` - 40 edges

## Surprising Connections (you probably didn't know these)
- `Auto-Mode Classifier (P3.13)` --implements--> `PermissionDecision()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/auto.go
- `Compaction Preservation (P2.10)` --implements--> `truncateToolResultsInHistory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/compaction.go
- `Durable Memory MEMORY.md (P1.7)` --implements--> `extractTaskMemory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/memory_md.go
- `Plan Mode Workflow (P1.4)` --implements--> `BeginPlanning()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/planmode.go
- `Deny-Rule Engine (P0.2)` --implements--> `CheckDenyRules()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → config/deny.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **MCP Marketplace 5-Layer Security Stack** — docs_build_plan_mcp_marketplace_layer1_source_only_registry, docs_build_plan_mcp_marketplace_layer2_ast_prompt_scan, docs_build_plan_mcp_marketplace_layer3_hermetic_ci, docs_build_plan_mcp_marketplace_sha256_pinning, readme_outbound_secret_redactor [EXTRACTED 1.00]
- **Scorp Resilient State & Memory Workflow (PlanMode + TaskLedger + Checkpoint + MEMORY.md)** — docs_implementation_plan_scorp_plan_mode, docs_implementation_plan_scorp_task_ledger_persistence, docs_implementation_plan_scorp_checkpoint_rewind, docs_implementation_plan_scorp_durable_memory [EXTRACTED 1.00]
- **Scorp Trust & Safety Gate Stack (Sandbox + DenyRules + Hooks + AutoClassifier)** — docs_implementation_plan_scorp_sandbox_bwrap, docs_implementation_plan_scorp_deny_rule_engine, docs_implementation_plan_scorp_hooks_lifecycle, docs_implementation_plan_scorp_auto_mode_classifier [EXTRACTED 1.00]
- **Scorp Verification & Integrity Pipeline (TestGate + ClaimGate + EvalArena)** — docs_implementation_plan_scorp_test_integrity_gate, docs_implementation_plan_scorp_claim_gate, docs_implementation_plan_scorp_eval_arena [EXTRACTED 1.00]
- **AI Transpiler Pipeline (Probe -> Generate -> Build -> Verify)** — docs_build_plan_mcp_marketplace_transpiler, docs_build_plan_mcp_marketplace_transpiler_probe_phase, docs_build_plan_mcp_marketplace_transpiler_generate_phase, docs_build_plan_mcp_marketplace_transpiler_build_phase, docs_build_plan_mcp_marketplace_transpiler_verify_phase [EXTRACTED 1.00]
- **Tri-Option Install Convergence onto ~/.scorp/mcp.json** — docs_build_plan_mcp_marketplace_tri_option_install, docs_build_plan_mcp_marketplace_mcp_json_config, docs_build_plan_mcp_marketplace_mcp_manage, docs_build_plan_mcp_marketplace_watchdog, docs_build_plan_mcp_marketplace_tool_registry [EXTRACTED 1.00]

## Communities (113 total, 6 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.05
Nodes (77): handleMCPCommand(), Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo (+69 more)

### Community 1 - "chat.go"
Cohesion: 0.16
Nodes (28): appendSessionHistory(), EnterAgentMode(), ExitAgentMode(), extractAndSaveMemory(), flushPendingMessages(), GetHistoryTokenEstimate(), getOrCreateSession(), getSession() (+20 more)

### Community 2 - "cost_router.go"
Cohesion: 0.07
Nodes (56): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+48 more)

### Community 3 - "acp.go"
Cohesion: 0.09
Nodes (30): checkACPAvailable(), launchACP(), listAvailableACP(), ACPError, ACPInitializeParams, ACPMessageNewParams, ACPMessagePart, ACPRequest (+22 more)

### Community 4 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (47): net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost(), getClient(), getShortClient() (+39 more)

### Community 5 - "RegisterTool"
Cohesion: 0.07
Nodes (37): init(), init(), init(), init(), ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive() (+29 more)

### Community 6 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 7 - "TaskPlan"
Cohesion: 0.11
Nodes (37): PlanItem, ApprovePlan(), BeginPlanning(), CancelPlan(), EndPlanning(), AgentMessage, PlanningState(), RevisePlan() (+29 more)

### Community 8 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.05
Nodes (58): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection, MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry (+50 more)

### Community 9 - "registerMCPToolsAsNative"
Cohesion: 0.17
Nodes (13): autoStats, sync.Mutex, TruncOutputTool(), buildArgDefsFromInputSchema(), MCPToolsDeferred(), registerMCPToolsAsNative(), StopMCPServers(), TestMCPToolsDeferredEnvParsing() (+5 more)

### Community 10 - "Benchmark"
Cohesion: 0.11
Nodes (33): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+25 more)

### Community 11 - "StartDaemon"
Cohesion: 0.13
Nodes (26): StopMCPServerMode(), Init(), StartServer(), StopServer(), EnsureBuiltinSkills(), runCommandLoop(), StartDaemon(), AnswerCallback() (+18 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.07
Nodes (33): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+25 more)

### Community 13 - "time.Time"
Cohesion: 0.08
Nodes (38): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, os.FileMode, time.Time, ScheduledTask, AddTask() (+30 more)

### Community 14 - "HandleTelegramAction"
Cohesion: 0.15
Nodes (25): RequestStop(), HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath() (+17 more)

### Community 15 - "ToolCall"
Cohesion: 0.23
Nodes (8): TestParseToolCalls(), formatMessagesForCLI(), init(), ClaudeCliProvider, ToolCall, ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls()

### Community 16 - "testing.T"
Cohesion: 0.07
Nodes (41): TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestIsDangerousCommand() (+33 more)

### Community 17 - "ScorpPath"
Cohesion: 0.16
Nodes (20): init(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), HomeDir(), Hostname() (+12 more)

### Community 18 - "CallOpenCodeStream"
Cohesion: 0.24
Nodes (5): CallOpenCodeStream(), init(), resolveOpenCodeBaseURL(), resolveOpenCodeKeyFromDisk(), OpenCodeProvider

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.13
Nodes (27): AgentMessage, confirmationDisplay(), ClearStopRequest(), ConsumeStopRequest(), cleanToolCallTags(), getSessionSearchContext(), maxIterations(), maxTurnTimeout() (+19 more)

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
Nodes (17): context.Context, AnthropicProvider, CallCommandCode(), CallCommandCodeStream(), AzureProvider, ChatMessage, ChatRequest, CohereProvider (+9 more)

### Community 24 - "GetAutonomyLevel"
Cohesion: 0.20
Nodes (15): ConfirmationRequired(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels() (+7 more)

### Community 25 - "collector_system.go"
Cohesion: 0.23
Nodes (18): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+10 more)

### Community 26 - "RunPreToolUseHooks"
Cohesion: 0.11
Nodes (35): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+27 more)

### Community 27 - "MCPServer"
Cohesion: 0.15
Nodes (13): bufio.Scanner, encoding/json.Encoder, sync.Once, FindMCPTool(), MCPServer, MCPTool, ProbeServer(), startMCPServer() (+5 more)

### Community 28 - "api_commandcode.go"
Cohesion: 0.17
Nodes (13): buildCommandCodePayload(), createCommandCodeRequest(), extractFallbackToolCalls(), init(), resolveCommandCodeKeyFromDisk(), commandCodeMsg, commandCodeParams, commandCodePayload (+5 more)

### Community 29 - "CreateCheckpoint"
Cohesion: 0.26
Nodes (19): Checkpoint and Rewind (P1.6), caseCheckpoint(), CheckpointDiffStat(), checkpointRepoRoot(), ckptGit(), ckptRefFor(), CreateCheckpoint(), DeleteCheckpoint() (+11 more)

### Community 30 - "eval/core.go"
Cohesion: 0.24
Nodes (9): TestExecuteToolDenyRulesFirst(), ReloadDenyRules(), caseDangerGate(), caseDenyInvalidSkipped(), caseDenyShellYOLO(), caseLedgerClear(), caseLedgerPersisted(), caseSensitivePath() (+1 more)

### Community 31 - "ExecuteSQL"
Cohesion: 0.83
Nodes (3): ExecuteSQL(), loadDBConnections(), dbConnection

### Community 32 - "TestIntegrityStatus"
Cohesion: 0.29
Nodes (16): Test-Integrity Gate (P0.3), IsTestRelatedPath(), IsTestRunCommand(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestIsTestRunCommand() (+8 more)

### Community 33 - "CheckServerContracts"
Cohesion: 0.20
Nodes (16): MCP Contract Watch (P3.14), Agent Failure Modes and Critiques (2026), caseMCPContract(), CheckServerContracts(), contractFile(), ContractStatusNotice(), ContractWarnings(), serverFingerprint() (+8 more)

### Community 34 - "client.go"
Cohesion: 0.19
Nodes (18): encoding/json.RawMessage, executeMCPServerTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError(), sendMCPResult() (+10 more)

### Community 35 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 36 - "LoadMCPConfig"
Cohesion: 0.31
Nodes (13): MCPConfigFilePath(), LoadMCPConfig(), rebuildMCPToolList(), ReloadMCPServers(), sanitizeMCPName(), StartMCPServers(), AddServerEntry(), ExecuteMCPManage() (+5 more)

### Community 37 - "helpers.go"
Cohesion: 0.67
Nodes (3): GetStringSliceArg(), getUnclosedTags(), SplitMessage()

### Community 38 - "ExecuteTool"
Cohesion: 0.16
Nodes (25): TestAutoAllowlistPrefixAnchoring(), TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), autoAllowlisted(), autoClassify(), autoClassifyWithModel(), AutoStatsSnapshot(), bumpAutoStat() (+17 more)

### Community 39 - "runSubagent"
Cohesion: 0.18
Nodes (21): runOpenCodeCLI(), runSubagentACP(), AgentMessage, TestParseDelegateParams(), TestParseDelegateParamsCapsAndDefaults(), TestValidateSubagentTools(), delegateBatchParams, delegateResult (+13 more)

### Community 40 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 41 - "HandleConfirmation"
Cohesion: 0.32
Nodes (11): TestStorePendingConfirmationArgs(), clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation() (+3 more)

### Community 42 - "SCORP — BRUTAL END-TO-END TEST PLAN"
Cohesion: 0.09
Nodes (20): 0. SCALE & SEVERITY DEFINITIONS, A. GATE-STACK ADVERSARIAL (S) — Layer Penetration Probes, B. LONG-HORIZON & CONTEXT COMPACTION (L), C. INFRASTRUCTURE CHAOS (M/L), D. CONCURRENCY & RACE CONDITIONS (M), E. TELEGRAM UX INTEGRITY (M), F. ADVERSARIAL SECURITY (M), G. EVALUATION & DEPLOYMENT GATE INTEGRITY (S) (+12 more)

### Community 43 - "GetStringArg"
Cohesion: 0.19
Nodes (12): init(), GetStringArg(), SendDocumentBytes(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteSystemInfo(), ExecuteWriteFile() (+4 more)

### Community 44 - "api_gemini.go"
Cohesion: 0.18
Nodes (18): callGemini(), callGeminiStream(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), resolveGeminiBaseURL(), geminiContent (+10 more)

### Community 45 - "ExecuteShell"
Cohesion: 0.23
Nodes (17): Sandbox Shell Bubblewrap (P0.1), caseSandbox(), ExecuteShell(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest() (+9 more)

### Community 46 - "prepareNewTurnHistory"
Cohesion: 0.29
Nodes (8): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary()

### Community 47 - "CallModel"
Cohesion: 0.14
Nodes (24): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestOpenCodeProvider_RegistrationAndKey(), ChatResponse, getCheapestModel(), RouteModelCostAware(), CostTracker (+16 more)

### Community 49 - "init"
Cohesion: 0.14
Nodes (15): init(), GetBoolArg(), GetIntArg(), unregisterMCPNativeTools(), GetAllTools(), ExecuteCompose(), countActiveTools(), countDeferredTools() (+7 more)

### Community 50 - "CheckDenyRules"
Cohesion: 0.24
Nodes (12): CheckDenyRules(), loadDenyRules(), ParseDenyRule(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped(), TestCheckDenyRulesMatches(), TestParseDenyRule() (+4 more)

### Community 51 - "collector_system_native.go"
Cohesion: 0.21
Nodes (18): FormatDuration(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+10 more)

### Community 52 - "FormatHourlyReport"
Cohesion: 0.18
Nodes (16): NetworkData, SystemData, CollectSystem(), GetTopProcesses(), TopProcess, PortInfo, Bar(), bar() (+8 more)

### Community 53 - "startCLI"
Cohesion: 0.27
Nodes (18): executeOneShot(), executeTurn(), formatTerminalText(), handleCLISOP(), hasDebugFlag(), printBanner(), printCLIHelp(), printCostUsage() (+10 more)

### Community 54 - "TruncateStr"
Cohesion: 0.18
Nodes (35): net/http.Request, TruncateStr(), anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic() (+27 more)

### Community 55 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "runLiveCase"
Cohesion: 0.22
Nodes (12): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, evalSandboxDir(), deploymentEnv(), LiveCases(), runLiveCase(), tail(), usageDelta() (+4 more)

### Community 57 - "maybeCompactHistory"
Cohesion: 0.27
Nodes (16): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, injectPreservationNote(), latestUserGoal(), maybeCompactHistory(), preservationNote() (+8 more)

### Community 58 - "RegisterProvider"
Cohesion: 0.12
Nodes (17): init(), init(), init(), init(), init(), init(), init(), init() (+9 more)

### Community 59 - "extractTaskMemory"
Cohesion: 0.21
Nodes (16): AppendMemoryMD(), extractTaskMemory(), AgentMessage, parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile(), memoryMDTestFile() (+8 more)

### Community 60 - "TruncOutput"
Cohesion: 0.34
Nodes (14): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+6 more)

### Community 61 - "RenameSession"
Cohesion: 0.30
Nodes (13): ClearChatSession(), FlushSessionHistory(), historyFilePath(), saveHistoryToDisk(), DeleteSession(), ListSessions(), RenameSession(), sanitizeSessionID() (+5 more)

### Community 62 - "CompactSessionHistory"
Cohesion: 0.30
Nodes (12): CompactSessionHistory(), FormatCompactStats(), CompactStats, BuildSessionMenuKeyboard(), FormatSessionMenuText(), GetActiveSessionID(), HandleSessionCallback(), init() (+4 more)

### Community 63 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 64 - "HandleUploadInAgentMode"
Cohesion: 0.18
Nodes (11): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), sendScorpReply(), contentPart, imageURL, TestBase64Encode(), cleanupAgentSessions() (+3 more)

### Community 65 - "ExecuteReadURL"
Cohesion: 0.36
Nodes (7): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), TestReadURL_LocalMock(), truncateURLOutput(), tryRemoteScrape()

### Community 66 - "testgate.go"
Cohesion: 0.14
Nodes (22): OperationalClaim, GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), dedupeOpObjects(), extractOpObjects() (+14 more)

### Community 67 - "HandleModelCallback"
Cohesion: 0.08
Nodes (57): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+49 more)

### Community 68 - "IsDangerousCommand"
Cohesion: 0.43
Nodes (5): devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), IsDangerousCommand(), isHarmlessDevSink()

### Community 69 - "os.File"
Cohesion: 0.19
Nodes (8): acquireSessionLock(), TestAcquireSessionLock(), lockFileExclusive(), unlockFile(), lockFileExclusive(), unlockFile(), os.File, sessionLockFile

### Community 70 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 71 - "collector_docker.go"
Cohesion: 0.27
Nodes (11): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+3 more)

### Community 72 - "eval.go"
Cohesion: 0.29
Nodes (9): Case, caseResult, CoreCases(), humanCount(), liveLabel(), report(), Run(), TestRunnerAggregationAndFilter() (+1 more)

### Community 73 - "IsContinuationDirective"
Cohesion: 0.40
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 75 - "TestGatewayEndpoints"
Cohesion: 0.27
Nodes (11): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), StartGateway() (+3 more)

### Community 76 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 78 - "deploy.sh"
Cohesion: 0.53
Nodes (4): die(), ok(), deploy.sh script, step()

### Community 79 - "readInteractiveInput"
Cohesion: 0.23
Nodes (11): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste(), getContextPill() (+3 more)

### Community 80 - "EscapeHTML"
Cohesion: 0.21
Nodes (11): collectTableLines(), convertInlineMarkdown(), convertTableToList(), TestMarkdownToTelegramHTMLEscaping(), isSeparatorRow(), markdownToTelegramHTML(), parseTableRow(), safeIndex() (+3 more)

### Community 81 - "collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 82 - "isolation.go"
Cohesion: 0.23
Nodes (9): cleanupSubagentSandbox(), createSubagentSandbox(), defaultIsolation(), formatIsolationInfo(), getSubagentIsolation(), isSubagentToolBlocked(), registerSubagentIsolation(), unregisterSubagentIsolation() (+1 more)

### Community 83 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 84 - "inline.go"
Cohesion: 0.33
Nodes (10): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+2 more)

### Community 85 - "InitDefaultSOPs"
Cohesion: 0.42
Nodes (8): SOP, Dir(), GetSOP(), InitDefaultSOPs(), ListSOPs(), SaveSOP(), TestSOPLifecycle(), ExecuteSOP()

### Community 86 - "clarify.go"
Cohesion: 0.29
Nodes (8): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), SetClarifyChatID(), PendingClarify

### Community 87 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 89 - "main"
Cohesion: 0.25
Nodes (6): RegisterAutonomous(), isCLIMode(), main(), FormatModelListWithHealth(), FormatUsageStats(), InitModelUsage()

### Community 90 - "wireCLICallbacks"
Cohesion: 0.31
Nodes (6): wireCLICallbacks(), formatFinalResponse(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 99 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 100 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 101 - "ScorpDir"
Cohesion: 0.48
Nodes (6): ScorpDir(), envKeyName, installSystemdService(), maskString(), RunQuickstart(), saveConfig()

### Community 102 - "probe_test.go"
Cohesion: 0.38
Nodes (6): contains(), goldenFilesystemBinary(), TestDetectRuntimeMessages(), TestProbeAndVerifyGoldenRoundTrip(), TestProbeUpstreamNpx(), TestStripCodeFences()

### Community 103 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 104 - "registerHookProbeTool"
Cohesion: 0.73
Nodes (5): registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks()

### Community 105 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 106 - "ScreenshotsDir"
Cohesion: 0.33
Nodes (5): RagVectorDBPath(), ScreenshotsDir(), TestConfigPaths_RagVectorDBPath(), TestConfigPaths_ScorpDir(), TestConfigPaths_ScreenshotsDir()

### Community 107 - "providers.go"
Cohesion: 0.40
Nodes (5): applyProviderDefaults(), ProviderPreset, hasAPIKey(), RegisterProviderPreset(), FormatModelList()

### Community 108 - "metasearch_test.go"
Cohesion: 0.40
Nodes (4): TestDeduplicateAndRank(), TestLiveWebSearch(), TestMetaSearchAggregator_ConcurrentSuccess(), TestMetaSearchAggregator_FaultTolerance()

### Community 109 - "TestGenerateNativeToolsSchema"
Cohesion: 0.40
Nodes (4): IsRateLimitError(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema(), TestIsRateLimitError()

### Community 110 - "TransformToolDefinitions"
Cohesion: 0.60
Nodes (5): TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions()

## Knowledge Gaps
- **39 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+34 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 190 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `TaskPlan`, `registerMCPToolsAsNative`, `StartDaemon`, `getAgentSystemPrompt`, `time.Time`, `RunAgentSessionLoop`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`, `CheckServerContracts`, `HandleConfirmation`, `ExecuteShell`, `FormatHourlyReport`, `RenameSession`, `CompactSessionHistory`, `HandleModelCallback`, `collector_docker.go`, `collector_coolify.go`, `InitDefaultSOPs`, `clarify.go`, `main`, `TestSteeringQueue`?**
  _High betweenness centrality (0.108) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `RegisterTool`, `rag_vector.go`, `TaskPlan`, `getAgentSystemPrompt`, `HandleTelegramAction`, `GetAutonomyLevel`, `CreateCheckpoint`, `TestIntegrityStatus`, `ExecuteTool`, `HandleConfirmation`, `GetStringArg`, `prepareNewTurnHistory`, `CallModel`, `startCLI`, `TruncateStr`, `maybeCompactHistory`, `extractTaskMemory`, `RenameSession`, `ExecuteTermuxAPI`, `HandleUploadInAgentMode`, `testgate.go`, `IsDangerousCommand`, `IsContinuationDirective`, `EscapeHTML`, `wireCLICallbacks`, `TestSteeringQueue`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Why does `TruncateStr()` connect `TruncateStr` to `chat.go`, `cost_router.go`, `TaskPlan`, `registerMCPToolsAsNative`, `getAgentSystemPrompt`, `time.Time`, `CallOpenCodeStream`, `RunAgentSessionLoop`, `context.Context`, `MCPServer`, `helpers.go`, `ExecuteTool`, `runSubagent`, `api_gemini.go`, `init`, `extractTaskMemory`, `TruncOutput`, `RenameSession`, `HandleModelCallback`, `EscapeHTML`?**
  _High betweenness centrality (0.066) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _39 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.05362614913176711 - nodes in this community are weakly interconnected._
- **Should `cost_router.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06663141195134849 - nodes in this community are weakly interconnected._