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

type Client struct {
	conn *redis.Client
}

func NewClient(url string) (*Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})

	if _, err := conn.Ping().Result(); err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.conn.Ping().Err()
}

func (c *Client) Close() error {
	return c.conn.Close()
}

type Stream struct {
	client *Client
}

func NewStream(client *Client) *Stream {
	return &Stream{
		client: client,
	}
}

func (s *Stream) Publish(ctx context.Context, topic string, message []byte) error {
	if _, err := s.client.conn.XAdd(&redis.XAddArgs{
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
			streams, err := s.client.conn.XRead(&redis.XReadArgs{
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
