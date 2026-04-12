package envfile

import (
	"fmt"
	"time"
)

// LineageEvent represents a single change event in the history of an env entry.
type LineageEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"` // "set", "delete", "import"
	Key       string    `json:"key"`
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty"`
	Source    string    `json:"source,omitempty"` // e.g. "profile:staging", "vault", "manual"
}

// Lineage tracks the change history of env entries.
type Lineage struct {
	Events []LineageEvent `json:"events"`
}

// NewLineage creates an empty Lineage.
func NewLineage() *Lineage {
	return &Lineage{Events: []LineageEvent{}}
}

// Record appends a new event to the lineage.
func (l *Lineage) Record(action, key, oldValue, newValue, source string) {
	l.Events = append(l.Events, LineageEvent{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		OldValue:  oldValue,
		NewValue:  newValue,
		Source:    source,
	})
}

// ForKey returns all events for a specific key.
func (l *Lineage) ForKey(key string) []LineageEvent {
	var result []LineageEvent
	for _, e := range l.Events {
		if e.Key == key {
			result = append(result, e)
		}
	}
	return result
}

// Summary returns a human-readable summary string for a lineage event.
func (e LineageEvent) Summary() string {
	ts := e.Timestamp.Format(time.RFC3339)
	switch e.Action {
	case "set":
		if e.OldValue == "" {
			return fmt.Sprintf("[%s] SET %s (source: %s)", ts, e.Key, e.Source)
		}
		return fmt.Sprintf("[%s] UPDATED %s (source: %s)", ts, e.Key, e.Source)
	case "delete":
		return fmt.Sprintf("[%s] DELETED %s (source: %s)", ts, e.Key, e.Source)
	case "import":
		return fmt.Sprintf("[%s] IMPORTED %s (source: %s)", ts, e.Key, e.Source)
	default:
		return fmt.Sprintf("[%s] %s %s", ts, e.Action, e.Key)
	}
}

// TrackDiff records lineage events derived from a Diff result.
func (l *Lineage) TrackDiff(d DiffResult, source string) {
	for _, e := range d.Added {
		l.Record("set", e.Key, "", e.Value, source)
	}
	for _, e := range d.Removed {
		l.Record("delete", e.Key, e.Value, "", source)
	}
	for _, c := range d.Changed {
		l.Record("set", c.Key, c.OldValue, c.NewValue, source)
	}
}
