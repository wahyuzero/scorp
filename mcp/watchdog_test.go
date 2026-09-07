package mcp

import (
	"os/exec"
	"sync"
	"testing"
)

// TestRegisterWatchdogReusesCounter pins the C6/F-13 fix: re-registering a
// watchdog after a restart must REUSE the existing watchdog so restartCount
// persists. Recreating it reset the counter and a fast-crashing server
// restarted forever (always "Attempt 1/5", max never reached).
func TestRegisterWatchdogReusesCounter(t *testing.T) {
	watchdogsMu.Lock()
	delete(watchdogs, "t_wd_reuse")
	watchdogsMu.Unlock()
	defer func() {
		watchdogsMu.Lock()
		wd := watchdogs["t_wd_reuse"]
		if wd != nil {
			wd.mu.Lock()
			wd.stopped = true
			wd.mu.Unlock()
		}
		delete(watchdogs, "t_wd_reuse")
		watchdogsMu.Unlock()
	}()

	srv1 := &MCPServer{Name: "t_wd_reuse", waitDone: make(chan struct{})}
	RegisterWatchdog("t_wd_reuse", srv1)

	wd1 := watchdogs["t_wd_reuse"]
	if wd1 == nil {
		t.Fatal("watchdog must be registered")
	}
	// Simulate a restart: monitor hand-off must reuse the SAME watchdog.
	wd1.mu.Lock()
	wd1.restartCount = 3
	wd1.mu.Unlock()

	srv2 := &MCPServer{Name: "t_wd_reuse", waitDone: make(chan struct{})}
	RegisterWatchdog("t_wd_reuse", srv2)

	wd2 := watchdogs["t_wd_reuse"]
	if wd2 != wd1 {
		t.Fatal("RegisterWatchdog must reuse the existing watchdog (counter persistence)")
	}
	wd2.mu.Lock()
	count := wd2.restartCount
	cur := wd2.current
	wd2.mu.Unlock()
	if count != 3 {
		t.Fatalf("restartCount must persist across re-registration, got %d", count)
	}
	if cur != srv2 {
		t.Fatal("watchdog must track the newest server instance")
	}

	wd2.mu.Lock()
	wd2.stopped = true
	wd2.mu.Unlock()
}

// TestReapWaitSingleFlight pins the C6/D7/F-21 fix: cmd.Wait must be reaped
// exactly once even when the watchdog monitor and Close race to wait. Run
// under -race: the old double-Wait tripped DATA RACE here.
func TestReapWaitSingleFlight(t *testing.T) {
	cmd := exec.Command("sleep", "0.05")
	if err := cmd.Start(); err != nil {
		t.Skipf("cannot spawn sleep: %v", err)
	}
	s := &MCPServer{Name: "t_reap", cmd: cmd, waitDone: make(chan struct{})}

	const waiters = 4
	var wg sync.WaitGroup
	errs := make([]error, waiters)
	for i := 0; i < waiters; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = s.reapWait()
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("waiter %d got error (double-Wait symptom): %v", i, err)
		}
	}
}
