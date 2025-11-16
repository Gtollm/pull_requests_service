package logger

type Field struct {
	key   string
	value interface{}
}

func F(key string, value interface{}) Field {
	return Field{key: key, value: value}
}

type Logger interface {
	Info(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(err error, msg string, fields ...Field)

	With(fields ...Field) Logger
}