package scheduler

// Callbacks for functions that live in the main (root) package.
// These are set during init in the root package to avoid import cycle.
var (
	// SendMessage sends a message to Telegram
	SendMessage func(text string, keyboard map[string]interface{})

	// SendMessageGetID sends a message and returns the message ID
	SendMessageGetID func(text string, chatID int64) int64

	// TgPost calls the Telegram Bot API (no return — scheduler ignores result)
	TgPost func(method string, payload map[string]interface{})

	// RunAgentLoop runs the full agent loop
	RunAgentLoop func(chatID int64, userMessage string, msgID int64)
)
