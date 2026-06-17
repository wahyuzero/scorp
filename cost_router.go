package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// S05: Cost-Aware Routing + Time-Based Auto-Switch
//
// Tracks cost per model (est. $/1M tokens), enforces daily budgets,
// and auto-switches to cheaper models during off-peak hours.
// ──────────────────────────────────────────────

// ModelCost tracks pricing and budget for cost-aware routing.
type ModelCost struct {
	InputPer1M  float64 `json:"input_per_1m"`  // $/1M input tokens
	OutputPer1M float64 `json:"output_per_1m"` // $/1M output tokens
}

// CostConfig holds cost-aware routing settings.
type CostConfig struct {
	DailyBudgetUSD float64                `json:"daily_budget_usd"` // 0 = unlimited
	ModelCosts     map[string]ModelCost   `json:"model_costs"`
	OffPeakModel   string                 `json:"offpeak_model"`   // cheaper model for off-peak
	OffPeakStart   int                    `json:"offpeak_start"`   // hour 0-23 (e.g., 22)
	OffPeakEnd     int                    `json:"offpeak_end"`     // hour 0-23 (e.g., 7)
	Enabled        bool                   `json:"enabled"`
}

// CostTracker accumulates daily cost.
type CostTracker struct {
	mu         sync.Mutex
	Date       string  `json:"date"`        // "2026-06-17"
	TotalUSD   float64 `json:"total_usd"`
	PerModel   map[string]float64 `json:"per_model"`
}

var (
	costCfg     *CostConfig
	costCfgMu   sync.RWMutex
	costTracker *CostTracker
	costPath    = os.ExpandEnv("$HOME") + "/.scorp-agent/cost_config.json"
	costLogPath = os.ExpandEnv("$HOME") + "/.scorp-agent/cost_daily.json"
)

func init() {
	costCfg = defaultCostConfig()
	costTracker = &CostTracker{PerModel: make(map[string]float64)}
}

func defaultCostConfig() *CostConfig {
	return &CostConfig{
		DailyBudgetUSD: 5.0,
		OffPeakModel:   "",
		OffPeakStart:   22,
		OffPeakEnd:     7,
		Enabled:        false,
		ModelCosts: map[string]ModelCost{
			"glm-4-flash":       {InputPer1M: 0.0, OutputPer1M: 0.0},
			"deepseek-chat":     {InputPer1M: 0.14, OutputPer1M: 0.28},
			"deepseek-coder":    {InputPer1M: 0.14, OutputPer1M: 0.28},
			"groq-llama-70b":    {InputPer1M: 0.59, OutputPer1M: 0.79},
			"gemini-flash":      {InputPer1M: 0.075, OutputPer1M: 0.30},
			"glm-4.6":           {InputPer1M: 0.60, OutputPer1M: 2.20},
			"glm-5.2":           {InputPer1M: 0.50, OutputPer1M: 2.00},
		},
	}
}

func loadCostConfig() {
	costCfgMu.Lock()
	defer costCfgMu.Unlock()

	data, err := os.ReadFile(costPath)
	if err != nil {
		log.Printf("[cost] No config found, using defaults")
		costCfg = defaultCostConfig()
		return
	}

	var cfg CostConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("[cost] Parse error: %v", err)
		costCfg = defaultCostConfig()
		return
	}
	costCfg = &cfg
	log.Printf("[cost] Config loaded: budget=$%.2f offpeak=%02d-%02d enabled=%v",
		costCfg.DailyBudgetUSD, costCfg.OffPeakStart, costCfg.OffPeakEnd, costCfg.Enabled)
}

func saveCostConfig() {
	costCfgMu.RLock()
	defer costCfgMu.RUnlock()
	if costCfg == nil {
		return
	}
	os.MkdirAll(os.ExpandEnv("$HOME")+"/.scorp-agent", 0755)
	data, _ := json.MarshalIndent(costCfg, "", "  ")
	os.WriteFile(costPath, data, 0644)
}

// loadCostTracker loads today's accumulated cost.
func loadCostTracker() {
	costTracker.mu.Lock()
	defer costTracker.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	data, err := os.ReadFile(costLogPath)
	if err != nil {
		costTracker = &CostTracker{Date: today, PerModel: make(map[string]float64)}
		return
	}

	var ct CostTracker
	if err := json.Unmarshal(data, &ct); err != nil {
		costTracker = &CostTracker{Date: today, PerModel: make(map[string]float64)}
		return
	}

	// Reset if new day
	if ct.Date != today {
		log.Printf("[cost] New day, resetting daily tracker (was %s)", ct.Date)
		costTracker = &CostTracker{Date: today, PerModel: make(map[string]float64)}
	} else {
		costTracker = &ct
		if costTracker.PerModel == nil {
			costTracker.PerModel = make(map[string]float64)
		}
	}
	log.Printf("[cost] Today's spend: $%.4f", costTracker.TotalUSD)
}

