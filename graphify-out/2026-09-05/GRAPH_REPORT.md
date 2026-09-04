# Graph Report - scorp  (2026-09-05)

## Corpus Check
- 178 files · ~119,586 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1402 nodes · 3435 edges · 71 communities (60 shown, 3 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 451 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `edd51c74`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CallCommandCodeWithTools
- MCPServer
- testing.T
- chat.go
- ScorpPath
- client.go
- execute.go
- init
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
- LoadModelConfig
- ModelConfig
- checker.go
- startCLI
- install.sh
- CallModelWithFallback
- scorp-agent
- getSession
- 🦂 Scorp
- skills.go
- RegisterTool
- metasearch_engines.go
- compaction_test.go
- renderStatusFooter
- prompt_test.go
- .CallWithTools
- RenameSession
- init
- acp.go
- RunAgentSessionLoop
- GetIntArg
- runSubagent
- StorePendingConfirmation
- ToolCall
- patch.go
- GetProvider
- ChatMessage
- TestSteeringQueue
- HandleUploadInAgentMode
- ExecuteTermuxAPI
- CredentialVault
- tools/monitor.go
- uptime.go
- TestPhase6_AllTools
- context.Context
- GetStringArg
- ExecuteReadURL
- GenerateContextualSessionTitle
- helpers.go
- GetRecentReceipts
- ExecuteScript
- .listenSSEStream
- ExecuteSQL
- ExecuteAutoLogin
- ExecuteAnalyzeImage

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 59 edges
2. `GetStringArg()` - 43 edges
3. `RunAgentSessionLoop()` - 42 edges
4. `StartDaemon()` - 41 edges
5. `startCLI()` - 37 edges
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
- `ExecuteGit()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/git.go → agent/confirmation.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go

## Import Cycles
- None detected.

## Communities (71 total, 3 thin omitted)

### Community 0 - "CallCommandCodeWithTools"
Cohesion: 0.17
Nodes (16): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), ChatRequest, commandCodeMsg, commandCodeParams (+8 more)

### Community 1 - "MCPServer"
Cohesion: 0.17
Nodes (13): bufio.Scanner, encoding/json.Encoder, FindMCPTool(), MCPServer, startMCPServer(), jsonRPCError, jsonRPCResponse, MCPServerConfig (+5 more)

### Community 2 - "testing.T"
Cohesion: 0.13
Nodes (22): TestSelfReviewCadence(), TestSelfReviewRateLimit(), TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty() (+14 more)

### Community 3 - "chat.go"
Cohesion: 0.18
Nodes (17): cleanupChatSessions(), CleanupSessionsLoop(), collectTableLines(), convertInlineMarkdown(), convertTableToList(), extractAndSaveMemory(), historyWriterLoop(), init() (+9 more)

### Community 4 - "ScorpPath"
Cohesion: 0.13
Nodes (23): init(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), Hostname(), MCPConfigFilePath() (+15 more)

### Community 5 - "client.go"
Cohesion: 0.16
Nodes (21): ACPRequest, encoding/json.RawMessage, getExposedTools(), GetMCPTools(), handleMCPRequest(), LoadMCPConfig(), MCPToolsForPrompt(), rebuildMCPToolList() (+13 more)

### Community 6 - "execute.go"
Cohesion: 0.26
Nodes (16): runOpenCodeCLI(), runSubagentACP(), AgentMessage, delegateBatchParams, delegateResult, delegateTaskParams, DefaultSubagentTools(), ExecuteDelegate() (+8 more)

### Community 7 - "init"
Cohesion: 0.06
Nodes (46): init(), activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir() (+38 more)

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
Cohesion: 0.19
Nodes (20): FormatDuration(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+12 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (43): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+35 more)

### Community 13 - "cost_router.go"
Cohesion: 0.06
Nodes (53): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+45 more)

### Community 14 - "time.Time"
Cohesion: 0.08
Nodes (42): printCronTasks(), time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath() (+34 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "TruncOutput"
Cohesion: 0.34
Nodes (14): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+6 more)

### Community 17 - "HandleTelegramAction"
Cohesion: 0.05
Nodes (76): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), HomeDir(), ProjectDir() (+68 more)

### Community 18 - "LoadModelConfig"
Cohesion: 0.24
Nodes (13): defaultModelConfig(), LoadModelConfig(), ModelRouterConfig, ModelUsage, ProviderPreset, migrateModelConfigs(), getProviderInfo(), HandleProviderCommand() (+5 more)

### Community 19 - "ModelConfig"
Cohesion: 0.19
Nodes (24): TruncateStr(), callAnthropic(), CallAnthropicWithTools(), CallCommandCodeStream(), RecordCost(), CallModelStream(), ModelConfig, applyProviderDefaults() (+16 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.05
Nodes (64): ExecuteTool(), FormatToolResult(), RegisterAutonomous(), wireCLICallbacks(), executeOneShot(), executeTurn(), formatFinalResponse(), formatTerminalText() (+56 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "CallModelWithFallback"
Cohesion: 0.30
Nodes (12): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModel(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), isVisionModelName() (+4 more)

### Community 31 - "getSession"
Cohesion: 0.18
Nodes (22): agentAutoStop(), appendSessionHistory(), EnterAgentMode(), ExitAgentMode(), flushPendingMessages(), GetHistoryTokenEstimate(), getOrCreateSession(), getSession() (+14 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "skills.go"
Cohesion: 0.08
Nodes (35): getSharedMemorySummary(), memoryFact, getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage, maybeRunSelfReview() (+27 more)

### Community 34 - "RegisterTool"
Cohesion: 0.07
Nodes (36): init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), executeMCPServerTool(), unregisterMCPNativeTools() (+28 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (39): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+31 more)

### Community 36 - "compaction_test.go"
Cohesion: 0.17
Nodes (23): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+15 more)

### Community 37 - "renderStatusFooter"
Cohesion: 0.23
Nodes (11): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste(), getContextPill() (+3 more)

### Community 38 - "prompt_test.go"
Cohesion: 0.14
Nodes (12): contentPart, imageURL, TestBase64Encode(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestIsDangerousCommand(), TestMaxIterations() (+4 more)

### Community 39 - ".CallWithTools"
Cohesion: 0.29
Nodes (4): AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool

### Community 40 - "RenameSession"
Cohesion: 0.33
Nodes (11): ClearChatSession(), historyFilePath(), saveHistoryToDisk(), DeleteSession(), ListSessions(), RenameSession(), sanitizeSessionID(), SessionExists() (+3 more)

### Community 41 - "init"
Cohesion: 0.22
Nodes (9): init(), SendDocumentBytes(), ExecuteListDir(), ExecuteReadFile(), ExecuteSendFile(), ExecuteSystemInfo(), ExecuteWriteFile(), isPathAllowed() (+1 more)

### Community 42 - "acp.go"
Cohesion: 0.21
Nodes (10): checkACPAvailable(), launchACP(), listAvailableACP(), ACPError, ACPInitializeParams, ACPMessageNewParams, ACPMessagePart, ACPResponse (+2 more)

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.18
Nodes (19): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask() (+11 more)

### Community 44 - "GetIntArg"
Cohesion: 0.18
Nodes (9): TestGetBoolArg(), GetBoolArg(), GetIntArg(), ExecuteGit(), ExecuteHTTP(), ExecuteLog(), TestLiveWebSearch(), ExecuteWebFetch() (+1 more)

### Community 45 - "runSubagent"
Cohesion: 0.24
Nodes (10): cleanupSubagentSandbox(), createSubagentSandbox(), defaultIsolation(), formatIsolationInfo(), getSubagentIsolation(), isSubagentToolBlocked(), registerSubagentIsolation(), unregisterSubagentIsolation() (+2 more)

### Community 46 - "StorePendingConfirmation"
Cohesion: 0.24
Nodes (11): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), StorePendingConfirmation() (+3 more)

### Community 47 - "ToolCall"
Cohesion: 0.23
Nodes (10): TestParseToolCalls(), CustomProvider, ToolCall, CallModelWithTools(), IsRateLimitError(), ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls() (+2 more)

### Community 48 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 49 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 50 - "ChatMessage"
Cohesion: 0.36
Nodes (7): CallOpenAI(), CallOpenAIWithTools(), formatOpenAIMessages(), ChatMessage, RecordCostWithCache(), OpenAIProvider, TrackModelUsageWithCache()

### Community 51 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 52 - "HandleUploadInAgentMode"
Cohesion: 0.29
Nodes (6): agentSession, scorpChat(), touchSession(), RunAgentLoop(), TGDocument, HandleUploadInAgentMode()

### Community 53 - "ExecuteTermuxAPI"
Cohesion: 0.73
Nodes (5): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification()

### Community 54 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 57 - "tools/monitor.go"
Cohesion: 0.38
Nodes (10): ExecuteMonitor(), InitMonitor(), loadMonitorTargets(), monitorCheckOne(), monitorLoop(), ragIngestText(), sanitizeFilename(), saveMonitorTargets() (+2 more)

### Community 58 - "uptime.go"
Cohesion: 0.38
Nodes (10): AddUptimeTarget(), checkTarget(), ExecuteUptime(), ListUptimeTargets(), RemoveUptimeTarget(), runUptimeCheck(), UptimeLoop(), UptimeMonitor (+2 more)

### Community 59 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 60 - "context.Context"
Cohesion: 0.23
Nodes (15): context.Context, callGemini(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), geminiContent, geminiFuncDecl (+7 more)

### Community 61 - "GetStringArg"
Cohesion: 0.50
Nodes (8): GetStringArg(), ReloadMCPServers(), sanitizeMCPName(), ExecuteMCPManage(), mcpManageAdd(), mcpManageList(), mcpManageReload(), mcpManageRemove()

### Community 62 - "ExecuteReadURL"
Cohesion: 0.36
Nodes (7): ExecuteReadURL(), ReadURL(), scrapeFirecrawl(), scrapeTavily(), TestReadURL_LocalMock(), truncateURLOutput(), tryRemoteScrape()

### Community 63 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 64 - "helpers.go"
Cohesion: 0.29
Nodes (7): TestGetStringSliceArg(), GetStringSliceArg(), getUnclosedTags(), SplitMessage(), TruncOutputTool(), buildArgDefsFromInputSchema(), registerMCPToolsAsNative()

### Community 65 - "GetRecentReceipts"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 66 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 67 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 68 - "ExecuteSQL"
Cohesion: 0.83
Nodes (3): ExecuteSQL(), loadDBConnections(), dbConnection

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 135 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `RenameSession`, `collector_system.go`, `collector_security.go`, `RunAgentSessionLoop`, `collector_system_native.go`, `HandleModelCallback`, `StorePendingConfirmation`, `time.Time`, `TestSteeringQueue`, `ModelConfig`, `startCLI`, `GetStringArg`, `getSession`?**
  _High betweenness centrality (0.132) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `skills.go`, `RegisterTool`, `chat.go`, `prompt_test.go`, `init`, `RenameSession`, `StorePendingConfirmation`, `time.Time`, `HandleTelegramAction`, `TestSteeringQueue`, `HandleUploadInAgentMode`, `startCLI`, `GetStringArg`, `ModelConfig`, `CallModelWithFallback`, `ExecuteTermuxAPI`, `GenerateContextualSessionTitle`, `getSession`?**
  _High betweenness centrality (0.112) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `skills.go`, `chat.go`, `client.go`, `prompt_test.go`, `init`, `collector_system.go`, `bg.go`, `collector_system_native.go`, `HandleModelCallback`, `cost_router.go`, `StorePendingConfirmation`, `time.Time`, `LoadModelConfig`, `ModelConfig`, `startCLI`?**
  _High betweenness centrality (0.093) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 25 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 25 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.12535612535612536 - nodes in this community are weakly interconnected._