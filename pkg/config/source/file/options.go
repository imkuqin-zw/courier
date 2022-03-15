package file

import (
	"context"

	"github.com/imkuqin-zw/courier/pkg/config/source"
)

type filePathKey struct{}
type watchKey struct{}

// WithPath sets the path to file
func WithPath(p string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, filePathKey{}, p)
	}
}

// WithWatch sets whether to watch file
func WithWatch(b bool) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, watchKey{}, b)
	}
}
