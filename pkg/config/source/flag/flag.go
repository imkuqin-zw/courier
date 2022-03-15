package flag

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/imkuqin-zw/courier/pkg/config/source"
)

type flagsrc struct {
	fs   *flag.FlagSet
	opts source.Options
}

func (fs *flagsrc) Read() (*source.ChangeSet, error) {
	if !fs.fs.Parsed() {
		_ = fs.fs.Parse(os.Args[1:])
		for len(fs.fs.Args()) != 0 {
			_ = fs.fs.Parse(fs.fs.Args()[1:])
		}
	}
	var changes map[string]interface{}

	visitFn := func(f *flag.Flag) {
		n := strings.ToLower(f.Name)
		keys := strings.FieldsFunc(n, split)
		reverse(keys)
		tmp := make(map[string]interface{})
		for i, k := range keys {
			if i == 0 {
				tmp[k] = f.Value
				continue
			}

			tmp = map[string]interface{}{k: tmp}
		}
		mergo.Map(&changes, tmp) // need to sort error handling
		return
	}
	unset, ok := fs.opts.Context.Value(includeUnsetKey{}).(bool)
	if ok && unset {
		fs.fs.VisitAll(visitFn)
	} else {
		fs.fs.Visit(visitFn)
	}

	b, err := fs.opts.Encoder.Encode(changes)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Priority:  source.PriorityFlag,
		Format:    fs.opts.Encoder.String(),
		Data:      b,
		Timestamp: time.Now(),
		Source:    fs.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func split(r rune) bool {
	return r == '-' || r == '_'
}

func reverse(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}

func (fs *flagsrc) Watch() (source.Watcher, error) {
	return nil, nil
}

func (fs *flagsrc) String() string {
	return "flag"
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)
	s := &flagsrc{opts: options, fs: flag.CommandLine}
	fs, ok := options.Context.Value(flagSetKey{}).(*flag.FlagSet)
	if ok {
		s.fs = fs
	}
	return s
}
