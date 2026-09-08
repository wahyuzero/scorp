# Graph Report - scorp  (2026-09-08)

## Corpus Check
- 274 files · ~192,942 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2146 nodes · 5645 edges · 121 communities (104 shown, 9 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 911 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2b90ab0f`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- chat.go
- cost_router.go
- runSubagent
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
- HomeDir
- ToolCall
- testing.T
- ScorpPath
- CallOpenCodeWithTools
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
- SaveModelConfig
- TruncOutput
- TestIntegrityStatus
- CheckServerContracts
- client.go
- skills.go
- LoadMCPConfig
- wizard.go
- PermissionDecision
- ResetNativeToolCache
- checker.go
- HandleConfirmation
- SCORP — BRUTAL END-TO-END TEST PLAN
- exec.go
- api_gemini.go
- ExecuteShell
- prepareNewTurnHistory
- ResolveAPIKey
- sync.Mutex
- GetStringArg
- eval/core.go
- collector_system_native.go
- FormatHourlyReport
- startCLI
- TruncateStr
- patch.go
- runLiveCase
- maybeCompactHistory
- RegisterProvider
- extractTaskMemory
- StreamChunk
- RenameSession
- session_ui.go
- ExecuteTermuxAPI
- HandleUploadInAgentMode
- ExecuteReadURL
- testgate.go
- HandleModelCallback
- IsDangerousCommand
- os.File
- TodoManager
- collector_docker.go
- time.Duration
- IsContinuationDirective
- net/http.Request
- .listenSSEStream
- install.sh
- deploy.sh
- readInteractiveInput
- markdownToTelegramHTML
- collector_coolify.go
- RecordToolReceipt
- HasGreenTestRun
- inline.go
- InitDefaultSOPs
- clarify.go
- model_catalog.go
- HandleTelegramAction
- wireCLICallbacks
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- serviceBridgeRequests
- net/http.Client
- main
- startMCPServer
- SectionSecurity
- ExecuteTool
- LoadConfig
- GetAllTools
- tools/callbacks.go
- SearchResult
- TestGenerateNativeToolsSchema
- TransformToolDefinitions
- NewDefaultMetaSearchAggregator
- ExecuteAnalyzeImage
- GenerateContextualSessionTitle
- TestToolSearchActivatesDeferredToolInStaticMode
- Auto-Mode Classifier (P3.13)
- BraveSearchEngine
- DuckDuckGoHTMLEngine
- GitHubEngine
- WikipediaEngine
- FormatCompactStats

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
- `MCP Contract Watch (P3.14)` --implements--> `CheckServerContracts()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → mcp/contract.go
- `Test-Integrity Gate (P0.3)` --implements--> `TestIntegrityStatus()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → tools/testgate.go
- `Auto-Mode Classifier (P3.13)` --implements--> `PermissionDecision()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/auto.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **MCP Marketplace 5-Layer Security Stack** — docs_build_plan_mcp_marketplace_layer1_source_only_registry, docs_build_plan_mcp_marketplace_layer2_ast_prompt_scan, docs_build_plan_mcp_marketplace_layer3_hermetic_ci, docs_build_plan_mcp_marketplace_sha256_pinning, readme_outbound_secret_redactor [EXTRACTED 1.00]
- **Scorp Resilient State & Memory Workflow (PlanMode + TaskLedger + Checkpoint + MEMORY.md)** — docs_implementation_plan_scorp_plan_mode, docs_implementation_plan_scorp_task_ledger_persistence, docs_implementation_plan_scorp_checkpoint_rewind, docs_implementation_plan_scorp_durable_memory [EXTRACTED 1.00]
- **Scorp Trust & Safety Gate Stack (Sandbox + DenyRules + Hooks + AutoClassifier)** — docs_implementation_plan_scorp_sandbox_bwrap, docs_implementation_plan_scorp_deny_rule_engine, docs_implementation_plan_scorp_hooks_lifecycle, docs_implementation_plan_scorp_auto_mode_classifier [EXTRACTED 1.00]
- **Scorp Verification & Integrity Pipeline (TestGate + ClaimGate + EvalArena)** — docs_implementation_plan_scorp_test_integrity_gate, docs_implementation_plan_scorp_claim_gate, docs_implementation_plan_scorp_eval_arena [EXTRACTED 1.00]
- **AI Transpiler Pipeline (Probe -> Generate -> Build -> Verify)** — docs_build_plan_mcp_marketplace_transpiler, docs_build_plan_mcp_marketplace_transpiler_probe_phase, docs_build_plan_mcp_marketplace_transpiler_generate_phase, docs_build_plan_mcp_marketplace_transpiler_build_phase, docs_build_plan_mcp_marketplace_transpiler_verify_phase [EXTRACTED 1.00]
- **Tri-Option Install Convergence onto ~/.scorp/mcp.json** — docs_build_plan_mcp_marketplace_tri_option_install, docs_build_plan_mcp_marketplace_mcp_json_config, docs_build_plan_mcp_marketplace_mcp_manage, docs_build_plan_mcp_marketplace_watchdog, docs_build_plan_mcp_marketplace_tool_registry [EXTRACTED 1.00]

## Communities (121 total, 9 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (74): Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo, RegistryIndex (+66 more)

### Community 1 - "chat.go"
Cohesion: 0.15
Nodes (30): appendSessionHistory(), EnterAgentMode(), ExitAgentMode(), extractAndSaveMemory(), flushPendingMessages(), FlushSessionHistory(), GetHistoryTokenEstimate(), getOrCreateSession() (+22 more)

### Community 2 - "cost_router.go"
Cohesion: 0.07
Nodes (56): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+48 more)

