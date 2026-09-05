# Graph Report - scorp  (2026-09-05)

## Corpus Check
- 182 files · ~121,601 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1420 nodes · 3500 edges · 75 communities (62 shown, 5 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 467 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0145d8c9`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CallCommandCodeWithTools
- truncateToolResultsInHistory
- init
- chat.go
- ScorpPath
- client.go
- runSubagent
- rag_vector.go
- collector_system.go
- collector_security.go
- browser_session.go
- collector_system_native.go
- HandleModelCallback
- cost_router.go
- time.Time
- session_search_fts5.go
- GetAllTools
- StartDaemon
- ResolveAPIKey
- main
- checker.go
- startCLI
- install.sh
- ModelConfig
- scorp-agent
- getSession
- 🦂 Scorp
- v2_skills.go
- RegisterTool
- metasearch_engines.go
- testing.T
- renderStatusFooter
- skills.go
- .CallWithTools
- RenameSession
- GetStringArg
- collector_native_test.go
- RunAgentSessionLoop
- inline.go
- patch.go
- CredentialVault
- ToolCall
- saveQuickstartConfig
- collector_docker.go
- context.Context
- TestSteeringQueue
- CleanupSessionsLoop
- TestPhase6_AllTools
- FormatHourlyReport
- collector_coolify.go
- ExecuteScript
- TodoManager
- api_gemini.go
- HandleTelegramAction
- clarify.go
- GenerateContextualSessionTitle
- read_url.go
- ExecuteTermuxAPI
- HomeDir
- upload.go
- telemetry.go
- ExecuteSQL
- echoPlugin
- ExecuteAutoLogin
- ExecuteAnalyzeImage
- GetRecentReceipts
- GetProvider

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 64 edges
2. `startCLI()` - 46 edges
3. `RunAgentSessionLoop()` - 43 edges
4. `GetStringArg()` - 43 edges
5. `StartDaemon()` - 43 edges
6. `ModelConfig` - 36 edges
7. `TruncateStr()` - 35 edges
8. `init()` - 32 edges
9. `ChatMessage` - 30 edges
10. `HandleModelCallback()` - 29 edges

## Surprising Connections (you probably didn't know these)
- `SetupInlineMode()` --calls--> `TgPost()`  [INFERRED]
  tools/inline.go → telegram/telegram.go
- `ExecuteSQL()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/db.go → agent/confirmation.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go
- `ExecuteShell()` --calls--> `IsDangerousCommand()`  [INFERRED]
  tools/exec.go → agent/safety.go

## Import Cycles
- None detected.

## Communities (75 total, 5 thin omitted)

### Community 0 - "CallCommandCodeWithTools"
Cohesion: 0.16
Nodes (17): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), ChatRequest, commandCodeMsg, commandCodeParams (+9 more)

### Community 1 - "truncateToolResultsInHistory"
Cohesion: 0.17
Nodes (23): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+15 more)

### Community 2 - "init"
Cohesion: 0.19
Nodes (9): init(), GetBoolArg(), GetIntArg(), ExecuteCompose(), ExecuteHTTP(), ExecuteLog(), ExecuteReadURL(), ExecuteWebFetch() (+1 more)

### Community 3 - "chat.go"
Cohesion: 0.13
Nodes (21): collectTableLines(), convertInlineMarkdown(), convertTableToList(), extractAndSaveMemory(), flushPendingMessages(), GetHistoryTokenEstimate(), historyWriterLoop(), init() (+13 more)

### Community 4 - "ScorpPath"
Cohesion: 0.12
Nodes (25): init(), setupCLILogging(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), Hostname() (+17 more)

### Community 5 - "client.go"
Cohesion: 0.06
Nodes (53): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool(), buildArgDefsFromInputSchema() (+45 more)

