package web

// TierDefinition represents a product tier with its features and limits
type TierDefinition struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Tier        string                 `json:"tier"`
	ProductID   string                 `json:"product_id"`
	Description string                 `json:"description"`
	PricePoint  string                 `json:"price_point"`
	Features    map[string]FeatureInfo `json:"features"`
}

// FeatureInfo contains details about a feature in a tier
type FeatureInfo struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Enabled        bool    `json:"enabled"`
	Description    string  `json:"description,omitempty"`
	RequiredTier   string  `json:"required_tier,omitempty"`
	Reason         string  `json:"reason,omitempty"`
	Quota          *Quota  `json:"quota,omitempty"`
	MaxTPS         float64 `json:"max_tps,omitempty"`
	MaxCapacity    int     `json:"max_capacity,omitempty"`
	MaxConcurrency int     `json:"max_concurrency,omitempty"`
}

// Quota represents quota configuration for a feature
type Quota struct {
	Max      int    `json:"max"`
	Window   string `json:"window"`
	ResetAt  string `json:"reset_at"`
}

var (
	BasicTier       *TierDefinition
	ProfessionalTier *TierDefinition
	EnterpriseTier  *TierDefinition
	AllTiers        []*TierDefinition
)

func init() {
	initializeTiers()
}

func initializeTiers() {
	BasicTier = &TierDefinition{
		ID:          "basic",
		Name:        "Basic Edition",
		Tier:        "basic",
		ProductID:   "data-insight-basic",
		Description: "Essential features for individual users and small projects",
		PricePoint:  "Free or $9/month",
		Features: map[string]FeatureInfo{
			"basic_reports": {
				ID:          "basic_reports",
				Name:        "Basic Reports",
				Enabled:     true,
				Description: "Generate basic statistical reports",
			},
			"ml_analytics": {
				ID:           "ml_analytics",
				Name:         "ML Analytics",
				Enabled:      false,
				Description:  "ML-powered analytics with predictive models",
				RequiredTier: "professional",
				Reason:       "requires_professional",
			},
			"pdf_export": {
				ID:           "pdf_export",
				Name:         "PDF Export",
				Enabled:      false,
				Description:  "Professional quality PDF reports",
				RequiredTier: "professional",
				Reason:       "requires_professional",
			},
			"excel_export": {
				ID:           "excel_export",
				Name:         "Excel Export",
				Enabled:      false,
				Description:  "Advanced Excel exports with templates",
				RequiredTier: "enterprise",
				Reason:       "requires_enterprise",
			},
			"custom_dashboard": {
				ID:           "custom_dashboard",
				Name:         "Custom Dashboard",
				Enabled:      false,
				Description:  "Build custom dashboards",
				RequiredTier: "enterprise",
				Reason:       "requires_enterprise",
			},
			"api_access": {
				ID:           "api_access",
				Name:         "API Access",
				Enabled:      false,
				Description:  "REST API access",
				RequiredTier: "professional",
				Reason:       "requires_professional",
			},
		},
	}

	ProfessionalTier = &TierDefinition{
		ID:          "professional",
		Name:        "Professional Edition",
		Tier:        "professional",
		ProductID:   "data-insight-pro",
		Description: "Advanced features for growing teams and businesses",
		PricePoint:  "$49/month or $490/year",
		Features: map[string]FeatureInfo{
			"basic_reports": {
				ID:          "basic_reports",
				Name:        "Basic Reports",
				Enabled:     true,
				Description: "Generate basic statistical reports",
			},
			"ml_analytics": {
				ID:          "ml_analytics",
				Name:        "ML Analytics",
				Enabled:     true,
				Description: "ML-powered analytics with predictive models",
				Quota: &Quota{
					Max:     10000,
					Window:  "daily",
					ResetAt: "00:00",
				},
				MaxTPS: 10.0,
			},
			"pdf_export": {
				ID:          "pdf_export",
				Name:        "PDF Export",
				Enabled:     true,
				Description: "Professional quality PDF reports",
				Quota: &Quota{
					Max:     200,
					Window:  "daily",
					ResetAt: "00:00",
				},
				MaxTPS: 5.0,
			},
			"excel_export": {
				ID:           "excel_export",
				Name:         "Excel Export",
				Enabled:      false,
				Description:  "Advanced Excel exports with templates",
				RequiredTier: "enterprise",
				Reason:       "requires_enterprise",
			},
			"custom_dashboard": {
				ID:           "custom_dashboard",
				Name:         "Custom Dashboard",
				Enabled:      false,
				Description:  "Build custom dashboards",
				RequiredTier: "enterprise",
				Reason:       "requires_enterprise",
			},
			"api_access": {
				ID:             "api_access",
				Name:           "API Access",
				Enabled:        true,
				Description:    "REST API access",
				MaxTPS:         100.0,
				MaxConcurrency: 10,
			},
		},
	}

	EnterpriseTier = &TierDefinition{
		ID:          "enterprise",
		Name:        "Enterprise Edition",
		Tier:        "enterprise",
		ProductID:   "data-insight-enterprise",
		Description: "Full-featured solution for large organizations",
		PricePoint:  "$299/month or $2,990/year",
		Features: map[string]FeatureInfo{
			"basic_reports": {
				ID:          "basic_reports",
				Name:        "Basic Reports",
				Enabled:     true,
				Description: "Generate basic statistical reports",
			},
			"ml_analytics": {
				ID:          "ml_analytics",
				Name:        "ML Analytics",
				Enabled:     true,
				Description: "ML-powered analytics with predictive models",
				Quota: &Quota{
					Max:     100000,
					Window:  "daily",
					ResetAt: "00:00",
				},
				MaxTPS: 50.0,
			},
			"pdf_export": {
				ID:          "pdf_export",
				Name:        "PDF Export",
				Enabled:     true,
				Description: "Professional quality PDF reports",
				Quota: &Quota{
					Max:     2000,
					Window:  "daily",
					ResetAt: "00:00",
				},
				MaxTPS: 20.0,
			},
			"excel_export": {
				ID:          "excel_export",
				Name:        "Excel Export",
				Enabled:     true,
				Description: "Advanced Excel exports with templates",
				Quota: &Quota{
					Max:     1000,
					Window:  "daily",
					ResetAt: "00:00",
				},
				MaxTPS: 10.0,
			},
			"custom_dashboard": {
				ID:          "custom_dashboard",
				Name:        "Custom Dashboard",
				Enabled:     true,
				Description: "Build custom dashboards",
				MaxCapacity: 100,
			},
			"api_access": {
				ID:             "api_access",
				Name:           "API Access",
				Enabled:        true,
				Description:    "REST API access",
				MaxTPS:         500.0,
				MaxConcurrency: 50,
			},
		},
	}

	AllTiers = []*TierDefinition{BasicTier, ProfessionalTier, EnterpriseTier}
}

