package logger

import (
	"github.com/rs/zerolog"
	"os"
)

type ZerologLogger struct {
	logger zerolog.Logger
}

func NewZerologLogger() Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &ZerologLogger{logger}
}

func attachFields(e *zerolog.Event, fields Fields) *zerolog.Event {
	for k, v := range fields {
		*e = *e.Interface(k, v)
	}
	return e
}

func (z *ZerologLogger) Info(msg string, fields Fields) {
	e := z.logger.Info()
	attachFields(e, fields).Msg(msg)
}

func (z *ZerologLogger) Debug(msg string, fields Fields) {
	e := z.logger.Debug()
	attachFields(e, fields).Msg(msg)
}

func (z *ZerologLogger) Warn(msg string, fields Fields) {
	e := z.logger.Warn()
	attachFields(e, fields).Msg(msg)
}

func (z *ZerologLogger) Error(err error, msg string, fields Fields) {
	e := z.logger.Error().Err(err)
	attachFields(e, fields).Msg(msg)
}

func (z *ZerologLogger) With(fields Fields) Logger {
	ctx := z.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}

	return &ZerologLogger{ctx.Logger()}
}