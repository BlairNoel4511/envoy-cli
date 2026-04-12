package envfile

// InjectFromVault merges decrypted vault secrets into a slice of Entry values.
// Only keys present in the vault that are also present in entries are updated.
// If overwrite is false, existing non-empty values are preserved.
func InjectFromVault(entries []Entry, v *Vault, passphrase string, overwrite bool) ([]Entry, []string, error) {
	result := make([]Entry, len(entries))
	copy(result, entries)

	var injected []string

	for i, e := range result {
		secret, err := v.Get(e.Key, passphrase)
		if err == ErrKeyNotFound {
			continue
		}
		if err != nil {
			return nil, nil, err
		}
		if !overwrite && e.Value != "" {
			continue
		}
		result[i].Value = secret
		injected = append(injected, e.Key)
	}

	return result, injected, nil
}

// ExportToVault encrypts all sensitive entries and stores them in the vault.
// Returns the list of keys that were exported.
func ExportToVault(entries []Entry, v *Vault, passphrase string) ([]string, error) {
	var exported []string
	for _, e := range entries {
		if !IsSensitive(e.Key) {
			continue
		}
		if err := v.Add(e.Key, e.Value, passphrase); err != nil {
			return nil, err
		}
		exported = append(exported, e.Key)
	}
	return exported, nil
}
