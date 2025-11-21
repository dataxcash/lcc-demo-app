package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/yourorg/lcc-sdk/pkg/auth"
	lccclient "github.com/yourorg/lcc-sdk/pkg/client"
	lccconfig "github.com/yourorg/lcc-sdk/pkg/config"
)

// Server provides Web UI and API for LCC demo simulation.
type Server struct {
	mux *http.ServeMux

	mu            sync.RWMutex
	lccURL        string
	publicBase    string
	clients       map[string]*lccclient.Client // productID -> client
	lastProducts  []PublicProduct              // cached latest products listing
}

func NewServer() *Server {
	s := &Server{
		mux:     http.NewServeMux(),
		clients: make(map[string]*lccclient.Client),
	}
	s.routes()
	s.loadConfig()
	return s
}

func (s *Server) Router() http.Handler { return s.mux }

func (s *Server) routes() {
	// New SPA UI
	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	s.mux.HandleFunc("/", s.handleSPA)
	
	// Old HTML pages (kept for backwards compatibility)
	s.mux.HandleFunc("/old/", s.handleIndex)
	s.mux.HandleFunc("/product/", s.handleProductPage)
	
	// API - Configuration
	s.mux.HandleFunc("/api/config", s.handleConfig)
	s.mux.HandleFunc("/api/config/validate", s.handleConfigValidate)
	
	// API - Tiers (Week 2)
	s.mux.HandleFunc("/api/tiers", s.handleGetTiers)
	s.mux.HandleFunc("/api/tiers/basic/license", s.handleGetTierLicense)
	s.mux.HandleFunc("/api/tiers/professional/license", s.handleGetTierLicense)
	s.mux.HandleFunc("/api/tiers/enterprise/license", s.handleGetTierLicense)
	s.mux.HandleFunc("/api/tiers/basic/yaml", s.handleGetTierYAML)
	s.mux.HandleFunc("/api/tiers/professional/yaml", s.handleGetTierYAML)
	s.mux.HandleFunc("/api/tiers/enterprise/yaml", s.handleGetTierYAML)
	s.mux.HandleFunc("/api/tiers/basic/check-feature", s.handleCheckTierFeature)
	s.mux.HandleFunc("/api/tiers/professional/check-feature", s.handleCheckTierFeature)
	s.mux.HandleFunc("/api/tiers/enterprise/check-feature", s.handleCheckTierFeature)
	
	// API - Products & Simulation
	s.mux.HandleFunc("/api/products", s.handleProducts)
	s.mux.HandleFunc("/api/features", s.handleFeatures)
	s.mux.HandleFunc("/api/sim/products", s.handleSimSelectProducts)
	s.mux.HandleFunc("/api/sim/registered", s.handleSimRegistered)
	// dynamic action routes: /api/sim/{product}/{action}
	s.mux.HandleFunc("/api/sim/", s.handleSimRoot)
}

// --- Handlers ---

type setConfigReq struct {
	LCCURL string `json:"lcc_url"`
}

type setConfigResp struct {
	OK     bool   `json:"ok"`
	LCCURL string `json:"lcc_url"`
}

func (s *Server) handleSPA(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		s.mu.RLock()
		url := s.lccURL
		if url == "" {
			url = "http://localhost:7086"
		}
		s.mu.RUnlock()
		_ = json.NewEncoder(w).Encode(map[string]any{
			"lcc_url":    url,
			"saved_at":   time.Now().Format(time.RFC3339),
			"is_default": s.lccURL == "",
		})
	case http.MethodPost:
		var req setConfigReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
			return
		}
		if req.LCCURL == "" {
			writeErr(w, http.StatusBadRequest, fmt.Errorf("lcc_url is required"))
			return
		}

		s.mu.Lock()
		s.lccURL = req.LCCURL
		s.publicBase = "/api/v1/public"
		s.mu.Unlock()

		if err := s.saveConfig(); err != nil {
			log.Printf("Failed to save config: %v", err)
		}

		_ = json.NewEncoder(w).Encode(&setConfigResp{OK: true, LCCURL: req.LCCURL})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type validateConfigResp struct {
	Reachable     bool   `json:"reachable"`
	Version       string `json:"version,omitempty"`
	ProductsCount int    `json:"products_count,omitempty"`
	Error         string `json:"error,omitempty"`
}

func (s *Server) handleConfigValidate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	s.mu.RLock()
	lccURL := s.lccURL
	s.mu.RUnlock()
	
	if lccURL == "" {
		lccURL = "http://localhost:7086"
	}

	base := strings.TrimRight(lccURL, "/") + "/api/v1/public"
	pc := NewPublicServiceClient(base)
	
	products, err := pc.ListProducts(r.Context())
	if err != nil {
		_ = json.NewEncoder(w).Encode(&validateConfigResp{
			Reachable: false,
			Error:     err.Error(),
		})
		return
	}
	
	// Try to determine version (you can enhance this based on actual LCC API)
	version := "v2.1.0" // default for now
	
	_ = json.NewEncoder(w).Encode(&validateConfigResp{
		Reachable:     true,
		Version:       version,
		ProductsCount: len(products),
	})
}

