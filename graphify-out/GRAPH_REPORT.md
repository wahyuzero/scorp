# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 171 files · ~117,928 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1372 nodes · 3361 edges · 54 communities (41 shown, 5 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 426 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9f874f5f`
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
- ConfigManager
- collector_native_test.go
- checker.go
- startCLI
- install.sh
- scorp-agent
- 🦂 Scorp
- skills.go
- RegisterTool
- net/http.Client
- EscapeHTML
- bg.go
- GetStringArg
- clarify.go
- RunAgentSessionLoop
- SearchResult
- uptime.go
- time.Duration
- metasearch_test.go
- runScriptTask
- RunAgentLoop
- StartServer
- GetProvider
- estimateHistoryTokens
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
- `serviceBridgeRequests()` --calls--> `ExecuteTool()`  [INFERRED]
  tools/exec_code.go → agent/prompt.go
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go
- `ExecuteShell()` --calls--> `IsDangerousCommand()`  [INFERRED]
  tools/exec.go → agent/safety.go

## Import Cycles
- None detected.

## Communities (54 total, 5 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.06
Nodes (86): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+78 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.17
Nodes (24): FormatUsageStats(), HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath() (+16 more)

### Community 2 - "testing.T"
Cohesion: 0.10
Nodes (24): TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestIsDangerousCommand(), TestMaxIterations() (+16 more)

### Community 3 - "chat.go"
Cohesion: 0.07
Nodes (59): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+51 more)

### Community 4 - "ScorpPath"
Cohesion: 0.05
Nodes (64): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+56 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (48): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, os/exec.Cmd, TruncOutputTool() (+40 more)

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
Nodes (40): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+32 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.23
Nodes (16): TestCollectorNative_StartCPUSampler_DoesNotBlock(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+8 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.08
Nodes (56): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+48 more)

### Community 13 - "cost_router.go"
Cohesion: 0.07
Nodes (56): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+48 more)

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
Cohesion: 0.17
Nodes (23): StopMCPServerMode(), Load(), runCommandLoop(), StartDaemon(), AnswerCallback(), DeleteWebhook(), EditMessage(), EditMessageByID() (+15 more)

### Community 18 - "ConfigManager"
Cohesion: 0.12
Nodes (16): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts() (+8 more)

### Community 19 - "collector_native_test.go"
Cohesion: 0.21
Nodes (12): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+4 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.06
Nodes (55): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+47 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "skills.go"
Cohesion: 0.07
Nodes (36): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+28 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (37): RegisterAutonomous(), init(), init(), init(), init(), unregisterMCPNativeTools(), TestGenerateNativeToolsSchema(), ActivateToolWithTTL() (+29 more)

### Community 35 - "net/http.Client"
Cohesion: 0.24
Nodes (9): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+1 more)

### Community 38 - "EscapeHTML"
Cohesion: 0.27
Nodes (12): EscapeHTML(), ShellTask(), AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker() (+4 more)

### Community 40 - "bg.go"
Cohesion: 0.13
Nodes (18): bytes.Buffer, sync.Mutex, CostTracker, bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait() (+10 more)

### Community 41 - "GetStringArg"
Cohesion: 0.05
Nodes (65): StorePendingConfirmation(), init(), init(), ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract() (+57 more)

### Community 42 - "clarify.go"
Cohesion: 0.29
Nodes (8): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), SetClarifyChatID(), PendingClarify

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.11
Nodes (30): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask() (+22 more)

### Community 45 - "SearchResult"
Cohesion: 0.20
Nodes (5): GoogleCSEEngine, deduplicateAndRank(), SearchResult, SearXNGEngine, WikipediaEngine

### Community 46 - "uptime.go"
Cohesion: 0.38
Nodes (10): AddUptimeTarget(), checkTarget(), ExecuteUptime(), ListUptimeTargets(), RemoveUptimeTarget(), runUptimeCheck(), UptimeLoop(), UptimeMonitor (+2 more)

### Community 47 - "time.Duration"
Cohesion: 0.27
Nodes (9): FormatDuration(), getUptime(), time.Duration, AgentMessage, AutonomousConfig, AutonomousLogEntry, ChatSession, pendingConfirmation (+1 more)

### Community 50 - "metasearch_test.go"
Cohesion: 0.24
Nodes (7): normalizeSearchURL(), TestDeduplicateAndRank(), TestLiveWebSearch(), TestMetaSearchAggregator_ConcurrentSuccess(), TestMetaSearchAggregator_FaultTolerance(), TestNormalizeSearchURL(), MockSearchEngine

### Community 51 - "runScriptTask"
Cohesion: 0.48
Nodes (5): isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage()

### Community 53 - "RunAgentLoop"
Cohesion: 0.67
Nodes (3): RunAgentLoop(), getChatLock(), StartTestEndpoint()

### Community 54 - "StartServer"
Cohesion: 0.67
Nodes (3): Init(), StartServer(), StopServer()

### Community 57 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 59 - "estimateHistoryTokens"
Cohesion: 0.52
Nodes (6): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, maybeCompactHistory(), TestEstimateTokens()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 133 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `chat.go`, `ScorpPath`, `client.go`, `collector_system.go`, `collector_security.go`, `clarify.go`, `RunAgentSessionLoop`, `collector_system_native.go`, `HandleModelCallback`, `time.Time`, `StartDaemon`, `RunAgentLoop`, `startCLI`?**
  _High betweenness centrality (0.131) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `skills.go`, `RegisterTool`, `ScorpPath`, `client.go`, `rag_vector.go`, `bg.go`, `runSubagent`, `RunAgentSessionLoop`?**
  _High betweenness centrality (0.084) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `HandleTelegramAction`, `testing.T`, `chat.go`, `ScorpPath`, `RegisterTool`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `skills.go`, `collector_system_native.go`, `HandleModelCallback`, `cost_router.go`, `time.Time`, `ConfigManager`, `RunAgentLoop`, `StartServer`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.05819079527420208 - nodes in this community are weakly interconnected._