### Community 6 - "runSubagent"
Cohesion: 0.06
Nodes (56): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+48 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.23
Nodes (19): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+11 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "browser_session.go"
Cohesion: 0.33
Nodes (13): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+5 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.23
Nodes (16): TestCollectorNative_StartCPUSampler_DoesNotBlock(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+8 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (44): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "cost_router.go"
Cohesion: 0.07
Nodes (56): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+48 more)

### Community 14 - "time.Time"
Cohesion: 0.08
Nodes (45): printCronTasks(), FormatDuration(), getUptime(), time.Duration, time.Time, EscapeHTML(), ScheduledTask, AddTask() (+37 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "GetAllTools"
Cohesion: 0.23
Nodes (12): ActivateToolWithTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), TestDynamicToolTTL(), TickToolTTL(), GetAllTools(), countActiveTools() (+4 more)

### Community 17 - "StartDaemon"
Cohesion: 0.10
Nodes (35): getUnclosedTags(), SplitMessage(), StopMCPServerMode(), Init(), StartServer(), StopServer(), InitModelUsage(), runCommandLoop() (+27 more)

### Community 18 - "ResolveAPIKey"
Cohesion: 0.17
Nodes (19): defaultModelConfig(), LoadModelConfig(), ModelRouterConfig, ModelUsage, applyProviderDefaults(), ProviderPreset, hasAPIKey(), migrateModelConfigs() (+11 more)

### Community 19 - "main"
Cohesion: 0.08
Nodes (34): hasDebugFlag(), Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), CM() (+26 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.07
Nodes (46): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+38 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "ModelConfig"
Cohesion: 0.26
Nodes (14): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModel(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), isVisionModelName() (+6 more)

### Community 31 - "getSession"
Cohesion: 0.28
Nodes (18): agentAutoStop(), appendSessionHistory(), EnterAgentMode(), ExitAgentMode(), getOrCreateSession(), getSession(), getSessionHistory(), getSessionMap() (+10 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.11
Nodes (17): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+9 more)

### Community 33 - "v2_skills.go"
Cohesion: 0.07
Nodes (34): getSharedMemorySummary(), memoryFact, getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage, maybeRunSelfReview() (+26 more)

### Community 34 - "RegisterTool"
Cohesion: 0.10
Nodes (21): RegisterAutonomous(), init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), TestGenerateNativeToolsSchema() (+13 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (40): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+32 more)

### Community 36 - "testing.T"
Cohesion: 0.09
Nodes (26): TestBase64Encode(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestIsDangerousCommand() (+18 more)

### Community 37 - "renderStatusFooter"
Cohesion: 0.23
Nodes (11): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste(), getContextPill() (+3 more)

### Community 38 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 39 - ".CallWithTools"
Cohesion: 0.33
Nodes (4): AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool

### Community 40 - "RenameSession"
Cohesion: 0.33
Nodes (11): ClearChatSession(), historyFilePath(), saveHistoryToDisk(), DeleteSession(), ListSessions(), RenameSession(), sanitizeSessionID(), SessionExists() (+3 more)

### Community 41 - "GetStringArg"
Cohesion: 0.19
Nodes (14): StorePendingConfirmation(), init(), GetStringArg(), TruncOutput(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteShell() (+6 more)

### Community 42 - "collector_native_test.go"
Cohesion: 0.23
Nodes (11): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+3 more)

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.17
Nodes (21): AgentMessage, setLoopActive(), countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation() (+13 more)

### Community 44 - "inline.go"
Cohesion: 0.33
Nodes (10): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+2 more)

### Community 45 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 46 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 47 - "ToolCall"
Cohesion: 0.24
Nodes (11): TestParseToolCalls(), CustomProvider, ToolCall, CallModelWithTools(), CallModelWithToolsAndFallback(), IsRateLimitError(), ParseAllToolCalls(), ParseCodeBlockFallback() (+3 more)

### Community 49 - "collector_docker.go"
Cohesion: 0.24
Nodes (12): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), TestCollectorNative_StartDockerStatsSampler_DoesNotBlock() (+4 more)

