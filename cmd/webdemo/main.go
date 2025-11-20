package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"demo-app/internal/analytics"
	"demo-app/internal/export"

	"github.com/yourorg/lcc-sdk/pkg/client"
	"github.com/yourorg/lcc-sdk/pkg/config"
)

//go:embed static/*
var staticFiles embed.FS

var (
	lccClient       *client.Client
	currentProduct  string
	simulationState SimulationState
	stateMu         sync.RWMutex
)

type SimulationState struct {
	Running          bool                   `json:"running"`
	ProductID        string                 `json:"product_id"`
	ProductTier      string                 `json:"product_tier"`
	Configuration    SimulationConfig       `json:"configuration"`
	Metrics          SimulationMetrics      `json:"metrics"`
	Events           []SimulationEvent      `json:"events"`
	LicenseInfo      map[string]FeatureInfo `json:"license_info"`
}

type SimulationConfig struct {
	EnabledControls []string `json:"enabled_controls"`
	LoopCount       int      `json:"loop_count"`
	IntervalMs      int      `json:"interval_ms"`
}

type SimulationMetrics struct {
	CurrentIteration int     `json:"current_iteration"`
	TotalIterations  int     `json:"total_iterations"`
	SuccessCount     int     `json:"success_count"`
	FailureCount     int     `json:"failure_count"`
	RateLimitHits    int     `json:"rate_limit_hits"`
	QuotaExceeded    int     `json:"quota_exceeded"`
	APICallRate      float64 `json:"api_call_rate"`
	Elapsed          int64   `json:"elapsed"`
}

type SimulationEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Code      string    `json:"code,omitempty"`
}

type FeatureInfo struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Enabled       bool    `json:"enabled"`
	Tier          string  `json:"tier"`
	Quota         int     `json:"quota,omitempty"`
	MaxCapacity   int     `json:"max_capacity,omitempty"`
	MaxTPS        float64 `json:"max_tps,omitempty"`
	MaxConcurrency int    `json:"max_concurrency,omitempty"`
	Reason        string  `json:"reason,omitempty"`
}

type Product struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Tier        string            `json:"tier"`
	Description string            `json:"description"`
	Features    []ProductFeature  `json:"features"`
	Limitations []ProductLimit    `json:"limitations"`
}

type ProductFeature struct {
	Name        string `json:"name"`
	Available   bool   `json:"available"`
	Description string `json:"description"`
}

