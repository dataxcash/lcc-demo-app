package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type StartSimulationRequest struct {
	InstanceID   string            `json:"instance_id"`
	Iterations   int               `json:"iterations"`
	IntervalMS   int               `json:"interval_ms"`
	FeaturesToCall []string        `json:"features_to_call"`
	CallPattern  map[string]int    `json:"call_pattern"`
}

type StartSimulationResponse struct {
	Success    bool   `json:"success"`
	InstanceID string `json:"instance_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
}

type StatusResponse struct {
	Success bool                `json:"success"`
	Status  string              `json:"status"`
	Metrics SimulationMetrics   `json:"metrics"`
	Error   string              `json:"error,omitempty"`
}

type EventsResponse struct {
	Success bool                `json:"success"`
	Events  []SimulationEvent   `json:"events"`
	Count   int                 `json:"count"`
	Error   string              `json:"error,omitempty"`
}

type ExportResponse struct {
	Success  bool                    `json:"success"`
	Summary  map[string]interface{}  `json:"summary"`
	Events   []SimulationEvent       `json:"events"`
	Metrics  SimulationMetrics       `json:"metrics"`
	Error    string                  `json:"error,omitempty"`
}

var simManager = NewSimulationManager()

func (s *Server) handleSimulationStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req StartSimulationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid json: %v", err),
		})
		return
	}

	if req.InstanceID == "" {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   "instance_id is required",
		})
		return
	}

	if req.Iterations <= 0 {
		req.Iterations = 100
	}
	if req.IntervalMS <= 0 {
		req.IntervalMS = 500
	}

	cli, err := s.getClient(req.InstanceID)
	if err != nil {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   fmt.Sprintf("instance not found: %v", err),
		})
		return
	}

	config := SimulationConfig{
		ProductID:      req.InstanceID,
		InstanceID:     req.InstanceID,
		Iterations:     req.Iterations,
		IntervalMS:     req.IntervalMS,
		FeaturesToCall: req.FeaturesToCall,
		CallPattern:    req.CallPattern,
	}

	engine := simManager.Create(config, cli)
	if engine == nil {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   "failed to create simulation engine",
		})
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if err := engine.Start(ctx); err != nil {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to start simulation: %v", err),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
		Success:    true,
		InstanceID: req.InstanceID,
		Status:     "running",
		Message:    "Simulation started successfully",
	})
}

func (s *Server) handleSimulationStop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	instanceID := r.URL.Query().Get("instance_id")
	if instanceID == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("instance_id is required"))
		return
	}

	engine := simManager.Get(instanceID)
	if engine == nil {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   "simulation not found",
		})
		return
	}

	if err := engine.Stop(); err != nil {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to stop simulation: %v", err),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
		Success:    true,
		InstanceID: instanceID,
		Status:     "stopped",
		Message:    "Simulation stopped",
	})
}

func (s *Server) handleSimulationPause(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	instanceID := r.URL.Query().Get("instance_id")
	if instanceID == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("instance_id is required"))
		return
	}

	engine := simManager.Get(instanceID)
	if engine == nil {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   "simulation not found",
		})
		return
	}

	if err := engine.Pause(); err != nil {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to pause: %v", err),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
		Success:    true,
		InstanceID: instanceID,
		Status:     "paused",
		Message:    "Simulation paused",
	})
}

func (s *Server) handleSimulationResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	instanceID := r.URL.Query().Get("instance_id")
	if instanceID == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("instance_id is required"))
		return
	}

	engine := simManager.Get(instanceID)
	if engine == nil {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   "simulation not found",
		})
		return
	}

	if err := engine.Resume(); err != nil {
		_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to resume: %v", err),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(&StartSimulationResponse{
		Success:    true,
		InstanceID: instanceID,
		Status:     "running",
		Message:    "Simulation resumed",
	})
}

func (s *Server) handleSimulationStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	instanceID := r.URL.Query().Get("instance_id")
	if instanceID == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("instance_id is required"))
		return
	}

	engine := simManager.Get(instanceID)
	if engine == nil {
		_ = json.NewEncoder(w).Encode(&StatusResponse{
			Success: false,
			Error:   "simulation not found",
		})
		return
	}

	status, metrics := engine.GetStatus()

	_ = json.NewEncoder(w).Encode(&StatusResponse{
		Success: true,
		Status:  string(status),
		Metrics: metrics,
	})
}

func (s *Server) handleSimulationEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	instanceID := r.URL.Query().Get("instance_id")
	if instanceID == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("instance_id is required"))
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	typeFilter := r.URL.Query().Get("type")

	engine := simManager.Get(instanceID)
	if engine == nil {
		_ = json.NewEncoder(w).Encode(&EventsResponse{
			Success: false,
			Events:  []SimulationEvent{},
			Error:   "simulation not found",
		})
		return
	}

	events := engine.GetEvents(limit)

	// Apply type filter if specified
	if typeFilter != "" {
		filtered := make([]SimulationEvent, 0, len(events))
		for _, e := range events {
			switch typeFilter {
			case "success":
				if e.Allowed && e.Type == EventTypeFeatureCall {
					filtered = append(filtered, e)
				}
			case "error":
				if !e.Allowed || e.Type == EventTypeError {
					filtered = append(filtered, e)
				}
			case "all":
				filtered = append(filtered, e)
			default:
				if string(e.Type) == typeFilter {
					filtered = append(filtered, e)
				}
			}
		}
		events = filtered
	}

	_ = json.NewEncoder(w).Encode(&EventsResponse{
		Success: true,
		Events:  events,
		Count:   len(events),
	})
}

func (s *Server) handleSimulationExport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	instanceID := r.URL.Query().Get("instance_id")
	if instanceID == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("instance_id is required"))
		return
	}

	engine := simManager.Get(instanceID)
	if engine == nil {
		_ = json.NewEncoder(w).Encode(&ExportResponse{
			Success: false,
			Error:   "simulation not found",
		})
		return
	}

	status, metrics := engine.GetStatus()
	events := engine.GetEvents(10000)

	summary := map[string]interface{}{
		"instance_id":      instanceID,
		"status":           string(status),
		"total_iterations": metrics.TotalIterations,
		"completed":        metrics.CompletedIterations,
		"success_count":    metrics.SuccessCount,
		"failure_count":    metrics.FailureCount,
		"success_rate":     float64(0),
		"elapsed_seconds":  metrics.ElapsedSeconds,
	}

	if metrics.CompletedIterations > 0 {
		successRate := float64(metrics.SuccessCount) / float64(metrics.SuccessCount+metrics.FailureCount) * 100
		summary["success_rate"] = successRate
	}

	_ = json.NewEncoder(w).Encode(&ExportResponse{
		Success: true,
		Summary: summary,
		Events:  events,
		Metrics: metrics,
	})
}

// handleSimulationRoot dispatches to specific simulation handlers
func (s *Server) handleSimulationRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path
	const prefix = "/api/simulation/"

	if !strings.HasPrefix(path, prefix) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rest := strings.TrimPrefix(path, prefix)
	parts := strings.Split(rest, "/")

	if len(parts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	action := parts[0]

	switch action {
	case "start":
		s.handleSimulationStart(w, r)
	case "stop":
		s.handleSimulationStop(w, r)
	case "pause":
		s.handleSimulationPause(w, r)
	case "resume":
		s.handleSimulationResume(w, r)
	case "status":
		s.handleSimulationStatus(w, r)
	case "events":
		s.handleSimulationEvents(w, r)
	case "export":
		s.handleSimulationExport(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
