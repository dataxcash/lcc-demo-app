# LCC Demo Application

A complete example application demonstrating integration with [lcc-sdk](https://github.com/yourorg/lcc-sdk).

This demo shows:
- ✅ Feature-based license control
- ✅ Automatic quota management
- ✅ Graceful degradation
- ✅ Tier-based authorization
- ✅ Usage reporting

## Architecture

```
Demo App
├── Analytics Module
│   ├── AdvancedAnalytics() [Professional]
│   └── BasicAnalytics() [Basic]
├── Export Module
│   ├── ExportToCloud() [Enterprise]
│   └── ExportToLocal() [Basic]
└── Reporting Module
    ├── PremiumReports() [Professional]
    └── StandardReports() [Basic]
```

## Prerequisites

1. **LCC Server running**
   ```bash
   # Start LCC server (from lcc repository)
   cd /home/fila/jqdDev_2025/lcc
   ./lcc_server
   ```

2. **Valid license file**
   - Place license file in `~/.lcc/license.lic`
   - Or configure path in environment variable

3. **Go 1.21+**

## Quick Start

### 1. Clone and Setup

```bash
cd /home/fila/jqdDev_2025/lcc-demo-app
go mod download
```

### 2. Run Demo

```bash
# Run with default config
make run

# Run in development mode (all features unlocked)
make run-dev

# Run specific demo scenario
go run cmd/demo/main.go --scenario advanced
```

## Demo Scenarios

### Scenario 1: Basic Tier (Free)

```bash
export LCC_LICENSE_TIER=basic
go run cmd/demo/main.go

# Output:
# ✓ BasicAnalytics: Success
# ✗ AdvancedAnalytics: Fallback to BasicAnalytics (tier insufficient)
# ✓ ExportToLocal: Success
# ✗ ExportToCloud: Denied (requires enterprise tier)
```

### Scenario 2: Professional Tier

```bash
export LCC_LICENSE_TIER=professional
go run cmd/demo/main.go

# Output:
# ✓ AdvancedAnalytics: Success (ML models enabled)
# ✓ PremiumReports: Success
# ✗ ExportToCloud: Denied (requires enterprise tier)
```

### Scenario 3: Enterprise Tier

```bash
export LCC_LICENSE_TIER=enterprise
go run cmd/demo/main.go

# Output:
# ✓ All features enabled
# ✓ Cloud export available
# ✓ Advanced ML models
```

### Scenario 4: Quota Limits

```bash
# AdvancedAnalytics has 1000/day quota
for i in {1..1001}; do
  curl http://localhost:8080/analytics/advanced
done

# After 1000 calls:
# ⚠ Quota exceeded, fallback to BasicAnalytics
```

## Project Structure

```
lcc-demo-app/
├── cmd/
│   └── demo/
│       └── main.go              # Application entry point
├── internal/
│   ├── analytics/
│   │   ├── advanced.go          # Professional tier feature
│   │   ├── basic.go             # Basic tier feature
│   │   └── lcc_gen.go           # Auto-generated (by lcc-sdk)
│   ├── export/
│   │   ├── cloud.go             # Enterprise tier feature
│   │   └── local.go             # Basic tier feature
│   └── reporting/
│       ├── premium.go           # Professional tier feature
│       └── standard.go          # Basic tier feature
├── lcc-features.yaml            # Feature manifest
├── configs/
│   ├── lcc-features.yaml        # Production config
│   ├── lcc-features.dev.yaml   # Development config
│   └── lcc-features.test.yaml  # Test config
└── docs/
    ├── TUTORIAL.md              # Step-by-step tutorial
    └── ARCHITECTURE.md          # Technical details
```

## Configuration

### lcc-features.yaml

```yaml
sdk:
  lcc_url: "http://localhost:7086"
  product_id: "demo-analytics-app"
  product_version: "1.0.0"

features:
  - id: advanced_analytics
    name: "Advanced Analytics with ML"
    tier: professional
    intercept:
      package: "github.com/yourorg/lcc-demo-app/internal/analytics"
      function: "AdvancedAnalytics"
    fallback:
      package: "github.com/yourorg/lcc-demo-app/internal/analytics"
      function: "BasicAnalytics"
    quota:
      limit: 1000
      period: daily

  - id: cloud_export
    name: "Cloud Export"
    tier: enterprise
    intercept:
      package: "github.com/yourorg/lcc-demo-app/internal/export"
      function: "ExportToCloud"
    on_deny:
      action: error
      message: "Cloud export requires Enterprise license"

  - id: premium_reports
    name: "Premium Reports"
    tier: professional
    intercept:
      package: "github.com/yourorg/lcc-demo-app/internal/reporting"
      function: "PremiumReports"
    fallback:
      package: "github.com/yourorg/lcc-demo-app/internal/reporting"
      function: "StandardReports"
```

## Building

```bash
# Generate license wrappers
make generate

# Build application
make build

# Run
./bin/demo-app
```

## Testing

```bash
# Run all tests
make test

# Test with different license tiers
make test-basic
make test-professional
make test-enterprise

# Test quota limits
make test-quota
```

## Development Workflow

1. **Write business logic** (no license code needed)
2. **Update lcc-features.yaml** to protect functions
3. **Run `make generate`** to create wrappers
4. **Build and test**

## Troubleshooting

### LCC Server Connection Failed

```bash
# Check if LCC is running
curl http://localhost:7086/health

# Check configuration
cat lcc-features.yaml
```

### Feature Always Denied

```bash
# Check license tier
lcc-sdk check --feature advanced_analytics

# Verify function mapping
lcc-sdk scan
```

### Code Generation Failed

```bash
# Check Go module
go mod tidy

# Verify package paths
lcc-sdk validate lcc-features.yaml
```

## Learn More

- [LCC SDK Documentation](https://github.com/yourorg/lcc-sdk)
- [Tutorial: Step-by-step Guide](docs/TUTORIAL.md)
- [Architecture Deep Dive](docs/ARCHITECTURE.md)

## License

MIT License - This is a demo application for educational purposes.
