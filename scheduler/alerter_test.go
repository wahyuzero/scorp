package scheduler

import (
	"scorp-agent/collectors"
	"scorp-agent/config"
	"sync"
	"testing"
	"time"
)

// ──────────────────────────────────────────────
// Helper: track  cfg  state
// ──────────────────────────────────────────────

func setAlertThresholds(cpu, ram, disk, load, swap float64, cooldownSec int) {
	alertMu.Lock()
	config.Cfg.CPUThreshold = cpu
	config.Cfg.RAMThreshold = ram
	config.Cfg.DiskThreshold = disk
	config.Cfg.LoadThreshold = load
	config.Cfg.SwapThresholdMB = swap
	config.Cfg.AlertCooldown = cooldownSec
	alertCooldowns = make(map[string]time.Time)
	alertCallCount = 0
	alertMu.Unlock()
}

func resetAlertCooldowns() {
	alertMu.Lock()
	alertCooldowns = make(map[string]time.Time)
	alertCallCount = 0
	alertMu.Unlock()
}

// ──────────────────────────────────────────────
// canFire — cooldown logic
// ──────────────────────────────────────────────

func TestCanFire(t *testing.T) {
	setAlertThresholds(90, 85, 85, 4, 512, 3600) // 1-hour cooldown
	defer resetAlertCooldowns()

	t.Run("first call always true", func(t *testing.T) {
		resetAlertCooldowns()
		if !CanFire("test_key") {
			t.Error("canFire should return true on first call")
		}
	})

	t.Run("second call blocked by cooldown", func(t *testing.T) {
		resetAlertCooldowns()
		CanFire("blocked_key") // first call, sets cooldown
		if CanFire("blocked_key") {
			t.Error("canFire should return false within cooldown period")
		}
	})

	t.Run("different keys independent", func(t *testing.T) {
		resetAlertCooldowns()
		CanFire("key_a")
		if !CanFire("key_b") {
			t.Error("canFire for key_b should return true independently of key_a")
		}
		if CanFire("key_a") {
			t.Error("key_a should still be in cooldown")
		}
	})

	t.Run("cooldown expires", func(t *testing.T) {
		resetAlertCooldowns()
		setAlertThresholds(90, 85, 85, 4, 512, 1) // 1-second cooldown
		CanFire("expire_key")
		if CanFire("expire_key") {
			t.Error("1-second cooldown should not expire immediately")
		}
		time.Sleep(1100 * time.Millisecond) // wait for cooldown expiry
		if !CanFire("expire_key") {
			t.Error("canFire should return true after cooldown expires")
		}
	})

	t.Run("cleanup after 100 calls", func(t *testing.T) {
		resetAlertCooldowns()
		setAlertThresholds(90, 85, 85, 4, 512, 60) // 60-second cooldown

		// Call with keys a,b,c to set cooldowns
		CanFire("cleanup_a")
		CanFire("cleanup_b")
		CanFire("cleanup_c")

		// Manually zero out cleanup_a's entry to simulate it being stale
		alertMu.Lock()
		alertCooldowns["cleanup_a"] = time.Now().Add(-120 * time.Second) // expired
		alertMu.Unlock()

		// Fire 97 more calls to trigger cleanup (total 100)
		for i := 0; i < 97; i++ {
			CanFire("cleanup_call")
		}

		// cleanup_a was expired → should have been deleted
		alertMu.Lock()
		_, aExists := alertCooldowns["cleanup_a"]
		_, bExists := alertCooldowns["cleanup_b"]
		alertMu.Unlock()
		if !aExists {
			t.Log("cleanup_a correctly deleted (expired)")
		}
		if !bExists {
			t.Error("cleanup_b should still exist (not expired yet)")
		}
	})
}

// ──────────────────────────────────────────────
// checkSystemAlerts
// ──────────────────────────────────────────────

