package ping_handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"tgram-slo-bot/internal"
	"tgram-slo-bot/internal/components/chat_storage"
	"tgram-slo-bot/internal/components/poll_storage"
)

const (
	handlerName = "ping_handler"
)

type Handler struct {
	log         internal.Logger
	pollStorage poll_storage.Storage
	chatStorage chat_storage.Storage
}

func New(log internal.Logger, storage poll_storage.Storage, chatStorage chat_storage.Storage) *Handler {
	return &Handler{
		log:         log,
		pollStorage: storage,
		chatStorage: chatStorage,
	}
}

func (s *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatId := update.FromChat().ID
	users, err := s.chatStorage.GetChatUsers(chatId)
	defer func() {
		if err != nil {
			ctx := s.log.WithFields(context.Background(), map[string]interface{}{
				"handlerName": handlerName,
			})
			s.log.Error(ctx, err)
		}
	}()

	pollId, err := s.pollStorage.GetActiveChatPoll(chatId)
	if err != nil {
		return
	}

	registeredUsers := make([]*tgbotapi.User, 0)
	for _, user := range users {
		voted, err := s.pollStorage.IsCurrentUserRegistered(pollId, user.ID)
		if err != nil {
			return
		}
		if voted {
			registeredUsers = append(registeredUsers, user)
		}
	}
	pingMessage := ""
	updateText := strings.Split(update.Message.Text, " ")
	if len(updateText) >= 2 || len(updateText[1]) >= 1 {
		pingMessage = strings.Join(updateText[1:], " ")
	}

	msg := tgbotapi.NewMessage(update.FromChat().ID, pingMessage+pingBuilder(registeredUsers...))
	_, err = bot.Send(msg)
}

func pingBuilder(opts ...*tgbotapi.User) string {
	var (
		sb         strings.Builder
		lineFormat = "@%s\n"
	)
	for _, v := range opts {
		sb.WriteString(fmt.Sprintf(lineFormat, v.UserName))
	}

	return sb.String()
}
