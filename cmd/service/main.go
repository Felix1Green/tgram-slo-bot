package main

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"tgram-slo-bot/internal"
	"tgram-slo-bot/internal/components/chat_storage"
	"tgram-slo-bot/internal/components/logger"
	"tgram-slo-bot/internal/components/poll_storage"
	"tgram-slo-bot/internal/rpc"
	"tgram-slo-bot/internal/rpc/flip_handler"
	"tgram-slo-bot/internal/rpc/go_handler"
	"tgram-slo-bot/internal/rpc/poll_update_handler"
	"tgram-slo-bot/internal/rpc/random_handler"
	"tgram-slo-bot/internal/rpc/reg_handler"
	"time"
)

func main() {
	// logger := MockLogger()
	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "redis:6379") },
	}
	var (
		log = logger.New()
		// without error until backoff realization
		chatStorage, _ = chat_storage.NewFromEnv(redisPool, log)
		pollStorage, _ = poll_storage.NewFromEnv(redisPool, log)
	)

	// init handlers
	var (
		goHandler         = go_handler.New(log, chatStorage, pollStorage)
		flipHandler       = flip_handler.New(log)
		randomHandler     = random_handler.New(log)
		regHandler        = reg_handler.New(log, chatStorage)
		pollUpdateHandler = poll_update_handler.New(log, pollStorage)
	)

	handlerComposer, err := rpc.NewFromEnv(log, map[string]internal.HandleFunc{
		"go":     goHandler.Handle,
		"flip":   flipHandler.Handle,
		"random": randomHandler.Handle,
		"reg":    regHandler.Handle,
	}, pollUpdateHandler.Handle, flipHandler.HandleChoice)
	if err != nil {
		log.Error(context.Background(), fmt.Sprintf("Handler composer initialization failed with error: %s", err.Error()))
		return
	}

	err = handlerComposer.Listen()
	if err != nil {
		log.Error(context.Background(), err)
	}
}
