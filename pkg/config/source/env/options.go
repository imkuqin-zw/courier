package env

import (
	"context"

	"strings"

	"github.com/imkuqin-zw/courier/pkg/config/source"
)

type strippedPrefixKey struct{}
type prefixKey struct{}
type keyReplaceRule struct{}
type separatorKey struct{}

// WithStrippedPrefix sets the environment variable prefixes to scope to.
// These prefixes will be removed from the actual conf entries.
func WithStrippedPrefix(p ...string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, strippedPrefixKey{}, p)
	}
}

// WithPrefix sets the environment variable prefixes to scope to.
// These prefixes will not be removed. Each prefix will be considered a top level conf entry.
func WithPrefix(p ...string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, prefixKey{}, p)
	}
}

func WithKeyReplaceRule(rule map[string]map[string]string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		replaceMap, ok := o.Context.Value(keyReplaceRule{}).(map[string]map[string]string)
		if ok {
			for key, replace := range rule {
				if len(replace) == 0 {
					delete(replaceMap, key)
				} else {
					if r, ok := replaceMap[key]; ok {
						for k, v := range replace {
							if len(v) == 0 {
								delete(r, k)
							} else {
								r[k] = v
							}
						}
						if len(r) > 0 {
							replaceMap[key] = r
						} else {
							delete(replaceMap, key)
						}
					} else {
						r := make(map[string]string)
						for k, v := range replace {
							if len(v) != 0 {
								r[k] = v
							}
						}
						if len(r) > 0 {
							replaceMap[key] = r
						}
					}
				}
			}
		} else {
			replaceMap := make(map[string]map[string]string)
			for key, replace := range rule {
				if len(replace) == 0 {
					delete(rule, key)
					continue
				}
				r := make(map[string]string)
				for k, v := range replace {
					if len(v) != 0 {
						r[k] = v
					}
				}
				if len(r) > 0 {
					replaceMap[key] = r
				}
			}
			o.Context = context.WithValue(o.Context, keyReplaceRule{}, replaceMap)
		}
	}
}

func WithSeparator(s string) source.Option {
	return func(o *source.Options) {
		if s == "" {
			return
		}
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, separatorKey{}, s)
	}
}

func appendUnderscore(prefixes []string, separator string) []string {
	//nolint:prealloc
	var result []string
	for _, p := range prefixes {
		if !strings.HasSuffix(p, "separator") {
			result = append(result, strings.ToLower(p)+separator)
			continue
		}

		result = append(result, strings.ToLower(p))
	}

	return result
}
