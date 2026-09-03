# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 159 files · ~109,386 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1279 nodes · 3168 edges · 37 communities (28 shown, 1 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 376 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9d3d413a`
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
- collector_system_native.go
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
- ExecuteAutonomous
- 🦂 Scorp
- RegisterTool
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
- `autoShowConfig()` --calls--> `SaveAutonomousConfig()`  [INFERRED]
  tools/autonomous.go → agent/autonomous.go
- `ExecuteAutonomous()` --calls--> `SaveAutonomousConfig()`  [INFERRED]
  tools/autonomous.go → agent/autonomous.go
- `ExecuteAutonomous()` --calls--> `SetKillSwitch()`  [INFERRED]
  tools/autonomous.go → agent/autonomous.go
- `ExecuteAutonomous()` --calls--> `RunAutonomousCycle()`  [INFERRED]
  tools/autonomous.go → agent/autonomous.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go

## Import Cycles
- None detected.

## Communities (37 total, 1 thin omitted)

### Community 0 - "TruncateStr"
Cohesion: 0.05
Nodes (91): context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic(), CallAnthropicWithTools() (+83 more)

### Community 1 - "HandleTelegramAction"
Cohesion: 0.05
Nodes (73): Config, EnvBool(), EnvFloat(), EnvInt(), EnvStr(), LoadConfig(), HomeDir(), PythonSitePackages() (+65 more)

### Community 2 - "testing.T"
Cohesion: 0.10
Nodes (26): TestBase64Encode(), TestBuildThinkingMessage(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg() (+18 more)

### Community 3 - "chat.go"
Cohesion: 0.08
Nodes (57): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+49 more)

### Community 4 - "agent/autonomous.go"
Cohesion: 0.14
Nodes (26): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+18 more)

### Community 5 - "client.go"
Cohesion: 0.07
Nodes (48): bufio.Scanner, context.CancelFunc, encoding/json.Encoder, encoding/json.RawMessage, io.ReadCloser, net/url.URL, buildArgDefsFromInputSchema(), executeMCPServerTool() (+40 more)

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
Nodes (37): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+29 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.06
Nodes (55): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+47 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (44): ProjectDir(), CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "startCLI"
Cohesion: 0.05
Nodes (71): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation (+63 more)

### Community 14 - "time.Time"
Cohesion: 0.16
Nodes (26): time.Time, EscapeHTML(), ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath() (+18 more)

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
Cohesion: 0.05
Nodes (58): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+50 more)

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
Cohesion: 0.25
Nodes (15): bytes.Buffer, io.WriteCloser, os/exec.Cmd, bgKill(), bgList(), bgPoll(), bgSpawn(), bgWait() (+7 more)

### Community 31 - "ExecuteAutonomous"
Cohesion: 0.60
Nodes (5): autoShowActions(), autoShowConfig(), autoShowLog(), autoStatus(), ExecuteAutonomous()

### Community 32 - "🦂 Scorp"
Cohesion: 0.11
Nodes (17): 1. Diet Ekstrem `main.go`, 2. Modularisasi `agent/loop.go`, 3. Pemisahan Provider AI Independen (`models/`), 🚀 Panduan Perintah Cepat (Quick Cheatsheet), 🧱 Rekapitulasi Refactoring Arsitektur Bersih (*Clean Architecture*), 📌 Ringkasan Pembaruan Utama (v2.0), 🦂 Scorp Agent v2.0 — Modernization & Architecture Upgrade Report, 🧪 Verifikasi & Kualitas Kode (+9 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (32): init(), init(), init(), init(), TestParseDelegateParams(), TestValidateSubagentTools(), TestGenerateNativeToolsSchema(), ActivateToolWithTTL() (+24 more)

### Community 41 - "GetStringArg"
Cohesion: 0.05
Nodes (59): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetInt64Arg(), GetIntArg(), GetStringArg(), GetStringSliceArg() (+51 more)

## Knowledge Gaps
- **26 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+21 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 117 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `chat.go`, `client.go`, `collector_system.go`, `collector_security.go`, `collector_system_native.go`, `wizard.go`, `startCLI`, `time.Time`, `RunAgentSessionLoop`?**
  _High betweenness centrality (0.145) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `HandleTelegramAction` to `TruncateStr`, `chat.go`, `agent/autonomous.go`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `collector_system_native.go`, `wizard.go`, `startCLI`, `time.Time`, `skills.go`, `ConfigMgr`, `RunAgentSessionLoop`?**
  _High betweenness centrality (0.100) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `HandleTelegramAction`, `RegisterTool`, `client.go`, `rag_vector.go`, `runSubagent`, `startCLI`, `skills.go`, `RunAgentSessionLoop`, `bg.go`?**
  _High betweenness centrality (0.098) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 21 inferred relationships involving `RunAgentSessionLoop()` (e.g. with `appendSessionHistory()` and `getSessionHistory()`) actually correct?**
  _`RunAgentSessionLoop()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _26 weakly-connected nodes found - possible documentation gaps or missing edges._