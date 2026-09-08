# Graph Report - scorp  (2026-09-08)

## Corpus Check
- 273 files · ~190,524 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2138 nodes · 5625 edges · 107 communities (93 shown, 6 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 908 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `983e2794`
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
- runPlanningTurns
- Scorp Agent (Go, ultra-light autonomous agent)
- SaveModelConfig
- Benchmark
- telegram.go
- getAgentSystemPrompt
- time.Time
- CallModelWithFallback
- HomeDir
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
- ToolCall
- CreateCheckpoint
- ExecuteTool
- TruncOutput
- TestIntegrityStatus
- CheckServerContracts
- client.go
- skills.go
- registerMCPToolsAsNative
- GetIntArg
- PermissionDecision
- wizard.go
- checker.go
- ExecuteSQL
- SCORP — BRUTAL END-TO-END TEST PLAN
- GetStringArg
- api_gemini.go
- SandboxActive
- prepareNewTurnHistory
- CallModel
- init
- cost_router.go
- CheckDenyRules
- StartDaemon
- LoadMCPConfig
- startCLI
- TruncateStr
- patch.go
- runLiveCase
- config/hooks.go
- RegisterProvider
- eval.go
- LoadModelConfig
- registerHookProbeTool
- SetTaskPlan
- ExecuteTermuxAPI
- CredentialVault
- ExecuteReadURL
- testgate.go
- HandleModelCallback
- IsDangerousCommand
- TestRegisterPlugin
- TodoManager
- TaskPlan
- ClearTaskPlan
- IsContinuationDirective
- inline.go
- .listenSSEStream
- install.sh
- deploy.sh
- readInteractiveInput
- clarify.go
- TestPhase6_AllTools
- CallModelWithToolsAndFallback
- HasGreenTestRun
- sync.Mutex
- tools/monitor.go
- ExecuteScript
- .CallWithTools
- LoadConfig
- renderStatusFooter
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- serviceBridgeRequests
- RedactSecrets
- ExecuteAutonomous
- TestSteeringQueue
- ShareToMarketplace
- ExecuteAutoLogin
- HandleTelegramAction
- GenerateNativeToolsSchema

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
- `Deny-Rule Engine (P0.2)` --implements--> `CheckDenyRules()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → config/deny.go
- `Compaction Preservation (P2.10)` --implements--> `truncateToolResultsInHistory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/compaction.go
- `Durable Memory MEMORY.md (P1.7)` --implements--> `extractTaskMemory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/memory_md.go
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

## Communities (107 total, 6 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (72): handleMCPCommand(), Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo (+64 more)

### Community 1 - "chat.go"
Cohesion: 0.06
Nodes (73): agentSession, appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown(), convertTableToList() (+65 more)

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
Cohesion: 0.16
Nodes (23): ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), clearDynamicEnv(), registerTempTool(), TestDynamicToolTTL() (+15 more)

### Community 6 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 7 - "runPlanningTurns"
Cohesion: 0.24
Nodes (12): ApprovePlan(), BeginPlanning(), EndPlanning(), AgentMessage, PlanningState(), RevisePlan(), RunPlanningLoop(), runPlanningTurns() (+4 more)

### Community 8 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.05
Nodes (58): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection, MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry (+50 more)

### Community 9 - "SaveModelConfig"
Cohesion: 0.21
Nodes (12): SaveModelConfig(), SwitchModel(), addToFallback(), FallbackText(), FormatAPIKeysList(), ModelInfoText(), ModelMenuText(), moveFallback() (+4 more)

### Community 10 - "Benchmark"
Cohesion: 0.10
Nodes (37): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+29 more)

### Community 11 - "telegram.go"
Cohesion: 0.17
Nodes (20): baseName(), DeleteWebhook(), EditMessage(), EditMessageByID(), MainMenuKeyboard(), MonitorMenuKeyboard(), PollUpdates(), SendChatAction() (+12 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.05
Nodes (50): getSharedMemorySummary(), AppendMemoryMD(), extractTaskMemory(), AgentMessage, parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile() (+42 more)

### Community 13 - "time.Time"
Cohesion: 0.08
Nodes (37): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, time.Time, ScheduledTask, AddTask(), AddTaskEx() (+29 more)

### Community 14 - "CallModelWithFallback"
Cohesion: 0.24
Nodes (12): SummarizeOldToolResult(), ChatResponse, getCheapestModel(), isBudgetExceeded(), RouteModelCostAware(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName() (+4 more)

### Community 15 - "HomeDir"
Cohesion: 0.29
Nodes (14): HomeDir(), ProjectDir(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), HumanSize() (+6 more)

### Community 16 - "testing.T"
Cohesion: 0.06
Nodes (43): TestBuildThinkingMessage(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestIsDangerousCommand(), TestMaxIterations() (+35 more)

### Community 17 - "ScorpPath"
Cohesion: 0.12
Nodes (23): init(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), Hostname(), MemoryFilePath() (+15 more)

### Community 18 - "CallOpenCodeWithTools"
Cohesion: 0.25
Nodes (9): CallOpenCode(), CallOpenCodeStream(), CallOpenCodeWithTools(), resolveOpenCodeBaseURL(), resolveOpenCodeKeyFromDisk(), resolveOpenCodeSessionID(), RecordCostWithCache(), OpenCodeProvider (+1 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.16
Nodes (22): AgentMessage, confirmationDisplay(), ClearStopRequest(), ConsumeStopRequest(), sendScorpReply(), confirmKeyboard(), cleanToolCallTags(), getSessionSearchContext() (+14 more)

### Community 20 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 21 - "compaction_test.go"
Cohesion: 0.14
Nodes (35): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, injectPreservationNote(), latestUserGoal(), maybeCompactHistory(), preservationNote() (+27 more)

### Community 22 - "collector_security.go"
Cohesion: 0.18
Nodes (25): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+17 more)

### Community 23 - "context.Context"
Cohesion: 0.10
Nodes (19): context.Context, AnthropicProvider, applyOpenAIHeaders(), buildOpenAIRequestBody(), CallOpenAI(), CallOpenAIStream(), CallOpenAIWithTools(), formatOpenAIMessages() (+11 more)

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
Cohesion: 0.14
Nodes (14): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, sync.Once, FindMCPTool(), MCPServer, MCPTool, ProbeServer() (+6 more)

### Community 28 - "ToolCall"
Cohesion: 0.15
Nodes (17): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), init(), resolveCommandCodeKeyFromDisk(), commandCodeMsg (+9 more)

### Community 29 - "CreateCheckpoint"
Cohesion: 0.26
Nodes (19): Checkpoint and Rewind (P1.6), caseCheckpoint(), CheckpointDiffStat(), checkpointRepoRoot(), ckptGit(), ckptRefFor(), CreateCheckpoint(), DeleteCheckpoint() (+11 more)

### Community 30 - "ExecuteTool"
Cohesion: 0.12
Nodes (22): TestAutoAllowlistPrefixAnchoring(), TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), ResetAutoAllowlist(), ResetAutoStats(), setAutoMode(), TestExecuteToolAutoDeniesOnNoChannelPath(), TestExecuteToolAutoPresetTrustedAndRecorded() (+14 more)

### Community 31 - "TruncOutput"
Cohesion: 0.34
Nodes (14): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+6 more)

### Community 32 - "TestIntegrityStatus"
Cohesion: 0.29
Nodes (16): Test-Integrity Gate (P0.3), IsTestRelatedPath(), IsTestRunCommand(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestIsTestRunCommand() (+8 more)

### Community 33 - "CheckServerContracts"
Cohesion: 0.20
Nodes (16): MCP Contract Watch (P3.14), Agent Failure Modes and Critiques (2026), caseMCPContract(), CheckServerContracts(), contractFile(), ContractStatusNotice(), ContractWarnings(), serverFingerprint() (+8 more)

### Community 34 - "client.go"
Cohesion: 0.14
Nodes (22): ACPRequest, encoding/json.RawMessage, executeMCPServerTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError() (+14 more)

### Community 35 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 36 - "registerMCPToolsAsNative"
Cohesion: 0.17
Nodes (13): TruncOutputTool(), buildArgDefsFromInputSchema(), MCPToolsDeferred(), rebuildMCPToolList(), registerMCPToolsAsNative(), StartMCPServers(), StopMCPServers(), TestMCPToolsDeferredEnvParsing() (+5 more)

### Community 37 - "GetIntArg"
Cohesion: 0.14
Nodes (11): TestGetBoolArg(), GetBoolArg(), GetIntArg(), getUnclosedTags(), SplitMessage(), ExecuteCompose(), ExecuteGit(), ExecuteHTTP() (+3 more)

### Community 38 - "PermissionDecision"
Cohesion: 0.24
Nodes (12): autoAllowlisted(), autoClassify(), autoClassifyWithModel(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand(), PermissionDecision(), TestAutoFallbackAfterRepeatedUncertain() (+4 more)

### Community 39 - "wizard.go"
Cohesion: 0.34
Nodes (16): AutoPopulateFromCatalog(), ProviderKeyEnv(), ModelMenuKeyboard(), modelWizard, askAPIKey(), ClearModelWizard(), EnvFilePath(), finalizeModelKeySave() (+8 more)

### Community 40 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 41 - "ExecuteSQL"
Cohesion: 0.83
Nodes (3): ExecuteSQL(), loadDBConnections(), dbConnection

### Community 42 - "SCORP — BRUTAL END-TO-END TEST PLAN"
Cohesion: 0.09
Nodes (20): 0. SCALE & SEVERITY DEFINITIONS, A. GATE-STACK ADVERSARIAL (S) — Layer Penetration Probes, B. LONG-HORIZON & CONTEXT COMPACTION (L), C. INFRASTRUCTURE CHAOS (M/L), D. CONCURRENCY & RACE CONDITIONS (M), E. TELEGRAM UX INTEGRITY (M), F. ADVERSARIAL SECURITY (M), G. EVALUATION & DEPLOYMENT GATE INTEGRITY (S) (+12 more)

### Community 43 - "GetStringArg"
Cohesion: 0.21
Nodes (11): init(), GetStringArg(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteSystemInfo(), ExecuteWriteFile(), isGuestLinuxRootfs() (+3 more)

### Community 44 - "api_gemini.go"
Cohesion: 0.20
Nodes (15): callGemini(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), resolveGeminiBaseURL(), geminiContent, geminiFuncDecl (+7 more)

### Community 45 - "SandboxActive"
Cohesion: 0.25
Nodes (15): caseSandbox(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest(), SandboxStatusNotice(), SandboxVersion() (+7 more)

### Community 46 - "prepareNewTurnHistory"
Cohesion: 0.29
Nodes (8): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary()

### Community 47 - "CallModel"
Cohesion: 0.16
Nodes (16): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestOpenCodeProvider_RegistrationAndKey(), CallModel(), CallModelStream(), GetProvider(), TestCoreAndExtendedProvidersRegistered() (+8 more)

### Community 48 - "init"
Cohesion: 0.26
Nodes (10): init(), PythonSitePackages(), GetAllTools(), countActiveTools(), countDeferredTools(), ExecuteToolCall(), ExecuteToolList(), ExecuteToolSearch() (+2 more)

### Community 49 - "cost_router.go"
Cohesion: 0.21
Nodes (17): LoadAutonomousConfig(), setupTestPaths(), TestPhase7_ConfigPersistence(), ConfigMgr(), defaultCostConfig(), formatCostReport(), FormatDailyCostSummary(), handleCostCommand() (+9 more)

### Community 50 - "CheckDenyRules"
Cohesion: 0.24
Nodes (12): TestExecuteToolDenyRulesFirst(), CheckDenyRules(), loadDenyRules(), ParseDenyRule(), ReloadDenyRules(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped() (+4 more)

### Community 51 - "StartDaemon"
Cohesion: 0.21
Nodes (11): UploadsDir(), Init(), StartServer(), StopServer(), InitModelUsage(), runCommandLoop(), StartDaemon(), InitTelegram() (+3 more)

### Community 52 - "LoadMCPConfig"
Cohesion: 0.38
Nodes (11): MCPConfigFilePath(), LoadMCPConfig(), ReloadMCPServers(), sanitizeMCPName(), AddServerEntry(), ExecuteMCPManage(), mcpManageAdd(), mcpManageList() (+3 more)

### Community 53 - "startCLI"
Cohesion: 0.05
Nodes (67): TestStorePendingConfirmationArgs(), clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), StorePendingConfirmation() (+59 more)

### Community 54 - "TruncateStr"
Cohesion: 0.21
Nodes (25): TruncateStr(), anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic(), callAnthropicStream() (+17 more)

### Community 55 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "runLiveCase"
Cohesion: 0.22
Nodes (12): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, evalSandboxDir(), deploymentEnv(), LiveCases(), runLiveCase(), tail(), usageDelta() (+4 more)

### Community 57 - "config/hooks.go"
Cohesion: 0.30
Nodes (13): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+5 more)

### Community 58 - "RegisterProvider"
Cohesion: 0.11
Nodes (19): init(), init(), init(), init(), init(), init(), init(), init() (+11 more)

### Community 59 - "eval.go"
Cohesion: 0.29
Nodes (9): Case, caseResult, CoreCases(), humanCount(), liveLabel(), report(), Run(), TestRunnerAggregationAndFilter() (+1 more)

### Community 60 - "LoadModelConfig"
Cohesion: 0.18
Nodes (14): defaultModelConfig(), LoadModelConfig(), CustomProvider, ModelRouterConfig, ModelUsage, getProviderInfo(), HandleProviderCommand(), listConfiguredModels() (+6 more)

### Community 61 - "registerHookProbeTool"
Cohesion: 0.73
Nodes (5): registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks()

### Community 62 - "SetTaskPlan"
Cohesion: 0.45
Nodes (11): mkTestPlan(), simulateRestart(), TestClearTaskPlanRemovesPersistedFile(), TestPlanFilePathSanitizesSessionID(), TestPlanPersistenceRoundtripAcrossRestart(), TestStalePlanExpiredOnLoad(), TestUpdateItemStatusPersistsThroughRestart(), uniqueSess() (+3 more)

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
Cohesion: 0.14
Nodes (22): OperationalClaim, GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), dedupeOpObjects(), extractOpObjects() (+14 more)

### Community 67 - "HandleModelCallback"
Cohesion: 0.18
Nodes (16): CatalogEntry, CatalogModels(), HasCatalog(), ProviderHasAPIKey(), RemoveProviderModels(), ApiFormatPickerKeyboard(), FallbackEditorKeyboard(), ModelInfoKeyboard() (+8 more)

### Community 68 - "IsDangerousCommand"
Cohesion: 0.43
Nodes (5): devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), IsDangerousCommand(), isHarmlessDevSink()

### Community 69 - "TestRegisterPlugin"
Cohesion: 0.24
Nodes (5): echoPlugin, RegisterPlugin(), TestRegisterPlugin(), ToolPlugin, ToolPluginWithSchema

### Community 70 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 71 - "TaskPlan"
Cohesion: 0.36
Nodes (4): PlanItem, TaskPlan, loadPlanFromDisk(), renderPlanItems()

### Community 72 - "ClearTaskPlan"
Cohesion: 0.44
Nodes (10): CancelPlan(), TestCancelPlanDropsLedger(), ClearTaskPlan(), execTaskPlanTool(), GetTaskPlan(), TestTaskPlanConcurrentAccess(), TestTaskPlanCreateUpdateLifecycle(), TestTaskPlanRender() (+2 more)

### Community 73 - "IsContinuationDirective"
Cohesion: 0.40
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 75 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 76 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 77 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 78 - "deploy.sh"
Cohesion: 0.53
Nodes (4): die(), ok(), deploy.sh script, step()

### Community 79 - "readInteractiveInput"
Cohesion: 0.43
Nodes (6): SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste()

### Community 80 - "clarify.go"
Cohesion: 0.24
Nodes (10): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+2 more)

### Community 81 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 82 - "CallModelWithToolsAndFallback"
Cohesion: 0.36
Nodes (7): TestParseToolCalls(), CallModelWithToolsAndFallback(), IsRateLimitError(), ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls(), TestIsRateLimitError()

### Community 83 - "HasGreenTestRun"
Cohesion: 0.36
Nodes (7): Evidence-Based Claim Gate (P4.16), caseClaimGate(), TestHasGreenTestRun(), TestLooksLikeTestPassClaim(), HasGreenTestRun(), LooksLikeTestPassClaim(), MarkTaskBoundary()

### Community 84 - "sync.Mutex"
Cohesion: 0.33
Nodes (5): autoStats, sync.Mutex, CostTracker, getChatLock(), StartTestEndpoint()

### Community 85 - "tools/monitor.go"
Cohesion: 0.38
Nodes (10): ExecuteMonitor(), InitMonitor(), loadMonitorTargets(), monitorCheckOne(), monitorLoop(), ragIngestText(), sanitizeFilename(), saveMonitorTargets() (+2 more)

### Community 86 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 89 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 90 - "renderStatusFooter"
Cohesion: 0.53
Nodes (5): GetDailyTotalUSD(), getContextPill(), getGitStatus(), getShortCwd(), renderStatusFooter()

### Community 99 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 101 - "ExecuteAutonomous"
Cohesion: 0.33
Nodes (9): SaveAutonomousConfig(), saveAutonomousConfigLocked(), SetKillSwitch(), TestPhase7_KillSwitch(), autoShowActions(), autoShowConfig(), autoShowLog(), autoStatus() (+1 more)

### Community 102 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 105 - "HandleTelegramAction"
Cohesion: 0.23
Nodes (12): RequestStop(), FormatUsageStats(), HandleTelegramAction(), GetPath(), handleTelegramMarketplaceInstall(), MarketplaceTriOptionKeyboard(), reportInstallOutcome(), BackAndRefreshKeyboard() (+4 more)

### Community 110 - "GenerateNativeToolsSchema"
Cohesion: 0.33
Nodes (8): ChatRequest, TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions(), GenerateNativeToolsSchema(), ToolSchema

## Knowledge Gaps
- **38 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+33 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 187 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `runPlanningTurns`, `SaveModelConfig`, `telegram.go`, `getAgentSystemPrompt`, `time.Time`, `HomeDir`, `RunAgentSessionLoop`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`, `CheckServerContracts`, `registerMCPToolsAsNative`, `wizard.go`, `SandboxActive`, `StartDaemon`, `startCLI`, `HandleModelCallback`, `ClearTaskPlan`, `clarify.go`, `TestSteeringQueue`?**
  _High betweenness centrality (0.098) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `registry/registry.go`, `rag_vector.go`, `runPlanningTurns`, `getAgentSystemPrompt`, `CallModelWithFallback`, `compaction_test.go`, `GetAutonomyLevel`, `CreateCheckpoint`, `ExecuteTool`, `TestIntegrityStatus`, `PermissionDecision`, `GetStringArg`, `prepareNewTurnHistory`, `startCLI`, `TruncateStr`, `ExecuteTermuxAPI`, `testgate.go`, `IsDangerousCommand`, `TaskPlan`, `ClearTaskPlan`, `IsContinuationDirective`, `CallModelWithToolsAndFallback`, `HasGreenTestRun`, `TestSteeringQueue`, `HandleTelegramAction`?**
  _High betweenness centrality (0.088) - this node is a cross-community bridge._
- **Why does `TruncateStr()` connect `TruncateStr` to `chat.go`, `agent/autonomous.go`, `runSubagent`, `registerMCPToolsAsNative`, `GetIntArg`, `PermissionDecision`, `runPlanningTurns`, `ExecuteAutonomous`, `getAgentSystemPrompt`, `api_gemini.go`, `time.Time`, `LoadModelConfig`, `CallOpenCodeWithTools`, `RunAgentSessionLoop`, `context.Context`, `MCPServer`, `ToolCall`, `TruncOutput`?**
  _High betweenness centrality (0.066) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _38 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.05730238025271819 - nodes in this community are weakly interconnected._
- **Should `chat.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05600722673893405 - nodes in this community are weakly interconnected._