# LCC SDK Demo App - Implementation Plan

**Version:** 1.0.0  
**Date:** 2025-01-21  
**Based on:** [UI_DESIGN_SPEC.md](./UI_DESIGN_SPEC.md)

---

## 📋 Executive Summary

This document outlines the implementation strategy for refactoring the LCC SDK Demo App according to the new UI design specification. We adopt a **progressive refactoring** approach to minimize risk while maximizing code reuse.

**Key Decision:** Progressive refactoring instead of complete rewrite

**Timeline:** 5 weeks (1 page per week)

**Risk Level:** Low (preserves stable infrastructure)

---

## 📊 Current Code Assessment

### Existing Codebase Analysis

**Total Lines of Code:**
- `internal/web/server.go`: 979 lines
- `cmd/webdemo/main.go`: 560 lines
- Total: ~1,500 lines

### Reusability Breakdown

#### ✅ Directly Reusable (40%)

**Infrastructure Layer:**
- `Server` struct and routing framework
- HTTP handler patterns
- Configuration management (`handleConfig`, `loadConfig`, `saveConfig`)
- KeyStore management (RSA key pair persistence)
- LCC Client initialization
- PublicServiceClient (LCC public API calls)

**Files to Keep Unchanged:**
```
internal/web/keystore.go           - 100% reuse
internal/web/public_client.go      - 100% reuse
internal/web/feature_source.go     - 100% reuse
internal/analytics/analytics.go    - 100% reuse
internal/export/export.go          - 100% reuse
internal/reporting/reporting.go    - 100% reuse
```

#### 🔄 Refactorable (30%)

**Business Logic Layer:**
- Product definition structures (`Product`, `ProductFeature`, `ProductLimit`)
- LCC SDK integration patterns
- Feature checking and Consume logic
- Existing `analytics` and `export` modules

**Files to Refactor:**
```
internal/web/server.go
  ├─ Keep: Server struct, Router, handleConfig, loadConfig, saveConfig
  ├─ Refactor: routes() - add new routes
  └─ Remove: handleIndex, handleProductPage (move to v2)
```

#### ❌ Must Rewrite (30%)

**Frontend (100% rewrite):**
- All HTML pages (new 5-page design)
- CSS (new design system)
- JavaScript (new SPA architecture)

**API Endpoints (60% rewrite):**
- `/api/tiers/*` - NEW
- `/api/limits/*` - NEW
- `/api/instances/*` - NEW (replaces `/api/simulation/*`)
- WebSocket support - NEW

**Simulation Engine (70% rewrite):**
- More complex scenarios
- Real-time event streaming
- Detailed metrics tracking

---

## 🎯 Implementation Strategy

### Why Progressive Refactoring?

✅ **Advantages:**
1. Preserve stable, tested infrastructure
2. Gradual migration reduces risk
3. Each step is testable and reversible
4. Reuse ~40% of existing code
5. Keep debugging tools available

❌ **Complete Rewrite Disadvantages:**
1. Wastes validated base code
2. Requires re-debugging all infrastructure
3. Higher time cost
4. Introduces more uncertainty

---

## 📅 5-Week Implementation Timeline

### Week 1: Foundation & Page 1 (Welcome)

**Goals:**
- Set up new directory structure
- Implement configuration page
- Reuse existing config management

**Tasks:**
1. Create new directories
2. Backup existing code (`git branch backup-old-ui`)
3. Implement Page 1 frontend (HTML/CSS/JS)
4. Reuse `/api/config` endpoint
5. Add `/api/config/validate` endpoint

**Deliverables:**
- Working Welcome page
- LCC connection test
- Configuration persistence

---

### Week 2: Page 2 (Tier Learning)

**Goals:**
- Implement tier comparison page
- Define three-tier product structure
- Interactive tier switching

**Tasks:**
1. Create `internal/web/products.go` with tier definitions
2. Implement `/api/tiers` endpoint
3. Implement `/api/tiers/{tier}/license` endpoint
4. Implement `/api/tiers/{tier}/yaml` endpoint
5. Implement `/api/tiers/{tier}/check-feature` endpoint
6. Build Page 2 frontend with tabs and code examples
7. Add Prism.js for syntax highlighting

