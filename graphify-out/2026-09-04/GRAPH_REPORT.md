# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 159 files · ~108,530 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1276 nodes · 3160 edges · 45 communities (36 shown, 1 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 375 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `af9049f2`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- TruncateStr
- HandleTelegramAction
- testing.T
- chat.go
- MCPServer
- client.go
- compaction_test.go
- rag_vector.go
- collector_system.go
- collector_security.go
- runSubagent
- uptime.go
- wizard.go
- startCLI
- time.Time
- session_search_fts5.go
- collector_system_native.go
- ConfigMgr
- ScorpPath
- skills.go
- checker.go
- GetProvider
- install.sh
- RunAgentSessionLoop
- scorp-agent
- TestRegisterPlugin
- 🦂 Scorp
- .listenSSEStream
- RegisterTool
- serviceBridgeRequests
- GetStringArg
- TestSessionManager
- TestSteeringQueue
- LoadMCPConfig
- ExecuteTermuxAPI
- collector_native_test.go
- HandleUploadInAgentMode

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 54 edges
2. `GetStringArg()` - 43 edges
3. `StartDaemon()` - 40 edges
4. `RunAgentSessionLoop()` - 36 edges
5. `TruncateStr()` - 36 edges
6. `ModelConfig` - 36 edges
7. `startCLI()` - 34 edges
8. `init()` - 32 edges
9. `ChatMessage` - 29 edges
10. `HandleModelCallback()` - 29 edges

## Surprising Connections (you probably didn't know these)
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `ExecuteShell()` --calls--> `IsDangerousCommand()`  [INFERRED]
  tools/exec.go → agent/prompt.go
- `serviceBridgeRequests()` --calls--> `ExecuteTool()`  [INFERRED]
  tools/exec_code.go → agent/prompt.go
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go
- `init()` --calls--> `RagVecProvider()`  [EXTRACTED]
  bootstrap/extended.go → rag/rag_vector.go

## Import Cycles
- None detected.

## Communities (45 total, 1 thin omitted)

### Community 0 - "TruncateStr"
Cohesion: 0.05
Nodes (92): AgentMessage, runSelfReview(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool (+84 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.05
Nodes (73): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), HomeDir(), PythonSitePackages() (+65 more)

### Community 2 - "testing.T"
Cohesion: 0.10
Nodes (26): TestBase64Encode(), TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg() (+18 more)

### Community 3 - "chat.go"
Cohesion: 0.12
Nodes (37): agentAutoStop(), appendSessionHistory(), collectTableLines(), convertInlineMarkdown(), convertTableToList(), EnterAgentMode(), ExitAgentMode(), extractAndSaveMemory() (+29 more)

### Community 4 - "MCPServer"
Cohesion: 0.20
Nodes (10): bufio.Scanner, encoding/json.Encoder, MCPServer, startMCPServer(), jsonRPCResponse, MCPServerConfig, isRemoteMCP(), startSSEServer() (+2 more)

### Community 5 - "client.go"
Cohesion: 0.18
Nodes (19): encoding/json.RawMessage, FindMCPTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt(), sendMCPError(), sendMCPResult() (+11 more)

### Community 6 - "compaction_test.go"
Cohesion: 0.19
Nodes (22): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+14 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (52): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+44 more)

### Community 9 - "collector_security.go"
Cohesion: 0.18
Nodes (25): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+17 more)

### Community 10 - "runSubagent"
Cohesion: 0.06
Nodes (56): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+48 more)

### Community 11 - "uptime.go"
Cohesion: 0.16
Nodes (21): FormatDuration(), getUptime(), net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost() (+13 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (44): ProjectDir(), CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "startCLI"
Cohesion: 0.05
Nodes (65): clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation, ExecuteTool() (+57 more)

### Community 14 - "time.Time"
Cohesion: 0.12
Nodes (32): time.Time, EscapeHTML(), ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath() (+24 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "collector_system_native.go"
Cohesion: 0.23
Nodes (16): TestCollectorNative_StartCPUSampler_DoesNotBlock(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+8 more)

### Community 17 - "ConfigMgr"
Cohesion: 0.05
Nodes (61): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+53 more)

### Community 18 - "ScorpPath"
Cohesion: 0.05
Nodes (62): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+54 more)

