package messaging

import (
	"context"

	"github.com/go-redis/redis"
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
