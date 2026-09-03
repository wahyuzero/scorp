package agent

import (
	"sync"
)

// ──────────────────────────────────────────────
// Real-Time Steering Queue (PicoClaw Parity)
// Enables users to intercept and redirect the agent mid-run
// without killing the session or waiting for stale tool calls.
// ──────────────────────────────────────────────

var (
	steeringQueues   = make(map[string][]string)
	steeringQueuesMu sync.Mutex
)

// QueueSteeringMessage enqueues an interruption/redirection instruction for an active session
func QueueSteeringMessage(chatIDStr, message string) {
	if message == "" {
		return
	}
	steeringQueuesMu.Lock()
	defer steeringQueuesMu.Unlock()
	steeringQueues[chatIDStr] = append(steeringQueues[chatIDStr], message)
}

// PopSteeringMessage retrieves and removes the oldest steering message for a session
func PopSteeringMessage(chatIDStr string) (string, bool) {
	steeringQueuesMu.Lock()
	defer steeringQueuesMu.Unlock()

	q, ok := steeringQueues[chatIDStr]
	if !ok || len(q) == 0 {
		return "", false
	}

	msg := q[0]
	steeringQueues[chatIDStr] = q[1:]
	if len(steeringQueues[chatIDStr]) == 0 {
		delete(steeringQueues, chatIDStr)
	}
	return msg, true
}

// HasSteeringMessage checks if there are pending steering messages
func HasSteeringMessage(chatIDStr string) bool {
	steeringQueuesMu.Lock()
	defer steeringQueuesMu.Unlock()
	q, ok := steeringQueues[chatIDStr]
	return ok && len(q) > 0
}

// ClearSteeringQueue empties the steering queue for a session
func ClearSteeringQueue(chatIDStr string) {
	steeringQueuesMu.Lock()
	defer steeringQueuesMu.Unlock()
	delete(steeringQueues, chatIDStr)
}
