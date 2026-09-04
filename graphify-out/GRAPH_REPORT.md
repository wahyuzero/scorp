# Graph Report - scorp  (2026-09-05)

## Corpus Check
- 176 files · ~118,586 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1392 nodes · 3411 edges · 58 communities (48 shown, 2 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 445 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b85bf168`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- context.Context
- HomeDir
- testing.T
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
- HandleTelegramAction
- StartDaemon
- GetAllTools
- MCPServer
- checker.go
- startCLI
- install.sh
- LoadMCPConfig
- scorp-agent
- setSession
- 🦂 Scorp
- skills.go
- RegisterTool
- metasearch_engines.go
- compaction_test.go
- renderStatusFooter
- prompt_test.go
- .listenSSEStream
- RenameSession
- GetStringArg
- serviceBridgeRequests
- RunAgentSessionLoop
- TestRegisterPlugin
- clarify.go
- LoadConfig
- session_ui.go
- inline.go
- GetProvider
- estimateHistoryTokens
- TestSteeringQueue
- CleanupSessionsLoop
- ExecuteTermuxAPI
- upload.go
- SummarizeOldToolResult

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 59 edges
2. `GetStringArg()` - 43 edges
3. `StartDaemon()` - 40 edges
4. `RunAgentSessionLoop()` - 39 edges
5. `ModelConfig` - 36 edges
6. `startCLI()` - 35 edges
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

## Communities (58 total, 2 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.06
Nodes (84): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+76 more)

### Community 1 - "HomeDir"
Cohesion: 0.26
Nodes (15): HomeDir(), ProjectDir(), PythonSitePackages(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo() (+7 more)

### Community 2 - "testing.T"
Cohesion: 0.13
Nodes (16): TestSelfReviewCadence(), TestSelfReviewRateLimit(), testing.T, TestBuildCommandCodePayload(), TestExtractFallbackToolCalls(), TestSwitchActiveModel(), IsRateLimitError(), TestCallModelWithToolsNilModel() (+8 more)

### Community 3 - "chat.go"
Cohesion: 0.15
Nodes (23): agentAutoStop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown(), convertTableToList(), ExitAgentMode(), extractAndSaveMemory(), flushPendingMessages() (+15 more)

### Community 4 - "ScorpPath"
Cohesion: 0.06
Nodes (47): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+39 more)

### Community 5 - "client.go"
Cohesion: 0.13
Nodes (24): encoding/json.RawMessage, TruncOutputTool(), buildArgDefsFromInputSchema(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), rebuildMCPToolList() (+16 more)

### Community 6 - "runSubagent"
Cohesion: 0.08
Nodes (37): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+29 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (52): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+44 more)

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
Cohesion: 0.07
Nodes (57): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+49 more)

### Community 13 - "cost_router.go"
Cohesion: 0.05
Nodes (60): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+52 more)

### Community 14 - "time.Time"
Cohesion: 0.13
Nodes (31): time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult() (+23 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "HandleTelegramAction"
Cohesion: 0.21
Nodes (17): HandleTelegramAction(), GetPath(), BackAndRefreshKeyboard(), BackButtonKeyboard(), baseName(), EditMessage(), MainMenuKeyboard(), MonitorMenuKeyboard() (+9 more)

### Community 17 - "StartDaemon"
Cohesion: 0.16
Nodes (17): UploadsDir(), runCommandLoop(), StartDaemon(), DeleteWebhook(), EditMessageByID(), InitTelegram(), SendChatAction(), SendMessageGetID() (+9 more)

### Community 18 - "GetAllTools"
Cohesion: 0.22
Nodes (13): ActivateToolWithTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), TestDynamicToolTTL(), TickToolTTL(), GetAllTools(), ResetNativeToolCache() (+5 more)

### Community 19 - "MCPServer"
Cohesion: 0.18
Nodes (12): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, FindMCPTool(), MCPServer, startMCPServer(), MCPServerConfig, MCPTool (+4 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.05
Nodes (64): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+56 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "LoadMCPConfig"
Cohesion: 0.50
Nodes (8): LoadMCPConfig(), ReloadMCPServers(), sanitizeMCPName(), ExecuteMCPManage(), mcpManageAdd(), mcpManageList(), mcpManageReload(), mcpManageRemove()

### Community 31 - "setSession"
Cohesion: 0.24
Nodes (17): appendSessionHistory(), EnterAgentMode(), GetHistoryTokenEstimate(), getOrCreateSession(), getSessionHistory(), getSessionMap(), AgentMessage, loadHistoryFromDisk() (+9 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "skills.go"
Cohesion: 0.07
Nodes (37): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+29 more)

### Community 34 - "RegisterTool"
Cohesion: 0.11
Nodes (16): RegisterAutonomous(), init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), unregisterMCPNativeTools() (+8 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (40): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+32 more)

### Community 36 - "compaction_test.go"
Cohesion: 0.27
Nodes (16): AgentMessage, makeHistory(), makeToolResult(), TestPrune_BoundaryAges(), TestPrune_DockerScenario_41Messages(), TestPrune_EmptyHistory(), TestPrune_HeadTailFormat(), TestPrune_NonToolMessages_Preserved() (+8 more)

### Community 37 - "renderStatusFooter"
Cohesion: 0.31
Nodes (9): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), getContextPill(), getGitStatus(), getShortCwd() (+1 more)

### Community 38 - "prompt_test.go"
Cohesion: 0.13
Nodes (13): TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestIsDangerousCommand(), TestMaxIterations() (+5 more)

### Community 39 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 40 - "RenameSession"
Cohesion: 0.35
Nodes (11): historyFilePath(), saveHistoryToDisk(), DeleteSession(), ListSessions(), RenameSession(), sanitizeSessionID(), SessionExists(), TestSessionManager() (+3 more)

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (60): StorePendingConfirmation(), init(), init(), ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract() (+52 more)

### Community 42 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.17
Nodes (20): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask() (+12 more)

### Community 44 - "TestRegisterPlugin"
Cohesion: 0.20
Nodes (7): executeMCPServerTool(), echoPlugin, RegisterPlugin(), TestRegisterPlugin(), ExecuteToolByName(), ToolPlugin, ToolPluginWithSchema

### Community 45 - "clarify.go"
Cohesion: 0.25
Nodes (9): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), SetClarifyChatID() (+1 more)

### Community 46 - "LoadConfig"
Cohesion: 0.33
Nodes (9): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), Init(), StartServer() (+1 more)

### Community 47 - "session_ui.go"
Cohesion: 0.42
Nodes (9): BuildSessionMenuKeyboard(), FormatSessionMenuText(), GetActiveSessionID(), HandleSessionCallback(), init(), loadTgSessionMapping(), saveTgSessionMapping(), SetActiveSessionID() (+1 more)

### Community 48 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 49 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 50 - "estimateHistoryTokens"
Cohesion: 0.43
Nodes (6): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, TestEstimateHistoryTokens(), TestEstimateTokens()

### Community 51 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 52 - "CleanupSessionsLoop"
Cohesion: 0.33
Nodes (5): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), cleanupAgentSessions(), TGDocument

### Community 53 - "ExecuteTermuxAPI"
Cohesion: 0.73
Nodes (5): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification()

### Community 54 - "upload.go"
Cohesion: 0.50
Nodes (4): contentPart, imageURL, TestBase64Encode(), base64Encode()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 135 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `HomeDir`, `chat.go`, `RenameSession`, `collector_system.go`, `collector_security.go`, `RunAgentSessionLoop`, `collector_system_native.go`, `clarify.go`, `time.Time`, `session_ui.go`, `HandleModelCallback`, `StartDaemon`, `TestSteeringQueue`, `startCLI`, `LoadMCPConfig`, `setSession`?**
  _High betweenness centrality (0.145) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `context.Context`, `skills.go`, `chat.go`, `prompt_test.go`, `rag_vector.go`, `GetStringArg`, `time.Time`, `HandleTelegramAction`, `GetAllTools`, `TestSteeringQueue`, `startCLI`, `ExecuteTermuxAPI`, `setSession`?**
  _High betweenness centrality (0.092) - this node is a cross-community bridge._
- **Why does `TruncateStr()` connect `context.Context` to `skills.go`, `chat.go`, `client.go`, `runSubagent`, `RenameSession`, `GetStringArg`, `RunAgentSessionLoop`, `cost_router.go`, `time.Time`, `MCPServer`?**
  _High betweenness centrality (0.083) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.05989010989010989 - nodes in this community are weakly interconnected._