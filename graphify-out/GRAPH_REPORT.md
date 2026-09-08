# Graph Report - scorp  (2026-09-08)

## Corpus Check
- 274 files · ~193,602 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2149 nodes · 5653 edges · 115 communities (101 shown, 6 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 912 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ff2fc604`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- chat.go
- agent/autonomous.go
- runSubagent
- metasearch_engines.go
- registry/registry.go
- rag_vector.go
- TaskPlan
- Scorp Agent (Go, ultra-light autonomous agent)
- registerMCPToolsAsNative
- Benchmark
- HandleTelegramAction
- getAgentSystemPrompt
- time.Time
- HomeDir
- ToolCall
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
- cost_router.go
- ExecuteSQL
- TestIntegrityStatus
- CheckServerContracts
- client.go
- skills.go
- LoadMCPConfig
- TruncOutput
- PermissionDecision
- ResetNativeToolCache
- checker.go
- HandleConfirmation
- SCORP — BRUTAL END-TO-END TEST PLAN
- GetStringArg
- api_gemini.go
- SandboxActive
- prepareNewTurnHistory
- ResolveAPIKey
- sync.Mutex
- GetIntArg
- eval/core.go
- collector_system_native.go
- FormatHourlyReport
- startCLI
- TruncateStr
- patch.go
- eval.go
- maybeCompactHistory
- RegisterProvider
- config/hooks.go
- CallModelWithFallback
- runLiveCase
- v2_skills.go
- ExecuteTermuxAPI
- CredentialVault
- ExecuteReadURL
- testgate.go
- HandleModelCallback
- IsDangerousCommand
- os.File
- TodoManager
- collector_docker.go
- tools/monitor.go
- IsContinuationDirective
- net/http.Request
- .listenSSEStream
- install.sh
- deploy.sh
- handleMCPCommand
- TestPhase6_AllTools
- collector_coolify.go
- RedactSecrets
- RecordToolReceipt
- inline.go
- main
- clarify.go
- ConfigMgr
- testgate_op_test.go
- wireCLICallbacks
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- serviceBridgeRequests
- ExecuteAutonomous
- ScorpDir
- ExecuteScript
- runSelfReview
- ExecuteTool
- LoadConfig
- init
- RegisterPlugin
- ShareToMarketplace
- IsRateLimitError
- TransformToolDefinitions
- ExecuteAutoLogin
- GenerateContextualSessionTitle
- TestToolSearchActivatesDeferredToolInStaticMode
- Auto-Mode Classifier (P3.13)

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

## Communities (115 total, 6 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (75): Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo, RegistryIndex (+67 more)

### Community 1 - "chat.go"
Cohesion: 0.06
Nodes (72): agentSession, appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), ClearStopRequest(), collectTableLines(), convertInlineMarkdown() (+64 more)

### Community 2 - "agent/autonomous.go"
Cohesion: 0.18
Nodes (19): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), makeDecision() (+11 more)

### Community 3 - "runSubagent"
Cohesion: 0.05
Nodes (59): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+51 more)

### Community 4 - "metasearch_engines.go"
Cohesion: 0.05
Nodes (51): net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost(), getClient(), getShortClient() (+43 more)

### Community 5 - "registry/registry.go"
Cohesion: 0.20
Nodes (9): echoPlugin, TestRegisterPlugin(), ExecuteToolByName(), GenerateSystemPromptDescriptions(), GetTool(), GetToolsByCategory(), ArgDef, ToolDef (+1 more)

### Community 6 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 7 - "TaskPlan"
Cohesion: 0.09
Nodes (45): PlanItem, ApprovePlan(), BeginPlanning(), CancelPlan(), EndPlanning(), AgentMessage, PlanningState(), presentPlanTelegram() (+37 more)

### Community 8 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.05
Nodes (58): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection, MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry (+50 more)

### Community 9 - "registerMCPToolsAsNative"
Cohesion: 0.15
Nodes (14): TruncOutputTool(), buildArgDefsFromInputSchema(), MCPToolsDeferred(), rebuildMCPToolList(), registerMCPToolsAsNative(), StartMCPServers(), StopMCPServers(), TestMCPToolsDeferredEnvParsing() (+6 more)

### Community 10 - "Benchmark"
Cohesion: 0.10
Nodes (37): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+29 more)

### Community 11 - "HandleTelegramAction"
Cohesion: 0.09
Nodes (41): RequestStop(), StopMCPServerMode(), Init(), StartServer(), StopServer(), FormatUsageStats(), InitModelUsage(), HandleTelegramAction() (+33 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.08
Nodes (33): getSharedMemorySummary(), AppendMemoryMD(), parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile(), memoryMDTestFile(), nonEmptyLines() (+25 more)

### Community 13 - "time.Time"
Cohesion: 0.08
Nodes (39): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, time.Time, EscapeHTML(), ScheduledTask, AddTask() (+31 more)

### Community 14 - "HomeDir"
Cohesion: 0.21
Nodes (19): HomeDir(), ProjectDir(), PythonSitePackages(), UploadsDir(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard() (+11 more)

### Community 15 - "ToolCall"
Cohesion: 0.50
Nodes (7): TestParseToolCalls(), ToolCall, CallModelWithTools(), CallModelWithToolsAndFallback(), ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls()

### Community 16 - "testing.T"
Cohesion: 0.06
Nodes (44): TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations() (+36 more)

### Community 17 - "ScorpPath"
Cohesion: 0.15
Nodes (19): init(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), Hostname(), MemoryFilePath() (+11 more)

### Community 18 - "RegisterTool"
Cohesion: 0.13
Nodes (14): TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), ResetAutoStats(), init(), init(), init(), init(), TestValidateSubagentTools() (+6 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.20
Nodes (19): AgentMessage, confirmationDisplay(), ConsumeStopRequest(), confirmKeyboard(), cleanToolCallTags(), getSessionSearchContext(), maxIterations(), maxTurnTimeout() (+11 more)

### Community 20 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 21 - "compaction_test.go"
Cohesion: 0.22
Nodes (19): AgentMessage, makeHistory(), makeToolResult(), mkCompactionHistory(), TestEstimateHistoryTokens(), TestPrune_BoundaryAges(), TestPrune_DockerScenario_41Messages(), TestPrune_EmptyHistory() (+11 more)

### Community 22 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 23 - "context.Context"
Cohesion: 0.08
Nodes (17): context.Context, AnthropicProvider, resolveOpenCodeKeyFromDisk(), AzureProvider, ChatMessage, CohereProvider, DashScopeProvider, FastInferenceProvider (+9 more)

### Community 24 - "GetAutonomyLevel"
Cohesion: 0.19
Nodes (16): GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels(), TestConfirmationRequired() (+8 more)

### Community 25 - "collector_system.go"
Cohesion: 0.19
Nodes (21): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+13 more)

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

### Community 30 - "cost_router.go"
Cohesion: 0.19
Nodes (17): defaultCostConfig(), formatCostReport(), FormatDailyCostSummary(), getCheapestModel(), handleCostCommand(), init(), isBudgetExceeded(), isOffPeak() (+9 more)

### Community 31 - "ExecuteSQL"
Cohesion: 0.83
Nodes (3): ExecuteSQL(), loadDBConnections(), dbConnection

### Community 32 - "TestIntegrityStatus"
Cohesion: 0.32
Nodes (15): IsTestRelatedPath(), IsTestRunCommand(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestIsTestRunCommand(), TestTestIntegrityStatus_FailingSuiteDoesNotCount() (+7 more)

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
Cohesion: 0.38
Nodes (11): MCPConfigFilePath(), LoadMCPConfig(), ReloadMCPServers(), sanitizeMCPName(), AddServerEntry(), ExecuteMCPManage(), mcpManageAdd(), mcpManageList() (+3 more)

### Community 37 - "TruncOutput"
Cohesion: 0.34
Nodes (14): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+6 more)

### Community 38 - "PermissionDecision"
Cohesion: 0.21
Nodes (18): TestAutoAllowlistPrefixAnchoring(), autoAllowlisted(), autoClassify(), autoClassifyWithModel(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand(), PermissionDecision() (+10 more)

### Community 39 - "ResetNativeToolCache"
Cohesion: 0.29
Nodes (13): ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), clearDynamicEnv(), registerTempTool(), TestDynamicToolTTL() (+5 more)

### Community 40 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 41 - "HandleConfirmation"
Cohesion: 0.28
Nodes (11): TestStorePendingConfirmationArgs(), clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), StorePendingConfirmation() (+3 more)

### Community 42 - "SCORP — BRUTAL END-TO-END TEST PLAN"
Cohesion: 0.09
Nodes (20): 0. SCALE & SEVERITY DEFINITIONS, A. GATE-STACK ADVERSARIAL (S) — Layer Penetration Probes, B. LONG-HORIZON & CONTEXT COMPACTION (L), C. INFRASTRUCTURE CHAOS (M/L), D. CONCURRENCY & RACE CONDITIONS (M), E. TELEGRAM UX INTEGRITY (M), F. ADVERSARIAL SECURITY (M), G. EVALUATION & DEPLOYMENT GATE INTEGRITY (S) (+12 more)

### Community 43 - "GetStringArg"
Cohesion: 0.18
Nodes (14): init(), Sandbox Shell Bubblewrap (P0.1), GetStringArg(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteShell(), ExecuteSystemInfo() (+6 more)

### Community 44 - "api_gemini.go"
Cohesion: 0.23
Nodes (19): callGemini(), callGeminiStream(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), parseDataURL(), resolveGeminiBaseURL() (+11 more)

### Community 45 - "SandboxActive"
Cohesion: 0.25
Nodes (15): caseSandbox(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest(), SandboxStatusNotice(), SandboxVersion() (+7 more)

### Community 46 - "prepareNewTurnHistory"
Cohesion: 0.29
Nodes (8): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary()

### Community 47 - "ResolveAPIKey"
Cohesion: 0.14
Nodes (21): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestGemini_MultimodalParsing(), TestOpenCodeProvider_RegistrationAndKey(), CallModel(), CallModelStream(), GetProvider() (+13 more)

### Community 48 - "sync.Mutex"
Cohesion: 0.50
Nodes (4): autoStats, sync.Mutex, getChatLock(), StartTestEndpoint()

### Community 49 - "GetIntArg"
Cohesion: 0.15
Nodes (10): GetBoolArg(), GetIntArg(), getUnclosedTags(), SplitMessage(), ExecuteCompose(), ExecuteGit(), ExecuteHTTP(), ExecuteLog() (+2 more)

### Community 50 - "eval/core.go"
Cohesion: 0.22
Nodes (15): CheckDenyRules(), loadDenyRules(), ParseDenyRule(), ReloadDenyRules(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped(), TestCheckDenyRulesMatches() (+7 more)

### Community 51 - "collector_system_native.go"
Cohesion: 0.21
Nodes (18): FormatDuration(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+10 more)

### Community 52 - "FormatHourlyReport"
Cohesion: 0.23
Nodes (12): SystemData, CollectSystem(), GetTopProcesses(), TopProcess, Bar(), bar(), FormatHourlyReport(), FormatStatusResponse() (+4 more)

### Community 53 - "startCLI"
Cohesion: 0.27
Nodes (18): executeOneShot(), executeTurn(), handleCLISession(), handleCLISOP(), hasDebugFlag(), printBanner(), printCLIHelp(), printCostUsage() (+10 more)

### Community 54 - "TruncateStr"
Cohesion: 0.18
Nodes (32): TruncateStr(), anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic(), callAnthropicStream() (+24 more)

### Community 55 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "eval.go"
Cohesion: 0.22
Nodes (12): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, Case, caseResult, CoreCases(), evalSandboxDir(), humanCount(), liveLabel() (+4 more)

### Community 57 - "maybeCompactHistory"
Cohesion: 0.27
Nodes (16): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, injectPreservationNote(), latestUserGoal(), maybeCompactHistory(), preservationNote() (+8 more)

### Community 58 - "RegisterProvider"
Cohesion: 0.09
Nodes (22): init(), init(), formatMessagesForCLI(), init(), init(), init(), init(), init() (+14 more)

### Community 59 - "config/hooks.go"
Cohesion: 0.30
Nodes (13): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+5 more)

### Community 60 - "CallModelWithFallback"
Cohesion: 0.20
Nodes (13): SummarizeOldToolResult(), ChatResponse, GetDailyTotalUSD(), CallModelWithFallback(), findFirstVisionModel(), isVisionModelName(), RouteModel(), ShouldFallbackOnError() (+5 more)

### Community 61 - "runLiveCase"
Cohesion: 0.31
Nodes (9): deploymentEnv(), LiveCases(), runLiveCase(), tail(), usageDelta(), usageSnapshot(), liveCase, TestUsageDeltaClampsNegative() (+1 more)

### Community 62 - "v2_skills.go"
Cohesion: 0.27
Nodes (9): SkillMeta, ActivateSkill(), ListSkillsOverview(), LoadAllSkills(), ParseSkillMetadata(), ReadSkillBody(), scanLegacyJSONSkills(), scanSkillsDirectory() (+1 more)

### Community 63 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 64 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 65 - "ExecuteReadURL"
Cohesion: 0.36
Nodes (7): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), TestReadURL_LocalMock(), truncateURLOutput(), tryRemoteScrape()

### Community 66 - "testgate.go"
Cohesion: 0.38
Nodes (9): OperationalClaim, dedupeOpObjects(), extractOpObjects(), isWriteToolName(), LooksLikeOperationalClaims(), opObjectInReceipt(), opReceiptMatchesClass(), splitClaimSentences() (+1 more)

### Community 67 - "HandleModelCallback"
Cohesion: 0.08
Nodes (56): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+48 more)

### Community 68 - "IsDangerousCommand"
Cohesion: 0.36
Nodes (6): TestIsDangerousCommand(), devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), IsDangerousCommand(), isHarmlessDevSink()

### Community 69 - "os.File"
Cohesion: 0.19
Nodes (8): acquireSessionLock(), TestAcquireSessionLock(), lockFileExclusive(), unlockFile(), lockFileExclusive(), unlockFile(), os.File, sessionLockFile

### Community 70 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 71 - "collector_docker.go"
Cohesion: 0.27
Nodes (11): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+3 more)

### Community 72 - "tools/monitor.go"
Cohesion: 0.38
Nodes (10): ExecuteMonitor(), InitMonitor(), loadMonitorTargets(), monitorCheckOne(), monitorLoop(), ragIngestText(), sanitizeFilename(), saveMonitorTargets() (+2 more)

### Community 73 - "IsContinuationDirective"
Cohesion: 0.40
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 75 - "net/http.Request"
Cohesion: 0.28
Nodes (12): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), TestGatewayEndpoints() (+4 more)

### Community 76 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 78 - "deploy.sh"
Cohesion: 0.53
Nodes (4): die(), ok(), deploy.sh script, step()

### Community 79 - "handleMCPCommand"
Cohesion: 0.31
Nodes (7): handleMCPCommand(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste()

### Community 80 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 81 - "collector_coolify.go"
Cohesion: 0.28
Nodes (12): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+4 more)

### Community 83 - "RecordToolReceipt"
Cohesion: 0.20
Nodes (14): Evidence-Based Claim Gate (P4.16), Test-Integrity Gate (P0.3), caseClaimGate(), TestHasGreenTestRun(), TestLooksLikeTestPassClaim(), GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt() (+6 more)

### Community 84 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 85 - "main"
Cohesion: 0.23
Nodes (12): RegisterAutonomous(), StartGateway(), isCLIMode(), main(), SOP, Dir(), GetSOP(), InitDefaultSOPs() (+4 more)

### Community 86 - "clarify.go"
Cohesion: 0.33
Nodes (7): executeClarify(), GetClarifyChatID(), HasPendingClarify(), init(), sendClarifyMessage(), SetClarifyChatID(), PendingClarify

### Community 87 - "ConfigMgr"
Cohesion: 0.36
Nodes (8): LoadAutonomousConfig(), saveAutonomousConfigLocked(), SetKillSwitch(), setupTestPaths(), TestPhase7_ConfigPersistence(), TestPhase7_KillSwitch(), ConfigMgr(), saveCostTracker()

### Community 89 - "testgate_op_test.go"
Cohesion: 0.32
Nodes (7): opReceipt(), seedOpReceipts(), TestLooksLikeOperationalClaims(), TestLooksLikeOperationalClaimsSkipsVague(), TestOperationalClaimsBounded(), TestRecordToolReceiptCapturesStructuredArgs(), TestUnverifiedOperationalClaims()

### Community 90 - "wireCLICallbacks"
Cohesion: 0.29
Nodes (8): formatFileSize(), wireCLICallbacks(), formatFinalResponse(), formatTerminalText(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 99 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 100 - "ExecuteAutonomous"
Cohesion: 0.52
Nodes (6): SaveAutonomousConfig(), autoShowActions(), autoShowConfig(), autoShowLog(), autoStatus(), ExecuteAutonomous()

### Community 101 - "ScorpDir"
Cohesion: 0.48
Nodes (6): ScorpDir(), envKeyName, installSystemdService(), maskString(), RunQuickstart(), saveConfig()

### Community 102 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 103 - "runSelfReview"
Cohesion: 0.50
Nodes (4): memoryFact, AgentMessage, maybeRunSelfReview(), runSelfReview()

### Community 104 - "ExecuteTool"
Cohesion: 0.44
Nodes (7): TestExecuteToolDenyRulesFirst(), registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks(), ExecuteTool()

### Community 105 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 106 - "init"
Cohesion: 0.26
Nodes (10): init(), unregisterMCPNativeTools(), GetAllTools(), countActiveTools(), countDeferredTools(), ExecuteToolCall(), ExecuteToolList(), ExecuteToolSearch() (+2 more)

### Community 107 - "RegisterPlugin"
Cohesion: 0.83
Nodes (3): RegisterPlugin(), ToolPlugin, ToolPluginWithSchema

### Community 109 - "IsRateLimitError"
Cohesion: 0.50
Nodes (3): IsRateLimitError(), TestCallModelWithToolsNilModel(), TestIsRateLimitError()

### Community 110 - "TransformToolDefinitions"
Cohesion: 0.39
Nodes (7): ChatRequest, TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions(), ToolSchema

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
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `TaskPlan`, `registerMCPToolsAsNative`, `time.Time`, `HomeDir`, `RunAgentSessionLoop`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`, `CheckServerContracts`, `HandleConfirmation`, `SandboxActive`, `FormatHourlyReport`, `v2_skills.go`, `HandleModelCallback`, `collector_docker.go`, `collector_coolify.go`, `main`, `clarify.go`?**
  _High betweenness centrality (0.116) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `rag_vector.go`, `TaskPlan`, `HandleTelegramAction`, `getAgentSystemPrompt`, `time.Time`, `ToolCall`, `RegisterTool`, `GetAutonomyLevel`, `CreateCheckpoint`, `TestIntegrityStatus`, `PermissionDecision`, `ResetNativeToolCache`, `HandleConfirmation`, `GetStringArg`, `prepareNewTurnHistory`, `startCLI`, `TruncateStr`, `maybeCompactHistory`, `CallModelWithFallback`, `ExecuteTermuxAPI`, `testgate.go`, `IsDangerousCommand`, `IsContinuationDirective`, `RecordToolReceipt`, `wireCLICallbacks`, `runSelfReview`, `ExecuteTool`, `GenerateContextualSessionTitle`?**
  _High betweenness centrality (0.085) - this node is a cross-community bridge._
- **Why does `TruncateStr()` connect `TruncateStr` to `chat.go`, `agent/autonomous.go`, `runSubagent`, `ExecuteAutonomous`, `TruncOutput`, `PermissionDecision`, `TaskPlan`, `runSelfReview`, `registerMCPToolsAsNative`, `HandleModelCallback`, `api_gemini.go`, `time.Time`, `ResolveAPIKey`, `GetIntArg`, `RunAgentSessionLoop`, `MCPServer`, `api_commandcode.go`?**
  _High betweenness centrality (0.057) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _39 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.05636114911080711 - nodes in this community are weakly interconnected._
- **Should `chat.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05709876543209876 - nodes in this community are weakly interconnected._