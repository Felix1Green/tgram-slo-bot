package rpc

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgram-slo-bot/internal"
)

type HandlerComposer struct {
	log               internal.Logger
	bot               *tgbotapi.BotAPI
	updateConfig      tgbotapi.UpdateConfig
	handlerConfig     map[string]internal.HandleFunc
	pollUpdateHandler internal.HandleFunc
}

type specifications struct {
	TelegramToken string `split_words:"true"`
}

func NewFromEnv(logger internal.Logger, config map[string]internal.HandleFunc, pollHandler internal.HandleFunc) (*HandlerComposer, error) {
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
	updates := t.bot.GetUpdatesChan(t.updateConfig)
	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}
		if update.PollAnswer != nil {
			t.pollUpdateHandler(&update, t.bot)
		} else if handler, ok := t.handlerConfig[update.Message.Text]; ok {
			handler(&update, t.bot)
		} else {
			msg := tgbotapi.NewPhoto(update.FromChat().ID, tgbotapi.FilePath("./media/no_way.png"))
			_, _ = t.bot.Send(msg)
		}
	}

	return nil
}
