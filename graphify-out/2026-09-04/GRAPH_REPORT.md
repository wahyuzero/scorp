# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 162 files · ~113,791 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1337 nodes · 3288 edges · 40 communities (31 shown, 1 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 379 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c8d91074`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- context.Context
- telegram.go
- testing.T
- chat.go
- browser_session.go
- client.go
- metasearch.go
- rag_vector.go
- collector_system.go
- collector_security.go
- runSubagent
- collector_system_native.go
- wizard.go
- time.Time
- scheduler.go
- session_search_fts5.go
- HandleTelegramAction
- cost_router.go
- ScorpPath
- skills.go
- checker.go
- install.sh
- bg.go
- scorp-agent
- StartDaemon
- 🦂 Scorp
- RegisterTool
- clarify.go
- RunAgentSessionLoop
- EscapeHTML
- runScriptTask
- GetStringArg

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 54 edges
2. `GetStringArg()` - 43 edges
3. `StartDaemon()` - 40 edges
4. `RunAgentSessionLoop()` - 37 edges
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

## Communities (40 total, 1 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.05
Nodes (94): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+86 more)

### Community 1 - "telegram.go"
Cohesion: 0.16
Nodes (18): BackAndRefreshKeyboard(), baseName(), DeleteWebhook(), EditMessage(), EditMessageByID(), MainMenuKeyboard(), PollUpdates(), SendChatAction() (+10 more)

### Community 2 - "testing.T"
Cohesion: 0.07
Nodes (53): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, maybeCompactHistory(), AgentMessage, makeHistory(), makeToolResult() (+45 more)

### Community 3 - "chat.go"
Cohesion: 0.08
Nodes (55): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+47 more)

### Community 4 - "browser_session.go"
Cohesion: 0.09
Nodes (34): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), ExecuteBrowser() (+26 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (47): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool(), buildArgDefsFromInputSchema() (+39 more)

### Community 6 - "metasearch.go"
Cohesion: 0.08
Nodes (27): net/http.Client, BraveSearchEngine, DuckDuckGoHTMLEngine, DuckDuckGoLiteEngine, GitHubEngine, deduplicateAndRank(), GetMetaSearchAggregator(), NewBraveSearchEngine() (+19 more)

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
Cohesion: 0.07
Nodes (48): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+40 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (43): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), SaveModelConfig() (+35 more)

### Community 13 - "time.Time"
Cohesion: 0.31
Nodes (8): time.Time, ModelUsage, AgentMessage, AutonomousLogEntry, ChatSession, PendingClarify, pendingConfirmation, TgResponse

### Community 14 - "scheduler.go"
Cohesion: 0.26
Nodes (17): ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), FormatTasksList(), GetTask(), LoadTasks(), Loop() (+9 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "HandleTelegramAction"
Cohesion: 0.23
Nodes (18): FormatUsageStats(), HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath() (+10 more)

### Community 17 - "cost_router.go"
Cohesion: 0.05
Nodes (61): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+53 more)

### Community 18 - "ScorpPath"
Cohesion: 0.05
Nodes (66): RegisterAutonomous(), init(), hasDebugFlag(), setupCLILogging(), Config, EnvBool(), EnvFloat(), EnvInt() (+58 more)

### Community 19 - "skills.go"
Cohesion: 0.07
Nodes (36): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+28 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "bg.go"
Cohesion: 0.25
Nodes (15): bytes.Buffer, io.WriteCloser, os/exec.Cmd, bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait() (+7 more)

### Community 31 - "StartDaemon"
Cohesion: 0.17
Nodes (13): RunAgentLoop(), Init(), StartServer(), StopServer(), InitModelUsage(), runCommandLoop(), StartDaemon(), InitTelegram() (+5 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Roadmap Pengembangan Scorp (Next Upgrades), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (Local-First Parsing vs Cloud Fallback), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (38): init(), init(), init(), init(), executeMCPServerTool(), unregisterMCPNativeTools(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema() (+30 more)

### Community 35 - "clarify.go"
Cohesion: 0.27
Nodes (9): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+1 more)

### Community 36 - "RunAgentSessionLoop"
Cohesion: 0.05
Nodes (71): AgentMessage, clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation() (+63 more)

### Community 39 - "EscapeHTML"
Cohesion: 0.30
Nodes (11): EscapeHTML(), ShellTask(), AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker() (+3 more)

### Community 40 - "runScriptTask"
Cohesion: 0.48
Nodes (5): isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage()

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (50): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage() (+42 more)

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 131 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `telegram.go`, `chat.go`, `RunAgentSessionLoop`, `client.go`, `clarify.go`, `collector_system.go`, `collector_security.go`, `collector_system_native.go`, `wizard.go`, `scheduler.go`, `ScorpPath`, `StartDaemon`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `context.Context`, `telegram.go`, `chat.go`, `RunAgentSessionLoop`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `collector_system_native.go`, `wizard.go`, `scheduler.go`, `HandleTelegramAction`, `cost_router.go`, `ScorpPath`, `skills.go`?**
  _High betweenness centrality (0.090) - this node is a cross-community bridge._
- **Why does `GetStringArg()` connect `GetStringArg` to `context.Context`, `testing.T`, `RegisterTool`, `RunAgentSessionLoop`, `client.go`, `browser_session.go`, `runSubagent`, `collector_system_native.go`, `scheduler.go`, `cost_router.go`, `ScorpPath`, `skills.go`?**
  _High betweenness centrality (0.079) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.05352323838080959 - nodes in this community are weakly interconnected._