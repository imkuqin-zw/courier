// package loader manages loading from multiple sources
package loader

import (
	"context"
	"errors"

	"github.com/imkuqin-zw/courier/pkg/config/reader"
	"github.com/imkuqin-zw/courier/pkg/config/source"
)

var ErrWatchStopped = errors.New("watcher stopped")

// Loader manages loading sources
type Loader interface {
	// Stop the loader
	Close() error
	// Load the sources
	Load(...source.Source) error
	// A Snapshot of loaded conf
	Snapshot() (*Snapshot, error)
	// Force sync of sources
	Sync() error
	// Watch for changes
	Watch(string) (Watcher, error)
	// Name of loader
	String() string
}

// Watcher lets you watch sources and returns a merged ChangeSet
type Watcher interface {
	// First call to next may return the current Snapshot
	// If you are watching a path then only the repo from
	// that path is returned.
	Next() (*Snapshot, error)
	// Stop watching for changes
	Stop() error
}

// Snapshot is a merged ChangeSet
type Snapshot struct {
	// The merged ChangeSet
	ChangeSet *source.ChangeSet
	// Deterministic and comparable version of the snapshot
	Version string
}

type Options struct {
	Reader reader.Reader
	Source []source.Source

	// for alternative repo
	Context context.Context
}

type Option func(o *Options)

// Copy snapshot
func Copy(s *Snapshot) *Snapshot {
	cs := *(s.ChangeSet)

	return &Snapshot{
		ChangeSet: &cs,
		Version:   s.Version,
	}
}
