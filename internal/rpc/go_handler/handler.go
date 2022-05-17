package go_handler

import tgbotapi "github.com/Syfaro/telegram-bot-api"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (t *Handler) Handle(update *tgbotapi.Update) {
	// some logic stuff
}