# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 160 files · ~110,902 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1291 nodes · 3184 edges · 48 communities (39 shown, 1 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 376 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a4b0e785`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- TruncateStr
- StartDaemon
- testing.T
- chat.go
- agent/autonomous.go
- client.go
- compaction_test.go
- rag_vector.go
- collector_system.go
- collector_security.go
- runSubagent
- collector_system_native.go
- wizard.go
- startCLI
- time.Time
- session_search_fts5.go
- HandleTelegramAction
- cost_router.go
- ScorpPath
- skills.go
- checker.go
- GetProvider
- install.sh
- bg.go
- scorp-agent
- main
- 🦂 Scorp
- TestGatewayEndpoints
- RegisterTool
- clarify.go
- RunAgentSessionLoop
- ExecuteTool
- HandleConfirmation
- inline.go
- EscapeHTML
- GetStringArg
- GetRecentReceipts
- ExecuteAutonomous
- formatFinalResponse
- serviceBridgeRequests

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

## Communities (48 total, 1 thin omitted)

### Community 0 - "TruncateStr"
Cohesion: 0.06
Nodes (81): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+73 more)

### Community 1 - "StartDaemon"
Cohesion: 0.13
Nodes (25): Init(), StartServer(), StopServer(), runCommandLoop(), StartDaemon(), baseName(), DeleteWebhook(), EditMessage() (+17 more)

### Community 2 - "testing.T"
Cohesion: 0.09
Nodes (32): TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations(), TestTruncOutput() (+24 more)

### Community 3 - "chat.go"
Cohesion: 0.07
Nodes (58): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+50 more)

### Community 4 - "agent/autonomous.go"
Cohesion: 0.14
Nodes (25): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+17 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (46): bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool(), buildArgDefsFromInputSchema(), FindMCPTool() (+38 more)

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
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "runSubagent"
Cohesion: 0.08
Nodes (39): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+31 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.07
Nodes (45): TestCollectorNative_CollectSystem_Structure(), FormatDuration(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes() (+37 more)

### Community 12 - "wizard.go"
Cohesion: 0.09
Nodes (52): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), LoadModelConfig() (+44 more)

### Community 13 - "startCLI"
Cohesion: 0.25
Nodes (19): executeOneShot(), executeTurn(), formatTerminalText(), handleCLISession(), handleCLISOP(), isTerminal(), printBanner(), printCLIHelp() (+11 more)

### Community 14 - "time.Time"
Cohesion: 0.21
Nodes (20): time.Time, ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), FormatTasksList(), GetTask() (+12 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "HandleTelegramAction"
Cohesion: 0.20
Nodes (20): HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath(), HumanSize() (+12 more)

### Community 17 - "cost_router.go"
Cohesion: 0.09
Nodes (30): CM(), ConfigMgr(), InitConfigManager(), NewConfigManager(), ConfigManager, os.FileMode, defaultCostConfig(), formatCostReport() (+22 more)

### Community 18 - "ScorpPath"
Cohesion: 0.05
Nodes (69): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+61 more)

### Community 19 - "skills.go"
Cohesion: 0.07
Nodes (36): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+28 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "bg.go"
Cohesion: 0.12
Nodes (20): bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, CostTracker, bgKill(), bgList(), bgPoll() (+12 more)

### Community 31 - "main"
Cohesion: 0.20
Nodes (14): RegisterAutonomous(), hasDebugFlag(), StartGateway(), isCLIMode(), main(), InitModelUsage(), SOP, Dir() (+6 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Roadmap Pengembangan Scorp (Next Upgrades), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (Local-First Parsing vs Cloud Fallback), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "TestGatewayEndpoints"
Cohesion: 0.31
Nodes (12): contextWithTimeout(), handleChat(), handleDashboard(), handleReceipts(), handleSOPs(), handleStatus(), handleTools(), TestGatewayEndpoints() (+4 more)

### Community 34 - "RegisterTool"
Cohesion: 0.07
Nodes (34): init(), init(), init(), init(), executeMCPServerTool(), unregisterMCPNativeTools(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema() (+26 more)

### Community 35 - "clarify.go"
Cohesion: 0.27
Nodes (9): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+1 more)

### Community 36 - "RunAgentSessionLoop"
Cohesion: 0.09
Nodes (34): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), looksLikeContinuation(), mentionsBrowserTask(), screenshotWasTaken() (+26 more)

### Community 37 - "ExecuteTool"
Cohesion: 0.23
Nodes (9): ExecuteTool(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), SetAutonomyLevel(), TestAutonomyLevels(), AutonomyLevel, RedactSecrets() (+1 more)

### Community 38 - "HandleConfirmation"
Cohesion: 0.36
Nodes (8): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation

### Community 39 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 40 - "EscapeHTML"
Cohesion: 0.36
Nodes (7): EscapeHTML(), isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage(), ShellTask()

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (51): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage() (+43 more)

### Community 42 - "GetRecentReceipts"
Cohesion: 0.46
Nodes (6): GetRecentReceipts(), loadReceiptsLocked(), RecordToolReceipt(), saveReceiptsLocked(), TestRecordToolReceipt(), ToolReceipt

### Community 43 - "ExecuteAutonomous"
Cohesion: 0.52
Nodes (6): SaveAutonomousConfig(), autoShowActions(), autoShowConfig(), autoShowLog(), autoStatus(), ExecuteAutonomous()

### Community 44 - "formatFinalResponse"
Cohesion: 0.50
Nodes (4): formatFinalResponse(), stripHTML(), TestFormatFinalResponse(), TestStripHTML()

### Community 45 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 123 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `StartDaemon`, `chat.go`, `RunAgentSessionLoop`, `client.go`, `HandleConfirmation`, `clarify.go`, `collector_system.go`, `collector_security.go`, `collector_system_native.go`, `wizard.go`, `startCLI`, `time.Time`, `main`?**
  _High betweenness centrality (0.145) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `chat.go`, `agent/autonomous.go`, `RunAgentSessionLoop`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `collector_system_native.go`, `wizard.go`, `time.Time`, `HandleTelegramAction`, `cost_router.go`, `skills.go`, `main`?**
  _High betweenness centrality (0.100) - this node is a cross-community bridge._
- **Why does `GetStringArg()` connect `GetStringArg` to `testing.T`, `RegisterTool`, `RunAgentSessionLoop`, `client.go`, `runSubagent`, `collector_system_native.go`, `wizard.go`, `time.Time`, `cost_router.go`, `ScorpPath`, `skills.go`, `bg.go`, `main`?**
  _High betweenness centrality (0.093) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 21 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._