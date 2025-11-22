# LCC SDK & Demo App Refactoring Task: Zero-Intrusion Design

## Context
We need to refactor the LCC (License Control Center) SDK and demo app to implement a **zero-intrusion, configuration-based limits control system**.

## Current Architecture Issues
1. ❌ Limits are defined at **feature-level** (incorrect - should be **product-level**)
2. ❌ SDK APIs require manual **Feature ID** parameters in business code (invasive)
3. ❌ Developers must manually call `CheckFeature()`, `Consume()`, `CheckTPS()`, etc.
4. ❌ No true "zero-intrusion" - requires code changes in business logic

## New Design Goals

### Core Principle: **Configuration-Based + Optional Helper Functions**

All limits should be:
- ✅ **Product-level** (shared across all features)
- ✅ **Configured in YAML** (not hardcoded in business logic)
- ✅ **Auto-injected by compiler** (zero code changes for developers)
- ✅ **Optionally customizable** via developer-provided helper functions

## Limits Control Strategy

| Limit Type | Auto-Injection | Developer Function | Required? |
|-----------|----------------|-------------------|-----------|
| **Quota** | ✅ Yes (default: 1 per call) | `QuotaConsumer` (custom amount) | Optional |
| **TPS** | ✅ Yes (SDK auto-tracks) | `TPSProvider` (custom measurement) | Optional |
| **Capacity** | ⚠️ Needs counter | `CapacityCounter` (get current usage) | **Required** |
| **Concurrency** | ✅ Yes (auto acquire/release) | None | Not needed |

## Target YAML Configuration Format

```yaml
sdk:
  product_id: data-insight-pro
  lcc_url: "http://localhost:7086"

# Product-level limits (all features share)
limits:
  quota:
    max: 50000
    window: monthly
    consumer: GetConsumeAmount  # Optional: custom amount calculator
  
  max_tps: 100.0
  tps_provider: GetCurrentTPS  # Optional: custom TPS measurement
  
  max_capacity: 100
  capacity_counter: GetCurrentProjectCount  # Required
  
  max_concurrency: 10

# Features (only define interception points)
features:
  - id: ml_analytics
    intercept:
      package: analytics
      function: ProcessAnalytics
    on_deny:
      action: error
      message: "ML Analytics requires Professional tier"
  
  - id: create_project
    intercept:
      package: projects
      function: CreateProject
    on_deny:
      action: error
      message: "Project creation limit reached"
```

## Developer Helper Functions (Optional/Required)

```go
// Optional: Quota - Custom consumption amount
func GetConsumeAmount(ctx context.Context, args ...interface{}) int {
    data := args[0].(Dataset)
    return data.Size() / 1024  // Charge by KB
}

// Optional: TPS - Custom rate measurement
func GetCurrentTPS() float64 {
    return myRateLimiter.GetRate()
}

// Required: Capacity - Current usage counter
func GetCurrentProjectCount() int {
    return db.Query("SELECT COUNT(*) FROM projects")
}
```

## Compiler Auto-Generated Code Example

### What Developer Writes (Clean Business Logic)
```go
func ProcessAnalytics(data Dataset) error {
    return analytics.Run(data)
}
```

### What Compiler Generates
```go
func ProcessAnalytics(data Dataset) error {
    // Auto-injected: Concurrency control
    release, ok := __lcc.AcquireSlot()
    if !ok { return ErrConcurrencyLimit }
    defer release()
    
    // Auto-injected: Quota consumption
    amount := GetConsumeAmount(context.Background(), data)  // Call dev function
    if !__lcc.Consume(amount) { return ErrQuotaExceeded }
    
    // Auto-injected: TPS check
    if !__lcc.CheckTPS() { return ErrRateLimitExceeded }
    
    // Original business logic
    return analytics.Run(data)
}
```

## SDK API Design Changes

### Old API (Feature-level, Invasive)
```go
// ❌ Requires featureID in every call
lccClient.Consume("ml_analytics", 1, nil)
lccClient.CheckTPS("api_access", currentTPS)
lccClient.CheckCapacity("projects", currentCount)
lccClient.AcquireSlot("concurrent_users", nil)
```

