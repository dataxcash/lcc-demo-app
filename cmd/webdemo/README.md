# LCC Web Demo - Interactive License Control Demonstration

An interactive web-based demonstration of the LCC SDK, showcasing license control capabilities through a three-stage workflow.

## Overview

This web demo provides a comprehensive, visual way to understand how the LCC SDK integrates into real applications:

1. **Discover Products** - Compare licensing tiers and their features
2. **Configure Simulation** - Set up license control scenarios with real code examples
3. **Runtime Dashboard** - Monitor live simulation with real-time metrics

## Features

### 🎯 Three-Stage Workflow

#### Stage 1: Product Discovery
- **Interactive Product Catalog**: Browse Basic, Professional, and Enterprise editions
- **Feature Comparison**: See available features and limitations for each tier
- **License Preview**: View JSON configuration for each product
- **Product Selection**: Automatically registers and loads the selected license

#### Stage 2: Simulation Configuration
- **Control Type Selection**:
  - Rate Limiting: API call frequency control
  - Quota Management: Consumption tracking against limits
  - Feature Gating: Tier-based access control
  - Capacity Control: Resource limits (projects, users)
  
- **Real Business Code Examples**:
  - `ProcessDataAnalytics()`: Shows quota control with `Consume()` API
  - `ExportReport()`: Demonstrates rate limiting
  - `CreateProject()`: Illustrates capacity checks
  
- **Configuration Impact Analysis**: Real-time preview of how settings affect behavior

#### Stage 3: Runtime Dashboard
- **Live Metrics**:
  - Progress tracking (current/total iterations)
  - Success/failure counters
  - Rate limit hits
  - Elapsed time
  
- **Visualization**:
  - Real-time line chart of success/failure trends
  - Dynamic code context showing active SDK calls
  
- **Event Log**: Chronological view of all license check operations

## Quick Start

### Prerequisites

1. **LCC Server Running**
   ```bash
   cd /home/fila/jqdDev_2025/lcc
   ./lcc_server
   ```

2. **License Files Configured**
   - Place appropriate `.lic` files in `~/.lcc/` or configure via environment

### Running the Web Demo

```bash
# Build and run
make demo

# Or manually
go build -o bin/webdemo ./cmd/webdemo
./bin/webdemo
```

### Access the Demo

Open your browser and navigate to:
```
http://localhost:9144/discover
```

The demo will guide you through:
1. Selecting a product tier
2. Configuring simulation parameters
3. Running and monitoring the simulation

### Custom Port

```bash
PORT=9145 ./bin/webdemo
```

## Architecture

### Backend (Go)

```
cmd/webdemo/main.go
├── HTTP Server (port 9144)
├── API Endpoints
│   ├── GET  /api/products
│   ├── POST /api/products/select
│   ├── POST /api/simulation/configure
│   ├── POST /api/simulation/start
│   ├── POST /api/simulation/stop
│   ├── GET  /api/simulation/status
│   └── GET  /api/simulation/events
├── Simulation Engine
│   ├── Rate Limit Scenario
│   ├── Quota Scenario
│   ├── Feature Gate Scenario
│   └── Capacity Scenario
└── LCC SDK Integration
    ├── client.Consume()
    ├── client.CheckFeature()
    ├── client.CheckCapacity()
    └── client.CheckTPS()
```

### Frontend (Vanilla JS + HTML/CSS)

```
cmd/webdemo/static/
├── discover.html    - Product catalog and selection
├── configure.html   - Simulation setup with code examples
└── runtime.html     - Live dashboard with Chart.js
```

## Demo Scenarios

### Scenario 1: Professional Tier Rate Limiting

1. Select "Professional Edition"
2. Enable "Rate Limiting" control
3. Set Loop Count: 250, Interval: 500ms
4. Start simulation
5. Observe: Rate limit triggers after ~200 PDF exports

**Code Context:**
```go
allowed, _, reason, err := lccClient.Consume(
    "pdf_export",
    1,
    metadata,
)
if !allowed {
    return ErrDailyQuotaExceeded
}
```

### Scenario 2: Quota Exhaustion

1. Select "Professional Edition"
2. Enable "Quota Management"
3. Set Loop Count: 1100, Interval: 200ms
4. Start simulation
5. Observe: Quota exceeded after ~1000 advanced analytics calls

