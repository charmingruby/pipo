package broker

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

const (
	batchSize = 100
	blockTime = 5 * time.Second
)

var (
	ErrInvalidMessageDataType = errors.New("invalid message data type")
)

type RedisStream struct {
	client *redis.Client
}

func NewRedisStream(client *redis.Client) *RedisStream {
	return &RedisStream{
		client: client,
	}
}

func (r *RedisStream) Publish(ctx context.Context, topic string, message []byte) error {
	if _, err := r.client.XAdd(&redis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"data": message,
		},
	}).Result(); err != nil {
		return err
	}

	return nil
}

func (r *RedisStream) Subscribe(ctx context.Context, topic string, handler func(message []byte) error) error {
	lastID := "$"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			streams, err := r.client.XRead(&redis.XReadArgs{
				Streams: []string{topic, lastID},
				Count:   batchSize,
				Block:   blockTime,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue
				}

				return err
			}

			if len(streams) == 0 || len(streams[0].Messages) == 0 {
				continue
			}

			messages := streams[0].Messages

			for _, message := range messages {
				data, ok := message.Values["data"].(string)
				if !ok {
					return ErrInvalidMessageDataType
				}

				if err := handler([]byte(data)); err != nil {
					return err
				}

				lastID = message.ID
			}
		}
	}
}
