# LCC Web Demo - Feature Showcase

## 🎯 Core Features

### 1. Product Discovery Page
**URL**: `http://localhost:9144/discover`

```
┌────────────────────────────────────────────────┐
│  🚀 License Control Center - Product Catalog  │
├────────────────────────────────────────────────┤
│                                                │
│  [Basic] [Professional] [Enterprise]          │
│  ┌──────────────────────────────────────────┐ │
│  │ Professional Edition                     │ │
│  │ Tier: PROFESSIONAL                       │ │
│  │                                          │ │
│  │ Features & Capabilities:                 │ │
│  │ ✓ Basic Analytics                        │ │
│  │ ✓ Advanced Analytics (ML-powered)        │ │
│  │ ✓ PDF Export (Professional quality)      │ │
│  │ ✓ Scheduled Reports                      │ │
│  │ ✗ Excel Export (Enterprise required)     │ │
│  │                                          │ │
│  │ Limitations & Quotas:                    │ │
│  │ API Calls:      10,000/day               │ │
│  │ PDF Exports:    200/day                  │ │
│  │ Projects:       50 max                   │ │
│  │ API Rate:       10 TPS                   │ │
│  │ Concurrent:     10 users                 │ │
│  │                                          │ │
│  │ License Configuration:                   │ │
│  │ {                                        │ │
│  │   "product_id": "demo-app-pro",          │ │
│  │   "tier": "professional",                │ │
│  │   "features": { ... }                    │ │
│  │ }                                        │ │
│  │                                          │ │
│  │  [Compare All]  [Select & Configure →]  │ │
│  └──────────────────────────────────────────┘ │
└────────────────────────────────────────────────┘
```

**Capabilities**:
- Tab-based product comparison
- Feature availability matrix
- Limitation visualization
- License JSON preview
- One-click product selection

---

### 2. Configuration Designer
**URL**: `http://localhost:9144/configure`

```
┌──────────────────────────────────────────────────────────────┐
│  ⚙️ Simulation Designer                                      │
├──────────────────────┬───────────────────────────────────────┤
│ Configuration        │ Code Integration Preview              │
│                      │                                       │
│ License Controls:    │ // ProcessDataAnalytics               │
│ ☑ Rate Limiting      │ func ProcessData(...) error {         │
│ ☑ Quota Management   │   allowed, remaining, reason, err :=  │
│ ☐ Feature Gating     │     lccClient.Consume(                │
│ ☐ Capacity Control   │       "advanced_analytics",           │
│                      │       1, nil,                         │
│ Runtime:             │     )                                 │
│ Loops:    [100]      │                                       │
│ Interval: [500ms]    │   if !allowed {                       │
│                      │     return ErrQuotaExceeded           │
│ [← Back] [Start →]   │   }                                   │
│                      │                                       │
│                      │   // Execute business logic           │
│                      │   analytics.RunMLModel(dataset)       │
│                      │ }                                     │
│                      │                                       │
│                      │ ⚠️ Configuration Impact:              │
│                      │ Rate limit triggers after ~200        │
│                      │ iterations (Pro tier quota)           │
└──────────────────────┴───────────────────────────────────────┘
```

**Capabilities**:
- Interactive control selection
- Real-time parameter adjustment
- SDK code examples in context
- Impact prediction
- Configuration validation

---

### 3. Runtime Dashboard
**URL**: `http://localhost:9144/runtime`

