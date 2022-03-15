package reader

import (
	"errors"
	"sort"
	"time"

	"github.com/imdario/mergo"
	"github.com/imkuqin-zw/courier/pkg/config/encoder"
	"github.com/imkuqin-zw/courier/pkg/config/encoder/json"
	"github.com/imkuqin-zw/courier/pkg/config/source"
)

type jsonReader struct {
	opts Options
	json encoder.Encoder
}

func (j *jsonReader) merge(changes ...*source.ChangeSet) (*source.ChangeSet, error) {
	var merged = make(map[string]interface{})
	for _, m := range changes {
		if m == nil {
			continue
		}

		if len(m.Data) == 0 {
			continue
		}

		codec, ok := j.opts.Encoding[m.Format]
		if !ok {
			// fallback
			codec = j.json
		}

		var data map[string]interface{}
		if err := codec.Decode(m.Data, &data); err != nil {
			return nil, err
		}
		if err := mergo.Map(&merged, data, mergo.WithOverride); err != nil {
			return nil, err
		}
	}
	b, err := j.json.Encode(merged)
	if err != nil {
		return nil, err
	}
	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Data:      b,
		Source:    "json",
		Format:    j.json.String(),
	}
	cs.Checksum = cs.Sum()
	return cs, nil
}

func (j *jsonReader) Merge(changes ...*source.ChangeSet) (*source.ChangeSet, error) {
	changesMap := make(map[source.Priority][]*source.ChangeSet)
	for _, change := range changes {
		if change == nil {
			continue
		}
		if len(change.Data) == 0 {
			continue
		}
		_, ok := changesMap[change.Priority]
		if !ok {
			changesMap[change.Priority] = []*source.ChangeSet{change}
		} else {
			changesMap[change.Priority] = append(changesMap[change.Priority], change)
		}
	}
	pChanges := make([]*source.ChangeSet, 0, len(changesMap))
	for _, ms := range changesMap {
		m, err := j.merge(ms...)
		if err != nil {
			return nil, err
		}
		pChanges = append(pChanges, m)
	}
	sort.SliceStable(pChanges, func(i, j int) bool {
		return changes[i].Priority < changes[j].Priority
	})
	return j.merge(pChanges...)
}

func (j *jsonReader) Values(ch *source.ChangeSet) (Values, error) {
	if ch == nil {
		return nil, errors.New("changeset is nil")
	}
	if ch.Format != "json" {
		return nil, errors.New("unsupported format")
	}
	return newValues(ch)
}

func (j *jsonReader) String() string {
	return "json"
}

// NewReader creates a json reader
func NewReader(opts ...Option) Reader {
	options := NewOptions(opts...)
	return &jsonReader{
		json: json.NewEncoder(),
		opts: options,
	}
}