func TestCheckSystemAlerts(t *testing.T) {
	t.Run("no alerts when all below threshold", func(t *testing.T) {
		setAlertThresholds(90, 85, 85, 4, 512, 3600)
		defer resetAlertCooldowns()

		d := collectors.SystemData{
			CPUPercent:  30,
			RAMPercent:  40,
			DiskPercent: 50,
			LoadAvg:     [3]float64{1.5, 1.2, 1.0},
			SwapUsedGB:  0.1,
			SwapTotalGB: 2.0,
			CPUCount:    4,
		}

		alerts := CheckSystemAlerts(d)
		if len(alerts) != 0 {
			t.Errorf("expected 0 alerts, got %d: %v", len(alerts), alerts)
		}
	})

	t.Run("CPU alert when above threshold", func(t *testing.T) {
		setAlertThresholds(50, 85, 85, 4, 512, 3600)
		defer resetAlertCooldowns()

		d := collectors.SystemData{
			CPUPercent: 85.5,
			LoadAvg:    [3]float64{1.5, 1.2, 1.0},
			CPUCount:   4,
		}

		alerts := CheckSystemAlerts(d)
		foundCPU := false
		for _, a := range alerts {
			if contains(a, "HIGH CPU") {
				foundCPU = true
			}
		}
		if !foundCPU {
			t.Errorf("expected HIGH CPU alert, got: %v", alerts)
		}
	})

	t.Run("RAM alert when above threshold", func(t *testing.T) {
		setAlertThresholds(90, 50, 85, 4, 512, 3600)
		defer resetAlertCooldowns()

		d := collectors.SystemData{
			RAMPercent: 72.3,
			RAMUsedGB:  6.0,
			RAMTotalGB: 8.0,
			RAMAvailGB: 2.0,
		}

		alerts := CheckSystemAlerts(d)
		foundRAM := false
		for _, a := range alerts {
			if contains(a, "HIGH MEMORY") {
				foundRAM = true
			}
		}
		if !foundRAM {
			t.Errorf("expected HIGH MEMORY alert, got: %v", alerts)
		}
	})

	t.Run("Disk alert when above threshold", func(t *testing.T) {
		setAlertThresholds(90, 85, 50, 4, 512, 3600)
		defer resetAlertCooldowns()

		d := collectors.SystemData{
			DiskPercent: 91.0,
			DiskUsedGB:  200.0,
			DiskTotalGB: 220.0,
		}

		alerts := CheckSystemAlerts(d)
		foundDisk := false
		for _, a := range alerts {
			if contains(a, "HIGH DISK") {
				foundDisk = true
			}
		}
		if !foundDisk {
			t.Errorf("expected HIGH DISK alert, got: %v", alerts)
		}
	})

	t.Run("Load alert when above threshold", func(t *testing.T) {
		setAlertThresholds(90, 85, 85, 3.0, 512, 3600)
		defer resetAlertCooldowns()

		d := collectors.SystemData{
			LoadAvg:  [3]float64{8.5, 6.2, 4.1},
			CPUCount: 4,
		}

		alerts := CheckSystemAlerts(d)
		foundLoad := false
		for _, a := range alerts {
			if contains(a, "HIGH LOAD") {
				foundLoad = true
			}
		}
		if !foundLoad {
			t.Errorf("expected HIGH LOAD alert, got: %v", alerts)
		}
	})

	t.Run("Swap alert when above threshold", func(t *testing.T) {
		setAlertThresholds(90, 85, 85, 4, 100, 3600) // 100 MB threshold
		defer resetAlertCooldowns()

		d := collectors.SystemData{
			SwapUsedGB:  2.0,
			SwapTotalGB: 4.0,
		}

		alerts := CheckSystemAlerts(d)
		foundSwap := false
		for _, a := range alerts {
			if contains(a, "SWAP USAGE") {
				foundSwap = true
			}
		}
		if !foundSwap {
			t.Errorf("expected SWAP USAGE alert, got: %v", alerts)
		}
	})

	t.Run("multiple alerts at once", func(t *testing.T) {
		setAlertThresholds(30, 30, 30, 1.0, 10, 3600)
		defer resetAlertCooldowns()

		d := collectors.SystemData{
			CPUPercent:  90,
			RAMPercent:  90,
			DiskPercent: 90,
			LoadAvg:     [3]float64{8.0, 6.0, 4.0},
			SwapUsedGB:  1.0,
			SwapTotalGB: 2.0,
			CPUCount:    4,
			RAMUsedGB:   7.2,
			RAMTotalGB:  8.0,
			RAMAvailGB:  0.8,
			DiskUsedGB:  200.0,
			DiskTotalGB: 220.0,
		}

		alerts := CheckSystemAlerts(d)
		if len(alerts) != 5 {
			t.Errorf("expected 5 alerts (CPU+RAM+Disk+Load+Swap), got %d: %v", len(alerts), alerts)
		}
	})

	t.Run("cooldown prevents duplicate CPU alert", func(t *testing.T) {
		setAlertThresholds(50, 85, 85, 4, 512, 600)
		defer resetAlertCooldowns()

		d := collectors.SystemData{
			CPUPercent: 85.5,
			LoadAvg:    [3]float64{1.5, 1.2, 1.0},
			CPUCount:   4,
		}

		alerts1 := CheckSystemAlerts(d) // 1st call triggers
		alerts2 := CheckSystemAlerts(d) // 2nd call blocked by cooldown

		if len(alerts1) != 1 {
			t.Errorf("expected 1 alert on first call, got %d", len(alerts1))
		}
		if len(alerts2) != 0 {
			t.Errorf("expected 0 alerts on second call (cooldown), got %d", len(alerts2))
		}
	})

	t.Run("top processes in CPU alert", func(t *testing.T) {
		setAlertThresholds(50, 85, 85, 4, 512, 3600)
		defer resetAlertCooldowns()

		d := collectors.SystemData{
			CPUPercent: 95.0,
			LoadAvg:    [3]float64{2.0, 1.5, 1.0},
			CPUCount:   4,
			TopProcesses: []collectors.TopProcess{
				{Name: "python3", CPU: 45.2, Mem: 12.1},
				{Name: "java", CPU: 30.0, Mem: 8.0},
			},
		}

		alerts := CheckSystemAlerts(d)
		if len(alerts) == 0 {
			t.Fatal("expected at least one alert")
		}
		if !contains(alerts[0], "python3") || !contains(alerts[0], "java") {
			t.Errorf("CPU alert should include top processes: %s", alerts[0])
		}
	})
}

