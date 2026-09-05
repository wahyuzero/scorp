# Graph Report - scorp  (2026-09-05)

## Corpus Check
- 184 files · ~123,397 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1445 nodes · 3524 edges · 78 communities (66 shown, 4 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 467 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `cb205454`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- context.Context
- TestGatewayEndpoints
- MCPServer
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
- agent/autonomous.go
- time.Time
- session_search_fts5.go
- cost_router.go
- StartDaemon
- ModelConfig
- ConfigManager
- checker.go
- startCLI
- install.sh
- ChatMessage
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
- 🦂 Scorp Native Go-MCP Marketplace & Transpiler Engine
- RunAgentSessionLoop
- inline.go
- main
- ExecuteTool
- ToolCall
- saveQuickstartConfig
- collector_docker.go
- TruncateStr
- TestSteeringQueue
- HandleUploadInAgentMode
- registerMCPToolsAsNative
- FormatHourlyReport
- collector_coolify.go
- wireCLICallbacks
- TodoManager
- api_gemini.go
- HandleTelegramAction
- clarify.go
- GenerateContextualSessionTitle
- HandleConfirmation
- tools/monitor.go
- HomeDir
- LoadMCPConfig
- .listenSSEStream
- ExecuteAutonomous
- echoPlugin
- serviceBridgeRequests
- LoadConfig
- providerTest
- GetRecentReceipts
- ConfigMgr
- GetProvider
- .CallWithTools

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
- `ExecuteAutonomous()` --calls--> `RunAutonomousCycle()`  [INFERRED]
  tools/autonomous.go → agent/autonomous.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `serviceBridgeRequests()` --calls--> `ExecuteTool()`  [INFERRED]
  tools/exec_code.go → agent/prompt.go
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go

## Import Cycles
- None detected.

## Communities (78 total, 4 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.18
Nodes (15): context.Context, buildCommandCodePayload(), CallCommandCode(), CallCommandCodeStream(), createCommandCodeRequest(), commandCodeMsg, commandCodeParams, commandCodePayload (+7 more)

### Community 1 - "TestGatewayEndpoints"
Cohesion: 0.33
Nodes (11): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), TestGatewayEndpoints() (+3 more)

### Community 2 - "MCPServer"
Cohesion: 0.19
Nodes (11): bufio.Scanner, encoding/json.Encoder, FindMCPTool(), MCPServer, startMCPServer(), MCPServerConfig, MCPTool, MCPServer (+3 more)

### Community 3 - "chat.go"
Cohesion: 0.20
Nodes (16): collectTableLines(), convertInlineMarkdown(), convertTableToList(), extractAndSaveMemory(), historyWriterLoop(), init(), isAnyModeActive(), isSeparatorRow() (+8 more)

### Community 4 - "ScorpPath"
Cohesion: 0.05
Nodes (57): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+49 more)

### Community 5 - "client.go"
Cohesion: 0.18
Nodes (19): ACPRequest, encoding/json.RawMessage, executeMCPServerTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError() (+11 more)

### Community 6 - "runSubagent"
Cohesion: 0.09
Nodes (36): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+28 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.23
Nodes (18): checkGDriveMount(), checkS3Gateway(), CollectNetwork(), CollectStorage(), detectNewPorts(), getDockerVolumeSizes(), getEstablishedConnections(), getListeningPorts() (+10 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "bg.go"
Cohesion: 0.25
Nodes (15): bytes.Buffer, io.WriteCloser, os/exec.Cmd, bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait() (+7 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.07
Nodes (47): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+39 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (44): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "agent/autonomous.go"
Cohesion: 0.18
Nodes (19): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), makeDecision() (+11 more)

### Community 14 - "time.Time"
Cohesion: 0.16
Nodes (26): printCronTasks(), time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath() (+18 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "cost_router.go"
Cohesion: 0.21
Nodes (15): TestGetFloatArg(), GetFloatArg(), defaultCostConfig(), formatCostReport(), FormatDailyCostSummary(), handleCostCommand(), init(), isBudgetExceeded() (+7 more)

### Community 17 - "StartDaemon"
Cohesion: 0.14
Nodes (26): getUnclosedTags(), SplitMessage(), StopMCPServerMode(), InitModelUsage(), runCommandLoop(), StartDaemon(), AnswerCallback(), BackAndRefreshKeyboard() (+18 more)

### Community 18 - "ModelConfig"
Cohesion: 0.26
Nodes (13): defaultModelConfig(), LoadModelConfig(), ModelConfig, ModelRouterConfig, ModelUsage, applyProviderDefaults(), hasAPIKey(), migrateModelConfigs() (+5 more)

### Community 19 - "ConfigManager"
Cohesion: 0.21
Nodes (5): CM(), InitConfigManager(), NewConfigManager(), ConfigManager, os.FileMode

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.35
Nodes (14): executeOneShot(), executeTurn(), handleCLISession(), handleCLISOP(), printBanner(), printCLIHelp(), printCostUsage(), printCurrentModel() (+6 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "ChatMessage"
Cohesion: 0.24
Nodes (16): ChatMessage, ChatRequest, ChatResponse, getCheapestModel(), RouteModelCostAware(), CostTracker, CallModel(), CallModelWithFallback() (+8 more)

### Community 31 - "getSession"
Cohesion: 0.16
Nodes (25): agentAutoStop(), appendSessionHistory(), EnterAgentMode(), ExitAgentMode(), flushPendingMessages(), GetHistoryTokenEstimate(), getOrCreateSession(), getSession() (+17 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "v2_skills.go"
Cohesion: 0.07
Nodes (34): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+26 more)

### Community 34 - "RegisterTool"
Cohesion: 0.11
Nodes (19): init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), unregisterMCPNativeTools(), RegisterPlugin() (+11 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (40): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+32 more)

### Community 36 - "testing.T"
Cohesion: 0.05
Nodes (58): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+50 more)

### Community 37 - "renderStatusFooter"
Cohesion: 0.23
Nodes (11): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), disableBracketedPaste(), enableBracketedPaste(), getContextPill() (+3 more)

