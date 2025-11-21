package web

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	lccclient "github.com/yourorg/lcc-sdk/pkg/client"
)

type SimulationStatus string

const (
	StatusIdle     SimulationStatus = "idle"
	StatusRunning  SimulationStatus = "running"
	StatusPaused   SimulationStatus = "paused"
	StatusStopped  SimulationStatus = "stopped"
	StatusError    SimulationStatus = "error"
	StatusCompleted SimulationStatus = "completed"
)

type EventType string

const (
	EventTypeStart       EventType = "simulation_start"
	EventTypeStop        EventType = "simulation_stop"
	EventTypePause       EventType = "simulation_pause"
	EventTypeResume      EventType = "simulation_resume"
	EventTypeComplete    EventType = "simulation_complete"
	EventTypeIterationStart EventType = "iteration_start"
	EventTypeFeatureCall EventType = "feature_call"
	EventTypeError       EventType = "error"
)

type SimulationEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	Type        EventType `json:"type"`
	Iteration   int       `json:"iteration,omitempty"`
	FeatureID   string    `json:"feature_id,omitempty"`
	Allowed     bool      `json:"allowed"`
	Reason      string    `json:"reason"`
	Details     string    `json:"details"`
	CallResult  map[string]interface{} `json:"call_result,omitempty"`
	Error       string    `json:"error,omitempty"`
}

type SimulationMetrics struct {
	TotalIterations    int            `json:"total_iterations"`
	CompletedIterations int           `json:"completed_iterations"`
	SuccessCount       int            `json:"success_count"`
	FailureCount       int            `json:"failure_count"`
	ElapsedSeconds     float64        `json:"elapsed_seconds"`
	EstimatedRemaining float64        `json:"estimated_remaining_seconds"`
	CurrentTPS         map[string]float64 `json:"current_tps"`
	QuotaRemaining     map[string]int `json:"quota_remaining"`
	FeatureCalls       map[string]int `json:"feature_calls"`
}

type SimulationConfig struct {
	ProductID        string `json:"product_id"`
	InstanceID       string `json:"instance_id"`
	Iterations       int    `json:"iterations"`
	IntervalMS       int    `json:"interval_ms"`
	FeaturesToCall   []string `json:"features_to_call"`
	CallPattern      map[string]int `json:"call_pattern"`
}

type SimulationEngine struct {
	mu              sync.RWMutex
	config          SimulationConfig
	client          *lccclient.Client
	status          SimulationStatus
	metrics         SimulationMetrics
	events          []SimulationEvent
	eventsChan      chan SimulationEvent
	stopChan        chan struct{}
	pauseChan       chan struct{}
	resumeChan      chan struct{}
	paused          bool
	startTime       time.Time
	pauseTime       time.Duration
	lastPauseStart  time.Time
}

func NewSimulationEngine(config SimulationConfig, client *lccclient.Client) *SimulationEngine {
	return &SimulationEngine{
		config:     config,
		client:     client,
		status:     StatusIdle,
		events:     make([]SimulationEvent, 0, 1000),
		eventsChan: make(chan SimulationEvent, 100),
		stopChan:   make(chan struct{}),
		pauseChan:  make(chan struct{}),
		resumeChan: make(chan struct{}),
		metrics: SimulationMetrics{
			TotalIterations: config.Iterations,
			CurrentTPS:      make(map[string]float64),
			QuotaRemaining:  make(map[string]int),
			FeatureCalls:    make(map[string]int),
		},
	}
}

func (e *SimulationEngine) Start(ctx context.Context) error {
	e.mu.Lock()
	if e.status == StatusRunning {
		e.mu.Unlock()
		return fmt.Errorf("simulation already running")
	}
	e.status = StatusRunning
	e.startTime = time.Now()
	e.mu.Unlock()

	go e.eventLoop(ctx)
	go e.simulationLoop(ctx)

	e.recordEvent(SimulationEvent{
		Timestamp: time.Now(),
		Type:      EventTypeStart,
		Details:   fmt.Sprintf("Starting simulation with %d iterations", e.config.Iterations),
	})

	return nil
}

