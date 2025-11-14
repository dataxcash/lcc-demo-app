# LCC Demo App Architecture

This document describes the high-level architecture of the LCC demo
application and how it integrates with the LCC SDK and server.

## Components

- **cmd/demo**
  - `main.go`: Interactive CLI demo entrypoint.
  - Starts the status HTTP server (`/status`, `/status/json`).
  - Uses the lcc-sdk client to register with LCC and perform feature checks,
    quota consumption, capacity checks, TPS checks, and concurrency slot
    acquisition.

- **cmd/regression**
  - `main.go`: End-to-end regression driver.
  - Starts the compiled demo binary, sends menu choices via stdin, and
    validates behavior using the status JSON API.

- **internal/analytics**
  - Contains basic and advanced analytics functions.
  - Demonstrates feature-based enablement and consumption-style limits.

- **internal/export**
  - Contains PDF / Excel export functions.
  - Demonstrates tier-based gating and quota-limited operations.

- **internal/reporting**
  - Contains scheduling/reporting functions.
  - Demonstrates tier-based feature control and usage reporting.

## Limitation Models

The demo showcases four primary limitation models:

1. **Consumption** (quota / event count)
2. **State capacity** (current count vs max capacity)
3. **Rate / TPS** (current TPS vs max allowed TPS)
4. **Concurrency** (in-flight jobs/users vs max slots)

These are surfaced both through the CLI flows and via `DemoStats` in the
status server.

## Status Server

The status server runs inside the demo process and exposes:

- `/status`: A small HTML dashboard for humans
- `/status/json`: A JSON endpoint suitable for tests and tooling

The `DemoStats` structure is updated by each demo scenario, and tests/regression
use the JSON endpoint as a stable contract.

## Integration with LCC

- The demo uses `lcc-sdk` to talk to the LCC server.
- Product ID, version, and LCC URL are configured via SDK config and
  feature manifests (YAML).
- Limits (quota, capacity, TPS, concurrency) are derived from license/
  product configuration on the LCC side and returned to the SDK.

This structure is intended as a reference for integrating real applications
with LCC using the same patterns.
