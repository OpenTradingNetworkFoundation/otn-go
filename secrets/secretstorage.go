package secrets

import "github.com/juju/errors"

// SecretStorage defines interface for reading from a confidential storage
type SecretStorage interface {
	ReadStringValue(path string) (string, error)
	ReadStringArray(path string) ([]string, error)
}

type StorageConfig struct {
	Vault *VaultStorageConfig `json:"vault"`
	Local *LocalStorageConfig `json:"local"`
}

func NewSecretStorage(cfg *StorageConfig) (SecretStorage, error) {
	if cfg.Vault != nil {
		vault, err := NewVaultStorage(cfg.Vault)
		if err != nil {
			return nil, err
		}
		return vault, nil
	}

	if cfg.Local != nil {
		storage := &LocalStorage{cfg.Local}
		return storage, nil
	}

	return nil, errors.New("Invalid configuration")
}
