# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 162 files · ~114,742 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1348 nodes · 3310 edges · 49 communities (40 shown, 1 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 380 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `caef1bb2`
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
- agent/autonomous.go
- time.Time
- session_search_fts5.go
- ModelConfig
- cost_router.go
- HandleTelegramAction
- skills.go
- checker.go
- startCLI
- install.sh
- bg.go
- scorp-agent
- api_gemini.go
- 🦂 Scorp
- compaction_test.go
- RegisterTool
- api_commandcode.go
- RunAgentSessionLoop
- ToolCall
- ExecuteAutonomous
- api_anthropic.go
- GetStringArg
- GetProvider
- HomeDir
- inline.go
- tools/callbacks.go
- clarify.go
- EscapeHTML

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

## Communities (49 total, 1 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.19
Nodes (24): context.Context, TruncateStr(), AnthropicProvider, callAnthropic(), CallAnthropicWithTools(), CallCommandCode(), CallCommandCodeStream(), CallCommandCodeWithTools() (+16 more)

### Community 1 - "StartDaemon"
Cohesion: 0.11
Nodes (29): UploadsDir(), Init(), StartServer(), StopServer(), Load(), runCommandLoop(), StartDaemon(), BackAndRefreshKeyboard() (+21 more)

### Community 2 - "testing.T"
Cohesion: 0.09
Nodes (26): TestBase64Encode(), TestGetBoolArg(), TestGetFloatArg(), TestGetInt64Arg(), TestGetIntArg(), TestGetStringArg(), TestGetStringSliceArg(), TestMaxIterations() (+18 more)

### Community 3 - "chat.go"
Cohesion: 0.08
Nodes (56): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+48 more)

### Community 4 - "ScorpPath"
Cohesion: 0.05
Nodes (65): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+57 more)

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
Cohesion: 0.08
Nodes (40): ProjectDir(), checkACPAvailable(), launchACP(), listAvailableACP(), runOpenCodeCLI(), runSubagentACP(), ACPError, ACPInitializeParams (+32 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.09
Nodes (41): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+33 more)

### Community 12 - "wizard.go"
Cohesion: 0.12
Nodes (43): CatalogEntry, AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels(), SaveModelConfig() (+35 more)

### Community 13 - "agent/autonomous.go"
Cohesion: 0.14
Nodes (26): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+18 more)

