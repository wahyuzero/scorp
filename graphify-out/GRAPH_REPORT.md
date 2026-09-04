# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 162 files · ~114,527 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1346 nodes · 3308 edges · 66 communities (52 shown, 6 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 380 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e44f9abb`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- context.Context
- StartDaemon
- testing.T
- chat.go
- ScorpPath
- client.go
- metasearch.go
- rag_vector.go
- collector_system.go
- collector_security.go
- execute.go
- collector_system_native.go
- wizard.go
- agent/autonomous.go
- time.Time
- session_search_fts5.go
- helpers.go
- cost_router.go
- HandleTelegramAction
- skills.go
- checker.go
- startCLI
- install.sh
- acp.go
- scorp-agent
- init
- 🦂 Scorp
- compaction_test.go
- RegisterTool
- browser_session.go
- RunAgentSessionLoop
- runSubagent
- net/http.Client
- SearchResult
- collector_native_test.go
- GetStringArg
- patch.go
- GetProvider
- HomeDir
- inline.go
- uptime.go
- time.Duration
- clarify.go
- metasearch_test.go
- TestPhase6_AllTools
- EscapeHTML
- GetRecentReceipts
- ExecuteScript
- LoadConfig
- read_url.go
- ExecuteTodo
- upload.go
- StartServer
- DuckDuckGoLiteEngine
- BraveSearchEngine
- DuckDuckGoHTMLEngine
- ExecuteAutoLogin
- TavilySearchEngine

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 54 edges
2. `GetStringArg()` - 43 edges
3. `StartDaemon()` - 40 edges
4. `RunAgentSessionLoop()` - 37 edges
5. `TruncateStr()` - 36 edges
6. `ModelConfig` - 36 edges
7. `startCLI()` - 34 edges
8. `init()` - 32 edges
9. `ChatMessage` - 29 edges
10. `HandleModelCallback()` - 29 edges

## Surprising Connections (you probably didn't know these)
- `SetupInlineMode()` --calls--> `TgPost()`  [INFERRED]
  tools/inline.go → telegram/telegram.go