### Community 19 - "skills.go"
Cohesion: 0.08
Nodes (33): getSharedMemorySummary(), FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), TestSelfReviewIntegration(), LoadJSON() (+25 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "RunAgentSessionLoop"
Cohesion: 0.16
Nodes (20): AgentMessage, confirmKeyboard(), countStepsInMessage(), AgentMessage, looksLikeContinuation(), mentionsBrowserTask(), screenshotWasTaken(), cleanToolCallTags() (+12 more)

### Community 31 - "TestRegisterPlugin"
Cohesion: 0.20
Nodes (7): executeMCPServerTool(), echoPlugin, RegisterPlugin(), TestRegisterPlugin(), ExecuteToolByName(), ToolPlugin, ToolPluginWithSchema

### Community 32 - "🦂 Scorp"
Cohesion: 0.11
Nodes (17): 1. Diet Ekstrem `main.go`, 2. Modularisasi `agent/loop.go`, 3. Pemisahan Provider AI Independen (`models/`), 🚀 Panduan Perintah Cepat (Quick Cheatsheet), 🧱 Rekapitulasi Refactoring Arsitektur Bersih (*Clean Architecture*), 📌 Ringkasan Pembaruan Utama (v2.0), 🦂 Scorp Agent v2.0 — Modernization & Architecture Upgrade Report, 🧪 Verifikasi & Kualitas Kode (+9 more)

### Community 33 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 34 - "RegisterTool"
Cohesion: 0.11
Nodes (20): init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), TruncOutputTool(), buildArgDefsFromInputSchema() (+12 more)

### Community 35 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 41 - "GetStringArg"
Cohesion: 0.05
Nodes (63): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage() (+55 more)

### Community 44 - "TestSessionManager"
Cohesion: 0.32
Nodes (12): ClearChatSession(), historyFilePath(), saveHistoryToDisk(), DeleteSession(), ListSessions(), RenameSession(), sanitizeSessionID(), SessionExists() (+4 more)

### Community 47 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 48 - "LoadMCPConfig"
Cohesion: 0.27
Nodes (12): GetStringSliceArg(), LoadMCPConfig(), rebuildMCPToolList(), ReloadMCPServers(), sanitizeMCPName(), StartMCPServers(), StopMCPServers(), ExecuteMCPManage() (+4 more)

### Community 50 - "ExecuteTermuxAPI"
Cohesion: 0.73
Nodes (5): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification()

### Community 51 - "collector_native_test.go"
Cohesion: 0.23
Nodes (11): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+3 more)

### Community 61 - "HandleUploadInAgentMode"
Cohesion: 0.22
Nodes (9): agentSession, cleanupChatSessions(), CleanupSessionsLoop(), contentPart, imageURL, cleanupAgentSessions(), TGDocument, base64Encode() (+1 more)

## Knowledge Gaps
- **25 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+20 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 115 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `chat.go`, `collector_system.go`, `collector_security.go`, `TestSessionManager`, `startCLI`, `time.Time`, `TestSteeringQueue`, `collector_system_native.go`, `LoadMCPConfig`, `wizard.go`, `RunAgentSessionLoop`?**
  _High betweenness centrality (0.128) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `HandleTelegramAction`, `RegisterTool`, `rag_vector.go`, `runSubagent`, `startCLI`, `LoadMCPConfig`, `ExecuteTermuxAPI`, `skills.go`?**
  _High betweenness centrality (0.104) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `TruncateStr`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `runSubagent`, `wizard.go`, `startCLI`, `time.Time`, `collector_system_native.go`, `ConfigMgr`, `LoadMCPConfig`, `skills.go`, `RunAgentSessionLoop`, `HandleUploadInAgentMode`?**
  _High betweenness centrality (0.097) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _25 weakly-connected nodes found - possible documentation gaps or missing edges._