```
┌────────────────────────────────────────────────────────────┐
│  📊 Live Simulation Dashboard                [Running]     │
├────────────────────────────────────────────────────────────┤
│  Progress: 47/100  |  Success: 45  |  Failures: 2         │
│  Rate Limits: 2    |  Elapsed: 23s                         │
├──────────────────────────────┬─────────────────────────────┤
│ Metrics & Visualization      │ Event Log                   │
│                              │                             │
│  ┌─────────────────────┐     │ [13:45:23] ✓ Iteration 47  │
│  │ Success/Failure     │     │   PDF export succeeded      │
│  │                     │     │                             │
│  │      /\  /\         │     │ [13:45:22] ⚠ Iteration 46  │
│  │     /  \/  \  /\    │     │   Rate limit hit            │
│  │    /        \/  \   │     │   RATE_LIMIT                │
│  │ __/              \_ │     │                             │
│  └─────────────────────┘     │ [13:45:22] ✓ Iteration 45  │
│                              │   Advanced analytics OK     │
│  Current SDK Calls:          │   remaining: 955            │
│  allowed, _, reason, err :=  │                             │
│    lccClient.Consume(...)    │ [13:45:21] ✓ Iteration 44  │
│  if !allowed {               │   PDF export succeeded      │
│    return ErrLimitExceeded   │                             │
│  }                           │ [13:45:21] ✓ Iteration 43  │
│                              │   Advanced analytics OK     │
│  [← Configure] [Stop]        │   remaining: 956            │
└──────────────────────────────┴─────────────────────────────┘
```

**Capabilities**:
- Real-time metric display
- Live chart updates
- Code context visualization
- Event log streaming
- Simulation control

---

## 📊 Simulation Scenarios

### Rate Limiting Demo
```
Iteration 1-199:   ✓ Success (within quota)
Iteration 200:     ⚠ Rate limit hit
Iteration 201+:    ✗ All denied (quota exhausted)
```

### Quota Management Demo
```
Iteration 1:       ✓ Success (remaining: 999/1000)
Iteration 100:     ✓ Success (remaining: 900/1000)
Iteration 1000:    ✓ Success (remaining: 0/1000)
Iteration 1001:    ⚠ Quota exceeded
```

### Feature Gating Demo
```
Product: Basic
Feature: Excel Export
Result:  ✗ Feature disabled (requires Enterprise)
```

### Capacity Control Demo
```
Iteration 1-50:    ✓ Project created (max: 50)
Iteration 51:      ⚠ Capacity limit reached
Iteration 52+:     ✗ All denied (max capacity)
```

---

## 🎨 UI/UX Highlights

