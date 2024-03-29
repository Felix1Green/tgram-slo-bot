package poll_update_handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgram-slo-bot/internal"
	"tgram-slo-bot/internal/components/poll_storage"
)

const (
	handlerName = "poll_update_handler"
)

type Handler struct {
	log         internal.Logger
	pollStorage poll_storage.Storage
}

func New(log internal.Logger, storage poll_storage.Storage) *Handler {
	return &Handler{
		log:         log,
		pollStorage: storage,
	}
}

func (s *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	answer := update.PollAnswer
	optionId := answer.OptionIDs[0]
	err := s.pollStorage.SetUserVoted(answer.PollID, answer.User.ID)
	if (optionId == 0 || optionId == 2) && err == nil {
		err = s.pollStorage.SetUserRegistered(answer.PollID, answer.User.ID)
	}
	if err != nil {
		ctx := s.log.WithFields(context.Background(), map[string]interface{}{
			"handlerName": handlerName,
		})
		s.log.Error(ctx, err)
	}
}
