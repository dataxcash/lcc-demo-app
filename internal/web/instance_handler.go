package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourorg/lcc-sdk/pkg/auth"
	lccclient "github.com/yourorg/lcc-sdk/pkg/client"
	lccconfig "github.com/yourorg/lcc-sdk/pkg/config"
)

type RegisterInstanceRequest struct {
	ProductID string `json:"product_id"`
	Version   string `json:"version"`
	LCCURL    string `json:"lcc_url,omitempty"`
}

type RegisterInstanceResponse struct {
	Success     bool   `json:"success"`
	InstanceID  string `json:"instance_id"`
	ProductID   string `json:"product_id"`
	Version     string `json:"version"`
	RegisteredAt string `json:"registered_at"`
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
}

type InstanceStatusResponse struct {
	ProductID  string              `json:"product_id"`
	InstanceID string              `json:"instance_id"`
	Status     string              `json:"status"`
	Features   []featureStatusDTO  `json:"features"`
	RegisteredAt string            `json:"registered_at,omitempty"`
}

type TestInstanceRequest struct {
	ProductID string `json:"product_id"`
	FeatureID string `json:"feature_id"`
}

type TestInstanceResponse struct {
	Success   bool   `json:"success"`
	Enabled   bool   `json:"enabled"`
	Reason    string `json:"reason"`
	FeatureID string `json:"feature_id"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
}

type ClearInstanceRequest struct {
	ProductID  string `json:"product_id"`
	InstanceID string `json:"instance_id,omitempty"`
}

type ClearInstanceResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type GenerateKeysResponse struct {
	Success    bool   `json:"success"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
}

func (s *Server) handleInstanceRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req RegisterInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(&RegisterInstanceResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid json: %v", err),
		})
		return
	}

	if req.ProductID == "" {
		_ = json.NewEncoder(w).Encode(&RegisterInstanceResponse{
			Success: false,
			Error:   "product_id is required",
		})
		return
	}

	if req.Version == "" {
		req.Version = "1.0.0"
	}

	s.mu.RLock()
	lccURL := s.lccURL
	s.mu.RUnlock()

	if req.LCCURL != "" {
		lccURL = req.LCCURL
	}

	if lccURL == "" {
		_ = json.NewEncoder(w).Encode(&RegisterInstanceResponse{
			Success: false,
			Error:   "lcc_url not configured",
		})
		return
	}

	cfg := &lccconfig.SDKConfig{
		LCCURL:         lccURL,
		ProductID:      req.ProductID,
		ProductVersion: req.Version,
		Timeout:        10 * time.Second,
		CacheTTL:       5 * time.Second,
	}

	ks, _ := NewKeyStore()
	var kp *auth.KeyPair

	if ks != nil {
		if loaded, err := ks.Load(req.ProductID); err == nil && loaded != nil {
			kp = loaded
		}
	}

	if kp == nil {
		var genErr error
		kp, genErr = auth.GenerateKeyPair()
		if genErr != nil {
			_ = json.NewEncoder(w).Encode(&RegisterInstanceResponse{
				Success: false,
				Error:   fmt.Sprintf("keypair generation failed: %v", genErr),
			})
			return
		}
		if ks != nil {
			_ = ks.Save(req.ProductID, kp)
		}
	}

	cli, err := lccclient.NewClientWithKeyPair(cfg, kp)
	if err != nil {
		_ = json.NewEncoder(w).Encode(&RegisterInstanceResponse{
			Success: false,
			Error:   fmt.Sprintf("client creation failed: %v", err),
		})
		return
	}

	if err := cli.Register(); err != nil {
		_ = json.NewEncoder(w).Encode(&RegisterInstanceResponse{
			Success: false,
			Error:   fmt.Sprintf("registration failed: %v", err),
		})
		return
	}

	instanceID := cli.GetInstanceID()
	registeredAt := time.Now()

	s.mu.Lock()
	// Store in old map for backward compatibility
	s.clients[req.ProductID] = cli
	// Store in new multi-instance map with unique key per product-version
	instanceKey := fmt.Sprintf("%s:%s:%s", req.ProductID, req.Version, instanceID)
	s.instances[instanceKey] = &Instance{
		InstanceID:   instanceID,
		ProductID:    req.ProductID,
		Version:      req.Version,
		RegisteredAt: registeredAt,
	}
	s.instanceKeys[instanceKey] = kp
	s.mu.Unlock()

	_ = json.NewEncoder(w).Encode(&RegisterInstanceResponse{
		Success:      true,
		InstanceID:   instanceID,
		ProductID:    req.ProductID,
		Version:      req.Version,
		RegisteredAt: registeredAt.Format(time.RFC3339),
		Message:      "Instance registered successfully",
	})
}

