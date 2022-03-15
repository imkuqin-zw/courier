package dubbo

import (
	"errors"
	"time"

	"dubbo.apache.org/dubbo-go/v3/config_center"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	"github.com/imkuqin-zw/courier/pkg/config/source"
	"github.com/imkuqin-zw/courier/pkg/utils/xstrings"
)

type dataListener struct {
	name string
	item *item
}

type watcher struct {
	client config_center.DynamicConfiguration
	item   *item
	name   string
	cs     chan *source.ChangeSet
}

func newWatcher(client config_center.DynamicConfiguration, name string, item *item) (*watcher, error) {
	w := &watcher{
		client: client,
		item:   item,
		name:   name,
		cs:     make(chan *source.ChangeSet, 1),
	}

	w.client.AddListener(w.item.DataID, w)
	return w, nil
}

func (w *watcher) Process(configType *config_center.ConfigChangeEvent) {
	if configType.ConfigType != remoting.EventTypeUpdate {
		return
	}

	cs := &source.ChangeSet{
		Priority:  source.PriorityRemote,
		Timestamp: time.Now(),
		Source:    w.name,
		Data:      xstrings.Str2bytes(configType.Value.(string)),
		Format:    w.item.Format,
	}
	cs.Checksum = cs.Sum()
	w.cs <- cs
	return
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	cs, ok := <-w.cs
	if !ok {
		return nil, errors.New("watch chan closed")
	}
	return cs, nil
}

func (w *watcher) Stop() error {
	close(w.cs)
	w.client.RemoveListener(w.name, w)
	return nil
}
