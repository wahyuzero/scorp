# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 159 files · ~109,724 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1280 nodes · 3171 edges · 43 communities (34 shown, 1 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 376 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4e8b17ca`
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
- ConfigMgr
- ScorpPath
- RunAgentSessionLoop
- checker.go
- GetProvider
- install.sh
- bg.go
- scorp-agent
- HomeDir
- 🦂 Scorp
- collector_system_native.go
- RegisterTool
- collector_native_test.go
- EscapeHTML
- clarify.go
- runScriptTask
- LoadConfig
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
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go
- `init()` --calls--> `RagVecProvider()`  [EXTRACTED]
  bootstrap/extended.go → rag/rag_vector.go
- `browserSessionScreenshot()` --calls--> `SendFile()`  [INFERRED]
  browser/browser_session.go → telegram/files.go

## Import Cycles
- None detected.

## Communities (43 total, 1 thin omitted)

### Community 0 - "TruncateStr"
Cohesion: 0.05
Nodes (92): context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic(), CallAnthropicWithTools() (+84 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.12
Nodes (32): Init(), StartServer(), StopServer(), HandleTelegramAction(), runCommandLoop(), StartDaemon(), AnswerCallback(), BackAndRefreshKeyboard() (+24 more)

### Community 2 - "testing.T"
Cohesion: 0.10
Nodes (26): TestBase64Encode(), TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg() (+18 more)

### Community 3 - "chat.go"
Cohesion: 0.08
Nodes (57): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+49 more)

### Community 4 - "agent/autonomous.go"
Cohesion: 0.12
Nodes (31): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+23 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (47): bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, buildArgDefsFromInputSchema(), executeMCPServerTool(), FindMCPTool() (+39 more)

### Community 6 - "compaction_test.go"
Cohesion: 0.19
Nodes (22): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+14 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.08
Nodes (51): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+43 more)

### Community 9 - "collector_security.go"
Cohesion: 0.18
Nodes (25): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+17 more)

### Community 10 - "runSubagent"
Cohesion: 0.08
Nodes (37): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+29 more)

### Community 11 - "uptime.go"
Cohesion: 0.11
Nodes (27): FormatDuration(), getUptime(), net/http.Client, net/http.Transport, sync.RWMutex, time.Duration, TransportPool, extractHost() (+19 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (43): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), SaveModelConfig() (+35 more)

### Community 13 - "startCLI"
Cohesion: 0.06
Nodes (62): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+54 more)

### Community 14 - "time.Time"
Cohesion: 0.21
Nodes (20): time.Time, ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), FormatTasksList(), GetTask() (+12 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "skills.go"
Cohesion: 0.24
Nodes (13): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+5 more)

### Community 17 - "ConfigMgr"
Cohesion: 0.08
Nodes (31): CM(), ConfigMgr(), InitConfigManager(), NewConfigManager(), ConfigManager, os.FileMode, GetFloatArg(), defaultCostConfig() (+23 more)

### Community 18 - "ScorpPath"
Cohesion: 0.06
Nodes (59): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+51 more)

### Community 19 - "RunAgentSessionLoop"
Cohesion: 0.05
Nodes (54): AgentMessage, getSharedMemorySummary(), countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), looksLikeContinuation(), mentionsBrowserTask() (+46 more)

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

### Community 31 - "HomeDir"
Cohesion: 0.20
Nodes (18): HomeDir(), ProjectDir(), PythonSitePackages(), UploadsDir(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard() (+10 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.11
Nodes (17): 1. Diet Ekstrem `main.go`, 2. Modularisasi `agent/loop.go`, 3. Pemisahan Provider AI Independen (`models/`), 🚀 Panduan Perintah Cepat (Quick Cheatsheet), 🧱 Rekapitulasi Refactoring Arsitektur Bersih (*Clean Architecture*), 📌 Ringkasan Pembaruan Utama (v2.0), 🦂 Scorp Agent v2.0 — Modernization & Architecture Upgrade Report, 🧪 Verifikasi & Kualitas Kode (+9 more)

### Community 33 - "collector_system_native.go"
Cohesion: 0.23
Nodes (16): TestCollectorNative_StartCPUSampler_DoesNotBlock(), CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg(), getMemInfo(), getNetBytes(), getProcesses() (+8 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (35): ExecuteTool(), init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), TestGenerateNativeToolsSchema() (+27 more)

### Community 35 - "collector_native_test.go"
Cohesion: 0.21
Nodes (12): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+4 more)

### Community 36 - "EscapeHTML"
Cohesion: 0.30
Nodes (11): EscapeHTML(), ShellTask(), AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker() (+3 more)

### Community 37 - "clarify.go"
Cohesion: 0.31
Nodes (8): executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage(), SetClarifyChatID()

### Community 38 - "runScriptTask"
Cohesion: 0.48
Nodes (5): isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage()

### Community 39 - "LoadConfig"
Cohesion: 0.67
Nodes (6): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig()

### Community 41 - "GetStringArg"
Cohesion: 0.05
Nodes (60): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetInt64Arg(), GetIntArg(), GetStringArg(), GetStringSliceArg() (+52 more)

## Knowledge Gaps
- **26 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+21 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 117 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `collector_system_native.go`, `chat.go`, `client.go`, `clarify.go`, `collector_system.go`, `collector_security.go`, `wizard.go`, `startCLI`, `time.Time`, `RunAgentSessionLoop`, `HomeDir`?**
  _High betweenness centrality (0.150) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `TruncateStr`, `collector_system_native.go`, `chat.go`, `agent/autonomous.go`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `wizard.go`, `startCLI`, `time.Time`, `skills.go`, `ConfigMgr`, `RunAgentSessionLoop`, `HomeDir`?**
  _High betweenness centrality (0.103) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `RegisterTool`, `client.go`, `rag_vector.go`, `runSubagent`, `startCLI`, `skills.go`, `RunAgentSessionLoop`, `bg.go`, `HomeDir`?**
  _High betweenness centrality (0.094) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 21 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _26 weakly-connected nodes found - possible documentation gaps or missing edges._