func (s *Server) handleProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	s.mu.RLock()
	lccURL := s.lccURL
	s.mu.RUnlock()
	if lccURL == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("lcc_url not configured; POST /api/config first"))
		return
	}
	base := strings.TrimRight(lccURL, "/") + "/api/v1/public"

	pc := NewPublicServiceClient(base)
	products, err := pc.ListProducts(r.Context())
	if err != nil {
		writeErr(w, http.StatusBadGateway, fmt.Errorf("failed to fetch products: %w", err))
		return
	}

	s.mu.Lock()
	s.lastProducts = products
	s.mu.Unlock()

	_ = json.NewEncoder(w).Encode(products)
}

type selectProductsReq struct {
	ProductIDs     []string `json:"product_ids"`
	DefaultVersion string   `json:"default_version,omitempty"`
}

type selectProductsResp struct {
	OK          bool              `json:"ok"`
	Registered  []string          `json:"registered"`
	Errors      map[string]string `json:"errors,omitempty"`
	InstanceIDs map[string]string `json:"instance_ids,omitempty"`
}

func (s *Server) handleSimSelectProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req selectProductsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
		return
	}
	if len(req.ProductIDs) == 0 {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("product_ids is required"))
		return
	}

	s.mu.RLock()
	lccURL := s.lccURL
	s.mu.RUnlock()
	if lccURL == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("lcc_url not configured"))
		return
	}

	registered := make([]string, 0, len(req.ProductIDs))
	errs := map[string]string{}
	instanceIDs := map[string]string{}

	ks, _ := NewKeyStore()

	for _, pid := range req.ProductIDs {
		cfg := &lccconfig.SDKConfig{
			LCCURL:         lccURL,
			ProductID:      pid,
			ProductVersion: nonEmpty(req.DefaultVersion, "1.0.0"),
			Timeout:        10 * time.Second,
			CacheTTL:       5 * time.Second,
		}

		var kp *auth.KeyPair
		if ks != nil {
			if loaded, err := ks.Load(pid); err == nil && loaded != nil {
				kp = loaded
			}
		}
		if kp == nil {
			// generate and persist
			var genErr error
			kp, genErr = auth.GenerateKeyPair()
			if genErr != nil {
				errs[pid] = fmt.Sprintf("keypair: %v", genErr)
				continue
			}
			if ks != nil { _ = ks.Save(pid, kp) }
		}

		cli, err := lccclient.NewClientWithKeyPair(cfg, kp)
		if err != nil {
			errs[pid] = fmt.Sprintf("new client: %v", err)
			continue
		}
		if err := cli.Register(); err != nil {
			errs[pid] = fmt.Sprintf("register: %v", err)
			continue
		}

		s.mu.Lock()
		s.clients[pid] = cli
		s.mu.Unlock()

		registered = append(registered, pid)
		instanceIDs[pid] = cli.GetInstanceID()
	}

	_ = json.NewEncoder(w).Encode(&selectProductsResp{
		OK:          len(registered) > 0 && len(errs) == 0,
		Registered:  registered,
		Errors:      ifNilMap(errs),
		InstanceIDs: ifNilMapStr(instanceIDs),
	})
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl := template.Must(template.New("index").Parse(indexHTML))
	_ = tpl.Execute(w, nil)
}

// handleSimRegistered lists registered product IDs
func (s *Server) handleSimRegistered(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	s.mu.RLock()
	ids := make([]string, 0, len(s.clients))
	for k := range s.clients {
		ids = append(ids, k)
	}
	s.mu.RUnlock()
	_ = json.NewEncoder(w).Encode(ids)
}

// handleSimRoot dispatches to per-action handlers: /api/sim/{product}/{action}
func (s *Server) handleSimRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path // e.g., /api/sim/<product>/<action>
	// strip prefix
	const prefix = "/api/sim/"
	if !strings.HasPrefix(path, prefix) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	rest := strings.TrimPrefix(path, prefix)
	parts := strings.Split(rest, "/")
	if len(parts) < 2 {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid path"))
		return
	}
	productID := parts[0]
	action := parts[1]

	cli, err := s.getClient(productID)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	switch action {
	case "consume":
		s.handleConsume(cli, w, r)
	case "tps-check":
		s.handleTPSCheck(cli, w, r)
	case "capacity-check":
		s.handleCapacityCheck(cli, w, r)
	case "concurrency":
		s.handleConcurrency(cli, w, r)
	case "status":
		s.handleStatus(cli, productID, w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) getClient(productID string) (*lccclient.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cli, ok := s.clients[productID]
	if !ok || cli == nil {
		return nil, fmt.Errorf("product not registered: %s", productID)
	}
	return cli, nil
}

// --- Simulation handlers ---

type consumeReq struct { FeatureID string `json:"feature_id"`; Amount int `json:"amount"` }
type consumeResp struct { Allowed bool `json:"allowed"`; Remaining int `json:"remaining"`; Reason string `json:"reason"` }

func (s *Server) handleConsume(cli *lccclient.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { w.WriteHeader(http.StatusMethodNotAllowed); return }
	var req consumeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err)); return }
	if req.FeatureID == "" || req.Amount <= 0 { writeErr(w, http.StatusBadRequest, fmt.Errorf("feature_id and positive amount required")); return }
	allowed, remaining, reason, err := cli.Consume(req.FeatureID, req.Amount, nil)
	if err != nil { writeErr(w, http.StatusBadGateway, err); return }
	_ = json.NewEncoder(w).Encode(&consumeResp{Allowed: allowed, Remaining: remaining, Reason: reason})
}

