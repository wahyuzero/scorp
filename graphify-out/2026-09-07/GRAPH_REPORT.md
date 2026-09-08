# Graph Report - scorp  (2026-09-07)

## Corpus Check
- 256 files · ~184,376 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2015 nodes · 5207 edges · 86 communities (75 shown, 3 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 818 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e3d0d848`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- chat.go
- cost_router.go
- GetIntArg
- metasearch_engines.go
- registry/registry.go
- init
- TaskPlan
- Scorp Agent (Go, ultra-light autonomous agent)
- HandleModelCallback
- Benchmark
- collector_system_native.go
- getAgentSystemPrompt
- time.Time
- ModelConfig
- HandleTelegramAction
- testing.T
- ScorpPath
- context.Context
- RunAgentSessionLoop
- session_search_fts5.go
- compaction_test.go
- collector_security.go
- LoadModelConfig
- GetAutonomyLevel
- collector_system.go
- RunPreToolUseHooks
- MCPServer
- CallCommandCodeWithTools
- CreateCheckpoint
- PermissionDecision
- ToolCall
- startCLI
- CheckServerContracts
- client.go
- skills.go
- registerMCPToolsAsNative
- GetBoolArg
- CheckDenyRules
- ExecuteShell
- checker.go
- HandleConfirmation
- SCORP — BRUTAL END-TO-END TEST PLAN
- GetStringArg
- api_gemini.go
- collector_docker.go
- prepareNewTurnHistory
- TruncateStr
- TestIntegrityStatus
- HasGreenTestRun
- bg.go
- FormatHourlyReport
- LoadMCPConfig
- config/hooks.go
- eval/core.go
- patch.go
- runLiveCase
- collector_coolify.go
- GetProvider
- collector_native_test.go
- saveQuickstartConfig
- ExecuteTool
- ExecuteTermuxAPI
- RecordToolReceipt
- ExecuteReadURL
- testgate.go
- IsDangerousCommand
- RegisterTool
- IsContinuationDirective
- .listenSSEStream
- install.sh
- deploy.sh
- readInteractiveInput
- GenerateContextualSessionTitle
- serviceBridgeRequests
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- ConfirmationRequired

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 81 edges
2. `RunAgentSessionLoop()` - 63 edges
3. `startCLI()` - 51 edges
4. `GetStringArg()` - 47 edges
5. `StartDaemon()` - 43 edges
6. `ScorpPath()` - 42 edges
7. `TruncateStr()` - 42 edges
8. `ModelConfig` - 42 edges
9. `resumeAgentLoop()` - 40 edges
10. `ChatMessage` - 35 edges

## Surprising Connections (you probably didn't know these)
- `Durable Memory MEMORY.md (P1.7)` --implements--> `extractTaskMemory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/memory_md.go
- `Auto-Mode Classifier (P3.13)` --implements--> `PermissionDecision()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/auto.go
- `Compaction Preservation (P2.10)` --implements--> `truncateToolResultsInHistory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/compaction.go
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

## Communities (86 total, 3 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (74): Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo, RegistryIndex (+66 more)

### Community 1 - "chat.go"
Cohesion: 0.05
Nodes (80): agentSession, appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown(), convertTableToList() (+72 more)

### Community 2 - "cost_router.go"
Cohesion: 0.06
Nodes (58): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+50 more)

### Community 3 - "GetIntArg"
Cohesion: 0.07
Nodes (43): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+35 more)

### Community 4 - "metasearch_engines.go"
Cohesion: 0.05
Nodes (50): net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost(), getClient(), getShortClient() (+42 more)

### Community 5 - "registry/registry.go"
Cohesion: 0.09
Nodes (37): executeMCPServerTool(), unregisterMCPNativeTools(), ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), clearDynamicEnv() (+29 more)

### Community 6 - "init"
Cohesion: 0.05
Nodes (49): init(), PythonSitePackages(), activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult (+41 more)

### Community 7 - "TaskPlan"
Cohesion: 0.10
Nodes (40): PlanItem, ApprovePlan(), BeginPlanning(), CancelPlan(), EndPlanning(), AgentMessage, PlanningState(), presentPlanTelegram() (+32 more)

### Community 8 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.05
Nodes (58): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection, MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry (+50 more)

### Community 9 - "HandleModelCallback"
Cohesion: 0.10
Nodes (45): ProjectDir(), CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv() (+37 more)