### New API (Product-level, Zero-intrusion)
```go
// ✅ Product-level APIs (auto-injected by compiler)
lccClient.Consume(amount)              // No featureID needed
lccClient.CheckTPS()                   // SDK auto-tracks
lccClient.CheckCapacity(currentCount)  // From helper function
lccClient.AcquireSlot()                // Auto acquire/release

// Optional: Tagged APIs for monitoring (not required)
lccClient.ConsumeWithTag("ml_analytics", amount)
lccClient.CheckTPSWithTag("api_access")
```

## Tasks to Complete

### Phase 1: lcc-demo-app (Backend & Frontend)

#### Backend Changes
- [ ] **products.go**: Refactor `GetLicenseJSON()` to ensure limits are product-level
  - ✅ Already done: limits moved to product level
  - [ ] Remove any remaining feature-level limit references
  
- [ ] **limits.go**: Update all 4 limit type examples
  - ✅ Already done: examples show product-level limits
  - [ ] Add helper function examples in code comments
  
- [ ] **API handlers**: Ensure all responses reflect product-level structure
  - [ ] `/api/tiers/{tier}/license` - verify limits structure
  - [ ] `/api/limits/{type}/example` - verify examples

#### Frontend Changes
- [ ] **tiers.js**: Update license display to show product-level limits
  - Current: Shows limits inside features (incorrect)
  - Target: Show limits as separate section at product level
  
- [ ] **limits.js**: Update limit examples display
  - [ ] Add helper function examples to UI
  - [ ] Show "optional vs required" indicator

#### Documentation
- [ ] **UI_DESIGN_SPEC.md**: Update license JSON structure examples
- [ ] Create new doc: **HELPER_FUNCTIONS_GUIDE.md** with examples

### Phase 2: lcc-sdk (If needed)

- [ ] Analyze current SDK implementation
- [ ] Refactor APIs to remove mandatory featureID
- [ ] Add helper function registration mechanism
- [ ] Implement auto-tracking for TPS
- [ ] Add compiler directive support

### Phase 3: Documentation & Examples

- [ ] Update YAML schema documentation
- [ ] Add complete helper function examples
- [ ] Create migration guide from old to new structure
- [ ] Update demo scenarios in UI

## File Paths

- **Demo app**: `/home/fila/jqdDev_2025/lcc-demo-app`
- **SDK**: `/home/fila/jqdDev_2025/lcc-sdk`
- **Main server**: `/home/fila/jqdDev_2025/lcc`

## Important Notes

- ⚠️ **All code comments and logs must be in English** (not Chinese)
- ✅ **Maintain backward compatibility** if possible
- 🎯 **Focus on zero-intrusion design** philosophy
- 📝 **Use Chinese for communication** with the user
- 🔄 **Git commits**: Use clear, descriptive messages in English

## Expected Outcome

A refactored system where:

1. ✅ Limits are truly **product-level** (shared by all features)
2. ✅ Developers write **clean business logic** without license check code
3. ✅ Compiler **auto-injects** all limit checks based on YAML config
4. ✅ Optional **helper functions** provide flexibility for custom logic
5. ✅ Demo app clearly demonstrates the zero-intrusion approach

## Success Criteria

- [ ] License JSON shows limits at product level (not per-feature)
- [ ] Demo UI correctly displays product-level limits
- [ ] All 4 limit types (Quota/TPS/Capacity/Concurrency) show correct structure
- [ ] Helper function examples are clear and documented
- [ ] Zero code changes needed in business logic (configuration only)

## Next Steps

1. Review current codebase structure
2. Identify all files that need changes
3. Create detailed refactoring plan
4. Implement changes incrementally
5. Test with demo scenarios
6. Update all documentation

---

**Last Updated**: 2025-01-22  
**Status**: Planning Phase  
**Related Discussion**: See conversation about zero-intrusion design
