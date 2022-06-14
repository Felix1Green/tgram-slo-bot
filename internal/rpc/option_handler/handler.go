package option_handler

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"tgram-slo-bot/internal"
)

var (
	optionGeneratorURL = "https://yesno.wtf/api"
	handlerName        = "optionHandler"
)

type Handler struct {
	log internal.Logger
}

func New(log internal.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

type Response struct {
	Answer string `json:"answer"`
	Forced bool   `json:"forced"`
	Image  string `json:"image"`
}

func (h *Handler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var (
		request, _ = http.NewRequest("GET", optionGeneratorURL, nil)
		client     = http.Client{
			Timeout: 0,
		}
		err error
	)
	defer func() {
		if err != nil {
			ctx := h.log.WithFields(context.Background(), map[string]interface{}{
				"handler": handlerName,
			})
			h.log.Error(ctx, err)
		}
	}()

	resp, err := h.getOptionImage(&client, request)
	if err != nil {
		return
	}

	msg := tgbotapi.NewPhoto(update.FromChat().ID, tgbotapi.FileBytes{
		Bytes: resp,
	})
	_, _ = bot.Send(msg)
}

func (h *Handler) getOptionImage(client *http.Client, request *http.Request) ([]byte, error) {
	var (
		responseBody []byte
		response     = &Response{}
	)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	_, err = resp.Body.Read(responseBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, response)
	if err != nil {
		return nil, err
	}
	imageRequest, err := http.NewRequest("GET", response.Image, nil)
	if err != nil {
		return nil, err
	}

	resp, err = client.Do(imageRequest)

	_, err = resp.Body.Read(responseBody)
	return responseBody, err
}