### Community 10 - "Benchmark"
Cohesion: 0.11
Nodes (33): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+25 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.20
Nodes (18): TestCollectorNative_StartCPUSampler_DoesNotBlock(), FormatDuration(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes() (+10 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.06
Nodes (46): getSharedMemorySummary(), AppendMemoryMD(), parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile(), memoryMDTestFile(), nonEmptyLines() (+38 more)

### Community 13 - "time.Time"
Cohesion: 0.08
Nodes (39): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, time.Time, EscapeHTML(), ScheduledTask, AddTask() (+31 more)

### Community 14 - "ModelConfig"
Cohesion: 0.18
Nodes (19): autoClassifyWithModel(), extractTaskMemory(), AgentMessage, ChatRequest, ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModel() (+11 more)

### Community 15 - "HandleTelegramAction"
Cohesion: 0.06
Nodes (74): RequestStop(), Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), HomeDir() (+66 more)

### Community 16 - "testing.T"
Cohesion: 0.08
Nodes (34): TestBuildThinkingMessage(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestIsDangerousCommand(), TestMaxIterations(), TestToolDescription() (+26 more)

### Community 17 - "ScorpPath"
Cohesion: 0.06
Nodes (55): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+47 more)

### Community 18 - "context.Context"
Cohesion: 0.14
Nodes (10): context.Context, AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, CallCommandCode(), ChatMessage, GeminiProvider (+2 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.24
Nodes (16): AgentMessage, confirmationDisplay(), ClearStopRequest(), ConsumeStopRequest(), cleanToolCallTags(), getSessionSearchContext(), maxIterations(), maxTurnTimeout() (+8 more)

### Community 20 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 21 - "compaction_test.go"
Cohesion: 0.14
Nodes (26): compactionTestSetup(), AgentMessage, makeHistory(), makeToolResult(), mkCompactionHistory(), TestActiveLoopCompactionPreservesContext(), TestCompactionNoopUnderThreshold(), TestEstimateHistoryTokens() (+18 more)

### Community 22 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 23 - "LoadModelConfig"
Cohesion: 0.17
Nodes (17): defaultModelConfig(), LoadModelConfig(), ModelRouterConfig, ModelUsage, applyProviderDefaults(), ProviderPreset, migrateModelConfigs(), ResolveBaseURL() (+9 more)

### Community 24 - "GetAutonomyLevel"
Cohesion: 0.18
Nodes (16): GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels(), TestConfirmationRequired() (+8 more)

### Community 25 - "collector_system.go"
Cohesion: 0.21
Nodes (20): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+12 more)

### Community 26 - "RunPreToolUseHooks"
Cohesion: 0.21
Nodes (20): PreToolUse & PostToolUse Hooks (P3.12), Community Praised Agent Patterns (2026), 2026 AI Coding Agent Competitor Landscape, hookPayload, appendHookContext(), runHookCommand(), RunPostToolUseHooks(), RunPreToolUseHooks() (+12 more)

### Community 27 - "MCPServer"
Cohesion: 0.14
Nodes (14): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, sync.Once, FindMCPTool(), MCPServer, MCPTool, ProbeServer() (+6 more)

### Community 28 - "CallCommandCodeWithTools"
Cohesion: 0.23
Nodes (13): buildCommandCodePayload(), CallCommandCodeStream(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), commandCodeMsg, commandCodeParams, commandCodePayload (+5 more)

### Community 29 - "CreateCheckpoint"
Cohesion: 0.26
Nodes (19): Checkpoint and Rewind (P1.6), caseCheckpoint(), CheckpointDiffStat(), checkpointRepoRoot(), ckptGit(), ckptRefFor(), CreateCheckpoint(), DeleteCheckpoint() (+11 more)

### Community 30 - "PermissionDecision"
Cohesion: 0.22
Nodes (17): TestAutoAllowlistPrefixAnchoring(), autoAllowlisted(), autoClassify(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand(), PermissionDecision(), ResetAutoAllowlist() (+9 more)

### Community 31 - "ToolCall"
Cohesion: 0.24
Nodes (11): TestParseToolCalls(), CustomProvider, ToolCall, CallModelWithTools(), CallModelWithToolsAndFallback(), IsRateLimitError(), ParseAllToolCalls(), ParseCodeBlockFallback() (+3 more)

### Community 32 - "startCLI"
Cohesion: 0.05
Nodes (61): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue(), RegisterAutonomous(), wireCLICallbacks(), executeOneShot() (+53 more)

### Community 33 - "CheckServerContracts"
Cohesion: 0.20
Nodes (16): MCP Contract Watch (P3.14), Agent Failure Modes and Critiques (2026), caseMCPContract(), CheckServerContracts(), contractFile(), ContractStatusNotice(), ContractWarnings(), serverFingerprint() (+8 more)

### Community 34 - "client.go"
Cohesion: 0.16
Nodes (19): ACPRequest, encoding/json.RawMessage, getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError(), sendMCPResult() (+11 more)

### Community 35 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 36 - "registerMCPToolsAsNative"
Cohesion: 0.18
Nodes (12): TruncOutputTool(), buildArgDefsFromInputSchema(), MCPToolsDeferred(), registerMCPToolsAsNative(), StopMCPServers(), TestMCPToolsDeferredEnvParsing(), ServerWatchdog, GetServerHealthStatus() (+4 more)

### Community 37 - "GetBoolArg"
Cohesion: 0.17
Nodes (8): sendScorpReply(), TestGetBoolArg(), GetBoolArg(), getUnclosedTags(), SplitMessage(), ExecuteCompose(), ExecuteHTTP(), ExecuteLog()

### Community 38 - "CheckDenyRules"
Cohesion: 0.20
Nodes (15): CheckDenyRules(), loadDenyRules(), ParseDenyRule(), ReloadDenyRules(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped(), TestCheckDenyRulesMatches() (+7 more)

### Community 39 - "ExecuteShell"
Cohesion: 0.23
Nodes (17): Sandbox Shell Bubblewrap (P0.1), caseSandbox(), ExecuteShell(), SandboxActive(), sandboxArgv(), SandboxModeEnabled(), sandboxRWPaths(), sandboxSmokeTest() (+9 more)

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
Cohesion: 0.23
Nodes (11): init(), GetStringArg(), SendDocumentBytes(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteSystemInfo(), ExecuteWriteFile() (+3 more)

### Community 44 - "api_gemini.go"
Cohesion: 0.32
Nodes (13): callGemini(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall (+5 more)

### Community 45 - "collector_docker.go"
Cohesion: 0.27
Nodes (11): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+3 more)

### Community 46 - "prepareNewTurnHistory"
Cohesion: 0.29
Nodes (8): AgentMessage, prepareNewTurnHistory(), AgentMessage, heavyHistory(), TestRollUpSkipsActivePlan(), TestRollUpSkipsContinuation(), TestRollUpSkipsSmallHistory(), TestRollUpTrimsAndMarksBoundary()

### Community 47 - "TruncateStr"
Cohesion: 0.22
Nodes (24): TruncateStr(), callAnthropic(), CallAnthropicWithTools(), CallOpenAI(), CallOpenAIWithTools(), formatOpenAIMessages(), CallOpenCode(), CallOpenCodeStream() (+16 more)

### Community 48 - "TestIntegrityStatus"
Cohesion: 0.29
Nodes (16): Test-Integrity Gate (P0.3), IsTestRelatedPath(), IsTestRunCommand(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestIsTestRunCommand() (+8 more)

### Community 49 - "HasGreenTestRun"
Cohesion: 0.36
Nodes (7): Evidence-Based Claim Gate (P4.16), caseClaimGate(), TestHasGreenTestRun(), TestLooksLikeTestPassClaim(), HasGreenTestRun(), LooksLikeTestPassClaim(), MarkTaskBoundary()

### Community 50 - "bg.go"
Cohesion: 0.11
Nodes (27): autoStats, bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, CostTracker, getChatLock(), StartTestEndpoint() (+19 more)

### Community 51 - "FormatHourlyReport"
Cohesion: 0.29
Nodes (12): SystemData, TopProcess, Bar(), bar(), FormatHourlyReport(), FormatStatusResponse(), SectionCoolify(), SectionDocker() (+4 more)

### Community 52 - "LoadMCPConfig"
Cohesion: 0.31
Nodes (13): MCPConfigFilePath(), LoadMCPConfig(), rebuildMCPToolList(), ReloadMCPServers(), sanitizeMCPName(), StartMCPServers(), AddServerEntry(), ExecuteMCPManage() (+5 more)

### Community 53 - "config/hooks.go"
Cohesion: 0.30
Nodes (13): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+5 more)

### Community 54 - "eval/core.go"
Cohesion: 0.53
Nodes (5): caseDenyInvalidSkipped(), caseDenyShellYOLO(), caseHooksBlockAndContext(), caseSensitivePath(), withEnv()

### Community 55 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "runLiveCase"
Cohesion: 0.14
Nodes (21): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, Case, caseResult, CoreCases(), evalSandboxDir(), humanCount(), liveLabel() (+13 more)

### Community 57 - "collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 58 - "GetProvider"
Cohesion: 0.53
Nodes (5): init(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 59 - "collector_native_test.go"
Cohesion: 0.23
Nodes (11): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+3 more)

### Community 61 - "ExecuteTool"
Cohesion: 0.44
Nodes (7): TestExecuteToolDenyRulesFirst(), registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks(), ExecuteTool()

### Community 63 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 64 - "RecordToolReceipt"
Cohesion: 0.26
Nodes (8): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), RedactSecrets(), TestRedactSecrets(), ToolReceipt

### Community 65 - "ExecuteReadURL"
Cohesion: 0.36
Nodes (7): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), TestReadURL_LocalMock(), truncateURLOutput(), tryRemoteScrape()

### Community 66 - "testgate.go"
Cohesion: 0.20
Nodes (16): OperationalClaim, dedupeOpObjects(), extractOpObjects(), isWriteToolName(), LooksLikeOperationalClaims(), opReceipt(), seedOpReceipts(), TestLooksLikeOperationalClaims() (+8 more)

### Community 68 - "IsDangerousCommand"
Cohesion: 0.43
Nodes (5): devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), IsDangerousCommand(), isHarmlessDevSink()

### Community 70 - "RegisterTool"
Cohesion: 0.13
Nodes (14): TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), ResetAutoStats(), init(), init(), init(), init(), TestValidateSubagentTools() (+6 more)