### Community 3 - "runSubagent"
Cohesion: 0.05
Nodes (60): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+52 more)

### Community 4 - "metasearch_engines.go"
Cohesion: 0.13
Nodes (11): BingEngine, DuckDuckGoLiteEngine, ghSearchResp, GoogleCSEEngine, NewBingEngine(), NewDuckDuckGoLiteEngine(), NewGoogleCSEEngine(), NewSearXNGEngine() (+3 more)

### Community 5 - "RegisterTool"
Cohesion: 0.08
Nodes (19): init(), init(), init(), init(), echoPlugin, RegisterPlugin(), TestRegisterPlugin(), ExecuteToolByName() (+11 more)

### Community 6 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 7 - "TaskPlan"
Cohesion: 0.09
Nodes (42): PlanItem, ApprovePlan(), BeginPlanning(), CancelPlan(), EndPlanning(), AgentMessage, PlanningState(), RevisePlan() (+34 more)

### Community 8 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.05
Nodes (58): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection, MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry (+50 more)

### Community 9 - "registerMCPToolsAsNative"
Cohesion: 0.18
Nodes (12): TruncOutputTool(), buildArgDefsFromInputSchema(), MCPToolsDeferred(), registerMCPToolsAsNative(), StopMCPServers(), TestMCPToolsDeferredEnvParsing(), ServerWatchdog, GetServerHealthStatus() (+4 more)

### Community 10 - "Benchmark"
Cohesion: 0.09
Nodes (39): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+31 more)

### Community 11 - "StartDaemon"
Cohesion: 0.15
Nodes (24): runCommandLoop(), StartDaemon(), BackAndRefreshKeyboard(), baseName(), DeleteWebhook(), EditMessage(), EditMessageByID(), InitTelegram() (+16 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.07
Nodes (34): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+26 more)

### Community 13 - "time.Time"
Cohesion: 0.09
Nodes (34): presentPlanTelegram(), CM(), InitConfigManager(), NewConfigManager(), ConfigManager, time.Time, EscapeHTML(), ScheduledTask (+26 more)

