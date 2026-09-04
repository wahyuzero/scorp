# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 168 files · ~115,386 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1358 nodes · 3333 edges · 54 communities (45 shown, 1 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 422 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2b6c54df`
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
- ResolveAPIKey
- ConfigManager
- StartDaemon
- runSelfReview
- checker.go
- startCLI
- install.sh
- browser_session.go
- scorp-agent
- ExecuteTermuxAPI
- 🦂 Scorp
- CallCommandCodeWithTools
- RegisterTool
- TestSteeringQueue
- resumeAgentLoop
- skills.go
- bg.go
- ModelConfig
- RunAgentLoop
- GetStringArg
- IsDangerousCommand
- ToolCall
- api_gemini.go
- clarify.go
- inline.go
- RunAgentSessionLoop
- EscapeHTML
- .CallWithTools
- StartServer
- telemetry.go

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 54 edges
2. `GetStringArg()` - 43 edges
3. `StartDaemon()` - 40 edges
4. `RunAgentSessionLoop()` - 37 edges
5. `ModelConfig` - 36 edges
6. `TruncateStr()` - 35 edges
7. `startCLI()` - 34 edges
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

## Communities (54 total, 1 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.21
Nodes (21): context.Context, TruncateStr(), callAnthropic(), CallAnthropicWithTools(), CallCommandCodeStream(), callGemini(), CallGeminiWithTools(), geminiDoRequest() (+13 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.17
Nodes (24): FormatUsageStats(), HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath() (+16 more)

### Community 2 - "testing.T"
Cohesion: 0.06
Nodes (54): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+46 more)

### Community 3 - "chat.go"
Cohesion: 0.08
Nodes (54): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+46 more)

### Community 4 - "ScorpPath"
Cohesion: 0.07
Nodes (49): init(), hasDebugFlag(), setupCLILogging(), Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr() (+41 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (46): bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool(), buildArgDefsFromInputSchema(), FindMCPTool() (+38 more)

### Community 6 - "metasearch.go"
Cohesion: 0.06
Nodes (39): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+31 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (53): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+45 more)

### Community 9 - "collector_security.go"
Cohesion: 0.18
Nodes (25): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+17 more)

### Community 10 - "runSubagent"
Cohesion: 0.08
Nodes (37): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+29 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.08
Nodes (44): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+36 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (44): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "cost_router.go"
Cohesion: 0.07
Nodes (57): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+49 more)

