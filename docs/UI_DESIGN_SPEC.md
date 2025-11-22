# LCC SDK Demo App - UI Design Specification

**Version:** 1.0.0  
**Date:** 2025-01-21  
**Purpose:** Educational interactive platform for LCC SDK learning and demonstration

---

## 📋 Table of Contents

1. [Design Philosophy](#design-philosophy)
2. [Overall Architecture](#overall-architecture)
3. [Page Specifications](#page-specifications)
4. [Three-Tier Product Design](#three-tier-product-design)
5. [API Specifications](#api-specifications)
6. [Design System](#design-system)
7. [Technical Implementation](#technical-implementation)

---

## Design Philosophy

### Core Principles

1. **Educational First**: This is a teaching tool, not just a demo. Every screen should explain concepts clearly.
2. **Progressive Learning**: Users advance from basic concepts (Tier) to advanced (Limits) to practical (Runtime).
3. **Code-Centric**: SDK code examples are the primary focus, with syntax highlighting and real-time annotations.
4. **Interactive**: Users should "learn by doing" with live simulations and immediate feedback.
5. **Self-Contained**: All necessary context is provided within the app - no external documentation required.

### Target Audience

- **Developers**: Learning how to integrate LCC SDK
- **Product Managers**: Understanding tier-based monetization
- **Sales Teams**: Demonstrating product capabilities
- **Technical Evaluators**: Assessing LCC platform features

---

## Overall Architecture

### Page Flow Diagram

```
┌─────────────────┐
│  1. Welcome     │  → Introduction + LCC Configuration
│     & Setup     │     Save/load settings, test connection
└────────┬────────┘
         ↓
┌─────────────────┐
│  2. Tier        │  → Lesson 1: Understanding License Tiers
│     Learning    │     Basic vs Pro vs Enterprise comparison
└────────┬────────┘
         ↓
┌─────────────────┐
│  3. Limits      │  → Lesson 2: Four Control Types
│     Learning    │     Quota / TPS / Capacity / Concurrency
└────────┬────────┘
         ↓
┌─────────────────┐
│  4. Instance    │  → Configure Simulation
│     Setup       │     Select license + controls + parameters
└────────┬────────┘
         ↓
┌─────────────────┐
│  5. Runtime     │  → Live Execution Dashboard
│     Dashboard   │     Charts / Metrics / Logs / Code tracing
└─────────────────┘
```

### Navigation Structure

```
Header (Persistent):
  - App Title: "LCC SDK Interactive Learning Demo"
  - Current Step Indicator: [1] → [2] → [3] → [4] → [5]
  - LCC Status Badge: ● Connected / ● Disconnected

Left Sidebar (Pages 2-5):
  - Quick Navigation
  - Progress Tracker
  - Mini Summary of Current Lesson

Main Content Area:
  - Primary content with cards
  - Code examples with syntax highlighting
  - Interactive demos

Footer:
  - Navigation buttons: [← Back] [Next →]
  - Action buttons: [Save] [Reset] [Help]
```

---

## Page Specifications

### Page 1: Welcome & Configuration

#### Layout

```
┌────────────────────────────────────────────────────────────┐
│  🎓 LCC SDK Interactive Learning Demo                     │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Welcome to the LCC SDK Tutorial Platform!                │
│                                                            │
│  This interactive demo will guide you through:            │
│  ✓ Understanding License Tiers (Basic/Pro/Enterprise)     │
│  ✓ Learning 4 Control Types (Quota/TPS/Capacity/Users)   │
│  ✓ Hands-on SDK Integration Examples                     │
│  ✓ Live Runtime Behavior Observation                     │
│                                                            │
│  ┌──────────────────────────────────────────────────────┐ │
│  │ 🔧 LCC Server Configuration                          │ │
│  │                                                      │ │
│  │  LCC URL: [http://localhost:7086            ]       │ │
│  │                                                      │ │
│  │  Status: ● Connected  |  Version: v2.1.0            │ │
│  │          3 products available                       │ │
│  │                                                      │ │
│  │  ⚠️ Configuration changed - click UPDATE to apply    │ │
│  │                                                      │ │
│  │  [Test Connection]  [UPDATE]  [Save & Continue →]   │ │
│  └──────────────────────────────────────────────────────┘ │
│                                                            │
│  💡 First time here? The default URL works if LCC is      │
│     running locally. Otherwise, update to your server.    │
│                                                            │
│  📚 What You'll Learn:                                     │
│  • How licenses control feature availability (Tiers)      │
│  • How limits control usage amounts (Quota/TPS/etc)       │
│  • How to integrate SDK into your application             │
│  • How to handle license checks and denials gracefully    │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

#### Behavior

- **On Load**: 
  - Fetch saved configuration from server (`GET /api/config`)
  - Auto-fill LCC URL input
  - Auto-test connection and display status
  
- **Configuration Change Detection**:
  - Compare current input with saved value
  - Show warning banner: "Configuration changed - click UPDATE to apply"
  - Disable "Continue" button until UPDATE is clicked
  
- **Test Connection**:
  - `GET /api/config/validate`
  - Show loading spinner
  - Display result: Connected ✓ / Failed ✗
  - If connected, show: version, product count
  
- **UPDATE Button**:
  - `POST /api/config` with new URL
  - Save to server-side persistent storage
  - Clear "changed" warning
  - Enable "Continue" button
  
- **Save & Continue**:
  - Navigate to Page 2 (Tier Learning)

#### API Endpoints

```
GET  /api/config
Response: {
  "lcc_url": "http://localhost:7086",
  "saved_at": "2025-01-21T10:30:00Z",
  "is_default": true
}

POST /api/config
Request: {
  "lcc_url": "http://localhost:7086"
}
Response: {
  "ok": true,
  "lcc_url": "http://localhost:7086"
}

GET  /api/config/validate
Response: {
  "reachable": true,
  "version": "v2.1.0",
  "products_count": 3,
  "error": null
}
```

---

### Page 2: Tier Feature Learning

#### Layout

```
┌────────────────────────────────────────────────────────────┐
│  📚 Lesson 1: Understanding License Tiers                  │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  🎯 Learning Objective                                     │
│  License tiers control WHICH features are available to    │
│  different customer segments. Higher tiers unlock more    │
│  capabilities - this is the foundation of monetization.   │
│                                                            │
│  Product: Data Insight Analytics Platform                 │
│                                                            │
│  ┌─ Tier Comparison ─────────────────────────────────────┐│
│  │  [Basic] [Professional] [Enterprise]  ← Tab Selection ││
│  │                                                        ││
│  │  ┌─────────────────┬───────────┬──────────┬──────────┐││
│  │  │ Feature         │   Basic   │   Pro    │   Ent    │││
│  │  ├─────────────────┼───────────┼──────────┼──────────┤││
│  │  │ Basic Reports   │     ✓     │    ✓     │    ✓     │││
│  │  │ ML Analytics    │     ✗     │    ✓     │    ✓     │││
│  │  │ PDF Export      │     ✗     │    ✓     │    ✓     │││
│  │  │ Excel Export    │     ✗     │    ✗     │    ✓     │││
│  │  │ Custom Dashbrd  │     ✗     │    ✗     │    ✓     │││
│  │  │ API Access      │     ✗     │    ✓     │    ✓     │││
│  │  └─────────────────┴───────────┴──────────┴──────────┘││
│  └────────────────────────────────────────────────────────┘│
│                                                            │
│  ┌─ How It Works (Professional Tier Selected) ───────────┐│
│  │                                                        ││
│  │  📄 License File Configuration:                       ││
│  │  {                                                    ││
│  │    "product_id": "data-insight-pro",                  ││
│  │    "tier": "professional",                            ││
│  │    "features": {                                      ││
│  │      "basic_reports": {                               ││
│  │        "enabled": true                                ││
│  │      },                                               ││
│  │      "ml_analytics": {                                ││
│  │        "enabled": true                                ││
│  │      },                                               ││
│  │      "pdf_export": {                                  ││
│  │        "enabled": true                                ││
│  │      },                                               ││
│  │      "excel_export": {                                ││
│  │        "enabled": false  ← ❌ Requires Enterprise     ││
│  │      }                                                ││
│  │    }                                                  ││
│  │  }                                                    ││
│  │                                                        ││
│  │  💻 Zero-Intrusion Approach:                          ││
│  │                                                        ││
│  │  ✍️ What Developers Write: (Original Source)          ││
│  │  ┌────────────────────────────────────────────────┐  ││
│  │  │ // Developer only writes business logic        │  ││
│  │  │ // No license check code needed!               │  ││
│  │  │ func ExportToExcel(reportID string) error {    │  ││
│  │  │     // Direct business logic                   │  ││
│  │  │     return generateExcelReport(reportID)       │  ││
│  │  │ }                                              │  ││
│  │  └────────────────────────────────────────────────┘  ││
│  │                                                        ││
│  │              ⬇️ Compiler Automatically Transforms      ││
│  │                                                        ││
│  │  🔧 After Compilation: (Auto-generated Code)           ││
│  │  ┌────────────────────────────────────────────────┐  ││
│  │  │ // Compiler automatically inserts checks       │  ││
│  │  │ func ExportToExcel(reportID string) error {    │  ││
│  │  │   // === Auto-generated by LCC compiler ===    │  ││
│  │  │   status, err := lccClient.CheckFeature(       │  ││
│  │  │     "excel_export")                            │  ││
│  │  │   if err != nil {                              │  ││
│  │  │     return fmt.Errorf("license: %w", err)     │  ││
│  │  │   }                                            │  ││
│  │  │   if !status.Enabled {                         │  ││
│  │  │     return fmt.Errorf("Requires Enterprise")  │  ││
│  │  │   }                                            │  ││
│  │  │   // ==================================        │  ││
│  │  │   // Original developer code                  │  ││
│  │  │   return generateExcelReport(reportID)        │  ││
│  │  │ }                                              │  ││
│  │  └────────────────────────────────────────────────┘  ││
│  │                                                        ││
│  │  🎯 Key Benefit: Developers write clean business      ││
│  │     logic. License enforcement is automatically       ││
│  │     injected during compilation.                      ││
│  │                                                        ││
│  │  ⚙️ Developer Configuration (lcc-features.yaml):      ││
│  │  features:                                            ││
│  │    - id: excel_export                                 ││
│  │      name: Excel Export                               ││
│  │      intercept:                                       ││
│  │        package: exports                               ││
│  │        function: ExportToExcel                        ││
│  │      on_deny:                                         ││
│  │        action: error                                  ││
│  │        message: "Upgrade to Enterprise tier"          ││
│  │                                                        ││
│  │  [View Complete YAML] [Download Sample]               ││
│  │                                                        ││
│  └────────────────────────────────────────────────────────┘│
│                                                            │
│  ┌─ Try It Yourself ──────────────────────────────────────┐│
│  │                                                        ││
│  │  Simulate calling different features:                 ││
│  │                                                        ││
│  │  Feature: [excel_export ▼]  [Call CheckFeature()]    ││
│  │                                                        ││
│  │  Result: ❌ Denied                                    ││
│  │  {                                                    ││
│  │    "enabled": false,                                  ││
│  │    "reason": "tier_insufficient",                     ││
│  │    "required_tier": "enterprise",                     ││
│  │    "current_tier": "professional"                     ││
│  │  }                                                    ││
│  │                                                        ││
│  └────────────────────────────────────────────────────────┘│
│                                                            │
│  📝 Key Concepts:                                          │
│  • Tier = Customer segment (Basic/Pro/Enterprise)         │
│  • Feature ID = Unique identifier (e.g., "excel_export")  │
│  • CheckFeature() = SDK API to gate business logic        │
│  • License file = Server-managed, ISV applies to customer │
│  • YAML config = Developer-defined, compiled into app     │
│  • Enabled flag = Authoritative control from license      │
│                                                            │
│  [← Back to Welcome]  [Next: Limits Learning →]           │
└────────────────────────────────────────────────────────────┘
```

#### Interactive Features

1. **Tier Tab Selection**:
   - Click Basic/Pro/Enterprise to switch
   - Comparison table updates instantly
   - License JSON updates to show selected tier
   - Code example remains same (demonstrates checking)
   
2. **Feature Hover**:
   - Hover over feature name shows tooltip
   - Explains what the feature does
   - Shows which tiers include it
   
3. **Try It Yourself Section**:
   - Dropdown lists all features
   - "Call CheckFeature()" button triggers simulation
   - Shows actual SDK response
   - Visual indicator: ✓ Allowed / ❌ Denied
   
4. **Code Syntax Highlighting**:
   - Use Prism.js with Go language
   - Highlight key API calls differently
   - Show comments explaining each step

#### API Endpoints

```
GET  /api/tiers
Response: [
  {
    "id": "basic",
    "name": "Basic Edition",
    "product_id": "data-insight-basic",
    "features": [
      {"id": "basic_reports", "enabled": true},
      {"id": "ml_analytics", "enabled": false},
      ...
    ]
  },
  ...
]

GET  /api/tiers/{tier}/license
Response: {
  "product_id": "data-insight-pro",
  "tier": "professional",
  "features": { ... }
}

GET  /api/tiers/{tier}/yaml
Response: {
  "yaml_content": "features:\n  - id: excel_export\n    ..."
}

POST /api/tiers/{tier}/check-feature
Request: {
  "feature_id": "excel_export"
}
Response: {
  "enabled": false,
  "reason": "tier_insufficient",
  "required_tier": "enterprise",
  "current_tier": "professional"
}
```

---

### Page 3: Limits Learning

#### Layout Structure

This page uses a **tab-based interface** to teach four limit types:

```
┌────────────────────────────────────────────────────────────┐
│  📚 Lesson 2: Understanding License Limits                 │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  🎯 Learning Objective                                     │
│  Beyond ON/OFF feature gating, limits control HOW MUCH    │
│  customers can consume. Four types serve different uses.  │
│                                                            │
│  [1️⃣ Quota] [2️⃣ TPS] [3️⃣ Capacity] [4️⃣ Concurrency]     │
│     ↑ Selected                                             │
│                                                            │
│  ┌─ Quota (配额控制) ────────────────────────────────────┐│
│  │                                                        ││
│  │  📖 What Is It?                                        ││
│  │  Cumulative consumption limit that resets on schedule ││
│  │  Server tracks total usage automatically              ││
│  │                                                        ││
│  │  🎯 Use Cases:                                         ││
│  │  • API call counting (10,000/day)                     ││
│  │  • Export operations (200 PDFs/month)                 ││
│  │  • License generation credits                         ││
│  │  • Any metered/consumable resource                    ││
│  │                                                        ││
│  │  ⏰ Time Dimension:                                    ││
│  │  Daily or Monthly window with auto-reset              ││
│  │                                                        ││
│  │  🔢 Who Tracks?                                        ││
│  │  Server-side (LCC) - developer just calls Consume()   ││
│  │                                                        ││
│  │  📄 License Configuration:                            ││
│  │  {                                                    ││
│  │    "ml_analytics": {                                  ││
│  │      "enabled": true,                                 ││
│  │      "quota": {                                       ││
│  │        "max": 10000,         ← Total allowed         ││
│  │        "window": "daily",    ← Reset period          ││
│  │        "reset_at": "00:00"   ← Reset time (UTC)      ││
│  │      }                                                ││
│  │    }                                                  ││
│  │  }                                                    ││
│  │                                                        ││
│  │  💻 SDK Integration:                                   ││
│  │  func ProcessAnalytics(data Dataset) error {          ││
│  │    // SDK reports usage to LCC automatically          ││
│  │    allowed, remaining, reason, err :=                 ││
│  │      lccClient.Consume(                               ││
│  │        "ml_analytics",  // Feature ID                 ││
│  │        1,               // Credits to consume         ││
│  │        nil,             // Optional metadata          ││
│  │      )                                                ││
│  │                                                        ││
│  │    if err != nil {                                    ││
│  │      return fmt.Errorf("license: %w", err)           ││
│  │    }                                                  ││
│  │                                                        ││
│  │    if !allowed {                                      ││
│  │      log.Warn("Quota exceeded",                       ││
│  │        "remaining", remaining,                        ││
│  │        "reason", reason)                              ││
│  │      return ErrQuotaExceeded                          ││
│  │    }                                                  ││
│  │                                                        ││
│  │    // Quota OK - execute expensive operation          ││
│  │    result := analytics.RunMLModel(data)               ││
│  │    log.Info("Success", "remaining", remaining)        ││
│  │    return nil                                         ││
│  │  }                                                    ││
│  │                                                        ││
│  │  🔄 Runtime Behavior:                                  ││
│  │  ┌────────────────────────────────────────────────┐  ││
│  │  │ Call #      Allowed    Remaining    Reason     │  ││
│  │  ├────────────────────────────────────────────────┤  ││
│  │  │ 1           ✓ Yes      9,999        ok         │  ││
│  │  │ 100         ✓ Yes      9,900        ok         │  ││
│  │  │ 9,999       ✓ Yes      1            ok         │  ││
│  │  │ 10,000      ✓ Yes      0            ok         │  ││
│  │  │ 10,001      ❌ No       0            exceeded   │  ││
│  │  │ (Next Day)  ✓ Yes      9,999        reset      │  ││
│  │  └────────────────────────────────────────────────┘  ││
│  │                                                        ││
│  │  ✨ Key Points:                                        ││
│  │  • Server tracks cumulative total automatically       ││
│  │  • Developer only needs to call Consume()             ││
│  │  • Remaining count returned for UI display            ││
│  │  • Auto-resets daily/monthly per license config       ││
│  │                                                        ││
│  │  [Run Interactive Simulation]                         ││
│  │                                                        ││
│  └────────────────────────────────────────────────────────┘│
│                                                            │
│  [← Back]  [Next Tab: TPS →]  [Skip to Instance Setup →] │
└────────────────────────────────────────────────────────────┘
```

#### Four Tabs Content Summary

**Tab 1: Quota**
- What: Cumulative limit with reset
- Tracking: Server-side automatic
- API: `Consume(featureID, amount)`
- Developer: Just call Consume()
- Example: 10,000 API calls per day

**Tab 2: TPS (Rate Limit)**
- What: Instantaneous throughput limit
- Tracking: Client calculates current rate
- API: `CheckTPS(featureID, currentTPS)`
- Developer: Measure TPS, pass to CheckTPS()
- Example: Max 10 requests per second

**Tab 3: Capacity**
- What: Maximum quantity of persistent resources
- Tracking: Client counts current usage
- API: `CheckCapacity(featureID, currentUsed)`
- Developer: Count resources, pass to CheckCapacity()
- Example: Max 50 projects

**Tab 4: Concurrency**
- What: Simultaneous execution slots
- Tracking: SDK internal counter
- API: `AcquireSlot(featureID)` → returns `release()`
- Developer: Call AcquireSlot(), defer release()
- Example: Max 10 concurrent users

#### Comparison Table (Shown on All Tabs)

```
┌─ Control Type Comparison ──────────────────────────────────┐
│                                                            │
│  Control      License Config      SDK API        Developer │
│  ────────────────────────────────────────────────────────  │
│  Quota        quota: {max,        Consume(id,    Just call │
│               window}             amount)        Consume() │
│                                                            │
│  TPS          max_tps: 10.0       CheckTPS(id,   Measure   │
│                                   current)       TPS       │
│                                                            │
│  Capacity     max_capacity: 50    CheckCapacity Count      │
│                                   (id, used)     resources │
│                                                            │
│  Concurrency  max_concurrency:    AcquireSlot   Call +     │
│               10                  (id)           release() │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

#### Interactive Simulation

Each tab has a "Run Interactive Simulation" button:

```
┌─ Mini Simulation: Quota ───────────────────────────────────┐
│                                                            │
│  Scenario: 10 API calls, Quota = 10,000                   │
│                                                            │
│  [Start Simulation]  [Reset]                               │
│                                                            │
│  Progress: [████████████████████] 10/10                    │
│                                                            │
│  Results:                                                  │
│  ✓ Call 1:  allowed=true,  remaining=9,999                │
│  ✓ Call 2:  allowed=true,  remaining=9,998                │
│  ...                                                       │
│  ✓ Call 10: allowed=true,  remaining=9,990                │
│                                                            │
│  All calls succeeded within quota!                         │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

#### API Endpoints

```
GET  /api/limits/types
Response: [
  {
    "type": "quota",
    "name": "Quota Control",
    "description": "Cumulative consumption with reset",
    "sdk_api": "Consume(featureID, amount)",
    ...
  },
  ...
]

GET  /api/limits/{type}/example
Response: {
  "license_config": { ... },
  "code_example": "...",
  "behavior_table": [ ... ]
}

POST /api/limits/{type}/simulate
Request: {
  "feature_id": "ml_analytics",
  "iterations": 10,
  "params": { ... }
}
Response: {
  "results": [
    {"iteration": 1, "allowed": true, "remaining": 9999},
    ...
  ]
}
```

---

### Page 4: Instance Setup

#### Layout

```
┌────────────────────────────────────────────────────────────┐
│  ⚙️ Configure Simulation Instance                          │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Now that you understand Tiers and Limits, let's create   │
│  a simulation to observe SDK behavior in real-time!       │
│                                                            │
│  ┌─ Step 1: Select License Tier ────────────────────────┐ │
│  │                                                       │ │
│  │  Which tier do you want to simulate?                 │ │
│  │                                                       │ │
│  │  ○ Basic Edition      (Limited features)             │ │
│  │  ● Professional       (Moderate limits) ← Selected   │ │
│  │  ○ Enterprise         (Highest quotas)               │ │
│  │                                                       │ │
│  │  License Details:                                     │ │
│  │  • Product ID: data-insight-pro                      │ │
│  │  • Version:    1.0.0                                 │ │
│  │  • Instance:   lcc-demo-abc123 (auto-generated)      │ │
│  │                                                       │ │
│  │  [View Full License JSON]                             │ │
│  │                                                       │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                            │
│  ┌─ Step 2: Choose Controls to Demonstrate ─────────────┐ │
│  │                                                       │ │
│  │  Select which SDK controls to test in simulation:    │ │
│  │                                                       │ │
│  │  ☑ Tier Gating         (CheckFeature API)            │ │
│  │  ☑ Quota Management    (Consume API)                 │ │
│  │  ☑ TPS Limiting        (CheckTPS API)                │ │
│  │  ☐ Capacity Control    (CheckCapacity API)           │ │
│  │  ☐ Concurrency Slots   (AcquireSlot API)             │ │
│  │                                                       │ │
│  │  💡 Tip: Start with Quota for first run              │ │
│  │                                                       │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                            │
│  ┌─ Step 3: Features & Their Limits ────────────────────┐ │
│  │                                                       │ │
│  │  Based on Professional license:                       │ │
│  │                                                       │ │
│  │  📊 ml_analytics                                      │ │
│  │  ├─ Enabled:    ✓ Yes                                │ │
│  │  ├─ Quota:      10,000 / day (resets daily)          │ │
│  │  └─ Max TPS:    10.0 requests/sec                    │ │
│  │                                                       │ │
│  │  📄 pdf_export                                        │ │
│  │  ├─ Enabled:    ✓ Yes                                │ │
│  │  ├─ Quota:      200 / day                            │ │
│  │  └─ Max TPS:    5.0 requests/sec                     │ │
│  │                                                       │ │
│  │  📊 excel_export                                      │ │
│  │  └─ Enabled:    ✗ No (Requires Enterprise)           │ │
│  │                                                       │ │
│  │  🔌 api_access                                        │ │
│  │  ├─ Enabled:    ✓ Yes                                │ │
│  │  ├─ Max TPS:    100.0 requests/sec                   │ │
│  │  └─ Max Users:  10 concurrent                        │ │
│  │                                                       │ │
│  │  [Expand All Details]                                 │ │
│  │                                                       │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                            │
│  ┌─ Step 4: Simulation Runtime Parameters ──────────────┐ │
│  │                                                       │ │
│  │  Loop Count:   [100  ] iterations                    │ │
│  │  Interval:     [500  ] ms between iterations         │ │
│  │                                                       │ │
│  │  Features to Call:                                    │ │
│  │  ☑ ml_analytics  (every iteration)                   │ │
│  │  ☑ pdf_export    (every 5th iteration)               │ │
│  │  ☐ api_access    (disabled - not selected above)     │ │
│  │                                                       │ │
│  │  📊 Prediction:                                       │ │
│  │  • Estimated Runtime: ~50 seconds                     │ │
│  │  • ml_analytics calls: 100 (within 10K quota ✓)      │ │
│  │  • pdf_export calls: 20 (within 200 quota ✓)         │ │
│  │  • Expected: All calls should succeed                 │ │
│  │                                                       │ │
│  │  ⚠️ To trigger quota limit, increase to 300 loops    │ │
│  │                                                       │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                            │
│  ┌─ Review Configuration ────────────────────────────────┐ │
│  │                                                       │ │
│  │  Instance Name: [My First Simulation        ]        │ │
│  │                                                       │ │
│  │  Summary:                                             │ │
│  │  • Tier:        Professional                          │ │
│  │  • Controls:    Tier Gating, Quota, TPS              │ │
│  │  • Iterations:  100 x 500ms = 50s                     │ │
│  │  • Features:    ml_analytics, pdf_export              │ │
│  │                                                       │ │
│  │  ☑ Save this configuration for later reuse            │ │
│  │                                                       │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                            │
│  [← Back to Learning]  [Reset]  [Start Simulation →]      │
└────────────────────────────────────────────────────────────┘
```

#### Behavior

**Step 1: License Selection**
- Radio buttons for Basic/Pro/Enterprise
- Auto-generates instance ID (UUID)
- Loads license JSON for selected tier
- "View Full License JSON" opens modal with complete config

**Step 2: Control Selection**
- Checkboxes for 5 control types
- Selecting a control auto-enables relevant features in Step 3
- Deselecting hides corresponding features

**Step 3: Feature Display**
- Read-only display of license-defined limits
- Expands/collapses for readability
- Icons differentiate feature types (📊 analytics, 📄 export, etc)
- Red ✗ for disabled features with upgrade hint

**Step 4: Runtime Config**
- Number inputs with validation
- Real-time prediction of behavior
- Warning if configuration will hit limits
- Feature call frequency (every iteration vs every Nth)

**Review Section**
- Named configuration for saving
- Clear summary of all choices
- Save checkbox for reuse

**Validation**
- "Start Simulation" disabled until:
  - At least one control selected
  - At least one feature selected
  - Valid iteration count (1-10000)
  - Valid interval (10-10000ms)

#### API Endpoints

```
POST /api/instances/create
Request: {
  "name": "My First Simulation",
  "tier": "professional",
  "enabled_controls": ["tier", "quota", "tps"],
  "features_to_call": ["ml_analytics", "pdf_export"],
  "iterations": 100,
  "interval_ms": 500,
  "call_pattern": {
    "ml_analytics": 1,    // Every iteration
    "pdf_export": 5       // Every 5th iteration
  },
  "save_config": true
}
Response: {
  "instance_id": "lcc-demo-abc123",
  "status": "created",
  "estimated_runtime_sec": 50,
  "predictions": {
    "ml_analytics_calls": 100,
    "pdf_export_calls": 20,
    "expected_quota_usage": {
      "ml_analytics": "100/10000",
      "pdf_export": "20/200"
    }
  }
}

GET  /api/instances
Response: [
  {
    "instance_id": "lcc-demo-abc123",
    "name": "My First Simulation",
    "tier": "professional",
    "created_at": "2025-01-21T11:00:00Z",
    "status": "ready"
  },
  ...
]

GET  /api/instances/{id}
Response: {
  "instance_id": "lcc-demo-abc123",
  "config": { ... },
  "status": "ready",
  "license_info": { ... }
}
```

---

### Page 5: Runtime Dashboard

#### Layout

```
┌────────────────────────────────────────────────────────────┐
│  📊 Live Simulation Dashboard                              │
│                                                            │
│  Instance: My First Simulation  |  Tier: Professional     │
│  Status: ● Running  [Pause] [Stop] [Reset]                │
├────────────────────────────────────────────────────────────┤
│  Progress: [████████████░░░░░░░░] 47/100                  │
│  Elapsed: 23.5s  |  Remaining: ~26s  |  Success: 45       │
├──────────────────────────────┬─────────────────────────────┤
│  📈 Real-Time Metrics        │  📋 Event Log               │
│                              │                             │
│  ┌─ Success vs Failure ───┐  │  [11:23:47.234] ✓ Iter 47  │
│  │                        │  │    ml_analytics            │
│  │      /\  /\            │  │    Consume() succeeded     │
│  │     /  \/  \  /\       │  │    remaining: 9,953        │
│  │    /        \/  \      │  │                            │
│  │ __/              \_    │  │  [11:23:47.134] ✓ Iter 47  │
│  └────────────────────────┘  │    pdf_export (skipped)    │
│                              │    not in pattern          │
│  ┌─ Quota Status ─────────┐  │                            │
│  │  ml_analytics:         │  │  [11:23:46.734] ⚠ Iter 46  │
│  │  [████████████░░] 9,953│  │    pdf_export              │
│  │                        │  │    CheckTPS() failed       │
│  │  pdf_export:           │  │    6.2 > 5.0 TPS          │
│  │  [████████████░░] 198  │  │                            │
│  └────────────────────────┘  │  [11:23:46.634] ✓ Iter 46  │
│                              │    ml_analytics            │
│  ┌─ TPS Monitor ──────────┐  │    Consume() succeeded     │
│  │  ml_analytics:         │  │    remaining: 9,954        │
│  │    Current:  8.5 TPS   │  │                            │
│  │    Max:      10.0 TPS  │  │  [11:23:46.234] ✓ Iter 45  │
│  │    Status:   ✓ OK      │  │    pdf_export              │
│  │                        │  │    Consume() succeeded     │
│  │  pdf_export:           │  │    remaining: 197          │
│  │    Current:  6.2 TPS   │  │                            │
│  │    Max:      5.0 TPS   │  │  [Filters: All ▼]          │
│  │    Status:   ⚠ LIMIT   │  │  [📥 Export Log]           │
│  └────────────────────────┘  │  [Auto-scroll: ON]         │
│                              │                             │
│  [Expand All Metrics]        │  [Clear] [Search...]       │
├──────────────────────────────┴─────────────────────────────┤
│  💻 Current SDK Call Context (Live Tracing)                │
│                                                            │
│  Iteration: 47  |  Timestamp: 11:23:47.234                │
│                                                            │
│  func RunIteration(iter int) error {                       │
│    // 1. Call ml_analytics (Consume API)                   │
│    allowed, remaining, reason, err :=                      │
│      lccClient.Consume("ml_analytics", 1, nil)             │
│    ✓ SUCCESS: allowed=true, remaining=9953                 │
│                                                            │
│    if !allowed {                                           │
│      log.Warn("Denied", "reason", reason)                  │
│      return ErrQuotaExceeded                               │
│    }                                                       │
│                                                            │
│    // 2. Check if should call pdf_export (every 5th)       │
│    if iter % 5 != 0 {                                      │
│      return nil  // Skip this iteration                    │
│    ✓ EXECUTED: Skipped (47 % 5 != 0)                       │
│    }                                                       │
│                                                            │
│    // 3. Call pdf_export (CheckTPS API)                    │
│    currentTPS := metrics.GetTPS("pdf_export")              │
│    allowed, maxTPS, reason, err :=                         │
│      lccClient.CheckTPS("pdf_export", currentTPS)          │
│    ⚠ Previous iteration: allowed=false (6.2 > 5.0 TPS)     │
│                                                            │
│    return nil                                              │
│  }                                                         │
│                                                            │
│  [View Full Execution Trace] [Export Code + Results]       │
├────────────────────────────────────────────────────────────┤
│  📊 Summary Statistics                                     │
│                                                            │
│  Total Calls:      94       (47 iterations × 2 features)   │
│  Success:          92       (97.9%)                        │
│  Failures:         2        (2.1%)                         │
│                                                            │
│  Failure Breakdown:                                        │
│  • tps_exceeded:   2        (pdf_export TPS limit)         │
│  • quota_exceeded: 0                                       │
│  • tier_denied:    0                                       │
│                                                            │
│  Average TPS:                                              │
│  • ml_analytics:   8.5 TPS  (max: 10.0)  ✓ OK             │
│  • pdf_export:     1.2 TPS  (max: 5.0)   ✓ OK             │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

#### Real-Time Updates

**Update Strategy:**
- WebSocket connection for events: `/api/instances/{id}/stream`
- HTTP polling for metrics: `/api/instances/{id}/status` every 500ms
- Chart.js with smooth line interpolation
- Event log with auto-scroll (toggleable)

**Update Elements:**

1. **Progress Bar**
   - Updates every iteration
   - Shows percentage and count
   - Estimated time remaining

2. **Charts** (Chart.js)
   - Success/Failure: Line chart, 2 datasets
   - Quota Remaining: Horizontal bar chart
   - TPS Monitor: Real-time gauge/number display
   - Updates every 5 iterations to reduce render overhead

3. **Event Log**
   - WebSocket pushes new events
   - Max 100 events displayed (rolling window)
   - Color-coded: ✓ Green / ⚠ Yellow / ✗ Red
   - Expandable details per event
   - Filter by: All / Success / Warnings / Errors

4. **Code Context**
   - Shows current iteration's code execution
   - Highlights executed lines in real-time
   - Displays actual SDK return values inline
   - Previous iteration results shown as comments

5. **Status Badges**
   - Running: Green badge with spinner animation
   - Paused: Yellow badge, static
   - Stopped: Red badge
   - Error: Red badge with error icon

#### Interactive Features

**Control Panel:**
- Pause: Freeze execution, resume later
- Stop: End simulation, keep results
- Reset: Clear and start over
- Export: Download results as JSON/PDF

**Metrics Panel:**
- Expand/Collapse sections
- Hover tooltips with detailed info
- Click feature name to highlight in log

**Event Log:**
- Click event to expand details
- Right-click to copy
- Filter dropdown: All/Success/Warnings/Errors
- Search box for text search
- Export button → CSV/JSON

**Code Context:**
- Toggle between: Compact / Detailed / Off
- Copy code button
- Link to SDK docs for each API

#### API Endpoints

```
POST /api/instances/{id}/start
Response: {
  "status": "running",
  "started_at": "2025-01-21T11:23:00Z"
}

POST /api/instances/{id}/stop
Response: {
  "status": "stopped",
  "completed_iterations": 47,
  "total_iterations": 100
}

POST /api/instances/{id}/pause
POST /api/instances/{id}/resume

GET  /api/instances/{id}/status
Response: {
  "instance_id": "lcc-demo-abc123",
  "status": "running",
  "progress": {
    "current": 47,
    "total": 100,
    "percent": 47
  },
  "metrics": {
    "elapsed_sec": 23.5,
    "success_count": 45,
    "failure_count": 2,
    "current_tps": {
      "ml_analytics": 8.5,
      "pdf_export": 1.2
    },
    "quota_remaining": {
      "ml_analytics": 9953,
      "pdf_export": 198
    }
  },
  "last_event": { ... }
}

WebSocket /api/instances/{id}/stream
Event Format: {
  "timestamp": "2025-01-21T11:23:47.234Z",
  "iteration": 47,
  "feature": "ml_analytics",
  "api": "Consume",
  "result": "success",
  "details": {
    "allowed": true,
    "remaining": 9953,
    "reason": "ok"
  }
}

GET /api/instances/{id}/report
Response: {
  "summary": { ... },
  "events": [ ... ],
  "charts_data": { ... }
}
```

---

## Three-Tier Product Design

### Product: Data Insight Analytics Platform

A comprehensive data analytics SaaS with ML capabilities, reporting, and export features.

---

### Tier 1: Basic Edition

**Product ID:** `data-insight-basic`

**Target Segment:** Individual users, hobbyists, small projects

**Price Point:** Free or $9/month

**Features:**

```yaml
features:
  basic_reports:
    enabled: true
    description: Generate basic statistical reports
    # No limitations

  ml_analytics:
    enabled: false
    reason: requires_professional
    required_tier: professional

  pdf_export:
    enabled: false
    reason: requires_professional
    required_tier: professional

  excel_export:
    enabled: false
    reason: requires_enterprise
    required_tier: enterprise

  custom_dashboard:
    enabled: false
    reason: requires_enterprise
    required_tier: enterprise

  api_access:
    enabled: false
    reason: requires_professional
    required_tier: professional
```

**License File (JSON):**

```json
{
  "product_id": "data-insight-basic",
  "product_name": "Data Insight Basic",
  "tier": "basic",
  "version": "1.0.0",
  "issued_at": "2025-01-21T00:00:00Z",
  "expires_at": "2026-01-21T00:00:00Z",
  "features": {
    "basic_reports": {
      "enabled": true
    },
    "ml_analytics": {
      "enabled": false
    },
    "pdf_export": {
      "enabled": false
    },
    "excel_export": {
      "enabled": false
    },
    "custom_dashboard": {
      "enabled": false
    },
    "api_access": {
      "enabled": false
    }
  }
}
```

**YAML Configuration (lcc-features.yaml):**

```yaml
sdk:
  product_id: data-insight-basic
  product_version: "1.0.0"
  lcc_url: "http://localhost:7086"

features:
  - id: basic_reports
    name: Basic Statistical Reports
    intercept:
      package: reports
      function: GenerateBasicReport
    on_deny:
      action: error
      message: "Report generation failed"

  - id: ml_analytics
    name: ML-Powered Analytics
    intercept:
      package: analytics
      function: RunMLAnalysis
    on_deny:
      action: error
      message: "ML Analytics requires Professional tier or higher"

  - id: pdf_export
    name: PDF Export
    intercept:
      package: exports
      function: ExportToPDF
    on_deny:
      action: error
      message: "PDF Export requires Professional tier or higher"

  - id: excel_export
    name: Excel Export
    intercept:
      package: exports
      function: ExportToExcel
    on_deny:
      action: error
      message: "Excel Export requires Enterprise tier"

  - id: custom_dashboard
    name: Custom Dashboard Builder
    intercept:
      package: dashboards
      function: CreateCustomDashboard
    on_deny:
      action: error
      message: "Custom dashboards require Enterprise tier"

  - id: api_access
    name: REST API Access
    intercept:
      package: api
      function: HandleAPIRequest
    on_deny:
      action: error
      message: "API access requires Professional tier or higher"
```

---

### Tier 2: Professional Edition

**Product ID:** `data-insight-pro`

**Target Segment:** Small teams, growing businesses, power users

**Price Point:** $49/month or $490/year

**Features:**

```yaml
features:
  basic_reports:
    enabled: true

  ml_analytics:
    enabled: true
    quota:
      max: 10000
      window: daily
      reset_at: "00:00"
    max_tps: 10.0

  pdf_export:
    enabled: true
    quota:
      max: 200
      window: daily
    max_tps: 5.0

  api_access:
    enabled: true
    max_tps: 100.0
    max_concurrency: 10

  excel_export:
    enabled: false
    reason: requires_enterprise
    required_tier: enterprise

  custom_dashboard:
    enabled: false
    reason: requires_enterprise
    required_tier: enterprise
```

**License File (JSON):**

```json
{
  "product_id": "data-insight-pro",
  "product_name": "Data Insight Professional",
  "tier": "professional",
  "version": "1.0.0",
  "issued_at": "2025-01-21T00:00:00Z",
  "expires_at": "2026-01-21T00:00:00Z",
  "features": {
    "basic_reports": {
      "enabled": true
    },
    "ml_analytics": {
      "enabled": true,
      "quota": {
        "max": 10000,
        "used": 0,
        "remaining": 10000,
        "window": "daily",
        "reset_at": "2025-01-22T00:00:00Z"
      },
      "max_tps": 10.0
    },
    "pdf_export": {
      "enabled": true,
      "quota": {
        "max": 200,
        "used": 0,
        "remaining": 200,
        "window": "daily",
        "reset_at": "2025-01-22T00:00:00Z"
      },
      "max_tps": 5.0
    },
    "api_access": {
      "enabled": true,
      "max_tps": 100.0,
      "max_concurrency": 10
    },
    "excel_export": {
      "enabled": false
    },
    "custom_dashboard": {
      "enabled": false
    }
  }
}
```

**YAML Configuration:** (Same structure as Basic, SDK reads limits from license)

---

### Tier 3: Enterprise Edition

**Product ID:** `data-insight-enterprise`

**Target Segment:** Large teams, enterprises, mission-critical deployments

**Price Point:** $299/month or $2,990/year

**Features:**

```yaml
features:
  basic_reports:
    enabled: true

  ml_analytics:
    enabled: true
    quota:
      max: 100000
      window: daily
      reset_at: "00:00"
    max_tps: 50.0

  pdf_export:
    enabled: true
    quota:
      max: 2000
      window: daily
    max_tps: 20.0

  excel_export:
    enabled: true
    quota:
      max: 1000
      window: daily
    max_tps: 10.0

  custom_dashboard:
    enabled: true
    max_capacity: 100

  api_access:
    enabled: true
    max_tps: 500.0
    max_concurrency: 50
```

**License File (JSON):**

```json
{
  "product_id": "data-insight-enterprise",
  "product_name": "Data Insight Enterprise",
  "tier": "enterprise",
  "version": "1.0.0",
  "issued_at": "2025-01-21T00:00:00Z",
  "expires_at": "2026-01-21T00:00:00Z",
  "features": {
    "basic_reports": {
      "enabled": true
    },
    "ml_analytics": {
      "enabled": true,
      "quota": {
        "max": 100000,
        "used": 0,
        "remaining": 100000,
        "window": "daily",
        "reset_at": "2025-01-22T00:00:00Z"
      },
      "max_tps": 50.0
    },
    "pdf_export": {
      "enabled": true,
      "quota": {
        "max": 2000,
        "used": 0,
        "remaining": 2000,
        "window": "daily",
        "reset_at": "2025-01-22T00:00:00Z"
      },
      "max_tps": 20.0
    },
    "excel_export": {
      "enabled": true,
      "quota": {
        "max": 1000,
        "used": 0,
        "remaining": 1000,
        "window": "daily",
        "reset_at": "2025-01-22T00:00:00Z"
      },
      "max_tps": 10.0
    },
    "custom_dashboard": {
      "enabled": true,
      "max_capacity": 100
    },
    "api_access": {
      "enabled": true,
      "max_tps": 500.0,
      "max_concurrency": 50
    }
  }
}
```

---

### Feature ID Summary

| Feature ID | Display Name | Basic | Pro | Enterprise | Control Types |
|-----------|--------------|-------|-----|------------|---------------|
| `basic_reports` | Basic Reports | ✓ | ✓ | ✓ | Tier only |
| `ml_analytics` | ML Analytics | ✗ | ✓ (10K/day, 10 TPS) | ✓ (100K/day, 50 TPS) | Tier + Quota + TPS |
| `pdf_export` | PDF Export | ✗ | ✓ (200/day, 5 TPS) | ✓ (2K/day, 20 TPS) | Tier + Quota + TPS |
| `excel_export` | Excel Export | ✗ | ✗ | ✓ (1K/day, 10 TPS) | Tier + Quota + TPS |
| `custom_dashboard` | Custom Dashboards | ✗ | ✗ | ✓ (100 max) | Tier + Capacity |
| `api_access` | API Access | ✗ | ✓ (100 TPS, 10 users) | ✓ (500 TPS, 50 users) | Tier + TPS + Concurrency |

---

## API Specifications

### Base URL Structure

```
http://localhost:9144/api
```

### Authentication

Demo app uses auto-generated RSA key pairs managed by SDK. No additional auth required for demo purposes.

### Common Response Format

**Success:**
```json
{
  "ok": true,
  "data": { ... }
}
```

**Error:**
```json
{
  "ok": false,
  "error": "error_code",
  "message": "Human-readable error message"
}
```

### Endpoint Index

#### Configuration (Page 1)

```
GET  /api/config
POST /api/config
GET  /api/config/validate
```

#### Tiers (Page 2)

```
GET  /api/tiers
GET  /api/tiers/{tier}/license
GET  /api/tiers/{tier}/yaml
POST /api/tiers/{tier}/check-feature
```

#### Limits (Page 3)

```
GET  /api/limits/types
GET  /api/limits/{type}/example
POST /api/limits/{type}/simulate
```

#### Instances (Page 4)

```
POST /api/instances/create
GET  /api/instances
GET  /api/instances/{id}
DELETE /api/instances/{id}
```

#### Runtime (Page 5)

```
POST   /api/instances/{id}/start
POST   /api/instances/{id}/stop
POST   /api/instances/{id}/pause
POST   /api/instances/{id}/resume
GET    /api/instances/{id}/status
GET    /api/instances/{id}/report
WebSocket /api/instances/{id}/stream
```

---

## Design System

### Color Palette

#### Dark Theme (Primary)

```css
/* Background Colors */
--bg-primary:   #0b1220;  /* Main background */
--bg-secondary: #1e293b;  /* Cards, panels */
--bg-tertiary:  #334155;  /* Hover states */

/* Border Colors */
--border-subtle: #1e293b;
--border:        #334155;
--border-accent: #475569;

/* Text Colors */
--text-primary:   #e2e8f0;  /* Main text */
--text-secondary: #cbd5e1;  /* Headings */
--text-muted:     #94a3b8;  /* Less important */
--text-disabled:  #64748b;  /* Disabled state */

/* Accent & Interactive */
--accent:       #60a5fa;  /* Primary blue */
--accent-hover: #3b82f6;
--accent-dim:   rgba(96, 165, 250, 0.15);

/* Status Colors */
--success:      #6ee7b7;  /* Green */
--success-bg:   rgba(110, 231, 183, 0.15);
--warning:      #fbbf24;  /* Yellow */
--warning-bg:   rgba(251, 191, 36, 0.15);
--error:        #f87171;  /* Red */
--error-bg:     rgba(248, 113, 113, 0.15);
--info:         #a78bfa;  /* Purple */
--info-bg:      rgba(167, 139, 250, 0.15);
```

#### Light Theme (Optional)

```css
/* Can be implemented later for accessibility */
```

### Typography

**Font Stack:**
```css
--font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, 
             "Helvetica Neue", Arial, sans-serif;
--font-mono: "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", 
             Consolas, monospace;
```

**Type Scale:**
```css
--text-xs:   12px;  /* Labels, captions */
--text-sm:   14px;  /* Body small */
--text-base: 16px;  /* Body text */
--text-lg:   18px;  /* Emphasized */
--text-xl:   20px;  /* Section titles */
--text-2xl:  24px;  /* Page titles */
--text-3xl:  30px;  /* Hero text */
```

**Font Weights:**
```css
--weight-normal: 400;
--weight-medium: 500;
--weight-semibold: 600;
--weight-bold: 700;
```

### Spacing System

**Scale (8px base):**
```css
--space-1:  4px;
--space-2:  8px;
--space-3:  12px;
--space-4:  16px;
--space-5:  20px;
--space-6:  24px;
--space-8:  32px;
--space-10: 40px;
--space-12: 48px;
--space-16: 64px;
```

### Layout

**Container Widths:**
```css
--container-sm: 640px;
--container-md: 768px;
--container-lg: 1024px;
--container-xl: 1280px;
```

**Layout Variables:**
```css
--header-height: 64px;
--sidebar-width: 260px;
--panel-width: 800px;
--code-panel-width: 600px;
```

**Border Radius:**
```css
--radius-sm: 4px;
--radius:    8px;
--radius-lg: 12px;
--radius-xl: 16px;
--radius-full: 9999px;
```

**Shadows:**
```css
--shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
--shadow:    0 1px 3px rgba(0, 0, 0, 0.1), 
             0 1px 2px rgba(0, 0, 0, 0.06);
--shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1), 
             0 2px 4px rgba(0, 0, 0, 0.06);
--shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1), 
             0 4px 6px rgba(0, 0, 0, 0.05);
--shadow-xl: 0 20px 25px rgba(0, 0, 0, 0.1), 
             0 10px 10px rgba(0, 0, 0, 0.04);
```

### Components

#### Buttons

**Primary:**
```css
.btn-primary {
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  color: white;
  padding: 10px 20px;
  border-radius: var(--radius);
  font-weight: 600;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}
.btn-primary:hover {
  filter: brightness(1.1);
  box-shadow: 0 8px 20px rgba(59, 130, 246, 0.4);
  transform: translateY(-2px);
}
```

**Secondary:**
```css
.btn-secondary {
  background: var(--bg-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border);
  padding: 10px 20px;
  border-radius: var(--radius);
  transition: all 0.2s ease;
}
.btn-secondary:hover {
  background: var(--bg-tertiary);
  border-color: var(--accent);
}
```

#### Cards

```css
.card {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow);
  transition: all 0.2s ease;
}
.card:hover {
  border-color: var(--border-accent);
  box-shadow: var(--shadow-md);
}
```

#### Code Blocks

```css
.code-block {
  background: #0d1117;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 16px;
  font-family: var(--font-mono);
  font-size: 14px;
  line-height: 1.6;
  overflow-x: auto;
}

/* Syntax highlighting via Prism.js */
.token.comment { color: #6a737d; }
.token.string { color: #9ecbff; }
.token.keyword { color: #f97583; }
.token.function { color: #b392f0; }
.token.number { color: #79b8ff; }
```

#### Badges

```css
.badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  border: 1px solid;
}

.badge-success {
  color: var(--success);
  background: var(--success-bg);
  border-color: var(--success);
}

.badge-running {
  color: var(--success);
  background: var(--success-bg);
  border-color: var(--success);
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
```

### Animations

```css
/* Page transitions */
.page-enter {
  opacity: 0;
  transform: translateX(20px);
}
.page-enter-active {
  opacity: 1;
  transform: translateX(0);
  transition: all 0.3s ease;
}

/* Pulse animation for running status */
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

/* Spinner */
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Smooth chart updates */
.chart-line {
  transition: all 0.2s ease;
}
```

---

## Technical Implementation

### Frontend Stack

**Core Technologies:**
- **HTML5 + CSS3**: Semantic markup, modern layouts
- **Vanilla JavaScript (ES6+)**: No framework overhead, direct DOM manipulation
- **WebSocket API**: Real-time event streaming

**Libraries:**
- **Prism.js**: Syntax highlighting for code examples
- **Chart.js**: Real-time charts and visualizations
- **No jQuery**: Modern DOM APIs are sufficient

**Build:**
- Go `embed` package to bundle all static assets
- Single binary deployment
- No build step required for development

### Backend Architecture

**Language:** Go 1.21+

**Structure:**
```
cmd/
  webdemo/
    main.go              # Entry point
internal/
  web/
    server.go            # HTTP server
    handlers.go          # Page handlers
    api_config.go        # Config API
    api_tiers.go         # Tier API
    api_limits.go        # Limits API
    api_instances.go     # Instance API
    api_runtime.go       # Runtime API
    websocket.go         # WebSocket handler
    products.go          # Product definitions
    simulation.go        # Simulation engine
  storage/
    config.go            # Config persistence
    instances.go         # Instance storage
static/
  index.html
  styles.css
  app.js
  prism.js
  chart.js
docs/
  UI_DESIGN_SPEC.md     # This document
```

**Data Storage:**
- Configuration: JSON file in `~/.lcc-demo/config.json`
- Instances: In-memory during runtime, optional JSON export
- No database required for demo

**Concurrency:**
- Goroutines for simulation execution
- Mutex-protected shared state
- WebSocket broadcast channel

### State Management

**Server-Side:**
```go
type Server struct {
  mu            sync.RWMutex
  config        *Config
  instances     map[string]*Instance
  wsClients     map[string][]*websocket.Conn
}

type Instance struct {
  ID            string
  Tier          string
  Status        InstanceStatus
  Config        *InstanceConfig
  Metrics       *Metrics
  EventLog      *RingBuffer
  cancelFunc    context.CancelFunc
}
```

**Client-Side:**
```javascript
const AppState = {
  config: null,
  currentPage: 'welcome',
  selectedTier: null,
  instanceConfig: null,
  runtimeData: null,
  wsConnection: null
};
```

### WebSocket Protocol

**Connection:**
```javascript
const ws = new WebSocket(`ws://localhost:9144/api/instances/${instanceId}/stream`);
```

**Event Types:**
```json
{
  "type": "iteration",
  "timestamp": "2025-01-21T11:23:47.234Z",
  "iteration": 47,
  "feature": "ml_analytics",
  "api": "Consume",
  "result": "success",
  "details": {
    "allowed": true,
    "remaining": 9953,
    "reason": "ok"
  }
}

{
  "type": "status",
  "status": "running",
  "progress": 47,
  "total": 100
}

{
  "type": "complete",
  "summary": { ... }
}
```

### Code Highlighting

**Prism.js Configuration:**
```javascript
Prism.highlightAll();

// Custom theme for Go
Prism.languages.go = {
  'comment': /\/\/.*/,
  'string': /(["'`])(?:\\[\s\S]|(?!\1)[^\\])*\1/,
  'keyword': /\b(?:func|if|else|return|var|const|type|struct|interface)\b/,
  'function': /\b\w+(?=\()/,
  'number': /\b0x[\da-f]+\b|(?:\b\d+\.?\d*|\B\.\d+)(?:e[+-]?\d+)?/i,
  'operator': /[*\/%^!=]=?|[\-+<>]=?|&&|\|\|/,
};
```

### Chart Configuration

**Chart.js Setup:**
```javascript
const successChart = new Chart(ctx, {
  type: 'line',
  data: {
    labels: [],
    datasets: [
      {
        label: 'Success',
        data: [],
        borderColor: '#6ee7b7',
        backgroundColor: 'rgba(110, 231, 183, 0.1)',
        tension: 0.4
      },
      {
        label: 'Failures',
        data: [],
        borderColor: '#f87171',
        backgroundColor: 'rgba(248, 113, 113, 0.1)',
        tension: 0.4
      }
    ]
  },
  options: {
    responsive: true,
    animation: {
      duration: 200
    },
    scales: {
      y: {
        beginAtZero: true
      }
    }
  }
});

// Update every 5 iterations
function updateChart(iteration, success, failures) {
  if (iteration % 5 === 0) {
    successChart.data.labels.push(iteration);
    successChart.data.datasets[0].data.push(success);
    successChart.data.datasets[1].data.push(failures);
    successChart.update();
  }
}
```

### Performance Optimization

**Frontend:**
- Debounce input handlers (300ms)
- Virtual scrolling for event log (only render visible items)
- Chart update throttling (every 5 iterations)
- WebSocket message batching

**Backend:**
- Connection pooling for LCC HTTP client
- Response caching with TTL
- Goroutine pool for simulations (max 10 concurrent)
- Memory-efficient ring buffer for event logs

---

## Future Enhancements

### Version 1.1
- [ ] Export simulation reports as PDF
- [ ] Share simulation config via URL
- [ ] Dark/Light theme toggle
- [ ] Responsive mobile layout

### Version 1.2
- [ ] Multiple instance comparison view
- [ ] Advanced simulation scenarios (burst traffic, etc)
- [ ] Custom feature definitions
- [ ] Simulation replay from saved state

### Version 2.0
- [ ] Multi-language support (i18n)
- [ ] Collaborative mode (multiple users)
- [ ] Integration with CI/CD for automated testing
- [ ] Plugin system for custom visualizations

---

## Appendix

### Design Inspirations

- **Stripe Dashboard**: Clean layout, excellent code examples
- **Tailwind UI**: Component patterns, dark theme
- **GitLab CI/CD**: Pipeline visualization, real-time logs
- **Vercel Dashboard**: Status badges, deployment logs

### Accessibility Considerations

- ARIA labels for interactive elements
- Keyboard navigation support
- Screen reader friendly
- High contrast mode support
- Focus indicators

### Browser Support

- Chrome/Edge: Latest 2 versions
- Firefox: Latest 2 versions
- Safari: Latest 2 versions
- No IE11 support (uses modern JS features)

---

**Document Version:** 1.0.0  
**Last Updated:** 2025-01-21  
**Maintained By:** LCC Demo App Team
