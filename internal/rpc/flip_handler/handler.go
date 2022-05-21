package flip_handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"math/rand"
	"tgram-slo-bot/internal"
)

var (
	handlerName = "flip_handler"
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

	var flipResult string
	flipIntegerResult := rand.Int() % 2
	switch flipIntegerResult {
	case 0:
		flipResult = "РЕШКА"
	case 1:
		flipResult = "ОРЕЛ"
	}

	message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Результат: %s", flipResult))
	_, err = bot.Send(message)
}
