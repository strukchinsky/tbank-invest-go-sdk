package investgo

// Logger interface for structural logging
type Logger interface {
	Debug(msg string, keyvals ...any)
	Info(msg string, keyvals ...any)
	Warn(msg string, keyvals ...any)
	Error(msg string, keyvals ...any)
}

// NoopLogger implements Logger interface but does nothing
type NoopLogger struct{}

func (NoopLogger) Debug(msg string, keyvals ...any) {}
func (NoopLogger) Info(msg string, keyvals ...any)  {}
func (NoopLogger) Warn(msg string, keyvals ...any)  {}
func (NoopLogger) Error(msg string, keyvals ...any) {}