### Community 38 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 39 - ".CallWithTools"
Cohesion: 0.29
Nodes (4): AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool

### Community 40 - "RenameSession"
Cohesion: 0.33
Nodes (11): ClearChatSession(), historyFilePath(), saveHistoryToDisk(), DeleteSession(), ListSessions(), RenameSession(), sanitizeSessionID(), SessionExists() (+3 more)

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (56): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), TruncOutput(), ActivateToolWithTTL() (+48 more)

### Community 42 - "🦂 Scorp Native Go-MCP Marketplace & Transpiler Engine"
Cohesion: 0.12
Nodes (15): 🎯 1. Executive Summary & Vision, 🧭 2. User Experience & Decision Flow, 🏗️ 3. Arsitektur Komponen Teknis, 🛡️ 4. Etika Open Source & Legal Compliance, 🚀 5. Roadmap Pelaksanaan, A. The Go-MCP Transpiler Pipeline (`mcp/transpiler/`), Architectural Blueprint & System Design Specification, B. Format Manifesto Marketplace (`mcp-manifest.json`) (+7 more)

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.15
Nodes (24): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask() (+16 more)

### Community 44 - "inline.go"
Cohesion: 0.23
Nodes (12): TestWebhookHandlerMalformedJSON(), WebhookHandler(), AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker() (+4 more)

### Community 45 - "main"
Cohesion: 0.21
Nodes (13): RegisterAutonomous(), hasDebugFlag(), StartGateway(), isCLIMode(), main(), SOP, Dir(), GetSOP() (+5 more)

### Community 46 - "ExecuteTool"
Cohesion: 0.21
Nodes (10): ExecuteTool(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), SetAutonomyLevel(), TestAutonomyLevels(), AutonomyLevel, SettingsMenuText() (+2 more)

### Community 47 - "ToolCall"
Cohesion: 0.31
Nodes (8): TestParseToolCalls(), extractFallbackToolCalls(), CustomProvider, ToolCall, CallModelWithTools(), ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls()

### Community 49 - "collector_docker.go"
Cohesion: 0.27
Nodes (11): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), ContainerInfo (+3 more)

### Community 50 - "TruncateStr"
Cohesion: 0.26
Nodes (18): TruncateStr(), callAnthropic(), CallAnthropicWithTools(), CallCommandCodeWithTools(), CallOpenAI(), CallOpenAIWithTools(), formatOpenAIMessages(), RecordCost() (+10 more)

### Community 51 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 52 - "HandleUploadInAgentMode"
Cohesion: 0.25
Nodes (7): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), RunAgentLoop(), cleanupAgentSessions(), TGDocument, HandleUploadInAgentMode()

### Community 53 - "registerMCPToolsAsNative"
Cohesion: 0.20
Nodes (11): sync.Mutex, TruncOutputTool(), buildArgDefsFromInputSchema(), registerMCPToolsAsNative(), ServerWatchdog, GetServerHealthStatus(), MCPServer, RegisterWatchdog() (+3 more)

### Community 54 - "FormatHourlyReport"
Cohesion: 0.31
Nodes (12): NetworkData, SystemData, Bar(), bar(), FormatHourlyReport(), FormatStatusResponse(), SectionCoolify(), SectionDocker() (+4 more)

### Community 57 - "collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 58 - "wireCLICallbacks"
Cohesion: 0.39
Nodes (6): wireCLICallbacks(), formatFinalResponse(), formatTerminalText(), isTerminal(), stripHTML(), os.File

