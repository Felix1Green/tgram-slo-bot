package rpc

import (
	"github.com/Syfaro/telegram-bot-api"
	"tgram-slo-bot/internal"
	"tgram-slo-bot/internal/rpc/go_handler"
)

type HandlerComposer struct {
	bot          *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig
	goHandler    *go_handler.Handler
}

type specifications struct {
	TelegramToken string `split_words:"true"`
}

func NewFromEnv() (*HandlerComposer, error) {
	options := &specifications{}
	err := internal.EnvOptions("", options)
	if err != nil {
		return nil, err
	}
	bot, err := tgbotapi.NewBotAPI(options.TelegramToken)
	if err != nil {
		return nil, err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return &HandlerComposer{
		bot:          bot,
		updateConfig: u,
	}, nil
}

func (t *HandlerComposer) Listen() error {
	updates, err := t.bot.GetUpdatesChan(t.updateConfig)
	if err != nil {
		return err
	}
	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}
		switch update.Message.Text {
		case "/go":
			t.goHandler.Handle(&update)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I cant do this shit")
			// навесить backoff на отправку сообщения
			_, _ = t.bot.Send(msg)
		}
	}

	return nil
}
