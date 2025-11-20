package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PublicProduct represents a product returned by LCC public service.
type PublicProduct struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type PublicServiceClient struct {
	base string
	cli  *http.Client
}

// PublicFeature describes a feature entry from public service
type PublicFeature struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

func NewPublicServiceClient(base string) *PublicServiceClient {
	return &PublicServiceClient{
		base: strings.TrimRight(base, "/"),
		cli:  &http.Client{Timeout: 5 * time.Second},
	}
}

// ListProducts calls the public endpoint to list available products.
// Expected base already points to public API root (e.g., http://host:port/api/v1/public)
func (p *PublicServiceClient) ListProducts(ctx context.Context) ([]PublicProduct, error) {
	u, err := url.Parse(p.base + "/products")
	if err != nil { return nil, err }
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil { return nil, err }
	resp, err := p.cli.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { return nil, fmt.Errorf("status=%d", resp.StatusCode) }
	var out []PublicProduct
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return nil, err }
	return out, nil
}

// ListFeatures lists features of a product via public service endpoint
// Expected base already points to public API root (e.g., http://host:port/api/v1/public)
func (p *PublicServiceClient) ListFeatures(ctx context.Context, productID string) ([]PublicFeature, error) {
	if productID == "" { return nil, fmt.Errorf("productID required") }
	u, err := url.Parse(p.base + "/products/" + url.PathEscape(productID) + "/features")
	if err != nil { return nil, err }
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil { return nil, err }
	resp, err := p.cli.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { return nil, fmt.Errorf("status=%d", resp.StatusCode) }
	var out []PublicFeature
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return nil, err }
	return out, nil
}
