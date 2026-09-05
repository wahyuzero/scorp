# Graph Report - scorp  (2026-09-05)

## Corpus Check
- 185 files · ~122,755 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1431 nodes · 3507 edges · 75 communities (62 shown, 5 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 463 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `01edf518`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- time.Time
- bg.go
- TestGatewayEndpoints
- chat.go
- ScorpPath
- client.go
- runSubagent
- rag_vector.go
- collector_system.go
- collector_security.go
- Incident Report: Runaway CLI Verify Session
- LoadModelConfig
- HandleModelCallback
- cost_router.go
- collector_system_native.go
- session_search_fts5.go
- context.Context
- StartDaemon
- CallModelWithFallback
- api_gemini.go
- checker.go
- startCLI
- install.sh
- ToolCall
- scorp-agent
- truncateToolResultsInHistory
- 🦂 Scorp
- v2_skills.go
- RegisterTool
- metasearch_engines.go
- testing.T
- renderStatusFooter
- skills.go
- browser_session.go
- TruncateStr
- GetStringArg
- ModelConfig
- RunAgentSessionLoop
- inline.go
- collector_native_test.go
- SearchResult
- wireCLICallbacks
- CredentialVault
- ExecuteTool
- CallOpenAIWithTools
- uptime.go
- HandleConfirmation
- main
- time.Duration
- InitDefaultSOPs
- TestPhase6_AllTools
- TodoManager
- net/http.Client
- EscapeHTML
- clarify.go
- GetRecentReceipts
- ExecuteScript
- estimateHistoryTokens
- HandleTelegramAction
- LoadConfig
- NewDefaultMetaSearchAggregator
- TestGenerateNativeToolsSchema
- StartServer
- BraveSearchEngine
- DuckDuckGoHTMLEngine
- WikipediaEngine
- ExecuteAutoLogin

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 64 edges
2. `startCLI()` - 46 edges
3. `GetStringArg()` - 43 edges
4. `StartDaemon()` - 43 edges
5. `RunAgentSessionLoop()` - 39 edges
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

## Communities (75 total, 5 thin omitted)

### Community 0 - "time.Time"
Cohesion: 0.23
Nodes (19): printCronTasks(), time.Time, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), FormatTasksList(), GetTask() (+11 more)

### Community 1 - "bg.go"
Cohesion: 0.25
Nodes (15): bytes.Buffer, io.WriteCloser, os/exec.Cmd, bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait() (+7 more)

### Community 2 - "TestGatewayEndpoints"
Cohesion: 0.25
Nodes (13): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), StartGateway() (+5 more)

### Community 3 - "chat.go"
Cohesion: 0.06
Nodes (72): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+64 more)

### Community 4 - "ScorpPath"
Cohesion: 0.11
Nodes (29): init(), setupCLILogging(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), HomeDir() (+21 more)

### Community 5 - "client.go"
Cohesion: 0.06
Nodes (55): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, sync.Mutex, TruncOutputTool() (+47 more)

### Community 6 - "runSubagent"
Cohesion: 0.08
Nodes (37): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+29 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.08
Nodes (51): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+43 more)

### Community 9 - "collector_security.go"
Cohesion: 0.18
Nodes (25): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+17 more)

### Community 10 - "Incident Report: Runaway CLI Verify Session"
Cohesion: 0.25
Nodes (7): Akar Masalah (dugaan), Ciri-ciri Runaway yang Terdeteksi, Incident Report: Runaway CLI Verify Session, Pelajaran Umum, Penanganan, Rekomendasi Perbaikan, Ringkasan

### Community 11 - "LoadModelConfig"
Cohesion: 0.24
Nodes (13): defaultModelConfig(), LoadModelConfig(), ModelRouterConfig, ModelUsage, ProviderPreset, migrateModelConfigs(), getProviderInfo(), HandleProviderCommand() (+5 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (44): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "cost_router.go"
Cohesion: 0.05
Nodes (60): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+52 more)

