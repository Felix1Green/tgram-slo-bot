package games_handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"tgram-slo-bot/internal"
	"tgram-slo-bot/internal/components/games_storage"
)

var (
	helpMessage       = "Чтобы добавить игру в голосование, необходимо после команды /add написать название игры"
	incorrectGameName = fmt.Errorf("incorrect game name")
)

type Handler struct {
	log     internal.Logger
	storage games_storage.Storage
}

func New(log internal.Logger, storage games_storage.Storage) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var err error
	defer func() {
		if err != nil {
			msg := tgbotapi.NewMessage(update.FromChat().ID, helpMessage)
			_, _ = bot.Send(msg)
		}
	}()
	updateText := strings.Split(update.Message.Text, " ")
	if len(updateText) < 2 || len(updateText[1]) < 1 {
		err = incorrectGameName
		return
	}

	gameName := strings.Join(updateText[1:], " ")
	err = h.storage.CreateNewChatGame(update.FromChat().ID, gameName)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Игра %s успешно добавлена", gameName))
	_, _ = bot.Send(msg)
}
