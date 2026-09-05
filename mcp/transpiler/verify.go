package transpiler

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"scorp-agent/mcp"
)

// ToolContract is the comparable projection of one MCP tool.
type ToolContract struct {
	Name         string   `json:"name"`
	Required     []string `json:"required,omitempty"`
	ParamTypes   map[string]string `json:"param_types,omitempty"`
}

// VerifyReport summarizes the contract diff between benchmark and build.
type VerifyReport struct {
	Matched        []string          `json:"matched"`
	Missing        []string          `json:"missing,omitempty"`
	ParamMismatches map[string][]string `json:"param_mismatches,omitempty"`
	Coverage       float64           `json:"coverage"` // 0..1
	Pass           bool              `json:"pass"`
	Details        string            `json:"details,omitempty"`
}

// VerifyContract boots the built binary (with the same arguments the probe
// used) and diffs its tools/list contract against the Phase-1 benchmark.
// Pass requires every benchmark tool present with matching parameter names
// and types.
func VerifyContract(ctx context.Context, bench *Benchmark, binPath string) (*VerifyReport, error) {
	srv, tools, err := mcp.ProbeServer("verify-"+bench.ServerName, mcp.MCPServerConfig{
		Command: binPath,
		Args:    bench.Args,
	})
	if err != nil {
		return nil, fmt.Errorf("transpiled binary failed to boot: %w", err)
	}
	defer srv.Close()

	want := contractsFromBenchmark(bench)
	got := contractsFromTools(tools)

	report := &VerifyReport{Matched: []string{}, Missing: []string{}, ParamMismatches: map[string][]string{}}
	for name, w := range want {
		g, ok := got[name]
		if !ok {
			report.Missing = append(report.Missing, name)
			continue
		}
		if diffs := diffContract(w, g); len(diffs) > 0 {
			report.ParamMismatches[name] = diffs
			continue
		}
		report.Matched = append(report.Matched, name)
	}

	sort.Strings(report.Matched)
	sort.Strings(report.Missing)
	for name := range report.ParamMismatches {
		sort.Strings(report.ParamMismatches[name])
	}

	total := len(want)
	if total == 0 {
		report.Pass = false
		report.Details = "benchmark had no tools"
		return report, nil
	}
	report.Coverage = float64(len(report.Matched)) / float64(total)
	report.Pass = len(report.Missing) == 0 && len(report.ParamMismatches) == 0

	var sb strings.Builder
	if len(report.Missing) > 0 {
		sb.WriteString("missing tools: " + strings.Join(report.Missing, ", ") + "; ")
	}
	for name, diffs := range report.ParamMismatches {
		sb.WriteString(fmt.Sprintf("%s: %s; ", name, strings.Join(diffs, ", ")))
	}
	report.Details = strings.TrimRight(sb.String(), "; ")
	return report, nil
}

func contractsFromBenchmark(bench *Benchmark) map[string]ToolContract {
	contracts := make(map[string]ToolContract)
	for _, t := range bench.Tools {
		contracts[t.Name] = ToolContract{
			Name:       t.Name,
			Required:   stringSlice(t.InputSchema["required"]),
			ParamTypes: paramTypes(t.InputSchema),
		}
	}
	return contracts
}

func contractsFromTools(tools []mcp.MCPTool) map[string]ToolContract {
	contracts := make(map[string]ToolContract)
	for _, t := range tools {
		contracts[t.Name] = ToolContract{
			Name:       t.Name,
			Required:   stringSlice(t.InputSchema["required"]),
			ParamTypes: paramTypes(t.InputSchema),
		}
	}
	return contracts
}

// diffContract returns human-readable mismatches between expected and actual.
func diffContract(want, got ToolContract) []string {
	var diffs []string

	wantReq, gotReq := setOf(want.Required), setOf(got.Required)
	for p := range wantReq {
		if !gotReq[p] {
			diffs = append(diffs, fmt.Sprintf("missing required param %q", p))
		}
	}
	for p := range gotReq {
		if !wantReq[p] {
			diffs = append(diffs, fmt.Sprintf("unexpected required param %q", p))
		}
	}
	for p, wt := range want.ParamTypes {
		gt, ok := got.ParamTypes[p]
		if !ok {
			diffs = append(diffs, fmt.Sprintf("missing param %q", p))
			continue
		}
		if wt != gt {
			diffs = append(diffs, fmt.Sprintf("param %q type %s != expected %s", p, gt, wt))
		}
	}
	return diffs
}

func paramTypes(schema map[string]interface{}) map[string]string {
	types := map[string]string{}
	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		return types
	}
	for name, raw := range props {
		prop, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		t, _ := prop["type"].(string)
		if t == "" {
			t = "string"
		}
		// mcp-go's WithNumber always emits "number"; treat JSON "integer" as
		// the same numeric contract (the SDK cannot express "integer").
		if t == "integer" {
			t = "number"
		}
		types[name] = t
	}
	return types
}

func stringSlice(raw interface{}) []string {
	list, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	var out []string
	for _, v := range list {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func setOf(items []string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, i := range items {
		m[i] = true
	}
	return m
}
