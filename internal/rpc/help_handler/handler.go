package help_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgram-slo-bot/internal"
)

var (
	helpMessage = "Чтобы принять участие в голосовании необходимо зарегистрироваться с помощью команды /reg"
)

type Handler struct {
	log internal.Logger
}

func New(log internal.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func (h *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, helpMessage)
	_, _ = bot.Send(msg)
}
