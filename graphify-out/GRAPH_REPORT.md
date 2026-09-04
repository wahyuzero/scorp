# Graph Report - scorp  (2026-09-04)

## Corpus Check
- 175 files · ~117,858 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1382 nodes · 3382 edges · 49 communities (38 shown, 3 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 441 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1cb7d18d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- context.Context
- init
- testing.T
- chat.go
- ScorpPath
- client.go
- execute.go
- rag_vector.go
- collector_system.go
- collector_security.go
- acp.go
- collector_system_native.go
- HandleModelCallback
- cost_router.go
- time.Time
- session_search_fts5.go
- TruncOutput
- HandleTelegramAction
- GetAllTools
- MCPServer
- checker.go
- startCLI
- install.sh
- LoadMCPConfig
- scorp-agent
- runSubagent
- 🦂 Scorp
- skills.go
- RegisterTool
- metasearch_engines.go
- patch.go
- renderStatusFooter
- StorePendingConfirmation
- .listenSSEStream
- read_url.go
- GetStringArg
- serviceBridgeRequests
- RunAgentSessionLoop
- echoPlugin
- RegisterPlugin
- ExecuteAutoLogin

## God Nodes (most connected - your core abstractions)
1. `HandleTelegramAction()` - 54 edges
2. `GetStringArg()` - 43 edges
3. `StartDaemon()` - 40 edges
4. `RunAgentSessionLoop()` - 38 edges
5. `ModelConfig` - 36 edges
6. `startCLI()` - 35 edges
7. `TruncateStr()` - 35 edges
8. `init()` - 32 edges
9. `ChatMessage` - 30 edges
10. `HandleModelCallback()` - 29 edges

## Surprising Connections (you probably didn't know these)
- `ExecuteShell()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/exec.go → agent/confirmation.go
- `ExecuteGit()` --calls--> `StorePendingConfirmation()`  [INFERRED]
  tools/git.go → agent/confirmation.go
- `runAgentTask()` --calls--> `RunAgentLoop()`  [INFERRED]
  scheduler/scheduler.go → agent/loop.go
- `serviceBridgeRequests()` --calls--> `ExecuteTool()`  [INFERRED]
  tools/exec_code.go → agent/prompt.go
- `FormatToolResult()` --references--> `ToolCall`  [EXTRACTED]
  agent/prompt.go → models/callback.go

## Import Cycles
- None detected.

## Communities (49 total, 3 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.05
Nodes (95): TestParseToolCalls(), context.Context, TruncateStr(), AnthropicProvider, anthropicRequest, anthropicResponse, anthropicTool, callAnthropic() (+87 more)

### Community 1 - "init"
Cohesion: 0.22
Nodes (8): init(), SendDocumentBytes(), ExecuteListDir(), ExecuteSendFile(), ExecuteSystemInfo(), ExecuteWriteFile(), isPathAllowed(), ExecuteSearchCode()

### Community 2 - "testing.T"
Cohesion: 0.06
Nodes (56): estimateHistoryTokens(), estimateTokens(), formatTokenEstimate(), AgentMessage, AgentMessage, makeHistory(), makeToolResult(), TestEstimateHistoryTokens() (+48 more)

### Community 3 - "chat.go"
Cohesion: 0.07
Nodes (59): agentSession, agentAutoStop(), appendSessionHistory(), cleanupChatSessions(), CleanupSessionsLoop(), ClearChatSession(), collectTableLines(), convertInlineMarkdown() (+51 more)

### Community 4 - "ScorpPath"
Cohesion: 0.05
Nodes (56): contains(), containsStr(), jsonToMap(), TestPhase6_AllTools(), TestPhase6_ScriptResult(), TestPhase6_VaultEncryption(), truncate(), init() (+48 more)

### Community 5 - "client.go"
Cohesion: 0.17
Nodes (21): ACPRequest, encoding/json.RawMessage, executeMCPServerTool(), FindMCPTool(), getExposedTools(), GetMCPTools(), handleMCPRequest(), MCPToolsForPrompt() (+13 more)

### Community 6 - "execute.go"
Cohesion: 0.26
Nodes (16): runOpenCodeCLI(), runSubagentACP(), AgentMessage, delegateBatchParams, delegateResult, delegateTaskParams, DefaultSubagentTools(), ExecuteDelegate() (+8 more)

### Community 7 - "rag_vector.go"
Cohesion: 0.06
Nodes (45): activeToolCall, getBoolArg(), getFloatArg(), getIntArg(), getStringArg(), hybridResult, homeDir(), ragDirPath() (+37 more)

### Community 8 - "collector_system.go"
Cohesion: 0.07
Nodes (52): cleanAppName(), CollectCoolify(), coolifyGet(), CoolifyData, jsonBool(), jsonStr(), parseStatus(), CollectDocker() (+44 more)

### Community 9 - "collector_security.go"
Cohesion: 0.19
Nodes (24): BruteForceAlert, CheckBruteForce(), CheckVNCConnections(), CollectSecurity(), CollectSecurityWithPeek(), DrainFailedSSHBuffer(), enrichLast10(), extractTime() (+16 more)

### Community 10 - "acp.go"
Cohesion: 0.08
Nodes (34): checkACPAvailable(), launchACP(), listAvailableACP(), ACPError, ACPInitializeParams, ACPMessageNewParams, ACPMessagePart, ACPResponse (+26 more)

### Community 11 - "collector_system_native.go"
Cohesion: 0.09
Nodes (41): TestCollectorNative_CollectSystem_Structure(), TestCollectorNative_GetTopProcesses(), TestCollectorNative_GetTopProcesses_Limit(), TestCollectorNative_NativeTopProcessStruct(), TestCollectorNative_SortByCPUDesc(), TestCollectorNative_SortByCPUDesc_Empty(), TestCollectorNative_SortByCPUDesc_EqualValues(), TestCollectorNative_SortByCPUDesc_Single() (+33 more)

### Community 12 - "HandleModelCallback"
Cohesion: 0.11
Nodes (44): CatalogEntry, SaveModelConfig(), AutoPopulateFromCatalog(), CatalogModels(), HasCatalog(), ProviderHasAPIKey(), ProviderKeyEnv(), RemoveProviderModels() (+36 more)

### Community 13 - "cost_router.go"
Cohesion: 0.05
Nodes (62): AppendAutoLog(), AutonomousLoop(), CheckKillSwitch(), executeAutonomousAction(), extractJSON(), gatherContext(), LoadAutoLog(), LoadAutonomousConfig() (+54 more)

### Community 14 - "time.Time"
Cohesion: 0.13
Nodes (31): time.Time, EscapeHTML(), ScheduledTask, AddTask(), AddTaskEx(), ExecuteSchedule(), isLikelyScriptPath(), notifyTaskResult() (+23 more)

### Community 15 - "session_search_fts5.go"
Cohesion: 0.10
Nodes (18): getIntArg(), getStringArg(), truncateString(), homeDir(), scorpDir(), scorpPath(), ExecuteSessionSearch(), SessionResult (+10 more)

### Community 16 - "TruncOutput"
Cohesion: 0.34
Nodes (14): ExecuteBrowser(), browserConsole(), browserSessionClick(), browserSessionEvaluate(), browserSessionExtract(), browserSessionFill(), browserSessionNavigate(), browserSessionScreenshot() (+6 more)

### Community 17 - "HandleTelegramAction"
Cohesion: 0.06
Nodes (70): StopMCPServerMode(), Init(), StartServer(), StopServer(), FormatUsageStats(), InitModelUsage(), Load(), HandleTelegramAction() (+62 more)

### Community 18 - "GetAllTools"
Cohesion: 0.22
Nodes (13): ActivateToolWithTTL(), IsDynamicModeEnabled(), IsToolActive(), ResetDynamicTools(), TestDynamicToolTTL(), TickToolTTL(), GetAllTools(), ResetNativeToolCache() (+5 more)

### Community 19 - "MCPServer"
Cohesion: 0.21
Nodes (9): bufio.Scanner, encoding/json.Encoder, MCPServer, startMCPServer(), MCPServerConfig, isRemoteMCP(), startSSEServer(), TestRemoteMCPServer_HTTP() (+1 more)

### Community 20 - "checker.go"
Cohesion: 0.25
Nodes (15): Asset, CheckForUpdate(), DownloadAsset(), FetchLatestRelease(), FindAssetForArch(), getRepo(), IsNewer(), isTermux() (+7 more)

### Community 21 - "startCLI"
Cohesion: 0.05
Nodes (65): clearPendingConfirmation(), confirmKeyboard(), getPendingConfirmation(), GetPendingConfirmationDetails(), HandleConfirmation(), HasPendingConfirmation(), pendingConfirmation, ExecuteTool() (+57 more)

### Community 23 - "install.sh"
Cohesion: 0.60
Nodes (5): ask(), die(), ok(), install.sh script, warn()

### Community 24 - "LoadMCPConfig"
Cohesion: 0.23
Nodes (14): TruncOutputTool(), buildArgDefsFromInputSchema(), LoadMCPConfig(), rebuildMCPToolList(), registerMCPToolsAsNative(), ReloadMCPServers(), sanitizeMCPName(), StartMCPServers() (+6 more)

### Community 31 - "runSubagent"
Cohesion: 0.24
Nodes (10): cleanupSubagentSandbox(), createSubagentSandbox(), defaultIsolation(), formatIsolationInfo(), getSubagentIsolation(), isSubagentToolBlocked(), registerSubagentIsolation(), unregisterSubagentIsolation() (+2 more)

### Community 32 - "🦂 Scorp"
Cohesion: 0.07
Nodes (25): 📊 1. Ringkasan Matriks Perbandingan, 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga", 🎯 3. Posisi & Rekomendasi Pilihan, 🚀 4. Status Implementasi Fitur Baru Scorp (Completed), 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw, 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback), 🦂 Scorp (The Native Embedded Metasearch Winner), 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch) (+17 more)

### Community 33 - "skills.go"
Cohesion: 0.07
Nodes (36): getSharedMemorySummary(), memoryFact, FormatToolResult(), getAgentSystemPrompt(), GetRepoMap(), InvalidateRepoMap(), TestGetRepoMap(), AgentMessage (+28 more)

### Community 34 - "RegisterTool"
Cohesion: 0.13
Nodes (15): init(), init(), init(), init(), unregisterMCPNativeTools(), TestRegisterPlugin(), ExecuteToolByName(), GenerateSystemPromptDescriptions() (+7 more)

### Community 35 - "metasearch_engines.go"
Cohesion: 0.06
Nodes (40): net/http.Client, net/http.Transport, sync.RWMutex, TransportPool, extractHost(), getClient(), getShortClient(), getTransport() (+32 more)

### Community 36 - "patch.go"
Cohesion: 0.35
Nodes (10): buildDiffPreview(), ExecutePatch(), ExecuteReplaceFileContent(), lineWindowMatch(), normalizeForCompare(), patchReplace(), scopedLineReplace(), splitLines() (+2 more)

### Community 37 - "renderStatusFooter"
Cohesion: 0.31
Nodes (9): GetDailyTotalUSD(), SlashCommand, filterCommands(), readInteractiveInput(), renderPopupBox(), getContextPill(), getGitStatus(), getShortCwd() (+1 more)

### Community 38 - "StorePendingConfirmation"
Cohesion: 0.32
Nodes (6): AgentMessage, StorePendingConfirmation(), ExecuteSQL(), loadDBConnections(), dbConnection, ExecuteProcess()

### Community 39 - ".listenSSEStream"
Cohesion: 0.33
Nodes (3): io.ReadCloser, net/url.URL, MCPServer

### Community 40 - "read_url.go"
Cohesion: 0.60
Nodes (5): ReadURL(), scrapeFirecrawl(), scrapeTavily(), truncateURLOutput(), tryRemoteScrape()

### Community 41 - "GetStringArg"
Cohesion: 0.14
Nodes (17): init(), GetBoolArg(), GetIntArg(), GetStringArg(), getUnclosedTags(), SplitMessage(), ExecuteCompose(), ExecuteReadFile() (+9 more)

### Community 42 - "serviceBridgeRequests"
Cohesion: 0.70
Nodes (4): executeCodeTool(), init(), serviceBridgeRequests(), writeBridgeResponse()

### Community 43 - "RunAgentSessionLoop"
Cohesion: 0.08
Nodes (35): AgentMessage, countStepsInMessage(), AgentMessage, hasCompletionIndicators(), hasForwardIntent(), isPureInformationalQuery(), looksLikeContinuation(), mentionsBrowserTask() (+27 more)

### Community 45 - "RegisterPlugin"
Cohesion: 0.83
Nodes (3): RegisterPlugin(), ToolPlugin, ToolPluginWithSchema

## Knowledge Gaps
- **32 isolated node(s):** `memoryFact`, `containerStats`, `ACPInitializeParams`, `AgentMessage`, `scorp-agent` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 135 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `HandleTelegramAction()` connect `HandleTelegramAction` to `chat.go`, `collector_system.go`, `collector_security.go`, `RunAgentSessionLoop`, `collector_system_native.go`, `HandleModelCallback`, `time.Time`, `startCLI`, `LoadMCPConfig`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Why does `TruncateStr()` connect `context.Context` to `skills.go`, `chat.go`, `execute.go`, `GetStringArg`, `RunAgentSessionLoop`, `cost_router.go`, `time.Time`, `TruncOutput`, `MCPServer`, `LoadMCPConfig`?**
  _High betweenness centrality (0.086) - this node is a cross-community bridge._
- **Why does `init()` connect `GetStringArg` to `skills.go`, `RegisterTool`, `ScorpPath`, `patch.go`, `execute.go`, `rag_vector.go`, `acp.go`, `RunAgentSessionLoop`, `GetAllTools`, `startCLI`, `LoadMCPConfig`?**
  _High betweenness centrality (0.080) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `HandleTelegramAction()` (e.g. with `DirKeyboard()` and `FileDetailKeyboard()`) actually correct?**
  _`HandleTelegramAction()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `StartDaemon()` (e.g. with `SendFile()` and `EditMessageByID()`) actually correct?**
  _`StartDaemon()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **What connects `memoryFact`, `containerStats`, `ACPInitializeParams` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.05098988748041589 - nodes in this community are weakly interconnected._