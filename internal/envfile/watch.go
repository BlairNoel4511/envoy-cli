package envfile

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchEvent describes a change detected in a watched .env file.
type WatchEvent struct {
	Path      string
	ChangedAt time.Time
	OldSum    string
	NewSum    string
}

// WatchOptions configures the file watcher behaviour.
type WatchOptions struct {
	// Interval between polls.
	Interval time.Duration
	// MaxChecks limits how many polls are performed (0 = unlimited).
	MaxChecks int
}

// DefaultWatchOptions returns sensible defaults.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Interval:  2 * time.Second,
		MaxChecks: 0,
	}
}

// fileChecksum returns a hex SHA-256 digest of the file at path.
func fileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Watch polls path at the given interval and sends a WatchEvent on the
// returned channel whenever the file checksum changes. The caller must
// close done to stop the watcher.
func Watch(path string, opts WatchOptions, done <-chan struct{}) (<-chan WatchEvent, error) {
	initialSum, err := fileChecksum(path)
	if err != nil {
		return nil, fmt.Errorf("watch: initial checksum failed: %w", err)
	}

	events := make(chan WatchEvent, 4)

	go func() {
		defer close(events)
		current := initialSum
		checks := 0

		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				newSum, err := fileChecksum(path)
				if err != nil {
					// file may have been temporarily unavailable; skip
					continue
				}
				if newSum != current {
					events <- WatchEvent{
						Path:      path,
						ChangedAt: t,
						OldSum:    current,
						NewSum:    newSum,
					}
					current = newSum
				}
				checks++
				if opts.MaxChecks > 0 && checks >= opts.MaxChecks {
					return
				}
			}
		}
	}()

	return events, nil
}
