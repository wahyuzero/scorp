# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 170 files · ~117,078 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1371 nodes · 3361 edges · 67 communities (54 shown, 5 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 426 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ff6c2f13`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- context.Context
- HandleTelegramAction
- testing.T
- chat.go
- ScorpPath
- client.go
- metasearch.go
- rag_vector.go
- collector_system.go
- collector_security.go
- runSubagent
- collector_system_native.go
- HandleModelCallback
- cost_router.go
- time.Time
- session_search_fts5.go
- compaction_test.go
- StartDaemon
- MCPServer
- runSelfReview
- checker.go
- startCLI
- install.sh
- renderStatusFooter
- scorp-agent
- prompt_test.go
- 🦂 Scorp
- skills.go
- RegisterTool
- net/http.Client
- RunAgentSessionLoop
- TestRegisterPlugin
- inline.go
- ResetNativeToolCache
- bg.go
- GetStringArg
- clarify.go
- resumeAgentLoop
- IsDangerousCommand
- SearchResult
- uptime.go
- time.Duration
- ExecuteTermuxAPI
- TestSteeringQueue
- metasearch_test.go
- EscapeHTML
- ExecuteMCPManage
- RunAgentLoop
- StartServer
- GetProvider
- GetAllTools
- estimateHistoryTokens
- .listenSSEStream
- serviceBridgeRequests
- config_helper.go
- DuckDuckGoLiteEngine
- BraveSearchEngine
- DuckDuckGoHTMLEngine
- TavilySearchEngine

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 54 edges
2. `GetStringArg()` - 43 edges
3. `StartDaemon()` - 40 edges
4. `RunAgentSessionLoop()` - 38 edges
5. `ModelConfig` - 36 edges
6. `startCLI()` - 35 edges
7. `TruncateStr()` - 35 edges
8. `init()` - 32 edges
9. `ChatMessage` - 30 edges
10. `HandleModelCallback()` - 29 edges

## Surprising Connections (you probably didn't know these)
- `SetupInlineMode()` --calls--> `TgPost()`  [INFERRED]
  tools/inline.go → telegram/telegram.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `ExecuteShell()` --calls--> `IsDangerousCommand()`  [INFERRED]
  tools/exec.go → agent/prompt.go
- `serviceBridgeRequests()` --calls--> `ExecuteTool()`  [INFERRED]
  tools/exec_code.go → agent/prompt.go
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go

## Import Cycles
- None detected.

## Communities (67 total, 5 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.06
Nodes (89): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+81 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.18
Nodes (23): HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath(), HumanSize() (+15 more)

### Community 2 - "testing.T"
Cohesion: 0.13
Nodes (22): TestSelfReviewCadence(), TestSelfReviewRateLimit(), TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty() (+14 more)

### Community 3 - "chat.go"
Cohesion: 0.07
Nodes (59): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+51 more)

### Community 4 - "ScorpPath"
Cohesion: 0.05
Nodes (69): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+61 more)

### Community 5 - "client.go"
Cohesion: 0.12
Nodes (26): encoding/json.RawMessage, getExposedTools(), GetMCPTools(), handleMCPRequest(), LoadMCPConfig(), MCPToolsForPrompt(), rebuildMCPToolList(), sendMCPError() (+18 more)

### Community 6 - "metasearch.go"
Cohesion: 0.23
Nodes (13): GitHubEngine, GetMetaSearchAggregator(), NewBingEngine(), NewBraveSearchEngine(), NewDefaultMetaSearchAggregator(), NewDuckDuckGoHTMLEngine(), NewGitHubEngine(), NewGoogleCSEEngine() (+5 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (52): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+44 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "runSubagent"
Cohesion: 0.08
Nodes (39): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+31 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.22
Nodes (17): CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses(), GetTopProcesses() (+9 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.07
Nodes (57): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+49 more)

### Community 13 - "cost_router.go"
Cohesion: 0.05
Nodes (60): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+52 more)

### Community 14 - "time.Time"
Cohesion: 0.25
Nodes (18): time.Time, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), FormatTasksList(), GetTask(), LoadTasks() (+10 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "compaction_test.go"
Cohesion: 0.25
Nodes (17): AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens(), TestPrune_BoundaryAges(), TestPrune_DockerScenario_41Messages(), TestPrune_EmptyHistory(), TestPrune_HeadTailFormat() (+9 more)

### Community 17 - "StartDaemon"
Cohesion: 0.19
Nodes (21): runCommandLoop(), StartDaemon(), AnswerCallback(), DeleteWebhook(), EditMessage(), EditMessageByID(), InitTelegram(), PollUpdates() (+13 more)

### Community 18 - "MCPServer"
Cohesion: 0.24
Nodes (8): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, FindMCPTool(), MCPServer, jsonRPCError, jsonRPCResponse, MCPTool

### Community 19 - "runSelfReview"
Cohesion: 0.23
Nodes (13): memoryFact, AgentMessage, maybeRunSelfReview(), runSelfReview(), TestSelfReviewIntegration(), LoadJSON(), deleteMemory(), ExecuteMemory() (+5 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.06
Nodes (64): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+56 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "renderStatusFooter"
Cohesion: 0.27
Nodes (10): GetHistoryTokenEstimate(), GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), getContextPill(), getGitStatus() (+2 more)

### Community 31 - "prompt_test.go"
Cohesion: 0.17
Nodes (11): TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations(), TestTruncOutput() (+3 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "skills.go"
Cohesion: 0.24
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 34 - "RegisterTool"
Cohesion: 0.14
Nodes (14): init(), init(), init(), init(), TruncOutputTool(), buildArgDefsFromInputSchema(), registerMCPToolsAsNative(), GenerateSystemPromptDescriptions() (+6 more)

### Community 35 - "net/http.Client"
Cohesion: 0.24
Nodes (9): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+1 more)

### Community 36 - "RunAgentSessionLoop"
Cohesion: 0.38
Nodes (9): countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask(), screenshotWasTaken() (+1 more)

### Community 37 - "TestRegisterPlugin"
Cohesion: 0.20
Nodes (7): executeMCPServerTool(), echoPlugin, RegisterPlugin(), TestRegisterPlugin(), ExecuteToolByName(), ToolPlugin, ToolPluginWithSchema

### Community 38 - "inline.go"
Cohesion: 0.33
Nodes (10): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+2 more)

