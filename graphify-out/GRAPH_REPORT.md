# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 169 files · ~116,521 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1364 nodes · 3344 edges · 58 communities (48 shown, 2 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 424 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ba4031bf`
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
- ConfigManager
- StartDaemon
- ResolveAPIKey
- runSelfReview
- checker.go
- startCLI
- install.sh
- CallCommandCodeWithTools
- scorp-agent
- api_gemini.go
- 🦂 Scorp
- skills.go
- RegisterTool
- browser_session.go
- RunAgentSessionLoop
- ModelConfig
- EscapeHTML
- ToolCall
- CredentialVault
- GetStringArg
- clarify.go
- resumeAgentLoop
- getAgentSystemPrompt
- TestPhase6_AllTools
- ExecuteTool
- GetRecentReceipts
- ExecuteTermuxAPI
- TestSteeringQueue
- ExecuteScript
- LoadConfig
- telemetry.go
- RunAgentLoop
- StartServer
- ExecuteAutoLogin

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

## Communities (58 total, 2 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.20
Nodes (19): context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic(), CallAnthropicWithTools() (+11 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.17
Nodes (24): FormatUsageStats(), HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath() (+16 more)

### Community 2 - "testing.T"
Cohesion: 0.06
Nodes (63): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+55 more)

### Community 3 - "chat.go"
Cohesion: 0.07
Nodes (58): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+50 more)

### Community 4 - "ScorpPath"
Cohesion: 0.12
Nodes (27): init(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), HomeDir(), Hostname() (+19 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (47): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool(), buildArgDefsFromInputSchema() (+39 more)

### Community 6 - "metasearch.go"
Cohesion: 0.06
Nodes (39): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+31 more)

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
Cohesion: 0.06
Nodes (56): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+48 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.09
Nodes (38): buildThinkingMessage(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), FormatDuration(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg() (+30 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (44): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "cost_router.go"
Cohesion: 0.07
Nodes (57): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+49 more)

### Community 14 - "time.Time"
Cohesion: 0.16
Nodes (24): time.Time, ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult() (+16 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "ConfigManager"
Cohesion: 0.12
Nodes (16): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts() (+8 more)

### Community 17 - "StartDaemon"
Cohesion: 0.18
Nodes (22): StopMCPServerMode(), InitModelUsage(), runCommandLoop(), StartDaemon(), AnswerCallback(), DeleteWebhook(), EditMessage(), EditMessageByID() (+14 more)

### Community 18 - "ResolveAPIKey"
Cohesion: 0.18
Nodes (18): defaultModelConfig(), LoadModelConfig(), ModelRouterConfig, applyProviderDefaults(), ProviderPreset, hasAPIKey(), migrateModelConfigs(), ResolveAPIKey() (+10 more)

### Community 19 - "runSelfReview"
Cohesion: 0.18
Nodes (16): memoryFact, AgentMessage, maybeRunSelfReview(), runSelfReview(), TestSelfReviewIntegration(), LoadJSON(), SaveJSON(), SaveJSONPerm() (+8 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.07
Nodes (54): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+46 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "CallCommandCodeWithTools"
Cohesion: 0.17
Nodes (16): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), ChatRequest, commandCodeMsg, commandCodeParams (+8 more)

### Community 31 - "api_gemini.go"
Cohesion: 0.23
Nodes (14): callGemini(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall (+6 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "skills.go"
Cohesion: 0.24
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (36): init(), init(), init(), init(), unregisterMCPNativeTools(), TestGenerateNativeToolsSchema(), ActivateToolWithTTL(), IsDynamicModeEnabled() (+28 more)

### Community 35 - "browser_session.go"
Cohesion: 0.33
Nodes (13): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+5 more)

### Community 36 - "RunAgentSessionLoop"
Cohesion: 0.38
Nodes (9): countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask(), screenshotWasTaken() (+1 more)

### Community 37 - "ModelConfig"
Cohesion: 0.26
Nodes (14): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModel(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), isVisionModelName() (+6 more)

### Community 38 - "EscapeHTML"
Cohesion: 0.24
Nodes (13): EscapeHTML(), runAgentTask(), ShellTask(), AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery() (+5 more)

### Community 39 - "ToolCall"
Cohesion: 0.27
Nodes (10): CustomProvider, ToolCall, CallModelWithTools(), CallModelWithToolsAndFallback(), IsRateLimitError(), ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls() (+2 more)

### Community 40 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (54): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetFloatArg(), GetInt64Arg(), GetIntArg(), GetStringArg() (+46 more)

### Community 42 - "clarify.go"
Cohesion: 0.31
Nodes (8): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID()

### Community 43 - "resumeAgentLoop"
Cohesion: 0.36
Nodes (7): AgentMessage, cleanToolCallTags(), maxIterations(), resumeAgentLoop(), shouldUpdateThinking(), toolCallSignature(), toolDescription()

### Community 44 - "getAgentSystemPrompt"
Cohesion: 0.25
Nodes (7): getSharedMemorySummary(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), GetPromptForMessage(), GetMemorySummary()

### Community 45 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 46 - "ExecuteTool"
Cohesion: 0.22
Nodes (5): ExecuteTool(), FormatToolResult(), IsDangerousCommand(), RedactSecrets(), TestRedactSecrets()

### Community 47 - "GetRecentReceipts"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 48 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 49 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 50 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 51 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 52 - "telemetry.go"
Cohesion: 0.47
Nodes (5): CheckModelHealth(), FormatModelList(), FormatModelListWithHealth(), TrackModelUsage(), TrackModelUsageWithCache()

### Community 53 - "RunAgentLoop"
Cohesion: 0.67
Nodes (3): RunAgentLoop(), getChatLock(), StartTestEndpoint()

### Community 54 - "StartServer"
Cohesion: 0.67
Nodes (3): Init(), StartServer(), StopServer()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 133 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `chat.go`, `client.go`, `collector_system.go`, `collector_security.go`, `clarify.go`, `collector_system_native.go`, `HandleModelCallback`, `time.Time`, `StartDaemon`, `TestSteeringQueue`, `RunAgentLoop`, `startCLI`?**
  _High betweenness centrality (0.130) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `skills.go`, `RegisterTool`, `ScorpPath`, `client.go`, `rag_vector.go`, `runSubagent`, `ExecuteTermuxAPI`, `runSelfReview`, `startCLI`?**
  _High betweenness centrality (0.090) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `skills.go`, `HandleTelegramAction`, `chat.go`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `collector_system_native.go`, `HandleModelCallback`, `cost_router.go`, `ExecuteTool`, `time.Time`, `ConfigManager`, `ResolveAPIKey`, `runSelfReview`, `startCLI`, `StartServer`, `RunAgentLoop`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.05516475379489078 - nodes in this community are weakly interconnected._