**Deliverables:**
- Working Tier Learning page
- Three product definitions (Basic/Pro/Enterprise)
- Interactive tier comparison
- Try It Yourself simulator

---

### Week 3: Page 3 (Limits Learning)

**Goals:**
- Implement four limit types education
- Interactive simulations for each type
- Code examples with highlighting

**Tasks:**
1. Implement `/api/limits/types` endpoint
2. Implement `/api/limits/{type}/example` endpoint
3. Implement `/api/limits/{type}/simulate` endpoint
4. Build Page 3 frontend with tab interface
5. Create mini-simulators for each limit type
6. Add comparison table

**Deliverables:**
- Working Limits Learning page
- Four limit type explanations (Quota/TPS/Capacity/Concurrency)
- Interactive mini-simulations
- Code examples for each type

---

### Week 4: Page 4 (Instance Setup)

**Goals:**
- Implement instance configuration page
- License selection and preview
- Runtime parameter configuration

**Tasks:**
1. Create `internal/web/instances.go` for instance management
2. Implement `/api/instances/create` endpoint
3. Implement `/api/instances` (list) endpoint
4. Implement `/api/instances/{id}` (get) endpoint
5. Build Page 4 frontend with 4-step wizard
6. Add real-time prediction logic

**Deliverables:**
- Working Instance Setup page
- Instance creation and storage
- Configuration validation
- Behavior prediction

---

### Week 5: Page 5 (Runtime Dashboard)

**Goals:**
- Implement live simulation dashboard
- Real-time metrics and charts
- WebSocket event streaming

**Tasks:**
1. Refactor simulation engine in `internal/web/simulation.go`
2. Implement WebSocket handler in `internal/web/websocket.go`
3. Implement `/api/instances/{id}/start` endpoint
4. Implement `/api/instances/{id}/stop` endpoint
5. Implement `/api/instances/{id}/status` endpoint
6. Implement WebSocket `/api/instances/{id}/stream`
7. Build Page 5 frontend with Chart.js
8. Add event log with filtering
9. Add code context tracing

**Deliverables:**
- Working Runtime Dashboard
- Real-time charts (Chart.js)
- WebSocket event streaming
- Live code execution tracing
- Export functionality

---

## 🗂️ New Directory Structure

```
lcc-demo-app/
├── cmd/
│   ├── demo/              # Keep: CLI demo tool
│   ├── webdemo/           # Refactor: Main web entry point
│   └── regression/        # Keep: Testing tool
├── internal/
│   ├── web/
│   │   ├── server.go            # Refactor: Add new routes
│   │   ├── keystore.go          # Keep: Unchanged
│   │   ├── public_client.go     # Keep: Unchanged
│   │   ├── feature_source.go    # Keep: Unchanged
│   │   ├── products.go          # NEW: Three-tier definitions
│   │   ├── tiers_handler.go     # NEW: Page 2 API
│   │   ├── limits_handler.go    # NEW: Page 3 API
│   │   ├── instances.go         # NEW: Instance management
│   │   ├── instances_handler.go # NEW: Page 4/5 API
│   │   ├── simulation.go        # NEW: Refactored engine
│   │   └── websocket.go         # NEW: WebSocket support
│   ├── analytics/         # Keep: Unchanged
│   ├── export/            # Keep: Unchanged
│   └── reporting/         # Keep: Unchanged
├── static/                # NEW: Complete rewrite
│   ├── index.html
│   ├── css/
│   │   └── styles.css
│   ├── js/
│   │   ├── app.js
│   │   ├── pages/
│   │   │   ├── welcome.js
│   │   │   ├── tiers.js
│   │   │   ├── limits.js
│   │   │   ├── setup.js
│   │   │   └── runtime.js
│   │   └── utils.js
│   └── vendor/
│       ├── prism.js
│       └── chart.js
├── docs/
│   ├── UI_DESIGN_SPEC.md        # Reference
│   └── IMPLEMENTATION_PLAN.md   # This document
└── README.md
```

---

## 🔧 Technical Implementation Details

