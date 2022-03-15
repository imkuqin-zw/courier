package reader

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/imkuqin-zw/courier/pkg/config/source"
)

func TestValues(t *testing.T) {
	emptyStr := ""
	testData := []struct {
		csdata   []byte
		path     []string
		accepter interface{}
		value    interface{}
	}{
		{
			[]byte(`{"foo": "bar", "baz": {"bar": "cat"}}`),
			[]string{"foo"},
			emptyStr,
			"bar",
		},
		{
			[]byte(`{"foo": "bar", "baz": {"bar": "cat"}}`),
			[]string{"baz", "bar"},
			emptyStr,
			"cat",
		},
	}

	for idx, test := range testData {
		values, err := newValues(&source.ChangeSet{
			Data: test.csdata,
		})
		if err != nil {
			t.Fatal(err)
		}

		err = values.Get(strings.Join(test.path, ".")).Scan(&test.accepter)
		if err != nil {
			t.Fatal(err)
		}
		if test.accepter != test.value {
			t.Fatalf("No.%d Expected %v got %v for path %v", idx, test.value, test.accepter, test.path)
		}
	}
}

func TestStructArray(t *testing.T) {
	type T struct {
		Foo string
	}

	emptyTSlice := []T{}

	testData := []struct {
		csdata   []byte
		accepter []T
		value    []T
	}{
		{
			[]byte(`[{"foo": "bar"}]`),
			emptyTSlice,
			[]T{{Foo: "bar"}},
		},
	}

	for idx, test := range testData {
		values, err := newValues(&source.ChangeSet{
			Data: test.csdata,
		})
		if err != nil {
			t.Fatal(err)
		}

		err = values.Get("").Scan(&test.accepter)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(test.accepter, test.value) {
			t.Fatalf("No.%d Expected %v got %v", idx, test.value, test.accepter)
		}
	}
}

func TestReplaceEnvVars(t *testing.T) {
	os.Setenv("myBar", "cat")
	os.Setenv("MYBAR", "cat")
	os.Setenv("my_Bar", "cat")
	os.Setenv("myBar_", "cat")

	testData := []struct {
		expected string
		data     []byte
	}{
		// Right use cases
		{
			`{"foo": "bar", "baz": {"bar": "cat"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${myBar}"}}`),
		},
		{
			`{"foo": "bar", "baz": {"bar": "cat"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${MYBAR}"}}`),
		},
		{
			`{"foo": "bar", "baz": {"bar": "cat"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${my_Bar}"}}`),
		},
		{
			`{"foo": "bar", "baz": {"bar": "cat"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${myBar_}"}}`),
		},
		// Wrong use cases
		{
			`{"foo": "bar", "baz": {"bar": "${myBar-}"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${myBar-}"}}`),
		},
		{
			`{"foo": "bar", "baz": {"bar": "${}"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${}"}}`),
		},
		{
			`{"foo": "bar", "baz": {"bar": "$sss}"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "$sss}"}}`),
		},
		{
			`{"foo": "bar", "baz": {"bar": "${sss"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${sss"}}`),
		},
		{
			`{"foo": "bar", "baz": {"bar": "{something}"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "{something}"}}`),
		},
		// Use cases without replace env vars
		{
			`{"foo": "bar", "baz": {"bar": "cat"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "cat"}}`),
		},
	}

	for _, test := range testData {
		res, err := ReplaceEnvVars(test.data)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Compare(test.expected, string(res)) != 0 {
			t.Fatalf("Expected %s got %s", test.expected, res)
		}
	}
}
