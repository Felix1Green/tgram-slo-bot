package flip_handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"tgram-slo-bot/internal"
)

var (
	handlerName = "flip_handler"
)

type Handler struct {
	log internal.Logger
}

var (
	eagle           = "ОРЕЛ"
	tails           = "РЕШКА"
	numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(eagle, eagle),
			tgbotapi.NewInlineKeyboardButtonData(tails, tails),
		),
	)
)

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
	message := tgbotapi.NewMessage(update.FromChat().ID, "Выбери предполагаемый результат")
	message.ReplyMarkup = numericKeyboard
	_, err = bot.Send(message)
}

func (h *Handler) HandleChoice(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var err error
	defer func() {
		ctx := h.log.WithFields(context.Background(), map[string]interface{}{
			"handler": handlerName,
		})
		if err != nil {
			h.log.Error(ctx, err)
		}
	}()

	var (
		flipResult       string
		winner           = "Выиграл"
		looser           = "Проиграл"
		CorrectlyGuessed = looser
	)
	flipIntegerResult := rand.Int() % 2
	switch flipIntegerResult {
	case 0:
		flipResult = tails
	case 1:
		flipResult = eagle
	}

	fmt.Println(update.CallbackQuery.Data)
	if flipResult == update.CallbackQuery.Data {
		CorrectlyGuessed = winner
	}

	message := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("@%s %s\nРезультат: %s", update.SentFrom().UserName, CorrectlyGuessed, flipResult))
	_, _ = bot.Send(message)
}
