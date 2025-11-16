package logger

type Field struct {
	key    string
	valuee interface{}
}

type Logger interface {
	Info(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(err error, msg string, fields ...Field)

	With(fields ...Field) Logger
}