### Backend Structure

#### New Files to Create

**1. `internal/web/products.go`**
```go
package web

// Define three-tier product structure
type TierDefinition struct {
    ID          string
    Name        string
    Tier        string
    Features    map[string]FeatureConfig
}

var (
    BasicTier       *TierDefinition
    ProfessionalTier *TierDefinition
    EnterpriseTier  *TierDefinition
)

func init() {
    // Initialize three tier definitions
    // Based on UI_DESIGN_SPEC.md section "Three-Tier Product Design"
}
```

**2. `internal/web/tiers_handler.go`**
```go
package web

// GET /api/tiers
func (s *Server) handleGetTiers(w http.ResponseWriter, r *http.Request)

// GET /api/tiers/{tier}/license
func (s *Server) handleGetTierLicense(w http.ResponseWriter, r *http.Request)

// GET /api/tiers/{tier}/yaml
func (s *Server) handleGetTierYAML(w http.ResponseWriter, r *http.Request)

// POST /api/tiers/{tier}/check-feature
func (s *Server) handleCheckTierFeature(w http.ResponseWriter, r *http.Request)
```

**3. `internal/web/limits_handler.go`**
```go
package web

// GET /api/limits/types
func (s *Server) handleGetLimitTypes(w http.ResponseWriter, r *http.Request)

// GET /api/limits/{type}/example
func (s *Server) handleGetLimitExample(w http.ResponseWriter, r *http.Request)

// POST /api/limits/{type}/simulate
func (s *Server) handleSimulateLimit(w http.ResponseWriter, r *http.Request)
```

**4. `internal/web/instances.go`**
```go
package web

import (
    "context"
    "sync"
)

type Instance struct {
    ID            string
    Name          string
    Tier          string
    Status        InstanceStatus
    Config        *InstanceConfig
    Metrics       *Metrics
    EventLog      *RingBuffer
    cancelFunc    context.CancelFunc
}

type InstanceManager struct {
    mu        sync.RWMutex
    instances map[string]*Instance
}

func NewInstanceManager() *InstanceManager
func (m *InstanceManager) Create(config *InstanceConfig) (*Instance, error)
func (m *InstanceManager) Get(id string) (*Instance, error)
func (m *InstanceManager) List() []*Instance
func (m *InstanceManager) Delete(id string) error
```

**5. `internal/web/simulation.go`**
```go
package web

// Refactored simulation engine with:
// - Real-time event streaming
// - Detailed metrics tracking
// - Support for all four control types
// - WebSocket integration

type SimulationEngine struct {
    instance   *Instance
    lccClient  *client.Client
    eventChan  chan SimulationEvent
    stopChan   chan struct{}
}

func NewSimulationEngine(instance *Instance) *SimulationEngine
func (e *SimulationEngine) Start(ctx context.Context) error
func (e *SimulationEngine) Stop() error
func (e *SimulationEngine) Pause() error
func (e *SimulationEngine) Resume() error
```

**6. `internal/web/websocket.go`**
```go
package web

import "github.com/gorilla/websocket"

type WebSocketHub struct {
    clients    map[string]map[*websocket.Conn]bool
    broadcast  chan WebSocketMessage
    register   chan *ClientRegistration
    unregister chan *ClientRegistration
}

func NewWebSocketHub() *WebSocketHub
func (h *WebSocketHub) Run()
func (h *WebSocketHub) HandleWebSocket(w http.ResponseWriter, r *http.Request)
```

### Frontend Structure

#### Single Page Application (SPA)

**Main Entry: `static/index.html`**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>LCC SDK Interactive Learning Demo</title>
    <link rel="stylesheet" href="/static/css/styles.css">
    <link rel="stylesheet" href="/static/vendor/prism.css">
</head>
<body>
    <div id="app">
        <!-- Dynamic content loaded here -->
    </div>
    
    <script src="/static/vendor/prism.js"></script>
    <script src="/static/vendor/chart.js"></script>
    <script src="/static/js/utils.js"></script>
    <script src="/static/js/pages/welcome.js"></script>
    <script src="/static/js/pages/tiers.js"></script>
    <script src="/static/js/pages/limits.js"></script>
    <script src="/static/js/pages/setup.js"></script>
    <script src="/static/js/pages/runtime.js"></script>
    <script src="/static/js/app.js"></script>
