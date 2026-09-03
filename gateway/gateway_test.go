package gateway

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "scorp-agent/bootstrap"
	"scorp-agent/config"
	"scorp-agent/models"
	"scorp-agent/sop"
)

func TestGatewayEndpoints(t *testing.T) {
	config.InitConfigManager()
	models.LoadModelConfig()
	sop.InitDefaultSOPs()

	// Test 1: Dashboard HTML
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handleDashboard(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Dashboard returned %d, expected 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Scorp Gateway") {
		t.Errorf("Dashboard HTML missing title")
	}

	// Test 2: API Status
	reqStatus := httptest.NewRequest("GET", "/api/status", nil)
	wStatus := httptest.NewRecorder()
	handleStatus(wStatus, reqStatus)

	if wStatus.Code != http.StatusOK {
		t.Fatalf("/api/status returned %d, expected 200", wStatus.Code)
	}
	if !strings.Contains(wStatus.Body.String(), `"status":"online"`) {
		t.Errorf("/api/status unexpected payload: %s", wStatus.Body.String())
	}

	// Test 3: API Tools
	reqTools := httptest.NewRequest("GET", "/api/tools", nil)
	wTools := httptest.NewRecorder()
	handleTools(wTools, reqTools)

	if wTools.Code != http.StatusOK {
		t.Fatalf("/api/tools returned %d, expected 200", wTools.Code)
	}
	if !strings.Contains(wTools.Body.String(), "read_file") {
		t.Logf("Returned tools: %s", wTools.Body.String())
		t.Errorf("/api/tools expected to contain read_file")
	}

	// Test 4: API SOPs
	reqSOPs := httptest.NewRequest("GET", "/api/sops", nil)
	wSOPs := httptest.NewRecorder()
	handleSOPs(wSOPs, reqSOPs)

	if wSOPs.Code != http.StatusOK {
		t.Fatalf("/api/sops returned %d, expected 200", wSOPs.Code)
	}
	if !strings.Contains(wSOPs.Body.String(), "health_audit") {
		t.Errorf("/api/sops expected to contain health_audit")
	}
}
