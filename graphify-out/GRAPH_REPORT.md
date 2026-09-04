# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 162 files · ~114,915 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1349 nodes · 3312 edges · 45 communities (35 shown, 2 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 380 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `669e59e5`
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
- wizard.go
- agent/autonomous.go
- time.Time
- session_search_fts5.go
- ExecuteTool
- cost_router.go
- HandleConfirmation
- skills.go
- checker.go
- startCLI
- install.sh
- TestPhase6_AllTools
- scorp-agent
- ExecuteTermuxAPI
- 🦂 Scorp
- compaction_test.go
- RegisterTool
- TestSteeringQueue
- RunAgentSessionLoop
- formatFinalResponse
- ExecuteAutonomous
- upload.go
- StartTestEndpoint
- GetStringArg
- GetProvider

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
- `ExecuteAutonomous()` --calls--> `SetKillSwitch()`  [INFERRED]
  tools/autonomous.go → agent/autonomous.go
- `ExecuteAutonomous()` --calls--> `RunAutonomousCycle()`  [INFERRED]
  tools/autonomous.go → agent/autonomous.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `ExecuteShell()` --calls--> `IsDangerousCommand()`  [INFERRED]
  tools/exec.go → agent/prompt.go
- `serviceBridgeRequests()` --calls--> `ExecuteTool()`  [INFERRED]
  tools/exec_code.go → agent/prompt.go

## Import Cycles
- None detected.

## Communities (45 total, 2 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.05
Nodes (97): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+89 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.06
Nodes (67): HomeDir(), ProjectDir(), UploadsDir(), Init(), StartServer(), StopServer(), HandleTelegramAction(), runCommandLoop() (+59 more)

### Community 2 - "testing.T"
Cohesion: 0.11
Nodes (23): TestBase64Encode(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations() (+15 more)

### Community 3 - "chat.go"
Cohesion: 0.08
Nodes (56): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+48 more)

### Community 4 - "ScorpPath"
Cohesion: 0.05
Nodes (64): init(), hasDebugFlag(), setupCLILogging(), Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr() (+56 more)

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
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "runSubagent"
Cohesion: 0.06
Nodes (56): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+48 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.09
Nodes (40): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+32 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (43): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), SaveModelConfig() (+35 more)

### Community 13 - "agent/autonomous.go"
Cohesion: 0.14
Nodes (25): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+17 more)

### Community 14 - "time.Time"
Cohesion: 0.12
Nodes (30): time.Time, ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult() (+22 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "ExecuteTool"
Cohesion: 0.23
Nodes (9): ExecuteTool(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), SetAutonomyLevel(), TestAutonomyLevels(), AutonomyLevel, RedactSecrets() (+1 more)

### Community 17 - "cost_router.go"
Cohesion: 0.09
Nodes (30): CM(), ConfigMgr(), InitConfigManager(), NewConfigManager(), ConfigManager, os.FileMode, defaultCostConfig(), formatCostReport() (+22 more)

### Community 18 - "HandleConfirmation"
Cohesion: 0.43
Nodes (7): clearPendingConfirmation(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation

### Community 19 - "skills.go"
Cohesion: 0.07
Nodes (36): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+28 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.25
Nodes (19): executeOneShot(), executeTurn(), formatTerminalText(), handleCLISession(), handleCLISOP(), isTerminal(), printBanner(), printCLIHelp() (+11 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "TestPhase6_AllTools"
Cohesion: 0.14
Nodes (17): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession() (+9 more)

### Community 31 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "compaction_test.go"
Cohesion: 0.19
Nodes (22): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+14 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (38): RegisterAutonomous(), init(), init(), init(), init(), executeMCPServerTool(), unregisterMCPNativeTools(), TestGenerateNativeToolsSchema() (+30 more)

### Community 35 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 36 - "RunAgentSessionLoop"
Cohesion: 0.16
Nodes (22): AgentMessage, confirmKeyboard(), countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), looksLikeContinuation(), mentionsBrowserTask() (+14 more)

### Community 37 - "formatFinalResponse"
Cohesion: 0.50
Nodes (4): formatFinalResponse(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 38 - "ExecuteAutonomous"
Cohesion: 0.52
Nodes (6): SaveAutonomousConfig(), autoShowActions(), autoShowConfig(), autoShowLog(), autoStatus(), ExecuteAutonomous()

### Community 39 - "upload.go"
Cohesion: 0.67
Nodes (3): contentPart, imageURL, base64Encode()

### Community 41 - "GetStringArg"
Cohesion: 0.05
Nodes (66): StorePendingConfirmation(), init(), init(), ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract() (+58 more)

### Community 43 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 133 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `TestSteeringQueue`, `chat.go`, `RunAgentSessionLoop`, `client.go`, `ScorpPath`, `collector_system.go`, `collector_security.go`, `collector_system_native.go`, `wizard.go`, `time.Time`, `HandleConfirmation`, `startCLI`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `context.Context`, `RegisterTool`, `chat.go`, `ScorpPath`, `RunAgentSessionLoop`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `StartTestEndpoint`, `collector_system_native.go`, `wizard.go`, `agent/autonomous.go`, `time.Time`, `cost_router.go`, `skills.go`?**
  _High betweenness centrality (0.087) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `RegisterTool`, `ScorpPath`, `client.go`, `rag_vector.go`, `runSubagent`, `skills.go`, `ExecuteTermuxAPI`?**
  _High betweenness centrality (0.083) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.05084033613445378 - nodes in this community are weakly interconnected._