func saveCostTracker() {
	costTracker.mu.Lock()
	defer costTracker.mu.Unlock()
	os.MkdirAll(os.ExpandEnv("$HOME")+"/.scorp-agent", 0755)
	data, _ := json.MarshalIndent(costTracker, "", "  ")
	os.WriteFile(costLogPath, data, 0644)
}

// isOffPeak checks if current hour is in off-peak window.
func isOffPeak() bool {
	costCfgMu.RLock()
	defer costCfgMu.RUnlock()
	if costCfg == nil || !costCfg.Enabled || costCfg.OffPeakModel == "" {
		return false
	}
	hour := time.Now().Hour()
	start, end := costCfg.OffPeakStart, costCfg.OffPeakEnd
	if start <= end {
		return hour >= start && hour < end
	}
	// Wraps midnight (e.g., 22→7)
	return hour >= start || hour < end
}

// recordCost adds to today's cost tracker.
func recordCost(modelName string, inputTokens, outputTokens int) {
	costCfgMu.RLock()
	cfg := costCfg
	costCfgMu.RUnlock()

	if cfg == nil || !cfg.Enabled {
		return
	}

	cost, ok := cfg.ModelCosts[modelName]
	if !ok {
		return // Unknown model cost, skip
	}

	callCost := (float64(inputTokens)/1_000_000.0)*cost.InputPer1M +
		(float64(outputTokens)/1_000_000.0)*cost.OutputPer1M

	costTracker.mu.Lock()
	costTracker.TotalUSD += callCost
	costTracker.PerModel[modelName] += callCost
	costTracker.mu.Unlock()

	saveCostTracker()

	log.Printf("[cost] +%s $%.6f (in=%d out=%d) total=$%.4f",
		modelName, callCost, inputTokens, outputTokens, costTracker.getTotal())
}

func (ct *CostTracker) getTotal() float64 {
	return ct.TotalUSD
}

// isBudgetExceeded checks if daily budget is exceeded.
func isBudgetExceeded() bool {
	costCfgMu.RLock()
	defer costCfgMu.RUnlock()
	if costCfg == nil || !costCfg.Enabled || costCfg.DailyBudgetUSD <= 0 {
		return false
	}
	costTracker.mu.Lock()
	defer costTracker.mu.Unlock()
	return costTracker.TotalUSD >= costCfg.DailyBudgetUSD
}

// routeModelCostAware wraps routeModel with cost logic.
// If off-peak, switches to off-peak model. If budget exceeded, uses cheapest model.
func routeModelCostAware(taskType string) *ModelConfig {
	// Check off-peak first
	if isOffPeak() {
		costCfgMu.RLock()
		offPeakName := costCfg.OffPeakModel
		costCfgMu.RUnlock()

		if offPeakName != "" {
			if m := getModelByName(offPeakName); m != nil {
				log.Printf("[cost] Off-peak routing → %s", offPeakName)
				return m
			}
		}
	}

	// Check budget exceeded
	if isBudgetExceeded() {
		cheapest := getCheapestModel()
		if cheapest != nil {
			log.Printf("[cost] Budget exceeded ($%.4f), using cheapest: %s",
				costTracker.getTotal(), cheapest.Model)
			return cheapest
		}
	}

	return routeModel(taskType)
}

// getCheapestModel returns the lowest-cost model.
func getCheapestModel() *ModelConfig {
	costCfgMu.RLock()
	defer costCfgMu.RUnlock()
	if costCfg == nil || costCfg.ModelCosts == nil {
		return nil
	}

	var cheapest string
	var minCost = -1.0
	for name, c := range costCfg.ModelCosts {
		avg := (c.InputPer1M + c.OutputPer1M) / 2
		if minCost < 0 || avg < minCost {
			if m := getModelByName(name); m != nil {
				minCost = avg
				cheapest = name
			}
		}
	}

	if cheapest != "" {
		return getModelByName(cheapest)
	}
	return nil
}

