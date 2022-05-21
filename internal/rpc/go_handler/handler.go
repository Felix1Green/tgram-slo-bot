package go_handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"tgram-slo-bot/internal"
	"tgram-slo-bot/internal/components/chat_storage"
	"tgram-slo-bot/internal/components/poll_storage"
)

var (
	handlerName = "go_handler"
)

type Handler struct {
	log         internal.Logger
	userStorage chat_storage.Storage
	pollStorage poll_storage.Storage
}

func New(logger internal.Logger) *Handler {
	return &Handler{
		log: logger,
	}
}

func (t *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var err error
	defer func() {
		if err != nil {
			ctx := t.log.WithFields(context.Background(), map[string]interface{}{
				"handler": handlerName,
			})
			t.log.Error(ctx, err)
		}
	}()

	chatUsers, err := t.userStorage.GetChatUsers(update.FromChat().ID)
	if err != nil {
		return
	}

	pollMessage := tgbotapi.SendPollConfig{
		IsAnonymous: false,
		BaseChat: tgbotapi.BaseChat{
			ChatID: update.FromChat().ID,
		},
		Question: questionBuilder("Идешь?", chatUsers...),
		Options: []string{
			"Да",
			"Нет",
		},
	}

	_, err = bot.Send(pollMessage)
	if err != nil {
		return
	}

	//t.pollStorage.
}

func questionBuilder(questionString string, opts ...*tgbotapi.User) string {
	var (
		sb         strings.Builder
		lineFormat = "%s\n"
	)
	sb.WriteString(fmt.Sprintf(lineFormat, questionString))
	for _, v := range opts {
		sb.WriteString(fmt.Sprintf(lineFormat, v.UserName))
	}

	return sb.String()
}