### Community 14 - "HomeDir"
Cohesion: 0.21
Nodes (19): HomeDir(), ProjectDir(), PythonSitePackages(), UploadsDir(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard() (+11 more)

### Community 15 - "ToolCall"
Cohesion: 0.26
Nodes (7): TestParseToolCalls(), formatMessagesForCLI(), ClaudeCliProvider, ToolCall, ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls()

### Community 16 - "testing.T"
Cohesion: 0.06
Nodes (45): TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestIsDangerousCommand() (+37 more)

### Community 17 - "ScorpPath"
Cohesion: 0.06
Nodes (51): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+43 more)

### Community 18 - "CallOpenCodeWithTools"
Cohesion: 0.25
Nodes (9): CallOpenCode(), CallOpenCodeStream(), CallOpenCodeWithTools(), resolveOpenCodeBaseURL(), resolveOpenCodeKeyFromDisk(), resolveOpenCodeSessionID(), RecordCostWithCache(), OpenCodeProvider (+1 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.22
Nodes (17): AgentMessage, confirmationDisplay(), ClearStopRequest(), ConsumeStopRequest(), confirmKeyboard(), cleanToolCallTags(), getSessionSearchContext(), maxIterations() (+9 more)

### Community 20 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 21 - "compaction_test.go"
Cohesion: 0.19
Nodes (20): AgentMessage, makeHistory(), makeToolResult(), mkCompactionHistory(), TestEstimateHistoryTokens(), TestPrune_BoundaryAges(), TestPrune_DockerScenario_41Messages(), TestPrune_EmptyHistory() (+12 more)

### Community 22 - "collector_security.go"
Cohesion: 0.24
Nodes (20): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+12 more)

### Community 23 - "context.Context"
Cohesion: 0.10
Nodes (18): context.Context, AnthropicProvider, applyOpenAIHeaders(), buildOpenAIRequestBody(), CallOpenAI(), CallOpenAIWithTools(), formatOpenAIMessages(), AzureProvider (+10 more)

### Community 24 - "GetAutonomyLevel"
Cohesion: 0.17
Nodes (17): GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels(), TestConfirmationRequired() (+9 more)

### Community 25 - "collector_system.go"
Cohesion: 0.23
Nodes (18): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+10 more)

### Community 26 - "RunPreToolUseHooks"
Cohesion: 0.13
Nodes (33): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+25 more)

### Community 27 - "MCPServer"
Cohesion: 0.18
Nodes (9): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, sync.Once, FindMCPTool(), MCPServer, MCPTool, jsonRPCError (+1 more)

### Community 28 - "api_commandcode.go"
Cohesion: 0.17
Nodes (15): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeStream(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), resolveCommandCodeKeyFromDisk(), commandCodeMsg (+7 more)

### Community 29 - "CreateCheckpoint"
Cohesion: 0.26
Nodes (19): Checkpoint and Rewind (P1.6), caseCheckpoint(), CheckpointDiffStat(), checkpointRepoRoot(), ckptGit(), ckptRefFor(), CreateCheckpoint(), DeleteCheckpoint() (+11 more)

### Community 30 - "SaveModelConfig"
Cohesion: 0.20
Nodes (15): defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, ModelRouterConfig, ModelUsage, SwitchModel(), getProviderInfo() (+7 more)

### Community 31 - "TruncOutput"
Cohesion: 0.26
Nodes (9): StorePendingConfirmation(), ConfirmationRequired(), TruncOutput(), ExecuteSQL(), loadDBConnections(), dbConnection, ExecuteSystemInfo(), ExecuteGit() (+1 more)

