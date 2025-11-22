package web

type LimitType struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SDKAPI      string `json:"sdk_api"`
	Tracking    string `json:"tracking"`
	UseCases    []string `json:"use_cases"`
	TimeDimension string `json:"time_dimension"`
	WhoTracks   string `json:"who_tracks"`
}

var AllLimitTypes = []LimitType{
	{
		Type:        "quota",
		Name:        "Quota Control",
		Title:       "Quota Control",
		Description: "Product-level cumulative consumption limit. All features share the same quota pool. Auto-resets on schedule.",
		SDKAPI:      "Consume(amount) - No featureID needed (product-level)",
		Tracking:    "Server-side automatic",
		UseCases: []string{
			"Total API calls across all features (50,000/month)",
			"Combined export operations (all formats)",
			"Shared license generation credits",
			"Any metered resource pool",
		},
		TimeDimension: "Daily or Monthly window with auto-reset",
		WhoTracks:     "Server-side (LCC) - compiler auto-injects Consume() calls",
	},
	{
		Type:        "tps",
		Name:        "TPS (Rate Limit)",
		Title:       "TPS (Rate Limit)",
		Description: "Product-level instantaneous throughput limit. All features share the same TPS budget.",
		SDKAPI:      "CheckTPS() - SDK auto-tracks or uses TPSProvider helper",
		Tracking:    "SDK automatic or custom TPSProvider function",
		UseCases: []string{
			"Product-wide API rate (100 req/sec total)",
			"Combined stream processing throughput",
			"Aggregate event ingestion rate",
			"Burst control across all operations",
		},
		TimeDimension: "Per-second instantaneous measurement",
		WhoTracks:     "SDK auto-tracks or developer provides TPSProvider function",
	},
	{
		Type:        "capacity",
		Name:        "Capacity Limit",
		Title:       "Capacity Limit",
		Description: "Product-level maximum quantity of persistent resources. Requires developer-provided CapacityCounter function.",
		SDKAPI:      "CheckCapacity() - Uses CapacityCounter helper (Required)",
		Tracking:    "Developer provides CapacityCounter function",
		UseCases: []string{
			"Total projects across product (100 projects max)",
			"Combined storage items (all resource types)",
			"Total user accounts per tenant",
			"Aggregate persistent resources",
		},
		TimeDimension: "Persistent - no time-based reset",
		WhoTracks:     "Developer must provide CapacityCounter function in config",
	},
	{
		Type:        "concurrency",
		Name:        "Concurrency Limit",
		Title:       "Concurrency Limit",
		Description: "Product-level simultaneous execution slots. All features share the concurrency pool. SDK manages automatically.",
		SDKAPI:      "AcquireSlot() → returns release() - Auto-managed by compiler",
		Tracking:    "SDK internal counter (automatic)",
		UseCases: []string{
			"Total concurrent sessions (10 max across product)",
			"Combined parallel operations",
			"Shared execution slot pool",
			"Product-wide connection pool",
		},
		TimeDimension: "Duration of operation (held then released)",
		WhoTracks:     "SDK tracks automatically - compiler injects acquire/release",
	},
}

func GetLimitTypeByType(limitType string) *LimitType {
	for _, lt := range AllLimitTypes {
		if lt.Type == limitType {
			return &lt
		}
	}
	return nil
}

type LimitExample struct {
	LicenseConfig string   `json:"license_config"`
	CodeExample   string   `json:"code_example"`
	BehaviorTable []BehaviorRow `json:"behavior_table"`
	KeyPoints     []string `json:"key_points"`
}

type BehaviorRow struct {
	Call      string `json:"call"`
	Allowed   string `json:"allowed"`
	Remaining string `json:"remaining"`
	Reason    string `json:"reason"`
}