func (e *SimulationEngine) Stop() error {
	e.mu.Lock()
	if e.status != StatusRunning && e.status != StatusPaused {
		e.mu.Unlock()
		return fmt.Errorf("simulation not running")
	}
	e.mu.Unlock()

	e.recordEvent(SimulationEvent{
		Timestamp: time.Now(),
		Type:      EventTypeStop,
		Details:   fmt.Sprintf("Stopped at iteration %d/%d", e.metrics.CompletedIterations, e.config.Iterations),
	})

	close(e.stopChan)
	e.setStatus(StatusStopped)

	return nil
}

func (e *SimulationEngine) Pause() error {
	e.mu.Lock()
	if e.status != StatusRunning {
		e.mu.Unlock()
		return fmt.Errorf("can only pause running simulation")
	}
	e.paused = true
	e.lastPauseStart = time.Now()
	e.mu.Unlock()

	e.recordEvent(SimulationEvent{
		Timestamp: time.Now(),
		Type:      EventTypePause,
		Details:   fmt.Sprintf("Paused at iteration %d", e.metrics.CompletedIterations),
	})

	e.setStatus(StatusPaused)
	return nil
}

func (e *SimulationEngine) Resume() error {
	e.mu.Lock()
	if e.status != StatusPaused {
		e.mu.Unlock()
		return fmt.Errorf("can only resume paused simulation")
	}
	e.paused = false
	e.pauseTime += time.Since(e.lastPauseStart)
	e.mu.Unlock()

	e.recordEvent(SimulationEvent{
		Timestamp: time.Now(),
		Type:      EventTypeResume,
		Details:   fmt.Sprintf("Resumed from iteration %d", e.metrics.CompletedIterations),
	})

	e.setStatus(StatusRunning)
	return nil
}

func (e *SimulationEngine) GetStatus() (SimulationStatus, SimulationMetrics) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	metrics := e.metrics
	if e.status == StatusRunning {
		elapsed := time.Since(e.startTime) - e.pauseTime
		metrics.ElapsedSeconds = elapsed.Seconds()
		
		if e.metrics.CompletedIterations > 0 {
			avgTime := elapsed.Seconds() / float64(e.metrics.CompletedIterations)
			remaining := float64(e.metrics.TotalIterations-e.metrics.CompletedIterations) * avgTime
			metrics.EstimatedRemaining = remaining
		}
	}

	return e.status, metrics
}

func (e *SimulationEngine) GetEvents(limit int) []SimulationEvent {
	e.mu.RLock()
	defer e.mu.RUnlock()

	events := e.events
	if limit > 0 && len(events) > limit {
		events = events[len(events)-limit:]
	}
	return events
}

