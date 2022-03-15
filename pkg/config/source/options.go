package source

import (
	"context"

	"github.com/imkuqin-zw/courier/pkg/config/encoder"
	"github.com/imkuqin-zw/courier/pkg/config/encoder/json"
)

const (
	PriorityKey = "_priority"
)

type Priority uint8

const (
	PriorityMin Priority = iota

	PriorityMemory
	PriorityFile
	PriorityEnv
	PriorityFlag
	PriorityCli
	PriorityRemote

	PriorityMax
)

type Options struct {
	// Encoder
	Encoder encoder.Encoder

	// for alternative repo
	Context context.Context
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encoder: json.NewEncoder(),
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithEncoder sets the source encoder
func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		o.Encoder = e
	}
}