</body>
</html>
```

**Main App Logic: `static/js/app.js`**
```javascript
const App = {
    state: {
        currentPage: 'welcome',
        config: null,
        selectedTier: null,
        instance: null,
        wsConnection: null
    },
    
    init() {
        this.setupRouter();
        this.loadConfig();
        this.renderCurrentPage();
    },
    
    setupRouter() {
        window.addEventListener('hashchange', () => {
            this.renderCurrentPage();
        });
    },
    
    renderCurrentPage() {
        const page = window.location.hash.slice(1) || 'welcome';
        this.state.currentPage = page;
        
        // Load page module and render
        switch(page) {
            case 'welcome':
                WelcomePage.render();
                break;
            case 'tiers':
                TiersPage.render();
                break;
            case 'limits':
                LimitsPage.render();
                break;
            case 'setup':
                SetupPage.render();
                break;
            case 'runtime':
                RuntimePage.render();
                break;
        }
    }
};

document.addEventListener('DOMContentLoaded', () => {
    App.init();
});
```

**Page Module Example: `static/js/pages/welcome.js`**
```javascript
const WelcomePage = {
    render() {
        const html = `
            <div class="page-welcome">
                <h1>🎓 LCC SDK Interactive Learning Demo</h1>
                <div class="card">
                    <h2>🔧 LCC Server Configuration</h2>
                    <input type="text" id="lcc-url" placeholder="http://localhost:7086">
                    <button onclick="WelcomePage.testConnection()">Test Connection</button>
                    <button onclick="WelcomePage.saveAndContinue()">Save & Continue →</button>
                </div>
            </div>
        `;
        document.getElementById('app').innerHTML = html;
        this.loadConfig();
    },
    
    async loadConfig() {
        const resp = await fetch('/api/config');
        const data = await resp.json();
        if (data.lcc_url) {
            document.getElementById('lcc-url').value = data.lcc_url;
        }
    },
    
    async testConnection() {
        const url = document.getElementById('lcc-url').value;
        const resp = await fetch('/api/config/validate?url=' + encodeURIComponent(url));
        const data = await resp.json();
        // Show result
    },
    
    async saveAndContinue() {
        const url = document.getElementById('lcc-url').value;
        await fetch('/api/config', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({lcc_url: url})
        });
        window.location.hash = 'tiers';
    }
};
```

---

## 🧪 Testing Strategy

### Per-Page Testing

**Week 1: Page 1**
- [ ] Config save/load works
- [ ] LCC connection validation works
- [ ] Navigation to Page 2 works

**Week 2: Page 2**
- [ ] All three tiers load correctly
- [ ] Tab switching works
- [ ] License JSON displays correctly
- [ ] Try It Yourself simulator works

**Week 3: Page 3**
- [ ] All four limit types display
- [ ] Tab switching works
- [ ] Code examples highlight correctly
- [ ] Mini-simulations run successfully

**Week 4: Page 4**
- [ ] Instance creation works
- [ ] License preview displays
- [ ] Parameter validation works
- [ ] Prediction displays correctly

**Week 5: Page 5**
- [ ] Simulation starts/stops correctly
- [ ] WebSocket events stream
- [ ] Charts update in real-time
- [ ] Event log displays correctly
- [ ] Code tracing works

### Integration Testing

After Week 5:
- [ ] Full user flow: Welcome → Tiers → Limits → Setup → Runtime
- [ ] All API endpoints respond correctly
- [ ] WebSocket reconnection works
- [ ] Browser refresh preserves state
- [ ] Multiple instances can run
- [ ] Export functionality works

---

## 🚨 Risk Mitigation

### Potential Risks & Solutions

**Risk 1: WebSocket complexity**
- **Mitigation:** Start with HTTP polling, add WebSocket in Week 5
- **Fallback:** Keep HTTP polling as backup

**Risk 2: Chart.js performance**
- **Mitigation:** Update charts every 5 iterations, not every event
- **Fallback:** Use simpler visualizations if needed

**Risk 3: State management complexity**
- **Mitigation:** Keep state simple, use vanilla JS patterns
- **Fallback:** Consider adding a lightweight state library if needed

**Risk 4: Browser compatibility**
- **Mitigation:** Test on Chrome, Firefox, Safari during development
- **Fallback:** Add polyfills if needed

---

## 📈 Success Criteria

### Functional Requirements
- ✅ All 5 pages work independently
- ✅ Navigation between pages works smoothly
- ✅ All API endpoints respond correctly
- ✅ Real-time updates work via WebSocket
- ✅ Code examples display with syntax highlighting
- ✅ Simulations run without errors

### Non-Functional Requirements
- ✅ Page load time < 1 second
- ✅ API response time < 100ms (excluding LCC calls)
- ✅ WebSocket latency < 50ms
- ✅ Chart updates smooth (60fps target)
- ✅ No memory leaks during long simulations

### User Experience
- ✅ Clear learning progression (5 pages)
- ✅ Interactive elements responsive
- ✅ Error messages helpful
- ✅ Code examples easy to read
- ✅ Simulations intuitive to configure

---

## 🔄 Deployment Strategy

### Local Development
```bash
# Run backend
cd cmd/webdemo
go run main.go

