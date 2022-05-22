package reg_handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgram-slo-bot/internal"
	"tgram-slo-bot/internal/components/chat_storage"
)

var (
	handlerName = "reg_handler"
)

type Handler struct {
	log         internal.Logger
	userStorage chat_storage.Storage
}

func New(logger internal.Logger, storage chat_storage.Storage) *Handler {
	return &Handler{
		log:         logger,
		userStorage: storage,
	}
}

func (t *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var (
		chatID = update.FromChat().ID
		err    error
	)
	defer func() {
		if err != nil {
			ctx := t.log.WithFields(context.Background(), map[string]interface{}{
				"handlerName": handlerName,
			})

			t.log.Error(ctx, err)
		}
	}()

	requestedUser := update.SentFrom()
	err = t.userStorage.RegistrateNewUser(chatID, requestedUser)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Пользователь %s успешно зарегистрирован", requestedUser.UserName))
	_, err = bot.Send(msg)
}
