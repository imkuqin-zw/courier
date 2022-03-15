package memory

import (
	"github.com/imkuqin-zw/courier/pkg/config/loader"
	"github.com/imkuqin-zw/courier/pkg/config/reader"
	"github.com/imkuqin-zw/courier/pkg/config/source"
)

// WithSource appends a source to list of sources
func WithSource(s source.Source) loader.Option {
	return func(o *loader.Options) {
		o.Source = append(o.Source, s)
	}
}

// WithReader sets the conf reader
func WithReader(r reader.Reader) loader.Option {
	return func(o *loader.Options) {
		o.Reader = r
	}
}
