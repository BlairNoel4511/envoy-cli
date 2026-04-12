package envfile

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// ErrKeyNotFound is returned when a vault lookup finds no matching entry.
var ErrKeyNotFound = errors.New("key not found in vault")

// VaultEntry holds an encrypted value alongside metadata.
type VaultEntry struct {
	Key       string    `json:"key"`
	Cipher    string    `json:"cipher"`
	CreatedAt time.Time `json:"created_at"`
}

// Vault is a collection of encrypted entries persisted as JSON.
type Vault struct {
	Entries []VaultEntry `json:"entries"`
}

// Add encrypts value and stores it under key, replacing any existing entry.
func (v *Vault) Add(key, value, passphrase string) error {
	cipher, err := Encrypt(value, passphrase)
	if err != nil {
		return err
	}
	for i, e := range v.Entries {
		if e.Key == key {
			v.Entries[i] = VaultEntry{Key: key, Cipher: cipher, CreatedAt: time.Now()}
			return nil
		}
	}
	v.Entries = append(v.Entries, VaultEntry{Key: key, Cipher: cipher, CreatedAt: time.Now()})
	return nil
}

// Get decrypts and returns the value stored under key.
func (v *Vault) Get(key, passphrase string) (string, error) {
	for _, e := range v.Entries {
		if e.Key == key {
			return Decrypt(e.Cipher, passphrase)
		}
	}
	return "", ErrKeyNotFound
}

// SaveVault writes the vault to path as JSON with restricted permissions.
func SaveVault(path string, v *Vault) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadVault reads a vault from path, returning an empty vault if not found.
func LoadVault(path string) (*Vault, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Vault{}, nil
	}
	if err != nil {
		return nil, err
	}
	var v Vault
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