### Community 14 - "collector_system_native.go"
Cohesion: 0.23
Nodes (16): TestCollectorNative_StartCPUSampler_DoesNotBlock(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+8 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "context.Context"
Cohesion: 0.20
Nodes (14): context.Context, AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic(), CallAnthropicWithTools(), CallCommandCode() (+6 more)

### Community 17 - "StartDaemon"
Cohesion: 0.15
Nodes (23): runCommandLoop(), StartDaemon(), BackAndRefreshKeyboard(), BackButtonKeyboard(), DeleteWebhook(), EditMessage(), EditMessageByID(), InitTelegram() (+15 more)

### Community 18 - "CallModelWithFallback"
Cohesion: 0.29
Nodes (12): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), isVisionModelName(), RouteModel() (+4 more)

### Community 19 - "api_gemini.go"
Cohesion: 0.35
Nodes (11): geminiBuildRequest(), geminiDoRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall, geminiGenConfig, geminiPart (+3 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.32
Nodes (15): executeOneShot(), executeTurn(), handleCLISession(), handleCLISOP(), printBanner(), printCLIHelp(), printCostUsage(), printCurrentModel() (+7 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "ToolCall"
Cohesion: 0.23
Nodes (11): TestParseToolCalls(), CustomProvider, LLMProvider, GetProvider(), init(), RegisterProviderAdapter(), ToolCall, CallModelWithTools() (+3 more)

### Community 31 - "truncateToolResultsInHistory"
Cohesion: 0.22
Nodes (18): AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens(), TestPrune_BoundaryAges(), TestPrune_DockerScenario_41Messages(), TestPrune_EmptyHistory(), TestPrune_HeadTailFormat() (+10 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.11
Nodes (17): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+9 more)

### Community 33 - "v2_skills.go"
Cohesion: 0.07
Nodes (35): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+27 more)

### Community 34 - "RegisterTool"
Cohesion: 0.07
Nodes (26): init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), unregisterMCPNativeTools(), echoPlugin (+18 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.11
Nodes (13): BingEngine, DuckDuckGoLiteEngine, ghSearchResp, GitHubEngine, GoogleCSEEngine, NewBingEngine(), NewDuckDuckGoLiteEngine(), NewGitHubEngine() (+5 more)

### Community 36 - "testing.T"
Cohesion: 0.10
Nodes (24): TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations(), TestTruncOutput() (+16 more)

### Community 37 - "renderStatusFooter"
Cohesion: 0.23
Nodes (11): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste(), getContextPill() (+3 more)

### Community 38 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 39 - "browser_session.go"
Cohesion: 0.33
Nodes (13): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+5 more)

### Community 40 - "TruncateStr"
Cohesion: 0.17
Nodes (17): TruncateStr(), buildCommandCodePayload(), CallCommandCodeStream(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), ChatRequest, commandCodeMsg (+9 more)

### Community 41 - "GetStringArg"
Cohesion: 0.05
Nodes (58): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage() (+50 more)

### Community 42 - "ModelConfig"
Cohesion: 0.28
Nodes (14): CallModel(), CallModelStream(), ModelConfig, applyProviderDefaults(), hasAPIKey(), KeySourceLabel(), ResolveAPIFormat(), ResolveAPIKey() (+6 more)

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.07
Nodes (36): AgentMessage, IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery(), cleanToolCallTags(), getSessionSearchContext(), maxIterations() (+28 more)

### Community 44 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 45 - "collector_native_test.go"
Cohesion: 0.21
Nodes (12): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+4 more)

### Community 46 - "SearchResult"
Cohesion: 0.23
Nodes (9): deduplicateAndRank(), normalizeSearchURL(), TestDeduplicateAndRank(), TestLiveWebSearch(), TestMetaSearchAggregator_ConcurrentSuccess(), TestMetaSearchAggregator_FaultTolerance(), TestNormalizeSearchURL(), MockSearchEngine (+1 more)