func (e *SimulationEngine) simulationLoop(ctx context.Context) {
	interval := time.Duration(e.config.IntervalMS) * time.Millisecond

	for i := 1; i <= e.config.Iterations; i++ {
		select {
		case <-ctx.Done():
			e.setStatus(StatusCompleted)
			return
		case <-e.stopChan:
			return
		default:
		}

		// Handle pause
		for {
			e.mu.RLock()
			paused := e.paused
			e.mu.RUnlock()

			if !paused {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		e.runIteration(i)

		e.mu.Lock()
		e.metrics.CompletedIterations = i
		e.mu.Unlock()

		if i < e.config.Iterations {
			time.Sleep(interval)
		}
	}

	e.setStatus(StatusCompleted)
	e.recordEvent(SimulationEvent{
		Timestamp: time.Now(),
		Type:      EventTypeComplete,
		Details:   fmt.Sprintf("Simulation completed: %d/%d iterations", e.metrics.CompletedIterations, e.config.Iterations),
	})
}

func (e *SimulationEngine) runIteration(iteration int) {
	e.recordEvent(SimulationEvent{
		Timestamp: time.Now(),
		Type:      EventTypeIterationStart,
		Iteration: iteration,
		Details:   fmt.Sprintf("Starting iteration %d", iteration),
	})

	for _, featureID := range e.config.FeaturesToCall {
		// Check if should call this feature based on pattern
		pattern := e.config.CallPattern[featureID]
		if pattern > 0 && iteration%pattern != 0 {
			e.recordEvent(SimulationEvent{
				Timestamp: time.Now(),
				Type:      EventTypeFeatureCall,
				Iteration: iteration,
				FeatureID: featureID,
				Allowed:   true,
				Reason:    "skipped",
				Details:   fmt.Sprintf("Skipped based on pattern (every %d iterations)", pattern),
			})
			continue
		}

		e.callFeature(iteration, featureID)
	}
}

func (e *SimulationEngine) callFeature(iteration int, featureID string) {
	if e.client == nil {
		e.recordEvent(SimulationEvent{
			Timestamp: time.Now(),
			Type:      EventTypeError,
			Iteration: iteration,
			FeatureID: featureID,
			Allowed:   false,
			Error:     "client is nil",
		})
		return
	}

	status, err := e.client.CheckFeature(featureID)
	if err != nil {
		e.recordEvent(SimulationEvent{
			Timestamp: time.Now(),
			Type:      EventTypeFeatureCall,
			Iteration: iteration,
			FeatureID: featureID,
			Allowed:   false,
			Reason:    "error",
			Error:     err.Error(),
			Details:   fmt.Sprintf("Feature check error: %v", err),
		})
		e.mu.Lock()
		e.metrics.FailureCount++
		e.mu.Unlock()
		return
	}

	e.mu.Lock()
	if status.Enabled {
		e.metrics.SuccessCount++
	} else {
		e.metrics.FailureCount++
	}
	e.metrics.FeatureCalls[featureID]++
	if status.Quota != nil {
		e.metrics.QuotaRemaining[featureID] = status.Quota.Remaining
	}
	e.mu.Unlock()

	callResult := map[string]interface{}{
		"enabled": status.Enabled,
		"reason":  status.Reason,
	}
	if status.Quota != nil {
		callResult["quota_remaining"] = status.Quota.Remaining
		callResult["quota_limit"] = status.Quota.Limit
	}

	e.recordEvent(SimulationEvent{
		Timestamp:  time.Now(),
		Type:       EventTypeFeatureCall,
		Iteration:  iteration,
		FeatureID:  featureID,
		Allowed:    status.Enabled,
		Reason:     status.Reason,
		CallResult: callResult,
		Details:    fmt.Sprintf("CheckFeature(%s) -> %v (%s)", featureID, status.Enabled, status.Reason),
	})
}

func (e *SimulationEngine) recordEvent(event SimulationEvent) {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	e.mu.Lock()
	if len(e.events) >= 10000 {
		e.events = e.events[1:]
	}
	e.events = append(e.events, event)
	e.mu.Unlock()

	select {
	case e.eventsChan <- event:
	default:
		log.Printf("Warning: event channel full, dropping event")
	}
}

func (e *SimulationEngine) eventLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-e.stopChan:
			return
		case event := <-e.eventsChan:
			_ = event
		}
	}
}

func (e *SimulationEngine) setStatus(status SimulationStatus) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.status = status
}

type SimulationManager struct {
	mu         sync.RWMutex
	simulations map[string]*SimulationEngine
}

func NewSimulationManager() *SimulationManager {
	return &SimulationManager{
		simulations: make(map[string]*SimulationEngine),
	}
}

func (m *SimulationManager) Create(config SimulationConfig, client *lccclient.Client) *SimulationEngine {
	engine := NewSimulationEngine(config, client)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.simulations[config.InstanceID] = engine
	return engine
}

func (m *SimulationManager) Get(instanceID string) *SimulationEngine {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.simulations[instanceID]
}

func (m *SimulationManager) Delete(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.simulations, instanceID)
}
