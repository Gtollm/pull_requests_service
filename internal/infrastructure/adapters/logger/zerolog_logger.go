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

func attachFields(e *zerolog.Event, fields ...Field) *zerolog.Event {
	for _, v := range fields {
		*e = *e.Interface(v.key, v.valuee)
	}
	return e
}

func (z *ZerologLogger) Info(msg string, fields ...Field) {
	e := z.logger.Info()
	attachFields(e, fields...).Msg(msg)
}

func (z *ZerologLogger) Debug(msg string, fields ...Field) {
	e := z.logger.Debug()
	attachFields(e, fields...).Msg(msg)
}

func (z *ZerologLogger) Warn(msg string, fields ...Field) {
	e := z.logger.Warn()
	attachFields(e, fields...).Msg(msg)
}

func (z *ZerologLogger) Error(err error, msg string, fields ...Field) {
	e := z.logger.Error().Err(err)
	attachFields(e, fields...).Msg(msg)
}

func (z *ZerologLogger) With(fields ...Field) Logger {
	ctx := z.logger.With()
	for _, v := range fields {
		ctx = ctx.Interface(v.key, v.valuee)
	}

	return &ZerologLogger{ctx.Logger()}
}