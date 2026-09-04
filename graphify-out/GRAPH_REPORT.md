# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 168 files · ~115,473 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1359 nodes · 3335 edges · 34 communities (25 shown, 1 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 423 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `994d4b7e`
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
- skills.go
- checker.go
- startCLI
- install.sh
- scorp-agent
- 🦂 Scorp
- RegisterTool
- RunAgentSessionLoop
- bg.go
- GetStringArg

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 54 edges
2. `GetStringArg()` - 43 edges
3. `StartDaemon()` - 40 edges
4. `RunAgentSessionLoop()` - 38 edges
5. `ModelConfig` - 36 edges
6. `TruncateStr()` - 35 edges
7. `startCLI()` - 34 edges
8. `init()` - 32 edges
9. `ChatMessage` - 30 edges
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

## Communities (34 total, 1 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.06
Nodes (87): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+79 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.06
Nodes (70): StopMCPServerMode(), Init(), StartServer(), StopServer(), FormatUsageStats(), InitModelUsage(), Load(), HandleTelegramAction() (+62 more)

### Community 2 - "testing.T"
Cohesion: 0.06
Nodes (56): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+48 more)

### Community 3 - "chat.go"
Cohesion: 0.08
Nodes (54): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+46 more)

### Community 4 - "ScorpPath"
Cohesion: 0.05
Nodes (58): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+50 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (46): bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool(), buildArgDefsFromInputSchema(), executeMCPServerTool() (+38 more)

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
Cohesion: 0.08
Nodes (39): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+31 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.09
Nodes (41): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+33 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.07
Nodes (56): CatalogEntry, defaultModelConfig(), LoadModelConfig(), SaveModelConfig(), CustomProvider, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog() (+48 more)

### Community 13 - "cost_router.go"
Cohesion: 0.05
Nodes (60): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+52 more)

### Community 14 - "time.Time"
Cohesion: 0.12
Nodes (31): time.Time, EscapeHTML(), ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath() (+23 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 19 - "skills.go"
Cohesion: 0.07
Nodes (36): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+28 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.06
Nodes (64): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+56 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (36): init(), init(), init(), init(), unregisterMCPNativeTools(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema(), ActivateToolWithTTL() (+28 more)

### Community 36 - "RunAgentSessionLoop"
Cohesion: 0.09
Nodes (35): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask() (+27 more)

### Community 38 - "bg.go"
Cohesion: 0.25
Nodes (15): bytes.Buffer, io.WriteCloser, os/exec.Cmd, bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait() (+7 more)

### Community 41 - "GetStringArg"
Cohesion: 0.05
Nodes (66): StorePendingConfirmation(), init(), init(), ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract() (+58 more)

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 133 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `chat.go`, `RunAgentSessionLoop`, `client.go`, `collector_system.go`, `collector_security.go`, `collector_system_native.go`, `HandleModelCallback`, `time.Time`, `startCLI`?**
  _High betweenness centrality (0.118) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `RegisterTool`, `ScorpPath`, `client.go`, `bg.go`, `rag_vector.go`, `RunAgentSessionLoop`, `runSubagent`, `skills.go`, `startCLI`?**
  _High betweenness centrality (0.088) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `chat.go`, `RunAgentSessionLoop`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `collector_system_native.go`, `HandleModelCallback`, `cost_router.go`, `time.Time`, `skills.go`, `startCLI`?**
  _High betweenness centrality (0.083) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.05801445953094692 - nodes in this community are weakly interconnected._