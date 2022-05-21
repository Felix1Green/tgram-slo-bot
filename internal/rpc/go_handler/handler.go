package go_handler

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"tgram-slo-bot/internal"
)

type Handler struct {
	log internal.Logger
}

func New(logger internal.Logger) *Handler {
	return &Handler{
		log: logger,
	}
}

func (t *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	//pollMes
}
