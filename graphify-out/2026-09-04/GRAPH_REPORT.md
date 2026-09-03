# Graph Report - scorp  (2026-09-04)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 1120 nodes · 2859 edges · 31 communities (23 shown, 2 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 266 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `78253b41`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- TruncateStr
- handleAction
- GetStringArg
- chat.go
- config_paths.go
- client.go
- testing.T
- rag_vector.go
- collector_system.go
- ConfigMgr
- runSubagent
- collector_system_native.go
- wizard.go
- startCLI
- time.Time
- session_search_fts5.go
- collector_security.go
- bg.go
- runSelfReview
- skills.go
- checker.go
- collector_coolify.go
- install.sh
- LoopController
- scorp-agent

## God Nodes (most connected - your core abstractions)
1. `handleAction()` - 59 edges
2. `main()` - 54 edges
3. `GetStringArg()` - 40 edges
4. `TruncateStr()` - 36 edges
5. `RunAgentLoop()` - 35 edges
6. `init()` - 30 edges
7. `HandleModelCallback()` - 29 edges
8. `ModelConfig` - 28 edges
9. `startCLI()` - 28 edges
10. `resumeAgentLoop()` - 23 edges

## Surprising Connections (you probably didn't know these)
- `handleAction()` --calls--> `MCPToolsSummary()`  [EXTRACTED]
  main.go → mcp/client.go
- `main()` --calls--> `hasDebugFlag()`  [INFERRED]
  main.go → cli.go
- `main()` --calls--> `startCLI()`  [INFERRED]
  main.go → cli.go
- `main()` --calls--> `StopMCPServerMode()`  [EXTRACTED]
  main.go → mcp/client.go
- `browserSessionScreenshot()` --calls--> `SendFile()`  [INFERRED]
  browser/browser_session.go → telegram/files.go

## Import Cycles
- None detected.

## Communities (31 total, 2 thin omitted)

### Community 0 - "TruncateStr"
Cohesion: 0.07
Nodes (85): context.Context, TruncateStr(), anthropicRequest, anthropicResponse, anthropicTool, callAnthropic(), CallAnthropicWithTools(), buildCommandCodePayload() (+77 more)

### Community 1 - "handleAction"
Cohesion: 0.06
Nodes (73): IsUserActive(), net/http.Request, net/http.ResponseWriter, commandLoop(), formatBasicAlert(), handleAction(), hourlyLoop(), isCLIMode() (+65 more)

### Community 2 - "GetStringArg"
Cohesion: 0.06
Nodes (60): init(), init(), ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill() (+52 more)

### Community 3 - "chat.go"
Cohesion: 0.06
Nodes (73): AgentMessage, agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines() (+65 more)

### Community 4 - "config_paths.go"
Cohesion: 0.05
Nodes (64): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+56 more)

### Community 5 - "client.go"
Cohesion: 0.05
Nodes (60): RegisterAutonomous(), init(), init(), init(), init(), bufio.Scanner, encoding/json.Encoder, encoding/json.RawMessage (+52 more)

### Community 6 - "testing.T"
Cohesion: 0.06
Nodes (60): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+52 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (59): CollectDocker(), CollectDockerFallback(), containerName(), dockerGet(), DockerData, InitDockerClient(), StartDockerStatsSampler(), checkGDriveMount() (+51 more)

### Community 9 - "ConfigMgr"
Cohesion: 0.08
Nodes (37): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+29 more)

### Community 10 - "runSubagent"
Cohesion: 0.08
Nodes (38): checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams, ACPMessageNewParams (+30 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.07
Nodes (46): TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), FormatDuration(), TopProcess, CollectSystem(), getCPUCount(), getDiskUsage(), getLoadAvg() (+38 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (43): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), SaveModelConfig() (+35 more)

### Community 13 - "startCLI"
Cohesion: 0.10
Nodes (37): HasPendingConfirmation(), executeOneShot(), executeTurn(), formatFinalResponse(), formatTerminalText(), hasDebugFlag(), isTerminal(), printBanner() (+29 more)

### Community 14 - "time.Time"
Cohesion: 0.15
Nodes (26): time.Time, EscapeHTML(), SplitMessage(), ModelUsage, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule() (+18 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 17 - "bg.go"
Cohesion: 0.17
Nodes (19): bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, CostTracker, getChatLock(), StartTestEndpoint(), bgKill() (+11 more)

### Community 18 - "runSelfReview"
Cohesion: 0.17
Nodes (16): saveToMemory(), memoryFact, AgentMessage, runSelfReview(), TestSelfReviewIntegration(), LoadJSON(), SaveJSON(), SaveJSONPerm() (+8 more)

### Community 19 - "skills.go"
Cohesion: 0.22
Nodes (14): ExecuteSkillManage(), executeSkillManageCreate(), executeSkillManageDelete(), executeSkillManageList(), executeSkillManageUpdate(), ExecuteSkillManageView(), Skill, Delete() (+6 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "collector_coolify.go"
Cohesion: 0.32
Nodes (11): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CoolifyApp (+3 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

## Knowledge Gaps
- **12 isolated node(s):** `textContentPart`, `toolCallContentPart`, `ACPInitializeParams`, `AgentMessage`, `TgResponse` (+7 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 75 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `handleAction` to `TruncateStr`, `chat.go`, `config_paths.go`, `client.go`, `rag_vector.go`, `collector_system.go`, `ConfigMgr`, `collector_system_native.go`, `wizard.go`, `startCLI`, `time.Time`, `collector_security.go`, `bg.go`, `runSelfReview`, `skills.go`, `checker.go`?**
  _High betweenness centrality (0.198) - this node is a cross-community bridge._
- **Why does `handleAction()` connect `handleAction` to `chat.go`, `config_paths.go`, `client.go`, `collector_system.go`, `collector_system_native.go`, `wizard.go`, `startCLI`, `time.Time`, `collector_security.go`, `skills.go`, `checker.go`, `collector_coolify.go`?**
  _High betweenness centrality (0.163) - this node is a cross-community bridge._
- **Why does `GetStringArg()` connect `GetStringArg` to `TruncateStr`, `chat.go`, `config_paths.go`, `client.go`, `testing.T`, `runSubagent`, `collector_system_native.go`, `startCLI`, `time.Time`, `runSelfReview`?**
  _High betweenness centrality (0.111) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `main()` (e.g. with `hasDebugFlag()` and `startCLI()`) actually correct?**
  _`main()` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `textContentPart`, `toolCallContentPart`, `ACPInitializeParams` to the rest of the system?**
  _12 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `TruncateStr` be split into smaller, more focused modules?**
  _Cohesion score 0.06808510638297872 - nodes in this community are weakly interconnected._
- **Should `handleAction` be split into smaller, more focused modules?**
  _Cohesion score 0.05727848101265823 - nodes in this community are weakly interconnected._