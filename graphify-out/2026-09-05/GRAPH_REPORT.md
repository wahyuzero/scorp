# Graph Report - scorp  (2026-09-05)

## Corpus Check
- 179 files · ~119,981 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1406 nodes · 3454 edges · 64 communities (55 shown, 1 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 459 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4af43eae`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CallCommandCodeWithTools
- ConfigManager
- HandleTelegramAction
- chat.go
- ScorpPath
- client.go
- runSubagent
- rag_vector.go
- collector_system.go
- collector_security.go
- bg.go
- collector_system_native.go
- HandleModelCallback
- cost_router.go
- time.Time
- session_search_fts5.go
- GetStringArg
- StartDaemon
- ResolveAPIKey
- telemetry.go
- checker.go
- startCLI
- install.sh
- ModelConfig
- scorp-agent
- getSession
- 🦂 Scorp
- runSelfReview
- RegisterTool
- metasearch_engines.go
- testing.T
- renderStatusFooter
- skills.go
- .CallWithTools
- RenameSession
- HomeDir
- telegram.go
- RunAgentSessionLoop
- clarify.go
- InitDefaultSOPs
- inline.go
- ToolCall
- main
- runCommandLoop
- context.Context
- TestSteeringQueue
- CleanupSessionsLoop
- ExecuteTermuxAPI
- CredentialVault
- LoadConfig
- ScreenshotsDir
- TestPhase6_AllTools
- api_gemini.go
- GenerateContextualSessionTitle
- GetRecentReceipts
- ExecuteScript

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 61 edges
2. `GetStringArg()` - 43 edges
3. `RunAgentSessionLoop()` - 42 edges
4. `StartDaemon()` - 41 edges
5. `startCLI()` - 40 edges
6. `ModelConfig` - 36 edges
7. `TruncateStr()` - 35 edges
8. `init()` - 32 edges
9. `ChatMessage` - 30 edges
10. `HandleModelCallback()` - 29 edges

## Surprising Connections (you probably didn't know these)
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `serviceBridgeRequests()` --calls--> `ExecuteTool()`  [INFERRED]
  tools/exec_code.go → agent/prompt.go
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go
- `ExecuteShell()` --calls--> `IsDangerousCommand()`  [INFERRED]
  tools/exec.go → agent/safety.go
- `init()` --calls--> `RagVecProvider()`  [EXTRACTED]
  bootstrap/extended.go → rag/rag_vector.go

## Import Cycles
- None detected.

## Communities (64 total, 1 thin omitted)

### Community 0 - "CallCommandCodeWithTools"
Cohesion: 0.16
Nodes (17): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), ChatRequest, commandCodeMsg, commandCodeParams (+9 more)

### Community 1 - "ConfigManager"
Cohesion: 0.12
Nodes (17): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts() (+9 more)

### Community 2 - "HandleTelegramAction"
Cohesion: 0.26
Nodes (15): ClearChatSession(), IsLoopActive(), FormatCompactStats(), CompactStats, DeleteSession(), HandleTelegramAction(), BuildSessionMenuKeyboard(), FormatSessionMenuText() (+7 more)

### Community 3 - "chat.go"
Cohesion: 0.14
Nodes (20): collectTableLines(), convertInlineMarkdown(), convertTableToList(), extractAndSaveMemory(), flushPendingMessages(), GetHistoryTokenEstimate(), historyWriterLoop(), init() (+12 more)

### Community 4 - "ScorpPath"
Cohesion: 0.17
Nodes (18): init(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), Hostname(), MCPConfigFilePath() (+10 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (46): bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool(), buildArgDefsFromInputSchema(), executeMCPServerTool() (+38 more)

### Community 6 - "runSubagent"
Cohesion: 0.08
Nodes (37): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+29 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (53): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+45 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "bg.go"
Cohesion: 0.11
Nodes (26): bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, CostTracker, getChatLock(), StartTestEndpoint(), bgKill() (+18 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.09
Nodes (41): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+33 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (44): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "cost_router.go"
Cohesion: 0.07
Nodes (56): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+48 more)

### Community 14 - "time.Time"
Cohesion: 0.13
Nodes (30): time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult() (+22 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "GetStringArg"
Cohesion: 0.05
Nodes (66): StorePendingConfirmation(), init(), init(), ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract() (+58 more)

### Community 17 - "StartDaemon"
Cohesion: 0.13
Nodes (20): getUnclosedTags(), SplitMessage(), StopMCPServerMode(), Init(), StartServer(), StopServer(), InitModelUsage(), StartDaemon() (+12 more)

### Community 18 - "ResolveAPIKey"
Cohesion: 0.17
Nodes (19): defaultModelConfig(), LoadModelConfig(), ModelRouterConfig, ModelUsage, applyProviderDefaults(), ProviderPreset, hasAPIKey(), migrateModelConfigs() (+11 more)

### Community 19 - "telemetry.go"
Cohesion: 0.67
Nodes (3): CheckModelHealth(), FormatModelList(), FormatModelListWithHealth()

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.08
Nodes (43): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+35 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "ModelConfig"
Cohesion: 0.26
Nodes (14): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModel(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), isVisionModelName() (+6 more)

### Community 31 - "getSession"
Cohesion: 0.26
Nodes (19): agentAutoStop(), appendSessionHistory(), EnterAgentMode(), ExitAgentMode(), getOrCreateSession(), getSession(), getSessionHistory(), getSessionMap() (+11 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "runSelfReview"
Cohesion: 0.10
Nodes (23): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+15 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (33): RegisterAutonomous(), init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), unregisterMCPNativeTools() (+25 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (39): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+31 more)

### Community 36 - "testing.T"
Cohesion: 0.05
Nodes (58): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+50 more)

### Community 37 - "renderStatusFooter"
Cohesion: 0.23
Nodes (11): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste(), getContextPill() (+3 more)

### Community 38 - "skills.go"
Cohesion: 0.24
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 39 - ".CallWithTools"
Cohesion: 0.33
Nodes (4): AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool

### Community 40 - "RenameSession"
Cohesion: 0.36
Nodes (9): historyFilePath(), saveHistoryToDisk(), ListSessions(), RenameSession(), sanitizeSessionID(), SessionExists(), TestSessionManager(), SessionMeta (+1 more)

### Community 41 - "HomeDir"
Cohesion: 0.29
Nodes (14): HomeDir(), ProjectDir(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), HumanSize() (+6 more)

### Community 42 - "telegram.go"
Cohesion: 0.19
Nodes (15): GetPath(), BackAndRefreshKeyboard(), BackButtonKeyboard(), baseName(), MainMenuKeyboard(), MonitorMenuKeyboard(), PollUpdates(), SendDocument() (+7 more)

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.17
Nodes (20): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask() (+12 more)

### Community 44 - "clarify.go"
Cohesion: 0.27
Nodes (9): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID() (+1 more)

### Community 45 - "InitDefaultSOPs"
Cohesion: 0.42
Nodes (8): SOP, Dir(), GetSOP(), InitDefaultSOPs(), ListSOPs(), SaveSOP(), TestSOPLifecycle(), ExecuteSOP()

### Community 46 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 47 - "ToolCall"
Cohesion: 0.24
Nodes (11): TestParseToolCalls(), CustomProvider, ToolCall, CallModelWithTools(), CallModelWithToolsAndFallback(), IsRateLimitError(), ParseAllToolCalls(), ParseCodeBlockFallback() (+3 more)

### Community 48 - "main"
Cohesion: 0.33
Nodes (7): hasDebugFlag(), setupCLILogging(), ScorpDir(), isCLIMode(), main(), RunQuickstart(), saveQuickstartConfig()

### Community 49 - "runCommandLoop"
Cohesion: 0.25
Nodes (6): UploadsDir(), runCommandLoop(), AnswerCallback(), SetupBotCommands(), TestWebhookHandlerMalformedJSON(), WebhookHandler()

### Community 50 - "context.Context"
Cohesion: 0.21
Nodes (21): context.Context, TruncateStr(), callAnthropic(), CallAnthropicWithTools(), CallCommandCodeStream(), callGemini(), CallGeminiWithTools(), geminiDoRequest() (+13 more)

### Community 51 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 52 - "CleanupSessionsLoop"
Cohesion: 0.33
Nodes (5): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), cleanupAgentSessions(), TGDocument

### Community 53 - "ExecuteTermuxAPI"
Cohesion: 0.73
Nodes (5): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification()

### Community 54 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 57 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 58 - "ScreenshotsDir"
Cohesion: 0.40
Nodes (4): ScreenshotsDir(), TestConfigPaths_RagVectorDBPath(), TestConfigPaths_ScorpDir(), TestConfigPaths_ScreenshotsDir()

### Community 59 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 60 - "api_gemini.go"
Cohesion: 0.36
Nodes (10): geminiBuildRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall, geminiGenConfig, geminiPart, geminiRequest (+2 more)

### Community 63 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 65 - "GetRecentReceipts"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 66 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 135 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `client.go`, `RenameSession`, `collector_system.go`, `collector_security.go`, `RunAgentSessionLoop`, `collector_system_native.go`, `HomeDir`, `time.Time`, `telegram.go`, `clarify.go`, `runCommandLoop`, `StartDaemon`, `TestSteeringQueue`, `InitDefaultSOPs`, `startCLI`, `HandleModelCallback`, `getSession`?**
  _High betweenness centrality (0.148) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `runSelfReview`, `RegisterTool`, `client.go`, `skills.go`, `runSubagent`, `rag_vector.go`, `bg.go`, `InitDefaultSOPs`, `ExecuteTermuxAPI`?**
  _High betweenness centrality (0.095) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `runSelfReview`, `RegisterTool`, `chat.go`, `testing.T`, `HandleTelegramAction`, `rag_vector.go`, `RenameSession`, `time.Time`, `ToolCall`, `GetStringArg`, `context.Context`, `TestSteeringQueue`, `startCLI`, `ExecuteTermuxAPI`, `GenerateContextualSessionTitle`, `getSession`?**
  _High betweenness centrality (0.092) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 25 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 25 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `ConfigManager` be split into smaller, more focused modules?**
  _Cohesion score 0.12183908045977011 - nodes in this community are weakly interconnected._