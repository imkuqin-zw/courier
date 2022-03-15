package xmap

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func Marshal(obj interface{}) map[string]interface{} {
	ot := reflect.TypeOf(obj)
	ov := reflect.ValueOf(obj)

	var out = make(map[string]interface{}, ot.NumField())
	for i := 0; i < ot.NumField(); i++ {
		out[ot.Field(i).Name] = ov.Field(i).Interface()
	}
	return out
}

func Unmarshal(in interface{}, out interface{}) error {
	return errors.WithStack(mapstructure.Decode(in, out))
}