type tpsReq struct { FeatureID string `json:"feature_id"`; TPS float64 `json:"tps"` }
type tpsResp struct { Allowed bool `json:"allowed"`; Max float64 `json:"max"`; Reason string `json:"reason"` }

func (s *Server) handleTPSCheck(cli *lccclient.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { w.WriteHeader(http.StatusMethodNotAllowed); return }
	var req tpsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err)); return }
	if req.FeatureID == "" { writeErr(w, http.StatusBadRequest, fmt.Errorf("feature_id required")); return }
	allowed, max, reason, err := cli.CheckTPS(req.FeatureID, req.TPS)
	if err != nil { writeErr(w, http.StatusBadGateway, err); return }
	_ = json.NewEncoder(w).Encode(&tpsResp{Allowed: allowed, Max: max, Reason: reason})
}

type capacityReq struct { FeatureID string `json:"feature_id"`; Current int `json:"current"` }
type capacityResp struct { Allowed bool `json:"allowed"`; Max int `json:"max"`; Reason string `json:"reason"` }

func (s *Server) handleCapacityCheck(cli *lccclient.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { w.WriteHeader(http.StatusMethodNotAllowed); return }
	var req capacityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err)); return }
	if req.FeatureID == "" { writeErr(w, http.StatusBadRequest, fmt.Errorf("feature_id required")); return }
	allowed, max, reason, err := cli.CheckCapacity(req.FeatureID, req.Current)
	if err != nil { writeErr(w, http.StatusBadGateway, err); return }
	_ = json.NewEncoder(w).Encode(&capacityResp{Allowed: allowed, Max: max, Reason: reason})
}

type concurrencyReq struct { FeatureID string `json:"feature_id"`; Slots int `json:"slots"`; HoldMS int `json:"hold_ms"`; Mode string `json:"mode"` }
type concurrencyResp struct { Accepted int `json:"accepted"`; Denied int `json:"denied"`; ReasonStats map[string]int `json:"reason_stats,omitempty"` }

func (s *Server) handleConcurrency(cli *lccclient.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { w.WriteHeader(http.StatusMethodNotAllowed); return }
	var req concurrencyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err)); return }
	if req.FeatureID == "" || req.Slots <= 0 { writeErr(w, http.StatusBadRequest, fmt.Errorf("feature_id and positive slots required")); return }
	mode := req.Mode
	if mode == "" { mode = "check-only" }
	reasonStats := map[string]int{}
	accepted, denied := 0, 0
	switch mode {
	case "signal-only":
		// lightweight signal via usage report
		if err := cli.ReportUsage(req.FeatureID, float64(req.Slots)); err != nil { writeErr(w, http.StatusBadGateway, err); return }
		accepted = req.Slots
	case "check-only":
		for i := 0; i < req.Slots; i++ {
			st, err := cli.CheckFeature(req.FeatureID)
			if err != nil { reasonStats["error"]++; denied++; continue }
			if st.Enabled { accepted++ } else { denied++; reasonStats[st.Reason]++ }
		}
	case "local-lock":
		hold := time.Duration(req.HoldMS) * time.Millisecond
		if hold <= 0 { hold = 100 * time.Millisecond }
		var wg sync.WaitGroup
		for i := 0; i < req.Slots; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				release, ok, reason, err := cli.AcquireSlot(req.FeatureID, nil)
				if err != nil { s.mu.Lock(); reasonStats["error"]++; denied++; s.mu.Unlock(); return }
				if !ok { s.mu.Lock(); reasonStats[reason]++; denied++; s.mu.Unlock(); return }
				time.Sleep(hold)
				release()
				s.mu.Lock(); accepted++; s.mu.Unlock()
			}()
		}
		wg.Wait()
	default:
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid mode"))
		return
	}
	_ = json.NewEncoder(w).Encode(&concurrencyResp{Accepted: accepted, Denied: denied, ReasonStats: ifNilStats(reasonStats)})
}

func ifNilStats(m map[string]int) map[string]int { if len(m) == 0 { return nil }; return m }

// --- Product status ---

type featureStatusDTO struct {
	ID             string                     `json:"id"`
	Name           string                     `json:"name,omitempty"`
	Enabled        bool                       `json:"enabled"`
	Reason         string                     `json:"reason,omitempty"`
	Quota          *lccclient.QuotaInfo       `json:"quota,omitempty"`
	MaxCapacity    int                        `json:"max_capacity,omitempty"`
	MaxTPS         float64                    `json:"max_tps,omitempty"`
	MaxConcurrency int                        `json:"max_concurrency,omitempty"`
}

type productStatusResp struct {
	ProductID  string              `json:"product_id"`
	InstanceID string              `json:"instance_id"`
	Features   []featureStatusDTO  `json:"features"`
}

