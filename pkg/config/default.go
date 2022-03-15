package config

import (
	"bytes"
	"sync"
	"time"

	"github.com/imkuqin-zw/courier/pkg/config/loader"
	"github.com/imkuqin-zw/courier/pkg/config/loader/memory"
	"github.com/imkuqin-zw/courier/pkg/config/reader"
	"github.com/imkuqin-zw/courier/pkg/config/source"
)

type watchManager struct {
	mu    sync.Mutex
	watch map[string]map[*WatcherProcess]Watcher
}

func newWatchManager() *watchManager {
	return &watchManager{
		watch: map[string]map[*WatcherProcess]Watcher{},
	}
}

func (m *watchManager) Add(key string, f WatcherProcess, w Watcher) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.watch[key]
	if !ok {
		m.watch[key] = map[*WatcherProcess]Watcher{&f: w}
	}
	m.watch[key][&f] = w
	go func() {
		for {
			e := ChangeEvent{Key: key}
			e.Value, e.Err = w.Next()
			if e.Err == loader.ErrWatchStopped {
				return
			}
			f(e)
		}
	}()
}

func (m *watchManager) Del(key string, p *WatcherProcess) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	watch, ok := m.watch[key]
	if !ok {
		return nil
	}
	if w, ok := watch[p]; ok {
		if err := w.Stop(); err != nil {
			return err
		}
		delete(watch, p)
	}
	if len(watch) == 0 {
		delete(m.watch, key)
	}
	return nil
}

type watcher struct {
	lw    loader.Watcher
	rd    reader.Reader
	path  string
	value reader.Value
}

func (w *watcher) Next() (reader.Value, error) {
	for {
		s, err := w.lw.Next()
		if err != nil {
			return nil, err
		}
		// only process changes
		if bytes.Equal(w.value.Bytes(), s.ChangeSet.Data) {
			continue
		}
		v, err := w.rd.Values(s.ChangeSet)
		if err != nil {
			return nil, err
		}
		w.value = v.Get("")
		return w.value, nil
	}
}

func (w *watcher) Stop() error {
	return w.lw.Stop()
}

type config struct {
	exit chan bool
	opts Options

	sync.RWMutex
	// the current snapshot
	snap *loader.Snapshot
	// the current values
	vals reader.Values

	watchProcess *watchManager
}

func newConfig(opts ...Option) Config {
	var c config

	if err := c.init(opts...); err != nil {
		panic(err)
	}
	//go c.run()

	return &c
}

func (c *config) init(opts ...Option) error {
	c.opts = Options{
		Reader: reader.NewReader(),
	}
	c.exit = make(chan bool)
	for _, o := range opts {
		o(&c.opts)
	}

	// default loader uses the configured reader
	if c.opts.Loader == nil {
		c.opts.Loader = memory.NewLoader(memory.WithReader(c.opts.Reader))
	}

	err := c.opts.Loader.Load(c.opts.Source...)
	if err != nil {
		return err
	}

	c.snap, err = c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.vals, err = c.opts.Reader.Values(c.snap.ChangeSet)
	if err != nil {
		return err
	}
	c.watchProcess = newWatchManager()
	return nil
}

func (c *config) Options() Options {
	return c.opts
}

func (c *config) run() {
	watch := func(w loader.Watcher) error {
		for {
			// get changeset
			snap, err := w.Next()
			if err != nil {
				return err
			}

			c.Lock()

			if c.snap.Version >= snap.Version {
				c.Unlock()
				continue
			}

			// save
			c.snap = snap

			// set values
			c.vals, _ = c.opts.Reader.Values(snap.ChangeSet)

			c.Unlock()
		}
	}

	for {
		w, err := c.opts.Loader.Watch("")
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		done := make(chan bool)

		// the stop watch func
		go func() {
			select {
			case <-done:
			case <-c.exit:
			}
			w.Stop()
		}()

		// block watch
		if err := watch(w); err != nil {
			// do something better
			time.Sleep(time.Second)
		}

		// close done chan
		close(done)

		// if the conf is closed exit
		select {
		case <-c.exit:
			return
		default:
		}
	}
}

func (c *config) Map() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.vals.Map()
}

func (c *config) Scan(v interface{}) error {
	c.RLock()
	defer c.RUnlock()
	return c.vals.Scan(v)
}

// sync loads all the sources, calls the parser and constsupdates the conf
func (c *config) Sync() error {
	if err := c.opts.Loader.Sync(); err != nil {
		return err
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	vals, err := c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.vals = vals

	return nil
}

func (c *config) Close() error {
	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
	}
	return nil
}

func (c *config) Get(key string) reader.Value {
	c.RLock()
	defer c.RUnlock()

	// did sync actually work?
	if c.vals != nil {
		return c.vals.Get(key)
	}

	// no value
	return value
}

func (c *config) Set(key string, val interface{}) {
	c.Lock()
	defer c.Unlock()
	if c.vals != nil {
		c.vals.Set(key, val)
	}

	return
}

func (c *config) Del(key string) {
	c.Lock()
	defer c.Unlock()

	if c.vals != nil {
		c.vals.Del(key)
	}

	return
}

func (c *config) Bytes() []byte {
	c.RLock()
	defer c.RUnlock()

	if c.vals == nil {
		return []byte{}
	}

	return c.vals.Bytes()
}

func (c *config) Load(sources ...source.Source) error {
	if err := c.opts.Loader.Load(sources...); err != nil {
		return err
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	vals, err := c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.vals = vals

	return nil
}

func (c *config) AddWatchProcess(prefix string, f WatcherProcess) error {
	w, err := c.Watch(prefix)
	if err != nil {
		return err
	}
	c.watchProcess.Add(prefix, f, w)
	return nil
}

func (c *config) DelWatchProcess(prefix string, f WatcherProcess) error {
	return c.watchProcess.Del(prefix, &f)
}

func (c *config) Watch(prefix string) (Watcher, error) {
	value := c.Get(prefix)
	w, err := c.opts.Loader.Watch(prefix)
	if err != nil {
		return nil, err
	}

	return &watcher{
		lw:    w,
		rd:    c.opts.Reader,
		path:  prefix,
		value: value,
	}, nil
}

func (c *config) String() string {
	return "conf"
}

var value = new(noValue)

type noValue struct{}

func (v *noValue) Bool(def ...bool) bool {
	return false
}
func (v *noValue) Int(def ...int) int {
	return 0
}
func (v *noValue) String(def ...string) string {
	return ""
}
func (v *noValue) Float64(def ...float64) float64 {
	return 0.0
}
func (v *noValue) Duration(def ...time.Duration) time.Duration {
	return time.Duration(0)
}
func (v *noValue) StringSlice(def ...string) []string {
	return nil
}
func (v *noValue) StringMap(def ...map[string]string) map[string]string {
	return map[string]string{}
}
func (v *noValue) Scan(val interface{}) error {
	return nil
}
func (v *noValue) Bytes() []byte {
	return nil
}