type ProductLimit struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9144"
	}

	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	// API endpoints
	mux.HandleFunc("/api/products", handleGetProducts)
	mux.HandleFunc("/api/products/select", handleSelectProduct)
	mux.HandleFunc("/api/simulation/configure", handleConfigureSimulation)
	mux.HandleFunc("/api/simulation/start", handleStartSimulation)
	mux.HandleFunc("/api/simulation/stop", handleStopSimulation)
	mux.HandleFunc("/api/simulation/status", handleSimulationStatus)
	mux.HandleFunc("/api/simulation/events", handleSimulationEvents)

	// Pages
	mux.HandleFunc("/", handleIndexPage)
	mux.HandleFunc("/discover", handleDiscoverPage)
	mux.HandleFunc("/configure", handleConfigurePage)
	mux.HandleFunc("/runtime", handleRuntimePage)

	addr := ":" + port
	log.Printf("LCC Web Demo starting on http://localhost%s", addr)
	log.Printf("Navigate to http://localhost%s/discover to begin", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{
		{
			ID:          "demo-app-basic",
			Name:        "Basic Edition",
			Tier:        "basic",
			Description: "Essential features for small teams and individual developers",
			Features: []ProductFeature{
				{Name: "Basic Analytics", Available: true, Description: "Standard data analysis and reporting"},
				{Name: "Local Export", Available: true, Description: "Export data to local files"},
				{Name: "Advanced Analytics", Available: false, Description: "ML-powered insights (Pro required)"},
				{Name: "PDF Export", Available: false, Description: "Professional PDF reports (Pro required)"},
				{Name: "Excel Export", Available: false, Description: "Advanced Excel exports (Enterprise required)"},
			},
			Limitations: []ProductLimit{
				{Name: "API Calls", Value: "100/day"},
				{Name: "Projects", Value: "3 max"},
				{Name: "Concurrent Users", Value: "1"},
			},
		},
		{
			ID:          "demo-app-pro",
			Name:        "Professional Edition",
			Tier:        "professional",
			Description: "Advanced features for growing teams and businesses",
			Features: []ProductFeature{
				{Name: "Basic Analytics", Available: true, Description: "Standard data analysis and reporting"},
				{Name: "Advanced Analytics", Available: true, Description: "ML-powered insights with predictive models"},
				{Name: "PDF Export", Available: true, Description: "Professional quality PDF reports"},
				{Name: "Scheduled Reports", Available: true, Description: "Automated report generation"},
				{Name: "Excel Export", Available: false, Description: "Advanced Excel exports (Enterprise required)"},
			},
			Limitations: []ProductLimit{
				{Name: "API Calls", Value: "10,000/day"},
				{Name: "PDF Exports", Value: "200/day"},
				{Name: "Projects", Value: "50 max"},
				{Name: "API Rate Limit", Value: "10 TPS"},
				{Name: "Concurrent Users", Value: "10"},
			},
		},
		{
			ID:          "demo-app-ent",
			Name:        "Enterprise Edition",
			Tier:        "enterprise",
			Description: "Full-featured solution for large organizations",
			Features: []ProductFeature{
				{Name: "All Pro Features", Available: true, Description: "Includes all Professional features"},
				{Name: "Excel Export", Available: true, Description: "Advanced Excel exports with templates"},
				{Name: "Cloud Integration", Available: true, Description: "Direct cloud storage integration"},
				{Name: "Custom Integrations", Available: true, Description: "REST API and webhooks"},
				{Name: "Priority Support", Available: true, Description: "24/7 dedicated support"},
			},
			Limitations: []ProductLimit{
				{Name: "API Calls", Value: "Unlimited"},
				{Name: "Exports", Value: "Unlimited"},
				{Name: "Projects", Value: "Unlimited"},
				{Name: "API Rate Limit", Value: "100 TPS"},
				{Name: "Concurrent Users", Value: "100"},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func handleSelectProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProductID string `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Initialize LCC SDK with selected product
	cfg := &config.SDKConfig{
		LCCURL:         "https://localhost:8088",
		ProductID:      req.ProductID,
		ProductVersion: "1.0.0",
		Timeout:        30 * time.Second,
		CacheTTL:       10 * time.Second,
	}

	var err error
	lccClient, err = client.NewClient(cfg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create client: %v", err), http.StatusInternalServerError)
		return
	}

	if err := lccClient.Register(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to register: %v", err), http.StatusInternalServerError)
		return
	}

	currentProduct = req.ProductID
	
	// Load license information
	licenseInfo := loadLicenseInfo()

	stateMu.Lock()
	simulationState.ProductID = req.ProductID
	simulationState.ProductTier = getTierFromProductID(req.ProductID)
	simulationState.LicenseInfo = licenseInfo
	stateMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success":      true,
		"product_id":   req.ProductID,
		"instance_id":  lccClient.GetInstanceID(),
		"license_info": licenseInfo,
	})
}

func handleConfigureSimulation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config SimulationConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	stateMu.Lock()
	simulationState.Configuration = config
	stateMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"config":  config,
	})
}

func handleStartSimulation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if lccClient == nil {
		http.Error(w, "No product selected", http.StatusBadRequest)
		return
	}

	stateMu.Lock()
	if simulationState.Running {
		stateMu.Unlock()
		http.Error(w, "Simulation already running", http.StatusConflict)
		return
	}
	simulationState.Running = true
	simulationState.Metrics = SimulationMetrics{TotalIterations: simulationState.Configuration.LoopCount}
	simulationState.Events = []SimulationEvent{}
	stateMu.Unlock()

	go runSimulation()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleStopSimulation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stateMu.Lock()
	simulationState.Running = false
	stateMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleSimulationStatus(w http.ResponseWriter, r *http.Request) {
	stateMu.RLock()
	state := simulationState
	stateMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func handleSimulationEvents(w http.ResponseWriter, r *http.Request) {
	stateMu.RLock()
	events := simulationState.Events
	stateMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func getTierFromProductID(productID string) string {
	if len(productID) > 4 {
		suffix := productID[len(productID)-3:]
		switch suffix {
		case "sic":
			return "basic"
		case "pro":
			return "professional"
		case "ent":
			return "enterprise"
		}
	}
	return "basic"
}

func loadLicenseInfo() map[string]FeatureInfo {
	features := []string{
		"advanced_analytics",
		"pdf_export",
		"excel_export",
		"scheduled_reports",
		"capacity.project.count",
		"api.v1.demo",
		"concurrent.user",
	}

	info := make(map[string]FeatureInfo)
	for _, featureID := range features {
		status, err := lccClient.CheckFeature(featureID)
		if err != nil {
			info[featureID] = FeatureInfo{
				ID:      featureID,
				Enabled: false,
				Reason:  err.Error(),
			}
			continue
		}

		info[featureID] = FeatureInfo{
			ID:             featureID,
			Enabled:        status.Enabled,
			Reason:         status.Reason,
			MaxCapacity:    status.MaxCapacity,
			MaxTPS:         status.MaxTPS,
			MaxConcurrency: status.MaxConcurrency,
		}
	}

	return info
}

func runSimulation() {
	startTime := time.Now()
	config := simulationState.Configuration

	addEvent("info", "Simulation started", "")

	for i := 1; i <= config.LoopCount; i++ {
		stateMu.RLock()
		if !simulationState.Running {
			stateMu.RUnlock()
			break
		}
		stateMu.RUnlock()

		stateMu.Lock()
		simulationState.Metrics.CurrentIteration = i
		simulationState.Metrics.Elapsed = int64(time.Since(startTime).Seconds())
		stateMu.Unlock()

		// Execute different scenarios based on configuration
		for _, control := range config.EnabledControls {
			switch control {
			case "rate_limit":
				executeRateLimitScenario(i)
			case "quota":
				executeQuotaScenario(i)
			case "feature_gate":
				executeFeatureGateScenario(i)
			case "capacity":
				executeCapacityScenario(i)
			}
		}

		time.Sleep(time.Duration(config.IntervalMs) * time.Millisecond)
	}

	stateMu.Lock()
	simulationState.Running = false
	stateMu.Unlock()

	addEvent("info", "Simulation completed", "")
}

func executeRateLimitScenario(iteration int) {
	allowed, _, reason, err := lccClient.Consume("pdf_export", 1, nil)
	if err != nil {
		addEvent("error", fmt.Sprintf("Iteration %d: Failed to check rate limit: %v", iteration, err), "")
		stateMu.Lock()
		simulationState.Metrics.FailureCount++
		stateMu.Unlock()
		return
	}

	if !allowed {
		addEvent("warning", fmt.Sprintf("Iteration %d: Rate limit hit - %s", iteration, reason), "RATE_LIMIT")
		stateMu.Lock()
		simulationState.Metrics.RateLimitHits++
		stateMu.Unlock()
		return
	}

	export.GeneratePDF(fmt.Sprintf("report-%d.pdf", iteration))
	addEvent("success", fmt.Sprintf("Iteration %d: PDF export succeeded", iteration), "")
	stateMu.Lock()
	simulationState.Metrics.SuccessCount++
	stateMu.Unlock()
}

func executeQuotaScenario(iteration int) {
	allowed, remaining, reason, err := lccClient.Consume("advanced_analytics", 1, nil)
	if err != nil {
		addEvent("error", fmt.Sprintf("Iteration %d: Failed to check quota: %v", iteration, err), "")
		stateMu.Lock()
		simulationState.Metrics.FailureCount++
		stateMu.Unlock()
		return
	}

	if !allowed {
		addEvent("warning", fmt.Sprintf("Iteration %d: Quota exceeded - %s", iteration, reason), "QUOTA_EXCEEDED")
		stateMu.Lock()
		simulationState.Metrics.QuotaExceeded++
		stateMu.Unlock()
		return
	}

	analytics.RunAdvanced()
	addEvent("success", fmt.Sprintf("Iteration %d: Advanced analytics completed (remaining: %d)", iteration, remaining), "")
	stateMu.Lock()
	simulationState.Metrics.SuccessCount++
	stateMu.Unlock()
}

func executeFeatureGateScenario(iteration int) {
	status, err := lccClient.CheckFeature("excel_export")
	if err != nil {
		addEvent("error", fmt.Sprintf("Iteration %d: Failed to check feature: %v", iteration, err), "")
		stateMu.Lock()
		simulationState.Metrics.FailureCount++
		stateMu.Unlock()
		return
	}

	if !status.Enabled {
		addEvent("warning", fmt.Sprintf("Iteration %d: Feature disabled - %s", iteration, status.Reason), "FEATURE_DISABLED")
		stateMu.Lock()
		simulationState.Metrics.FailureCount++
		stateMu.Unlock()
		return
	}

	export.GenerateExcel(fmt.Sprintf("report-%d.xlsx", iteration))
	addEvent("success", fmt.Sprintf("Iteration %d: Excel export succeeded", iteration), "")
	stateMu.Lock()
	simulationState.Metrics.SuccessCount++
	stateMu.Unlock()
}

func executeCapacityScenario(iteration int) {
	allowed, max, reason, err := lccClient.CheckCapacity("capacity.project.count", iteration)
	if err != nil {
		addEvent("error", fmt.Sprintf("Iteration %d: Failed to check capacity: %v", iteration, err), "")
		stateMu.Lock()
		simulationState.Metrics.FailureCount++
		stateMu.Unlock()
		return
	}

	if !allowed {
		addEvent("warning", fmt.Sprintf("Iteration %d: Capacity limit reached (max: %d) - %s", iteration, max, reason), "CAPACITY_EXCEEDED")
		stateMu.Lock()
		simulationState.Metrics.FailureCount++
		stateMu.Unlock()
		return
	}

	addEvent("success", fmt.Sprintf("Iteration %d: Project created (current: %d, max: %d)", iteration, iteration, max), "")
	stateMu.Lock()
	simulationState.Metrics.SuccessCount++
	stateMu.Unlock()
}

func addEvent(level, message, code string) {
	event := SimulationEvent{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Code:      code,
	}

	stateMu.Lock()
	simulationState.Events = append(simulationState.Events, event)
	// Keep only last 100 events
	if len(simulationState.Events) > 100 {
		simulationState.Events = simulationState.Events[len(simulationState.Events)-100:]
	}
	stateMu.Unlock()
}

// Page handlers
func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/discover", http.StatusTemporaryRedirect)
}

func handleDiscoverPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "cmd/webdemo/static/discover.html")
}

func handleConfigurePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "cmd/webdemo/static/configure.html")
}

func handleRuntimePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "cmd/webdemo/static/runtime.html")
}
