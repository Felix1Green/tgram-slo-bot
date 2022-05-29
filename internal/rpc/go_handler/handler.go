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

func New(logger internal.Logger, userStorage chat_storage.Storage, pollStorage poll_storage.Storage) *Handler {
	return &Handler{
		log:         logger,
		userStorage: userStorage,
		pollStorage: pollStorage,
	}
}

func (t *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var (
		err    error
		chatID = update.FromChat().ID
	)
	defer func() {
		if err != nil {
			ctx := t.log.WithFields(context.Background(), map[string]interface{}{
				"handler": handlerName,
			})
			t.log.Error(ctx, err)
		}
	}()

	chatUsers, err := t.userStorage.GetChatUsers(chatID)
	if err != nil {
		return
	}

	pollMessage := tgbotapi.SendPollConfig{
		IsAnonymous: false,
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatID,
		},
		Question:    "Идешь?",
		Explanation: questionBuilder(chatUsers...),
		Options: []string{
			"Да",
			"Нет",
		},
	}

	msg, err := bot.Send(pollMessage)
	if err != nil {
		return
	}

	err = t.pollStorage.CreateNewPoll(chatID, msg.Poll.ID)
}

func questionBuilder(opts ...*tgbotapi.User) string {
	var (
		sb         strings.Builder
		lineFormat = "@%s\n"
	)
	for _, v := range opts {
		sb.WriteString(fmt.Sprintf(lineFormat, v.UserName))
	}

	return sb.String()
}