// ──────────────────────────────────────────────
// checkDockerAlerts
// ──────────────────────────────────────────────

func TestCheckDockerAlerts(t *testing.T) {
	setAlertThresholds(90, 85, 85, 4, 512, 3600)
	defer resetAlertCooldowns()

	t.Run("healthy containers no alerts", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.DockerData{
			Containers: []collectors.ContainerInfo{
				{Name: "web", Status: "running", Health: "healthy", Image: "nginx"},
				{Name: "db", Status: "running", Health: "healthy", Image: "postgres"},
			},
		}

		alerts := CheckDockerAlerts(d)
		if len(alerts) != 0 {
			t.Errorf("expected 0 alerts for healthy containers, got %d: %v", len(alerts), alerts)
		}
	})

	t.Run("container down alert", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.DockerData{
			Containers: []collectors.ContainerInfo{
				{Name: "web", Status: "exited", Health: "healthy", Image: "nginx"},
			},
		}

		alerts := CheckDockerAlerts(d)
		foundDown := false
		for _, a := range alerts {
			if contains(a, "CONTAINER DOWN") && contains(a, "web") {
				foundDown = true
			}
		}
		if !foundDown {
			t.Errorf("expected CONTAINER DOWN alert for 'web', got: %v", alerts)
		}
	})

	t.Run("container unhealthy alert", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.DockerData{
			Containers: []collectors.ContainerInfo{
				{Name: "web", Status: "running", Health: "unhealthy", Image: "nginx"},
			},
		}

		alerts := CheckDockerAlerts(d)
		foundUnhealthy := false
		for _, a := range alerts {
			if contains(a, "CONTAINER UNHEALTHY") {
				foundUnhealthy = true
			}
		}
		if !foundUnhealthy {
			t.Errorf("expected CONTAINER UNHEALTHY alert, got: %v", alerts)
		}
	})

	t.Run("both down and unhealthy", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.DockerData{
			Containers: []collectors.ContainerInfo{
				{Name: "web", Status: "exited", Health: "unhealthy", Image: "nginx"},
			},
		}

		alerts := CheckDockerAlerts(d)
		if len(alerts) != 2 {
			t.Errorf("expected 2 alerts (down + unhealthy), got %d: %v", len(alerts), alerts)
		}
	})

	t.Run("multiple containers with issues", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.DockerData{
			Containers: []collectors.ContainerInfo{
				{Name: "web", Status: "exited", Health: "unhealthy", Image: "nginx"},
				{Name: "worker", Status: "exited", Health: "", Image: "redis"},
				{Name: "healthy_db", Status: "running", Health: "healthy", Image: "postgres"},
			},
		}

		alerts := CheckDockerAlerts(d)
		// web: down + unhealthy = 2, worker: down = 1, healthy_db: 0 = 3 total
		if len(alerts) != 3 {
			t.Errorf("expected 3 alerts (web down, web unhealthy, worker down), got %d: %v", len(alerts), alerts)
		}
	})

	t.Run("cooldown per container", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.DockerData{
			Containers: []collectors.ContainerInfo{
				{Name: "web", Status: "exited", Health: "healthy", Image: "nginx"},
			},
		}

		alerts1 := CheckDockerAlerts(d) // fires
		alerts2 := CheckDockerAlerts(d) // cooldown

		if len(alerts1) != 1 {
			t.Errorf("expected 1 alert on first call, got %d", len(alerts1))
		}
		if len(alerts2) != 0 {
			t.Errorf("expected 0 alerts on second call (cooldown), got %d", len(alerts2))
		}
	})
}

