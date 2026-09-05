# Graph Report - scorp  (2026-09-05)

## Corpus Check
- 183 files · ~122,351 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1429 nodes · 3509 edges · 74 communities (62 shown, 4 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 467 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2a0c5acc`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CallCommandCodeWithTools
- TestGatewayEndpoints
- MCPServer
- chat.go
- ScorpPath
- client.go
- GetIntArg
- rag_vector.go
- collector_system.go
- collector_security.go
- bg.go
- collector_system_native.go
- HandleModelCallback
- cost_router.go
- time.Time
- session_search_fts5.go
- TruncOutput
- HandleTelegramAction
- ResolveAPIKey
- telemetry.go
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
- acp.go
- RunAgentSessionLoop
- ResetNativeToolCache
- InitDefaultSOPs
- ExecuteTool
- ToolCall
- main
- runSubagent
- context.Context
- TestSteeringQueue
- HandleUploadInAgentMode
- LoadMCPConfig
- GetBoolArg
- patch.go
- wireCLICallbacks
- TodoManager
- api_gemini.go
- StorePendingConfirmation
- sync.Mutex
- GenerateContextualSessionTitle
- HandleConfirmation
- ExecuteReadURL
- init
- ReloadMCPServers
- .listenSSEStream
- GetAllTools
- echoPlugin
- serviceBridgeRequests
- ExecuteAutoLogin
- ExecuteAnalyzeImage

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
- `ExecuteShell()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/exec.go → agent/confirmation.go
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

## Communities (74 total, 4 thin omitted)

### Community 0 - "CallCommandCodeWithTools"
Cohesion: 0.16
Nodes (17): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), ChatRequest, commandCodeMsg, commandCodeParams (+9 more)

### Community 1 - "TestGatewayEndpoints"
Cohesion: 0.25
Nodes (13): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), TestGatewayEndpoints() (+5 more)

### Community 2 - "MCPServer"
Cohesion: 0.19
Nodes (11): bufio.Scanner, encoding/json.Encoder, FindMCPTool(), MCPServer, startMCPServer(), MCPServerConfig, MCPTool, MCPServer (+3 more)

### Community 3 - "chat.go"
Cohesion: 0.20
Nodes (16): collectTableLines(), convertInlineMarkdown(), convertTableToList(), extractAndSaveMemory(), historyWriterLoop(), init(), isAnyModeActive(), isSeparatorRow() (+8 more)

### Community 4 - "ScorpPath"
Cohesion: 0.06
Nodes (48): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+40 more)

### Community 5 - "client.go"
Cohesion: 0.16
Nodes (20): ACPRequest, encoding/json.RawMessage, executeMCPServerTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError() (+12 more)

### Community 6 - "GetIntArg"
Cohesion: 0.25
Nodes (17): runOpenCodeCLI(), runSubagentACP(), AgentMessage, delegateBatchParams, delegateResult, delegateTaskParams, DefaultSubagentTools(), ExecuteDelegate() (+9 more)

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
Cohesion: 0.38
Nodes (11): bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait(), bgWrite(), closeStdin(), ExecuteBgProcess() (+3 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.08
Nodes (46): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+38 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (44): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "cost_router.go"
Cohesion: 0.05
Nodes (61): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+53 more)

### Community 14 - "time.Time"
Cohesion: 0.16
Nodes (26): printCronTasks(), time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath() (+18 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "TruncOutput"
Cohesion: 0.34
Nodes (14): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+6 more)

### Community 17 - "HandleTelegramAction"
Cohesion: 0.05
Nodes (83): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), HomeDir(), ProjectDir() (+75 more)

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
Cohesion: 0.32
Nodes (15): executeOneShot(), executeTurn(), handleCLISession(), handleCLISOP(), printBanner(), printCLIHelp(), printCostUsage(), printCurrentModel() (+7 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "ModelConfig"
Cohesion: 0.26
Nodes (14): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModel(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), isVisionModelName() (+6 more)

### Community 31 - "getSession"
Cohesion: 0.16
Nodes (25): agentAutoStop(), appendSessionHistory(), EnterAgentMode(), ExitAgentMode(), flushPendingMessages(), GetHistoryTokenEstimate(), getOrCreateSession(), getSession() (+17 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "v2_skills.go"
Cohesion: 0.08
Nodes (33): getSharedMemorySummary(), memoryFact, getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage, maybeRunSelfReview() (+25 more)

### Community 34 - "RegisterTool"
Cohesion: 0.13
Nodes (16): init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), unregisterMCPNativeTools(), TestRegisterPlugin() (+8 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (39): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+31 more)

### Community 36 - "testing.T"
Cohesion: 0.06
Nodes (56): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+48 more)

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
Cohesion: 0.23
Nodes (11): init(), GetStringArg(), SendDocumentBytes(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteShell(), ExecuteSystemInfo() (+3 more)

### Community 42 - "acp.go"
Cohesion: 0.21
Nodes (10): checkACPAvailable(), launchACP(), listAvailableACP(), ACPError, ACPInitializeParams, ACPMessageNewParams, ACPMessagePart, ACPResponse (+2 more)

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.15
Nodes (24): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask() (+16 more)

### Community 44 - "ResetNativeToolCache"
Cohesion: 0.20
Nodes (12): TestGenerateNativeToolsSchema(), ActivateToolWithTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), TestDynamicToolTTL(), TickToolTTL(), RegisterPlugin() (+4 more)

### Community 45 - "InitDefaultSOPs"
Cohesion: 0.42
Nodes (8): SOP, Dir(), GetSOP(), InitDefaultSOPs(), ListSOPs(), SaveSOP(), TestSOPLifecycle(), ExecuteSOP()

