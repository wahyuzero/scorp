# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 159 files · ~108,832 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1277 nodes · 3162 edges · 49 communities (40 shown, 1 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 375 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ff01118c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- TruncateStr
- StartDaemon
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
- bg.go
- HandleTelegramAction
- ComputeSimhash
- inline.go
- prompt_test.go
- GetStringArg
- clarify.go
- EscapeHTML
- ResetNativeToolCache
- StartServer
- LoadMCPConfig

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

## Communities (49 total, 1 thin omitted)

### Community 0 - "TruncateStr"
Cohesion: 0.05
Nodes (94): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+86 more)

### Community 1 - "StartDaemon"
Cohesion: 0.18
Nodes (21): runCommandLoop(), StartDaemon(), AnswerCallback(), DeleteWebhook(), EditMessage(), EditMessageByID(), InitTelegram(), PollUpdates() (+13 more)

### Community 2 - "testing.T"
Cohesion: 0.13
Nodes (21): TestSelfReviewCadence(), TestSelfReviewRateLimit(), TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty() (+13 more)

### Community 3 - "chat.go"
Cohesion: 0.07
Nodes (58): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+50 more)

### Community 4 - "MCPServer"
Cohesion: 0.18
Nodes (12): bufio.Scanner, encoding/json.Encoder, FindMCPTool(), MCPServer, startMCPServer(), jsonRPCResponse, MCPServerConfig, MCPTool (+4 more)

### Community 5 - "client.go"
Cohesion: 0.13
Nodes (24): ACPRequest, encoding/json.RawMessage, TruncOutputTool(), buildArgDefsFromInputSchema(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt() (+16 more)

### Community 6 - "compaction_test.go"
Cohesion: 0.19
Nodes (22): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+14 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.07
Nodes (34): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+26 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (53): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+45 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "runSubagent"
Cohesion: 0.09
Nodes (36): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+28 more)

### Community 11 - "uptime.go"
Cohesion: 0.12
Nodes (26): FormatDuration(), net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost(), getClient() (+18 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (43): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), SaveModelConfig() (+35 more)

### Community 13 - "startCLI"
Cohesion: 0.06
Nodes (66): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+58 more)

### Community 14 - "time.Time"
Cohesion: 0.25
Nodes (18): time.Time, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), FormatTasksList(), GetTask(), LoadTasks() (+10 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "collector_system_native.go"
Cohesion: 0.21
Nodes (18): CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses(), GetTopProcesses() (+10 more)

### Community 17 - "ConfigMgr"
Cohesion: 0.05
Nodes (62): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+54 more)

### Community 18 - "ScorpPath"
Cohesion: 0.05
Nodes (69): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+61 more)

### Community 19 - "skills.go"
Cohesion: 0.07
Nodes (37): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+29 more)

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
Cohesion: 0.09
Nodes (33): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), looksLikeContinuation(), mentionsBrowserTask(), screenshotWasTaken(), cleanToolCallTags() (+25 more)

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
Cohesion: 0.12
Nodes (15): RegisterAutonomous(), init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), unregisterMCPNativeTools() (+7 more)

### Community 35 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 36 - "bg.go"
Cohesion: 0.12
Nodes (20): bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, CostTracker, bgKill(), bgList(), bgPoll() (+12 more)

### Community 37 - "HandleTelegramAction"
Cohesion: 0.17
Nodes (22): HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath(), HumanSize() (+14 more)

### Community 38 - "ComputeSimhash"
Cohesion: 0.23
Nodes (11): tokenize(), ComputeSimhash(), hammingDistance(), simhashSimilarity(), TestSimHash_ComputeSimhash(), TestSimHash_HammingDistance(), TestSimHash_Similarity(), TestSimHash_TokenHash() (+3 more)

### Community 39 - "inline.go"
Cohesion: 0.33
Nodes (10): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+2 more)

### Community 40 - "prompt_test.go"
Cohesion: 0.20
Nodes (9): TestGetBoolArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations(), TestTruncOutput(), GetInt64Arg() (+1 more)

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (58): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage() (+50 more)

### Community 42 - "clarify.go"
Cohesion: 0.29
Nodes (8): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), SetClarifyChatID(), PendingClarify

### Community 43 - "EscapeHTML"
Cohesion: 0.36
Nodes (7): EscapeHTML(), isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage(), ShellTask()

### Community 44 - "ResetNativeToolCache"
Cohesion: 0.31
Nodes (6): TestGenerateNativeToolsSchema(), IsToolActive(), ResetDynamicTools(), TestDynamicToolTTL(), TickToolTTL(), ResetNativeToolCache()

### Community 45 - "StartServer"
Cohesion: 0.67
Nodes (3): Init(), StartServer(), StopServer()

### Community 48 - "LoadMCPConfig"
Cohesion: 0.50
Nodes (8): LoadMCPConfig(), ReloadMCPServers(), sanitizeMCPName(), ExecuteMCPManage(), mcpManageAdd(), mcpManageList(), mcpManageReload(), mcpManageRemove()

## Knowledge Gaps
- **25 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+20 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 116 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `TruncateStr`, `StartDaemon`, `chat.go`, `collector_system.go`, `collector_security.go`, `clarify.go`, `wizard.go`, `startCLI`, `time.Time`, `collector_system_native.go`, `LoadMCPConfig`, `RunAgentSessionLoop`?**
  _High betweenness centrality (0.154) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `RegisterTool`, `bg.go`, `rag_vector.go`, `runSubagent`, `startCLI`, `LoadMCPConfig`, `ScorpPath`, `skills.go`, `RunAgentSessionLoop`?**
  _High betweenness centrality (0.107) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `TruncateStr`, `RegisterTool`, `chat.go`, `client.go`, `HandleTelegramAction`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `wizard.go`, `startCLI`, `StartServer`, `time.Time`, `collector_system_native.go`, `ConfigMgr`, `skills.go`, `RunAgentSessionLoop`?**
  _High betweenness centrality (0.095) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _25 weakly-connected nodes found - possible documentation gaps or missing edges._