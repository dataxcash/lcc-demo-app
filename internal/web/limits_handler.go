package web

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

func (s *Server) handleGetLimitTypes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	_ = json.NewEncoder(w).Encode(AllLimitTypes)
}

func (s *Server) handleGetLimitExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	limitType := extractLimitTypeFromPath(r.URL.Path, "/api/limits/", "/example")
	if limitType == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid limit type"))
		return
	}

	example := GetLimitExample(limitType)
	if example == nil {
		writeErr(w, http.StatusNotFound, fmt.Errorf("limit type not found"))
		return
	}

	_ = json.NewEncoder(w).Encode(example)
}

type SimulateRequest struct {
	FeatureID  string                 `json:"feature_id"`
	Iterations int                    `json:"iterations"`
	Params     map[string]interface{} `json:"params"`
}

type SimulationResult struct {
	Iteration int    `json:"iteration"`
	Allowed   bool   `json:"allowed"`
	Remaining string `json:"remaining"`
	Reason    string `json:"reason"`
	Details   string `json:"details,omitempty"`
}

type SimulateResponse struct {
	Success bool               `json:"success"`
	Type    string             `json:"type"`
	Results []SimulationResult `json:"results"`
	Summary string             `json:"summary"`
}

func (s *Server) handleSimulateLimitType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	limitType := extractLimitTypeFromPath(r.URL.Path, "/api/limits/", "/simulate")
	if limitType == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid limit type"))
		return
	}

	var req SimulateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
		return
	}

	if req.Iterations <= 0 {
		req.Iterations = 10
	}
	if req.Iterations > 100 {
		req.Iterations = 100
	}

	resp := simulateLimitType(limitType, req)
	_ = json.NewEncoder(w).Encode(resp)
}

func simulateLimitType(limitType string, req SimulateRequest) SimulateResponse {
	switch limitType {
	case "quota":
		return simulateQuota(req)
	case "tps":
		return simulateTPS(req)
	case "capacity":
		return simulateCapacity(req)
	case "concurrency":
		return simulateConcurrency(req)
	default:
		return SimulateResponse{
			Success: false,
			Type:    limitType,
			Summary: "Unknown limit type",
		}
	}
}

func simulateQuota(req SimulateRequest) SimulateResponse {
	maxQuota := 10000
	if max, ok := req.Params["max"].(float64); ok {
		maxQuota = int(max)
	}

	consumed := 0
	results := make([]SimulationResult, 0, req.Iterations)
	successCount := 0

	for i := 1; i <= req.Iterations; i++ {
		amount := 1
		if amt, ok := req.Params["amount"].(float64); ok {
			amount = int(amt)
		}

		allowed := (consumed + amount) <= maxQuota
		if allowed {
			consumed += amount
			successCount++
		}

		remaining := maxQuota - consumed
		if remaining < 0 {
			remaining = 0
		}

		reason := "ok"
		if !allowed {
			reason = "exceeded"
		}

		results = append(results, SimulationResult{
			Iteration: i,
			Allowed:   allowed,
			Remaining: fmt.Sprintf("%d", remaining),
			Reason:    reason,
			Details:   fmt.Sprintf("consumed=%d/%d", consumed, maxQuota),
		})
	}

	summary := fmt.Sprintf("Completed %d iterations. Success: %d, Failed: %d", 
		req.Iterations, successCount, req.Iterations-successCount)

	return SimulateResponse{
		Success: true,
		Type:    "quota",
		Results: results,
		Summary: summary,
	}
}

