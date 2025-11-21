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
		Title:       "Quota (配额控制)",
		Description: "Cumulative consumption limit that resets on schedule. Server tracks total usage automatically.",
		SDKAPI:      "Consume(featureID, amount)",
		Tracking:    "Server-side automatic",
		UseCases: []string{
			"API call counting (10,000/day)",
			"Export operations (200 PDFs/month)",
			"License generation credits",
			"Any metered/consumable resource",
		},
		TimeDimension: "Daily or Monthly window with auto-reset",
		WhoTracks:     "Server-side (LCC) - developer just calls Consume()",
	},
	{
		Type:        "tps",
		Name:        "TPS (Rate Limit)",
		Title:       "TPS (速率限制)",
		Description: "Instantaneous throughput limit. Controls requests per second.",
		SDKAPI:      "CheckTPS(featureID, currentTPS)",
		Tracking:    "Client calculates current rate",
		UseCases: []string{
			"API rate limiting (10 req/sec)",
			"Stream processing throughput",
			"Real-time event ingestion",
			"Burst control for expensive operations",
		},
		TimeDimension: "Per-second instantaneous measurement",
		WhoTracks:     "Client measures TPS and passes to CheckTPS()",
	},
	{
		Type:        "capacity",
		Name:        "Capacity Limit",
		Title:       "Capacity (容量限制)",
		Description: "Maximum quantity of persistent resources. Controls how many items can exist.",
		SDKAPI:      "CheckCapacity(featureID, currentUsed)",
		Tracking:    "Client counts current usage",
		UseCases: []string{
			"Maximum projects (50 projects)",
			"Storage items (1000 documents)",
			"User accounts per tenant",
			"Persistent database records",
		},
		TimeDimension: "Persistent - no time-based reset",
		WhoTracks:     "Client counts resources and passes to CheckCapacity()",
	},
	{
		Type:        "concurrency",
		Name:        "Concurrency Limit",
		Title:       "Concurrency (并发限制)",
		Description: "Simultaneous execution slots. Controls how many operations can run at the same time.",
		SDKAPI:      "AcquireSlot(featureID) → returns release()",
		Tracking:    "SDK internal counter",
		UseCases: []string{
			"Concurrent user sessions (10 users)",
			"Parallel batch jobs",
			"Export generation slots",
			"License server connection pool",
		},
		TimeDimension: "Duration of operation (held then released)",
		WhoTracks:     "SDK tracks automatically with AcquireSlot/release pattern",
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
  "ml_analytics": {
    "enabled": true,
    "quota": {
      "max": 10000,         // Total allowed
      "window": "daily",    // Reset period
      "reset_at": "00:00"   // Reset time (UTC)
    }
  }
}`,
			CodeExample: `func ProcessAnalytics(data Dataset) error {
    // SDK reports usage to LCC automatically
    allowed, remaining, reason, err :=
      lccClient.Consume(
        "ml_analytics",  // Feature ID
        1,               // Credits to consume
        nil,             // Optional metadata
      )

    if err != nil {
      return fmt.Errorf("license: %w", err)
    }

    if !allowed {
      log.Warn("Quota exceeded",
        "remaining", remaining,
        "reason", reason)
      return ErrQuotaExceeded
    }

    // Quota OK - execute expensive operation
    result := analytics.RunMLModel(data)
    log.Info("Success", "remaining", remaining)
    return nil
}`,
			BehaviorTable: []BehaviorRow{
				{Call: "1", Allowed: "✓ Yes", Remaining: "9,999", Reason: "ok"},
				{Call: "100", Allowed: "✓ Yes", Remaining: "9,900", Reason: "ok"},
				{Call: "9,999", Allowed: "✓ Yes", Remaining: "1", Reason: "ok"},
				{Call: "10,000", Allowed: "✓ Yes", Remaining: "0", Reason: "ok"},
				{Call: "10,001", Allowed: "❌ No", Remaining: "0", Reason: "exceeded"},
				{Call: "(Next Day)", Allowed: "✓ Yes", Remaining: "9,999", Reason: "reset"},
			},
			KeyPoints: []string{
				"Server tracks cumulative total automatically",
				"Developer only needs to call Consume()",
				"Remaining count returned for UI display",
				"Auto-resets daily/monthly per license config",
			},
		}
	case "tps":
		return &LimitExample{
			LicenseConfig: `{
  "api_access": {
    "enabled": true,
    "max_tps": 10.0   // Maximum requests per second
  }
}`,
			CodeExample: `func HandleAPIRequest() error {
    // Measure current TPS (last 1 second)
    currentTPS := rateLimiter.GetCurrentTPS()

    // Check against license limit
    allowed, maxTPS, reason, err :=
      lccClient.CheckTPS(
        "api_access",   // Feature ID
        currentTPS,     // Current rate
      )

    if err != nil {
      return fmt.Errorf("license: %w", err)
    }

    if !allowed {
      log.Warn("TPS limit exceeded",
        "current", currentTPS,
        "max", maxTPS,
        "reason", reason)
      return ErrRateLimitExceeded
    }

    // Rate OK - process request
    return processRequest()
}`,
			BehaviorTable: []BehaviorRow{
				{Call: "TPS=5.2", Allowed: "✓ Yes", Remaining: "max=10.0", Reason: "ok"},
				{Call: "TPS=9.8", Allowed: "✓ Yes", Remaining: "max=10.0", Reason: "ok"},
				{Call: "TPS=10.0", Allowed: "✓ Yes", Remaining: "max=10.0", Reason: "ok"},
				{Call: "TPS=10.5", Allowed: "❌ No", Remaining: "max=10.0", Reason: "exceeded"},
				{Call: "TPS=15.0", Allowed: "❌ No", Remaining: "max=10.0", Reason: "exceeded"},
				{Call: "(Next Sec)", Allowed: "✓ Yes", Remaining: "max=10.0", Reason: "ok (rate dropped)"},
			},
			KeyPoints: []string{
				"Client measures instantaneous rate (req/sec)",
				"CheckTPS() validates against license limit",
				"No server-side state accumulation",
				"Useful for burst control and API throttling",
			},
		}
	case "capacity":
		return &LimitExample{
			LicenseConfig: `{
  "projects": {
    "enabled": true,
    "max_capacity": 50   // Maximum total projects
  }
}`,
			CodeExample: `func CreateProject(name string) error {
    // Count current projects in database
    currentCount, err := db.CountProjects()
    if err != nil {
      return err
    }

    // Check capacity before creating
    allowed, maxCap, reason, err :=
      lccClient.CheckCapacity(
        "projects",      // Feature ID
        currentCount,    // Current count
      )

    if err != nil {
      return fmt.Errorf("license: %w", err)
    }

    if !allowed {
      log.Warn("Capacity limit reached",
        "current", currentCount,
        "max", maxCap,
        "reason", reason)
      return ErrCapacityExceeded
    }

    // Capacity OK - create new project
    return db.CreateProject(name)
}`,
			BehaviorTable: []BehaviorRow{
				{Call: "count=10", Allowed: "✓ Yes", Remaining: "max=50", Reason: "ok"},
				{Call: "count=25", Allowed: "✓ Yes", Remaining: "max=50", Reason: "ok"},
				{Call: "count=49", Allowed: "✓ Yes", Remaining: "max=50", Reason: "ok"},
				{Call: "count=50", Allowed: "❌ No", Remaining: "max=50", Reason: "at_limit"},
				{Call: "count=51", Allowed: "❌ No", Remaining: "max=50", Reason: "exceeded"},
				{Call: "(After Delete)", Allowed: "✓ Yes", Remaining: "max=50", Reason: "ok (space freed)"},
			},
			KeyPoints: []string{
				"Client counts persistent resources",
				"No time-based reset - persistent limit",
				"Requires client to track current usage",
				"Ideal for database records, storage items",
			},
		}
	case "concurrency":
		return &LimitExample{
			LicenseConfig: `{
  "concurrent_users": {
    "enabled": true,
    "max_concurrency": 10   // Max simultaneous slots
  }
}`,
			CodeExample: `func HandleUserSession(userID string) error {
    // Acquire slot at session start
    release, allowed, reason, err :=
      lccClient.AcquireSlot(
        "concurrent_users",  // Feature ID
        nil,                 // Optional metadata
      )

    if err != nil {
      return fmt.Errorf("license: %w", err)
    }

    if !allowed {
      log.Warn("Concurrency limit reached",
        "reason", reason)
      return ErrTooManyConcurrentUsers
    }

    // MUST release when done
    defer release()

    // Slot acquired - handle session
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
				"SDK tracks slots internally",
				"Must call release() when done (use defer)",
				"Real-time slot availability",
				"Ideal for sessions, jobs, connections",
			},
		}
	default:
		return nil
	}
}
