package metrics

import (
	"scorp-agent/config"
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ──────────────────────────────────────────────
// Prometheus Metrics
// ──────────────────────────────────────────────

var (
	// Counters
	metricAgentIterations = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "scorp_agent_iterations_total",
		Help: "Total number of agent iterations.",
	})
	metricToolCalls = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "scorp_tool_calls_total",
		Help: "Total number of tool calls, by tool name.",
	}, []string{"tool"})
	metricMessagesReceived = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "scorp_messages_received_total",
		Help: "Total number of messages received from Telegram.",
	})
	metricMessagesSent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "scorp_messages_sent_total",
		Help: "Total number of messages sent to Telegram.",
	})
	metricSubagentsCompleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "scorp_subagents_completed_total",
		Help: "Total number of subagent executions completed.",
	})
	metricErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "scorp_errors_total",
		Help: "Total errors by type.",
	}, []string{"type"})

	// Gauges (current state)
	metricActiveSessions = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scorp_active_sessions",
		Help: "Current number of active sessions.",
	})
	metricMemoryItems = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scorp_memory_items",
		Help: "Current number of memory cache items.",
	})
	metricScheduledTasks = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scorp_scheduled_tasks",
		Help: "Current number of scheduled tasks.",
	})
	metricRAGChunks = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scorp_rag_chunks",
		Help: "Current number of RAG index chunks.",
	})
	metricUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scorp_up",
		Help: "1 if the scorp-agent is running, 0 otherwise.",
	})

	// Histograms
	metricResponseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "scorp_response_time_seconds",
		Help:    "Response time for AI model calls.",
		Buckets: prometheus.DefBuckets, // .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
	})
	metricToolExecTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "scorp_tool_execution_seconds",
		Help:    "Execution time for tool calls.",
		Buckets: []float64{.01, .05, .1, .25, .5, 1, 2.5, 5, 10, 30, 60},
	})

	metricsMu     sync.Mutex
	metricsServer *http.Server
	metricsInit   bool
)

func Init() {
	if metricsInit {
		return
	}
	metricsInit = true

	prometheus.MustRegister(metricAgentIterations)
	prometheus.MustRegister(metricToolCalls)
	prometheus.MustRegister(metricMessagesReceived)
	prometheus.MustRegister(metricMessagesSent)
	prometheus.MustRegister(metricSubagentsCompleted)
	prometheus.MustRegister(metricErrors)
	prometheus.MustRegister(metricActiveSessions)
	prometheus.MustRegister(metricMemoryItems)
	prometheus.MustRegister(metricScheduledTasks)
	prometheus.MustRegister(metricRAGChunks)
	prometheus.MustRegister(metricUp)
	prometheus.MustRegister(metricResponseTime)
	prometheus.MustRegister(metricToolExecTime)

	metricUp.Set(1)
	log.Println("[metrics] Prometheus metrics registered")
}

// startMetricsServer starts the /metrics HTTP endpoint.
// Listen on :9090 by default (configurable via METRICS_PORT env).
func StartServer() {
	Init()

	port := config.EnvStr("METRICS_PORT", "9091")
	addr := "127.0.0.1:" + port

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	metricsServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Printf("[metrics] Listening on %s/metrics", addr)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[metrics] Server error: %v", err)
		}
	}()
}

// stopMetricsServer gracefully stops the metrics server.
func StopServer() {
	if metricsServer != nil {
		metricsServer.Close()
		metricUp.Set(0)
		log.Println("[metrics] Server stopped")
	}
}