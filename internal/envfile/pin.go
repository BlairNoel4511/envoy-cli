package envfile

import (
	"fmt"
	"strings"
)

// PinEntry represents a key pinned to a specific environment.
type PinEntry struct {
	Key     string `json:"key"`
	Env     string `json:"env"`
	Reason  string `json:"reason,omitempty"`
}

// PinStore holds pinned keys per environment.
type PinStore struct {
	Pins []PinEntry `json:"pins"`
}

// NewPinStore creates an empty PinStore.
func NewPinStore() *PinStore {
	return &PinStore{Pins: []PinEntry{}}
}

// Add pins a key to an environment with an optional reason.
func (p *PinStore) Add(key, env, reason string) {
	for i, pin := range p.Pins {
		if pin.Key == key && pin.Env == env {
			p.Pins[i].Reason = reason
			return
		}
	}
	p.Pins = append(p.Pins, PinEntry{Key: key, Env: env, Reason: reason})
}

// Remove unpins a key from an environment.
func (p *PinStore) Remove(key, env string) bool {
	for i, pin := range p.Pins {
		if pin.Key == key && pin.Env == env {
			p.Pins = append(p.Pins[:i], p.Pins[i+1:]...)
			return true
		}
	}
	return false
}

// IsPinned returns true if the key is pinned to the given environment.
func (p *PinStore) IsPinned(key, env string) bool {
	for _, pin := range p.Pins {
		if pin.Key == key && pin.Env == env {
			return true
		}
	}
	return false
}

// ForEnv returns all pinned keys for a given environment.
func (p *PinStore) ForEnv(env string) []PinEntry {
	var result []PinEntry
	for _, pin := range p.Pins {
		if pin.Env == env {
			result = append(result, pin)
		}
	}
	return result
}

// FormatPinList formats pinned entries for display.
func FormatPinList(pins []PinEntry) string {
	if len(pins) == 0 {
		return "  (no pinned keys)"
	}
	var sb strings.Builder
	for _, p := range pins {
		line := fmt.Sprintf("  [%s] %s", p.Env, p.Key)
		if p.Reason != "" {
			line += fmt.Sprintf(" — %s", p.Reason)
		}
		sb.WriteString(line + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
