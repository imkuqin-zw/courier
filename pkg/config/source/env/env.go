package env

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/imkuqin-zw/courier/pkg/config/source"
)

type env struct {
	prefixes         []string
	strippedPrefixes []string
	separator        string
	replaceRule      map[string]map[string]string
	opts             source.Options
}

func (e *env) Read() (*source.ChangeSet, error) {
	var changes map[string]interface{}
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		value := pair[1]
		key := strings.ToLower(pair[0])
		if len(e.prefixes) > 0 || len(e.strippedPrefixes) > 0 {
			notFound := true

			if _, ok := matchPrefix(e.prefixes, key); ok {
				notFound = false
			}
			if match, ok := matchPrefix(e.strippedPrefixes, key); ok {
				if key = strings.TrimPrefix(key, match); key != "" {
					notFound = false
				}
			}
			if notFound {
				continue
			}

		}
		keys := strings.Split(key, e.separator)
		e.replaceKey(key, keys)
		reverse(keys)
		tmp := make(map[string]interface{})
		for i, k := range keys {
			if i == 0 {
				if intValue, err := strconv.Atoi(value); err == nil {
					tmp[k] = intValue
				} else if boolValue, err := strconv.ParseBool(value); err == nil {
					tmp[k] = boolValue
				} else {
					tmp[k] = value
				}
				continue
			}

			tmp = map[string]interface{}{k: tmp}
		}

		if err := mergo.Map(&changes, tmp); err != nil {
			return nil, err
		}
	}

	b, err := e.opts.Encoder.Encode(changes)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Priority:  source.PriorityEnv,
		Format:    e.opts.Encoder.String(),
		Data:      b,
		Timestamp: time.Now(),
		Source:    e.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func (e *env) replaceKey(key string, keys []string) {
	i := strings.Index(key, e.separator)
	if i < 0 {
		return
	}
	if replace, ok := e.replaceRule[key[:i]]; ok {
		for oldStr, newStr := range replace {
			for _, item := range keys {
				item = strings.ReplaceAll(item, oldStr, newStr)
			}
		}
	}
}

func matchPrefix(pre []string, s string) (string, bool) {
	for _, p := range pre {
		if strings.HasPrefix(s, p) {
			return p, true
		}
	}

	return "", false
}

func reverse(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}

func (e *env) Watch() (source.Watcher, error) {
	return nil, nil
}

func (e *env) String() string {
	return "env"
}

// NewSource returns a conf source for parsing ENV variables.
// Underscores are delimiters for nesting, and all randkeys are lowercased.
//
// Example:
//      "DATABASE_SERVER_HOST=localhost" will convert to
//
//      {
//          "database": {
//              "pub": {
//                  "host": "localhost"
//              }
//          }
//      }
func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)

	var sp []string
	var pre []string
	var rr map[string]map[string]string
	var separator = "__"
	if p, ok := options.Context.Value(separatorKey{}).(string); ok {
		separator = p
	}
	if p, ok := options.Context.Value(strippedPrefixKey{}).([]string); ok {
		sp = appendUnderscore(p, separator)
	}

	if p, ok := options.Context.Value(prefixKey{}).([]string); ok {
		pre = appendUnderscore(p, separator)
	}

	if p, ok := options.Context.Value(keyReplaceRule{}).(map[string]map[string]string); ok {
		rr = p
	} else {
		rr = make(map[string]map[string]string)
	}
	return &env{prefixes: pre, strippedPrefixes: sp, replaceRule: rr, separator: separator, opts: options}
}
