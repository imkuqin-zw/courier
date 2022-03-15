package gopool

var defaultLogger Logger

func init() {
	defaultLogger = &logger{}
}

func SetLogger(lg Logger) {
	if lg == nil {
		return
	}
	defaultLogger = lg
}

// Logger is a logger interface that provides logging function with levels.
type Logger interface {
	Errorf(format string, v ...interface{})
}

type logger struct{}

func (l *logger) Errorf(format string, v ...interface{}) {
	return
}
