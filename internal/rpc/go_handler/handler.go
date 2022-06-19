package go_handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"tgram-slo-bot/internal"
	"tgram-slo-bot/internal/components/chat_storage"
	"tgram-slo-bot/internal/components/games_storage"
	"tgram-slo-bot/internal/components/poll_storage"
)

var (
	handlerName        = "go_handler"
	refreshText        = "🔄"
	incorrectGameError = fmt.Errorf("incorrect game chosen in handler: %s", handlerName)
)

type Handler struct {
	log          internal.Logger
	userStorage  chat_storage.Storage
	pollStorage  poll_storage.Storage
	gamesStorage games_storage.Storage
}

func New(logger internal.Logger, userStorage chat_storage.Storage, pollStorage poll_storage.Storage, gamesStorage games_storage.Storage) *Handler {
	return &Handler{
		log:          logger,
		userStorage:  userStorage,
		pollStorage:  pollStorage,
		gamesStorage: gamesStorage,
	}
}

func (h *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var (
		err         error
		messageText = "Выбери игру"
	)
	defer func() {
		ctx := h.log.WithFields(context.Background(), map[string]interface{}{
			"handler": handlerName,
		})
		if err != nil {
			h.log.Error(ctx, err)
		}
	}()

	keyboardRow := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData(refreshText, createOptionData(refreshText))},
	}
	games, err := h.gamesStorage.GetChatGames(update.FromChat().ID)
	if err != nil {
		return
	}

	for _, game := range games {
		keyboardRow = append(
			keyboardRow,
			[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(game, createOptionData(game))},
		)
	}

	optionKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		keyboardRow...,
	)

	if update.CallbackQuery == nil {
		message := tgbotapi.NewMessage(update.FromChat().ID, messageText)
		message.ReplyMarkup = &optionKeyboard
		_, err = bot.Send(message)
		return
	}

	editMessage := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, messageText)
	editMessage.ReplyMarkup = &optionKeyboard
	_, err = bot.Send(editMessage)
}

func (h *Handler) HandleChoice(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var (
		err    error
		chatID = update.FromChat().ID
	)
	defer func() {
		if err != nil {
			ctx := h.log.WithFields(context.Background(), map[string]interface{}{
				"handler": handlerName,
			})
			h.log.Error(ctx, err)
		}
	}()

	chatUsers, err := h.userStorage.GetChatUsers(chatID)
	if err != nil {
		return
	}

	chosenGame, err := getChosenGame(update.CallbackQuery.Data)
	if err != nil {
		return
	}

	if chosenGame == refreshText {
		h.Handle(update, bot)
		return
	}

	pollMessage := tgbotapi.SendPollConfig{
		IsAnonymous: false,
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatID,
		},
		Question: fmt.Sprintf("Идешь в %s?", chosenGame),
		Options: []string{
			"Да",
			"Нет",
		},
	}
	_, _ = bot.Send(tgbotapi.NewMessage(chatID, pollBuilder(chatUsers...)))
	msg, err := bot.Send(pollMessage)
	if err != nil {
		return
	}

	err = h.pollStorage.CreateNewPoll(chatID, msg.Poll.ID)

	editMessage := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("Пользователем @%s выбрана игра: %s", update.SentFrom().UserName, chosenGame))
	_, err = bot.Send(editMessage)
	if err != nil {
		return
	}
}

func (h *Handler) IsRightCommand(inputCmd string) bool {
	cmd := strings.Split(inputCmd, ":")
	if len(cmd) < 1 {
		return false
	}
	return cmd[0] == handlerName
}

func pollBuilder(opts ...*tgbotapi.User) string {
	var (
		sb         strings.Builder
		lineFormat = "@%s\n"
	)
	for _, v := range opts {
		sb.WriteString(fmt.Sprintf(lineFormat, v.UserName))
	}

	return sb.String()
}

func createOptionData(option string) string {
	return fmt.Sprintf("%s:%s", handlerName, option)
}

func getChosenGame(inputCmd string) (string, error) {
	game := strings.Split(inputCmd, ":")
	if len(game) < 2 || len(game[1]) < 1 {
		return "", incorrectGameError
	}
	return strings.Join(game[1:], ":"), nil
}