// formatCostReport generates a cost summary for the models tool.
func formatCostReport() string {
	costCfgMu.RLock()
	defer costCfgMu.RUnlock()
	costTracker.mu.Lock()
	defer costTracker.mu.Unlock()

	var sb strings.Builder
	sb.WriteString("💰 **Cost Report**\n\n")

	if !costCfg.Enabled {
		sb.WriteString("Cost tracking: **disabled**\n")
		sb.WriteString("Enable: `models cost enable`\n")
		return sb.String()
	}

	budget := costCfg.DailyBudgetUSD
	used := costTracker.TotalUSD
	remaining := budget - used
	pct := 0.0
	if budget > 0 {
		pct = (used / budget) * 100
	}

	bar := makeBudgetBar(pct)
	sb.WriteString(fmt.Sprintf("**Budget:** $%.2f/day\n", budget))
	sb.WriteString(fmt.Sprintf("**Used:** $%.4f (%.1f%%)\n", used, pct))
	sb.WriteString(fmt.Sprintf("**Remaining:** $%.4f\n\n", remaining))
	sb.WriteString(fmt.Sprintf("```\n%s %.1f%%\n```\n\n", bar, pct))

	if len(costTracker.PerModel) > 0 {
		sb.WriteString("**Per-Model Spend:**\n")
		for model, amt := range costTracker.PerModel {
			sb.WriteString(fmt.Sprintf("  • `%s`: $%.4f\n", model, amt))
		}
	}

	// Off-peak info
	if costCfg.OffPeakModel != "" {
		sb.WriteString(fmt.Sprintf("\n**Off-Peak:** %s (%02d:00–%02d:00)",
			costCfg.OffPeakModel, costCfg.OffPeakStart, costCfg.OffPeakEnd))
		if isOffPeak() {
			sb.WriteString(" 🟢 active")
		}
	}

	return sb.String()
}

func makeBudgetBar(pct float64) string {
	width := 20
	filled := int(pct / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := filled; i < width; i++ {
		bar += "░"
	}
	return bar
}

// handleCostCommand processes cost subcommands in the models tool.
// Actions: enable, disable, budget <amount>, offpeak <model> <start> <end>,
// rate <model> <in> <out>, report
func handleCostCommand(args map[string]interface{}) (string, bool) {
	sub := getStringArg(args, "cost_action", "")

	switch sub {
	case "enable":
		costCfgMu.Lock()
		costCfg.Enabled = true
		costCfgMu.Unlock()
		saveCostConfig()
		loadCostTracker()
		return "✅ Cost tracking enabled.", true

	case "disable":
		costCfgMu.Lock()
		costCfg.Enabled = false
		costCfgMu.Unlock()
		saveCostConfig()
		return "⏸ Cost tracking disabled.", true

	case "budget":
		amount := getFloatArg(args, "amount", 5.0)
		costCfgMu.Lock()
		costCfg.DailyBudgetUSD = amount
		costCfgMu.Unlock()
		saveCostConfig()
		return fmt.Sprintf("✅ Daily budget set to $%.2f", amount), true

	case "offpeak":
		model := getStringArg(args, "offpeak_model", "")
		start := int(getFloatArg(args, "offpeak_start", 22))
		end := int(getFloatArg(args, "offpeak_end", 7))

		costCfgMu.Lock()
		costCfg.OffPeakModel = model
		costCfg.OffPeakStart = start
		costCfg.OffPeakEnd = end
		if !costCfg.Enabled {
			costCfg.Enabled = true
		}
		costCfgMu.Unlock()
		saveCostConfig()
		return fmt.Sprintf("✅ Off-peak model: %s (%02d:00–%02d:00)", model, start, end), true

	case "rate":
		model := getStringArg(args, "rate_model", "")
		inCost := getFloatArg(args, "input_cost", 0)
		outCost := getFloatArg(args, "output_cost", 0)

		costCfgMu.Lock()
		if costCfg.ModelCosts == nil {
			costCfg.ModelCosts = make(map[string]ModelCost)
		}
		costCfg.ModelCosts[model] = ModelCost{InputPer1M: inCost, OutputPer1M: outCost}
		costCfgMu.Unlock()
		saveCostConfig()
		return fmt.Sprintf("✅ Rate set: %s = $%.3f/$%.3f per 1M (in/out)", model, inCost, outCost), true

	case "report", "":
		return formatCostReport(), true

	default:
		return "Usage: `models cost <enable|disable|budget|offpeak|rate|report>`", false
	}
}
