package flag

import (
	"context"
	"flag"

	"github.com/imkuqin-zw/courier/pkg/config/source"
)

type includeUnsetKey struct{}
type flagSetKey struct{}

// IncludeUnset toggles the loading of unset flags and their respective default values.
// Default behavior is to ignore any unset flags.
func IncludeUnset(b bool) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, includeUnsetKey{}, true)
	}
}

// IncludeUnset toggles the loading of unset flags and their respective default values.
// Default behavior is to ignore any unset flags.
func WithFlagSet(fs *flag.FlagSet) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, flagSetKey{}, fs)
	}
}
