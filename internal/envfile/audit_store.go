package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const auditFileMode = 0o600

// SaveAuditLog writes the audit log to a JSON file at the given path.
func SaveAuditLog(path string, log *AuditLog) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, auditFileMode)
}

// LoadAuditLog reads and parses an audit log from a JSON file.
// If the file does not exist, an empty AuditLog is returned.
func LoadAuditLog(path string) (*AuditLog, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewAuditLog(), nil
	}
	if err != nil {
		return nil, err
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, err
	}
	if log.Entries == nil {
		log.Entries = []AuditEntry{}
	}
	return &log, nil
}

// AppendAuditEntry loads an existing log (or creates one), appends the entry, and saves it.
func AppendAuditEntry(path string, action AuditAction, key, profile string, sensitive bool, note string) error {
	log, err := LoadAuditLog(path)
	if err != nil {
		return err
	}
	log.Record(action, key, profile, sensitive, note)
	return SaveAuditLog(path, log)
}