### Community 50 - "context.Context"
Cohesion: 0.21
Nodes (21): context.Context, TruncateStr(), callAnthropic(), CallAnthropicWithTools(), CallCommandCodeStream(), callGemini(), CallGeminiWithTools(), geminiDoRequest() (+13 more)

### Community 51 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 52 - "CleanupSessionsLoop"
Cohesion: 0.33
Nodes (5): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), cleanupAgentSessions(), TGDocument

### Community 53 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 54 - "FormatHourlyReport"
Cohesion: 0.33
Nodes (11): SystemData, Bar(), bar(), FormatHourlyReport(), FormatStatusResponse(), SectionCoolify(), SectionDocker(), SectionNetwork() (+3 more)

### Community 57 - "collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 58 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 59 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 60 - "api_gemini.go"
Cohesion: 0.36
Nodes (10): geminiBuildRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall, geminiGenConfig, geminiPart, geminiRequest (+2 more)

### Community 61 - "HandleTelegramAction"
Cohesion: 0.30
Nodes (13): FormatCompactStats(), CompactStats, FormatUsageStats(), HandleTelegramAction(), BuildSessionMenuKeyboard(), FormatSessionMenuText(), GetActiveSessionID(), HandleSessionCallback() (+5 more)

### Community 62 - "clarify.go"
Cohesion: 0.24
Nodes (10): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+2 more)

### Community 63 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 64 - "read_url.go"
Cohesion: 0.60
Nodes (5): ReadURL(), scrapeFirecrawl(), scrapeTavily(), truncateURLOutput(), tryRemoteScrape()

### Community 65 - "ExecuteTermuxAPI"
Cohesion: 0.73
Nodes (5): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification()

### Community 66 - "HomeDir"
Cohesion: 0.22
Nodes (17): HomeDir(), ProjectDir(), PythonSitePackages(), UploadsDir(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard() (+9 more)

### Community 67 - "upload.go"
Cohesion: 0.67
Nodes (3): contentPart, imageURL, base64Encode()

### Community 68 - "telemetry.go"
Cohesion: 0.67
Nodes (3): CheckModelHealth(), FormatModelList(), FormatModelListWithHealth()

### Community 69 - "ExecuteSQL"
Cohesion: 0.83
Nodes (3): ExecuteSQL(), loadDBConnections(), dbConnection

### Community 74 - "GetRecentReceipts"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 76 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

## Knowledge Gaps
- **26 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+21 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 132 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `chat.go`, `client.go`, `collector_system.go`, `collector_security.go`, `collector_system_native.go`, `HandleModelCallback`, `time.Time`, `StartDaemon`, `main`, `startCLI`, `getSession`, `v2_skills.go`, `RenameSession`, `RunAgentSessionLoop`, `collector_docker.go`, `TestSteeringQueue`, `FormatHourlyReport`, `collector_coolify.go`, `clarify.go`, `HomeDir`?**
  _High betweenness centrality (0.151) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `v2_skills.go`, `ExecuteTermuxAPI`, `chat.go`, `testing.T`, `rag_vector.go`, `RenameSession`, `GetStringArg`, `time.Time`, `ToolCall`, `GetAllTools`, `context.Context`, `TestSteeringQueue`, `HandleTelegramAction`, `startCLI`, `GenerateContextualSessionTitle`, `getSession`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `ScorpPath`, `client.go`, `runSubagent`, `rag_vector.go`, `collector_system_native.go`, `HandleModelCallback`, `cost_router.go`, `time.Time`, `ResolveAPIKey`, `main`, `startCLI`, `v2_skills.go`, `RegisterTool`, `testing.T`, `skills.go`, `GetStringArg`, `collector_docker.go`, `CleanupSessionsLoop`, `HomeDir`?**
  _High betweenness centrality (0.085) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 5 inferred relationships involving `startCLI()` (e.g. with `wireCLICallbacks()` and `formatTerminalText()`) actually correct?**
  _`startCLI()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 25 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 25 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _26 weakly-connected nodes found - possible documentation gaps or missing edges._