### Community 39 - "ResetNativeToolCache"
Cohesion: 0.27
Nodes (9): TestGenerateNativeToolsSchema(), ActivateToolWithTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), TestDynamicToolTTL(), TickToolTTL(), ResetNativeToolCache() (+1 more)

### Community 40 - "bg.go"
Cohesion: 0.12
Nodes (20): bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, CostTracker, bgKill(), bgList(), bgPoll() (+12 more)

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (51): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage() (+43 more)

### Community 42 - "clarify.go"
Cohesion: 0.29
Nodes (8): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), SetClarifyChatID(), PendingClarify

### Community 43 - "resumeAgentLoop"
Cohesion: 0.26
Nodes (10): AgentMessage, cleanToolCallTags(), maxIterations(), resumeAgentLoop(), TestBuildThinkingMessage(), TestToolDescription(), buildThinkingMessage(), shouldUpdateThinking() (+2 more)

### Community 44 - "IsDangerousCommand"
Cohesion: 0.15
Nodes (10): getSharedMemorySummary(), FormatToolResult(), getAgentSystemPrompt(), IsDangerousCommand(), TestIsDangerousCommand(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap() (+2 more)

### Community 45 - "SearchResult"
Cohesion: 0.20
Nodes (5): GoogleCSEEngine, deduplicateAndRank(), SearchResult, SearXNGEngine, WikipediaEngine

### Community 46 - "uptime.go"
Cohesion: 0.38
Nodes (10): AddUptimeTarget(), checkTarget(), ExecuteUptime(), ListUptimeTargets(), RemoveUptimeTarget(), runUptimeCheck(), UptimeLoop(), UptimeMonitor (+2 more)

### Community 47 - "time.Duration"
Cohesion: 0.27
Nodes (9): FormatDuration(), getUptime(), time.Duration, AgentMessage, AutonomousConfig, AutonomousLogEntry, ChatSession, pendingConfirmation (+1 more)

### Community 48 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 49 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 50 - "metasearch_test.go"
Cohesion: 0.24
Nodes (7): normalizeSearchURL(), TestDeduplicateAndRank(), TestLiveWebSearch(), TestMetaSearchAggregator_ConcurrentSuccess(), TestMetaSearchAggregator_FaultTolerance(), TestNormalizeSearchURL(), MockSearchEngine

### Community 51 - "EscapeHTML"
Cohesion: 0.36
Nodes (7): EscapeHTML(), isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage(), ShellTask()

### Community 52 - "ExecuteMCPManage"
Cohesion: 0.54
Nodes (7): ReloadMCPServers(), sanitizeMCPName(), ExecuteMCPManage(), mcpManageAdd(), mcpManageList(), mcpManageReload(), mcpManageRemove()

### Community 53 - "RunAgentLoop"
Cohesion: 0.67
Nodes (3): RunAgentLoop(), getChatLock(), StartTestEndpoint()

### Community 54 - "StartServer"
Cohesion: 0.67
Nodes (3): Init(), StartServer(), StopServer()

### Community 57 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 58 - "GetAllTools"
Cohesion: 0.48
Nodes (6): unregisterMCPNativeTools(), GetAllTools(), UnregisterTool(), countActiveTools(), countDeferredTools(), ExecuteToolList()

### Community 59 - "estimateHistoryTokens"
Cohesion: 0.53
Nodes (5): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, TestEstimateTokens()

### Community 60 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 61 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 62 - "config_helper.go"
Cohesion: 0.67
Nodes (3): SaveJSON(), SaveJSONPerm(), os.FileMode

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 133 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `context.Context`, `chat.go`, `collector_system.go`, `collector_security.go`, `clarify.go`, `collector_system_native.go`, `HandleModelCallback`, `time.Time`, `StartDaemon`, `TestSteeringQueue`, `ExecuteMCPManage`, `RunAgentLoop`, `startCLI`?**
  _High betweenness centrality (0.125) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `context.Context`, `skills.go`, `HandleTelegramAction`, `chat.go`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `collector_system_native.go`, `IsDangerousCommand`, `cost_router.go`, `HandleModelCallback`, `time.Time`, `runSelfReview`, `startCLI`, `StartServer`, `RunAgentLoop`?**
  _High betweenness centrality (0.085) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `skills.go`, `RegisterTool`, `ScorpPath`, `rag_vector.go`, `bg.go`, `ResetNativeToolCache`, `runSubagent`, `ExecuteTermuxAPI`, `runSelfReview`, `ExecuteMCPManage`, `startCLI`, `GetAllTools`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.05554628857381151 - nodes in this community are weakly interconnected._