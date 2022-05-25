package poll_storage

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
	"tgram-slo-bot/internal"
	"time"
)

var (
	componentName     = "redis_poll_storage"
	mainKey           = fmt.Sprintf("%s:poll_keys", componentName)
	IncorrectInputKey = fmt.Errorf("inserted key is incorrect")
)

// TODO: add backoff on redis commands
//type redisOptions struct {
//	BackoffMaxValue int   `split_words:"true"`
//	BackoffMaxTries int64 `split_words:"true"`
//}

type Poll struct {
	ChatID           int64
	PollID           string
	CreatedTimeStamp int64
}

type RedisPollStorage struct {
	log  internal.Logger
	pool *redis.Pool
	//backoffMaxValue time.Duration
	//backoffMaxTries int64
}

func NewFromEnv(pool *redis.Pool, logger internal.Logger) (*RedisPollStorage, error) {
	//options := &redisOptions{}
	//err := internal.EnvOptions(componentName, options)
	//if err != nil {
	//	return nil, err
	//}

	return &RedisPollStorage{
		pool: pool,
		log:  logger,
		//backoffMaxTries: options.BackoffMaxTries,
		//backoffMaxValue: time.Duration(options.BackoffMaxValue) * time.Second,
	}, nil
}

func (s *RedisPollStorage) CreateNewPoll(chatID int64, pollID string) error {
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

	_, err = conn.Do("SADD", mainKey, s.createPollTimestampKey(chatID, pollID))

	return err
}

func (s *RedisPollStorage) GetPollInfoFromKey(key string) (*Poll, error) {
	splitter := strings.Split(key, ":")
	if len(splitter) < 4 {
		return nil, IncorrectInputKey
	}

	var (
		chatID, _    = strconv.Atoi(splitter[1])
		pollID       = splitter[2]
		timestamp, _ = strconv.Atoi(splitter[3])
	)

	return &Poll{
		ChatID:           int64(chatID),
		PollID:           pollID,
		CreatedTimeStamp: int64(timestamp),
	}, nil
}

func (s *RedisPollStorage) RemovePoll(pollKey string) error {
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

	_, err = conn.Do("SREM", mainKey, pollKey)
	return err
}

func (s *RedisPollStorage) SetUserVoted(chatID int64, pollID string, userId int64) error {
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

	_, err = conn.Do("SADD", s.createSimplePollKey(chatID, pollID), userId)
	return err
}

func (s *RedisPollStorage) IsCurrentUserVoted(chatID int64, pollID string, userID int64) (bool, error) {
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

	isMember, err := redis.Bool(conn.Do("SISMEMBER", s.createSimplePollKey(chatID, pollID), userID))
	return isMember, err
}

func (s *RedisPollStorage) GetActivePollKeys() ([]string, error) {
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

	polls, err := redis.Strings(conn.Do("SMEMBERS", mainKey))

	return polls, err
}

func (s *RedisPollStorage) createPollTimestampKey(chatID int64, pollID string) string {
	return fmt.Sprintf("%s:%d:%s:%d", componentName, chatID, pollID, time.Now().Unix())
}

func (s *RedisPollStorage) createSimplePollKey(chatID int64, pollID string) string {
	return fmt.Sprintf("%s:%d:%s", componentName, chatID, pollID)
}
