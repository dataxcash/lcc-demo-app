package web

import (
	"encoding/json"
	"net/http"
	"strings"
)

// handleGetTiers returns all available tiers
func (s *Server) handleGetTiers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	_ = json.NewEncoder(w).Encode(AllTiers)
}

// handleGetTierLicense returns the license JSON for a specific tier
func (s *Server) handleGetTierLicense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tierID := extractTierFromPath(r.URL.Path, "/api/tiers/", "/license")
	if tierID == "" {
		writeErr(w, http.StatusBadRequest, nil)
		return
	}

	tier := GetTierByID(tierID)
	if tier == nil {
		writeErr(w, http.StatusNotFound, nil)
		return
	}

	license := GetLicenseJSON(tier)
	_ = json.NewEncoder(w).Encode(license)
}

// handleGetTierYAML returns the YAML configuration for a specific tier
func (s *Server) handleGetTierYAML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tierID := extractTierFromPath(r.URL.Path, "/api/tiers/", "/yaml")
	if tierID == "" {
		writeErr(w, http.StatusBadRequest, nil)
		return
	}

	tier := GetTierByID(tierID)
	if tier == nil {
		writeErr(w, http.StatusNotFound, nil)
		return
	}

	yamlContent := GetYAMLConfig()
	_ = json.NewEncoder(w).Encode(map[string]string{
		"yaml_content": yamlContent,
	})
}

// handleCheckTierFeature checks if a feature is enabled for a specific tier
func (s *Server) handleCheckTierFeature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tierID := extractTierFromPath(r.URL.Path, "/api/tiers/", "/check-feature")
	if tierID == "" {
		writeErr(w, http.StatusBadRequest, nil)
		return
	}

	tier := GetTierByID(tierID)
	if tier == nil {
		writeErr(w, http.StatusNotFound, nil)
		return
	}

	var req struct {
		FeatureID string `json:"feature_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	if req.FeatureID == "" {
		writeErr(w, http.StatusBadRequest, nil)
		return
	}

	result := CheckFeatureForTier(tier, req.FeatureID)
	_ = json.NewEncoder(w).Encode(result)
}

// extractTierFromPath extracts tier ID from URL path
// Example: /api/tiers/professional/license -> "professional"
func extractTierFromPath(path, prefix, suffix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	if !strings.HasSuffix(path, suffix) {
		return ""
	}

	path = strings.TrimPrefix(path, prefix)
	path = strings.TrimSuffix(path, suffix)
	
	return path
}
