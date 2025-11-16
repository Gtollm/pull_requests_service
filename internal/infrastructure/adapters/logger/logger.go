package logger

type Fields map[string]interface{}

type Logger interface {
	Info(msg string, fields Fields)
	Debug(msg string, fields Fields)
	Warn(msg string, fields Fields)
	Error(err error, msg string, fields Fields)

	With(fields Fields) Logger
}