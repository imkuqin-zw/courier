package dubbo

import (
	"fmt"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common/config"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	"github.com/imkuqin-zw/courier/pkg/config/source"
	"github.com/imkuqin-zw/courier/pkg/utils/xstrings"
)

func NewDubbov3Source(Key string, opts ...Option) source.Source {
	dynamicCfg := config.GetEnvInstance().GetDynamicConfiguration()
	if dynamicCfg == nil {
		return nil
	}
	o := defaultOpts()
	for _, f := range opts {
		f(o)
	}
	return &dubbov3Source{
		item: &item{
			DataID: Key,
			Format: o.Format,
		},
		client: dynamicCfg,
	}
}

type item struct {
	DataID string
	Format string
}

type dubbov3Source struct {
	item   *item
	client config_center.DynamicConfiguration
}

func (n *dubbov3Source) Read() (*source.ChangeSet, error) {
	content, err := n.client.GetProperties(n.item.DataID)
	if err != nil {
		return nil, fmt.Errorf("fault to get config (dataID: %s): %s", n.item.DataID, err.Error())
	}
	cs := &source.ChangeSet{
		Priority:  source.PriorityRemote,
		Timestamp: time.Now(),
		Source:    n.String(),
		Data:      xstrings.Str2bytes(content),
		Format:    n.item.Format,
	}
	cs.Checksum = cs.Sum()
	return cs, nil
}

func (n *dubbov3Source) Watch() (source.Watcher, error) {
	return newWatcher(n.client, n.String(), n.item)
}

func (n *dubbov3Source) String() string {
	return "dubbov3"
}
