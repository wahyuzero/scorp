package tools

import (
	"sync"
	"time"

	"scorp-agent/collectors"
	"scorp-agent/models"
)

// ──────────────────────────────────────────────
// Callbacks — set from main.go to avoid import cycles
// ──────────────────────────────────────────────

// Telegram callbacks
var (
	SendMessage       func(text string, keyboard map[string]interface{}) bool
	SendMessageGetID  func(text string, chatID int64) int64
	EditMessageByID   func(chatID int64, messageID int64, text string, keyboard map[string]interface{}) bool
	SendChatAction    func(chatID int64, action string)
	AnswerCallback    func(callbackID string, text string)
	TgPost            func(method string, payload map[string]interface{}) (TgResponse, error)
)

// File/bridge callbacks
var (
	SendDocumentBytes func(chatID string, data []byte, filename string, caption string) bool
)

// Confirmation/danger callbacks
var (
	StorePendingConfirmation func(chatID, toolName, command string, messages []AgentMessage)
	IsDangerousCommand       func(cmd string) bool
)

// Agent session/history callbacks
var (
	GetOrCreateSession   func(chatID string) *ChatSession
	GetSessionHistory    func(chatID string) []AgentMessage
	AppendSessionHistory func(chatID string, msgs ...AgentMessage)
	LoadHistoryFromDisk  func(chatID string) []AgentMessage
	SaveHistoryToDisk    func(chatID string, msgs []AgentMessage)
	ScheduleHistorySave  func(chatID string)
	EnterAgentMode       func(chatID string)
	ExitAgentMode        func(chatID string) bool
	IsUserActive         func() bool
	FlushPendingMessages func()
	SummarizeHistory     func(chatID string)
	ClearChatSession     func(chatID string)
	HistoryFilePath      func(chatID string) string
	MarkdownToTelegramHTML func(md string) string
	RunAgentLoop         func(chatID int64, userMessage string, msgID int64)
	ScorpChat            func(chatID, userMessage string) (string, error)
	ScorpChatMultiTurn   func(messages []AgentMessage) (string, error)
	SendScorpReply       func(chatID int64, msgID int64, reply string)
	SendMessageSmart     func(text string, keyboard map[string]interface{})
	SendFile             func(chatID, filePath string) bool
	MainMenuKeyboard     func() map[string]interface{}
)

// Autonomous callbacks
var (
	AutoConfig         *AutonomousConfig
	AutoMu             *sync.Mutex
	AutoLog            *[]AutonomousLogEntry
	AutoKillFile       string
	AutoCycleNum       *int
	RunAutonomousCycle func()
	SaveAutonomousConfig func()
	SetKillSwitch        func(bool)
)

// ExecuteTool callback for execute_code bridge
var ExecuteTool func(tc models.ToolCall, chatID int64) (string, bool)

// Collector callbacks
var (
	CollectSystem func() *collectors.SystemData
	CollectDocker func() *collectors.DockerData
)

// ──────────────────────────────────────────────
// Types needed by tools (mirrored from root)
// ──────────────────────────────────────────────

// tgResponse mirrors the struct in telegram.go
type TgResponse struct {
	OK          bool   `json:"ok"`
	Result      any    `json:"result"`
	Description string `json:"description"`
}

// AgentMessage mirrors the struct in models/tools.go
type AgentMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// pendingConfirmation mirrors the struct in agent_loop.go
type pendingConfirmation struct {
	ChatID    string
	ToolName  string
	Command   string
	Messages  []AgentMessage
	Timestamp int64
}

// ChatSession holds per-chat session state (minimal for callbacks)
type ChatSession struct {
	History     []AgentMessage
	AgentMode   bool
	LoopActive  bool
	PendingMsg  []AgentMessage
	LastActive  int64
	PendingConf *pendingConfirmation
}

// AutonomousConfig mirrors the struct in autonomous.go
type AutonomousConfig struct {
	Enabled        bool
	Interval       time.Duration
	ApprovalLevel  string
	MaxActions     int
	KillSwitch     bool
	LastCycle      time.Time
	TotalCycles    int
	TotalActions   int
}

// AutonomousLogEntry mirrors the struct in autonomous.go
type AutonomousLogEntry struct {
	Timestamp  time.Time
	Cycle      int
	Analysis   string
	Tool       string
	Args       map[string]interface{}
	Reason     string
	Risk       string
	Result     string
	Success    bool
	Approved   bool
}