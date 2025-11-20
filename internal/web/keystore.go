package web

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourorg/lcc-sdk/pkg/auth"
)

// KeyStore stores per-product private keys to keep instance IDs stable.
type KeyStore struct {
	baseDir string
}

func NewKeyStore() (*KeyStore, error) {
	h, err := os.UserHomeDir()
	if err != nil { return nil, err }
	root := filepath.Join(h, ".lcc-demo", "keys")
	if err := os.MkdirAll(root, 0700); err != nil { return nil, err }
	return &KeyStore{ baseDir: root }, nil
}

func (ks *KeyStore) pathFor(productID string) (string, error) {
	if productID == "" { return "", fmt.Errorf("empty productID") }
	return filepath.Join(ks.baseDir, productID+".pem"), nil
}

func (ks *KeyStore) Load(productID string) (*auth.KeyPair, error) {
	p, err := ks.pathFor(productID)
	if err != nil { return nil, err }
	if _, err := os.Stat(p); err != nil {
		return nil, err
	}
	return auth.LoadKeyPairFromPEMFile(p)
}

func (ks *KeyStore) Save(productID string, kp *auth.KeyPair) error {
	p, err := ks.pathFor(productID)
	if err != nil { return err }
	return kp.SavePrivateKeyPEMFile(p)
}
