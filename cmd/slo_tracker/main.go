package main

import (
	"github.com/gomodule/redigo/redis"
	"tgram-slo-bot/internal/components/chat_storage"
	"tgram-slo-bot/internal/components/logger"
	"tgram-slo-bot/internal/components/poll_storage"
	"tgram-slo-bot/internal/workers/monitoring"
	"time"
)

func main() {
	log := logger.New()
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

	worker := monitoring.New(log, pollStorage, chatStorage)
	// add cron
	worker.Run()
}
