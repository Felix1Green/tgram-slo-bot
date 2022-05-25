package monitoring

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
	componentName = "monitoring_worker"
)

type Worker struct {
	log         internal.Logger
	pollStorage poll_storage.Storage
	chatStorage chat_storage.Storage
	bot         *tgbotapi.BotAPI
}

func New(log internal.Logger, poll poll_storage.Storage, chat chat_storage.Storage) *Worker {
	return &Worker{
		log:         log,
		pollStorage: poll,
		chatStorage: chat,
	}
}

func (w *Worker) Run() {
	var (
		err error
	)
	defer func() {
		if err != nil {
			ctx := w.log.WithFields(context.Background(), map[string]interface{}{
				"componentName": componentName,
			})
			w.log.Error(ctx, err)
		}
	}()

	polls, err := w.pollStorage.GetActivePollKeys()
	if err != nil {
		return
	}

	for _, key := range polls {
		unvotedUsers := make([]string, 0)
		poll, err := w.pollStorage.GetPollInfoFromKey(key)
		if err != nil {
			return
		}

		chatUsers, err := w.chatStorage.GetChatUsers(poll.ChatID)
		if err != nil {
			return
		}
		// TODO: get all not voted users from one request
		for _, user := range chatUsers {
			voted, err := w.pollStorage.IsCurrentUserVoted(poll.ChatID, poll.PollID, user.ID)
			if err != nil {
				return
			}

			if !voted {
				unvotedUsers = append(unvotedUsers, user.UserName)
			}
		}

		_ = w.pollStorage.RemovePoll(key)
		if len(unvotedUsers) > 0 {
			message := tgbotapi.NewMessage(poll.ChatID, w.BuildNotifyMessage(unvotedUsers))
			_, _ = w.bot.Send(message)
		}
	}
}

func (w *Worker) BuildNotifyMessage(users []string) string {
	var (
		sb         strings.Builder
		lineFormat = "%s\n"
	)
	sb.WriteString(fmt.Sprintf(lineFormat, "Пора ответить за игнорирование!"))
	for _, v := range users {
		sb.WriteString(fmt.Sprintf(lineFormat, v))
	}

	return sb.String()
}