### Community 32 - "TestIntegrityStatus"
Cohesion: 0.32
Nodes (15): IsTestRelatedPath(), IsTestRunCommand(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestIsTestRunCommand(), TestTestIntegrityStatus_FailingSuiteDoesNotCount() (+7 more)

### Community 33 - "CheckServerContracts"
Cohesion: 0.24
Nodes (14): caseMCPContract(), CheckServerContracts(), contractFile(), ContractStatusNotice(), ContractWarnings(), serverFingerprint(), SetContractPath(), contractTestServer() (+6 more)

### Community 34 - "client.go"
Cohesion: 0.19
Nodes (17): encoding/json.RawMessage, executeMCPServerTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError(), sendMCPResult() (+9 more)

### Community 35 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 36 - "LoadMCPConfig"
Cohesion: 0.31
Nodes (13): MCPConfigFilePath(), LoadMCPConfig(), rebuildMCPToolList(), ReloadMCPServers(), sanitizeMCPName(), StartMCPServers(), AddServerEntry(), ExecuteMCPManage() (+5 more)

### Community 37 - "wizard.go"
Cohesion: 0.31
Nodes (17): AutoPopulateFromCatalog(), ProviderKeyEnv(), modelWizard, askAPIKey(), ClearModelWizard(), EnvFilePath(), finalizeModelKeySave(), finalizeProviderKeySave() (+9 more)

### Community 38 - "PermissionDecision"
Cohesion: 0.22
Nodes (16): autoAllowlisted(), autoClassify(), autoClassifyWithModel(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand(), PermissionDecision(), setAutoMode() (+8 more)

### Community 39 - "ResetNativeToolCache"
Cohesion: 0.29
Nodes (13): ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), clearDynamicEnv(), registerTempTool(), TestDynamicToolTTL() (+5 more)

### Community 40 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 41 - "HandleConfirmation"
Cohesion: 0.35
Nodes (10): TestStorePendingConfirmationArgs(), clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), StorePendingConfirmationArgs() (+2 more)

### Community 42 - "SCORP — BRUTAL END-TO-END TEST PLAN"
Cohesion: 0.09
Nodes (20): 0. SCALE & SEVERITY DEFINITIONS, A. GATE-STACK ADVERSARIAL (S) — Layer Penetration Probes, B. LONG-HORIZON & CONTEXT COMPACTION (L), C. INFRASTRUCTURE CHAOS (M/L), D. CONCURRENCY & RACE CONDITIONS (M), E. TELEGRAM UX INTEGRITY (M), F. ADVERSARIAL SECURITY (M), G. EVALUATION & DEPLOYMENT GATE INTEGRITY (S) (+12 more)

### Community 43 - "exec.go"
Cohesion: 0.22
Nodes (8): init(), ExecuteListDir(), ExecuteSendFile(), ExecuteWriteFile(), formatSendFileSize(), isGuestLinuxRootfs(), isPathAllowed(), TestExecuteSendFile_SafetyAndDispatch()

### Community 44 - "api_gemini.go"
Cohesion: 0.18
Nodes (18): callGemini(), callGeminiStream(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), resolveGeminiBaseURL(), geminiContent (+10 more)