// ──────────────────────────────────────────────
// checkStorageAlerts
// ──────────────────────────────────────────────

func TestCheckStorageAlerts(t *testing.T) {
	setAlertThresholds(90, 85, 85, 4, 512, 3600)
	defer resetAlertCooldowns()

	t.Run("no alerts when mounted and reachable", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.StorageData{
			GDriveMount: collectors.GDriveInfo{Mounted: true, Path: "/gdrive"},
			S3Gateway:   collectors.S3Info{Reachable: true, URL: "https://s3.example.com"},
		}
		alerts := CheckStorageAlerts(d)
		if len(alerts) != 0 {
			t.Errorf("expected 0 alerts, got %d: %v", len(alerts), alerts)
		}
	})

	t.Run("gdrive unmounted alert", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.StorageData{
			GDriveMount: collectors.GDriveInfo{Mounted: false, Path: "/gdrive"},
			S3Gateway:   collectors.S3Info{Reachable: true, URL: "https://s3.example.com"},
		}
		alerts := CheckStorageAlerts(d)
		foundGDrive := false
		for _, a := range alerts {
			if contains(a, "GDRIVE UNMOUNTED") {
				foundGDrive = true
			}
		}
		if !foundGDrive {
			t.Errorf("expected GDRIVE UNMOUNTED alert, got: %v", alerts)
		}
	})

	t.Run("s3 unreachable alert", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.StorageData{
			GDriveMount: collectors.GDriveInfo{Mounted: true, Path: "/gdrive"},
			S3Gateway:   collectors.S3Info{Reachable: false, URL: "https://s3.example.com", Error: "connection refused"},
		}
		alerts := CheckStorageAlerts(d)
		foundS3 := false
		for _, a := range alerts {
			if contains(a, "S3 GATEWAY DOWN") {
				foundS3 = true
			}
		}
		if !foundS3 {
			t.Errorf("expected S3 GATEWAY DOWN alert, got: %v", alerts)
		}
	})

	t.Run("both alerts when both down", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.StorageData{
			GDriveMount: collectors.GDriveInfo{Mounted: false, Path: "/gdrive"},
			S3Gateway:   collectors.S3Info{Reachable: false, URL: "https://s3.example.com", Error: "timeout"},
		}
		alerts := CheckStorageAlerts(d)
		if len(alerts) != 2 {
			t.Errorf("expected 2 alerts, got %d: %v", len(alerts), alerts)
		}
	})
}

// ──────────────────────────────────────────────
// checkNetworkAlerts
// ──────────────────────────────────────────────

