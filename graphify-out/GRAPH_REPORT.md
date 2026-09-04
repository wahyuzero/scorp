# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 162 files · ~114,124 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1341 nodes · 3296 edges · 50 communities (41 shown, 1 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 379 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `85107415`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- context.Context
- StartDaemon
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
- time.Time
- scheduler.go
- session_search_fts5.go
- HandleTelegramAction
- cost_router.go
- main
- skills.go
- checker.go
- startCLI
- install.sh
- bg.go
- scorp-agent
- RunAgentLoop
- 🦂 Scorp
- compaction_test.go
- RegisterTool
- clarify.go
- RunAgentSessionLoop
- ExecuteTool
- TestGatewayEndpoints
- EscapeHTML
- collector_native_test.go
- GetStringArg
- HandleConfirmation
- GetProvider
- maybeCompactHistory
- LoadConfig
- upload.go
- StartServer

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
- `SetupInlineMode()` --calls--> `TgPost()`  [INFERRED]
  tools/inline.go → telegram/telegram.go
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

## Communities (50 total, 1 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.05
Nodes (96): scorpChat(), TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool (+88 more)

### Community 1 - "StartDaemon"
Cohesion: 0.17
Nodes (22): HandleUploadInAgentMode(), runCommandLoop(), StartDaemon(), AnswerCallback(), DeleteWebhook(), EditMessage(), EditMessageByID(), InitTelegram() (+14 more)

### Community 2 - "testing.T"
Cohesion: 0.11
Nodes (23): TestBase64Encode(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations() (+15 more)

### Community 3 - "chat.go"
Cohesion: 0.09
Nodes (49): agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown(), convertTableToList() (+41 more)

### Community 4 - "ScorpPath"
Cohesion: 0.06
Nodes (59): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+51 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (47): bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool(), buildArgDefsFromInputSchema(), executeMCPServerTool() (+39 more)

### Community 6 - "metasearch.go"
Cohesion: 0.05
Nodes (49): FormatDuration(), getUptime(), net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost() (+41 more)

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
Cohesion: 0.23
Nodes (16): TestCollectorNative_StartCPUSampler_DoesNotBlock(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+8 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (43): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), SaveModelConfig() (+35 more)

### Community 13 - "time.Time"
Cohesion: 0.21
Nodes (11): agentSession, TGDocument, time.Time, ModelUsage, AgentMessage, AutonomousConfig, AutonomousLogEntry, ChatSession (+3 more)

### Community 14 - "scheduler.go"
Cohesion: 0.19
Nodes (21): ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig() (+13 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "HandleTelegramAction"
Cohesion: 0.14
Nodes (26): HomeDir(), ProjectDir(), PythonSitePackages(), UploadsDir(), HandleTelegramAction(), BackKB(), createZip(), DirKeyboard() (+18 more)

### Community 17 - "cost_router.go"
Cohesion: 0.05
Nodes (61): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+53 more)

### Community 18 - "main"
Cohesion: 0.20
Nodes (14): RegisterAutonomous(), hasDebugFlag(), StartGateway(), isCLIMode(), main(), InitModelUsage(), SOP, Dir() (+6 more)

### Community 19 - "skills.go"
Cohesion: 0.07
Nodes (36): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+28 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.20
Nodes (22): executeOneShot(), executeTurn(), formatFinalResponse(), formatTerminalText(), handleCLISession(), handleCLISOP(), isTerminal(), printBanner() (+14 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "bg.go"
Cohesion: 0.12
Nodes (20): bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, CostTracker, bgKill(), bgList(), bgPoll() (+12 more)

### Community 31 - "RunAgentLoop"
Cohesion: 0.50
Nodes (4): RunAgentLoop(), runAgentTask(), getChatLock(), StartTestEndpoint()

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Roadmap Pengembangan Scorp (Next Upgrades), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (Local-First Parsing vs Cloud Fallback), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "compaction_test.go"
Cohesion: 0.25
Nodes (17): AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens(), TestPrune_BoundaryAges(), TestPrune_DockerScenario_41Messages(), TestPrune_EmptyHistory(), TestPrune_HeadTailFormat() (+9 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (36): init(), init(), init(), init(), unregisterMCPNativeTools(), TestGenerateNativeToolsSchema(), ActivateToolWithTTL(), IsDynamicModeEnabled() (+28 more)

### Community 35 - "clarify.go"
Cohesion: 0.31
Nodes (8): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID()

### Community 36 - "RunAgentSessionLoop"
Cohesion: 0.10
Nodes (31): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), looksLikeContinuation(), mentionsBrowserTask(), screenshotWasTaken() (+23 more)

### Community 37 - "ExecuteTool"
Cohesion: 0.18
Nodes (11): ExecuteTool(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), SetAutonomyLevel(), TestAutonomyLevels(), AutonomyLevel, RedactSecrets() (+3 more)

### Community 38 - "TestGatewayEndpoints"
Cohesion: 0.33
Nodes (11): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), TestGatewayEndpoints() (+3 more)

### Community 39 - "EscapeHTML"
Cohesion: 0.27
Nodes (12): EscapeHTML(), ShellTask(), AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker() (+4 more)

### Community 40 - "collector_native_test.go"
Cohesion: 0.23
Nodes (11): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+3 more)

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (50): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage() (+42 more)

### Community 42 - "HandleConfirmation"
Cohesion: 0.33
Nodes (9): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+1 more)

### Community 43 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 44 - "maybeCompactHistory"
Cohesion: 0.52
Nodes (6): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, maybeCompactHistory(), TestEstimateTokens()

### Community 45 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 46 - "upload.go"
Cohesion: 0.67
Nodes (3): contentPart, imageURL, base64Encode()

### Community 47 - "StartServer"
Cohesion: 0.67
Nodes (3): Init(), StartServer(), StopServer()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 132 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `StartDaemon`, `chat.go`, `RunAgentSessionLoop`, `client.go`, `clarify.go`, `collector_system.go`, `collector_security.go`, `HandleConfirmation`, `collector_system_native.go`, `wizard.go`, `scheduler.go`, `main`, `startCLI`, `RunAgentLoop`?**
  _High betweenness centrality (0.133) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `context.Context`, `chat.go`, `RunAgentSessionLoop`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `collector_system_native.go`, `wizard.go`, `scheduler.go`, `StartServer`, `HandleTelegramAction`, `cost_router.go`, `main`, `skills.go`, `RunAgentLoop`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `RegisterTool`, `RunAgentSessionLoop`, `client.go`, `rag_vector.go`, `runSubagent`, `HandleTelegramAction`, `main`, `skills.go`, `bg.go`?**
  _High betweenness centrality (0.081) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.0515596068936049 - nodes in this community are weakly interconnected._