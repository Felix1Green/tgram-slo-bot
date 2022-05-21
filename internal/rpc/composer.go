package rpc

import (
	"github.com/Syfaro/telegram-bot-api"
	"tgram-slo-bot/internal"
)

type HandlerComposer struct {
	log           internal.Logger
	bot           *tgbotapi.BotAPI
	updateConfig  tgbotapi.UpdateConfig
	handlerConfig map[string]internal.HandleFunc
}

type specifications struct {
	TelegramToken string `split_words:"true"`
}

func NewFromEnv(logger internal.Logger, config map[string]internal.HandleFunc) (*HandlerComposer, error) {
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
		log:           logger,
		handlerConfig: config,
		bot:           bot,
		updateConfig:  u,
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
		if handler, ok := t.handlerConfig[update.Message.Text]; ok {
			handler(&update, t.bot)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I cant do this")
			_, _ = t.bot.Send(msg)
		}
	}

	return nil
}
