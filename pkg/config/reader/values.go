package reader

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	simple "github.com/bitly/go-simplejson"
	"github.com/imkuqin-zw/courier/pkg/config/source"
	"github.com/mitchellh/mapstructure"
)

type jsonValues struct {
	ch *source.ChangeSet
	sj *simple.Json
}

type jsonValue struct {
	*simple.Json
}

func newValues(ch *source.ChangeSet) (Values, error) {
	sj := simple.New()
	data, _ := ReplaceEnvVars(ch.Data)
	if err := sj.UnmarshalJSON(data); err != nil {
		sj.SetPath(nil, string(ch.Data))
	}
	return &jsonValues{ch, sj}, nil
}

func (j *jsonValues) Get(key string) Value {
	if key == "" {
		return &jsonValue{j.sj.GetPath()}
	}
	return &jsonValue{j.sj.GetPath(strings.Split(key, Separator)...)}
}

func (j *jsonValues) Del(key string) {
	if len(key) == 0 {
		j.sj = simple.New()
		return
	}
	path := strings.Split(key, Separator)
	if len(path) == 1 {
		j.sj.Del(path[0])
		return
	}

	vals := j.sj.GetPath(path[:len(path)-1]...)
	vals.Del(path[len(path)-1])
	j.sj.SetPath(path[:len(path)-1], vals.Interface())
	return
}

func (j *jsonValues) Set(key string, val interface{}) {
	if len(key) == 0 {
		j.sj.SetPath([]string{}, val)
		return
	}
	j.sj.SetPath(strings.Split(key, Separator), val)
}

func (j *jsonValues) Bytes() []byte {
	b, _ := j.sj.MarshalJSON()
	return b
}

func (j *jsonValues) Map() map[string]interface{} {
	m, _ := j.sj.Map()
	return m
}

func (j *jsonValues) Scan(v interface{}) error {
	b, err := j.sj.Map()
	if err != nil {
		return err
	}
	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     v,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	return decoder.Decode(b)
}

func (j *jsonValues) String() string {
	return "json"
}

func (j *jsonValue) Bool(def ...bool) bool {
	b, err := j.Json.Bool()
	if err == nil {
		return b
	}

	str, ok := j.Interface().(string)
	if !ok {
		if len(def) == 0 {
			return false
		}
		return def[0]
	}

	b, err = strconv.ParseBool(str)
	if err != nil {
		if len(def) == 0 {
			return false
		}
		return def[0]
	}

	return b
}

func (j *jsonValue) Int(def ...int) int {
	i, err := j.Json.Int()
	if err == nil {
		return i
	}

	str, ok := j.Interface().(string)
	if !ok {
		if len(def) == 0 {
			return 0
		}
		return def[0]
	}

	i, err = strconv.Atoi(str)
	if err != nil {
		if len(def) == 0 {
			return 0
		}
		return def[0]
	}

	return i
}

func (j *jsonValue) String(def ...string) string {
	return j.Json.MustString(def...)
}

func (j *jsonValue) Float64(def ...float64) float64 {
	f, err := j.Json.Float64()
	if err == nil {
		return f
	}

	str, ok := j.Interface().(string)
	if !ok {
		if len(def) == 0 {
			return 0
		}
		return def[0]
	}

	f, err = strconv.ParseFloat(str, 64)
	if err != nil {
		if len(def) == 0 {
			return 0
		}
		return def[0]
	}

	return f
}

func (j *jsonValue) Duration(def ...time.Duration) time.Duration {
	v, err := j.Json.String()
	if err != nil {
		if len(def) == 0 {
			return 0
		}
		return def[0]
	}

	value, err := time.ParseDuration(v)
	if err != nil {
		if len(def) == 0 {
			return 0
		}
		return def[0]
	}

	return value
}

func (j *jsonValue) StringSlice(def ...string) []string {
	v, err := j.Json.String()
	if err == nil {
		sl := strings.Split(v, ",")
		if len(sl) > 1 {
			return sl
		}
	}
	return j.Json.MustStringArray(def)
}

func (j *jsonValue) StringMap(def ...map[string]string) map[string]string {
	res := map[string]string{}
	m, err := j.Json.Map()
	if err != nil {
		if len(def) == 0 {
			return res
		}
		return def[0]
	}

	for k, v := range m {
		res[k] = fmt.Sprintf("%v", v)
	}

	return res
}

func (j *jsonValue) Scan(v interface{}) error {
	if j.Interface() == nil {
		return nil
	}
	b, err := j.Json.Map()
	if err != nil {
		return err
	}
	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     v,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	return decoder.Decode(b)
}

func (j *jsonValue) Bytes() []byte {
	b, err := j.Json.Bytes()
	if err != nil {
		// try return marshalled
		b, err = j.Json.MarshalJSON()
		if err != nil {
			return []byte{}
		}
		return b
	}
	return b
}

func ReplaceEnvVars(raw []byte) ([]byte, error) {
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_]+)\}`)
	if re.Match(raw) {
		dataS := string(raw)
		res := re.ReplaceAllStringFunc(dataS, replaceEnvVars)
		return []byte(res), nil
	} else {
		return raw, nil
	}
}

func replaceEnvVars(element string) string {
	v := element[2 : len(element)-1]
	el := os.Getenv(v)
	return el
}