**Code Context:**
```go
allowed, remaining, reason, err := lccClient.Consume(
    "advanced_analytics",
    1,
    nil,
)
if !allowed {
    log.Warn("Quota exceeded", "reason", reason)
}
```

### Scenario 3: Feature Gate Blocking

1. Select "Basic Edition"
2. Enable "Feature Gating"
3. Start simulation
4. Observe: Excel export denied (requires Enterprise)

**Code Context:**
```go
status, err := lccClient.CheckFeature("excel_export")
if !status.Enabled {
    return ErrFeatureDisabled
}
```

### Scenario 4: Capacity Limits

1. Select "Professional Edition"
2. Enable "Capacity Control"
3. Set Loop Count: 60
4. Start simulation
5. Observe: Project creation blocked after reaching 50 (Pro tier limit)

**Code Context:**
```go
allowed, max, reason, err := lccClient.CheckCapacity(
    "capacity.project.count",
    currentCount + 1,
)
if !allowed {
    return fmt.Errorf("limit reached: %d/%d", currentCount, max)
}
```

## Product Tiers

### Basic Edition
- **API Calls**: 100/day
- **Projects**: 3 max
- **Concurrent Users**: 1
- **Features**: Basic analytics, local export only

### Professional Edition
- **API Calls**: 10,000/day
- **PDF Exports**: 200/day
- **Projects**: 50 max
- **API Rate Limit**: 10 TPS
- **Concurrent Users**: 10
- **Features**: Advanced analytics, PDF export, scheduled reports

### Enterprise Edition
- **API Calls**: Unlimited
- **Exports**: Unlimited
- **Projects**: Unlimited
- **API Rate Limit**: 100 TPS
- **Concurrent Users**: 100
- **Features**: All Pro features + Excel export + cloud integration

## Development

### Code Structure

```go
// Product selection initializes LCC client
cfg := &config.SDKConfig{
    LCCURL:         "https://localhost:8088",
    ProductID:      selectedProduct,
    ProductVersion: "1.0.0",
}
lccClient, _ = client.NewClient(cfg)
lccClient.Register()

// Simulation runs scenarios based on configuration
func runSimulation() {
    for i := 1; i <= config.LoopCount; i++ {
        for _, control := range config.EnabledControls {
            switch control {
            case "rate_limit":
                executeRateLimitScenario(i)
            case "quota":
                executeQuotaScenario(i)
            // ...
            }
        }
    }
}
```

### Adding New Scenarios

1. Add scenario function in `main.go`:
   ```go
   func executeCustomScenario(iteration int) {
       // Your LCC SDK calls here
   }
   ```

2. Add to simulation runner:
   ```go
   case "custom":
       executeCustomScenario(i)
   ```

3. Update frontend configuration UI

### Extending the Dashboard

- Modify `runtime.html` to add new metrics
- Update `SimulationMetrics` struct in `main.go`
- Add Chart.js visualizations as needed

## Troubleshooting

### "No product selected" Error
- Ensure you selected a product on the Discover page
- Check that LCC server is running and accessible

### License Registration Failed
- Verify LCC server is running: `curl http://localhost:8088/health`
- Check license files exist in `~/.lcc/`
- Review server logs for authentication errors

### Simulation Not Starting
- Open browser console (F12) for JavaScript errors
- Check backend logs for Go errors
- Verify configuration was saved (check /api/simulation/status)

### Metrics Not Updating
- Check browser console for fetch errors
- Verify simulation is running (status badge should show "Running")
- Ensure no CORS issues (frontend and backend on same origin)

## Educational Value

This demo is designed to teach:

1. **SDK Integration Patterns**
   - How to initialize and register with LCC
   - When to call different SDK APIs (Consume vs CheckFeature vs CheckCapacity)
   - Error handling and fallback strategies

2. **License Control Models**
   - Consumption-based (quota tracking)
   - Rate limiting (TPS/QPM controls)
   - Feature gating (tier-based access)
   - State capacity (resource limits)

3. **Real-World Application**
   - Code examples show actual business functions
   - Demonstrates where to place license checks
   - Shows impact of license decisions on user experience

## License

MIT License - Educational demonstration purposes.
