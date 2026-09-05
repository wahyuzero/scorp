# Graph Report - scorp  (2026-09-05)

## Corpus Check
- 187 files · ~123,363 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1438 nodes · 3521 edges · 55 communities (44 shown, 3 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 466 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4cdfb45f`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- time.Time
- bg.go
- main
- HandleTelegramAction
- ScorpPath
- client.go
- execute.go
- rag_vector.go
- collector_system.go
- collector_security.go
- Incident Report: Runaway CLI Verify Session
- MCPServer
- HandleModelCallback
- cost_router.go
- collector_system_native.go
- session_search_fts5.go
- context.Context
- StartDaemon
- ACPSession
- runSubagent
- checker.go
- startCLI
- install.sh
- LoadMCPConfig
- scorp-agent
- TestRegisterPlugin
- 🦂 Scorp
- v2_skills.go
- RegisterTool
- metasearch_engines.go
- testing.T
- registerMCPToolsAsNative
- skills.go
- browser_session.go
- acp.go
- GetStringArg
- .listenSSEStream
- RunAgentSessionLoop
- serviceBridgeRequests
- StartTestEndpoint
- CredentialVault
- TestPhase6_AllTools
- TodoManager
- GetRecentReceipts
- ExecuteScript
- HomeDir
- ExecuteAutoLogin

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 64 edges
2. `startCLI()` - 47 edges
3. `GetStringArg()` - 43 edges
4. `StartDaemon()` - 43 edges
5. `RunAgentSessionLoop()` - 40 edges
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
- `init()` --calls--> `RagVecProvider()`  [EXTRACTED]
  bootstrap/extended.go → rag/rag_vector.go

## Import Cycles
- None detected.

## Communities (55 total, 3 thin omitted)

### Community 0 - "time.Time"
Cohesion: 0.17
Nodes (25): time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult() (+17 more)

### Community 1 - "bg.go"
Cohesion: 0.38
Nodes (11): bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait(), bgWrite(), closeStdin(), ExecuteBgProcess() (+3 more)

### Community 2 - "main"
Cohesion: 0.09
Nodes (28): hasDebugFlag(), CM(), InitConfigManager(), NewConfigManager(), ConfigManager, contextWithTimeout(), handleChat(), handleDashboard() (+20 more)

### Community 3 - "HandleTelegramAction"
Cohesion: 0.07
Nodes (68): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+60 more)

### Community 4 - "ScorpPath"
Cohesion: 0.12
Nodes (25): init(), setupCLILogging(), BrowserSessionPath(), BrowserSessionsDir(), CostConfigFilePath(), CostLogFilePath(), DBConnectionsPath(), Hostname() (+17 more)

### Community 5 - "client.go"
Cohesion: 0.18
Nodes (19): ACPRequest, encoding/json.RawMessage, FindMCPTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError() (+11 more)

### Community 6 - "execute.go"
Cohesion: 0.26
Nodes (16): runOpenCodeCLI(), runSubagentACP(), AgentMessage, delegateBatchParams, delegateResult, delegateTaskParams, DefaultSubagentTools(), ExecuteDelegate() (+8 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.08
Nodes (51): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+43 more)

### Community 9 - "collector_security.go"
Cohesion: 0.18
Nodes (25): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+17 more)

### Community 10 - "Incident Report: Runaway CLI Verify Session"
Cohesion: 0.25
Nodes (7): Akar Masalah (dugaan), Ciri-ciri Runaway yang Terdeteksi, Incident Report: Runaway CLI Verify Session, Pelajaran Umum, Penanganan, Rekomendasi Perbaikan & Status Implementasi, Ringkasan

### Community 11 - "MCPServer"
Cohesion: 0.17
Nodes (12): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, MCPServer, startMCPServer(), jsonRPCError, jsonRPCResponse, MCPServerConfig (+4 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.08
Nodes (55): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+47 more)

### Community 13 - "cost_router.go"
Cohesion: 0.07
Nodes (56): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+48 more)

### Community 14 - "collector_system_native.go"
Cohesion: 0.07
Nodes (47): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+39 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "context.Context"
Cohesion: 0.06
Nodes (88): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicTool, callAnthropic(), CallAnthropicWithTools() (+80 more)

### Community 17 - "StartDaemon"
Cohesion: 0.05
Nodes (68): UploadsDir(), getUnclosedTags(), SplitMessage(), StopMCPServerMode(), Init(), StartServer(), StopServer(), runCommandLoop() (+60 more)

### Community 18 - "ACPSession"
Cohesion: 0.23
Nodes (8): launchACP(), ACPSession, bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, CostTracker, BGProcess

### Community 19 - "runSubagent"
Cohesion: 0.24
Nodes (10): cleanupSubagentSandbox(), createSubagentSandbox(), defaultIsolation(), formatIsolationInfo(), getSubagentIsolation(), isSubagentToolBlocked(), registerSubagentIsolation(), unregisterSubagentIsolation() (+2 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.06
Nodes (52): clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation, wireCLICallbacks() (+44 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "LoadMCPConfig"
Cohesion: 0.31
Nodes (11): LoadMCPConfig(), rebuildMCPToolList(), ReloadMCPServers(), sanitizeMCPName(), StartMCPServers(), StopMCPServers(), ExecuteMCPManage(), mcpManageAdd() (+3 more)

### Community 31 - "TestRegisterPlugin"
Cohesion: 0.20
Nodes (7): executeMCPServerTool(), echoPlugin, RegisterPlugin(), TestRegisterPlugin(), ExecuteToolByName(), ToolPlugin, ToolPluginWithSchema

### Community 32 - "🦂 Scorp"
Cohesion: 0.11
Nodes (17): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+9 more)

### Community 33 - "v2_skills.go"
Cohesion: 0.07
Nodes (33): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+25 more)

### Community 34 - "RegisterTool"
Cohesion: 0.10
Nodes (19): RegisterAutonomous(), init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), unregisterMCPNativeTools() (+11 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (40): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+32 more)

### Community 36 - "testing.T"
Cohesion: 0.06
Nodes (57): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+49 more)

### Community 37 - "registerMCPToolsAsNative"
Cohesion: 0.29
Nodes (8): TruncOutputTool(), buildArgDefsFromInputSchema(), registerMCPToolsAsNative(), ServerWatchdog, GetServerHealthStatus(), MCPServer, RegisterWatchdog(), RestartServer()

### Community 38 - "skills.go"
Cohesion: 0.22
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 39 - "browser_session.go"
Cohesion: 0.33
Nodes (13): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+5 more)

### Community 40 - "acp.go"
Cohesion: 0.33
Nodes (8): checkACPAvailable(), listAvailableACP(), ACPError, ACPInitializeParams, ACPMessageNewParams, ACPMessagePart, ACPResponse, ACPUserMessage

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (53): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), TruncOutput(), ActivateToolWithTTL() (+45 more)

### Community 42 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.06
Nodes (41): AgentMessage, sendScorpReply(), confirmKeyboard(), IsContinuationDirective(), IsPureInformationalQuery(), TestContinuationDirectives(), TestPureInformationalQuery(), cleanToolCallTags() (+33 more)

### Community 44 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 48 - "CredentialVault"
Cohesion: 0.31
Nodes (3): CredentialEntry, CredentialVault, ExecuteVault()

### Community 58 - "TestPhase6_AllTools"
Cohesion: 0.47
Nodes (8): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession()

### Community 59 - "TodoManager"
Cohesion: 0.33
Nodes (7): ExecuteTodo(), formatTodoList(), GetDefaultTodoManager(), getStringArgFromMap(), NewTodoManager(), TodoItem, TodoManager

### Community 63 - "GetRecentReceipts"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 64 - "ExecuteScript"
Cohesion: 0.52
Nodes (6): ExecuteScript(), ExecuteScriptList(), executeStep(), formatScriptResult(), ScriptResult, ScriptStep

### Community 67 - "HomeDir"
Cohesion: 0.24
Nodes (12): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), HomeDir(), ProjectDir() (+4 more)

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 142 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `time.Time`, `v2_skills.go`, `main`, `registerMCPToolsAsNative`, `collector_system.go`, `collector_security.go`, `RunAgentSessionLoop`, `HandleModelCallback`, `collector_system_native.go`, `context.Context`, `StartDaemon`, `startCLI`?**
  _High betweenness centrality (0.154) - this node is a cross-community bridge._
- **Why does `GetStringArg()` connect `GetStringArg` to `time.Time`, `v2_skills.go`, `ExecuteScript`, `main`, `testing.T`, `execute.go`, `browser_session.go`, `ExecuteAutoLogin`, `RunAgentSessionLoop`, `HandleModelCallback`, `cost_router.go`, `collector_system_native.go`, `CredentialVault`, `LoadMCPConfig`?**
  _High betweenness centrality (0.088) - this node is a cross-community bridge._
- **Why does `RunAgentSessionLoop()` connect `RunAgentSessionLoop` to `time.Time`, `v2_skills.go`, `HandleTelegramAction`, `rag_vector.go`, `GetStringArg`, `context.Context`, `startCLI`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 6 inferred relationships involving `startCLI()` (e.g. with `wireCLICallbacks()` and `formatTerminalText()`) actually correct?**
  _`startCLI()` has 6 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._