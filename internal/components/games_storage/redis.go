package games_storage

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"tgram-slo-bot/internal"
	"time"
)

const (
	componentName = "redis_games_storage"
)

type redisOptions struct {
	BackoffMaxValue int   `split_words:"true"`
	BackoffMaxTries int64 `split_words:"true"`
}

type RedisGamesStorage struct {
	log             internal.Logger
	pool            *redis.Pool
	backoffMaxValue time.Duration
	backoffMaxTries int64
}

func NewFromEnv(pool *redis.Pool, logger internal.Logger) (*RedisGamesStorage, error) {
	options := &redisOptions{}
	err := internal.EnvOptions(componentName, options)
	if err != nil {
		return nil, err
	}

	return &RedisGamesStorage{
		pool:            pool,
		log:             logger,
		backoffMaxTries: options.BackoffMaxTries,
		backoffMaxValue: time.Duration(options.BackoffMaxValue) * time.Second,
	}, nil
}

func (s *RedisGamesStorage) GetChatGames(chatID int64) ([]string, error) {
	var (
		conn = s.pool.Get()
		err  error
	)
	defer func() {
		_ = conn.Close()
		ctx := s.log.WithFields(context.Background(), map[string]interface{}{
			"componentName": componentName,
		})
		if err != nil {
			s.log.Error(ctx, err)
		}
	}()

	values, err := redis.Strings(conn.Do("SMEMBERS", s.getKey(chatID)))
	if err != nil {
		return nil, err
	}

	userArr := make([]string, len(values))
	for _, val := range values {
		userArr = append(userArr, val)
	}

	return userArr, nil
}

func (s *RedisGamesStorage) CreateNewChatGame(chatID int64, game string) error {
	var (
		conn = s.pool.Get()
		err  error
	)
	defer func() {
		_ = conn.Close()
		ctx := s.log.WithFields(context.Background(), map[string]interface{}{
			"componentName": componentName,
		})
		if err != nil {
			s.log.Error(ctx, err)
		}
	}()

	_, err = conn.Do("SADD", s.getKey(chatID), game)
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisGamesStorage) getKey(chatID int64) string {
	return fmt.Sprintf("%s:%d", componentName, chatID)
}

//TODO: install backoff package
//func (s *RedisPollStorage) backoffDo(conn redis.Conn, commandName string, args ...interface{}) (reply interface{}, err error) {
//	backoffCfg := backoff.NewExponentialBackOff()
//	backoffCfg.MaxInterval = s.backoffMaxValue
//	retryCount := int64(0)
//
//	_ = backoff.Retry(func() error {
//		if retryCount > s.backoffMaxTries {
//			return nil
//		}
//
//		reply, err = conn.Do(commandName, args...)
//		retryCount++
//
//		return err
//	}, backoffCfg)
//
//	return reply, err
//}