### Community 73 - "IsContinuationDirective"
Cohesion: 0.40
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

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
Cohesion: 0.23
Nodes (11): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste(), getContextPill() (+3 more)

### Community 81 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 84 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 99 - "ConfirmationRequired"
Cohesion: 0.28
Nodes (6): ConfirmationRequired(), ExecuteSQL(), loadDBConnections(), dbConnection, ExecuteGit(), ExecuteProcess()

## Knowledge Gaps
- **35 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+30 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 173 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `TaskPlan`, `HandleModelCallback`, `collector_system_native.go`, `getAgentSystemPrompt`, `time.Time`, `RunAgentSessionLoop`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`, `startCLI`, `CheckServerContracts`, `registerMCPToolsAsNative`, `ExecuteShell`, `HandleConfirmation`, `collector_docker.go`, `FormatHourlyReport`, `collector_coolify.go`?**
  _High betweenness centrality (0.127) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `registry/registry.go`, `init`, `TaskPlan`, `getAgentSystemPrompt`, `time.Time`, `ModelConfig`, `HandleTelegramAction`, `GetAutonomyLevel`, `CreateCheckpoint`, `PermissionDecision`, `ToolCall`, `startCLI`, `GetBoolArg`, `HandleConfirmation`, `GetStringArg`, `prepareNewTurnHistory`, `TruncateStr`, `TestIntegrityStatus`, `HasGreenTestRun`, `ExecuteTool`, `ExecuteTermuxAPI`, `testgate.go`, `IsDangerousCommand`, `RegisterTool`, `IsContinuationDirective`, `GenerateContextualSessionTitle`, `ConfirmationRequired`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Why does `init()` connect `init` to `startCLI`, `ExecuteReadURL`, `skills.go`, `GetIntArg`, `GetBoolArg`, `RegisterTool`, `registry/registry.go`, `getAgentSystemPrompt`, `time.Time`, `bg.go`, `LoadMCPConfig`, `patch.go`, `ExecuteTermuxAPI`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **Are the 34 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `confirmationDisplay()` and `PermissionDecision()`) actually correct?**
  _`RunAgentSessionLoop()` has 34 INFERRED edges - model-reasoned connections that need verification._
- **Are the 7 inferred relationships involving `startCLI()` (e.g. with `wireCLICallbacks()` and `formatTerminalText()`) actually correct?**
  _`startCLI()` has 7 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _35 weakly-connected nodes found - possible documentation gaps or missing edges._