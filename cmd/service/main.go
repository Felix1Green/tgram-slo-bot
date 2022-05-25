package main

import (
	"context"
	"fmt"
	"tgram-slo-bot/internal"
	"tgram-slo-bot/internal/components/chat_storage"
	"tgram-slo-bot/internal/components/poll_storage"
	"tgram-slo-bot/internal/rpc"
	"tgram-slo-bot/internal/rpc/flip_handler"
	"tgram-slo-bot/internal/rpc/go_handler"
	"tgram-slo-bot/internal/rpc/poll_update_handler"
	"tgram-slo-bot/internal/rpc/random_handler"
	"tgram-slo-bot/internal/rpc/reg_handler"
)

func main() {
	// logger := MockLogger()
	var (
		logger      internal.Logger
		chatStorage chat_storage.Storage
		pollStorage poll_storage.Storage
	)

	// init handlers
	var (
		goHandler         = go_handler.New(logger, chatStorage, pollStorage)
		flipHandler       = flip_handler.New(logger)
		randomHandler     = random_handler.New(logger)
		regHandler        = reg_handler.New(logger, chatStorage)
		pollUpdateHandler = poll_update_handler.New(logger, pollStorage)
	)

	handlerComposer, err := rpc.NewFromEnv(logger, map[string]internal.HandleFunc{
		"/go":     goHandler.Handle,
		"/flip":   flipHandler.Handle,
		"/random": randomHandler.Handle,
		"/reg":    regHandler.Handle,
	}, pollUpdateHandler.Handle)
	if err != nil {
		logger.Error(context.Background(), fmt.Sprintf("Handler composer initialization failed with error: %s", err.Error()))
		return
	}

	err = handlerComposer.Listen()
	if err != nil {
		logger.Error(context.Background(), err)
	}
}
