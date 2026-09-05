# Graph Report - scorp  (2026-09-05)

## Corpus Check
- 185 files · ~122,561 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1430 nodes · 3503 edges · 59 communities (50 shown, 1 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 462 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2f4aebd9`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- time.Time
- bg.go
- ConfigManager
- chat.go
- ScorpPath
- client.go
- runSubagent
- init
- collector_system.go
- collector_security.go
- Incident Report: Runaway CLI Verify Session
- ResolveAPIKey
- HandleModelCallback
- cost_router.go
- collector_system_native.go
- session_search_fts5.go
- context.Context
- StartDaemon
- ModelConfig
- api_gemini.go
- checker.go
- startCLI
- install.sh
- ToolCall
- scorp-agent
- collector_docker.go
- 🦂 Scorp
- runSelfReview
- RegisterTool
- metasearch_engines.go
- testing.T
- renderStatusFooter
- skills.go
- FormatHourlyReport
- api_commandcode.go
- GetStringArg
- collector_coolify.go
- RunAgentSessionLoop
- inline.go
- v2_skills.go
- getAgentSystemPrompt
- continuation.go
- ExecuteTermuxAPI
- ExecuteTool
- TruncateStr
- TestSteeringQueue
- RunAgentLoop
- api_anthropic.go
- TodoManager
- clarify.go
- HandleTelegramAction

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
- `browserSessionScreenshot()` --calls--> `SendFile()`  [INFERRED]
  browser/browser_session.go → telegram/files.go

## Import Cycles
- None detected.

## Communities (59 total, 1 thin omitted)

### Community 0 - "time.Time"
Cohesion: 0.17
Nodes (25): time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult() (+17 more)

### Community 1 - "bg.go"
Cohesion: 0.25
Nodes (15): bytes.Buffer, io.WriteCloser, os/exec.Cmd, bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait() (+7 more)

### Community 2 - "ConfigManager"
Cohesion: 0.13
Nodes (16): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts() (+8 more)

### Community 3 - "chat.go"
Cohesion: 0.07
Nodes (67): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+59 more)

### Community 4 - "ScorpPath"
Cohesion: 0.05
Nodes (58): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+50 more)

### Community 5 - "client.go"
Cohesion: 0.06
Nodes (54): ACPRequest, bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, sync.Mutex, TruncOutputTool() (+46 more)

### Community 6 - "runSubagent"
Cohesion: 0.09
Nodes (36): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+28 more)

### Community 7 - "init"
Cohesion: 0.06
Nodes (46): init(), activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir() (+38 more)

### Community 8 - "collector_system.go"
Cohesion: 0.23
Nodes (18): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+10 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "Incident Report: Runaway CLI Verify Session"
Cohesion: 0.25
Nodes (7): Akar Masalah (dugaan), Ciri-ciri Runaway yang Terdeteksi, Incident Report: Runaway CLI Verify Session, Pelajaran Umum, Penanganan, Rekomendasi Perbaikan, Ringkasan

### Community 11 - "ResolveAPIKey"
Cohesion: 0.17
Nodes (19): defaultModelConfig(), LoadModelConfig(), ModelRouterConfig, ModelUsage, applyProviderDefaults(), ProviderPreset, hasAPIKey(), migrateModelConfigs() (+11 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (44): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "cost_router.go"
Cohesion: 0.07
Nodes (57): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+49 more)

### Community 14 - "collector_system_native.go"
Cohesion: 0.08
Nodes (46): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+38 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "context.Context"
Cohesion: 0.18
Nodes (7): context.Context, AnthropicProvider, CallCommandCode(), ChatMessage, CommandCodeProvider, GeminiProvider, OpenAIProvider

### Community 17 - "StartDaemon"
Cohesion: 0.14
Nodes (24): Init(), StartServer(), StopServer(), StartDaemon(), BackAndRefreshKeyboard(), BackButtonKeyboard(), DeleteWebhook(), EditMessage() (+16 more)

### Community 18 - "ModelConfig"
Cohesion: 0.25
Nodes (16): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModel(), CallModelWithFallback(), findFirstVisionModel(), GetModelByName(), isVisionModelName() (+8 more)

### Community 19 - "api_gemini.go"
Cohesion: 0.32
Nodes (13): callGemini(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall (+5 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.07
Nodes (51): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+43 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "ToolCall"
Cohesion: 0.21
Nodes (10): TestParseToolCalls(), extractFallbackToolCalls(), CustomProvider, ToolCall, IsRateLimitError(), ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls() (+2 more)

### Community 31 - "collector_docker.go"
Cohesion: 0.27
Nodes (11): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+3 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.11
Nodes (17): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+9 more)

### Community 33 - "runSelfReview"
Cohesion: 0.18
Nodes (16): memoryFact, AgentMessage, maybeRunSelfReview(), runSelfReview(), TestSelfReviewIntegration(), LoadJSON(), SaveJSON(), SaveJSONPerm() (+8 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (41): init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), executeMCPServerTool(), unregisterMCPNativeTools() (+33 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (40): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+32 more)

### Community 36 - "testing.T"
Cohesion: 0.06
Nodes (57): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+49 more)

### Community 37 - "renderStatusFooter"
Cohesion: 0.23
Nodes (11): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste(), getContextPill() (+3 more)

### Community 38 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 39 - "FormatHourlyReport"
Cohesion: 0.31
Nodes (12): NetworkData, SystemData, Bar(), bar(), FormatHourlyReport(), FormatStatusResponse(), SectionCoolify(), SectionDocker() (+4 more)

### Community 40 - "api_commandcode.go"
Cohesion: 0.24
Nodes (12): buildCommandCodePayload(), CallCommandCodeStream(), createCommandCodeRequest(), ChatRequest, commandCodeMsg, commandCodeParams, commandCodePayload, commandCodeTool (+4 more)

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (59): StorePendingConfirmation(), init(), ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill() (+51 more)

### Community 42 - "collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.16
Nodes (18): AgentMessage, cleanToolCallTags(), getSessionSearchContext(), maxIterations(), resumeAgentLoop(), RunAgentSessionLoop(), TestBuildThinkingMessage(), TestIsDangerousCommand() (+10 more)

### Community 44 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 45 - "v2_skills.go"
Cohesion: 0.24
Nodes (10): SkillMeta, ActivateSkill(), ListSkillsOverview(), LoadAllSkills(), ParseSkillMetadata(), ReadSkillBody(), scanLegacyJSONSkills(), scanSkillsDirectory() (+2 more)

### Community 46 - "getAgentSystemPrompt"
Cohesion: 0.22
Nodes (8): getSharedMemorySummary(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), FormatSkillsIndexForSystemPrompt(), GetActiveSkillsContext(), GetMemorySummary()

### Community 47 - "continuation.go"
Cohesion: 0.29
Nodes (4): IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery()

### Community 48 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 49 - "ExecuteTool"
Cohesion: 0.29
Nodes (4): ExecuteTool(), FormatToolResult(), RedactSecrets(), TestRedactSecrets()

### Community 50 - "TruncateStr"
Cohesion: 0.26
Nodes (18): TruncateStr(), callAnthropic(), CallAnthropicWithTools(), CallCommandCodeWithTools(), CallOpenAI(), CallOpenAIWithTools(), formatOpenAIMessages(), RecordCost() (+10 more)

### Community 51 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 52 - "RunAgentLoop"
Cohesion: 0.67
Nodes (3): RunAgentLoop(), getChatLock(), StartTestEndpoint()

### Community 53 - "api_anthropic.go"
Cohesion: 0.67
Nodes (3): anthropicRequest, anthropicResponse, anthropicTool

### Community 59 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 62 - "clarify.go"
Cohesion: 0.27
Nodes (9): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID() (+1 more)

### Community 66 - "HandleTelegramAction"
Cohesion: 0.14
Nodes (27): FormatUsageStats(), HandleTelegramAction(), runCommandLoop(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo() (+19 more)

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 141 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `time.Time`, `chat.go`, `client.go`, `FormatHourlyReport`, `collector_system.go`, `collector_security.go`, `collector_coolify.go`, `RunAgentSessionLoop`, `HandleModelCallback`, `v2_skills.go`, `collector_system_native.go`, `StartDaemon`, `TestSteeringQueue`, `startCLI`, `clarify.go`, `collector_docker.go`?**
  _High betweenness centrality (0.144) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `time.Time`, `runSelfReview`, `RegisterTool`, `chat.go`, `HandleTelegramAction`, `init`, `GetStringArg`, `v2_skills.go`, `getAgentSystemPrompt`, `continuation.go`, `ExecuteTermuxAPI`, `ExecuteTool`, `TruncateStr`, `TestSteeringQueue`, `RunAgentLoop`, `startCLI`, `ModelConfig`?**
  _High betweenness centrality (0.090) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `time.Time`, `ConfigManager`, `chat.go`, `ScorpPath`, `client.go`, `init`, `ResolveAPIKey`, `HandleModelCallback`, `cost_router.go`, `collector_system_native.go`, `startCLI`, `collector_docker.go`, `runSelfReview`, `skills.go`, `GetStringArg`, `RunAgentSessionLoop`, `v2_skills.go`, `ExecuteTool`, `RunAgentLoop`, `HandleTelegramAction`?**
  _High betweenness centrality (0.081) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 5 inferred relationships involving `startCLI()` (e.g. with `wireCLICallbacks()` and `formatTerminalText()`) actually correct?**
  _`startCLI()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._