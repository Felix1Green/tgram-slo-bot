package main

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/robfig/cron"
	"os"
	"os/signal"
	"syscall"
	"tgram-slo-bot/internal/components/chat_storage"
	"tgram-slo-bot/internal/components/logger"
	"tgram-slo-bot/internal/components/poll_storage"
	"tgram-slo-bot/internal/workers/monitoring"
	"time"
)

func main() {
	log := logger.New()
	scheduler := cron.New()
	defer scheduler.Stop()

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "redis:6379") },
	}
	var (
		// without error until backoff realization
		chatStorage, _ = chat_storage.NewFromEnv(redisPool, log)
		pollStorage, _ = poll_storage.NewFromEnv(redisPool, log)
	)

	worker, err := monitoring.NewFromEnv(log, pollStorage, chatStorage)
	if err != nil {
		log.Error(context.Background(), err)
		return
	}

	// add cron
	err = scheduler.AddFunc("@every 1m", worker.Run)
	if err != nil {
		log.Error(context.Background(), err)
		return
	}
	go scheduler.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
