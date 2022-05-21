package main

import (
	"context"
	"fmt"
	"tgram-slo-bot/internal"
	"tgram-slo-bot/internal/rpc"
	"tgram-slo-bot/internal/rpc/flip_handler"
	"tgram-slo-bot/internal/rpc/go_handler"
	"tgram-slo-bot/internal/rpc/random_handler"
)

func main() {
	// logger := MockLogger()
	var logger internal.Logger

	// init handlers
	var (
		goHandler     = go_handler.New(logger)
		flipHandler   = flip_handler.New(logger)
		randomHandler = random_handler.New(logger)
	)

	handlerComposer, err := rpc.NewFromEnv(logger, map[string]internal.HandleFunc{
		"/go":     goHandler.Handle,
		"/flip":   flipHandler.Handle,
		"/random": randomHandler.Handle,
	})
	if err != nil {
		logger.Error(context.Background(), fmt.Sprintf("Handler composer initialization failed with error: %s", err.Error()))
		return
	}

	err = handlerComposer.Listen()
	if err != nil {
		logger.Error(context.Background(), err)
	}
}
