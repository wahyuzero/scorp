# Graph Report - scorp  (2026-09-05)

## Corpus Check
- 177 files · ~119,241 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1397 nodes · 3423 edges · 63 communities (53 shown, 2 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 448 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a6a0583e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- ToolCall
- HandleTelegramAction
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
- StartDaemon
- ResolveAPIKey
- TruncateStr
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
- context.Context
- RenameSession
- GetStringArg
- RunAgentSessionLoop
- clarify.go
- LoadConfig
- session_ui.go
- inline.go
- GetProvider
- estimateHistoryTokens
- TestSteeringQueue
- HandleUploadInAgentMode
- ExecuteTermuxAPI
- upload.go
- collector_docker.go
- FormatHourlyReport
- collector_coolify.go
- api_gemini.go
- GenerateContextualSessionTitle
- session_fallback_test.go
- api_anthropic.go
- TestSwitchActiveModel
- RedactSecrets

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 59 edges
2. `GetStringArg()` - 43 edges
3. `RunAgentSessionLoop()` - 42 edges
4. `StartDaemon()` - 41 edges
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

## Communities (63 total, 2 thin omitted)

### Community 0 - "ToolCall"
Cohesion: 0.15
Nodes (20): TestParseToolCalls(), buildCommandCodePayload(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), ChatRequest, commandCodeMsg, commandCodeParams (+12 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.15
Nodes (24): HomeDir(), ProjectDir(), PythonSitePackages(), UploadsDir(), FormatUsageStats(), HandleTelegramAction(), runCommandLoop(), BackKB() (+16 more)

### Community 2 - "testing.T"
Cohesion: 0.13
Nodes (19): TestSelfReviewCadence(), TestSelfReviewRateLimit(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues() (+11 more)

### Community 3 - "chat.go"
Cohesion: 0.18
Nodes (17): cleanupChatSessions(), CleanupSessionsLoop(), collectTableLines(), convertInlineMarkdown(), convertTableToList(), extractAndSaveMemory(), historyWriterLoop(), init() (+9 more)

### Community 4 - "ScorpPath"
Cohesion: 0.05
Nodes (61): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+53 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (47): bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool(), buildArgDefsFromInputSchema(), executeMCPServerTool() (+39 more)

### Community 6 - "runSubagent"
Cohesion: 0.08
Nodes (37): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+29 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.21
Nodes (20): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+12 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "bg.go"
Cohesion: 0.11
Nodes (26): bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, CostTracker, getChatLock(), StartTestEndpoint(), bgKill() (+18 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.10
Nodes (35): TestCollectorNative_CollectSystem_Structure(), FormatDuration(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes() (+27 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (44): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "cost_router.go"
Cohesion: 0.05
Nodes (60): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+52 more)

### Community 14 - "time.Time"
Cohesion: 0.17
Nodes (25): time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult() (+17 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 17 - "StartDaemon"
Cohesion: 0.09
Nodes (35): getUnclosedTags(), SplitMessage(), StopMCPServerMode(), Init(), StartServer(), StopServer(), InitModelUsage(), Load() (+27 more)

### Community 18 - "ResolveAPIKey"
Cohesion: 0.15
Nodes (20): defaultModelConfig(), LoadModelConfig(), CustomProvider, ModelRouterConfig, ModelUsage, applyProviderDefaults(), ProviderPreset, hasAPIKey() (+12 more)

### Community 19 - "TruncateStr"
Cohesion: 0.21
Nodes (19): TruncateStr(), callAnthropic(), CallAnthropicWithTools(), CallCommandCodeStream(), CallOpenAI(), CallOpenAIWithTools(), formatOpenAIMessages(), RecordCost() (+11 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.06
Nodes (64): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+56 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "CallModelWithFallback"
Cohesion: 0.20
Nodes (16): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModel(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), isVisionModelName() (+8 more)

### Community 31 - "getSession"
Cohesion: 0.18
Nodes (22): agentAutoStop(), appendSessionHistory(), EnterAgentMode(), ExitAgentMode(), flushPendingMessages(), GetHistoryTokenEstimate(), getOrCreateSession(), getSession() (+14 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "skills.go"
Cohesion: 0.07
Nodes (36): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+28 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (34): init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema() (+26 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (40): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+32 more)

### Community 36 - "compaction_test.go"
Cohesion: 0.22
Nodes (18): AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens(), TestPrune_BoundaryAges(), TestPrune_DockerScenario_41Messages(), TestPrune_EmptyHistory(), TestPrune_HeadTailFormat() (+10 more)

### Community 37 - "renderStatusFooter"
Cohesion: 0.31
Nodes (9): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), getContextPill(), getGitStatus(), getShortCwd() (+1 more)

### Community 38 - "prompt_test.go"
Cohesion: 0.13
Nodes (13): TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestIsDangerousCommand(), TestMaxIterations() (+5 more)

### Community 39 - "context.Context"
Cohesion: 0.28
Nodes (9): context.Context, AnthropicProvider, CallCommandCode(), callGemini(), CallGeminiWithTools(), geminiDoRequest(), ChatMessage, GeminiProvider (+1 more)

### Community 40 - "RenameSession"
Cohesion: 0.33
Nodes (11): ClearChatSession(), historyFilePath(), saveHistoryToDisk(), DeleteSession(), ListSessions(), RenameSession(), sanitizeSessionID(), SessionExists() (+3 more)

### Community 41 - "GetStringArg"
Cohesion: 0.07
Nodes (48): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), TruncOutput(), GetAllTools() (+40 more)

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.18
Nodes (19): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask() (+11 more)

### Community 45 - "clarify.go"
Cohesion: 0.24
Nodes (10): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+2 more)

### Community 46 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 47 - "session_ui.go"
Cohesion: 0.44
Nodes (9): BuildSessionMenuKeyboard(), FormatSessionMenuText(), GetActiveSessionID(), HandleSessionCallback(), init(), loadTgSessionMapping(), saveTgSessionMapping(), SetActiveSessionID() (+1 more)

### Community 48 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 49 - "GetProvider"
Cohesion: 0.70
Nodes (4): LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 50 - "estimateHistoryTokens"
Cohesion: 0.53
Nodes (5): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, TestEstimateTokens()

### Community 51 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 52 - "HandleUploadInAgentMode"
Cohesion: 0.29
Nodes (6): agentSession, scorpChat(), touchSession(), RunAgentLoop(), TGDocument, HandleUploadInAgentMode()

### Community 53 - "ExecuteTermuxAPI"
Cohesion: 0.73
Nodes (5): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification()

### Community 54 - "upload.go"
Cohesion: 0.50
Nodes (4): contentPart, imageURL, TestBase64Encode(), base64Encode()

### Community 57 - "collector_docker.go"
Cohesion: 0.27
Nodes (11): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+3 more)

### Community 58 - "FormatHourlyReport"
Cohesion: 0.36
Nodes (10): SystemData, Bar(), bar(), FormatHourlyReport(), FormatStatusResponse(), SectionCoolify(), SectionDocker(), SectionSecurity() (+2 more)

### Community 59 - "collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 60 - "api_gemini.go"
Cohesion: 0.36
Nodes (10): geminiBuildRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall, geminiGenConfig, geminiPart, geminiRequest (+2 more)

### Community 63 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 64 - "session_fallback_test.go"
Cohesion: 0.40
Nodes (4): TestSessionFallback_ExecuteSessionSearch_RequiresQuery(), TestSessionFallback_InitAndIndex(), TestSessionFallback_LIKEQueryWithSpaces(), TestSessionFallback_SessionResultStruct()

### Community 65 - "api_anthropic.go"
Cohesion: 0.67
Nodes (3): anthropicRequest, anthropicResponse, anthropicTool

### Community 66 - "TestSwitchActiveModel"
Cohesion: 0.50
Nodes (3): TestBuildCommandCodePayload(), TestExtractFallbackToolCalls(), TestSwitchActiveModel()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 135 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `client.go`, `RenameSession`, `collector_security.go`, `collector_system.go`, `RunAgentSessionLoop`, `collector_system_native.go`, `clarify.go`, `time.Time`, `session_ui.go`, `HandleModelCallback`, `StartDaemon`, `TestSteeringQueue`, `startCLI`, `collector_docker.go`, `FormatHourlyReport`, `collector_coolify.go`, `getSession`?**
  _High betweenness centrality (0.152) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `skills.go`, `RegisterTool`, `chat.go`, `HandleTelegramAction`, `prompt_test.go`, `rag_vector.go`, `RenameSession`, `GetStringArg`, `time.Time`, `TestSteeringQueue`, `HandleUploadInAgentMode`, `startCLI`, `TruncateStr`, `ExecuteTermuxAPI`, `CallModelWithFallback`, `GenerateContextualSessionTitle`, `getSession`?**
  _High betweenness centrality (0.108) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `HandleTelegramAction`, `skills.go`, `chat.go`, `client.go`, `prompt_test.go`, `rag_vector.go`, `GetStringArg`, `bg.go`, `collector_system_native.go`, `HandleModelCallback`, `cost_router.go`, `time.Time`, `ResolveAPIKey`, `startCLI`, `collector_docker.go`?**
  _High betweenness centrality (0.080) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 25 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 25 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `ToolCall` be split into smaller, more focused modules?**
  _Cohesion score 0.14624505928853754 - nodes in this community are weakly interconnected._