func TestCheckNetworkAlerts(t *testing.T) {
	setAlertThresholds(90, 85, 85, 4, 512, 3600)
	defer resetAlertCooldowns()

	t.Run("no alerts when ports stable and traefik fine", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.NetworkData{
			NewPorts: []collectors.PortInfo{},
			Traefik:  collectors.TraefikInfo{Error5xx: 0, TotalRequests: 100},
		}
		alerts := CheckNetworkAlerts(d)
		if len(alerts) != 0 {
			t.Errorf("expected 0 alerts, got %d: %v", len(alerts), alerts)
		}
	})

	t.Run("new port detected alert", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.NetworkData{
			NewPorts: []collectors.PortInfo{
				{Port: 8443, Process: "nginx", Address: "0.0.0.0"},
			},
			Traefik: collectors.TraefikInfo{Error5xx: 0, TotalRequests: 100},
		}
		alerts := CheckNetworkAlerts(d)
		foundPort := false
		for _, a := range alerts {
			if contains(a, "NEW PORT") && contains(a, "8443") {
				foundPort = true
			}
		}
		if !foundPort {
			t.Errorf("expected NEW PORT alert for 8443, got: %v", alerts)
		}
	})

	t.Run("multiple new ports", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.NetworkData{
			NewPorts: []collectors.PortInfo{
				{Port: 8080, Process: "java", Address: "0.0.0.0"},
				{Port: 9090, Process: "metrics", Address: "127.0.0.1"},
			},
			Traefik: collectors.TraefikInfo{Error5xx: 0, TotalRequests: 100},
		}
		alerts := CheckNetworkAlerts(d)
		if len(alerts) != 2 {
			t.Errorf("expected 2 alerts for 2 new ports, got %d: %v", len(alerts), alerts)
		}
	})

	t.Run("traefik high 5xx alert", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.NetworkData{
			NewPorts: []collectors.PortInfo{},
			Traefik:  collectors.TraefikInfo{Error5xx: 25, TotalRequests: 500},
		}
		alerts := CheckNetworkAlerts(d)
		found5xx := false
		for _, a := range alerts {
			if contains(a, "HIGH 5xx") {
				found5xx = true
			}
		}
		if !found5xx {
			t.Errorf("expected HIGH 5xx alert, got: %v", alerts)
		}
	})

	t.Run("low 5xx below threshold", func(t *testing.T) {
		resetAlertCooldowns()
		d := collectors.NetworkData{
			NewPorts: []collectors.PortInfo{},
			Traefik:  collectors.TraefikInfo{Error5xx: 5, TotalRequests: 500}, // below 10
		}
		alerts := CheckNetworkAlerts(d)
		for _, a := range alerts {
			if contains(a, "5xx") {
				t.Errorf("unexpected 5xx alert when Error5xx=5 < 10: %s", a)
			}
		}
	})
}

// ──────────────────────────────────────────────
// formatTopProcs
// ──────────────────────────────────────────────

func TestFormatTopProcs(t *testing.T) {
	t.Run("empty returns N/A", func(t *testing.T) {
		result := formatTopProcs([]collectors.TopProcess{})
		if result != "N/A" {
			t.Errorf("expected 'N/A', got %q", result)
		}
	})

	t.Run("single process", func(t *testing.T) {
		result := formatTopProcs([]collectors.TopProcess{
			{Name: "python3", CPU: 45.2, Mem: 12.1},
		})
		if !contains(result, "python3") || !contains(result, "45.2") || !contains(result, "12.1") {
			t.Errorf("unexpected format: %q", result)
		}
	})

	t.Run("truncated to 3 processes", func(t *testing.T) {
		procs := []collectors.TopProcess{
			{Name: "p1", CPU: 50.0, Mem: 10.0},
			{Name: "p2", CPU: 30.0, Mem: 8.0},
			{Name: "p3", CPU: 20.0, Mem: 5.0},
			{Name: "p4", CPU: 10.0, Mem: 3.0},
		}
		result := formatTopProcs(procs)
		if !contains(result, "p1") || !contains(result, "p2") || !contains(result, "p3") {
			t.Errorf("expected p1, p2, p3 in result, got: %q", result)
		}
		if contains(result, "p4") {
			t.Errorf("p4 should not appear (truncated to 3): %q", result)
		}
		// Should have 3 lines
		lines := 0
		for _, c := range result {
			if c == '\n' {
				lines++
			}
		}
		if lines != 2 { // 3 lines = 2 newlines
			t.Errorf("expected 3 lines (2 newlines), got %d lines in: %q", lines+1, result)
		}
	})
}

// ──────────────────────────────────────────────
// TestContainsHelper
// ──────────────────────────────────────────────

func TestContainsHelper(t *testing.T) {
	if !contains("hello world", "world") {
		t.Error("contains should find 'world' in 'hello world'")
	}
	if contains("hello world", "xyz") {
		t.Error("contains should not find 'xyz' in 'hello world'")
	}
}

// ──────────────────────────────────────────────
// Integration: concurrent canFire safety
// ──────────────────────────────────────────────

func TestCanFireConcurrentSafety(t *testing.T) {
	setAlertThresholds(90, 85, 85, 4, 512, 10)
	defer resetAlertCooldowns()

	var wg sync.WaitGroup
	results := make([]bool, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = CanFire("concurrent_key")
		}(i)
	}
	wg.Wait()

	// Only one should have gotten true
	trueCount := 0
	for _, r := range results {
		if r {
			trueCount++
		}
	}
	if trueCount != 1 {
		t.Errorf("expected exactly 1 true from concurrent canFire, got %d", trueCount)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
