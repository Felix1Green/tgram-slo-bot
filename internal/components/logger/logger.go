package logger

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func New() *Logger {
	return &Logger{
		log: logrus.New(),
	}
}

func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	l.log.Error(args...)
}

func (l *Logger) WithFields(ctx context.Context, fields map[string]interface{}) context.Context {
	ownCtx := context.WithValue(ctx, "", "")
	for k, v := range fields {
		ownCtx = context.WithValue(ownCtx, k, v)
	}

	return ownCtx
}