### Community 46 - "ExecuteTool"
Cohesion: 0.20
Nodes (9): ExecuteTool(), FormatToolResult(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), TestAutonomyLevels(), AutonomyLevel, RedactSecrets() (+1 more)

### Community 47 - "ToolCall"
Cohesion: 0.24
Nodes (11): TestParseToolCalls(), CustomProvider, ToolCall, CallModelWithTools(), CallModelWithToolsAndFallback(), IsRateLimitError(), ParseAllToolCalls(), ParseCodeBlockFallback() (+3 more)

### Community 48 - "main"
Cohesion: 0.21
Nodes (9): RegisterAutonomous(), hasDebugFlag(), SetAutonomyLevel(), StartGateway(), isCLIMode(), main(), InitModelUsage(), RunQuickstart() (+1 more)

### Community 49 - "runSubagent"
Cohesion: 0.24
Nodes (10): cleanupSubagentSandbox(), createSubagentSandbox(), defaultIsolation(), formatIsolationInfo(), getSubagentIsolation(), isSubagentToolBlocked(), registerSubagentIsolation(), unregisterSubagentIsolation() (+2 more)

### Community 50 - "context.Context"
Cohesion: 0.21
Nodes (21): context.Context, TruncateStr(), callAnthropic(), CallAnthropicWithTools(), CallCommandCodeStream(), callGemini(), CallGeminiWithTools(), geminiDoRequest() (+13 more)

### Community 51 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 52 - "HandleUploadInAgentMode"
Cohesion: 0.25
Nodes (7): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), RunAgentLoop(), cleanupAgentSessions(), TGDocument, HandleUploadInAgentMode()

### Community 53 - "LoadMCPConfig"
Cohesion: 0.22
Nodes (11): TruncOutputTool(), buildArgDefsFromInputSchema(), LoadMCPConfig(), rebuildMCPToolList(), registerMCPToolsAsNative(), StartMCPServers(), StopMCPServers(), GetServerHealthStatus() (+3 more)

### Community 54 - "GetBoolArg"
Cohesion: 0.18
Nodes (7): TestGetBoolArg(), GetBoolArg(), getUnclosedTags(), SplitMessage(), ExecuteCompose(), ExecuteHTTP(), ExecuteLog()

### Community 57 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 58 - "wireCLICallbacks"
Cohesion: 0.27
Nodes (8): wireCLICallbacks(), formatFinalResponse(), formatTerminalText(), isTerminal(), stripHTML(), TestFormatFinalResponse(), TestStripHTML(), os.File

### Community 59 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 60 - "api_gemini.go"
Cohesion: 0.36
Nodes (10): geminiBuildRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall, geminiGenConfig, geminiPart, geminiRequest (+2 more)

### Community 61 - "StorePendingConfirmation"
Cohesion: 0.24
Nodes (7): AgentMessage, StorePendingConfirmation(), ExecuteSQL(), loadDBConnections(), dbConnection, ExecuteGit(), ExecuteProcess()

### Community 62 - "sync.Mutex"
Cohesion: 0.22
Nodes (9): bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, ServerWatchdog, CostTracker, getChatLock(), StartTestEndpoint() (+1 more)

### Community 63 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 64 - "HandleConfirmation"
Cohesion: 0.39
Nodes (8): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation, printStatus()

### Community 65 - "ExecuteReadURL"
Cohesion: 0.36
Nodes (7): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), TestReadURL_LocalMock(), truncateURLOutput(), tryRemoteScrape()

### Community 66 - "init"
Cohesion: 0.29
Nodes (5): init(), PythonSitePackages(), TestLiveWebSearch(), ExecuteWebFetch(), ExecuteWebSearch()

### Community 67 - "ReloadMCPServers"
Cohesion: 0.54
Nodes (7): ReloadMCPServers(), sanitizeMCPName(), ExecuteMCPManage(), mcpManageAdd(), mcpManageList(), mcpManageReload(), mcpManageRemove()

### Community 68 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 69 - "GetAllTools"
Cohesion: 0.67
Nodes (5): GetAllTools(), countActiveTools(), countDeferredTools(), ExecuteToolList(), ExecuteToolSearch()

### Community 71 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 138 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `HandleConfirmation`, `v2_skills.go`, `RenameSession`, `collector_system.go`, `collector_security.go`, `RunAgentSessionLoop`, `collector_system_native.go`, `InitDefaultSOPs`, `time.Time`, `HandleModelCallback`, `TestSteeringQueue`, `LoadMCPConfig`, `startCLI`, `getSession`?**
  _High betweenness centrality (0.148) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `chat.go`, `rag_vector.go`, `time.Time`, `HandleTelegramAction`, `startCLI`, `getSession`, `v2_skills.go`, `testing.T`, `RenameSession`, `GetStringArg`, `ResetNativeToolCache`, `ExecuteTool`, `ToolCall`, `context.Context`, `TestSteeringQueue`, `HandleUploadInAgentMode`, `StorePendingConfirmation`, `GenerateContextualSessionTitle`, `HandleConfirmation`?**
  _High betweenness centrality (0.087) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `ScorpPath`, `client.go`, `rag_vector.go`, `collector_system.go`, `collector_system_native.go`, `HandleModelCallback`, `cost_router.go`, `time.Time`, `ResolveAPIKey`, `v2_skills.go`, `testing.T`, `skills.go`, `InitDefaultSOPs`, `ExecuteTool`, `main`, `HandleUploadInAgentMode`, `LoadMCPConfig`, `StorePendingConfirmation`, `sync.Mutex`?**
  _High betweenness centrality (0.083) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 5 inferred relationships involving `startCLI()` (e.g. with `wireCLICallbacks()` and `formatTerminalText()`) actually correct?**
  _`startCLI()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 25 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 25 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._