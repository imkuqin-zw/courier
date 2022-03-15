// Package file is a file source. Expected format is json
package file

import (
	"io/ioutil"
	"os"

	"github.com/imkuqin-zw/courier/pkg/config/source"
)

type file struct {
	path  string
	watch bool
	opts  source.Options
}

const defaultPath = "conf.json"

func (f *file) Read() (*source.ChangeSet, error) {
	fh, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}
	info, err := fh.Stat()
	if err != nil {
		return nil, err
	}
	cs := &source.ChangeSet{
		Priority:  source.PriorityFile,
		Format:    format(f.path, f.opts.Encoder),
		Source:    f.String(),
		Timestamp: info.ModTime(),
		Data:      b,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func (f *file) String() string {
	return "file"
}

func (f *file) Watch() (source.Watcher, error) {
	if _, err := os.Stat(f.path); err != nil {
		return nil, err
	}
	if !f.watch {
		return nil, nil
	}
	return newWatcher(f)
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)
	path := defaultPath
	f, ok := options.Context.Value(filePathKey{}).(string)
	if ok {
		path = f
	}
	watch := false
	b, ok := options.Context.Value(watchKey{}).(bool)
	if ok {
		watch = b
	}
	return &file{opts: options, path: path, watch: watch}
}
