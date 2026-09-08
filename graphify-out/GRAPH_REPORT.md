# Graph Report - scorp  (2026-09-08)

## Corpus Check
- 273 files · ~191,033 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2140 nodes · 5629 edges · 87 communities (75 shown, 4 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 908 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `983e2794`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Manifest
- chat.go
- cost_router.go
- runSubagent
- metasearch_engines.go
- ResetNativeToolCache
- init
- TaskPlan
- Scorp Agent (Go, ultra-light autonomous agent)
- main
- Benchmark
- HandleTelegramAction
- getAgentSystemPrompt
- time.Time
- ToolCall
- eval/core.go
- testing.T
- ScorpPath
- TruncateStr
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
- RegisterTool
- os.File
- TestIntegrityStatus
- CheckServerContracts
- client.go
- skills.go
- LoadMCPConfig
- prompt_test.go
- ExecuteTool
- v2_skills.go
- checker.go
- HandleConfirmation
- SCORP — BRUTAL END-TO-END TEST PLAN
- GetStringArg
- api_gemini.go
- ExecuteShell
- prepareNewTurnHistory
- ResolveAPIKey
- sync.Mutex
- TestGatewayEndpoints
- CheckDenyRules
- InitDefaultSOPs
- wireCLICallbacks
- startCLI
- CallAnthropicWithTools
- patch.go
- runLiveCase
- config/hooks.go
- RegisterProvider
- Auto-Mode Classifier (P3.13)
- EnsureBuiltinSkills
- ExecuteTermuxAPI
- ExecuteReadURL
- testgate.go
- HandleModelCallback
- IsDangerousCommand
- TestRegisterPlugin
- TodoManager
- IsContinuationDirective
- .listenSSEStream
- install.sh
- deploy.sh
- readInteractiveInput
- HasGreenTestRun
- MCP Deferred-by-Default (P2.9)
- scorp-agent
- serviceBridgeRequests
- RedactSecrets
- TransformToolDefinitions

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
- `Durable Memory MEMORY.md (P1.7)` --implements--> `extractTaskMemory()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → agent/memory_md.go
- `Deny-Rule Engine (P0.2)` --implements--> `CheckDenyRules()`  [EXTRACTED]
  docs/IMPLEMENTATION_PLAN_SCORP.md → config/deny.go
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

## Communities (87 total, 4 thin omitted)

### Community 0 - "Manifest"
Cohesion: 0.06
Nodes (72): handleMCPCommand(), Artifact, Build, Contributor, Drift, Health, InstallOption, PortInfo (+64 more)

### Community 1 - "chat.go"
Cohesion: 0.06
Nodes (71): agentSession, appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), ClearStopRequest(), collectTableLines(), convertInlineMarkdown() (+63 more)

### Community 2 - "cost_router.go"
Cohesion: 0.07
Nodes (55): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+47 more)

### Community 3 - "runSubagent"
Cohesion: 0.05
Nodes (59): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+51 more)

### Community 4 - "metasearch_engines.go"
Cohesion: 0.05
Nodes (51): FormatDuration(), net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost(), getClient() (+43 more)

### Community 5 - "ResetNativeToolCache"
Cohesion: 0.20
Nodes (19): ActivateToolWithTTL(), ClearToolTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), clearDynamicEnv(), registerTempTool(), TestDynamicToolTTL() (+11 more)

### Community 6 - "init"
Cohesion: 0.05
Nodes (49): init(), PythonSitePackages(), activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult (+41 more)

### Community 7 - "TaskPlan"
Cohesion: 0.09
Nodes (42): PlanItem, ApprovePlan(), BeginPlanning(), CancelPlan(), EndPlanning(), AgentMessage, PlanningState(), RevisePlan() (+34 more)

### Community 8 - "Scorp Agent (Go, ultra-light autonomous agent)"
Cohesion: 0.05
Nodes (58): CGO + FTS5 Build Requirement, CI Workflow (vet, build, test on push/PR), Linux amd64/arm64 Cross-Compile Matrix, GitHub Release Publish (softprops/action-gh-release), Release Workflow (tag-triggered cross-compile + GitHub release), Updater Version ldflags Injection, MCP Marketplace Build Plan & Execution Roadmap, Security Layer 1: Source-Only Registry (+50 more)

### Community 9 - "main"
Cohesion: 0.22
Nodes (7): RegisterAutonomous(), isCLIMode(), main(), FormatModelListWithHealth(), FormatUsageStats(), InitModelUsage(), SwitchModel()

### Community 10 - "Benchmark"
Cohesion: 0.11
Nodes (33): NewSandbox(), runIn(), sanitizeModuleName(), tail(), Generate(), Repair(), stripCodeFences(), detectRuntime() (+25 more)

### Community 11 - "HandleTelegramAction"
Cohesion: 0.05
Nodes (81): RequestStop(), Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), HomeDir() (+73 more)

### Community 12 - "getAgentSystemPrompt"
Cohesion: 0.07
Nodes (38): getSharedMemorySummary(), AppendMemoryMD(), extractTaskMemory(), AgentMessage, parseMemoryEntries(), ReadMemoryMD(), readMemoryMDLocked(), SetMemoryMDFile() (+30 more)

### Community 13 - "time.Time"
Cohesion: 0.08
Nodes (41): presentPlanTelegram(), CM(), InitConfigManager(), NewConfigManager(), ConfigManager, os.FileMode, time.Time, EscapeHTML() (+33 more)

### Community 14 - "ToolCall"
Cohesion: 0.16
Nodes (21): TestParseToolCalls(), ChatResponse, getCheapestModel(), isBudgetExceeded(), RouteModelCostAware(), CostTracker, CallModelWithFallback(), findFirstVisionModel() (+13 more)

### Community 15 - "eval/core.go"
Cohesion: 0.21
Nodes (13): TestAutoAllowlistPrefixAnchoring(), TestExecuteToolAutoAllowlistReachesExec(), TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(), autoAllowlisted(), ResetAutoStats(), caseAutoClassifier(), caseDenyInvalidSkipped(), caseDenyShellYOLO() (+5 more)

### Community 16 - "testing.T"
Cohesion: 0.08
Nodes (39): registerHookProbeTool(), setHookEnvAgent(), TestExecuteToolHookContextAppended(), TestExecuteToolHookScopedToOtherToolDoesNotFire(), TestExecuteToolPreHookBlocks(), TestSelfReviewCadence(), TestSelfReviewRateLimit(), TestCollectorNative_CollectSystem_Structure() (+31 more)

### Community 17 - "ScorpPath"
Cohesion: 0.06
Nodes (55): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+47 more)

### Community 18 - "TruncateStr"
Cohesion: 0.24
Nodes (22): TruncateStr(), applyAzureHeaders(), callAzure(), callAzureStream(), callAzureWithTools(), resolveAzureEndpoint(), CallCommandCodeWithTools(), applyOpenAIHeaders() (+14 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.16
Nodes (22): AgentMessage, confirmationDisplay(), ConsumeStopRequest(), confirmKeyboard(), cleanToolCallTags(), getSessionSearchContext(), maxIterations(), maxTurnTimeout() (+14 more)

### Community 20 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 21 - "compaction_test.go"
Cohesion: 0.13
Nodes (36): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, injectPreservationNote(), latestUserGoal(), maybeCompactHistory(), preservationNote() (+28 more)

### Community 22 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 23 - "context.Context"
Cohesion: 0.08
Nodes (18): context.Context, AnthropicProvider, CallCommandCode(), CallCommandCodeStream(), resolveOpenCodeKeyFromDisk(), AzureProvider, ChatMessage, CohereProvider (+10 more)

### Community 24 - "GetAutonomyLevel"
Cohesion: 0.19
Nodes (16): GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), PlanningModeActive(), SetAutonomyLevel(), SetPlanningMode(), TestAutonomyLevels(), TestConfirmationRequired() (+8 more)

### Community 25 - "collector_system.go"
Cohesion: 0.05
Nodes (73): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+65 more)

### Community 26 - "RunPreToolUseHooks"
Cohesion: 0.21
Nodes (20): PreToolUse & PostToolUse Hooks (P3.12), Community Praised Agent Patterns (2026), 2026 AI Coding Agent Competitor Landscape, hookPayload, appendHookContext(), runHookCommand(), RunPostToolUseHooks(), RunPreToolUseHooks() (+12 more)

### Community 27 - "MCPServer"
Cohesion: 0.13
Nodes (16): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, sync.Once, FindMCPTool(), MCPServer, MCPTool, ProbeServer() (+8 more)

### Community 28 - "api_commandcode.go"
Cohesion: 0.14
Nodes (16): contextWithTimeout(), handleChat(), net/http.Request, buildCommandCodePayload(), createCommandCodeRequest(), extractFallbackToolCalls(), init(), resolveCommandCodeKeyFromDisk() (+8 more)

### Community 29 - "CreateCheckpoint"
Cohesion: 0.26
Nodes (19): Checkpoint and Rewind (P1.6), caseCheckpoint(), CheckpointDiffStat(), checkpointRepoRoot(), ckptGit(), ckptRefFor(), CreateCheckpoint(), DeleteCheckpoint() (+11 more)

### Community 30 - "RegisterTool"
Cohesion: 0.11
Nodes (14): init(), init(), init(), init(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema(), GenerateSystemPromptDescriptions(), GetToolsByCategory() (+6 more)

### Community 31 - "os.File"
Cohesion: 0.19
Nodes (8): acquireSessionLock(), TestAcquireSessionLock(), lockFileExclusive(), unlockFile(), lockFileExclusive(), unlockFile(), os.File, sessionLockFile

### Community 32 - "TestIntegrityStatus"
Cohesion: 0.29
Nodes (16): Test-Integrity Gate (P0.3), IsTestRelatedPath(), IsTestRunCommand(), shellTouchesTestFile(), mkReceipt(), setTestReceipts(), TestIsTestRelatedPath(), TestIsTestRunCommand() (+8 more)

### Community 33 - "CheckServerContracts"
Cohesion: 0.20
Nodes (16): MCP Contract Watch (P3.14), Agent Failure Modes and Critiques (2026), caseMCPContract(), CheckServerContracts(), contractFile(), ContractStatusNotice(), ContractWarnings(), serverFingerprint() (+8 more)

### Community 34 - "client.go"
Cohesion: 0.15
Nodes (20): ACPRequest, encoding/json.RawMessage, executeMCPServerTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsDeferred(), MCPToolsForPrompt() (+12 more)

### Community 35 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 36 - "LoadMCPConfig"
Cohesion: 0.17
Nodes (21): MCPConfigFilePath(), buildArgDefsFromInputSchema(), LoadMCPConfig(), rebuildMCPToolList(), registerMCPToolsAsNative(), ReloadMCPServers(), sanitizeMCPName(), StartMCPServers() (+13 more)

### Community 37 - "prompt_test.go"
Cohesion: 0.15
Nodes (15): sendScorpReply(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations(), TestTruncOutput() (+7 more)

### Community 38 - "ExecuteTool"
Cohesion: 0.22
Nodes (17): autoClassify(), autoClassifyWithModel(), AutoStatsSnapshot(), bumpAutoStat(), IsReadOnlyShellCommand(), PermissionDecision(), ResetAutoAllowlist(), setAutoMode() (+9 more)

### Community 39 - "v2_skills.go"
Cohesion: 0.27
Nodes (9): SkillMeta, ActivateSkill(), ListSkillsOverview(), LoadAllSkills(), ParseSkillMetadata(), ReadSkillBody(), scanLegacyJSONSkills(), scanSkillsDirectory() (+1 more)

### Community 40 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 41 - "HandleConfirmation"
Cohesion: 0.17
Nodes (16): TestStorePendingConfirmationArgs(), clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), StorePendingConfirmation() (+8 more)

### Community 42 - "SCORP — BRUTAL END-TO-END TEST PLAN"
Cohesion: 0.09
Nodes (20): 0. SCALE & SEVERITY DEFINITIONS, A. GATE-STACK ADVERSARIAL (S) — Layer Penetration Probes, B. LONG-HORIZON & CONTEXT COMPACTION (L), C. INFRASTRUCTURE CHAOS (M/L), D. CONCURRENCY & RACE CONDITIONS (M), E. TELEGRAM UX INTEGRITY (M), F. ADVERSARIAL SECURITY (M), G. EVALUATION & DEPLOYMENT GATE INTEGRITY (S) (+12 more)

### Community 43 - "GetStringArg"
Cohesion: 0.11
Nodes (19): TestGetIntArg(), init(), GetIntArg(), GetStringArg(), SendDocumentBytes(), ExecuteCompose(), ExecuteListDir(), ExecuteReadFile() (+11 more)

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
Cohesion: 0.17
Nodes (20): TestGemini_LiveModels(), TestGemini_LiveStreaming(), TestGemini_LiveToolCalling(), TestOpenCodeProvider_RegistrationAndKey(), CallModel(), CallModelStream(), GetProvider(), TestCoreAndExtendedProvidersRegistered() (+12 more)

### Community 48 - "sync.Mutex"
Cohesion: 0.26
Nodes (10): autoStats, sync.Mutex, unregisterMCPNativeTools(), GetAllTools(), getChatLock(), StartTestEndpoint(), countActiveTools(), countDeferredTools() (+2 more)

### Community 49 - "TestGatewayEndpoints"
Cohesion: 0.36
Nodes (8): handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), StartGateway(), TestGatewayEndpoints(), net/http.ResponseWriter

### Community 50 - "CheckDenyRules"
Cohesion: 0.24
Nodes (12): TestExecuteToolDenyRulesFirst(), CheckDenyRules(), loadDenyRules(), ParseDenyRule(), ReloadDenyRules(), resetDenyRules(), TestCheckDenyRulesHoldInYOLO(), TestCheckDenyRulesInvalidSpecsSkipped() (+4 more)

### Community 51 - "InitDefaultSOPs"
Cohesion: 0.42
Nodes (8): SOP, Dir(), GetSOP(), InitDefaultSOPs(), ListSOPs(), SaveSOP(), TestSOPLifecycle(), ExecuteSOP()

### Community 52 - "wireCLICallbacks"
Cohesion: 0.31
Nodes (6): wireCLICallbacks(), formatFinalResponse(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 53 - "startCLI"
Cohesion: 0.26
Nodes (19): executeOneShot(), executeTurn(), formatTerminalText(), handleCLISession(), handleCLISOP(), hasDebugFlag(), printBanner(), printCLIHelp() (+11 more)

### Community 54 - "CallAnthropicWithTools"
Cohesion: 0.27
Nodes (11): anthropicResponse, anthropicTool, applyAnthropicHeaders(), buildAnthropicMessages(), buildAnthropicRequestBody(), callAnthropic(), callAnthropicStream(), CallAnthropicWithTools() (+3 more)

### Community 55 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 56 - "runLiveCase"
Cohesion: 0.14
Nodes (21): Scorp Eval Arena (P4.15), Merge-Rate Verification Mindset, Case, caseResult, CoreCases(), evalSandboxDir(), humanCount(), liveLabel() (+13 more)

### Community 57 - "config/hooks.go"
Cohesion: 0.30
Nodes (13): HookEntry, HookMatches(), loadHooks(), parseHookEnv(), ParseHookSpec(), PostToolHooks(), PreToolHooks(), ReloadHooks() (+5 more)

### Community 58 - "RegisterProvider"
Cohesion: 0.10
Nodes (20): init(), init(), formatMessagesForCLI(), init(), init(), init(), init(), init() (+12 more)

### Community 59 - "Auto-Mode Classifier (P3.13)"
Cohesion: 0.50
Nodes (4): Auto-Mode Classifier (P3.13), Deny-Rule Engine (P0.2), Durable Memory MEMORY.md (P1.7), Claude Code 2026 Reference Architecture

### Community 63 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 65 - "ExecuteReadURL"
Cohesion: 0.36
Nodes (7): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), TestReadURL_LocalMock(), truncateURLOutput(), tryRemoteScrape()

### Community 66 - "testgate.go"
Cohesion: 0.14
Nodes (22): OperationalClaim, GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), dedupeOpObjects(), extractOpObjects() (+14 more)

### Community 67 - "HandleModelCallback"
Cohesion: 0.07
Nodes (57): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+49 more)

### Community 68 - "IsDangerousCommand"
Cohesion: 0.36
Nodes (6): TestIsDangerousCommand(), devOverwriteTarget(), TestDevNullRedirectionNotDangerous(), TestDevTcpPseudoDeviceAllowed(), IsDangerousCommand(), isHarmlessDevSink()

### Community 69 - "TestRegisterPlugin"
Cohesion: 0.22
Nodes (6): echoPlugin, RegisterPlugin(), TestRegisterPlugin(), ExecuteToolByName(), ToolPlugin, ToolPluginWithSchema

### Community 70 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

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

### Community 83 - "HasGreenTestRun"
Cohesion: 0.36
Nodes (7): Evidence-Based Claim Gate (P4.16), caseClaimGate(), TestHasGreenTestRun(), TestLooksLikeTestPassClaim(), HasGreenTestRun(), LooksLikeTestPassClaim(), MarkTaskBoundary()

### Community 99 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 110 - "TransformToolDefinitions"
Cohesion: 0.39
Nodes (7): ChatRequest, TestToolSchemaTransform(), NormalizeToolSchemaTransform(), sanitizeProperty(), SanitizeToolSchema(), TransformToolDefinitions(), ToolSchema

## Knowledge Gaps
- **38 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+33 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 187 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `Manifest`, `chat.go`, `CheckServerContracts`, `HandleModelCallback`, `LoadMCPConfig`, `TaskPlan`, `v2_skills.go`, `HandleConfirmation`, `main`, `time.Time`, `ExecuteShell`, `RunAgentSessionLoop`, `InitDefaultSOPs`, `collector_security.go`, `GetAutonomyLevel`, `collector_system.go`, `CreateCheckpoint`?**
  _High betweenness centrality (0.106) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `ResetNativeToolCache`, `init`, `TaskPlan`, `HandleTelegramAction`, `getAgentSystemPrompt`, `time.Time`, `ToolCall`, `eval/core.go`, `TruncateStr`, `compaction_test.go`, `GetAutonomyLevel`, `CreateCheckpoint`, `TestIntegrityStatus`, `prompt_test.go`, `ExecuteTool`, `HandleConfirmation`, `GetStringArg`, `prepareNewTurnHistory`, `wireCLICallbacks`, `startCLI`, `ExecuteTermuxAPI`, `testgate.go`, `IsDangerousCommand`, `IsContinuationDirective`, `HasGreenTestRun`?**
  _High betweenness centrality (0.086) - this node is a cross-community bridge._
- **Why does `TruncateStr()` connect `TruncateStr` to `chat.go`, `cost_router.go`, `runSubagent`, `HandleModelCallback`, `prompt_test.go`, `ExecuteTool`, `TaskPlan`, `GetStringArg`, `getAgentSystemPrompt`, `api_gemini.go`, `time.Time`, `ResolveAPIKey`, `ScorpPath`, `RunAgentSessionLoop`, `CallAnthropicWithTools`, `context.Context`, `MCPServer`?**
  _High betweenness centrality (0.063) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _38 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Manifest` be split into smaller, more focused modules?**
  _Cohesion score 0.05730238025271819 - nodes in this community are weakly interconnected._
- **Should `chat.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05759493670886076 - nodes in this community are weakly interconnected._