### Community 45 - "ExecuteShell"
Cohesion: 0.23
Nodes (17): Sandbox Shell Bubblewrap (P0.1), caseSandbox(), ExecuteShell(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest() (+9 more)

### Community 46 - "prepareNewTurnHistory"
Cohesion: 0.29
Nodes (8): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary()

### Community 47 - "ResolveAPIKey"
Cohesion: 0.15
Nodes (21): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestOpenCodeProvider_RegistrationAndKey(), CallModel(), CallModelStream(), GetProvider(), TestCoreAndExtendedProvidersRegistered() (+13 more)

### Community 48 - "sync.Mutex"
Cohesion: 0.50
Nodes (4): autoStats, sync.Mutex, getChatLock(), StartTestEndpoint()

### Community 49 - "GetStringArg"
Cohesion: 0.15
Nodes (15): init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage(), ExecuteCompose(), ExecuteToolCall() (+7 more)

### Community 50 - "eval/core.go"
Cohesion: 0.16
Nodes (18): TestExecuteToolDenyRulesFirst(), CheckDenyRules(), loadDenyRules(), ParseDenyRule(), ReloadDenyRules(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped() (+10 more)

### Community 51 - "collector_system_native.go"
Cohesion: 0.23
Nodes (17): CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses(), GetTopProcesses() (+9 more)

### Community 52 - "FormatHourlyReport"
Cohesion: 0.21
Nodes (13): NetworkData, SystemData, CollectSystem(), GetTopProcesses(), TopProcess, PortInfo, Bar(), bar() (+5 more)

### Community 53 - "startCLI"
Cohesion: 0.27
Nodes (17): executeOneShot(), executeTurn(), formatTerminalText(), handleCLISession(), handleCLISOP(), handleMCPCommand(), printBanner(), printCLIHelp() (+9 more)

### Community 54 - "TruncateStr"
Cohesion: 0.27
Nodes (20): TruncateStr(), anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic(), callAnthropicStream() (+12 more)

### Community 55 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "runLiveCase"
Cohesion: 0.14
Nodes (21): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, Case, caseResult, CoreCases(), evalSandboxDir(), humanCount(), liveLabel() (+13 more)

### Community 57 - "maybeCompactHistory"
Cohesion: 0.27
Nodes (16): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, injectPreservationNote(), latestUserGoal(), maybeCompactHistory(), preservationNote() (+8 more)

### Community 58 - "RegisterProvider"
Cohesion: 0.11
Nodes (20): init(), init(), init(), init(), init(), init(), init(), init() (+12 more)

### Community 59 - "extractTaskMemory"
Cohesion: 0.25
Nodes (14): AppendMemoryMD(), extractTaskMemory(), AgentMessage, parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile(), memoryMDTestFile() (+6 more)

### Community 60 - "StreamChunk"
Cohesion: 0.26
Nodes (13): ChatResponse, getCheapestModel(), RouteModelCostAware(), CostTracker, CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), isVisionModelName() (+5 more)

### Community 61 - "RenameSession"
Cohesion: 0.33
Nodes (11): ClearChatSession(), historyFilePath(), saveHistoryToDisk(), DeleteSession(), ListSessions(), RenameSession(), sanitizeSessionID(), SessionExists() (+3 more)

### Community 62 - "session_ui.go"
Cohesion: 0.44
Nodes (9): BuildSessionMenuKeyboard(), FormatSessionMenuText(), GetActiveSessionID(), HandleSessionCallback(), init(), loadTgSessionMapping(), saveTgSessionMapping(), SetActiveSessionID() (+1 more)

### Community 63 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 64 - "HandleUploadInAgentMode"
Cohesion: 0.20
Nodes (10): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), contentPart, imageURL, TestBase64Encode(), cleanupAgentSessions(), TGDocument (+2 more)

### Community 65 - "ExecuteReadURL"
Cohesion: 0.36
Nodes (7): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), TestReadURL_LocalMock(), truncateURLOutput(), tryRemoteScrape()

### Community 66 - "testgate.go"
Cohesion: 0.20
Nodes (16): OperationalClaim, dedupeOpObjects(), extractOpObjects(), isWriteToolName(), LooksLikeOperationalClaims(), opReceipt(), seedOpReceipts(), TestLooksLikeOperationalClaims() (+8 more)

### Community 67 - "HandleModelCallback"
Cohesion: 0.20
Nodes (16): ApiFormatPickerKeyboard(), FallbackEditorKeyboard(), ModelInfoKeyboard(), ModelMenuKeyboard(), ModelPickerKeyboard(), addToFallback(), FallbackText(), FormatAPIKeysList() (+8 more)

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
Cohesion: 0.24
Nodes (12): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+4 more)

### Community 72 - "time.Duration"
Cohesion: 0.27
Nodes (13): FormatDuration(), time.Duration, AutonomousConfig, AddUptimeTarget(), checkTarget(), ExecuteUptime(), ListUptimeTargets(), RemoveUptimeTarget() (+5 more)

### Community 73 - "IsContinuationDirective"
Cohesion: 0.40
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 75 - "net/http.Request"
Cohesion: 0.31
Nodes (12): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), StartGateway() (+4 more)

