package main

import (
	"testing"
)

func TestAcquireSessionLock(t *testing.T) {
	testSess := "test_lock_session_123"

	// First acquisition must succeed
	lock1, err := acquireSessionLock(testSess)
	if err != nil {
		t.Fatalf("First acquireSessionLock failed: %v", err)
	}
	if lock1 == nil {
		t.Fatalf("Expected non-nil lock")
	}
	defer lock1.Release()

	// Second acquisition from same/another caller on same session must fail
	lock2, err := acquireSessionLock(testSess)
	if err == nil {
		lock2.Release()
		t.Fatalf("Expected second acquireSessionLock to fail, but succeeded")
	}

	// Different session must succeed concurrently
	diffSess := "test_diff_session_456"
	lockDiff, err := acquireSessionLock(diffSess)
	if err != nil {
		t.Fatalf("Acquire different session lock failed: %v", err)
	}
	lockDiff.Release()

	// After lock1 release, acquiring testSess again must succeed
	lock1.Release()
	lock3, err := acquireSessionLock(testSess)
	if err != nil {
		t.Fatalf("Re-acquiring lock after release failed: %v", err)
	}
	lock3.Release()
}