### Color Coding
- **Success**: 🟢 Green (#6ee7b7) - Operations succeeded
- **Warning**: 🟡 Yellow (#fbbf24) - Rate limits hit
- **Error**: 🔴 Red (#f87171) - Operations failed
- **Info**: 🔵 Blue (#60a5fa) - General information

### Status Indicators
- **Running**: Green badge with spinner
- **Stopped**: Red badge
- **Idle**: Blue badge

### Animations
- Page transitions: Fade in with slide
- Chart updates: Smooth line interpolation
- Event log: Auto-scroll to latest
- Button hovers: Lift effect with shadow

---

## 🔌 API Endpoints

### Product Management
```
GET  /api/products
Response: [
  {
    "id": "demo-app-basic",
    "name": "Basic Edition",
    "tier": "basic",
    "features": [...],
    "limitations": [...]
  },
  ...
]

POST /api/products/select
Request: { "product_id": "demo-app-pro" }
Response: {
  "success": true,
  "instance_id": "uuid-1234",
  "license_info": { ... }
}
```

### Simulation Control
```
POST /api/simulation/configure
Request: {
  "enabled_controls": ["rate_limit", "quota"],
  "loop_count": 100,
  "interval_ms": 500
}

POST /api/simulation/start
POST /api/simulation/stop

GET  /api/simulation/status
Response: {
  "running": true,
  "metrics": {
    "current_iteration": 47,
    "success_count": 45,
    "failure_count": 2,
    ...
  },
  "events": [...]
}
```

---

## 💻 Code Examples in Demo

### Consume() API - Quota Control
```go
// Business Function: Process expensive analytics
func ProcessDataAnalytics(ctx context.Context, dataset *Dataset) error {
    // Check license BEFORE expensive operation
    allowed, remaining, reason, err := lccClient.Consume(
        "advanced_analytics",  // Feature ID from license
        1,                     // Consume 1 credit
        nil,                   // Optional metadata
    )
    
    if err != nil {
        return fmt.Errorf("license check failed: %w", err)
    }
    
    if !allowed {
        // Denied - log and return error
        log.Warn("Analytics denied", "reason", reason)
        return ErrFeatureNotAvailable
    }
    
    // License OK - proceed with business logic
    result := analytics.RunMLModel(dataset)
    log.Info("Analytics completed", "remaining", remaining)
    
    return nil
}
```

### CheckFeature() API - Feature Gating
```go
// Business Function: Export to Excel
func ExportToExcel(reportID string) error {
    // Check if feature is enabled for this tier
    status, err := lccClient.CheckFeature("excel_export")
    if err != nil {
        return fmt.Errorf("feature check failed: %w", err)
    }
    
    if !status.Enabled {
        // Feature not available in current tier
        return fmt.Errorf(
            "Excel export requires %s tier (reason: %s)",
            status.RequiredTier,
            status.Reason,
        )
    }
    
    // Feature enabled - generate Excel
    export.GenerateExcel(reportID)
    
    // Report usage for analytics
    lccClient.ReportUsage("excel_export", 1)
    
    return nil
}
```

### CheckCapacity() API - State Limits
```go
// Business Function: Create new project
func CreateProject(name string) (*Project, error) {
    currentCount := db.CountProjects()
    
    // Check if we're within capacity limits
    allowed, maxCapacity, reason, err := lccClient.CheckCapacity(
        "capacity.project.count",
        currentCount + 1,  // Desired new count
    )
    
    if err != nil {
        return nil, fmt.Errorf("capacity check failed: %w", err)
    }
    
    if !allowed {
        // Capacity limit reached
        return nil, fmt.Errorf(
            "project limit reached: %d/%d (%s)",
            currentCount, maxCapacity, reason,
        )
    }
    
    // Within limits - create project
    project := &Project{Name: name}
    db.Save(project)
    
    return project, nil
}
```

---

## 📈 Real-Time Updates

### Polling Strategy
```javascript
// Update every 500ms for smooth experience
setInterval(async () => {
    const response = await fetch('/api/simulation/status');
    const state = await response.json();
    
    // Update metrics
    updateMetrics(state.metrics);
    
    // Update chart (every 5 iterations)
    if (state.metrics.current_iteration % 5 === 0) {
        updateChart(state.metrics);
    }
    
    // Refresh event log
    updateEventLog();
}, 500);
```

### Chart Configuration
```javascript
new Chart(ctx, {
    type: 'line',
    data: {
        labels: iterations,
        datasets: [
            { label: 'Success', data: successData, color: '#6ee7b7' },
            { label: 'Failures', data: failureData, color: '#f87171' }
        ]
    }
});
```

---

## 🚀 Performance

- **Backend**: Go HTTP server, sub-millisecond response times
- **Frontend**: Vanilla JS, no framework overhead
- **Updates**: 500ms polling interval (2 requests/sec)
- **Chart**: Updates every 5 iterations (reduces render overhead)
- **Event Log**: Last 20 events only (memory efficient)
- **Static Assets**: Embedded in binary (no external dependencies)

---

## 📚 Educational Impact

### For Developers
- **Pattern Recognition**: See where to place license checks
- **Error Handling**: Learn graceful degradation strategies
- **API Usage**: Understand Consume vs CheckFeature vs CheckCapacity
- **Metadata**: See how to attach context to license calls

### For Product Managers
- **Tier Differentiation**: Visual comparison of product levels
- **Feature Gating**: Understand license-based feature control
- **Quota Management**: See how consumption tracking works
- **User Experience**: Observe behavior when limits are hit

### For Sales Teams
- **Live Demonstration**: Show product capabilities interactively
- **Value Proposition**: Visualize differences between tiers
- **Scalability**: Demonstrate unlimited plans vs limited
- **Integration**: Prove ease of SDK integration

---

**Built with ❤️ for the LCC SDK ecosystem**
