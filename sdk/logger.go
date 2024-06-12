package investgo

// Logger interface for structural logging
type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}

// NoopLogger implements Logger interface but does nothing
type NoopLogger struct{}

func (NoopLogger) Debug(msg string, keyvals ...interface{}) {}
func (NoopLogger) Info(msg string, keyvals ...interface{})  {}
func (NoopLogger) Warn(msg string, keyvals ...interface{})  {}
func (NoopLogger) Error(msg string, keyvals ...interface{}) {}