func (s *Server) handleStatus(cli *lccclient.Client, productID string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost { w.WriteHeader(http.StatusMethodNotAllowed); return }
	s.mu.RLock(); lccURL := s.lccURL; s.mu.RUnlock()
	if lccURL == "" { writeErr(w, http.StatusBadRequest, fmt.Errorf("lcc_url not configured")); return }
features, _ := LoadFeaturesForProduct(productID)
	if len(features) == 0 {
		features, _ = LoadFeatureUnion()
	}
	out := make([]featureStatusDTO, 0, len(features))
	for _, f := range features {
		st, err := cli.CheckFeature(f.ID)
		if err != nil {
			out = append(out, featureStatusDTO{ ID: f.ID, Name: f.Name, Enabled: false, Reason: "check_error" })
			continue
		}
		out = append(out, featureStatusDTO{
			ID: f.ID, Name: f.Name, Enabled: st.Enabled, Reason: st.Reason,
			Quota: st.Quota, MaxCapacity: st.MaxCapacity, MaxTPS: st.MaxTPS, MaxConcurrency: st.MaxConcurrency,
		})
	}
	_ = json.NewEncoder(w).Encode(&productStatusResp{ ProductID: productID, InstanceID: cli.GetInstanceID(), Features: out })
}

func (s *Server) handleProductPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	path := r.URL.Path
	const prefix = "/product/"
	if !strings.HasPrefix(path, prefix) { w.WriteHeader(http.StatusNotFound); return }
	pid := strings.TrimPrefix(path, prefix)
	if pid == "" { w.WriteHeader(http.StatusNotFound); return }
	tpl := template.Must(template.New("product").Parse(productHTML))
	_ = tpl.Execute(w, map[string]any{"ProductID": pid})
}

const productHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>LCC Product Status</title>
  <style>
    body { font-family: sans-serif; background: #0f172a; color: #e5e7eb; padding: 20px; }
    table { border-collapse: collapse; width: 100%; }
    th, td { border: 1px solid #1f2937; padding: 6px 8px; }
    .ok { color: #86efac; }
    .no { color: #fca5a5; }
  </style>
</head>
<body>
  <h2>Product: {{.ProductID}}</h2>
  <button onclick="loadStatus()">Refresh</button>
  <div id="meta"></div>
  <table style="margin-top:10px">
    <thead>
      <tr><th>Feature</th><th>Enabled</th><th>Reason</th><th>Quota</th><th>Capacity</th><th>TPS</th><th>Concurrency</th></tr>
    </thead>
    <tbody id="rows"></tbody>
  </table>
<script>
async function loadStatus(){
  const r = await fetch(window.location.origin + '/api/sim/' + encodeURIComponent('{{.ProductID}}') + '/status');
  const j = await r.json();
  if (j && !j.error){
    document.getElementById('meta').textContent = 'InstanceID: '+(j.instance_id||'');
    const tb = document.getElementById('rows'); tb.innerHTML='';
    (j.features||[]).forEach(function(f){
      var tr = document.createElement('tr');
      function td(t){ var d=document.createElement('td'); d.textContent=t; return d; }
      tr.appendChild(td((f.name||f.id)));
      var en = document.createElement('td'); en.textContent = f.enabled ? '✓' : '✗'; en.className = f.enabled ? 'ok' : 'no'; tr.appendChild(en);
      tr.appendChild(td(f.reason||''));
      var q=''; if (f.quota){ q = 'lim='+f.quota.limit+', used='+f.quota.used+', rem='+f.quota.remaining; }
      tr.appendChild(td(q));
      tr.appendChild(td(f.max_capacity>0?String(f.max_capacity):''));
      tr.appendChild(td(f.max_tps>0?String(f.max_tps):''));
      tr.appendChild(td(f.max_concurrency>0?String(f.max_concurrency):''));
      tb.appendChild(tr);
    });
  }
}
loadStatus();
</script>
</body>
</html>`

func writeErr(w http.ResponseWriter, code int, err error) {
	log.Printf("error: %v", err)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": err.Error(),
	})
}

// handleFeatures returns features for a product from local manifests (SDK-side knowledge)
func (s *Server) handleFeatures(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet { w.WriteHeader(http.StatusMethodNotAllowed); return }
	pid := r.URL.Query().Get("product_id")
	if pid == "" { writeErr(w, http.StatusBadRequest, fmt.Errorf("product_id required")); return }
	features, _ := LoadFeaturesForProduct(pid)
	if len(features) == 0 {
		// fallback: union of all known features in local manifests
		features, _ = LoadFeatureUnion()
	}
	_ = json.NewEncoder(w).Encode(features)
}

func nonEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func ifNilMap(m map[string]string) map[string]string {
	if m == nil || len(m) == 0 {
		return nil
	}
	return m
}

func ifNilMapStr(m map[string]string) map[string]string {
	if m == nil || len(m) == 0 {
		return nil
	}
	return m
}

// --- Config persistence ---

type persistedConfig struct {
	LCCURL string `json:"lcc_url"`
}

func (s *Server) configPath() (string, error) {
	h, err := os.UserHomeDir()
	if err != nil { return "", err }
	root := filepath.Join(h, ".lcc-demo")
	if err := os.MkdirAll(root, 0700); err != nil { return "", err }
	return filepath.Join(root, "config.json"), nil
}

func (s *Server) loadConfig() {
	p, err := s.configPath()
	if err != nil { return }
	data, err := os.ReadFile(p)
	if err != nil { return }
	var cfg persistedConfig
	if err := json.Unmarshal(data, &cfg); err != nil { return }
	s.mu.Lock()
	s.lccURL = cfg.LCCURL
	s.publicBase = "/api/v1/public"
	s.mu.Unlock()
}

func (s *Server) saveConfig() error {
	p, err := s.configPath()
	if err != nil { return err }
	s.mu.RLock()
	cfg := persistedConfig{LCCURL: s.lccURL}
	s.mu.RUnlock()
	data, err := json.MarshalIndent(&cfg, "", "  ")
	if err != nil { return err }
	return os.WriteFile(p, data, 0600)
}

const indexHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>LCC SDK Demo App</title>
  <style>
    :root {
      --bg: #0f172a;
      --bg-soft: #0b1224;
      --panel: #020617;
      --border: #1f2937;
      --muted: #9ca3af;
      --text: #e5e7eb;
      --accent: #60a5fa;
      --accent-hover: #3b82f6;
      --shadow: 0 8px 20px rgba(0,0,0,0.25);
      --radius: 12px;
      --header-h: 70px;
    }
    * { box-sizing: border-box; }
    html, body { height: 100%; }
    body { font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Ubuntu, Cantarell, Noto Sans, Arial, "Apple Color Emoji", "Segoe UI Emoji"; background: radial-gradient(1200px 800px at 20% 0%, #0b132e 0%, var(--bg) 40%), var(--bg); color: var(--text); margin: 0; }
    a { color: inherit; }

    .header { position: sticky; top: 0; z-index: 10; height: var(--header-h); display: flex; align-items: center; gap: 12px; padding: 0 30px; background: linear-gradient(180deg, #0b1220 0%, #0f172a 100%); border-bottom: 1px solid var(--border); box-shadow: 0 4px 20px rgba(0,0,0,0.3); }
    .brand { font-size: 24px; font-weight: 700; background: linear-gradient(to right, #60a5fa, #a78bfa); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
    .badge { font-size: 12px; color: #a78bfa; background: rgba(167,139,250,0.15); padding: 4px 12px; border: 1px solid rgba(167,139,250,0.3); border-radius: 999px; font-weight: 600; }

    .container { display: flex; min-height: calc(100vh - var(--header-h)); }
    .nav { width: 260px; background: var(--panel); border-right: 1px solid var(--border); padding: 24px 16px; position: sticky; top: var(--header-h); height: calc(100vh - var(--header-h)); overflow: auto; }
    .nav a { display: block; color: var(--text); text-decoration: none; padding: 12px 16px; border-radius: 8px; margin-bottom: 8px; transition: all .2s ease; font-weight: 500; }
    .nav a:hover { background: #1e293b; transform: translateX(4px); }
    .nav a.active { background: rgba(96,165,250,0.15); border-left: 3px solid var(--accent); padding-left: 13px; color: var(--accent); }

    .main { flex: 1; padding: 24px; }
    section { scroll-margin-top: calc(var(--header-h) + 12px); }

    .card { background: var(--panel); border-radius: var(--radius); padding: 18px; border: 1px solid var(--border); box-shadow: 0 1px 0 rgba(255,255,255,0.02) inset; }
    .card + .card { margin-top: 14px; }

    .tabs { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 8px; }
    .tab { background: transparent; border: 1px solid var(--border); color: var(--text); padding: 6px 10px; border-radius: 999px; cursor: pointer; }
    .tab.active { background: rgba(37,99,235,0.15); border-color: var(--accent); }

    h1, h2, h3, h4 { margin: 0 0 10px 0; }
    h3 { color: #cbd5e1; font-weight: 600; }

    input, select, button { padding: 8px 10px; border-radius: 8px; border: 1px solid var(--border); background: #0b1224; color: var(--text); }
    input::placeholder { color: #6b7280; }
    button { background: linear-gradient(135deg, #3b82f6, #8b5cf6); border: none; font-weight: 600; cursor: pointer; transition: all .2s ease; box-shadow: 0 4px 12px rgba(59,130,246,0.3); }
    button:hover { filter: brightness(1.1); box-shadow: 0 8px 20px rgba(59,130,246,0.4); transform: translateY(-2px); }
    button:active { transform: translateY(0); box-shadow: 0 2px 8px rgba(59,130,246,0.3); }

    .muted { color: var(--muted); }
    code { color: #93c5fd; }
    .error { color: #fca5a5; }
    .hidden { display: none; }
  </style>
</head>
<body>
  <header class="header">
    <div class="brand">LCC SDK Demo App</div>
    <div class="badge">Demo</div>
  </header>
  <div class="container">
    <aside class="nav">
      <a href="#configure">Configure LCC</a>
      <a href="#discover">Discover Products</a>
      <a href="#simulate">Simulate</a>
    </aside>
    <main class="main">
      <section id="configure" class="card">
        <h3>Configure LCC</h3>
        <div style="display:flex; gap:8px; flex-wrap:wrap">
          <input type="text" id="lcc_url" placeholder="http://localhost:7086" style="min-width:360px" />
          <button onclick="setConfig()">Save</button>
        </div>
        <span id="cfg_status" class="muted" style="margin-top:8px; display:block"></span>
      </section>

      <section id="discover" class="card hidden">
        <h3>Discover Products - License Details</h3>
        <button onclick="loadProducts()">Load Products</button>
        <div id="products_tabs" class="tabs" style="margin-top:12px"></div>
        <div id="product_details" style="margin-top:12px"></div>
      </section>

      <section id="simulate" class="card hidden">
        <h3>Simulate</h3>
        <div>
          <div class="tabs" id="sim_tabs"></div>
          <div style="display:flex; gap:6px; align-items:center; flex-wrap:wrap; margin-bottom:8px">
            <button onclick="refreshRegisteredTabs()">Refresh</button>
            <select id="sim_add_pid" title="Product ID"></select>
            <input type="text" id="sim_add_ver" value="1.0.0" title="Version" style="width:110px" />
            <button onclick="simLoadAvailableProducts()">Load Products</button>
            <button onclick="simRegister()">Register</button>
          </div>
          <div style="display:none">
            <select id="sim_product"></select>
          </div>
          <div style="display:flex; gap:6px; align-items:center; flex-wrap:wrap; margin-top:4px">
            <button onclick="loadFeatures()">Load Features</button>
            <button onclick="openStatus()">Open Status</button>
            <datalist id="features_list"></datalist>
          </div>
        </div>
        <div style="margin-top:12px">
          <h4 class="muted">CALL (Consume)</h4>
          <input type="text" id="c_feature" placeholder="feature_id" list="features_list" />
          <input type="number" id="c_amount" placeholder="amount" value="1" />
          <button onclick="doConsume()">Send</button>
        </div>
        <div style="margin-top:12px">
          <h4 class="muted">TPS Check</h4>
          <input type="text" id="t_feature" placeholder="feature_id" list="features_list" />
          <input type="number" id="t_value" placeholder="tps" value="1" step="0.1" />
          <button onclick="doTPS()">Check</button>
        </div>
        <div style="margin-top:12px">
          <h4 class="muted">CAPACITY Check</h4>
          <input type="text" id="k_feature" placeholder="feature_id" list="features_list" />
          <input type="number" id="k_current" placeholder="current" value="1" />
          <button onclick="doCapacity()">Check</button>
        </div>
        <div style="margin-top:12px">
          <h4 class="muted">USERS / Concurrency</h4>
          <input type="text" id="u_feature" placeholder="feature_id" list="features_list" />
          <input type="number" id="u_slots" placeholder="slots" value="5" />
          <input type="number" id="u_hold" placeholder="hold_ms" value="100" />
          <select id="u_mode">
            <option value="check-only">check-only</option>
            <option value="signal-only">signal-only</option>
            <option value="local-lock">local-lock</option>
          </select>
          <button onclick="doConcurrency()">Run</button>
        </div>
        <pre id="sim_log" style="white-space:pre-wrap;background:#0b1224;color:#d1d5db;padding:10px;margin-top:12px;max-height:240px;overflow:auto;border:1px solid var(--border); border-radius: 8px;"></pre>
      </section>
    </main>
  </div>

<script>
// Simple in-page navigation: show a single section at a time
(function(){
  var sections = ['configure','discover','simulate'];
  var links = Array.prototype.slice.call(document.querySelectorAll('.nav a'));
  function setActive(id){
    sections.forEach(function(s){
      var el = document.getElementById(s);
      if (!el) return;
      if (s === id) { el.classList.remove('hidden'); } else { el.classList.add('hidden'); }
    });
    links.forEach(function(a){
      var t = (a.getAttribute('href')||'').replace('#','');
      if (t === id) { a.classList.add('active'); } else { a.classList.remove('active'); }
    });
  }
  function navigateTo(id, replace){
    if (sections.indexOf(id) === -1) id = sections[0];
    setActive(id);
    if (replace) { history.replaceState(null,'','#'+id); } else { history.pushState(null,'','#'+id); }
  }
  window.addEventListener('hashchange', function(){ navigateTo((location.hash||'').replace('#',''), true); });
  links.forEach(function(a){ a.addEventListener('click', function(e){ e.preventDefault(); navigateTo((this.getAttribute('href')||'').replace('#','')); }); });
  var initial = (location.hash||'').replace('#','') || sections[0];
  navigateTo(initial, true);
  // Initialize simulate area
  refreshRegisteredTabs();
})();

async function setConfig() {
  const url = document.getElementById('lcc_url').value.trim();
  const r = await fetch('/api/config', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({lcc_url:url})});
  const j = await r.json();
  document.getElementById('cfg_status').textContent = j.error ? ('Error: '+j.error) : ('OK: '+j.lcc_url);
}

async function loadCurrentConfig() {
  const r = await fetch('/api/config');
  if (r.ok) {
    const j = await r.json();
    if (j && j.lcc_url) document.getElementById('lcc_url').value = j.lcc_url;
  }
}
window.addEventListener('DOMContentLoaded', loadCurrentConfig);

let ALL_PRODUCTS = [];
let SELECTED_PRODUCT_IDX = 0;

async function loadProducts() {
  const tabsEl = document.getElementById('products_tabs');
  const detailsEl = document.getElementById('product_details');
  tabsEl.innerHTML = '<span class=\"muted\">Loading...</span>';
  detailsEl.innerHTML = '';
  
  const r = await fetch('/api/products');
  const j = await r.json();
  if (j.error) { 
    tabsEl.innerHTML = '<span class=\"error\">'+j.error+'</span>'; 
    return; 
  }
  if (!Array.isArray(j) || j.length === 0) { 
    tabsEl.innerHTML = '<span class=\"muted\">No products found</span>'; 
    return; 
  }
  
  ALL_PRODUCTS = j;
  renderProductTabs();
  await showProductDetails(0);
}

function renderProductTabs() {
  const tabsEl = document.getElementById('products_tabs');
  tabsEl.innerHTML = '';
  ALL_PRODUCTS.forEach(function(p, idx) {
    var btn = document.createElement('button');
    btn.className = 'tab' + (idx === SELECTED_PRODUCT_IDX ? ' active' : '');
    btn.textContent = p.name || p.id;
    btn.onclick = function() { showProductDetails(idx); };
    tabsEl.appendChild(btn);
  });
}

async function showProductDetails(idx) {
  SELECTED_PRODUCT_IDX = idx;
  renderProductTabs();
  
  const product = ALL_PRODUCTS[idx];
  const detailsEl = document.getElementById('product_details');
  detailsEl.innerHTML = '<div class=\"muted\">Loading license details...</div>';
  
  // Register product to get SDK client
  const regResp = await fetch('/api/sim/products', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({product_ids: [product.id], default_version: '1.0.0'})
  });
  const regData = await regResp.json();
  
  if (regData.error || !regData.ok) {
    detailsEl.innerHTML = '<div class=\"error\">Failed to load license: ' + (regData.error || 'Unknown error') + '</div>';
    return;
  }
  
  // Get license details via status endpoint
  const statusResp = await fetch('/api/sim/' + encodeURIComponent(product.id) + '/status');
  const statusData = await statusResp.json();
  
  if (statusData.error) {
    detailsEl.innerHTML = '<div class=\"error\">Failed to load license details</div>';
    return;
  }
  
  // Render license details
  var html = '<div style=\"background:#0b1224;padding:16px;border-radius:8px;border:1px solid var(--border)\"><h4>Product: ' + product.id + '</h4>';
  html += '<div class=\"muted\" style=\"margin:8px 0\">Instance ID: ' + (statusData.instance_id || '') + '</div>';
  
  if (statusData.features && statusData.features.length > 0) {
    html += '<h4 style=\"margin-top:16px\">Features & Limitations:</h4>';
    html += '<table style=\"width:100%;border-collapse:collapse;margin-top:8px\">';
    html += '<thead><tr><th style=\"border:1px solid var(--border);padding:6px;text-align:left\">Feature</th><th style=\"border:1px solid var(--border);padding:6px\">Status</th><th style=\"border:1px solid var(--border);padding:6px\">Details</th></tr></thead><tbody>';
    
    statusData.features.forEach(function(f) {
      var statusText = f.enabled ? '<span style=\"color:#86efac\">✓ Enabled</span>' : '<span style=\"color:#fca5a5\">✗ Disabled</span>';
      var details = [];
      if (f.reason) details.push('Reason: ' + f.reason);
      if (f.quota) details.push('Quota: ' + f.quota.limit + ' (used: ' + f.quota.used + ', remaining: ' + f.quota.remaining + ')');
      if (f.max_capacity > 0) details.push('Max Capacity: ' + f.max_capacity);
      if (f.max_tps > 0) details.push('Max TPS: ' + f.max_tps);
      if (f.max_concurrency > 0) details.push('Max Concurrency: ' + f.max_concurrency);
      
      html += '<tr>';
      html += '<td style=\"border:1px solid var(--border);padding:6px\">' + (f.name || f.id) + '</td>';
      html += '<td style=\"border:1px solid var(--border);padding:6px;text-align:center\">' + statusText + '</td>';
      html += '<td style=\"border:1px solid var(--border);padding:6px;font-size:0.9em\">' + (details.length > 0 ? details.join('<br>') : '-') + '</td>';
      html += '</tr>';
    });
    
    html += '</tbody></table>';
  }
  
  html += '<div style=\"margin-top:16px;text-align:center\">';
  html += '<button onclick=\"selectProductForSimulation(\\'' + product.id + '\\')\" style=\"padding:10px 24px\">Select & Configure Simulation →</button>';
  html += '</div>';
  html += '</div>';
  
  detailsEl.innerHTML = html;
}

// --- Simulate tabs & registration ---
let CURRENT_PRODUCT = '';
function renderTabs(ids){
  const c = document.getElementById('sim_tabs');
  c.innerHTML = '';
  (ids||[]).forEach(function(id){
    var b = document.createElement('button');
    b.className = 'tab' + (id===CURRENT_PRODUCT?' active':'');
    b.textContent = id;
    b.onclick = function(){ simSelectProduct(id); };
    c.appendChild(b);
  });
}
function simSelectProduct(id){
  CURRENT_PRODUCT = id || '';
  renderTabs(Array.from(document.querySelectorAll('#sim_tabs .tab')).map(function(el){return el.textContent;}));
}
async function refreshRegisteredTabs(){
  const r = await fetch('/api/sim/registered');
  const j = await r.json();
  if (Array.isArray(j)){
    if (!CURRENT_PRODUCT && j.length>0) CURRENT_PRODUCT = j[0];
    renderTabs(j);
  }
}
async function simLoadAvailableProducts(){
  const sel = document.getElementById('sim_add_pid'); sel.innerHTML='';
  const r = await fetch('/api/products');
  const j = await r.json();
  if (Array.isArray(j)){
    j.forEach(function(p){ var o=document.createElement('option'); o.value=p.id; o.text=p.id + (p.name?(' - '+p.name):''); sel.appendChild(o); });
  }
}
async function simRegister(){
  const pidSel = document.getElementById('sim_add_pid');
  const pid = pidSel && pidSel.value ? pidSel.value : '';
  const ver = (document.getElementById('sim_add_ver').value||'').trim() || '1.0.0';
  if (!pid){ appendLog('No product selected for register'); return; }
  const r = await fetch('/api/sim/products', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({product_ids:[pid], default_version: ver})});
  const j = await r.json();
  if (j && !j.error){ CURRENT_PRODUCT = pid; await refreshRegisteredTabs(); appendLog('Registered: '+(j.registered||[]).join(', ')); }
}

async function registerSelected() {
  const boxes = Array.from(document.querySelectorAll('input[name=\"pid\"]:checked'));
  const ids = boxes.map(b => b.value);
  const defVer = document.getElementById('def_ver').value.trim();
  if (ids.length === 0) { document.getElementById('reg_status').textContent = 'No products selected'; return; }
  const r = await fetch('/api/sim/products', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({product_ids: ids, default_version: defVer})});
  const j = await r.json();
  if (j.error) { document.getElementById('reg_status').textContent = 'Error: '+j.error; return; }
  document.getElementById('reg_status').textContent = 'Registered: '+(j.registered||[]).join(', ');
}

async function loadRegistered() { await refreshRegisteredTabs(); }

function pickProduct(){ if (CURRENT_PRODUCT) return CURRENT_PRODUCT; var sel=document.getElementById('sim_product'); return sel && sel.value ? sel.value : ''; }
function appendLog(msg){ var el=document.getElementById('sim_log'); el.textContent = (el.textContent?el.textContent+'\n':'') + msg; }

async function loadFeatures(){
  const pid = pickProduct(); if(!pid){ appendLog('No product selected'); return; }
  const r = await fetch('/api/features?product_id='+encodeURIComponent(pid));
  const j = await r.json();
  if (j && !j.error && Array.isArray(j)) {
    var dl = document.getElementById('features_list');
    dl.innerHTML = '';
    j.forEach(function(f){ var opt = document.createElement('option'); opt.value = f.id; opt.label = f.name || f.id; dl.appendChild(opt); });
  } else {
    appendLog('load features error: '+(j && j.error ? j.error : 'unexpected'));
  }
}

function openStatus(){ var pid = pickProduct(); if(!pid){ appendLog('No product selected'); return; } window.open('/product/'+encodeURIComponent(pid), '_blank'); }

async function doConsume(){
  const pid = pickProduct(); if(!pid){ appendLog('No product selected'); return; }
  const f = document.getElementById('c_feature').value.trim();
  const a = parseInt(document.getElementById('c_amount').value,10)||0;
  const r = await fetch('/api/sim/'+encodeURIComponent(pid)+'/consume', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({feature_id:f, amount:a})});
  appendLog('consume => '+await r.text());
}

async function doTPS(){
  const pid = pickProduct(); if(!pid){ appendLog('No product selected'); return; }
  const f = document.getElementById('t_feature').value.trim();
  const v = parseFloat(document.getElementById('t_value').value)||0;
  const r = await fetch('/api/sim/'+encodeURIComponent(pid)+'/tps-check', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({feature_id:f, tps:v})});
  appendLog('tps-check => '+await r.text());
}

async function doCapacity(){
  const pid = pickProduct(); if(!pid){ appendLog('No product selected'); return; }
  const f = document.getElementById('k_feature').value.trim();
  const c = parseInt(document.getElementById('k_current').value,10)||0;
  const r = await fetch('/api/sim/'+encodeURIComponent(pid)+'/capacity-check', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({feature_id:f, current:c})});
  appendLog('capacity-check => '+await r.text());
}

function selectProductForSimulation(productId) {
  CURRENT_PRODUCT = productId;
  refreshRegisteredTabs();
  window.location.hash = 'simulate';
}

async function doConcurrency(){
  const pid = pickProduct(); if(!pid){ appendLog('No product selected'); return; }
  const f = document.getElementById('u_feature').value.trim();
  const s = parseInt(document.getElementById('u_slots').value,10)||0;
  const h = parseInt(document.getElementById('u_hold').value,10)||0;
  const m = document.getElementById('u_mode').value;
  const r = await fetch('/api/sim/'+encodeURIComponent(pid)+'/concurrency', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({feature_id:f, slots:s, hold_ms:h, mode:m})});
  appendLog('concurrency => '+await r.text());
}
</script>
</body>
</html>`