### Community 76 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 78 - "deploy.sh"
Cohesion: 0.53
Nodes (4): die(), ok(), deploy.sh script, step()

### Community 79 - "readInteractiveInput"
Cohesion: 0.23
Nodes (11): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste(), getContextPill() (+3 more)

### Community 80 - "markdownToTelegramHTML"
Cohesion: 0.24
Nodes (9): collectTableLines(), convertInlineMarkdown(), convertTableToList(), TestMarkdownToTelegramHTMLEscaping(), isSeparatorRow(), markdownToTelegramHTML(), parseTableRow(), safeIndex() (+1 more)

### Community 81 - "collector_coolify.go"
Cohesion: 0.28
Nodes (12): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+4 more)

### Community 82 - "RecordToolReceipt"
Cohesion: 0.26
Nodes (8): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), RedactSecrets(), TestRedactSecrets(), ToolReceipt

### Community 83 - "HasGreenTestRun"
Cohesion: 0.24
Nodes (10): Evidence-Based Claim Gate (P4.16), MCP Contract Watch (P3.14), Test-Integrity Gate (P0.3), Agent Failure Modes and Critiques (2026), caseClaimGate(), TestHasGreenTestRun(), TestLooksLikeTestPassClaim(), HasGreenTestRun() (+2 more)

### Community 84 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 85 - "InitDefaultSOPs"
Cohesion: 0.42
Nodes (8): SOP, Dir(), GetSOP(), InitDefaultSOPs(), ListSOPs(), SaveSOP(), TestSOPLifecycle(), ExecuteSOP()

### Community 86 - "clarify.go"
Cohesion: 0.24
Nodes (10): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+2 more)

### Community 87 - "model_catalog.go"
Cohesion: 0.31
Nodes (8): CatalogEntry, CatalogModels(), HasCatalog(), ProviderHasAPIKey(), RemoveProviderModels(), ProviderListKeyboard(), formatProvidersList(), AllProviderNames()

### Community 89 - "HandleTelegramAction"
Cohesion: 0.25
Nodes (10): RequestStop(), FormatUsageStats(), HandleTelegramAction(), GetPath(), MainMenuKeyboard(), MonitorMenuKeyboard(), SettingsMenuKeyboard(), SystemMenuKeyboard() (+2 more)

### Community 90 - "wireCLICallbacks"
Cohesion: 0.29
Nodes (7): formatFileSize(), wireCLICallbacks(), formatFinalResponse(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 99 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 100 - "net/http.Client"
Cohesion: 0.36
Nodes (8): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport()

### Community 101 - "main"
Cohesion: 0.20
Nodes (12): RegisterAutonomous(), hasDebugFlag(), setupCLILogging(), ScorpDir(), envKeyName, isCLIMode(), main(), InitModelUsage() (+4 more)

### Community 102 - "startMCPServer"
Cohesion: 0.33
Nodes (7): ProbeServer(), startMCPServer(), MCPServerConfig, MCPServer, isRemoteMCP(), startSSEServer(), TestRemoteMCPServer_HTTP()

### Community 103 - "SectionSecurity"
Cohesion: 0.53
Nodes (5): EnrichWithGeo(), Flag(), LookupIP(), GeoInfo, SectionSecurity()

### Community 104 - "ExecuteTool"
Cohesion: 0.29
Nodes (13): TestAutoAllowlistPrefixAnchoring(), TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), ResetAutoAllowlist(), ResetAutoStats(), registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended() (+5 more)

### Community 105 - "LoadConfig"
Cohesion: 0.33
Nodes (9): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), Init(), StartServer() (+1 more)

### Community 106 - "GetAllTools"
Cohesion: 0.60
Nodes (5): unregisterMCPNativeTools(), GetAllTools(), countActiveTools(), countDeferredTools(), ExecuteToolList()

### Community 107 - "tools/callbacks.go"
Cohesion: 0.53
Nodes (5): AgentMessage, AutonomousLogEntry, ChatSession, pendingConfirmation, TgResponse