// GetTierByID returns a tier definition by its ID
func GetTierByID(tierID string) *TierDefinition {
	switch tierID {
	case "basic":
		return BasicTier
	case "professional", "pro":
		return ProfessionalTier
	case "enterprise", "ent":
		return EnterpriseTier
	default:
		return nil
	}
}

// GetLicenseJSON returns the license JSON for a tier
func GetLicenseJSON(tier *TierDefinition) map[string]interface{} {
	features := make(map[string]interface{})
	limits := make(map[string]interface{})
	
	// Features only contain enabled/disabled status
	for id, feature := range tier.Features {
		features[id] = map[string]interface{}{
			"enabled": feature.Enabled,
		}
	}
	
	// Limits are product-level configurations
	// Multiple limits can exist at product level
	switch tier.ID {
	case "basic":
		// Basic tier has no limits
		limits = map[string]interface{}{}
		
	case "professional":
		limits = map[string]interface{}{
			"quota": map[string]interface{}{
				"ml_analytics": map[string]interface{}{
					"max":      10000,
					"used":     0,
					"remaining": 10000,
					"window":   "daily",
					"reset_at": "2025-01-22T00:00:00Z",
				},
				"pdf_export": map[string]interface{}{
					"max":      200,
					"used":     0,
					"remaining": 200,
					"window":   "daily",
					"reset_at": "2025-01-22T00:00:00Z",
				},
			},
			"max_tps": map[string]interface{}{
				"ml_analytics": 10.0,
				"pdf_export":   5.0,
				"api_access":   100.0,
			},
			"max_concurrency": map[string]interface{}{
				"api_access": 10,
			},
		}
		
	case "enterprise":
		limits = map[string]interface{}{
			"quota": map[string]interface{}{
				"ml_analytics": map[string]interface{}{
					"max":      100000,
					"used":     0,
					"remaining": 100000,
					"window":   "daily",
					"reset_at": "2025-01-22T00:00:00Z",
				},
				"pdf_export": map[string]interface{}{
					"max":      2000,
					"used":     0,
					"remaining": 2000,
					"window":   "daily",
					"reset_at": "2025-01-22T00:00:00Z",
				},
				"excel_export": map[string]interface{}{
					"max":      1000,
					"used":     0,
					"remaining": 1000,
					"window":   "daily",
					"reset_at": "2025-01-22T00:00:00Z",
				},
			},
			"max_tps": map[string]interface{}{
				"ml_analytics": 50.0,
				"pdf_export":   20.0,
				"excel_export":  10.0,
				"api_access":   500.0,
			},
			"max_capacity": map[string]interface{}{
				"custom_dashboard": 100,
			},
			"max_concurrency": map[string]interface{}{
				"api_access": 50,
			},
		}
	}
	
	return map[string]interface{}{
		"product_id":   tier.ProductID,
		"product_name": tier.Name,
		"tier":         tier.Tier,
		"version":      "1.0.0",
		"issued_at":    "2025-01-21T00:00:00Z",
		"expires_at":   "2026-01-21T00:00:00Z",
		"features":     features,
		"limits":       limits,
	}
}

// GetYAMLConfig returns the YAML configuration template for a tier
func GetYAMLConfig() string {
	return `sdk:
  product_id: data-insight-pro
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
      message: "API access requires Professional tier or higher"`
}

// CheckFeatureForTier simulates checking a feature for a specific tier
func CheckFeatureForTier(tier *TierDefinition, featureID string) map[string]interface{} {
	feature, exists := tier.Features[featureID]
	if !exists {
		return map[string]interface{}{
			"enabled": false,
			"reason":  "feature_not_found",
		}
	}
	
	result := map[string]interface{}{
		"enabled": feature.Enabled,
	}
	
	if !feature.Enabled {
		result["reason"] = feature.Reason
		if feature.RequiredTier != "" {
			result["required_tier"] = feature.RequiredTier
			result["current_tier"] = tier.Tier
		}
	} else {
		result["reason"] = "ok"
	}
	
	return result
}
