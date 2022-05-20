package random_handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"math/rand"
	"tgram-slo-bot/internal"
)

var (
	handlerName = "random_handler"
)

type Handler struct {
	log internal.Logger
}

func New(logger internal.Logger) *Handler {
	return &Handler{
		log: logger,
	}
}

func (h *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var err error
	defer func() {
		ctx := h.log.WithFields(context.Background(), map[string]interface{}{
			"handler": handlerName,
		})
		if err != nil {
			h.log.Error(ctx, err)
		}
	}()
	_, max := h.getRequestBounds(update)
	randomResult := rand.Intn(max)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Случайное число: %d", randomResult))
	_, err = bot.Send(msg)
}

func (h *Handler) getRequestBounds(update *tgbotapi.Update) (min int, max int) {
	min, max = 0, 10
	// TODO: fix this handler func
	return
}