### Community 59 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 60 - "api_gemini.go"
Cohesion: 0.23
Nodes (14): callGemini(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall (+6 more)

### Community 61 - "HandleTelegramAction"
Cohesion: 0.39
Nodes (11): FormatUsageStats(), HandleTelegramAction(), BuildSessionMenuKeyboard(), FormatSessionMenuText(), GetActiveSessionID(), HandleSessionCallback(), init(), loadTgSessionMapping() (+3 more)

### Community 62 - "clarify.go"
Cohesion: 0.27
Nodes (9): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID() (+1 more)

### Community 63 - "GenerateContextualSessionTitle"
Cohesion: 0.60
Nodes (4): fallbackTitleFromText(), GenerateContextualSessionTitle(), sanitizeSessionTitle(), ShouldAutoTitleSession()

### Community 64 - "HandleConfirmation"
Cohesion: 0.33
Nodes (9): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+1 more)

### Community 65 - "tools/monitor.go"
Cohesion: 0.38
Nodes (10): ExecuteMonitor(), InitMonitor(), loadMonitorTargets(), monitorCheckOne(), monitorLoop(), ragIngestText(), sanitizeFilename(), saveMonitorTargets() (+2 more)

### Community 66 - "HomeDir"
Cohesion: 0.15
Nodes (24): HomeDir(), ProjectDir(), PythonSitePackages(), UploadsDir(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard() (+16 more)

### Community 67 - "LoadMCPConfig"
Cohesion: 0.31
Nodes (11): LoadMCPConfig(), rebuildMCPToolList(), ReloadMCPServers(), sanitizeMCPName(), StartMCPServers(), StopMCPServers(), ExecuteMCPManage(), mcpManageAdd() (+3 more)

### Community 68 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 69 - "ExecuteAutonomous"
Cohesion: 0.33
Nodes (9): SaveAutonomousConfig(), saveAutonomousConfigLocked(), SetKillSwitch(), TestPhase7_KillSwitch(), autoShowActions(), autoShowConfig(), autoShowLog(), autoStatus() (+1 more)

### Community 71 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 72 - "LoadConfig"
Cohesion: 0.33
Nodes (9): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), Init(), StartServer() (+1 more)

### Community 73 - "providerTest"
Cohesion: 0.42
Nodes (8): ProviderPreset, getProviderInfo(), HandleProviderCommand(), listConfiguredModels(), listProvidersWithKeyStatus(), providerAddInteractive(), providerRemove(), providerTest()

### Community 74 - "GetRecentReceipts"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 75 - "ConfigMgr"
Cohesion: 0.60
Nodes (5): LoadAutonomousConfig(), setupTestPaths(), TestPhase7_ConfigPersistence(), ConfigMgr(), saveCostTracker()

### Community 76 - "GetProvider"
Cohesion: 0.70
Nodes (4): LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

## Knowledge Gaps
- **43 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+38 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 150 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `collector_system.go`, `collector_security.go`, `collector_system_native.go`, `HandleModelCallback`, `time.Time`, `StartDaemon`, `getSession`, `v2_skills.go`, `RenameSession`, `RunAgentSessionLoop`, `main`, `ExecuteTool`, `collector_docker.go`, `TestSteeringQueue`, `registerMCPToolsAsNative`, `FormatHourlyReport`, `collector_coolify.go`, `clarify.go`, `HandleConfirmation`, `HomeDir`?**
  _High betweenness centrality (0.142) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `HandleConfirmation`, `v2_skills.go`, `chat.go`, `testing.T`, `rag_vector.go`, `RenameSession`, `GetStringArg`, `ExecuteTool`, `time.Time`, `TruncateStr`, `TestSteeringQueue`, `HandleUploadInAgentMode`, `startCLI`, `HandleTelegramAction`, `ChatMessage`, `GenerateContextualSessionTitle`, `getSession`?**
  _High betweenness centrality (0.085) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `ScorpPath`, `client.go`, `rag_vector.go`, `collector_system_native.go`, `HandleModelCallback`, `agent/autonomous.go`, `time.Time`, `cost_router.go`, `ModelConfig`, `ConfigManager`, `v2_skills.go`, `testing.T`, `skills.go`, `GetStringArg`, `main`, `ExecuteTool`, `collector_docker.go`, `HandleUploadInAgentMode`, `registerMCPToolsAsNative`, `HomeDir`, `LoadMCPConfig`, `LoadConfig`?**
  _High betweenness centrality (0.081) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 5 inferred relationships involving `startCLI()` (e.g. with `wireCLICallbacks()` and `formatTerminalText()`) actually correct?**
  _`startCLI()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 25 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 25 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _43 weakly-connected nodes found - possible documentation gaps or missing edges._