### Community 14 - "time.Time"
Cohesion: 0.25
Nodes (18): time.Time, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), FormatTasksList(), GetTask(), LoadTasks() (+10 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "ResolveAPIKey"
Cohesion: 0.17
Nodes (19): defaultModelConfig(), LoadModelConfig(), ModelRouterConfig, ModelUsage, applyProviderDefaults(), ProviderPreset, hasAPIKey(), migrateModelConfigs() (+11 more)

### Community 17 - "ConfigManager"
Cohesion: 0.11
Nodes (18): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts() (+10 more)

### Community 18 - "StartDaemon"
Cohesion: 0.20
Nodes (20): runCommandLoop(), StartDaemon(), AnswerCallback(), DeleteWebhook(), EditMessage(), EditMessageByID(), InitTelegram(), PollUpdates() (+12 more)

### Community 19 - "runSelfReview"
Cohesion: 0.19
Nodes (15): memoryFact, AgentMessage, maybeRunSelfReview(), runSelfReview(), TestSelfReviewIntegration(), LoadJSON(), SaveJSON(), SaveJSONPerm() (+7 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.08
Nodes (46): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+38 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "browser_session.go"
Cohesion: 0.09
Nodes (34): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), ExecuteBrowser() (+26 more)

### Community 31 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "CallCommandCodeWithTools"
Cohesion: 0.16
Nodes (17): buildCommandCodePayload(), CallCommandCode(), CallCommandCodeWithTools(), createCommandCodeRequest(), extractFallbackToolCalls(), ChatRequest, commandCodeMsg, commandCodeParams (+9 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (40): RegisterAutonomous(), init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), executeMCPServerTool() (+32 more)

### Community 35 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 36 - "resumeAgentLoop"
Cohesion: 0.26
Nodes (10): AgentMessage, cleanToolCallTags(), maxIterations(), resumeAgentLoop(), TestBuildThinkingMessage(), TestToolDescription(), buildThinkingMessage(), shouldUpdateThinking() (+2 more)

### Community 37 - "skills.go"
Cohesion: 0.24
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 38 - "bg.go"
Cohesion: 0.25
Nodes (15): bytes.Buffer, io.WriteCloser, os/exec.Cmd, bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait() (+7 more)

### Community 39 - "ModelConfig"
Cohesion: 0.26
Nodes (14): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModel(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), isVisionModelName() (+6 more)

### Community 40 - "RunAgentLoop"
Cohesion: 0.67
Nodes (3): RunAgentLoop(), getChatLock(), StartTestEndpoint()

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (50): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage() (+42 more)

### Community 42 - "IsDangerousCommand"
Cohesion: 0.15
Nodes (10): getSharedMemorySummary(), FormatToolResult(), getAgentSystemPrompt(), IsDangerousCommand(), TestIsDangerousCommand(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap() (+2 more)

### Community 43 - "ToolCall"
Cohesion: 0.24
Nodes (11): TestParseToolCalls(), CustomProvider, ToolCall, CallModelWithTools(), CallModelWithToolsAndFallback(), IsRateLimitError(), ParseAllToolCalls(), ParseCodeBlockFallback() (+3 more)

### Community 44 - "api_gemini.go"
Cohesion: 0.36
Nodes (10): geminiBuildRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall, geminiGenConfig, geminiPart, geminiRequest (+2 more)

### Community 45 - "clarify.go"
Cohesion: 0.27
Nodes (9): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID() (+1 more)

### Community 46 - "inline.go"
Cohesion: 0.33
Nodes (10): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+2 more)

### Community 47 - "RunAgentSessionLoop"
Cohesion: 0.42
Nodes (8): countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), looksLikeContinuation(), mentionsBrowserTask(), screenshotWasTaken(), RunAgentSessionLoop()

### Community 48 - "EscapeHTML"
Cohesion: 0.36
Nodes (7): EscapeHTML(), isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage(), ShellTask()

### Community 49 - ".CallWithTools"
Cohesion: 0.33
Nodes (4): AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool

### Community 50 - "StartServer"
Cohesion: 0.67
Nodes (3): Init(), StartServer(), StopServer()

### Community 51 - "telemetry.go"
Cohesion: 0.67
Nodes (3): CheckModelHealth(), FormatModelList(), FormatModelListWithHealth()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 133 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `TestSteeringQueue`, `chat.go`, `client.go`, `ScorpPath`, `RunAgentLoop`, `collector_system.go`, `collector_security.go`, `collector_system_native.go`, `HandleModelCallback`, `clarify.go`, `time.Time`, `StartDaemon`, `startCLI`?**
  _High betweenness centrality (0.125) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `RegisterTool`, `ScorpPath`, `client.go`, `skills.go`, `rag_vector.go`, `bg.go`, `runSubagent`, `runSelfReview`, `ExecuteTermuxAPI`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Why does `TruncateStr()` connect `context.Context` to `CallCommandCodeWithTools`, `chat.go`, `resumeAgentLoop`, `client.go`, `GetStringArg`, `runSubagent`, `cost_router.go`, `time.Time`, `RunAgentSessionLoop`, `EscapeHTML`, `ResolveAPIKey`, `runSelfReview`, `telemetry.go`, `browser_session.go`?**
  _High betweenness centrality (0.076) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.062003968253968256 - nodes in this community are weakly interconnected._