func (s *Server) handleInstanceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("product_id is required"))
		return
	}

	cli, err := s.getClient(productID)
	if err != nil {
		_ = json.NewEncoder(w).Encode(&InstanceStatusResponse{
			ProductID: productID,
			Status:    "not_registered",
			Features:  []featureStatusDTO{},
		})
		return
	}

	features, _ := LoadFeaturesForProduct(productID)
	if len(features) == 0 {
		features, _ = LoadFeatureUnion()
	}

	out := make([]featureStatusDTO, 0, len(features))
	for _, f := range features {
		st, err := cli.CheckFeature(f.ID)
		if err != nil {
			out = append(out, featureStatusDTO{
				ID:      f.ID,
				Name:    f.Name,
				Enabled: false,
				Reason:  "check_error",
			})
			continue
		}
		out = append(out, featureStatusDTO{
			ID:             f.ID,
			Name:           f.Name,
			Enabled:        st.Enabled,
			Reason:         st.Reason,
			Quota:          st.Quota,
			MaxCapacity:    st.MaxCapacity,
			MaxTPS:         st.MaxTPS,
			MaxConcurrency: st.MaxConcurrency,
		})
	}

	_ = json.NewEncoder(w).Encode(&InstanceStatusResponse{
		ProductID:  productID,
		InstanceID: cli.GetInstanceID(),
		Status:     "active",
		Features:   out,
	})
}

func (s *Server) handleInstanceTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req TestInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(&TestInstanceResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid json: %v", err),
		})
		return
	}

	if req.ProductID == "" || req.FeatureID == "" {
		_ = json.NewEncoder(w).Encode(&TestInstanceResponse{
			Success: false,
			Error:   "product_id and feature_id are required",
		})
		return
	}

	cli, err := s.getClient(req.ProductID)
	if err != nil {
		_ = json.NewEncoder(w).Encode(&TestInstanceResponse{
			Success: false,
			Error:   fmt.Sprintf("product not registered: %v", err),
		})
		return
	}

	status, err := cli.CheckFeature(req.FeatureID)
	if err != nil {
		_ = json.NewEncoder(w).Encode(&TestInstanceResponse{
			Success:   false,
			FeatureID: req.FeatureID,
			Error:     fmt.Sprintf("check failed: %v", err),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(&TestInstanceResponse{
		Success:   true,
		Enabled:   status.Enabled,
		Reason:    status.Reason,
		FeatureID: req.FeatureID,
		Message:   "Feature check successful",
	})
}

func (s *Server) handleInstanceClear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req ClearInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(&ClearInstanceResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid json: %v", err),
		})
		return
	}

	if req.ProductID == "" {
		_ = json.NewEncoder(w).Encode(&ClearInstanceResponse{
			Success: false,
			Error:   "product_id is required",
		})
		return
	}

	s.mu.Lock()
	// Delete from old single-instance map
	delete(s.clients, req.ProductID)

	// Delete from new multi-instance map if instanceID provided
	if req.InstanceID != "" {
		for key := range s.instances {
			if s.instances[key].InstanceID == req.InstanceID {
				delete(s.instances, key)
				delete(s.instanceKeys, key)
				break
			}
		}
	} else {
		// If no instanceID, delete all instances for this product
		for key := range s.instances {
			if s.instances[key].ProductID == req.ProductID {
				delete(s.instances, key)
				delete(s.instanceKeys, key)
			}
		}
	}
	s.mu.Unlock()

	_ = json.NewEncoder(w).Encode(&ClearInstanceResponse{
		Success: true,
		Message: "Instance cleared successfully",
	})
}

func (s *Server) handleInstanceGenerateKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	kp, err := auth.GenerateKeyPair()
	if err != nil {
		_ = json.NewEncoder(w).Encode(&GenerateKeysResponse{
			Success: false,
			Error:   fmt.Sprintf("key generation failed: %v", err),
		})
		return
	}

	publicKeyPEM, err := kp.GetPublicKeyPEM()
	if err != nil {
		_ = json.NewEncoder(w).Encode(&GenerateKeysResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to export public key: %v", err),
		})
		return
	}

	privateKeyPEM, err := kp.ExportPrivateKeyPEM()
	if err != nil {
		_ = json.NewEncoder(w).Encode(&GenerateKeysResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to export private key: %v", err),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(&GenerateKeysResponse{
		Success:    true,
		PublicKey:  publicKeyPEM,
		PrivateKey: privateKeyPEM,
		Message:    "Keys generated successfully",
	})
}

// List all registered instances
type ListInstancesResponse struct {
	Instances []Instance `json:"instances"`
}

func (s *Server) handleListInstances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	s.mu.RLock()
	instances := make([]Instance, 0, len(s.instances))
	for _, inst := range s.instances {
		instances = append(instances, *inst)
	}
	s.mu.RUnlock()

	_ = json.NewEncoder(w).Encode(&ListInstancesResponse{Instances: instances})
}
