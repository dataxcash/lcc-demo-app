package web

import (
	"os"

	lccconfig "github.com/yourorg/lcc-sdk/pkg/config"
)

var candidateManifests = []string{
	"lcc-features.yaml",
	"configs/lcc-features.basic.yaml",
	"configs/lcc-features.pro.yaml",
	"configs/lcc-features.ent.yaml",
}

// LoadFeaturesForProduct loads features from local manifests matching sdk.product_id.
func LoadFeaturesForProduct(productID string) ([]PublicFeature, error) {
	for _, p := range candidateManifests {
		if !fileExists(p) { continue }
		mf, err := lccconfig.LoadManifest(p)
		if err != nil { continue }
		if mf.SDK.ProductID == productID {
			out := make([]PublicFeature, 0, len(mf.Features))
			for _, f := range mf.Features {
				out = append(out, PublicFeature{ID: f.ID, Name: f.Name})
			}
			return out, nil
		}
	}
	return nil, nil
}

// LoadFeatureUnion returns the union of all features across local manifests.
func LoadFeatureUnion() ([]PublicFeature, error) {
	seen := map[string]bool{}
	var out []PublicFeature
	for _, p := range candidateManifests {
		if !fileExists(p) { continue }
		mf, err := lccconfig.LoadManifest(p)
		if err != nil { continue }
		for _, f := range mf.Features {
			if seen[f.ID] { continue }
			seen[f.ID] = true
			out = append(out, PublicFeature{ID: f.ID, Name: f.Name})
		}
	}
	return out, nil
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}