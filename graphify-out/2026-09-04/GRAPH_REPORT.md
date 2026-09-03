# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 159 files · ~108,991 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1278 nodes · 3164 edges · 57 communities (46 shown, 3 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 376 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7ea991cc`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- TruncateStr
- HandleTelegramAction
- testing.T
- chat.go
- TruncOutput
- client.go
- compaction_test.go
- rag_vector.go
- collector_system.go
- collector_security.go
- runSubagent
- uptime.go
- wizard.go
- RunAgentSessionLoop
- time.Time
- session_search_fts5.go
- collector_system_native.go
- ConfigMgr
- ScorpPath
- skills.go
- checker.go
- GetProvider
- install.sh
- RunAgentLoop
- scorp-agent
- TestPhase6_AllTools
- 🦂 Scorp
- GetStringArg
- RegisterTool
- TestGatewayEndpoints
- CredentialVault
- files.go
- main
- EscapeHTML
- prompt_test.go
- init
- clarify.go
- runScriptTask
- TestGenerateNativeToolsSchema
- GetBoolArg
- patch.go
- GetRecentReceipts
- ExecuteReadURL
- estimateHistoryTokens
- LoadConfig
- ExecuteTodo
- ExecuteSQL
- ExecuteAutoLogin
- saveQuickstartConfig

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
- `ExecuteSQL()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/db.go → agent/confirmation.go
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

## Communities (57 total, 3 thin omitted)

### Community 0 - "TruncateStr"
Cohesion: 0.05
Nodes (94): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+86 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.14
Nodes (29): StopMCPServerMode(), HandleTelegramAction(), runCommandLoop(), StartDaemon(), AnswerCallback(), BackAndRefreshKeyboard(), BackButtonKeyboard(), DeleteWebhook() (+21 more)

### Community 2 - "testing.T"
Cohesion: 0.15
Nodes (19): TestSelfReviewCadence(), TestSelfReviewRateLimit(), TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty() (+11 more)

### Community 3 - "chat.go"
Cohesion: 0.07
Nodes (58): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+50 more)

### Community 4 - "TruncOutput"
Cohesion: 0.34
Nodes (14): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+6 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (45): bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, buildArgDefsFromInputSchema(), executeMCPServerTool(), FindMCPTool() (+37 more)

### Community 6 - "compaction_test.go"
Cohesion: 0.25
Nodes (17): AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens(), TestPrune_BoundaryAges(), TestPrune_DockerScenario_41Messages(), TestPrune_EmptyHistory(), TestPrune_HeadTailFormat() (+9 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (44): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+36 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (52): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+44 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "runSubagent"
Cohesion: 0.06
Nodes (54): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+46 more)

### Community 11 - "uptime.go"
Cohesion: 0.11
Nodes (27): FormatDuration(), getUptime(), net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost() (+19 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (43): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), SaveModelConfig() (+35 more)

### Community 13 - "RunAgentSessionLoop"
Cohesion: 0.05
Nodes (68): AgentMessage, clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), HandleConfirmation(), HasPendingConfirmation(), countStepsInMessage() (+60 more)

### Community 14 - "time.Time"
Cohesion: 0.21
Nodes (20): time.Time, ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), FormatTasksList(), GetTask() (+12 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "collector_system_native.go"
Cohesion: 0.22
Nodes (17): CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses(), GetTopProcesses() (+9 more)

### Community 17 - "ConfigMgr"
Cohesion: 0.05
Nodes (60): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+52 more)

### Community 18 - "ScorpPath"
Cohesion: 0.11
Nodes (29): init(), hasDebugFlag(), setupCLILogging(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath() (+21 more)

### Community 19 - "skills.go"
Cohesion: 0.07
Nodes (37): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+29 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "RunAgentLoop"
Cohesion: 0.29
Nodes (6): RunAgentLoop(), Init(), StartServer(), StopServer(), getChatLock(), StartTestEndpoint()

### Community 31 - "TestPhase6_AllTools"
Cohesion: 0.26
Nodes (14): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession() (+6 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.11
Nodes (17): 1. Diet Ekstrem `main.go`, 2. Modularisasi `agent/loop.go`, 3. Pemisahan Provider AI Independen (`models/`), 🚀 Panduan Perintah Cepat (Quick Cheatsheet), 🧱 Rekapitulasi Refactoring Arsitektur Bersih (*Clean Architecture*), 📌 Ringkasan Pembaruan Utama (v2.0), 🦂 Scorp Agent v2.0 — Modernization & Architecture Upgrade Report, 🧪 Verifikasi & Kualitas Kode (+9 more)

### Community 33 - "GetStringArg"
Cohesion: 0.23
Nodes (11): init(), GetStringArg(), SendDocumentBytes(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteSystemInfo(), ExecuteWriteFile() (+3 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (38): RegisterAutonomous(), init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), unregisterMCPNativeTools() (+30 more)

### Community 35 - "TestGatewayEndpoints"
Cohesion: 0.28
Nodes (13): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), StartGateway() (+5 more)

### Community 36 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 37 - "files.go"
Cohesion: 0.28
Nodes (14): BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath(), HumanSize(), PathID() (+6 more)

### Community 38 - "main"
Cohesion: 0.29
Nodes (11): isCLIMode(), main(), InitModelUsage(), SOP, Dir(), GetSOP(), InitDefaultSOPs(), ListSOPs() (+3 more)

### Community 39 - "EscapeHTML"
Cohesion: 0.30
Nodes (11): EscapeHTML(), ShellTask(), AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker() (+3 more)

### Community 40 - "prompt_test.go"
Cohesion: 0.17
Nodes (11): TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations(), TestTruncOutput() (+3 more)

### Community 41 - "init"
Cohesion: 0.16
Nodes (10): init(), GetIntArg(), RagVecProvider(), ExecuteCompose(), ExecuteHTTP(), ExecuteLog(), ExecuteAnalyzeImage(), decodeURL() (+2 more)

### Community 42 - "clarify.go"
Cohesion: 0.31
Nodes (8): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID()

### Community 43 - "runScriptTask"
Cohesion: 0.48
Nodes (5): isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage()

### Community 44 - "TestGenerateNativeToolsSchema"
Cohesion: 0.40
Nodes (4): IsRateLimitError(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema(), TestIsRateLimitError()

### Community 45 - "GetBoolArg"
Cohesion: 0.20
Nodes (9): AgentMessage, StorePendingConfirmation(), GetBoolArg(), getUnclosedTags(), SplitMessage(), TruncOutputTool(), ExecuteShell(), ExecuteGit() (+1 more)

### Community 46 - "patch.go"
Cohesion: 0.36
Nodes (9): buildDiffPreview(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines(), TestExecuteReplaceFileContent() (+1 more)

### Community 47 - "GetRecentReceipts"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 48 - "ExecuteReadURL"
Cohesion: 0.52
Nodes (6): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), truncateURLOutput(), tryRemoteScrape()

### Community 49 - "estimateHistoryTokens"
Cohesion: 0.53
Nodes (5): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, TestEstimateTokens()

### Community 50 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 51 - "ExecuteTodo"
Cohesion: 0.60
Nodes (5): ExecuteTodo(), formatTodoList(), formatTodoListLocked(), getStringArgFromMap(), TodoItem

### Community 52 - "ExecuteSQL"
Cohesion: 0.83
Nodes (3): ExecuteSQL(), loadDBConnections(), dbConnection

## Knowledge Gaps
- **26 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+21 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 117 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `TruncateStr`, `chat.go`, `client.go`, `files.go`, `main`, `collector_system.go`, `collector_security.go`, `clarify.go`, `wizard.go`, `RunAgentSessionLoop`, `time.Time`, `collector_system_native.go`, `RunAgentLoop`?**
  _High betweenness centrality (0.145) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `TruncateStr`, `RegisterTool`, `chat.go`, `client.go`, `main`, `rag_vector.go`, `collector_system.go`, `files.go`, `wizard.go`, `RunAgentSessionLoop`, `GetBoolArg`, `time.Time`, `collector_system_native.go`, `ConfigMgr`, `skills.go`, `RunAgentLoop`?**
  _High betweenness centrality (0.113) - this node is a cross-community bridge._
- **Why does `init()` connect `init` to `GetStringArg`, `RegisterTool`, `client.go`, `main`, `rag_vector.go`, `runSubagent`, `RunAgentSessionLoop`, `ExecuteReadURL`, `ScorpPath`, `skills.go`, `ExecuteTodo`?**
  _High betweenness centrality (0.097) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 21 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _26 weakly-connected nodes found - possible documentation gaps or missing edges._