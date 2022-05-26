package chat_storage

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gomodule/redigo/redis"
	"tgram-slo-bot/internal"
	"time"
)

const (
	componentName = "redis_chat_storage"
)

type redisOptions struct {
	BackoffMaxValue int   `split_words:"true"`
	BackoffMaxTries int64 `split_words:"true"`
}

type RedisChatStorage struct {
	log             internal.Logger
	pool            *redis.Pool
	backoffMaxValue time.Duration
	backoffMaxTries int64
}

func NewFromEnv(pool *redis.Pool, logger internal.Logger) (*RedisChatStorage, error) {
	options := &redisOptions{}
	err := internal.EnvOptions(componentName, options)
	if err != nil {
		return nil, err
	}

	return &RedisChatStorage{
		pool:            pool,
		log:             logger,
		backoffMaxTries: options.BackoffMaxTries,
		backoffMaxValue: time.Duration(options.BackoffMaxValue) * time.Second,
	}, nil
}

func (s *RedisChatStorage) GetChatUsers(chatID int64) ([]*tgbotapi.User, error) {
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

	userArr := make([]*tgbotapi.User, len(values))
	for index, val := range values {
		err = json.Unmarshal([]byte(val), &userArr[index])
	}

	return userArr, nil
}

func (s *RedisChatStorage) RegistrateNewUser(chatID int64, user *tgbotapi.User) error {
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

	rawUser, _ := json.Marshal(user)
	_, err = conn.Do("SADD", s.getKey(chatID), string(rawUser))
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisChatStorage) getKey(chatID int64) string {
	return fmt.Sprintf("%s:%d", componentName, chatID)
}