# Backend serves static files via embed
# Access at http://localhost:9144
```

### Production Build
```bash
# Build single binary with embedded static files
go build -o lcc-demo-app cmd/webdemo/main.go

# Run
./lcc-demo-app
```

### Docker (Optional)
```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o lcc-demo-app cmd/webdemo/main.go

FROM alpine:latest
COPY --from=builder /app/lcc-demo-app /usr/local/bin/
EXPOSE 9144
CMD ["lcc-demo-app"]
```

---

## 📚 Documentation Updates

### Documents to Update

**README.md**
- Add quick start guide
- Add architecture overview
- Add screenshots (after implementation)

**New Documents**
- API_REFERENCE.md (detailed API docs)
- DEVELOPER_GUIDE.md (how to extend)
- USER_GUIDE.md (end-user documentation)

---

## 🎓 Learning Resources

### For Developers Working on This

**Go Backend:**
- [Go HTTP server patterns](https://golang.org/pkg/net/http/)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [Go embed directive](https://pkg.go.dev/embed)

**Frontend:**
- [Vanilla JS SPA patterns](https://developer.mozilla.org/en-US/docs/Web/API/History_API)
- [Prism.js documentation](https://prismjs.com/)
- [Chart.js documentation](https://www.chartjs.org/)
- [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)

---

## 🚀 Next Steps

### Before Starting New Conversation

1. ✅ Read UI_DESIGN_SPEC.md thoroughly
2. ✅ Read this IMPLEMENTATION_PLAN.md
3. ✅ Backup existing code: `git branch backup-old-ui`
4. ✅ Create feature branch: `git checkout -b feature/new-ui-v2`

### In New Conversation

**Start with:**
> "I want to implement the LCC Demo App according to:
> - UI_DESIGN_SPEC.md at `/home/fila/jqdDev_2025/lcc-demo-app/docs/UI_DESIGN_SPEC.md`
> - IMPLEMENTATION_PLAN.md at `/home/fila/jqdDev_2025/lcc-demo-app/docs/IMPLEMENTATION_PLAN.md`
> 
> Let's start with Week 1: Foundation & Page 1 (Welcome)"

**Focus on:**
- One week/page at a time
- Complete each page before moving to next
- Test thoroughly before proceeding
- Ask questions if anything is unclear

---

## 📞 Support & Questions

If you encounter issues during implementation:

1. **Design Questions:** Refer to UI_DESIGN_SPEC.md
2. **Implementation Questions:** Refer to this document
3. **API Questions:** Check existing `/api/config` for patterns
4. **SDK Questions:** Check existing `internal/web/server.go` for LCC SDK usage

---

**Document Version:** 1.0.0  
**Last Updated:** 2025-01-21  
**Status:** Ready for Implementation

---

**Let's build an amazing educational tool! 🚀**
