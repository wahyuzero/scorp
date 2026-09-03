package agent

import (
	"testing"
)

func TestSteeringQueue(t *testing.T) {
	chatID := "test_steer_123"
	ClearSteeringQueue(chatID)

	if HasSteeringMessage(chatID) {
		t.Errorf("Expected empty steering queue initially")
	}

	QueueSteeringMessage(chatID, "Stop and search for file Y")
	QueueSteeringMessage(chatID, "Second instruction")

	if !HasSteeringMessage(chatID) {
		t.Errorf("Expected queue to have steering messages")
	}

	msg1, ok1 := PopSteeringMessage(chatID)
	if !ok1 || msg1 != "Stop and search for file Y" {
		t.Errorf("Expected first message 'Stop and search for file Y', got '%s'", msg1)
	}

	msg2, ok2 := PopSteeringMessage(chatID)
	if !ok2 || msg2 != "Second instruction" {
		t.Errorf("Expected second message 'Second instruction', got '%s'", msg2)
	}

	if HasSteeringMessage(chatID) {
		t.Errorf("Expected queue to be empty after popping all messages")
	}
}