### Community 14 - "time.Time"
Cohesion: 0.25
Nodes (18): time.Time, ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), FormatTasksList(), GetTask(), LoadTasks() (+10 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "ModelConfig"
Cohesion: 0.19
Nodes (21): ChatResponse, getCheapestModel(), RouteModelCostAware(), CallModel(), CallModelWithFallback(), CheckModelHealth(), defaultModelConfig(), findFirstVisionModel() (+13 more)

### Community 17 - "cost_router.go"
Cohesion: 0.09
Nodes (30): CM(), ConfigMgr(), InitConfigManager(), NewConfigManager(), ConfigManager, os.FileMode, defaultCostConfig(), formatCostReport() (+22 more)

### Community 18 - "HandleTelegramAction"
Cohesion: 0.23
Nodes (18): HandleTelegramAction(), BackKB(), createZip(), DirKeyboard(), FileDetailKeyboard(), FolderZipInfo(), GetPath(), HumanSize() (+10 more)

### Community 19 - "skills.go"
Cohesion: 0.07
Nodes (35): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+27 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.06
Nodes (58): ExecuteTool(), RegisterAutonomous(), executeOneShot(), executeTurn(), formatFinalResponse(), formatTerminalText(), handleCLISession(), handleCLISOP() (+50 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "bg.go"
Cohesion: 0.12
Nodes (20): bytes.Buffer, io.WriteCloser, os/exec.Cmd, sync.Mutex, CostTracker, bgKill(), bgList(), bgPoll() (+12 more)

### Community 31 - "api_gemini.go"
Cohesion: 0.23
Nodes (14): callGemini(), CallGeminiWithTools(), geminiBuildRequest(), geminiDoRequest(), geminiMessages(), geminiContent, geminiFuncDecl, geminiFunctionCall (+6 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "compaction_test.go"
Cohesion: 0.19
Nodes (22): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+14 more)

### Community 34 - "RegisterTool"
Cohesion: 0.06
Nodes (36): init(), init(), init(), init(), executeMCPServerTool(), unregisterMCPNativeTools(), TestGenerateNativeToolsSchema(), ActivateToolWithTTL() (+28 more)

### Community 35 - "api_commandcode.go"
Cohesion: 0.17
Nodes (13): buildCommandCodePayload(), extractFallbackToolCalls(), ChatRequest, commandCodeMsg, commandCodeParams, commandCodePayload, CommandCodeProvider, commandCodeTool (+5 more)

### Community 36 - "RunAgentSessionLoop"
Cohesion: 0.07
Nodes (42): AgentMessage, clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), AgentMessage, HandleConfirmation(), HasPendingConfirmation() (+34 more)

### Community 37 - "ToolCall"
Cohesion: 0.23
Nodes (9): TestParseToolCalls(), CustomProvider, ToolCall, IsRateLimitError(), ParseAllToolCalls(), ParseCodeBlockFallback(), ParseToolCalls(), TestCallModelWithToolsNilModel() (+1 more)

### Community 38 - "ExecuteAutonomous"
Cohesion: 0.60
Nodes (5): autoShowActions(), autoShowConfig(), autoShowLog(), autoStatus(), ExecuteAutonomous()

### Community 39 - "api_anthropic.go"
Cohesion: 0.67
Nodes (3): anthropicRequest, anthropicResponse, anthropicTool

### Community 41 - "GetStringArg"
Cohesion: 0.06
Nodes (51): StorePendingConfirmation(), init(), init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage() (+43 more)

### Community 43 - "GetProvider"
Cohesion: 0.36
Nodes (6): TestCallOpenAI_MockServer(), TestLLMProviderFactory(), LLMProvider, GetProvider(), init(), RegisterProviderAdapter()

### Community 44 - "HomeDir"
Cohesion: 0.18
Nodes (16): HomeDir(), PythonSitePackages(), LoadModelConfig(), applyProviderDefaults(), ProviderPreset, migrateModelConfigs(), ResolveBaseURL(), resolveCommandCodeKeyFromDisk() (+8 more)

### Community 45 - "inline.go"
Cohesion: 0.38
Nodes (9): AnswerInlineQuery(), buildInlineResults(), firstN(), TGInlineQuery, HandleInlineQuery(), quickDocker(), quickStatus(), quickStorage() (+1 more)

### Community 47 - "tools/callbacks.go"
Cohesion: 0.53
Nodes (5): AgentMessage, AutonomousLogEntry, ChatSession, pendingConfirmation, TgResponse

### Community 48 - "clarify.go"
Cohesion: 0.24
Nodes (10): AnswerCallback(), executeClarify(), GetClarifyChatID(), handleClarifyResponse(), HasPendingClarify(), init(), ResolveClarify(), sendClarifyMessage() (+2 more)

### Community 51 - "EscapeHTML"
Cohesion: 0.36
Nodes (7): EscapeHTML(), isLikelyScriptPath(), notifyTaskResult(), runScriptTask(), runShellTaskConfig(), splitMessage(), ShellTask()

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 133 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `StartDaemon`, `chat.go`, `RunAgentSessionLoop`, `client.go`, `collector_system.go`, `collector_security.go`, `collector_system_native.go`, `wizard.go`, `time.Time`, `ModelConfig`, `clarify.go`, `startCLI`?**
  _High betweenness centrality (0.135) - this node is a cross-community bridge._
- **Why does `StartDaemon()` connect `StartDaemon` to `chat.go`, `RunAgentSessionLoop`, `client.go`, `rag_vector.go`, `collector_system.go`, `GetStringArg`, `collector_system_native.go`, `HomeDir`, `agent/autonomous.go`, `time.Time`, `wizard.go`, `cost_router.go`, `HandleTelegramAction`, `skills.go`, `startCLI`?**
  _High betweenness centrality (0.090) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `RegisterTool`, `RunAgentSessionLoop`, `client.go`, `rag_vector.go`, `runSubagent`, `HomeDir`, `skills.go`, `startCLI`, `bg.go`?**
  _High betweenness centrality (0.083) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `StartDaemon` be split into smaller, more focused modules?**
  _Cohesion score 0.10795454545454546 - nodes in this community are weakly interconnected._