func simulateTPS(req SimulateRequest) SimulateResponse {
	maxTPS := 10.0
	if max, ok := req.Params["max_tps"].(float64); ok {
		maxTPS = max
	}

	results := make([]SimulationResult, 0, req.Iterations)
	successCount := 0

	for i := 1; i <= req.Iterations; i++ {
		currentTPS := rand.Float64() * maxTPS * 1.5

		allowed := currentTPS <= maxTPS
		if allowed {
			successCount++
		}

		reason := "ok"
		if !allowed {
			reason = "exceeded"
		}

		results = append(results, SimulationResult{
			Iteration: i,
			Allowed:   allowed,
			Remaining: fmt.Sprintf("max=%.1f", maxTPS),
			Reason:    reason,
			Details:   fmt.Sprintf("current_tps=%.2f", currentTPS),
		})
	}

	summary := fmt.Sprintf("Completed %d iterations. Success: %d, Failed: %d", 
		req.Iterations, successCount, req.Iterations-successCount)

	return SimulateResponse{
		Success: true,
		Type:    "tps",
		Results: results,
		Summary: summary,
	}
}

func simulateCapacity(req SimulateRequest) SimulateResponse {
	maxCapacity := 50
	if max, ok := req.Params["max_capacity"].(float64); ok {
		maxCapacity = int(max)
	}

	currentCount := 0
	results := make([]SimulationResult, 0, req.Iterations)
	successCount := 0

	for i := 1; i <= req.Iterations; i++ {
		action := "create"
		if i > req.Iterations/2 && currentCount > 0 && rand.Float64() > 0.6 {
			action = "delete"
		}

		var allowed bool
		var reason string

		if action == "create" {
			allowed = currentCount < maxCapacity
			if allowed {
				currentCount++
				successCount++
				reason = "ok"
			} else {
				reason = "at_limit"
			}
		} else {
			currentCount--
			if currentCount < 0 {
				currentCount = 0
			}
			allowed = true
			reason = "deleted"
		}

		results = append(results, SimulationResult{
			Iteration: i,
			Allowed:   allowed,
			Remaining: fmt.Sprintf("max=%d", maxCapacity),
			Reason:    reason,
			Details:   fmt.Sprintf("current=%d, action=%s", currentCount, action),
		})
	}

	summary := fmt.Sprintf("Completed %d iterations. Success: %d, Failed: %d, Final count: %d", 
		req.Iterations, successCount, req.Iterations-successCount, currentCount)

	return SimulateResponse{
		Success: true,
		Type:    "capacity",
		Results: results,
		Summary: summary,
	}
}

func simulateConcurrency(req SimulateRequest) SimulateResponse {
	maxSlots := 10
	if max, ok := req.Params["max_concurrency"].(float64); ok {
		maxSlots = int(max)
	}

	currentSlots := 0
	results := make([]SimulationResult, 0, req.Iterations)
	successCount := 0

	for i := 1; i <= req.Iterations; i++ {
		action := "acquire"
		if currentSlots > 0 && rand.Float64() > 0.6 {
			action = "release"
		}

		var allowed bool
		var reason string

		if action == "acquire" {
			allowed = currentSlots < maxSlots
			if allowed {
				currentSlots++
				successCount++
				reason = "ok"
			} else {
				reason = "max_reached"
			}
		} else {
			currentSlots--
			if currentSlots < 0 {
				currentSlots = 0
			}
			allowed = true
			reason = "released"
		}

		freeSlots := maxSlots - currentSlots

		results = append(results, SimulationResult{
			Iteration: i,
			Allowed:   allowed,
			Remaining: fmt.Sprintf("%d free", freeSlots),
			Reason:    reason,
			Details:   fmt.Sprintf("slots=%d/%d, action=%s", currentSlots, maxSlots, action),
		})
	}

	summary := fmt.Sprintf("Completed %d iterations. Acquired: %d, Rejected: %d, Active slots: %d", 
		req.Iterations, successCount, req.Iterations-successCount, currentSlots)

	return SimulateResponse{
		Success: true,
		Type:    "concurrency",
		Results: results,
		Summary: summary,
	}
}

func extractLimitTypeFromPath(path, prefix, suffix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	if !strings.HasSuffix(path, suffix) {
		return ""
	}

	path = strings.TrimPrefix(path, prefix)
	path = strings.TrimSuffix(path, suffix)

	return path
}