### Community 108 - "SearchResult"
Cohesion: 0.23
Nodes (9): deduplicateAndRank(), normalizeSearchURL(), TestDeduplicateAndRank(), TestLiveWebSearch(), TestMetaSearchAggregator_ConcurrentSuccess(), TestMetaSearchAggregator_FaultTolerance(), TestNormalizeSearchURL(), MockSearchEngine (+1 more)

### Community 109 - "TestGenerateNativeToolsSchema"
Cohesion: 0.40
Nodes (4): IsRateLimitError(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema(), TestIsRateLimitError()

### Community 110 - "TransformToolDefinitions"
Cohesion: 0.39
Nodes (7): ChatRequest, TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions(), ToolSchema

### Community 111 - "NewDefaultMetaSearchAggregator"
Cohesion: 0.73
Nodes (4): GetMetaSearchAggregator(), NewDefaultMetaSearchAggregator(), MetaSearchAggregator, SearchEngine

### Community 113 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 114 - "TestToolSearchActivatesDeferredToolInStaticMode"
Cohesion: 0.70
Nodes (4): deferredToolActive(), osUnsetDynamic(), registerSearchTarget(), TestToolSearchActivatesDeferredToolInStaticMode()

### Community 115 - "Auto-Mode Classifier (P3.13)"
Cohesion: 0.50
Nodes (4): Auto-Mode Classifier (P3.13), Deny-Rule Engine (P0.2), Durable Memory MEMORY.md (P1.7), Claude Code 2026 Reference Architecture

## Knowledge Gaps
- **39 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+34 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 190 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `TaskPlan`, `registerMCPToolsAsNative`, `StartDaemon`, `getAgentSystemPrompt`, `time.Time`, `HomeDir`, `RunAgentSessionLoop`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`, `CheckServerContracts`, `wizard.go`, `HandleConfirmation`, `ExecuteShell`, `FormatHourlyReport`, `RenameSession`, `session_ui.go`, `HandleModelCallback`, `collector_docker.go`, `collector_coolify.go`, `InitDefaultSOPs`, `clarify.go`, `SectionSecurity`, `FormatCompactStats`?**
  _High betweenness centrality (0.111) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `rag_vector.go`, `TaskPlan`, `getAgentSystemPrompt`, `time.Time`, `GetAutonomyLevel`, `CreateCheckpoint`, `TruncOutput`, `TestIntegrityStatus`, `PermissionDecision`, `ResetNativeToolCache`, `HandleConfirmation`, `prepareNewTurnHistory`, `GetStringArg`, `startCLI`, `TruncateStr`, `maybeCompactHistory`, `extractTaskMemory`, `StreamChunk`, `RenameSession`, `ExecuteTermuxAPI`, `testgate.go`, `IsDangerousCommand`, `IsContinuationDirective`, `markdownToTelegramHTML`, `HasGreenTestRun`, `HandleTelegramAction`, `wireCLICallbacks`, `ExecuteTool`, `GenerateContextualSessionTitle`?**
  _High betweenness centrality (0.087) - this node is a cross-community bridge._
- **Why does `TruncateStr()` connect `TruncateStr` to `chat.go`, `cost_router.go`, `runSubagent`, `TaskPlan`, `registerMCPToolsAsNative`, `getAgentSystemPrompt`, `time.Time`, `ScorpPath`, `CallOpenCodeWithTools`, `RunAgentSessionLoop`, `context.Context`, `MCPServer`, `api_commandcode.go`, `SaveModelConfig`, `PermissionDecision`, `api_gemini.go`, `ResolveAPIKey`, `GetStringArg`, `extractTaskMemory`, `RenameSession`?**
  _High betweenness centrality (0.058) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _39 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.05714285714285714 - nodes in this community are weakly interconnected._
- **Should `cost_router.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06663141195134849 - nodes in this community are weakly interconnected._