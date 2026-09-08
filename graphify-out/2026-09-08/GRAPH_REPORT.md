# Graph Report - scorp  (2026-09-08)

## Corpus Check
- 274 files · ~193,817 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2150 nodes · 5655 edges · 109 communities (93 shown, 8 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 912 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0e37b735`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- chat.go
- cost_router.go
- acp.go
- metasearch_engines.go
- execute.go
- rag_vector.go
- TaskPlan
- Scorp Agent (Go, ultra-light autonomous agent)
- LoadMCPConfig
- Benchmark
- HandleTelegramAction
- getAgentSystemPrompt
- time.Time
- HomeDir
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
- CallOpenCodeWithTools
- ExecuteSQL
- TestIntegrityStatus
- CheckServerContracts
- client.go
- skills.go
- SearchResult
- TruncOutput
- PermissionDecision
- RegisterTool
- checker.go
- HandleConfirmation
- SCORP — BRUTAL END-TO-END TEST PLAN
- GetStringArg
- api_gemini.go
- ExecuteShell
- runSubagent
- ResolveAPIKey
- probe_test.go
- init
- CheckDenyRules
- collector_system_native.go
- session_ui.go
- startCLI
- TruncateStr
- patch.go
- runLiveCase
- maybeCompactHistory
- RegisterProvider
- config/hooks.go
- time.Duration
- eval/core.go
- v2_skills.go
- ExecuteTermuxAPI
- CredentialVault
- ExecuteReadURL
- testgate.go
- HandleModelCallback
- IsDangerousCommand
- os.File
- net/http.Client
- collector_docker.go
- VerifyContract
- IsContinuationDirective
- net/http.Request
- .listenSSEStream
- install.sh
- deploy.sh
- metasearch_test.go
- TestPhase6_AllTools
- helpers.go
- RedactSecrets
- RecordToolReceipt
- inline.go
- InitDefaultSOPs
- clarify.go
- WikipediaEngine
- MCPToolsDeferred
- wireCLICallbacks
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- serviceBridgeRequests
- main
- registerHookProbeTool
- LoadConfig
- ExecuteAnalyzeImage
- handleMCPCommand
- TestGenerateNativeToolsSchema
- TransformToolDefinitions
- ExecuteAutoLogin
- GenerateContextualSessionTitle

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
- `Test-Integrity Gate (P0.3)` --implements--> `TestIntegrityStatus()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → tools/testgate.go
- `Auto-Mode Classifier (P3.13)` --implements--> `PermissionDecision()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/auto.go
- `Compaction Preservation (P2.10)` --implements--> `truncateToolResultsInHistory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/compaction.go
- `Plan Mode Workflow (P1.4)` --implements--> `BeginPlanning()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/planmode.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **MCP Marketplace 5-Layer Security Stack** — docs_build_plan_mcp_marketplace_layer1_source_only_registry, docs_build_plan_mcp_marketplace_layer2_ast_prompt_scan, docs_build_plan_mcp_marketplace_layer3_hermetic_ci, docs_build_plan_mcp_marketplace_sha256_pinning, readme_outbound_secret_redactor [EXTRACTED 1.00]
- **Scorp Resilient State & Memory Workflow (PlanMode + TaskLedger + Checkpoint + MEMORY.md)** — docs_implementation_plan_scorp_plan_mode, docs_implementation_plan_scorp_task_ledger_persistence, docs_implementation_plan_scorp_checkpoint_rewind, docs_implementation_plan_scorp_durable_memory [EXTRACTED 1.00]
- **Scorp Trust & Safety Gate Stack (Sandbox + DenyRules + Hooks + AutoClassifier)** — docs_implementation_plan_scorp_sandbox_bwrap, docs_implementation_plan_scorp_deny_rule_engine, docs_implementation_plan_scorp_hooks_lifecycle, docs_implementation_plan_scorp_auto_mode_classifier [EXTRACTED 1.00]
- **Scorp Verification & Integrity Pipeline (TestGate + ClaimGate + EvalArena)** — docs_implementation_plan_scorp_test_integrity_gate, docs_implementation_plan_scorp_claim_gate, docs_implementation_plan_scorp_eval_arena [EXTRACTED 1.00]
- **AI Transpiler Pipeline (Probe -> Generate -> Build -> Verify)** — docs_build_plan_mcp_marketplace_transpiler, docs_build_plan_mcp_marketplace_transpiler_probe_phase, docs_build_plan_mcp_marketplace_transpiler_generate_phase, docs_build_plan_mcp_marketplace_transpiler_build_phase, docs_build_plan_mcp_marketplace_transpiler_verify_phase [EXTRACTED 1.00]
- **Tri-Option Install Convergence onto ~/.scorp/mcp.json** — docs_build_plan_mcp_marketplace_tri_option_install, docs_build_plan_mcp_marketplace_mcp_json_config, docs_build_plan_mcp_marketplace_mcp_manage, docs_build_plan_mcp_marketplace_watchdog, docs_build_plan_mcp_marketplace_tool_registry [EXTRACTED 1.00]

## Communities (109 total, 8 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (74): Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo, RegistryIndex (+66 more)

### Community 1 - "chat.go"
Cohesion: 0.07
Nodes (59): agentSession, appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown(), convertTableToList() (+51 more)

### Community 2 - "cost_router.go"
Cohesion: 0.05
Nodes (66): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+58 more)

### Community 3 - "acp.go"
Cohesion: 0.06
Nodes (40): autoStats, checkACPAvailable(), launchACP(), listAvailableACP(), ACPError, ACPInitializeParams, ACPMessageNewParams, ACPMessagePart (+32 more)

### Community 4 - "metasearch_engines.go"
Cohesion: 0.10
Nodes (15): BingEngine, BraveSearchEngine, DuckDuckGoHTMLEngine, DuckDuckGoLiteEngine, ghSearchResp, GitHubEngine, GoogleCSEEngine, NewBingEngine() (+7 more)

### Community 5 - "execute.go"
Cohesion: 0.19
Nodes (20): runOpenCodeCLI(), runSubagentACP(), AgentMessage, TestParseDelegateParams(), TestParseDelegateParamsCapsAndDefaults(), TestValidateSubagentTools(), delegateBatchParams, delegateResult (+12 more)

### Community 6 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 7 - "TaskPlan"
Cohesion: 0.07
Nodes (50): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary() (+42 more)

### Community 8 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.05
Nodes (58): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection, MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry (+50 more)

### Community 9 - "LoadMCPConfig"
Cohesion: 0.16
Nodes (23): MCPConfigFilePath(), TruncOutputTool(), buildArgDefsFromInputSchema(), LoadMCPConfig(), rebuildMCPToolList(), registerMCPToolsAsNative(), ReloadMCPServers(), sanitizeMCPName() (+15 more)

### Community 10 - "Benchmark"
Cohesion: 0.21
Nodes (19): Generate(), Repair(), stripCodeFences(), detectRuntime(), DumpBenchmark(), isKnownRuntime(), MarshalBenchmark(), Probe() (+11 more)

### Community 11 - "HandleTelegramAction"
Cohesion: 0.10
Nodes (39): RequestStop(), Init(), StartServer(), StopServer(), FormatUsageStats(), HandleTelegramAction(), runCommandLoop(), StartDaemon() (+31 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.07
Nodes (38): getSharedMemorySummary(), AppendMemoryMD(), extractTaskMemory(), AgentMessage, parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile() (+30 more)

### Community 13 - "time.Time"
Cohesion: 0.08
Nodes (40): presentPlanTelegram(), CM(), InitConfigManager(), NewConfigManager(), ConfigManager, time.Time, EscapeHTML(), ScheduledTask (+32 more)

### Community 14 - "HomeDir"
Cohesion: 0.20
Nodes (20): HomeDir(), ProjectDir(), PythonSitePackages(), UploadsDir(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard() (+12 more)

### Community 15 - "ToolCall"
Cohesion: 0.23
Nodes (8): TestParseToolCalls(), formatMessagesForCLI(), init(), ClaudeCliProvider, ToolCall, ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls()

### Community 16 - "testing.T"
Cohesion: 0.06
Nodes (43): TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestIsDangerousCommand() (+35 more)

### Community 17 - "ScorpPath"
Cohesion: 0.12
Nodes (22): init(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), Hostname(), MemoryFilePath() (+14 more)

### Community 18 - "ExecuteTool"
Cohesion: 0.38
Nodes (9): TestAutoAllowlistPrefixAnchoring(), TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), ResetAutoAllowlist(), ResetAutoStats(), ExecuteTool(), caseAutoClassifier(), caseHooksBlockAndContext() (+1 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.21
Nodes (18): AgentMessage, confirmationDisplay(), ClearStopRequest(), ConsumeStopRequest(), sendScorpReply(), cleanToolCallTags(), getSessionSearchContext(), maxIterations() (+10 more)

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
Nodes (23): context.Context, AnthropicProvider, applyOpenAIHeaders(), buildOpenAIRequestBody(), CallOpenAI(), CallOpenAIStream(), CallOpenAIWithTools(), formatOpenAIMessages() (+15 more)

### Community 24 - "GetAutonomyLevel"
Cohesion: 0.20
Nodes (15): ConfirmationRequired(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels() (+7 more)

### Community 25 - "collector_system.go"
Cohesion: 0.11
Nodes (34): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+26 more)

### Community 26 - "RunPreToolUseHooks"
Cohesion: 0.21
Nodes (20): PreToolUse & PostToolUse Hooks (P3.12), Community Praised Agent Patterns (2026), 2026 AI Coding Agent Competitor Landscape, hookPayload, appendHookContext(), runHookCommand(), RunPostToolUseHooks(), RunPreToolUseHooks() (+12 more)

### Community 27 - "MCPServer"
Cohesion: 0.14
Nodes (14): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, sync.Once, FindMCPTool(), MCPServer, MCPTool, ProbeServer() (+6 more)

### Community 28 - "api_commandcode.go"
Cohesion: 0.14
Nodes (18): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeStream(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), resolveCommandCodeKeyFromDisk(), ChatRequest (+10 more)

### Community 29 - "CreateCheckpoint"
Cohesion: 0.14
Nodes (30): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+22 more)

### Community 30 - "CallOpenCodeWithTools"
Cohesion: 0.21
Nodes (10): CallOpenCode(), CallOpenCodeStream(), CallOpenCodeWithTools(), resolveOpenCodeBaseURL(), resolveOpenCodeKeyFromDisk(), resolveOpenCodeSessionID(), RecordCostWithCache(), CostTracker (+2 more)

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
Cohesion: 0.15
Nodes (21): ACPRequest, encoding/json.RawMessage, executeMCPServerTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError() (+13 more)

### Community 35 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 36 - "SearchResult"
Cohesion: 0.24
Nodes (10): deduplicateAndRank(), NewSearXNGEngine(), GetMetaSearchAggregator(), NewDefaultMetaSearchAggregator(), normalizeSearchURL(), TestNormalizeSearchURL(), MetaSearchAggregator, SearchEngine (+2 more)

### Community 37 - "TruncOutput"
Cohesion: 0.34
Nodes (14): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+6 more)

### Community 38 - "PermissionDecision"
Cohesion: 0.22
Nodes (16): autoAllowlisted(), autoClassify(), autoClassifyWithModel(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand(), PermissionDecision(), setAutoMode() (+8 more)

### Community 39 - "RegisterTool"
Cohesion: 0.07
Nodes (35): init(), init(), init(), init(), ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive() (+27 more)

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
Cohesion: 0.20
Nodes (12): init(), GetStringArg(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteSystemInfo(), ExecuteWriteFile(), formatSendFileSize() (+4 more)

### Community 44 - "api_gemini.go"
Cohesion: 0.23
Nodes (18): callGemini(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), parseDataURL(), resolveGeminiBaseURL(), geminiContent (+10 more)

### Community 45 - "ExecuteShell"
Cohesion: 0.23
Nodes (17): Sandbox Shell Bubblewrap (P0.1), caseSandbox(), ExecuteShell(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest() (+9 more)

### Community 46 - "runSubagent"
Cohesion: 0.24
Nodes (10): cleanupSubagentSandbox(), createSubagentSandbox(), defaultIsolation(), formatIsolationInfo(), getSubagentIsolation(), isSubagentToolBlocked(), registerSubagentIsolation(), unregisterSubagentIsolation() (+2 more)

### Community 47 - "ResolveAPIKey"
Cohesion: 0.13
Nodes (28): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestGemini_MultimodalParsing(), TestOpenCodeProvider_RegistrationAndKey(), getCheapestModel(), isBudgetExceeded(), RouteModelCostAware() (+20 more)

### Community 48 - "probe_test.go"
Cohesion: 0.21
Nodes (10): NewSandbox(), runIn(), sanitizeModuleName(), tail(), contains(), goldenFilesystemBinary(), TestDetectRuntimeMessages(), TestProbeAndVerifyGoldenRoundTrip() (+2 more)

### Community 49 - "init"
Cohesion: 0.13
Nodes (16): init(), GetBoolArg(), GetIntArg(), unregisterMCPNativeTools(), GetAllTools(), ExecuteCompose(), countActiveTools(), countDeferredTools() (+8 more)

### Community 50 - "CheckDenyRules"
Cohesion: 0.20
Nodes (14): CheckDenyRules(), loadDenyRules(), ParseDenyRule(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped(), TestCheckDenyRulesMatches(), TestParseDenyRule() (+6 more)

### Community 51 - "collector_system_native.go"
Cohesion: 0.23
Nodes (17): CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses(), GetTopProcesses() (+9 more)

### Community 52 - "session_ui.go"
Cohesion: 0.31
Nodes (11): FormatCompactStats(), CompactStats, BuildSessionMenuKeyboard(), FormatSessionMenuText(), GetActiveSessionID(), HandleSessionCallback(), init(), loadTgSessionMapping() (+3 more)

### Community 53 - "startCLI"
Cohesion: 0.27
Nodes (18): executeOneShot(), executeTurn(), handleCLISession(), handleCLISOP(), hasDebugFlag(), printBanner(), printCLIHelp(), printCostUsage() (+10 more)

### Community 54 - "TruncateStr"
Cohesion: 0.22
Nodes (22): TruncateStr(), anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic(), callAnthropicStream() (+14 more)

### Community 55 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "runLiveCase"
Cohesion: 0.14
Nodes (21): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, Case, caseResult, CoreCases(), evalSandboxDir(), humanCount(), liveLabel() (+13 more)

### Community 57 - "maybeCompactHistory"
Cohesion: 0.25
Nodes (17): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, injectPreservationNote(), latestUserGoal(), maybeCompactHistory(), preservationNote() (+9 more)

### Community 58 - "RegisterProvider"
Cohesion: 0.11
Nodes (19): init(), init(), init(), init(), init(), init(), init(), init() (+11 more)

### Community 59 - "config/hooks.go"
Cohesion: 0.30
Nodes (13): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+5 more)

### Community 60 - "time.Duration"
Cohesion: 0.31
Nodes (12): FormatDuration(), time.Duration, AddUptimeTarget(), checkTarget(), ExecuteUptime(), ListUptimeTargets(), RemoveUptimeTarget(), runUptimeCheck() (+4 more)

### Community 61 - "eval/core.go"
Cohesion: 0.24
Nodes (9): TestExecuteToolDenyRulesFirst(), ReloadDenyRules(), caseDangerGate(), caseDenyInvalidSkipped(), caseDenyShellYOLO(), caseLedgerClear(), caseLedgerPersisted(), caseSensitivePath() (+1 more)

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
Cohesion: 0.52
Nodes (6): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), truncateURLOutput(), tryRemoteScrape()

### Community 66 - "testgate.go"
Cohesion: 0.20
Nodes (16): OperationalClaim, dedupeOpObjects(), extractOpObjects(), isWriteToolName(), LooksLikeOperationalClaims(), opReceipt(), seedOpReceipts(), TestLooksLikeOperationalClaims() (+8 more)

### Community 67 - "HandleModelCallback"
Cohesion: 0.08
Nodes (57): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+49 more)

### Community 68 - "IsDangerousCommand"
Cohesion: 0.43
Nodes (5): devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), IsDangerousCommand(), isHarmlessDevSink()

### Community 69 - "os.File"
Cohesion: 0.19
Nodes (8): acquireSessionLock(), TestAcquireSessionLock(), lockFileExclusive(), unlockFile(), lockFileExclusive(), unlockFile(), os.File, sessionLockFile

### Community 70 - "net/http.Client"
Cohesion: 0.36
Nodes (8): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport()

### Community 71 - "collector_docker.go"
Cohesion: 0.27
Nodes (11): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+3 more)

### Community 72 - "VerifyContract"
Cohesion: 0.53
Nodes (8): contractsFromBenchmark(), contractsFromTools(), diffContract(), paramTypes(), setOf(), stringSlice(), VerifyContract(), ToolContract

### Community 73 - "IsContinuationDirective"
Cohesion: 0.40
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 75 - "net/http.Request"
Cohesion: 0.25
Nodes (13): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), StartGateway() (+5 more)

### Community 76 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 78 - "deploy.sh"
Cohesion: 0.53
Nodes (4): die(), ok(), deploy.sh script, step()

### Community 79 - "metasearch_test.go"
Cohesion: 0.32
Nodes (5): TestDeduplicateAndRank(), TestLiveWebSearch(), TestMetaSearchAggregator_ConcurrentSuccess(), TestMetaSearchAggregator_FaultTolerance(), MockSearchEngine

### Community 80 - "TestPhase6_AllTools"
Cohesion: 0.26
Nodes (14): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession() (+6 more)

### Community 81 - "helpers.go"
Cohesion: 0.67
Nodes (3): GetStringSliceArg(), getUnclosedTags(), SplitMessage()

### Community 83 - "RecordToolReceipt"
Cohesion: 0.20
Nodes (14): Evidence-Based Claim Gate (P4.16), Test-Integrity Gate (P0.3), caseClaimGate(), TestHasGreenTestRun(), TestLooksLikeTestPassClaim(), GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt() (+6 more)

### Community 84 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 85 - "InitDefaultSOPs"
Cohesion: 0.42
Nodes (8): SOP, Dir(), GetSOP(), InitDefaultSOPs(), ListSOPs(), SaveSOP(), TestSOPLifecycle(), ExecuteSOP()

### Community 86 - "clarify.go"
Cohesion: 0.33
Nodes (7): executeClarify(), GetClarifyChatID(), HasPendingClarify(), init(), sendClarifyMessage(), SetClarifyChatID(), PendingClarify

### Community 90 - "wireCLICallbacks"
Cohesion: 0.29
Nodes (8): formatFileSize(), wireCLICallbacks(), formatFinalResponse(), formatTerminalText(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 99 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 101 - "main"
Cohesion: 0.23
Nodes (10): RegisterAutonomous(), ScorpDir(), envKeyName, isCLIMode(), main(), InitModelUsage(), installSystemdService(), maskString() (+2 more)

### Community 104 - "registerHookProbeTool"
Cohesion: 0.73
Nodes (5): registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks()

### Community 105 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 108 - "handleMCPCommand"
Cohesion: 0.50
Nodes (3): handleMCPCommand(), ShareToMarketplace(), toolsForManifest()

### Community 109 - "TestGenerateNativeToolsSchema"
Cohesion: 0.40
Nodes (4): IsRateLimitError(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema(), TestIsRateLimitError()

### Community 110 - "TransformToolDefinitions"
Cohesion: 0.60
Nodes (5): TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions()

### Community 113 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

## Knowledge Gaps
- **39 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+34 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 190 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `TaskPlan`, `LoadMCPConfig`, `time.Time`, `HomeDir`, `RunAgentSessionLoop`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`, `CheckServerContracts`, `HandleConfirmation`, `ExecuteShell`, `session_ui.go`, `v2_skills.go`, `HandleModelCallback`, `collector_docker.go`, `InitDefaultSOPs`, `clarify.go`?**
  _High betweenness centrality (0.105) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `rag_vector.go`, `TaskPlan`, `HandleTelegramAction`, `getAgentSystemPrompt`, `time.Time`, `ExecuteTool`, `GetAutonomyLevel`, `CreateCheckpoint`, `TestIntegrityStatus`, `PermissionDecision`, `RegisterTool`, `HandleConfirmation`, `GetStringArg`, `ResolveAPIKey`, `startCLI`, `TruncateStr`, `maybeCompactHistory`, `ExecuteTermuxAPI`, `testgate.go`, `IsDangerousCommand`, `IsContinuationDirective`, `RecordToolReceipt`, `wireCLICallbacks`, `GenerateContextualSessionTitle`?**
  _High betweenness centrality (0.093) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `chat.go`, `cost_router.go`, `acp.go`, `rag_vector.go`, `LoadMCPConfig`, `getAgentSystemPrompt`, `time.Time`, `HomeDir`, `ScorpPath`, `ExecuteTool`, `client.go`, `skills.go`, `HandleConfirmation`, `v2_skills.go`, `HandleModelCallback`, `IsDangerousCommand`, `collector_docker.go`, `InitDefaultSOPs`, `main`?**
  _High betweenness centrality (0.055) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _39 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.05714285714285714 - nodes in this community are weakly interconnected._
- **Should `chat.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0703962703962704 - nodes in this community are weakly interconnected._