func GetLimitExample(limitType string) *LimitExample {
	switch limitType {
	case "quota":
		return &LimitExample{
			LicenseConfig: `{
  "product_id": "data-insight-pro",
  "tier": "professional",
  "features": {
    "ml_analytics": { "enabled": true },
    "pdf_export": { "enabled": true },
    "api_access": { "enabled": true }
  },
  "limits": {
    "quota": {
      "max": 50000,           // Product-level quota
      "used": 0,              // Shared by all features
      "remaining": 50000,
      "window": "monthly",    // Reset period
      "reset_at": "2025-02-01T00:00:00Z"
    }
  }
}`,
			CodeExample: `// ========== Developer Code (Clean Business Logic) ==========
func ProcessAnalytics(data Dataset) error {
    // No license code needed - compiler auto-injects!
    return analytics.RunMLModel(data)
}

// ========== Optional: Custom Quota Calculator ==========
// Define in YAML: limits.quota.consumer = GetConsumeAmount
func GetConsumeAmount(ctx context.Context, args ...interface{}) int {
    data := args[0].(Dataset)
    return data.SizeKB()  // Charge by data size
}

// ========== Compiler Auto-Generated Code ==========
func ProcessAnalytics__generated(data Dataset) error {
    // Auto-injected quota check
    amount := GetConsumeAmount(context.Background(), data)
    allowed, remaining, err := __lcc.Consume(amount)  // No featureID!
    
    if err != nil || !allowed {
        log.Warn("Quota exceeded", "remaining", remaining)
        return ErrQuotaExceeded
    }
    
    // Original business logic
    return analytics.RunMLModel(data)
}`,
			BehaviorTable: []BehaviorRow{
				{Call: "1", Allowed: "✓ Yes", Remaining: "49,999", Reason: "ok"},
				{Call: "1,000", Allowed: "✓ Yes", Remaining: "49,000", Reason: "ok"},
				{Call: "25,000", Allowed: "✓ Yes", Remaining: "25,000", Reason: "ok"},
				{Call: "49,999", Allowed: "✓ Yes", Remaining: "1", Reason: "ok"},
				{Call: "50,000", Allowed: "✓ Yes", Remaining: "0", Reason: "ok"},
				{Call: "50,001", Allowed: "❌ No", Remaining: "0", Reason: "exceeded"},
				{Call: "(Next Month)", Allowed: "✓ Yes", Remaining: "49,999", Reason: "reset"},
			},
			KeyPoints: []string{
				"✅ Product-level limit (shared across all features)",
				"✅ Zero-intrusion: compiler auto-injects Consume() calls",
				"🔧 Optional helper: QuotaConsumer for custom amount calculation",
				"📊 Server tracks cumulative total automatically",
				"🔄 Auto-resets daily/monthly per license config",
			},
		}
	case "tps":
		return &LimitExample{
			LicenseConfig: `{
  "product_id": "data-insight-pro",
  "tier": "professional",
  "features": {
    "ml_analytics": { "enabled": true },
    "pdf_export": { "enabled": true },
    "api_access": { "enabled": true }
  },
  "limits": {
    "max_tps": 100.0   // Product-level TPS (shared by all features)
  }
}`,
			CodeExample: `// ========== Developer Code (Clean Business Logic) ==========
func HandleAPIRequest(ctx context.Context) error {
    // No license code needed - compiler auto-injects!
    return processRequest(ctx)
}

// ========== Optional: Custom TPS Provider ==========
// Define in YAML: limits.tps_provider = GetCurrentTPS
func GetCurrentTPS() float64 {
    return myRateLimiter.GetCurrentRate()  // Custom rate measurement
}

// ========== Compiler Auto-Generated Code ==========
func HandleAPIRequest__generated(ctx context.Context) error {
    // Auto-injected TPS check
    currentTPS := GetCurrentTPS()  // Or SDK auto-tracks if not provided
    allowed, maxTPS, err := __lcc.CheckTPS(currentTPS)  // No featureID!
    
    if err != nil || !allowed {
        log.Warn("TPS exceeded", "current", currentTPS, "max", maxTPS)
        return ErrRateLimitExceeded
    }
    
    // Original business logic
    return processRequest(ctx)
}`,
			BehaviorTable: []BehaviorRow{
				{Call: "TPS=50.5", Allowed: "✓ Yes", Remaining: "max=100.0", Reason: "ok"},
				{Call: "TPS=95.3", Allowed: "✓ Yes", Remaining: "max=100.0", Reason: "ok"},
				{Call: "TPS=100.0", Allowed: "✓ Yes", Remaining: "max=100.0", Reason: "ok"},
				{Call: "TPS=100.5", Allowed: "❌ No", Remaining: "max=100.0", Reason: "exceeded"},
				{Call: "TPS=150.0", Allowed: "❌ No", Remaining: "max=100.0", Reason: "exceeded"},
				{Call: "(Next Sec)", Allowed: "✓ Yes", Remaining: "max=100.0", Reason: "ok (rate dropped)"},
			},
			KeyPoints: []string{
				"✅ Product-level limit (shared TPS budget for all features)",
				"✅ Zero-intrusion: compiler auto-injects CheckTPS() calls",
				"🔧 Optional helper: TPSProvider for custom rate measurement",
				"📊 SDK can auto-track TPS if helper not provided",
				"⚡ Instantaneous rate check (no cumulative state)",
			},
		}
	case "capacity":
		return &LimitExample{
			LicenseConfig: `{
  "product_id": "data-insight-enterprise",
  "tier": "enterprise",
  "features": {
    "custom_dashboard": { "enabled": true },
    "projects": { "enabled": true }
  },
  "limits": {
    "max_capacity": 100   // Product-level capacity (shared resources)
  }
}`,
			CodeExample: `// ========== Developer Code (Clean Business Logic) ==========
func CreateProject(name string) error {
    // No license code needed - compiler auto-injects!
    return db.CreateProject(name)
}

// ========== Required: Capacity Counter Function ==========
// MUST define in YAML: limits.capacity_counter = GetCurrentProjectCount
func GetCurrentProjectCount() int {
    count, _ := db.Query("SELECT COUNT(*) FROM projects")
    return count  // Return current resource usage
}

// ========== Compiler Auto-Generated Code ==========
func CreateProject__generated(name string) error {
    // Auto-injected capacity check
    currentCount := GetCurrentProjectCount()  // Call required helper
    allowed, maxCap, err := __lcc.CheckCapacity(currentCount)  // No featureID!
    
    if err != nil || !allowed {
        log.Warn("Capacity exceeded", "current", currentCount, "max", maxCap)
        return ErrCapacityExceeded
    }
    
    // Original business logic
    return db.CreateProject(name)
}`,
			BehaviorTable: []BehaviorRow{
				{Call: "count=10", Allowed: "✓ Yes", Remaining: "max=100", Reason: "ok"},
				{Call: "count=50", Allowed: "✓ Yes", Remaining: "max=100", Reason: "ok"},
				{Call: "count=99", Allowed: "✓ Yes", Remaining: "max=100", Reason: "ok"},
				{Call: "count=100", Allowed: "❌ No", Remaining: "max=100", Reason: "at_limit"},
				{Call: "count=101", Allowed: "❌ No", Remaining: "max=100", Reason: "exceeded"},
				{Call: "(After Delete)", Allowed: "✓ Yes", Remaining: "max=100", Reason: "ok (space freed)"},
			},
			KeyPoints: []string{
				"✅ Product-level limit (total resources across all features)",
				"✅ Zero-intrusion: compiler auto-injects CheckCapacity() calls",
				"⚠️ REQUIRED helper: CapacityCounter function MUST be provided",
				"📊 Developer provides counter to query current resource usage",
				"♻️ Persistent limit - no time-based reset",
			},
		}
	case "concurrency":
		return &LimitExample{
			LicenseConfig: `{
  "product_id": "data-insight-pro",
  "tier": "professional",
  "features": {
    "api_access": { "enabled": true },
    "concurrent_sessions": { "enabled": true }
  },
  "limits": {
    "max_concurrency": 10   // Product-level concurrency (all features share)
  }
}`,
			CodeExample: `// ========== Developer Code (Clean Business Logic) ==========
func HandleUserSession(userID string) error {
    // No license code needed - compiler auto-injects!
    return handleUserActions(userID)
}

// ========== Compiler Auto-Generated Code ==========
func HandleUserSession__generated(userID string) error {
    // Auto-injected concurrency control
    release, allowed, err := __lcc.AcquireSlot()  // No featureID!
    if err != nil || !allowed {
        log.Warn("Concurrency limit reached")
        return ErrTooManyConcurrentUsers
    }
    
    // Auto-injected defer to release slot
    defer release()
    
    // Original business logic
    log.Info("Session started", "user", userID)
    return handleUserActions(userID)
}`,
			BehaviorTable: []BehaviorRow{
				{Call: "Slot 1", Allowed: "✓ Yes", Remaining: "9 free", Reason: "ok"},
				{Call: "Slot 5", Allowed: "✓ Yes", Remaining: "5 free", Reason: "ok"},
				{Call: "Slot 10", Allowed: "✓ Yes", Remaining: "0 free", Reason: "ok"},
				{Call: "Slot 11", Allowed: "❌ No", Remaining: "0 free", Reason: "max_reached"},
				{Call: "(Slot Released)", Allowed: "✓ Yes", Remaining: "1 free", Reason: "ok (slot freed)"},
				{Call: "Slot 10 again", Allowed: "✓ Yes", Remaining: "0 free", Reason: "ok"},
			},
			KeyPoints: []string{
				"✅ Product-level limit (shared slot pool across all features)",
				"✅ Zero-intrusion: compiler auto-injects acquire/release",
				"🔧 No helper needed - SDK tracks slots automatically",
				"♻️ Compiler ensures defer release() for safety",
				"⚡ Real-time slot management (instant acquire/release)",
			},
		}
	default:
		return nil
	}
}
