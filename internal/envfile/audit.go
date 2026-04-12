package envfile

import (
	"fmt"
	"time"
)

// AuditAction represents the type of action recorded in an audit log entry.
type AuditAction string

const (
	AuditActionSet    AuditAction = "SET"
	AuditActionDelete AuditAction = "DELETE"
	AuditActionSync   AuditAction = "SYNC"
	AuditActionExport AuditAction = "EXPORT"
)

// AuditEntry represents a single recorded event in the audit log.
type AuditEntry struct {
	Timestamp time.Time   `json:"timestamp"`
	Action    AuditAction `json:"action"`
	Key       string      `json:"key"`
	Profile   string      `json:"profile,omitempty"`
	Redacted  bool        `json:"redacted"`
	Note      string      `json:"note,omitempty"`
}

// AuditLog holds a list of audit entries.
type AuditLog struct {
	Entries []AuditEntry `json:"entries"`
}

// NewAuditLog creates an empty AuditLog.
func NewAuditLog() *AuditLog {
	return &AuditLog{Entries: []AuditEntry{}}
}

// Record appends a new entry to the audit log.
func (a *AuditLog) Record(action AuditAction, key, profile string, sensitive bool, note string) {
	a.Entries = append(a.Entries, AuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		Profile:   profile,
		Redacted:  sensitive,
		Note:      note,
	})
}

// FilterByAction returns all entries matching the given action.
func (a *AuditLog) FilterByAction(action AuditAction) []AuditEntry {
	var result []AuditEntry
	for _, e := range a.Entries {
		if e.Action == action {
			result = append(result, e)
		}
	}
	return result
}

// FilterByKey returns all entries for a specific key.
func (a *AuditLog) FilterByKey(key string) []AuditEntry {
	var result []AuditEntry
	for _, e := range a.Entries {
		if e.Key == key {
			result = append(result, e)
		}
	}
	return result
}

// Summary returns a human-readable summary line for an audit entry.
func AuditEntrySummary(e AuditEntry) string {
	redactedTag := ""
	if e.Redacted {
		redactedTag = " [redacted]"
	}
	profileTag := ""
	if e.Profile != "" {
		profileTag = fmt.Sprintf(" (profile: %s)", e.Profile)
	}
	return fmt.Sprintf("%s  %-8s  %s%s%s",
		e.Timestamp.Format(time.RFC3339),
		string(e.Action),
		e.Key,
		redactedTag,
		profileTag,
	)
}
