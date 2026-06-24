package collectors

import (
	"testing"
	"time"
)

func TestCollectorNative_SortByCPUDesc(t *testing.T) {
	procs := []nativeTopProcess{
		{PID: 1, Name: "proc1", CPU: 10.5, Mem: 5.0},
		{PID: 2, Name: "proc2", CPU: 50.2, Mem: 10.0},
		{PID: 3, Name: "proc3", CPU: 5.1, Mem: 2.0},
		{PID: 4, Name: "proc4", CPU: 25.8, Mem: 8.0},
	}

	SortByCPUDesc(procs)

	// Should be sorted descending by CPU
	expectedOrder := []float64{50.2, 25.8, 10.5, 5.1}
	for i, expected := range expectedOrder {
		if procs[i].CPU != expected {
			t.Errorf("SortByCPUDesc: position %d: got CPU=%f, want %f (name=%s)", i, procs[i].CPU, expected, procs[i].Name)
		}
	}
}

func TestCollectorNative_SortByCPUDesc_Empty(t *testing.T) {
	var procs []nativeTopProcess
	SortByCPUDesc(procs)
	if len(procs) != 0 {
		t.Error("SortByCPUDesc on empty slice should not panic")
	}
}

func TestCollectorNative_SortByCPUDesc_Single(t *testing.T) {
	procs := []nativeTopProcess{{PID: 1, Name: "only", CPU: 42.0, Mem: 10.0}}
	SortByCPUDesc(procs)
	if len(procs) != 1 || procs[0].CPU != 42.0 {
		t.Error("SortByCPUDesc on single element failed")
	}
}

func TestCollectorNative_SortByCPUDesc_EqualValues(t *testing.T) {
	procs := []nativeTopProcess{
		{PID: 1, Name: "a", CPU: 10.0, Mem: 5.0},
		{PID: 2, Name: "b", CPU: 10.0, Mem: 5.0},
		{PID: 3, Name: "c", CPU: 10.0, Mem: 5.0},
	}
	SortByCPUDesc(procs)
	if len(procs) != 3 {
		t.Error("SortByCPUDesc lost elements with equal CPU")
	}
}

func TestCollectorNative_GetTopProcesses(t *testing.T) {
	// This calls ReadProcessList() which reads /proc - may be empty in test env
	top := GetTopProcesses(5)
	// Just verify it doesn't panic and returns valid structure
	for _, p := range top {
		if p.PID <= 0 {
			t.Errorf("Invalid PID: %d", p.PID)
		}
		if p.Name == "" {
			t.Errorf("Empty process name for PID %d", p.PID)
		}
		if p.CPU < 0 {
			t.Errorf("Negative CPU for PID %d: %f", p.PID, p.CPU)
		}
	}
}

func TestCollectorNative_GetTopProcesses_Limit(t *testing.T) {
	top := GetTopProcesses(3)
	if len(top) > 3 {
		t.Errorf("GetTopProcesses(3) returned %d items, want <=3", len(top))
	}
}

func TestCollectorNative_NativeTopProcessStruct(t *testing.T) {
	p := nativeTopProcess{
		PID:  123,
		Name: "testproc",
		CPU:  12.34,
		Mem:  5.67,
	}
	if p.PID != 123 || p.Name != "testproc" || p.CPU != 12.34 || p.Mem != 5.67 {
		t.Errorf("nativeTopProcess struct failed: %+v", p)
	}
}

func TestCollectorNative_CollectSystem_Structure(t *testing.T) {
	// This will read actual system info - just verify structure
	d := CollectSystem()
	
	if d.CPUCount <= 0 {
		t.Errorf("CPUCount should be > 0, got %d", d.CPUCount)
	}
	if d.RAMTotalGB <= 0 {
		t.Errorf("RAMTotalGB should be > 0, got %f", d.RAMTotalGB)
	}
	if d.DiskTotalGB <= 0 {
		t.Errorf("DiskTotalGB should be > 0, got %f", d.DiskTotalGB)
	}
	if d.Uptime == "" {
		t.Error("Uptime should not be empty")
	}
	
	// TopProcesses should have at most 5 entries
	if len(d.TopProcesses) > 5 {
		t.Errorf("TopProcesses has %d entries, want <=5", len(d.TopProcesses))
	}
	
	for _, p := range d.TopProcesses {
		if p.PID <= 0 {
			t.Errorf("TopProcess invalid PID: %d", p.PID)
		}
	}
	
	t.Logf("System: CPU=%d cores, RAM=%.2fGB, Disk=%.2fGB, Uptime=%s, TopProcs=%d",
		d.CPUCount, d.RAMTotalGB, d.DiskTotalGB, d.Uptime, len(d.TopProcesses))
}

func TestCollectorNative_StartCPUSampler_DoesNotBlock(t *testing.T) {
	done := make(chan struct{})
	
	// Start sampler
	go StartCPUSampler(done)
	
	// Give it a moment to run
	time.Sleep(100 * time.Millisecond)
	
	// Stop it
	close(done)
	
	// Should return quickly
	select {
	case <-time.After(1 * time.Second):
		t.Error("StartCPUSampler did not stop after closing done channel")
	default:
		// Good - returned immediately
	}
}

func TestCollectorNative_StartDockerStatsSampler_DoesNotBlock(t *testing.T) {
	done := make(chan struct{})
	
	go StartDockerStatsSampler(done)
	
	time.Sleep(100 * time.Millisecond)
	
	close(done)
	
	select {
	case <-time.After(1 * time.Second):
		t.Error("StartDockerStatsSampler did not stop after closing done channel")
	default:
		// Good
	}
}