### Community 47 - "wireCLICallbacks"
Cohesion: 0.27
Nodes (8): wireCLICallbacks(), formatFinalResponse(), formatTerminalText(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML(), os.File

### Community 48 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 49 - "ExecuteTool"
Cohesion: 0.23
Nodes (9): ExecuteTool(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), SetAutonomyLevel(), TestAutonomyLevels(), AutonomyLevel, RedactSecrets() (+1 more)

### Community 50 - "CallOpenAIWithTools"
Cohesion: 0.31
Nodes (7): CallOpenAI(), CallOpenAIWithTools(), formatOpenAIMessages(), RecordCostWithCache(), OpenAIProvider, TrackModelUsageWithCache(), GetAIClient()

### Community 51 - "uptime.go"
Cohesion: 0.38
Nodes (10): AddUptimeTarget(), checkTarget(), ExecuteUptime(), ListUptimeTargets(), RemoveUptimeTarget(), runUptimeCheck(), UptimeLoop(), UptimeMonitor (+2 more)

### Community 52 - "HandleConfirmation"
Cohesion: 0.33
Nodes (9): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+1 more)

### Community 53 - "main"
Cohesion: 0.24
Nodes (7): RegisterAutonomous(), hasDebugFlag(), isCLIMode(), main(), InitModelUsage(), RunQuickstart(), saveQuickstartConfig()

### Community 54 - "time.Duration"
Cohesion: 0.27
Nodes (9): FormatDuration(), getUptime(), time.Duration, AgentMessage, AutonomousConfig, AutonomousLogEntry, ChatSession, pendingConfirmation (+1 more)

### Community 57 - "InitDefaultSOPs"
Cohesion: 0.42
Nodes (8): SOP, Dir(), GetSOP(), InitDefaultSOPs(), ListSOPs(), SaveSOP(), TestSOPLifecycle(), ExecuteSOP()

### Community 58 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 59 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 60 - "net/http.Client"
Cohesion: 0.36
Nodes (8): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport()

### Community 61 - "EscapeHTML"
Cohesion: 0.36
Nodes (7): EscapeHTML(), isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage(), ShellTask()

### Community 62 - "clarify.go"
Cohesion: 0.24
Nodes (10): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+2 more)

### Community 63 - "GetRecentReceipts"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 64 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 65 - "estimateHistoryTokens"
Cohesion: 0.53
Nodes (5): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, TestEstimateTokens()

### Community 66 - "HandleTelegramAction"
Cohesion: 0.21
Nodes (21): HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath(), HumanSize() (+13 more)

### Community 67 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 68 - "NewDefaultMetaSearchAggregator"
Cohesion: 0.73
Nodes (4): GetMetaSearchAggregator(), NewDefaultMetaSearchAggregator(), MetaSearchAggregator, SearchEngine

### Community 69 - "TestGenerateNativeToolsSchema"
Cohesion: 0.40
Nodes (4): IsRateLimitError(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema(), TestIsRateLimitError()

### Community 70 - "StartServer"
Cohesion: 0.67
Nodes (3): Init(), StartServer(), StopServer()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 141 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `time.Time`, `v2_skills.go`, `chat.go`, `client.go`, `collector_system.go`, `collector_security.go`, `RunAgentSessionLoop`, `HandleModelCallback`, `collector_system_native.go`, `StartDaemon`, `HandleConfirmation`, `startCLI`, `InitDefaultSOPs`, `clarify.go`?**
  _High betweenness centrality (0.153) - this node is a cross-community bridge._
- **Why does `GetStringArg()` connect `GetStringArg` to `time.Time`, `v2_skills.go`, `ExecuteScript`, `testing.T`, `client.go`, `runSubagent`, `browser_session.go`, `ExecuteAutoLogin`, `RunAgentSessionLoop`, `LoadModelConfig`, `cost_router.go`, `CredentialVault`, `uptime.go`, `InitDefaultSOPs`?**
  _High betweenness centrality (0.096) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `time.Time`, `chat.go`, `ScorpPath`, `client.go`, `rag_vector.go`, `collector_system.go`, `LoadModelConfig`, `HandleModelCallback`, `cost_router.go`, `collector_system_native.go`, `v2_skills.go`, `skills.go`, `GetStringArg`, `RunAgentSessionLoop`, `ExecuteTool`, `main`, `InitDefaultSOPs`, `HandleTelegramAction`, `StartServer`?**
  _High betweenness centrality (0.087) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 5 inferred relationships involving `startCLI()` (e.g. with `wireCLICallbacks()` and `formatTerminalText()`) actually correct?**
  _`startCLI()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._