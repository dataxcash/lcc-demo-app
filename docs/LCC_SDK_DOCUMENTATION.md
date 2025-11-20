# LCC SDK Documentation (Demo Overview)

This document provides a short, demo-oriented overview of how the LCC SDK is
used from this application. It is **not** the full SDK reference, but a
practical guide to the patterns demonstrated here.

## SDK Configuration

The SDK is configured via `config.SDKConfig`:

- `LCCURL`: URL of the LCC server
- `ProductID`: product identifier
- `ProductVersion`: product version string
- `Timeout`: HTTP timeout for SDK calls
- `CacheTTL`: how long feature checks are cached on the client side

Example (see `cmd/demo/main.go`):

```go
cfg := &config.SDKConfig{
    LCCURL:         "https://localhost:8088",
    ProductID:      "demo-app",
    ProductVersion: "1.0.0",
    Timeout:        30 * time.Second,
    CacheTTL:       10 * time.Second,
}
```

## Client Lifecycle

1. **Create client**
2. **Register** with LCC
3. Use the client for feature checks and reporting
4. Call `Close()` when shutting down

```go
client, err := client.NewClient(cfg)
if err != nil {
    // handle error
}
if err := client.Register(); err != nil {
    // handle error
}

defer client.Close()
```

## Core Operations

The demo focuses on four limitation models. The SDK exposes helpers for each.

### 1. Consumption (Quota / Event Count)

Use `Consume(featureID, amount, meta)` to check and consume quota in one step.

```go
allowed, remaining, reason, err := client.Consume("advanced_analytics", 1, nil)
if err != nil {
    // network or server error
}
if !allowed {
    // denied by license; reason describes why
}
// proceed with the operation
```

The helper returns:

- `allowed`: whether the operation is permitted
- `remaining`: approximate remaining quota (for UX only)
- `reason`: denial reason (`insufficient_tier`, `quota_exceeded`, etc.)

### 2. State Capacity (Current Count vs Max Capacity)

Use `CheckCapacity(featureID, currentValue)` when **your app** tracks the
current count (e.g., number of projects) and the license defines a max.

```go
allowed, max, reason, err := client.CheckCapacity("capacity.project.count", currentProjects)
if !allowed {
    // stop creating more entities; max reached
}
```

### 3. Rate / TPS (Current TPS vs Max TPS)

If your app already measures TPS, you can pass the current TPS directly to
`CheckTPS`.

```go
currentTPS := 7.5 // measured by your app
allowed, maxTPS, reason, err := client.CheckTPS("api.v1.demo", currentTPS)
if !allowed {
    // throttle or reject requests
}
```

The SDK itself does **not** compute TPS; it only compares against the license
limit using the value you provide.

### 4. Concurrency (Slots / In-flight Jobs)

Use `AcquireSlot(featureID, meta)` to guard concurrent operations (jobs,
users, connections) with a max concurrency limit.

```go
release, allowed, reason, err := client.AcquireSlot("concurrent.user", map[string]any{
    "job_id": 123,
})
if err != nil {
    // error talking to LCC
}
if !allowed {
    // no slot available; respect the limit
}
defer release()

// run the job here
```

The `release` function must be called to free the slot when the work finishes.

## Mapping Features to Application Code

Feature IDs and their mapping to concrete functions are defined in
`lcc-features.yaml` and the `configs/` variants. These manifests are
processed by `lcc-codegen` to generate wrappers that integrate licensing into
business functions.

For more details, see:

- `lcc-features.yaml`
- `configs/lcc-features.*.yaml`
- `internal/*/lcc_gen.go` (generated after running `make generate`)

A full SDK reference will be provided in the dedicated LCC SDK repository.
This document is meant to be a minimal, self-contained guide for this demo.
