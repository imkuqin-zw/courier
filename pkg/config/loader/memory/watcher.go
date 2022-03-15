package memory

import (
	"bytes"
	"time"

	"github.com/imkuqin-zw/courier/pkg/config/loader"
	"github.com/imkuqin-zw/courier/pkg/config/reader"
	"github.com/imkuqin-zw/courier/pkg/config/source"
)

type watcher struct {
	exit    chan bool
	path    string
	value   reader.Value
	reader  reader.Reader
	version string
	updates chan updateValue
}

func (w *watcher) Next() (*loader.Snapshot, error) {
	update := func(v reader.Value) *loader.Snapshot {
		w.value = v
		cs := &source.ChangeSet{
			Priority:  source.PriorityMax,
			Data:      v.Bytes(),
			Format:    w.reader.String(),
			Source:    "watcher",
			Timestamp: time.Now(),
		}
		cs.Checksum = cs.Sum()
		return &loader.Snapshot{
			ChangeSet: cs,
			Version:   w.version,
		}
	}

	for {
		select {
		case <-w.exit:
			return nil, loader.ErrWatchStopped

		case uv := <-w.updates:
			if uv.version <= w.version {
				continue
			}

			v := uv.value

			w.version = uv.version

			if bytes.Equal(w.value.Bytes(), v.Bytes()) {
				continue
			}
			return update(v), nil
		}
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
	default:
		close(w.exit)
		close(w.updates)
	}

	return nil
}
