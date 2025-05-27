package redis

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

func NewClient(url string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}

type Stream struct {
	client *redis.Client
}

func NewStream(client *redis.Client) *Stream {
	return &Stream{
		client: client,
	}
}

func (s *Stream) Publish(ctx context.Context, topic string, message []byte) error {
	if _, err := s.client.XAdd(&redis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"data": message,
		},
	}).Result(); err != nil {
		return err
	}

	return nil
}

func (s *Stream) Subscribe(ctx context.Context, topic string, handler func(message []byte) error) error {
	lastID := "$"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			streams, err := s.client.XRead(&redis.XReadArgs{
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
