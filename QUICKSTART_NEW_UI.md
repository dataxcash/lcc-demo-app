# LCC Demo App - New UI Quick Start Guide

## Overview

This is the new interactive learning platform for LCC SDK, featuring:
- ✨ Modern dark-themed UI
- 🎓 5-step progressive learning path
- 💻 Live code examples and simulations
- 📊 Real-time metrics and monitoring

## Prerequisites

- Go 1.21+ installed
- LCC Server running (optional for Week 1, but recommended)

## Quick Start

### 1. Build the Application

```bash
cd /home/fila/jqdDev_2025/lcc-demo-app
go build -o bin/web ./cmd/web
```

### 2. Start the Server

```bash
./bin/web
```

You should see:
```
LCC Demo Web UI listening on http://localhost:9144
```

### 3. Access the Application

Open your browser and navigate to:
```
http://localhost:9144/
```

## Current Status (Week 1)

### ✅ Implemented
- **Welcome Page**: LCC server configuration and connection testing
- **Design System**: Complete dark theme with animations
- **Navigation**: 5-step indicator and routing system
- **API Endpoints**:
  - `GET /api/config` - Get configuration
  - `POST /api/config` - Save configuration  
  - `GET /api/config/validate` - Test LCC connection

### 🚧 Coming Soon
- **Week 2**: Tier Learning page (Basic/Pro/Enterprise comparison)
- **Week 3**: Limits Learning page (Quota/TPS/Capacity/Concurrency)
- **Week 4**: Instance Setup page (Simulation configuration)
- **Week 5**: Runtime Dashboard (Live metrics and monitoring)

## Usage Guide

### Welcome Page

1. **Configure LCC Server**
   - Enter your LCC server URL (default: `http://localhost:7086`)
   - Click "Test Connection" to verify connectivity
   - If URL is changed, click "UPDATE" to save

2. **Connection Status**
   - ✅ Green indicator: Connected to LCC server
   - ❌ Red indicator: Connection failed
   - Status shows version and product count when connected

3. **Navigation**
   - Click "Save & Continue →" to proceed to next page
   - Use step indicators (1-5) at top to jump between pages
   - Use "← Back" and "Next →" buttons in footer

## Development

### Project Structure

```
lcc-demo-app/
├── cmd/
│   └── web/
│       └── main.go          # Entry point
├── internal/
│   └── web/
│       └── server.go        # Backend API server
├── static/
│   ├── index.html          # SPA entry
│   ├── css/
│   │   └── styles.css      # Design system
│   └── js/
│       ├── app.js          # Router
│       ├── utils.js        # Helpers
│       └── pages/
│           ├── welcome.js  # Week 1 ✅
│           ├── tiers.js    # Week 2 🚧
│           ├── limits.js   # Week 3 🚧
│           ├── setup.js    # Week 4 🚧
│           └── runtime.js  # Week 5 🚧
└── docs/
    ├── UI_DESIGN_SPEC.md      # Complete UI specification
    ├── IMPLEMENTATION_PLAN.md # 5-week implementation plan
    └── WEEK1_SUMMARY.md       # Week 1 completion summary
```

### Adding New Pages

To add or modify pages:

1. Edit the corresponding file in `static/js/pages/`
2. Implement the `render()` method
3. The app will automatically route to your page

Example:
```javascript
const MyPage = {
    async render() {
        const app = document.getElementById('app');
        app.innerHTML = `
            <div class="page-container">
                <h1>My Page Title</h1>
                <div class="card">
                    <!-- Your content here -->
                </div>
            </div>
        `;
        
        // Initialize page-specific logic
        this.init();
    },
    
    init() {
        // Setup event listeners, load data, etc.
    }
};
```

### API Development

To add new API endpoints, modify `internal/web/server.go`:

```go
// Add route in routes()
s.mux.HandleFunc("/api/myendpoint", s.handleMyEndpoint)

// Implement handler
func (s *Server) handleMyEndpoint(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // Your logic here
    json.NewEncoder(w).Encode(response)
}
```

## Troubleshooting

### Port Already in Use
If port 9144 is already in use:
```bash
# Kill existing process
pkill -f 'bin/web'

# Or use a different port
PORT=9145 ./bin/web
```

### LCC Connection Failed
- Make sure LCC server is running on the configured URL
- Check firewall settings
- Verify the URL format (include http:// or https://)

### Static Files Not Loading
- Ensure you're in the project root directory when running `./bin/web`
- Check that `static/` directory exists with all files
- Try rebuilding: `go build -o bin/web ./cmd/web`

## Accessing Old UI

The previous UI is preserved at:
```
http://localhost:9144/old/
```

## Testing

### Manual Testing
1. Start the server
2. Open browser to http://localhost:9144/
3. Test configuration form
4. Test connection validation
5. Test navigation between pages

### API Testing
```bash
# Get configuration
curl http://localhost:9144/api/config | jq

# Validate connection
curl http://localhost:9144/api/config/validate | jq

# Save configuration
curl -X POST http://localhost:9144/api/config \
  -H "Content-Type: application/json" \
  -d '{"lcc_url": "http://localhost:7086"}' | jq
```

## Next Steps

See [IMPLEMENTATION_PLAN.md](docs/IMPLEMENTATION_PLAN.md) for the complete 5-week roadmap.

Week 2 will add the **Tier Learning** page with:
- Three-tier product comparison (Basic/Pro/Enterprise)
- Interactive tier switching
- SDK code examples
- Feature availability matrix
- "Try It Yourself" simulator

## Support

For issues or questions:
1. Check [UI_DESIGN_SPEC.md](docs/UI_DESIGN_SPEC.md) for design details
2. Review [WEEK1_SUMMARY.md](docs/WEEK1_SUMMARY.md) for implementation notes
3. Examine the code - it's well-commented and structured

---

**Version:** 1.0.0 (Week 1)  
**Last Updated:** 2025-11-21  
**Status:** Foundation Complete ✅
