# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 159 files · ~110,027 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1282 nodes · 3175 edges · 46 communities (36 shown, 2 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 376 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ae66f821`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- TruncateStr
- HandleTelegramAction
- testing.T
- chat.go
- agent/autonomous.go
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
- skills.go
- cost_router.go
- ScorpPath
- runSelfReview
- checker.go
- GetProvider
- install.sh
- bg.go
- scorp-agent
- TestPhase6_AllTools
- 🦂 Scorp
- collector_system_native.go
- RegisterTool
- collector_native_test.go
- RunAgentSessionLoop
- ExecuteTool
- StorePendingConfirmation
- ExecuteTermuxAPI
- TestSteeringQueue
- GetStringArg
- upload.go
- StartTestEndpoint

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
- `ExecuteSQL()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/db.go → agent/confirmation.go
- `ExecuteShell()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/exec.go → agent/confirmation.go
- `ExecuteGit()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/git.go → agent/confirmation.go
- `ExecuteProcess()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/process.go → agent/confirmation.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go

## Import Cycles
- None detected.

## Communities (46 total, 2 thin omitted)

### Community 0 - "TruncateStr"
Cohesion: 0.05
Nodes (94): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+86 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.06
Nodes (66): HomeDir(), PythonSitePackages(), Init(), StartServer(), StopServer(), HandleTelegramAction(), runCommandLoop(), StartDaemon() (+58 more)

### Community 2 - "testing.T"
Cohesion: 0.11
Nodes (23): TestBase64Encode(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations() (+15 more)

### Community 3 - "chat.go"
Cohesion: 0.08
Nodes (54): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+46 more)

### Community 4 - "agent/autonomous.go"
Cohesion: 0.12
Nodes (31): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+23 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (47): bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, TruncOutputTool(), buildArgDefsFromInputSchema(), executeMCPServerTool() (+39 more)

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

### Community 11 - "uptime.go"
Cohesion: 0.11
Nodes (27): FormatDuration(), getUptime(), net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost() (+19 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (44): ProjectDir(), CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "startCLI"
Cohesion: 0.20
Nodes (22): executeOneShot(), executeTurn(), formatFinalResponse(), formatTerminalText(), handleCLISession(), handleCLISOP(), isTerminal(), printBanner() (+14 more)

### Community 14 - "time.Time"
Cohesion: 0.16
Nodes (26): time.Time, EscapeHTML(), ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath() (+18 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "skills.go"
Cohesion: 0.24
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 17 - "cost_router.go"
Cohesion: 0.09
Nodes (30): CM(), ConfigMgr(), InitConfigManager(), NewConfigManager(), ConfigManager, os.FileMode, defaultCostConfig(), formatCostReport() (+22 more)

### Community 18 - "ScorpPath"
Cohesion: 0.05
Nodes (63): init(), hasDebugFlag(), setupCLILogging(), Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr() (+55 more)

### Community 19 - "runSelfReview"
Cohesion: 0.10
Nodes (23): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+15 more)

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
Cohesion: 0.25
Nodes (15): bytes.Buffer, io.WriteCloser, os/exec.Cmd, bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait() (+7 more)

### Community 31 - "TestPhase6_AllTools"
Cohesion: 0.13
Nodes (19): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), CloseBrowserSession() (+11 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.11
Nodes (17): 1. Diet Ekstrem `main.go`, 2. Modularisasi `agent/loop.go`, 3. Pemisahan Provider AI Independen (`models/`), 🚀 Panduan Perintah Cepat (Quick Cheatsheet), 🧱 Rekapitulasi Refactoring Arsitektur Bersih (*Clean Architecture*), 📌 Ringkasan Pembaruan Utama (v2.0), 🦂 Scorp Agent v2.0 — Modernization & Architecture Upgrade Report, 🧪 Verifikasi & Kualitas Kode (+9 more)

### Community 33 - "collector_system_native.go"
Cohesion: 0.23
Nodes (16): TestCollectorNative_StartCPUSampler_DoesNotBlock(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+8 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (39): RegisterAutonomous(), init(), init(), init(), init(), unregisterMCPNativeTools(), TestCallModelWithToolsNilModel(), TestGenerateNativeToolsSchema() (+31 more)

### Community 35 - "collector_native_test.go"
Cohesion: 0.21
Nodes (12): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+4 more)

### Community 36 - "RunAgentSessionLoop"
Cohesion: 0.16
Nodes (21): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), looksLikeContinuation(), mentionsBrowserTask(), screenshotWasTaken() (+13 more)

### Community 37 - "ExecuteTool"
Cohesion: 0.23
Nodes (9): ExecuteTool(), GetAutonomyLevel(), IsPathRestricted(), IsToolAllowed(), SetAutonomyLevel(), TestAutonomyLevels(), AutonomyLevel, RedactSecrets() (+1 more)

### Community 38 - "StorePendingConfirmation"
Cohesion: 0.31
Nodes (10): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), StorePendingConfirmation() (+2 more)

### Community 39 - "ExecuteTermuxAPI"
Cohesion: 0.46
Nodes (6): AcquireTermuxWakeLock(), ExecuteTermuxAPI(), IsTermux(), ReleaseTermuxWakeLock(), SendTermuxNotification(), TestExecuteTermuxAPI_Simulation()

### Community 40 - "TestSteeringQueue"
Cohesion: 0.43
Nodes (5): ClearSteeringQueue(), HasSteeringMessage(), PopSteeringMessage(), QueueSteeringMessage(), TestSteeringQueue()

### Community 41 - "GetStringArg"
Cohesion: 0.05
Nodes (65): init(), init(), ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill() (+57 more)

### Community 42 - "upload.go"
Cohesion: 0.67
Nodes (3): contentPart, imageURL, base64Encode()

## Knowledge Gaps
- **26 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+21 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 117 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `collector_system_native.go`, `chat.go`, `RunAgentSessionLoop`, `client.go`, `StorePendingConfirmation`, `TestSteeringQueue`, `collector_system.go`, `collector_security.go`, `wizard.go`, `startCLI`, `time.Time`, `ScorpPath`?**
  _High betweenness centrality (0.147) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `TruncateStr`, `collector_system_native.go`, `RegisterTool`, `chat.go`, `agent/autonomous.go`, `RunAgentSessionLoop`, `StorePendingConfirmation`, `client.go`, `collector_system.go`, `rag_vector.go`, `StartTestEndpoint`, `wizard.go`, `time.Time`, `skills.go`, `cost_router.go`, `ScorpPath`, `runSelfReview`?**
  _High betweenness centrality (0.100) - this node is a cross-community bridge._
- **Why does `GetStringArg()` connect `GetStringArg` to `TruncateStr`, `testing.T`, `RegisterTool`, `RunAgentSessionLoop`, `client.go`, `ExecuteTermuxAPI`, `runSubagent`, `uptime.go`, `time.Time`, `cost_router.go`, `ScorpPath`, `runSelfReview`, `TestPhase6_AllTools`?**
  _High betweenness centrality (0.094) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 21 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _26 weakly-connected nodes found - possible documentation gaps or missing edges._