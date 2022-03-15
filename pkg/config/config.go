// Package conf is an interface for dynamic configuration.
package config

import (
	"github.com/imkuqin-zw/courier/pkg/config/reader"
	"github.com/imkuqin-zw/courier/pkg/config/source"
)

// Config is an interface abstraction for dynamic configuration
type Config interface {
	// provide the reader.Values interface
	reader.Values
	// Options in the conf
	Options() Options
	// Stop the conf loader/watcher
	Close() error
	// Load conf sources
	Load(source ...source.Source) error
	// Force a source changeset sync
	Sync() error
	// Watch a value for changes
	Watch(prefix string) (Watcher, error)
	AddWatchProcess(prefix string, f WatcherProcess) error
	DelWatchProcess(prefix string, f WatcherProcess) error
}

type ChangeEvent struct {
	Key   string
	Value reader.Value
	Err   error
}

type WatcherProcess func(ChangeEvent)

// Watcher is the conf watcher
type Watcher interface {
	Next() (reader.Value, error)
	Stop() error
}

var defaultConfig = newConfig()

// Return conf as raw json
func Bytes() []byte {
	return defaultConfig.Bytes()
}

// Return conf as a map
func Map() map[string]interface{} {
	return defaultConfig.Map()
}

// Scan values to a go type
func Scan(v interface{}) error {
	return defaultConfig.Scan(v)
}

// Force a source changeset sync
func Sync() error {
	return defaultConfig.Sync()
}

// Get a value from the conf
func Get(key string) reader.Value {
	return defaultConfig.Get(key)
}

// Set a value to the conf
func Set(key string, val interface{}) {
	defaultConfig.Set(key, val)
}

// Del a value from the conf
func Del(key string) {
	defaultConfig.Del(key)
}

// Load conf sources
func Load(source ...source.Source) error {
	return defaultConfig.Load(source...)
}

// Watch a value for changes
func Watch(prefix string) (Watcher, error) {
	return defaultConfig.Watch(prefix)
}

func AddWatchProcess(prefix string, f WatcherProcess) error {
	return defaultConfig.AddWatchProcess(prefix, f)
}

func DelWatchProcess(prefix string, f WatcherProcess) error {
	return defaultConfig.DelWatchProcess(prefix, f)
}