- `ExecuteGit()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/git.go → agent/confirmation.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `ExecuteShell()` --calls--> `IsDangerousCommand()`  [INFERRED]
  tools/exec.go → agent/prompt.go
- `serviceBridgeRequests()` --calls--> `ExecuteTool()`  [INFERRED]
  tools/exec_code.go → agent/prompt.go

## Import Cycles
- None detected.

## Communities (66 total, 6 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.06
Nodes (82): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+74 more)

### Community 1 - "StartDaemon"
Cohesion: 0.19
Nodes (20): runCommandLoop(), StartDaemon(), AnswerCallback(), DeleteWebhook(), EditMessage(), EditMessageByID(), InitTelegram(), PollUpdates() (+12 more)

### Community 2 - "testing.T"
Cohesion: 0.11
Nodes (22): TestBase64Encode(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations() (+14 more)

### Community 3 - "chat.go"
Cohesion: 0.08
Nodes (54): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+46 more)

### Community 4 - "ScorpPath"
Cohesion: 0.13
Nodes (24): init(), setupCLILogging(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), Hostname() (+16 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (48): ACPRequest, bufio.Scanner, context.CancelFunc, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool() (+40 more)

### Community 6 - "metasearch.go"
Cohesion: 0.23
Nodes (13): GitHubEngine, GetMetaSearchAggregator(), NewBingEngine(), NewBraveSearchEngine(), NewDefaultMetaSearchAggregator(), NewDuckDuckGoHTMLEngine(), NewGitHubEngine(), NewGoogleCSEEngine() (+5 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (53): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+45 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "execute.go"
Cohesion: 0.26
Nodes (16): runOpenCodeCLI(), runSubagentACP(), AgentMessage, delegateBatchParams, delegateResult, delegateTaskParams, DefaultSubagentTools(), ExecuteDelegate() (+8 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.23
Nodes (16): TestCollectorNative_StartCPUSampler_DoesNotBlock(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+8 more)

### Community 12 - "wizard.go"
Cohesion: 0.09
Nodes (52): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), LoadModelConfig() (+44 more)

### Community 13 - "agent/autonomous.go"
Cohesion: 0.12
Nodes (31): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+23 more)

### Community 14 - "time.Time"
Cohesion: 0.21
Nodes (20): time.Time, ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), FormatTasksList(), GetTask() (+12 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "helpers.go"
Cohesion: 0.67
Nodes (3): GetStringSliceArg(), getUnclosedTags(), SplitMessage()

### Community 17 - "cost_router.go"
Cohesion: 0.09
Nodes (30): CM(), ConfigMgr(), InitConfigManager(), NewConfigManager(), ConfigManager, os.FileMode, defaultCostConfig(), formatCostReport() (+22 more)

### Community 18 - "HandleTelegramAction"
Cohesion: 0.19
Nodes (21): HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath(), HumanSize() (+13 more)

### Community 19 - "skills.go"
Cohesion: 0.07
Nodes (36): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+28 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.06
Nodes (65): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+57 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "acp.go"
Cohesion: 0.08
Nodes (30): checkACPAvailable(), launchACP(), listAvailableACP(), ACPError, ACPInitializeParams, ACPMessageNewParams, ACPMessagePart, ACPResponse (+22 more)

### Community 31 - "init"
Cohesion: 0.17
Nodes (11): init(), GetBoolArg(), GetIntArg(), ExecuteCompose(), ExecuteToolSearch(), ExecuteGit(), ExecuteHTTP(), ExecuteLog() (+3 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "compaction_test.go"
Cohesion: 0.19
Nodes (22): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+14 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (39): init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), executeMCPServerTool(), unregisterMCPNativeTools() (+31 more)

### Community 35 - "browser_session.go"
Cohesion: 0.33
Nodes (13): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+5 more)

### Community 36 - "RunAgentSessionLoop"
Cohesion: 0.09
Nodes (34): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), looksLikeContinuation(), mentionsBrowserTask(), screenshotWasTaken() (+26 more)

### Community 37 - "runSubagent"
Cohesion: 0.24
Nodes (10): cleanupSubagentSandbox(), createSubagentSandbox(), defaultIsolation(), formatIsolationInfo(), getSubagentIsolation(), isSubagentToolBlocked(), registerSubagentIsolation(), unregisterSubagentIsolation() (+2 more)

### Community 38 - "net/http.Client"
Cohesion: 0.24
Nodes (9): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+1 more)

### Community 39 - "SearchResult"
Cohesion: 0.20
Nodes (5): GoogleCSEEngine, deduplicateAndRank(), SearchResult, SearXNGEngine, WikipediaEngine

### Community 40 - "collector_native_test.go"
Cohesion: 0.23
Nodes (11): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+3 more)

### Community 41 - "GetStringArg"
Cohesion: 0.15
Nodes (18): StorePendingConfirmation(), init(), GetStringArg(), TruncOutput(), SendDocumentBytes(), ExecuteSQL(), loadDBConnections(), dbConnection (+10 more)

### Community 42 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 43 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 44 - "HomeDir"
Cohesion: 0.22
Nodes (10): HomeDir(), ProjectDir(), PythonSitePackages(), UploadsDir(), applyProviderDefaults(), ProviderPreset, migrateModelConfigs(), ResolveBaseURL() (+2 more)

### Community 45 - "inline.go"
Cohesion: 0.33
Nodes (10): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+2 more)

### Community 46 - "uptime.go"
Cohesion: 0.38
Nodes (10): AddUptimeTarget(), checkTarget(), ExecuteUptime(), ListUptimeTargets(), RemoveUptimeTarget(), runUptimeCheck(), UptimeLoop(), UptimeMonitor (+2 more)

### Community 47 - "time.Duration"
Cohesion: 0.27
Nodes (9): FormatDuration(), getUptime(), time.Duration, AgentMessage, AutonomousConfig, AutonomousLogEntry, ChatSession, pendingConfirmation (+1 more)

### Community 48 - "clarify.go"
Cohesion: 0.31
Nodes (8): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID()

### Community 49 - "metasearch_test.go"
Cohesion: 0.24
Nodes (7): normalizeSearchURL(), TestDeduplicateAndRank(), TestLiveWebSearch(), TestMetaSearchAggregator_ConcurrentSuccess(), TestMetaSearchAggregator_FaultTolerance(), TestNormalizeSearchURL(), MockSearchEngine

### Community 50 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 51 - "EscapeHTML"
Cohesion: 0.36
Nodes (7): EscapeHTML(), isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage(), ShellTask()

### Community 52 - "GetRecentReceipts"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 53 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 54 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 57 - "read_url.go"
Cohesion: 0.60
Nodes (5): ReadURL(), scrapeFirecrawl(), scrapeTavily(), truncateURLOutput(), tryRemoteScrape()

### Community 58 - "ExecuteTodo"
Cohesion: 0.60
Nodes (5): ExecuteTodo(), formatTodoList(), formatTodoListLocked(), getStringArgFromMap(), TodoItem

### Community 59 - "upload.go"
Cohesion: 0.67
Nodes (3): contentPart, imageURL, base64Encode()

### Community 60 - "StartServer"
Cohesion: 0.67
Nodes (3): Init(), StartServer(), StopServer()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 133 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `StartDaemon`, `chat.go`, `RunAgentSessionLoop`, `client.go`, `collector_system.go`, `collector_security.go`, `collector_system_native.go`, `wizard.go`, `time.Time`, `clarify.go`, `startCLI`?**
  _High betweenness centrality (0.134) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `chat.go`, `RunAgentSessionLoop`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `collector_system_native.go`, `wizard.go`, `agent/autonomous.go`, `time.Time`, `cost_router.go`, `HandleTelegramAction`, `skills.go`, `startCLI`, `StartServer`?**
  _High betweenness centrality (0.087) - this node is a cross-community bridge._
- **Why does `init()` connect `init` to `RegisterTool`, `RunAgentSessionLoop`, `client.go`, `rag_vector.go`, `GetStringArg`, `execute.go`, `patch.go`, `HomeDir`, `skills.go`, `startCLI`, `acp.go`, `ExecuteTodo`?**
  _High betweenness centrality (0.084) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.061488673139158574 - nodes in this community are weakly interconnected._