package dubbo

type Options struct {
	Format string
}

func defaultOpts() *Options {
	return &Options{Format: "yaml"}
}

type Option func(o *Options)

func WithFormat(format string) Option {
	return func(o *Options) {
		o.Format = format
	}
}
