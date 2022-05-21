package internal

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Logger interface {
	Error(ctx context.Context, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
	WithError(ctx context.Context, err error) context.Context
	WithFields(ctx context.Context, fields map[string]interface{}) context.Context
}

type HandleFunc func(update *tgbotapi.Update, bot *tgbotapi.BotAPI)
