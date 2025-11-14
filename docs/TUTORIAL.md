# LCC Demo App Tutorial

This tutorial walks you through running the LCC demo application and exploring
how license-based limitations (consumption, capacity, TPS, concurrency) affect
behavior.

## 1. Prerequisites

- LCC server running and reachable
- A valid license configured for the demo product(s)
- Go toolchain installed

## 2. Build and Run the Demo

```bash
cd lcc-demo-app
make build
./bin/demo-app
```

You will see an interactive menu in the terminal with several options,
including basic/advanced analytics, exports, state capacity demo, TPS demo,
and concurrency demo.

## 3. Status UI

While the demo is running, open the status UI in your browser:

- HTML: `http://localhost:8080/status`
- JSON: `http://localhost:8080/status/json`

Trigger different menu options and watch the metrics update.

## 4. End-to-end Regression

You can run an automated regression that drives the demo and validates that
all four limitation types are exercised and surfaced via the status JSON:

```bash
make regression
```

This will:

1. Build the demo binary
2. Start the demo on a test status port
3. Send menu choices to exercise features
4. Assert that the status JSON reflects usage of all four limitation models

## 5. Next Steps

- Integrate this demo into your CI as a smoke test for license behavior
- Extend the demo with additional features and limits
- Use the same patterns in your own applications with the lcc-sdk
