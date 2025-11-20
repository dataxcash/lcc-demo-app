# Web Demo Quick Start Guide

## 5-Minute Setup

### Step 1: Start LCC Server

```bash
cd /home/fila/jqdDev_2025/lcc
./lcc_server
```

Keep this terminal running.

### Step 2: Launch Web Demo

Open a new terminal:

```bash
cd /home/fila/jqdDev_2025/lcc-demo-app
make demo
```

Output:
```
Building web demo...
Starting web demo server...
Navigate to http://localhost:9144/discover
LCC Web Demo starting on http://localhost:9144
```

### Step 3: Open Browser

Navigate to: **http://localhost:9144/discover**

---

## Demo Walkthrough

### Page 1: Discover Products

You'll see three product tiers displayed as tabs:

- **Basic Edition**: Entry-level features
- **Professional Edition**: ← **Start here!**
- **Enterprise Edition**: Full-featured

**Action**: 
1. Click on "Professional Edition" tab
2. Review the features and limitations
3. Click **"Select & Configure Demo →"**

The system will:
- Register with LCC server
- Load Professional tier license
- Navigate to configuration page

---

### Page 2: Configure Simulation

Left panel shows configuration options:

**License Control Types** (pre-selected):
- ✅ Rate Limiting
- ✅ Quota Management  
- ⬜ Feature Gating
- ⬜ Capacity Control

**Runtime Parameters**:
- Loop Count: `100`
- Interval: `500ms`

Right panel shows **real SDK code examples** with impact analysis.

**Action**: 
1. Keep default settings (Rate Limiting + Quota)
2. Click **"Start Simulation →"**

---

### Page 3: Runtime Dashboard

Watch the simulation in action!

**Top Status Bar** shows:
- Progress: `47/100`
- Success: `45` ✓
- Failures: `2` ✗
- Rate Limits: `2` ⚠
- Elapsed: `23s`

**Main Panel**:
- **Chart**: Real-time success/failure trends
- **Code Context**: Current SDK calls being executed

**Right Panel**:
- **Event Log**: Live stream of license checks
  ```
  [13:45:23] ✓ Iteration 47: PDF export succeeded
  [13:45:22] ⚠ Iteration 46: Rate limit hit - quota_exceeded
  [13:45:22] ✓ Iteration 45: Advanced analytics completed (remaining: 955)
  ```

**Controls**:
- **Stop**: Halt simulation early
- **← Configure**: Return to setup

---

## What You'll Learn

### 1. Rate Limiting in Action

Around **iteration 200**, you'll see:
```
⚠ Rate limit hit - quota_exceeded
```

This demonstrates the Professional tier's 200/day PDF export limit.

**Code Behind It**:
```go
allowed, _, reason, err := lccClient.Consume("pdf_export", 1, nil)
if !allowed {
    return ErrDailyQuotaExceeded  // ← This happens
}
```

### 2. Quota Tracking

Watch the "remaining" counter decrease:
```
✓ Iteration 10: Advanced analytics completed (remaining: 990)
✓ Iteration 20: Advanced analytics completed (remaining: 980)
...
⚠ Iteration 1001: Quota exceeded - quota_exceeded
```

**Code Behind It**:
```go
allowed, remaining, reason, err := lccClient.Consume("advanced_analytics", 1, nil)
// remaining shows how many calls are left
```

---

## Try Different Scenarios

### Scenario A: Test Feature Gating

1. Click **← Configure**
2. Uncheck "Rate Limiting" and "Quota"
3. Check **"Feature Gating"**
4. Start simulation

**Result**: All iterations fail because Professional tier doesn't have Excel export (Enterprise only).

### Scenario B: Test Capacity Limits

1. Configure → Select "Capacity Control"
2. Set Loop Count: `60`
3. Start simulation

**Result**: Fails after 50 iterations (Professional tier project limit).

### Scenario C: Compare Product Tiers

1. Go back to **Discover** page
2. Select **"Basic Edition"**
3. Configure with same settings
4. Start simulation

**Result**: Much more restrictive limits - only 100 API calls/day.

---

## Architecture At A Glance

```
┌─────────────┐
│   Browser   │ ← You interact here
└──────┬──────┘
       │ HTTP
       ↓
┌─────────────────┐
│  webdemo (Go)   │ ← Handles simulation
│  :9144          │
└──────┬──────────┘
       │ lcc-sdk
       ↓
┌─────────────────┐
│  LCC Server     │ ← License validation
│  :8088          │
└─────────────────┘
```

**Data Flow**:
1. Product selection → SDK initialization → License registration
2. Configuration → Simulation setup
3. Runtime → SDK API calls → License checks → Metrics updates

---

## Troubleshooting

### Problem: "Failed to select product"

**Solution**: 
- Check LCC server is running: `curl http://localhost:8088/health`
- Review server logs in LCC terminal

### Problem: Metrics not updating

**Solution**:
- Press F12 to open browser console
- Check for JavaScript errors
- Verify simulation is "Running" (green badge)

### Problem: Port already in use

**Solution**:
```bash
PORT=9145 ./bin/webdemo
# Then navigate to http://localhost:9145/discover
```

---

## Next Steps

1. **Explore the Code**:
   - `cmd/webdemo/main.go` - Backend logic
   - `cmd/webdemo/static/*.html` - Frontend UI

2. **Read Full Documentation**:
   - [Web Demo README](cmd/webdemo/README.md)
   - [Main Demo README](README.md)

3. **Try Custom Scenarios**:
   - Modify loop count and interval
   - Enable multiple control types simultaneously
   - Compare all three product tiers

4. **Integrate Into Your App**:
   - Study the SDK code examples
   - Apply patterns to your own business logic

---

## Quick Reference

| Command | Description |
|---------|-------------|
| `make demo` | Build and start web demo |
| `./bin/webdemo` | Run pre-built binary |
| `PORT=9145 ./bin/webdemo` | Use custom port |
| `make clean` | Remove build artifacts |

| URL | Page |
|-----|------|
| `http://localhost:9144/discover` | Product catalog |
| `http://localhost:9144/configure` | Simulation setup |
| `http://localhost:9144/runtime` | Live dashboard |
| `http://localhost:9144/api/simulation/status` | JSON API endpoint |

---

**Enjoy exploring LCC SDK capabilities!** 🚀
