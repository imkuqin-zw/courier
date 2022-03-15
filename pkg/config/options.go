package config

import (
	"github.com/imkuqin-zw/courier/pkg/config/loader"
	"github.com/imkuqin-zw/courier/pkg/config/reader"
	"github.com/imkuqin-zw/courier/pkg/config/source"
)

type Options struct {
	Loader loader.Loader
	Reader reader.Reader
	Source []source.Source
}

type Option func(o *Options)

// WithLoader sets the loader for manager conf
func WithLoader(l loader.Loader) Option {
	return func(o *Options) {
		o.Loader = l
	}
}

// WithSource appends a source to list of sources
func WithSource(s ...source.Source) Option {
	return func(o *Options) {
		o.Source = append(o.Source, s...)
	}
}

// WithReader sets the conf reader
func WithReader(r reader.Reader) Option {
	return func(o *Options) {
		o.Reader = r
	}
}
