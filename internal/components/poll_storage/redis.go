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

type redisOptions struct {
	BackoffMaxValue int   `split_words:"true"`
	BackoffMaxTries int64 `split_words:"true"`
}

type Poll struct {
	ChatID           int64
	PollID           string
	MsgID            int
	CreatedTimeStamp int64
}

type RedisPollStorage struct {
	log             internal.Logger
	pool            *redis.Pool
	backoffMaxValue time.Duration
	backoffMaxTries int64
}

func NewFromEnv(pool *redis.Pool, logger internal.Logger) (*RedisPollStorage, error) {
	options := &redisOptions{}
	err := internal.EnvOptions(componentName, options)
	if err != nil {
		return nil, err
	}

	return &RedisPollStorage{
		pool:            pool,
		log:             logger,
		backoffMaxTries: options.BackoffMaxTries,
		backoffMaxValue: time.Duration(options.BackoffMaxValue) * time.Second,
	}, nil
}

func (s *RedisPollStorage) CreateNewPoll(chatID int64, pollID string, msgID int) error {
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

	_, err = conn.Do("SADD", mainKey, s.createPollTimestampKey(chatID, pollID, msgID))
	if err == nil {
		err = s.SetActiveChatPoll(chatID, pollID)
	}

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
		msgID, _     = strconv.Atoi(splitter[3])
		timestamp, _ = strconv.Atoi(splitter[4])
	)

	return &Poll{
		ChatID:           int64(chatID),
		PollID:           pollID,
		MsgID:            msgID,
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

func (s *RedisPollStorage) SetUserVoted(pollID string, userId int64) error {
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

	_, err = conn.Do("SADD", s.createSimplePollKey(pollID), userId)
	return err
}

func (s *RedisPollStorage) IsCurrentUserVoted(pollID string, userID int64) (bool, error) {
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

	isMember, err := redis.Bool(conn.Do("SISMEMBER", s.createSimplePollKey(pollID), userID))
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

func (s *RedisPollStorage) SetUserRegistered(pollID string, userId int64) error {
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

	_, err = conn.Do("SADD", s.createRegisteredPollKey(pollID), userId)
	return err
}

func (s *RedisPollStorage) IsCurrentUserRegistered(pollID string, userID int64) (bool, error) {
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

	isMember, err := redis.Bool(conn.Do("SISMEMBER", s.createRegisteredPollKey(pollID), userID))
	return isMember, err
}

func (s *RedisPollStorage) SetActiveChatPoll(chatId int64, pollId string) error {
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

	_, err = conn.Do("SET", chatId, pollId)
	return err
}

func (s *RedisPollStorage) GetActiveChatPoll(chatId int64) (string, error) {
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

	return redis.String(conn.Do("GET", chatId))
}

func (s *RedisPollStorage) createPollTimestampKey(chatID int64, pollID string, msgID int) string {
	return fmt.Sprintf("%s:%d:%s:%d:%d", componentName, chatID, pollID, msgID, time.Now().Unix())
}

func (s *RedisPollStorage) createSimplePollKey(pollID string) string {
	return fmt.Sprintf("%s:%s", componentName, pollID)
}

func (s *RedisPollStorage) createRegisteredPollKey(pollId string) string {
	return fmt.Sprintf("%s:%s:%s", componentName, pollId, "yes")
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
